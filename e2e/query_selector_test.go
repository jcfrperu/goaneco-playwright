//go:build e2e

// E2E tests for: TestQuerySelector.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPageQuerySelectorReturnsNilWhenNotFound verifies QuerySelector returns nil when selector has no match.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestPageQuerySelectorReturnsNilWhenNotFound(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err)
	is.Nil(el, "QuerySelector should return nil for non-existing element")
}

// TestPageQuerySelectorAllReturnsEmptySlice verifies QuerySelectorAll returns empty slice when no match.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestPageQuerySelectorAllReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	handles, err := page.QuerySelectorAll(ctx, "span.nonexistent")
	must.NoError(err)
	is.Empty(handles)
}

// TestPageQuerySelectorAllWithNestedSelector verifies QuerySelectorAll with descendant selectors.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestPageQuerySelectorAllWithNestedSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="container">
			<p>one</p>
			<p>two</p>
		</div>
		<p>three</p>
	`))

	handles, err := page.QuerySelectorAll(ctx, ".container p")
	must.NoError(err)
	is.Len(handles, 2)
}

// TestQuerySelectorAllFindsMultipleEx2 verifies QuerySelectorAll returns multiple elements.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorAllFindsMultipleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">A</div>
		<div class="item">B</div>
		<div class="item">C</div>
	`))

	elements, err := page.QuerySelectorAll(ctx, ".item")
	must.NoError(err)
	is.Len(elements, 3)
}

// TestQuerySelectorByCSSClassEx2 verifies QuerySelector finds element by class.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorByCSSClassEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p class="highlight">Highlighted text</p>`))

	el, err := page.QuerySelector(ctx, ".highlight")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorAllReturnsMultiple verifies QuerySelectorAll returns all matching elements.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorAllReturnsMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">a</div>
		<div class="item">b</div>
		<div class="item">c</div>
	`))

	handles, err := page.QuerySelectorAll(ctx, ".item")
	must.NoError(err)
	is.Len(handles, 3)
}

// TestQuerySelectorAllReturnsEmpty verifies QuerySelectorAll returns empty for no match.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorAllReturnsEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))

	handles, err := page.QuerySelectorAll(ctx, ".nonexistent")
	must.NoError(err)
	is.Empty(handles)
}

// TestQuerySelectorAllTextContent verifies TextContent from each QuerySelectorAll result.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorAllTextContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>first</p>
		<p>second</p>
	`))

	handles, err := page.QuerySelectorAll(ctx, "p")
	must.NoError(err)
	is.Len(handles, 2)

	text0, err := handles[0].TextContent(ctx)
	must.NoError(err)
	is.Equal("first", text0)

	text1, err := handles[1].TextContent(ctx)
	must.NoError(err)
	is.Equal("second", text1)
}

// TestQuerySelectorAllEmptyEx3 verifies QuerySelectorAll returns empty slice when no match.
// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorAllEmptyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>No divs here</p>`))

	els, err := page.QuerySelectorAll(ctx, "div")
	must.NoError(err)
	is.Empty(els)
}

// Extra: no equivalent method in playwright-java; related class: TestQuerySelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestPageQuerySelectorAllEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>no items</div>`)
	must.NoError(err, "SetContent failed")

	handles, err := page.QuerySelectorAll(ctx, ".nonexistent")
	must.NoError(err, "QuerySelectorAll failed")
	is.Len(handles, 0)
}

// TestQuerySelectorShouldQueryExistingElementWithTextSelector verifies QuerySelector works with text= selector syntax.
// Ref: TestQuerySelector.java#shouldQueryExistingElementWithTextSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldQueryExistingElementWithTextSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "text='test'")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorShouldQueryExistingElementWithXpathSelector verifies QuerySelector works with xpath= selector syntax.
// Ref: TestQuerySelector.java#shouldQueryExistingElementWithXpathSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldQueryExistingElementWithXpathSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "xpath=/html/body/section")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorShouldAutoDetectXpathSelector verifies QuerySelector auto-detects XPath selectors starting with //.
// Ref: TestQuerySelector.java#shouldAutoDetectXpathSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldAutoDetectXpathSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "//html/body/section")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorShouldAutoDetectXpathSelectorWithStartingParenthesis verifies QuerySelector auto-detects XPath selectors starting with (.
// Ref: TestQuerySelector.java#shouldAutoDetectXpathSelectorWithStartingParenthesis  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldAutoDetectXpathSelectorWithStartingParenthesis(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "(//section)[1]")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorShouldAutoDetectTextSelector verifies QuerySelector auto-detects text selectors enclosed in quotes.
// Ref: TestQuerySelector.java#shouldAutoDetectTextSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldAutoDetectTextSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "'test'")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorShouldSupportSyntax verifies QuerySelector supports the >> combinator syntax.
// Ref: TestQuerySelector.java#shouldSupportSyntax  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldSupportSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div>test</div></section>`))

	el, err := page.QuerySelector(ctx, "css=section >> css=div")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorXpathShouldQueryExistingElement verifies QuerySelectorAll with xpath= returns a single match.
