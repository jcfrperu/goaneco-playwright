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
type fallbackOverrides struct {
	url      *string
	method   *string
	headers  map[string]string
	postData []byte
}

// NetworkRequest represents an HTTP request intercepted or observed by the browser.
type NetworkRequest struct {
	owner             ChannelOwner
	guid              string
	url               string
	method            string
	headers           map[string]string
	resourceType      string
	isNavigation      bool
	postData          string
	frameGuid         string
	frame             *Frame
	failure           *string
	page              *Page
	existingResponse  *NetworkResponse
	didFailOrFinish   bool
	fallbackOverrides *fallbackOverrides
}

func (r *NetworkRequest) applyFallbackOverrides(opts *RouteFallbackOptions) {
	if opts == nil {
		return
	}
	if r.fallbackOverrides == nil {
		r.fallbackOverrides = &fallbackOverrides{}
	}
	if opts.URL != nil {
		r.fallbackOverrides.url = opts.URL
	}
	if opts.Method != nil {
		r.fallbackOverrides.method = opts.Method
	}
	if opts.Headers != nil {
		h := make(map[string]string, len(opts.Headers))
		for k, v := range opts.Headers {
			h[k] = v
		}
		r.fallbackOverrides.headers = h
	}
	if opts.PostData != nil {
		switch p := opts.PostData.(type) {
		case string:
			r.fallbackOverrides.postData = []byte(p)
		case []byte:
			r.fallbackOverrides.postData = p
		}
	}
}

// URL returns the full request URL.
func (r *NetworkRequest) URL() string {
	if r.fallbackOverrides != nil && r.fallbackOverrides.url != nil {
		return *r.fallbackOverrides.url
	}
	return r.url
}

// Method returns the HTTP method (GET, POST, etc.).
func (r *NetworkRequest) Method() string {
	if r.fallbackOverrides != nil && r.fallbackOverrides.method != nil {
		return *r.fallbackOverrides.method
	}
	return r.method
}

// Headers returns request headers as a name→value map.
func (r *NetworkRequest) Headers() map[string]string {
	if r.fallbackOverrides != nil && r.fallbackOverrides.headers != nil {
		res := make(map[string]string, len(r.fallbackOverrides.headers))
		for k, v := range r.fallbackOverrides.headers {
			res[strings.ToLower(k)] = v
		}
		return res
	}
	return r.headers
}

// ResourceType returns the resource type (document, stylesheet, image, etc.).
func (r *NetworkRequest) ResourceType() string { return r.resourceType }

// IsNavigationRequest reports whether this is a navigation request.
func (r *NetworkRequest) IsNavigationRequest() bool { return r.isNavigation }

// PostData returns the request body as a string, or empty string if there is none.
func (r *NetworkRequest) PostData() string {
	if r.fallbackOverrides != nil && r.fallbackOverrides.postData != nil {
		return string(r.fallbackOverrides.postData)
	}
	return r.postData
}

// PostDataBuffer returns the request body as raw bytes, or nil if there is none.
func (r *NetworkRequest) PostDataBuffer() []byte {
	if r.fallbackOverrides != nil && r.fallbackOverrides.postData != nil {
		return r.fallbackOverrides.postData
	}
	if r.postData == "" {
		return nil
	}
	return []byte(r.postData)
}

// ExistingResponse returns the cached response if already received, or nil.
func (r *NetworkRequest) ExistingResponse() *NetworkResponse {
	return r.existingResponse
}

// Frame returns the frame that initiated the request, or the page's main frame.
func (r *NetworkRequest) Frame() *Frame {
	if r.isNavigation && r.page == nil {
		panic("Frame for this navigation request is not available, because the request\nwas issued before the frame is created. You can check whether the request\nis a navigation request by calling isNavigationRequest() method.")
	}
	if r.frame != nil {
		return r.frame
	}
	if r.page != nil {
		if r.frameGuid != "" {
			if f := r.page.findFrameByGUID(r.frameGuid); f != nil {
				r.frame = f
				return f
			}
		}
		return r.page.MainFrame()
	}
	panic("Frame for this navigation request is not available, because the request\nwas issued before the frame is created. You can check whether the request\nis a navigation request by calling isNavigationRequest() method.")
}

