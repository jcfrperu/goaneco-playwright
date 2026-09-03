package playwright

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// FilePayload represents an in-memory file for use with SetInputFilesPayload.
type FilePayload struct {
	Name     string
	MimeType string
	Buffer   []byte
}

// subSelectorArg JSON-encodes a selector string as required by internal:and=,
// internal:or=, internal:has=, and internal:has-not=. Playwright expects the
// value to be a JSON string literal (e.g. "button.highlighted"), matching what
// JSON.stringify(locator._selector) produces in the TypeScript client.
func subSelectorArg(sel string) string {
	b, _ := json.Marshal(sel)
	return string(b)
}

// AriaRole represents an ARIA role for use with GetByRole.
type AriaRole string

const (
	AriaRoleAlert            AriaRole = "alert"
	AriaRoleAlertDialog      AriaRole = "alertdialog"
	AriaRoleApplication      AriaRole = "application"
	AriaRoleArticle          AriaRole = "article"
	AriaRoleBanner           AriaRole = "banner"
	AriaRoleBlockquote       AriaRole = "blockquote"
	AriaRoleButton           AriaRole = "button"
	AriaRoleCaption          AriaRole = "caption"
	AriaRoleCell             AriaRole = "cell"
	AriaRoleCheckbox         AriaRole = "checkbox"
	AriaRoleCode             AriaRole = "code"
	AriaRoleColumnHeader     AriaRole = "columnheader"
	AriaRoleCombobox         AriaRole = "combobox"
	AriaRoleComplementary    AriaRole = "complementary"
	AriaRoleContentInfo      AriaRole = "contentinfo"
	AriaRoleDefinition       AriaRole = "definition"
	AriaRoleDeletion         AriaRole = "deletion"
	AriaRoleDialog           AriaRole = "dialog"
	AriaRoleDirectory        AriaRole = "directory"
	AriaRoleDocument         AriaRole = "document"
	AriaRoleEmphasis         AriaRole = "emphasis"
	AriaRoleFeed             AriaRole = "feed"
	AriaRoleFigure           AriaRole = "figure"
	AriaRoleForm             AriaRole = "form"
	AriaRoleGeneric          AriaRole = "generic"
	AriaRoleGrid             AriaRole = "grid"
	AriaRoleGridCell         AriaRole = "gridcell"
	AriaRoleGroup            AriaRole = "group"
	AriaRoleHeading          AriaRole = "heading"
	AriaRoleImg              AriaRole = "img"
	AriaRoleInsertion        AriaRole = "insertion"
	AriaRoleLink             AriaRole = "link"
	AriaRoleList             AriaRole = "list"
	AriaRoleListBox          AriaRole = "listbox"
	AriaRoleListItem         AriaRole = "listitem"
	AriaRoleLog              AriaRole = "log"
	AriaRoleMain             AriaRole = "main"
	AriaRoleMarquee          AriaRole = "marquee"
	AriaRoleMath             AriaRole = "math"
	AriaRoleMenu             AriaRole = "menu"
	AriaRoleMenuBar          AriaRole = "menubar"
	AriaRoleMenuItem         AriaRole = "menuitem"
	AriaRoleMenuItemCheckbox AriaRole = "menuitemcheckbox"
	AriaRoleMenuItemRadio    AriaRole = "menuitemradio"
	AriaRoleMeter            AriaRole = "meter"
	AriaRoleNavigation       AriaRole = "navigation"
	AriaRoleNone             AriaRole = "none"
	AriaRoleNote             AriaRole = "note"
	AriaRoleOption           AriaRole = "option"
	AriaRoleParagraph        AriaRole = "paragraph"
	AriaRolePresentation     AriaRole = "presentation"
	AriaRoleProgressBar      AriaRole = "progressbar"
	AriaRoleRadio            AriaRole = "radio"
	AriaRoleRadioGroup       AriaRole = "radiogroup"
	AriaRoleRegion           AriaRole = "region"
	AriaRoleRow              AriaRole = "row"
	AriaRoleRowGroup         AriaRole = "rowgroup"
	AriaRoleRowHeader        AriaRole = "rowheader"
	AriaRoleScrollBar        AriaRole = "scrollbar"
	AriaRoleSearch           AriaRole = "search"
	AriaRoleSearchBox        AriaRole = "searchbox"
	AriaRoleSeparator        AriaRole = "separator"
	AriaRoleSlider           AriaRole = "slider"
	AriaRoleSpinButton       AriaRole = "spinbutton"
	AriaRoleStatus           AriaRole = "status"
	AriaRoleStrong           AriaRole = "strong"
	AriaRoleSubscript        AriaRole = "subscript"
	AriaRoleSuperscript      AriaRole = "superscript"
	AriaRoleSwitch           AriaRole = "switch"
	AriaRoleTab              AriaRole = "tab"
	AriaRoleTable            AriaRole = "table"
	AriaRoleTabList          AriaRole = "tablist"
	AriaRoleTabPanel         AriaRole = "tabpanel"
	AriaRoleTerm             AriaRole = "term"
	AriaRoleTextbox          AriaRole = "textbox"
	AriaRoleTime             AriaRole = "time"
	AriaRoleTimer            AriaRole = "timer"
	AriaRoleToolBar          AriaRole = "toolbar"
	AriaRoleToolTip          AriaRole = "tooltip"
	AriaRoleTree             AriaRole = "tree"
	AriaRoleTreeGrid         AriaRole = "treegrid"
	AriaRoleTreeItem         AriaRole = "treeitem"
)

