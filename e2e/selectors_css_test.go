//go:build e2e

// E2E tests for: TestSelectorsCss.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestSelectorsCssShouldWorkWithLargeDOM verifies CSS selectors work correctly on a large DOM tree.
// Ref: TestSelectorsCss.java#shouldWorkWithLargeDOM  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithLargeDOM(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithGreaterThanCombinatorAndSpaces verifies the > combinator with various spacing.
// Ref: TestSelectorsCss.java#shouldWorkWithGreaterThanCombinatorAndSpaces  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithGreaterThanCombinatorAndSpaces(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithCommaSeparatedList verifies comma-separated selector lists pierce shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithCommaSeparatedList  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithCommaSeparatedList(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithCommaInsideText verifies attribute values containing commas don't confuse the parser.
// Ref: TestSelectorsCss.java#shouldWorkWithCommaInsideText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithCommaInsideText(t *testing.T) {
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

// TestSelectorsCssShouldNotMatchRootAfterGreaterGreaterThan verifies >> doesn't match the root element.
// Ref: TestSelectorsCss.java#shouldNotMatchRootAfterGreaterGreaterThan  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldNotMatchRootAfterGreaterGreaterThan(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithNumericalId verifies CSS escaped numerical IDs work.
// Ref: TestSelectorsCss.java#shouldWorkWithNumericalId  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithNumericalId(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithWrongCaseId verifies CSS ID selectors match case-insensitively.
// Ref: TestSelectorsCss.java#shouldWorkWithWrongCaseId  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithWrongCaseId(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithAsterisk verifies the wildcard * selector.
// Ref: TestSelectorsCss.java#shouldWorkWithAsterisk  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithAsterisk(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithColonNthChild verifies :nth-child pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonNthChild  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithColonNthChild(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithColonNot verifies :not pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithColonNot(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithTilde verifies the ~ general sibling combinator.
// Ref: TestSelectorsCss.java#shouldWorkWithTilde  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithTilde(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithPlus verifies the + adjacent sibling combinator.
// Ref: TestSelectorsCss.java#shouldWorkWithPlus  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithPlus(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithSpacesInColonNthChildAndColonNot verifies spaces inside :nth-child() and :not() arguments.
// Ref: TestSelectorsCss.java#shouldWorkWithSpacesInColonNthChildAndColonNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithSpacesInColonNthChildAndColonNot(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithColonIs verifies :is() pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonIs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithColonIs(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithColonHas verifies :has() pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonHas  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithColonHas(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithColonScope verifies :scope pseudo-class with shadow DOM.
// Ref: TestSelectorsCss.java#shouldWorkWithColonScope  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithColonScope(t *testing.T) {
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

// TestSelectorsCssShouldWorkForOpenShadowRoots verifies that CSS selectors pierce open shadow roots for element lookup.
// Ref: TestSelectorsCss.java#shouldWorkForOpenShadowRoots  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkForOpenShadowRoots(t *testing.T) {
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

// Ref: TestSelectorsCss.java#shouldWorkWithCommaSeparatedListInVariousPosition  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithCommaSeparatedListInVariousPosition(t *testing.T) {
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

// TestSelectorsCssShouldWorkWithAttributeSelectors verifies attribute selectors work correctly.
// Ref: TestSelectorsCss.java#shouldWorkWithAttributeSelectors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldWorkWithAttributeSelectors(t *testing.T) {
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

// TestSelectorsCssShouldKeepDomOrderWithCommaSeparatedList verifies comma-separated selectors match multiple elements.
// Ref: TestSelectorsCss.java#shouldKeepDomOrderWithCommaSeparatedList  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
func TestSelectorsCssShouldKeepDomOrderWithCommaSeparatedList(t *testing.T) {
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

// TestCSSAsteriskSelector verifies the * wildcard selector matches all elements.
// Ref: TestSelectorsCss.java (wildcard)
// Extra: no equivalent method in playwright-java; related class: TestSelectorsCss.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsCss.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsCss.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsCss.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsCss.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsCss.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsCss.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsCss.java
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
