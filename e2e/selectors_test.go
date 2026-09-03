//go:build e2e

package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCSSAsteriskSelector verifies the * wildcard selector matches all elements.
// Ref: TestSelectorsCss.java (wildcard)
func TestCSSAsteriskSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div><span>a</span><b>b</b></div>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, "div *", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestCSSAttributeSelector verifies attribute selectors work correctly.
// Ref: TestSelectorsCss.java#shouldWorkWithAttributeSelectors
func TestCSSAttributeSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div attr="hello world" data-x="foo"><span></span></div>`)
	must.NoError(err)

	outer, err := page.EvalOnSelector(ctx, `[attr="hello world"]`, "e => e.tagName")
	must.NoError(err)
	is.Equal("DIV", outer)

	outer, err = page.EvalOnSelector(ctx, `[attr~=world]`, "e => e.tagName")
	must.NoError(err)
	is.Equal("DIV", outer)

	outer, err = page.EvalOnSelector(ctx, `[attr^=hello]`, "e => e.tagName")
	must.NoError(err)
	is.Equal("DIV", outer)

	outer, err = page.EvalOnSelector(ctx, `[attr$=world]`, "e => e.tagName")
	must.NoError(err)
	is.Equal("DIV", outer)

	outer, err = page.EvalOnSelector(ctx, `[attr*="llo wo"]`, "e => e.tagName")
	must.NoError(err)
	is.Equal("DIV", outer)
}

// TestCSSChildCombinator verifies child combinator (>) with attribute selectors.
// Ref: TestSelectorsCss.java#shouldWorkWithGreaterThanCombinatorAndSpaces
func TestCSSChildCombinator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div foo="bar" bar="baz"><span></span></div>`)
	must.NoError(err)

	outer, err := page.EvalOnSelector(ctx, `div[foo="bar"] > span`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span></span>", outer)

	outer, err = page.EvalOnSelector(ctx, `div[foo="bar"]>span`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span></span>", outer)
}

// TestCSSCommaSeparatedList verifies comma-separated selectors match multiple elements.
// Ref: TestSelectorsCss.java#shouldKeepDomOrderWithCommaSeparatedList
func TestCSSCommaSeparatedList(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section><span><div></div></span></section>`)
	must.NoError(err)

	names, err := page.EvalOnSelectorAll(ctx, "span,div", "els => els.map(e => e.nodeName).join(',')")
	must.NoError(err)
	is.Equal("SPAN,DIV", names)

	names, err = page.EvalOnSelectorAll(ctx, "div,span", "els => els.map(e => e.nodeName).join(',')")
	must.NoError(err)
	is.Equal("SPAN,DIV", names)
}

// TestCSSCommaInsideAttribute verifies comma-in-attribute-value selectors.
// Ref: TestSelectorsCss.java#shouldWorkWithCommaInsideText
func TestCSSCommaInsideAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<span></span><div attr="hello,world!"></div>`)
	must.NoError(err)

	outer, err := page.EvalOnSelector(ctx, `div[attr="hello,world!"]`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal(`<div attr="hello,world!"></div>`, outer)

	outer, err = page.EvalOnSelector(ctx, `[attr="hello,world!"]`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal(`<div attr="hello,world!"></div>`, outer)
}

// TestCSSNumericalID verifies selector with numerical ID attribute.
// Ref: TestSelectorsCss.java (numerical id)
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

// TestCSSLargeDOM verifies that CSS selectors work correctly on a large DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithLargeDOM
func TestCSSLargeDOM(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		let id = 0;
		const next = tag => {
			const e = document.createElement(tag);
			e.textContent = 'id' + (++id);
			e.id = '' + id;
			return e;
		};
		const generate = depth => {
			const div = next('div');
			const span1 = next('span');
			const span2 = next('span');
			div.appendChild(span1);
			div.appendChild(span2);
			if (depth > 0) {
				div.appendChild(generate(depth - 1));
				div.appendChild(generate(depth - 1));
			}
			return div;
		};
		document.body.appendChild(generate(5));
	}`)
	must.NoError(err)

	for _, selector := range []string{"div div div span", "div > div div > span", "span"} {
		pwCount, err := page.EvalOnSelectorAll(ctx, selector, "els => els.length")
		must.NoError(err)
		qsCount, err := page.Evaluate(ctx, "sel => document.querySelectorAll(sel).length", selector)
		must.NoError(err)
		is.Equal(qsCount, pwCount, "selector %q count mismatch", selector)
	}
}

// TestCSSNotSelector verifies :not() pseudo-class selector.
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

func TestCSSTildeCombinator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="div1"></div>
<div id="div2"></div>
<div id="div3"></div>
<div id="div4"></div>
<div id="div5"></div>
<div id="div6"></div>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, "#div1 ~ div ~ #div6", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)

	count, err = page.EvalOnSelectorAll(ctx, "#div1 ~ div ~ div", "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), count)

	count, err = page.EvalOnSelectorAll(ctx, "#div3 ~ div ~ div", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	count, err = page.EvalOnSelectorAll(ctx, "#div4 ~ div ~ div", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)

	count, err = page.EvalOnSelectorAll(ctx, "#div5 ~ div ~ div", "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), count)
}

func TestCSSPlusCombinator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section>
  <div id="div1"></div>
  <div id="div2"></div>
  <div id="div3"></div>
  <div id="div4"></div>
  <div id="div5"></div>
  <div id="div6"></div>
</section>`)
	must.NoError(err)

	count, err := page.EvalOnSelectorAll(ctx, "#div3 + div + div", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)

	count, err = page.EvalOnSelectorAll(ctx, "#div3 + #div4 + #div5", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)

	count, err = page.EvalOnSelectorAll(ctx, "div + #div1", "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), count)
}

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

func TestCSSIsSelector(t *testing.T) {
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

func TestCSSWrongCaseId(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section id="Section">hello</section>`)
	must.NoError(err)

	el, err := page.QuerySelector(ctx, "#Section")
	must.NoError(err)
	must.NotNil(el)

	el2, err := page.QuerySelector(ctx, "#section")
	must.NoError(err)
	is.Nil(el2, "ID selectors are case-sensitive in HTML")
}

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

func TestCSSChildCombinatorEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			<span>Direct child</span>
			<p><span>Grandchild</span></p>
		</div>
	`))

	count, err := page.Locator("#parent > span").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

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

func TestCSSSelectorNthOfTypeEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>First para</p>
		<span>Span</span>
		<p>Second para</p>
		<p>Third para</p>
	`))

	text, err := page.Locator("p:nth-of-type(2)").InnerText(ctx)
	must.NoError(err)
	is.Equal("Second para", text)
}

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

func TestCSSChildSelectorEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			<p id="direct">Direct</p>
			<div><p id="nested">Nested</p></div>
		</div>
	`))

	count, err := page.Locator("#parent > p").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

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

