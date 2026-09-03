//go:build e2e

// SauceDemo UI Automation Examples.
// Demonstrates the Page Object Model (POM) pattern with goaneco-playwright.
//
// Prerequisites:
//   - PLAYWRIGHT_CLI_PATH environment variable pointing to playwright-core/cli.js
//   - Internet access to reach https://www.saucedemo.com
//
// Run all scenarios:
//
//	go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
package uiautomation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
)

var (
	globalPW      *playwright.Playwright
	globalBrowser *playwright.Browser
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	cliPath := resolvePlaywrightCLI()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pw, err := playwright.Run(ctx, cliPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ui-automation: failed to start Playwright: %v\n", err)
		return 1
	}
	globalPW = pw
	defer func() { _ = pw.Stop() }()

	bt, err := pw.Chromium()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ui-automation: chromium not available: %v\n", err)
		return 1
	}

	headless := true
	browser, err := bt.Launch(ctx, &playwright.BrowserTypeLaunchOptions{Headless: &headless})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ui-automation: failed to launch browser: %v\n", err)
		return 1
	}
	globalBrowser = browser
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = browser.Close(closeCtx)
	}()

	return m.Run()
}

// newPage creates a new browser context and page, registering cleanup on test teardown.
func newPage(t *testing.T) *playwright.Page {
	t.Helper()
	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	require.NoError(t, err, "NewContext failed")
	page, err := bCtx.NewPage(ctx)
	if err != nil {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
		require.NoError(t, err, "NewPage failed")
	}
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	return page
}

// testCtx returns a 90-second context suitable for remote site interactions.
func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func resolvePlaywrightCLI() string {
	if path := os.Getenv("PLAYWRIGHT_CLI_PATH"); path != "" {
		return path
	}
	for _, p := range npmGlobalCandidates() {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	fmt.Fprintln(os.Stderr, "ui-automation: PLAYWRIGHT_CLI_PATH not set and no Playwright CLI found.")
	os.Exit(1)
	return ""
}

func npmGlobalCandidates() []string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		return []string{
			filepath.Join(appData, "npm", "node_modules", "playwright", "node_modules", "playwright-core", "cli.js"),
			filepath.Join(appData, "npm", "node_modules", "playwright-core", "cli.js"),
		}
	case "darwin":
		home := os.Getenv("HOME")
		return []string{
			"/usr/local/lib/node_modules/playwright/node_modules/playwright-core/cli.js",
			filepath.Join(home, ".npm-global", "lib", "node_modules", "playwright-core", "cli.js"),
		}
	default:
		return []string{
			"/usr/local/lib/node_modules/playwright/node_modules/playwright-core/cli.js",
			"/usr/lib/node_modules/playwright/node_modules/playwright-core/cli.js",
		}
	}
}
