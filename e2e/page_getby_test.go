//go:build e2e

// E2E tests for Page.GetByAltText, Page.GetByTitle, and Page.QuerySelectorAll.
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

func TestPageGetByAltText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<img id="logo" alt="company logo" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==">`)
	must.NoError(err, "SetContent failed")

	loc := page.GetByAltText("company logo")
	attr, err := loc.GetAttribute(ctx, "id")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "logo" {
		t.Errorf("GetByAltText: id = %v, want 'logo'", attr)
	}
}

func TestPageGetByAltTextExact(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<img id="exact" alt="cat">
		<img id="partial" alt="catfish">
	`)
	must.NoError(err, "SetContent failed")

	exact := true
	count, err := page.GetByAltText("cat", &playwright.GetByAltTextOptions{Exact: &exact}).Count(ctx)
	must.NoError(err, "Count failed")
	is.Equal(1, count)
}

func TestPageGetByTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button title="Submit form">Submit</button>`)
	must.NoError(err, "SetContent failed")

	text, err := page.GetByTitle("Submit form").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("Submit", text)
}

func TestPageGetByTitlePartial(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<span title="Save document">Save</span>
		<span title="Save as">Save as</span>
	`)
	must.NoError(err, "SetContent failed")

	count, err := page.GetByTitle("Save").Count(ctx)
	must.NoError(err, "Count failed")
	if count < 2 {
		t.Errorf("partial GetByTitle count = %d, want >= 2", count)
	}
}

func TestPageQuerySelectorAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<ul>
			<li class="item">one</li>
			<li class="item">two</li>
			<li class="item">three</li>
		</ul>
	`)
	must.NoError(err, "SetContent failed")

	handles, err := page.QuerySelectorAll(ctx, ".item")
	must.NoError(err, "QuerySelectorAll failed")
	is.Len(handles, 3)

	texts := []string{"one", "two", "three"}
	for i, el := range handles {
		text, err := el.InnerText(ctx)
		must.NoErrorf(err, "InnerText[%d] failed", i)
		if text != texts[i] {
			t.Errorf("handles[%d].InnerText = %q, want %q", i, text, texts[i])
		}
	}
}

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

func TestGetByTextPartialMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>Hello World</div>
		<div>Hello Playwright</div>
		<div>Goodbye</div>
	`))

	count, err := page.GetByText("Hello").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

func TestGetByTextExactMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	exact := true
	must.NoError(page.SetContent(ctx, `
		<div>Hello</div>
		<div>Hello World</div>
	`))

	count, err := page.GetByText("Hello", &playwright.GetByTextOptions{Exact: &exact}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByLabelExactMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	exact := true
	must.NoError(page.SetContent(ctx, `
		<label for="a">Name</label><input id="a" value="Alice">
		<label for="b">Full Name</label><input id="b" value="Bob">
	`))

	count, err := page.GetByLabel("Name", &playwright.GetByLabelOptions{Exact: &exact}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count, "exact label 'Name' should match only one element")

	val, err := page.GetByLabel("Name", &playwright.GetByLabelOptions{Exact: &exact}).InputValue(ctx)
	must.NoError(err)
	is.Equal("Alice", val)
}

func TestGetByPlaceholderPartial(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input placeholder="Enter your name">
		<input placeholder="Enter email">
	`))

	count, err := page.GetByPlaceholder("Enter").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

