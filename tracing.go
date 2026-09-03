package playwright

import (
	"context"
	"encoding/json"
	"fmt"
)

// TracingStartOptions configures Tracing.Start.
type TracingStartOptions struct {
	Name        *string
	Screenshots *bool
	Snapshots   *bool
	Sources     *bool
}

type tracingStartWire struct {
	Name           *string `json:"name,omitempty"`
	Screencast     *bool   `json:"screencast,omitempty"`
	SnapshotDom    *bool   `json:"snapshotDom,omitempty"`
	SnapshotScreen *bool   `json:"snapshotScreen,omitempty"`
}

// Tracing provides methods to start and stop Playwright traces.
// Traces capture browser interactions as a zip archive for later analysis.
type Tracing struct {
	owner ChannelOwner
}

// Start begins a trace recording.
func (t *Tracing) Start(ctx context.Context, opts ...*TracingStartOptions) error {
	wire := tracingStartWire{}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		wire.Name = o.Name
		wire.Screencast = o.Screenshots
		wire.SnapshotDom = o.Snapshots
		wire.SnapshotScreen = o.Screenshots
	}
	_, err := t.owner.SendMessageRequest(ctx, "tracingStart", wire)
	if err != nil {
		return fmt.Errorf("tracing.start failed: %w", err)
	}
	return nil
}

// Stop ends the trace recording without saving.
func (t *Tracing) Stop(ctx context.Context) error {
	_, err := t.owner.SendMessageRequest(ctx, "tracingStop", struct{}{})
	if err != nil {
		return fmt.Errorf("tracing.stop failed: %w", err)
	}
	return nil
}

// StartChunk begins a new chunk of an ongoing trace recording.
// Returns the generated trace name.
func (t *Tracing) StartChunk(ctx context.Context, name ...string) (string, error) {
	req := struct {
		Name *string `json:"name,omitempty"`
	}{}
	if len(name) > 0 {
		req.Name = &name[0]
	}
	result, err := t.owner.SendMessageRequest(ctx, "tracingStartChunk", req)
	if err != nil {
		return "", fmt.Errorf("tracing.startChunk failed: %w", err)
	}
	var resp struct {
		TraceName string `json:"traceName"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse startChunk response: %w", err)
	}
	return resp.TraceName, nil
}

// StopChunk ends the current trace chunk.
// If path is non-empty the trace is saved to that path; otherwise it is discarded.
func (t *Tracing) StopChunk(ctx context.Context, path string) error {
	mode := "discard"
	if path != "" {
		mode = "archive"
	}
	result, err := t.owner.SendMessageRequest(ctx, "tracingStopChunk", map[string]string{"mode": mode})
	if err != nil {
		return fmt.Errorf("tracing.stopChunk failed: %w", err)
	}
	if path == "" {
		return nil
	}
	var resp struct {
		Artifact *struct {
			Guid string `json:"guid"`
		} `json:"artifact,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil || resp.Artifact == nil {
		return nil
	}
	artifactOwner := t.owner.child(resp.Artifact.Guid)
	if _, err := artifactOwner.SendMessageRequest(ctx, "saveAs", map[string]string{"path": path}); err != nil {
		return fmt.Errorf("tracing.stopChunk: saveAs to %q failed: %w", path, err)
	}
	return nil
}

// Tracing returns the Tracing object for this BrowserContext.
// The Tracing GUID comes from the BrowserContext initializer.
func (c *BrowserContext) Tracing() *Tracing {
	var init struct {
		Tracing struct {
			Guid string `json:"guid"`
		} `json:"tracing"`
	}
	if raw := c.initializer; len(raw) > 0 {
		_ = json.Unmarshal(raw, &init) // best-effort; zero-value fallback if server data is malformed
	}
	if init.Tracing.Guid == "" {
		return nil
	}
	return &Tracing{owner: c.owner.child(init.Tracing.Guid)}
}
