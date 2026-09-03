package playwright

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// APIRequestContextOptions configures a new APIRequestContext.
type APIRequestContextOptions struct {
	BaseURL           *string
	ExtraHTTPHeaders  map[string]string
	HttpCredentials   *HttpCredentials
	UserAgent         *string
	IgnoreHTTPSErrors *bool
	Timeout           *float64
	MaxRedirects      *int
	FailOnStatusCode  *bool
}

// APIRequestContext enables making HTTP requests outside a browser page,
// useful for API testing and setup/teardown operations.
type APIRequestContext struct {
	owner          ChannelOwner
	defaultTimeout *float64
}

// APIResponseTiming holds network timing metrics for an APIResponse.
type APIResponseTiming struct {
	StartTime             float64
	DomainLookupStart     float64
	DomainLookupEnd       float64
	ConnectStart          float64
	SecureConnectionStart float64
	ConnectEnd            float64
	RequestStart          float64
	ResponseStart         float64
	ResponseEnd           float64
}

// APIResponse represents an HTTP response received by APIRequestContext.
type APIResponse struct {
	owner             ChannelOwner
	fetchUID          string
	method            string
	url               string
	status            int
	statusText        string
	headers           map[string]string
	serverAddr        *protocol.RemoteAddr
	securityDetails   *protocol.SecurityDetails
	timing            *protocol.ResourceTiming
	responseEndTiming *float64
	mu                sync.Mutex
	body              []byte
	bodyLoaded        bool
}

// URL returns the response URL.
func (r *APIResponse) URL() string { return r.url }

// Status returns the HTTP status code.
func (r *APIResponse) Status() int { return r.status }

// StatusText returns the HTTP status text.
func (r *APIResponse) StatusText() string { return r.statusText }

// Headers returns response headers as a name→value map.
func (r *APIResponse) Headers() map[string]string { return r.headers }

// OK returns true if the status code is in the 200-299 range.
func (r *APIResponse) OK() bool { return r.status >= 200 && r.status < 300 }

// ServerAddr returns the remote address of the server that served the response,
// or nil if the information is not available.
func (r *APIResponse) ServerAddr() *protocol.RemoteAddr { return r.serverAddr }

// SecurityDetails returns TLS certificate details for HTTPS responses, or nil for HTTP.
func (r *APIResponse) SecurityDetails() *protocol.SecurityDetails { return r.securityDetails }

// Timing returns network timing information for the request, or nil if unavailable.
func (r *APIResponse) Timing() *APIResponseTiming {
	if r.timing == nil {
		return nil
	}
	t := &APIResponseTiming{
		StartTime:             r.timing.StartTime,
		DomainLookupStart:     r.timing.DomainLookupStart,
		DomainLookupEnd:       r.timing.DomainLookupEnd,
		ConnectStart:          r.timing.ConnectStart,
		SecureConnectionStart: r.timing.SecureConnectionStart,
		ConnectEnd:            r.timing.ConnectEnd,
		RequestStart:          r.timing.RequestStart,
		ResponseStart:         r.timing.ResponseStart,
		ResponseEnd:           -1,
	}
	if r.responseEndTiming != nil {
		t.ResponseEnd = *r.responseEndTiming
	}
	return t
}

