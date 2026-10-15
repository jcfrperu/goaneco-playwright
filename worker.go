package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Worker represents a Web Worker associated with a page.
type Worker struct {
	page  *Page
	owner ChannelOwner
	url   string

	mu            sync.Mutex
	nextListener  uint64
	closeCh       chan struct{}
	closeOnce     sync.Once
	closed        bool
	closeHandlers map[uint64]func(*Worker)
}

// Page returns the Page that this worker belongs to, if any.
func (w *Worker) Page() *Page { return w.page }

// workerInitializer is the wire format of the Worker channel object initializer.
type workerInitializer struct {
	URL string `json:"url"`
}

// URL returns the URL of the worker script.
func (w *Worker) URL() string { return w.url }

// OnClose registers a handler called when this worker is closed.
func (w *Worker) OnClose(handler func(*Worker)) func() {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		handler(w)
		return func() {}
	}
	w.nextListener++
	id := w.nextListener
	if w.closeHandlers == nil {
		w.closeHandlers = make(map[uint64]func(*Worker))
	}
	w.closeHandlers[id] = handler
	w.mu.Unlock()
	return func() {
		w.mu.Lock()
		delete(w.closeHandlers, id)
		w.mu.Unlock()
	}
}

// WaitForClose waits until this worker closes, running trigger first when provided.
func (w *Worker) WaitForClose(ctx context.Context, trigger func() error) (*Worker, error) {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return w, nil
	}
	ch := w.closeCh
	if ch == nil {
		ch = make(chan struct{})
		w.closeCh = ch
	}
	w.mu.Unlock()

	if trigger != nil {
		if err := trigger(); err != nil {
			return nil, fmt.Errorf("worker.waitForClose: trigger failed: %w", err)
		}
	}
	select {
	case <-ch:
		return w, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("worker.waitForClose: %w", ctx.Err())
	}
}

func (w *Worker) markClosed() {
	w.closeOnce.Do(func() {
		w.mu.Lock()
		w.closed = true
		if w.closeCh == nil {
			w.closeCh = make(chan struct{})
		}
		close(w.closeCh)
		handlers := make([]func(*Worker), 0, len(w.closeHandlers))
		for _, handler := range w.closeHandlers {
			handlers = append(handlers, handler)
		}
		w.mu.Unlock()
		for _, handler := range handlers {
			handler(w)
		}
	})
}

// Evaluate executes a JavaScript expression in the worker's context.
func (w *Worker) Evaluate(ctx context.Context, expression string, arg ...any) (any, error) {
	var input any
	if len(arg) > 0 {
		input = arg[0]
	}
	result, err := w.owner.SendMessageRequest(ctx, "evaluateExpression", jsHandleEvalParams{
		Expression: expression,
		Arg:        serializeArgument(input),
	})
	if err != nil {
		return nil, fmt.Errorf("worker.evaluate failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse worker.evaluate response: %w", err)
	}
	return deserializeValue(resp.Value)
}

