package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// frameInitializer contains protocol initializer data for a Frame.
// Corresponds to the `initializer` field in the `__create__` event for Frame.
// Ref: packages/protocol/spec/frame.yml
type frameInitializer struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type frameGotoParams struct {
	URL       string                   `json:"url"`
	Timeout   float64                  `json:"timeout"`
	WaitUntil *protocol.LifecycleEvent `json:"waitUntil,omitempty"`
	Referer   *string                  `json:"referer,omitempty"`
}

type frameSetContentParams struct {
	HTML      string                   `json:"html"`
	Timeout   float64                  `json:"timeout"`
	WaitUntil *protocol.LifecycleEvent `json:"waitUntil,omitempty"`
}

type frameEvaluateExpressionParams struct {
	Expression string                `json:"expression"`
	World      string                `json:"world"`
	Arg        serializedArgumentRaw `json:"arg"`
}

// Page provides methods to interact with a single tab in a Browser,
// or an extension background page in Chromium.
type Page struct {
	owner                   ChannelOwner
	mainFrame               ChannelOwner
	frameInit               frameInitializer // main frame data (URL, name) — updated via events
	Keyboard                *Keyboard        // keyboard input controller for this page
	Mouse                   *Mouse           // mouse input controller for this page
	frames                  []*Frame         // child frames tracked via frameAttached/frameDetached; guarded by mu
	framesByGUID            map[string]*Frame
	initializer             json.RawMessage // full page initializer (for future usage)
	openerGUID              string          // GUID of the page that opened this one via window.open; empty for top-level pages
	browserContext          *BrowserContext // context to which this page belongs
	mu                      sync.RWMutex    // protects frameInit, closed, errorWriter, workers, pageRouter, video, and collectedPageErrors
	closed                  bool
	navListenID             connection.ListenerID // navigation listener ID (for cancellation)
	pageLoadListenID        connection.ListenerID // loadstate listener from subscribeToPageLoad
	pageErrorListenID       connection.ListenerID // pageError listener from subscribeToPageErrors
	pageRouter              routeRouter           // manages Page-level route handler chain
	pageRouteListenerID     connection.ListenerID // single event listener for "route" events
	pageRouteListenerActive bool                  // true once the listener is registered
	routePatterns           []string              // accumulated glob patterns for network interception
	errorWriter             io.Writer             // destination for diagnostic messages (nil -> os.Stderr)
	workers                 []*Worker             // currently active workers; guarded by mu
	workerListenersByID     map[int]func(*Worker)
	workerNextID            int
	workerListenerActive    bool
	workerListenerID        connection.ListenerID
	video                   *Video        // set when the page's video recording artifact is ready
	collectedPageErrors     []string      // uncaught exceptions accumulated via context pageerror event; guarded by mu
	viewportSize            *ViewportSize // nil when context was created with NoDefaultViewport
}

// SetErrorWriter sets the destination for internal diagnostic messages (e.g., route handler panics).
// If w is nil, os.Stderr is used. Safe to call concurrently.
func (p *Page) SetErrorWriter(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.errorWriter = w
}

// getErrorWriter returns the configured error writer or os.Stderr under read lock.
func (p *Page) getErrorWriter() io.Writer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.errorWriter != nil {
		return p.errorWriter
	}
	return os.Stderr
}

// Context returns the BrowserContext to which this page belongs.
func (p *Page) Context() *BrowserContext {
	return p.browserContext
}

// URL returns the current URL of the main frame.
// Initialized from the Frame initializer and updated via `navigated` events.
func (p *Page) URL() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.frameInit.URL
}

// IsClosed returns true if the page has been closed.
func (p *Page) IsClosed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.closed
}

// markClosed sets the page's closed flag without triggering IPC or parent-context notification.
// Called by BrowserContext.Close() to propagate closed state to all owned pages.
func (p *Page) markClosed() {
	p.mu.Lock()
	p.closed = true
	p.mu.Unlock()
}

// PageGotoOptions configures the behaviour of Page.Goto.
type PageGotoOptions struct {
	// WaitUntil is the load state to wait for: "load", "domcontentloaded", "networkidle", or "commit".
	WaitUntil *string
	// Referer overrides the HTTP Referer header.
	Referer *string
	// Timeout overrides the default navigation timeout (milliseconds).
	Timeout *float64
}

// Goto navigates the page to the specified URL.
func (p *Page) Goto(ctx context.Context, url string, opts ...*PageGotoOptions) error {
	params := frameGotoParams{
		URL:     url,
		Timeout: defaultActionTimeoutMs,
	}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.Timeout != nil {
			params.Timeout = *o.Timeout
		}
		if o.WaitUntil != nil {
			le := protocol.LifecycleEvent(*o.WaitUntil)
			params.WaitUntil = &le
		}
		params.Referer = o.Referer
	}

	_, err := p.mainFrame.SendMessageRequest(ctx, "goto", params)
	if err != nil {
		return fmt.Errorf("page.goto failed: %w", err)
	}
	return nil
}

// Title returns the document title of the current page.
func (p *Page) Title(ctx context.Context) (string, error) {
	req := protocol.FrameTitleRequest{}
	result, err := p.mainFrame.SendMessageRequest(ctx, "title", req)
	if err != nil {
		return "", fmt.Errorf("page.title failed: %w", err)
	}

	var resp protocol.FrameTitleResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse page.title response: %w", err)
	}
	return resp.Value, nil
}

// Content returns the full HTML contents of the page.
func (p *Page) Content(ctx context.Context) (string, error) {
	req := protocol.FrameContentRequest{}
	result, err := p.mainFrame.SendMessageRequest(ctx, "content", req)
	if err != nil {
		return "", fmt.Errorf("page.content failed: %w", err)
	}

	var resp protocol.FrameContentResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse page.content response: %w", err)
	}
	return resp.Value, nil
}

