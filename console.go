package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// ConsoleMessage represents a message logged to the browser console.
type ConsoleMessage struct {
	page      *Page
	worker    *Worker
	msgType   string
	text      string
	timestamp float64
	location  string
	args      []*JSHandle
}

// Page returns the Page that emitted this console message, if any.
func (m *ConsoleMessage) Page() *Page { return m.page }

// Worker returns the Worker that emitted this console message, if any.
func (m *ConsoleMessage) Worker() *Worker { return m.worker }

// Type returns the log level: "log", "error", "warn", "info", "debug", etc.
func (m *ConsoleMessage) Type() string { return m.msgType }

// Text returns the concatenated text of all console arguments.
func (m *ConsoleMessage) Text() string { return m.text }

// Timestamp returns the timestamp of the message in milliseconds since the Unix epoch.
func (m *ConsoleMessage) Timestamp() float64 { return m.timestamp }

// Location returns the source location of the message as "url:line:column".
func (m *ConsoleMessage) Location() string { return m.location }

// Args returns JavaScript handles for the console message arguments.
func (m *ConsoleMessage) Args() []*JSHandle {
	return append([]*JSHandle(nil), m.args...)
}

// consoleLocation is the source location embedded in a console event.
type consoleLocation struct {
	URL          string `json:"url"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`
}

// consoleEventPayload represents the wire payload of a BrowserContext "console" event.
// The $mixin fields (type, text, args, timestamp, location) are inlined alongside the page reference.
type consoleEventPayload struct {
	Page      *protocol.Page   `json:"page,omitempty"`
	Worker    *protocol.Worker `json:"worker,omitempty"`
	Type      string           `json:"type"`
	Text      string           `json:"text"`
	Message   string           `json:"message"`
	Args      json.RawMessage  `json:"args"`
	Timestamp float64          `json:"timestamp"`
	Location  consoleLocation  `json:"location"`
}

func consoleMessageFrom(owner ChannelOwner, payload consoleEventPayload, page *Page, worker ...*Worker) *ConsoleMessage {
	text := payload.Text
	if text == "" && payload.Message != "" {
		text = payload.Message
	}
	if payload.Type == "" {
		payload.Type = "log"
	}
	location := fmt.Sprintf("%s:%d:%d", payload.Location.URL, payload.Location.LineNumber, payload.Location.ColumnNumber)
	scope := owner
	if payload.Worker != nil && payload.Worker.Guid != "" {
		scope = owner.child(payload.Worker.Guid)
	}
	var handles []protocol.JSHandle
	if len(payload.Args) > 0 {
		_ = json.Unmarshal(payload.Args, &handles)
	}
	args := make([]*JSHandle, 0, len(handles))
	for _, handle := range handles {
		if handle.Guid != "" {
			args = append(args, &JSHandle{owner: scope.child(handle.Guid)})
		}
	}
	var w *Worker
	if len(worker) > 0 && worker[0] != nil {
		w = worker[0]
	}
	return &ConsoleMessage{
		page:      page,
		worker:    w,
		msgType:   payload.Type,
		text:      text,
		timestamp: payload.Timestamp,
		location:  location,
		args:      args,
	}
}

type consoleHandlerEntry struct {
	page    *Page
	worker  *Worker
	handler func(*ConsoleMessage)
}

func (c *BrowserContext) registerConsoleHandler(page *Page, worker *Worker, handler func(*ConsoleMessage)) func() {
	if c == nil {
		return func() {}
	}
	c.mu.Lock()
	if c.consoleHandlers == nil {
		c.consoleHandlers = make(map[int]consoleHandlerEntry)
	}
	c.nextHandlerID++
	id := c.nextHandlerID
	c.consoleHandlers[id] = consoleHandlerEntry{page: page, worker: worker, handler: handler}
	c.mu.Unlock()

	c.ensureConsoleSubscription()

	var once sync.Once
	return func() {
		once.Do(func() {
			c.mu.Lock()
			delete(c.consoleHandlers, id)
			c.mu.Unlock()
			c.removeConsoleSubscription()
		})
	}
}

