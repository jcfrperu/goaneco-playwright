//go:build e2e

// E2E tests for: TestSelectorsMisc.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestCSSNumericalID verifies selector with numerical ID attribute.
// Ref: TestSelectorsCss.java (numerical id)
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSNumericalID(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="123">hello</div>`)
	must.NoError(err)

	text, err := page.EvalOnSelector(ctx, `[id="123"]`, "e => e.textContent")
	must.NoError(err)
	is.Equal("hello", text)
}

// TestCSSNotSelector verifies :not() pseudo-class selector.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSNotSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="a">A</div><div class="b">B</div><div class="c">C</div>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, "div:not(.a)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestCSSScope verifies scoped CSS queries work correctly.
// Ref: TestSelectorsCss.java (:scope)
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSScope(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section><div id="inner"><span>hello</span></div></section>`)
	must.NoError(err)

	el, err := page.QuerySelector(ctx, "#inner > span")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSNthChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul>
  <li>item1</li>
  <li>item2</li>
  <li>item3</li>
  <li>item4</li>
</ul>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, "li:nth-child(2)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)

	count, err = page.EvalOnSelectorAll(ctx, "li:nth-child(odd)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	count, err = page.EvalOnSelectorAll(ctx, "li:nth-child(even)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	count, err = page.EvalOnSelectorAll(ctx, "li:nth-child(2n+1)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSFirstAndLastChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul>
  <li id="a">first</li>
  <li id="b">middle</li>
  <li id="c">last</li>
</ul>`)
	must.NoError(err)

	firstID, err := page.EvalOnSelector(ctx, "li:first-child", "e => e.id")
	must.NoError(err)
	is.Equal("a", firstID)

	lastID, err := page.EvalOnSelector(ctx, "li:last-child", "e => e.id")
	must.NoError(err)
	is.Equal("c", lastID)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSHasSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="root">
  <div id="has-span"><span id="target" data-testid="foo"></span></div>
  <div id="no-span"></div>
  <div id="nested"><div><span data-testid="foo"></span></div></div>
</div>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, "div:has(span)", "els => els.length")
	must.NoError(err)
	c, _ := count.(float64)
	is.GreaterOrEqual(c, float64(2))

	count2, err := page.EvalOnSelectorAll(ctx, "div:has([data-testid=foo])", "els => els.length")
	must.NoError(err)
	c2, _ := count2.(float64)
	is.GreaterOrEqual(c2, float64(2))
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSCheckedAndDisabledPseudo(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<form>
  <input type="checkbox" id="checked" checked>
  <input type="checkbox" id="unchecked">
  <input type="text" id="disabled" disabled>
  <input type="text" id="enabled">
</form>`)
	must.NoError(err)

	checkedID, err := page.EvalOnSelector(ctx, "input:checked", "e => e.id")
	must.NoError(err)
	is.Equal("checked", checkedID)

	disabledID, err := page.EvalOnSelector(ctx, "input:disabled", "e => e.id")
	must.NoError(err)
	is.Equal("disabled", disabledID)

	enabledCount, err := page.EvalOnSelectorAll(ctx, "input:enabled", "els => els.length")
	must.NoError(err)
	is.Equal(float64(3), enabledCount)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSDescendantCombinator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section><div><span>deep</span></div><span>shallow</span></section>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, "section span", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	count, err = page.EvalOnSelectorAll(ctx, "div span", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSiblingCombinatorExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="trigger"></div>
		<span id="sibling1">sibling</span>
		<span id="sibling2">sibling2</span>
	`))

	count, err := page.Locator("#trigger ~ span").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAdjacentSiblingCombinatorExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="first"></div>
		<span id="adjacent">adjacent</span>
		<span id="far">far</span>
	`))

	count, err := page.Locator("#first + span").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)

	text, err := page.Locator("#first + span").InnerText(ctx)
	must.NoError(err)
	is.Equal("adjacent", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSMultipleClassSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="foo bar">both</div>
		<div class="foo">only foo</div>
		<div class="bar">only bar</div>
	`))

	count, err := page.Locator(".foo.bar").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSNotSelectorMatchesCorrectly(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="excluded">skip</p>
		<p class="included">include me</p>
		<p class="included">include me too</p>
	`))

	count, err := page.Locator("p:not(.excluded)").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSNthOfType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>div1</div>
		<p>para1</p>
		<div>div2</div>
		<p>para2</p>
	`))

	text, err := page.Locator("p:nth-of-type(2)").InnerText(ctx)
	must.NoError(err)
	is.Equal("para2", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSLastChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>first</li>
			<li>second</li>
			<li>last</li>
		</ul>
	`))

	text, err := page.Locator("li:last-child").InnerText(ctx)
	must.NoError(err)
	is.Equal("last", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSPseudoClassChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="c1" checked>
		<input type="checkbox" id="c2">
		<input type="checkbox" id="c3" checked>
	`))

	count, err := page.Locator("input:checked").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSPseudoClassDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="b1">Active</button>
		<button id="b2" disabled>Disabled</button>
		<button id="b3" disabled>Also disabled</button>
	`))

	count, err := page.Locator("button:disabled").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSPseudoClassEnabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="i1" type="text">
		<input id="i2" type="text" disabled>
		<input id="i3" type="text">
	`))

	count, err := page.Locator("input:enabled").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAttributeContainsSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div data-role="primary-button">A</div>
		<div data-role="secondary-button">B</div>
		<div data-role="primary-link">C</div>
	`))

	count, err := page.Locator(`[data-role*="primary"]`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAttributeStartsWithSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a href="/home">Home</a>
		<a href="/about">About</a>
		<a href="https://external.com">External</a>
	`))

	count, err := page.Locator(`a[href^="/"]`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAttributeEndsWithSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img src="photo.jpg">
		<img src="logo.png">
		<img src="banner.jpg">
	`))

	count, err := page.Locator(`img[src$=".jpg"]`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSPseudoFirstChildEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="first">First</li>
			<li>Second</li>
			<li>Third</li>
		</ul>
	`))

	count, err := page.Locator("li:first-child").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSPseudoLastChildEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>First</li>
			<li>Second</li>
			<li class="last">Third</li>
		</ul>
	`))

	count, err := page.Locator("li:last-child").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSPseudoNthChildEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>One</li>
			<li>Two</li>
			<li>Three</li>
			<li>Four</li>
		</ul>
	`))

	count, err := page.Locator("li:nth-child(2)").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)

	text, err := page.Locator("li:nth-child(2)").InnerText(ctx)
	must.NoError(err)
	is.Equal("Two", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSDescendantCombinatorEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			<span>Inside</span>
			<p><span>Deep inside</span></p>
		</div>
		<span>Outside</span>
	`))

	count, err := page.Locator("#parent span").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAdjacentSiblingEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h2>Heading</h2>
		<p>Paragraph after heading</p>
		<p>Another paragraph</p>
	`))

	count, err := page.Locator("h2 + p").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorAttributePrefixEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a href="https://example.com">External</a>
		<a href="/internal">Internal</a>
	`))

	count, err := page.Locator(`a[href^="https://"]`).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorAttributeSuffixEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img src="photo.jpg">
		<img src="icon.png">
		<img src="banner.jpg">
	`))

	count, err := page.Locator(`img[src$=".jpg"]`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorAttributeContainsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="primary-button">Click</div>
		<div class="secondary-text">Text</div>
		<div class="primary-header">Header</div>
	`))

	count, err := page.Locator(`[class*="primary"]`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorNotPseudoEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="t1">
		<input type="checkbox" id="c1">
		<input type="text" id="t2">
	`))

	count, err := page.Locator(`input:not([type="checkbox"])`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorGeneralSiblingEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>
			<h2>Heading</h2>
			<p>Para 1</p>
			<span>Span</span>
			<p>Para 2</p>
		</div>
	`))

	count, err := page.Locator(`h2 ~ p`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorMultipleClassesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="btn primary">Primary</div>
		<div class="btn secondary">Secondary</div>
		<div class="primary">Just Primary</div>
	`))

	count, err := page.Locator(`.btn.primary`).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorFirstOfTypeEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>First div</div>
		<span>Span</span>
		<div>Second div</div>
	`))

	text, err := page.Locator("div:first-of-type").InnerText(ctx)
	must.NoError(err)
	is.Equal("First div", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorLastOfTypeEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>First div</div>
		<div>Middle div</div>
		<div>Last div</div>
		<span>Span</span>
	`))

	text, err := page.Locator("div:last-of-type").InnerText(ctx)
	must.NoError(err)
	is.Equal("Last div", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorEmptyPseudoEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="empty"></div>
		<div class="nonempty">Content</div>
		<div class="empty"></div>
	`))

	count, err := page.Locator("div:empty").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSelectorCheckedPseudoEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="c1" checked>
		<input type="checkbox" id="c2">
		<input type="checkbox" id="c3" checked>
	`))

	count, err := page.Locator("input:checked").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSFirstOfTypeEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div><p>First para</p><p>Second para</p><span>A span</span></div>
	`))

	text, err := page.Locator("p:first-of-type").TextContent(ctx)
	must.NoError(err)
	is.Equal("First para", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSLastOfTypeEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div><p>First para</p><p>Second para</p><p>Third para</p></div>
	`))

	text, err := page.Locator("p:last-of-type").TextContent(ctx)
	must.NoError(err)
	is.Equal("Third para", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSNthChildEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
			<li>Item 3</li>
		</ul>
	`))

	text, err := page.Locator("li:nth-child(2)").TextContent(ctx)
	must.NoError(err)
	is.Equal("Item 2", text)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSEmptyPseudoEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="empty"></div>
		<div id="full">Content</div>
	`))

	count, err := page.Locator("div:empty").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSCheckedPseudoEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="c1" checked>
		<input type="checkbox" id="c2">
		<input type="checkbox" id="c3" checked>
	`))

	count, err := page.Locator("input:checked").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSFocusPseudoEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp1" type="text">
		<input id="inp2" type="text">
	`))

	must.NoError(page.Locator("#inp2").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.querySelector(':focus').id`)
	must.NoError(err)
	is.Equal("inp2", focused)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSDisabledPseudoEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="enabled">
		<input type="text" id="disabled1" disabled>
		<input type="text" id="disabled2" disabled>
	`))

	count, err := page.Locator("input:disabled").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSEnabledPseudoEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="b1">Active</button>
		<button id="b2" disabled>Disabled</button>
		<button id="b3">Active 2</button>
	`))

	count, err := page.Locator("button:enabled").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAdjacentSiblingEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="inp">Name:</label>
		<input id="inp" type="text">
	`))

	el, err := page.QuerySelector(ctx, `label + input`)
	must.NoError(err)
	must.NotNil(el)

	id, err := el.GetAttribute(ctx, "id")
	must.NoError(err)
	is.Equal("inp", id)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSSiblingGeneralEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>
			<h2 id="h">Heading</h2>
			<p>Para 1</p>
			<p>Para 2</p>
		</div>
	`))

	count, err := page.Locator("h2 ~ p").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSNotSelectorEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="active">Active</li>
			<li>Inactive 1</li>
			<li>Inactive 2</li>
		</ul>
	`))

	count, err := page.Locator("li:not(.active)").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAttributeContainsEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a id="a1" href="/path/to/page">Link 1</a>
		<a id="a2" href="/other">Link 2</a>
		<a id="a3" href="/path/more">Link 3</a>
	`))

	count, err := page.Locator("a[href*='/path']").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSAttributeEndsWithEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img src="photo.jpg" id="i1">
		<img src="icon.png" id="i2">
		<img src="banner.jpg" id="i3">
	`))

	count, err := page.Locator("img[src$='.jpg']").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestCSSOnlyChildEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d1"><p>Sole child</p></div>
		<div id="d2"><p>First</p><p>Second</p></div>
	`))

	count, err := page.Locator("p:only-child").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextExactStringEx3 verifies GetByText with exact=true matches only exact strings.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestGetByTextExactStringEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Hello World</p>
		<p>Hello</p>
	`))

	count, err := page.GetByText("Hello").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByTextWithButtonEx3 verifies GetByText finds button by text.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestGetByTextWithButtonEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Save Changes</button>
		<button>Cancel</button>
	`))

	count, err := page.GetByText("Save Changes").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextWithLinkEx3 verifies GetByText finds anchor link by text.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestGetByTextWithLinkEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav>
			<a href="/">Home</a>
			<a href="/about">About Us</a>
			<a href="/contact">Contact</a>
		</nav>
	`))

	count, err := page.GetByText("About Us").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextNumbersEx3 verifies GetByText can find numeric content.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestGetByTextNumbersEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<td>42</td>
		<td>100</td>
		<td>42</td>
	`))

	count, err := page.GetByText("42").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestTextSelectorExactMatchEx4 verifies GetByText with exact=true.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestTextSelectorExactMatchEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Buy</button>
		<button>Buy Now</button>
	`))

	count, err := page.GetByText("Buy", &playwright.GetByTextOptions{Exact: localBoolPtrST4(true)}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestTextSelectorCaseInsensitiveEx4 verifies GetByText matches case-insensitively.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestTextSelectorCaseInsensitiveEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Hello World</p>
		<p>HELLO WORLD</p>
	`))

	count, err := page.GetByText("hello world", &playwright.GetByTextOptions{Exact: localBoolPtrST4(false)}).Count(ctx)
	must.NoError(err)
	is.GreaterOrEqual(count, 1)
}

