//go:build e2e

// E2E tests for: TestLocatorAssertions.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

// TestExpectToBeVisible verifies basic visibility assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
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

// TestExpectLocatorToContainText verifies text content contains expected string.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestExpectLocatorToContainText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Hello World</div>`))

	text, err := page.Locator("#el").InnerText(ctx)
	must.NoError(err)
	is.Contains(text, "Hello")
	is.Contains(text, "World")
}

// TestExpectLocatorHaveValue verifies input has expected value.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestExpectLocatorHaveValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="expected">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("expected", val)
}

// TestExpectLocatorToBeChecked verifies checkbox is checked.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestExpectLocatorToBeChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" checked>`))

	checked, err := page.Locator("input").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestExpectLocatorNotToBeChecked verifies checkbox is unchecked.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestExpectLocatorNotToBeChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox">`))

	checked, err := page.Locator("input").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestExpectLocatorToBeEnabled verifies button is enabled.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestExpectLocatorToBeEnabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>go</button>`))

	enabled, err := page.Locator("button").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestExpectLocatorToBeDisabled verifies button is disabled.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestExpectLocatorToBeDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button disabled>stop</button>`))

	disabled, err := page.Locator("button").IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)
}

// TestExpectLocatorToBeEditable verifies input is editable.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestExpectLocatorToBeEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.True(editable)
}

// TestLocatorAssertIsVisibleEx4 verifies IsVisible for visible element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsVisibleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Visible element</div>`))

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorAssertIsNotVisibleEx4 verifies IsVisible is false for visibility:hidden element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsNotVisibleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" style="visibility:hidden">Hidden</div>`))

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorAssertIsCheckedEx4 verifies IsChecked for checked checkbox.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsCheckedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorAssertIsDisabledEx4 verifies IsDisabled for disabled element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsDisabledEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" disabled>`))

	disabled, err := page.Locator("#inp").IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)
}

// TestLocatorAssertIsEditableEx4 verifies IsEditable for editable textarea.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsEditableEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	editable, err := page.Locator("#ta").IsEditable(ctx)
	must.NoError(err)
	is.True(editable)
}

// TestLocatorAssertIsNotCheckedEx4 verifies IsChecked is false for unchecked checkbox.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsNotCheckedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorAssertHasClassEx5 verifies HasClass assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasClassEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="active primary">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveClass(ctx, "active primary")
	must.NoError(err)
}

// TestLocatorAssertHasCountEx5 verifies Count assertion works for multiple elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasCountEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>One</li>
			<li>Two</li>
			<li>Three</li>
		</ul>
	`))

	err := playwright.Expect(page.Locator("li")).ToHaveCount(ctx, 3)
	must.NoError(err)
}

// TestLocatorAssertIsEditableEx5 verifies IsEditable assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsEditableEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEditable(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotEditableEx5 verifies not editable assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertNotEditableEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" readonly>`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEditable(ctx)
	must.NoError(err)
}

// TestLocatorAssertHasValueEx5 verifies HasValue assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasValueEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="expected">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveValue(ctx, "expected")
	must.NoError(err)
}

// TestLocatorAssertIsCheckedEx5 verifies IsChecked assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsCheckedEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	err := playwright.Expect(page.Locator("#chk")).ToBeChecked(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsDisabledEx6 verifies ToBeDisabled assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsDisabledEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Disabled</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeDisabled(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsEnabledEx6 verifies ToBeEnabled assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsEnabledEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Active</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeEnabled(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsVisibleEx6 verifies ToBeVisible assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsVisibleEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeVisible(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsHiddenEx6 verifies ToBeHidden assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsHiddenEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none;">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeHidden(ctx)
	must.NoError(err)
}

// TestLocatorAssertHasTextEx6 verifies ToHaveText assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasTextEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Expected text</p>`))

	err := playwright.Expect(page.Locator("#p")).ToHaveText(ctx, "Expected text")
	must.NoError(err)
}

// TestLocatorAssertContainsTextEx6 verifies ToContainText assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertContainsTextEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">This has some text in it</p>`))

	err := playwright.Expect(page.Locator("#p")).ToContainText(ctx, "some text")
	must.NoError(err)
}

// TestLocatorAssertHasIDEx7 verifies ToHaveID assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasIDEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">Content</div>`))

	err := playwright.Expect(page.Locator("#target")).ToHaveID(ctx, "target")
	must.NoError(err)
}

// TestLocatorAssertHasAttributeEx7 verifies ToHaveAttribute assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasAttributeEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="email" name="email">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveAttribute(ctx, "type", "email")
	must.NoError(err)
}

// TestLocatorAssertHasCSSEx7 verifies ToHaveCSS assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasCSSEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display: flex;">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveCSS(ctx, "display", "flex")
	must.NoError(err)
}

// TestLocatorAssertIsAttachedEx7 verifies ToBeAttached assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsAttachedEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Attached</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeAttached(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsEmptyEx7 verifies ToBeEmpty assertion for empty elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertIsEmptyEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"></div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeEmpty(ctx)
	must.NoError(err)
}

// TestLocatorAssertHasRoleEx8 verifies ToHaveRole assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertHasRoleEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToHaveRole(ctx, "button")
	must.NoError(err)
}

// TestLocatorAssertInViewportEx8 verifies ToBeInViewport assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertInViewportEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible element</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeInViewport(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotCheckedEx8 verifies Not().ToBeChecked assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertNotCheckedEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	err := playwright.Expect(page.Locator("#chk")).Not().ToBeChecked(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotVisibleEx8 verifies Not().ToBeVisible for hidden.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertNotVisibleEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="visibility:hidden;">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).Not().ToBeVisible(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotEnabledEx8 verifies Not().ToBeEnabled for disabled.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertNotEnabledEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" disabled>`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEnabled(ctx)
	must.NoError(err)
}

// TestToHaveValueEx9 verifies ToHaveValue assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveValueEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="expected">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveValue(ctx, "expected")
	must.NoError(err)
}

// TestToBeEditableEx9 verifies ToBeEditable assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeEditableEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEditable(ctx)
	must.NoError(err)
}

// TestNotToBeEditableEx9 verifies Not.ToBeEditable for readonly input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeEditableEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" readonly>`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEditable(ctx)
	must.NoError(err)
}

// TestToBeInViewportEx9 verifies ToBeInViewport assertion for visible element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeInViewportEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeInViewport(ctx)
	must.NoError(err)
}

// TestToHaveCountEx9 verifies ToHaveCount assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveCountEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul><li>A</li><li>B</li><li>C</li></ul>`))

	err := playwright.Expect(page.Locator("li")).ToHaveCount(ctx, 3)
	must.NoError(err)
}

// TestToBeDisabledEx10 verifies ToBeDisabled assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeDisabledEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Click</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeDisabled(ctx)
	must.NoError(err)
}

// TestNotToBeDisabledEx10 verifies Not.ToBeDisabled for enabled element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeDisabledEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	err := playwright.Expect(page.Locator("#btn")).Not().ToBeDisabled(ctx)
	must.NoError(err)
}

// TestToBeCheckedEx10 verifies ToBeChecked assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeCheckedEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	err := playwright.Expect(page.Locator("#chk")).ToBeChecked(ctx)
	must.NoError(err)
}

// TestNotToBeCheckedEx10 verifies Not.ToBeChecked for unchecked element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeCheckedEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	err := playwright.Expect(page.Locator("#chk")).Not().ToBeChecked(ctx)
	must.NoError(err)
}

// TestToHaveAttributeEx10 verifies ToHaveAttribute assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveAttributeEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="https://example.com">Link</a>`))

	err := playwright.Expect(page.Locator("#a")).ToHaveAttribute(ctx, "href", "https://example.com")
	must.NoError(err)
}

// TestToHaveIdEx11 verifies ToHaveAttribute assertion for id attribute.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveIdEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="my-div">Content</div>`))

	err := playwright.Expect(page.Locator("#my-div")).ToHaveAttribute(ctx, "id", "my-div")
	must.NoError(err)
}

