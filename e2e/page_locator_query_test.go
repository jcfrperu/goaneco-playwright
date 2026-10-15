//go:build e2e

// E2E tests for: TestPageLocatorQuery.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

// TestLocatorFilterByHasNotText verifies that Filter(HasNotText) excludes items by text.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorFilterByHasNotText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<ul>
			<li>In stock</li>
			<li>Out of stock</li>
			<li>In stock</li>
		</ul>
	`)
	must.NoError(err, "SetContent failed")

	items := page.Locator("li")
	out := "Out"
	inStock := items.Filter(&playwright.LocatorFilterOptions{HasNotText: &out})

	count, err := inStock.Count(ctx)
	must.NoError(err, "Filter(HasNotText).Count() failed")
	is.Equal(2, count)
}

// TestLocatorFilterByHasNotTextCount verifies Filter(HasNotText) returns the right count.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorFilterByHasNotTextCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Active</li>
			<li>Inactive</li>
			<li>Active</li>
		</ul>
	`))

	exclude := "Inactive"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{
		HasNotText: &exclude,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterByHasChild verifies Filter(Has) with child locator.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorFilterByHasChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="card">
			<span class="badge">New</span>
			<p>Product A</p>
		</div>
		<div class="card">
			<p>Product B</p>
		</div>
		<div class="card">
			<span class="badge">Sale</span>
			<p>Product C</p>
		</div>
	`))

	badge := page.Locator(".badge")
	count, err := page.Locator(".card").Filter(&playwright.LocatorFilterOptions{
		Has: badge,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterByHasNotChild verifies Filter(HasNot) excludes elements with child.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorFilterByHasNotChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">
			<button class="action">Buy</button>
			Product X
		</div>
		<div class="item">Product Y</div>
		<div class="item">Product Z</div>
	`))

	button := page.Locator(".action")
	count, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{
		HasNot: button,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorSelectorById verifies page.Locator finds element by ID.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorSelectorById(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">found</div>`))

	text, err := page.Locator("#target").InnerText(ctx)
	must.NoError(err)
	is.Equal("found", text)
}

// TestLocatorSelectorByClass verifies page.Locator finds element by class.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorSelectorByClass(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span class="highlight">span text</span>`))

	text, err := page.Locator(".highlight").InnerText(ctx)
	must.NoError(err)
	is.Equal("span text", text)
}

// TestLocatorSelectorByAttribute verifies page.Locator finds element by attribute.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorSelectorByAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a data-test="mylink" href="#">link</a>`))

	text, err := page.Locator("[data-test='mylink']").InnerText(ctx)
	must.NoError(err)
	is.Equal("link", text)
}

// TestLocatorSelectorByTagName verifies page.Locator finds element by tag.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorSelectorByTagName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h2>heading two</h2>`))

	text, err := page.Locator("h2").InnerText(ctx)
	must.NoError(err)
	is.Equal("heading two", text)
}

// TestLocatorCSSDescendantCombinator verifies Locator with descendant combinator.
// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorCSSDescendantCombinator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<article>
			<p>inside</p>
		</article>
		<p>outside</p>
	`))

	count, err := page.Locator("article p").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)

	text, err := page.Locator("article p").InnerText(ctx)
	must.NoError(err)
	is.Equal("inside", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageLocatorQuery.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestLocatorPseudoFirstChild(t *testing.T) {
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

	text, err := page.Locator("li:first-child").InnerText(ctx)
	must.NoError(err)
	is.Equal("first", text)
}

// TestPageLocatorQueryShouldFilterByRegex verifies Filter(HasTextRegex) narrows to elements matching the pattern.
// Ref: TestPageLocatorQuery.java#shouldFilterByRegex  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Foobar</div><div>Bar</div>`))

	re := regexp.MustCompile(`Foo.*`)
	tc, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re}).TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal("Foobar", *tc)
}