// Failure returns the error text if the request failed, or nil if successful.
func (r *NetworkRequest) Failure() *string {
	return r.failure
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
	owner       ChannelOwner
	url         string
	status      int
	statusText  string
	headers     map[string]string
	requestGuid string
	request     *NetworkRequest
	page        *Page
}

// URL returns the response URL (may differ from request URL after redirects).
func (r *NetworkResponse) URL() string { return r.url }

// Status returns the HTTP status code.
func (r *NetworkResponse) Status() int { return r.status }

// StatusText returns the HTTP status text (e.g. "OK", "Not Found").
func (r *NetworkResponse) StatusText() string { return r.statusText }

// Headers returns response headers as a name→value map.
func (r *NetworkResponse) Headers() map[string]string { return r.headers }

// Request returns the NetworkRequest associated with this response.
func (r *NetworkResponse) Request() *NetworkRequest {
	if r.request != nil {
		if r.request.page == nil && r.page != nil {
			r.request.page = r.page
		}
		if r.request.existingResponse == nil {
			r.request.existingResponse = r
		}
		return r.request
	}
	if r.requestGuid != "" {
		raw := r.owner.Initializer(r.requestGuid)
		var bCtx *BrowserContext
		if r.page != nil {
			bCtx = r.page.browserContext
		}
		if bCtx != nil {
			r.request = bCtx.requestFromGUID(r.requestGuid, raw)
		} else {
			r.request = networkRequestFrom(r.owner, r.requestGuid, raw)
		}
		if r.request != nil {
			if r.page != nil {
				r.request.page = r.page
			}
			r.request.existingResponse = r
		}
	}
	return r.request
}

// Finished waits for the response to finish and returns nil on success, or an error if failed or closed.
func (r *NetworkResponse) Finished(ctxs ...context.Context) error {
	ctx := context.Background()
	if len(ctxs) > 0 && ctxs[0] != nil {
		ctx = ctxs[0]
	}

	page := r.page
	if page == nil && r.Request() != nil {
		page = r.Request().page
	}
	var bCtx *BrowserContext
	if page != nil {
		bCtx = page.browserContext
	}

	if page != nil && page.IsClosed() {
		return fmt.Errorf("response.Finished: %s", page.effectiveCloseReason())
	}
	if bCtx != nil && bCtx.IsClosed() {
		return fmt.Errorf("response.Finished: %s", bCtx.effectiveCloseReason())
	}

	req := r.Request()
	if req != nil && req.didFailOrFinish {
		if req.failure != nil {
			return fmt.Errorf("request failed: %s", *req.failure)
		}
		return nil
	}

	reqGuid := r.requestGuid
	if req != nil && req.guid != "" {
		reqGuid = req.guid
	}

	done := make(chan error, 1)
	var unsubList []func()
	defer func() {
		for _, u := range unsubList {
			if u != nil {
				u()
			}
		}
	}()

	if page != nil {
		unsubReqFin := page.OnRequestFinished(func(q *NetworkRequest) {
			if q.guid == reqGuid {
				select {
				case done <- nil:
				default:
				}
			}
		})
		unsubList = append(unsubList, unsubReqFin)

		unsubReqFail := page.OnRequestFailed(func(q *NetworkRequest) {
			if q.guid == reqGuid {
				var err error
				if q.failure != nil {
					err = fmt.Errorf("request failed: %s", *q.failure)
				} else {
					err = fmt.Errorf("request failed")
				}
				select {
				case done <- err:
				default:
				}
			}
		})
		unsubList = append(unsubList, unsubReqFail)

		unsubClose := page.OnClose(func(*Page) {
			select {
			case done <- fmt.Errorf("response.Finished: %s", page.effectiveCloseReason()):
			default:
			}
		})
		unsubList = append(unsubList, unsubClose)
	}

	if bCtx != nil {
		unsubCtxClose := bCtx.OnClose(func(*BrowserContext) {
			reason := "Target closed"
			if page != nil {
				reason = page.effectiveCloseReason()
			} else {
				reason = bCtx.effectiveCloseReason()
			}
			select {
			case done <- fmt.Errorf("response.Finished: %s", reason):
			default:
			}
		})
		unsubList = append(unsubList, unsubCtxClose)
	}

	if page != nil && page.IsClosed() {
		return fmt.Errorf("response.Finished: %s", page.effectiveCloseReason())
	}
	if bCtx != nil && bCtx.IsClosed() {
		return fmt.Errorf("response.Finished: %s", bCtx.effectiveCloseReason())
	}
	if req != nil && req.didFailOrFinish {
		if req.failure != nil {
			return fmt.Errorf("request failed: %s", *req.failure)
		}
		return nil
	}

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ServerAddr returns the remote address of the server that served this response.
func (r *NetworkResponse) ServerAddr(ctx context.Context) (*protocol.RemoteAddr, error) {
	result, err := r.owner.SendMessageRequest(ctx, "serverAddr", protocol.ResponseServerAddrRequest{})
	if err != nil {
		return nil, fmt.Errorf("networkResponse.serverAddr failed: %w", err)
	}
	var response protocol.ResponseServerAddrResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("failed to parse networkResponse.serverAddr response: %w", err)
	}
	return response.Value, nil
}

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

