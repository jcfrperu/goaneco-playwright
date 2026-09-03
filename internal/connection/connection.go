package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
)

// eventKey is the composite key used to look up handlers for a specific object+method pair.
type eventKey struct {
	GUID   string
	Method string
}

// ListenerID is an opaque identifier returned by OnEvent that can be passed to OffEvent
// to deregister a specific handler without disturbing others registered on the same event.
type ListenerID uint64

// eventHandler pairs a handler function with its unique ListenerID so it can be removed
// individually even when multiple handlers are registered for the same event key.
type eventHandler struct {
	id ListenerID
	fn func(json.RawMessage)
}

// Connection manages the IPC state between the Go client and the Playwright Node.js driver process.
// It tracks in-flight RPC callbacks, registered protocol objects, and per-object event listeners.
// All public methods are safe for concurrent use.
type Connection struct {
	mu              sync.RWMutex
	closed          bool
	callbacks       map[int]chan Message
	objects         map[string]*ObjectRef
	eventListeners  map[eventKey][]eventHandler
	guidListeners   map[string][]eventKey // secondary index: GUID → registered event keys for O(1) dispose
	dispatchWg      sync.WaitGroup        // tracks in-flight async event handler goroutines
	logger          *slog.Logger
	nextID          atomic.Int32
	nextListenerID  atomic.Uint64
	transportSend   func([]byte) error
	playwrightOnce  sync.Once
	playwrightObj   *ObjectRef
	playwrightReady chan struct{}
}

// MessageError carries the error detail returned by the Playwright server in a response envelope.
type MessageError struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

