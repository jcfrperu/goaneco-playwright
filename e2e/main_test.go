//go:build e2e

// Package e2e contains end-to-end integration tests for goaneco-playwright.
//
// Running tests in this package requires:
// 1. An accessible Playwright Node.js driver (PLAYWRIGHT_CLI_PATH environment variable).
// 2. Installed browsers (chromium by default, selectable via PLAYWRIGHT_BROWSER).
// 3. Explicit execution with the "e2e" build tag:
//
//	go test -tags=e2e ./e2e/... -v -timeout=120s
//
// Fixture architecture:
//   - TestMain starts a shared Playwright instance for the entire package.
//   - Each test receives a Browser created in TestMain.
//   - Each test creates its own BrowserContext + Page using newPage() or newContext().
//   - The HTTP test server is created per-test (using testserver.New).
//
// Equivalent to TestBase.java + @BeforeAll/@AfterAll structure in the Java reference project.
package e2e

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
	"go.uber.org/goleak"
)

// Global fixtures shared across all tests in the package.
// Initialized in TestMain, accessed in individual tests (read-only).
var (
	globalPW      *playwright.Playwright
	globalBrowser *playwright.Browser
	globalBTName  string // selected browser name (chromium/firefox/webkit)
)

// TestMain configures the global E2E environment: launches Playwright, starts the browser,
// runs all tests, verifies absence of goroutine leaks via goleak, and shuts down cleanly.
func TestMain(m *testing.M) {
	exitCode := run(m)
	// goleak is checked only when all tests pass (exitCode == 0) to prevent
	// false positives caused by lingering goroutines from failed assertions or timeouts.
	if exitCode == 0 {
		if err := goleak.Find(); err != nil {
			fmt.Fprintf(os.Stderr, "goleak: %v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(exitCode)
}

// run encapsulates setup/teardown logic in a separate function so os.Exit does not bypass defers.
func run(m *testing.M) (exitCode int) {
	cliPath := resolvePlaywrightCLI()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pw, err := playwright.Run(ctx, cliPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2e: failed to start Playwright: %v\n", err)
		return 1
	}
	globalPW = pw
	defer func() {
		if stopErr := pw.Stop(); stopErr != nil {
			fmt.Fprintf(os.Stderr, "e2e: failed to stop Playwright: %v\n", stopErr)
		}
	}()

	bt, err := resolveBrowserType(pw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2e: failed to resolve browser type: %v\n", err)
		return 1
	}
	globalBTName = bt.Name()

	launchCtx, launchCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer launchCancel()

	opts := &playwright.BrowserTypeLaunchOptions{}
	browser, err := bt.Launch(launchCtx, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2e: failed to launch browser %q: %v\n", globalBTName, err)
		return 1
	}
	globalBrowser = browser
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer closeCancel()
		if closeErr := browser.Close(closeCtx); closeErr != nil {
			fmt.Fprintf(os.Stderr, "e2e: failed to close browser: %v\n", closeErr)
		}
	}()

	return m.Run()
}

// resolvePlaywrightCLI returns the filesystem path to the Playwright Node.js CLI.
// Lookup precedence:
//  1. PLAYWRIGHT_CLI_PATH environment variable
//  2. Well-known global npm install paths
//  3. Fatal error if not found
func resolvePlaywrightCLI() string {
	if path := os.Getenv("PLAYWRIGHT_CLI_PATH"); path != "" {
		return path
	}

	// Global npm install paths (Windows / macOS / Linux)
	candidates := npmGlobalCandidates()
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	fmt.Fprintf(os.Stderr, "e2e: PLAYWRIGHT_CLI_PATH not set and no Playwright CLI found in standard paths.\n")
	fmt.Fprintf(os.Stderr, "     Set PLAYWRIGHT_CLI_PATH to the path of playwright-core/cli.js\n")
	os.Exit(1)
	return ""
}

// npmGlobalCandidates returns candidate paths for the Playwright CLI based on OS platform.
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
			"/usr/local/lib/node_modules/playwright-core/cli.js",
			filepath.Join(home, ".npm-global", "lib", "node_modules", "playwright-core", "cli.js"),
		}
	default: // linux
		return []string{
			"/usr/lib/node_modules/playwright/node_modules/playwright-core/cli.js",
			"/usr/local/lib/node_modules/playwright-core/cli.js",
			"/usr/local/lib/node_modules/playwright/node_modules/playwright-core/cli.js",
		}
	}
}

// resolveBrowserType resolves the BrowserType instance to use in tests based on PLAYWRIGHT_BROWSER.
// Supported values: "chromium", "firefox", "webkit". Defaults to "chromium".
func resolveBrowserType(pw *playwright.Playwright) (*playwright.BrowserType, error) {
	browserName := os.Getenv("PLAYWRIGHT_BROWSER")
	if browserName == "" {
		browserName = "chromium"
	}

	switch browserName {
	case "chromium":
		return pw.Chromium()
	case "firefox":
		return pw.Firefox()
	case "webkit":
		return pw.WebKit()
	default:
		return nil, fmt.Errorf("unsupported browser %q (use: chromium, firefox, webkit)", browserName)
	}
}

// newPage is a convenience helper that creates a new page in the global browser
// and registers automatic cleanup via t.Cleanup.
func newPage(t *testing.T) *playwright.Page {
	t.Helper()
	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	require.NoError(t, err, "newPage: NewContext failed")
	page, err := bCtx.NewPage(ctx)
	if err != nil {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
		require.NoError(t, err, "newPage: NewPage failed")
	}
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	return page
}

// newContext creates a new BrowserContext in the global browser and registers cleanup.
func newContext(t *testing.T) *playwright.BrowserContext {
	t.Helper()
	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	require.NoError(t, err, "newContext: NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	return bCtx
}

// testCtx creates a standard 30s timeout context.Context for test operations.
func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}
