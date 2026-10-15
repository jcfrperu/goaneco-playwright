//go:build e2e

// E2E tests for: TestEvalOnSelectorAll.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestEvalOnSelectorAllShouldWorkWithTextSelector verifies text= engine matches all elements with the given text.
// Ref: TestEvalOnSelectorAll.java#shouldWorkWithTextSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldWorkWithTextSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div><div>beautiful</div><div>beautiful</div><div>world!</div>`))

	val, err := page.EvalOnSelectorAll(ctx, "text='beautiful'", "divs => divs.length")
	must.NoError(err)
	is.Equal(float64(2), val)
}

// TestEvalOnSelectorAllShouldWorkWithXpathSelector verifies xpath= engine returns all matching elements.
// Ref: TestEvalOnSelectorAll.java#shouldWorkWithXpathSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldWorkWithXpathSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div><div>beautiful</div><div>world!</div>`))

	val, err := page.EvalOnSelectorAll(ctx, "xpath=/html/body/div", "divs => divs.length")
	must.NoError(err)
	is.Equal(float64(3), val)
}

// TestEvalOnSelectorAllShouldSupportSyntax verifies css=A >> css=B chaining collects all descendant matches.
// Ref: TestEvalOnSelectorAll.java#shouldSupportSyntax  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldSupportSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>hello</span></div><div>beautiful</div><div><span>wo</span><span>rld!</span></div><span>Not this one</span>`))

	val, err := page.EvalOnSelectorAll(ctx, "css=div >> css=span", "spans => spans.length")
	must.NoError(err)
	is.Equal(float64(3), val)
}

// TestEvalOnSelectorAllShouldSupportCapture verifies the * capture modifier works with EvalOnSelectorAll.
// Ref: TestEvalOnSelectorAll.java#shouldSupportCapture  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldSupportCapture(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span>a</span></div></section><section><div><span>b</span></div></section>`))

	val, err := page.EvalOnSelectorAll(ctx, "*css=div >> 'b'", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, "section >> *css=div >> 'b'", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, "section >> *", "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), val)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllWithCSSSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>A</div><div>B</div><div>C</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelectorAll(ctx, "div", "els => els.map(e => e.textContent)")
	must.NoError(err, "EvalOnSelectorAll failed")
	arr, ok := val.([]any)
	must.Truef(ok, "expected []any, got %T: %v", val, val)
	is.Len(arr, 3)
	if arr[0] != "A" || arr[1] != "B" || arr[2] != "C" {
		t.Errorf("unexpected values: %v", arr)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllReturnsEmptyArrayWhenNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>content</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelectorAll(ctx, "span", "els => els.length")
	must.NoError(err, "EvalOnSelectorAll failed")
	is.Equal(float64(0), val)
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllReturnsTextContents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>one</div>
		<div>two</div>
		<div>three</div>
	`))

	result, err := page.EvalOnSelectorAll(ctx, "div", `els => els.map(e => e.textContent.trim())`)
	must.NoError(err)

	arr, ok := result.([]any)
	is.True(ok)
	is.Len(arr, 3)
	is.Equal("one", arr[0])
	is.Equal("two", arr[1])
	is.Equal("three", arr[2])
}

// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllCountElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>a</li><li>b</li><li>c</li><li>d</li>
	`))

	result, err := page.EvalOnSelectorAll(ctx, "li", `els => els.length`)
	must.NoError(err)
	is.Equal(float64(4), result)
}

// TestEvalOnSelectorAllSumsValuesExtra verifies EvalOnSelectorAll can aggregate.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllSumsValuesExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="number" value="10">
		<input type="number" value="20">
		<input type="number" value="30">
	`))

	result, err := page.EvalOnSelectorAll(ctx, "input", `els => els.reduce((s, el) => s + parseInt(el.value), 0)`)
	must.NoError(err)
	is.Equal(float64(60), result)
}

// TestEvalOnSelectorAllMapClassNames verifies EvalOnSelectorAll maps class names.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllMapClassNames(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="a">1</div>
		<div class="b">2</div>
		<div class="c">3</div>
	`))

	result, err := page.EvalOnSelectorAll(ctx, "div", `els => els.map(el => el.className)`)
	must.NoError(err)
	must.NotNil(result)
	list, ok := result.([]any)
	is.True(ok)
	is.Len(list, 3)
	is.Equal("a", list[0])
}