// TestTextSelectorWithSpecialCharsEx4 verifies GetByText handles special chars.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestTextSelectorWithSpecialCharsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="s">Price: $9.99</span>`))

	count, err := page.GetByText("$9.99", &playwright.GetByTextOptions{Exact: localBoolPtrST4(false)}).Count(ctx)
	must.NoError(err)
	is.GreaterOrEqual(count, 1)
}

// TestTextSelectorMultipleMatchesEx4 verifies GetByText can match multiple elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestTextSelectorMultipleMatchesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>Item</li>
		<li>Item</li>
		<li>Other</li>
	`))

	count, err := page.GetByText("Item", &playwright.GetByTextOptions{Exact: localBoolPtrST4(true)}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestTextSelectorNestedEx4 verifies GetByText finds text in nested elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestTextSelectorNestedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>
			<article>
				<p>Article content here</p>
			</article>
		</div>
	`))

	count, err := page.GetByText("Article content", &playwright.GetByTextOptions{Exact: localBoolPtrST4(false)}).Count(ctx)
	must.NoError(err)
	is.GreaterOrEqual(count, 1)
}

// TestGetByTextCaseInsensitiveEx5 verifies GetByText is case-insensitive by default.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestGetByTextCaseInsensitiveEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>HELLO WORLD</p>`))

	count, err := page.GetByText("hello world", &playwright.GetByTextOptions{Exact: localBoolPtrST5(false)}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextInListEx5 verifies GetByText finds item in list.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestGetByTextInListEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Apple</li>
			<li>Banana</li>
			<li>Cherry</li>
		</ul>
	`))

	text, err := page.GetByText("Banana", &playwright.GetByTextOptions{Exact: localBoolPtrST5(true)}).TextContent(ctx)
	must.NoError(err)
	is.Equal("Banana", text)
}

// TestXPathSelectsByTagNameEx3 verifies XPath selects elements by tag name.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathSelectsByTagNameEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
			<li>Item 3</li>
		</ul>
	`))

	count, err := page.Locator(`xpath=//li`).Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestXPathSelectsByAttributeEx3 verifies XPath selects elements by attribute.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathSelectsByAttributeEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="t1">
		<input type="submit" id="s1">
	`))

	count, err := page.Locator(`xpath=//input[@type='text']`).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestXPathSelectsByTextEx3 verifies XPath selects elements by text content.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathSelectsByTextEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Hello World</p>
		<p>Goodbye World</p>
	`))

	count, err := page.Locator(`xpath=//p[contains(text(),'Hello')]`).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestXPathSelectsParentEx3 verifies XPath can navigate to parent.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathSelectsParentEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			<span id="child">Child</span>
		</div>
	`))

	text, err := page.Locator(`xpath=//span[@id='child']/..`).GetAttribute(ctx, "id")
	must.NoError(err)
	is.Equal("parent", text)
}

