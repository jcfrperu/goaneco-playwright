package playwright

import (
	"context"
	"fmt"
)

// Mouse provides mouse input methods for a page, addressing elements by viewport coordinates.
// Access via Page.Mouse.
type Mouse struct {
	page *Page
}

// Move moves the mouse to the given (x, y) coordinates.
func (m *Mouse) Move(ctx context.Context, x, y float64) error {
	_, err := m.page.owner.SendMessageRequest(ctx, "mouseMove", map[string]any{
		"x": x,
		"y": y,
	})
	if err != nil {
		return fmt.Errorf("mouse.move failed: %w", err)
	}
	return nil
}

// Click moves to (x, y) and performs a mouse click.
// button is optional and defaults to "left".
func (m *Mouse) Click(ctx context.Context, x, y float64, button ...string) error {
	btn := "left"
	if len(button) > 0 && button[0] != "" {
		btn = button[0]
	}
	_, err := m.page.owner.SendMessageRequest(ctx, "mouseClick", map[string]any{
		"x":      x,
		"y":      y,
		"button": btn,
	})
	if err != nil {
		return fmt.Errorf("mouse.click failed: %w", err)
	}
	return nil
}

// DblClick moves to (x, y) and performs a double-click.
// button is optional and defaults to "left".
func (m *Mouse) DblClick(ctx context.Context, x, y float64, button ...string) error {
	btn := "left"
	if len(button) > 0 && button[0] != "" {
		btn = button[0]
	}
	count := 2
	_, err := m.page.owner.SendMessageRequest(ctx, "mouseClick", map[string]any{
		"x":          x,
		"y":          y,
		"button":     btn,
		"clickCount": count,
	})
	if err != nil {
		return fmt.Errorf("mouse.dblClick failed: %w", err)
	}
	return nil
}

// Down presses a mouse button at the current position.
// button is optional and defaults to "left".
func (m *Mouse) Down(ctx context.Context, button ...string) error {
	btn := "left"
	if len(button) > 0 && button[0] != "" {
		btn = button[0]
	}
	_, err := m.page.owner.SendMessageRequest(ctx, "mouseDown", map[string]any{
		"button": btn,
	})
	if err != nil {
		return fmt.Errorf("mouse.down failed: %w", err)
	}
	return nil
}

// Up releases a mouse button at the current position.
// button is optional and defaults to "left".
func (m *Mouse) Up(ctx context.Context, button ...string) error {
	btn := "left"
	if len(button) > 0 && button[0] != "" {
		btn = button[0]
	}
	_, err := m.page.owner.SendMessageRequest(ctx, "mouseUp", map[string]any{
		"button": btn,
	})
	if err != nil {
		return fmt.Errorf("mouse.up failed: %w", err)
	}
	return nil
}

// Wheel dispatches a wheel event at the current mouse position.
func (m *Mouse) Wheel(ctx context.Context, deltaX, deltaY float64) error {
	_, err := m.page.owner.SendMessageRequest(ctx, "mouseWheel", map[string]any{
		"deltaX": deltaX,
		"deltaY": deltaY,
	})
	if err != nil {
		return fmt.Errorf("mouse.wheel failed: %w", err)
	}
	return nil
}
