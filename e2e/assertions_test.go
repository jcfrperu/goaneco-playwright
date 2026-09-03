//go:build e2e

// E2E tests for Web-First Assertions (E2E-LOC-04).
// Covers ToBeVisible, ToBeEnabled, ToBeDisabled, ToBeEditable,
// ToBeChecked, ToHaveText, ToHaveValue, ToHaveCount, and Not() negation.
package e2e

import (
	"regexp"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExpectToBeVisible verifies basic visibility assertion.
func TestExpectToBeVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="visible">Hello</div>
		<div id="hidden" style="display:none">Hidden</div>
	`)
	must.NoError(err, "SetContent failed")

	visibleLoc := page.Locator("#visible")
	must.NoError(playwright.Expect(visibleLoc).ToBeVisible(ctx), "Expect(#visible).ToBeVisible() failed")

	hiddenLoc := page.Locator("#hidden")
	must.NoError(playwright.Expect(hiddenLoc).Not().ToBeVisible(ctx), "Expect(#hidden).Not().ToBeVisible() failed")
}

// TestExpectToBeEnabled verifies enabled state assertion.
func TestExpectToBeEnabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="enabled-btn">Enabled</button>
		<button id="disabled-btn" disabled>Disabled</button>
		<input id="enabled-input" type="text" />
		<input id="disabled-input" type="text" disabled />
	`)
	must.NoError(err, "SetContent failed")

	enabledBtn := page.Locator("#enabled-btn")
	must.NoError(playwright.Expect(enabledBtn).ToBeEnabled(ctx), "Expect(#enabled-btn).ToBeEnabled() failed")

	disabledBtn := page.Locator("#disabled-btn")
	must.NoError(playwright.Expect(disabledBtn).Not().ToBeEnabled(ctx), "Expect(#disabled-btn).Not().ToBeEnabled() failed")

	enabledInput := page.Locator("#enabled-input")
	must.NoError(playwright.Expect(enabledInput).ToBeEnabled(ctx), "Expect(#enabled-input).ToBeEnabled() failed")

	disabledInput := page.Locator("#disabled-input")
	must.NoError(playwright.Expect(disabledInput).ToBeDisabled(ctx), "Expect(#disabled-input).ToBeDisabled() failed")
}

// TestExpectToBeDisabled verifies disabled state assertion.
func TestExpectToBeDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="disabled-btn" disabled>Disabled</button>
		<button id="enabled-btn">Enabled</button>
	`)
	must.NoError(err, "SetContent failed")

	disabledLoc := page.Locator("#disabled-btn")
	must.NoError(playwright.Expect(disabledLoc).ToBeDisabled(ctx), "Expect(#disabled-btn).ToBeDisabled() failed")

	enabledLoc := page.Locator("#enabled-btn")
	must.NoError(playwright.Expect(enabledLoc).Not().ToBeDisabled(ctx), "Expect(#enabled-btn).Not().ToBeDisabled() failed")
}

// TestExpectToHaveText verifies text assertion (substring match with whitespace normalization).
func TestExpectToHaveText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<p id="msg">Hello   World</p>
		<h1 id="title">Page Title</h1>
	`)
	must.NoError(err, "SetContent failed")

	msgLoc := page.Locator("#msg")
	must.NoError(playwright.Expect(msgLoc).ToHaveText(ctx, "Hello World"), "Expect(#msg).ToHaveText(Hello World) failed")

	titleLoc := page.Locator("#title")
	must.NoError(playwright.Expect(titleLoc).ToHaveText(ctx, "Page Title"), "Expect(#title).ToHaveText(Page Title) failed")
}