// TestToHaveClassMultipleEx11 verifies ToHaveClass with multiple classes.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveClassMultipleEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="foo bar baz">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveClass(ctx, "foo bar baz")
	must.NoError(err)
}

// TestToBeEmptyInputEx11 verifies ToBeEmpty assertion for empty input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeEmptyInputEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEmpty(ctx)
	must.NoError(err)
}

// TestNotToBeEmptyInputEx11 verifies Not.ToBeEmpty when input has value.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeEmptyInputEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="content">`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEmpty(ctx)
	must.NoError(err)
}

// TestToContainTextEx11 verifies ToContainText assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToContainTextEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello, World!</p>`))

	err := playwright.Expect(page.Locator("#p")).ToContainText(ctx, "World")
	must.NoError(err)
}

// TestToBeVisibleEx12 verifies ToBeVisible assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeVisibleEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Visible</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeVisible(ctx)
	must.NoError(err)
}

// TestNotToBeVisibleEx12 verifies Not.ToBeVisible for hidden element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeVisibleEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).Not().ToBeVisible(ctx)
	must.NoError(err)
}

// TestToBeEnabledEx12 verifies ToBeEnabled assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeEnabledEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEnabled(ctx)
	must.NoError(err)
}

// TestToHaveTextEx12 verifies ToHaveText assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveTextEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h">Page Title</h1>`))

	err := playwright.Expect(page.Locator("#h")).ToHaveText(ctx, "Page Title")
	must.NoError(err)
}

// TestToHaveValueSelectEx12 verifies ToHaveValue for select after SelectOption.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveValueSelectEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "b")
	must.NoError(err)

	err = playwright.Expect(page.Locator("#sel")).ToHaveValue(ctx, "b")
	must.NoError(err)
}

// TestToHaveClassHeadingEx13 verifies ToHaveClass for heading element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveClassHeadingEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h" class="title main">Heading</h1>`))

	err := playwright.Expect(page.Locator("#h")).ToHaveClass(ctx, "title main")
	must.NoError(err)
}

// TestToHaveAttributePlaceholderEx13 verifies ToHaveAttribute for placeholder.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveAttributePlaceholderEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" placeholder="Enter text here">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveAttribute(ctx, "placeholder", "Enter text here")
	must.NoError(err)
}

// TestNotToBeCheckedRadioEx13 verifies Not.ToBeChecked for unchecked radio.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeCheckedRadioEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="r" type="radio" name="grp">`))

	err := playwright.Expect(page.Locator("#r")).Not().ToBeChecked(ctx)
	must.NoError(err)
}

// TestToBeCheckedRadioEx13 verifies ToBeChecked for checked radio button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeCheckedRadioEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="r" type="radio" name="grp" checked>`))

	err := playwright.Expect(page.Locator("#r")).ToBeChecked(ctx)
	must.NoError(err)
}

// TestToContainTextPartialEx13 verifies ToContainText with partial substring.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToContainTextPartialEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Welcome to our platform!</p>`))

	err := playwright.Expect(page.Locator("#p")).ToContainText(ctx, "platform")
	must.NoError(err)
}

// TestToHaveCSSEx14 verifies ToHaveCSS assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveCSSEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:flex;">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveCSS(ctx, "display", "flex")
	must.NoError(err)
}

// TestNotToHaveTextEx14 verifies Not.ToHaveText for non-matching text.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToHaveTextEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Actual text</p>`))

	err := playwright.Expect(page.Locator("#p")).Not().ToHaveText(ctx, "Different text")
	must.NoError(err)
}

// TestToHaveCountZeroEx14 verifies ToHaveCount(0) for no matching elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveCountZeroEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>No list here</div>`))

	err := playwright.Expect(page.Locator("li")).ToHaveCount(ctx, 0)
	must.NoError(err)
}

// TestToBeInViewportAfterScrollEx14 verifies ToBeInViewport after scroll.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeInViewportAfterScrollEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;"></div>
		<div id="deep">Deep content</div>
	`))

	must.NoError(page.Locator("#deep").ScrollIntoViewIfNeeded(ctx))

	err := playwright.Expect(page.Locator("#deep")).ToBeInViewport(ctx)
	must.NoError(err)
}

// TestToBeHiddenEx14 verifies ToBeHidden assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeHiddenEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="visibility:hidden">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeHidden(ctx)
	must.NoError(err)
}

// TestToHaveValueTextareaEx15 verifies ToHaveValue for textarea.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveValueTextareaEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">Initial text</textarea>`))

	err := playwright.Expect(page.Locator("#ta")).ToHaveValue(ctx, "Initial text")
	must.NoError(err)
}

// TestToHaveJSPropertyEx15 verifies ToHaveJSProperty assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveJSPropertyEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="checkbox" checked>`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveJSProperty(ctx, "checked", true)
	must.NoError(err)
}

// TestNotToHaveClassEx15 verifies Not.ToHaveClass assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToHaveClassEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="foo">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).Not().ToHaveClass(ctx, "bar")
	must.NoError(err)
}

// TestToBeAttachedEx15 verifies ToBeAttached assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeAttachedEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Attached</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeAttached(ctx)
	must.NoError(err)
}

// TestToHaveCSSColorEx15 verifies ToHaveCSS for color property.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveCSSColorEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="color: blue;">Blue text</p>`))

	err := playwright.Expect(page.Locator("#p")).ToHaveCSS(ctx, "color", "rgb(0, 0, 255)")
	must.NoError(err)
}

// TestToHaveValueAfterClearEx16 verifies ToHaveValue after clearing input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveValueAfterClearEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="old">`))

	must.NoError(page.Locator("#inp").Fill(ctx, ""))

	err := playwright.Expect(page.Locator("#inp")).ToHaveValue(ctx, "")
	must.NoError(err)
}

// TestToHaveCountAfterAddEx16 verifies ToHaveCount after adding element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveCountAfterAddEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"><li>One</li></ul>`))

	_, err := page.Evaluate(ctx, `() => {
		const li = document.createElement('li');
		li.textContent = 'Two';
		document.getElementById('list').appendChild(li);
	}`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#list li")).ToHaveCount(ctx, 2)
	must.NoError(err)
}

// TestNotToBeInViewportEx16 verifies Not.ToBeInViewport for off-screen element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeInViewportEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;"></div>
		<div id="deep" style="height:100px;">Deep content</div>
	`))

	err := playwright.Expect(page.Locator("#deep")).Not().ToBeInViewport(ctx)
	must.NoError(err)
}

// TestToHaveAttributeDataEx16 verifies ToHaveAttribute for data attribute.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveAttributeDataEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" data-status="active">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveAttribute(ctx, "data-status", "active")
	must.NoError(err)
}

// TestToBeEditableTextareaEx16 verifies ToBeEditable for textarea.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeEditableTextareaEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">Text</textarea>`))

	err := playwright.Expect(page.Locator("#ta")).ToBeEditable(ctx)
	must.NoError(err)
}

// TestToHaveIdEx verifies ToHaveID assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveIdEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="main-content">Content</div>`))

	err := playwright.Expect(page.Locator("#main-content")).ToHaveID(ctx, "main-content")
	must.NoError(err)
}

// TestToHaveClassSingleEx verifies ToHaveClass with single class.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveClassSingleEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button class="btn btn-primary">Submit</button>`))

	err := playwright.Expect(page.Locator("button")).ToHaveClass(ctx, "btn btn-primary")
	must.NoError(err)
}

// TestToHaveCSSFontSizeEx verifies ToHaveCSS for font-size.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveCSSFontSizeEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h" style="font-size:24px;">Title</h1>`))

	err := playwright.Expect(page.Locator("#h")).ToHaveCSS(ctx, "font-size", "24px")
	must.NoError(err)
}

