package playwright

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// RouteHandler is the function called when a request matches a route pattern.
type RouteHandler func(route *Route)

// Route represents an intercepted network request that can be fulfilled,
// continued, or aborted.
type Route struct {
	owner          ChannelOwner
	mu             sync.Mutex
	doneCh         chan<- bool // non-nil when managed by routeRouter; receives false=handled, true=fallback
	reqOnce        sync.Once
	cachedReq      *NetworkRequest
	browserContext *BrowserContext
}

// setDoneCh arms the route with a completion channel (called by routeRouter before each handler).
func (r *Route) setDoneCh(ch chan<- bool) {
	r.mu.Lock()
	r.doneCh = ch
	r.mu.Unlock()
}

// signalDone notifies the managing router that this route has been resolved.
// fallback=false means the route was handled; fallback=true means Fallback was called.
func (r *Route) signalDone(fallback bool) {
	r.mu.Lock()
	ch := r.doneCh
	r.mu.Unlock()
	if ch != nil {
		select {
		case ch <- fallback:
		default:
		}
	}
}

// routeInit is the wire format of the Route channel object initializer.
type routeInit struct {
	Request struct {
		Guid string `json:"guid"`
	} `json:"request"`
}

// Request returns the NetworkRequest associated with this route.
// The result is cached after the first call.
func (r *Route) Request() *NetworkRequest {
	r.reqOnce.Do(func() {
		raw := r.owner.Initializer(r.owner.guid)
		var init routeInit
		if err := json.Unmarshal(raw, &init); err != nil || init.Request.Guid == "" {
			return
		}
		reqRaw := r.owner.Initializer(init.Request.Guid)
		if r.browserContext != nil {
			r.cachedReq = r.browserContext.requestFromGUID(init.Request.Guid, reqRaw)
		} else {
			r.cachedReq = networkRequestFrom(r.owner, init.Request.Guid, reqRaw)
		}
	})
	return r.cachedReq
}

// RouteFulfillOptions specifies the response for Route.Fulfill.
type RouteFulfillOptions struct {
	Status      *int
	Headers     map[string]string
	Body        *string
	BodyBytes   []byte
	ContentType *string
	Response    *APIResponse
}

// RouteContinueOptions specifies modifications for Route.Continue.
type RouteContinueOptions struct {
	Headers  map[string]string
	Method   *string
	PostData []byte
	URL      *string
}

// SetURL sets the URL override.
func (o *RouteContinueOptions) SetURL(url string) *RouteContinueOptions {
	o.URL = &url
	return o
}

// SetMethod sets the Method override.
func (o *RouteContinueOptions) SetMethod(method string) *RouteContinueOptions {
	o.Method = &method
	return o
}

// SetHeaders sets the Headers override.
func (o *RouteContinueOptions) SetHeaders(headers map[string]string) *RouteContinueOptions {
	o.Headers = headers
	return o
}

// SetPostData sets the PostData override.
func (o *RouteContinueOptions) SetPostData(postData []byte) *RouteContinueOptions {
	o.PostData = postData
	return o
}

