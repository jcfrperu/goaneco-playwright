package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// TracingStartOptions configures Tracing.Start.
type TracingStartOptions struct {
	Name        *string
	Screenshots *bool
	Snapshots   *bool
	Sources     *bool
}

// TracingGroupLocation identifies the source location associated with a trace group.
type TracingGroupLocation struct {
	File   string
	Line   *int
	Column *int
}

// TracingGroupOptions configures Tracing.Group.
type TracingGroupOptions struct {
	Location *TracingGroupLocation
}

// TracingStartHarOptions configures HAR recording through Tracing.
type TracingStartHarOptions struct {
	Content *string
	Mode    *string
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
	owner        ChannelOwner
	mu           sync.Mutex
	currentHarID string
	harPath      string
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
	chunk := struct {
		Name *string `json:"name,omitempty"`
	}{Name: wire.Name}
	if _, err := t.owner.SendMessageRequest(ctx, "tracingStartChunk", chunk); err != nil {
		return fmt.Errorf("tracing.start failed: %w", err)
	}
	return nil
}

// Stop ends the trace recording without saving.
func (t *Tracing) Stop(ctx context.Context) error {
	// Mirrors the official clients: stop the current chunk (discarding it) before stopping the trace.
	if _, err := t.owner.SendMessageRequest(ctx, "tracingStopChunk", map[string]string{"mode": "discard"}); err != nil {
		return fmt.Errorf("tracing.stop failed: %w", err)
	}
	_, err := t.owner.SendMessageRequest(ctx, "tracingStop", struct{}{})
	if err != nil {
		return fmt.Errorf("tracing.stop failed: %w", err)
	}
	return nil
}

// Group starts a named group in the trace.
func (t *Tracing) Group(ctx context.Context, name string, opts ...*TracingGroupOptions) error {
	var location any
	if len(opts) > 0 && opts[0] != nil && opts[0].Location != nil {
		loc := opts[0].Location
		location = struct {
			File   string `json:"file"`
			Line   *int   `json:"line,omitempty"`
			Column *int   `json:"column,omitempty"`
		}{File: loc.File, Line: loc.Line, Column: loc.Column}
	}
	_, err := t.owner.SendMessageRequest(ctx, "tracingGroup", protocol.TracingTracingGroupRequest{
		Name:     name,
		Location: location,
	})
	if err != nil {
		return fmt.Errorf("tracing.group failed: %w", err)
	}
	return nil
}

// GroupEnd closes the current trace group.
func (t *Tracing) GroupEnd(ctx context.Context) error {
	_, err := t.owner.SendMessageRequest(ctx, "tracingGroupEnd", protocol.TracingTracingGroupEndRequest{})
	if err != nil {
		return fmt.Errorf("tracing.groupEnd failed: %w", err)
	}
	return nil
}

// StartHar begins recording network activity to a HAR file.
func (t *Tracing) StartHar(ctx context.Context, path string, opts ...*TracingStartHarOptions) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.currentHarID != "" {
		return fmt.Errorf("tracing.startHar: HAR recording has already been started")
	}

	content := "embed"
	if filepath.Ext(path) == ".zip" {
		content = "attach"
	}
	mode := "full"
	if len(opts) > 0 && opts[0] != nil {
		if opts[0].Content != nil {
			content = *opts[0].Content
		}
		if opts[0].Mode != nil {
			mode = *opts[0].Mode
		}
	}
	harPath := path
	recordOptions := protocol.RecordHarOptions{Content: content, Mode: mode}
	if filepath.Ext(path) != ".zip" {
		recordOptions.HarPath = &harPath
	}
	result, err := t.owner.SendMessageRequest(ctx, "harStart", protocol.TracingHarStartRequest{Options: recordOptions})
	if err != nil {
		return fmt.Errorf("tracing.startHar failed: %w", err)
	}
	var response protocol.TracingHarStartResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return fmt.Errorf("failed to parse tracing.startHar response: %w", err)
	}
	if response.HarId == "" {
		return fmt.Errorf("tracing.startHar: server returned empty HAR id")
	}
	t.currentHarID = response.HarId
	t.harPath = path
	return nil
}

// StopHar stops the active HAR recording and writes the HAR artifact to its path.
func (t *Tracing) StopHar(ctx context.Context) error {
	t.mu.Lock()
	harID, harPath := t.currentHarID, t.harPath
	t.currentHarID = ""
	t.harPath = ""
	t.mu.Unlock()
	if harID == "" {
		return fmt.Errorf("tracing.stopHar: HAR recording has not been started")
	}
	exportMode := "archive"
	if filepath.Ext(harPath) != ".zip" {
		exportMode = "entries"
	}
	result, err := t.owner.SendMessageRequest(ctx, "harExport", protocol.TracingHarExportRequest{
		HarId: &harID,
		Mode:  exportMode,
	})
	if err != nil {
		return fmt.Errorf("tracing.stopHar failed: %w", err)
	}
	var response protocol.TracingHarExportResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return fmt.Errorf("failed to parse tracing.stopHar response: %w", err)
	}
	if response.Artifact != nil && response.Artifact.Guid != "" {
		artifact := t.owner.child(response.Artifact.Guid)
		if _, err := artifact.SendMessageRequest(ctx, "saveAs", map[string]string{"path": harPath}); err != nil {
			return fmt.Errorf("tracing.stopHar: saveAs to %q failed: %w", harPath, err)
		}
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
