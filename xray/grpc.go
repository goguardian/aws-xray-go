package xray

import (
	"github.com/goguardian/aws-xray-go/handlers"
	"time"

	"google.golang.org/grpc"
)

// NewGRPCServer creates a new gRPC server with a unary interceptor that makes
// segments for all handled requests.
func NewGRPCServer(name string, opt ...grpc.ServerOption) *grpc.Server {
	options := append(opt, grpc.UnaryInterceptor(
		handlers.GRPCServerUnaryInterceptor(name)))
	return grpc.NewServer(options...)
}

// SetSegmentCacheDuration updates the duration of segments in in-memory cache.
func SetSegmentCacheDuration(duration time.Duration) {
	handlers.SetSegmentCacheDuration(duration)
}