// GetByRoleOptions contains optional filters for GetByRole.
type GetByRoleOptions struct {
	Name *string
	// Exact controls whether Name matching is case-sensitive and exact. Does not affect role matching.
	Exact         *bool
	Checked       *bool
	Disabled      *bool
	Pressed       *bool
	Selected      *bool
	Expanded      *bool
	Level         *int
	IncludeHidden *bool
}

// GetByTextOptions contains optional filters for GetByText.
type GetByTextOptions struct {
	Exact *bool
}

// GetByLabelOptions contains optional filters for GetByLabel.
type GetByLabelOptions struct {
	Exact *bool
}

// GetByPlaceholderOptions contains optional filters for GetByPlaceholder.
type GetByPlaceholderOptions struct {
	Exact *bool
}

// GetByAltTextOptions contains optional filters for GetByAltText.
type GetByAltTextOptions struct {
	Exact *bool
}

// GetByTitleOptions contains optional filters for GetByTitle.
type GetByTitleOptions struct {
	Exact *bool
}

// Locator provides a way to find and interact with elements on a page.
// Locator objects are strict: they enforce that exactly one element matches
// the selector when performing actions (unless Count is used).
type Locator struct {
	frame    ChannelOwner
	selector string
}

// escapeForSelector escapes a string value for embedding inside Playwright internal selector
// attribute expressions (e.g. [name="..."]). Backslashes and double-quotes are escaped.
func escapeForSelector(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

// escapeForAttributeSelector wraps a value in double quotes for use in attribute selector expressions,
// escaping backslashes and double-quotes.
func escapeForAttributeSelector(value string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(value) + `"`
}

// locatorFromRole builds a locator using the internal:role selector.
func locatorFromRole(frame ChannelOwner, role AriaRole, opts *GetByRoleOptions) *Locator {
	var parts []string
	parts = append(parts, "internal:role="+string(role))

	if opts != nil {
		if opts.Name != nil {
			exact := false
			if opts.Exact != nil {
				exact = *opts.Exact
			}
			escaped := escapeForSelector(*opts.Name)
			if exact {
				parts = append(parts, fmt.Sprintf(`[name="%s"]`, escaped))
			} else {
				parts = append(parts, fmt.Sprintf(`[name="%s"i]`, escaped))
			}
		}
		if opts.Checked != nil {
			val := "[checked=false]"
			if *opts.Checked {
				val = "[checked]"
			}
			parts = append(parts, val)
		}
		if opts.Disabled != nil {
			val := "[disabled=false]"
			if *opts.Disabled {
				val = "[disabled]"
			}
			parts = append(parts, val)
		}
		if opts.Pressed != nil {
			val := "[pressed=false]"
			if *opts.Pressed {
				val = "[pressed]"
			}
			parts = append(parts, val)
		}
		if opts.Selected != nil {
			val := "[selected=false]"
			if *opts.Selected {
				val = "[selected]"
			}
			parts = append(parts, val)
		}
		if opts.Expanded != nil {
			val := "[expanded=false]"
			if *opts.Expanded {
				val = "[expanded]"
			}
			parts = append(parts, val)
		}
		if opts.Level != nil {
			parts = append(parts, fmt.Sprintf("[level=%d]", *opts.Level))
		}
		if opts.IncludeHidden != nil && *opts.IncludeHidden {
			parts = append(parts, "[include-hidden]")
		}
	}

	selector := strings.Join(parts, "")
	return &Locator{frame: frame, selector: selector}
}

// locatorFromText builds a locator using the internal:text selector.
// For partial (non-exact) matching, it uses case-insensitive mode with the "i" flag.
// For exact matching, the flag is omitted for case-sensitive full-string matching.
func locatorFromText(frame ChannelOwner, text string, exact bool) *Locator {
	escaped := escapeForSelector(text)
	var selector string
	if exact {
		// Exact match: no case-insensitive flag — matches the full text exactly
		selector = fmt.Sprintf(`internal:text="%s"`, escaped)
	} else {
		// Partial match: case-insensitive with "i" flag
		selector = fmt.Sprintf(`internal:text="%s"i`, escaped)
	}
	return &Locator{frame: frame, selector: selector}
}

// locatorFromLabel builds a locator using the internal:label selector.
func locatorFromLabel(frame ChannelOwner, text string, exact bool) *Locator {
	escaped := escapeForSelector(text)
	var selector string
	if exact {
		selector = fmt.Sprintf(`internal:label="%s"`, escaped)
	} else {
		selector = fmt.Sprintf(`internal:label="%s"i`, escaped)
	}
	return &Locator{frame: frame, selector: selector}
}

// locatorFromAttr builds a locator using the internal:attr selector.
func locatorFromAttr(frame ChannelOwner, attr, value string, exact bool) *Locator {
	escaped := escapeForSelector(value)
	var selector string
	if exact {
		selector = fmt.Sprintf(`internal:attr=[%s="%s"]`, attr, escaped)
	} else {
		selector = fmt.Sprintf(`internal:attr=[%s="%s"i]`, attr, escaped)
	}
	return &Locator{frame: frame, selector: selector}
}

// IsVisible returns true if the element is visible.
func (l *Locator) IsVisible(ctx context.Context) (bool, error) {
	req := protocol.FrameIsVisibleRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
	}
	result, err := l.frame.SendMessageRequest(ctx, "isVisible", req)
	if err != nil {
		return false, fmt.Errorf("locator.isVisible failed: %w", err)
	}
	var resp protocol.FrameIsVisibleResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isVisible response: %w", err)
	}
	return resp.Value, nil
}

