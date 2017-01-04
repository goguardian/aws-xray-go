package segment

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/goguardian/aws-xray-go/attributes"
	"github.com/goguardian/aws-xray-go/utils"
	"os"
	"sync"
)

var (
	emitter      = NewEmitter()
	sampler      = utils.NewSampler(10, 0.05)
	samplerMutex = &sync.RWMutex{}
)

// Segment represents a segment.
type Segment struct {
	StartTime   float64                `json:"start_time"`
	EndTime     float64                `json:"end_time"`
	InProgress  bool                   `json:"in_progress"`
	Throttle    bool                   `json:"throttle"`
	Fault       bool                   `json:"fault"`
	Traced      bool                   `json:"-"`
	Counter     int32                  `json:"-"`
	ID          string                 `json:"id"`
	TraceID     string                 `json:"trace_id"`
	ParentID    string                 `json:"parent_id,omitempty"`
	Name        string                 `json:"name"`
	HTTP        *attributes.Local      `json:"http,omitempty"`
	Annotations map[string]interface{} `json:"annotations,omitempty"`
	Subsegments []*Subsegment          `json:"subsegments,omitempty"`
	Service     *service               `json:"service,omitempty"`
	Cause       *cause                 `json:"cause,omitempty"`
	exception   *exception

	sync.RWMutex
}

type exception struct {
	Ex    string `json:"-"`
	Cause *cause `json:"-"`
}

type service struct {
	Version string `json:"version,omitempty"`
}

type cause struct {
	ID               *cause                       `json:"cause"`
	WorkingDirectory string                       `json:"working_directory"`
	Paths            []string                     `json:"paths"`
	Exceptions       []*attributes.LocalException `json:"exceptions"`
}

// New creates a new segment.
func New(name string, ctx context.Context) *Segment {
	startTime := utils.CurrentTimeSecond()

	traceID, parentID, sampled := utils.GetIDsFromContext(ctx)

	if traceID == "" {
		traceIDSuffix := make([]byte, 12)
		rand.Read(traceIDSuffix)
		traceID = fmt.Sprintf("1-%x-%x", int64(startTime), traceIDSuffix)
	}

	idBytes := make([]byte, 8)
	rand.Read(idBytes)
	id := fmt.Sprintf("%x", idBytes)

	seg := &Segment{
		ID:         id,
		TraceID:    traceID,
		ParentID:   parentID,
		Name:       name,
		StartTime:  startTime,
		InProgress: true,
	}

	seg.resolveSampling(sampled)

	return seg
}

// AddAnnotation adds a key-value pair that can be queried with
// GetTraceSummaries.  Acceptable value types are string, numbers, and boolean.
func (s *Segment) AddAnnotation(key string, value interface{}) error {
	switch value.(type) {
	case bool, string, int, int16, int32, int64, uint, uint16, uint32, uint64,
		float32, float64:
	case fmt.Stringer:
		value = value.(fmt.Stringer).String()
	default:
		return fmt.Errorf("Failed to add annotation key: %s value: %v to "+
			"segment %s. Value must be a string, number, or boolean.",
			key, value, s.Name)
	}

	s.Lock()
	defer s.Unlock()

	if s.Annotations == nil {
		s.Annotations = map[string]interface{}{}
	}

	s.Annotations[key] = value

	return nil
}

// AddError adds error data into the segment.
func (s *Segment) AddError(err error) {
	s.AddFault()

	s.Lock()
	defer s.Unlock()

	if s.exception != nil {
		if err.Error() == s.exception.Ex {
			s.Cause = &cause{ID: s.exception.Cause}
			s.exception = nil
			return
		}

		s.exception = nil
	}

	if s.Cause == nil {
		wd, _ := os.Getwd()
		s.Cause = &cause{
			WorkingDirectory: wd,
			Paths:            []string{},
			Exceptions:       []*attributes.LocalException{},
		}
	}

	s.Cause.Exceptions = append(s.Cause.Exceptions,
		attributes.NewLocalException(err))
}

// AddFault adds fault flag to the segment.
func (s *Segment) AddFault() {
	s.Lock()
	defer s.Unlock()
	s.Fault = true
}

// AddHTTPAttribute adds an HTTP property with associated data to the segment.
func (s *Segment) AddHTTPAttribute(localHTTP *attributes.Local) {
	s.Lock()
	defer s.Unlock()
	s.HTTP = localHTTP
}

// AddNewSubsegment adds a new subsegment to the slice of subsegments,
func (s *Segment) AddNewSubsegment(name string) *Subsegment {
	subseg := NewSubsegment(name)
	s.AddSubsegment(subseg)
	return subseg
}

// AddServiceVersion adds a service with associated version data to the
// segment.
func (s *Segment) AddServiceVersion(version string) {
	s.Lock()
	defer s.Unlock()
	s.Service = &service{Version: version}
}

// AddSubsegment adds a subsegment to the slice of subsegments.
func (s *Segment) AddSubsegment(subseg *Subsegment) {
	s.Lock()
	defer s.Unlock()

	subseg.Lock()
	defer subseg.Unlock()

	subseg.Segment = s

	s.Subsegments = append(s.Subsegments, subseg)

	if subseg.EndTime == 0 {
		s.Counter++
	}
}

// AddThrottle adds throttle flag to the segment.
func (s *Segment) AddThrottle() {
	s.Lock()
	defer s.Unlock()
	s.Throttle = true
}

// Bytes returns the segment as a JSON encoded byte slice
func (s *Segment) Bytes() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	return json.Marshal(s)
}

// Close closes the segment and sets the end time.
func (s *Segment) Close() error {
	s.Lock()

	if s.EndTime == 0 {
		s.EndTime = utils.CurrentTimeSecond()
	}

	s.InProgress = false

	s.Unlock()

	s.RLock()
	counter := s.Counter
	s.RUnlock()

	if counter <= 0 {
		return s.Flush()
	}

	return nil
}

// DecrementCounter decrements the count of open subsegments for the segment.
func (s *Segment) DecrementCounter() error {
	s.Lock()
	s.Counter--
	counter := s.Counter
	endTime := s.EndTime
	s.Unlock()

	if counter <= 0 && endTime > 0 {
		return s.Flush()
	}

	return nil
}

// Flush read locks and sends the segment to the daemon.
func (s *Segment) Flush() error {
	s.RLock()
	defer s.RUnlock()

	if !s.Traced {
		return nil
	}

	return emitter.Send(s)
}

// resolveSampling determines whether to sample the segment
func (s *Segment) resolveSampling(sampled string) {
	s.Lock()
	defer s.Unlock()

	if sampled == "1" {
		s.Traced = true
		return
	}

	if sampled == "0" {
		s.Traced = false
		return
	}

	s.Traced = sampler.IsSampled()
}

// String returns the segment as a JSON encode string
func (s *Segment) String() (string, error) {
	segBytes, err := s.Bytes()
	return string(segBytes), err
}

// SetSampler updates the sampler used for segment sampling.
func SetSampler(s *utils.Sampler) {
	samplerMutex.Lock()
	defer samplerMutex.Unlock()
	sampler = s
}
