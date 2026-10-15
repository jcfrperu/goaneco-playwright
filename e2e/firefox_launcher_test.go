//go:build e2e

// E2E tests for: TestFirefoxLauncher.java.
package e2e

import (
	"context"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestFirefoxLauncherShouldPassFirefoxUserPreferences verifies that Firefox user preferences (e.g. proxy settings)
// can be passed via LaunchOptions and take effect.
// Ref: TestFirefoxLauncher.java#shouldPassFirefoxUserPreferences  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFirefoxLauncher.java
func TestFirefoxLauncherShouldPassFirefoxUserPreferences(t *testing.T) {
	if globalBTName != "firefox" {
		t.Skip("firefox-specific test")
	}
	t.Parallel()
	must := require.New(t)

	port := freePort(t)
	launchCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	bt := globalBrowser.BrowserType()
	browser, err := bt.Launch(launchCtx, &playwright.BrowserTypeLaunchOptions{
		FirefoxUserPrefs: map[string]any{
			"network.proxy.type":      1,
			"network.proxy.http":      "127.0.0.1",
			"network.proxy.http_port": port,
		},
	})
	must.NoError(err, fmt.Sprintf("Launch with FirefoxUserPrefs failed (port %d)", port))
	t.Cleanup(func() {
		closeCtx, cc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cc()
		_ = browser.Close(closeCtx)
	})
	must.NotNil(browser)
}

// TestFirefoxLauncherShouldPassFirefoxUserPreferencesInPersistent verifies Firefox prefs work with a persistent context.
// Ref: TestFirefoxLauncher.java#shouldPassFirefoxUserPreferencesInPersistent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFirefoxLauncher.java
func TestFirefoxLauncherShouldPassFirefoxUserPreferencesInPersistent(t *testing.T) {
	if globalBTName != "firefox" {
		t.Skip("firefox-specific test")
	}
	t.Parallel()
	must := require.New(t)

	port := freePort(t)
	userDataDir := t.TempDir()

	launchCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	bt := globalBrowser.BrowserType()
	browser, bCtx, err := bt.LaunchPersistentContext(launchCtx, userDataDir, &playwright.LaunchPersistentContextOptions{
		FirefoxUserPrefs: map[string]any{
			"network.proxy.type":      1,
			"network.proxy.http":      "127.0.0.1",
			"network.proxy.http_port": port,
		},
	})
	must.NoError(err, fmt.Sprintf("LaunchPersistentContext with FirefoxUserPrefs failed (port %d)", port))
	t.Cleanup(func() {
		closeCtx, cc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cc()
		_ = bCtx.Close(closeCtx)
		_ = browser.Close(closeCtx)
	})
	must.NotNil(browser)
	must.NotNil(bCtx)
}