// defaultLocatorTimeout is the default timeout for locator actions in milliseconds.
// Using 30000ms matches Playwright's default action timeout.
const defaultLocatorTimeout = defaultActionTimeoutMs

// locatorTimeoutRequest is the shared wire format for locator methods (isEnabled, isDisabled,
// isEditable, isChecked, innerText, innerHTML, inputValue) that require selector, strict, and timeout.
type locatorTimeoutRequest struct {
	Selector string  `json:"selector"`
	Strict   *bool   `json:"strict,omitempty"`
	Timeout  float64 `json:"timeout"`
}

// locatorGetAttributeRequest is the wire format for getAttribute, including the required timeout.
type locatorGetAttributeRequest struct {
	Selector string  `json:"selector"`
	Name     string  `json:"name"`
	Strict   *bool   `json:"strict,omitempty"`
	Timeout  float64 `json:"timeout"`
}

// locatorFillRequest is the wire format for fill, including the required timeout.
type locatorFillRequest struct {
	Selector string  `json:"selector"`
	Value    string  `json:"value"`
	Force    *bool   `json:"force,omitempty"`
	Strict   *bool   `json:"strict,omitempty"`
	Timeout  float64 `json:"timeout"`
}

// LocatorClickOptions specifies optional parameters for Locator.Click.
type LocatorClickOptions struct {
	Button      *string  // "left" (default), "right", or "middle"
	ClickCount  *int     // number of clicks (default 1)
	Delay       *float64 // ms between mousedown and mouseup
	Force       *bool    // bypass actionability checks
	Modifiers   []string
	NoWaitAfter *bool
	Position    *struct{ X, Y float64 } // click offset relative to element top-left
	Steps       *int
	Scroll      *string
	Timeout     *float64 // override default timeout in ms
	Trial       *bool    // perform actionability checks without clicking
}

// locatorClickRequest is the wire format for click, including required fields.
type locatorClickRequest struct {
	Selector    string   `json:"selector"`
	Button      string   `json:"button"`
	ClickCount  *int     `json:"clickCount,omitempty"`
	Delay       *float64 `json:"delay,omitempty"`
	Force       *bool    `json:"force,omitempty"`
	Modifiers   []any    `json:"modifiers"`
	NoWaitAfter *bool    `json:"noWaitAfter,omitempty"`
	Position    *struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position,omitempty"`
	Scroll  *string `json:"scroll,omitempty"`
	Strict  *bool   `json:"strict,omitempty"`
	Steps   *int    `json:"steps,omitempty"`
	Timeout float64 `json:"timeout"`
	Trial   *bool   `json:"trial,omitempty"`
}

// locatorHoverRequest is the wire format for hover.
// Modifiers and Scroll must be present (not omitted) per the protocol.
type locatorHoverRequest struct {
	Modifiers []any   `json:"modifiers"`
	Scroll    any     `json:"scroll"`
	Selector  string  `json:"selector"`
	Strict    *bool   `json:"strict,omitempty"`
	Timeout   float64 `json:"timeout"`
}

// locatorCheckRequest is the wire format for check/uncheck.
// Scroll must be present (not omitted) per the protocol.
type locatorCheckRequest struct {
	Scroll   any     `json:"scroll"`
	Selector string  `json:"selector"`
	Strict   *bool   `json:"strict,omitempty"`
	Timeout  float64 `json:"timeout"`
}

// locatorPressRequest is the wire format for press.
type locatorPressRequest struct {
	Key      string  `json:"key"`
	Selector string  `json:"selector"`
	Strict   *bool   `json:"strict,omitempty"`
	Timeout  float64 `json:"timeout"`
}

// locatorSelectOptionRequest is the wire format for selectOption.
type locatorSelectOptionRequest struct {
	Elements []any   `json:"elements"`
	Options  []any   `json:"options"`
	Selector string  `json:"selector"`
	Strict   *bool   `json:"strict,omitempty"`
	Timeout  float64 `json:"timeout"`
}

// IsEnabled returns true if the element is enabled.
func (l *Locator) IsEnabled(ctx context.Context) (bool, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "isEnabled", req)
	if err != nil {
		return false, fmt.Errorf("locator.isEnabled failed: %w", err)
	}
	var resp protocol.FrameIsEnabledResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isEnabled response: %w", err)
	}
	return resp.Value, nil
}

// IsDisabled returns true if the element is disabled.
func (l *Locator) IsDisabled(ctx context.Context) (bool, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "isDisabled", req)
	if err != nil {
		return false, fmt.Errorf("locator.isDisabled failed: %w", err)
	}
	var resp protocol.FrameIsDisabledResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isDisabled response: %w", err)
	}
	return resp.Value, nil
}

// IsEditable returns true if the element is editable.
func (l *Locator) IsEditable(ctx context.Context) (bool, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "isEditable", req)
	if err != nil {
		return false, fmt.Errorf("locator.isEditable failed: %w", err)
	}
	var resp protocol.FrameIsEditableResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isEditable response: %w", err)
	}
	return resp.Value, nil
}

