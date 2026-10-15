//go:build e2e

// E2E tests for: TestLocatorAnyFrame.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestFrameChildFrameLocatorClickEx verifies click inside iframe.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAnyFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestFrameChildFrameLocatorClickEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	must.NoError(frame.SetContent(ctx, `
		<button id="btn" onclick="this.dataset.clicked='yes'">Click me</button>
	`))

	must.NoError(frame.Locator("#btn").Click(ctx))

	attr, err := frame.Locator("#btn").GetAttribute(ctx, "data-clicked")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("yes", *attr)
}

// TestPageFrameLocatorLocator verifies that FrameLocator.Locator() finds elements inside an iframe.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAnyFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestPageFrameLocatorLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Serve the inner iframe content
	srv.ServeWithBody("/inner-frame", "text/html", `<div id="greeting">Hello from iframe</div>`)

	// Serve the outer page with an iframe
	srv.ServeWithBody("/frame-host", "text/html",
		`<iframe id="my-frame" src="`+srv.Prefix()+`/inner-frame"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-host")
	must.NoError(err, "Goto failed")

	// Locate the element inside the iframe
	fl := page.FrameLocator("#my-frame")
	greeting := fl.Locator("#greeting")

	text, err := greeting.InnerText(ctx)
	must.NoError(err, "FrameLocator.Locator().InnerText() failed")
	if text != "Hello from iframe" {
		t.Errorf("expected 'Hello from iframe', got %q", text)
	}
}

// TestFrameLocatorFrameLocator verifies nested FrameLocator for iframes within iframes.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAnyFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestFrameLocatorFrameLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/fl-nested-inner", "text/html", `<p id="deep">deep content</p>`)
	srv.ServeWithBody("/fl-nested-mid", "text/html",
		`<iframe id="inner-frame" src="`+srv.Prefix()+`/fl-nested-inner"></iframe>`)
	srv.ServeWithBody("/fl-nested-outer", "text/html",
		`<iframe id="mid-frame" src="`+srv.Prefix()+`/fl-nested-mid"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/fl-nested-outer")
	must.NoError(err, "Goto failed")

	// page → #mid-frame → #inner-frame → #deep
	text, err := page.FrameLocator("#mid-frame").FrameLocator("#inner-frame").Locator("#deep").InnerText(ctx)
	must.NoError(err, "nested FrameLocator.InnerText() failed")
	if text != "deep content" {
		t.Errorf("nested FrameLocator text = %q, want 'deep content'", text)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorAnyFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestFrameLocatorScopedCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul><li>a</li><li>b</li><li>c</li></ul>
	`))

	count, err := page.MainFrame().Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorAnyFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestFrameLocatorClickWorks(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button onclick="document.getElementById('r').textContent='ok'">go</button>
		<div id="r"></div>
	`))

	must.NoError(page.MainFrame().Locator("button").Click(ctx))

	text, err := page.MainFrame().Locator("#r").InnerText(ctx)
	must.NoError(err)
	is.Equal("ok", text)
}

// TestLocatorCountMatchesMultipleEx2 verifies Count returns number of matched elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorAnyFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorCountMatchesMultipleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li><li>Item 4</li></ul>
	`))

	count, err := page.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(4, count)
}

// Ref: TestLocatorAnyFrame.java#shouldClickAButtonInsideAnIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldClickAButtonInsideAnIframe(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldClickAButtonInTheMainFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldClickAButtonInTheMainFrame(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldFailClickWhenElementsMatchInMultipleFrames  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldFailClickWhenElementsMatchInMultipleFrames(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldFailClickUponStrictModeViolationInsideASingleFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldFailClickUponStrictModeViolationInsideASingleFrame(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldTimeOutOnClickWhenThereAreNoMatches  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldTimeOutOnClickWhenThereAreNoMatches(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldCountElementsInASingleFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldCountElementsInASingleFrame(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldFailCountWhenElementsMatchInMultipleFrames  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldFailCountWhenElementsMatchInMultipleFrames(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportHasCount  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportHasCount(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldWaitForAFrameToAppearWithHasCount  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldWaitForAFrameToAppearWithHasCount(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldFailHasCountWhenElementsMatchInMultipleFrames  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldFailHasCountWhenElementsMatchInMultipleFrames(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportHasText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportHasText(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportEvaluate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportEvaluate(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportFirstLastNth  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportFirstLastNth(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldNotAllowNthFrameOnAnyFrameLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldNotAllowNthFrameOnAnyFrameLocator(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportTwoFrameLocators  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportTwoFrameLocators(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportLocatorBeforeFrameLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportLocatorBeforeFrameLocator(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportOwnerOfAFrameLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportOwnerOfAFrameLocator(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldSupportContentFrameAfterAnyFrameLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldSupportContentFrameAfterAnyFrameLocator(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}

// Ref: TestLocatorAnyFrame.java#shouldWorkOnFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorAnyFrame.java
func TestLocatorAnyFrameShouldWorkOnFrame(t *testing.T) {
	t.Skip("not implemented: Page.FrameLocator() without selector (any-frame locator)")
}
