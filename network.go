package playwright

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// NetworkRequest represents an HTTP request intercepted or observed by the browser.
type NetworkRequest struct {
	owner        ChannelOwner
	guid         string
	url          string
	method       string
	headers      map[string]string
	resourceType string
	isNavigation bool
	postData     string
}

// URL returns the full request URL.
func (r *NetworkRequest) URL() string { return r.url }

// Method returns the HTTP method (GET, POST, etc.).
func (r *NetworkRequest) Method() string { return r.method }

// Headers returns request headers as a name→value map.
func (r *NetworkRequest) Headers() map[string]string { return r.headers }

// ResourceType returns the resource type (document, stylesheet, image, etc.).
func (r *NetworkRequest) ResourceType() string { return r.resourceType }

// IsNavigationRequest reports whether this is a navigation request.
func (r *NetworkRequest) IsNavigationRequest() bool { return r.isNavigation }

// PostData returns the request body as a string, or empty string if there is none.
func (r *NetworkRequest) PostData() string { return r.postData }

// PostDataBuffer returns the request body as raw bytes, or nil if there is none.
func (r *NetworkRequest) PostDataBuffer() []byte {
	if r.postData == "" {
		return nil
	}
	return []byte(r.postData)
}

// RequestSizes holds byte-level size information for a request and its response.
type RequestSizes struct {
	RequestBodySize     int `json:"requestBodySize"`
	RequestHeadersSize  int `json:"requestHeadersSize"`
	ResponseBodySize    int `json:"responseBodySize"`
	ResponseHeadersSize int `json:"responseHeadersSize"`
}

// Sizes returns byte-level size information for this request and its response.
func (r *NetworkRequest) Sizes(ctx context.Context) (*RequestSizes, error) {
	if r.guid == "" {
		return nil, nil
	}
	result, err := r.owner.SendMessageRequest(ctx, "requestSizes", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("networkRequest.sizes failed: %w", err)
	}
	var resp struct {
		Sizes RequestSizes `json:"sizes"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse requestSizes response: %w", err)
	}
	return &resp.Sizes, nil
}

// Response returns the matching NetworkResponse for this request, or nil if none.
func (r *NetworkRequest) Response(ctx context.Context) (*NetworkResponse, error) {
	if r.guid == "" {
		return nil, nil
	}
	result, err := r.owner.SendMessageRequest(ctx, "response", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("networkRequest.response failed: %w", err)
	}
	var resp struct {
		Response *protocol.Response `json:"response,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse networkRequest.response: %w", err)
	}
	if resp.Response == nil {
		return nil, nil
	}
	raw := r.owner.Initializer(resp.Response.Guid)
	return networkResponseFrom(r.owner, resp.Response.Guid, raw), nil
}

// NetworkResponse represents an HTTP response observed by the browser.
type NetworkResponse struct {
	owner      ChannelOwner
	url        string
	status     int
	statusText string
	headers    map[string]string
}

// URL returns the response URL (may differ from request URL after redirects).
func (r *NetworkResponse) URL() string { return r.url }

// Status returns the HTTP status code.
func (r *NetworkResponse) Status() int { return r.status }

// StatusText returns the HTTP status text (e.g. "OK", "Not Found").
func (r *NetworkResponse) StatusText() string { return r.statusText }

// Headers returns response headers as a name→value map.
func (r *NetworkResponse) Headers() map[string]string { return r.headers }

// Body returns the response body as raw bytes.
func (r *NetworkResponse) Body(ctx context.Context) ([]byte, error) {
	result, err := r.owner.SendMessageRequest(ctx, "body", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("networkResponse.body failed: %w", err)
	}
	var resp struct {
		Binary []byte `json:"binary"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse networkResponse.body response: %w", err)
	}
	return resp.Binary, nil
}

// OK returns true if the HTTP status code is in the range 200–299.
func (r *NetworkResponse) OK() bool {
	return r.status >= 200 && r.status < 300
}

// Text returns the response body as a UTF-8 string.
func (r *NetworkResponse) Text(ctx context.Context) (string, error) {
	body, err := r.Body(ctx)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// networkRequestInit is the wire format of the Request channel-object initializer.
type networkRequestInit struct {
	URL          string               `json:"url"`
	Method       string               `json:"method"`
	Headers      []protocol.NameValue `json:"headers"`
	ResourceType string               `json:"resourceType"`
	IsNavigation bool                 `json:"isNavigationRequest"`
	PostData     string               `json:"postData,omitempty"`
}

// networkResponseInit is the wire format of the Response channel-object initializer.
type networkResponseInit struct {
	URL        string               `json:"url"`
	Status     int                  `json:"status"`
	StatusText string               `json:"statusText"`
	Headers    []protocol.NameValue `json:"headers"`
}

func nameValuesToMap(pairs []protocol.NameValue) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, nv := range pairs {
		m[strings.ToLower(nv.Name)] = nv.Value
	}
	return m
}

// networkRequestFrom creates a NetworkRequest with an IPC channel for Response() access.
// Returns nil if the initializer payload cannot be parsed.
func networkRequestFrom(parent ChannelOwner, guid string, raw json.RawMessage) *NetworkRequest {
	var init networkRequestInit
	if err := json.Unmarshal(raw, &init); err != nil {
		return nil
	}
	postData := init.PostData
	// The Playwright wire protocol encodes postData as base64 when non-empty.
	if decoded, err := base64.StdEncoding.DecodeString(postData); err == nil && postData != "" {
		postData = string(decoded)
	}
	r := &NetworkRequest{
		url:          init.URL,
		method:       init.Method,
		headers:      nameValuesToHeaders(init.Headers),
		resourceType: init.ResourceType,
		isNavigation: init.IsNavigation,
		postData:     postData,
	}
	if guid != "" {
		r.guid = guid
		r.owner = parent.child(guid)
	}
	return r
}

// networkResponseFrom creates a NetworkResponse with an IPC channel for body/text access.
// Returns nil if the initializer payload cannot be parsed.
func networkResponseFrom(parent ChannelOwner, guid string, raw json.RawMessage) *NetworkResponse {
	var init networkResponseInit
	if err := json.Unmarshal(raw, &init); err != nil {
		return nil
	}
	r := &NetworkResponse{
		url:        init.URL,
		status:     init.Status,
		statusText: init.StatusText,
		headers:    nameValuesToHeaders(init.Headers),
	}
	if guid != "" {
		r.owner = parent.child(guid)
	}
	return r
}