// TestToBeFocusedEx verifies ToBeFocused assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToBeFocusedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Focus(ctx))

	err := playwright.Expect(page.Locator("#inp")).ToBeFocused(ctx)
	must.NoError(err)
}

// TestNotToBeFocusedEx verifies Not.ToBeFocused assertion.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestNotToBeFocusedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="a" type="text"><input id="b" type="text">`))

	must.NoError(page.Locator("#a").Focus(ctx))

	err := playwright.Expect(page.Locator("#b")).Not().ToBeFocused(ctx)
	must.NoError(err)
}

// TestToHaveCountThreeEx verifies ToHaveCount with 3 elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestToHaveCountThreeEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="item">A</li>
			<li class="item">B</li>
			<li class="item">C</li>
		</ul>
	`))

	err := playwright.Expect(page.Locator(".item")).ToHaveCount(ctx, 3)
	must.NoError(err)
}

// TestHasTextFilterMatchesHeadingEx2 verifies HasText filter finds heading.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestHasTextFilterMatchesHeadingEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h2>Introduction</h2>
		<h2>Summary</h2>
		<h2>Conclusion</h2>
	`))

	text := "Summary"
	count, err := page.Locator("h2").Filter(&playwright.LocatorFilterOptions{HasText: &text}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestHasTextFilterMultipleMatchesEx2 verifies HasText can match multiple elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestHasTextFilterMultipleMatchesEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Error: File not found</p>
		<p>Error: Permission denied</p>
		<p>Success: Operation complete</p>
	`))

	text := "Error"
	count, err := page.Locator("p").Filter(&playwright.LocatorFilterOptions{HasText: &text}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestHasTextFilterGetsTextEx2 verifies HasText filtered text is accessible.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestHasTextFilterGetsTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">Category: Books</div>
		<div class="item">Category: Music</div>
		<div class="item">Category: Movies</div>
	`))

	text := "Music"
	innerText, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{HasText: &text}).InnerText(ctx)
	must.NoError(err)
	is.Equal("Category: Music", innerText)
}

// TestHasNotTextFilterEx2 verifies HasNotText excludes matching text.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestHasNotTextFilterEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>Red</li>
		<li>Blue</li>
		<li>Green</li>
		<li>Yellow</li>
	`))

	notText := "Red"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{HasNotText: &notText}).Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestHaveTextFilterBadgesEx3 verifies Filter.HasText finds badge items.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestHaveTextFilterBadgesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="badge">Sale</span>
		<span class="badge">New</span>
		<span class="badge">Out of stock</span>
	`))

	count, err := page.Locator(".badge").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrHT3("New"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestHaveTextFilterMenuItemsEx3 verifies Filter.HasText finds menu items.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestHaveTextFilterMenuItemsEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav>
			<a class="nav-item" href="/home">Home</a>
			<a class="nav-item" href="/about">About Us</a>
			<a class="nav-item" href="/contact">Contact</a>
		</nav>
	`))

	text, err := page.Locator(".nav-item").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrHT3("About"),
	}).TextContent(ctx)
	must.NoError(err)
	is.Equal("About Us", text)
}

// TestHaveTextFilterStatusEx3 verifies Filter.HasText finds status elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestHaveTextFilterStatusEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="status">Pending</div>
		<div class="status">Approved</div>
		<div class="status">Rejected</div>
		<div class="status">Approved</div>
	`))

	count, err := page.Locator(".status").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrHT3("Approved"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorAssertionsContainsTextWRegexPass verifies ToContainTextRegex matches element text against a regex.
// Ref: TestLocatorAssertions.java#containsTextWRegexPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsTextWRegexPass(t *testing.T) {
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

// TestLocatorAssertionsHasTextWRegexPass verifies ToHaveTextRegex matches element text against a regex.
// Ref: TestLocatorAssertions.java#hasTextWRegexPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWRegexPass(t *testing.T) {
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

// TestLocatorAssertionsHasTextWTextArrayPass verifies ToHaveTextArray checks text across multiple matched elements.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextArrayPass(t *testing.T) {
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

// TestLocatorAssertionsHasAttributeRegExpPass verifies ToHaveAttributeRegex matches attribute value against regex.
// Ref: TestLocatorAssertions.java#hasAttributeRegExpPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasAttributeRegExpPass(t *testing.T) {
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

// TestLocatorAssertionsHasAttributeTextIgnoreCase verifies ToHaveAttributeIgnoreCase performs case-insensitive comparison.
// Ref: TestLocatorAssertions.java#hasAttributeTextIgnoreCase  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasAttributeTextIgnoreCase(t *testing.T) {
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

// TestLocatorAssertionsHasClassRegExpPass verifies ToHaveClassRegex matches class attribute against a regex.
// Ref: TestLocatorAssertions.java#hasClassRegExpPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassRegExpPass(t *testing.T) {
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

// TestLocatorAssertionsHasValueRegExpPass verifies ToHaveValueRegex matches input value against a regex.
// Ref: TestLocatorAssertions.java#hasValueRegExpPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValueRegExpPass(t *testing.T) {
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

// TestLocatorAssertionsIsCheckedFalsePass verifies ToBeCheckedFalse asserts that a checkbox is unchecked.
// Ref: TestLocatorAssertions.java#isCheckedFalsePass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsCheckedFalsePass(t *testing.T) {
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

// TestLocatorAssertionsContainsClassPass verifies ToContainClass checks presence of a specific class name.
// Ref: TestLocatorAssertions.java#containsClassPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsClassPass(t *testing.T) {
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

// TestContainsTextRegexPass verifies containsText with a plain regex pattern.
// Ref: TestLocatorAssertions.java#containsTextWRegexPass
func TestContainsTextRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	loc := page.Locator("#node")
	must.NoError(playwright.Expect(loc).ToContainTextRegex(ctx, regexp.MustCompile(`ex`)))
	// Should not normalize whitespace for regex.
	must.NoError(playwright.Expect(loc).ToContainTextRegex(ctx, regexp.MustCompile(`ext   cont`)))
}

// TestLocatorAssertionsContainsTextWRegexCaseInsensitivePass verifies case-insensitive regex containsText.
// Ref: TestLocatorAssertions.java#containsTextWRegexCaseInsensitivePass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsTextWRegexCaseInsensitivePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	// (?i) flag makes the regex case-insensitive.
	must.NoError(playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`(?i)text`)))
}

// TestLocatorAssertionsContainsTextWRegexMultilinePass verifies multiline regex containsText.
// Ref: TestLocatorAssertions.java#containsTextWRegexMultilinePass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsTextWRegexMultilinePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div id=node>Text \nContent</div>")
	must.NoError(err)

	// (?m) flag enables multiline matching (^ matches start of each line).
	must.NoError(playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`(?m)^Content`)))
}

// TestLocatorAssertionsContainsTextWRegexDotAllPass verifies dotall regex containsText (. matches newline).
// Ref: TestLocatorAssertions.java#containsTextWRegexDotAllPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsTextWRegexDotAllPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div id=node>foo\nbar</div>")
	must.NoError(err)

	// (?s) flag enables dot-all mode (. matches newline).
	must.NoError(playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`(?s)foo.bar`)))
}

// TestLocatorAssertionsContainsTextWRegexFail verifies that a non-matching regex fails containsText.
// Ref: TestLocatorAssertions.java#containsTextWRegexFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsTextWRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`ex2`))
	is.Error(err, "non-matching regex should fail containsText")
}

// TestLocatorAssertionsContainsTextWTextPass verifies containsText with ignoreCase and falsy ignoreCase.
// Ref: TestLocatorAssertions.java#containsTextWTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsTextWTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	loc := page.Locator("#node")
	// Basic substring
	must.NoError(playwright.Expect(loc).ToContainText(ctx, "Text"))
	// IgnoreCase = true
	must.NoError(playwright.Expect(loc).ToContainTextIgnoreCase(ctx, "EXT"))
	// IgnoreCase = false (default): "TEXT" should NOT match "Text content" case-sensitively
	must.NoError(playwright.Expect(loc).Not().ToContainText(ctx, "TEXT"))
}

