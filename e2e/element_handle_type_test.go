//go:build e2e

// E2E tests for: TestElementHandleType.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandleType.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleType.java
func TestElementHandleTypeInInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.Type(ctx, "hello"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Equal("hello", val)
}

// Ref: TestElementHandleType.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleType.java
func TestElementHandleTypeShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='text' />"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Type(ctx, "hello"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "hello", v)
}

// Ref: TestElementHandleType.java#shouldNotSelectExistingValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleType.java
func TestElementHandleTypeShouldNotSelectExistingValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='text' value='hello' />"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Type(ctx, "world"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "worldhello", v)
}

// Ref: TestElementHandleType.java#shouldResetSelectionWhenNotFocus  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleType.java
func TestElementHandleTypeShouldResetSelectionWhenNotFocus(t *testing.T) {
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
	must.NoError(input.Type(ctx, "world"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "worldhello", v)
}

// Ref: TestElementHandleType.java#shouldNotModifySelectionWhenFocus  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleType.java
func TestElementHandleTypeShouldNotModifySelectionWhenFocus(t *testing.T) {
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
	must.NoError(input.Type(ctx, "world"))
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, "heworldo", v)
}

// Ref: TestElementHandleType.java#shouldWorkWithNumberInput  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleType.java
func TestElementHandleTypeShouldWorkWithNumberInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='number' value=2 />"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Type(ctx, "13"))
	expected := "132"
	if a3IsWebKit(page) {
		expected = "13"
	}
	v, err := a3InputValue(ctx, page)
	must.NoError(err)
	assert.Equal(t, expected, v)
}
