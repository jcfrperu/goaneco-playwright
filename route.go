package playwright

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// RouteHandler is the function called when a request matches a route pattern.
type RouteHandler func(route *Route)

// Route represents an intercepted network request that can be fulfilled,
// continued, or aborted.
type Route struct {
	owner     ChannelOwner
	mu        sync.Mutex
	doneCh    chan bool // non-nil when managed by routeRouter; receives false=handled, true=fallback
	reqOnce   sync.Once
	cachedReq *NetworkRequest
}

// setDoneCh arms the route with a completion channel (called by routeRouter before each handler).
func (r *Route) setDoneCh(ch chan bool) {
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
		r.cachedReq = networkRequestFrom(r.owner, init.Request.Guid, reqRaw)
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
}

// RouteContinueOptions specifies modifications for Route.Continue.
type RouteContinueOptions struct {
	Headers  map[string]string
	Method   *string
	PostData []byte
	URL      *string
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
		req.Status = opts.Status
		req.Body = opts.Body

		if len(opts.BodyBytes) > 0 {
			encoded := base64.StdEncoding.EncodeToString(opts.BodyBytes)
			req.Body = &encoded
			req.IsBase64 = protocol.Bool(true)
		}

		headers := make([]protocol.NameValue, 0)
		for k, v := range opts.Headers {
			headers = append(headers, protocol.NameValue{Name: k, Value: v})
		}
		if opts.ContentType != nil {
			headers = append(headers, protocol.NameValue{Name: "content-type", Value: *opts.ContentType})
		}
		if len(headers) > 0 {
			req.Headers = headers
		}
	}

	_, err := r.owner.SendMessageRequest(ctx, "fulfill", req)
	if err != nil {
		return fmt.Errorf("route.fulfill failed: %w", err)
	}
	r.signalDone(false)
	return nil
}

// Continue passes the request to the network with optional modifications.
func (r *Route) Continue(ctx context.Context, opts *RouteContinueOptions) error {
	req := protocol.RouteContinueRequest{
		IsFallback: false,
		Headers:    []protocol.NameValue{},
	}

	if opts != nil {
		req.Method = opts.Method
		req.Url = opts.URL
		if opts.PostData != nil {
			encoded := base64.StdEncoding.EncodeToString(opts.PostData)
			req.PostData = &encoded
		}

		headers := make([]protocol.NameValue, 0, len(opts.Headers))
		for k, v := range opts.Headers {
			headers = append(headers, protocol.NameValue{Name: k, Value: v})
		}
		if len(headers) > 0 {
			req.Headers = headers
		}
	}

	_, err := r.owner.SendMessageRequest(ctx, "continue", req)
	if err != nil {
		return fmt.Errorf("route.continue failed: %w", err)
	}
	r.signalDone(false)
	return nil
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
	req := protocol.RouteContinueRequest{
		IsFallback: true,
		Headers:    []protocol.NameValue{},
	}
	_, err := r.owner.SendMessageRequest(ctx, "continue", req)
	if err != nil {
		return fmt.Errorf("route.continueAsFallback failed: %w", err)
	}
	return nil
}

// RouteFallbackOptions configures optional overrides when falling back a route.
type RouteFallbackOptions struct {
	// URL overrides the request URL when falling back to the next handler or the network.
	URL *string
}

// Fallback defers handling of this route to the next registered handler. When managed by a
// routeRouter, it signals the router to proceed to the next entry. When called outside a
// router context, it sends "continue" with isFallback=true directly to the server.
func (r *Route) Fallback(ctx context.Context, opts ...*RouteFallbackOptions) error {
	var urlOverride *string
	if len(opts) > 0 && opts[0] != nil {
		urlOverride = opts[0].URL
	}

	if urlOverride != nil {
		// URL override bypasses the router fallback chain and sends a modified continue directly.
		req := protocol.RouteContinueRequest{
			IsFallback: true,
			Headers:    []protocol.NameValue{},
			Url:        urlOverride,
		}
		_, err := r.owner.SendMessageRequest(ctx, "continue", req)
		if err != nil {
			return fmt.Errorf("route.fallback (url override) failed: %w", err)
		}
		r.signalDone(false)
		return nil
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
			route := &Route{owner: p.owner.child(event.Route.Guid)}
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
