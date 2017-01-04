package segment

import (
	"bytes"
	"context"
	"errors"
	"github.com/goguardian/aws-xray-go/attributes"
	"github.com/goguardian/aws-xray-go/utils"
	"net/http"
	"net/url"
	"testing"
)

func TestNewSegment(t *testing.T) {
	req1 := &http.Request{
		Header: http.Header{
			utils.XRayHeader: []string{"Root=123; Parent=456; Sampled=1"},
		},
	}
	req1 = req1.WithContext(utils.ContextFromHeaders(req1))

	req2 := &http.Request{
		Header: http.Header{
			utils.XRayHeader: []string{"Root=123; Parent=456; Sampled=0"},
		},
	}
	req2 = req2.WithContext(utils.ContextFromHeaders(req2))

	tests := []struct {
		name           string
		context        context.Context
		expectTraceID  string
		expectParentID string
	}{
		{
			name:    "test",
			context: nil,
		},
		{
			name:           "test",
			context:        req1.Context(),
			expectTraceID:  "123",
			expectParentID: "456",
		},
		{
			name:           "test",
			context:        req2.Context(),
			expectTraceID:  "123",
			expectParentID: "456",
		},
	}

	for _, test := range tests {
		seg := New(test.name, test.context)

		if test.expectTraceID != "" && seg.TraceID != test.expectTraceID {
			t.Errorf("Segment trace ID is %s, expected %s", seg.TraceID, test.expectTraceID)
		}

		if test.expectParentID != "" && seg.ParentID != test.expectParentID {
			t.Errorf("Segment parent ID is %s, expected %s", seg.ParentID, test.expectParentID)
		}

		if seg.Name != test.name {
			t.Errorf("Segment name %s, expected: %s", seg.Name, test.name)
		}

		if !seg.InProgress {
			t.Error("Segment should be in progress")
		}

		if _, err := seg.Bytes(); err != nil {
			t.Error(err)
		}

		if _, err := seg.String(); err != nil {
			t.Error(err)
		}

		if seg.AddFault(); seg.Fault != true {
			t.Error("Segment should be marked fault")
		}

		if seg.AddThrottle(); seg.Throttle != true {
			t.Error("Segment should be marked throttle")
		}

		seg.Fault = false
		seg.AddError(errors.New("An error"))
		if !seg.Fault {
			t.Error("Segment should be marked fault")
		}

		seg.exception = &exception{
			Ex: "An error",
		}
		seg.AddError(errors.New("An error"))
		if seg.exception != nil {
			t.Error("Segment exception should be nil")
		}

		seg.exception = &exception{
			Ex: "An error 2",
		}
		seg.AddError(errors.New("An error"))
		if seg.exception != nil {
			t.Error("Segment exception should be nil")
		}

		key, value := "key", "value"
		seg.AddAnnotation(key, value)
		if val, ok := seg.Annotations[key]; !ok || val != value {
			t.Errorf("Annotation '%s' should be set with value '%s', not '%s'",
				key, value, val)
		}

		badkey, badvalue := "bad", bytes.Buffer{}
		seg.AddAnnotation(badkey, badvalue)
		if _, ok := seg.Annotations[badkey]; ok {
			t.Errorf("Annotation '%s' should not have a value", badkey)
		}

		seg.AddNewSubsegment("subsegment")
		if len(seg.Subsegments) != 1 {
			t.Error("Segment should have one Subsegment")
		}

		seg.AddServiceVersion("1.2.3.4")
		if seg.Service.Version != "1.2.3.4" {
			t.Errorf("Segment should have service version '1.2.3.4', but got: '%s'",
				seg.Service.Version)
		}

		if seg.HTTP != nil {
			t.Error("Segment HTTP should be nil")
		}

		reqURL := &url.URL{
			Scheme: "https",
			Host:   "www.example.com",
			Path:   "/",
		}
		seg.AddHTTPAttribute(attributes.NewLocal(&http.Request{
			Method: http.MethodGet,
			URL:    reqURL,
		}))
		if seg.HTTP.Request.Method != http.MethodGet {
			t.Errorf("Segment HTTP attribute request method should be '%s', "+
				"but got: '%s'", http.MethodGet, seg.HTTP.Request.Method)
		}
		if seg.HTTP.Request.URL != reqURL.String() {
			t.Errorf("Segment HTTP attribute request URL should be '%s', "+
				"but got: '%s'", reqURL.String(), seg.HTTP.Request.URL)
		}

		seg.Close()

		seg.DecrementCounter()

		if seg.InProgress {
			t.Error("Segment should not be in progress")
		}
	}
}

func TestSetSampler(t *testing.T) {
	SetSampler(utils.NewSampler(1, 0))

	traceCount := 0
	for i := 0; i < 10; i++ {
		seg := New("test", nil)
		if seg.Traced {
			traceCount++
		}
	}

	if traceCount != 1 {
		t.Errorf("Trace count should be %d, not %d", 1, traceCount)
	}
}
