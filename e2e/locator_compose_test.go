//go:build e2e

// E2E tests for Locator composition operators: And, Or, Filter.
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorAnd verifies that And() narrows to elements matching both locators.
func TestLocatorAnd(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button class="primary">Save</button>
		<button class="secondary">Cancel</button>
		<button class="primary highlighted">Submit</button>
	`)
	must.NoError(err, "SetContent failed")

	// All primary buttons
	primary := page.Locator("button.primary")
	count, err := primary.Count(ctx)
	if err != nil || count != 2 {
		t.Fatalf("expected 2 primary buttons, got %d (err: %v)", count, err)
	}

	// Primary AND highlighted → only the one with both classes
	highlighted := page.Locator("button.highlighted")
	both := primary.And(highlighted)

	count, err = both.Count(ctx)
	must.NoError(err, "And().Count() failed")
	is.Equal(1, count)

	text, err := both.InnerText(ctx)
	must.NoError(err, "And().InnerText() failed")
	is.Equal("Submit", text)
}

// TestLocatorOr verifies that Or() matches elements satisfying either locator.
func TestLocatorOr(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button>Click me</button>
		<a href="#">Link one</a>
		<a href="#">Link two</a>
		<span>Just text</span>
	`)
	must.NoError(err, "SetContent failed")

	buttons := page.Locator("button")
	links := page.Locator("a")
	buttonsOrLinks := buttons.Or(links)

	count, err := buttonsOrLinks.Count(ctx)
	must.NoError(err, "Or().Count() failed")
	// 1 button + 2 links = 3
	is.Equal(3, count)
}

// TestLocatorFilterByText verifies that Filter(HasText) narrows by visible text.
func TestLocatorFilterByText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<ul>
			<li>Apple</li>
			<li>Banana</li>
			<li>Apricot</li>
			<li>Cherry</li>
		</ul>
	`)
	must.NoError(err, "SetContent failed")

	items := page.Locator("li")
	total, _ := items.Count(ctx)
	if total != 4 {
		t.Fatalf("expected 4 list items, got %d", total)
	}

	needle := "ap"
	filtered := items.Filter(&playwright.LocatorFilterOptions{HasText: &needle})
	count, err := filtered.Count(ctx)
	must.NoError(err, "Filter(HasText).Count() failed")
	is.Equal(2, count)
}

// TestLocatorFilterByHasNotText verifies that Filter(HasNotText) excludes items by text.
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

// TestLocatorFilterByHas verifies that Filter(Has) selects elements containing a child locator.
func TestLocatorFilterByHas(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<ul>
			<li><span class="badge">new</span> Widget A</li>
			<li>Widget B</li>
			<li><span class="badge">sale</span> Widget C</li>
		</ul>
	`)
	must.NoError(err, "SetContent failed")

	items := page.Locator("li")
	badge := page.Locator(".badge")
	withBadge := items.Filter(&playwright.LocatorFilterOptions{Has: badge})

	count, err := withBadge.Count(ctx)
	must.NoError(err, "Filter(Has).Count() failed")
	is.Equal(2, count)
}

// TestLocatorFilterByHasNot verifies that Filter(HasNot) excludes elements containing a child locator.
func TestLocatorFilterByHasNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<ul>
			<li><span class="badge">new</span> Widget A</li>
			<li>Widget B</li>
			<li><span class="badge">sale</span> Widget C</li>
		</ul>
	`)
	must.NoError(err, "SetContent failed")

	items := page.Locator("li")
	badge := page.Locator(".badge")
	withoutBadge := items.Filter(&playwright.LocatorFilterOptions{HasNot: badge})

	count, err := withoutBadge.Count(ctx)
	must.NoError(err, "Filter(HasNot).Count() failed")
	is.Equal(1, count)

	text, err := withoutBadge.InnerText(ctx)
	must.NoError(err, "Filter(HasNot).InnerText() failed")
	is.Equal("Widget B", text)
}

// TestLocatorAndReturnsIntersection verifies Locator.And returns elements matching both selectors.
// Ref: TestLocatorCompose.java#shouldAndReturnIntersection
func TestLocatorAndReturnsIntersectionExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="blue">blue</div>
		<div class="red">red</div>
		<div class="blue red">both</div>
	`))

	count, err := page.Locator(".blue").And(page.Locator(".red")).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)

	text, err := page.Locator(".blue").And(page.Locator(".red")).InnerText(ctx)
	must.NoError(err)
	is.Equal("both", text)
}

// TestLocatorOrReturnsUnion verifies Locator.Or returns elements matching either selector.
// Ref: TestLocatorCompose.java#shouldOrReturnUnion
func TestLocatorOrReturnsUnionExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn1">Button 1</button>
		<input id="inp1" type="text">
	`))

	count, err := page.Locator("button").Or(page.Locator("input")).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterWithHasTextChained verifies Filter by text can be chained with locators.
// Ref: TestLocatorCompose.java#shouldFilterByTextChained
func TestLocatorFilterWithHasTextChained(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="item">apple</li>
			<li class="item">banana</li>
			<li class="item">cherry</li>
		</ul>
	`))

	hasText := "banana"
	item, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{
		HasText: &hasText,
	}).InnerText(ctx)
	must.NoError(err)
	is.Equal("banana", item)
}

// TestLocatorFilterWithHasLocatorNested verifies Filter with Has finds nested elements.
// Ref: TestLocatorCompose.java#shouldFilterByHasLocator
func TestLocatorFilterWithHasLocatorNested(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">
			<span class="badge">featured</span>
			<p>First item</p>
		</div>
		<div class="item">
			<p>Second item</p>
		</div>
	`))

	count, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{
		Has: page.Locator(".badge"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}
