package playwright

import (
	"context"
	"fmt"
)

// Touchscreen provides methods for interacting with a page's touch screen.
// Access via Page.Touchscreen().
type Touchscreen struct {
	page *Page
}

// Touchscreen returns the Touchscreen input controller for this page.
func (p *Page) Touchscreen() *Touchscreen {
	return &Touchscreen{page: p}
}

// Tap taps the given (x, y) coordinates.
func (t *Touchscreen) Tap(ctx context.Context, x, y float64) error {
	_, err := t.page.owner.SendMessageRequest(ctx, "touchscreenTap", map[string]any{
		"x": x,
		"y": y,
	})
	if err != nil {
		return fmt.Errorf("touchscreen.tap failed: %w", err)
	}
	return nil
}
