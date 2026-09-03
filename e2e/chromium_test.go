//go:build e2e

// E2E tests for Chromium-specific features (CDP connections, tracing over CDP).
// Migration of: TestChromium.java
//
// These tests only run when the active browser is Chromium (PLAYWRIGHT_BROWSER=chromium).
package e2e

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// freePort finds an available TCP port by briefly listening on :0.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "freePort: listen failed")
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

// launchChromiumWithCDP launches a new Chromium browser exposing the Chrome DevTools Protocol
// on the given port. The browser is closed automatically via t.Cleanup.
func launchChromiumWithCDP(t *testing.T, port int) *playwright.Browser {
	t.Helper()
	must := require.New(t)

	launchCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	bt := globalBrowser.BrowserType()
	browser, err := bt.Launch(launchCtx, &playwright.BrowserTypeLaunchOptions{
		Args: []string{fmt.Sprintf("--remote-debugging-port=%d", port)},
	})
	must.NoError(err, "Launch with CDP port failed")

	t.Cleanup(func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer closeCancel()
		_ = browser.Close(closeCtx)
	})

	return browser
}

// TestChromiumConnectOverCDP verifies connecting to an existing Chromium session via CDP.
// Ref: TestChromium.java#shouldConnectToAnExistingCdpSession
func TestChromiumConnectOverCDP(t *testing.T) {
	if globalBTName != "chromium" {
		t.Skip("chromium-specific test")
	}
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	port := freePort(t)
	_ = launchChromiumWithCDP(t, port)

	ctx := testCtx(t)
	bt := globalBrowser.BrowserType()

	cdpBrowser, err := bt.ConnectOverCDP(ctx, fmt.Sprintf("http://localhost:%d", port))
	must.NoError(err, "ConnectOverCDP failed")
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = cdpBrowser.Close(closeCtx)
	}()

	contexts := cdpBrowser.Contexts()
	is.GreaterOrEqual(len(contexts), 1, "expected at least one browser context via CDP")
}

// TestChromiumConnectOverCDPTwice verifies two simultaneous CDP connections to the same browser.
// Ref: TestChromium.java#shouldConnectToAnExistingCdpSessionTwice
func TestChromiumConnectOverCDPTwice(t *testing.T) {
	if globalBTName != "chromium" {
		t.Skip("chromium-specific test")
	}
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	port := freePort(t)
	_ = launchChromiumWithCDP(t, port)

	ctx := testCtx(t)
	bt := globalBrowser.BrowserType()
	endpoint := fmt.Sprintf("http://localhost:%d", port)

	cdp1, err := bt.ConnectOverCDP(ctx, endpoint)
	must.NoError(err, "ConnectOverCDP (1) failed")
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = cdp1.Close(closeCtx)
	}()

	cdp2, err := bt.ConnectOverCDP(ctx, endpoint)
	must.NoError(err, "ConnectOverCDP (2) failed")
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = cdp2.Close(closeCtx)
	}()

	is.GreaterOrEqual(len(cdp1.Contexts()), 1, "cdp1: expected at least one context")
	is.GreaterOrEqual(len(cdp2.Contexts()), 1, "cdp2: expected at least one context")
}

// TestChromiumTracingOverCDP verifies that tracing works on a browser connected via CDP.
// Ref: TestChromium.java#shouldSupportTracingOverCDP
func TestChromiumTracingOverCDP(t *testing.T) {
	if globalBTName != "chromium" {
		t.Skip("chromium-specific test")
	}
	t.Parallel()
	must := require.New(t)

	port := freePort(t)
	_ = launchChromiumWithCDP(t, port)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bt := globalBrowser.BrowserType()

	cdpBrowser, err := bt.ConnectOverCDP(ctx, fmt.Sprintf("http://localhost:%d", port))
	must.NoError(err, "ConnectOverCDP failed")
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = cdpBrowser.Close(closeCtx)
	}()

	contexts := cdpBrowser.Contexts()
	must.NotEmpty(contexts, "expected at least one browser context")

	bCtx := contexts[0]
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage via CDP context failed")

	traceCtx, traceCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer traceCancel()

	must.NoError(bCtx.Tracing().Start(traceCtx, &playwright.TracingStartOptions{
		Screenshots: boolPtr(true),
		Snapshots:   boolPtr(true),
	}), "Tracing.Start failed")

	must.NoError(page.Goto(ctx, srv.EmptyPage()), "Goto failed")
	must.NoError(page.SetContent(ctx, `<button>Click</button>`))
	must.NoError(page.Locator("button").Click(ctx))

	must.NoError(bCtx.Tracing().Stop(traceCtx), "Tracing.Stop failed")
}
