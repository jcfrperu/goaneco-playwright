//go:build e2e

// E2E tests for Locator.AriaSnapshot and Expect.ToMatchAriaSnapshot.
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

// TestLocatorAriaSnapshot verifies that AriaSnapshot returns a non-empty YAML
// accessibility tree string for a simple element.
func TestLocatorAriaSnapshot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/aria", "text/html", `<button>Click me</button>`)
	err := page.Goto(ctx, srv.Prefix()+"/aria")
	must.NoError(err, "Goto failed")

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err, "AriaSnapshot failed")
	if snapshot == "" {
		t.Fatal("expected non-empty snapshot")
	}
	t.Logf("aria snapshot: %s", snapshot)

	is.Contains(snapshot, "button")
}

// TestExpectToMatchAriaSnapshot verifies the ToMatchAriaSnapshot assertion.
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

func TestAriaSnapshotList(t *testing.T) {
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

func TestAriaSnapshotListWithAccessibleName(t *testing.T) {
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

func TestAriaSnapshotComplex(t *testing.T) {
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

func TestAriaSnapshotDetailsVisibility(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<details><summary>Summary</summary><div>Details</div></details>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("body")).ToMatchAriaSnapshot(ctx, "- group: Summary")
	must.NoError(err)
}

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

func TestAriaSnapshotIntegration(t *testing.T) {
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

// TestPageAriaSnapshotHandleTopLevelDeepEqual verifies that a wrong deep-equal pattern at the
// top level causes ToMatchAriaSnapshot to return an error.
// Ref: TestPageAriaSnapshot.java#shouldHandleTopLevelDeepEqual
func TestPageAriaSnapshotHandleTopLevelDeepEqual(t *testing.T) {
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

// TestPageAriaSnapshotMatchValuesBothRegexAndString verifies that URL values containing
// regex-special characters are matched as literal strings in aria snapshots.
// Ref: TestPageAriaSnapshot.java#matchValuesBothAgainstRegexAndString
func TestPageAriaSnapshotMatchValuesBothRegexAndString(t *testing.T) {
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
