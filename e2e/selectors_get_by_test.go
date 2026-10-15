//go:build e2e

// E2E tests for: TestSelectorsGetBy.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"regexp"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetByRoleButtonEx5 verifies GetByRole for buttons.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleButtonEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Submit</button>
		<button>Cancel</button>
	`))

	count, err := page.GetByRole("button").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByRoleLinkEx5 verifies GetByRole for links.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleLinkEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a href="/home">Home</a>
		<a href="/about">About</a>
	`))

	count, err := page.GetByRole("link").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByRoleCheckboxEx5 verifies GetByRole for checkboxes.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleCheckboxEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="chk1">
		<input type="checkbox" id="chk2" checked>
	`))

	count, err := page.GetByRole("checkbox").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByLabelTextEx5 verifies GetByLabel finds associated input.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByLabelTextEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="username">Username</label>
		<input id="username" type="text">
	`))

	must.NoError(page.GetByLabel("Username").Fill(ctx, "john_doe"))

	val, err := page.GetByLabel("Username").InputValue(ctx)
	must.NoError(err)
	is.Equal("john_doe", val)
}

// TestGetByPlaceholderEx5 verifies GetByPlaceholder finds input by placeholder.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByPlaceholderEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" placeholder="Enter email">`))

	must.NoError(page.GetByPlaceholder("Enter email").Fill(ctx, "test@example.com"))

	val, err := page.GetByPlaceholder("Enter email").InputValue(ctx)
	must.NoError(err)
	is.Equal("test@example.com", val)
}

// TestGetByRoleNamedButtonEx6 verifies GetByRole with exact name.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleNamedButtonEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Submit</button>
		<button>Cancel</button>
	`))

	count, err := page.GetByRole("button", &playwright.GetByRoleOptions{
		Name:  localStringPtrGB6("Submit"),
		Exact: localBoolPtrGB6(true),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextExactEx6 verifies GetByText with exact match.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTextExactEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Hello</p>
		<p>Hello World</p>
	`))

	count, err := page.GetByText("Hello", &playwright.GetByTextOptions{
		Exact: localBoolPtrGB6(true),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByAltTextExactEx6 verifies GetByAltText locator.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByAltTextExactEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img alt="Profile picture" src="profile.jpg">
		<img alt="Logo" src="logo.png">
	`))

	count, err := page.GetByAltText("Profile picture").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTitleEx6 verifies GetByTitle locator.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTitleEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span title="Tooltip text">Hover me</span>
	`))

	count, err := page.GetByTitle("Tooltip text").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleComboboxEx7 verifies GetByRole finds combobox (select).
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleComboboxEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel"><option>Alpha</option></select>`))

	count, err := page.GetByRole("combobox").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByLabelExactEx7 verifies GetByLabel with exact match.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByLabelExactEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="inp">Email</label>
		<input id="inp" type="email">
	`))

	count, err := page.GetByLabel("Email").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextPartialEx7 verifies GetByText with partial match.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTextPartialEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Welcome to the site</p>
		<p>Other content</p>
	`))

	count, err := page.GetByText("Welcome", &playwright.GetByTextOptions{Exact: localBoolPtrGB7(false)}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleWithNameEx7 verifies GetByRole with Name option.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleWithNameEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Cancel</button>
		<button>Confirm</button>
	`))

	name := localStringPtrGB7("Confirm")
	count, err := page.GetByRole("button", &playwright.GetByRoleOptions{Name: name}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleHeadingEx8 verifies GetByRole finds heading elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleHeadingEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h1>Main Heading</h1>
		<h2>Sub Heading 1</h2>
		<h2>Sub Heading 2</h2>
	`))

	count, err := page.GetByRole("heading").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestGetByRoleImgEx8 verifies GetByRole finds img elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleImgEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7" alt="Photo 1">
		<img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7" alt="Photo 2">
	`))

	count, err := page.GetByRole("img").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByTextExactCaseSensitiveEx8 verifies GetByText exact match is case sensitive.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTextExactCaseSensitiveEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span>HELLO</span>
		<span>hello</span>
		<span>Hello</span>
	`))

	count, err := page.GetByText("Hello", &playwright.GetByTextOptions{Exact: localBoolPtrGB8(true)}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleButtonWithNameEx8 verifies GetByRole button with specific name.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleButtonWithNameEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Save</button>
		<button>Delete</button>
		<button>Save Draft</button>
	`))

	name := localStringPtrGB8("Delete")
	count, err := page.GetByRole("button", &playwright.GetByRoleOptions{Name: name}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByPlaceholderEx9 verifies GetByPlaceholder finds input.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByPlaceholderEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" placeholder="Enter your name">
	`))

	must.NoError(page.GetByPlaceholder("Enter your name").Fill(ctx, "Alice"))

	val, err := page.GetByPlaceholder("Enter your name").InputValue(ctx)
	must.NoError(err)
	is.Equal("Alice", val)
}

