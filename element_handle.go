package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// BoundingBox represents the bounding rectangle of an element on screen.
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ElementHandle represents a reference to a DOM element in the browser.
type ElementHandle struct {
	owner ChannelOwner
}

// elementHandleFromGUID creates an ElementHandle from a GUID.
func elementHandleFromGUID(parent ChannelOwner, guid string) *ElementHandle {
	return &ElementHandle{owner: parent.child(guid)}
}

// AsJSHandle returns this element handle as a JSHandle.
func (e *ElementHandle) AsJSHandle() *JSHandle {
	return &JSHandle{owner: e.owner}
}

// BoundingBox returns the bounding box of the element, or nil if not visible.
func (e *ElementHandle) BoundingBox(ctx context.Context) (*BoundingBox, error) {
	result, err := e.owner.SendMessageRequest(ctx, "boundingBox", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("elementHandle.boundingBox failed: %w", err)
	}
	var resp struct {
		Value *protocol.Rect `json:"value,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse boundingBox response: %w", err)
	}
	if resp.Value == nil {
		return nil, nil
	}
	return &BoundingBox{
		X:      resp.Value.X,
		Y:      resp.Value.Y,
		Width:  resp.Value.Width,
		Height: resp.Value.Height,
	}, nil
}

// Click performs a click on the element.
func (e *ElementHandle) Click(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "click", struct{}{})
	if err != nil {
		return fmt.Errorf("elementHandle.click failed: %w", err)
	}
	return nil
}

// DblClick performs a double-click on the element.
func (e *ElementHandle) DblClick(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "dblclick", struct{}{})
	if err != nil {
		return fmt.Errorf("elementHandle.dblclick failed: %w", err)
	}
	return nil
}

// Check checks a checkbox or radio button element.
func (e *ElementHandle) Check(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "check", map[string]any{"timeout": defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("elementHandle.check failed: %w", err)
	}
	return nil
}

// Uncheck unchecks a checkbox element.
func (e *ElementHandle) Uncheck(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "uncheck", map[string]any{"timeout": defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("elementHandle.uncheck failed: %w", err)
	}
	return nil
}

// Fill sets the value of an input element.
func (e *ElementHandle) Fill(ctx context.Context, value string) error {
	_, err := e.owner.SendMessageRequest(ctx, "fill", map[string]any{
		"value":   value,
		"timeout": defaultActionTimeoutMs,
	})
	if err != nil {
		return fmt.Errorf("elementHandle.fill failed: %w", err)
	}
	return nil
}

// Focus gives focus to the element.
func (e *ElementHandle) Focus(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "focus", struct{}{})
	if err != nil {
		return fmt.Errorf("elementHandle.focus failed: %w", err)
	}
	return nil
}

// Hover moves the mouse to the center of the element.
func (e *ElementHandle) Hover(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "hover", map[string]any{"timeout": defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("elementHandle.hover failed: %w", err)
	}
	return nil
}

// Tap taps the element (mobile-like gesture).
func (e *ElementHandle) Tap(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "tap", map[string]any{"timeout": defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("elementHandle.tap failed: %w", err)
	}
	return nil
}

// Press presses a key while the element is focused.
func (e *ElementHandle) Press(ctx context.Context, key string) error {
	_, err := e.owner.SendMessageRequest(ctx, "press", map[string]string{"key": key})
	if err != nil {
		return fmt.Errorf("elementHandle.press failed: %w", err)
	}
	return nil
}

// Type types text into the element, one character at a time.
func (e *ElementHandle) Type(ctx context.Context, text string) error {
	_, err := e.owner.SendMessageRequest(ctx, "type", map[string]any{"text": text, "timeout": defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("elementHandle.type failed: %w", err)
	}
	return nil
}

// GetAttribute returns the value of the named attribute, or nil if the attribute does not exist.
func (e *ElementHandle) GetAttribute(ctx context.Context, name string) (*string, error) {
	result, err := e.owner.SendMessageRequest(ctx, "getAttribute", map[string]string{"name": name})
	if err != nil {
		return nil, fmt.Errorf("elementHandle.getAttribute failed: %w", err)
	}
	var resp struct {
		Value *string `json:"value,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse getAttribute response: %w", err)
	}
	return resp.Value, nil
}

// InnerHTML returns the inner HTML of the element.
func (e *ElementHandle) InnerHTML(ctx context.Context) (string, error) {
	result, err := e.owner.SendMessageRequest(ctx, "innerHTML", struct{}{})
	if err != nil {
		return "", fmt.Errorf("elementHandle.innerHTML failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse innerHTML response: %w", err)
	}
	return resp.Value, nil
}

// InnerText returns the inner text of the element.
func (e *ElementHandle) InnerText(ctx context.Context) (string, error) {
	result, err := e.owner.SendMessageRequest(ctx, "innerText", struct{}{})
	if err != nil {
		return "", fmt.Errorf("elementHandle.innerText failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse innerText response: %w", err)
	}
	return resp.Value, nil
}

