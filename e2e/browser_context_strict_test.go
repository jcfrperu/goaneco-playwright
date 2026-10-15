//go:build e2e

// E2E tests for: TestBrowserContextStrict.java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestBrowserContextStrictShouldNotFailPageTextContentInNonStrictMode verifies that without StrictSelectors,
// page.Locator("span").First().TextContent() returns the first match without error.
// Ref: TestBrowserContextStrict.java#shouldNotFailPageTextContentInNonStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStrict.java
func TestBrowserContextStrictShouldNotFailPageTextContentInNonStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	must.NoError(page.SetContent(ctx, "<span>span1</span><div><span>target</span></div>"))

	// Without context-level strict selectors, First() resolves the ambiguity.
	text, err := page.Locator("span").First().TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("span1", *text)
}

// TestBrowserContextStrictShouldFailPageTextContentInStrictMode verifies that when a context is created with
// StrictSelectors: true, textContent on a multi-match selector throws "strict mode violation".
// Ref: TestBrowserContextStrict.java#shouldFailPageTextContentInStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStrict.java
func TestBrowserContextStrictShouldFailPageTextContentInStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	trueVal := true
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StrictSelectors: &trueVal,
	})
	must.NoError(err, "NewContext with StrictSelectors failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	must.NoError(page.SetContent(ctx, "<span>span1</span><div><span>target</span></div>"))

	_, err = page.Locator("span").TextContent(ctx)
	is.Error(err)
	is.Contains(err.Error(), "strict mode violation")
}

// TestBrowserContextStrictShouldOptOutOfStrictMode verifies that even in a strict-selectors context,
// using First() opts out of strictness and returns the first match.
// Ref: TestBrowserContextStrict.java#shouldOptOutOfStrictMode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStrict.java
func TestBrowserContextStrictShouldOptOutOfStrictMode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	trueVal := true
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StrictSelectors: &trueVal,
	})
	must.NoError(err, "NewContext with StrictSelectors failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	must.NoError(page.SetContent(ctx, "<span>span1</span><div><span>target</span></div>"))

	// Nth(0) / First() narrows to one element, opting out of strict-mode violation.
	text, err := page.Locator("span").First().TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("span1", *text)
}