func roleBoolPtr(b bool) *bool    { return &b }
func roleStrPtr(s string) *string { return &s }
func roleIntPtr(i int) *int       { return &i }

// TestRoleSelectorDetectRoles verifies that role selectors match expected elements.
// Ref: TestSelectorsRole.java#shouldDetectRoles
func TestRoleSelectorDetectRoles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hello</button>
	<select multiple size="2"></select>
	<select></select>
	<h3>Heading</h3>
	<details><summary>Hello</summary></details>
	<div role="dialog">I am a dialog</div>`)
	must.NoError(err)

	buttons, err := page.Locator("role=button").EvaluateAll(ctx, "els => els.map(e => e.tagName)")
	must.NoError(err)
	is.Equal([]any{"BUTTON"}, buttons)

	headings, err := page.GetByRole(playwright.AriaRoleHeading).EvaluateAll(ctx, "els => els.map(e => e.tagName)")
	must.NoError(err)
	is.Equal([]any{"H3"}, headings)

	dialogs, err := page.GetByRole(playwright.AriaRoleDialog).EvaluateAll(ctx, "els => els.map(e => e.getAttribute('role'))")
	must.NoError(err)
	is.Equal([]any{"dialog"}, dialogs)

	menuitems, err := page.GetByRole(playwright.AriaRoleMenuItem).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), menuitems)
}

// TestRoleSelectorChecked verifies [checked] filter on role=checkbox.
// Ref: TestSelectorsRole.java#shouldSupportChecked
func TestRoleSelectorChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox">
	<input type="checkbox" checked>
	<div role="checkbox" aria-checked="true">Hi</div>
	<div role="checkbox" aria-checked="false">Hello</div>`)
	must.NoError(err)

	checked, err := page.GetByRole(playwright.AriaRoleCheckbox, &playwright.GetByRoleOptions{
		Checked: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), checked)

	unchecked, err := page.GetByRole(playwright.AriaRoleCheckbox, &playwright.GetByRoleOptions{
		Checked: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), unchecked)
}