// SetContent sets the HTML markup of the page.
func (p *Page) SetContent(ctx context.Context, html string) error {
	params := frameSetContentParams{
		HTML:    html,
		Timeout: defaultActionTimeoutMs,
	}
	_, err := p.mainFrame.SendMessageRequest(ctx, "setContent", params)
	if err != nil {
		return fmt.Errorf("page.setContent failed: %w", err)
	}
	return nil
}

// Evaluate executes a JavaScript expression or function in the context of the main frame and returns the deserialized result.
func (p *Page) Evaluate(ctx context.Context, expression string, arg ...any) (any, error) {
	var inputArg any
	if len(arg) > 0 {
		inputArg = arg[0]
	}

	params := frameEvaluateExpressionParams{
		Expression: expression,
		World:      "main",
		Arg:        serializeArgument(inputArg),
	}

	result, err := p.mainFrame.SendMessageRequest(ctx, "evaluateExpression", params)
	if err != nil {
		return nil, fmt.Errorf("page.evaluate failed: %w", err)
	}

	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse evaluate response: %w", err)
	}

	return deserializeValue(resp.Value)
}

// DispatchEvent dispatches an event of the given type on the first element matching selector.
// eventInit can optionally provide event initialization data (passed as-is to the event constructor).
func (p *Page) DispatchEvent(ctx context.Context, selector, eventType string, eventInit ...any) error {
	var initArg any
	if len(eventInit) > 0 {
		initArg = eventInit[0]
	}
	_, err := p.mainFrame.SendMessageRequest(ctx, "dispatchEvent", map[string]any{
		"selector":  selector,
		"type":      eventType,
		"eventInit": serializeArgument(initArg),
		"strict":    false,
		"timeout":   defaultActionTimeoutMs,
	})
	if err != nil {
		return fmt.Errorf("page.dispatchEvent failed: %w", err)
	}
	return nil
}

// EvalOnSelector evaluates an expression on the first element matching selector and returns the result.
// Equivalent to document.querySelector(selector) followed by element.evaluate(expression, arg).
func (p *Page) EvalOnSelector(ctx context.Context, selector, expression string, arg ...any) (any, error) {
	var inputArg any
	if len(arg) > 0 {
		inputArg = arg[0]
	}
	result, err := p.mainFrame.SendMessageRequest(ctx, "evalOnSelector", map[string]any{
		"selector":   selector,
		"expression": expression,
		"arg":        serializeArgument(inputArg),
		"strict":     false,
	})
	if err != nil {
		return nil, fmt.Errorf("page.evalOnSelector failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse evalOnSelector response: %w", err)
	}
	return deserializeValue(resp.Value)
}

// EvalOnSelectorAll evaluates an expression on all elements matching selector and returns the result.
// Equivalent to document.querySelectorAll(selector) followed by NodeList.evaluate(expression, arg).
func (p *Page) EvalOnSelectorAll(ctx context.Context, selector, expression string, arg ...any) (any, error) {
	var inputArg any
	if len(arg) > 0 {
		inputArg = arg[0]
	}
	result, err := p.mainFrame.SendMessageRequest(ctx, "evalOnSelectorAll", map[string]any{
		"selector":   selector,
		"expression": expression,
		"arg":        serializeArgument(inputArg),
	})
	if err != nil {
		return nil, fmt.Errorf("page.evalOnSelectorAll failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse evalOnSelectorAll response: %w", err)
	}
	return deserializeValue(resp.Value)
}

// EmulateMediaOptions contains options for Page.EmulateMedia.
type EmulateMediaOptions struct {
	// Media specifies the CSS media type: "screen", "print", or "null" to reset.
	Media *string
	// ColorScheme specifies the color scheme: "dark", "light", "no-preference", or "null" to reset.
	ColorScheme *string
	// ReducedMotion specifies the reduced motion preference: "reduce", "no-preference", or "null" to reset.
	ReducedMotion *string
	// ForcedColors specifies the forced colors mode: "active", "none", or "null" to reset.
	ForcedColors *string
	// Contrast specifies the contrast preference: "more", "no-preference", or "null" to reset.
	Contrast *string
}

// EmulateMedia overrides the CSS media type and/or color scheme for the page.
func (p *Page) EmulateMedia(ctx context.Context, opts ...*EmulateMediaOptions) error {
	req := map[string]any{}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.Media != nil {
			req["media"] = *o.Media
		}
		if o.ColorScheme != nil {
			req["colorScheme"] = *o.ColorScheme
		}
		if o.ReducedMotion != nil {
			req["reducedMotion"] = *o.ReducedMotion
		}
		if o.ForcedColors != nil {
			req["forcedColors"] = *o.ForcedColors
		}
		if o.Contrast != nil {
			req["contrast"] = *o.Contrast
		}
	}
	_, err := p.owner.SendMessageRequest(ctx, "emulateMedia", req)
	if err != nil {
		return fmt.Errorf("page.emulateMedia failed: %w", err)
	}
	return nil
}

// Close closes the page and marks its internal state as closed.
// It is idempotent: later calls perform no action and return nil.
//
// Fail-fast design: local state p.closed is updated to true immediately to prevent
// redundant concurrent calls. If the underlying IPC call fails due to driver disconnection,
// the Playwright server automatically reaps remote resources on process exit.
func (p *Page) Close(ctx context.Context) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	listenID := p.navListenID
	loadListenID := p.pageLoadListenID
	errListenID := p.pageErrorListenID
	bCtx := p.browserContext
	p.mu.Unlock()

	// Detach from parent BrowserContext if present
	if bCtx != nil {
		bCtx.removePage(p)
	}

	if p.owner.conn != nil {
		if listenID != 0 {
			p.owner.conn.OffEvent(p.mainFrame.guid, "navigated", listenID)
		}
		if loadListenID != 0 {
			p.owner.conn.OffEvent(p.mainFrame.guid, "loadstate", loadListenID)
		}
		if errListenID != 0 && bCtx != nil {
			p.owner.conn.OffEvent(bCtx.owner.guid, "pageError", errListenID)
		}
	}

	req := protocol.PageCloseRequest{}
	_, err := p.owner.SendMessageRequest(ctx, "close", req)
	if err != nil {
		return fmt.Errorf("failed to close page: %w", err)
	}
	return nil
}