// TestPageQuerySelectorAllReturnsAll verifies QuerySelectorAll returns all matches.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestPageQuerySelectorAllReturnsAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="tag">A</span>
		<span class="tag">B</span>
		<span class="tag">C</span>
	`))

	handles, err := page.QuerySelectorAll(ctx, "span.tag")
	must.NoError(err)
	is.Len(handles, 3)

	expected := []string{"A", "B", "C"}
	for i, el := range handles {
		text, err := el.InnerText(ctx)
		must.NoError(err)
		is.Equal(expected[i], text)
	}
}

// TestQuerySelectorAllInputValues verifies InputValue from QuerySelectorAll inputs.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestQuerySelectorAllInputValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input value="one">
		<input value="two">
	`))

	handles, err := page.QuerySelectorAll(ctx, "input")
	must.NoError(err)
	is.Len(handles, 2)

	v0, err := handles[0].InputValue(ctx)
	must.NoError(err)
	is.Equal("one", v0)

	v1, err := handles[1].InputValue(ctx)
	must.NoError(err)
	is.Equal("two", v1)
}

// TestQuerySelectorAllSingleMatch verifies QuerySelectorAll with single match.
// Extra: no equivalent method in playwright-java; related class: TestEvalOnSelectorAll.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestQuerySelectorAllSingleMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="title">hello</h1>`))

	handles, err := page.QuerySelectorAll(ctx, "h1")
	must.NoError(err)
	is.Len(handles, 1)

	text, err := handles[0].InnerText(ctx)
	must.NoError(err)
	is.Equal("hello", text)
}

// Ref: TestEvalOnSelectorAll.java#shouldWorkWithCssSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldWorkWithCssSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>hello</div><div>beautiful</div><div>world!</div>"))
	divsCount, err := page.EvalOnSelectorAll(ctx, "css=div", "divs => divs.length")
	must.NoError(err)
	assert.Equal(t, float64(3), divsCount)
}

// Ref: TestEvalOnSelectorAll.java#shouldAutoDetectCssSelector  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldAutoDetectCssSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>hello</div><div>beautiful</div><div>world!</div>"))
	divsCount, err := page.EvalOnSelectorAll(ctx, "div", "divs => divs.length")
	must.NoError(err)
	assert.Equal(t, float64(3), divsCount)
}

// Ref: TestEvalOnSelectorAll.java#shouldSupportCaptureWhenMultiplePathsMatch  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldSupportCaptureWhenMultiplePathsMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div><div><span></span></div></div><div></div>"))
	n, err := page.EvalOnSelectorAll(ctx, "*css=div >> span", "els => els.length")
	must.NoError(err)
	assert.Equal(t, float64(2), n)
	must.NoError(page.SetContent(ctx, "<div><div><span></span></div><span></span><span></span></div><div></div>"))
	n, err = page.EvalOnSelectorAll(ctx, "*css=div >> span", "els => els.length")
	must.NoError(err)
	assert.Equal(t, float64(2), n)
}

// Ref: TestEvalOnSelectorAll.java#shouldReturnComplexValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestEvalOnSelectorAll.java
func TestEvalOnSelectorAllShouldReturnComplexValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>hello</div><div>beautiful</div><div>world!</div>"))
	texts, err := page.EvalOnSelectorAll(ctx, "css=div", "divs => divs.map(div => div.textContent)")
	must.NoError(err)
	assert.Equal(t, []any{"hello", "beautiful", "world!"}, texts)
}
