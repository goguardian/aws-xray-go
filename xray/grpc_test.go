package xray

import "testing"

func TestNewGRPCServer(t *testing.T) {
	server := NewGRPCServer("name")
	if server == nil {
		t.Error("gRPC Server instance should not be nil")
	}
}