// TestExpectToHaveValue verifies input value assertion.
func TestExpectToHaveValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="username" type="text" value="testuser" />
		<textarea id="bio">My bio text</textarea>
	`)
	must.NoError(err, "SetContent failed")

	inputLoc := page.Locator("#username")
	must.NoError(playwright.Expect(inputLoc).ToHaveValue(ctx, "testuser"), "Expect(#username).ToHaveValue(testuser) failed")

	textareaLoc := page.Locator("#bio")
	must.NoError(playwright.Expect(textareaLoc).ToHaveValue(ctx, "My bio text"), "Expect(#bio).ToHaveValue(My bio text) failed")
}

// TestExpectNot verifies that negated assertions work correctly.
func TestExpectNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="visible">Visible Element</div>
		<div id="hidden" style="display:none">Hidden</div>
		<button id="active">Active</button>
		<button id="inactive" disabled>Inactive</button>
	`)
	must.NoError(err, "SetContent failed")

	// Not().ToBeVisible on a hidden element should pass
	hiddenLoc := page.Locator("#hidden")
	must.NoError(playwright.Expect(hiddenLoc).Not().ToBeVisible(ctx), "Expect(#hidden).Not().ToBeVisible() failed")

	// Not().ToBeDisabled on an enabled button should pass
	activeLoc := page.Locator("#active")
	must.NoError(playwright.Expect(activeLoc).Not().ToBeDisabled(ctx), "Expect(#active).Not().ToBeDisabled() failed")

	// Not().ToBeEnabled on a disabled button should pass
	inactiveLoc := page.Locator("#inactive")
	must.NoError(playwright.Expect(inactiveLoc).Not().ToBeEnabled(ctx), "Expect(#inactive).Not().ToBeEnabled() failed")
}

// TestExpectToHaveCount verifies the count assertion.
func TestExpectToHaveCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
			<li>Item 3</li>
		</ul>
		<p class="note">Note</p>
	`)
	must.NoError(err, "SetContent failed")

	liLoc := page.Locator("li")
	must.NoError(playwright.Expect(liLoc).ToHaveCount(ctx, 3), "Expect(li).ToHaveCount(3) failed")

	noteLoc := page.Locator("p.note")
	must.NoError(playwright.Expect(noteLoc).ToHaveCount(ctx, 1), "Expect(p.note).ToHaveCount(1) failed")

	// Zero count for non-existent elements
	missingLoc := page.Locator("#missing")
	must.NoError(playwright.Expect(missingLoc).ToHaveCount(ctx, 0), "Expect(#missing).ToHaveCount(0) failed")
}

// TestExpectToBeEditable verifies that editable (not readonly, not disabled) inputs pass.
func TestExpectToBeEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="editable" type="text" />
		<input id="readonly" type="text" readonly />
		<input id="disabled" type="text" disabled />
	`)
	must.NoError(err, "SetContent failed")

	editable := page.Locator("#editable")
	must.NoError(playwright.Expect(editable).ToBeEditable(ctx), "Expect(#editable).ToBeEditable() failed")

	readonly := page.Locator("#readonly")
	must.NoError(playwright.Expect(readonly).Not().ToBeEditable(ctx), "Expect(#readonly).Not().ToBeEditable() failed")

	disabled := page.Locator("#disabled")
	must.NoError(playwright.Expect(disabled).Not().ToBeEditable(ctx), "Expect(#disabled).Not().ToBeEditable() failed")
}

// TestExpectToHaveAttribute verifies that an element has a specific attribute value.
func TestExpectToHaveAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<a id="link" href="https://example.com" target="_blank">link</a>
		<button id="btn" type="submit">Submit</button>
	`)
	must.NoError(err, "SetContent failed")

	link := page.Locator("#link")
	must.NoError(playwright.Expect(link).ToHaveAttribute(ctx, "href", "https://example.com"), "Expect(#link).ToHaveAttribute(href) failed")
	must.NoError(playwright.Expect(link).ToHaveAttribute(ctx, "target", "_blank"), "Expect(#link).ToHaveAttribute(target) failed")

	btn := page.Locator("#btn")
	must.NoError(playwright.Expect(btn).ToHaveAttribute(ctx, "type", "submit"), "Expect(#btn).ToHaveAttribute(type) failed")
}

// TestExpectToBeChecked verifies checked state assertion for checkboxes.
func TestExpectToBeChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="checked-box" type="checkbox" checked />
		<input id="unchecked-box" type="checkbox" />
	`)
	must.NoError(err, "SetContent failed")

	checkedLoc := page.Locator("#checked-box")
	must.NoError(playwright.Expect(checkedLoc).ToBeChecked(ctx), "Expect(#checked-box).ToBeChecked() failed")

	uncheckedLoc := page.Locator("#unchecked-box")
	must.NoError(playwright.Expect(uncheckedLoc).Not().ToBeChecked(ctx), "Expect(#unchecked-box).Not().ToBeChecked() failed")
}

