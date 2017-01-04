package utils

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	// XRayHeader represents an AWS X-Ray trace HTTP header key.
	XRayHeader   = "X-Amzn-Trace-Id"
	mdParentKey  = "xray-parentid"
	mdRootKey    = "xray-rootid"
	mdSampledKey = "xray-sampled"
	mdSegmentKey = "xray-segment"
)

// GetIDsFromContext returns the root ID, parent ID, and whether the request
// was sampled based on context metadata.
func GetIDsFromContext(ctx context.Context) (rootID, parentID, sampled string) {
	if ctx == nil {
		return
	}

	data, ok := metadata.FromContext(ctx)
	if !ok {
		return
	}

	rootIDs, ok := data[mdRootKey]
	if ok && len(rootIDs) > 0 {
		rootID = rootIDs[0]
	}

	parentIDs, ok := data[mdParentKey]
	if ok && len(parentIDs) > 0 {
		parentID = parentIDs[0]
	}

	sampleds, ok := data[mdSampledKey]
	if ok && len(sampleds) > 0 {
		sampled = sampleds[0]
	}

	return
}

// ContextFromHeaders parses the root ID, parent ID, and whether the
// request was sampled based on the HTTP request X-Ray header and returns a
// new context containing the metadata.
func ContextFromHeaders(r *http.Request) context.Context {
	header, rootID, parentID, sampled := "", "", "", ""

	headers, ok := r.Header[XRayHeader]
	if ok && len(headers) > 0 {
		header = headers[0]
	}

	if header == "" {
		headers, ok = r.Header[strings.ToLower(XRayHeader)]
		if ok && len(headers) > 0 {
			header = headers[0]
		}
	}

	if header == "" {
		return r.Context()
	}

	pairs := strings.Split(strings.Replace(header, " ", "", -1), ";")
	for _, pair := range pairs {
		set := strings.Split(pair, "=")
		if len(set) == 2 {
			if set[0] == "Root" {
				rootID = set[1]
			} else if set[0] == "Parent" {
				parentID = set[1]
			} else if set[0] == "Sampled" {
				sampled = set[1]
			}
		}
	}

	newMD := metadata.New(map[string]string{
		mdRootKey:    rootID,
		mdParentKey:  parentID,
		mdSampledKey: sampled,
	})

	oldMD, found := metadata.FromContext(r.Context())
	if found {
		newMD = metadata.Join(newMD, oldMD)
	}

	return metadata.NewContext(r.Context(), newMD)
}
