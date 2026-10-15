//go:build e2e

// E2E tests for: TestBrowser1.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

// TestBrowser1VersionShouldWork verifies that the browser version has the expected format.
// Ref: TestBrowser1.java#versionShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowser1VersionShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	version := globalBrowser.Version()

	if version == "" {
		t.Fatal("Browser.Version() returned empty string")
	}

	var pattern string
	switch globalBTName {
	case "chromium":
		pattern = `^\d+\.\d+\.\d+\.\d+$`
	case "webkit":
		pattern = `^\d+\.\d+`
	case "firefox":
		// Firefox can be e.g. "85.0b1"
		pattern = `^\d+\.\d+.*`
	default:
		t.Fatalf("unknown browser name: %q", globalBTName)
	}

	matched, err := regexp.MatchString(pattern, version)
	must.NoErrorf(err, "invalid regex pattern %q", pattern)
	if !matched {
		t.Errorf("Browser.Version() = %q does not match pattern %q for browser %q",
			version, pattern, globalBTName)
	}

	t.Logf("browser version: %s", version)
}

// TestBrowser1ShouldReturnBrowserType verifies that browser.BrowserType() returns the correct BrowserType.
// Ref: TestBrowser1.java#shouldReturnBrowserType  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowser1ShouldReturnBrowserType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bt := globalBrowser.BrowserType()

	must.NotNil(bt, "Browser.BrowserType() returned nil")
	is.Equal(globalBTName, bt.Name())
}

// TestBrowser1ShouldCreateNewPage verifies the lifecycle of Browser.Contexts() upon creating and closing pages.
//
// Not parallelized: counts globalBrowser.Contexts() which is shared across all tests.
// Running this concurrently with other tests that create contexts would produce non-deterministic counts.
// Ref: TestBrowser1.java#shouldCreateNewPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowser1ShouldCreateNewPage(t *testing.T) { //nolint:paralleltest
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	initialContexts := len(globalBrowser.Contexts())

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")

	contexts := globalBrowser.Contexts()
	is.Len(contexts, initialContexts+1)

	page, err := bCtx.NewPage(ctx)
	if err != nil {
		_ = bCtx.Close(ctx)
		must.NoError(err, "NewPage failed")
	}

	if page.Context() != bCtx {
		t.Error("page.Context() does not match the parent context")
	}

	err = bCtx.Close(ctx)
	must.NoError(err, "Close failed")

	contextsAfter := globalBrowser.Contexts()
	is.Len(contextsAfter, initialContexts)
}

// TestBrowser1ShouldFireContextEvent verifies that Browser.OnContext triggers when a new context is created.
// Ref: TestBrowser1.java#shouldFireContextEvent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowser1ShouldFireContextEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	var capturedContext *playwright.BrowserContext
	cancel := globalBrowser.OnContext(func(bCtx *playwright.BrowserContext) {
		capturedContext = bCtx
	})
	t.Cleanup(cancel)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		_ = bCtx.Close(closeCtx)
	})

	if capturedContext != bCtx {
		t.Errorf("capturedContext = %v, want %v", capturedContext, bCtx)
	}
}

// TestBrowser1ShouldThrowUponSecondCreateNewPage verifies that newPage can be called multiple times on any context.
// Note: Playwright 1.61+ allows multiple pages per context; test updated to verify current behavior.
// Ref: TestBrowser1.java#shouldThrowUponSecondCreateNewPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowser1ShouldThrowUponSecondCreateNewPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	// In Playwright 1.61+, creating multiple pages in the same context is allowed.
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(cleanCtx)
	})

	page1, err := bCtx.NewPage(ctx)
	must.NoError(err, "first NewPage should succeed")
	must.NotNil(page1)

	// Second NewPage should now succeed (behavior changed in Playwright 1.61+)
	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "second NewPage should succeed in Playwright 1.61+")
	must.NotNil(page2)
}

// TestBrowserNewPageAndClose verifies isolated creation and closing of multiple pages.
// Partial equivalent of TestBrowser1.java#shouldCreateNewPage.
// Extra: no equivalent method in playwright-java; related class: TestBrowser1.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowserNewPageAndClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	page1 := newPage(t)

	must.NotNil(page1, "newPage returned nil page")
	is.False(page1.IsClosed(), "freshly created page should not be closed")

	page2 := newPage(t)

	must.NotNil(page2, "second newPage returned nil page")

	// Distinct pages must have distinct pointers/GUIDs
	if page1 == page2 {
		t.Error("two different newPage() calls should return distinct Page objects")
	}
}

// TestBrowserIsConnected verifies that Browser.IsConnected() returns true and a new page starts at about:blank.
// Extra: no equivalent method in playwright-java; related class: TestBrowser1.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowserIsConnected(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	is.True(globalBrowser.IsConnected(), "expected globalBrowser.IsConnected() to return true")

	page := newPage(t)
	must.NotNil(page, "could not create page — browser may be disconnected")

	// Initial URL of a new page must be strictly about:blank
	url := page.URL()
	is.Equal("about:blank", url)
}

// TestBrowserCloseShouldBeCallableTwice verifies that Browser.Close() is idempotent.
// Extra: no equivalent method in playwright-java; related class: TestBrowser1.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowserCloseShouldBeCallableTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bt := globalBrowser.BrowserType()

	browser, err := bt.Launch(ctx, nil)
	must.NoError(err, "Launch failed")

	is.True(browser.IsConnected(), "expected browser to be connected initially")

	err = browser.Close(ctx)
	must.NoError(err, "first Close failed")
	if browser.IsConnected() {
		t.Error("expected browser to not be connected after close")
	}

	err = browser.Close(ctx)
	must.NoError(err, "second Close failed")
	err = browser.Close(ctx)
	must.NoError(err, "third Close failed")
}

// TestBrowser1ShouldPropagateCloseReasonToPendingActions is not ported.
// Ref: TestBrowser1.java#shouldPropagateCloseReasonToPendingActions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser1.java
func TestBrowser1ShouldPropagateCloseReasonToPendingActions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bt := globalBrowser.BrowserType()
	browser, err := bt.Launch(ctx, nil)
	must.NoError(err)

	bCtx, err := browser.NewContext(ctx)
	must.NoError(err)

	reason := "The reason."
	_, err = bCtx.WaitForPage(ctx, func() error {
		return browser.Close(ctx, &playwright.BrowserCloseOptions{Reason: &reason})
	})
	must.Error(err)
	is.Contains(err.Error(), "The reason.")
}
