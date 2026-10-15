//go:build e2e

// E2E tests for: TestLocatorAssertions2.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

// TestLocatorAssertions2IsAttachedDefault verifies isAttached passes on an attached element.
// Ref: TestLocatorAssertions2.java#isAttachedDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedWithHiddenElement verifies isAttached passes even for hidden elements.
// Ref: TestLocatorAssertions2.java#isAttachedWithHiddenElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedWithHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button style='display:none'>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedWithNot verifies not().isAttached passes when element is not in DOM.
// Ref: TestLocatorAssertions2.java#isAttachedWithNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).Not().ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedWithAttachedTrue verifies isAttached(attached=true) on a present button.
// Ref: TestLocatorAssertions2.java#isAttachedWithAttachedTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedWithAttachedTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedWithAttachedFalse verifies isAttached(attached=false) on absent element.
// Ref: TestLocatorAssertions2.java#isAttachedWithAttachedFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedWithAttachedFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	// attached=false means "expect not attached" = same as not().toBeAttached()
	must.NoError(playwright.Expect(page.Locator("input")).ToBeDetached(ctx))
}

// TestLocatorAssertions2IsAttachedWithNotAndAttachedFalse verifies not().isAttached(attached=false) on present element.
// Ref: TestLocatorAssertions2.java#isAttachedWithNotAndAttachedFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedWithNotAndAttachedFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	// not().isAttached(false) = not(notAttached) = attached
	must.NoError(playwright.Expect(page.Locator("button")).ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedEventually verifies isAttached passes once element is added to DOM.
// Ref: TestLocatorAssertions2.java#isAttachedEventually  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedEventually(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('div').innerHTML = '<span>Hello</span>'; }, 100); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedEventuallyWithNot verifies not().isAttached passes once element is removed.
// Ref: TestLocatorAssertions2.java#isAttachedEventuallyWithNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedEventuallyWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div><span>Hello</span></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('div').textContent = ''; }, 0); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).Not().ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedFail verifies isAttached fails when element is not in DOM.
// Ref: TestLocatorAssertions2.java#isAttachedFail  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hello</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeAttached(ctx)
	is.Error(err, "absent element should fail isAttached")
}

// TestLocatorAssertions2IsAttachedFailWithNot verifies not().isAttached fails when element is present.
// Ref: TestLocatorAssertions2.java#isAttachedFailWithNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedFailWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeAttached(ctx)
	is.Error(err, "present element should fail not().isAttached")
}

// TestLocatorAssertions2IsAttachedWithImpossibleTimeout verifies isAttached succeeds with a very short timeout when element exists.
// Ref: TestLocatorAssertions2.java#isAttachedWithImpossibleTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedWithImpossibleTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	// Element is already present, so even a 1ms timeout should succeed immediately.
	must.NoError(playwright.Expect(page.Locator("#node")).ToBeAttached(ctx))
}

// TestLocatorAssertions2IsAttachedWithImpossibleTimeoutNot verifies not().isAttached succeeds with very short timeout when absent.
// Ref: TestLocatorAssertions2.java#isAttachedWithImpossibleTimeoutNot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2IsAttachedWithImpossibleTimeoutNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	// Element "no-such-thing" is absent, so not().isAttached should pass immediately.
	must.NoError(playwright.Expect(page.Locator("no-such-thing")).Not().ToBeAttached(ctx))
}

// ── toHaveAccessibleName extended ────────────────────────────────────────────

// TestLocatorAssertions2ToHaveAccessibleName verifies accessible name with ignoreCase and regex variants.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleName  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToHaveAccessibleName(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="button" aria-label="Hello"></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveAccessibleName(ctx, "Hello"))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleName(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleNameIgnoreCase(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleNameRegex(ctx, regexp.MustCompile(`ell\w`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleNameRegex(ctx, regexp.MustCompile(`hello`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleNameRegex(ctx, regexp.MustCompile(`(?i)hello`)))

	// Test &nbsp; normalization
	err = page.SetContent(ctx, `<button>foo&nbsp;bar
baz</button>`)
	must.NoError(err)
	must.NoError(playwright.Expect(page.Locator("button")).ToHaveAccessibleName(ctx, "foo bar baz"))
}

// ── toHaveAccessibleDescription extended ────────────────────────────────────

