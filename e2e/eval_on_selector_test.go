//go:build e2e

// E2E tests for: TestEvalOnSelector.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorWithCSSSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section id="testAttribute">43543</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "css=section", "e => e.id")
	must.NoError(err, "EvalOnSelector failed")
	is.Equal("testAttribute", val)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorWithIDSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section id="testAttribute">43543</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "#testAttribute", "e => e.id")
	must.NoError(err, "EvalOnSelector failed")
	is.Equal("testAttribute", val)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorAcceptsArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section>hello</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "section", "(e, suffix) => e.textContent + suffix", " world")
	must.NoError(err, "EvalOnSelector failed")
	is.Equal("hello world", val)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorReturnsComplexValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section>hello</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "section", "e => ({ tag: e.tagName, text: e.textContent })")
	must.NoError(err, "EvalOnSelector failed")
	m, ok := val.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", val, val)
	if m["tag"] != "SECTION" {
		t.Errorf("tag = %v, want 'SECTION'", m["tag"])
	}
	if m["text"] != "hello" {
		t.Errorf("text = %v, want 'hello'", m["text"])
	}
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorThrowsWhenNoElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section>value</section>`)
	must.NoError(err, "SetContent failed")

	_, err = page.EvalOnSelector(ctx, "div", "e => e.id")
	is.Error(err)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorReturnsUndefinedForMissingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id="testAttribute">43543</section>`))

	result, err := page.EvalOnSelector(ctx, "css=div", `e => e.id`)
	must.NoError(err)
	is.Nil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorReturnsInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>hello</section>`))

	result, err := page.EvalOnSelector(ctx, "css=section", `e => e.innerText`)
	must.NoError(err)
	is.Equal("hello", result)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorWithAttributeSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div data-testid="box">content</div>`))

	result, err := page.EvalOnSelector(ctx, `[data-testid="box"]`, `e => e.textContent`)
	must.NoError(err)
	is.Equal("content", result)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorWithArgument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class="myDiv">original</div>`))

	result, err := page.EvalOnSelector(ctx, ".myDiv", `(el, value) => el.textContent + value`, " suffix")
	must.NoError(err)
	is.Equal("original suffix", result)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorReturnsBooleanTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" checked>`))

	result, err := page.EvalOnSelector(ctx, "input", `e => e.checked`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorReturnsBooleanFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox">`))

	result, err := page.EvalOnSelector(ctx, "input", `e => e.checked`)
	must.NoError(err)
	is.Equal(false, result)
}

// TestEvalOnSelectorGetsBoundingRect verifies EvalOnSelector can query layout.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorGetsBoundingRect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:100px;height:50px">box</div>`))

	result, err := page.EvalOnSelector(ctx, "div", `el => el.getBoundingClientRect().width`)
	must.NoError(err)
	is.Equal(float64(100), result)
}

// TestEvalOnSelectorSetsAttribute verifies EvalOnSelector can modify DOM.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorSetsAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">text</div>`))

	_, err := page.EvalOnSelector(ctx, "#el", `el => el.setAttribute('data-x', 'yes')`)
	must.NoError(err)

	val, err := page.Locator("#el").GetAttribute(ctx, "data-x")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("yes", *val)
}

// TestEvalOnSelectorReturnsNthChildText verifies EvalOnSelector works on nth child.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorReturnsNthChildText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>first</li>
			<li>second</li>
			<li>third</li>
		</ul>
	`))

	result, err := page.EvalOnSelector(ctx, "li:nth-child(2)", `el => el.textContent`)
	must.NoError(err)
	is.Equal("second", result)
}

// TestPageQuerySelectorFindsById verifies QuerySelector works with ID selector.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestPageQuerySelectorFindsById(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">found</div>`))

	el, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(el, "QuerySelector should find element by id")

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("found", text)
}

