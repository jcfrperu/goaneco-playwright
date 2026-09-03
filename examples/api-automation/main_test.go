//go:build e2e

// Petstore API Automation Examples.
// Demonstrates APIRequestContext usage with goaneco-playwright.
//
// Prerequisites:
//   - PLAYWRIGHT_CLI_PATH environment variable pointing to playwright-core/cli.js
//   - Internet access to reach https://petstore.swagger.io/v2
//
// Run all scenarios:
//
//	go test -tags e2e -v -timeout 300s ./examples/api-automation/...
package apiautomation

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

const baseURL = "https://petstore.swagger.io/v2"

var globalPW *playwright.Playwright

func TestMain(m *testing.M) {
	exitCode := run(m)
	if exitCode == 0 {
		if err := goleak.Find(); err != nil {
			fmt.Fprintf(os.Stderr, "goleak: %v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(exitCode)
}

func run(m *testing.M) int {
	cliPath := resolvePlaywrightCLI()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pw, err := playwright.Run(ctx, cliPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "api-automation: failed to start Playwright: %v\n", err)
		return 1
	}
	globalPW = pw
	defer func() { _ = pw.Stop() }()

	return m.Run()
}

// newAPICtx creates a new APIRequestContext with Petstore as the base URL
// and registers cleanup on test teardown.
func newAPICtx(t *testing.T) *playwright.APIRequestContext {
	t.Helper()
	ctx := testCtx(t)
	base := baseURL
	apiCtx, err := globalPW.NewAPIRequestContext(ctx, &playwright.APIRequestContextOptions{
		BaseURL: &base,
	})
	require.NoError(t, err, "NewAPIRequestContext failed")
	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = apiCtx.Dispose(cleanCtx)
	})
	return apiCtx
}

// testCtx returns a 30-second context suitable for API interactions.
func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	fmt.Fprintln(os.Stderr, "api-automation: PLAYWRIGHT_CLI_PATH not set and no Playwright CLI found.")
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
