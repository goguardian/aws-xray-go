package handlers

import (
	"github.com/goguardian/aws-xray-go/segment"
	"github.com/goguardian/aws-xray-go/utils"
	"net/http"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	seg := segment.New("segment", nil)

	client := NewHTTPClient(seg)
	if client == nil {
		t.Error("HTTP Client should not be nil")
	}

	_, err := client.Get("https://aws.amazon.com")
	if err != nil {
		t.Error(err)
	}
}

func TestAddXRayHeader(t *testing.T) {
	req := &http.Request{}

	seg := segment.New("segment", nil)
	seg.Traced = false
	subseg := segment.NewSubsegment("subsegment")
	subseg.Segment = seg

	addXRayHeader(req, subseg)

	header := req.Header.Get(utils.XRayHeader)
	if header == "" {
		t.Errorf("HTTP request header '%s' should be set", utils.XRayHeader)
	}
}