func TestGetByPlaceholderExact(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	exact := true
	must.NoError(page.SetContent(ctx, `
		<input placeholder="Enter your name">
		<input placeholder="Enter email">
	`))

	count, err := page.GetByPlaceholder("Enter your name", &playwright.GetByPlaceholderOptions{Exact: &exact}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByRoleWithName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Save</button>
		<button>Cancel</button>
		<button>Submit</button>
	`))

	name := "Save"
	loc := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{Name: &name})
	count, err := loc.Count(ctx)
	must.NoError(err)
	is.Equal(1, count)

	text, err := loc.InnerText(ctx)
	must.NoError(err)
	is.Equal("Save", text)
}

func TestGetByRoleWithExactName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	exact := true
	must.NoError(page.SetContent(ctx, `
		<button>Save</button>
		<button>Save All</button>
	`))

	name := "Save"
	count, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name:  &name,
		Exact: &exact,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count, "exact name 'Save' should match only one button")
}

func TestGetByRoleChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	checked := true
	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" checked>
		<input type="checkbox">
		<input type="checkbox" checked>
	`))

	count, err := page.GetByRole(playwright.AriaRoleCheckbox, &playwright.GetByRoleOptions{
		Checked: &checked,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count, "should find 2 checked checkboxes")
}

func TestGetByRoleDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	disabled := true
	must.NoError(page.SetContent(ctx, `
		<button>Active</button>
		<button disabled>Disabled</button>
		<button>Also Active</button>
	`))

	count, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Disabled: &disabled,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count, "should find 1 disabled button")

	text, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Disabled: &disabled,
	}).InnerText(ctx)
	must.NoError(err)
	is.Equal("Disabled", text)
}

func TestGetByTestIdDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button data-testid="submit-btn">Submit</button>
		<button data-testid="cancel-btn">Cancel</button>
	`))

	text, err := page.GetByTestId("submit-btn").InnerText(ctx)
	must.NoError(err)
	is.Equal("Submit", text)

	count, err := page.GetByTestId("cancel-btn").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByAltTextFindsImage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<img alt="a cat" src="cat.jpg">`))

	count, err := page.GetByAltText("a cat").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByAltTextPartialMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img alt="a beautiful sunset" src="sunset.jpg">
		<img alt="a cat" src="cat.jpg">
	`))

	count, err := page.GetByAltText("sunset").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByTitleFindsElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button title="Submit form">Submit</button>`))

	count, err := page.GetByTitle("Submit form").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByTitlePartialMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<abbr title="World Health Organization">WHO</abbr>
		<abbr title="United Nations">UN</abbr>
	`))

	count, err := page.GetByTitle("World Health").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByTestIdFindsElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div data-testid="submit-button">Click me</div>
		<div data-testid="cancel-button">Cancel</div>
	`))

	count, err := page.GetByTestId("submit-button").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)

	text, err := page.GetByTestId("submit-button").InnerText(ctx)
	must.NoError(err)
	is.Equal("Click me", text)
}

func TestGetByPlaceholderFindsTwoInputs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input placeholder="enter name">
		<input placeholder="enter email">
		<input placeholder="different">
	`))

	count, err := page.GetByPlaceholder("enter").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

func TestGetByRoleWithNameOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Save</button>
		<button>Cancel</button>
		<button>Delete</button>
	`))

	name := "Save"
	count, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name: &name,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByRoleWithExactNameOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Save</button>
		<button>Save all</button>
	`))

	name := "Save"
	exact := true
	count, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name:  &name,
		Exact: &exact,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByLabelExactMatchEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="fn">First name</label>
		<input id="fn" type="text">
		<label for="fn2">First name (required)</label>
		<input id="fn2" type="text">
	`))

	exact := true
	count, err := page.GetByLabel("First name", &playwright.GetByLabelOptions{Exact: &exact}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByPlaceholderExactMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input placeholder="Search here">
		<input placeholder="Search here or press Enter">
	`))

	exact := true
	count, err := page.GetByPlaceholder("Search here", &playwright.GetByPlaceholderOptions{Exact: &exact}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestGetByRoleHeadingWithLevel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h1>Top heading</h1>
		<h2>Section A</h2>
		<h2>Section B</h2>
		<h3>Sub section</h3>
	`))

	level := 2
	count, err := page.GetByRole(playwright.AriaRoleHeading, &playwright.GetByRoleOptions{
		Level: &level,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}