// Body returns the response body as bytes.
func (r *APIResponse) Body(ctx context.Context) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.bodyLoaded {
		return r.body, nil
	}
	result, err := r.owner.SendMessageRequest(ctx, "fetchResponseBody", map[string]string{
		"fetchUid": r.fetchUID,
	})
	if err != nil {
		return nil, fmt.Errorf("apiResponse.body failed: %w", err)
	}
	var resp struct {
		Binary []byte `json:"binary,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse apiResponse.body response: %w", err)
	}
	r.body = resp.Binary
	r.bodyLoaded = true
	return r.body, nil
}

// Text returns the response body as a UTF-8 string.
func (r *APIResponse) Text(ctx context.Context) (string, error) {
	body, err := r.Body(ctx)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// JSON parses the response body as JSON and returns the result.
func (r *APIResponse) JSON(ctx context.Context) (any, error) {
	body, err := r.Body(ctx)
	if err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("apiResponse.json: failed to parse JSON: %w", err)
	}
	return result, nil
}

// Dispose releases the response body resources on the Playwright server.
func (r *APIResponse) Dispose(ctx context.Context) error {
	_, err := r.owner.SendMessageRequest(ctx, "disposeAPIResponse", map[string]string{
		"fetchUid": r.fetchUID,
	})
	if err != nil {
		return fmt.Errorf("apiResponse.dispose failed: %w", err)
	}
	return nil
}

// FormDataField is a single name-value pair for application/x-www-form-urlencoded requests.
type FormDataField struct {
	Name  string
	Value string
}

// MultipartField is a field in a multipart/form-data request.
// Exactly one of Value or FilePath must be set.
// When FilePath is set the file is read from disk; the MIME type is inferred from the extension.
type MultipartField struct {
	Name     string
	Value    string // text value (mutually exclusive with FilePath)
	FilePath string // path to file to upload (mutually exclusive with Value)
}

// APIFetchOptions configures a single fetch request.
type APIFetchOptions struct {
	Method           *string
	Headers          map[string]string
	Data             []byte
	Params           map[string]string
	FormData         []FormDataField  // application/x-www-form-urlencoded fields
	MultipartData    []MultipartField // multipart/form-data fields
	MaxRedirects     *int
	MaxRetries       *int
	FailOnStatusCode *bool
	Timeout          *float64
}

// apiFetchWire is the IPC wire format for a fetch request.
type apiFetchWire struct {
	FailOnStatusCode *bool                `json:"failOnStatusCode,omitempty"`
	FormData         []protocol.NameValue `json:"formData"`
	Headers          []protocol.NameValue `json:"headers"`
	MaxRedirects     *int                 `json:"maxRedirects,omitempty"`
	MaxRetries       *int                 `json:"maxRetries,omitempty"`
	Method           *string              `json:"method,omitempty"`
	MultipartData    []protocol.FormField `json:"multipartData"`
	Params           []protocol.NameValue `json:"params"`
	PostData         []byte               `json:"postData,omitempty"`
	Timeout          float64              `json:"timeout"`
	Url              string               `json:"url"`
}

// Fetch makes an HTTP request and returns the response.
func (c *APIRequestContext) Fetch(ctx context.Context, url string, opts ...*APIFetchOptions) (*APIResponse, error) {
	req := apiFetchWire{
		Url:           url,
		FormData:      []protocol.NameValue{},
		Headers:       []protocol.NameValue{},
		MultipartData: []protocol.FormField{},
		Params:        []protocol.NameValue{},
		Timeout:       defaultActionTimeoutMs,
	}
	if c.defaultTimeout != nil {
		req.Timeout = *c.defaultTimeout
	}

	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		req.Method = o.Method
		for k, v := range o.Headers {
			req.Headers = append(req.Headers, protocol.NameValue{Name: k, Value: v})
		}
		req.PostData = o.Data
		for k, v := range o.Params {
			req.Params = append(req.Params, protocol.NameValue{Name: k, Value: v})
		}
		for _, f := range o.FormData {
			req.FormData = append(req.FormData, protocol.NameValue{Name: f.Name, Value: f.Value})
		}
		for _, f := range o.MultipartData {
			if f.FilePath != "" {
				fileBytes, readErr := os.ReadFile(f.FilePath)
				if readErr != nil {
					return nil, fmt.Errorf("apiRequestContext.fetch: failed to read file %q: %w", f.FilePath, readErr)
				}
				filename := filepath.Base(f.FilePath)
				mimeType := mime.TypeByExtension(filepath.Ext(filename))
				if idx := strings.Index(mimeType, ";"); idx >= 0 {
					mimeType = strings.TrimSpace(mimeType[:idx])
				}
				if mimeType == "" {
					mimeType = "application/octet-stream"
				}
				req.MultipartData = append(req.MultipartData, protocol.FormField{
					Name: f.Name,
					File: map[string]any{
						"name":     filename,
						"mimeType": mimeType,
						"buffer":   base64.StdEncoding.EncodeToString(fileBytes),
					},
				})
			} else {
				v := f.Value
				req.MultipartData = append(req.MultipartData, protocol.FormField{Name: f.Name, Value: &v, File: nil})
			}
		}
		req.MaxRedirects = o.MaxRedirects
		req.MaxRetries = o.MaxRetries
		req.FailOnStatusCode = o.FailOnStatusCode
		if o.Timeout != nil {
			req.Timeout = *o.Timeout
		}
	}

	effectiveMethod := "GET"
	if req.Method != nil {
		effectiveMethod = *req.Method
	}

	result, err := c.owner.sendWithTimeout(ctx, "fetch", req, req.Timeout)
	if err != nil {
		return nil, fmt.Errorf("apiRequestContext.fetch failed: %w", err)
	}

	var resp struct {
		Response protocol.APIResponse `json:"response"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse fetch response: %w", err)
	}

	return &APIResponse{
		owner:             c.owner,
		fetchUID:          resp.Response.FetchUid,
		method:            effectiveMethod,
		url:               resp.Response.Url,
		status:            resp.Response.Status,
		statusText:        resp.Response.StatusText,
		headers:           nameValuesToHeaders(resp.Response.Headers),
		serverAddr:        resp.Response.ServerAddr,
		securityDetails:   resp.Response.SecurityDetails,
		timing:            resp.Response.Timing,
		responseEndTiming: resp.Response.ResponseEndTiming,
	}, nil
}

