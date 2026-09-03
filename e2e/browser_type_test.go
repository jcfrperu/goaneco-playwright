//go:build e2e

// Tests for BrowserType (Priority 1 - Smoke Tests).
// Migration of: TestBrowserTypeBasic.java
package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrowserTypeExecutablePathShouldWork verifies that the browser binary exists on disk.
// Ref: TestBrowserTypeBasic.java#browserTypeExecutablePathShouldWork
func TestBrowserTypeExecutablePathShouldWork(t *testing.T) {
	t.Parallel()
	// Skip if custom channel is configured (equivalent to Assumptions.assumeTrue in Java)
	if os.Getenv("PLAYWRIGHT_BROWSER_CHANNEL") != "" {
		t.Skip("skipping: PLAYWRIGHT_BROWSER_CHANNEL is set")
	}

	bt := globalBrowser.BrowserType()
	execPath := bt.ExecutablePath()

	if execPath == "" {
		t.Fatal("BrowserType.ExecutablePath() returned empty string")
	}

	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		t.Fatalf("browser executable not found at path: %q", execPath)
	}

	t.Logf("browser executable: %s", filepath.ToSlash(execPath))
}

// TestBrowserTypeNameShouldWork verifies that the browser type name matches
// the browser selected via PLAYWRIGHT_BROWSER.
// Ref: TestBrowserTypeBasic.java#browserTypeNameShouldWork
func TestBrowserTypeNameShouldWork(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	bt := globalBrowser.BrowserType()
	name := bt.Name()

	if name == "" {
		t.Fatal("BrowserType.Name() returned empty string")
	}

	expectedName := globalBTName // configured in TestMain
	is.Equal(expectedName, name)

	t.Logf("browser type name: %s", name)
}

// TestBrowserTypeNameIsKnown verifies that the returned name is one of the supported
// Playwright browser types.
func TestBrowserTypeNameIsKnown(t *testing.T) {
	t.Parallel()
	bt := globalBrowser.BrowserType()
	name := bt.Name()

	knownBrowsers := map[string]bool{
		"chromium": true,
		"firefox":  true,
		"webkit":   true,
	}

	if !knownBrowsers[name] {
		t.Errorf("BrowserType.Name() = %q is not a known Playwright browser", name)
	}
}

func TestBrowserTypeExecutablePathNotEmptyExtra(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	bt := globalBrowser.BrowserType()
	is.NotEmpty(bt.ExecutablePath())
}

func TestBrowserTypeNameIsChromiumExtra(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	bt := globalBrowser.BrowserType()
	name := bt.Name()
	is.NotEmpty(name)
	is.Equal("chromium", name)
}

func TestBrowserLaunchAndClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	bt := globalBrowser.BrowserType()
	browser, err := bt.Launch(ctx, nil)
	must.NoError(err)
	must.NotNil(browser)

	must.NoError(browser.Close(ctx))
}

func TestBrowserNewContextCreatesContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	must.NotNil(bc)

	must.NoError(bc.Close(ctx))
}

func TestBrowserNewPageCreatesPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	page, err := globalBrowser.NewPage(ctx)
	must.NoError(err)
	must.NotNil(page)
	is.False(page.IsClosed())

	must.NoError(page.Close(ctx))
}