// InputValue returns the current value of an input, textarea, or select element.
func (e *ElementHandle) InputValue(ctx context.Context) (string, error) {
	result, err := e.owner.SendMessageRequest(ctx, "inputValue", struct{}{})
	if err != nil {
		return "", fmt.Errorf("elementHandle.inputValue failed: %w", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse inputValue response: %w", err)
	}
	return resp.Value, nil
}

// TextContent returns the text content of the element and its descendants.
func (e *ElementHandle) TextContent(ctx context.Context) (string, error) {
	result, err := e.owner.SendMessageRequest(ctx, "textContent", struct{}{})
	if err != nil {
		return "", fmt.Errorf("elementHandle.textContent failed: %w", err)
	}
	var resp struct {
		Value *string `json:"value,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse textContent response: %w", err)
	}
	if resp.Value == nil {
		return "", nil
	}
	return *resp.Value, nil
}

// IsChecked returns true if the element is checked.
func (e *ElementHandle) IsChecked(ctx context.Context) (bool, error) {
	result, err := e.owner.SendMessageRequest(ctx, "isChecked", struct{}{})
	if err != nil {
		return false, fmt.Errorf("elementHandle.isChecked failed: %w", err)
	}
	var resp struct {
		Value bool `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isChecked response: %w", err)
	}
	return resp.Value, nil
}

// IsDisabled returns true if the element is disabled.
func (e *ElementHandle) IsDisabled(ctx context.Context) (bool, error) {
	result, err := e.owner.SendMessageRequest(ctx, "isDisabled", struct{}{})
	if err != nil {
		return false, fmt.Errorf("elementHandle.isDisabled failed: %w", err)
	}
	var resp struct {
		Value bool `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isDisabled response: %w", err)
	}
	return resp.Value, nil
}

// IsEditable returns true if the element is editable.
func (e *ElementHandle) IsEditable(ctx context.Context) (bool, error) {
	result, err := e.owner.SendMessageRequest(ctx, "isEditable", struct{}{})
	if err != nil {
		return false, fmt.Errorf("elementHandle.isEditable failed: %w", err)
	}
	var resp struct {
		Value bool `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isEditable response: %w", err)
	}
	return resp.Value, nil
}

// IsEnabled returns true if the element is enabled.
func (e *ElementHandle) IsEnabled(ctx context.Context) (bool, error) {
	result, err := e.owner.SendMessageRequest(ctx, "isEnabled", struct{}{})
	if err != nil {
		return false, fmt.Errorf("elementHandle.isEnabled failed: %w", err)
	}
	var resp struct {
		Value bool `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isEnabled response: %w", err)
	}
	return resp.Value, nil
}

// IsHidden returns true if the element is hidden.
func (e *ElementHandle) IsHidden(ctx context.Context) (bool, error) {
	result, err := e.owner.SendMessageRequest(ctx, "isHidden", struct{}{})
	if err != nil {
		return false, fmt.Errorf("elementHandle.isHidden failed: %w", err)
	}
	var resp struct {
		Value bool `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isHidden response: %w", err)
	}
	return resp.Value, nil
}

// IsVisible returns true if the element is visible.
func (e *ElementHandle) IsVisible(ctx context.Context) (bool, error) {
	result, err := e.owner.SendMessageRequest(ctx, "isVisible", struct{}{})
	if err != nil {
		return false, fmt.Errorf("elementHandle.isVisible failed: %w", err)
	}
	var resp struct {
		Value bool `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isVisible response: %w", err)
	}
	return resp.Value, nil
}

// ScrollIntoViewIfNeeded scrolls the element into view if it is not already visible.
func (e *ElementHandle) ScrollIntoViewIfNeeded(ctx context.Context) error {
	_, err := e.owner.SendMessageRequest(ctx, "scrollIntoViewIfNeeded", map[string]any{"timeout": defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("elementHandle.scrollIntoViewIfNeeded failed: %w", err)
	}
	return nil
}

// WaitForElementState waits until the element reaches the given state.
// Valid states: "visible", "hidden", "stable", "enabled", "disabled", "editable".
func (e *ElementHandle) WaitForElementState(ctx context.Context, state string) error {
	_, err := e.owner.SendMessageRequest(ctx, "waitForElementState", map[string]any{
		"state":   state,
		"timeout": defaultActionTimeoutMs,
	})
	if err != nil {
		return fmt.Errorf("elementHandle.waitForElementState(%q) failed: %w", state, err)
	}
	return nil
}

// QuerySelector returns the first child element matching the selector, or nil.
func (e *ElementHandle) QuerySelector(ctx context.Context, selector string) (*ElementHandle, error) {
	result, err := e.owner.SendMessageRequest(ctx, "querySelector", map[string]string{"selector": selector})
	if err != nil {
		return nil, fmt.Errorf("elementHandle.querySelector failed: %w", err)
	}
	var resp struct {
		Element *struct {
			Guid string `json:"guid"`
		} `json:"element,omitempty"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse querySelector response: %w", err)
	}
	if resp.Element == nil {
		return nil, nil
	}
	return elementHandleFromGUID(e.owner, resp.Element.Guid), nil
}

// QuerySelectorAll returns all child elements matching the selector.
func (e *ElementHandle) QuerySelectorAll(ctx context.Context, selector string) ([]*ElementHandle, error) {
	result, err := e.owner.SendMessageRequest(ctx, "querySelectorAll", map[string]string{"selector": selector})
	if err != nil {
		return nil, fmt.Errorf("elementHandle.querySelectorAll failed: %w", err)
	}
	var resp struct {
		Elements []struct {
			Guid string `json:"guid"`
		} `json:"elements"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse querySelectorAll response: %w", err)
	}
	handles := make([]*ElementHandle, len(resp.Elements))
	for i, el := range resp.Elements {
		handles[i] = elementHandleFromGUID(e.owner, el.Guid)
	}
	return handles, nil
}

// SetChecked sets the checked state of the element.
func (e *ElementHandle) SetChecked(ctx context.Context, checked bool) error {
	method := "uncheck"
	if checked {
		method = "check"
	}
	_, err := e.owner.SendMessageRequest(ctx, method, map[string]any{"timeout": defaultActionTimeoutMs})
	if err != nil {
		return fmt.Errorf("elementHandle.setChecked(%v) failed: %w", checked, err)
	}
	return nil
}

// SelectOption selects options in a <select> element by value.
// Returns the list of selected option values.
func (e *ElementHandle) SelectOption(ctx context.Context, values ...string) ([]string, error) {
	opts := make([]any, len(values))
	for i, v := range values {
		opts[i] = map[string]string{"value": v}
	}
	result, err := e.owner.SendMessageRequest(ctx, "selectOption", map[string]any{
		"options":  opts,
		"elements": []any{},
		"timeout":  defaultActionTimeoutMs,
	})
	if err != nil {
		return nil, fmt.Errorf("elementHandle.selectOption failed: %w", err)
	}
	var resp struct {
		Values []string `json:"values"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse selectOption response: %w", err)
	}
	return resp.Values, nil
}

// Evaluate executes a JavaScript expression with this element as the argument.
func (e *ElementHandle) Evaluate(ctx context.Context, expression string, arg ...any) (any, error) {
	return e.AsJSHandle().Evaluate(ctx, expression, arg...)
}

// EvaluateHandle executes a JavaScript expression with this element as the argument and returns a JSHandle.
func (e *ElementHandle) EvaluateHandle(ctx context.Context, expression string, arg ...any) (*JSHandle, error) {
	return e.AsJSHandle().EvaluateHandle(ctx, expression, arg...)
}

// JSONValue returns the JSON value of this element handle.
func (e *ElementHandle) JSONValue(ctx context.Context) (any, error) {
	return e.AsJSHandle().JSONValue(ctx)
}

// DispatchEvent dispatches an event of the given type on the element.
// eventInit can optionally provide event initialization data.
func (e *ElementHandle) DispatchEvent(ctx context.Context, eventType string, eventInit ...any) error {
	var initArg any
	if len(eventInit) > 0 {
		initArg = eventInit[0]
	}
	_, err := e.owner.SendMessageRequest(ctx, "dispatchEvent", map[string]any{
		"type":      eventType,
		"eventInit": serializeArgument(initArg),
	})
	if err != nil {
		return fmt.Errorf("elementHandle.dispatchEvent(%q) failed: %w", eventType, err)
	}
	return nil
}

// Dispose releases the element handle from the browser.
func (e *ElementHandle) Dispose(ctx context.Context) error {
	return e.AsJSHandle().Dispose(ctx)
}

// Screenshot captures a screenshot of this element.
// Returns raw PNG bytes by default; pass ScreenshotOptions to select JPEG or quality.
// If opts.Path is set, the image is also written to that file path.
func (e *ElementHandle) Screenshot(ctx context.Context, opts ...*ScreenshotOptions) ([]byte, error) {
	req := map[string]any{"timeout": defaultActionTimeoutMs}
	var savePath string
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.Type != "" {
			req["type"] = o.Type
		}
		if o.Quality != nil {
			req["quality"] = *o.Quality
		}
		savePath = o.Path
	}
	result, err := e.owner.SendMessageRequest(ctx, "screenshot", req)
	if err != nil {
		return nil, fmt.Errorf("elementHandle.screenshot failed: %w", err)
	}
	var resp struct {
		Binary []byte `json:"binary"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("elementHandle.screenshot: parse response failed: %w", err)
	}
	if savePath != "" {
		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			return resp.Binary, fmt.Errorf("elementHandle.screenshot: failed to create directory: %w", err)
		}
		if err := os.WriteFile(savePath, resp.Binary, 0644); err != nil {
			return resp.Binary, fmt.Errorf("elementHandle.screenshot: failed to write file %q: %w", savePath, err)
		}
	}
	return resp.Binary, nil
}
