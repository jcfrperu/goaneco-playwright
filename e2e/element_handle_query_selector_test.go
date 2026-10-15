//go:build e2e

// E2E tests for: TestElementHandleQuerySelector.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/el", "text/html", `<div id="target">hello world</div>`)
	err := page.Goto(ctx, srv.Prefix()+"/el")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "#target")
	must.NoError(err, "QuerySelector failed")
	must.NotNil(el, "QuerySelector returned nil for existing element")

	text, err := el.InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("hello world", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorReturnsNilWhenNotFound(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/el", "text/html", `<div>nothing</div>`)
	err := page.Goto(ctx, srv.Prefix()+"/el")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err, "QuerySelector failed")
	is.Nil(el)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/list", "text/html", `<ul><li>one</li><li>two</li><li>three</li></ul>`)
	err := page.Goto(ctx, srv.Prefix()+"/list")
	must.NoError(err, "Goto failed")

	elements, err := page.QuerySelectorAll(ctx, "li")
	must.NoError(err, "QuerySelectorAll failed")
	is.Len(elements, 3)

	wantTexts := []string{"one", "two", "three"}
	for i, el := range elements {
		text, err := el.InnerText(ctx)
		must.NoErrorf(err, "InnerText[%d] failed", i)
		if text != wantTexts[i] {
			t.Errorf("element[%d] text = %q, want %q", i, text, wantTexts[i])
		}
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorFindsChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="parent"><span id="child">child</span></div>`))

	parent, err := page.QuerySelector(ctx, "#parent")
	must.NoError(err)
	must.NotNil(parent)

	child, err := parent.QuerySelector(ctx, "#child")
	must.NoError(err)
	must.NotNil(child)

	text, err := child.InnerText(ctx)
	must.NoError(err)
	is.Equal("child", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorAllExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			<span>a</span>
			<span>b</span>
			<span>c</span>
		</div>
	`))

	parent, err := page.QuerySelector(ctx, "#parent")
	must.NoError(err)
	must.NotNil(parent)

	children, err := parent.QuerySelectorAll(ctx, "span")
	must.NoError(err)
	is.Len(children, 3)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="container">
			<span id="child">Child element</span>
		</div>
	`))

	container, err := page.QuerySelector(ctx, "#container")
	must.NoError(err)
	must.NotNil(container)

	child, err := container.QuerySelector(ctx, "#child")
	must.NoError(err)
	must.NotNil(child)

	text, err := child.TextContent(ctx)
	must.NoError(err)
	is.Equal("Child element", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorAllEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul id="list">
			<li>A</li>
			<li>B</li>
			<li>C</li>
		</ul>
	`))

	list, err := page.QuerySelector(ctx, "#list")
	must.NoError(err)
	must.NotNil(list)

	items, err := list.QuerySelectorAll(ctx, "li")
	must.NoError(err)
	is.Len(items, 3)
}

// --- From element_handle_bounding_box_extra_test.go ---

// TestElementHandleQuerySelectorXpathShouldQueryExistingElement verifies ElementHandle.QuerySelectorAll with an XPath selector.
// Ref: TestElementHandleQuerySelector.java#xpathShouldQueryExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorXpathShouldQueryExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><div class="second"><div class="inner">A</div></div></body></html>`))

	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)

	second, err := html.QuerySelectorAll(ctx, `xpath=./body/div[contains(@class, 'second')]`)
	must.NoError(err)
	is.Len(second, 1)

	inner, err := second[0].QuerySelectorAll(ctx, `xpath=./div[contains(@class, 'inner')]`)
	must.NoError(err)
	is.Len(inner, 1)

	text, err := inner[0].TextContent(ctx)
	must.NoError(err)
	is.Equal("A", text)
}

// TestElementHandleQuerySelectorXpathShouldReturnNullForNonExistingElement verifies QuerySelectorAll with XPath returns empty for no matches.
// Ref: TestElementHandleQuerySelector.java#xpathShouldReturnNullForNonExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorXpathShouldReturnNullForNonExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><div class="second"><div class="inner">B</div></div></body></html>`))

	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)

	elements, err := html.QuerySelectorAll(ctx, `xpath=/div[contains(@class, 'third')]`)
	must.NoError(err)
	is.Empty(elements)
}

