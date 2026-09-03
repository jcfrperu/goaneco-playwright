package playwright

import (
	"context"
	"encoding/json"
	"fmt"
)

// RecordVideoOptions configures video recording for a BrowserContext.
type RecordVideoOptions struct {
	// Dir is the directory in which to store recorded videos.
	Dir string `json:"dir"`
	// Size is the video frame size; defaults to the viewport size.
	Size *ViewportSize `json:"size,omitempty"`
}

// Video represents a recorded video for a page.
// The video file is only available after the page (and its recording) is closed.
type Video struct {
	owner ChannelOwner
}

// Path waits for the video recording to finish and returns the path to the video file.
func (v *Video) Path(ctx context.Context) (string, error) {
	result, err := v.owner.SendMessageRequest(ctx, "pathAfterFinished", struct{}{})
	if err != nil {
		return "", fmt.Errorf("video.path failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse video.path response: %w", err)
	}
	return resp.Value, nil
}

// SaveAs copies the recorded video to the destination path.
func (v *Video) SaveAs(ctx context.Context, path string) error {
	_, err := v.owner.SendMessageRequest(ctx, "saveAs", map[string]string{"path": path})
	if err != nil {
		return fmt.Errorf("video.saveAs failed: %w", err)
	}
	return nil
}

// Delete removes the recorded video file.
func (v *Video) Delete(ctx context.Context) error {
	_, err := v.owner.SendMessageRequest(ctx, "delete", struct{}{})
	if err != nil {
		return fmt.Errorf("video.delete failed: %w", err)
	}
	return nil
}

// Video returns the Video associated with this page, or nil if recording is not enabled.
// The video file is only accessible after the page is closed.
func (p *Page) Video() *Video {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.video
}

// subscribeToVideo listens for the "video" event on the page, which fires once
// the video artifact is ready (typically after the page is closed when recording).
func (p *Page) subscribeToVideo() {
	p.owner.conn.OnEventOnce(p.owner.guid, "video", func(params json.RawMessage) {
		var event struct {
			Artifact struct {
				Guid string `json:"guid"`
			} `json:"artifact"`
		}
		if err := json.Unmarshal(params, &event); err != nil || event.Artifact.Guid == "" {
			return
		}
		p.mu.Lock()
		p.video = &Video{owner: p.owner.child(event.Artifact.Guid)}
		p.mu.Unlock()
	})
}