// IsChecked returns true if the element is checked (checkbox or radio button).
func (l *Locator) IsChecked(ctx context.Context) (bool, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "isChecked", req)
	if err != nil {
		return false, fmt.Errorf("locator.isChecked failed: %w", err)
	}
	var resp protocol.FrameIsCheckedResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return false, fmt.Errorf("failed to parse isChecked response: %w", err)
	}
	return resp.Value, nil
}

// InnerText returns the element's inner text.
func (l *Locator) InnerText(ctx context.Context) (string, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "innerText", req)
	if err != nil {
		return "", fmt.Errorf("locator.innerText failed: %w", err)
	}
	var resp protocol.FrameInnerTextResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse innerText response: %w", err)
	}
	return resp.Value, nil
}

// InnerHTML returns the element's inner HTML.
func (l *Locator) InnerHTML(ctx context.Context) (string, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "innerHTML", req)
	if err != nil {
		return "", fmt.Errorf("locator.innerHTML failed: %w", err)
	}
	var resp protocol.FrameInnerHTMLResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse innerHTML response: %w", err)
	}
	return resp.Value, nil
}

// InputValue returns the value of an input, textarea, or select element.
func (l *Locator) InputValue(ctx context.Context) (string, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "inputValue", req)
	if err != nil {
		return "", fmt.Errorf("locator.inputValue failed: %w", err)
	}
	var resp protocol.FrameInputValueResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse inputValue response: %w", err)
	}
	return resp.Value, nil
}

// GetAttribute returns the value of an attribute on the element.
// Returns nil if the attribute does not exist.
func (l *Locator) GetAttribute(ctx context.Context, name string) (*string, error) {
	req := locatorGetAttributeRequest{
		Selector: l.selector,
		Name:     name,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "getAttribute", req)
	if err != nil {
		return nil, fmt.Errorf("locator.getAttribute failed: %w", err)
	}
	var resp protocol.FrameGetAttributeResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse getAttribute response: %w", err)
	}
	return resp.Value, nil
}

// Count returns the number of elements matching the locator's selector.
func (l *Locator) Count(ctx context.Context) (int, error) {
	req := protocol.FrameQueryCountRequest{
		Selector: l.selector,
	}
	result, err := l.frame.SendMessageRequest(ctx, "queryCount", req)
	if err != nil {
		return 0, fmt.Errorf("locator.count failed: %w", err)
	}
	var resp protocol.FrameQueryCountResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return 0, fmt.Errorf("failed to parse queryCount response: %w", err)
	}
	return resp.Value, nil
}

// Click clicks the element.
func (l *Locator) Click(ctx context.Context, opts ...*LocatorClickOptions) error {
	req := locatorClickRequest{
		Selector:  l.selector,
		Button:    "left",
		Modifiers: []any{},
		Strict:    protocol.Bool(true),
		Timeout:   defaultLocatorTimeout,
	}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.Button != nil {
			req.Button = *o.Button
		}
		req.ClickCount = o.ClickCount
		req.Delay = o.Delay
		req.Force = o.Force
		for _, mod := range o.Modifiers {
			req.Modifiers = append(req.Modifiers, mod)
		}
		req.NoWaitAfter = o.NoWaitAfter
		if o.Position != nil {
			req.Position = &struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			}{X: o.Position.X, Y: o.Position.Y}
		}
		req.Steps = o.Steps
		req.Scroll = o.Scroll
		if o.Timeout != nil {
			req.Timeout = *o.Timeout
		}
		req.Trial = o.Trial
	}
	_, err := l.frame.SendMessageRequest(ctx, "click", req)
	if err != nil {
		return fmt.Errorf("locator.click failed: %w", err)
	}
	return nil
}

// DblClick performs a double-click on the matched element.
func (l *Locator) DblClick(ctx context.Context) error {
	req := struct {
		Button    string  `json:"button"`
		Modifiers []any   `json:"modifiers"`
		Scroll    any     `json:"scroll"`
		Selector  string  `json:"selector"`
		Strict    *bool   `json:"strict,omitempty"`
		Timeout   float64 `json:"timeout"`
	}{
		Button:    "left",
		Modifiers: []any{},
		Scroll:    nil,
		Selector:  l.selector,
		Strict:    protocol.Bool(true),
		Timeout:   defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "dblclick", req)
	if err != nil {
		return fmt.Errorf("locator.dblclick failed: %w", err)
	}
	return nil
}

// Fill sets the value of an input element.
func (l *Locator) Fill(ctx context.Context, value string) error {
	req := locatorFillRequest{
		Selector: l.selector,
		Value:    value,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "fill", req)
	if err != nil {
		return fmt.Errorf("locator.fill failed: %w", err)
	}
	return nil
}

// Selector returns the CSS/internal selector string used by this locator.
func (l *Locator) Selector() string {
	return l.selector
}

// Hover moves the mouse pointer over the matched element.
func (l *Locator) Hover(ctx context.Context) error {
	req := locatorHoverRequest{
		Modifiers: []any{},
		Scroll:    nil,
		Selector:  l.selector,
		Strict:    protocol.Bool(true),
		Timeout:   defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "hover", req)
	if err != nil {
		return fmt.Errorf("locator.hover failed: %w", err)
	}
	return nil
}