// Fulfill completes the route with a custom HTTP response.
func (r *Route) Fulfill(ctx context.Context, opts *RouteFulfillOptions) error {
	req := protocol.RouteFulfillRequest{
		Headers: []protocol.NameValue{},
	}

	if opts != nil {
		if opts.Body != nil && len(opts.BodyBytes) > 0 {
			return fmt.Errorf("route.Fulfill: cannot specify both Body and BodyBytes")
		}

		status := opts.Status
		var headersMap map[string]string
		if opts.Response != nil {
			if status == nil {
				s := opts.Response.Status()
				status = &s
			}
			headersMap = make(map[string]string, len(opts.Response.Headers()))
			for k, v := range opts.Response.Headers() {
				headersMap[strings.ToLower(k)] = v
			}
			if opts.Response.fetchUID != "" {
				uid := opts.Response.fetchUID
				req.FetchResponseUid = &uid
			} else {
				b, err := opts.Response.Body(ctx)
				if err != nil {
					return fmt.Errorf("route.Fulfill: failed to read response body: %w", err)
				}
				encoded := base64.StdEncoding.EncodeToString(b)
				req.Body = &encoded
				req.IsBase64 = protocol.Bool(true)
			}
		}

		if opts.Headers != nil {
			if headersMap == nil {
				headersMap = make(map[string]string, len(opts.Headers))
			}
			for k, v := range opts.Headers {
				headersMap[strings.ToLower(k)] = v
			}
		}

		req.Status = status
		if opts.Body != nil {
			req.Body = opts.Body
			req.IsBase64 = nil
			req.FetchResponseUid = nil
		}
		if len(opts.BodyBytes) > 0 {
			encoded := base64.StdEncoding.EncodeToString(opts.BodyBytes)
			req.Body = &encoded
			req.IsBase64 = protocol.Bool(true)
			req.FetchResponseUid = nil
		}

		if opts.ContentType != nil {
			if headersMap == nil {
				headersMap = make(map[string]string)
			}
			headersMap["content-type"] = *opts.ContentType
		}

		headers := make([]protocol.NameValue, 0, len(headersMap))
		for k, v := range headersMap {
			headers = append(headers, protocol.NameValue{Name: k, Value: v})
		}
		req.Headers = headers
	}

	_, err := r.owner.SendMessageRequest(ctx, "fulfill", req)
	if err != nil {
		return fmt.Errorf("route.fulfill failed: %w", err)
	}
	r.signalDone(false)
	return nil
}

// RouteFetchOptions configures Route.Fetch.
type RouteFetchOptions struct {
	Headers      map[string]string
	MaxRedirects *int
	Method       *string
	PostData     any // string or []byte
	Timeout      *float64
	URL          *string
}

// Fetch performs the intercepted request and returns the resulting APIResponse.
func (r *Route) Fetch(ctx context.Context, opts ...*RouteFetchOptions) (*APIResponse, error) {
	bCtx := r.browserContext
	if bCtx == nil && r.cachedReq != nil && r.cachedReq.page != nil {
		bCtx = r.cachedReq.page.browserContext
	}
	if bCtx == nil {
		return nil, fmt.Errorf("route.Fetch: no browser context available")
	}

	req := r.Request()
	fetchURL := req.URL()
	method := req.Method()
	headers := req.Headers()
	var postData []byte
	if req.postData != "" {
		postData = []byte(req.postData)
	}
	var maxRedirects *int
	var timeout *float64

	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.URL != nil {
			fetchURL = *o.URL
		}
		if o.Method != nil {
			method = *o.Method
		}
		if o.Headers != nil {
			headers = o.Headers
		}
		if o.PostData != nil {
			switch p := o.PostData.(type) {
			case string:
				postData = []byte(p)
			case []byte:
				postData = p
			}
		}
		if o.MaxRedirects != nil {
			maxRedirects = o.MaxRedirects
		}
		if o.Timeout != nil {
			timeout = o.Timeout
		}
	}

	apiOptions := &APIFetchOptions{
		Method:       &method,
		Headers:      headers,
		Data:         postData,
		MaxRedirects: maxRedirects,
		Timeout:      timeout,
	}

	return bCtx.Request().Fetch(ctx, fetchURL, apiOptions)
}