// TestExpectToBeHidden verifies that hidden elements pass ToBeHidden and visible elements fail.
// Ref: TestLocatorAssertionsToBeHidden
func TestExpectToBeHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<details>
			<summary>click to open</summary>
			<ul>
				<li>hidden item 1</li>
				<li>hidden item 2</li>
			</ul>
		</details>
	`)
	must.NoError(err, "SetContent failed")

	// <summary> is visible → ToBeHidden should fail, Not().ToBeHidden should pass
	summary := page.Locator("summary")
	is.Error(playwright.Expect(summary).ToBeHidden(ctx), "visible element should fail ToBeHidden")
	must.NoError(playwright.Expect(summary).Not().ToBeHidden(ctx), "visible element should pass Not().ToBeHidden")

	// <ul> inside closed <details> is hidden → ToBeHidden should pass
	ul := page.Locator("ul")
	must.NoError(playwright.Expect(ul).ToBeHidden(ctx), "hidden element should pass ToBeHidden")
	is.Error(playwright.Expect(ul).Not().ToBeHidden(ctx), "hidden element should fail Not().ToBeHidden")
}

// TestExpectToBeFocused verifies focused state assertion.
// Ref: TestLocatorAssertionsToBeFocused
func TestExpectToBeFocused(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="input1"><input id="input2">`)
	must.NoError(err, "SetContent failed")

	input1 := page.Locator("#input1")
	must.NoError(input1.Focus(ctx), "Focus failed")
	must.NoError(playwright.Expect(input1).ToBeFocused(ctx), "focused input should pass ToBeFocused")
	is.Error(playwright.Expect(input1).Not().ToBeFocused(ctx), "focused input should fail Not().ToBeFocused")

	input2 := page.Locator("#input2")
	is.Error(playwright.Expect(input2).ToBeFocused(ctx), "unfocused input should fail ToBeFocused")
	must.NoError(playwright.Expect(input2).Not().ToBeFocused(ctx), "unfocused input should pass Not().ToBeFocused")
}

// TestExpectToBeEmpty verifies empty state assertion for inputs and divs.
// Ref: TestLocatorAssertionsToBeEmpty
func TestExpectToBeEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<textarea id="empty-ta"></textarea>
		<textarea id="filled-ta">text</textarea>
		<div id="empty-div"></div>
		<div id="filled-div">content</div>
	`)
	must.NoError(err, "SetContent failed")

	must.NoError(playwright.Expect(page.Locator("#empty-ta")).ToBeEmpty(ctx))
	is.Error(playwright.Expect(page.Locator("#empty-ta")).Not().ToBeEmpty(ctx))

	is.Error(playwright.Expect(page.Locator("#filled-ta")).ToBeEmpty(ctx))
	must.NoError(playwright.Expect(page.Locator("#filled-ta")).Not().ToBeEmpty(ctx))

	must.NoError(playwright.Expect(page.Locator("#empty-div")).ToBeEmpty(ctx))
	is.Error(playwright.Expect(page.Locator("#filled-div")).ToBeEmpty(ctx))
}

// TestExpectToContainText verifies substring text assertion.
// Ref: TestLocatorAssertionsToContainText
func TestExpectToContainText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Hello World test123</div>`)
	must.NoError(err, "SetContent failed")

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainText(ctx, "Hello"), "should contain 'Hello'")
	must.NoError(playwright.Expect(loc).ToContainText(ctx, "test123"), "should contain 'test123'")
	is.Error(playwright.Expect(loc).ToContainText(ctx, "missing"), "should not contain 'missing'")
	must.NoError(playwright.Expect(loc).Not().ToContainText(ctx, "missing"), "Not().ToContainText('missing') should pass")
}

// TestExpectToHaveClass verifies CSS class assertion.
// Ref: TestLocatorAssertionsToHaveClass
func TestExpectToHaveClass(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="a" class="foo bar">a</div>
		<div id="b" class="baz">b</div>
	`)
	must.NoError(err, "SetContent failed")

	a := page.Locator("#a")
	must.NoError(playwright.Expect(a).ToHaveClass(ctx, "foo bar"), "should have class 'foo bar'")
	is.Error(playwright.Expect(a).Not().ToHaveClass(ctx, "foo bar"), "Not().ToHaveClass should fail for matching class")

	b := page.Locator("#b")
	is.Error(playwright.Expect(b).ToHaveClass(ctx, "foo"), "div#b should not have class 'foo'")
	must.NoError(playwright.Expect(b).Not().ToHaveClass(ctx, "foo"), "Not().ToHaveClass should pass for non-matching")
}

// TestExpectToHaveID verifies element id assertion.
func TestExpectToHaveID(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="my-element">content</div>`)
	must.NoError(err, "SetContent failed")

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveID(ctx, "my-element"), "should have id 'my-element'")
	is.Error(playwright.Expect(loc).ToHaveID(ctx, "other-id"), "should fail for wrong id")
	must.NoError(playwright.Expect(loc).Not().ToHaveID(ctx, "other-id"), "Not().ToHaveID should pass for wrong id")
}