// Check checks a checkbox or radio element.
func (l *Locator) Check(ctx context.Context) error {
	req := locatorCheckRequest{
		Scroll:   nil,
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "check", req)
	if err != nil {
		return fmt.Errorf("locator.check failed: %w", err)
	}
	return nil
}

// Uncheck unchecks a checkbox element.
func (l *Locator) Uncheck(ctx context.Context) error {
	req := locatorCheckRequest{
		Scroll:   nil,
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "uncheck", req)
	if err != nil {
		return fmt.Errorf("locator.uncheck failed: %w", err)
	}
	return nil
}

// Press presses a key on the focused element matching the locator.
func (l *Locator) Press(ctx context.Context, key string) error {
	req := locatorPressRequest{
		Key:      key,
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "press", req)
	if err != nil {
		return fmt.Errorf("locator.press failed: %w", err)
	}
	return nil
}

// SelectOption selects one or more options in a <select> element by value.
// Returns the list of values that were selected.
func (l *Locator) SelectOption(ctx context.Context, values ...string) ([]string, error) {
	opts := make([]any, len(values))
	for i, v := range values {
		opts[i] = map[string]string{"value": v}
	}
	req := locatorSelectOptionRequest{
		Elements: []any{},
		Options:  opts,
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "selectOption", req)
	if err != nil {
		return nil, fmt.Errorf("locator.selectOption failed: %w", err)
	}
	var resp protocol.FrameSelectOptionResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse selectOption response: %w", err)
	}
	return resp.Values, nil
}

// Nth returns a locator that narrows the match to the element at the given index (0-based).
// Use -1 to select the last element.
func (l *Locator) Nth(index int) *Locator {
	return &Locator{frame: l.frame, selector: fmt.Sprintf("%s >> nth=%d", l.selector, index)}
}

// First returns a locator that narrows the match to the first matching element.
func (l *Locator) First() *Locator {
	return l.Nth(0)
}

// Locator returns a new Locator scoped to descendants of this locator that match selector.
func (l *Locator) Locator(selector string) *Locator {
	return &Locator{frame: l.frame, selector: l.selector + " >> " + selector}
}

// Last returns a locator that narrows the match to the last matching element.
func (l *Locator) Last() *Locator {
	return l.Nth(-1)
}

// And returns a new Locator that matches elements satisfying both this locator
// and the given locator simultaneously.
func (l *Locator) And(other *Locator) *Locator {
	return &Locator{
		frame:    l.frame,
		selector: fmt.Sprintf("%s >> internal:and=%s", l.selector, subSelectorArg(other.selector)),
	}
}

// Or returns a new Locator that matches elements satisfying either this locator
// or the given locator.
func (l *Locator) Or(other *Locator) *Locator {
	return &Locator{
		frame:    l.frame,
		selector: fmt.Sprintf("%s >> internal:or=%s", l.selector, subSelectorArg(other.selector)),
	}
}

// LocatorFilterOptions contains options for narrowing a Locator's match set.
type LocatorFilterOptions struct {
	// HasText restricts to elements whose text content includes this string (case-insensitive substring).
	HasText *string
	// HasTextRegex restricts to elements whose text content matches this regular expression.
	HasTextRegex *regexp.Regexp
	// HasNotText restricts to elements whose text content does NOT include this string.
	HasNotText *string
	// HasNotTextRegex restricts to elements whose text content does NOT match this regular expression.
	HasNotTextRegex *regexp.Regexp
	// Has restricts to elements that contain a descendant matching the given Locator.
	Has *Locator
	// HasNot restricts to elements that do NOT contain a descendant matching the given Locator.
	HasNot *Locator
	// Visible, when true, restricts to visible elements; when false, to hidden elements.
	Visible *bool
}

// regexToSelectorPattern converts a Go regexp to a Playwright /pattern/flags string.
func regexToSelectorPattern(re *regexp.Regexp) string {
	src, flags := extractRegexpInfo(re)
	return fmt.Sprintf("/%s/%s", src, flags)
}

// Filter narrows the Locator to elements that satisfy the given filter options.
// Multiple options are applied cumulatively (AND semantics).
func (l *Locator) Filter(opts *LocatorFilterOptions) *Locator {
	if opts == nil {
		return &Locator{frame: l.frame, selector: l.selector}
	}
	sel := l.selector
	if opts.HasText != nil {
		sel = fmt.Sprintf(`%s >> internal:has-text="%s"i`, sel, escapeForSelector(*opts.HasText))
	}
	if opts.HasTextRegex != nil {
		sel = fmt.Sprintf(`%s >> internal:has-text=%s`, sel, regexToSelectorPattern(opts.HasTextRegex))
	}
	if opts.HasNotText != nil {
		sel = fmt.Sprintf(`%s >> internal:has-not-text="%s"i`, sel, escapeForSelector(*opts.HasNotText))
	}
	if opts.HasNotTextRegex != nil {
		sel = fmt.Sprintf(`%s >> internal:has-not-text=%s`, sel, regexToSelectorPattern(opts.HasNotTextRegex))
	}
	if opts.Has != nil {
		sel = fmt.Sprintf("%s >> internal:has=%s", sel, subSelectorArg(opts.Has.selector))
	}
	if opts.HasNot != nil {
		sel = fmt.Sprintf("%s >> internal:has-not=%s", sel, subSelectorArg(opts.HasNot.selector))
	}
	if opts.Visible != nil {
		// visible=true → visible elements only; visible=false → hidden elements only.
		sel = fmt.Sprintf("%s >> visible=%v", sel, *opts.Visible)
	}
	return &Locator{frame: l.frame, selector: sel}
}

