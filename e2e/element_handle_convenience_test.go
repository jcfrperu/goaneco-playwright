//go:build e2e

// E2E tests for: TestElementHandleConvenience.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Ref: TestElementHandleConvenience.java#isEnabledAndIsDisabledShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceIsEnabledAndIsDisabledShouldWork(t *testing.T) {
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

// Ref: TestElementHandleConvenience.java#isVisibleAndIsHiddenShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceIsVisibleAndIsHiddenShouldWork(t *testing.T) {
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

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleTextContentReturnsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello <span>World</span></p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)
	must.NotNil(el)

	tc, err := el.TextContent(ctx)
	must.NoError(err)
	is.Equal("Hello World", tc)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleIsVisibleForVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="visible">I'm visible</div>`))

	el, err := page.QuerySelector(ctx, "#visible")
	must.NoError(err)
	must.NotNil(el)

	visible, err := el.IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleInnerHTMLExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"><span>child</span></div>`))

	handle, err := page.QuerySelector(ctx, "#container")
	must.NoError(err)
	must.NotNil(handle)

	html, err := handle.InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<span>child</span>")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleIsEditableForInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="editable" type="text">
		<input id="readonly" type="text" readonly>
	`))

	editable, err := page.QuerySelector(ctx, "#editable")
	must.NoError(err)
	isEditable, err := editable.IsEditable(ctx)
	must.NoError(err)
	is.True(isEditable)

	readonly, err := page.QuerySelector(ctx, "#readonly")
	must.NoError(err)
	isReadOnly, err := readonly.IsEditable(ctx)
	must.NoError(err)
	is.False(isReadOnly)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleInnerHTMLEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><span>child</span></div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	html, err := el.InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<span>child</span>", html)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleInnerHTMLEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><strong>Bold</strong></div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	html, err := el.InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<strong>")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleInnerHTML(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el"><b>hello</b></div>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	html, err := el.InnerHTML(ctx)
	must.NoError(err, "InnerHTML failed")
	is.Equal("<b>hello</b>", html)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleIsEnabledAndDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="ok">OK</button><button id="no" disabled>No</button>`)
	must.NoError(err, "SetContent failed")

	elOK, _ := page.QuerySelector(ctx, "#ok")
	elNo, _ := page.QuerySelector(ctx, "#no")

	enabled, err := elOK.IsEnabled(ctx)
	if err != nil || !enabled {
		t.Errorf("#ok should be enabled, got %v err=%v", enabled, err)
	}

	disabled, err := elNo.IsEnabled(ctx)
	if err != nil || disabled {
		t.Errorf("#no should be disabled (IsEnabled=false), got %v err=%v", disabled, err)
	}
}

// TestElementHandleInnerHTMLReturnsContent verifies InnerHTML returns the inner HTML string.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleInnerHTMLReturnsContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"><span>hello</span></div>`))

	el, err := page.QuerySelector(ctx, "#container")
	must.NoError(err)
	must.NotNil(el)

	html, err := el.InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<span>hello</span>", html)
}

// Ref: TestElementHandleConvenience.java#shouldHaveANicePreview  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceShouldHaveANicePreview(t *testing.T) {
	t.Skip("not implemented: JSHandle.toString (handle preview)")
}

// Ref: TestElementHandleConvenience.java#shouldHaveANicePreviewForNonAsciiAttributesChildren  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceShouldHaveANicePreviewForNonAsciiAttributesChildren(t *testing.T) {
	t.Skip("not implemented: JSHandle.toString (handle preview)")
}

// Ref: TestElementHandleConvenience.java#getAttributeShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceGetAttributeShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))
	handle, err := page.QuerySelector(ctx, "#outer")
	must.NoError(err)
	must.NotNil(handle)
	v, err := handle.GetAttribute(ctx, "name")
	must.NoError(err)
	must.NotNil(v)
	is.Equal("value", *v)
	v, err = handle.GetAttribute(ctx, "foo")
	must.NoError(err)
	is.Nil(v)
	// page-level variants (Page.getAttribute) go through the equivalent locator API
	v, err = page.Locator("#outer").GetAttribute(ctx, "name")
	must.NoError(err)
	must.NotNil(v)
	is.Equal("value", *v)
	v, err = page.Locator("#outer").GetAttribute(ctx, "foo")
	must.NoError(err)
	is.Nil(v)
}

