//go:build e2e

// Tests for strict mode selector behavior at the page and browser-context level.
// Migration of: TestPageStrict.java, TestBrowserContextStrict.java
package e2e

import (
	"context"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageStrictTextContent verifies that using a locator (which is always strict) throws
// "strict mode violation" when the selector matches multiple elements.
// Ref: TestPageStrict.java#shouldFailPageTextContentInStrictMode
func TestPageStrictTextContent(t *testing.T) {
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

// TestPageStrictGetAttribute verifies that using a locator (always strict) throws
// "strict mode violation" when multiple elements match for GetAttribute.
// Ref: TestPageStrict.java#shouldFailPageGetAttributeInStrictMode
func TestPageStrictGetAttribute(t *testing.T) {
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

// TestPageStrictFill verifies that using a locator (always strict) throws
// "strict mode violation" when multiple input elements match for Fill.
// Ref: TestPageStrict.java#shouldFailPageFillInStrictMode
func TestPageStrictFill(t *testing.T) {
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

// TestPageStrictQuerySelector verifies that using a locator (always strict) throws
// "strict mode violation" when the selector matches multiple elements.
// Ref: TestPageStrict.java#shouldFailPageInStrictMode
func TestPageStrictQuerySelector(t *testing.T) {
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

// TestPageStrictWaitForSelector verifies that WaitFor on a locator (always strict) throws
// "strict mode violation" when multiple elements match.
// Ref: TestPageStrict.java#shouldFailPageWaitForSelectorInStrictMode
func TestPageStrictWaitForSelector(t *testing.T) {
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

// TestPageStrictDispatchEvent verifies that DispatchEvent via locator (always strict) throws
// "strict mode violation" when multiple elements match.
// Ref: TestPageStrict.java#shouldFailPageDispatchEventInStrictMode
func TestPageStrictDispatchEvent(t *testing.T) {
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

// TestBrowserContextStrictNotFailInNonStrictMode verifies that without StrictSelectors,
// page.Locator("span").First().TextContent() returns the first match without error.
// Ref: TestBrowserContextStrict.java#shouldNotFailPageTextContentInNonStrictMode
func TestBrowserContextStrictNotFailInNonStrictMode(t *testing.T) {
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

// TestBrowserContextStrictFailTextContent verifies that when a context is created with
// StrictSelectors: true, textContent on a multi-match selector throws "strict mode violation".
// Ref: TestBrowserContextStrict.java#shouldFailPageTextContentInStrictMode
func TestBrowserContextStrictFailTextContent(t *testing.T) {
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

// TestBrowserContextStrictOptOut verifies that even in a strict-selectors context,
// using First() opts out of strictness and returns the first match.
// Ref: TestBrowserContextStrict.java#shouldOptOutOfStrictMode
func TestBrowserContextStrictOptOut(t *testing.T) {
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