// TestXPathSelectsNthEx3 verifies XPath [n] position selector.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathSelectsNthEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ol>
			<li>First</li>
			<li>Second</li>
			<li>Third</li>
		</ol>
	`))

	text, err := page.Locator(`xpath=//li[2]`).InnerText(ctx)
	must.NoError(err)
	is.Equal("Second", text)
}

// TestXPathDescendantEx4 verifies XPath descendant axis selection.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathDescendantEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="container">
			<section><p id="target">Target paragraph</p></section>
		</div>
	`))

	text, err := page.Locator("xpath=//div[@id='container']//p").TextContent(ctx)
	must.NoError(err)
	is.Equal("Target paragraph", text)
}

// TestXPathSiblingEx4 verifies XPath following-sibling axis.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathSiblingEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li id="first">First</li>
			<li id="second">Second</li>
		</ul>
	`))

	text, err := page.Locator("xpath=//li[@id='first']/following-sibling::li").TextContent(ctx)
	must.NoError(err)
	is.Equal("Second", text)
}

// TestXPathParentEx4 verifies XPath parent axis.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathParentEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent"><span id="child">Child text</span></div>
	`))

	tag, err := page.Evaluate(ctx, `() => document.querySelector('#child').parentElement.id`)
	must.NoError(err)
	is.Equal("parent", tag)
}

// TestXPathContainsTextEx4 verifies XPath contains() text function.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathContainsTextEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Buy now</p>
		<p>Read more</p>
		<p>Buy later</p>
	`))

	count, err := page.Locator("xpath=//p[contains(text(),'Buy')]").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestXPathPositionEx5 verifies XPath position() function.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathPositionEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
			<li>Item 3</li>
		</ul>
	`))

	text, err := page.Locator("xpath=//li[position()=2]").TextContent(ctx)
	must.NoError(err)
	is.Equal("Item 2", text)
}

// TestXPathLastEx5 verifies XPath last() function.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathLastEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ol>
			<li>First</li>
			<li>Middle</li>
			<li>Last item</li>
		</ol>
	`))

	text, err := page.Locator("xpath=//li[last()]").TextContent(ctx)
	must.NoError(err)
	is.Equal("Last item", text)
}

