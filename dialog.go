package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// dialogInitializer contains the initializer data for a Dialog channel object.
type dialogInitializer struct {
	Type         string `json:"type"`
	Message      string `json:"message"`
	DefaultValue string `json:"defaultValue"`
	Page         struct {
		Guid string `json:"guid"`
	} `json:"page"`
}

// Dialog represents a browser dialog (alert, confirm, prompt, or beforeunload).
type Dialog struct {
	owner        ChannelOwner
	dialogType   string
	message      string
	defaultValue string
}

// Type returns the dialog type: "alert", "confirm", "prompt", or "beforeunload".
func (d *Dialog) Type() string { return d.dialogType }

// Message returns the message shown in the dialog.
func (d *Dialog) Message() string { return d.message }

// DefaultValue returns the default value for prompt dialogs (empty for others).
func (d *Dialog) DefaultValue() string { return d.defaultValue }

// Accept accepts the dialog. For prompt dialogs, promptText sets the input value.
func (d *Dialog) Accept(ctx context.Context, promptText ...string) error {
	req := protocol.DialogAcceptRequest{}
	if len(promptText) > 0 {
		req.PromptText = &promptText[0]
	}
	_, err := d.owner.SendMessageRequest(ctx, "accept", req)
	if err != nil {
		return fmt.Errorf("dialog.accept failed: %w", err)
	}
	return nil
}

// Dismiss dismisses the dialog (equivalent to clicking Cancel).
func (d *Dialog) Dismiss(ctx context.Context) error {
	_, err := d.owner.SendMessageRequest(ctx, "dismiss", protocol.DialogDismissRequest{})
	if err != nil {
		return fmt.Errorf("dialog.dismiss failed: %w", err)
	}
	return nil
}

// OnDialog registers a handler for browser dialogs on this page.
// The handler is called in a goroutine when the browser emits a "dialog" event.
// The handler MUST call dialog.Accept or dialog.Dismiss; otherwise the page hangs.
// Note: the BrowserContext dialog event does not carry a page reference in the wire
// protocol, so the handler fires for dialogs from ANY page in the same context.
// The returned function cancels the listener when called.
func (p *Page) OnDialog(handler func(*Dialog)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	if _, err := p.browserContext.owner.SendMessageRequest(subCtx, "updateSubscription", protocol.BrowserContextUpdateSubscriptionRequest{
		Event:   "dialog",
		Enabled: true,
	}); err != nil {
		slog.Default().Warn("page.OnDialog: failed to enable dialog subscription", "error", err)
	}
	subCancel()

	handleDialogEvent := func(params json.RawMessage) {
		var event protocol.BrowserContextDialogEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Dialog.Guid == "" {
			return
		}
		dialogOwner := p.owner.child(event.Dialog.Guid)
		var init dialogInitializer
		if raw := p.owner.Initializer(event.Dialog.Guid); len(raw) > 0 {
			_ = json.Unmarshal(raw, &init) // best-effort; zero-value fallback if server data is malformed
		}
		if init.Page.Guid != "" && init.Page.Guid != p.owner.guid {
			return // dialog belongs to a different page
		}

		d := &Dialog{
			owner:        dialogOwner,
			dialogType:   init.Type,
			message:      init.Message,
			defaultValue: init.DefaultValue,
		}
		go handler(d)
	}

	id := p.owner.conn.OnEvent(p.browserContext.owner.guid, "dialog", handleDialogEvent)
	return func() {
		p.owner.conn.OffEvent(p.browserContext.owner.guid, "dialog", id)
	}
}
