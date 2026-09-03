//go:build e2e

// Ref: TestSelectorsGetBy.java
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetByTestIdEscapesId verifies GetByTestId works with attribute values containing double quotes.
// Ref: TestSelectorsGetBy.java#getByTestIdShouldEscapeId
func TestGetByTestIdEscapesId(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><div data-testid='He"llo'>Hello world</div></div>`))

	loc := page.GetByTestId(`He"llo`)
	text, err := loc.TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Hello world", *text)
}

// TestGetByTestIdWithCustomAttribute verifies GetByTestId uses the attribute set by SetTestIdAttribute.
// Ref: TestSelectorsGetBy.java#getByTestIdWithCustomTestIdShouldWork
func TestGetByTestIdWithCustomAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx := newContext(t)
	must.NoError(bCtx.SetTestIdAttribute(ctx, "data-my-custom-testid"))
	defer func() {
		_ = bCtx.SetTestIdAttribute(ctx, "data-testid")
	}()

	pg, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(pg.SetContent(ctx, `<div><div data-my-custom-testid='Hello'>Hello world</div></div>`))

	loc := pg.GetByTestId("Hello")
	text, err := loc.TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Hello world", *text)
}

// TestGetByTestIdWithMultipleAttributes verifies GetByTestId matches any of a comma-separated list of test-id attributes.
// Ref: TestSelectorsGetBy.java#getByTestIdWithCommaSeparatedTestIdAttributesShouldMatchAny
func TestGetByTestIdWithMultipleAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx := newContext(t)
	must.NoError(bCtx.SetTestIdAttribute(ctx, "data-pw,data-ti"))
	defer func() {
		_ = bCtx.SetTestIdAttribute(ctx, "data-testid")
	}()

	pg, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(pg.SetContent(ctx, `<section><div data-pw='Hello'>first</div><div data-ti='Hello'>second</div><div data-testid='Hello'>third</div></section>`))

	count, err := pg.GetByTestId("Hello").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByLabelNestedElements verifies GetByLabel works with labels containing nested elements.
// Ref: TestSelectorsGetBy.java#getByLabelShouldWorkWithNestedElements
func TestGetByLabelNestedElements(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	page := newPage(t)
	is := assert.New(t)
	must := require.New(t)

	must.NoError(page.SetContent(ctx, `<label for=target>Last <span>Name</span></label><input id=target type=text>`))

	trueVal := true

	checkAttr := func(loc *playwright.Locator) {
		t.Helper()
		attr, err := loc.GetAttribute(ctx, "id")
		must.NoError(err)
		must.NotNil(attr)
		is.Equal("target", *attr)
	}
	checkCount := func(loc *playwright.Locator, want int) {
		t.Helper()
		count, err := loc.Count(ctx)
		must.NoError(err)
		is.Equal(want, count)
	}

	checkAttr(page.GetByLabel("last name"))
	checkAttr(page.GetByLabel("st na"))
	checkAttr(page.GetByLabel("Name"))
	checkAttr(page.GetByLabel("Last Name", &playwright.GetByLabelOptions{Exact: &trueVal}))

	checkCount(page.GetByLabel("Last", &playwright.GetByLabelOptions{Exact: &trueVal}), 0)
	checkCount(page.GetByLabel("last name", &playwright.GetByLabelOptions{Exact: &trueVal}), 0)
	checkCount(page.GetByLabel("Name", &playwright.GetByLabelOptions{Exact: &trueVal}), 0)
	checkCount(page.GetByLabel("what?"), 0)
}