// TestExpectToHaveCSS verifies computed CSS property assertion.
// Ref: TestLocatorAssertionsToHaveCSS
func TestExpectToHaveCSS(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" style="display:flex;color:red">hi</div>`)
	must.NoError(err, "SetContent failed")

	loc := page.Locator("#el")
	must.NoError(playwright.Expect(loc).ToHaveCSS(ctx, "display", "flex"), "should have display:flex")
	is.Error(playwright.Expect(loc).ToHaveCSS(ctx, "display", "block"), "should fail for wrong display value")
	must.NoError(playwright.Expect(loc).Not().ToHaveCSS(ctx, "display", "block"), "Not().ToHaveCSS should pass for wrong value")
}

// TestExpectToHaveJSProperty verifies JavaScript property assertion.
// Ref: TestLocatorAssertionsToHaveJSProperty
func TestExpectToHaveJSProperty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">element</div>`)
	must.NoError(err, "SetContent failed")

	// Set a custom JS property on the element
	_, err = page.Evaluate(ctx, `() => { document.getElementById('el').myProp = 42; }`)
	must.NoError(err, "Evaluate failed")

	loc := page.Locator("#el")
	must.NoError(playwright.Expect(loc).ToHaveJSProperty(ctx, "myProp", float64(42)), "should have myProp=42")
	is.Error(playwright.Expect(loc).ToHaveJSProperty(ctx, "myProp", float64(99)), "should fail for wrong value")
}

// TestExpectToHaveValues verifies multi-select value assertion.
// TestExpectToBeAttached verifies that ToBeAttached passes when the element is in the DOM
// and Not().ToBeAttached passes when it has been removed.
// Ref: TestLocatorAssertionsToBeAttached
func TestExpectToBeAttached(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">attached</div>`)
	must.NoError(err, "SetContent failed")

	loc := page.Locator("#el")
	must.NoError(playwright.Expect(loc).ToBeAttached(ctx), "should be attached")
	is.Error(playwright.Expect(loc).Not().ToBeAttached(ctx), "Not().ToBeAttached should fail when element is in DOM")

	_, err = page.Evaluate(ctx, `() => document.getElementById('el').remove()`)
	must.NoError(err, "Evaluate remove failed")

	must.NoError(playwright.Expect(loc).Not().ToBeAttached(ctx), "should not be attached after remove")
}

// TestExpectToBeInViewport verifies that ToBeInViewport passes for elements within the visible area
// and fails for elements positioned far outside the viewport.
// Ref: TestLocatorAssertionsToBeInViewport
func TestExpectToBeInViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="in-view">Visible element</div>
		<div id="out-view" style="position:absolute;top:10000px">Offscreen element</div>
	`)
	must.NoError(err, "SetContent failed")

	inView := page.Locator("#in-view")
	must.NoError(playwright.Expect(inView).ToBeInViewport(ctx), "visible element should be in viewport")

	outView := page.Locator("#out-view")
	must.NoError(playwright.Expect(outView).Not().ToBeInViewport(ctx), "offscreen element should not be in viewport")
}

// TestExpectToHaveAccessibleName verifies ARIA accessible name assertion.
// Ref: TestLocatorAssertionsToHaveAccessibleName
func TestExpectToHaveAccessibleName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button aria-label="Submit form">Click</button>
		<input id="email" aria-label="Email address" type="email">
	`)
	must.NoError(err, "SetContent failed")

	btn := page.Locator("button")
	must.NoError(playwright.Expect(btn).ToHaveAccessibleName(ctx, "Submit form"), "button should have accessible name")
	is.Error(playwright.Expect(btn).ToHaveAccessibleName(ctx, "wrong name"), "should fail for wrong name")
	must.NoError(playwright.Expect(btn).Not().ToHaveAccessibleName(ctx, "wrong name"), "Not() should pass for wrong name")

	input := page.Locator("#email")
	must.NoError(playwright.Expect(input).ToHaveAccessibleName(ctx, "Email address"), "input should have accessible name")
}

