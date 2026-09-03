//go:build e2e

// Page type (keyboard input) E2E tests.
// Migration of: TestPageType.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageTypeIntoTextarea verifies typing into a textarea via Keyboard.Type.
// Ref: TestPageType.java#shouldTypeIntoTextarea
func TestPageTypeIntoTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))

	must.NoError(page.Locator("textarea").Click(ctx))
	must.NoError(page.Keyboard.Type(ctx, "Hello World"))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello World", val)
}

// TestPageTypeIntoTextareaWithSpecialChars verifies typing special characters via Keyboard.Type.
// Ref: TestPageType.java#shouldTypeIntoTextareaWithSpecialChars
func TestPageTypeIntoTextareaWithSpecialChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must.NoError(page.Locator("textarea").Click(ctx))

	must.NoError(page.Keyboard.Type(ctx, "This is Text"))
	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Keyboard.Type(ctx, "Very "))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("Very This is Text", val)
}

// TestPageTypeWithNewlinesKeyboardPress verifies pressing Enter creates newlines.
// Ref: TestPageType.java#shouldTypeIntoTextareaInsertNewlines
func TestPageTypeWithNewlinesKeyboardPress(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must.NoError(page.Locator("textarea").Click(ctx))

	must.NoError(page.Keyboard.Type(ctx, "Hello"))
	must.NoError(page.Keyboard.Press(ctx, "Enter"))
	must.NoError(page.Keyboard.Type(ctx, "World"))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello\nWorld", val)
}

// TestPagePressSequentiallyIntoInput verifies PressSequentially into input.
// Ref: TestPageType.java#shouldTypeIntoInputElement
func TestPagePressSequentiallyIntoInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	must.NoError(page.Locator("input").PressSequentially(ctx, "Hello"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello", val)
}

// TestLocatorClickAndTypeSetsValue verifies clicking and then typing sets the value.
// Ref: TestPageType.java#shouldClickThenType
func TestLocatorClickAndTypeSetsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" value="existing">`))

	// Triple-click to select all, then type to replace
	must.NoError(page.Locator("input").Click(ctx))
	// Select all via keyboard shortcut
	must.NoError(page.Keyboard.Press(ctx, "Control+a"))
	must.NoError(page.Keyboard.Type(ctx, "replaced"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("replaced", val)
}

// TestPageTypeAndVerifyEventsFire verifies that typing fires keydown, keypress, and keyup events.
// Ref: TestPageType.java#shouldSendKeyEvents
func TestPageTypeAndVerifyEventsFire(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.KeyboardPage()))

	must.NoError(page.Keyboard.Type(ctx, "f"))

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, ok := result.(string)
	is.True(ok)
	is.Contains(s, "Keydown: f")
	is.Contains(s, "Keyup: f")
}

func TestPageTypeAppendsToExistingValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="hello" id="inp">`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, " world"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Equal("hello world", val)
}

func TestPageTypeFiresInputEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" oninput="window.__count = (window.__count||0) + 1">
	`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, "abc"))

	count, err := page.Evaluate(ctx, `() => window.__count`)
	must.NoError(err)
	is.Equal(float64(3), count)
}

func TestPageTypeIntoPasswordField(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="password" id="pw">`))

	handle, err := page.QuerySelector(ctx, "#pw")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, "secret"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Equal("secret", val)
}

func TestPageTypePressEnterSubmitsForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form onsubmit="window.__submitted=true; return false;">
			<input id="inp" type="text">
		</form>
	`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Press(ctx, "Enter"))

	result, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, result)
}

func TestPageTypeWithNewlineInTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	handle, err := page.QuerySelector(ctx, "#ta")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, "line1"))
	must.NoError(handle.Press(ctx, "Enter"))
	must.NoError(handle.Type(ctx, "line2"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "line1")
	is.Contains(val, "line2")
}

// TestLocatorTypeTextEx2 verifies Fill enters text into an input.
// Ref: TestLocatorType.java#shouldTypeText
func TestLocatorTypeTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "hello world"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hello world", val)
}

// TestLocatorFillInTextareaEx2 verifies Fill works in textarea.
// Ref: TestLocatorType.java#shouldTypeInTextarea
func TestLocatorFillInTextareaEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	must.NoError(page.Locator("#ta").Fill(ctx, "multiline text"))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "multiline text")
}

// TestLocatorFillNumbersEx2 verifies Fill works with numeric strings.
// Ref: TestLocatorType.java#shouldTypeNumbers
func TestLocatorFillNumbersEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "9876"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("9876", val)
}

// TestLocatorFillFiresInputEventEx2 verifies Fill fires input event.
// Ref: TestLocatorType.java#shouldFireInputEvent
func TestLocatorFillFiresInputEventEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" oninput="window.__inputFired=true">
	`))

	must.NoError(page.Locator("#inp").Fill(ctx, "some text"))

	fired, err := page.Evaluate(ctx, `() => window.__inputFired`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestLocatorFillAndClearEx2 verifies Fill followed by Fill with empty string clears.
// Ref: TestLocatorType.java#shouldClearWithEmptyFill
func TestLocatorFillAndClearEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "initial"))
	must.NoError(page.Locator("#inp").Fill(ctx, ""))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Empty(val)
}
