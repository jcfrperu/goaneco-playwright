package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// Rect represents a rectangle for clipping or bounding-box operations.
type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// ScreenshotOptions controls the output format of Page.Screenshot.
type ScreenshotOptions struct {
	// Path is the file path to save the screenshot to. If empty, the screenshot is only returned as bytes.
	Path string
	// Type is "png" (default) or "jpeg".
	Type string
	// Quality is JPEG quality 0–100. Only applies when Type is "jpeg".
	Quality *int
	// FullPage captures the full scrollable page.
	FullPage *bool
	// Clip restricts the screenshot to a sub-rectangle of the viewport.
	Clip *Rect
	// OmitBackground hides the default white background and allows capturing screenshots with transparency.
	OmitBackground *bool
	// Scale is "css" (default) or "device". Use "device" to capture at device pixel ratio.
	Scale string
	// Mask is a list of locators that should be masked (blacked out) in the screenshot.
	Mask []*Locator
}

type pageScreenshotRequestWire struct {
	Timeout        float64        `json:"timeout"`
	Clip           *protocol.Rect `json:"clip,omitempty"`
	FullPage       *bool          `json:"fullPage,omitempty"`
	Quality        *int           `json:"quality,omitempty"`
	Type           *string        `json:"type,omitempty"`
	OmitBackground *bool          `json:"omitBackground,omitempty"`
	Scale          *string        `json:"scale,omitempty"`
	Mask           []maskEntry    `json:"mask,omitempty"`
}

type maskEntry struct {
	Frame    string `json:"frame"`
	Selector string `json:"selector"`
}

// Screenshot captures a screenshot of the current page viewport.
// Returns the raw image bytes (PNG by default).
func (p *Page) Screenshot(ctx context.Context, opts ...*ScreenshotOptions) ([]byte, error) {
	req := pageScreenshotRequestWire{
		Timeout: defaultLocatorTimeout,
	}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.Type != "" {
			req.Type = &o.Type
		}
		req.Quality = o.Quality
		req.FullPage = o.FullPage
		req.OmitBackground = o.OmitBackground
		if o.Scale != "" {
			req.Scale = &o.Scale
		}
		if o.Clip != nil {
			req.Clip = &protocol.Rect{
				X:      o.Clip.X,
				Y:      o.Clip.Y,
				Width:  o.Clip.Width,
				Height: o.Clip.Height,
			}
		}
		for _, loc := range o.Mask {
			req.Mask = append(req.Mask, maskEntry{
				Frame:    loc.frame.guid,
				Selector: loc.selector,
			})
		}
	}

	result, err := p.owner.SendMessageRequest(ctx, "screenshot", req)
	if err != nil {
		return nil, fmt.Errorf("page.screenshot failed: %w", err)
	}

	var resp protocol.PageScreenshotResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse screenshot response: %w", err)
	}
	if len(opts) > 0 && opts[0] != nil && opts[0].Path != "" {
		if err := os.WriteFile(opts[0].Path, resp.Binary, 0644); err != nil {
			return nil, fmt.Errorf("page.screenshot: failed to write file %q: %w", opts[0].Path, err)
		}
	}
	return resp.Binary, nil
}