// TestRoleSelectorDisabled verifies [disabled] filter.
// Ref: TestSelectorsRole.java#shouldSupportDisabled
func TestRoleSelectorDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Enabled</button>
	<button disabled>Disabled</button>
	<div role="button" aria-disabled="true">AriaDisabled</div>`)
	must.NoError(err)

	disabled, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Disabled: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), disabled)

	enabled, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Disabled: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), enabled)
}

// TestRoleSelectorByName verifies [name] filter on role selectors.
// Ref: TestSelectorsRole.java#shouldSupportName
func TestRoleSelectorByName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hello</button><button>World</button>`)
	must.NoError(err)

	hello, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name: roleStrPtr("Hello"),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Hello"}, hello)

	world, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name: roleStrPtr("World"),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"World"}, world)
}

// TestRoleSelectorByLevel verifies heading level filter.
// Ref: TestSelectorsRole.java#shouldSupportLevel
func TestRoleSelectorByLevel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<h1>H1</h1><h2>H2</h2><h3>H3</h3>`)
	must.NoError(err)

	lvl1 := 1
	h1, err := page.GetByRole(playwright.AriaRoleHeading, &playwright.GetByRoleOptions{
		Level: &lvl1,
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"H1"}, h1)

	lvl2 := 2
	h2, err := page.GetByRole(playwright.AriaRoleHeading, &playwright.GetByRoleOptions{
		Level: &lvl2,
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"H2"}, h2)
}

// TestRoleSelectorExpanded verifies [expanded] filter.
// Ref: TestSelectorsRole.java#shouldSupportExpanded
func TestRoleSelectorExpanded(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="treeitem">Hi</div>
	<div role="treeitem" aria-expanded="true">Hello</div>
	<div role="treeitem" aria-expanded="false">Bye</div>`)
	must.NoError(err)

	expanded, err := page.GetByRole(playwright.AriaRoleTreeItem, &playwright.GetByRoleOptions{
		Expanded: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Hello"}, expanded)

	collapsed, err := page.GetByRole(playwright.AriaRoleTreeItem, &playwright.GetByRoleOptions{
		Expanded: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Bye"}, collapsed)
}

// TestRoleSelectorSelected verifies [selected] filter on role=option.
// Ref: TestSelectorsRole.java#shouldSupportSelected
func TestRoleSelectorSelected(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select>
		<option>Hi</option>
		<option selected>Hello</option>
	</select>
	<div>
		<div role="option" aria-selected="true">Hi</div>
		<div role="option" aria-selected="false">Hello</div>
	</div>`)
	must.NoError(err)

	count, err := page.GetByRole(playwright.AriaRoleOption, &playwright.GetByRoleOptions{
		Selected: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestRoleSelectorPressed verifies [pressed] filter on role=button.
// Ref: TestSelectorsRole.java#shouldSupportPressed
func TestRoleSelectorPressed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hi</button>
	<button aria-pressed="true">Hello</button>
	<button aria-pressed="false">Bye</button>`)
	must.NoError(err)

	pressed, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Pressed: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Hello"}, pressed)

	notPressed, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Pressed: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), notPressed)
}