// TestExpectToHaveAccessibleDescription verifies ARIA accessible description assertion.
// Ref: TestLocatorAssertionsToHaveAccessibleDescription
func TestExpectToHaveAccessibleDescription(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="pwd" type="password" aria-describedby="hint">
		<div id="hint">Must be at least 8 characters</div>
	`)
	must.NoError(err, "SetContent failed")

	input := page.Locator("#pwd")
	must.NoError(playwright.Expect(input).ToHaveAccessibleDescription(ctx, "Must be at least 8 characters"), "should have accessible description")
	is.Error(playwright.Expect(input).ToHaveAccessibleDescription(ctx, "wrong hint"), "should fail for wrong description")
	must.NoError(playwright.Expect(input).Not().ToHaveAccessibleDescription(ctx, "wrong hint"), "Not() should pass for wrong description")
}

// TestExpectToHaveRole verifies ARIA role assertion.
// Ref: TestLocatorAssertionsToHaveRole
func TestExpectToHaveRole(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div role="button" tabindex="0">Custom button</div>
		<nav role="navigation">Nav</nav>
		<button>Real button</button>
	`)
	must.NoError(err, "SetContent failed")

	customBtn := page.Locator("div[role='button']")
	must.NoError(playwright.Expect(customBtn).ToHaveRole(ctx, "button"), "div with role=button should have role button")
	is.Error(playwright.Expect(customBtn).ToHaveRole(ctx, "link"), "should fail for wrong role")

	nav := page.Locator("nav")
	must.NoError(playwright.Expect(nav).ToHaveRole(ctx, "navigation"), "nav element should have role navigation")

	realBtn := page.Locator("button")
	must.NoError(playwright.Expect(realBtn).ToHaveRole(ctx, "button"), "button element should have implicit role button")
	must.NoError(playwright.Expect(realBtn).Not().ToHaveRole(ctx, "link"), "Not().ToHaveRole should pass for wrong role")
}

// Ref: TestLocatorAssertionsToHaveValues
func TestExpectToHaveValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<select id="sel" multiple>
			<option value="a" selected>A</option>
			<option value="b" selected>B</option>
			<option value="c">C</option>
		</select>
	`)
	must.NoError(err, "SetContent failed")

	loc := page.Locator("#sel")
	must.NoError(playwright.Expect(loc).ToHaveValues(ctx, []string{"a", "b"}), "should have values [a, b]")
	is.Error(playwright.Expect(loc).ToHaveValues(ctx, []string{"a", "b", "c"}), "should fail when 'c' is not selected")
	must.NoError(playwright.Expect(loc).Not().ToHaveValues(ctx, []string{"c"}), "Not().ToHaveValues should pass for unselected value")
}

// TestExpectToContainTextRegex verifies ToContainTextRegex matches element text against a regex.
// Ref: TestLocatorAssertions.java#containsTextWRegexPass
func TestExpectToContainTextRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<div id="el">Hello World Playwright</div>`))

	loc := page.Locator("#el")
	must.NoError(playwright.Expect(loc).ToContainTextRegex(ctx, regexp.MustCompile(`World`)), "should match 'World'")
	must.NoError(playwright.Expect(loc).ToContainTextRegex(ctx, regexp.MustCompile(`(?i)hello`)), "case-insensitive match should pass")
	is.Error(playwright.Expect(loc).ToContainTextRegex(ctx, regexp.MustCompile(`NOMATCH`)), "should fail for non-matching regex")
}

