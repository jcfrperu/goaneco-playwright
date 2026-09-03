package playwright

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// FileChooser represents a browser file chooser dialog opened when an <input type="file">
// element is activated. Use SetFiles to supply the chosen files programmatically.
type FileChooser struct {
	page       *Page
	element    ChannelOwner
	isMultiple bool
}

// Page returns the page that owns this file chooser.
func (f *FileChooser) Page() *Page { return f.page }

// IsMultiple reports whether the file chooser accepts multiple files.
func (f *FileChooser) IsMultiple() bool { return f.isMultiple }

// SetFiles sets the chosen files by their local filesystem paths.
func (f *FileChooser) SetFiles(ctx context.Context, paths []string) error {
	req := map[string]any{
		"localPaths": paths,
		"timeout":    30000.0,
	}
	_, err := f.element.SendMessageRequest(ctx, "setInputFiles", req)
	if err != nil {
		return fmt.Errorf("fileChooser.setFiles failed: %w", err)
	}
	return nil
}

// SetFilesPayload sets the chosen files from in-memory payloads.
func (f *FileChooser) SetFilesPayload(ctx context.Context, payloads []FilePayload) error {
	wirePayloads := make([]map[string]any, len(payloads))
	for i, p := range payloads {
		wirePayloads[i] = map[string]any{
			"name":     p.Name,
			"mimeType": p.MimeType,
			"buffer":   base64.StdEncoding.EncodeToString(p.Buffer),
		}
	}
	req := map[string]any{
		"payloads": wirePayloads,
		"timeout":  30000.0,
	}
	_, err := f.element.SendMessageRequest(ctx, "setInputFiles", req)
	if err != nil {
		return fmt.Errorf("fileChooser.setFilesPayload failed: %w", err)
	}
	return nil
}

// WaitForFileChooser waits for a file chooser to be opened, executes the trigger function,
// and returns the resulting FileChooser. It times out after the given deadline from ctx.
func (p *Page) WaitForFileChooser(ctx context.Context, trigger func()) (*FileChooser, error) {
	ch := make(chan *FileChooser, 1)
	cancel := p.OnFileChooser(func(fc *FileChooser) {
		select {
		case ch <- fc:
		default:
		}
	})
	defer cancel()

	if trigger != nil {
		trigger()
	}

	select {
	case fc := <-ch:
		return fc, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("WaitForFileChooser: %w", ctx.Err())
	}
}

// OnFileChooser registers a handler that fires when a file chooser dialog is opened on this page.
// The returned function cancels the listener.
func (p *Page) OnFileChooser(handler func(*FileChooser)) func() {
	// Enable file-chooser events on the page channel.
	subCtx, subCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := p.owner.SendMessageRequest(subCtx, "updateSubscription", protocol.PageUpdateSubscriptionRequest{
		Event:   "fileChooser",
		Enabled: true,
	}); err != nil {
		slog.Default().Warn("page.OnFileChooser: failed to enable fileChooser subscription", "error", err)
	}
	subCancel()

	process := func(params json.RawMessage) {
		var event protocol.PageFileChooserEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Element.Guid == "" {
			return
		}
		fc := &FileChooser{
			page:       p,
			element:    p.owner.child(event.Element.Guid),
			isMultiple: event.IsMultiple,
		}
		go handler(fc)
	}
	id := p.owner.conn.OnEvent(p.owner.guid, "fileChooser", process)
	return func() { p.owner.conn.OffEvent(p.owner.guid, "fileChooser", id) }
}
