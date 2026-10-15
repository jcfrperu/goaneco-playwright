//go:build e2e

// E2E tests for: TestPageStrict.java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPageStrictShouldFailPageTextContentInStrictMode verifies that using a locator (which is always strict) throws
// "strict mode violation" when the selector matches multiple elements.
// Ref: TestPageStrict.java#shouldFailPageTextContentInStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageStrict.java
func TestPageStrictShouldFailPageTextContentInStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	ctx := testCtx(t)

	must.NoError(page.SetContent(ctx, "<span>span1</span><div><span>target</span></div>"))

	_, err := page.Locator("span").TextContent(ctx)
	is.Error(err)
	is.Contains(err.Error(), "strict mode violation")
}

// TestPageStrictShouldFailPageGetAttributeInStrictMode verifies that using a locator (always strict) throws
// "strict mode violation" when multiple elements match for GetAttribute.
// Ref: TestPageStrict.java#shouldFailPageGetAttributeInStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageStrict.java
func TestPageStrictShouldFailPageGetAttributeInStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	ctx := testCtx(t)

	must.NoError(page.SetContent(ctx, "<span>span1</span><div><span>target</span></div>"))

	_, err := page.Locator("span").GetAttribute(ctx, "id")
	is.Error(err)
	is.Contains(err.Error(), "strict mode violation")
}

// TestPageStrictShouldFailPageFillInStrictMode verifies that using a locator (always strict) throws
// "strict mode violation" when multiple input elements match for Fill.
// Ref: TestPageStrict.java#shouldFailPageFillInStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageStrict.java
func TestPageStrictShouldFailPageFillInStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	ctx := testCtx(t)

	must.NoError(page.SetContent(ctx, "<input></input><div><input></input></div>"))

	err := page.Locator("input").Fill(ctx, "text")
	is.Error(err)
	is.Contains(err.Error(), "strict mode violation")
}

// TestPageStrictShouldFailPageInStrictMode verifies that using a locator (always strict) throws
// "strict mode violation" when the selector matches multiple elements.
// Ref: TestPageStrict.java#shouldFailPageInStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageStrict.java
func TestPageStrictShouldFailPageInStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	ctx := testCtx(t)

	must.NoError(page.SetContent(ctx, "<span>span1</span><div><span>target</span></div>"))

	// Locator.IsVisible uses strict mode internally; any action on a multi-match locator errors.
	_, err := page.Locator("span").IsVisible(ctx)
	is.Error(err)
	is.Contains(err.Error(), "strict mode violation")
}

// TestPageStrictShouldFailPageWaitForSelectorInStrictMode verifies that WaitFor on a locator (always strict) throws
// "strict mode violation" when multiple elements match.
// Ref: TestPageStrict.java#shouldFailPageWaitForSelectorInStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageStrict.java
func TestPageStrictShouldFailPageWaitForSelectorInStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	ctx := testCtx(t)

	must.NoError(page.SetContent(ctx, "<span>span1</span><div><span>target</span></div>"))

	err := page.Locator("span").WaitFor(ctx)
	is.Error(err)
	is.Contains(err.Error(), "strict mode violation")
}

// TestPageStrictShouldFailPageDispatchEventInStrictMode verifies that DispatchEvent via locator (always strict) throws
// "strict mode violation" when multiple elements match.
// Ref: TestPageStrict.java#shouldFailPageDispatchEventInStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageStrict.java
func TestPageStrictShouldFailPageDispatchEventInStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	ctx := testCtx(t)

	must.NoError(page.SetContent(ctx, "<span></span><div><span></span></div>"))

	err := page.Locator("span").Click(ctx)
	is.Error(err)
	is.Contains(err.Error(), "strict mode violation")
}