// TestXPathAttributeStartsWithEx5 verifies XPath starts-with() function.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathAttributeStartsWithEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="user-name" type="text">
		<input id="user-email" type="email">
		<input id="password" type="password">
	`))

	count, err := page.Locator("xpath=//input[starts-with(@id,'user')]").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestXPathOrConditionEx5 verifies XPath or condition.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestXPathOrConditionEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p id="a">Alpha</p>
		<p id="b">Beta</p>
		<p id="c">Gamma</p>
	`))

	count, err := page.Locator("xpath=//p[@id='a' or @id='c']").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestSelectorsMiscShouldWorkWithLayoutSelectors verifies :right-of, :left-of, :above, :below, and :near CSS pseudo-selectors.
// Ref: TestSelectorsMisc.java#shouldWorkWithLayoutSelectors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestSelectorsMiscShouldWorkWithLayoutSelectors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<container style='width: 500px; height: 500px; position: relative;'></container>`)
	must.NoError(err)

	setupScript := `(container, boxes) => {
    for (let i = 0; i < boxes.length; i++) {
      const div = document.createElement('div');
      div.style.position = 'absolute';
      div.style.overflow = 'hidden';
      div.style.boxSizing = 'border-box';
      div.style.border = '1px solid black';
      div.id = 'id' + i;
      div.textContent = 'id' + i;
      const box = boxes[i];
      div.style.left = box[0] + 'px';
      // Note: top is a flipped y coordinate
      div.style.top = (250 - box[1] - box[3]) + 'px';
      div.style.width = box[2] + 'px';
      div.style.height = box[3] + 'px';
      container.appendChild(div);
      const span = document.createElement('span');
      span.textContent = '' + i;
      div.appendChild(span);
    }
  }`

	boxes := [][]any{
		{0, 0, 150, 150},
		{100, 200, 50, 50},
		{200, 200, 50, 50},
		{100, 150, 50, 50},
		{201, 150, 50, 50},
		{200, 100, 50, 50},
		{50, 50, 50, 50},
		{150, 50, 50, 50},
		{150, -51, 50, 50},
		{201, -101, 50, 50},
	}

	_, err = page.EvalOnSelector(ctx, "container", setupScript, boxes)
	must.NoError(err)

	// :right-of assertions
	val, err := page.EvalOnSelector(ctx, "div:right-of(#id6)", "e => e.id")
	must.NoError(err)
	is.Equal("id7", val)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id1)", "e => e.id")
	must.NoError(err)
	is.Equal("id2", val)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id3)", "e => e.id")
	must.NoError(err)
	is.Equal("id4", val)

	el, err := page.QuerySelector(ctx, "div:right-of(#id4)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id0)", "e => e.id")
	must.NoError(err)
	is.Equal("id7", val)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id9", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:right-of(#id3)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id4,id2,id5,id7,id8,id9", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:right-of(#id3, 50)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id2,id5,id7,id8", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:right-of(#id3, 49)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id7,id8", val)

	// :left-of assertions
	val, err = page.EvalOnSelector(ctx, "div:left-of(#id2)", "e => e.id")
	must.NoError(err)
	is.Equal("id1", val)

	el, err = page.QuerySelector(ctx, "div:left-of(#id0)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelector(ctx, "div:left-of(#id5)", "e => e.id")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelector(ctx, "div:left-of(#id9)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	val, err = page.EvalOnSelector(ctx, "div:left-of(#id4)", "e => e.id")
	must.NoError(err)
	is.Equal("id3", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:left-of(#id5)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id7,id3,id1,id6,id8", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:left-of(#id5, 3)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id7,id8", val)

	// :above assertions
	val, err = page.EvalOnSelector(ctx, "div:above(#id0)", "e => e.id")
	must.NoError(err)
	is.Equal("id3", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id5)", "e => e.id")
	must.NoError(err)
	is.Equal("id4", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id7)", "e => e.id")
	must.NoError(err)
	is.Equal("id5", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id9)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	el, err = page.QuerySelector(ctx, "div:above(#id2)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelectorAll(ctx, "div:above(#id5)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id4,id2,id3,id1", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:above(#id5, 20)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id4,id3", val)

	// :below assertions
	val, err = page.EvalOnSelector(ctx, "div:below(#id4)", "e => e.id")
	must.NoError(err)
	is.Equal("id5", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id3)", "e => e.id")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id2)", "e => e.id")
	must.NoError(err)
	is.Equal("id4", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id6)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id7)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id9", val)

	el, err = page.QuerySelector(ctx, "div:below(#id9)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelectorAll(ctx, "div:below(#id3)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id5,id6,id7,id8,id9", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:below(#id3, 105)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id5,id6,id7", val)

	// :near assertions
	val, err = page.EvalOnSelector(ctx, "div:near(#id0)", "e => e.id")
	must.NoError(err)
	is.Equal("id3", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id7)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id5,id3,id6", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id0)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id3,id6,id7,id8,id1,id5", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id6)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id3,id7", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id6, 10)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id0, 100)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id3,id6,id7,id8,id1,id5,id4,id2", val)

	// Combined layout
	val, err = page.EvalOnSelectorAll(ctx, "div:below(#id5):above(#id8)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id7,id6", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id5):above(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id7", val)

	// Error cases
	_, nearErr := page.QuerySelector(ctx, ":near(50)")
	is.Error(nearErr)
	is.ErrorContains(nearErr, `"near" engine expects a selector list`)

	_, leftOfErr := page.QuerySelector(ctx, "div >> left-of=abc")
	is.Error(leftOfErr)
	is.ErrorContains(leftOfErr, "Malformed selector: left-of=abc")
}

// TestSelectorsMiscShouldWorkWithInternalHasNot verifies internal:has-not selector filters sections without matching descendants.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalHasNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestSelectorsMiscShouldWorkWithInternalHasNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section><span></span><div></div></section><section><br></section>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `section >> internal:has-not="span"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="span, div, br"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="br"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="span, div"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="article"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), val)
}

// TestSelectorsMiscShouldWorkWithInternalAnd verifies internal:and selector intersects element sets.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalAnd  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestSelectorsMiscShouldWorkWithInternalAnd(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class=foo>hello</div><div class=bar>world</div><span class=foo>hello2</span><span class=bar>world2</span>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `div >> internal:and="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	arr, ok := val.([]any)
	is.True(ok)
	is.Empty(arr)

	val, err = page.EvalOnSelectorAll(ctx, `div >> internal:and=".foo"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `div >> internal:and=".bar"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"world"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `span >> internal:and="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello2", "world2"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `.foo >> internal:and="div"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `.bar >> internal:and="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"world2"}, val)
}

// TestSelectorsMiscShouldWorkWithInternalOr verifies internal:or selector unions element sets.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalOr  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestSelectorsMiscShouldWorkWithInternalOr(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>hello</div><span>world</span>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `div >> internal:or="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello", "world"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `span >> internal:or="div"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello", "world"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `article >> internal:or="something"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), val)

	text, err := page.Locator(`article >> internal:or="div"`).TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text)

	text, err = page.Locator(`article >> internal:or="span"`).TextContent(ctx)
	must.NoError(err)
	is.Equal("world", text)
}

// TestSelectorsMiscShouldWorkWithInternalChain verifies internal:chain selector scopes results within a parent match.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalChain  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsMisc.java
func TestSelectorsMiscShouldWorkWithInternalChain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>one <span>two</span> <button>three</button> </div><span>four</span><button>five</button>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `div >> internal:chain="button"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"three"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `div >> internal:chain="span >> internal:or=\"button\""`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"two", "three"}, val)
}
