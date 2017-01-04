package handlers

import (
	"bytes"
	"github.com/goguardian/aws-xray-go/segment"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"
)

func TestAddSegmentToContext(t *testing.T) {
	seg := segment.New("test", nil)
	seg.AddFault()
	seg.AddThrottle()

	ctx := AddSegmentToContext(seg, context.TODO())

	resultSeg, err := GetSegmentFromContext(ctx)
	if err != nil {
		t.Errorf("Error retrieving segment from context: %s", err.Error())
	}

	if !resultSeg.Fault {
		t.Error("Segment should be marked fault")
	}

	if !resultSeg.Throttle {
		t.Error("Segment should be marked throttle")
	}
}

func TestGetSegmentFromContext(t *testing.T) {
	ctx := context.TODO()

	_, err := GetSegmentFromContext(ctx)
	if err == nil {
		t.Error("A context with no segment should error")
	}

	ctx = metadata.NewContext(context.TODO(), metadata.New(map[string]string{
		mdSegmentKey: "",
	}))

	_, err = GetSegmentFromContext(ctx)
	if err == nil {
		t.Error("A context with no segment metadata key should error")
	}

	ctx = metadata.NewContext(context.TODO(), metadata.New(map[string]string{
		mdSegmentKey: "123",
	}))

	_, err = GetSegmentFromContext(ctx)
	if err == nil {
		t.Error("A context with a segment key not found in cache should error")
	}

	segmentCache.Add("123", bytes.Buffer{}, 1*time.Minute)
	_, err = GetSegmentFromContext(ctx)
	if err == nil {
		t.Error("A bad cache entry should not be able to be type asserted")
	}
}