// TestGetByTitleEx9 verifies GetByTitle finds element.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTitleEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<abbr title="HyperText Markup Language">HTML</abbr>`))

	text, err := page.GetByTitle("HyperText Markup Language").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("HTML", *text)
}

// TestGetByRoleCheckboxEx9 verifies GetByRole for checkbox.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleCheckboxEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="agree" aria-label="I agree">
		<label for="agree">I agree</label>
	`))

	name := localStringPtrGB9("I agree")
	count, err := page.GetByRole("checkbox", &playwright.GetByRoleOptions{Name: name}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextExactMatch verifies GetByText with exact text.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTextExactMatchLoc(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>hello world</div>
		<div>hello</div>
	`))

	count, err := page.GetByText("hello", &playwright.GetByTextOptions{Exact: boolPtr(true)}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextPartialMatch verifies GetByText with partial text match.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTextPartialMatchLoc(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span>apple pie</span>
		<span>apple juice</span>
		<span>orange</span>
	`))

	count, err := page.GetByText("apple").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByLabelFindsInput verifies GetByLabel finds associated input.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByLabelFindsInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="email">Email address</label>
		<input id="email" type="email">
	`))

	must.NoError(page.GetByLabel("Email address").Fill(ctx, "test@example.com"))

	val, err := page.GetByLabel("Email address").InputValue(ctx)
	must.NoError(err)
	is.Equal("test@example.com", val)
}

// TestGetByPlaceholderFindsInput verifies GetByPlaceholder finds input by placeholder.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByPlaceholderFindsInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input placeholder="Enter your name">`))

	must.NoError(page.GetByPlaceholder("Enter your name").Fill(ctx, "Alice"))

	val, err := page.GetByPlaceholder("Enter your name").InputValue(ctx)
	must.NoError(err)
	is.Equal("Alice", val)
}

// TestGetByTextFindsHeadingEx4 verifies GetByText locates a heading.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTextFindsHeadingEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1>Welcome to the Site</h1>`))

	count, err := page.GetByText("Welcome to the Site").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTextCountMultipleEx4 verifies GetByText counts multiple matching elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTextCountMultipleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Hello</p>
		<span>Hello</span>
		<div>Hello</div>
	`))

	count, err := page.GetByText("Hello").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestGetByAltTextFindsImageEx4 verifies GetByAltText locates an image by alt attribute.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByAltTextFindsImageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<img src="logo.png" alt="Company Logo">`))

	count, err := page.GetByAltText("Company Logo").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTitleFindsElementEx4 verifies GetByTitle locates element by title attribute.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTitleFindsElementEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<abbr title="HyperText Markup Language">HTML</abbr>`))

	count, err := page.GetByTitle("HyperText Markup Language").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByTestIdFindsElementEx4 verifies GetByTestId locates element by data-testid.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByTestIdFindsElementEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button data-testid="submit-btn">Submit</button>
	`))

	count, err := page.GetByTestId("submit-btn").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleLink verifies GetByRole finds anchor elements with link role.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleLink(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a href="#">Go here</a>`))

	count, err := page.GetByRole(playwright.AriaRoleLink).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleHeading verifies GetByRole finds heading elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleHeading(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h1>Main Title</h1>
		<h2>Subtitle</h2>
	`))

	count, err := page.GetByRole(playwright.AriaRoleHeading).Count(ctx)
	must.NoError(err)
	is.GreaterOrEqual(count, 2)
}

