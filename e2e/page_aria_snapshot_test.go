//go:build e2e

// E2E tests for: TestPageAriaSnapshot.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestExpectToMatchAriaSnapshot verifies the ToMatchAriaSnapshot assertion.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestExpectToMatchAriaSnapshot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/aria-assert", "text/html", `<button>Submit</button>`)
	err := page.Goto(ctx, srv.Prefix()+"/aria-assert")
	must.NoError(err, "Goto failed")

	// The snapshot should contain a button role with name "Submit"
	expected := "- button \"Submit\""
	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx, expected)
	must.NoError(err, "ToMatchAriaSnapshot failed")
}

// TestExpectToMatchAriaSnapshotNegated verifies the negated form of ToMatchAriaSnapshot.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestExpectToMatchAriaSnapshotNegated(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/aria-neg", "text/html", `<button>Save</button>`)
	err := page.Goto(ctx, srv.Prefix()+"/aria-neg")
	must.NoError(err, "Goto failed")

	// Should NOT match "Delete" (which is not present)
	notExpected := "- button \"Delete\""
	err = playwright.Expect(page.Locator("body")).Not().ToMatchAriaSnapshot(ctx, notExpected)
	must.NoError(err, "Not().ToMatchAriaSnapshot failed")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotHeading(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<h1>title</h1>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "heading")
	is.Contains(snapshot, "title")

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx, `- heading "title" [level=1]`)
	must.NoError(err)
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotTextNodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<h1>Microsoft</h1><div>Open source projects and samples from Microsoft</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx,
		"- heading \"Microsoft\" [level=1]\n- text: Open source projects and samples from Microsoft")
	must.NoError(err)
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotMatchURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<a href="https://example.com">Link</a>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "link")
	is.Contains(snapshot, "example.com")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotNegated(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>OK</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).Not().ToMatchAriaSnapshot(ctx, `- heading "title" [level=1]`)
	must.NoError(err)
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotConcatenateSpanText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<span>One</span> <span>Two</span> <span>Three</span>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "One")
	is.Contains(snapshot, "Two")
	is.Contains(snapshot, "Three")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotConcatenateSpanText2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<span>One </span><span>Two </span><span>Three</span>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "One")
	is.Contains(snapshot, "Two")
	is.Contains(snapshot, "Three")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotConcatenateDivText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>One</div><div>Two</div><div>Three</div>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "One")
	is.Contains(snapshot, "Two")
	is.Contains(snapshot, "Three")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotMultilineText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<p>\n      Line 1\n      Line 2\n      Line 3\n    </p>")
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "Line 1")
	is.Contains(snapshot, "Line 2")
	is.Contains(snapshot, "Line 3")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotAriaHiddenText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<p><span>hello</span><span aria-hidden="true">world</span></p>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "hello")
	must.NotContains(snapshot, "world")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotPresentationAndNoneRoles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul><li role="presentation">hello</li><li role="none">world</li></ul>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "hello")
	is.Contains(snapshot, "world")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotCheckboxNotUseOnValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox"><input type="radio">`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	must.NotContains(snapshot, `"on"`)
	is.Contains(snapshot, "checkbox")
	is.Contains(snapshot, "radio")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotTextareaShowsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<textarea>Before</textarea>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "Before")

	_, err = page.Evaluate(ctx, `() => document.querySelector('textarea').value = 'After'`)
	must.NoError(err)

	snapshot2, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot2, "After")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotHiddenElementNoChildren(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div style="visibility: hidden;"><div style="visibility: visible;"><button>Button</button></div></div>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	must.NotContains(snapshot, "button")
	must.NotContains(snapshot, "Button")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotAriaHiddenNoChildren(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div aria-hidden="true"><div aria-hidden="false"><button>Button</button></div></div>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	must.NotContains(snapshot, "button")
	must.NotContains(snapshot, "Button")
}

// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestAriaSnapshotSlotsUseAssigned(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button><div>foo</div></button>
<script>(() => {
  const container = document.querySelector('div');
  const shadow = container.attachShadow({ mode: 'open' });
  const slot = document.createElement('slot');
  shadow.appendChild(slot);
})()</script>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "foo")
	is.Contains(snapshot, "button")
}