type pageNavParams struct {
	Timeout float64 `json:"timeout"`
}

// Opener returns the Page that opened this page via window.open(), or nil if this page
// was not opened as a popup. Returns (nil, nil) for top-level pages.
func (p *Page) Opener(_ context.Context) (*Page, error) {
	if p.openerGUID == "" || p.browserContext == nil {
		return nil, nil
	}
	return p.browserContext.pageByGUID(p.openerGUID), nil
}

// ExposeBinding exposes a Go function as a global JavaScript function on this page only.
// When JavaScript calls window[name](...args), handler is invoked with the deserialized args.
// Unlike BrowserContext.ExposeBinding, this is scoped to the single page.
func (p *Page) ExposeBinding(ctx context.Context, name string, handler BindingHandler) error {
	if p.browserContext == nil {
		return fmt.Errorf("page.ExposeBinding: page has no browser context")
	}
	_, err := p.owner.SendMessageRequest(ctx, "exposeBinding", map[string]string{"name": name})
	if err != nil {
		return fmt.Errorf("page.exposeBinding(%q) failed: %w", name, err)
	}
	c := p.browserContext
	c.mu.Lock()
	if c.bindingListenerIDs == nil {
		c.bindingListenerIDs = make(map[string]connection.ListenerID)
	}
	// Use a page-scoped key so page and context bindings with the same name don't collide.
	key := p.owner.guid + "\x00" + name
	if prev, ok := c.bindingListenerIDs[key]; ok {
		p.owner.conn.OffEvent(p.owner.guid, "bindingCall", prev)
	}
	id := p.owner.conn.OnEvent(p.owner.guid, "bindingCall", func(params json.RawMessage) {
		var event struct {
			Binding struct {
				Guid string `json:"guid"`
			} `json:"binding"`
		}
		if err := json.Unmarshal(params, &event); err != nil || event.Binding.Guid == "" {
			return
		}
		raw := p.owner.Initializer(event.Binding.Guid)
		var init struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(raw, &init); err != nil || init.Name != name {
			return
		}
		go c.dispatchBindingCall(event.Binding.Guid, handler)
	})
	c.bindingListenerIDs[key] = id
	c.mu.Unlock()
	return nil
}

// SetDefaultTimeout sets the default timeout for all actions on this page (in milliseconds).
func (p *Page) SetDefaultTimeout(ctx context.Context, timeout float64) error {
	_, err := p.owner.SendMessageRequest(ctx, "setDefaultTimeout", map[string]float64{"timeout": timeout})
	if err != nil {
		return fmt.Errorf("page.setDefaultTimeout failed: %w", err)
	}
	return nil
}

// SetDefaultNavigationTimeout sets the default timeout for navigation operations (in milliseconds).
func (p *Page) SetDefaultNavigationTimeout(ctx context.Context, timeout float64) error {
	_, err := p.owner.SendMessageRequest(ctx, "setDefaultNavigationTimeout", map[string]float64{"timeout": timeout})
	if err != nil {
		return fmt.Errorf("page.setDefaultNavigationTimeout failed: %w", err)
	}
	return nil
}