// Ref: TestElementHandleConvenience.java#inputValueShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceInputValueShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	must.NoError(page.Locator("#textarea").Fill(ctx, "text value"))
	v, err := page.Locator("#textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("text value", v)

	must.NoError(page.Locator("#input").Fill(ctx, "input value"))
	v, err = page.Locator("#input").InputValue(ctx)
	must.NoError(err)
	is.Equal("input value", v)
	handle, err := page.QuerySelector(ctx, "#input")
	must.NoError(err)
	must.NotNil(handle)
	v, err = handle.InputValue(ctx)
	must.NoError(err)
	is.Equal("input value", v)

	_, err = page.Locator("#inner").InputValue(ctx)
	must.Error(err)
	is.Contains(err.Error(), "Node is not an <input>, <textarea> or <select> element")

	handle2, err := page.QuerySelector(ctx, "#inner")
	must.NoError(err)
	must.NotNil(handle2)
	_, err = handle2.InputValue(ctx)
	must.Error(err)
	is.Contains(err.Error(), "Node is not an <input>, <textarea> or <select> element")
}

// Ref: TestElementHandleConvenience.java#innerHTMLShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceInnerHTMLShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))
	handle, err := page.QuerySelector(ctx, "#outer")
	must.NoError(err)
	must.NotNil(handle)
	v, err := handle.InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<div id=\"inner\">Text,\nmore text</div>", v)
	v, err = page.Locator("#outer").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<div id=\"inner\">Text,\nmore text</div>", v)
}

// Ref: TestElementHandleConvenience.java#innerTextShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceInnerTextShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))
	handle, err := page.QuerySelector(ctx, "#inner")
	must.NoError(err)
	must.NotNil(handle)
	v, err := handle.InnerText(ctx)
	must.NoError(err)
	is.Equal("Text, more text", v)
	v, err = page.Locator("#inner").InnerText(ctx)
	must.NoError(err)
	is.Equal("Text, more text", v)
}

// Ref: TestElementHandleConvenience.java#innerTextShouldThrow  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceInnerTextShouldThrow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<svg>text</svg>"))
	_, err := page.Locator("svg").InnerText(ctx)
	must.Error(err)
	is.Contains(err.Error(), "Node is not an HTMLElement")
	handle, err := page.QuerySelector(ctx, "svg")
	must.NoError(err)
	must.NotNil(handle)
	_, err = handle.InnerText(ctx)
	must.Error(err)
	is.Contains(err.Error(), "Node is not an HTMLElement")
}

// Ref: TestElementHandleConvenience.java#textContentShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceTextContentShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))
	handle, err := page.QuerySelector(ctx, "#inner")
	must.NoError(err)
	must.NotNil(handle)
	v, err := handle.TextContent(ctx)
	must.NoError(err)
	is.Equal("Text,\nmore text", v)
	tc, err := page.Locator("#inner").TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal("Text,\nmore text", *tc)
}

// Ref: TestElementHandleConvenience.java#textContentShouldBeAtomic  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceTextContentShouldBeAtomic(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Context().RegisterSelectorEngine(ctx, "textContent", a3AtomicSelector))
	must.NoError(page.SetContent(ctx, "<div>Hello</div>"))
	tc, err := page.Locator("textContent=div").TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal("Hello", *tc)
	res, err := page.Evaluate(ctx, "() => document.querySelector('div').textContent")
	must.NoError(err)
	is.Equal("modified", res)
}