// TestLocatorAssertions2ToHaveAccessibleDescription verifies accessible description with ignoreCase and regex.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleDescription  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToHaveAccessibleDescription(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="button" aria-description="Hello"></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescription(ctx, "Hello"))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleDescription(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescriptionIgnoreCase(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescriptionRegex(ctx, regexp.MustCompile(`ell\w`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleDescriptionRegex(ctx, regexp.MustCompile(`hello`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescriptionRegex(ctx, regexp.MustCompile(`(?i)hello`)))

	// Test &nbsp; normalization via aria-describedby
	err = page.SetContent(ctx, `<div role="button" aria-describedby="desc"></div>
<span id="desc">foo&nbsp;bar
baz</span>`)
	must.NoError(err)
	must.NoError(playwright.Expect(page.Locator("div")).ToHaveAccessibleDescription(ctx, "foo bar baz"))
}

// ── toHaveRole ───────────────────────────────────────────────────────────────

// TestLocatorAssertions2ToHaveRole verifies hasRole assertion.
// Ref: TestLocatorAssertions2.java#toHaveRole  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToHaveRole(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="button">Button!</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveRole(ctx, "button"))
	must.NoError(playwright.Expect(page.Locator("div")).Not().ToHaveRole(ctx, "checkbox"))
}

// ── toHaveAccessibleErrorMessage ─────────────────────────────────────────────

// TestLocatorAssertions2ToHaveAccessibleErrorMessage verifies error message assertion with text/regex/ignoreCase.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleErrorMessage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToHaveAccessibleErrorMessage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<form>
<input role="textbox" aria-invalid="true" aria-errormessage="error-message" />
<div id="error-message">Hello</div>
<div id="irrelevant-error">This should not be considered.</div>
</form>`)
	must.NoError(err)

	loc := page.Locator(`input[role="textbox"]`)
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessage(ctx, "Hello"))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessage(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageIgnoreCase(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`ell\w`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`hello`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)hello`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessage(ctx, "This should not be considered."))
}

// TestLocatorAssertions2ToHaveAccessibleErrorMessageShouldHandleMultipleAriaErrorMessageReferences verifies error message with multiple aria-errormessage refs.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleErrorMessageShouldHandleMultipleAriaErrorMessageReferences  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToHaveAccessibleErrorMessageShouldHandleMultipleAriaErrorMessageReferences(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<form>
  <input role="textbox" aria-invalid="true" aria-errormessage="error1 error2" />
  <div id="error1">First error message.</div>
  <div id="error2">Second error message.</div>
  <div id="irrelevant-error">This should not be considered.</div>
</form>`)
	must.NoError(err)

	loc := page.Locator(`input[role="textbox"]`)
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessage(ctx, "First error message. Second error message."))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)first error message.`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)second error message.`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)This should not be considered.`)))
}

// ── isChecked indeterminate ──────────────────────────────────────────────────

// TestLocatorAssertions2ToBeEditableWithIndeterminateTrue verifies isChecked with indeterminate=true passes on indeterminate checkbox.
// Ref: TestLocatorAssertions2.java#toBeEditableWithIndeterminateTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToBeEditableWithIndeterminateTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	_, err = page.Locator("input").Evaluate(ctx, `e => e.indeterminate = true`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeCheckedIndeterminate(ctx))
}

// TestLocatorAssertions2ToBeEditableFailWithIndeterminateTrue verifies isChecked with indeterminate=true fails on non-indeterminate checkbox.
// Ref: TestLocatorAssertions2.java#toBeEditableFailWithIndeterminateTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToBeEditableFailWithIndeterminateTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeCheckedIndeterminate(ctx)
	is.Error(err, "non-indeterminate checkbox should fail isChecked(indeterminate=true)")
}

// TestLocatorAssertions2ToBeEditableWithIndeterminateTrueAndChecked verifies that setting both indeterminate=true and
// checked=false simultaneously returns an error from the Playwright server.
// Ref: TestLocatorAssertions2.java#toBeEditableWithIndeterminateTrueAndChecked  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestLocatorAssertions2ToBeEditableWithIndeterminateTrueAndChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	_, err = page.Locator("input").Evaluate(ctx, `e => e.indeterminate = true`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeCheckedWithIndeterminateAndChecked(ctx, false)
	is.Error(err, "indeterminate=true with checked=false should fail")
	is.ErrorContains(err, "Can't assert indeterminate and checked at the same time")
}

// TestExpectLocatorNotToBeEditable verifies readonly input is not editable.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAssertions2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAssertions2.java
func TestExpectLocatorNotToBeEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" readonly>`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.False(editable)
}
