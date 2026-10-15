//go:build e2e

// E2E tests for: TestElementHandleSelectText.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandleSelectText.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectOptionExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
			<option value="c">Gamma</option>
		</select>
	`))

	handle, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(handle)

	selected, err := handle.SelectOption(ctx, "b")
	must.NoError(err)
	is.Contains(selected, "b")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleSelectText.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectOptionEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
	`))

	el, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(el)

	_, err = el.SelectOption(ctx, "b")
	must.NoError(err)

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("b", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleSelectText.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectOptionEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
			<option value="c">Gamma</option>
		</select>
	`))

	el, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(el)

	_, err = el.SelectOption(ctx, "b")
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => document.getElementById('sel').value`)
	must.NoError(err)
	is.Equal("b", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleSelectText.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectOptionEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
			<option value="c">Gamma</option>
		</select>
	`))

	el, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(el)

	_, err = el.SelectOption(ctx, "b")
	must.NoError(err)

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("b", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleSelectText.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select id="sel"><option value="foo">Foo</option><option value="bar">Bar</option></select>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#sel")

	selected, err := el.SelectOption(ctx, "bar")
	must.NoError(err, "SelectOption failed")
	if len(selected) != 1 || selected[0] != "bar" {
		t.Errorf("SelectOption returned %v, want ['bar']", selected)
	}
}

// TestElementHandleInnerTextReturnsText verifies InnerText returns visible text.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleSelectText.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleInnerTextReturnsText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello <b>World</b></p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello World", text)
}

// Ref: TestElementHandleSelectText.java#shouldSelectTextarea  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectTextShouldSelectTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<textarea id='text'>hello world</textarea>"))
	handle, err := page.QuerySelector(ctx, "#text")
	must.NoError(err)
	must.NotNil(handle)
	must.NoError(handle.SelectText(ctx))
	selected, err := page.Evaluate(ctx, "() => window.getSelection().toString()")
	must.NoError(err)
	assert.Equal(t, "hello world", selected)
}

// Ref: TestElementHandleSelectText.java#shouldSelectInput  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectTextShouldSelectInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input id='input' value='hello world' />"))
	handle, err := page.QuerySelector(ctx, "#input")
	must.NoError(err)
	must.NotNil(handle)
	must.NoError(handle.SelectText(ctx))
	selected, err := page.Evaluate(ctx, "() => document.getElementById('input').value.substring(document.getElementById('input').selectionStart, document.getElementById('input').selectionEnd)")
	must.NoError(err)
	assert.Equal(t, "hello world", selected)
}

// Ref: TestElementHandleSelectText.java#shouldSelectPlainDiv  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectTextShouldSelectPlainDiv(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div id='div'>hello world</div>"))
	handle, err := page.QuerySelector(ctx, "#div")
	must.NoError(err)
	must.NotNil(handle)
	must.NoError(handle.SelectText(ctx))
	selected, err := page.Evaluate(ctx, "() => window.getSelection().toString()")
	must.NoError(err)
	assert.Equal(t, "hello world", selected)
}

// Ref: TestElementHandleSelectText.java#shouldTimeoutWaitingForInvisibleElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectTextShouldTimeoutWaitingForInvisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div id='div' style='display:none'>hello world</div>"))
	handle, err := page.QuerySelector(ctx, "#div")
	must.NoError(err)
	must.NotNil(handle)
	timeout := 100.0
	err = handle.SelectText(ctx, &playwright.ElementHandleSelectTextOptions{Timeout: &timeout})
	must.Error(err)
	assert.Contains(t, err.Error(), "Timeout")
}

// Ref: TestElementHandleSelectText.java#shouldWaitForVisible  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleSelectText.java
func TestElementHandleSelectTextShouldWaitForVisible(t *testing.T) {
	t.Skip("not implemented: ElementHandle.SelectText (the Java test is an empty, disabled placeholder)")
}
