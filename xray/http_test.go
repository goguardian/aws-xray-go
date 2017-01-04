package xray

import "testing"

func TestGetHTTPClient(t *testing.T) {
	ctx := NewContext(name, nil)

	httpClient, err := GetHTTPClient(ctx)
	if httpClient == nil {
		t.Error("Context HTTP Client should not be nil")
	}
	if err != nil {
		t.Error(err)
	}
}