// Reload reloads the current page and waits for the load event.
func (p *Page) Reload(ctx context.Context) error {
	_, err := p.owner.SendMessageRequest(ctx, "reload", pageNavParams{Timeout: defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("page.reload failed: %w", err)
	}
	return nil
}

// GoBack navigates to the previous page in the session history.
func (p *Page) GoBack(ctx context.Context) error {
	_, err := p.owner.SendMessageRequest(ctx, "goBack", pageNavParams{Timeout: defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("page.goBack failed: %w", err)
	}
	return nil
}

// GoForward navigates to the next page in the session history.
func (p *Page) GoForward(ctx context.Context) error {
	_, err := p.owner.SendMessageRequest(ctx, "goForward", pageNavParams{Timeout: defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("page.goForward failed: %w", err)
	}
	return nil
}

// Locator returns a new Locator for the given CSS/attribute selector.
func (p *Page) Locator(selector string) *Locator {
	return &Locator{frame: p.mainFrame, selector: selector}
}

// MainFrame returns the main (top-level) frame of the page.
func (p *Page) MainFrame() *Frame {
	return &Frame{
		owner: p.mainFrame,
		page:  p,
		name:  p.frameInit.Name,
		url:   p.URL(),
	}
}

// Frames returns all frames currently attached to the page, including the main frame as the first element.
func (p *Page) Frames() []*Frame {
	p.mu.RLock()
	// Read frameInit fields directly to avoid re-acquiring the lock via p.URL() or p.MainFrame().
	main := &Frame{owner: p.mainFrame, name: p.frameInit.Name, url: p.frameInit.URL}
	result := make([]*Frame, 0, 1+len(p.frames))
	result = append(result, main)
	result = append(result, p.frames...)
	p.mu.RUnlock()
	return result
}

// Frame returns the first child frame with the given name attribute, or nil if not found.
func (p *Page) Frame(name string) *Frame {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, f := range p.frames {
		if f.name == name {
			return f
		}
	}
	return nil
}

// FrameByURL returns the first child frame whose URL matches the given pattern, or nil if not found.
// A trailing * means prefix matching; otherwise exact match.
func (p *Page) FrameByURL(urlPattern string) *Frame {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, f := range p.frames {
		if urlMatchesPattern(f.URL(), urlPattern) {
			return f
		}
	}
	return nil
}

// DragTo performs a drag-and-drop from the element matching sourceSelector to targetSelector.
func (p *Page) DragTo(ctx context.Context, sourceSelector, targetSelector string) error {
	_, err := p.mainFrame.SendMessageRequest(ctx, "dragAndDrop", map[string]any{
		"source":  sourceSelector,
		"target":  targetSelector,
		"timeout": defaultActionTimeoutMs,
		"strict":  false,
	})
	if err != nil {
		return fmt.Errorf("page.dragTo failed: %w", err)
	}
	return nil
}

// GetByRole returns a Locator for elements matching the given ARIA role.
func (p *Page) GetByRole(role AriaRole, opts ...*GetByRoleOptions) *Locator {
	var o *GetByRoleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	return locatorFromRole(p.mainFrame, role, o)
}

// GetByText returns a Locator for elements containing the given text.
func (p *Page) GetByText(text string, opts ...*GetByTextOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return locatorFromText(p.mainFrame, text, exact)
}

// GetByLabel returns a Locator for form elements associated with the given label text.
func (p *Page) GetByLabel(text string, opts ...*GetByLabelOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return locatorFromLabel(p.mainFrame, text, exact)
}

// GetByPlaceholder returns a Locator for input elements with the given placeholder text.
func (p *Page) GetByPlaceholder(text string, opts ...*GetByPlaceholderOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return locatorFromAttr(p.mainFrame, "placeholder", text, exact)
}

// GetByAltText returns a Locator for elements with the given alt text.
func (p *Page) GetByAltText(text string, opts ...*GetByAltTextOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return locatorFromAttr(p.mainFrame, "alt", text, exact)
}

// GetByTitle returns a Locator for elements with the given title attribute.
func (p *Page) GetByTitle(text string, opts ...*GetByTitleOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return locatorFromAttr(p.mainFrame, "title", text, exact)
}

// GetByTestId returns a Locator for elements with the configured testid attribute.
// The attribute name defaults to "data-testid" and can be changed with BrowserContext.SetTestIdAttribute.
func (p *Page) GetByTestId(testId string) *Locator {
	attr := "data-testid"
	if p.browserContext != nil {
		attr = p.browserContext.TestIdAttributeName()
	}
	escaped := escapeForAttributeSelector(testId)
	return &Locator{frame: p.mainFrame, selector: fmt.Sprintf(`internal:testid=[%s=%s]`, attr, escaped)}
}

// FrameLocator returns a FrameLocator scoped to the first iframe matching selector.
// All locator operations on the returned FrameLocator target elements inside the iframe.
func (p *Page) FrameLocator(selector string) *FrameLocator {
	return &FrameLocator{
		frame:          p.mainFrame,
		selector:       selector + " >> internal:control=enter-frame",
		browserContext: p.browserContext,
	}
}

// OnRequest registers a handler that fires for every network request made by this page.
// The handler receives a NetworkRequest with URL, method, headers, and resource type.
// The returned function cancels the listener.
func (p *Page) OnRequest(handler func(*NetworkRequest)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	p.browserContext.retainSubscription(subCtx, "request")
	subCancel()
	process := func(params json.RawMessage) {
		var event struct {
			Page    *protocol.Page   `json:"page,omitempty"`
			Request protocol.Request `json:"request"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Page != nil && event.Page.Guid != p.owner.guid {
			return
		}
		raw := p.owner.Initializer(event.Request.Guid)
		if req := networkRequestFrom(p.owner, event.Request.Guid, raw); req != nil {
			go handler(req)
		}
	}
	id := p.owner.conn.OnEvent(p.browserContext.owner.guid, "request", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			p.owner.conn.OffEvent(p.browserContext.owner.guid, "request", id)
			p.browserContext.releaseSubscription("request")
		})
	}
}

// OnResponse registers a handler that fires for every HTTP response received by this page.
// The handler receives a NetworkResponse with URL, status, and headers.
// The returned function cancels the listener.
func (p *Page) OnResponse(handler func(*NetworkResponse)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	p.browserContext.retainSubscription(subCtx, "response")
	subCancel()
	process := func(params json.RawMessage) {
		var event struct {
			Page     *protocol.Page    `json:"page,omitempty"`
			Response protocol.Response `json:"response"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Page != nil && event.Page.Guid != p.owner.guid {
			return
		}
		raw := p.owner.Initializer(event.Response.Guid)
		if resp := networkResponseFrom(p.owner, event.Response.Guid, raw); resp != nil {
			go handler(resp)
		}
	}
	id := p.owner.conn.OnEvent(p.browserContext.owner.guid, "response", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			p.owner.conn.OffEvent(p.browserContext.owner.guid, "response", id)
			p.browserContext.releaseSubscription("response")
		})
	}
}

// OnRequestFinished registers a handler for network requests that have completed successfully.
// The handler receives a NetworkRequest. The returned function cancels the listener.
func (p *Page) OnRequestFinished(handler func(*NetworkRequest)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	p.browserContext.retainSubscription(subCtx, "requestFinished")
	subCancel()
	process := func(params json.RawMessage) {
		var event protocol.BrowserContextRequestFinishedEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Page != nil && event.Page.Guid != p.owner.guid {
			return
		}
		raw := p.owner.Initializer(event.Request.Guid)
		if req := networkRequestFrom(p.owner, event.Request.Guid, raw); req != nil {
			go handler(req)
		}
	}
	id := p.owner.conn.OnEvent(p.browserContext.owner.guid, "requestFinished", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			p.owner.conn.OffEvent(p.browserContext.owner.guid, "requestFinished", id)
			p.browserContext.releaseSubscription("requestFinished")
		})
	}
}