// Ref: TestElementHandleConvenience.java#innerTextShouldBeAtomic  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceInnerTextShouldBeAtomic(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Context().RegisterSelectorEngine(ctx, "innerText", a3AtomicSelector))
	must.NoError(page.SetContent(ctx, "<div>Hello</div>"))
	tc, err := page.Locator("innerText=div").InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello", tc)
	res, err := page.Evaluate(ctx, "() => document.querySelector('div').innerText")
	must.NoError(err)
	is.Equal("modified", res)
}

// Ref: TestElementHandleConvenience.java#innerHTMLShouldBeAtomic  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceInnerHTMLShouldBeAtomic(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Context().RegisterSelectorEngine(ctx, "innerHTML", a3AtomicSelector))
	must.NoError(page.SetContent(ctx, "<div>Hello<span>world</span></div>"))
	tc, err := page.Locator("innerHTML=div").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("Hello<span>world</span>", tc)
	res, err := page.Evaluate(ctx, "() => document.querySelector('div').innerHTML")
	must.NoError(err)
	is.Equal("modified", res)
}

// Ref: TestElementHandleConvenience.java#getAttributeShouldBeAtomic  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceGetAttributeShouldBeAtomic(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	script := `{
  query(root, selector) {
    const result = root.querySelector(selector);
    if (result)
      Promise.resolve().then(() => result.setAttribute('foo', 'modified'));
    return result;
  },
  queryAll(root, selector) {
    const result = Array.from(root.querySelectorAll(selector));
    for (const e of result)
      Promise.resolve().then(() => e.setAttribute('foo', 'modified'));
    return result;
  }
}`
	must.NoError(page.Context().RegisterSelectorEngine(ctx, "getAttribute", script))
	must.NoError(page.SetContent(ctx, "<div foo=hello></div>"))
	tc, err := page.Locator("getAttribute=div").GetAttribute(ctx, "foo")
	must.NoError(err)
	must.NotNil(tc)
	is.Equal("hello", *tc)
	res, err := page.Evaluate(ctx, "() => document.querySelector('div').getAttribute('foo')")
	must.NoError(err)
	is.Equal("modified", res)
}

// Ref: TestElementHandleConvenience.java#isEditableShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceIsEditableShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input id=input1 disabled><textarea></textarea><input id=input2>"))
	_, err := page.EvalOnSelector(ctx, "textarea", "t => t.readOnly = true")
	must.NoError(err)

	input1, err := page.QuerySelector(ctx, "#input1")
	must.NoError(err)
	must.NotNil(input1)
	v, err := input1.IsEditable(ctx)
	must.NoError(err)
	is.False(v)
	v, err = page.Locator("#input1").IsEditable(ctx)
	must.NoError(err)
	is.False(v)

	input2, err := page.QuerySelector(ctx, "#input2")
	must.NoError(err)
	must.NotNil(input2)
	v, err = input2.IsEditable(ctx)
	must.NoError(err)
	is.True(v)
	v, err = page.Locator("#input2").IsEditable(ctx)
	must.NoError(err)
	is.True(v)

	textarea, err := page.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NotNil(textarea)
	v, err = textarea.IsEditable(ctx)
	must.NoError(err)
	is.False(v)
	v, err = page.Locator("textarea").IsEditable(ctx)
	must.NoError(err)
	is.False(v)
}

// Ref: TestElementHandleConvenience.java#isCheckedShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleConvenience.java
func TestElementHandleConvenienceIsCheckedShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input type='checkbox' checked><div>Not a checkbox</div>"))
	handle, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(handle)
	v, err := handle.IsChecked(ctx)
	must.NoError(err)
	is.True(v)
	v, err = page.Locator("input").IsChecked(ctx)
	must.NoError(err)
	is.True(v)
	_, err = handle.Evaluate(ctx, "input => input.checked = false")
	must.NoError(err)
	v, err = handle.IsChecked(ctx)
	must.NoError(err)
	is.False(v)
	v, err = page.Locator("input").IsChecked(ctx)
	must.NoError(err)
	is.False(v)
	_, err = page.Locator("div").IsChecked(ctx)
	must.Error(err)
	is.Contains(err.Error(), "Not a checkbox or radio button")
}