// Get makes a GET request to the given URL.
func (c *APIRequestContext) Get(ctx context.Context, url string, opts ...*APIFetchOptions) (*APIResponse, error) {
	return c.Fetch(ctx, url, mergeMethod("GET", opts...))
}

// Post makes a POST request to the given URL.
func (c *APIRequestContext) Post(ctx context.Context, url string, opts ...*APIFetchOptions) (*APIResponse, error) {
	return c.Fetch(ctx, url, mergeMethod("POST", opts...))
}

// Put makes a PUT request to the given URL.
func (c *APIRequestContext) Put(ctx context.Context, url string, opts ...*APIFetchOptions) (*APIResponse, error) {
	return c.Fetch(ctx, url, mergeMethod("PUT", opts...))
}

// Patch makes a PATCH request to the given URL.
func (c *APIRequestContext) Patch(ctx context.Context, url string, opts ...*APIFetchOptions) (*APIResponse, error) {
	return c.Fetch(ctx, url, mergeMethod("PATCH", opts...))
}

// Delete makes a DELETE request to the given URL.
func (c *APIRequestContext) Delete(ctx context.Context, url string, opts ...*APIFetchOptions) (*APIResponse, error) {
	return c.Fetch(ctx, url, mergeMethod("DELETE", opts...))
}

// Head makes a HEAD request to the given URL.
func (c *APIRequestContext) Head(ctx context.Context, url string, opts ...*APIFetchOptions) (*APIResponse, error) {
	return c.Fetch(ctx, url, mergeMethod("HEAD", opts...))
}

// mergeMethod returns options with Method set, copying caller's options if present.
func mergeMethod(method string, opts ...*APIFetchOptions) *APIFetchOptions {
	if len(opts) > 0 && opts[0] != nil {
		o := *opts[0]
		o.Method = &method
		return &o
	}
	return &APIFetchOptions{Method: &method}
}