// TestElementHandleQuerySelectorShouldWorkForAdoptedElements verifies ElementHandle remains valid after cross-document adoption.
// Ref: TestElementHandleQuerySelector.java#shouldWorkForAdoptedElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorShouldWorkForAdoptedElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup := waitForPopup(t, page, func() {
		_, err := page.Evaluate(ctx, `url => window['__popup'] = window.open(url)`, srv.EmptyPage())
		must.NoError(err)
	})
	must.NotNil(popup)

	// Create a div with a span in the main page.
	_, err := page.Evaluate(ctx, `() => {
		const div = document.createElement('div');
		document.body.appendChild(div);
		const span = document.createElement('span');
		span.textContent = 'hello';
		div.appendChild(span);
	}`)
	must.NoError(err)

	elHandle, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(elHandle)

	// Span is accessible before adoption.
	span, err := elHandle.QuerySelector(ctx, "span")
	must.NoError(err)
	must.NotNil(span)

	text, err := span.TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text)

	// Adopt the div into the popup's document.
	must.NoError(popup.WaitForLoadState(ctx, "domcontentloaded"))
	_, err = page.Evaluate(ctx, `() => {
		const div = document.querySelector('div');
		window['__popup'].document.body.appendChild(div);
	}`)
	must.NoError(err)

	// ElementHandle still resolves the span after adoption.
	span2, err := elHandle.QuerySelector(ctx, "span")
	must.NoError(err)
	must.NotNil(span2)

	text2, err := span2.TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text2)

	// The span is also queryable from the popup.
	spanInPopup, err := popup.QuerySelector(ctx, "span")
	must.NoError(err)
	must.NotNil(spanInPopup)

	text3, err := spanInPopup.TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text3)
}

// TestPageQuerySelectorWithAttributeSelector verifies attribute-based CSS selectors.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestPageQuerySelectorWithAttributeSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" data-test="first">
		<input type="password" data-test="second">
	`))

	el, err := page.QuerySelector(ctx, `input[type="password"]`)
	must.NoError(err)
	must.NotNil(el)

	attr, err := el.GetAttribute(ctx, "data-test")
	must.NoError(err)
	is.Equal("second", attr)
}

// TestFrameQuerySelectorReturnsElement verifies Frame.QuerySelector finds elements.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestFrameQuerySelectorReturnsElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">found</div>`))

	el, err := page.MainFrame().QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("found", text)
}

// TestFrameQuerySelectorReturnsElementExtra2 verifies Frame.QuerySelector finds element.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestFrameQuerySelectorReturnsElementExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">click</button>`))

	frame := page.MainFrame()
	handle, err := frame.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(handle)
}

// Ref: TestElementHandleQuerySelector.java#shouldQueryExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorShouldQueryExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><div class="second"><div class="inner">A</div></div></body></html>`))
	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)
	second, err := html.QuerySelector(ctx, ".second")
	must.NoError(err)
	must.NotNil(second)
	inner, err := second.QuerySelector(ctx, ".inner")
	must.NoError(err)
	must.NotNil(inner)
	content, err := page.Evaluate(ctx, "e => e.textContent", inner)
	must.NoError(err)
	assert.Equal(t, "A", content)
}

// Ref: TestElementHandleQuerySelector.java#shouldReturnNullForNonExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorShouldReturnNullForNonExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><div class="second"><div class="inner">B</div></div></body></html>`))
	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)
	second, err := html.QuerySelector(ctx, ".third")
	must.NoError(err)
	assert.Nil(t, second)
}

// Ref: TestElementHandleQuerySelector.java#shouldQueryExistingElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorShouldQueryExistingElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<html><body><div>A</div><br/><div>B</div></body></html>"))
	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)
	elements, err := html.QuerySelectorAll(ctx, "div")
	must.NoError(err)
	must.Len(elements, 2)
	var result []string
	for _, element := range elements {
		v, err := page.Evaluate(ctx, "e => e.textContent", element)
		must.NoError(err)
		s, ok := v.(string)
		must.True(ok)
		result = append(result, s)
	}
	assert.Equal(t, []string{"A", "B"}, result)
}

// Ref: TestElementHandleQuerySelector.java#shouldReturnEmptyArrayForNonExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleQuerySelector.java
func TestElementHandleQuerySelectorShouldReturnEmptyArrayForNonExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<html><body><span>A</span><br/><span>B</span></body></html>"))
	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)
	elements, err := html.QuerySelectorAll(ctx, "div")
	must.NoError(err)
	assert.Len(t, elements, 0)
}