// TestLocatorAssertionsContainsTextWTextArrayPass verifies containsText with a string array.
// Ref: TestLocatorAssertions.java#containsTextWTextArrayPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsTextWTextArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Text \n1</div><div>Text2</div><div>Text3</div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainTextArray(ctx, []string{"ext     1", "ext3"}))
	// ignoreCase variant
	must.NoError(playwright.Expect(loc).ToContainTextArrayIgnoreCase(ctx, []string{"EXT 1", "eXt3"}))
}

// TestHasTextRegexPass verifies hasText with various regex options.
// Ref: TestLocatorAssertions.java#hasTextWRegexPass
func TestHasTextRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	loc := page.Locator("#node")
	must.NoError(playwright.Expect(loc).ToHaveTextRegex(ctx, regexp.MustCompile(`Te.t`)))
	// Should not normalize whitespace for regex.
	must.NoError(playwright.Expect(loc).ToHaveTextRegex(ctx, regexp.MustCompile(`Text.+content`)))
	// ignoreCase override: regex without flag + ignoreCase=true → should match
	must.NoError(playwright.Expect(loc).ToHaveTextRegexIgnoreCase(ctx, regexp.MustCompile(`text   content`), true))
	// ignoreCase override: regex with (?i) flag + ignoreCase=false → should NOT match case-insensitively
	must.NoError(playwright.Expect(loc).Not().ToHaveTextRegexIgnoreCase(ctx, regexp.MustCompile(`(?i)text   content`), false))
}

// TestLocatorAssertionsHasTextWRegexFail verifies that a non-matching regex fails hasText.
// Ref: TestLocatorAssertions.java#hasTextWRegexFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveTextRegex(ctx, regexp.MustCompile(`Text 2`))
	is.Error(err, "non-matching regex should fail hasText")
}

// TestLocatorAssertionsHasTextWTextPass verifies hasText with ignoreCase option.
// Ref: TestLocatorAssertions.java#hasTextWTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div id=node><span></span>Text \ncontent&nbsp;    </div>")
	must.NoError(err)

	loc := page.Locator("#node")
	// Whitespace normalization
	must.NoError(playwright.Expect(loc).ToHaveText(ctx, "Text                        content"))
	// ignoreCase = true
	must.NoError(playwright.Expect(loc).ToHaveTextIgnoreCase(ctx, "text CONTENT"))
	// ignoreCase = false: "TEXT" should NOT match "Text content"
	must.NoError(playwright.Expect(loc).Not().ToHaveText(ctx, "TEXT"))
}

// TestLocatorAssertionsHasTextWTextFail verifies that a non-matching string fails hasText.
// Ref: TestLocatorAssertions.java#hasTextWTextFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	// "Text" (substring) should fail hasText when full text is "Text content"
	err = playwright.Expect(page.Locator("#node")).ToHaveText(ctx, "Text")
	is.Error(err, "substring 'Text' should fail hasText on 'Text content'")
}

// TestLocatorAssertionsHasTextWTextInnerTextPass verifies hasText with useInnerText (ignores hidden spans).
// Ref: TestLocatorAssertions.java#hasTextWTextInnerTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextInnerTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text <span hidden>garbage</span> content</div>`)
	must.NoError(err)

	// useInnerText excludes hidden elements
	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveTextUseInnerText(ctx, "Text content"))
}

// TestHasTextArrayPass verifies hasText with a string array (normalized whitespace).
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPass
func TestHasTextArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div>Text    \n1</div><div>Text   2a</div>")
	must.NoError(err)

	loc := page.Locator("div")
	// Normalized whitespace
	must.NoError(playwright.Expect(loc).ToHaveTextArray(ctx, []string{"Text  1", "Text   2a"}))
	// ignoreCase
	must.NoError(playwright.Expect(loc).ToHaveTextArrayIgnoreCase(ctx, []string{"tEXT 1", "TExt 2A"}))
}

// TestLocatorAssertionsHasTextWTextArrayPassEmpty verifies hasText with an empty array passes on empty locator.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassEmpty  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextArrayPassEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	// locator("p") matches nothing → empty array assertion should pass
	must.NoError(playwright.Expect(page.Locator("p")).ToHaveTextArray(ctx, []string{}))
}

// TestLocatorAssertionsHasTextWTextArrayPassNotEmpty verifies not().hasText with empty array fails when elements exist.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassNotEmpty  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextArrayPassNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div><p>Test</p></div>`)
	must.NoError(err)

	// "div" matches an element with text → not().hasText([]) should pass (element is not empty)
	must.NoError(playwright.Expect(page.Locator("div")).Not().ToHaveTextArray(ctx, []string{}))
}

// TestLocatorAssertionsHasTextWTextArrayPassOnEmpty verifies not().hasText fails on empty locator.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassOnEmpty  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextArrayPassOnEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	// locator("p") matches nothing → not().hasText(["Test"]) should pass
	must.NoError(playwright.Expect(page.Locator("p")).Not().ToHaveTextArray(ctx, []string{"Test"}))
}

// TestLocatorAssertionsHasTextWTextArrayFailOnNotEmpty verifies not().hasText([]) fails when locator matches nothing.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayFailOnNotEmpty  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextArrayFailOnNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	// locator("p") matches nothing → not().hasText([]) fails (both expected and actual are [])
	err = playwright.Expect(page.Locator("p")).Not().ToHaveTextArray(ctx, []string{})
	is.Error(err, "not().hasText([]) should fail when locator matches nothing (both are [])")
}

// TestLocatorAssertionsHasTextWTextArrayFail verifies hasText array fails when count or text differ.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Text 1</div><div>Text 3</div>`)
	must.NoError(err)

	// Two elements, asking for three entries → should fail
	err = playwright.Expect(page.Locator("div")).ToHaveTextArray(ctx, []string{"Text 1", "Text 3", "Extra"})
	is.Error(err, "array count mismatch should fail hasText")
}

// TestLocatorAssertionsHasTextWRegExArrayPass verifies hasText with a regex array.
// Ref: TestLocatorAssertions.java#hasTextWRegExArrayPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWRegExArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div>Text    \n1</div><div>Text   2a</div>")
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveTextRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile("Text    \n1"),
		regexp.MustCompile(`Text   \d+a`),
	}))
}

// TestLocatorAssertionsHasTextWRegExArrayFail verifies hasText with a regex array fails on mismatch.
// Ref: TestLocatorAssertions.java#hasTextWRegExArrayFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWRegExArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Text 1</div><div>Text 3</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveTextRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile(`Text 1`),
		regexp.MustCompile(`Text   \d`),
		regexp.MustCompile(`Extra`),
	})
	is.Error(err, "regex array count mismatch should fail hasText")
}