// OnRequestFailed registers a handler for network requests that have failed.
// The handler receives a NetworkRequest. The returned function cancels the listener.
func (p *Page) OnRequestFailed(handler func(*NetworkRequest)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	p.browserContext.retainSubscription(subCtx, "requestFailed")
	subCancel()
	process := func(params json.RawMessage) {
		var event protocol.BrowserContextRequestFailedEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Page != nil && event.Page.Guid != p.owner.guid {
			return
		}
		raw := p.owner.Initializer(event.Request.Guid)
		if req := networkRequestFrom(p.owner, event.Request.Guid, raw); req != nil {
			go handler(req)
		}
	}
	id := p.owner.conn.OnEvent(p.browserContext.owner.guid, "requestFailed", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			p.owner.conn.OffEvent(p.browserContext.owner.guid, "requestFailed", id)
			p.browserContext.releaseSubscription("requestFailed")
		})
	}
}

// WaitForURL waits until the page URL matches the given pattern.
// The pattern can be an exact URL or end with * for prefix matching.
// Uses the provided context for cancellation; default timeout is 30s.
func (p *Page) WaitForURL(ctx context.Context, pattern string, timeout ...time.Duration) error {
	to := 30 * time.Second
	if len(timeout) > 0 {
		to = timeout[0]
	}

	matched := make(chan struct{}, 1)
	// Subscribe before checking current URL to avoid missing a navigation that
	// fires between the check and the subscribe.
	id := p.owner.conn.OnEvent(p.mainFrame.guid, "navigated", func(params json.RawMessage) {
		var nav protocol.FrameNavigatedEvent
		if err := json.Unmarshal(params, &nav); err != nil {
			return
		}
		if urlMatchesPattern(nav.Url, pattern) {
			select {
			case matched <- struct{}{}:
			default:
			}
		}
	})
	defer p.owner.conn.OffEvent(p.mainFrame.guid, "navigated", id)

	if urlMatchesPattern(p.URL(), pattern) {
		return nil
	}

	timer := time.NewTimer(to)
	defer timer.Stop()

	select {
	case <-matched:
		return nil
	case <-timer.C:
		return fmt.Errorf("page.waitForURL: timeout after %v waiting for %q", to, pattern)
	case <-ctx.Done():
		return ctx.Err()
	}
}

// urlMatchesPattern returns true if url matches the given pattern.
// A trailing * means prefix match; otherwise exact match is used.
func urlMatchesPattern(url, pattern string) bool {
	if prefix, ok := strings.CutSuffix(pattern, "*"); ok {
		return strings.HasPrefix(url, prefix)
	}
	return url == pattern
}

// WaitForLoadState waits until the page reaches the given load state.
// Valid states: "load" (default), "domcontentloaded", "networkidle".
// Delegates to the server-side Frame.waitForLoadState which handles already-reached states correctly,
// including "networkidle" which has no equivalent DOM readyState to check client-side.
func (p *Page) WaitForLoadState(ctx context.Context, state ...string) error {
	want := "load"
	if len(state) > 0 && state[0] != "" {
		want = state[0]
	}
	timeout := defaultActionTimeoutMs
	if d, hasDeadline := ctx.Deadline(); hasDeadline {
		if remaining := float64(time.Until(d).Milliseconds()); remaining > 0 {
			timeout = remaining
		}
	}
	_, err := p.mainFrame.SendMessageRequest(ctx, "waitForLoadState", map[string]any{
		"state":   want,
		"timeout": timeout,
	})
	if err != nil {
		return fmt.Errorf("page.waitForLoadState(%q) failed: %w", want, err)
	}
	return nil
}

// OnLoad registers a handler that fires when the page's "load" event fires.
// The returned function cancels the listener.
func (p *Page) OnLoad(handler func(*Page)) func() {
	process := func(params json.RawMessage) {
		var ev protocol.FrameLoadstateEvent
		if err := json.Unmarshal(params, &ev); err != nil {
			return
		}
		if ev.Add != nil && string(*ev.Add) == "load" {
			go handler(p)
		}
	}
	id := p.owner.conn.OnEvent(p.mainFrame.guid, "loadstate", process)
	return func() { p.owner.conn.OffEvent(p.mainFrame.guid, "loadstate", id) }
}

// OnDOMContentLoaded registers a handler that fires when the page's "DOMContentLoaded" event fires.
// The returned function cancels the listener.
func (p *Page) OnDOMContentLoaded(handler func(*Page)) func() {
	process := func(params json.RawMessage) {
		var ev protocol.FrameLoadstateEvent
		if err := json.Unmarshal(params, &ev); err != nil {
			return
		}
		if ev.Add != nil && string(*ev.Add) == "domcontentloaded" {
			go handler(p)
		}
	}
	id := p.owner.conn.OnEvent(p.mainFrame.guid, "loadstate", process)
	return func() { p.owner.conn.OffEvent(p.mainFrame.guid, "loadstate", id) }
}

// WaitForPopup registers a popup listener, calls trigger, then waits for a popup to open.
// trigger is called after the listener is registered to avoid race conditions.
func (p *Page) WaitForPopup(ctx context.Context, trigger func() error) (*Page, error) {
	ch := make(chan *Page, 1)
	unsub := p.OnPopup(func(pg *Page) {
		select {
		case ch <- pg:
		default:
		}
	})
	defer unsub()

	if trigger != nil {
		if err := trigger(); err != nil {
			return nil, fmt.Errorf("page.WaitForPopup: trigger failed: %w", err)
		}
	}

	select {
	case pg := <-ch:
		return pg, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("page.WaitForPopup: %w", ctx.Err())
	}
}

