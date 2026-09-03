//go:build e2e

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorIsVisibleTrueForVisibleElementEx2 verifies IsVisible returns true for visible element.
// Ref: TestLocatorIsVisible.java#shouldReturnTrue
func TestLocatorIsVisibleTrueForVisibleElementEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Visible</p>`))

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleFalseForDisplayNoneEx2 verifies IsVisible returns false for display:none.
// Ref: TestLocatorIsVisible.java#shouldReturnFalseForDisplayNone
func TestLocatorIsVisibleFalseForDisplayNoneEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="display:none">Hidden</p>`))

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsVisibleFalseForVisibilityHiddenEx2 verifies IsVisible returns false for visibility:hidden.
// Ref: TestLocatorIsVisible.java#shouldReturnFalseForVisibilityHidden
func TestLocatorIsVisibleFalseForVisibilityHiddenEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="visibility:hidden">Hidden</p>`))

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsVisibleTrueAfterShowEx2 verifies IsVisible returns true after making element visible.
// Ref: TestLocatorIsVisible.java#shouldReturnTrueAfterShow
func TestLocatorIsVisibleTrueAfterShowEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="display:none">Initially hidden</p>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('p').style.display = 'block'`)
	must.NoError(err)

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleButton verifies IsVisible returns true for button.
// Ref: TestLocatorIsVisible.java#shouldReturnTrueForButton
func TestLocatorIsVisibleButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click me</button>`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleAfterShowEx3 verifies IsVisible returns true after display change.
// Ref: TestLocatorIsVisible.java#shouldReturnTrueAfterShow
func TestIsVisibleAfterShowEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none;">Hidden</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').style.display = ''`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleAfterHideEx3 verifies IsVisible returns false after hiding.
// Ref: TestLocatorIsVisible.java#shouldReturnFalseAfterHide
func TestIsVisibleAfterHideEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').style.visibility = 'hidden'`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleButtonEx3 verifies IsVisible for button elements.
// Ref: TestLocatorIsVisible.java#shouldCheckButton
func TestIsVisibleButtonEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click me</button>`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleOpacityZeroEx3 verifies IsVisible for opacity:0 element.
// Ref: TestLocatorIsVisible.java#shouldReturnFalseForOpacityZero
func TestIsVisibleOpacityZeroEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="opacity:0;">Invisible</div>`))

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleAfterShowEx4 verifies IsVisible reflects JS-shown element.
// Ref: TestLocatorIsVisible.java#shouldReflectJSShow
func TestIsVisibleAfterShowEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none">Hidden</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').style.display = 'block'`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleButtonEx4 verifies IsVisible for visible button.
// Ref: TestLocatorIsVisible.java#shouldReturnTrueForButton
func TestIsVisibleButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleHiddenByOpacityEx4 verifies IsVisible for opacity-zero element.
// Ref: TestLocatorIsVisible.java#shouldReturnFalseForOpacityZero
func TestIsVisibleHiddenByOpacityEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="opacity:0">Invisible</div>`))

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleLinkEx4 verifies IsVisible for anchor element.
// Ref: TestLocatorIsVisible.java#shouldReturnTrueForLink
func TestIsVisibleLinkEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="#">Link</a>`))

	visible, err := page.Locator("#a").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleDisplayNoneEx verifies IsVisible returns false for display:none.
// Ref: TestLocatorIsVisible.java#shouldReturnFalseForDisplayNone
func TestIsVisibleDisplayNoneEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none">Hidden</div>`))

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleAfterClassRemoveEx verifies IsVisible after removing hidden class.
// Ref: TestLocatorIsVisible.java#shouldBeTrueAfterClassRemove
func TestIsVisibleAfterClassRemoveEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>.hidden { display:none; }</style>
		<div id="d" class="hidden">Content</div>
	`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').classList.remove('hidden')`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleSelectElementEx verifies IsVisible returns true for select.
// Ref: TestLocatorIsVisible.java#shouldReturnTrueForSelect
func TestIsVisibleSelectElementEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">A</option>
		</select>
	`))

	visible, err := page.Locator("#sel").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleImageEx verifies IsVisible returns true for img element.
// Ref: TestLocatorIsVisible.java#shouldReturnTrueForImage
func TestIsVisibleImageEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img id="img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			style="width:20px;height:20px;">
	`))

	visible, err := page.Locator("#img").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleButtonInHiddenContainerEx verifies IsVisible for button inside hidden div.
