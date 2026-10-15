package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// Frame represents an iframe or the main frame of a page.
// Methods mirror Page's navigation/evaluation API but are scoped to the frame.
type Frame struct {
	owner        ChannelOwner
	page         *Page
	mu           sync.RWMutex
	name         string
	url          string
	detached     bool
	loadStates   map[string]bool
	parentGUID   string
	navListenID  connection.ListenerID // "navigated" listener ID registered in subscribeToFrames; cleaned up on detach
	loadListenID connection.ListenerID // "loadstate" listener ID registered in subscribeToFrames
}

// Page returns the page that this frame belongs to.
func (f *Frame) Page() *Page {
	return f.page
}

// ParentFrame returns the parent frame of this frame, or nil if this is the main frame.
func (f *Frame) ParentFrame() *Frame {
	f.mu.RLock()
	parentGUID := f.parentGUID
	f.mu.RUnlock()
	if parentGUID == "" || f.page == nil {
		return nil
	}
	return f.page.findFrameByGUID(parentGUID)
}

// URL returns the frame's current URL.
func (f *Frame) URL() string {
	if f.page != nil && f.owner.guid == f.page.mainFrame.guid {
		return f.page.URL()
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.url
}

// Name returns the frame's name attribute from the embedding element.
func (f *Frame) Name() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.name
}

// IsDetached reports whether this frame has been detached from the page.
func (f *Frame) IsDetached() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.detached
}

// Goto navigates the frame to the given URL.
func (f *Frame) Goto(ctx context.Context, url string) error {
	params := frameGotoParams{URL: url, Timeout: defaultActionTimeoutMs}
	_, err := f.owner.SendMessageRequest(ctx, "goto", params)
	if err != nil {
		return fmt.Errorf("frame.goto failed: %w", err)
	}
	return nil
}

// SetContent sets the HTML markup of the frame.
func (f *Frame) SetContent(ctx context.Context, html string, opts ...*SetContentOptions) error {
	timeout := f.owner.effectiveTimeout("setContent")
	var waitUntil *protocol.LifecycleEvent
	if len(opts) > 0 && opts[0] != nil {
		if opts[0].Timeout != nil {
			timeout = *opts[0].Timeout
		}
		if opts[0].WaitUntil != nil {
			le := protocol.LifecycleEvent(*opts[0].WaitUntil)
			waitUntil = &le
		}
	}
	params := frameSetContentParams{HTML: html, Timeout: timeout, WaitUntil: waitUntil}
	_, err := f.owner.SendMessageRequest(ctx, "setContent", params)
	if err != nil {
		return fmt.Errorf("frame.setContent failed: %w", err)
	}
	return nil
}

// Content returns the full HTML of the frame.
func (f *Frame) Content(ctx context.Context, opts ...*PageContentOptions) (string, error) {
	req := protocol.FrameContentRequest{}
	if len(opts) > 0 && opts[0] != nil {
		req.IncludeShadow = opts[0].IncludeShadow
	}
	result, err := f.owner.SendMessageRequest(ctx, "content", req)
	if err != nil {
		return "", fmt.Errorf("frame.content failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse frame.content response: %w", err)
	}
	return resp.Value, nil
}

// Title returns the document title of the frame.
func (f *Frame) Title(ctx context.Context) (string, error) {
	result, err := f.owner.SendMessageRequest(ctx, "title", struct{}{})
	if err != nil {
		return "", fmt.Errorf("frame.title failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse frame.title response: %w", err)
	}
	return resp.Value, nil
}

// Evaluate executes a JavaScript expression in the frame and returns the deserialized result.
func (f *Frame) Evaluate(ctx context.Context, expression string, arg ...any) (any, error) {
	var inputArg any
	if len(arg) > 0 {
		inputArg = arg[0]
	}
	params := frameEvaluateExpressionParams{
		Expression: expression,
		World:      "main",
		Arg:        serializeArgument(inputArg),
	}
	result, err := f.owner.SendMessageRequest(ctx, "evaluateExpression", params)
	if err != nil {
		return nil, fmt.Errorf("frame.evaluate failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse frame.evaluate response: %w", err)
	}
	return deserializeValue(resp.Value)
}

