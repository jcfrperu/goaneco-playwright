package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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

type dialogHandlerEntry struct {
	page    *Page
	handler func(*Dialog)
}

// Dialog represents a browser dialog (alert, confirm, prompt, or beforeunload).
type Dialog struct {
	owner        ChannelOwner
	dialogType   string
	message      string
	defaultValue string
	page         *Page
	context      *BrowserContext
	closeOnce    sync.Once
}

// Type returns the dialog type: "alert", "confirm", "prompt", or "beforeunload".
func (d *Dialog) Type() string { return d.dialogType }

// Message returns the message shown in the dialog.
func (d *Dialog) Message() string { return d.message }

// DefaultValue returns the default value for prompt dialogs (empty for others).
func (d *Dialog) DefaultValue() string { return d.defaultValue }

// Page returns the page that owns this dialog, or nil if the page is unknown.
func (d *Dialog) Page() *Page { return d.page }

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
	d.notifyClosed()
	return nil
}

// Dismiss dismisses the dialog (equivalent to clicking Cancel).
func (d *Dialog) Dismiss(ctx context.Context) error {
	_, err := d.owner.SendMessageRequest(ctx, "dismiss", protocol.DialogDismissRequest{})
	if err != nil {
		return fmt.Errorf("dialog.dismiss failed: %w", err)
	}
	d.notifyClosed()
	return nil
}

func (d *Dialog) notifyClosed() {
	d.closeOnce.Do(func() {
		if d.context != nil {
			d.context.dispatchDialogClosed(d)
		}
	})
}

func (c *BrowserContext) dialogFromGUID(guid string) *Dialog {
	c.mu.RLock()
	existing := c.dialogsByGUID[guid]
	c.mu.RUnlock()
	if existing != nil {
		return existing
	}

	var init dialogInitializer
	if raw := c.owner.Initializer(guid); len(raw) > 0 {
		_ = json.Unmarshal(raw, &init)
	}
	dialog := &Dialog{
		owner:        c.owner.child(guid),
		dialogType:   init.Type,
		message:      init.Message,
		defaultValue: init.DefaultValue,
		context:      c,
	}
	if init.Page.Guid != "" {
		dialog.page = c.pageByGUID(init.Page.Guid)
	}

	c.mu.Lock()
	if c.dialogsByGUID == nil {
		c.dialogsByGUID = make(map[string]*Dialog)
	}
	if existing = c.dialogsByGUID[guid]; existing == nil {
		c.dialogsByGUID[guid] = dialog
		existing = dialog
	}
	c.mu.Unlock()
	return existing
}

