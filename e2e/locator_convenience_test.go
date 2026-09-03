//go:build e2e

// Locator convenience method tests.
// Migration of: TestLocatorConvenience.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocatorTextContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">hello <b>world</b></div>`)
	must.NoError(err, "SetContent failed")

	tc, err := page.Locator("#el").TextContent(ctx)
	must.NoError(err, "TextContent failed")
	if tc == nil || *tc != "hello world" {
		t.Errorf("TextContent = %v, want 'hello world'", tc)
	}
}

func TestLocatorInnerHTML(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el"><b>hello</b></div>`)
	must.NoError(err, "SetContent failed")

	html, err := page.Locator("#el").InnerHTML(ctx)
	must.NoError(err, "InnerHTML failed")
	is.Equal("<b>hello</b>", html)
}

func TestLocatorInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">hello</div><div id="hidden" style="display:none">hidden</div>`)
	must.NoError(err, "SetContent failed")

	text, err := page.Locator("#el").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("hello", text)
}

func TestLocatorInputValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="in" value="initial"/>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator("#in").InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("initial", val)
}

func TestLocatorIsEnabledAndDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="ok">OK</button><button id="no" disabled>No</button>`)
	must.NoError(err, "SetContent failed")

	enabled, err := page.Locator("#ok").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled, "expected #ok to be enabled")

	disabled, err := page.Locator("#no").IsEnabled(ctx)
	must.NoError(err)
	is.False(disabled, "expected #no to be disabled (IsEnabled=false)")
}

func TestLocatorIsEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="ed"/><input id="ro" readonly/>`)
	must.NoError(err, "SetContent failed")

	editable, err := page.Locator("#ed").IsEditable(ctx)
	must.NoError(err)
	is.True(editable, "expected #ed to be editable")

	readonly, err := page.Locator("#ro").IsEditable(ctx)
	must.NoError(err)
	is.False(readonly, "expected #ro to NOT be editable")
}

func TestLocatorIsChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox" id="ch" checked/><input type="checkbox" id="un"/>`)
	must.NoError(err, "SetContent failed")

	checked, err := page.Locator("#ch").IsChecked(ctx)
	must.NoError(err)
	is.True(checked, "expected #ch to be checked")

	unchecked, err := page.Locator("#un").IsChecked(ctx)
	must.NoError(err)
	is.False(unchecked, "expected #un to NOT be checked")
}

func TestLocatorAllTextContents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul><li>one</li><li>two</li><li>three</li></ul>`)
	must.NoError(err, "SetContent failed")

	texts, err := page.Locator("li").AllTextContents(ctx)
	must.NoError(err, "AllTextContents failed")
	is.Len(texts, 3)
	if texts[0] != "one" || texts[1] != "two" || texts[2] != "three" {
		t.Errorf("unexpected texts: %v", texts)
	}
}

func TestLocatorAllInnerTexts(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="item">alpha</div><div class="item">beta</div>`)
	must.NoError(err, "SetContent failed")

	texts, err := page.Locator(".item").AllInnerTexts(ctx)
	must.NoError(err, "AllInnerTexts failed")
	is.Len(texts, 2)
	if texts[0] != "alpha" || texts[1] != "beta" {
		t.Errorf("unexpected texts: %v", texts)
	}
}

func TestLocatorIsVisibleAndHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="vis">visible</div><div id="hid" style="display:none">hidden</div>`)
	must.NoError(err, "SetContent failed")

	visible, err := page.Locator("#vis").IsVisible(ctx)
	must.NoError(err)
	is.True(visible, "expected #vis to be visible")

	hidden, err := page.Locator("#hid").IsVisible(ctx)
	must.NoError(err)
	is.False(hidden, "expected #hid to NOT be visible")
}

func TestLocatorGetAttributeConvenience(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" data-value="42">text</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator("#el").GetAttribute(ctx, "data-value")
	must.NoError(err, "GetAttribute failed")
	if val == nil || *val != "42" {
		t.Errorf("GetAttribute(data-value) = %v, want '42'", val)
	}
}

func TestLocatorAllShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul><li>a</li><li>b</li><li>c</li></ul>`)
	must.NoError(err, "SetContent failed")

	items, err := page.Locator("li").All(ctx)
	must.NoError(err, "All failed")
	is.Len(items, 3)

	var texts []string
	for _, loc := range items {
		tc, err := loc.TextContent(ctx)
		if err != nil || tc == nil {
			t.Fatalf("TextContent on nth locator failed: %v", err)
		}
		texts = append(texts, *tc)
	}
	if texts[0] != "a" || texts[1] != "b" || texts[2] != "c" {
		t.Errorf("unexpected texts from All(): %v", texts)
	}
}

func TestLocatorAllEmptyWhenNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>nothing</div>`)
	must.NoError(err, "SetContent failed")

	items, err := page.Locator("span").All(ctx)
	must.NoError(err, "All failed")
	if len(items) != 0 {
		t.Errorf("expected empty All() for non-matching selector, got %d", len(items))
	}
}

func TestLocatorHasLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="parent"><span class="child">yes</span></div><div class="parent">no</div>`)
	must.NoError(err, "SetContent failed")

	loc := page.Locator(".parent").Filter(&playwright.LocatorFilterOptions{
		Has: page.Locator(".child"),
	})
	count, err := loc.Count(ctx)
	must.NoError(err, "Count failed")
	is.Equal(1, count)
}