// TestPageLocatorQueryShouldFilterByTextWithQuotes verifies Filter(HasText) works with quoted strings.
// Ref: TestPageLocatorQuery.java#shouldFilterByTextWithQuotes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByTextWithQuotes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Hello "world"</div><div>Hello world</div>`))

	needle := `Hello "world"`
	tc, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasText: &needle}).TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal(`Hello "world"`, *tc)
}

// TestPageLocatorQueryShouldFilterByRegexWithQuotes verifies Filter(HasTextRegex) works with quoted patterns.
// Ref: TestPageLocatorQuery.java#shouldFilterByRegexWithQuotes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByRegexWithQuotes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Hello "world"</div><div>Hello world</div>`))

	re := regexp.MustCompile(`Hello "world"`)
	tc, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re}).TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal(`Hello "world"`, *tc)
}

// TestPageLocatorQueryShouldFilterByRegexAndRegexpFlags verifies Filter(HasTextRegex) with case-insensitive flag.
// Ref: TestPageLocatorQuery.java#shouldFilterByRegexAndRegexpFlags  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByRegexAndRegexpFlags(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Hello "world"</div><div>Hello world</div>`))

	re := regexp.MustCompile(`(?i)hElLo "world"`)
	tc, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re}).TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal(`Hello "world"`, *tc)
}

// TestPageLocatorQueryShouldFilterByCaseInsensitiveRegexInAChild verifies case-insensitive regex filter across child elements.
// Ref: TestPageLocatorQuery.java#shouldFilterByCaseInsensitiveRegexInAChild  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByCaseInsensitiveRegexInAChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class="test"><h5>Title Text</h5></div>`))

	re := regexp.MustCompile(`(?i)^title text$`)
	assertions := playwright.Expect(page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re}))
	must.NoError(assertions.ToHaveText(ctx, "Title Text"))
}

// TestPageLocatorQueryShouldFilterByCaseInsensitiveRegexInMultipleChildren verifies regex matching across nested children.
// Ref: TestPageLocatorQuery.java#shouldFilterByCaseInsensitiveRegexInMultipleChildren  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByCaseInsensitiveRegexInMultipleChildren(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div class=\"test\"><h5>Title</h5> <h2><i>Text</i></h2></div>"))

	re := regexp.MustCompile(`(?i)^title text$`)
	assertions := playwright.Expect(page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re}))
	must.NoError(assertions.ToHaveClass(ctx, "test"))
}

// TestPageLocatorQueryShouldSupportLocatorAnd verifies And() narrows to elements matching both locators.
// Ref: TestPageLocatorQuery.java#shouldSupportLocatorAnd  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldSupportLocatorAnd(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div data-testid="foo">hello</div><div data-testid="bar">world</div>
    <span data-testid="foo">hello2</span><span data-testid="bar">world2</span>`))

	// div.and(div) → 2 divs
	count, err := page.Locator("div").And(page.Locator("div")).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)

	// div.and(testid=foo) → only "hello"
	andFoo := page.Locator("div").And(page.GetByTestId("foo"))
	must.NoError(playwright.Expect(andFoo).ToHaveTextArray(ctx, []string{"hello"}))

	// div.and(testid=bar) → only "world"
	andBar := page.Locator("div").And(page.GetByTestId("bar"))
	must.NoError(playwright.Expect(andBar).ToHaveTextArray(ctx, []string{"world"}))

	// testid=foo.and(div) → only "hello"
	fooAndDiv := page.GetByTestId("foo").And(page.Locator("div"))
	must.NoError(playwright.Expect(fooAndDiv).ToHaveTextArray(ctx, []string{"hello"}))

	// testid=bar.and(span) → only "world2"
	barAndSpan := page.GetByTestId("bar").And(page.Locator("span"))
	must.NoError(playwright.Expect(barAndSpan).ToHaveTextArray(ctx, []string{"world2"}))
}

// TestPageLocatorQueryShouldSupportLocatorOr verifies Or() selects elements from either locator.
// Ref: TestPageLocatorQuery.java#shouldSupportLocatorOr  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldSupportLocatorOr(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div><span>world</span>`))

	// div.or(span) → 2 elements
	divOrSpan := page.Locator("div").Or(page.Locator("span"))
	count, err := divOrSpan.Count(ctx)
	must.NoError(err)
	is.Equal(2, count)

	// div.or(span) has texts hello and world
	must.NoError(playwright.Expect(divOrSpan).ToHaveTextArray(ctx, []string{"hello", "world"}))

	// triple or: span.or(article).or(div)
	tripleOr := page.Locator("span").Or(page.Locator("article")).Or(page.Locator("div"))
	must.NoError(playwright.Expect(tripleOr).ToHaveTextArray(ctx, []string{"hello", "world"}))

	// article.or(something) → 0
	noMatch, err := page.Locator("article").Or(page.Locator("something")).Count(ctx)
	must.NoError(err)
	is.Equal(0, noMatch)

	// article.or(div) → "hello"
	must.NoError(playwright.Expect(page.Locator("article").Or(page.Locator("div"))).ToHaveText(ctx, "hello"))

	// article.or(span) → "world"
	must.NoError(playwright.Expect(page.Locator("article").Or(page.Locator("span"))).ToHaveText(ctx, "world"))
}

