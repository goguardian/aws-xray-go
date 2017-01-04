package segment

import (
	"crypto/rand"
	"fmt"
	"github.com/goguardian/aws-xray-go/attributes"
	"github.com/goguardian/aws-xray-go/utils"
	"sync"
)

// Subsegment represents a subsegment
type Subsegment struct {
	Segment      *Segment               `json:"-"`
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	StartTime    float64                `json:"start_time"`
	EndTime      float64                `json:"end_time"`
	PrecursorIDs []string               `json:"precursor_ids,omitempty"`
	Namespace    string                 `json:"namespace,omitempty"`
	Throttle     bool                   `json:"throttle"`
	Fault        bool                   `json:"fault"`
	Error        bool                   `json:"error"`
	Annotations  map[string]interface{} `json:"annotations,omitempty"`
	RemoteData   *attributes.Remote     `json:"http,omitempty"`
	Metadata     *metadata              `json:"metadata,omitempty"`
	Subsegments  []*Subsegment          `json:"subsegments,omitempty"`

	sync.RWMutex
}

type metadata struct {
	Default map[string]interface{} `json:"default,omitempty"`
}

// NewSubsegment creates a new default subsegment.
func NewSubsegment(name string) *Subsegment {
	startTime := utils.CurrentTimeSecond()

	idBytes := make([]byte, 8)
	rand.Read(idBytes)
	id := fmt.Sprintf("%x", idBytes)

	return &Subsegment{
		ID:        id,
		Name:      name,
		StartTime: startTime,
	}
}

// AddAnnotation adds a key-value pair that can be queryable through
// GetTraceSummaries.  Only accepted value types are strings, numbers, and
// boolean.
func (s *Subsegment) AddAnnotation(key string, value interface{}) error {
	switch value.(type) {
	case bool, string, int, int16, int32, int64, uint, uint16, uint32, uint64,
		float32, float64:
	case fmt.Stringer:
		value = value.(fmt.Stringer).String()
	default:
		return fmt.Errorf("Failed to add annotation key: %s value: %v to "+
			"subsegment %s. Value must be a string, number, or boolean.",
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

// AddError adds an error with associated data into the subsegment.
func (s *Subsegment) AddError(err error, errType string) {
	if err == nil {
		return
	}

	if errType != utils.FaultType && errType != utils.ErrorType {
		errType = utils.FaultType
	}

	if errType == utils.FaultType {
		s.AddFault()
		return
	}

	s.Lock()
	defer s.Unlock()

	s.Error = true
}

// AddFault adds fault flag into the subsegment.
func (s *Subsegment) AddFault() {
	s.Lock()
	defer s.Unlock()
	s.Fault = true
}

// AddMetadata adds a key-value pair to the subsegment.  Metadata is not
// queryable, but is recorded.
func (s *Subsegment) AddMetadata(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	if s.Metadata == nil {
		s.Metadata = &metadata{Default: map[string]interface{}{}}
	}

	s.Metadata.Default[key] = value
}

// AddNewSubsegment adds a new subsegment to the slice of subsegments.
func (s *Subsegment) AddNewSubsegment(name string) *Subsegment {
	subseg := NewSubsegment(name)
	s.AddSubsegment(subseg)
	return subseg
}

// AddPrecursorID adds a subsegment ID to record ordering.
func (s *Subsegment) AddPrecursorID(id string) {
	s.Lock()
	defer s.Unlock()
	s.PrecursorIDs = append(s.PrecursorIDs, id)
}

// AddRemote adds remote flag into the subsegment.
func (s *Subsegment) AddRemote() {
	s.Lock()
	defer s.Unlock()
	s.Namespace = "remote"
}

// AddRemoteData adds data for an outgoing HTTP/HTTPS call.
func (s *Subsegment) AddRemoteData(remote *attributes.Remote) {
	s.Lock()
	defer s.Unlock()
	s.RemoteData = remote
}

// AddSubsegment adds a subsegment to the slice of subsegment.
func (s *Subsegment) AddSubsegment(subseg *Subsegment) {
	s.Lock()
	defer s.Unlock()

	if subseg.EndTime == 0 {
		subseg.Segment = s.Segment

		if s.Segment != nil {
			s.Segment.Counter++
		}
	}

	s.Subsegments = append(s.Subsegments, subseg)
}

// AddThrottle adds throttled flag into the subsegment.
func (s *Subsegment) AddThrottle() {
	s.Lock()
	defer s.Unlock()
	s.Throttle = true
}

// Close closes the current subsegment.  Additionally, it captures any exception
// and sets the end time.
func (s *Subsegment) Close(err error, errType string) {
	s.Lock()
	if s.EndTime == 0 {
		s.EndTime = utils.CurrentTimeSecond()
	}
	s.Unlock()

	if err != nil {
		s.AddError(err, errType)
	}

	s.RLock()
	segment := s.Segment
	s.RUnlock()

	if segment != nil {
		s.Segment.DecrementCounter()
	}
}
