package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// ScreenshotAnimations specifies animation handling for screenshots.
type ScreenshotAnimations string

const (
	ScreenshotAnimationsDisabled ScreenshotAnimations = "disabled"
	ScreenshotAnimationsAllow    ScreenshotAnimations = "allow"
)

// ScreenshotCaret specifies caret handling for screenshots.
type ScreenshotCaret string

const (
	ScreenshotCaretHide    ScreenshotCaret = "hide"
	ScreenshotCaretInitial ScreenshotCaret = "initial"
)

// ScreenshotType specifies the image format.
type ScreenshotType string

const (
	ScreenshotTypePNG  ScreenshotType = "png"
	ScreenshotTypeJPEG ScreenshotType = "jpeg"
	ScreenshotTypeWebP ScreenshotType = "webp"
)

// Rect represents a rectangle for clipping or bounding-box operations.
type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// ScreenshotOptions controls the output format of Page.Screenshot and Locator.Screenshot.
type ScreenshotOptions struct {
	// Path is the file path to save the screenshot to. If empty, the screenshot is only returned as bytes.
	Path string
	// Type is "png" (default), "jpeg", or "webp".
	Type string
	// Quality is JPEG/WebP quality 0–100.
	Quality *int
	// FullPage captures the full scrollable page.
	FullPage *bool
	// Clip restricts the screenshot to a sub-rectangle of the viewport.
	Clip *Rect
	// OmitBackground hides the default white background and allows capturing screenshots with transparency.
	OmitBackground *bool
	// Scale is "css" (default) or "device". Use "device" to capture at device pixel ratio.
	Scale string
	// Caret is "hide" (default) or "initial".
	Caret string
	// Animations is "disabled" or "allow".
	Animations string
	// Mask is a list of locators that should be masked (blacked out) in the screenshot.
	Mask []*Locator
	// MaskColor is the custom color for the mask overlay, e.g. "#00FF00".
	MaskColor string
	// Style is custom stylesheet text to inject while taking the screenshot.
	Style string
	// Timeout in milliseconds.
	Timeout *float64
}

func (o *ScreenshotOptions) SetPath(p string) *ScreenshotOptions    { o.Path = p; return o }
func (o *ScreenshotOptions) SetType(t string) *ScreenshotOptions    { o.Type = t; return o }
func (o *ScreenshotOptions) SetQuality(q int) *ScreenshotOptions    { o.Quality = &q; return o }
func (o *ScreenshotOptions) SetFullPage(fp bool) *ScreenshotOptions { o.FullPage = &fp; return o }
func (o *ScreenshotOptions) SetClip(c *Rect) *ScreenshotOptions     { o.Clip = c; return o }
func (o *ScreenshotOptions) SetOmitBackground(ob bool) *ScreenshotOptions {
	o.OmitBackground = &ob
	return o
}
func (o *ScreenshotOptions) SetScale(s string) *ScreenshotOptions      { o.Scale = s; return o }
func (o *ScreenshotOptions) SetCaret(c string) *ScreenshotOptions      { o.Caret = c; return o }
func (o *ScreenshotOptions) SetAnimations(a string) *ScreenshotOptions { o.Animations = a; return o }
func (o *ScreenshotOptions) SetMask(m []*Locator) *ScreenshotOptions   { o.Mask = m; return o }
func (o *ScreenshotOptions) SetMaskColor(mc string) *ScreenshotOptions { o.MaskColor = mc; return o }
func (o *ScreenshotOptions) SetStyle(s string) *ScreenshotOptions      { o.Style = s; return o }
func (o *ScreenshotOptions) SetTimeout(t float64) *ScreenshotOptions   { o.Timeout = &t; return o }

func determineScreenshotType(t, path string) *string {
	if t != "" {
		return &t
	}
	if path != "" {
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".png":
			res := "png"
			return &res
		case ".jpeg", ".jpg":
			res := "jpeg"
			return &res
		case ".webp":
			res := "webp"
			return &res
		}
	}
	return nil
}

type pageScreenshotRequestWire struct {
	Timeout        float64        `json:"timeout"`
	Clip           *protocol.Rect `json:"clip,omitempty"`
	FullPage       *bool          `json:"fullPage,omitempty"`
	Quality        *int           `json:"quality,omitempty"`
	Type           *string        `json:"type,omitempty"`
	OmitBackground *bool          `json:"omitBackground,omitempty"`
	Scale          *string        `json:"scale,omitempty"`
	Caret          *string        `json:"caret,omitempty"`
	Animations     *string        `json:"animations,omitempty"`
	Mask           []maskEntry    `json:"mask,omitempty"`
	MaskColor      *string        `json:"maskColor,omitempty"`
	Style          *string        `json:"style,omitempty"`
}

type maskEntry struct {
	Frame struct {
		Guid string `json:"guid"`
	} `json:"frame"`
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
		if o.Timeout != nil {
			req.Timeout = *o.Timeout
		}
		imgType := determineScreenshotType(o.Type, o.Path)
		if imgType != nil {
			req.Type = imgType
		}
		req.Quality = o.Quality
		req.FullPage = o.FullPage
		req.OmitBackground = o.OmitBackground
		if o.Scale != "" {
			req.Scale = &o.Scale
		}
		if o.Caret != "" {
			req.Caret = &o.Caret
		}
		if o.Animations != "" {
			req.Animations = &o.Animations
		}
		if o.MaskColor != "" {
			req.MaskColor = &o.MaskColor
		}
		if o.Style != "" {
			req.Style = &o.Style
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
				Frame: struct {
					Guid string `json:"guid"`
				}{Guid: loc.frame.guid},
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
		savePath := opts[0].Path
		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			return nil, fmt.Errorf("page.screenshot: failed to create directory: %w", err)
		}
		if err := os.WriteFile(savePath, resp.Binary, 0644); err != nil {
			return nil, fmt.Errorf("page.screenshot: failed to write file %q: %w", savePath, err)
		}
	}
	return resp.Binary, nil
}