// TestPageLocatorQueryShouldRespectFirstAndLast verifies First() and Last() narrow the match set correctly.
// Ref: TestPageLocatorQuery.java#shouldRespectFirstAndLast  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldRespectFirstAndLast(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>
    <div><p>A</p></div>
    <div><p>A</p><p>A</p></div>
    <div><p>A</p><p>A</p><p>A</p></div>
  </section>`))

	total, err := page.Locator("div >> p").Count(ctx)
	must.NoError(err)
	is.Equal(6, total)

	total2, err := page.Locator("div").Locator("p").Count(ctx)
	must.NoError(err)
	is.Equal(6, total2)

	firstCount, err := page.Locator("div").First().Locator("p").Count(ctx)
	must.NoError(err)
	is.Equal(1, firstCount)

	lastCount, err := page.Locator("div").Last().Locator("p").Count(ctx)
	must.NoError(err)
	is.Equal(3, lastCount)
}

// TestPageLocatorQueryShouldRespectNth verifies Nth() selects the element at a given zero-based index.
// Ref: TestPageLocatorQuery.java#shouldRespectNth  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldRespectNth(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>
    <div><p>A</p></div>
    <div><p>A</p><p>A</p></div>
    <div><p>A</p><p>A</p><p>A</p></div>
  </section>`))

	c0, err := page.Locator("div >> p").Nth(0).Count(ctx)
	must.NoError(err)
	is.Equal(1, c0)

	c1, err := page.Locator("div").Nth(1).Locator("p").Count(ctx)
	must.NoError(err)
	is.Equal(2, c1)

	c2, err := page.Locator("div").Nth(2).Locator("p").Count(ctx)
	must.NoError(err)
	is.Equal(3, c2)
}

// TestPageLocatorQueryShouldThrowOnDueToStrictness verifies that calling a single-element method on a
// multi-match locator produces a strict mode violation error.
// Ref: TestPageLocatorQuery.java#shouldThrowOnDueToStrictness  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldThrowOnDueToStrictness(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>A</div><div>B</div>`))

	_, err := page.Locator("div").IsVisible(ctx)
	is.Error(err, "expected strict mode error for multi-match locator")
	is.Contains(err.Error(), "strict mode violation")
}

// TestPageLocatorQueryShouldThrowOnDueToStrictness2 verifies strict mode violation for evaluate on multi-match locator.
// Ref: TestPageLocatorQuery.java#shouldThrowOnDueToStrictness2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldThrowOnDueToStrictness2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select><option>One</option><option>Two</option></select>`))

	_, err := page.Locator("option").Evaluate(ctx, "e => {}")
	is.Error(err, "expected strict mode error for multi-match option locator")
	is.Contains(err.Error(), "strict mode violation")
}