// Ref: TestQuerySelector.java#xpathShouldQueryExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorXpathShouldQueryExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	elements, err := page.QuerySelectorAll(ctx, "xpath=/html/body/section")
	must.NoError(err)
	is.Len(elements, 1)
	must.NotNil(elements[0])
}

// TestQuerySelectorXpathShouldReturnEmptyArrayForNonExistingElement verifies QuerySelectorAll returns empty for a non-existing XPath.
// Ref: TestQuerySelector.java#xpathShouldReturnEmptyArrayForNonExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorXpathShouldReturnEmptyArrayForNonExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	elements, err := page.QuerySelectorAll(ctx, "//html/body/non-existing-element")
	must.NoError(err)
	is.Empty(elements)
}

// TestQuerySelectorXpathShouldReturnMultipleElements verifies QuerySelectorAll with xpath= returns multiple matches.
// Ref: TestQuerySelector.java#xpathShouldReturnMultipleElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorXpathShouldReturnMultipleElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div><div></div>`))

	elements, err := page.QuerySelectorAll(ctx, "xpath=/html/body/div")
	must.NoError(err)
	is.Len(elements, 2)
}

// Ref: TestQuerySelector.java#shouldThrowForNonStringSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldThrowForNonStringSelector(t *testing.T) {
	t.Skip("skipped: Go statically typed signature requires a string; null cannot be represented as string")
}

// Ref: TestQuerySelector.java#shouldQueryExistingElementWithCssSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldQueryExistingElementWithCssSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section>test</section>"))
	element, err := page.QuerySelector(ctx, "css=section")
	must.NoError(err)
	assert.NotNil(t, element)
}

// Ref: TestQuerySelector.java#shouldReturnNullForNonExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldReturnNullForNonExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	element, err := page.QuerySelector(ctx, "non-existing-element")
	must.NoError(err)
	assert.Nil(t, element)
}

// Ref: TestQuerySelector.java#shouldAutoDetectCssSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldAutoDetectCssSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section>test</section>"))
	element, err := page.QuerySelector(ctx, "section")
	must.NoError(err)
	assert.NotNil(t, element)
}

// Ref: TestQuerySelector.java#shouldQueryExistingElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldQueryExistingElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>A</div><br/><div>B</div>"))
	elements, err := page.QuerySelectorAll(ctx, "div")
	must.NoError(err)
	must.Len(elements, 2)
	var results []any
	for _, element := range elements {
		v, err := page.Evaluate(ctx, "e => e.textContent", element)
		must.NoError(err)
		results = append(results, v)
	}
	assert.Equal(t, []any{"A", "B"}, results)
}

// Ref: TestQuerySelector.java#shouldReturnEmptyArrayIfNothingIsFound  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorShouldReturnEmptyArrayIfNothingIsFound(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	elements, err := page.QuerySelectorAll(ctx, "div")
	must.NoError(err)
	assert.Len(t, elements, 0)
}

// Ref: TestQuerySelector.java#querySelectorAllShouldWorkWithBogusArrayFrom  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestQuerySelector.java
func TestQuerySelectorQuerySelectorAllShouldWorkWithBogusArrayFrom(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>hello</div><div></div>"))
	div1, err := page.EvaluateHandle(ctx, `() => {
  Array.from = () => [];
  return document.querySelector('div');
}`)
	must.NoError(err)
	elements, err := page.QuerySelectorAll(ctx, "div")
	must.NoError(err)
	must.Len(elements, 2)
	same, err := elements[0].Evaluate(ctx, "(div, div1) => div === div1", div1)
	must.NoError(err)
	is.Equal(true, same)
}
