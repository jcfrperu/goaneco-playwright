package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// WebSocketFrame represents a single frame transmitted over a WebSocket connection.
// Opcode 1 = text frame, 2 = binary frame.
type WebSocketFrame struct {
	Data   string
	Opcode int
}

// WebSocketWaitForFrameOptions configures a frame wait predicate.
type WebSocketWaitForFrameOptions struct {
	Predicate func(WebSocketFrame) bool
}

// Text returns the frame payload as a UTF-8 string.
// For binary frames this is still the raw string representation of the data.
func (f WebSocketFrame) Text() string { return f.Data }

// WebSocket represents an active WebSocket connection opened from a page.
// Event handlers are goroutine-safe; all On* calls may be called from any goroutine.
type WebSocket struct {
	owner ChannelOwner
	url   string
	page  *Page

	mu       sync.Mutex
	isClosed bool
	nextID   int

	frameReceivedHandlers map[int]func(WebSocketFrame)
	frameSentHandlers     map[int]func(WebSocketFrame)
	closeHandlers         map[int]func()
	errorHandlers         map[int]func(string)
	listenerIDs           [4]connection.ListenerID
}

// WaitForFrameReceived waits for a received frame matching the optional predicate.
// The optional trigger runs after the listener is registered.
func (ws *WebSocket) WaitForFrameReceived(ctx context.Context, trigger func() error, opts ...*WebSocketWaitForFrameOptions) (WebSocketFrame, error) {
	return ws.waitForFrame(ctx, "frameReceived", trigger, opts...)
}

// WaitForFrameSent waits for a sent frame matching the optional predicate.
// The optional trigger runs after the listener is registered.
func (ws *WebSocket) WaitForFrameSent(ctx context.Context, trigger func() error, opts ...*WebSocketWaitForFrameOptions) (WebSocketFrame, error) {
	return ws.waitForFrame(ctx, "frameSent", trigger, opts...)
}

func (ws *WebSocket) waitForFrame(ctx context.Context, event string, trigger func() error, opts ...*WebSocketWaitForFrameOptions) (WebSocketFrame, error) {
	var zero WebSocketFrame
	name := "webSocket.waitForFrameReceived"
	subscribe := ws.OnFrameReceived
	if event == "frameSent" {
		name = "webSocket.waitForFrameSent"
		subscribe = ws.OnFrameSent
	}
	var predicate func(WebSocketFrame) bool
	if len(opts) > 0 && opts[0] != nil {
		predicate = opts[0].Predicate
	}
	if predicate == nil {
		predicate = func(WebSocketFrame) bool { return true }
	}

	type waitResult struct {
		frame WebSocketFrame
		err   error
	}
	resultCh := make(chan waitResult, 1)
	sendResult := func(result waitResult) {
		select {
		case resultCh <- result:
		default:
		}
	}
	offFrame := subscribe(func(frame WebSocketFrame) {
		if predicate(frame) {
			sendResult(waitResult{frame: frame})
		}
	})
	defer offFrame()
	offClose := ws.OnClose(func() {
		sendResult(waitResult{err: fmt.Errorf("Socket closed")})
	})
	defer offClose()
	offError := ws.OnError(func(message string) {
		sendResult(waitResult{err: fmt.Errorf("Socket error: %s", message)})
	})
	defer offError()

	if ws.IsClosed() {
		sendResult(waitResult{err: fmt.Errorf("Socket closed")})
	}
	if trigger != nil {
		if err := trigger(); err != nil {
			return zero, fmt.Errorf("%s: %w", name, err)
		}
	}
	if ws.page != nil {
		ws.page.mu.RLock()
		pageClosed := ws.page.closed
		ws.page.mu.RUnlock()
		if pageClosed {
			return zero, fmt.Errorf("%s: Target page, context or browser has been closed", name)
		}
	}

	select {
	case result := <-resultCh:
		if result.err != nil {
			return zero, fmt.Errorf("%s: %w", name, result.err)
		}
		return result.frame, nil
	case <-ctx.Done():
		return zero, fmt.Errorf("%s: Timeout exceeded: %w", name, ctx.Err())
	}
}

// webSocketInitializer is the wire format of the WebSocket channel object initializer.
type webSocketInitializer struct {
	URL string `json:"url"`
}

// URL returns the WebSocket endpoint URL (e.g. "ws://localhost:1234/ws").
func (ws *WebSocket) URL() string { return ws.url }

// IsClosed reports whether the WebSocket connection has been closed.
func (ws *WebSocket) IsClosed() bool {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return ws.isClosed
}

// OnFrameReceived registers a handler called for every frame received from the server.
// The returned function cancels the listener.
func (ws *WebSocket) OnFrameReceived(handler func(WebSocketFrame)) func() {
	ws.mu.Lock()
	if ws.frameReceivedHandlers == nil {
		ws.frameReceivedHandlers = make(map[int]func(WebSocketFrame))
	}
	id := ws.nextID
	ws.nextID++
	ws.frameReceivedHandlers[id] = handler
	ws.mu.Unlock()
	return func() {
		ws.mu.Lock()
		defer ws.mu.Unlock()
		delete(ws.frameReceivedHandlers, id)
	}
}

