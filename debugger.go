package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

// DebuggerLocation identifies a source position (file + optional line/column).
type DebuggerLocation struct {
	File   string `json:"file"`
	Line   *int   `json:"line,omitempty"`
	Column *int   `json:"column,omitempty"`
}

// DebuggerPausedDetails contains information about why execution is paused.
type DebuggerPausedDetails struct {
	Location *DebuggerLocation `json:"location,omitempty"`
	Title    string            `json:"title"`
	Stack    *string           `json:"stack,omitempty"`
}

// Debugger provides step-through debugging for a BrowserContext.
// Obtain via BrowserContext.Debugger().
type Debugger struct {
	owner ChannelOwner
}

// RequestPause requests that execution pauses at the next Playwright action.
func (d *Debugger) RequestPause(ctx context.Context) error {
	_, err := d.owner.SendMessageRequest(ctx, "requestPause", struct{}{})
	if err != nil {
		return fmt.Errorf("debugger.requestPause failed: %w", err)
	}
	return nil
}

// Resume resumes execution after a pause.
func (d *Debugger) Resume(ctx context.Context) error {
	_, err := d.owner.SendMessageRequest(ctx, "resume", struct{}{})
	if err != nil {
		return fmt.Errorf("debugger.resume failed: %w", err)
	}
	return nil
}

// Next advances one step when paused.
func (d *Debugger) Next(ctx context.Context) error {
	_, err := d.owner.SendMessageRequest(ctx, "next", struct{}{})
	if err != nil {
		return fmt.Errorf("debugger.next failed: %w", err)
	}
	return nil
}

// RunTo resumes execution until the given source location is reached.
func (d *Debugger) RunTo(ctx context.Context, location DebuggerLocation) error {
	_, err := d.owner.SendMessageRequest(ctx, "runTo", map[string]any{
		"location": location,
	})
	if err != nil {
		return fmt.Errorf("debugger.runTo failed: %w", err)
	}
	return nil
}

// PausedDetails returns the current paused state details, or nil if not paused.
// This is a client-side read of the last received pausedStateChanged event.
// Use OnPausedStateChanged to track state changes asynchronously.
func (d *Debugger) PausedDetails(ctx context.Context) (*DebuggerPausedDetails, error) {
	result, err := d.owner.SendMessageRequest(ctx, "pausedDetails", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("debugger.pausedDetails failed: %w", err)
	}
	var resp struct {
		PausedDetails *DebuggerPausedDetails `json:"pausedDetails,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse pausedDetails response: %w", err)
	}
	return resp.PausedDetails, nil
}

// OnPausedStateChanged registers a handler called each time the debugger pause state changes.
// The handler receives the new paused details, or nil when execution is resumed.
// The returned function cancels the listener.
func (d *Debugger) OnPausedStateChanged(handler func(*DebuggerPausedDetails)) func() {
	id := d.owner.conn.OnEvent(d.owner.guid, "pausedStateChanged", func(params json.RawMessage) {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Default().Error("debugger pausedStateChanged handler panic", "panic", r)
				}
			}()
			var event struct {
				PausedDetails *DebuggerPausedDetails `json:"pausedDetails,omitempty"`
			}
			if err := json.Unmarshal(params, &event); err != nil {
				slog.Default().Error("debugger: failed to parse pausedStateChanged event", "err", err)
				return
			}
			handler(event.PausedDetails)
		}()
	})
	return func() { d.owner.conn.OffEvent(d.owner.guid, "pausedStateChanged", id) }
}

// Debugger returns the Debugger object for this BrowserContext.
// The Debugger GUID comes from the BrowserContext initializer.
func (c *BrowserContext) Debugger() *Debugger {
	var init struct {
		Debugger struct {
			Guid string `json:"guid"`
		} `json:"debugger"`
	}
	if raw := c.initializer; len(raw) > 0 {
		_ = json.Unmarshal(raw, &init)
	}
	if init.Debugger.Guid == "" {
		return nil
	}
	return &Debugger{owner: c.owner.child(init.Debugger.Guid)}
}