// EvaluateHandle executes a JavaScript expression in the worker and returns the result as a JSHandle.
func (w *Worker) EvaluateHandle(ctx context.Context, expression string, arg ...any) (*JSHandle, error) {
	var input any
	if len(arg) > 0 {
		input = arg[0]
	}
	result, err := w.owner.SendMessageRequest(ctx, "evaluateExpressionHandle", jsHandleEvalParams{
		Expression: expression,
		Arg:        serializeArgument(input),
	})
	if err != nil {
		return nil, fmt.Errorf("worker.evaluateHandle failed: %w", err)
	}
	var resp struct {
		Handle struct {
			Guid string `json:"guid"`
		} `json:"handle"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse worker.evaluateHandle response: %w", err)
	}
	return &JSHandle{owner: w.owner.child(resp.Handle.Guid)}, nil
}

// Workers returns the currently active web workers associated with this page.
func (p *Page) Workers() []*Worker {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]*Worker, len(p.workers))
	copy(result, p.workers)
	return result
}

// subscribeToWorkers sets up a single centralized listener for worker events on this page.
// Must be called once at page construction.
func (p *Page) subscribeToWorkers() {
	id := p.owner.conn.OnEvent(p.owner.guid, "worker", func(params json.RawMessage) {
		var event struct {
			Worker struct {
				Guid string `json:"guid"`
			} `json:"worker"`
		}
		if err := json.Unmarshal(params, &event); err != nil || event.Worker.Guid == "" {
			return
		}
		var init workerInitializer
		if raw := p.owner.Initializer(event.Worker.Guid); len(raw) > 0 {
			_ = json.Unmarshal(raw, &init)
		}
		w := &Worker{page: p, owner: p.owner.child(event.Worker.Guid), url: init.URL, closeCh: make(chan struct{})}
		p.mu.Lock()
		p.workers = append(p.workers, w)
		listeners := make([]func(*Worker), 0, len(p.workerListenersByID))
		for _, fn := range p.workerListenersByID {
			listeners = append(listeners, fn)
		}
		p.mu.Unlock()
		p.owner.conn.OnEventOnce(event.Worker.Guid, "close", func(_ json.RawMessage) {
			w.markClosed()
			p.mu.Lock()
			for i, pw := range p.workers {
				if pw == w {
					copy(p.workers[i:], p.workers[i+1:])
					p.workers[len(p.workers)-1] = nil
					p.workers = p.workers[:len(p.workers)-1]
					break
				}
			}
			p.mu.Unlock()
		})
		for _, fn := range listeners {
			fn(w)
		}
	})
	p.mu.Lock()
	p.workerListenerID = id
	p.workerListenerActive = true
	p.mu.Unlock()
}

// workerByGUID returns the worker with the given GUID associated with this page, or nil.
func (p *Page) workerByGUID(guid string) *Worker {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, w := range p.workers {
		if w.owner.guid == guid {
			return w
		}
	}
	return nil
}

// OnWorker registers a handler that is called when a new worker is created in this page.
// The handler receives the new Worker and is called in a goroutine.
// The returned function cancels the listener.
func (p *Page) OnWorker(handler func(*Worker)) func() {
	p.mu.Lock()
	p.workerNextID++
	id := p.workerNextID
	if p.workerListenersByID == nil {
		p.workerListenersByID = make(map[int]func(*Worker))
	}
	p.workerListenersByID[id] = handler
	p.mu.Unlock()
	return func() {
		p.mu.Lock()
		delete(p.workerListenersByID, id)
		p.mu.Unlock()
	}
}

// WaitForWorker waits for a worker to be created on this page, running trigger if provided.
func (p *Page) WaitForWorker(ctx context.Context, args ...any) (*Worker, error) {
	var predicate func(*Worker) bool
	var trigger func() error
	var explicitTimeout time.Duration
	var sawTimeout bool
	for _, arg := range args {
		switch v := arg.(type) {
		case nil:
		case func(*Worker) bool:
			predicate = v
		case func() error:
			trigger = v
		case func():
			trigger = func() error { v(); return nil }
		case time.Duration:
			explicitTimeout = v
			sawTimeout = true
		case int:
			explicitTimeout = time.Duration(v) * time.Millisecond
			sawTimeout = true
		case float64:
			explicitTimeout = time.Duration(v) * time.Millisecond
			sawTimeout = true
		default:
			return nil, fmt.Errorf("page.WaitForWorker: unsupported argument type %T", arg)
		}
	}

	deadlineTimeout := time.Duration(p.owner.effectiveTimeout("waitForWorker")) * time.Millisecond
	if sawTimeout {
		deadlineTimeout = explicitTimeout
	}
	if deadlineTimeout > 0 {
		if d, ok := ctx.Deadline(); ok {
			remaining := time.Until(d)
			if remaining > 0 && remaining < deadlineTimeout {
				deadlineTimeout = remaining
			}
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, deadlineTimeout)
		defer cancel()
	}

	return waitForOne(ctx, "page.WaitForWorker", func(send func(*Worker)) func() {
		return p.OnWorker(func(w *Worker) {
			if predicate != nil && !predicate(w) {
				return
			}
			send(w)
		})
	}, trigger)
}

// OnConsole registers a handler for console messages emitted by this worker.
func (w *Worker) OnConsole(handler func(*ConsoleMessage)) func() {
	if w.page == nil || w.page.browserContext == nil {
		return func() {}
	}
	return w.page.browserContext.registerConsoleHandler(nil, w, handler)
}

// WaitForConsoleMessage waits for a console message from this worker matching predicate, running trigger if provided.
func (w *Worker) WaitForConsoleMessage(ctx context.Context, args ...any) (*ConsoleMessage, error) {
	var predicate func(*ConsoleMessage) bool
	var trigger func() error
	var explicitTimeout time.Duration
	var sawTimeout bool
	for _, arg := range args {
		switch v := arg.(type) {
		case nil:
		case func(*ConsoleMessage) bool:
			predicate = v
		case func() error:
			trigger = v
		case func():
			trigger = func() error { v(); return nil }
		case time.Duration:
			explicitTimeout = v
			sawTimeout = true
		case int:
			explicitTimeout = time.Duration(v) * time.Millisecond
			sawTimeout = true
		case float64:
			explicitTimeout = time.Duration(v) * time.Millisecond
			sawTimeout = true
		default:
			return nil, fmt.Errorf("worker.WaitForConsoleMessage: unsupported argument type %T", arg)
		}
	}

	deadlineTimeout := time.Duration(w.owner.effectiveTimeout("waitForConsoleMessage")) * time.Millisecond
	if sawTimeout {
		deadlineTimeout = explicitTimeout
	}
	if deadlineTimeout > 0 {
		if d, ok := ctx.Deadline(); ok {
			remaining := time.Until(d)
			if remaining > 0 && remaining < deadlineTimeout {
				deadlineTimeout = remaining
			}
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, deadlineTimeout)
		defer cancel()
	}

	return waitForOne(ctx, "worker.WaitForConsoleMessage", func(send func(*ConsoleMessage)) func() {
		return w.OnConsole(func(msg *ConsoleMessage) {
			if predicate != nil && !predicate(msg) {
				return
			}
			send(msg)
		})
	}, trigger)
}