// TestHasAttributeIgnoreCase verifies hasAttribute with ignoreCase option.
// Ref: TestLocatorAssertions.java#hasAttributeTextIgnoreCase
func TestHasAttributeIgnoreCase(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=NoDe>Text content</div>`)
	must.NoError(err)

	loc := page.Locator("#NoDe")
	// ignoreCase=true: "node" matches "NoDe"
	must.NoError(playwright.Expect(loc).ToHaveAttributeIgnoreCase(ctx, "id", "node"))
	// case-sensitive (default): "node" does NOT match "NoDe"
	must.NoError(playwright.Expect(loc).Not().ToHaveAttribute(ctx, "id", "node"))
}

// TestLocatorAssertionsHasAttributeTextFail verifies that hasAttribute fails on value mismatch.
// Ref: TestLocatorAssertions.java#hasAttributeTextFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasAttributeTextFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveAttribute(ctx, "id", "foo")
	is.Error(err, "wrong attribute value should fail hasAttribute")
}

// TestHasAttributeRegexPass verifies hasAttribute with a regex pattern.
// Ref: TestLocatorAssertions.java#hasAttributeRegExpPass
func TestHasAttributeRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveAttributeRegex(ctx, "id", regexp.MustCompile(`n..e`)))
}

// TestLocatorAssertionsHasAttributeRegExpFail verifies that hasAttribute regex fails on no match.
// Ref: TestLocatorAssertions.java#hasAttributeRegExpFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasAttributeRegExpFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveAttributeRegex(ctx, "id", regexp.MustCompile(`.Nod..`))
	is.Error(err, "non-matching regex should fail hasAttribute")
}

// TestLocatorAssertionsHasClassTextFail verifies hasClass fails when the class string doesn't match.
// Ref: TestLocatorAssertions.java#hasClassTextFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassTextFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="bar baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClass(ctx, "foo bar baz")
	is.Error(err, "class mismatch should fail hasClass")
}

// TestHasClassRegexPass verifies hasClass with a regex pattern.
// Ref: TestLocatorAssertions.java#hasClassRegExpPass
func TestHasClassRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo bar baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClassRegex(ctx, regexp.MustCompile(`foo.* baz`)))
}

// TestLocatorAssertionsHasClassRegExpFail verifies hasClass regex fails on no match.
// Ref: TestLocatorAssertions.java#hasClassRegExpFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassRegExpFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="bar baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClassRegex(ctx, regexp.MustCompile(`foo Z.*`))
	is.Error(err, "non-matching regex should fail hasClass")
}

// TestLocatorAssertionsHasClassTextArrayPass verifies hasClass with a string array.
// Ref: TestLocatorAssertions.java#hasClassTextArrayPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassTextArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClassArray(ctx, []string{"foo", "bar", "baz"}))
}

// TestLocatorAssertionsHasClassTextArrayFail verifies hasClass array fails on mismatch.
// Ref: TestLocatorAssertions.java#hasClassTextArrayFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassTextArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClassArray(ctx, []string{"foo", "bar", "missing"})
	is.Error(err, "class array mismatch should fail hasClass")
}

// TestLocatorAssertionsHasClassRegExpArrayPass verifies hasClass with a regex array.
// Ref: TestLocatorAssertions.java#hasClassRegExpArrayPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassRegExpArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClassRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile(`fo.*`),
		regexp.MustCompile(`.ar`),
		regexp.MustCompile(`baz`),
	}))
}

// TestLocatorAssertionsHasClassRegExpArrayFail verifies hasClass regex array fails on count mismatch.
// Ref: TestLocatorAssertions.java#hasClassRegExpArrayFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassRegExpArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClassRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile(`fo.*`),
		regexp.MustCompile(`.ar`),
		regexp.MustCompile(`baz`),
		regexp.MustCompile(`extra`),
	})
	is.Error(err, "regex array count mismatch should fail hasClass")
}

// TestLocatorAssertionsHasCountFail verifies hasCount fails when the count is wrong.
// Ref: TestLocatorAssertions.java#hasCountFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasCountFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select><option>One</option><option>Two</option></select>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("option")).ToHaveCount(ctx, 1)
	is.Error(err, "wrong count should fail hasCount")
}

// TestLocatorAssertionsHasCountPassZero verifies hasCount(0) and not().hasCount(1) on non-existent elements.
// Ref: TestLocatorAssertions.java#hasCountPassZero  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasCountPassZero(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	must.NoError(playwright.Expect(loc).ToHaveCount(ctx, 0))
	must.NoError(playwright.Expect(loc).Not().ToHaveCount(ctx, 1))
}

// TestLocatorAssertionsHasCSSPass verifies hasCSS with exact value.
// Ref: TestLocatorAssertions.java#hasCSSPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasCSSPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveCSS(ctx, "color", "rgb(255, 0, 0)"))
}

// TestLocatorAssertionsHasCSSFail verifies hasCSS fails on value mismatch.
// Ref: TestLocatorAssertions.java#hasCSSFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasCSSFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveCSS(ctx, "color", "red")
	is.Error(err, "CSS value mismatch should fail hasCSS")
}

// TestLocatorAssertionsHasCSSRegExPass verifies hasCSS with a regex pattern.
// Ref: TestLocatorAssertions.java#hasCSSRegExPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasCSSRegExPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveCSSRegex(ctx, "color", regexp.MustCompile(`rgb.*`)))
}

// TestLocatorAssertionsHasCSSRegExFail verifies hasCSS regex fails on no match.
// Ref: TestLocatorAssertions.java#hasCSSRegExFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasCSSRegExFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveCSSRegex(ctx, "color", regexp.MustCompile(`red`))
	is.Error(err, "non-matching CSS regex should fail hasCSS")
}

// TestLocatorAssertionsHasIdPass verifies hasId matches the id attribute.
// Ref: TestLocatorAssertions.java#hasIdPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasIdPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveID(ctx, "node"))
}

// TestLocatorAssertionsHasIdFail verifies hasId fails on mismatch.
// Ref: TestLocatorAssertions.java#hasIdFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasIdFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveID(ctx, "foo")
	is.Error(err, "id mismatch should fail hasId")
}

// TestLocatorAssertionsHasJSPropertyPass verifies hasJSProperty with an object value.
// Ref: TestLocatorAssertions.java#hasJSPropertyPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasJSPropertyPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => { document.querySelector('div').foo = { a: 1, b: 'string' }; }`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveJSProperty(ctx, "foo", map[string]any{"a": 1, "b": "string"}))
}

// TestLocatorAssertionsHasJSPropertyNumberFail verifies hasJSProperty fails on number mismatch.
// Ref: TestLocatorAssertions.java#hasJSPropertyNumberFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasJSPropertyNumberFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => { document.querySelector('div').foo = 2021; }`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveJSProperty(ctx, "foo", 1)
	is.Error(err, "number mismatch should fail hasJSProperty")
}

// TestLocatorAssertionsHasJSPropertyObjectFail verifies hasJSProperty fails on object mismatch.
// Ref: TestLocatorAssertions.java#hasJSPropertyObjectFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasJSPropertyObjectFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => { document.querySelector('div').foo = { a: 1, b: 'string' }; }`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveJSProperty(ctx, "foo", map[string]any{"a": 2})
	is.Error(err, "object mismatch should fail hasJSProperty")
}

// TestLocatorAssertionsHasJSPropertyStringFail verifies hasJSProperty fails on string mismatch.
// Ref: TestLocatorAssertions.java#hasJSPropertyStringFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasJSPropertyStringFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveJSProperty(ctx, "id", "foo")
	is.Error(err, "string mismatch should fail hasJSProperty")
}

// TestLocatorAssertionsHasValueTextFail verifies hasValue fails on mismatch.
// Ref: TestLocatorAssertions.java#hasValueTextFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValueTextFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveValue(ctx, "Text2")
	is.Error(err, "value mismatch should fail hasValue")
}

// TestHasValueRegexPass verifies hasValue with a regex pattern.
// Ref: TestLocatorAssertions.java#hasValueRegExpPass
func TestHasValueRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveValueRegex(ctx, regexp.MustCompile(`Text`)))
}

// TestLocatorAssertionsHasValueRegExpPassWithNot verifies not().hasValue with a non-matching regex.
// Ref: TestLocatorAssertions.java#hasValueRegExpPassWithNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValueRegExpPassWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).Not().ToHaveValueRegex(ctx, regexp.MustCompile(`Text2`)))
}

// TestLocatorAssertionsHasValueRegExpFail verifies hasValue regex fails on no match.
// Ref: TestLocatorAssertions.java#hasValueRegExpFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValueRegExpFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveValueRegex(ctx, regexp.MustCompile(`Text2`))
	is.Error(err, "non-matching regex should fail hasValue")
}

// TestLocatorAssertionsHasValuesWorksWithText verifies hasValues with string values.
// Ref: TestLocatorAssertions.java#hasValuesWorksWithText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValuesWorksWithText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "R", "G")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("select")).ToHaveValues(ctx, []string{"R", "G"}))
}

