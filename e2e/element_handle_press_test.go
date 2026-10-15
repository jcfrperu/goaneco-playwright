//go:build e2e

// E2E tests for: TestElementHandlePress.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandlePress.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandlePress.java
func TestElementHandlePressEnterSubmitsForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form onsubmit="window.__submitted=true;return false;">
			<input id="inp" type="text">
			<button type="submit">Go</button>
		</form>
	`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Press(ctx, "Enter"))

	result, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandlePress.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandlePress.java
func TestElementHandlePressEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Type(ctx, "test"))
	must.NoError(el.Press(ctx, "Backspace"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("tes", val)
}

// Ref: TestElementHandlePress.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandlePress.java
func TestElementHandlePressShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='text' />"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Press(ctx, "h"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "h", v)
}

// Ref: TestElementHandlePress.java#shouldNotSelectExistingValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandlePress.java
func TestElementHandlePressShouldNotSelectExistingValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='text' value='hello' />"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Press(ctx, "w"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "whello", v)
}

// Ref: TestElementHandlePress.java#shouldResetSelectionWhenNotFocused  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandlePress.java
func TestElementHandlePressShouldResetSelectionWhenNotFocused(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='text' value='hello' /><div tabIndex=2>text</div>"))
	_, err := page.EvalOnSelector(ctx, "input", `input => {
    input.selectionStart = 2;
    input.selectionEnd = 4;
    document.querySelector('div').focus();
  }`)
	must.NoError(err)
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Press(ctx, "w"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "whello", v)
}

// Ref: TestElementHandlePress.java#shouldNotModifySelectionWhenFocused  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandlePress.java
func TestElementHandlePressShouldNotModifySelectionWhenFocused(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='text' value='hello' />"))
	_, err := page.EvalOnSelector(ctx, "input", `input => {
    input.focus();
    input.selectionStart = 2;
    input.selectionEnd = 4;
  }`)
	must.NoError(err)
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Press(ctx, "w"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "hewo", v)
}

// Ref: TestElementHandlePress.java#shouldWorkWithNumberInput  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandlePress.java
func TestElementHandlePressShouldWorkWithNumberInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='number' value=2 />"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Press(ctx, "1"))
	expected := "12"
	if a3IsWebKit(page) {
		expected = "1"
	}
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, expected, v)
}
