package utils

import (
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func TestGetIDsFromContext(t *testing.T) {
	tests := []struct {
		ctx            context.Context
		request        *http.Request
		expectRootID   string
		expectParentID string
		expectSampled  string
	}{
		{
			ctx:            nil,
			expectRootID:   "",
			expectParentID: "",
		},
		{
			request:        &http.Request{},
			expectRootID:   "",
			expectParentID: "",
		},
		{
			request: &http.Request{
				Header: http.Header{
					XRayHeader: []string{"Root=123; Parent=456"},
				},
			},
			expectRootID:   "123",
			expectParentID: "456",
			expectSampled:  "",
		},
		{
			request: &http.Request{
				Header: http.Header{
					strings.ToLower(XRayHeader): []string{"Root=123; Parent=456; Sampled=1"},
				},
			},
			expectRootID:   "123",
			expectParentID: "456",
			expectSampled:  "1",
		},
		{
			ctx:            context.TODO(),
			expectRootID:   "",
			expectParentID: "",
		},
		{
			ctx: metadata.NewContext(context.TODO(), metadata.New(map[string]string{
				"xray-rootid":   "123",
				"xray-parentid": "456",
				"xray-sampled":  "0",
			})),
			expectRootID:   "123",
			expectParentID: "456",
			expectSampled:  "0",
		},
	}

	for _, test := range tests {
		if test.request != nil {
			test.ctx = ContextFromHeaders(test.request)
		}

		rootID, parentID, sampled := GetIDsFromContext(test.ctx)
		if rootID != test.expectRootID {
			t.Errorf("Expected root ID '%s', got '%s'", test.expectRootID, rootID)
		}

		if parentID != test.expectParentID {
			t.Errorf("Expected parent ID '%s', got '%s'",
				test.expectParentID, parentID)
		}

		if sampled != test.expectSampled {
			t.Errorf("Expected sampled '%s', got '%s'", test.expectSampled, sampled)
		}
	}
}
