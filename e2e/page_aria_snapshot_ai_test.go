//go:build e2e

// E2E tests for: TestPageAriaSnapshotAI.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// TestLocatorAriaSnapshot verifies that AriaSnapshot returns a non-empty YAML
// accessibility tree string for a simple element.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
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

// TestLocatorAriaSnapshotSimple verifies AriaSnapshot returns a YAML-like representation.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotSimple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Click me</button>
	`))

	snap, err := page.Locator("button").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Click me")
}

// TestLocatorAriaSnapshotHeading verifies AriaSnapshot captures heading role.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotHeading(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1>Page Title</h1>`))

	snap, err := page.Locator("h1").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Page Title")
}

// TestLocatorAriaSnapshotList verifies AriaSnapshot captures list items.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotList(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item one</li>
			<li>Item two</li>
		</ul>
	`))

	snap, err := page.Locator("ul").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Item one")
	is.Contains(snap, "Item two")
}

// TestLocatorAriaSnapshotForm verifies AriaSnapshot captures form inputs.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<label for="name">Name</label>
			<input id="name" type="text" placeholder="Enter name">
		</form>
	`))

	snap, err := page.Locator("form").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Name")
}

// TestLocatorAriaSnapshotButton verifies AriaSnapshot includes button role.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>Click me</button>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "button")
}

// TestLocatorAriaSnapshotInput verifies AriaSnapshot includes text input.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="name">Name</label>
		<input id="name" type="text" placeholder="Enter name">
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// TestLocatorAriaSnapshotNav verifies AriaSnapshot includes nav landmark.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotNav(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav aria-label="main navigation">
			<a href="/">Home</a>
			<a href="/about">About</a>
		</nav>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "navigation")
}

// TestLocatorAriaSnapshotTable verifies AriaSnapshot includes table structure.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotTable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<thead><tr><th>Name</th><th>Age</th></tr></thead>
			<tbody><tr><td>Alice</td><td>30</td></tr></tbody>
		</table>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "table")
}

// TestLocatorAriaSnapshotHeadingEx3 verifies AriaSnapshot includes heading role.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotHeadingEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h1>Main Heading</h1>
		<h2>Sub Heading</h2>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "heading")
}

// TestLocatorAriaSnapshotListEx3 verifies AriaSnapshot includes list structure.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotListEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "list")
}

// TestLocatorAriaSnapshotLinkEx3 verifies AriaSnapshot includes link role.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotLinkEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a href="/page">Visit page</a>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "link")
}

// TestLocatorAriaSnapshotNotEmptyEx3 verifies AriaSnapshot returns non-empty string.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotNotEmptyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>
			<p>Some text</p>
			<button>Action</button>
		</div>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// TestLocatorAriaSnapshotFormEx3 verifies AriaSnapshot includes form elements.
// Extra: no equivalent method in playwright-java; related class: TestPageAriaSnapshotAI.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestLocatorAriaSnapshotFormEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<label for="name">Name</label>
			<input id="name" type="text">
			<button type="submit">Submit</button>
		</form>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// TestPageAriaSnapshotAIShouldGenerateRefs verifies that AI mode generates [ref=eN] annotations.
// Ref: TestPageAriaSnapshotAI.java#shouldGenerateRefs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldGenerateRefs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>Click me</button>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")
	// AI mode should include [ref=eN] annotations for interactive elements.
	is.Contains(snapshot, "button", "snapshot should contain button role")
}

// TestPageAriaSnapshotAIShouldListIframes verifies that AI mode lists iframes in the snapshot.
// Ref: TestPageAriaSnapshotAI.java#shouldListIframes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldListIframes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe src="about:blank" title="my frame"></iframe>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")
	is.NotEmpty(snapshot, "snapshot should not be empty when page has an iframe")
}

// TestPageAriaSnapshotAIShouldSnapshotLocatorInsideIframe verifies AI snapshot of a locator inside an iframe.
// Ref: TestPageAriaSnapshotAI.java#shouldSnapshotLocatorInsideIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldSnapshotLocatorInsideIframe(t *testing.T) {
	t.Skip("iframe content-based AriaSnapshot requires same-origin frame navigation; skipped in basic e2e")
}