// TestGetByEscaping verifies GetBy* methods correctly handle special characters in values.
// Ref: TestSelectorsGetBy.java#getByEscaping
func TestGetByEscaping(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must := require.New(t)

	// Part 1: value with embedded newline and double-quote
	must.NoError(page.SetContent(ctx, "<label id=label for=control>Hello my\nwo\"rld</label><input id=control />"))
	_, err := page.EvalOnSelector(ctx, "input", `e => {
		e.placeholder = "hello my\nwo\"rld";
		e.title = "hello my\nwo\"rld";
		e.alt = "hello my\nwo\"rld";
	}`)
	must.NoError(err)

	checkID := func(loc *playwright.Locator, wantID string) {
		t.Helper()
		attr, err := loc.GetAttribute(ctx, "id")
		must.NoError(err)
		must.NotNil(attr)
		is.Equal(wantID, *attr)
	}

	val1 := "hello my\nwo\"rld"
	checkID(page.GetByText(val1), "label")
	checkID(page.GetByLabel(val1), "control")
	checkID(page.GetByPlaceholder(val1), "control")
	checkID(page.GetByAltText(val1), "control")
	checkID(page.GetByTitle(val1), "control")

	// Part 2: value with embedded newline but no double-quote
	must.NoError(page.SetContent(ctx, "<label id=label for=control>Hello my\nworld</label><input id=control />"))
	_, err = page.EvalOnSelector(ctx, "input", `e => {
		e.placeholder = "hello my\nworld";
		e.title = "hello my\nworld";
		e.alt = "hello my\nworld";
	}`)
	must.NoError(err)

	val2 := "hello my\nworld"
	checkID(page.GetByText(val2), "label")
	checkID(page.GetByLabel(val2), "control")
	checkID(page.GetByPlaceholder(val2), "control")
	checkID(page.GetByAltText(val2), "control")
	checkID(page.GetByTitle(val2), "control")

	// Part 3: HTML entities and special characters
	must.NoError(page.SetContent(ctx, `<label for=target>foo &gt;&gt; bar</label><input id=target>`))
	_, err = page.EvalOnSelector(ctx, "input", `e => {
		e.placeholder = "foo >> bar";
		e.title = "foo >> bar";
		e.alt = "foo >> bar";
	}`)
	must.NoError(err)

	val3 := "foo >> bar"

	text, err := page.GetByText(val3).TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal(val3, *text)

	checkID(page.GetByLabel(val3), "target")
	checkID(page.GetByPlaceholder(val3), "target")
	checkID(page.GetByAltText(val3), "target")
	checkID(page.GetByTitle(val3), "target")
}

// TestGetByRoleEscaping verifies GetByRole correctly handles names with special whitespace.
// Ref: TestSelectorsGetBy.java#getByRoleEscaping
func TestGetByRoleEscaping(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	page := newPage(t)
	must := require.New(t)
	is := assert.New(t)

	must.NoError(page.SetContent(ctx, `
<a href="https://playwright.dev">issues 123</a>
<a href="https://playwright.dev">he llo 56</a>
<button>Click me</button>
`))

	trueVal := true

	toSlice := func(v any) []any {
		t.Helper()
		arr, ok := v.([]any)
		if !ok {
			t.Fatalf("EvaluateAll result is not []any: %T", v)
		}
		return arr
	}

	// All buttons
	res, err := page.GetByRole(playwright.AriaRoleButton).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)
	arr := toSlice(res)
	is.Len(arr, 1)
	is.Equal("<button>Click me</button>", arr[0])

	// All links
	res, err = page.GetByRole(playwright.AriaRoleLink).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)
	arr = toSlice(res)
	is.Len(arr, 2)

	// Link by exact name
	issues123Name := "issues 123"
	res, err = page.GetByRole(playwright.AriaRoleLink, &playwright.GetByRoleOptions{Name: &issues123Name}).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)
	arr = toSlice(res)
	is.Len(arr, 1)

	// Partial name match
	suesName := "sues"
	res, err = page.GetByRole(playwright.AriaRoleLink, &playwright.GetByRoleOptions{Name: &suesName}).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)
	arr = toSlice(res)
	is.Len(arr, 1)

	// Whitespace normalization
	wspaceName := "  he    \n  llo "
	res, err = page.GetByRole(playwright.AriaRoleLink, &playwright.GetByRoleOptions{Name: &wspaceName}).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)
	arr = toSlice(res)
	is.Len(arr, 1)

	// No match: button with link name
	res, err = page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{Name: &issues123Name}).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)
	arr = toSlice(res)
	is.Len(arr, 0)

	// Exact partial — no match
	res, err = page.GetByRole(playwright.AriaRoleLink, &playwright.GetByRoleOptions{Name: &suesName, Exact: &trueVal}).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)
	arr = toSlice(res)
	is.Len(arr, 0)
}

// TestLocatorGetByRoleScoped verifies GetByRole can be scoped within a locator subtree.
// Ref: TestSelectorsGetBy.java#locatorGetByRole
func TestLocatorGetByRoleScoped(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><button>Click me</button></div>`))

	roleSelector := page.GetByRole(playwright.AriaRoleButton).Selector()
	scoped := page.Locator("div").Locator(roleSelector)

	res, err := scoped.EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
	must.NoError(err)

	arr, ok := res.([]any)
	is.True(ok, "EvaluateAll result is not []any: %T", res)
	is.Len(arr, 1)
	is.Equal("<button>Click me</button>", arr[0])
}
