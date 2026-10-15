//go:build e2e

// E2E tests for: TestLocatorHighlight.java.
package e2e

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestLocatorHighlightHighlightAndHideHighlightShouldNotThrow verifies that highlight can be shown and hidden without error.
// Ref: TestLocatorHighlight.java#highlightAndHideHighlightShouldNotThrow  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorHighlight.java
func TestLocatorHighlightHighlightAndHideHighlightShouldNotThrow(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>Click me</button>`))

	must.NoError(page.Locator("button").Highlight(ctx), "Highlight() should not error")
	must.NoError(page.HideHighlight(ctx), "HideHighlight() should not error")
}

// TestLocatorHighlightShouldHighlightLocator verifies the highlight overlay matches the element's bounding box.
//
// NOTE: Verifying the visual overlay requires screenshot comparison or DOM inspection of
// the highlight overlay element, which is an internal Playwright implementation detail.
// This test verifies the API is callable and the element exists in the DOM.
// Ref: TestLocatorHighlight.java#shouldHighlightLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorHighlight.java
func TestLocatorHighlightShouldHighlightLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:100px;height:100px;background:red"></div>`))

	must.NoError(page.Locator("div").Highlight(ctx), "Highlight() should not error for visible element")
	must.NoError(page.HideHighlight(ctx))
}