// HeaderValue represents an HTTP header name and value pair.
type HeaderValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type rawHeaders struct {
	headersArray []HeaderValue
	headersMap   map[string][]string
}

func newRawHeaders(headers []HeaderValue) *rawHeaders {
	r := &rawHeaders{
		headersArray: headers,
		headersMap:   make(map[string][]string),
	}
	for _, h := range headers {
		name := strings.ToLower(h.Name)
		r.headersMap[name] = append(r.headersMap[name], h.Value)
	}
	return r
}

func (r *rawHeaders) get(name string) *string {
	values := r.headersMap[strings.ToLower(name)]
	if len(values) == 0 {
		return nil
	}
	sep := ", "
	if strings.ToLower(name) == "set-cookie" {
		sep = "\n"
	}
	res := strings.Join(values, sep)
	return &res
}

func (r *rawHeaders) getAll(name string) []string {
	vals, ok := r.headersMap[strings.ToLower(name)]
	if !ok {
		return []string{}
	}
	res := make([]string, len(vals))
	copy(res, vals)
	return res
}

func (r *rawHeaders) headers() map[string]string {
	result := make(map[string]string, len(r.headersMap))
	for name := range r.headersMap {
		if val := r.get(name); val != nil {
			result[name] = *val
		}
	}
	return result
}

// AllHeaders returns all request headers as a map, fetching raw headers if needed.
func (r *NetworkRequest) AllHeaders(ctx context.Context) (map[string]string, error) {
	if r.fallbackOverrides != nil && r.fallbackOverrides.headers != nil {
		res := make(map[string]string, len(r.fallbackOverrides.headers))
		for k, v := range r.fallbackOverrides.headers {
			res[strings.ToLower(k)] = v
		}
		return res, nil
	}
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.headers(), nil
}

// HeadersArray returns all request headers as an array of name-value pairs.
func (r *NetworkRequest) HeadersArray(ctx context.Context) ([]HeaderValue, error) {
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.headersArray, nil
}

// HeaderValue returns the value of the named header, or nil if not present.
func (r *NetworkRequest) HeaderValue(ctx context.Context, name string) (*string, error) {
	if r.fallbackOverrides != nil && r.fallbackOverrides.headers != nil {
		lowerName := strings.ToLower(name)
		for k, v := range r.fallbackOverrides.headers {
			if strings.ToLower(k) == lowerName {
				val := v
				return &val, nil
			}
		}
		return nil, nil
	}
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.get(name), nil
}

// HeaderValues returns all values of the named header.
func (r *NetworkRequest) HeaderValues(ctx context.Context, name string) ([]string, error) {
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.getAll(name), nil
}