// TestExpectToHaveTextRegex verifies ToHaveTextRegex matches element text against a regex.
// Ref: TestLocatorAssertions.java#hasTextWRegexPass
func TestExpectToHaveTextRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<div id="msg">hello world</div>`))

	loc := page.Locator("#msg")
	must.NoError(playwright.Expect(loc).ToHaveTextRegex(ctx, regexp.MustCompile(`hello world`)), "exact regex match should pass")
	must.NoError(playwright.Expect(loc).ToHaveTextRegex(ctx, regexp.MustCompile(`(?i)HELLO`)), "case-insensitive regex match should pass")
	is.Error(playwright.Expect(loc).ToHaveTextRegex(ctx, regexp.MustCompile(`^world$`)), "non-matching regex should fail")
}

// TestExpectToHaveTextArray verifies ToHaveTextArray checks text across multiple matched elements.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPass
func TestExpectToHaveTextArray(t *testing.T) {
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

	loc := page.Locator("li")
	must.NoError(playwright.Expect(loc).ToHaveTextArray(ctx, []string{"Apple", "Banana", "Cherry"}), "all texts should match")
	is.Error(playwright.Expect(loc).ToHaveTextArray(ctx, []string{"Apple", "Banana", "Durian"}), "wrong third element should fail")
}

// TestExpectToHaveAttributeRegex verifies ToHaveAttributeRegex matches attribute value against regex.
// Ref: TestLocatorAssertions.java#hasAttributeRegExpPass
func TestExpectToHaveAttributeRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<a id="lnk" href="https://example.com/path?query=1">link</a>`))

	loc := page.Locator("#lnk")
	must.NoError(playwright.Expect(loc).ToHaveAttributeRegex(ctx, "href", regexp.MustCompile(`example\.com`)), "regex should match href")
	is.Error(playwright.Expect(loc).ToHaveAttributeRegex(ctx, "href", regexp.MustCompile(`^nomatch`)), "non-matching regex should fail")
}

// TestExpectToHaveAttributeIgnoreCase verifies ToHaveAttributeIgnoreCase performs case-insensitive comparison.
// Ref: TestLocatorAssertions.java#hasAttributeTextIgnoreCase
func TestExpectToHaveAttributeIgnoreCase(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<div id="el" data-val="HelloWorld">x</div>`))

	loc := page.Locator("#el")
	must.NoError(playwright.Expect(loc).ToHaveAttributeIgnoreCase(ctx, "data-val", "helloworld"), "ignore-case match should pass")
	is.Error(playwright.Expect(loc).ToHaveAttribute(ctx, "data-val", "helloworld"), "case-sensitive match should fail for different case")
}

// TestExpectToHaveClassRegex verifies ToHaveClassRegex matches class attribute against a regex.
// Ref: TestLocatorAssertions.java#hasClassRegExpPass
func TestExpectToHaveClassRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<div class="btn btn-primary active">x</div>`))

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveClassRegex(ctx, regexp.MustCompile(`btn-primary`)), "regex matching class substring should pass")
	is.Error(playwright.Expect(loc).ToHaveClassRegex(ctx, regexp.MustCompile(`^btn-secondary`)), "non-matching class regex should fail")
}

// TestExpectToHaveValueRegex verifies ToHaveValueRegex matches input value against a regex.
// Ref: TestLocatorAssertions.java#hasValueRegExpPass
func TestExpectToHaveValueRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<input id="inp" value="foo@example.com">`))

	loc := page.Locator("#inp")
	must.NoError(playwright.Expect(loc).ToHaveValueRegex(ctx, regexp.MustCompile(`@example\.com$`)), "regex matching email suffix should pass")
	is.Error(playwright.Expect(loc).ToHaveValueRegex(ctx, regexp.MustCompile(`^nomatch`)), "non-matching regex should fail")
}

// TestExpectToBeCheckedFalse verifies ToBeCheckedFalse asserts that a checkbox is unchecked.
// Ref: TestLocatorAssertions.java#isCheckedFalsePass
func TestExpectToBeCheckedFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="unchecked">
		<input type="checkbox" id="checked" checked>
	`))

	unchecked := page.Locator("#unchecked")
	checked := page.Locator("#checked")

	must.NoError(playwright.Expect(unchecked).ToBeCheckedFalse(ctx), "unchecked box should pass ToBeCheckedFalse")
	is.Error(playwright.Expect(checked).ToBeCheckedFalse(ctx), "checked box should fail ToBeCheckedFalse")
}

// TestExpectToContainClass verifies ToContainClass checks presence of a specific class name.
// Ref: TestLocatorAssertions.java#containsClassPass
func TestExpectToContainClass(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<div class="foo bar baz">x</div>`))

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "foo"), "should contain class 'foo'")
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "bar"), "should contain class 'bar'")
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "baz"), "should contain class 'baz'")
	is.Error(playwright.Expect(loc).ToContainClass(ctx, "qux"), "should fail for absent class 'qux'")
	must.NoError(playwright.Expect(loc).Not().ToContainClass(ctx, "qux"), "Not().ToContainClass should pass for absent class")
}