// TestLocatorAriaSnapshotCheckbox verifies AriaSnapshot includes checkbox.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestLocatorAriaSnapshotCheckbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" aria-label="Accept terms">
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotList  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotList(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<h1>title</h1><h1>title 2</h1>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx,
		"- heading \"title\" [level=1]\n- heading \"title 2\" [level=1]")
	must.NoError(err)
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotListWithAccessibleName  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotListWithAccessibleName(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul aria-label="my list"><li>one</li><li>two</li></ul>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx,
		"- list \"my list\":\n  - listitem: one\n  - listitem: two")
	must.NoError(err)
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotComplex  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotComplex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul><li><a href="about:blank">link</a></li></ul>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx,
		"- list:\n  - listitem:\n    - link \"link\"")
	must.NoError(err)
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotDetailsVisibility  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotDetailsVisibility(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<details><summary>Summary</summary><div>Details</div></details>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx, "- group: Summary")
	must.NoError(err)
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotIntegration  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotIntegration(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<h1>Microsoft</h1>
<div>Open source projects and samples from Microsoft</div>
<ul>
<li><details><summary>Verified</summary><div><div>
<p>We've verified that the organization <strong>microsoft</strong> controls the domain:</p>
<ul><li class="mb-1"><strong>opensource.microsoft.com</strong></li></ul>
<div><a href="about:blank">Learn more about verified organizations</a></div>
</div></div></details></li>
<li><a href="about:blank"><summary title="Label: GitHub Sponsor">Sponsor</summary></a></li>
</ul>`)
	must.NoError(err)

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "Microsoft")
	is.Contains(snapshot, "Open source projects and samples from Microsoft")
	is.Contains(snapshot, "Verified")
	is.Contains(snapshot, "Sponsor")
}

// TestPageAriaSnapshotShouldHandleTopLevelDeepEqual verifies that a wrong deep-equal pattern at the
// top level causes ToMatchAriaSnapshot to return an error.
// Ref: TestPageAriaSnapshot.java#shouldHandleTopLevelDeepEqual  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldHandleTopLevelDeepEqual(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul>
  <li>
    <ul>
      <li>1.1</li>
      <li>1.2</li>
    </ul>
  </li>
</ul>`)
	must.NoError(err)

	// The actual structure is list > listitem > list > listitem, but the pattern
	// expects list > listitem > listitem (flat). Deep-equal comparison must fail.
	pattern := "- /children: deep-equal\n- list:\n  - listitem:\n    - listitem: \"1.1\"\n    - listitem: \"1.2\""
	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx, pattern)
	is.Error(err, "expected ToMatchAriaSnapshot to fail for wrong deep-equal pattern")
}

// TestPageAriaSnapshotMatchValuesBothAgainstRegexAndString verifies that URL values containing
// regex-special characters are matched as literal strings in aria snapshots.
// Ref: TestPageAriaSnapshot.java#matchValuesBothAgainstRegexAndString  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotMatchValuesBothAgainstRegexAndString(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<a href="/auth?r=/">Log in</a>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx,
		"- link \"Log in\":\n  - /url: /auth?r=/")
	must.NoError(err)
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshot(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<h1>title</h1>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- heading \"title\" [level=1]")
}

// Ref: TestPageAriaSnapshot.java#shouldAllowTextNodes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldAllowTextNodes(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<h1>Microsoft</h1><div>Open source projects and samples from Microsoft</div>")
	a3CheckAndMatchSnapshot(ctx, t, body, ""+
		"    - heading \"Microsoft\" [level=1]\n"+
		"    - text: Open source projects and samples from Microsoft")
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotChildren  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotChildren(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<ul><li><img />One</li><li>Two</li><li>Three</li></ul>")
	must.NoError(playwright.Expect(body).ToMatchAriaSnapshot(ctx, "- list:\n  - /children: equal\n  - listitem\n  - listitem: Two\n  - listitem: Three"))
	must.NoError(playwright.Expect(body).Not().ToMatchAriaSnapshot(ctx, "- list:\n  - /children: equal\n  - listitem\n  - listitem: Two"))

	must.NoError(playwright.Expect(body).ToMatchAriaSnapshot(ctx, "- list:\n  - /children: deep-equal\n  - listitem:\n    - img\n    - text: One\n  - listitem: Two\n  - listitem: Three"))
	must.NoError(playwright.Expect(body).Not().ToMatchAriaSnapshot(ctx, "- list:\n  - /children: deep-equal\n  - listitem:\n    - text: One\n  - listitem: Two\n  - listitem: Three"))
	must.NoError(playwright.Expect(body).ToMatchAriaSnapshot(ctx, "- list:\n  - /children: deep-equal\n  - listitem:\n    - /children: contain\n    - text: One\n  - listitem: Two\n  - listitem: Three"))
}

// Ref: TestPageAriaSnapshot.java#shouldMatchUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldMatchUrl(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<a href='https://example.com'>Link</a>")
	require.NoError(t, playwright.Expect(body).ToMatchAriaSnapshot(ctx, ""+
		"- link:\n"+
		"  - /url: /.*example.com/"))
}