// FrameLocator scopes subsequent locator lookups within the first iframe matching selector.
func (l *Locator) FrameLocator(selector string) *FrameLocator {
	return &FrameLocator{
		frame:    l.frame,
		selector: fmt.Sprintf("%s >> %s >> internal:control=enter-frame", l.selector, selector),
	}
}

// SetInputFiles sets the files on the file input element matched by this locator.
// paths must be absolute local filesystem paths.
func (l *Locator) SetInputFiles(ctx context.Context, paths []string) error {
	req := map[string]any{
		"localPaths": paths,
		"selector":   l.selector,
		"strict":     true,
		"timeout":    defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "setInputFiles", req)
	if err != nil {
		return fmt.Errorf("locator.setInputFiles failed: %w", err)
	}
	return nil
}

// SetInputFilesPayload sets in-memory file payloads on the file input matched by this locator.
func (l *Locator) SetInputFilesPayload(ctx context.Context, payloads []FilePayload) error {
	wirePayloads := make([]map[string]any, len(payloads))
	for i, p := range payloads {
		wirePayloads[i] = map[string]any{
			"name":     p.Name,
			"mimeType": p.MimeType,
			"buffer":   base64.StdEncoding.EncodeToString(p.Buffer),
		}
	}
	req := map[string]any{
		"payloads": wirePayloads,
		"selector": l.selector,
		"strict":   true,
		"timeout":  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "setInputFiles", req)
	if err != nil {
		return fmt.Errorf("locator.setInputFilesPayload failed: %w", err)
	}
	return nil
}

// Evaluate evaluates an expression on the first element matching this locator and returns the result.
func (l *Locator) Evaluate(ctx context.Context, expression string, arg ...any) (any, error) {
	var inputArg any
	if len(arg) > 0 {
		inputArg = arg[0]
	}
	result, err := l.frame.SendMessageRequest(ctx, "evalOnSelector", map[string]any{
		"selector":   l.selector,
		"expression": expression,
		"arg":        serializeArgument(inputArg),
		"strict":     true,
	})
	if err != nil {
		return nil, fmt.Errorf("locator.evaluate failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("locator.evaluate: parse response failed: %w", err)
	}
	return deserializeValue(resp.Value)
}

// EvaluateAll evaluates an expression on all elements matching this locator and returns the result.
func (l *Locator) EvaluateAll(ctx context.Context, expression string, arg ...any) (any, error) {
	var inputArg any
	if len(arg) > 0 {
		inputArg = arg[0]
	}
	result, err := l.frame.SendMessageRequest(ctx, "evalOnSelectorAll", map[string]any{
		"selector":   l.selector,
		"expression": expression,
		"arg":        serializeArgument(inputArg),
	})
	if err != nil {
		return nil, fmt.Errorf("locator.evaluateAll failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("locator.evaluateAll: parse response failed: %w", err)
	}
	return deserializeValue(resp.Value)
}

// TextContent returns the text content of the element matched by this locator.
// Returns nil if the element has no text content (e.g. a void element).
func (l *Locator) TextContent(ctx context.Context) (*string, error) {
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	result, err := l.frame.SendMessageRequest(ctx, "textContent", req)
	if err != nil {
		return nil, fmt.Errorf("locator.textContent failed: %w", err)
	}
	var resp protocol.FrameTextContentResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse textContent response: %w", err)
	}
	return resp.Value, nil
}

// AllTextContents returns the text content of all elements matching this locator.
// Uses evalOnSelectorAll IPC to support compound selectors (e.g. "article >> p").
func (l *Locator) AllTextContents(ctx context.Context) ([]string, error) {
	return l.allStringsViaEvalOnSelectorAll(ctx, "els => Array.from(els).map(el => el.textContent ?? '')", "allTextContents")
}

// AllInnerTexts returns the inner text of all elements matching this locator.
// Uses evalOnSelectorAll IPC to support compound selectors (e.g. "article >> p").
func (l *Locator) AllInnerTexts(ctx context.Context) ([]string, error) {
	return l.allStringsViaEvalOnSelectorAll(ctx, "els => Array.from(els).map(el => el.innerText ?? '')", "allInnerTexts")
}

// allStringsViaEvalOnSelectorAll uses the evalOnSelectorAll IPC method which supports
// Playwright's full selector syntax including >> compound selectors.
func (l *Locator) allStringsViaEvalOnSelectorAll(ctx context.Context, expression, opName string) ([]string, error) {
	result, err := l.frame.SendMessageRequest(ctx, "evalOnSelectorAll", map[string]any{
		"selector":   l.selector,
		"expression": expression,
		"arg":        serializeArgument(nil),
	})
	if err != nil {
		return nil, fmt.Errorf("locator.%s failed: %w", opName, err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("locator.%s: parse response failed: %w", opName, err)
	}
	val, err := deserializeValue(resp.Value)
	if err != nil {
		return nil, fmt.Errorf("locator.%s: deserialize failed: %w", opName, err)
	}
	raw, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("locator.%s: unexpected type %T", opName, val)
	}
	texts := make([]string, len(raw))
	for i, v := range raw {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("locator.%s: element %d has unexpected type %T", opName, i, v)
		}
		texts[i] = s
	}
	return texts, nil
}