func (r *NetworkRequest) getRawHeaders(ctx context.Context) (*rawHeaders, error) {
	if r.guid == "" {
		return newRawHeaders(nil), nil
	}
	result, err := r.owner.SendMessageRequest(ctx, "rawRequestHeaders", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("networkRequest.rawRequestHeaders failed: %w", err)
	}
	var resp struct {
		Headers []HeaderValue `json:"headers"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse rawRequestHeaders: %w", err)
	}
	return newRawHeaders(resp.Headers), nil
}

// AllHeaders returns all response headers as a map.
func (r *NetworkResponse) AllHeaders(ctx context.Context) (map[string]string, error) {
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.headers(), nil
}

// HeadersArray returns all response headers as an array of name-value pairs.
func (r *NetworkResponse) HeadersArray(ctx context.Context) ([]HeaderValue, error) {
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.headersArray, nil
}

// HeaderValue returns the value of the named response header, or nil if not present.
func (r *NetworkResponse) HeaderValue(ctx context.Context, name string) (*string, error) {
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.get(name), nil
}

// HeaderValues returns all values of the named response header.
func (r *NetworkResponse) HeaderValues(ctx context.Context, name string) ([]string, error) {
	raw, err := r.getRawHeaders(ctx)
	if err != nil {
		return nil, err
	}
	return raw.getAll(name), nil
}

// HttpVersion returns the HTTP version of this response (e.g. "HTTP/1.1").
func (r *NetworkResponse) HttpVersion(ctx context.Context) (string, error) {
	if r.owner.guid == "" {
		return "", nil
	}
	result, err := r.owner.SendMessageRequest(ctx, "httpVersion", struct{}{})
	if err != nil {
		return "", fmt.Errorf("networkResponse.httpVersion failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse httpVersion response: %w", err)
	}
	return resp.Value, nil
}

func (r *NetworkResponse) getRawHeaders(ctx context.Context) (*rawHeaders, error) {
	if r.owner.guid == "" {
		return newRawHeaders(nil), nil
	}
	result, err := r.owner.SendMessageRequest(ctx, "rawResponseHeaders", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("networkResponse.rawResponseHeaders failed: %w", err)
	}
	var resp struct {
		Headers []HeaderValue `json:"headers"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse rawResponseHeaders: %w", err)
	}
	return newRawHeaders(resp.Headers), nil
}

// networkRequestInit is the wire format of the Request channel-object initializer.
type networkRequestInit struct {
	URL          string               `json:"url"`
	Method       string               `json:"method"`
	Headers      []protocol.NameValue `json:"headers"`
	ResourceType string               `json:"resourceType"`
	IsNavigation bool                 `json:"isNavigationRequest"`
	PostData     string               `json:"postData,omitempty"`
	Frame        struct {
		Guid string `json:"guid"`
	} `json:"frame"`
}

// networkResponseInit is the wire format of the Response channel-object initializer.
type networkResponseInit struct {
	URL        string               `json:"url"`
	Status     int                  `json:"status"`
	StatusText string               `json:"statusText"`
	Headers    []protocol.NameValue `json:"headers"`
	Request    struct {
		Guid string `json:"guid"`
	} `json:"request"`
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
		frameGuid:    init.Frame.Guid,
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
		url:         init.URL,
		status:      init.Status,
		statusText:  init.StatusText,
		headers:     nameValuesToHeaders(init.Headers),
		requestGuid: init.Request.Guid,
	}
	if guid != "" {
		r.owner = parent.child(guid)
	}
	return r
} // requestFromGUID returns a cached or newly created NetworkRequest for the given GUID.
func (c *BrowserContext) requestFromGUID(guid string, raw json.RawMessage) *NetworkRequest {
	if c == nil {
		return networkRequestFrom(ChannelOwner{}, guid, raw)
	}
	c.mu.RLock()
	existing := c.requestsByGUID[guid]
	c.mu.RUnlock()
	if existing != nil {
		return existing
	}

	req := networkRequestFrom(c.owner, guid, raw)
	if req == nil {
		return nil
	}

	c.mu.Lock()
	if c.requestsByGUID == nil {
		c.requestsByGUID = make(map[string]*NetworkRequest)
	}
	if existing = c.requestsByGUID[guid]; existing == nil {
		c.requestsByGUID[guid] = req
		existing = req
	}
	c.mu.Unlock()
	return existing
}