// TestPageLocatorQueryShouldFilterByText verifies Filter(HasText) narrows to elements containing the text.
// Ref: TestPageLocatorQuery.java#shouldFilterByText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Foobar</div><div>Bar</div>`))

	needle := "Foo"
	tc, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasText: &needle}).TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal("Foobar", *tc)
}

// TestPageLocatorQueryShouldFilterByText2 verifies Filter(HasText) with text nested inside a span.
// Ref: TestPageLocatorQuery.java#shouldFilterByText2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByText2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>foo <span>hello world</span> bar</div>`))

	needle := "hello world"
	tc, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasText: &needle}).TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Contains(*tc, "hello world")
}

// TestPageLocatorQueryShouldFilterByTextWithAmpersand verifies Filter(HasText) with an ampersand in text.
// Ref: TestPageLocatorQuery.java#shouldFilterByTextWithAmpersand  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByTextWithAmpersand(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Save &amp; Continue</div>`))

	needle := "Save & Continue"
	tc, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasText: &needle}).TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Equal("Save & Continue", *tc)
}

// TestPageLocatorQueryShouldSupportHasLocator verifies Filter(Has) selects elements that contain a child locator.
// Ref: TestPageLocatorQuery.java#shouldSupportHasLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldSupportHasLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>hello</span></div><div><span>world</span></div>`))

	// div:has(text=world) should find 1
	worldSpan := page.Locator("text=world")
	count, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{Has: worldSpan}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)

	// div:has(text=hello) should find 1
	helloSpan := page.Locator("text=hello")
	count2, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{Has: helloSpan}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count2)

	// div:has(span) should find 2
	span := page.Locator("span")
	count3, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{Has: span}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count3)
}

// TestPageLocatorQueryShouldSupportLocatorFilter verifies chained Filter options including HasText, HasNot, HasNotText.
// Ref: TestPageLocatorQuery.java#shouldSupportLocatorFilter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldSupportLocatorFilter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span>hello</span></div><div><span>world</span></div></section>`))

	hello := "hello"
	world := "world"

	// filter(hasText:"hello") → 1
	c1, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasText: &hello}).Count(ctx)
	must.NoError(err)
	is.Equal(1, c1)

	// filter(hasNotText:"hello") → 1 (world)
	c2, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasNotText: &hello}).Count(ctx)
	must.NoError(err)
	is.Equal(1, c2)

	// filter(hasText:"hello").filter(hasText:"world") → 0
	c3, err := page.Locator("div").
		Filter(&playwright.LocatorFilterOptions{HasText: &hello}).
		Filter(&playwright.LocatorFilterOptions{HasText: &world}).
		Count(ctx)
	must.NoError(err)
	is.Equal(0, c3)

	// section filter(hasText:"hello").filter(hasText:"world") → 1
	c4, err := page.Locator("section").
		Filter(&playwright.LocatorFilterOptions{HasText: &hello}).
		Filter(&playwright.LocatorFilterOptions{HasText: &world}).
		Count(ctx)
	must.NoError(err)
	is.Equal(1, c4)

	// div.filter(has:span.filter(hasText:"world")) → 1
	spanWorld := page.Locator("span").Filter(&playwright.LocatorFilterOptions{HasText: &world})
	c5, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{Has: spanWorld}).Count(ctx)
	must.NoError(err)
	is.Equal(1, c5)

	// div.filter(hasNot:span.filter(hasText:"world")) → 1
	c6, err := page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasNot: spanWorld}).Count(ctx)
	must.NoError(err)
	is.Equal(1, c6)
}

