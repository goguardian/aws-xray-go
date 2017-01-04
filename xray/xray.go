package xray

import (
	"github.com/goguardian/aws-xray-go/attributes"
	"github.com/goguardian/aws-xray-go/handlers"
	"github.com/goguardian/aws-xray-go/segment"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewContext creates a new segment and adds it to the request context.
func NewContext(name string, ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	segment := segment.New(name, ctx)

	return handlers.AddSegmentToContext(segment, ctx)
}

// AddLocalHTTP adds local HTTP request attributes to the request segment.
func AddLocalHTTP(r *http.Request) error {
	seg, err := handlers.GetSegmentFromContext(r.Context())
	if err != nil {
		return err
	}

	seg.AddHTTPAttribute(attributes.NewLocal(r))

	return nil
}

// Close closes the context segment.
func Close(ctx context.Context) error {
	seg, err := handlers.GetSegmentFromContext(ctx)
	if err != nil {
		return err
	}

	return seg.Close()
}

// GetGRPCClientConn returns a gRPC connection for performing X-Ray subsegment
// traced gRPC requests.
func GetGRPCClientConn(
	address string,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {

	return handlers.NewGRPCClientConn(address, opts...)
}

// GetSegment returns the segment of the X-Ray context.
func GetSegment(ctx context.Context) (*segment.Segment, error) {
	return handlers.GetSegmentFromContext(ctx)
}
