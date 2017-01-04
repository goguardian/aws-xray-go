package utils

import (
	"net/http"
	"testing"
)

func TestGetCauseFromHTTPStatus(t *testing.T) {
	tests := []struct {
		statusCode  int
		expectCause string
	}{
		{200, ""},
		{301, ""},
		{403, ErrorType},
		{500, FaultType},
	}

	for _, test := range tests {
		cause := GetCauseFromHTTPStatus(test.statusCode)
		if cause != test.expectCause {
			t.Errorf("Expected cause '%s', got '%s'", test.expectCause, cause)
		}
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		req      *http.Request
		expectIP string
	}{
		{
			req:      &http.Request{},
			expectIP: "",
		},
		{
			req: &http.Request{
				Header: http.Header{"X-Forwarded-For": []string{"127.0.0.1"}},
			},
			expectIP: "127.0.0.1",
		},
		{
			req: &http.Request{
				RemoteAddr: "127.0.0.1:2000",
			},
			expectIP: "127.0.0.1",
		},
	}

	for _, test := range tests {
		ip := GetClientIP(test.req)
		if ip != test.expectIP {
			t.Errorf("Expected IP '%s', got '%s'", test.expectIP, ip)
		}
	}
}

func TestGetContentLength(t *testing.T) {
	tests := []struct {
		res                 *http.Response
		expectContentLength int
		expectError         bool
	}{
		{
			res: &http.Response{
				Header: http.Header{"Content-Length": []string{"123"}},
			},
			expectContentLength: 123,
			expectError:         false,
		},
		{
			res:                 &http.Response{},
			expectContentLength: 0,
			expectError:         true,
		},
	}

	for _, test := range tests {
		contentLength, err := GetContentLength(test.res)
		if contentLength != test.expectContentLength {
			t.Errorf("Expected content length '%d', got '%d'",
				test.expectContentLength, contentLength)
		}

		if err != nil && !test.expectError {
			t.Error("Expected error checking content length")
		}

		if err == nil && test.expectError {
			t.Error("Unexpected error checking content length", err)
		}
	}
}