// Ref: TestPageAriaSnapshot.java#shouldSupportMultilineText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSupportMultilineText(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<p>\n      Line 1\n      Line 2\n      Line 3\n    </p>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- paragraph: Line 1 Line 2 Line 3")
	require.NoError(t, playwright.Expect(body).ToMatchAriaSnapshot(ctx, ""+
		"    - paragraph: |\n"+
		"        Line 1\n"+
		"        Line 2\n"+
		"        Line 3"))
}

// Ref: TestPageAriaSnapshot.java#shouldConcatenateSpanText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldConcatenateSpanText(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<span>One</span> <span>Two</span> <span>Three</span>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- text: One Two Three")
}

// Ref: TestPageAriaSnapshot.java#shouldConcatenateSpanText2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldConcatenateSpanText2(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<span>One </span><span>Two </span><span>Three</span>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- text: One Two Three")
}

// Ref: TestPageAriaSnapshot.java#shouldConcatenateDivTextWithSpaces  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldConcatenateDivTextWithSpaces(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<div>One</div><div>Two</div><div>Three</div>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- text: One Two Three")
}

// Ref: TestPageAriaSnapshot.java#shouldIncludePseudoInText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldIncludePseudoInText(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<style>span:before { content: 'world'; } div:after { content: 'bye'; }</style>"+
		"<a href=\"about:blank\"><span>hello</span><div>hello</div></a>")
	a3CheckAndMatchSnapshot(ctx, t, body, ""+
		"- link \"worldhello hellobye\":\n"+
		"  - /url: about:blank")
}

// Ref: TestPageAriaSnapshot.java#shouldNotIncludeHiddenPseudoInText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldNotIncludeHiddenPseudoInText(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<style>span:before { content: 'world'; display: none; } div:after { content: 'bye'; visibility: hidden; }</style>"+
		"<a href=\"about:blank\"><span>hello</span><div>hello</div></a>")
	a3CheckAndMatchSnapshot(ctx, t, body, ""+
		"- link \"hello hello\":\n"+
		"  - /url: about:blank")
}

// Ref: TestPageAriaSnapshot.java#shouldIncludeNewLineForBlockPseudo  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldIncludeNewLineForBlockPseudo(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<style>span:before { content: 'world'; display: block; } div:after { content: 'bye'; display: block; }</style>"+
		"<a href=\"about:blank\"><span>hello</span><div>hello</div></a>")
	a3CheckAndMatchSnapshot(ctx, t, body, ""+
		"- link \"world hello hello bye\":\n"+
		"  - /url: about:blank")
}

// Ref: TestPageAriaSnapshot.java#shouldWorkWithSlots  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldWorkWithSlots(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page, body := a3AriaBody(ctx, t, "<button><div>foo</div></button>"+
		"<script>(() => {"+
		"  const container = document.querySelector('div');"+
		"  const shadow = container.attachShadow({ mode: 'open' });"+
		"  const slot = document.createElement('slot');"+
		"  shadow.appendChild(slot);"+
		"})()</script>")
	// Text "foo" is assigned to the slot, should not be used twice.
	a3CheckAndMatchSnapshot(ctx, t, body, "- button \"foo\"")

	// Text "foo" is assigned to the slot, should be used instead of slot content.
	must.NoError(page.SetContent(ctx, "<div>foo</div>"+
		"<script>(() => {"+
		"  const container = document.querySelector('div');"+
		"  const shadow = container.attachShadow({ mode: 'open' });"+
		"  const button = document.createElement('button');"+
		"  shadow.appendChild(button);"+
		"  const slot = document.createElement('slot');"+
		"  button.appendChild(slot);"+
		"  const span = document.createElement('span');"+
		"  span.textContent = 'pre';"+
		"  slot.appendChild(span);"+
		"})()</script>"))
	a3CheckAndMatchSnapshot(ctx, t, body, "- button \"foo\"")

	// Nothing is assigned to the slot, should use slot content.
	must.NoError(page.SetContent(ctx, "<div></div>"+
		"<script>(() => {"+
		"  const container = document.querySelector('div');"+
		"  const shadow = container.attachShadow({ mode: 'open' });"+
		"  const button = document.createElement('button');"+
		"  shadow.appendChild(button);"+
		"  const slot = document.createElement('slot');"+
		"  button.appendChild(slot);"+
		"  const span = document.createElement('span');"+
		"  span.textContent = 'pre';"+
		"  slot.appendChild(span);"+
		"})()</script>"))
	a3CheckAndMatchSnapshot(ctx, t, body, "- button \"pre\"")
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotInnerText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotInnerText(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<div role=\"listitem\"><div><div><span title=\"a.test.ts\">a.test.ts</span></div>"+
		"<div><button title=\"Run\"></button><button title=\"Show source\"></button><button title=\"Watch\"></button></div></div></div>"+
		"<div role=\"listitem\"><div><div><span title=\"snapshot\">snapshot</span></div>"+
		"<div class=\"ui-mode-list-item-time\">30ms</div>"+
		"<div><button title=\"Run\"></button><button title=\"Show source\"></button><button title=\"Watch\"></button></div></div></div>")
	a3CheckAndMatchSnapshot(ctx, t, body, ""+
		"    - listitem:\n"+
		"      - text: a.test.ts\n"+
		"      - button \"Run\"\n"+
		"      - button \"Show source\"\n"+
		"      - button \"Watch\"\n"+
		"    - listitem:\n"+
		"      - text: snapshot 30ms\n"+
		"      - button \"Run\"\n"+
		"      - button \"Show source\"\n"+
		"      - button \"Watch\"")
}