// Locator returns a new Locator scoped to this frame.
func (f *Frame) Locator(selector string) *Locator {
	var bc *BrowserContext
	if f.page != nil {
		bc = f.page.browserContext
	}
	return &Locator{frame: f.owner, selector: selector, browserContext: bc}
}

// GetByTestId returns a Locator scoped to this frame matching elements by the configured testid attribute.
func (f *Frame) GetByTestId(testId any) *Locator {
	attr := "data-testid"
	var bc *BrowserContext
	if f.page != nil {
		bc = f.page.browserContext
		if bc != nil {
			attr = bc.TestIdAttributeName()
		}
	}
	return &Locator{frame: f.owner, selector: getByTestIdSelector(attr, testId), browserContext: bc}
}

// SelectOption selects one or more options in a <select> element found by selector.
func (f *Frame) SelectOption(ctx context.Context, selector string, values ...any) ([]string, error) {
	return f.Locator(selector).SelectOption(ctx, values...)
}

// QuerySelector returns the first element matching the selector in this frame.
func (f *Frame) QuerySelector(ctx context.Context, selector string) (*ElementHandle, error) {
	result, err := f.owner.SendMessageRequest(ctx, "querySelector", map[string]string{"selector": selector})
	if err != nil {
		return nil, fmt.Errorf("frame.querySelector failed: %w", err)
	}
	var resp struct {
		Element *struct {
			Guid string `json:"guid"`
		} `json:"element,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse frame.querySelector response: %w", err)
	}
	if resp.Element == nil {
		return nil, nil
	}
	return elementHandleFromGUID(f.owner, resp.Element.Guid, f.page), nil
}

// WaitForLoadState waits until the frame reaches the given load state.
// Valid states: "load" (default), "domcontentloaded", "networkidle", "commit".
func (f *Frame) WaitForLoadState(ctx context.Context, state ...string) error {
	want := "load"
	if len(state) > 0 && state[0] != "" {
		want = state[0]
	}
	if want == "commit" {
		return nil
	}

	f.mu.Lock()
	if f.loadStates == nil {
		f.loadStates = make(map[string]bool)
	}
	if f.loadStates[want] {
		f.mu.Unlock()
		return nil
	}
	f.mu.Unlock()

	matched := make(chan struct{}, 1)
	id := f.owner.conn.OnEvent(f.owner.guid, "loadstate", func(params json.RawMessage) {
		var ev protocol.FrameLoadstateEvent
		if err := json.Unmarshal(params, &ev); err != nil {
			return
		}
		if ev.Add != nil && string(*ev.Add) == want {
			select {
			case matched <- struct{}{}:
			default:
			}
		}
	})
	defer f.owner.conn.OffEvent(f.owner.guid, "loadstate", id)

	f.mu.Lock()
	if f.loadStates[want] {
		f.mu.Unlock()
		return nil
	}
	f.mu.Unlock()

	select {
	case <-matched:
		f.mu.Lock()
		f.loadStates[want] = true
		f.mu.Unlock()
		return nil
	case <-ctx.Done():
		return fmt.Errorf("frame.waitForLoadState(%q) failed: %w", want, ctx.Err())
	}
}

// Press presses a keyboard key while the element matching selector is focused.
func (f *Frame) Press(ctx context.Context, selector, key string) error {
	return f.Locator(selector).Press(ctx, key)
}

// frameInitializerFull extends frameInitializer to also capture the parentFrame GUID.
type frameInitializerFull struct {
	URL         string   `json:"url"`
	Name        string   `json:"name"`
	LoadStates  []string `json:"loadStates"`
	ParentFrame *struct {
		GUID string `json:"guid"`
	} `json:"parentFrame,omitempty"`
}

// subscribeToFrames listens for frameAttached/frameDetached events on the page,
// maintaining the frames and framesByGUID maps on the page.
func (p *Page) subscribeToFrames() {
	p.owner.conn.OnEvent(p.owner.guid, "frameAttached", func(params json.RawMessage) {
		var event struct {
			Frame struct {
				Guid string `json:"guid"`
			} `json:"frame"`
		}
		if err := json.Unmarshal(params, &event); err != nil || event.Frame.Guid == "" {
			return
		}
		guid := event.Frame.Guid

		var fi frameInitializerFull
		if raw := p.owner.Initializer(guid); len(raw) > 0 {
			_ = json.Unmarshal(raw, &fi)
		}

		loadStates := make(map[string]bool)
		for _, ls := range fi.LoadStates {
			loadStates[ls] = true
		}

		parentGUID := ""
		if fi.ParentFrame != nil {
			parentGUID = fi.ParentFrame.GUID
		}

		frame := &Frame{
			owner:      p.owner.child(guid),
			page:       p,
			name:       fi.Name,
			url:        fi.URL,
			loadStates: loadStates,
			parentGUID: parentGUID,
		}

		// Track URL and name updates via navigated events.
		// The frame name is often empty in the initializer and populated via the first navigated event.
		navID := p.owner.conn.OnEvent(guid, "navigated", func(params json.RawMessage) {
			var nav protocol.FrameNavigatedEvent
			if err := json.Unmarshal(params, &nav); err != nil {
				return
			}
			frame.mu.Lock()
			frame.url = nav.Url
			if nav.Name != "" {
				frame.name = nav.Name
			}
			frame.loadStates = make(map[string]bool)
			frame.loadStates["commit"] = true
			frame.mu.Unlock()

			p.callFrameNavigatedHandlers(frame)
			if p.browserContext != nil {
				p.browserContext.callFrameNavigatedHandlers(frame)
			}
		})
		frame.navListenID = navID

		loadID := p.owner.conn.OnEvent(guid, "loadstate", func(params json.RawMessage) {
			var ev protocol.FrameLoadstateEvent
			if err := json.Unmarshal(params, &ev); err != nil {
				return
			}
			frame.mu.Lock()
			if frame.loadStates == nil {
				frame.loadStates = make(map[string]bool)
			}
			if ev.Add != nil {
				frame.loadStates[string(*ev.Add)] = true
			}
			if ev.Remove != nil {
				delete(frame.loadStates, string(*ev.Remove))
			}
			frame.mu.Unlock()
		})
		frame.loadListenID = loadID

		p.mu.Lock()
		p.frames = append(p.frames, frame)
		if p.framesByGUID == nil {
			p.framesByGUID = make(map[string]*Frame)
		}
		p.framesByGUID[guid] = frame
		p.mu.Unlock()

		p.callFrameAttachedHandlers(frame)
		if p.browserContext != nil {
			p.browserContext.callFrameAttachedHandlers(frame)
		}
	})

	p.owner.conn.OnEvent(p.owner.guid, "frameDetached", func(params json.RawMessage) {
		var event struct {
			Frame struct {
				Guid string `json:"guid"`
			} `json:"frame"`
		}
		if err := json.Unmarshal(params, &event); err != nil || event.Frame.Guid == "" {
			return
		}
		guid := event.Frame.Guid

		p.mu.Lock()
		var detachedFrame *Frame
		if f, ok := p.framesByGUID[guid]; ok {
			detachedFrame = f
			if f.navListenID != 0 {
				p.owner.conn.OffEvent(guid, "navigated", f.navListenID)
			}
			if f.loadListenID != 0 {
				p.owner.conn.OffEvent(guid, "loadstate", f.loadListenID)
			}
			f.mu.Lock()
			f.detached = true
			f.mu.Unlock()
			delete(p.framesByGUID, guid)
		}
		for i, f := range p.frames {
			if f.owner.guid == guid {
				if detachedFrame == nil {
					detachedFrame = f
				}
				copy(p.frames[i:], p.frames[i+1:])
				p.frames[len(p.frames)-1] = nil
				p.frames = p.frames[:len(p.frames)-1]
				break
			}
		}
		p.mu.Unlock()

		if detachedFrame != nil {
			p.callFrameDetachedHandlers(detachedFrame)
			if p.browserContext != nil {
				p.browserContext.callFrameDetachedHandlers(detachedFrame)
			}
		}
	})
}

// WaitForURL waits until the frame URL matches the given pattern.
// WaitForURL waits until the frame's URL matches pattern.
// Options can be a time.Duration or a *FrameWaitForURLOptions.
func (f *Frame) WaitForURL(ctx context.Context, pattern string, opts ...any) error {
	to := 30 * time.Second
	var waitUntil *protocol.LifecycleEvent
	if len(opts) > 0 {
		switch o := opts[0].(type) {
		case time.Duration:
			to = o
		case *FrameWaitForURLOptions:
			if o != nil {
				if o.Timeout != nil {
					to = time.Duration(*o.Timeout) * time.Millisecond
				}
				waitUntil = o.WaitUntil
			}
		}
	}

	matched := make(chan struct{}, 1)
	id := f.owner.conn.OnEvent(f.owner.guid, "navigated", func(params json.RawMessage) {
		var nav protocol.FrameNavigatedEvent
		if err := json.Unmarshal(params, &nav); err != nil {
			return
		}
		if urlMatchesPattern(nav.Url, pattern) {
			f.mu.Lock()
			f.url = nav.Url
			f.mu.Unlock()
			select {
			case matched <- struct{}{}:
			default:
			}
		}
	})
	defer f.owner.conn.OffEvent(f.owner.guid, "navigated", id)

	tctx, tcancel := context.WithTimeout(ctx, to)
	defer tcancel()

	if !urlMatchesPattern(f.URL(), pattern) {
		select {
		case <-matched:
		case <-tctx.Done():
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("frame.WaitForURL: Timeout %dms exceeded after %v waiting for %q (current: %q)", to.Milliseconds(), to, pattern, f.URL())
		}
	}

	if waitUntil != nil {
		state := string(*waitUntil)
		if state != "commit" {
			if err := f.WaitForLoadState(tctx, state); err != nil {
				return err
			}
		}
	}
	return nil
}

// EvaluateHandle executes a JavaScript expression in the context of the frame and returns the result as a JSHandle.
func (f *Frame) EvaluateHandle(ctx context.Context, expression string, arg ...any) (*JSHandle, error) {
	var inputArg any
	if len(arg) > 0 {
		inputArg = arg[0]
	}

	params := frameEvaluateExpressionParams{
		Expression: expression,
		World:      "main",
		Arg:        serializeArgument(inputArg),
	}

	result, err := f.owner.SendMessageRequest(ctx, "evaluateExpressionHandle", params)
	if err != nil {
		return nil, fmt.Errorf("frame.evaluateHandle failed: %w", err)
	}

	var resp struct {
		Handle struct {
			Guid string `json:"guid"`
		} `json:"handle"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse evaluateHandle response: %w", err)
	}

	return &JSHandle{owner: f.owner.child(resp.Handle.Guid), page: f.page}, nil
}

// ChildFrames returns all child frames of this frame.
func (f *Frame) ChildFrames() []*Frame {
	if f.page == nil {
		return nil
	}
	f.page.mu.RLock()
	defer f.page.mu.RUnlock()
	var children []*Frame
	isMain := f.page.mainFrameObj != nil && f == f.page.mainFrameObj
	for _, child := range f.page.frames {
		if child.parentGUID == f.owner.guid || (child.parentGUID == "" && isMain) {
			children = append(children, child)
		}
	}
	return children
}
