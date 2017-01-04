package segment

import (
	"bytes"
	"errors"
	"github.com/goguardian/aws-xray-go/attributes"
	"github.com/goguardian/aws-xray-go/utils"
	"testing"
)

func TestNewSubsegment(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "test"},
		{name: "test2"},
	}

	seg := New("segment", nil)

	for _, test := range tests {
		subseg := NewSubsegment(test.name)

		subseg.Segment = seg

		if subseg.Name != test.name {
			t.Errorf("Subsegment name is %s, expected: %s", subseg.Name, test.name)
		}

		if subseg.ID == "" {
			t.Error("SubSgment ID should not be empty")
		}

		if len(subseg.Subsegments) > 0 {
			t.Error("Subsegment should have no Subsegments")
		}

		subseg.AddNewSubsegment("subsegment")
		if len(subseg.Subsegments) != 1 {
			t.Error("Subsegment should have 1 Subsegment")
		}

		if len(subseg.PrecursorIDs) > 0 {
			t.Error("Subsegment should have no PrecursorIDs")
		}

		subseg.AddPrecursorID("id")
		if len(subseg.PrecursorIDs) != 1 {
			t.Error("Subsegment should have 1 PrecursorID")
		}

		subseg.AddAnnotation("key", "value")
		if val, ok := subseg.Annotations["key"]; !ok || val != "value" {
			t.Error("Subsegment Annotation 'key' should have value 'value'")
		}

		subseg.AddAnnotation("badkey", bytes.Buffer{})
		if _, ok := subseg.Annotations["badkey"]; ok {
			t.Error("Subsegment Annotation 'badkey' should not have a value")
		}

		if subseg.Fault {
			t.Error("Subsegment should not be marked fault")
		}

		subseg.AddFault()
		if !subseg.Fault {
			t.Error("Subsegment should be marked fault")
		}

		if subseg.Throttle {
			t.Error("Subsegment should not be marked throttle")
		}

		subseg.AddThrottle()
		if !subseg.Throttle {
			t.Error("Subsegment should be marked throttle")
		}

		if subseg.Namespace != "" {
			t.Error("Subsegment should have an empty namespace")
		}

		subseg.AddRemote()
		if subseg.Namespace != "remote" {
			t.Error("Subsegment should have Namespace 'remote'")
		}

		if subseg.Metadata != nil {
			t.Error("Subsegment Metadata should be nil")
		}

		subseg.AddMetadata("key", "value")
		if subseg.Metadata == nil {
			t.Error("Subsegment Metadata should not be nil")
		}

		if val, ok := subseg.Metadata.Default["key"]; !ok || val.(string) != "value" {
			t.Error("Subsegment Metadata Default key 'key' should have value 'value'")
		}

		subseg.AddError(nil, "")
		if subseg.Error {
			t.Error("Subsegment should not be marked error")
		}

		subseg.AddError(errors.New(""), utils.ErrorType)
		if !subseg.Error {
			t.Error("Subsegment should be marked error")
		}

		subseg.Fault = false
		subseg.AddError(errors.New(""), utils.FaultType)
		if !subseg.Fault {
			t.Error("Subsegment should be marked fault")
		}

		if subseg.RemoteData != nil {
			t.Error("Subsegment should not have remote data")
		}

		remoteData := &attributes.Remote{}
		subseg.AddRemoteData(remoteData)
		if subseg.RemoteData != remoteData {
			t.Error("Subsegment should have remote data")
		}

		subseg.Close(nil, "")

		subseg.Fault = false
		subseg.Close(errors.New(""), utils.FaultType)
		if !subseg.Fault {
			t.Error("Subsegment should be marked fault")
		}
	}
}
