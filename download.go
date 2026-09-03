package playwright

import (
	"context"
	"encoding/json"
	"fmt"
)

// Download represents a file download started by the browser.
type Download struct {
	artifactOwner ChannelOwner
	url           string
	suggestedName string
}

// URL returns the download URL.
func (d *Download) URL() string { return d.url }

// SuggestedFilename returns the browser-suggested filename for the download.
func (d *Download) SuggestedFilename() string { return d.suggestedName }

// Path returns the path to the completed download file.
// Blocks until the download is finished.
func (d *Download) Path(ctx context.Context) (string, error) {
	result, err := d.artifactOwner.SendMessageRequest(ctx, "pathAfterFinished", struct{}{})
	if err != nil {
		return "", fmt.Errorf("download.path failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse download.path response: %w", err)
	}
	return resp.Value, nil
}

// SaveAs copies the downloaded file to the specified path.
func (d *Download) SaveAs(ctx context.Context, path string) error {
	_, err := d.artifactOwner.SendMessageRequest(ctx, "saveAs", map[string]string{"path": path})
	if err != nil {
		return fmt.Errorf("download.saveAs failed: %w", err)
	}
	return nil
}

// Delete removes the downloaded file.
func (d *Download) Delete(ctx context.Context) error {
	_, err := d.artifactOwner.SendMessageRequest(ctx, "delete", struct{}{})
	if err != nil {
		return fmt.Errorf("download.delete failed: %w", err)
	}
	return nil
}

// Cancel cancels the in-progress download.
func (d *Download) Cancel(ctx context.Context) error {
	_, err := d.artifactOwner.SendMessageRequest(ctx, "cancel", struct{}{})
	if err != nil {
		return fmt.Errorf("download.cancel failed: %w", err)
	}
	return nil
}

// Failure returns an error if the download failed, or nil if it succeeded.
func (d *Download) Failure(ctx context.Context) error {
	result, err := d.artifactOwner.SendMessageRequest(ctx, "failure", struct{}{})
	if err != nil {
		return fmt.Errorf("download.failure IPC call failed: %w", err)
	}
	var resp struct {
		Error *string `json:"error,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to parse download.failure response: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("download failed: %s", *resp.Error)
	}
	return nil
}

// OnDownload registers a handler for file downloads initiated by this page.
// The handler is called in a goroutine for each download.
// The returned function cancels the listener.
func (p *Page) OnDownload(handler func(*Download)) func() {
	process := func(params json.RawMessage) {
		var event struct {
			Artifact struct {
				Guid string `json:"guid"`
			} `json:"artifact"`
			SuggestedFilename string `json:"suggestedFilename"`
			URL               string `json:"url"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		if event.Artifact.Guid == "" {
			return
		}
		d := &Download{
			artifactOwner: p.owner.child(event.Artifact.Guid),
			url:           event.URL,
			suggestedName: event.SuggestedFilename,
		}
		go handler(d)
	}
	id := p.owner.conn.OnEvent(p.owner.guid, "download", process)
	return func() { p.owner.conn.OffEvent(p.owner.guid, "download", id) }
}
