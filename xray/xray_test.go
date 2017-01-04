package xray

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc"
)

const name = "xray"

func TestNewContext(t *testing.T) {
	ctx := NewContext(name, context.Background())

	seg, err := GetSegment(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if seg == nil {
		t.Error("Context Segment should not be nil")
		return
	}
	if seg.Name != name {
		t.Errorf("Context Segment should equal '%s'", name)
	}

	Close(ctx)

	if seg.InProgress {
		t.Error("Context Segment should not be in progress")
	}
}

func TestAddLocalHTTP(t *testing.T) {
	ctx := NewContext(name, context.Background())

	req := &http.Request{}
	req = req.WithContext(ctx)

	if err := AddLocalHTTP(req); err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {
	ctx := context.TODO()

	if err := Close(ctx); err == nil {
		t.Error("Context without segment should error on close")
	}
}

func TestGetGRPCClientConn(t *testing.T) {
	_, err := GetGRPCClientConn("127.0.0.1:1000", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
}