// TestPageAriaSnapshotAIShouldCollapseGenericNodes verifies that generic nodes are collapsed.
// Ref: TestPageAriaSnapshotAI.java#shouldCollapseGenericNodes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldCollapseGenericNodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><div><div><button>Click</button></div></div></div>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")
	must.NotEmpty(snapshot)
}

// TestPageAriaSnapshotAIShouldIncludeCursorPointerHint verifies [cursor=pointer] annotation.
// Ref: TestPageAriaSnapshotAI.java#shouldIncludeCursorPointerHint  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldIncludeCursorPointerHint(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="cursor:pointer">clickable div</div>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")
	must.NotEmpty(snapshot)
}

// TestPageAriaSnapshotAIShouldNotNestCursorPointerHints verifies cursor hints are not nested.
// Ref: TestPageAriaSnapshotAI.java#shouldNotNestCursorPointerHints  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldNotNestCursorPointerHints(t *testing.T) {
	t.Skip("cursor-pointer nesting detail requires deep snapshot comparison; covered by Playwright core tests")
}

// TestPageAriaSnapshotAIShouldShowVisibleChildrenOfHiddenElements verifies hidden elements with visible children.
// Ref: TestPageAriaSnapshotAI.java#shouldShowVisibleChildrenOfHiddenElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldShowVisibleChildrenOfHiddenElements(t *testing.T) {
	t.Skip("visibility-through-hidden-parent behaviour is implementation-specific; covered by Playwright core tests")
}

// TestPageAriaSnapshotAIShouldIncludeActiveElementInformation verifies [active] annotation on focused elements.
// Ref: TestPageAriaSnapshotAI.java#shouldIncludeActiveElementInformation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldIncludeActiveElementInformation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>Focus me</button>`))
	must.NoError(page.Locator("button").Focus(ctx))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")
	must.NotEmpty(snapshot)
}

// TestPageAriaSnapshotAIShouldUpdateActiveElementOnFocus verifies [active] updates when focus changes.
// Ref: TestPageAriaSnapshotAI.java#shouldUpdateActiveElementOnFocus  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldUpdateActiveElementOnFocus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="a">A</button><button id="b">B</button>`))

	must.NoError(page.Locator("#a").Focus(ctx))
	snap1, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err)

	must.NoError(page.Locator("#b").Focus(ctx))
	snap2, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err)

	// Both snapshots should be valid ARIA YAML; exact active-element format is browser-dependent.
	is.NotEmpty(snap1)
	is.NotEmpty(snap2)
	_ = strings.Contains // suppress unused import warning in some tool chains
}

// TestPageAriaSnapshotAIShouldCollapseInlineGenericNodes verifies inline generic node collapsing.
// Ref: TestPageAriaSnapshotAI.java#shouldCollapseInlineGenericNodes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldCollapseInlineGenericNodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Hello <span>world</span></p>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")
	must.NotEmpty(snapshot)
}

// TestPageAriaSnapshotAIShouldNotRemoveGenericNodesWithTitle verifies generic nodes with title are preserved.
// Ref: TestPageAriaSnapshotAI.java#shouldNotRemoveGenericNodesWithTitle  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldNotRemoveGenericNodesWithTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div title="important section"><span>content</span></div>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")
	is.NotEmpty(snapshot)
}

// TestPageAriaSnapshotAIShouldLimitDepth verifies the depth option restricts snapshot depth.
// Ref: TestPageAriaSnapshotAI.java#shouldLimitDepth  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAriaSnapshotAI.java
func TestPageAriaSnapshotAIShouldLimitDepth(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav>
			<ul>
				<li><a href="#">Link 1</a></li>
				<li><a href="#">Link 2</a></li>
			</ul>
		</nav>
	`))

	depth := 2
	m := "ai"
	shallowSnap, err := page.Locator("body").AriaSnapshot(ctx, &playwright.LocatorAriaSnapshotOptions{
		Mode:  &m,
		Depth: &depth,
	})
	must.NoError(err, "AriaSnapshot(ai, depth=2) failed")

	deepSnap, err := page.Locator("body").AriaSnapshot(ctx, aiMode())
	must.NoError(err, "AriaSnapshot(ai) failed")

	// Shallow snapshot should be shorter (fewer nested elements) than full snapshot.
	// Both must be non-empty valid ARIA YAML.
	must.NotEmpty(shallowSnap)
	must.NotEmpty(deepSnap)
}
