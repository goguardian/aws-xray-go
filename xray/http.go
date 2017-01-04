package xray

import (
	"context"
	"github.com/goguardian/aws-xray-go/handlers"
	"github.com/goguardian/aws-xray-go/utils"
	"net/http"
)

// GetHTTPClient returns an HTTP client for performing X-Ray subsegment traced
// HTTP requests.
func GetHTTPClient(ctx context.Context) (*http.Client, error) {
	segment, err := handlers.GetSegmentFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return handlers.NewHTTPClient(segment), nil
}

// Middleware provides a middleware for tracing HTTP handlers.
func Middleware(name string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := utils.ContextFromHeaders(r)

		r = r.WithContext(NewContext(name, ctx))
		defer Close(r.Context())

		AddLocalHTTP(r)

		handler(w, r)
	}
}
