//go:build e2e

// E2E tests for: TestBrowser2.java.
package e2e

import (
	"os"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
)

// TestBrowser2ShouldSupportDeprecatedChannelEnum verifies launching with a browser channel taken from the environment.
// Ref: TestBrowser2.java#shouldSupportDeprecatedChannelEnum  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowser2.java
func TestBrowser2ShouldSupportDeprecatedChannelEnum(t *testing.T) {
	channel := os.Getenv("BROWSER_CHANNEL")
	if channel == "" {
		t.Skip("BROWSER_CHANNEL not set (Java: Assumptions.assumeTrue(channel != null))")
	}
	t.Parallel()
	ctx := testCtx(t)
	bt := globalBrowser.BrowserType()
	browser, err := bt.Launch(ctx, &playwright.BrowserTypeLaunchOptions{Channel: &channel})
	require.NoError(t, err)
	require.NotNil(t, browser)
	require.NoError(t, browser.Close(ctx))
}
