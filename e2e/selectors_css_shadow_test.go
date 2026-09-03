//go:build e2e

package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCSSLargeDOMPerformance verifies CSS selectors work correctly on a large DOM tree.
// Ref: TestSelectorsCss.java#shouldWorkWithLargeDOM
func TestCSSLargeDOMPerformance(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		let id = 0;
		const next = (tag) => {
			const e = document.createElement(tag);
			const eid = ++id;
			e.textContent = 'id' + eid;
			e.id = '' + eid;
			return e;
		};
		const generate = (depth) => {
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
		document.body.appendChild(generate(12));
	}`)
	must.NoError(err)

	selectors := []string{
		"div div div span",
		"div > div div > span",
		"div + div div div span + span",
		"div ~ div div > span ~ span",
		"div > div > div + div > div + div > span ~ span",
		"div div div div div div div div div div span",
		"div > div > div > div > div > div > div > div > div > div > span",
		"div ~ div div ~ div div ~ div div ~ div div ~ div span",
		"span",
	}

	for _, sel := range selectors {
		pwCount, err := page.EvalOnSelectorAll(ctx, sel, "els => els.length")
		must.NoError(err, "EvalOnSelectorAll failed for selector %q", sel)
		qsCount, err := page.Evaluate(ctx, "selector => document.querySelectorAll(selector).length", sel)
		must.NoError(err, "Evaluate querySelectorAll failed for selector %q", sel)
		is.Equal(qsCount, pwCount, "count mismatch for selector %q", sel)
	}
}

// TestCSSGreaterThanCombinatorSpaces verifies the > combinator with various spacing.
// Ref: TestSelectorsCss.java#shouldWorkWithGreaterThanCombinatorAndSpaces
func TestCSSGreaterThanCombinatorSpaces(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div foo="bar" bar="baz"><span></span></div>`))

	cases := []string{
		`div[foo="bar"] > span`,
		`div[foo="bar"]> span`,
		`div[foo="bar"] >span`,
		`div[foo="bar"]>span`,
		`div[foo="bar"]   >    span`,
		`div[foo="bar"]>    span`,
		`div[foo="bar"]     >span`,
		`div[foo="bar"][bar="baz"] > span`,
		`div[foo="bar"][bar="baz"]> span`,
		`div[foo="bar"][bar="baz"] >span`,
		`div[foo="bar"][bar="baz"]>span`,
		`div[foo="bar"][bar="baz"]   >    span`,
		`div[foo="bar"][bar="baz"]>    span`,
		`div[foo="bar"][bar="baz"]     >span`,
	}
	for _, sel := range cases {
		val, err := page.EvalOnSelector(ctx, sel, "e => e.outerHTML")
		must.NoError(err, "EvalOnSelector failed for %q", sel)
		is.Equal("<span></span>", val, "outerHTML mismatch for %q", sel)
	}
}

// TestCSSCommaSeparatedListShadow verifies comma-separated selector lists pierce shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithCommaSeparatedList
func TestCSSCommaSeparatedListShadow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	count, err := page.EvalOnSelectorAll(ctx, "css=span,section #root1", "els => els.length")
	must.NoError(err)
	is.Equal(float64(5), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=section #root1, div span", "els => els.length")
	must.NoError(err)
	is.Equal(float64(5), count)

	id, err := page.EvalOnSelector(ctx, "css=doesnotexist , section #root1", "e => e.id")
	must.NoError(err)
	is.Equal("root1", id)

	count, err = page.EvalOnSelectorAll(ctx, "css=doesnotexist ,section #root1", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=span,div span", "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=span,div span,div div span", "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), count)

	count, err = page.EvalOnSelectorAll(ctx, `css=#target,[attr="value\\ space"]`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	count, err = page.EvalOnSelectorAll(ctx, `css=#target,[data-testid="foo"],[attr="value\\ space"]`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), count)

	count, err = page.EvalOnSelectorAll(ctx, `css=#target,[data-testid="foo"],[attr="value\\ space"],span`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), count)
}