func (r *Route) buildContinueRequest(opts *RouteContinueOptions, isFallback bool) protocol.RouteContinueRequest {
	req := protocol.RouteContinueRequest{
		IsFallback: isFallback,
	}

	var url *string
	var method *string
	var headers map[string]string
	var postData []byte

	if rReq := r.Request(); rReq != nil && rReq.fallbackOverrides != nil {
		url = rReq.fallbackOverrides.url
		method = rReq.fallbackOverrides.method
		if rReq.fallbackOverrides.headers != nil {
			headers = make(map[string]string, len(rReq.fallbackOverrides.headers))
			for k, v := range rReq.fallbackOverrides.headers {
				headers[k] = v
			}
		}
		if rReq.fallbackOverrides.postData != nil {
			postData = rReq.fallbackOverrides.postData
		}
	}

	if opts != nil {
		if opts.URL != nil {
			url = opts.URL
		}
		if opts.Method != nil {
			method = opts.Method
		}
		if opts.Headers != nil {
			headers = opts.Headers
		}
		if opts.PostData != nil {
			postData = opts.PostData
		}
	}

	req.Url = url
	req.Method = method
	if postData != nil {
		encoded := base64.StdEncoding.EncodeToString(postData)
		req.PostData = &encoded
	}
	if headers != nil {
		hvs := make([]protocol.NameValue, 0, len(headers))
		for k, v := range headers {
			hvs = append(hvs, protocol.NameValue{Name: k, Value: v})
		}
		req.Headers = &hvs
	}

	return req
}

// Continue passes the request to the network with optional modifications.
func (r *Route) Continue(ctx context.Context, options ...*RouteContinueOptions) error {
	var opts *RouteContinueOptions
	if len(options) > 0 {
		opts = options[0]
	}
	req := r.buildContinueRequest(opts, false)

	_, err := r.owner.SendMessageRequest(ctx, "continue", req)
	if err != nil {
		return fmt.Errorf("route.continue failed: %w", err)
	}
	r.signalDone(false)
	return nil
}

// Resume is an alias for Continue(ctx), matching playwright-java.
func (r *Route) Resume(ctx context.Context) error {
	return r.Continue(ctx)
}

// Abort cancels the intercepted request.
// An optional errorCode can be provided (defaults to "failed").
func (r *Route) Abort(ctx context.Context, errorCode ...string) error {
	req := protocol.RouteAbortRequest{}
	if len(errorCode) > 0 && errorCode[0] != "" {
		code := errorCode[0]
		req.ErrorCode = &code
	} else {
		code := "failed"
		req.ErrorCode = &code
	}

	_, err := r.owner.SendMessageRequest(ctx, "abort", req)
	if err != nil {
		return fmt.Errorf("route.abort failed: %w", err)
	}
	r.signalDone(false)
	return nil
}

// continueAsFallback sends "continue" with isFallback=true to the server, handing off to the
// next routing layer (e.g., page routes falling back to context routes, or context routes
// falling back to the network).
func (r *Route) continueAsFallback(ctx context.Context) error {
	req := r.buildContinueRequest(nil, true)
	_, err := r.owner.SendMessageRequest(ctx, "continue", req)
	if err != nil {
		return fmt.Errorf("route.continueAsFallback failed: %w", err)
	}
	return nil
}

// RouteFallbackOptions configures optional overrides when falling back a route.
type RouteFallbackOptions struct {
	Headers  map[string]string
	Method   *string
	PostData any // string or []byte
	URL      *string
}

// SetURL sets the URL override.
func (o *RouteFallbackOptions) SetURL(url string) *RouteFallbackOptions {
	o.URL = &url
	return o
}

// SetMethod sets the Method override.
func (o *RouteFallbackOptions) SetMethod(method string) *RouteFallbackOptions {
	o.Method = &method
	return o
}

// SetHeaders sets the Headers override.
func (o *RouteFallbackOptions) SetHeaders(headers map[string]string) *RouteFallbackOptions {
	o.Headers = headers
	return o
}

// SetPostData sets the PostData override.
func (o *RouteFallbackOptions) SetPostData(postData any) *RouteFallbackOptions {
	o.PostData = postData
	return o
}

