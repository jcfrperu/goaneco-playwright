package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// ConsoleMessage represents a message logged to the browser console.
type ConsoleMessage struct {
	msgType   string
	text      string
	timestamp float64
	location  string
}

// Type returns the log level: "log", "error", "warn", "info", "debug", etc.
func (m *ConsoleMessage) Type() string { return m.msgType }

// Text returns the concatenated text of all console arguments.
func (m *ConsoleMessage) Text() string { return m.text }

// Timestamp returns the timestamp of the message in milliseconds since the Unix epoch.
func (m *ConsoleMessage) Timestamp() float64 { return m.timestamp }

// Location returns the source location of the message as "url:line:column".
func (m *ConsoleMessage) Location() string { return m.location }

// consoleLocation is the source location embedded in a console event.
type consoleLocation struct {
	URL          string `json:"url"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`
}

// consoleEventPayload represents the wire payload of a BrowserContext "console" event.
// The $mixin fields (type, text, args, timestamp, location) are inlined alongside the page reference.
type consoleEventPayload struct {
	Page      *protocol.Page  `json:"page,omitempty"`
	Type      string          `json:"type"`
	Text      string          `json:"text"`
	Message   string          `json:"message"`
	Args      json.RawMessage `json:"args"`
	Timestamp float64         `json:"timestamp"`
	Location  consoleLocation `json:"location"`
}

// OnConsole registers a handler for all console messages in this context,
// regardless of which page emitted them. Unlike Page.OnConsole, this handler
// fires for messages from popups and other pages within the same context.
// The returned function cancels the listener.
func (c *BrowserContext) OnConsole(handler func(*ConsoleMessage)) func() {
	subCtx, subCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := c.owner.SendMessageRequest(subCtx, "updateSubscription", protocol.BrowserContextUpdateSubscriptionRequest{
		Event:   "console",
		Enabled: true,
	}); err != nil {
		slog.Default().Warn("browserContext.OnConsole: failed to enable console subscription", "error", err)
	}
	subCancel()

	processEvent := func(params json.RawMessage) {
		var payload consoleEventPayload
		if err := json.Unmarshal(params, &payload); err != nil {
			return
		}
		text := payload.Text
		if text == "" && payload.Message != "" {
			text = payload.Message
		}
		if payload.Type == "" {
			payload.Type = "log"
		}
		loc := fmt.Sprintf("%s:%d:%d", payload.Location.URL, payload.Location.LineNumber, payload.Location.ColumnNumber)
		go handler(&ConsoleMessage{msgType: payload.Type, text: text, timestamp: payload.Timestamp, location: loc})
	}

	id := c.owner.conn.OnEvent(c.owner.guid, "console", processEvent)
	return func() { c.owner.conn.OffEvent(c.owner.guid, "console", id) }
}

// OnConsole registers a handler for console messages emitted by this page.
// The handler is called from a goroutine for each message.
// Console events are emitted on the BrowserContext channel; events from other
// pages in the same context are filtered out.
// The returned function cancels the listener when called.
func (p *Page) OnConsole(handler func(*ConsoleMessage)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	subCtx, subCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := p.browserContext.owner.SendMessageRequest(subCtx, "updateSubscription", protocol.BrowserContextUpdateSubscriptionRequest{
		Event:   "console",
		Enabled: true,
	}); err != nil {
		slog.Default().Warn("page.OnConsole: failed to enable console subscription", "error", err)
	}
	subCancel()

	processEvent := func(params json.RawMessage) {
		var payload consoleEventPayload
		if err := json.Unmarshal(params, &payload); err != nil {
			return
		}
		if payload.Page != nil && payload.Page.Guid != p.owner.guid {
			return
		}
		text := payload.Text
		if text == "" && payload.Message != "" {
			text = payload.Message
		}
		if payload.Type == "" {
			payload.Type = "log"
		}
		loc := fmt.Sprintf("%s:%d:%d", payload.Location.URL, payload.Location.LineNumber, payload.Location.ColumnNumber)
		go handler(&ConsoleMessage{msgType: payload.Type, text: text, timestamp: payload.Timestamp, location: loc})
	}

	id := p.owner.conn.OnEvent(p.browserContext.owner.guid, "console", processEvent)
	return func() {
		p.owner.conn.OffEvent(p.browserContext.owner.guid, "console", id)
	}
}