// WaitForClose waits for this page to be closed. trigger is called after the listener
// is registered. Returns when the page fires its "close" event.
func (p *Page) WaitForClose(ctx context.Context, trigger func() error) error {
	closed := make(chan struct{}, 1)
	p.owner.conn.OnEventOnce(p.owner.guid, "close", func(_ json.RawMessage) {
		select {
		case closed <- struct{}{}:
		default:
		}
	})

	if trigger != nil {
		if err := trigger(); err != nil {
			return fmt.Errorf("page.WaitForClose: trigger failed: %w", err)
		}
	}

	select {
	case <-closed:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("page.WaitForClose: %w", ctx.Err())
	}
}

// Pause pauses script execution. Playwright will stop executing the script and wait
// for the user to either press 'Resume' button in the page overlay or call playwright.resume().
func (p *Page) Pause(ctx context.Context) error {
	if p.browserContext == nil {
		return nil
	}
	_, err := p.browserContext.owner.SendMessageRequest(ctx, "pause", protocol.BrowserContextPauseRequest{})
	if err != nil {
		return fmt.Errorf("page.pause failed: %w", err)
	}
	return nil
}

// RequestGC triggers a garbage collection in the page's JavaScript engine.
// Useful in tests to force GC before checking memory-sensitive behaviour.
func (p *Page) RequestGC(ctx context.Context) error {
	_, err := p.owner.SendMessageRequest(ctx, "requestGC", struct{}{})
	if err != nil {
		return fmt.Errorf("page.requestGC failed: %w", err)
	}
	return nil
}

// HideHighlight removes all active element highlights on the page.
func (p *Page) HideHighlight(ctx context.Context) error {
	_, err := p.owner.SendMessageRequest(ctx, "hideHighlight", struct{}{})
	if err != nil {
		return fmt.Errorf("page.hideHighlight failed: %w", err)
	}
	return nil
}

// OnPopup registers a handler that fires whenever this page opens a popup via window.open()
// or a link with target="_blank". The returned function cancels the listener.
func (p *Page) OnPopup(fn func(*Page)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	return p.browserContext.OnPage(func(pg *Page) {
		if pg.openerGUID == p.owner.guid {
			fn(pg)
		}
	})
}

// OnWebSocket registers a handler that fires whenever the page opens a WebSocket connection.
// The handler receives the new WebSocket and is called in a goroutine.
// The returned function cancels the listener.
func (p *Page) OnWebSocket(handler func(*WebSocket)) func() {
	process := func(params json.RawMessage) {
		var event protocol.PageWebSocketEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.WebSocket.Guid == "" {
			return
		}
		raw := p.owner.Initializer(event.WebSocket.Guid)
		ws := webSocketFromGUID(p.owner, event.WebSocket.Guid, raw)
		go handler(ws)
	}
	id := p.owner.conn.OnEvent(p.owner.guid, "webSocket", process)
	return func() { p.owner.conn.OffEvent(p.owner.guid, "webSocket", id) }
}

// WaitForWebSocket waits for a WebSocket connection to be created during the execution
// of the trigger function. trigger is called after the listener is registered to avoid
// race conditions. Returns the first WebSocket created, or an error if ctx is canceled.
func (p *Page) WaitForWebSocket(ctx context.Context, trigger func() error) (*WebSocket, error) {
	ch := make(chan *WebSocket, 1)
	unsub := p.OnWebSocket(func(ws *WebSocket) {
		select {
		case ch <- ws:
		default:
		}
	})
	defer unsub()

	if trigger != nil {
		if err := trigger(); err != nil {
			return nil, fmt.Errorf("page.WaitForWebSocket: trigger failed: %w", err)
		}
	}

	select {
	case ws := <-ch:
		return ws, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("page.WaitForWebSocket: %w", ctx.Err())
	}
}

// WaitForEvent waits for a page-level event to fire and returns its payload.
// Supported events: "popup" (*Page), "request" (*NetworkRequest),
// "response" (*NetworkResponse), "download" (*Download), "websocket" (*WebSocket).
func (p *Page) WaitForEvent(ctx context.Context, event string) (any, error) {
	ch := make(chan any, 1)
	send := func(v any) {
		select {
		case ch <- v:
		default:
		}
	}

	var unsub func()
	switch event {
	case "popup":
		unsub = p.OnPopup(func(pg *Page) { send(pg) })
	case "request":
		unsub = p.OnRequest(func(r *NetworkRequest) { send(r) })
	case "response":
		unsub = p.OnResponse(func(r *NetworkResponse) { send(r) })
	case "download":
		unsub = p.OnDownload(func(d *Download) { send(d) })
	case "websocket":
		unsub = p.OnWebSocket(func(ws *WebSocket) { send(ws) })
	default:
		return nil, fmt.Errorf("page.WaitForEvent: unsupported event %q", event)
	}
	defer unsub()

	select {
	case val := <-ch:
		return val, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("page.WaitForEvent(%q): %w", event, ctx.Err())
	}
}

// QuerySelector returns the first element matching the CSS selector, or nil if not found.
func (p *Page) QuerySelector(ctx context.Context, selector string) (*ElementHandle, error) {
	result, err := p.mainFrame.SendMessageRequest(ctx, "querySelector", map[string]string{"selector": selector})
	if err != nil {
		return nil, fmt.Errorf("page.querySelector failed: %w", err)
	}
	var resp struct {
		Element *struct {
			Guid string `json:"guid"`
		} `json:"element,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse querySelector response: %w", err)
	}
	if resp.Element == nil {
		return nil, nil
	}
	return elementHandleFromGUID(p.mainFrame, resp.Element.Guid), nil
}

// QuerySelectorAll returns all elements matching the CSS selector.
func (p *Page) QuerySelectorAll(ctx context.Context, selector string) ([]*ElementHandle, error) {
	result, err := p.mainFrame.SendMessageRequest(ctx, "querySelectorAll", map[string]string{"selector": selector})
	if err != nil {
		return nil, fmt.Errorf("page.querySelectorAll failed: %w", err)
	}
	var resp struct {
		Elements []struct {
			Guid string `json:"guid"`
		} `json:"elements"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse querySelectorAll response: %w", err)
	}
	handles := make([]*ElementHandle, len(resp.Elements))
	for i, el := range resp.Elements {
		handles[i] = elementHandleFromGUID(p.mainFrame, el.Guid)
	}
	return handles, nil
}