// TestPageQuerySelectorFindsByClass verifies QuerySelector works with class selector.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestPageQuerySelectorFindsByClass(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">first</div>
		<div class="item">second</div>
	`))

	// QuerySelector returns the FIRST match
	el, err := page.QuerySelector(ctx, ".item")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("first", text)
}

// TestQuerySelectorFindsButtonEx2 verifies QuerySelector finds button by id.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestQuerySelectorFindsButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Submit</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorReturnsNilForMissingEx2 verifies QuerySelector returns nil for no match.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestQuerySelectorReturnsNilForMissingEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>No button here</div>`))

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err)
	is.Nil(el)
}

// TestQuerySelectorByTagEx2 verifies QuerySelector finds element by tag.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelector.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestQuerySelectorByTagEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h2>A heading</h2>`))

	el, err := page.QuerySelector(ctx, "h2")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("A heading", text)
}

// TestEvalOnSelectorShouldWorkWithDataTestSelector verifies the data-test= attribute selector engine.
// Ref: TestEvalOnSelector.java#shouldWorkWithDataTestSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithDataTestSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section data-test=foo id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "data-test=foo", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorShouldWorkWithDataTestidSelector verifies the data-testid= attribute selector engine.
// Ref: TestEvalOnSelector.java#shouldWorkWithDataTestidSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithDataTestidSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section data-testid=foo id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "data-testid=foo", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorShouldWorkWithDataTestIdSelector verifies the data-test-id= attribute selector engine.
// Ref: TestEvalOnSelector.java#shouldWorkWithDataTestIdSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithDataTestIdSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section data-test-id=foo id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "data-test-id=foo", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorShouldWorkWithTextSelector1 verifies text= selector with a single-quoted value.
// Ref: TestEvalOnSelector.java#shouldWorkWithTextSelector1  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithTextSelector1(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "text='43543'", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorShouldWorkWithXpathSelector verifies the xpath= selector engine prefix.
// Ref: TestEvalOnSelector.java#shouldWorkWithXpathSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithXpathSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "xpath=/html/body/section", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorShouldWorkWithTextSelector2 verifies text= selector without quotes does a substring match.
// Ref: TestEvalOnSelector.java#shouldWorkWithTextSelector2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithTextSelector2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "text=43543", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorShouldAutoDetectCssSelectorWithAttributes verifies CSS attribute selectors are auto-detected without an engine prefix.
// Ref: TestEvalOnSelector.java#shouldAutoDetectCssSelectorWithAttributes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldAutoDetectCssSelectorWithAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "section[id='testAttribute']", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorShouldAutoDetectNestedSelectors verifies >> chains with auto-detected CSS and quoted text selectors.
// Ref: TestEvalOnSelector.java#shouldAutoDetectNestedSelectors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldAutoDetectNestedSelectors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div foo=bar><section>43543<span>Hello<div id=target></div></span></section></div>`))

	val, err := page.EvalOnSelector(ctx, "div[foo=bar] > section >> 'Hello' >> div", "e => e.id")
	must.NoError(err)
	is.Equal("target", val)
}

// TestEvalOnSelectorShouldAcceptElementHandlesAsArguments verifies an ElementHandle can be passed as an extra argument.
// Ref: TestEvalOnSelector.java#shouldAcceptElementHandlesAsArguments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldAcceptElementHandlesAsArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>hello</section><div> world</div>`))

	divHandle, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(divHandle)

	val, err := page.EvalOnSelector(ctx, "section", "(e, div) => e.textContent + div.textContent", divHandle)
	must.NoError(err)
	is.Equal("hello world", val)
}

// TestEvalOnSelectorShouldSupportSyntax verifies the css=A >> css=B chained selector syntax.
// Ref: TestEvalOnSelector.java#shouldSupportSyntax  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldSupportSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div>hello</div></section>`))

	val, err := page.EvalOnSelector(ctx, "css=section >> css=div", "(e, suffix) => e.textContent + suffix", " world!")
	must.NoError(err)
	is.Equal("hello world!", val)
}

// TestEvalOnSelectorShouldSupportSyntaxWithDifferentEngines verifies mixed xpath/css/text engines chained with >>.
// Ref: TestEvalOnSelector.java#shouldSupportSyntaxWithDifferentEngines  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldSupportSyntaxWithDifferentEngines(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span>hello</span></div></section>`))

	val, err := page.EvalOnSelector(ctx, "xpath=/html/body/section >> css=div >> text='hello'", "(e, suffix) => e.textContent + suffix", " world!")
	must.NoError(err)
	is.Equal("hello world!", val)
}

