package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
)

// BindingHandler is called when JavaScript invokes an exposed binding.
// args contains the deserialized JavaScript arguments.
// The return value is serialized and sent back to JavaScript.
type BindingHandler func(args ...any) any

// BindingSource provides contextual information about the page and frame that invoked an exposed binding.
type BindingSource struct {
	context *BrowserContext
	page    *Page
	frame   *Frame
}

func (s *BindingSource) Context() *BrowserContext { return s.context }
func (s *BindingSource) Page() *Page              { return s.page }
func (s *BindingSource) Frame() *Frame            { return s.frame }

// BindingSourceHandler is called when JavaScript invokes an exposed binding with BindingSource.
type BindingSourceHandler func(source *BindingSource, args ...any) any

// ExposeBinding exposes a Go function as a global JavaScript function in all pages of this context.
// When JavaScript calls window[name](...args), handler is invoked with the deserialized args.
// Returns an error if the IPC registration fails. Re-registering the same name replaces the handler.
func (c *BrowserContext) ExposeBinding(ctx context.Context, name string, handler any) error {
	_, err := c.owner.SendMessageRequest(ctx, "exposeBinding", map[string]string{"name": name})
	if err != nil {
		return fmt.Errorf("browserContext.exposeBinding(%q) failed: %w", name, err)
	}

	c.mu.Lock()
	// Remove any previous listener for this binding name to avoid duplicate dispatch.
	if c.bindingListenerIDs == nil {
		c.bindingListenerIDs = make(map[string]connection.ListenerID)
	}
	if prev, ok := c.bindingListenerIDs[name]; ok {
		c.owner.conn.OffEvent(c.owner.guid, "bindingCall", prev)
	}
	id := c.owner.conn.OnEvent(c.owner.guid, "bindingCall", func(params json.RawMessage) {
		var event struct {
			Binding struct {
				Guid string `json:"guid"`
			} `json:"binding"`
		}
		if err := json.Unmarshal(params, &event); err != nil || event.Binding.Guid == "" {
			return
		}
		// Filter: only handle calls for this binding name.
		raw := c.owner.Initializer(event.Binding.Guid)
		var init struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(raw, &init); err != nil || init.Name != name {
			return
		}
		guid := event.Binding.Guid
		c.dispatchBindingCall(guid, handler)
	})
	c.bindingListenerIDs[name] = id
	c.mu.Unlock()
	return nil
}

// dispatchBindingCall is called in a goroutine when the Playwright driver fires a "bindingCall" event.
// It deserializes the JS arguments from the binding initializer, invokes handler with a timeout,
// then sends either "resolve" (with the serialized return value) or "reject" (with the error)
// back to the driver so the JavaScript promise can settle.
func (c *BrowserContext) dispatchBindingCall(guid string, handler any) {
	raw := c.owner.Initializer(guid)
	var init struct {
		Frame struct {
			Guid string `json:"guid"`
		} `json:"frame"`
		Args []serializedValueRaw `json:"args"`
	}
	if err := json.Unmarshal(raw, &init); err != nil {
		return
	}

	args := make([]any, len(init.Args))
	for i, arg := range init.Args {
		v, err := deserializeValue(arg)
		if err != nil {
			slog.Default().Warn("dispatchBindingCall: failed to deserialize arg", "index", i, "error", err)
		}
		args[i] = v
	}

	var frame *Frame
	var page *Page
	if init.Frame.Guid != "" {
		c.mu.RLock()
		pages := c.pages
		c.mu.RUnlock()
		for _, p := range pages {
			if p.mainFrame.guid == init.Frame.Guid {
				frame = p.MainFrame()
				page = p
				break
			}
			p.mu.RLock()
			if f, ok := p.framesByGUID[init.Frame.Guid]; ok {
				frame = f
				page = p
			}
			p.mu.RUnlock()
			if frame != nil {
				break
			}
		}
	}
	if page == nil {
		c.mu.RLock()
		if len(c.pages) > 0 {
			page = c.pages[0]
		}
		c.mu.RUnlock()
	}
	if frame == nil && page != nil {
		frame = page.MainFrame()
	}

	source := &BindingSource{
		context: c,
		page:    page,
		frame:   frame,
	}

	bindingOwner := c.owner.child(guid)

	type outcome struct {
		value    any
		panicErr error
	}
	done := make(chan outcome, 1)
	handlerTimeout := time.Duration(defaultBindingHandlerTimeoutMs * float64(time.Millisecond))
	handlerCtx, handlerCancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer handlerCancel()

	go func() {
		var o outcome
		func() {
			defer func() {
				if r := recover(); r != nil {
					if e, ok := r.(error); ok {
						o.panicErr = e
					} else {
						o.panicErr = fmt.Errorf("%v", r)
					}
				}
			}()
			switch h := handler.(type) {
			case BindingHandler:
				o.value = h(args...)
			case BindingSourceHandler:
				o.value = h(source, args...)
			case func(source *BindingSource, args ...any) any:
				o.value = h(source, args...)
			case func(args ...any) any:
				o.value = h(args...)
			default:
				o.panicErr = fmt.Errorf("unsupported binding handler type: %T", handler)
			}
		}()
		done <- o
	}()

	var o outcome
	select {
	case o = <-done:
	case <-handlerCtx.Done():
		o.panicErr = fmt.Errorf("binding handler timed out after %v", handlerTimeout)
	}

	// Use a fresh context for IPC: the handler timeout may have already expired.
	ipcCtx, ipcCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	defer ipcCancel()

	if o.panicErr != nil {
		if _, err := bindingOwner.SendMessageRequest(ipcCtx, "reject", map[string]any{
			"error": map[string]any{
				"error": map[string]any{
					"message": o.panicErr.Error(),
					"name":    "Error",
				},
			},
		}); err != nil {
			slog.Default().Error("dispatchBindingCall: failed to send reject", "error", err)
		}
		return
	}
	if _, err := bindingOwner.SendMessageRequest(ipcCtx, "resolve", map[string]any{
		"result": serializeArgument(o.value),
	}); err != nil {
		slog.Default().Error("dispatchBindingCall: failed to send resolve", "error", err)
	}
}