// OnFrameSent registers a handler called for every frame sent to the server.
// The returned function cancels the listener.
func (ws *WebSocket) OnFrameSent(handler func(WebSocketFrame)) func() {
	ws.mu.Lock()
	if ws.frameSentHandlers == nil {
		ws.frameSentHandlers = make(map[int]func(WebSocketFrame))
	}
	id := ws.nextID
	ws.nextID++
	ws.frameSentHandlers[id] = handler
	ws.mu.Unlock()
	return func() {
		ws.mu.Lock()
		defer ws.mu.Unlock()
		delete(ws.frameSentHandlers, id)
	}
}

// OnClose registers a handler called when the WebSocket connection closes.
// The returned function cancels the listener.
func (ws *WebSocket) OnClose(handler func()) func() {
	ws.mu.Lock()
	if ws.closeHandlers == nil {
		ws.closeHandlers = make(map[int]func())
	}
	id := ws.nextID
	ws.nextID++
	ws.closeHandlers[id] = handler
	ws.mu.Unlock()
	return func() {
		ws.mu.Lock()
		defer ws.mu.Unlock()
		delete(ws.closeHandlers, id)
	}
}

// OnError registers a handler called when the WebSocket encounters a socket error.
// The error string is passed to the handler. The returned function cancels the listener.
func (ws *WebSocket) OnError(handler func(string)) func() {
	ws.mu.Lock()
	if ws.errorHandlers == nil {
		ws.errorHandlers = make(map[int]func(string))
	}
	id := ws.nextID
	ws.nextID++
	ws.errorHandlers[id] = handler
	ws.mu.Unlock()
	return func() {
		ws.mu.Lock()
		defer ws.mu.Unlock()
		delete(ws.errorHandlers, id)
	}
}

// subscribe wires up the IPC event listeners for this WebSocket object.
// Must be called once after the WebSocket is created.
func (ws *WebSocket) subscribe() {
	guid := ws.owner.guid
	conn := ws.owner.conn

	ws.listenerIDs[0] = conn.OnEvent(guid, "frameReceived", func(params json.RawMessage) {
		var event protocol.WebSocketFrameReceivedEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		frame := WebSocketFrame{Data: event.Data, Opcode: event.Opcode}
		ws.mu.Lock()
		handlers := make([]func(WebSocketFrame), 0, len(ws.frameReceivedHandlers))
		for _, h := range ws.frameReceivedHandlers {
			handlers = append(handlers, h)
		}
		ws.mu.Unlock()
		for _, h := range handlers {
			h(frame)
		}
	})

	ws.listenerIDs[1] = conn.OnEvent(guid, "frameSent", func(params json.RawMessage) {
		var event protocol.WebSocketFrameSentEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		frame := WebSocketFrame{Data: event.Data, Opcode: event.Opcode}
		ws.mu.Lock()
		handlers := make([]func(WebSocketFrame), 0, len(ws.frameSentHandlers))
		for _, h := range ws.frameSentHandlers {
			handlers = append(handlers, h)
		}
		ws.mu.Unlock()
		for _, h := range handlers {
			h(frame)
		}
	})

	ws.listenerIDs[2] = conn.OnEvent(guid, "socketError", func(params json.RawMessage) {
		var event protocol.WebSocketSocketErrorEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		ws.mu.Lock()
		handlers := make([]func(string), 0, len(ws.errorHandlers))
		for _, h := range ws.errorHandlers {
			handlers = append(handlers, h)
		}
		ws.mu.Unlock()
		for _, h := range handlers {
			h(event.Error)
		}
	})

	ws.listenerIDs[3] = conn.OnEvent(guid, "close", func(_ json.RawMessage) {
		conn.OffEvent(guid, "frameReceived", ws.listenerIDs[0])
		conn.OffEvent(guid, "frameSent", ws.listenerIDs[1])
		conn.OffEvent(guid, "socketError", ws.listenerIDs[2])
		ws.mu.Lock()
		ws.isClosed = true
		handlers := make([]func(), 0, len(ws.closeHandlers))
		for _, h := range ws.closeHandlers {
			handlers = append(handlers, h)
		}
		ws.mu.Unlock()
		for _, h := range handlers {
			h()
		}
	})
}

// webSocketFromGUID constructs a WebSocket from a known GUID and its raw initializer.
func webSocketFromGUID(parent ChannelOwner, guid string, initRaw json.RawMessage) *WebSocket {
	var init webSocketInitializer
	if len(initRaw) > 0 {
		_ = json.Unmarshal(initRaw, &init) // best-effort; zero-value fallback
	}
	ws := &WebSocket{
		owner: parent.child(guid),
		url:   init.URL,
	}
	ws.subscribe()
	return ws
}
