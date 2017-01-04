package attributes

import (
	"github.com/goguardian/aws-xray-go/utils"
	"net/http"
)

// Local represents an incoming HTTP/HTTPS call.
type Local struct {
	Request  *LocalRequest  `json:"request,omitempty"`
	Response *LocalResponse `json:"response,omitempty"`
}

// LocalRequest represents the request of an incoming HTTP/HTTPS call.
type LocalRequest struct {
	URL       string `json:"url,omitempty"`
	Method    string `json:"method,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	ClientIP  string `json:"client_ip,omitempty"`
}

// LocalResponse represents the response to an incoming HTTP/HTTPS call.
type LocalResponse struct {
	Status        int `json:"status"`
	ContentLength int `json:"content_length"`
}

// NewLocal creates a new Local object from an HTTP request.
func NewLocal(req *http.Request) *Local {
	local := &Local{
		Request:  &LocalRequest{},
		Response: &LocalResponse{},
	}

	if req == nil {
		return local
	}

	urlStr := ""
	if req.URL != nil {
		urlStr = req.URL.String()
	}

	local.Request = &LocalRequest{
		Method:    req.Method,
		UserAgent: req.UserAgent(),
		ClientIP:  utils.GetClientIP(req),
		URL:       urlStr,
	}

	return local
}

// Close closes a local HTTP/HTTPS request, adding response data.
func (l *Local) Close(res *http.Response) {
	if res == nil {
		return
	}

	l.Response.Status = res.StatusCode

	contentLength, _ := utils.GetContentLength(res)
	l.Response.ContentLength = contentLength
}