func (c *BrowserContext) ensureConsoleSubscription() {
	c.retainSubscription(context.Background(), "console")
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.consoleListenerID != 0 {
		return
	}
	c.consoleListenerID = c.owner.conn.OnEvent(c.owner.guid, "console", func(params json.RawMessage) {
		var payload consoleEventPayload
		if err := json.Unmarshal(params, &payload); err != nil {
			return
		}
		var p *Page
		if payload.Page != nil && payload.Page.Guid != "" {
			p = c.pageByGUID(payload.Page.Guid)
		}
		var w *Worker
		if payload.Worker != nil && payload.Worker.Guid != "" {
			if p != nil {
				w = p.workerByGUID(payload.Worker.Guid)
			}
			if w == nil {
				w = c.workerByGUID(payload.Worker.Guid)
			}
		}
		msg := consoleMessageFrom(c.owner, payload, p, w)

		c.mu.RLock()
		var matching []func(*ConsoleMessage)
		for _, entry := range c.consoleHandlers {
			if entry.worker != nil {
				if payload.Worker != nil && (entry.worker == w || entry.worker.owner.guid == payload.Worker.Guid) {
					matching = append(matching, entry.handler)
				}
			} else if entry.page != nil {
				if payload.Page != nil && (entry.page == p || entry.page.owner.guid == payload.Page.Guid) {
					matching = append(matching, entry.handler)
				}
			} else {
				matching = append(matching, entry.handler)
			}
		}
		c.mu.RUnlock()

		for _, fn := range matching {
			func() {
				defer func() {
					if r := recover(); r != nil {
						slog.Default().Error("connection: panic in event handler", "guid", c.owner.guid, "method", "console", "panic", r)
					}
				}()
				fn(msg)
			}()
		}
	})
}

func (c *BrowserContext) removeConsoleSubscription() {
	c.mu.Lock()
	if len(c.consoleHandlers) > 0 || c.consoleListenerID == 0 {
		c.mu.Unlock()
		return
	}
	id := c.consoleListenerID
	c.consoleListenerID = 0
	c.mu.Unlock()

	c.owner.conn.OffEvent(c.owner.guid, "console", id)
	c.releaseSubscription("console")
}

// OnConsole registers a handler for all console messages in this context,
// regardless of which page emitted them. Unlike Page.OnConsole, this handler
// fires for messages from popups and other pages within the same context.
// The returned function cancels the listener.
func (c *BrowserContext) OnConsole(handler func(*ConsoleMessage)) func() {
	return c.registerConsoleHandler(nil, nil, handler)
}

// OnConsole registers a handler for console messages emitted by this page.
// Handlers are invoked sequentially, in the order messages were emitted.
// Console events are emitted on the BrowserContext channel; events from other
// pages in the same context are filtered out.
// The returned function cancels the listener when called.
func (p *Page) OnConsole(handler func(*ConsoleMessage)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	return p.browserContext.registerConsoleHandler(p, nil, handler)
}

// ConsoleMessagesFilter specifies the filter for console messages.
type ConsoleMessagesFilter = protocol.ConsoleMessagesFilter

const (
	ConsoleMessagesFilterAll             ConsoleMessagesFilter = protocol.ConsoleMessagesFilterAll
	ConsoleMessagesFilterSinceNavigation ConsoleMessagesFilter = protocol.ConsoleMessagesFilterSincenavigation
)

// ConsoleMessagesOptions specifies options for Page.ConsoleMessages.
type ConsoleMessagesOptions struct {
	Filter *ConsoleMessagesFilter
}

// SetFilter sets the filter for console messages: "all" or "since-navigation".
func (o *ConsoleMessagesOptions) SetFilter(f ConsoleMessagesFilter) *ConsoleMessagesOptions {
	o.Filter = &f
	return o
}

// ConsoleMessages returns all console messages matching the options.
func (p *Page) ConsoleMessages(ctx context.Context, opts ...*ConsoleMessagesOptions) ([]*ConsoleMessage, error) {
	req := protocol.PageConsoleMessagesRequest{}
	if len(opts) > 0 && opts[0] != nil && opts[0].Filter != nil {
		req.Filter = opts[0].Filter
	}
	result, err := p.owner.SendMessageRequest(ctx, "consoleMessages", req)
	if err != nil {
		return nil, fmt.Errorf("page.consoleMessages failed: %w", err)
	}
	var resp struct {
		Messages []consoleEventPayload `json:"messages"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse page.consoleMessages response: %w", err)
	}
	messages := make([]*ConsoleMessage, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		messages = append(messages, consoleMessageFrom(p.owner, m, p))
	}
	return messages, nil
}

// ClearConsoleMessages clears the accumulated console messages for this page.
func (p *Page) ClearConsoleMessages(ctx context.Context) error {
	_, err := p.owner.SendMessageRequest(ctx, "clearConsoleMessages", struct{}{})
	if err != nil {
		return fmt.Errorf("page.clearConsoleMessages failed: %w", err)
	}
	return nil
}

// WaitForConsoleMessage waits for a console message matching predicate, running trigger if provided.
func (p *Page) WaitForConsoleMessage(ctx context.Context, args ...any) (*ConsoleMessage, error) {
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
			return nil, fmt.Errorf("page.WaitForConsoleMessage: unsupported argument type %T", arg)
		}
	}

	deadlineTimeout := time.Duration(p.owner.effectiveTimeout("waitForConsoleMessage")) * time.Millisecond
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

	return waitForOne(ctx, "page.WaitForConsoleMessage", func(send func(*ConsoleMessage)) func() {
		return p.OnConsole(func(msg *ConsoleMessage) {
			if predicate != nil && !predicate(msg) {
				return
			}
			send(msg)
		})
	}, trigger)
}