// TestLocatorIsEnabledTrueForInputEx2 verifies IsEnabled returns true for normal input.
// Ref: TestLocatorIsEnabled.java#shouldReturnTrueForInput
func TestLocatorIsEnabledTrueForInputEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestLocatorIsEnabledFalseForDisabledInputEx2 verifies IsEnabled returns false for disabled input.
// Ref: TestLocatorIsEnabled.java#shouldReturnFalseForDisabledInput
func TestLocatorIsEnabledFalseForDisabledInputEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" disabled>`))

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestLocatorIsEnabledTrueForButtonEx2 verifies IsEnabled returns true for normal button.
// Ref: TestLocatorIsEnabled.java#shouldReturnTrueForButton
func TestLocatorIsEnabledTrueForButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	enabled, err := page.Locator("#btn").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestLocatorIsEnabledFalseForDisabledButtonEx2 verifies IsEnabled returns false for disabled button.
// Ref: TestLocatorIsEnabled.java#shouldReturnFalseForDisabledButton
func TestLocatorIsEnabledFalseForDisabledButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Click</button>`))

	enabled, err := page.Locator("#btn").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestLocatorIsEnabledSelectEx2 verifies IsEnabled returns true for select.
// Ref: TestLocatorIsEnabled.java#shouldReturnTrueForSelect
func TestLocatorIsEnabledSelectEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel"><option value="a">A</option></select>`))

	enabled, err := page.Locator("#sel").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledSelectEx3 verifies IsEnabled for select element.
// Ref: TestLocatorIsEnabled.java#shouldCheckSelect
func TestIsEnabledSelectEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option>A</option>
		</select>
	`))

	enabled, err := page.Locator("#s").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledSelectEx3 verifies IsEnabled returns false for disabled select.
// Ref: TestLocatorIsEnabled.java#shouldReturnFalseForDisabledSelect
func TestIsEnabledDisabledSelectEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s" disabled>
			<option>A</option>
		</select>
	`))

	enabled, err := page.Locator("#s").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledTextareaEx3 verifies IsEnabled for textarea.
// Ref: TestLocatorIsEnabled.java#shouldCheckTextarea
func TestIsEnabledTextareaEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">text</textarea>`))

	enabled, err := page.Locator("#ta").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledTextareaEx3 verifies IsEnabled returns false for disabled textarea.
// Ref: TestLocatorIsEnabled.java#shouldReturnFalseForDisabledTextarea
func TestIsEnabledDisabledTextareaEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta" disabled>text</textarea>`))

	enabled, err := page.Locator("#ta").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledAfterEnableEx3 verifies IsEnabled reflects dynamic enable change.
// Ref: TestLocatorIsEnabled.java#shouldReflectDynamicEnable
func TestIsEnabledAfterEnableEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" disabled>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('inp').disabled = false`)
	must.NoError(err)

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledSelectEx4 verifies IsEnabled for select element.
// Ref: TestLocatorIsEnabled.java#shouldReturnTrueForSelect
func TestIsEnabledSelectEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel"><option>One</option></select>`))

	enabled, err := page.Locator("#sel").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledSelectEx4 verifies IsEnabled returns false for disabled select.
// Ref: TestLocatorIsEnabled.java#shouldReturnFalseForDisabledSelect
func TestIsEnabledDisabledSelectEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel" disabled><option>One</option></select>`))

	enabled, err := page.Locator("#sel").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledTextareaEx4 verifies IsEnabled for textarea element.
// Ref: TestLocatorIsEnabled.java#shouldReturnTrueForTextarea
func TestIsEnabledTextareaEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">Content</textarea>`))

	enabled, err := page.Locator("#ta").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledAfterJSEnableEx4 verifies IsEnabled reflects JS-enabled element.
// Ref: TestLocatorIsEnabled.java#shouldReflectJSEnabled
func TestIsEnabledAfterJSEnableEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" disabled>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('inp').disabled = false`)
	must.NoError(err)

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledAfterJSEnableFieldsetEx5 verifies IsEnabled after JS enables fieldset child.
// Ref: TestLocatorIsEnabled.java#shouldReflectFieldsetDisabled
func TestIsEnabledAfterJSEnableFieldsetEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<fieldset disabled>
			<input id="inp" type="text">
		</fieldset>
	`))

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledButtonWithTypeEx5 verifies IsEnabled for submit button.
// Ref: TestLocatorIsEnabled.java#shouldReturnTrueForSubmitButton
func TestIsEnabledButtonWithTypeEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="submit" type="submit">Submit</button>`))

	enabled, err := page.Locator("#submit").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledButtonWithTypeEx5 verifies IsEnabled for disabled submit button.
// Ref: TestLocatorIsEnabled.java#shouldReturnFalseForDisabledSubmit
func TestIsEnabledDisabledButtonWithTypeEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="submit" type="submit" disabled>Submit</button>`))

	enabled, err := page.Locator("#submit").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}