// AddInitScript injects a JavaScript script that runs before any page script on every navigation.
// Useful for mocking globals, stubbing APIs, or seeding state.
func (p *Page) AddInitScript(ctx context.Context, script string) error {
	_, err := p.owner.SendMessageRequest(ctx, "addInitScript", protocol.PageAddInitScriptRequest{Source: script})
	if err != nil {
		return fmt.Errorf("page.addInitScript failed: %w", err)
	}
	return nil
}

// PagePdfOptions specifies optional parameters for Page.Pdf.
type PagePdfOptions struct {
	// DisplayHeaderFooter displays header and footer. Defaults to false.
	DisplayHeaderFooter *bool
	// FooterTemplate is the HTML template for the print footer.
	FooterTemplate *string
	// Format specifies the paper format (e.g. "A4", "Letter"). Takes priority over Width/Height.
	Format *string
	// HeaderTemplate is the HTML template for the print header.
	HeaderTemplate *string
	// Height specifies the paper height, accepts values labeled with units.
	Height *string
	// Landscape renders the page in landscape orientation. Defaults to false.
	Landscape *bool
	// Outline embeds the document outline into the PDF. Defaults to false.
	Outline *bool
	// PageRanges is the paper ranges to print, e.g., "1-5, 8, 11-13". Defaults to all pages.
	PageRanges *string
	// PreferCSSPageSize gives priority to CSS page size declared in the page. Defaults to false.
	PreferCSSPageSize *bool
	// PrintBackground prints background graphics. Defaults to false.
	PrintBackground *bool
	// Scale scales the rendering of the web page. Amount must be between 0.1 and 2. Defaults to 1.
	Scale *float64
	// Tagged generates tagged (accessible) PDF. Defaults to false.
	Tagged *bool
	// Width specifies the paper width, accepts values labeled with units.
	Width *string
}

// PDF generates a PDF of the page and returns the raw bytes.
// Only works in headless Chromium.
func (p *Page) PDF(ctx context.Context, opts ...*PagePdfOptions) ([]byte, error) {
	req := protocol.PagePdfRequest{}
	if len(opts) > 0 && opts[0] != nil {
		opt := opts[0]
		req.DisplayHeaderFooter = opt.DisplayHeaderFooter
		req.FooterTemplate = opt.FooterTemplate
		req.Format = opt.Format
		req.HeaderTemplate = opt.HeaderTemplate
		req.Height = opt.Height
		req.Landscape = opt.Landscape
		req.Outline = opt.Outline
		req.PageRanges = opt.PageRanges
		req.PreferCSSPageSize = opt.PreferCSSPageSize
		req.PrintBackground = opt.PrintBackground
		req.Scale = opt.Scale
		req.Tagged = opt.Tagged
		req.Width = opt.Width
	}
	result, err := p.owner.SendMessageRequest(ctx, "pdf", req)
	if err != nil {
		return nil, fmt.Errorf("page.Pdf: %w", err)
	}
	var resp protocol.PagePdfResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("page.Pdf: parse response: %w", err)
	}
	return resp.Pdf, nil
}

// BringToFront brings the page to the front (activates the tab).
func (p *Page) BringToFront(ctx context.Context) error {
	_, err := p.owner.SendMessageRequest(ctx, "bringToFront", protocol.PageBringToFrontRequest{})
	if err != nil {
		return fmt.Errorf("page.bringToFront failed: %w", err)
	}
	return nil
}

// SetExtraHTTPHeaders sets extra HTTP headers sent on every request from this page.
// Calling with an empty map clears previously set headers.
func (p *Page) SetExtraHTTPHeaders(ctx context.Context, headers map[string]string) error {
	pairs := make([]protocol.NameValue, 0, len(headers))
	for k, v := range headers {
		pairs = append(pairs, protocol.NameValue{Name: k, Value: v})
	}
	_, err := p.owner.SendMessageRequest(ctx, "setExtraHTTPHeaders", protocol.PageSetExtraHTTPHeadersRequest{Headers: pairs})
	if err != nil {
		return fmt.Errorf("page.setExtraHTTPHeaders failed: %w", err)
	}
	return nil
}

// SetViewportSize changes the page viewport to the given dimensions.
func (p *Page) SetViewportSize(ctx context.Context, width, height int) error {
	_, err := p.owner.SendMessageRequest(ctx, "setViewportSize", map[string]any{
		"viewportSize": map[string]int{"width": width, "height": height},
	})
	if err != nil {
		return fmt.Errorf("page.setViewportSize failed: %w", err)
	}
	p.mu.Lock()
	p.viewportSize = &ViewportSize{Width: width, Height: height}
	p.mu.Unlock()
	return nil
}

// ViewportSize returns the current viewport dimensions.
func (p *Page) ViewportSize(_ context.Context) (*ViewportSize, error) {
	p.mu.RLock()
	vs := p.viewportSize
	p.mu.RUnlock()
	if vs == nil {
		return nil, nil
	}
	return &ViewportSize{Width: vs.Width, Height: vs.Height}, nil
}

// WaitForFunctionOptions configures WaitForFunction behavior.
type WaitForFunctionOptions struct {
	PollingInterval float64 // polling interval in ms (0 = requestAnimationFrame)
	Timeout         float64 // timeout in ms (0 = server default ~30s; negative = disable)
}