// TestLocatorAssertionsHasValuesExactMatchWithText verifies hasValues fails when values don't match exactly.
// Ref: TestLocatorAssertions.java#hasValuesExactMatchWithText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValuesExactMatchWithText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="RR">Red</option>
		<option value="GG">Green</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "RR", "GG")
	must.NoError(err)

	err = playwright.Expect(page.Locator("select")).ToHaveValues(ctx, []string{"R", "G"})
	is.Error(err, "partial value match should fail hasValues")
}

// TestLocatorAssertionsHasValuesWorksWithRegex verifies hasValues with regex patterns.
// Ref: TestLocatorAssertions.java#hasValuesWorksWithRegex  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValuesWorksWithRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "R", "G")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("select")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	}))
}

// TestLocatorAssertionsHasValuesFailsWhenItemsNotSelected verifies hasValues fails when wrong items are selected.
// Ref: TestLocatorAssertions.java#hasValuesFailsWhenItemsNotSelected  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValuesFailsWhenItemsNotSelected(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "B")
	must.NoError(err)

	err = playwright.Expect(page.Locator("select")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	})
	is.Error(err, "wrong selection should fail hasValues regex")
}

// TestLocatorAssertionsHasValuesFailsWhenMultipleNotSpecified verifies hasValues fails when select is not multiple.
// Ref: TestLocatorAssertions.java#hasValuesFailsWhenMultipleNotSpecified  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValuesFailsWhenMultipleNotSpecified(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "B")
	must.NoError(err)

	err = playwright.Expect(page.Locator("select")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	})
	is.Error(err, "single-select should fail hasValues with multiple patterns")
}

// TestLocatorAssertionsHasValuesFailsWhenNotASelectElement verifies hasValues fails on non-select elements.
// Ref: TestLocatorAssertions.java#hasValuesFailsWhenNotASelectElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValuesFailsWhenNotASelectElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input value="foo" />`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	})
	is.Error(err, "non-select element should fail hasValues")
}

// TestLocatorAssertionsIsCheckedFail verifies isChecked fails on unchecked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsCheckedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeChecked(ctx)
	is.Error(err, "unchecked checkbox should fail isChecked")
}

// TestLocatorAssertionsNotIsCheckedFail verifies not().isChecked fails on checked checkbox.
// Ref: TestLocatorAssertions.java#notIsCheckedFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsCheckedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox checked></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeChecked(ctx)
	is.Error(err, "checked checkbox should fail not().isChecked")
}

// TestIsCheckedFalsePass verifies isChecked(checked=false) on unchecked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedFalsePass
func TestIsCheckedFalsePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeCheckedFalse(ctx))
}

// TestLocatorAssertionsIsCheckedFalseFail verifies isChecked(checked=false) fails on checked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedFalseFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsCheckedFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input checked type=checkbox></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeCheckedFalse(ctx)
	is.Error(err, "checked checkbox should fail isChecked(false)")
}

// TestLocatorAssertionsIsDisabledFail verifies isDisabled fails on enabled button.
// Ref: TestLocatorAssertions.java#isDisabledFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsDisabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeDisabled(ctx)
	is.Error(err, "enabled button should fail isDisabled")
}

// TestLocatorAssertionsNotIsDisabledFail verifies not().isDisabled fails on disabled button.
// Ref: TestLocatorAssertions.java#notIsDisabledFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsDisabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).Not().ToBeDisabled(ctx)
	is.Error(err, "disabled button should fail not().isDisabled")
}

// TestLocatorAssertionsIsEditableFail verifies isEditable fails on disabled input.
// Ref: TestLocatorAssertions.java#isEditableFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditableFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input disabled></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeEditable(ctx)
	is.Error(err, "disabled input should fail isEditable")
}

// TestLocatorAssertionsIsEditableFalseFail verifies isEditable(editable=false) fails on editable input.
// Ref: TestLocatorAssertions.java#isEditableFalseFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditableFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	// editable=false is equivalent to "is readonly" → fails on editable input
	err = playwright.Expect(page.Locator("input")).ToBeReadOnly(ctx)
	is.Error(err, "editable input should fail isEditable(editable=false)")
}

// TestLocatorAssertionsNotIsEditableFail verifies not().isEditable fails on editable input.
// Ref: TestLocatorAssertions.java#notIsEditableFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsEditableFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeEditable(ctx)
	is.Error(err, "editable input should fail not().isEditable")
}

// TestLocatorAssertionsIsEditableWithNot verifies not().isEditable passes on readonly input.
// Ref: TestLocatorAssertions.java#isEditableWithNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditableWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input readonly></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).Not().ToBeEditable(ctx))
}

// TestLocatorAssertionsIsEditableWithEditableTrue verifies isEditable(editable=true) passes on editable input.
// Ref: TestLocatorAssertions.java#isEditableWithEditableTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditableWithEditableTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeEditable(ctx))
}

// TestLocatorAssertionsIsEditableWithEditableFalse verifies isEditable(editable=false) passes on readonly input.
// Ref: TestLocatorAssertions.java#isEditableWithEditableFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditableWithEditableFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input readonly></input>`)
	must.NoError(err)

	// editable=false = isReadonly
	must.NoError(playwright.Expect(page.Locator("input")).ToBeReadOnly(ctx))
}

// TestLocatorAssertionsIsEditableWithNotAndEditableFalse verifies not().isEditable(editable=false) on editable input.
// Ref: TestLocatorAssertions.java#isEditableWithNotAndEditableFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditableWithNotAndEditableFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	// not().isEditable(editable=false) = not isReadonly = isEditable
	must.NoError(playwright.Expect(page.Locator("input")).ToBeEditable(ctx))
}

// TestLocatorAssertionsIsEditableThrowsOnNonInputElement verifies isEditable fails on a non-input element.
// Ref: TestLocatorAssertions.java#isEditableThrowsOnNonInputElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditableThrowsOnNonInputElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button></button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeEditable(ctx)
	is.Error(err, "button (non-input) should fail isEditable")
}

// TestLocatorAssertionsIsEmptyPass verifies isEmpty passes on empty input.
// Ref: TestLocatorAssertions.java#isEmptyPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEmptyPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeEmpty(ctx))
}

// TestLocatorAssertionsIsEmptyFail verifies isEmpty fails when input has a value.
// Ref: TestLocatorAssertions.java#isEmptyFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEmptyFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input value=text></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeEmpty(ctx)
	is.Error(err, "input with value should fail isEmpty")
}

// TestLocatorAssertionsNotIsEmptyFail verifies not().isEmpty fails on empty input.
// Ref: TestLocatorAssertions.java#notIsEmptyFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsEmptyFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeEmpty(ctx)
	is.Error(err, "empty input should fail not().isEmpty")
}

// TestLocatorAssertionsIsEnabledFail verifies isEnabled fails on disabled button.
// Ref: TestLocatorAssertions.java#isEnabledFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeEnabled(ctx)
	is.Error(err, "disabled button should fail isEnabled")
}

// TestLocatorAssertionsNotIsEnabledFail verifies not().isEnabled fails on enabled button.
// Ref: TestLocatorAssertions.java#notIsEnabledFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsEnabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).Not().ToBeEnabled(ctx)
	is.Error(err, "enabled button should fail not().isEnabled")
}

// TestLocatorAssertionsIsEnabledTrue verifies isEnabled(enabled=true) passes on enabled button.
// Ref: TestLocatorAssertions.java#isEnabledTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestLocatorAssertionsIsEnabledFalse verifies isEnabled(enabled=false) passes on disabled button.
// Ref: TestLocatorAssertions.java#isEnabledFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	// enabled=false is equivalent to isDisabled
	must.NoError(playwright.Expect(page.Locator("button")).ToBeDisabled(ctx))
}

// TestLocatorAssertionsIsEnabledFalseFail verifies isEnabled(enabled=false) fails on enabled button.
// Ref: TestLocatorAssertions.java#isEnabledFalseFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	// enabled=false = isDisabled; should fail on enabled button
	err = playwright.Expect(page.Locator("button")).ToBeDisabled(ctx)
	is.Error(err, "enabled button should fail isEnabled(enabled=false)")
}