// Fallback defers handling of this route to the next registered handler. When managed by a
// routeRouter, it signals the router to proceed to the next entry. When called outside a
// router context, it sends "continue" with isFallback=true directly to the server.
func (r *Route) Fallback(ctx context.Context, opts ...*RouteFallbackOptions) error {
	if len(opts) > 0 && opts[0] != nil {
		if req := r.Request(); req != nil {
			req.applyFallbackOverrides(opts[0])
		}
	}

	r.mu.Lock()
	ch := r.doneCh
	r.mu.Unlock()
	if ch != nil {
		select {
		case ch <- true:
		default:
		}
		return nil
	}
	return r.continueAsFallback(ctx)
}

// PageRouteOptions configures optional parameters for Page.Route.
type PageRouteOptions struct {
	// Times limits how many times the handler is invoked. Zero (default) means unlimited.
	Times *int
}

// Route registers a network interception handler for requests matching the given glob pattern.
// Multiple calls stack handlers in LIFO order; a handler may call Route.Fallback to defer to
// the next registered handler.
func (p *Page) Route(ctx context.Context, pattern string, handler RouteHandler, opts ...*PageRouteOptions) error {
	p.mu.Lock()
	found := false
	for _, existing := range p.routePatterns {
		if existing == pattern {
			found = true
			break
		}
	}
	if !found {
		p.routePatterns = append(p.routePatterns, pattern)
	}
	allPatterns := make([]any, len(p.routePatterns))
	for i, pat := range p.routePatterns {
		allPatterns[i] = map[string]string{"glob": pat}
	}
	p.mu.Unlock()

	req := protocol.PageSetNetworkInterceptionPatternsRequest{
		Patterns: allPatterns,
	}
	_, err := p.owner.SendMessageRequest(ctx, "setNetworkInterceptionPatterns", req)
	if err != nil {
		return fmt.Errorf("page.route: failed to set network interception patterns: %w", err)
	}

	times := 0
	if len(opts) > 0 && opts[0] != nil && opts[0].Times != nil {
		times = *opts[0].Times
	}

	p.mu.Lock()
	p.pageRouter.add(pattern, handler, times)
	if !p.pageRouteListenerActive {
		// When all page-level handlers fall back, pass the route to the context router
		// (mirrors Java PageImpl: calls browserContext.handleRoute when page routes fall back).
		if p.browserContext != nil {
			bc := p.browserContext
			p.pageRouter.mu.Lock()
			p.pageRouter.onAllFallback = func(route *Route) {
				bc.ctxRouter.dispatch(route)
			}
			p.pageRouter.mu.Unlock()
		}
		p.pageRouteListenerID = p.owner.conn.OnEvent(p.owner.guid, "route", func(params json.RawMessage) {
			var event protocol.PageRouteEvent
			if err := json.Unmarshal(params, &event); err != nil {
				return
			}
			if event.Route.Guid == "" {
				return
			}
			route := &Route{owner: p.owner.child(event.Route.Guid), browserContext: p.browserContext}
			p.pageRouter.dispatch(route)
		})
		p.pageRouteListenerActive = true
	}
	p.mu.Unlock()

	return nil
}

// UnrouteAll removes all registered route handlers and clears network interception patterns.
func (p *Page) UnrouteAll(ctx context.Context, behavior ...string) error {
	return p.Unroute(ctx)
}

// Unroute clears all registered network route handlers on the page.
func (p *Page) Unroute(ctx context.Context) error {
	p.mu.Lock()
	p.pageRouter.clear()
	p.routePatterns = nil
	listenerID := p.pageRouteListenerID
	active := p.pageRouteListenerActive
	p.pageRouteListenerActive = false
	p.mu.Unlock()

	if active {
		p.owner.conn.OffEvent(p.owner.guid, "route", listenerID)
	}

	req := protocol.PageSetNetworkInterceptionPatternsRequest{
		Patterns: []any{},
	}
	_, err := p.owner.SendMessageRequest(ctx, "setNetworkInterceptionPatterns", req)
	if err != nil {
		return fmt.Errorf("page.unroute failed: %w", err)
	}
	return nil
}
