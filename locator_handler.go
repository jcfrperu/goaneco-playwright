package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// AddLocatorHandlerOptions configures the behavior of AddLocatorHandler.
type AddLocatorHandlerOptions struct {
	// NoWaitAfter disables waiting for the handler's locator to become hidden after
	// the handler returns. Useful when the handler itself waits for the element to hide.
	NoWaitAfter *bool
	// Times limits how many times the handler is invoked. After Times calls the
	// handler is automatically unregistered. Zero or nil means unlimited.
	Times *int
}

// AddLocatorHandler registers a handler that Playwright calls automatically when the
// given locator becomes visible during a page interaction (e.g. a cookie banner overlay).
// The handler must dismiss or hide the element before returning.
// The returned cancel function unregisters the handler.
func (p *Page) AddLocatorHandler(ctx context.Context, locator *Locator, handler func(*Locator), opts ...AddLocatorHandlerOptions) (func(), error) {
	var opt AddLocatorHandlerOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	req := protocol.PageRegisterLocatorHandlerRequest{
		Selector:    locator.selector,
		NoWaitAfter: opt.NoWaitAfter,
	}
	result, err := p.owner.SendMessageRequest(ctx, "registerLocatorHandler", req)
	if err != nil {
		return nil, fmt.Errorf("page.addLocatorHandler: register failed: %w", err)
	}
	var resp protocol.PageRegisterLocatorHandlerResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("page.addLocatorHandler: failed to parse response: %w", err)
	}
	uid := resp.Uid

	var callCount atomic.Int32
	var maxCalls int32
	if opt.Times != nil && *opt.Times > 0 {
		maxCalls = int32(*opt.Times)
	}

	var cancelFuncPtr atomic.Pointer[func()]

	id := p.owner.conn.OnEvent(p.owner.guid, "locatorHandlerTriggered", func(params json.RawMessage) {
		var event protocol.PageLocatorHandlerTriggeredEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Uid != uid {
			return
		}
		go func() {
			// Always resolve the handler to unblock the Playwright server, even if quota is
			// exceeded or the handler panics. Without this, the page action that triggered
			// the handler (e.g. click) would hang until its own timeout expires.
			defer func() {
				_ = recover()
				resolveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				resolveReq := protocol.PageResolveLocatorHandlerNoReplyRequest{Uid: uid}
				_, _ = p.owner.SendMessageRequest(resolveCtx, "resolveLocatorHandlerNoReply", resolveReq)
			}()

			if maxCalls > 0 && callCount.Add(1) > maxCalls {
				return // quota already reached; resolve sent by defer above
			}
			handler(locator)
			if maxCalls > 0 && callCount.Load() >= maxCalls {
				if f := cancelFuncPtr.Load(); f != nil {
					(*f)()
				}
			}
		}()
	})

	cancelFunc := func() {
		p.owner.conn.OffEvent(p.owner.guid, "locatorHandlerTriggered", id)
		unregCtx, ucancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer ucancel()
		req := protocol.PageUnregisterLocatorHandlerRequest{Uid: uid}
		_, _ = p.owner.SendMessageRequest(unregCtx, "unregisterLocatorHandler", req)
	}
	cancelFuncPtr.Store(&cancelFunc)
	return cancelFunc, nil
}
