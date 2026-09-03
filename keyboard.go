package playwright

import (
	"context"
	"fmt"
)

// Keyboard provides keyboard input methods for a page.
// Access via Page.Keyboard.
type Keyboard struct {
	page *Page
}

// Press presses a key and releases it (down + up). Supports modifier combos like "Shift+A".
func (k *Keyboard) Press(ctx context.Context, key string) error {
	_, err := k.page.owner.SendMessageRequest(ctx, "keyboardPress", map[string]any{
		"key": key,
	})
	if err != nil {
		return fmt.Errorf("keyboard.press(%q) failed: %w", key, err)
	}
	return nil
}

// Down dispatches a keydown event for the given key.
func (k *Keyboard) Down(ctx context.Context, key string) error {
	_, err := k.page.owner.SendMessageRequest(ctx, "keyboardDown", map[string]any{
		"key": key,
	})
	if err != nil {
		return fmt.Errorf("keyboard.down(%q) failed: %w", key, err)
	}
	return nil
}

// Up dispatches a keyup event for the given key.
func (k *Keyboard) Up(ctx context.Context, key string) error {
	_, err := k.page.owner.SendMessageRequest(ctx, "keyboardUp", map[string]any{
		"key": key,
	})
	if err != nil {
		return fmt.Errorf("keyboard.up(%q) failed: %w", key, err)
	}
	return nil
}

// Type types the given text by sending keydown, keypress/input, and keyup events for each character.
func (k *Keyboard) Type(ctx context.Context, text string) error {
	_, err := k.page.owner.SendMessageRequest(ctx, "keyboardType", map[string]any{
		"text": text,
	})
	if err != nil {
		return fmt.Errorf("keyboard.type failed: %w", err)
	}
	return nil
}

// InsertText dispatches an input event with the given text, bypassing key events.
func (k *Keyboard) InsertText(ctx context.Context, text string) error {
	_, err := k.page.owner.SendMessageRequest(ctx, "keyboardInsertText", map[string]any{
		"text": text,
	})
	if err != nil {
		return fmt.Errorf("keyboard.insertText failed: %w", err)
	}
	return nil
}