// SetChecked sets the checked state of a checkbox or radio element.
// Uses check/uncheck IPC since Frame.setChecked is not available in Playwright v1.61.1.
func (l *Locator) SetChecked(ctx context.Context, checked bool) error {
	method := "uncheck"
	if checked {
		method = "check"
	}
	req := locatorTimeoutRequest{
		Selector: l.selector,
		Strict:   protocol.Bool(true),
		Timeout:  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, method, req)
	if err != nil {
		return fmt.Errorf("locator.setChecked(%v) failed: %w", checked, err)
	}
	return nil
}

// Focus focuses the element matched by this locator.
func (l *Locator) Focus(ctx context.Context) error {
	req := map[string]any{
		"selector": l.selector,
		"strict":   true,
		"timeout":  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "focus", req)
	if err != nil {
		return fmt.Errorf("locator.focus failed: %w", err)
	}
	return nil
}

// Clear clears the value of an editable element (input, textarea) matched by this locator.
func (l *Locator) Clear(ctx context.Context) error {
	_, err := l.frame.SendMessageRequest(ctx, "fill", map[string]any{
		"value":    "",
		"selector": l.selector,
		"strict":   true,
		"timeout":  defaultLocatorTimeout,
	})
	if err != nil {
		return fmt.Errorf("locator.clear failed: %w", err)
	}
	return nil
}

// PressSequentially types text into the element one character at a time,
// firing keydown/keyup/input events for each character.
// Uses the "type" IPC (Frame.type) which is the v1.61.1 equivalent.
func (l *Locator) PressSequentially(ctx context.Context, text string) error {
	req := map[string]any{
		"text":     text,
		"selector": l.selector,
		"strict":   true,
		"timeout":  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "type", req)
	if err != nil {
		return fmt.Errorf("locator.pressSequentially failed: %w", err)
	}
	return nil
}

// WaitFor waits for the element matched by this locator to satisfy a given state.
// state can be "attached", "detached", "visible", "hidden". Defaults to "visible".
func (l *Locator) WaitFor(ctx context.Context, state ...string) error {
	s := "visible"
	if len(state) > 0 && state[0] != "" {
		s = state[0]
	}
	req := map[string]any{
		"selector": l.selector,
		"strict":   true,
		"state":    s,
		"timeout":  defaultLocatorTimeout,
	}
	_, err := l.frame.SendMessageRequest(ctx, "waitForSelector", req)
	if err != nil {
		return fmt.Errorf("locator.waitFor(%q) failed: %w", s, err)
	}
	return nil
}

// LocatorWaitForFunctionOptions configures Locator.WaitForFunction behavior.
type LocatorWaitForFunctionOptions struct {
	// Timeout in milliseconds. 0 uses the server default; negative disables the timeout.
	Timeout float64
}

// WaitForFunction polls the given JavaScript expression until it returns a truthy value,
// passing the resolved element as the first argument. An optional extra arg is passed as the second.
// The locator must resolve to exactly one element (strict mode).
func (l *Locator) WaitForFunction(ctx context.Context, expression string, arg any, opts ...*LocatorWaitForFunctionOptions) error {
	req := map[string]any{
		"selector":   l.selector,
		"strict":     true,
		"expression": expression,
		"arg":        serializeArgument(arg),
		"timeout":    defaultLocatorTimeout,
	}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.Timeout < 0 {
			req["timeout"] = 0.0
		} else if o.Timeout > 0 {
			req["timeout"] = o.Timeout
		}
	}
	_, err := l.frame.SendMessageRequest(ctx, "waitForFunction", req)
	if err != nil {
		return fmt.Errorf("locator.waitForFunction failed: %w", err)
	}
	return nil
}

// All returns all elements matching this locator as a slice of individual Locators.
func (l *Locator) All(ctx context.Context) ([]*Locator, error) {
	count, err := l.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("locator.all: count failed: %w", err)
	}
	result := make([]*Locator, count)
	for i := range count {
		result[i] = l.Nth(i)
	}
	return result, nil
}

// AriaSnapshot returns the ARIA snapshot of the element matched by this locator.
// The snapshot is a YAML string representing the accessibility tree.
// DragTo performs a drag-and-drop from this locator's element to the target locator's element.
func (l *Locator) DragTo(ctx context.Context, target *Locator) error {
	_, err := l.frame.SendMessageRequest(ctx, "dragAndDrop", map[string]any{
		"source":  l.selector,
		"target":  target.selector,
		"timeout": defaultActionTimeoutMs,
		"strict":  false,
	})
	if err != nil {
		return fmt.Errorf("locator.dragTo failed: %w", err)
	}
	return nil
}

// LocatorAriaSnapshotOptions configures Locator.AriaSnapshot.
type LocatorAriaSnapshotOptions struct {
	// Mode controls snapshot generation: "default" (standard ARIA tree) or "ai" (AI-optimized with refs/hints).
	Mode *string
	// Depth limits how many levels deep the snapshot traverses (AI mode only).
	Depth *int
}

