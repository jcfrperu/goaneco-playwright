package playwright

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// CDPSession provides a raw Chrome DevTools Protocol (CDP) session
// for direct protocol communication with Chromium.
type CDPSession struct {
	owner ChannelOwner
}

// Send sends a CDP command and returns the result.
// Method is the CDP domain and method name (e.g. "Runtime.evaluate").
// Params are the command parameters (or nil for no params).
func (s *CDPSession) Send(ctx context.Context, method string, params any) (json.RawMessage, error) {
	result, err := s.owner.SendMessageRequest(ctx, "send", protocol.CDPSessionSendRequest{
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, fmt.Errorf("cdpSession.send(%q) failed: %w", method, err)
	}
	var resp protocol.CDPSessionSendResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse cdpSession.send response: %w", err)
	}
	// Re-marshal the result so callers get raw JSON they can unmarshal into their types.
	return json.Marshal(resp.Result)
}

// On registers a handler for CDP events with the given method name.
// The handler receives raw event parameters; method "event" receives the complete event.
// The returned function cancels the listener.
func (s *CDPSession) On(method string, handler func(json.RawMessage)) func() {
	process := func(params json.RawMessage) {
		if method == "event" {
			handler(append(json.RawMessage(nil), params...))
			return
		}
		var event protocol.CDPSessionEventEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Method != method {
			return
		}
		raw, err := json.Marshal(event.Params)
		if err != nil {
			return
		}
		// Runs on the listener's serial queue: ordered and panic-safe.
		handler(raw)
	}
	id := s.owner.conn.OnEvent(s.owner.guid, "event", process)
	return func() { s.owner.conn.OffEvent(s.owner.guid, "event", id) }
}

// Detach detaches the CDPSession from the target.
func (s *CDPSession) Detach(ctx context.Context) error {
	_, err := s.owner.SendMessageRequest(ctx, "detach", protocol.CDPSessionDetachRequest{})
	if err != nil {
		return fmt.Errorf("cdpSession.detach failed: %w", err)
	}
	return nil
}

// OnClose registers a one-time handler called when the CDPSession is closed.
// The returned ListenerID can be used to cancel the listener.
func (s *CDPSession) OnClose(handler func()) connection.ListenerID {
	return s.owner.conn.OnEventOnce(s.owner.guid, "close", func(_ json.RawMessage) {
		handler()
	})
}

// NewCDPSession creates a new CDP session for the given page or frame (Chromium only).
func (c *BrowserContext) NewCDPSession(ctx context.Context, target any) (*CDPSession, error) {
	req := protocol.BrowserContextNewCDPSessionRequest{}
	switch target := target.(type) {
	case *Page:
		if target == nil {
			return nil, fmt.Errorf("browserContext.newCDPSession: target must be a page or frame")
		}
		req.Page = &protocol.Page{Guid: target.owner.guid}
	case *Frame:
		if target == nil {
			return nil, fmt.Errorf("browserContext.newCDPSession: target must be a page or frame")
		}
		req.Frame = &protocol.Frame{Guid: target.owner.guid}
	default:
		return nil, fmt.Errorf("browserContext.newCDPSession: target must be a page or frame, got %T", target)
	}
	result, err := c.owner.SendMessageRequest(ctx, "newCDPSession", req)
	if err != nil {
		return nil, fmt.Errorf("browserContext.newCDPSession failed: %w", err)
	}
	var resp protocol.BrowserContextNewCDPSessionResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse newCDPSession response: %w", err)
	}
	if resp.Session.Guid == "" {
		return nil, fmt.Errorf("newCDPSession: server returned empty session GUID")
	}
	return &CDPSession{owner: c.owner.child(resp.Session.Guid)}, nil
}
