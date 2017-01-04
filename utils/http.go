package utils

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	// ErrorType represents a 4XX status code
	ErrorType = "error"
	// FaultType represents a 5XX status code
	FaultType = "fault"
)

// GetCauseFromHTTPStatus returns a string representing an error, fault, or no
// cause based on HTTP response status codes
func GetCauseFromHTTPStatus(status int) string {
	if status >= 400 && status <= 499 {
		return ErrorType
	}

	if status >= 500 && status <= 599 {
		return FaultType
	}

	return ""
}

// GetClientIP returns the IP address of a request from an X-Forwarded-For
// header or remote address
func GetClientIP(req *http.Request) string {
	remoteAddr := string(req.Header.Get("X-Forwarded-For"))

	if remoteAddr == "" {
		remoteAddr = strings.Split(req.RemoteAddr, ":")[0]
	}

	return strings.Split(remoteAddr, ",")[0]
}

// GetContentLength returns the content length of a response if provided as
// a header.
func GetContentLength(res *http.Response) (int, error) {
	contentLength := string(res.Header.Get("Content-Length"))

	contentLen, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(contentLen), nil
}