// TestGetByRoleTextbox verifies GetByRole finds text input as textbox.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleTextbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" placeholder="type here">`))

	count, err := page.GetByRole(playwright.AriaRoleTextbox).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleList verifies GetByRole finds list elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleList(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul aria-label="items">
			<li>item 1</li>
			<li>item 2</li>
		</ul>
	`))

	count, err := page.GetByRole(playwright.AriaRoleList).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleButtonWithPressed verifies GetByRole can filter by pressed state.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleButtonWithPressed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button aria-pressed="true">On</button>
		<button aria-pressed="false">Off</button>
	`))

	pressedTrue := true
	count, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Pressed: &pressedTrue,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleCombobox verifies GetByRole finds select/combobox elements.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleCombobox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select>
			<option>A</option>
			<option>B</option>
		</select>
	`))

	count, err := page.GetByRole(playwright.AriaRoleCombobox).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleCheckbox verifies GetByRole for checkbox.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleCheckbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" aria-label="Accept">`))

	count, err := page.GetByRole(playwright.AriaRoleCheckbox).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleRadio verifies GetByRole for radio button.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleRadio(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="radio" name="choice" value="a">
		<input type="radio" name="choice" value="b">
	`))

	count, err := page.GetByRole(playwright.AriaRoleRadio).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByRoleTextbox verifies GetByRole for textbox.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleTextboxExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" aria-label="First name">
		<input type="text" aria-label="Last name">
	`))

	count, err := page.GetByRole(playwright.AriaRoleTextbox).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByRoleImg verifies GetByRole for images.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleImg(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img src="" alt="Product image">
		<img src="" alt="Banner">
	`))

	count, err := page.GetByRole(playwright.AriaRoleImg).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByRoleArticle verifies GetByRole for article landmark.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleArticle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<article>first article</article>
		<article>second article</article>
	`))

	count, err := page.GetByRole(playwright.AriaRoleArticle).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByRoleTab verifies GetByRole for tab.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleTab(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div role="tablist">
			<button role="tab" aria-selected="true">Tab 1</button>
			<button role="tab">Tab 2</button>
			<button role="tab">Tab 3</button>
		</div>
	`))

	count, err := page.GetByRole(playwright.AriaRoleTab).Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestGetByRoleNavigation verifies GetByRole for navigation landmark.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav aria-label="Primary">primary nav</nav>
		<nav aria-label="Secondary">secondary nav</nav>
	`))

	count, err := page.GetByRole(playwright.AriaRoleNavigation).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestGetByRoleDialog verifies GetByRole for dialog.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleDialog(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div role="dialog" aria-label="Confirm">
			<p>Are you sure?</p>
		</div>
	`))

	count, err := page.GetByRole(playwright.AriaRoleDialog).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleTable verifies GetByRole for table.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleTable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr><th>Name</th><th>Age</th></tr>
			<tr><td>Alice</td><td>30</td></tr>
		</table>
	`))

	count, err := page.GetByRole(playwright.AriaRoleTable).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestGetByRoleForm verifies GetByRole for form.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestGetByRoleForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form aria-label="Login form">
			<input type="text"><input type="password">
		</form>
	`))

	count, err := page.GetByRole(playwright.AriaRoleForm).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// TestGetByTextPartialMatchEx3 verifies GetByText partial match works.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// TestGetByTextExactMatchEx5 verifies GetByText with exact match finds only exact.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// TestGetByTextPartialMatchEx5 verifies GetByText partial match (non-exact).
// Extra: no equivalent method in playwright-java; related class: TestSelectorsGetBy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
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

// TestSelectorsGetByGetByPlaceholderShouldWork verifies GetByPlaceholder finds input by placeholder.
// Ref: TestSelectorsGetBy.java#getByPlaceholderShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByPlaceholderShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" placeholder="Enter your name">`))

	count, err := page.GetByPlaceholder("Enter your name").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestSelectorsGetByGetByLabelShouldWork verifies GetByLabel finds input by label.
// Ref: TestSelectorsGetBy.java#getByLabelShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByLabelShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="email">Email address</label>
		<input id="email" type="email">
	`))

	must.NoError(page.GetByLabel("Email address").Fill(ctx, "test@example.com"))

	val, err := page.Locator("#email").InputValue(ctx)
	must.NoError(err)
	is.Equal("test@example.com", val)
}

// TestSelectorsGetByGetByAltTextShouldWork verifies GetByAltText finds image.
// Ref: TestSelectorsGetBy.java#getByAltTextShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByAltTextShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			alt="Company logo" id="logo">
	`))

	attr, err := page.GetByAltText("Company logo").GetAttribute(ctx, "id")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("logo", *attr)
}

// TestSelectorsGetByGetByTextShouldWork verifies GetByText selects elements by text content.
// Ref: TestSelectorsGetBy.java#getByTextShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByTextShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<p>Hello World</p>
		<p>Hello Playwright</p>
		<span>Different content</span>
	`)
	must.NoError(err, "SetContent failed")

	// Partial match (default): both "Hello" paragraphs
	helloLocator := page.GetByText("Hello")
	count, err := helloLocator.Count(ctx)
	must.NoError(err, "GetByText(Hello).Count() failed")
	is.Equal(2, count)

	// Exact match for "Hello World"
	exact := true
	exactLocator := page.GetByText("Hello World", &playwright.GetByTextOptions{Exact: &exact})
	exactCount, err := exactLocator.Count(ctx)
	must.NoError(err, "GetByText(Hello World, exact).Count() failed")
	is.Equal(1, exactCount)
}

// TestSelectorsGetByGetByTestIdShouldWork verifies GetByTestId selects elements by data-testid attribute.
// Ref: TestSelectorsGetBy.java#getByTestIdShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByTestIdShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button data-testid="submit-btn">Submit</button>
		<input data-testid="search-input" type="text" value="playwright" />
		<div data-testid="result-container"><p>Result here</p></div>
	`)
	must.NoError(err, "SetContent failed")

	// Locate by data-testid
	btnLocator := page.GetByTestId("submit-btn")
	visible, err := btnLocator.IsVisible(ctx)
	must.NoError(err, "GetByTestId(submit-btn).IsVisible() failed")
	if !visible {
		t.Error("expected submit-btn to be visible")
	}

	text, err := btnLocator.InnerText(ctx)
	must.NoError(err, "GetByTestId(submit-btn).InnerText() failed")
	is.Equal("Submit", text)

	inputLocator := page.GetByTestId("search-input")
	val, err := inputLocator.InputValue(ctx)
	must.NoError(err, "GetByTestId(search-input).InputValue() failed")
	is.Equal("playwright", val)
}

