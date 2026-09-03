package playwright

import (
	"context"
	"encoding/json"
	"fmt"
)

// Worker represents a Web Worker associated with a page.
type Worker struct {
	owner ChannelOwner
	url   string
}

// workerInitializer is the wire format of the Worker channel object initializer.
type workerInitializer struct {
	URL string `json:"url"`
}

// URL returns the URL of the worker script.
func (w *Worker) URL() string { return w.url }

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
		w := &Worker{owner: p.owner.child(event.Worker.Guid), url: init.URL}
		p.mu.Lock()
		p.workers = append(p.workers, w)
		listeners := make([]func(*Worker), 0, len(p.workerListenersByID))
		for _, fn := range p.workerListenersByID {
			listeners = append(listeners, fn)
		}
		p.mu.Unlock()
		p.owner.conn.OnEventOnce(event.Worker.Guid, "close", func(_ json.RawMessage) {
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
			go fn(w)
		}
	})
	p.mu.Lock()
	p.workerListenerID = id
	p.workerListenerActive = true
	p.mu.Unlock()
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