// TestCSSCommaInsideText verifies attribute values containing commas don't confuse the parser.
// Ref: TestSelectorsCss.java#shouldWorkWithCommaInsideText
func TestCSSCommaInsideText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span></span><div attr="hello,world!"></div>`))

	for _, sel := range []string{
		`css=div[attr="hello,world!"]`,
		`css=[attr="hello,world!"]`,
		`css=div[attr='hello,world!']`,
		`css=[attr='hello,world!']`,
	} {
		val, err := page.EvalOnSelector(ctx, sel, "e => e.outerHTML")
		must.NoError(err, "EvalOnSelector failed for %q", sel)
		is.Equal(`<div attr="hello,world!"></div>`, val, "mismatch for %q", sel)
	}

	val, err := page.EvalOnSelector(ctx, `css=div[attr="hello,world!"],span`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span></span>", val)
}

// TestCSSNotMatchRootAfterGreaterGreater verifies >> doesn't match the root element.
// Ref: TestSelectorsCss.java#shouldNotMatchRootAfterGreaterGreaterThan
func TestCSSNotMatchRootAfterGreaterGreater(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div>test</div></section>`))

	el, err := page.QuerySelector(ctx, "css=section >> css=section")
	must.NoError(err)
	is.Nil(el)
}

// TestCSSNumericalIDEscaped verifies CSS escaped numerical IDs work.
// Ref: TestSelectorsCss.java#shouldWorkWithNumericalId
func TestCSSNumericalIDEscaped(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id="123"></section>`))

	el, err := page.QuerySelector(ctx, `#\31\32\33`)
	must.NoError(err)
	is.NotNil(el)
}

// TestCSSWrongCaseID verifies CSS ID selectors match case-insensitively.
// Ref: TestSelectorsCss.java#shouldWorkWithWrongCaseId
func TestCSSWrongCaseID(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id="Hello"></section>`))

	for _, sel := range []string{"#Hello", "#hello", "#HELLO", "#helLO"} {
		val, err := page.EvalOnSelector(ctx, sel, "e => e.tagName")
		must.NoError(err, "EvalOnSelector failed for %q", sel)
		is.Equal("SECTION", val, "tagName mismatch for %q", sel)
	}
}

// TestCSSAsterisk verifies the wildcard * selector.
// Ref: TestSelectorsCss.java#shouldWorkWithAsterisk
func TestCSSAsterisk(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id=div1></div><div id=div2><span><span></span></span></div>`))

	cases := []struct {
		sel      string
		expected float64
	}{
		{"*", 7},
		{"*#div1", 1},
		{"*:not(#div1)", 6},
		{"*:not(div)", 5},
		{"*:not(span)", 5},
		{"*:not(*)", 0},
		{"*:is(*)", 7},
		{"* *", 6},
		{"* *:not(span)", 4},
		{"div > *", 1},
		{"div *", 2},
		{"* > *", 6},
	}
	for _, tc := range cases {
		count, err := page.EvalOnSelectorAll(ctx, tc.sel, "els => els.length")
		must.NoError(err, "EvalOnSelectorAll failed for %q", tc.sel)
		is.Equal(tc.expected, count, "count mismatch for %q", tc.sel)
	}
}

// TestCSSNthChildShadow verifies :nth-child pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonNthChild
func TestCSSNthChildShadow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	cases := []struct {
		sel      string
		expected float64
	}{
		{"css=span:nth-child(odd)", 3},
		{"css=span:nth-child(even)", 1},
		{"css=span:nth-child(n+1)", 4},
		{"css=span:nth-child(n+2)", 1},
		{"css=span:nth-child(2n)", 1},
		{"css=span:nth-child(2n+1)", 3},
		{"css=span:nth-child(-n)", 0},
		{"css=span:nth-child(-n+1)", 3},
		{"css=span:nth-child(-n+2)", 4},
		{"css=span:nth-child(23n+2)", 1},
	}
	for _, tc := range cases {
		count, err := page.EvalOnSelectorAll(ctx, tc.sel, "els => els.length")
		must.NoError(err, "EvalOnSelectorAll failed for %q", tc.sel)
		is.Equal(tc.expected, count, "count mismatch for %q", tc.sel)
	}
}

// TestCSSNot verifies :not pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonNot
func TestCSSNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	count, err := page.EvalOnSelectorAll(ctx, "css=div:not(#root1)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=body :not(span)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=div > :not(span):not(div)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), count)
}