// Message is the top-level JSON envelope exchanged over the IPC pipe with the Playwright driver.
// Outbound messages carry an ID and a method; inbound messages either carry an ID (RPC response)
// or a method and GUID (server-initiated event such as "__create__", "__dispose__", or any domain event).
type Message struct {
	ID     *int            `json:"id"`
	GUID   string          `json:"guid"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Error MessageError `json:"error"`
	} `json:"error"`
}

// NewConnection creates and returns a ready-to-use Connection.
// Call SetTransportSend before sending any requests, and Start reading from the transport
// to have incoming messages delivered via Dispatch.
func NewConnection() *Connection {
	return &Connection{
		callbacks:       make(map[int]chan Message),
		objects:         make(map[string]*ObjectRef),
		eventListeners:  make(map[eventKey][]eventHandler),
		guidListeners:   make(map[string][]eventKey),
		logger:          slog.Default(),
		playwrightReady: make(chan struct{}),
	}
}

// SetTransportSend injects the transport sending function to decouple transport from connection.
func (c *Connection) SetTransportSend(send func([]byte) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.transportSend = send
}

// Dispatch parses a raw IPC message from the Playwright driver and routes it to the
// appropriate handler:
//   - ID-bearing messages wake the waiting SendRequest goroutine.
//   - "__create__" registers a new protocol object and signals the Playwright root when it arrives.
//   - "__dispose__" removes the object and all its event listeners.
//   - All other messages are dispatched to registered OnEvent handlers, each in its own goroutine.
func (c *Connection) Dispatch(msg []byte) {
	var m Message
	if err := json.Unmarshal(msg, &m); err != nil {
		c.logger.Error("connection: failed to unmarshal message", "error", err)
		return
	}

	if m.ID != nil {
		c.mu.Lock()
		cb, ok := c.callbacks[*m.ID]
		if ok {
			delete(c.callbacks, *m.ID)
		}
		c.mu.Unlock()

		if ok {
			select {
			case cb <- m:
			default:
				// Channel already has a message or was drained by a ctx.Done() branch;
				// the caller already returned an error so the response can be dropped.
				c.logger.Warn("connection: response arrived after request context was cancelled, dropping", "id", *m.ID)
			}
		}
	} else if m.Method == "__create__" {
		var createParams struct {
			Type        string          `json:"type"`
			GUID        string          `json:"guid"`
			Initializer json.RawMessage `json:"initializer"`
		}
		if err := json.Unmarshal(m.Params, &createParams); err != nil {
			c.logger.Error("connection: failed to parse __create__ params", "error", err)
			return
		}
		if createParams.GUID == "" {
			c.logger.Warn("connection: received __create__ with empty GUID, ignoring")
			return
		}
		obj := NewObjectRef(createParams.GUID, createParams.Initializer)
		c.RegisterObject(createParams.GUID, obj)
		if createParams.Type == "Playwright" {
			c.playwrightOnce.Do(func() {
				c.playwrightObj = obj
				close(c.playwrightReady)
			})
		}
	} else if m.Method == "__dispose__" {
		if m.GUID == "" {
			c.logger.Warn("connection: received __dispose__ with empty GUID, ignoring")
			return
		}
		c.mu.Lock()
		delete(c.objects, m.GUID)
		// O(1) cleanup using secondary index instead of scanning all event keys.
		for _, key := range c.guidListeners[m.GUID] {
			delete(c.eventListeners, key)
		}
		delete(c.guidListeners, m.GUID)
		c.mu.Unlock()
	} else {
		key := eventKey{GUID: m.GUID, Method: m.Method}
		c.mu.RLock()
		src := c.eventListeners[key]
		handlers := make([]eventHandler, len(src))
		copy(handlers, src)
		c.mu.RUnlock()
		// Each handler runs in its own goroutine so a slow handler cannot stall the IPC read loop.
		// All handlers in the library use internal mutexes, making concurrent invocation safe.
		params := m.Params
		for _, handler := range handlers {
			fn := handler.fn
			if fn != nil {
				c.dispatchWg.Add(1)
				go func() {
					defer c.dispatchWg.Done()
					fn(params)
				}()
			}
		}
	}
}

// Close signals that the connection is permanently broken.
// All goroutines blocked in SendRequest will receive err and unblock immediately.
// Subsequent SendRequest calls return err without sending.
// Also unblocks any WaitPlaywright callers that are still waiting.
func (c *Connection) Close(err error) {
	if err == nil {
		err = fmt.Errorf("connection closed")
	}
	c.mu.Lock()
	c.closed = true
	errMsg := Message{
		Error: &struct {
			Error MessageError `json:"error"`
		}{Error: MessageError{Name: "Error", Message: err.Error()}},
	}
	for id, ch := range c.callbacks {
		select {
		case ch <- errMsg:
		default:
		}
		delete(c.callbacks, id)
	}
	c.mu.Unlock()

	// Unblock WaitPlaywright if the Playwright root object was never received.
	// playwrightOnce guarantees close(playwrightReady) is called exactly once.
	c.playwrightOnce.Do(func() { close(c.playwrightReady) })
}

// SendRequest registers a pending callback for id, calls send(payload), then waits for
// Dispatch to deliver the response or ctx to be canceled. The callback entry is always
// removed from the map regardless of the outcome, preventing memory leaks from orphaned requests.
func (c *Connection) SendRequest(ctx context.Context, id int, payload []byte) (Message, error) {
	ch := make(chan Message, 1)
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return Message{}, fmt.Errorf("connection is closed")
	}
	send := c.transportSend
	c.callbacks[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.callbacks, id)
		c.mu.Unlock()
	}()

	if send == nil {
		return Message{}, fmt.Errorf("transport send function not initialized")
	}

	if err := send(payload); err != nil {
		return Message{}, err
	}

	select {
	case msg := <-ch:
		if msg.Error != nil {
			name, detail := msg.Error.Error.Name, msg.Error.Error.Message
			if name == "" && detail == "" {
				return Message{}, fmt.Errorf("playwright server error (no detail)")
			}
			return Message{}, fmt.Errorf("%s: %s", name, detail)
		}
		return msg, nil
	case <-ctx.Done():
		return Message{}, ctx.Err()
	}
}

// RegisterObject stores a protocol object in the connection registry under its GUID.
// Called automatically by Dispatch when a "__create__" event arrives; callers can also
// pre-register objects that are known before the driver sends a "__create__" message.
func (c *Connection) RegisterObject(guid string, owner *ObjectRef) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.objects[guid] = owner
}

// GetObject retrieves a registered protocol object by GUID.
// Returns (nil, false) if the object has not been created or has already been disposed.
func (c *Connection) GetObject(guid string) (*ObjectRef, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	owner, ok := c.objects[guid]
	return owner, ok
}

// NextID atomically generates and returns the next message ID for IPC requests.
func (c *Connection) NextID() int {
	return int(c.nextID.Add(1))
}

// WaitPlaywright blocks until the Playwright root object has been received from the driver
// (i.e., a "__create__" message with type "Playwright" arrives via Dispatch) or ctx is canceled.
// Returns the ObjectRef for the root object, which holds the initializer data used to
// construct BrowserType instances (Chromium, Firefox, WebKit).
func (c *Connection) WaitPlaywright(ctx context.Context) (*ObjectRef, error) {
	select {
	case <-c.playwrightReady:
		c.mu.RLock()
		obj := c.playwrightObj
		c.mu.RUnlock()
		if obj == nil {
			return nil, fmt.Errorf("connection closed before playwright was initialized")
		}
		return obj, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// OnEvent registers an event handler for a specific object and method.
func (c *Connection) OnEvent(guid, method string, handler func(json.RawMessage)) ListenerID {
	id := ListenerID(c.nextListenerID.Add(1))
	c.mu.Lock()
	defer c.mu.Unlock()
	key := eventKey{GUID: guid, Method: method}
	if _, exists := c.eventListeners[key]; !exists {
		// First handler for this key: record it in the secondary index.
		c.guidListeners[guid] = append(c.guidListeners[guid], key)
	}
	c.eventListeners[key] = append(c.eventListeners[key], eventHandler{id: id, fn: handler})
	return id
}

// OnEventOnce registers an event listener that executes strictly once and automatically unregisters itself.
// The returned ListenerID can be used to cancel the listener before it fires via OffEvent.
// After the handler has executed, calling OffEvent with the returned ID is a safe no-op.
func (c *Connection) OnEventOnce(guid, method string, handler func(json.RawMessage)) ListenerID {
	// Pre-allocate the ID before registering the handler so the closure captures
	// an immutable value, eliminating the data race where the dispatch goroutine
	// could read a zero onceID before the calling goroutine finishes assigning it.
	id := ListenerID(c.nextListenerID.Add(1))
	var once sync.Once

	wrapper := func(payload json.RawMessage) {
		once.Do(func() {
			c.OffEvent(guid, method, id)
			if handler != nil {
				handler(payload)
			}
		})
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	key := eventKey{GUID: guid, Method: method}
	if _, exists := c.eventListeners[key]; !exists {
		c.guidListeners[guid] = append(c.guidListeners[guid], key)
	}
	c.eventListeners[key] = append(c.eventListeners[key], eventHandler{id: id, fn: wrapper})
	return id
}

// OffEvent removes an event handler by its ListenerID.
func (c *Connection) OffEvent(guid, method string, id ListenerID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := eventKey{GUID: guid, Method: method}
	handlers := c.eventListeners[key]
	for i, h := range handlers {
		if h.id == id {
			copy(handlers[i:], handlers[i+1:])
			handlers[len(handlers)-1] = eventHandler{} // allow GC of removed fn
			updated := handlers[:len(handlers)-1]
			if len(updated) == 0 {
				delete(c.eventListeners, key)
				// Remove from secondary index when last handler for this key is gone.
				keys := c.guidListeners[guid]
				for j, k := range keys {
					if k == key {
						c.guidListeners[guid] = append(keys[:j], keys[j+1:]...)
						break
					}
				}
				if len(c.guidListeners[guid]) == 0 {
					delete(c.guidListeners, guid)
				}
			} else {
				c.eventListeners[key] = updated
			}
			break
		}
	}
}

// WaitDispatch blocks until all goroutines spawned by async event dispatch have completed.
// Primarily useful in tests that need to observe handler side-effects after Dispatch returns.
func (c *Connection) WaitDispatch() {
	c.dispatchWg.Wait()
}
