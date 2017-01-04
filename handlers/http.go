package handlers

import (
	"fmt"
	"github.com/goguardian/aws-xray-go/attributes"
	"github.com/goguardian/aws-xray-go/segment"
	"github.com/goguardian/aws-xray-go/utils"
	"net/http"
	"strings"
)

// NewHTTPClient creates a new HTTP client for a segment.
func NewHTTPClient(seg *segment.Segment) *http.Client {
	hc := &http.Client{}
	hc.Transport = NewHTTPInterceptor(seg)
	return hc
}

// HTTPInterceptor represents an http.RoundTripper that also performs tracing
// of requests by creating subsegments of the segment.
type HTTPInterceptor struct {
	transport *http.Transport
	segment   *segment.Segment
}

// NewHTTPInterceptor creates a new HTTPInterceptor.
func NewHTTPInterceptor(seg *segment.Segment) http.RoundTripper {
	return &HTTPInterceptor{transport: &http.Transport{}, segment: seg}
}

// RoundTrip implements http.RoundTripper interface.  Adds a new subsegment,
// adds trace header to HTTP requests, performs the requests, and records
// remote response data.
func (h HTTPInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	subseg := h.segment.AddNewSubsegment(req.URL.Host)
	subseg.AddRemote()

	addXRayHeader(req, subseg)

	if strings.HasSuffix(req.URL.Host, ".amazonaws.com") {
		// TODO: determine consistent means of interpreting AWS requests for
		// tracing AWS
	}

	res, err := h.transport.RoundTrip(req)
	if err != nil {
		subseg.Close(err, utils.ErrorType)
		return res, err
	}

	errType := utils.GetCauseFromHTTPStatus(res.StatusCode)

	// 4XX/5XX Status Code
	if errType != "" {
		err = fmt.Errorf("%s: status code %d", errType, res.StatusCode)
	}

	if res.StatusCode == http.StatusTooManyRequests {
		subseg.AddThrottle()
	}

	contentLength, _ := utils.GetContentLength(res)

	remoteData := &attributes.Remote{
		Request: &attributes.RemoteRequest{
			Method: req.Method,
			URL:    req.URL.String(),
			Traced: subseg.Segment.Traced,
		},
		Response: &attributes.RemoteResponse{
			Status:        res.StatusCode,
			ContentLength: contentLength,
		},
	}

	subseg.AddRemoteData(remoteData)

	subseg.Close(err, errType)

	return res, err
}

// addXRayHeader adds X-Ray trace header to an HTTP request.
func addXRayHeader(req *http.Request, subseg *segment.Subsegment) {
	traceID := ""
	sampled := "0"
	if subseg.Segment != nil {
		traceID = subseg.Segment.TraceID
		if subseg.Segment.Traced {
			sampled = "1"
		}
	}

	if req.Header == nil {
		req.Header = http.Header{}
	}

	req.Header.Add(utils.XRayHeader, fmt.Sprintf(
		"Root=%s; Parent=%s; Sampled=%s", traceID, subseg.ID, sampled))
}