// TestCSSTilde verifies the ~ general sibling combinator.
// Ref: TestSelectorsCss.java#shouldWorkWithTilde
func TestCSSTilde(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id=div1></div>
<div id=div2></div>
<div id=div3></div>
<div id=div4></div>
<div id=div5></div>
<div id=div6></div>`))

	cases := []struct {
		sel      string
		expected float64
	}{
		{"css=#div1 ~ div ~ #div6", 1},
		{"css=#div1 ~ div ~ div", 4},
		{"css=#div3 ~ div ~ div", 2},
		{"css=#div4 ~ div ~ div", 1},
		{"css=#div5 ~ div ~ div", 0},
		{"css=#div3 ~ #div2 ~ #div6", 0},
		{"css=#div3 ~ #div4 ~ #div5", 1},
	}
	for _, tc := range cases {
		count, err := page.EvalOnSelectorAll(ctx, tc.sel, "els => els.length")
		must.NoError(err, "EvalOnSelectorAll failed for %q", tc.sel)
		is.Equal(tc.expected, count, "count mismatch for %q", tc.sel)
	}
}

// TestCSSPlus verifies the + adjacent sibling combinator.
// Ref: TestSelectorsCss.java#shouldWorkWithPlus
func TestCSSPlus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>
  <div id=div1></div>
  <div id=div2></div>
  <div id=div3></div>
  <div id=div4></div>
  <div id=div5></div>
  <div id=div6></div>
</section>`))

	cases := []struct {
		sel      string
		expected float64
	}{
		{"css=#div1 ~ div + #div6", 1},
		{"css=#div1 ~ div + div", 4},
		{"css=#div3 + div + div", 1},
		{"css=#div4 ~ #div5 + div", 1},
		{"css=#div5 + div + div", 0},
		{"css=#div3 ~ #div2 + #div6", 0},
		{"css=#div3 + #div4 + #div5", 1},
		{"css=div + #div1", 0},
		{"css=section > div + div ~ div", 4},
		{"css=section > div + #div4 ~ div", 2},
		{"css=section:has(:scope > div + #div2)", 1},
		{"css=section:has(:scope > div + #div1)", 0},
	}
	for _, tc := range cases {
		count, err := page.EvalOnSelectorAll(ctx, tc.sel, "els => els.length")
		must.NoError(err, "EvalOnSelectorAll failed for %q", tc.sel)
		is.Equal(tc.expected, count, "count mismatch for %q", tc.sel)
	}
}

// TestCSSSpacesInPseudoClasses verifies spaces inside :nth-child() and :not() arguments.
// Ref: TestSelectorsCss.java#shouldWorkWithSpacesInColonNthChildAndColonNot
func TestCSSSpacesInPseudoClasses(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	cases := []struct {
		sel      string
		expected float64
	}{
		{"css=span:nth-child(23n +2)", 1},
		{"css=span:nth-child(23n+ 2)", 1},
		{"css=span:nth-child( 23n + 2 )", 1},
		{"css=span:not(#root1 #target)", 3},
		{"css=span:not(:not(#root1 #target))", 1},
		{"css=span:not(span:not(#root1 #target))", 1},
		{"css=div > :not(span)", 2},
		{"css=body :not(span, div)", 1},
		{"css=span, section:not(span, div)", 5},
	}
	for _, tc := range cases {
		count, err := page.EvalOnSelectorAll(ctx, tc.sel, "els => els.length")
		must.NoError(err, "EvalOnSelectorAll failed for %q", tc.sel)
		is.Equal(tc.expected, count, "count mismatch for %q", tc.sel)
	}
}

// TestCSSIs verifies :is() pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonIs
func TestCSSIs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	cases := []struct {
		sel      string
		expected float64
	}{
		{"css=div:is(#root1)", 1},
		{"css=div:is(#root1, #target)", 1},
		{"css=div:is(span, #target)", 0},
		{"css=div:is(span, #root1 > *)", 2},
		{"css=div:is(section div)", 3},
		{"css=:is(div, span)", 7},
		{"css=section:is(section) div:is(section div)", 3},
		{"css=:is(div, span) > *", 6},
		{"css=#root1:has(:is(#root1))", 0},
		{"css=#root1:has(:is(:scope, #root1))", 1},
	}
	for _, tc := range cases {
		count, err := page.EvalOnSelectorAll(ctx, tc.sel, "els => els.length")
		must.NoError(err, "EvalOnSelectorAll failed for %q", tc.sel)
		is.Equal(tc.expected, count, "count mismatch for %q", tc.sel)
	}
}

