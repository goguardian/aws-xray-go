package attributes

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNewLocal(t *testing.T) {
	tests := []struct {
		req                 *http.Request
		res                 *http.Response
		expectMethod        string
		expectURL           string
		expectIP            string
		expectAgent         string
		expectContentLength int
		expectStatus        int
	}{
		{},
		{
			req: &http.Request{},
			res: &http.Response{},
		},
		{
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "https",
					Host:   "www.example.com",
					Path:   "/test",
				},
				RemoteAddr: "127.0.0.1:2000",
				Header: http.Header{
					"User-Agent": []string{"user-agent"},
				},
			},
			res: &http.Response{
				Header: http.Header{
					"Content-Length": []string{"123"},
				},
				StatusCode: 200,
			},
			expectMethod:        http.MethodGet,
			expectURL:           "https://www.example.com/test",
			expectIP:            "127.0.0.1",
			expectAgent:         "user-agent",
			expectContentLength: 123,
			expectStatus:        200,
		},
	}

	for _, test := range tests {
		local := NewLocal(test.req)
		local.Close(test.res)

		if local.Request.ClientIP != test.expectIP {
			t.Errorf("Expected client IP '%s', got '%s'",
				test.expectIP, local.Request.ClientIP)
		}

		if local.Request.Method != test.expectMethod {
			t.Errorf("Expected method '%s', got '%s'",
				test.expectMethod, local.Request.Method)
		}

		if local.Request.URL != test.expectURL {
			t.Errorf("Expected URL '%s', got '%s'",
				test.expectURL, local.Request.URL)
		}

		if local.Request.UserAgent != test.expectAgent {
			t.Errorf("Expected user agent '%s', got '%s'",
				test.expectAgent, local.Request.UserAgent)
		}

		if local.Response.ContentLength != test.expectContentLength {
			t.Errorf("Expected content length '%d', got '%d'",
				test.expectContentLength, local.Response.ContentLength)
		}

		if local.Response.Status != test.expectStatus {
			t.Errorf("Expected status '%d', got '%d'",
				test.expectStatus, local.Response.Status)
		}
	}
}