// Ref: TestLocatorIsVisible.java#shouldReturnFalseForButtonInHiddenContainer
func TestIsVisibleButtonInHiddenContainerEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="display:none">
			<button id="btn">Hidden button</button>
		</div>
	`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsVisibleTrueForVisibleElement verifies IsVisible returns true for visible element.
// Ref: TestLocatorVisibility.java#shouldReturnTrueForVisible
func TestLocatorIsVisibleTrueForVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="visible">I'm here</div>`))

	visible, err := page.Locator("#visible").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleFalseForHiddenElement verifies IsVisible returns false for display:none.
// Ref: TestLocatorVisibility.java#shouldReturnFalseForHidden
func TestLocatorIsVisibleFalseForHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="display:none">hidden</div>`))

	visible, err := page.Locator("#hidden").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsHiddenTrueForHiddenElement verifies IsVisible returns false for visibility:hidden element.
// Ref: TestLocatorVisibility.java#shouldReturnTrueForIsHidden
func TestLocatorIsHiddenTrueForHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="visibility:hidden">hidden</div>`))

	visible, err := page.Locator("#hidden").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsHiddenFalseForVisibleElement verifies IsVisible returns true for visible button.
// Ref: TestLocatorVisibility.java#shouldReturnFalseForIsHidden
func TestLocatorIsHiddenFalseForVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>visible button</button>`))

	visible, err := page.Locator("button").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsEnabledFalseForDisabledButton verifies IsEnabled returns false for disabled button.
// Ref: TestLocatorVisibility.java#shouldReturnFalseForDisabledButton
func TestLocatorIsEnabledFalseForDisabledButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button disabled>Can't click</button>`))

	enabled, err := page.Locator("button").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestLocatorIsDisabledFalseForEnabledInput verifies IsDisabled returns false for enabled input.
// Ref: TestLocatorVisibility.java#shouldReturnFalseForEnabledInput
func TestLocatorIsDisabledFalseForEnabledInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" placeholder="type here">`))

	disabled, err := page.Locator("input").IsDisabled(ctx)
	must.NoError(err)
	is.False(disabled)
}

// TestLocatorIsEditableTrueForInput verifies IsEditable returns true for a standard input.
// Ref: TestLocatorVisibility.java#shouldReturnTrueForEditableInput
func TestLocatorIsEditableTrueForInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.True(editable)
}

// TestLocatorIsEditableFalseForReadOnly verifies IsEditable returns false for readonly input.
// Ref: TestLocatorVisibility.java#shouldReturnFalseForReadonly
func TestLocatorIsEditableFalseForReadOnly(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" readonly>`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.False(editable)
}

// TestLocatorIsVisibleAfterStyleChange verifies IsVisible reflects style changes.
// Ref: TestLocatorVisibility.java#shouldReflectStyleChange
func TestLocatorIsVisibleAfterStyleChange(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" style="display:none">content</div>`))

	hidden, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.False(hidden)

	_, err = page.Evaluate(ctx, `() => document.getElementById('el').style.display = 'block'`)
	must.NoError(err)

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleForOpacityZero verifies opacity:0 element is visible.
// Ref: TestLocatorVisibility.java#shouldConsiderOpacityZeroVisible
func TestLocatorIsVisibleForOpacityZero(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// opacity:0 elements are visible (take up space)
	must.NoError(page.SetContent(ctx, `<div id="el" style="opacity:0;width:100px;height:100px">text</div>`))

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsEnabledForNonFormElement verifies IsEnabled for div is true.
// Ref: TestLocatorVisibility.java#shouldReturnEnabledForDiv
func TestLocatorIsEnabledForNonFormElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">content</div>`))

	enabled, err := page.Locator("#el").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestLocatorIsDisabledFalseForEnabledButton verifies IsDisabled false for enabled button.
// Ref: TestLocatorVisibility.java#shouldReturnFalseForEnabledButton
func TestLocatorIsDisabledFalseForEnabledButtonExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">go</button>`))

	disabled, err := page.Locator("#btn").IsDisabled(ctx)
	must.NoError(err)
	is.False(disabled)
}

// TestLocatorIsEditableFalseForDisabledInput verifies IsEditable false for disabled input.
// Ref: TestLocatorVisibility.java#shouldReturnFalseForDisabledInput
func TestLocatorIsEditableFalseForDisabledInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" disabled>`))

	editable, err := page.Locator("#inp").IsEditable(ctx)
	must.NoError(err)
	is.False(editable)
}