func (c *BrowserContext) ensureDialogSubscription() {
	c.mu.Lock()
	if c.dialogListenerID != 0 {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()

	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	defer subCancel()
	c.retainSubscription(subCtx, "dialog")

	c.mu.Lock()
	if c.dialogListenerID == 0 {
		c.dialogListenerID = c.owner.conn.OnEvent(c.owner.guid, "dialog", func(params json.RawMessage) {
			var event protocol.BrowserContextDialogEvent
			if err := json.Unmarshal(params, &event); err != nil || event.Dialog.Guid == "" {
				return
			}
			dialog := c.dialogFromGUID(event.Dialog.Guid)
			c.dispatchDialog(dialog)
		})
	}
	c.mu.Unlock()
}

func (c *BrowserContext) removeDialogSubscription() {
	c.mu.Lock()
	if len(c.dialogHandlers) > 0 || len(c.dialogClosedHandlers) > 0 || c.dialogListenerID == 0 {
		c.mu.Unlock()
		return
	}
	id := c.dialogListenerID
	c.dialogListenerID = 0
	c.mu.Unlock()

	c.owner.conn.OffEvent(c.owner.guid, "dialog", id)
	c.releaseSubscription("dialog")
}

func (c *BrowserContext) registerDialogHandler(page *Page, handler func(*Dialog)) func() {
	c.mu.Lock()
	if c.dialogHandlers == nil {
		c.dialogHandlers = make(map[int]dialogHandlerEntry)
	}
	c.nextHandlerID++
	id := c.nextHandlerID
	c.dialogHandlers[id] = dialogHandlerEntry{page: page, handler: handler}
	c.mu.Unlock()

	c.ensureDialogSubscription()

	var once sync.Once
	return func() {
		once.Do(func() {
			c.mu.Lock()
			delete(c.dialogHandlers, id)
			c.mu.Unlock()
			c.removeDialogSubscription()
		})
	}
}

func (c *BrowserContext) registerDialogClosedHandler(page *Page, handler func(*Dialog)) func() {
	c.mu.Lock()
	if c.dialogClosedHandlers == nil {
		c.dialogClosedHandlers = make(map[int]dialogHandlerEntry)
	}
	c.nextHandlerID++
	id := c.nextHandlerID
	c.dialogClosedHandlers[id] = dialogHandlerEntry{page: page, handler: handler}
	c.mu.Unlock()

	c.ensureDialogSubscription()

	var once sync.Once
	return func() {
		once.Do(func() {
			c.mu.Lock()
			delete(c.dialogClosedHandlers, id)
			c.mu.Unlock()
			c.removeDialogSubscription()
		})
	}
}

func (c *BrowserContext) dispatchDialog(dialog *Dialog) {
	c.mu.RLock()
	var matching []func(*Dialog)
	for _, entry := range c.dialogHandlers {
		if entry.page == nil || entry.page == dialog.page {
			matching = append(matching, entry.handler)
		}
	}
	c.mu.RUnlock()

	if len(matching) > 0 {
		for _, h := range matching {
			h := h
			safeGo(func() { h(dialog) })
		}
	} else {
		// Auto-dismiss dialogs when no listeners are attached, exactly matching Java BrowserContextImpl
		safeGo(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if dialog.dialogType == "beforeunload" {
				_ = dialog.Accept(ctx)
			} else {
				_ = dialog.Dismiss(ctx)
			}
		})
	}
}

func (c *BrowserContext) dispatchDialogClosed(dialog *Dialog) {
	c.mu.RLock()
	var matching []func(*Dialog)
	for _, entry := range c.dialogClosedHandlers {
		if entry.page == nil || entry.page == dialog.page {
			matching = append(matching, entry.handler)
		}
	}
	c.mu.RUnlock()

	for _, h := range matching {
		h := h
		safeGo(func() { h(dialog) })
	}
}

// OnDialog registers a handler for browser dialogs on this page.
// The handler is called in a goroutine when the browser emits a "dialog" event.
// The handler MUST call dialog.Accept or dialog.Dismiss; otherwise the page hangs.
// Dialogs from other pages in the same context are ignored.
// The returned function cancels the listener when called.
func (p *Page) OnDialog(handler func(*Dialog)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	return p.browserContext.registerDialogHandler(p, handler)
}

// OnDialog registers a handler for dialogs from any page in this context.
// The handler must accept or dismiss each dialog so the owning page can continue.
// The returned function unregisters the handler.
func (c *BrowserContext) OnDialog(handler func(*Dialog)) func() {
	return c.registerDialogHandler(nil, handler)
}

// OnDialogClosed registers a handler invoked when any dialog in the context is accepted or dismissed.
// The returned function unregisters the handler.
func (c *BrowserContext) OnDialogClosed(handler func(*Dialog)) func() {
	return c.registerDialogClosedHandler(nil, handler)
}

// OnDialogClosed registers a handler invoked when a dialog on this page is accepted or dismissed.
// The returned function unregisters the handler.
func (p *Page) OnDialogClosed(handler func(*Dialog)) func() {
	if p.browserContext == nil {
		return func() {}
	}
	return p.browserContext.registerDialogClosedHandler(p, handler)
}