// WaitForFunction polls the given JavaScript expression or function until it returns a truthy value,
// then returns the result wrapped in a JSHandle. If the expression throws, the error propagates.
func (p *Page) WaitForFunction(ctx context.Context, expression string, arg any, opts ...*WaitForFunctionOptions) (*JSHandle, error) {
	req := map[string]any{
		"expression": expression,
		"arg":        serializeArgument(arg),
		"timeout":    defaultActionTimeoutMs, // required by server; 0 disables timeout
	}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.PollingInterval != 0 {
			req["pollingInterval"] = o.PollingInterval
		}
		if o.Timeout < 0 {
			req["timeout"] = 0.0 // negative → send 0, which disables the timeout on the server
		} else if o.Timeout > 0 {
			req["timeout"] = o.Timeout
		}
	}
	result, err := p.mainFrame.SendMessageRequest(ctx, "waitForFunction", req)
	if err != nil {
		return nil, fmt.Errorf("page.waitForFunction failed: %w", err)
	}
	var resp struct {
		Handle struct {
			Guid string `json:"guid"`
		} `json:"handle"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse waitForFunction response: %w", err)
	}
	return &JSHandle{owner: p.mainFrame.child(resp.Handle.Guid)}, nil
}

// WaitForTimeout pauses execution for the specified number of milliseconds.
func (p *Page) WaitForTimeout(ctx context.Context, ms float64) error {
	_, err := p.mainFrame.SendMessageRequest(ctx, "waitForTimeout", map[string]any{
		"waitTimeout": ms,
	})
	if err != nil {
		return fmt.Errorf("page.waitForTimeout failed: %w", err)
	}
	return nil
}

// PageErrors returns all uncaught exceptions collected since the last ClearPageErrors call.
// Each entry is formatted as "Name: message" (e.g. "Error: something went wrong").
// Errors are tracked via the context-level "pageerror" event started by subscribeToPageErrors.
func (p *Page) PageErrors(_ context.Context) ([]string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]string, len(p.collectedPageErrors))
	copy(out, p.collectedPageErrors)
	return out, nil
}

// ClearPageErrors resets the accumulated page error list.
func (p *Page) ClearPageErrors(_ context.Context) error {
	p.mu.Lock()
	p.collectedPageErrors = p.collectedPageErrors[:0]
	p.mu.Unlock()
	return nil
}

// subscribeToPageErrors listens for the context-level "pageError" event and
// accumulates uncaught JS exceptions for this page in collectedPageErrors.
// Wire format: { error: { error: { name, message, stack } }, page: { guid }, location: {...} }
func (p *Page) subscribeToPageErrors() {
	if p.browserContext == nil {
		return
	}
	id := p.owner.conn.OnEvent(p.browserContext.owner.guid, "pageError", func(params json.RawMessage) {
		var event struct {
			Page struct {
				Guid string `json:"guid"`
			} `json:"page"`
			Error struct {
				Error struct {
					Name    string `json:"name"`
					Message string `json:"message"`
				} `json:"error"`
			} `json:"error"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Page.Guid != p.owner.guid {
			return
		}
		name := event.Error.Error.Name
		if name == "" {
			name = "Error"
		}
		msg := name + ": " + event.Error.Error.Message
		p.mu.Lock()
		p.collectedPageErrors = append(p.collectedPageErrors, msg)
		p.mu.Unlock()
	})
	p.mu.Lock()
	p.pageErrorListenID = id
	p.mu.Unlock()
}

// subscribeToPageClose listens to the page-level "close" event, cleans up state,
// and calls context-level OnPageClose handlers.
func (p *Page) subscribeToPageClose() {
	if p.browserContext == nil {
		return
	}
	p.owner.conn.OnEventOnce(p.owner.guid, "close", func(_ json.RawMessage) {
		p.mu.Lock()
		p.closed = true
		navID := p.navListenID
		loadID := p.pageLoadListenID
		errID := p.pageErrorListenID
		bCtx := p.browserContext
		p.mu.Unlock()

		if bCtx != nil {
			bCtx.removePage(p)
			if navID != 0 {
				p.owner.conn.OffEvent(p.mainFrame.guid, "navigated", navID)
			}
			if loadID != 0 {
				p.owner.conn.OffEvent(p.mainFrame.guid, "loadstate", loadID)
			}
			if errID != 0 {
				p.owner.conn.OffEvent(bCtx.owner.guid, "pageError", errID)
			}
			bCtx.callPageCloseHandlers(p)
		}
	})
}

// subscribeToPageLoad listens to the frame-level "loadstate" event and calls
// context-level OnPageLoad handlers when the "load" lifecycle event fires.
func (p *Page) subscribeToPageLoad() {
	if p.browserContext == nil {
		return
	}
	id := p.owner.conn.OnEvent(p.mainFrame.guid, "loadstate", func(params json.RawMessage) {
		var ev protocol.FrameLoadstateEvent
		if err := json.Unmarshal(params, &ev); err != nil {
			return
		}
		if ev.Add != nil && string(*ev.Add) == "load" {
			p.browserContext.callPageLoadHandlers(p)
		}
	})
	p.mu.Lock()
	p.pageLoadListenID = id
	p.mu.Unlock()
}

// subscribeToNavigation registers a listener for the main frame's "navigated" event,
// keeping p.frameInit.URL in sync without additional IPC calls.
// Called internally by BrowserContext.NewPage and OnPage after page construction.
func (p *Page) subscribeToNavigation() {
	listenerID := p.owner.conn.OnEvent(p.mainFrame.guid, "navigated", func(params json.RawMessage) {
		var nav protocol.FrameNavigatedEvent
		if err := json.Unmarshal(params, &nav); err != nil {
			return
		}
		p.mu.Lock()
		p.frameInit.URL = nav.Url
		p.mu.Unlock()
	})
	p.mu.Lock()
	if p.closed {
		// Close() ran between OnEvent and here — unregister immediately.
		p.owner.conn.OffEvent(p.mainFrame.guid, "navigated", listenerID)
	} else {
		p.navListenID = listenerID
	}
	p.mu.Unlock()
}