// Ref: TestPageAriaSnapshot.java#checkAriaHiddenText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotCheckAriaHiddenText(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<p><span>hello</span><span aria-hidden=\"true\">world</span></p>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- paragraph: hello")
}

// Ref: TestPageAriaSnapshot.java#shouldIgnorePresentationAndNoneRoles  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldIgnorePresentationAndNoneRoles(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<ul><li role=\"presentation\">hello</li><li role=\"none\">world</li></ul>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- list: hello world")
}

// Ref: TestPageAriaSnapshot.java#shouldNotUseOnAsCheckboxValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldNotUseOnAsCheckboxValue(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<input type=\"checkbox\"><input type=\"radio\">")
	a3CheckAndMatchSnapshot(ctx, t, body, "- checkbox\n- radio")
}

// Ref: TestPageAriaSnapshot.java#shouldNotReportTextareaTextContent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldNotReportTextareaTextContent(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	page, body := a3AriaBody(ctx, t, "<textarea>Before</textarea>")
	a3CheckAndMatchSnapshot(ctx, t, body, "- textbox: Before")
	_, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value = 'After'")
	require.NoError(t, err)
	a3CheckAndMatchSnapshot(ctx, t, body, "- textbox: After")
}

// Ref: TestPageAriaSnapshot.java#shouldNotShowVisibleChildrenOfHiddenElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldNotShowVisibleChildrenOfHiddenElements(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<div style=\"visibility: hidden;\">"+
		"<div style=\"visibility: visible;\"><button>Button</button></div>"+
		"</div>")
	got, err := body.AriaSnapshot(ctx)
	require.NoError(t, err)
	assert.Equal(t, "", got)
}

// Ref: TestPageAriaSnapshot.java#shouldNotShowUnhiddenChildrenOfAriaHiddenElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldNotShowUnhiddenChildrenOfAriaHiddenElements(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	_, body := a3AriaBody(ctx, t, "<div aria-hidden=\"true\">"+
		"<div aria-hidden=\"false\"><button>Button</button></div>"+
		"</div>")
	got, err := body.AriaSnapshot(ctx)
	require.NoError(t, err)
	assert.Equal(t, "", got)
}

// Ref: TestPageAriaSnapshot.java#shouldSnapshotPlaceholderWhenDifferentFromName  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotShouldSnapshotPlaceholderWhenDifferentFromName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page, body := a3AriaBody(ctx, t, "<input placeholder=\"Placeholder\">")
	must.NoError(playwright.Expect(body).ToMatchAriaSnapshot(ctx, "- textbox \"Placeholder\""))

	must.NoError(page.SetContent(ctx, "<input placeholder=\"Placeholder\" aria-label=\"Label\">"))
	must.NoError(playwright.Expect(body).ToMatchAriaSnapshot(ctx, ""+
		"- textbox \"Label\":\n"+
		"  - /placeholder: Placeholder"))
}

// Ref: TestPageAriaSnapshot.java#pageMatchesAriaSnapshot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshot.java
func TestPageAriaSnapshotPageMatchesAriaSnapshot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<h1>hello</h1>"))
	must.NoError(playwright.ExpectPage(page).ToMatchAriaSnapshot(ctx, `- heading "hello" [level=1]`))

	timeout := 1000.0
	err := playwright.ExpectPage(page).ToMatchAriaSnapshot(ctx, `- heading "world"`, &playwright.MatchesAriaSnapshotOptions{
		Timeout: &timeout,
	})
	must.Error(err)
	is.Contains(err.Error(), "to.match.aria")
}
