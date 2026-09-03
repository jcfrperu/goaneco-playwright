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
	owner       ChannelOwner
	page        *Page
	mu          sync.RWMutex
	name        string
	url         string
	detached    bool
	navListenID connection.ListenerID // "navigated" listener ID registered in subscribeToFrames; cleaned up on detach
}

// URL returns the frame's current URL.
func (f *Frame) URL() string {
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
func (f *Frame) SetContent(ctx context.Context, html string) error {
	params := frameSetContentParams{HTML: html, Timeout: defaultActionTimeoutMs}
	_, err := f.owner.SendMessageRequest(ctx, "setContent", params)
	if err != nil {
		return fmt.Errorf("frame.setContent failed: %w", err)
	}
	return nil
}

// Content returns the full HTML of the frame.
func (f *Frame) Content(ctx context.Context) (string, error) {
	result, err := f.owner.SendMessageRequest(ctx, "content", struct{}{})
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
	return &Locator{frame: f.owner, selector: selector}
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
	return elementHandleFromGUID(f.owner, resp.Element.Guid), nil
}

// WaitForLoadState waits until the frame reaches the given load state.
// Valid states: "load" (default), "domcontentloaded", "networkidle".
func (f *Frame) WaitForLoadState(ctx context.Context, state ...string) error {
	want := "load"
	if len(state) > 0 && state[0] != "" {
		want = state[0]
	}
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

	readyState, err := f.Evaluate(ctx, "document.readyState")
	if err == nil {
		switch want {
		case "load":
			if readyState == "complete" {
				return nil
			}
		case "domcontentloaded":
			if readyState == "interactive" || readyState == "complete" {
				return nil
			}
		}
	}

	timer := time.NewTimer(time.Duration(defaultActionTimeoutMs) * time.Millisecond)
	defer timer.Stop()

	select {
	case <-matched:
		return nil
	case <-timer.C:
		return fmt.Errorf("frame.waitForLoadState(%q): timeout", want)
	case <-ctx.Done():
		return fmt.Errorf("frame.waitForLoadState(%q): %w", want, ctx.Err())
	}
}

// Press presses a keyboard key while the element matching selector is focused.
func (f *Frame) Press(ctx context.Context, selector, key string) error {
	return f.Locator(selector).Press(ctx, key)
}

// frameInitializerFull extends frameInitializer to also capture the parentFrame GUID.
type frameInitializerFull struct {
	URL         string `json:"url"`
	Name        string `json:"name"`
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

		frame := &Frame{
			owner: p.owner.child(guid),
			page:  p,
			name:  fi.Name,
			url:   fi.URL,
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
			frame.mu.Unlock()
		})
		frame.navListenID = navID

		p.mu.Lock()
		p.frames = append(p.frames, frame)
		if p.framesByGUID == nil {
			p.framesByGUID = make(map[string]*Frame)
		}
		p.framesByGUID[guid] = frame
		p.mu.Unlock()
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
		if f, ok := p.framesByGUID[guid]; ok {
			if f.navListenID != 0 {
				p.owner.conn.OffEvent(guid, "navigated", f.navListenID)
			}
			f.mu.Lock()
			f.detached = true
			f.mu.Unlock()
			delete(p.framesByGUID, guid)
		}
		for i, f := range p.frames {
			if f.owner.guid == guid {
				p.frames = append(p.frames[:i], p.frames[i+1:]...)
				break
			}
		}
		p.mu.Unlock()
	})
}