func (l *Locator) AriaSnapshot(ctx context.Context, opts ...*LocatorAriaSnapshotOptions) (string, error) {
	req := map[string]any{
		"mode":     "default",
		"selector": l.selector,
		"timeout":  defaultLocatorTimeout,
	}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.Mode != nil {
			req["mode"] = *o.Mode
		}
		if o.Depth != nil {
			req["depth"] = *o.Depth
		}
	}
	result, err := l.frame.SendMessageRequest(ctx, "ariaSnapshot", req)
	if err != nil {
		return "", fmt.Errorf("locator.ariaSnapshot failed: %w", err)
	}
	var resp protocol.FrameAriaSnapshotResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse ariaSnapshot response: %w", err)
	}
	return resp.Snapshot, nil
}

// Highlight visually highlights the element matched by this locator in the browser.
// Useful during debugging and authoring. The highlight is temporary and browser-internal.
func (l *Locator) Highlight(ctx context.Context) error {
	_, err := l.frame.SendMessageRequest(ctx, "highlight", map[string]any{
		"selector": l.selector,
	})
	if err != nil {
		return fmt.Errorf("locator.highlight failed: %w", err)
	}
	return nil
}

// Tap performs a tap gesture on the element matched by this locator.
func (l *Locator) Tap(ctx context.Context) error {
	req := protocol.FrameTapRequest{
		Modifiers: []any{},
		Scroll:    nil,
		Selector:  l.selector,
		Strict:    protocol.Bool(true),
	}
	_, err := l.frame.SendMessageRequest(ctx, "tap", req)
	if err != nil {
		return fmt.Errorf("locator.tap failed: %w", err)
	}
	return nil
}

// ScrollIntoViewIfNeeded scrolls the element matched by this locator into the viewport if needed.
// Uses a two-step querySelector + ElementHandle call because the Playwright v1.61.1 wire protocol
// does not expose Frame.scrollIntoViewIfNeeded directly.
func (l *Locator) ScrollIntoViewIfNeeded(ctx context.Context) error {
	result, err := l.frame.SendMessageRequest(ctx, "querySelector", map[string]any{
		"selector": l.selector,
		"strict":   true,
	})
	if err != nil {
		return fmt.Errorf("locator.scrollIntoViewIfNeeded: querySelector failed: %w", err)
	}
	var qresp struct {
		Element *struct {
			Guid string `json:"guid"`
		} `json:"element,omitempty"`
	}
	if err := json.Unmarshal(result, &qresp); err != nil {
		return fmt.Errorf("locator.scrollIntoViewIfNeeded: parse querySelector failed: %w", err)
	}
	if qresp.Element == nil {
		return fmt.Errorf("locator.scrollIntoViewIfNeeded: element not found for selector %q", l.selector)
	}
	el := elementHandleFromGUID(l.frame, qresp.Element.Guid)
	defer func() { _ = el.Dispose(context.Background()) }()
	return el.ScrollIntoViewIfNeeded(ctx)
}

// BoundingBox returns the bounding box of the element matched by this locator,
// or nil if the element is not visible.
// Uses a two-step querySelector + ElementHandle call because the Playwright v1.61.1 wire protocol
// does not expose Frame.boundingBox directly.
func (l *Locator) BoundingBox(ctx context.Context) (*BoundingBox, error) {
	result, err := l.frame.SendMessageRequest(ctx, "querySelector", map[string]any{
		"selector": l.selector,
		"strict":   true,
	})
	if err != nil {
		return nil, fmt.Errorf("locator.boundingBox: querySelector failed: %w", err)
	}
	var qresp struct {
		Element *struct {
			Guid string `json:"guid"`
		} `json:"element,omitempty"`
	}
	if err := json.Unmarshal(result, &qresp); err != nil {
		return nil, fmt.Errorf("locator.boundingBox: parse querySelector failed: %w", err)
	}
	if qresp.Element == nil {
		return nil, fmt.Errorf("locator.boundingBox: element not found for selector %q", l.selector)
	}
	el := elementHandleFromGUID(l.frame, qresp.Element.Guid)
	defer func() { _ = el.Dispose(context.Background()) }()
	return el.BoundingBox(ctx)
}

// Screenshot captures a screenshot of the element matched by this locator.
// Returns raw image bytes (PNG by default). Accepts the same ScreenshotOptions as Page.Screenshot.
func (l *Locator) Screenshot(ctx context.Context, opts ...*ScreenshotOptions) ([]byte, error) {
	// Frame.screenshot is not available in Playwright v1.61.1; resolve to ElementHandle first.
	result, err := l.frame.SendMessageRequest(ctx, "querySelector", map[string]any{
		"selector": l.selector,
		"strict":   true,
	})
	if err != nil {
		return nil, fmt.Errorf("locator.screenshot: querySelector failed: %w", err)
	}
	var qresp struct {
		Element *struct {
			Guid string `json:"guid"`
		} `json:"element,omitempty"`
	}
	if err := json.Unmarshal(result, &qresp); err != nil {
		return nil, fmt.Errorf("locator.screenshot: parse querySelector failed: %w", err)
	}
	if qresp.Element == nil {
		return nil, fmt.Errorf("locator.screenshot: element not found for selector %q", l.selector)
	}
	el := elementHandleFromGUID(l.frame, qresp.Element.Guid)
	defer func() { _ = el.Dispose(context.Background()) }()
	return el.Screenshot(ctx, opts...)
}
