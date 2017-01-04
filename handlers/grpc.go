package handlers

import (
	"errors"
	"fmt"
	"github.com/goguardian/aws-xray-go/segment"
	"github.com/goguardian/aws-xray-go/utils"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	mdParentKey  = "xray-parentid"
	mdRootKey    = "xray-rootid"
	mdSampledKey = "xray-sampled"
	mdSegmentKey = "xray-segment"
)

var (
	segmentCache          *cache.Cache
	cacheCleanupFrequency = 30 * time.Second
	cacheDuration         = 10 * time.Minute
	cacheDurationMutex    = &sync.RWMutex{}
)

func init() {
	segmentCache = cache.New(cache.NoExpiration, cacheCleanupFrequency)
}

// SetSegmentCacheDuration updates the duration of segments in in-memory cache
func SetSegmentCacheDuration(duration time.Duration) {
	cacheDurationMutex.Lock()
	defer cacheDurationMutex.Unlock()
	cacheDuration = duration
}

// NewGRPCClientConn creates a new gRPC client connection with a unary
// interceptor that performs traces of gRPC requests.
func NewGRPCClientConn(
	target string,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {

	opts = append(opts,
		grpc.WithUnaryInterceptor(GRPCClientUnaryInterceptor))

	return grpc.Dial(target, opts...)
}

// GRPCClientUnaryInterceptor is a gRPC unary interceptor for tracing
// outbound gRPC requests
func GRPCClientUnaryInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	seg, err := GetSegmentFromContext(ctx)
	if err != nil {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	subseg := seg.AddNewSubsegment(method)
	subseg.AddRemote()

	sampled := "0"
	if seg.Traced {
		sampled = "1"
	}

	mdctx := metadata.NewContext(ctx, metadata.New(map[string]string{
		mdRootKey:    seg.TraceID,
		mdParentKey:  subseg.ID,
		mdSampledKey: sampled,
	}))

	err = invoker(mdctx, method, req, reply, cc, opts...)

	subseg.Close(err, utils.ErrorType)

	return err
}

// GRPCServerUnaryInterceptor is a gRPC unary interceptor for tracing
// inbound gRPC requests.
func GRPCServerUnaryInterceptor(name string) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		seg := segment.New(name, ctx)
		ctx = AddSegmentToContext(seg, ctx)
		defer seg.Close()

		return handler(ctx, req)
	}
}

// AddSegmentToContext adds a segment reference to a context.Context instance.
func AddSegmentToContext(seg *segment.Segment, ctx context.Context) context.Context {
	sampled := "0"
	if seg.Traced {
		sampled = "1"
	}

	newctx := metadata.NewContext(ctx, metadata.MD{
		mdParentKey:  []string{seg.ParentID},
		mdRootKey:    []string{seg.TraceID},
		mdSampledKey: []string{sampled},
		mdSegmentKey: []string{seg.TraceID},
	})

	cacheDurationMutex.RLock()
	defer cacheDurationMutex.RUnlock()

	segmentCache.Add(seg.TraceID, seg, cacheDuration)

	return newctx
}

// GetSegmentFromContext retrieves a segment from a context.Context instance.
func GetSegmentFromContext(ctx context.Context) (*segment.Segment, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, errors.New("Unable to load metadata from context")
	}

	str := md[mdSegmentKey]
	if len(str) < 1 || len(str[0]) < 1 {
		return nil, errors.New("Segment not found in context")
	}

	segmentTraceID := str[0]

	result, found := segmentCache.Get(segmentTraceID)
	if !found {
		return nil, fmt.Errorf("Segment not found in cache: %s", segmentTraceID)
	}

	seg, ok := result.(*segment.Segment)
	if !ok {
		return nil, errors.New("Unable to assert type for segment from cache")
	}

	return seg, nil
}
