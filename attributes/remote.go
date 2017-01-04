package attributes

// Remote represents a remote call with request and response data.
type Remote struct {
	Request  *RemoteRequest  `json:"request,omitempty"`
	Response *RemoteResponse `json:"response,omitemtpy"`
}

// RemoteRequest represents remote requests details.
type RemoteRequest struct {
	URL    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
	Traced bool   `json:"traced"`
}

// RemoteResponse represents remote response details.
type RemoteResponse struct {
	Status        int `json:"status"`
	ContentLength int `json:"content_length"`
}