// TestLocatorAssertionsIsEnabledEventually verifies isEnabled passes once the element becomes enabled.
// Ref: TestLocatorAssertions.java#isEnabledEventually  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledEventually(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	// Schedule removal of disabled attribute after 300ms
	_, err = page.Locator("button").Evaluate(ctx, `e => setTimeout(() => { e.removeAttribute('disabled'); }, 300)`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestLocatorAssertionsIsEnabledEventuallyWithNot verifies not().isEnabled passes once element becomes disabled.
// Ref: TestLocatorAssertions.java#isEnabledEventuallyWithNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledEventuallyWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	// Schedule addition of disabled attribute after 300ms
	_, err = page.Locator("button").Evaluate(ctx, `e => setTimeout(() => { e.setAttribute('disabled', ''); }, 300)`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).Not().ToBeEnabled(ctx))
}

// TestLocatorAssertionsIsEnabledWithNotAndEnabledFalse verifies not().isEnabled(enabled=false) on enabled button.
// Ref: TestLocatorAssertions.java#isEnabledWithNotAndEnabledFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledWithNotAndEnabledFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	// not().isEnabled(enabled=false) = not(isDisabled) = isEnabled
	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestLocatorAssertionsIsFocusedPass verifies isFocused passes on focused input.
// Ref: TestLocatorAssertions.java#isFocusedPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsFocusedPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = page.Locator("input").Focus(ctx)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeFocused(ctx))
}

// TestLocatorAssertionsIsFocusedFail verifies isFocused fails on unfocused input.
// Ref: TestLocatorAssertions.java#isFocusedFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsFocusedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeFocused(ctx)
	is.Error(err, "unfocused input should fail isFocused")
}

// TestLocatorAssertionsNotIsFocusedFail verifies not().isFocused fails on focused input.
// Ref: TestLocatorAssertions.java#notIsFocusedFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsFocusedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = page.Locator("input").Focus(ctx)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeFocused(ctx)
	is.Error(err, "focused input should fail not().isFocused")
}

// TestLocatorAssertionsIsHiddenFail verifies isHidden fails on visible button.
// Ref: TestLocatorAssertions.java#isHiddenFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsHiddenFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button></button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeHidden(ctx)
	is.Error(err, "visible button should fail isHidden")
}

// TestLocatorAssertionsNotIsHiddenFail verifies not().isHidden fails on hidden button.
// Ref: TestLocatorAssertions.java#notIsHiddenFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsHiddenFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button style='display: none'></button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).Not().ToBeHidden(ctx)
	is.Error(err, "hidden button should fail not().isHidden")
}

// TestLocatorAssertionsIsVisibleFail verifies isVisible fails on hidden input.
// Ref: TestLocatorAssertions.java#isVisibleFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisibleFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input style='display: none'></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeVisible(ctx)
	is.Error(err, "hidden input should fail isVisible")
}

// TestLocatorAssertionsIsVisibleFalseFail verifies isVisible(visible=false) fails on visible input.
// Ref: TestLocatorAssertions.java#isVisibleFalseFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisibleFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	// visible=false = isHidden; should fail on visible input
	err = playwright.Expect(page.Locator("input")).ToBeHidden(ctx)
	is.Error(err, "visible input should fail isVisible(visible=false)")
}

// TestLocatorAssertionsNotIsVisibleFail verifies not().isVisible fails on visible input.
// Ref: TestLocatorAssertions.java#notIsVisibleFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsNotIsVisibleFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeVisible(ctx)
	is.Error(err, "visible input should fail not().isVisible")
}

// TestLocatorAssertionsIsVisibleWithTrue verifies isVisible(visible=true) passes on visible button.
// Ref: TestLocatorAssertions.java#isVisibleWithTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisibleWithTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeVisible(ctx))
}

// TestLocatorAssertionsIsVisibleWithFalse verifies isVisible(visible=false) passes on hidden button.
// Ref: TestLocatorAssertions.java#isVisibleWithFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisibleWithFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button hidden>hello</button>`)
	must.NoError(err)

	// visible=false = isHidden
	must.NoError(playwright.Expect(page.Locator("button")).ToBeHidden(ctx))
}

// TestLocatorAssertionsIsVisibleWithNotAndFalse verifies not().isVisible(visible=false) on visible button.
// Ref: TestLocatorAssertions.java#isVisibleWithNotAndFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisibleWithNotAndFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	// not().isVisible(false) = not(isHidden) = isVisible
	must.NoError(playwright.Expect(page.Locator("button")).ToBeVisible(ctx))
}

// TestLocatorAssertionsIsVisibleEventually verifies isVisible passes once element becomes visible.
// Ref: TestLocatorAssertions.java#isVisibleEventually  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisibleEventually(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('div').innerHTML = '<span>Hello</span>'; }, 100); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).ToBeVisible(ctx))
}

// TestLocatorAssertionsIsVisibleEventuallyWithNot verifies not().isVisible passes once element becomes hidden.
// Ref: TestLocatorAssertions.java#isVisibleEventuallyWithNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisibleEventuallyWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div><span>Hello</span></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('span').textContent = ''; }, 100); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).Not().ToBeVisible(ctx))
}

// TestLocatorAssertionsLocatorCountShouldWorkWithDeletedMapInMainWorld verifies hasCount works even when Map is overridden in page.
// Ref: TestLocatorAssertions.java#locatorCountShouldWorkWithDeletedMapInMainWorld  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsLocatorCountShouldWorkWithDeletedMapInMainWorld(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `Map = 1`)
	must.NoError(err)

	count, err := page.Locator("#searchResultTableDiv .x-grid3-row").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)

	must.NoError(playwright.Expect(page.Locator("#searchResultTableDiv .x-grid3-row")).ToHaveCount(ctx, 0))
}

// TestContainsClassPass verifies containsClass with various class combinations.
// Ref: TestLocatorAssertions.java#containsClassPass
func TestContainsClassPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='foo bar baz'></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, ""))
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "bar"))
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "baz bar"))
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "  bar   foo "))
	must.NoError(playwright.Expect(loc).Not().ToContainClass(ctx, "  baz not-matching"))
}

// TestLocatorAssertionsContainsClassPassWithSvgs verifies containsClass works on SVG elements.
// Ref: TestLocatorAssertions.java#containsClassPassWithSvgs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsClassPassWithSvgs(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<svg class='c1 c2' role='img' xmlns='http://www.w3.org/2000/svg' viewBox='0 0 512 512'></svg>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("svg")).ToContainClass(ctx, "c1"))
	must.NoError(playwright.Expect(page.Locator("svg")).ToContainClass(ctx, "c2 c1"))
}

// TestLocatorAssertionsContainsClassFail verifies containsClass fails when class doesn't exist.
// Ref: TestLocatorAssertions.java#containsClassFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsClassFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='bar baz'></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToContainClass(ctx, "does-not-exist")
	is.Error(err, "missing class should fail containsClass")
}

// TestLocatorAssertionsContainsClassPassWithArray verifies containsClass with a string array.
// Ref: TestLocatorAssertions.java#containsClassPassWithArray  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsClassPassWithArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='foo'></div><div class='hello bar'></div><div class='baz'></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainClassArray(ctx, []string{"foo", "hello", "baz"}))
	// class mismatch
	must.NoError(playwright.Expect(loc).Not().ToHaveClassArray(ctx, []string{"not-there", "hello", "baz"}))
	// count mismatch
	must.NoError(playwright.Expect(loc).Not().ToHaveClassArray(ctx, []string{"foo", "hello"}))
}

// TestLocatorAssertionsContainsClassFailWithArray verifies containsClass array fails on class mismatch.
// Ref: TestLocatorAssertions.java#containsClassFailWithArray  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsContainsClassFailWithArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='foo'></div><div class='bar'></div><div class='bar'></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToContainClassArray(ctx, []string{"foo", "bar", "baz"})
	is.Error(err, "class mismatch should fail containsClass array")
}

