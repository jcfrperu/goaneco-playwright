//go:build e2e

package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetByRoleButtonEx5 verifies GetByRole for buttons.
// Ref: TestLocatorGetBy.java#shouldGetButton
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
// Ref: TestLocatorGetBy.java#shouldGetLink
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
// Ref: TestLocatorGetBy.java#shouldGetCheckbox
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
// Ref: TestLocatorGetBy.java#shouldGetByLabel
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
// Ref: TestLocatorGetBy.java#shouldGetByPlaceholder
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

func localBoolPtrGB6(b bool) *bool { return &b }

func localStringPtrGB6(s string) *string { return &s }

// TestGetByRoleNamedButtonEx6 verifies GetByRole with exact name.
// Ref: TestLocatorGetBy.java#shouldGetNamedButton
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
// Ref: TestLocatorGetBy.java#shouldGetByTextExact
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
// Ref: TestLocatorGetBy.java#shouldGetByAltText
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
// Ref: TestLocatorGetBy.java#shouldGetByTitle
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

func localBoolPtrGB7(b bool) *bool { return &b }

func localStringPtrGB7(s string) *string { return &s }

// TestGetByPlaceholderEx7 verifies GetByPlaceholder finds input by placeholder.
// Ref: TestLocatorGetBy.java#shouldGetByPlaceholder
func TestGetByPlaceholderEx7(t *testing.T) {
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

// TestGetByRoleComboboxEx7 verifies GetByRole finds combobox (select).
// Ref: TestLocatorGetBy.java#shouldGetByRoleCombobox
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
// Ref: TestLocatorGetBy.java#shouldGetByLabelExact
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
// Ref: TestLocatorGetBy.java#shouldGetByTextPartial
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
// Ref: TestLocatorGetBy.java#shouldGetByRoleWithName
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

func localBoolPtrGB8(b bool) *bool { return &b }

func localStringPtrGB8(s string) *string { return &s }

// TestGetByRoleHeadingEx8 verifies GetByRole finds heading elements.
// Ref: TestLocatorGetBy.java#shouldGetByRoleHeading
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
// Ref: TestLocatorGetBy.java#shouldGetByRoleImg
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
// Ref: TestLocatorGetBy.java#shouldBeExactCaseSensitive
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
// Ref: TestLocatorGetBy.java#shouldGetButtonWithName
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

func localStringPtrGB9(s string) *string { return &s }

// TestGetByLabelEx9 verifies GetByLabel finds input by label.
// Ref: TestLocatorGetBy.java#shouldGetByLabel
func TestGetByLabelEx9(t *testing.T) {
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

// TestGetByPlaceholderEx9 verifies GetByPlaceholder finds input.
// Ref: TestLocatorGetBy.java#shouldGetByPlaceholder
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
// Ref: TestLocatorGetBy.java#shouldGetByTitle
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

// TestGetByAltTextEx9 verifies GetByAltText finds image.
// Ref: TestLocatorGetBy.java#shouldGetByAltText
func TestGetByAltTextEx9(t *testing.T) {
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

// TestGetByRoleCheckboxEx9 verifies GetByRole for checkbox.
// Ref: TestLocatorGetBy.java#shouldGetByRoleCheckbox
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
// Ref: TestLocatorGetBy.java#shouldGetByTextExact
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
// Ref: TestLocatorGetBy.java#shouldGetByTextPartial
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
// Ref: TestLocatorGetBy.java#shouldGetByLabel
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
// Ref: TestLocatorGetBy.java#shouldGetByPlaceholder
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
// Ref: TestLocatorGetBy.java#shouldFindHeading
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
// Ref: TestLocatorGetBy.java#shouldCountMultipleText
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
// Ref: TestLocatorGetBy.java#shouldFindByAlt
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
// Ref: TestLocatorGetBy.java#shouldFindByTitle
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
// Ref: TestLocatorGetBy.java#shouldFindByTestId
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
// Ref: TestPageGetByRole.java#shouldFindLink
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
// Ref: TestPageGetByRole.java#shouldFindHeading
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
// Ref: TestPageGetByRole.java#shouldFindTextbox
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
// Ref: TestPageGetByRole.java#shouldFindList
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
// Ref: TestPageGetByRole.java#shouldFilterByPressedState
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
// Ref: TestPageGetByRole.java#shouldFindCombobox
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
// Ref: TestLocatorGetByRole.java#shouldGetCheckbox
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
// Ref: TestLocatorGetByRole.java#shouldGetRadio
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
// Ref: TestLocatorGetByRole.java#shouldGetTextbox
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
// Ref: TestLocatorGetByRole.java#shouldGetImg
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
// Ref: TestLocatorGetByRole.java#shouldGetArticle
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
// Ref: TestLocatorGetByRole.java#shouldGetTab
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
// Ref: TestLocatorGetByRole.java#shouldGetNavigation
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
// Ref: TestLocatorGetByRole.java#shouldGetDialog
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
// Ref: TestLocatorGetByRole.java#shouldGetTable
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
// Ref: TestLocatorGetByRole.java#shouldGetForm
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