// TestCSSHas verifies :has() pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonHas
func TestCSSHas(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	count, err := page.EvalOnSelectorAll(ctx, "css=div:has(#target)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=div:has([data-testid=foo])", "els => els.length")
	must.NoError(err)
	is.Equal(float64(3), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=div:has([attr*=value])", "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestCSSScopeShadow verifies :scope pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonScope
func TestCSSScopeShadow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	count, err := page.EvalOnSelectorAll(ctx, "css=div:is(:scope#root1)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=div:is(:scope #root1)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)

	count, err = page.EvalOnSelectorAll(ctx, "css=div:has(:scope > #target)", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), count)
}

// TestCSSOpenShadowRoots verifies that CSS selectors pierce open shadow roots for element lookup.
// Ref: TestSelectorsCss.java#shouldWorkForOpenShadowRoots
func TestCSSOpenShadowRoots(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/deep-shadow.html"))

	val, err := page.EvalOnSelector(ctx, "css=span", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root1", val)

	val, err = page.EvalOnSelector(ctx, `css=[attr="value\\ space"]`, "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root3 #2", val)

	val, err = page.EvalOnSelector(ctx, `css=[attr='value\\ \\space']`, "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root3 #2", val)

	val, err = page.EvalOnSelector(ctx, "css=div div span", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)

	val, err = page.EvalOnSelector(ctx, `css=div span + span`, "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root3 #2", val)

	val, err = page.EvalOnSelector(ctx, `css=span + [attr*="value"]`, "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root3 #2", val)

	val, err = page.EvalOnSelector(ctx, `css=[data-testid="foo"] + [attr*="value"]`, "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root3 #2", val)

	val, err = page.EvalOnSelector(ctx, "css=#target", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)

	val, err = page.EvalOnSelector(ctx, "css=div #target", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)

	val, err = page.EvalOnSelector(ctx, "css=div div #target", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)

	el, err := page.QuerySelector(ctx, "css=div div div #target")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelector(ctx, "css=section > div div span", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)

	val, err = page.EvalOnSelector(ctx, "css=section > div div span:nth-child(2)", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root3 #2", val)

	el, err = page.QuerySelector(ctx, "css=section div div div div")
	must.NoError(err)
	is.Nil(el)

	root2, err := page.QuerySelector(ctx, "css=div div")
	must.NoError(err)
	must.NotNil(root2)

	target, err := root2.QuerySelector(ctx, "css=#target")
	must.NoError(err)
	must.NotNil(target)

	val, err = target.Evaluate(ctx, "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)

	lightEl, err := root2.QuerySelector(ctx, "css:light=#target")
	must.NoError(err)
	is.Nil(lightEl)
}

// Ref: TestSelectorsCss.java#shouldWorkWithCommaSeparatedListInVariousPosition
func TestShouldWorkWithCommaSeparatedListInVariousPosition(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><span><div><x></x><y></y></div></span></section>`))

	eval := func(sel, expr string) string {
		v, err := page.EvalOnSelectorAll(ctx, sel, expr)
		must.NoError(err)
		return v.(string)
	}
	is.Equal("X,Y", eval("css=span,div >> css=x,y", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("X", eval("css=span,div >> css=x", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("X,Y", eval("css=div >> css=x,y", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("X", eval("css=div >> css=x", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("X", eval("css=section >> css=div >> css=x", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("Y", eval("css=section >> css=span >> css=div >> css=y", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("X,Y", eval("css=section >> css=div >> css=x,y", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("X,Y", eval("css=section >> css=div,span >> css=x,y", "els => els.map(e => e.nodeName).join(',')"))
	is.Equal("X,Y", eval("css=section >> css=span >> css=x,y", "els => els.map(e => e.nodeName).join(',')"))
}