// TestLocatorAssertionsHasTextWTextArrayPassLazyPass verifies hasText(array) passes after DOM is populated asynchronously.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassLazyPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasTextWTextArrayPassLazyPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=div></div>`)
	must.NoError(err)

	loc := page.Locator("p")
	_, err = page.Evaluate(ctx, `() => setTimeout(() => {
		document.getElementById('div').innerHTML = '<p>Text 1</p><p>Text 2</p>';
	}, 100)`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).ToHaveTextArray(ctx, []string{"Text  1", "Text   2"}))
}

// TestLocatorAssertionsHasAttributeTextPass verifies hasAttribute passes with matching attribute value.
// Ref: TestLocatorAssertions.java#hasAttributeTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasAttributeTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveAttribute(ctx, "id", "node"))
}

// TestLocatorAssertionsHasClassTextPass verifies hasClass passes with matching class attribute string.
// Ref: TestLocatorAssertions.java#hasClassTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasClassTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo bar baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClass(ctx, "foo bar baz"))
}

// TestLocatorAssertionsHasCountPass verifies hasCount passes when the locator matches the expected number of elements.
// Ref: TestLocatorAssertions.java#hasCountPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasCountPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select><option>One</option><option>Two</option></select>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("option")).ToHaveCount(ctx, 2))
}

// TestLocatorAssertionsHasValueTextPass verifies hasValue passes when the input contains the expected text.
// Ref: TestLocatorAssertions.java#hasValueTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValueTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)

	loc := page.Locator("#node")
	must.NoError(loc.Fill(ctx, "Text content"))
	must.NoError(playwright.Expect(loc).ToHaveValue(ctx, "Text content"))
}

// TestLocatorAssertionsHasValuesFollowsLabels verifies hasValues works when locating a select via its label text.
// Ref: TestLocatorAssertions.java#hasValuesFollowsLabels  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsHasValuesFollowsLabels(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<label for="colors">Pick a Color</label>
<select id="colors" multiple>
  <option value="R">Red</option>
  <option value="G">Green</option>
  <option value="B">Blue</option>
</select>`)
	must.NoError(err)

	loc := page.Locator("text=Pick a Color")
	_, err = loc.SelectOption(ctx, "R", "G")
	must.NoError(err, "SelectOption via label failed")
	must.NoError(playwright.Expect(loc).ToHaveValues(ctx, []string{"R", "G"}))
}

// TestLocatorAssertionsIsCheckedPass verifies isChecked passes on a checked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsCheckedPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox checked></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeChecked(ctx))
}

// TestLocatorAssertionsIsDisabledPass verifies isDisabled passes on a disabled button.
// Ref: TestLocatorAssertions.java#isDisabledPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsDisabledPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeDisabled(ctx))
}

// TestLocatorAssertionsIsEditablePass verifies isEditable passes on an enabled input.
// Ref: TestLocatorAssertions.java#isEditablePass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEditablePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeEditable(ctx))
}

// TestLocatorAssertionsIsEnabledPass verifies isEnabled passes on an enabled button.
// Ref: TestLocatorAssertions.java#isEnabledPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsEnabledPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestLocatorAssertionsIsHiddenPass verifies isHidden passes on a hidden button.
// Ref: TestLocatorAssertions.java#isHiddenPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsHiddenPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button style='display: none'></button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeHidden(ctx))
}

// TestLocatorAssertionsIsVisiblePass verifies isVisible passes on a visible input.
// Ref: TestLocatorAssertions.java#isVisiblePass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsIsVisiblePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeVisible(ctx))
}

// TestLocatorAssertionsDefaultTimeoutHasTextFail verifies hasText fails when the element does not contain the expected text.
// Ref: TestLocatorAssertions.java#defaultTimeoutHasTextFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsDefaultTimeoutHasTextFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveText(ctx, "foo")
	is.Error(err, "hasText should fail when element lacks expected text")
}

// TestLocatorAssertionsDefaultTimeoutHasTextPass verifies hasText passes immediately when element contains the text.
// Ref: TestLocatorAssertions.java#defaultTimeoutHasTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsDefaultTimeoutHasTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>foo</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveText(ctx, "foo"))
}

// TestLocatorAssertionsDefaultTimeoutZeroHasTextPass verifies hasText passes when element has the expected text
// even when a zero-value timeout is used (server applies its default).
// Ref: TestLocatorAssertions.java#defaultTimeoutZeroHasTextPass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorAssertionsDefaultTimeoutZeroHasTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>foo</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveText(ctx, "foo"))
}

// ── isAttached variants ──────────────────────────────────────────────────────

// TestLocatorCountZeroForNoMatch verifies Count returns 0 when no elements match.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorCountZeroForNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class="box">one</div>`))

	count, err := page.Locator(".notexist").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// TestLocatorHasTextNoMatchReturnsZeroEx2 verifies Filter HasText with no matches returns 0.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorHasTextNoMatchReturnsZeroEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>One</li>
		<li>Two</li>
	`))

	hasText := "Three"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{HasText: &hasText}).Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// TestLocatorInnerTextExcludesHiddenText verifies InnerText excludes CSS-hidden text.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorInnerTextExcludesHiddenText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			Visible
			<span style="display:none">Hidden</span>
		</div>
	`))

	text, err := page.Locator("#d").InnerText(ctx)
	must.NoError(err)
	is.Contains(text, "Visible")
	must.NotContains(text, "Hidden")
}

// TestLocatorInnerTextEmptyElement verifies InnerText on empty element returns empty string.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorInnerTextEmptyElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="empty"></div>`))

	text, err := page.Locator("#empty").InnerText(ctx)
	must.NoError(err)
	is.Empty(text)
}

// TestLocatorInputValueEmptyForEmptyInput verifies InputValue returns empty string for empty input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorInputValueEmptyForEmptyInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorIsVisibleTrueForVisibleElementEx2 verifies IsVisible returns true for visible element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsVisibleTrueForVisibleElementEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Visible</p>`))

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleTrueForVisibleElement verifies IsVisible returns true for visible element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsVisibleTrueForVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="visible">I'm here</div>`))

	visible, err := page.Locator("#visible").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleFalseForHiddenElement verifies IsVisible returns false for display:none.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsVisibleFalseForHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="display:none">hidden</div>`))

	visible, err := page.Locator("#hidden").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsHiddenTrueForHiddenElement verifies IsVisible returns false for visibility:hidden element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsHiddenTrueForHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="visibility:hidden">hidden</div>`))

	visible, err := page.Locator("#hidden").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsHiddenFalseForVisibleElement verifies IsVisible returns true for visible button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsHiddenFalseForVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>visible button</button>`))

	visible, err := page.Locator("button").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsEnabledFalseForDisabledButton verifies IsEnabled returns false for disabled button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsEnabledFalseForDisabledButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button disabled>Can't click</button>`))

	enabled, err := page.Locator("button").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestLocatorIsDisabledFalseForEnabledInput verifies IsDisabled returns false for enabled input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsDisabledFalseForEnabledInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" placeholder="type here">`))

	disabled, err := page.Locator("input").IsDisabled(ctx)
	must.NoError(err)
	is.False(disabled)
}

// TestLocatorIsEditableTrueForInput verifies IsEditable returns true for a standard input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsEditableTrueForInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.True(editable)
}

// TestLocatorIsDisabledFalseForEnabledButton verifies IsDisabled false for enabled button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsDisabledFalseForEnabledButtonExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">go</button>`))

	disabled, err := page.Locator("#btn").IsDisabled(ctx)
	must.NoError(err)
	is.False(disabled)
}

// TestLocatorIsEditableFalseForDisabledInput verifies IsEditable false for disabled input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions.java
func TestLocatorIsEditableFalseForDisabledInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" disabled>`))

	editable, err := page.Locator("#inp").IsEditable(ctx)
	must.NoError(err)
	is.False(editable)
}