// TestEvalOnSelectorShouldSupportSpacesWithSyntax verifies extra whitespace around >> and engine= prefixes is ignored.
// Ref: TestEvalOnSelector.java#shouldSupportSpacesWithSyntax  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldSupportSpacesWithSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	err := page.Goto(ctx, srv.Prefix()+"/deep-shadow.html")
	must.NoError(err)

	val, err := page.EvalOnSelector(ctx, " css = div >>css=div>>css   = span  ", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)
}

// TestEvalOnSelectorShouldNotStopAtFirstFailureWithSyntax verifies >> chaining scans all siblings, not just the first.
// Ref: TestEvalOnSelector.java#shouldNotStopAtFirstFailureWithSyntax  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldNotStopAtFirstFailureWithSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>Next</span><button>Previous</button><button>Next</button></div>`))

	val, err := page.EvalOnSelector(ctx, "button >> 'Next'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<button>Next</button>", val)
}

// TestEvalOnSelectorShouldSupportCapture verifies the * capture modifier returns the captured element rather than the final match.
// Ref: TestEvalOnSelector.java#shouldSupportCapture  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldSupportCapture(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span>a</span></div></section><section><div><span>b</span></div></section>`))

	val, err := page.EvalOnSelector(ctx, "*css=div >> 'b'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<div><span>b</span></div>", val)

	val, err = page.EvalOnSelector(ctx, "section >> *css=div >> 'b'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<div><span>b</span></div>", val)

	val, err = page.EvalOnSelector(ctx, "css=div >> *text='b'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span>b</span>", val)
}

// TestEvalOnSelectorShouldThrowOnMultipleCaptures verifies that using * on more than one part of a chain returns an error.
// Ref: TestEvalOnSelector.java#shouldThrowOnMultipleCaptures  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldThrowOnMultipleCaptures(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span></span></div></section>`))

	_, err := page.EvalOnSelector(ctx, "*css=div >> *css=span", "e => e.outerHTML")
	is.Error(err)
	is.Contains(err.Error(), "Only one of the selectors can capture using * modifier")
}

// TestEvalOnSelectorShouldThrowOnMalformedCapture verifies that *=engine syntax without a valid engine name returns an error.
// Ref: TestEvalOnSelector.java#shouldThrowOnMalformedCapture  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldThrowOnMalformedCapture(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div></div></section>`))

	_, err := page.EvalOnSelector(ctx, "*=div", "e => e.outerHTML")
	is.Error(err)
	is.Contains(err.Error(), "Unknown engine")
}

// Ref: TestEvalOnSelector.java#shouldWorkWithCssSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithCssSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="root">
  <div id="inner1" class="child">
    <span id="target"></span>
  </div>
  <div id="inner2" class="child"></div>
</div>
<span id="standalone"></span>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, ":is(div, span)", "els => els.length")
	must.NoError(err)

	c, _ := count.(float64)
	is.GreaterOrEqual(c, float64(4))

	count2, err := page.EvalOnSelectorAll(ctx, "div:is(.child)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count2)
}

// Ref: TestEvalOnSelector.java#shouldWorkWithIdSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithIdSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section id='testAttribute'>43543</section>"))
	idAttribute, err := page.EvalOnSelector(ctx, "id=testAttribute", "e => e.id")
	must.NoError(err)
	assert.Equal(t, "testAttribute", idAttribute)
}

// Ref: TestEvalOnSelector.java#shouldAutoDetectCssSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldAutoDetectCssSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section id='testAttribute'>43543</section>"))
	idAttribute, err := page.EvalOnSelector(ctx, "section", "e => e.id")
	must.NoError(err)
	assert.Equal(t, "testAttribute", idAttribute)
}

// Ref: TestEvalOnSelector.java#shouldAcceptArguments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldAcceptArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section>hello</section>"))
	text, err := page.EvalOnSelector(ctx, "section", "(e, suffix) => e.textContent + suffix", " world!")
	must.NoError(err)
	assert.Equal(t, "hello world!", text)
}