// TestPageLocatorQueryShouldSupportLocatorLocatorWithAndOr verifies And() and Or() combined with nested locators.
// Ref: TestPageLocatorQuery.java#shouldSupportLocatorLocatorWithAndOr  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldSupportLocatorLocatorWithAndOr(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
    <div>one <span>two</span> <button>three</button> </div>
    <span>four</span>
    <button>five</button>
  `))

	// div.locator("button") → only "three"
	divButton := page.Locator("div").Locator("button")
	assertions := playwright.Expect(divButton)
	err := assertions.ToHaveText(ctx, "three")
	must.NoError(err)

	// div.locator("button, span") → "two" and "three" (span and button within div)
	divButtonOrSpan := page.Locator("div").Locator("button, span")
	count, err := divButtonOrSpan.Count(ctx)
	must.NoError(err)
	is.Equal(2, count)

	// button.or(span) → "two", "three", "four", "five"
	buttonOrSpan := page.Locator("button").Or(page.Locator("span"))
	count2, err := buttonOrSpan.Count(ctx)
	must.NoError(err)
	is.Equal(4, count2)
}

// TestLocatorPseudoFirstChild verifies Locator with :first-child pseudo.
// Ref: TestPageLocatorQuery.java#shouldThrowOnCaptureWNth  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldThrowOnCaptureWNth(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must := require.New(t)
	must.NoError(page.SetContent(ctx, `<section><div><p>A</p></div></section>`))

	err := page.Locator("*css=div >> p").Nth(1).Click(ctx)
	is.Error(err)
	is.ErrorContains(err, "n-th")
}

// Ref: TestPageLocatorQuery.java#shouldFilterByRegexWithASingleQuote  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByRegexWithASingleQuote(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>let's let's<span>hello</span></button>`))

	re := regexp.MustCompile(`(?i)let's`)
	tc, err := page.Locator("button").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re}).Locator("span").TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", *tc)

	re2 := regexp.MustCompile(`(?i)'s`)
	tc, err = page.Locator("button").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re2}).Locator("span").TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", *tc)
}

// Ref: TestPageLocatorQuery.java#shouldFilterByRegexWithSpecialSymbols  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldFilterByRegexWithSpecialSymbols(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class="test"><h5>First/"and"</h5><h2><i>Second\</i></h2></div>`))

	re := regexp.MustCompile(`(?i)first\/".*"second\\$`)
	assertions := playwright.Expect(page.Locator("div").Filter(&playwright.LocatorFilterOptions{HasTextRegex: re}))
	must.NoError(assertions.ToHaveClass(ctx, "test"))
	is.True(true)
}

// Ref: TestPageLocatorQuery.java#shouldSupportLocatorWithin  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorQuery.java
func TestPageLocatorQueryShouldSupportLocatorWithin(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
    <table>
      <tr><td>a1</td><td>a2</td><td>a3</td></tr>
      <tr><td>b1</td><td>b2</td><td>b3</td></tr>
      <tr><td>c1</td><td>c2</td><td>c3</td></tr>
    </table>
    <span>outside</span>
  `))

	must.NoError(playwright.Expect(page.GetByRole(playwright.AriaRoleCell).Within(page.GetByRole(playwright.AriaRoleRow))).ToHaveTextArray(ctx, []string{"a1", "a2", "a3", "b1", "b2", "b3", "c1", "c2", "c3"}))
	must.NoError(playwright.Expect(page.GetByRole(playwright.AriaRoleCell).Nth(1).Within(page.GetByRole(playwright.AriaRoleRow))).ToHaveTextArray(ctx, []string{"a2", "b2", "c2"}))
	must.NoError(playwright.Expect(page.GetByRole(playwright.AriaRoleCell).Last().Within(page.GetByRole(playwright.AriaRoleRow))).ToHaveTextArray(ctx, []string{"a3", "b3", "c3"}))
	must.NoError(playwright.Expect(page.GetByRole(playwright.AriaRoleCell).Nth(1).Within(page.GetByRole(playwright.AriaRoleRow).Nth(2))).ToHaveTextArray(ctx, []string{"c2"}))
	must.NoError(playwright.Expect(page.Locator("span").Within(page.GetByRole(playwright.AriaRoleRow))).ToHaveCount(ctx, 0))
	must.NoError(playwright.Expect(page.GetByRole(playwright.AriaRoleCell).Nth(1).Within(page.GetByRole(playwright.AriaRoleRow)).Nth(1)).ToHaveText(ctx, "b2"))
}