// StorageState returns the current cookies stored in this request context.
func (c *APIRequestContext) StorageState(ctx context.Context) (*StorageState, error) {
	result, err := c.owner.SendMessageRequest(ctx, "storageState", protocol.APIRequestContextStorageStateRequest{})
	if err != nil {
		return nil, fmt.Errorf("apiRequestContext.storageState failed: %w", err)
	}
	var resp protocol.APIRequestContextStorageStateResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse storageState response: %w", err)
	}
	state := &StorageState{
		Cookies: make([]Cookie, 0, len(resp.Cookies)),
	}
	for _, nc := range resp.Cookies {
		state.Cookies = append(state.Cookies, cookieFromProtocol(nc))
	}
	return state, nil
}

// Dispose releases all resources associated with this request context.
// An optional reason string is included in subsequent-call error messages.
func (c *APIRequestContext) Dispose(ctx context.Context, reason ...string) error {
	req := protocol.APIRequestContextDisposeRequest{}
	if len(reason) > 0 && reason[0] != "" {
		r := reason[0]
		req.Reason = &r
	}
	_, err := c.owner.SendMessageRequest(ctx, "dispose", req)
	if err != nil {
		return fmt.Errorf("apiRequestContext.dispose failed: %w", err)
	}
	return nil
}

// newAPIRequestContextWire is the IPC wire format for the newRequest call.
type newAPIRequestContextWire struct {
	BaseURL            *string                    `json:"baseURL,omitempty"`
	ClientCertificates []any                      `json:"clientCertificates"`
	ExtraHTTPHeaders   []protocol.NameValue       `json:"extraHTTPHeaders"`
	FailOnStatusCode   *bool                      `json:"failOnStatusCode,omitempty"`
	HttpCredentials    []protocol.HttpCredentials `json:"httpCredentials"`
	IgnoreHTTPSErrors  *bool                      `json:"ignoreHTTPSErrors,omitempty"`
	MaxRedirects       *int                       `json:"maxRedirects,omitempty"`
	Proxy              any                        `json:"proxy"`
	StorageState       any                        `json:"storageState"`
	UserAgent          *string                    `json:"userAgent,omitempty"`
}

// NewAPIRequestContext creates a standalone HTTP request context.
func (p *Playwright) NewAPIRequestContext(ctx context.Context, opts ...*APIRequestContextOptions) (*APIRequestContext, error) {
	req := newAPIRequestContextWire{
		ClientCertificates: []any{},
		ExtraHTTPHeaders:   []protocol.NameValue{},
		HttpCredentials:    []protocol.HttpCredentials{},
		Proxy:              nil,
		StorageState:       nil,
	}

	var defaultTimeout *float64
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		req.BaseURL = o.BaseURL
		req.UserAgent = o.UserAgent
		req.IgnoreHTTPSErrors = o.IgnoreHTTPSErrors
		req.FailOnStatusCode = o.FailOnStatusCode
		req.MaxRedirects = o.MaxRedirects
		for k, v := range o.ExtraHTTPHeaders {
			req.ExtraHTTPHeaders = append(req.ExtraHTTPHeaders, protocol.NameValue{Name: k, Value: v})
		}
		if o.HttpCredentials != nil {
			cred := protocol.HttpCredentials{
				Username: o.HttpCredentials.Username,
				Password: o.HttpCredentials.Password,
				Origin:   o.HttpCredentials.Origin,
			}
			if o.HttpCredentials.Send != nil {
				cred.Send = *o.HttpCredentials.Send
			}
			req.HttpCredentials = append(req.HttpCredentials, cred)
		}
		defaultTimeout = o.Timeout
	}

	result, err := p.owner.SendMessageRequest(ctx, "newRequest", req)
	if err != nil {
		return nil, fmt.Errorf("playwright.newAPIRequestContext failed: %w", err)
	}

	var resp protocol.PlaywrightNewRequestResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse newRequest response: %w", err)
	}
	if resp.Request.Guid == "" {
		return nil, fmt.Errorf("newAPIRequestContext: server returned empty GUID")
	}

	return &APIRequestContext{
		owner:          p.owner.child(resp.Request.Guid),
		defaultTimeout: defaultTimeout,
	}, nil
}