// Ref: TestSelectorsGetBy.java#getByTitleShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByTitleShouldWork(t *testing.T) {
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

// TestSelectorsGetByGetByTestIdShouldEscapeId verifies GetByTestId works with attribute values containing double quotes.
// Ref: TestSelectorsGetBy.java#getByTestIdShouldEscapeId  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByTestIdShouldEscapeId(t *testing.T) {
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

// TestSelectorsGetByGetByTestIdWithCustomTestIdShouldWork verifies GetByTestId uses the attribute set by SetTestIdAttribute.
// Ref: TestSelectorsGetBy.java#getByTestIdWithCustomTestIdShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByTestIdWithCustomTestIdShouldWork(t *testing.T) {
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

// TestSelectorsGetByGetByTestIdWithCommaSeparatedTestIdAttributesShouldMatchAny verifies GetByTestId matches any of a comma-separated list of test-id attributes.
// Ref: TestSelectorsGetBy.java#getByTestIdWithCommaSeparatedTestIdAttributesShouldMatchAny  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByTestIdWithCommaSeparatedTestIdAttributesShouldMatchAny(t *testing.T) {
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

// TestSelectorsGetByGetByLabelShouldWorkWithNestedElements verifies GetByLabel works with labels containing nested elements.
// Ref: TestSelectorsGetBy.java#getByLabelShouldWorkWithNestedElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByLabelShouldWorkWithNestedElements(t *testing.T) {
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

// TestSelectorsGetByGetByEscaping verifies GetBy* methods correctly handle special characters in values.
// Ref: TestSelectorsGetBy.java#getByEscaping  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByEscaping(t *testing.T) {
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

// TestSelectorsGetByGetByRoleEscaping verifies GetByRole correctly handles names with special whitespace.
// Ref: TestSelectorsGetBy.java#getByRoleEscaping  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByRoleEscaping(t *testing.T) {
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

// TestSelectorsGetByLocatorGetByRole verifies GetByRole can be scoped within a locator subtree.
// Ref: TestSelectorsGetBy.java#locatorGetByRole  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByLocatorGetByRole(t *testing.T) {
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

// Ref: TestSelectorsGetBy.java#shouldUseDataTestidInStrictErrors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByShouldUseDataTestidInStrictErrors(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Context().SetTestIdAttribute(ctx, "data-custom-id"))
	must.NoError(page.SetContent(ctx, ""+
		"      <div>\n"+
		"      <div></div>\n"+
		"        <div>\n"+
		"          <div></div>\n"+
		"          <div></div>\n"+
		"        </div>\n"+
		"      </div>\n"+
		"      <div>\n"+
		"        <div class='foo bar:0' data-custom-id='One'>\n"+
		"        </div>\n"+
		"        <div class='foo bar:1' data-custom-id='Two'>\n"+
		"        </div>\n"+
		"      </div>"))
	err := page.Locator(".foo").Hover(ctx)
	must.Error(err)
	msg := a3NormalizeQuotes(err.Error())
	assert.Contains(t, msg, "strict mode violation", msg)
	assert.Contains(t, msg, "<div class=\"foo bar:0", msg)
	assert.Contains(t, msg, "<div class=\"foo bar:1", msg)
	assert.Contains(t, msg, "aka getByTestId(\"One\")", msg)
	assert.Contains(t, msg, "aka getByTestId(\"Two\")", msg)
}

// Ref: TestSelectorsGetBy.java#getByTestIdShouldWorkForRegex  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsGetBy.java
func TestSelectorsGetByGetByTestIdShouldWorkForRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div><div data-testid='Hello'>Hello world</div></div>"))
	must.NoError(playwright.Expect(page.GetByTestId(regexp.MustCompile("He[l]*o"))).ToHaveText(ctx, "Hello world"))
	must.NoError(playwright.Expect(page.MainFrame().GetByTestId("Hello")).ToHaveText(ctx, "Hello world"))
	must.NoError(playwright.Expect(page.Locator("div").GetByTestId("Hello")).ToHaveText(ctx, "Hello world"))
}