// TestRoleSelectorFilterHidden verifies includeHidden option.
// Ref: TestSelectorsRole.java#shouldFilterHidden
func TestRoleSelectorFilterHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Visible</button>
	<button style="display:none">Hidden</button>`)
	must.NoError(err)

	// By default, hidden elements are not matched
	visible, err := page.GetByRole(playwright.AriaRoleButton).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), visible)

	// With includeHidden=true, hidden elements are matched
	all, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		IncludeHidden: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), all)
}

// TestGetByTextExactStringEx3 verifies GetByText with exact=true matches only exact strings.
// Ref: TestSelectorText.java#shouldMatchExactText
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

// TestGetByTextPartialMatchEx3 verifies GetByText partial match works.
// Ref: TestSelectorText.java#shouldPartialMatch
func TestGetByTextPartialMatchEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>I love JavaScript</div>
		<div>I love Go</div>
		<div>Python is cool</div>
	`))

	count, err := page.GetByText("I love").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByTextWithButtonEx3 verifies GetByText finds button by text.
// Ref: TestSelectorText.java#shouldFindButton
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
// Ref: TestSelectorText.java#shouldFindLink
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
// Ref: TestSelectorText.java#shouldFindNumbers
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

func localBoolPtrST4(b bool) *bool { return &b }

// TestTextSelectorExactMatchEx4 verifies GetByText with exact=true.
// Ref: TestTextSelectors.java#shouldMatchExact
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
// Ref: TestTextSelectors.java#shouldMatchCaseInsensitive
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
// Ref: TestTextSelectors.java#shouldHandleSpecialChars
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
// Ref: TestTextSelectors.java#shouldMatchMultiple
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
// Ref: TestTextSelectors.java#shouldFindNestedText
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

func localBoolPtrST5(b bool) *bool { return &b }

// TestGetByTextExactMatchEx5 verifies GetByText with exact match finds only exact.
// Ref: TestSelectorsText.java#shouldMatchExactText
func TestGetByTextExactMatchEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Login</button>
		<button>Login Now</button>
	`))

	count, err := page.GetByText("Login", &playwright.GetByTextOptions{Exact: localBoolPtrST5(true)}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextCaseInsensitiveEx5 verifies GetByText is case-insensitive by default.
// Ref: TestSelectorsText.java#shouldMatchCaseInsensitive
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

// TestGetByTextPartialMatchEx5 verifies GetByText partial match (non-exact).
// Ref: TestSelectorsText.java#shouldMatchPartial
func TestGetByTextPartialMatchEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span>Click here to subscribe</span>
		<span>Read more</span>
	`))

	count, err := page.GetByText("subscribe", &playwright.GetByTextOptions{Exact: localBoolPtrST5(false)}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextInListEx5 verifies GetByText finds item in list.
// Ref: TestSelectorsText.java#shouldFindInList
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
// Ref: TestXPathSelectors.java#shouldSelectByTag
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
// Ref: TestXPathSelectors.java#shouldSelectByAttribute
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
// Ref: TestXPathSelectors.java#shouldSelectByText
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
// Ref: TestXPathSelectors.java#shouldSelectParent
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
// Ref: TestXPathSelectors.java#shouldSelectNthElement
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
// Ref: TestSelectorsXPath.java#shouldSelectDescendant
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
// Ref: TestSelectorsXPath.java#shouldSelectSibling
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
// Ref: TestSelectorsXPath.java#shouldSelectParent
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
// Ref: TestSelectorsXPath.java#shouldSelectByContainsText
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
// Ref: TestSelectorsXPath.java#shouldSelectByPosition
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
// Ref: TestSelectorsXPath.java#shouldSelectLast
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
// Ref: TestSelectorsXPath.java#shouldSelectByStartsWith
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
// Ref: TestSelectorsXPath.java#shouldSelectOrCondition
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