// Ref: TestEvalOnSelector.java#shouldThrowErrorIfNoElementIsFound  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldThrowErrorIfNoElementIsFound(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.EvalOnSelector(ctx, "section", "e => e.id")
	must.Error(err)
	assert.Contains(t, err.Error(), "Failed to find element matching selector \"section\"")
}

// Ref: TestEvalOnSelector.java#shouldWorkWithSpacesInCssAttributes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithSpacesInCssAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div><input placeholder='Select date'></div>"))
	for _, sel := range []string{
		`[placeholder="Select date"]`,
		`[placeholder='Select date']`,
		`input[placeholder="Select date"]`,
		`input[placeholder='Select date']`,
	} {
		must.NoError(page.Locator(sel).WaitFor(ctx), sel)
		el, err := page.QuerySelector(ctx, sel)
		must.NoError(err)
		is.NotNil(el, sel)
	}
	for _, sel := range []string{
		`[placeholder="Select date"]`,
		`[placeholder='Select date']`,
		`input[placeholder="Select date"]`,
		`input[placeholder='Select date']`,
		`css=[placeholder="Select date"]`,
		`css=[placeholder='Select date']`,
		`css=input[placeholder="Select date"]`,
		`css=input[placeholder='Select date']`,
		`div >> [placeholder="Select date"]`,
		`div >> [placeholder='Select date']`,
	} {
		html, err := page.EvalOnSelector(ctx, sel, "e => e.outerHTML")
		must.NoError(err, sel)
		is.Equal(`<input placeholder="Select date">`, html, sel)
	}
}

// Ref: TestEvalOnSelector.java#shouldWorkWithQuotesInCssAttributes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithQuotesInCssAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	check := func(html string, selectors ...string) {
		must.NoError(page.SetContent(ctx, html))
		for _, sel := range selectors {
			el, err := page.QuerySelector(ctx, sel)
			must.NoError(err, sel)
			must.NotNil(el, sel)
		}
	}
	check(`<div><input placeholder="Select&quot;date"></div>`, `[placeholder="Select\"date"]`, `[placeholder='Select"date']`)
	check(`<div><input placeholder="Select &quot; date"></div>`, `[placeholder="Select \" date"]`, `[placeholder='Select " date']`)
	check(`<div><input placeholder="Select&apos;date"></div>`, `[placeholder="Select'date"]`, `[placeholder='Select\'date']`)
	check(`<div><input placeholder="Select &apos; date"></div>`, `[placeholder="Select ' date"]`, `[placeholder='Select \' date']`)
}

// Ref: TestEvalOnSelector.java#shouldWorkWithSpacesInCssAttributesWhenMissing  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithSpacesInCssAttributesWhenMissing(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	el, err := page.QuerySelector(ctx, "[placeholder='Select date']")
	must.NoError(err)
	assert.Nil(t, el)
	must.NoError(page.SetContent(ctx, "<div><input placeholder='Select date'></div>"))
	must.NoError(page.Locator("[placeholder='Select date']").WaitFor(ctx))
}

// Ref: TestEvalOnSelector.java#shouldWorkWithQuotesInCssAttributesWhenMissing  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldWorkWithQuotesInCssAttributesWhenMissing(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	el, err := page.QuerySelector(ctx, `[placeholder='Select\"date']`)
	must.NoError(err)
	assert.Nil(t, el)
	must.NoError(page.SetContent(ctx, "<div><input placeholder='Select&quot;date'></div>"))
	must.NoError(page.Locator(`[placeholder='Select\"date']`).WaitFor(ctx))
}

// Ref: TestEvalOnSelector.java#shouldReturnComplexValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelector.java
func TestEvalOnSelectorShouldReturnComplexValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section id='testAttribute'>43543</section>"))
	idAttribute, err := page.EvalOnSelector(ctx, "css=section", "e => [{ id: e.id }]")
	must.NoError(err)
	assert.Equal(t, []any{map[string]any{"id": "testAttribute"}}, idAttribute)
}
