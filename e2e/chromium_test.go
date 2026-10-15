//go:build e2e

// E2E tests for: TestChromium.java.
package e2e

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"testing"
	"time"
)

// TestChromiumShouldConnectToAnExistingCdpSession verifies connecting to an existing Chromium session via CDP.
// Ref: TestChromium.java#shouldConnectToAnExistingCdpSession  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestChromium.java
func TestChromiumShouldConnectToAnExistingCdpSession(t *testing.T) {
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

// TestChromiumShouldConnectToAnExistingCdpSessionTwice verifies two simultaneous CDP connections to the same browser.
// Ref: TestChromium.java#shouldConnectToAnExistingCdpSessionTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestChromium.java
func TestChromiumShouldConnectToAnExistingCdpSessionTwice(t *testing.T) {
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

// TestChromiumShouldSupportTracingOverCDP verifies that tracing works on a browser connected via CDP.
// Ref: TestChromium.java#shouldSupportTracingOverCDP  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestChromium.java
func TestChromiumShouldSupportTracingOverCDP(t *testing.T) {
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

// TestChromiumShouldConnectOverAWsEndpoint verifies connecting over CDP using a ws:// endpoint taken from /json/version.
// Ref: TestChromium.java#shouldConnectOverAWsEndpoint  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestChromium.java
func TestChromiumShouldConnectOverAWsEndpoint(t *testing.T) {
	if globalBTName != "chromium" {
		t.Skip("chromium-specific test")
	}
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	port := freePort(t)
	_ = launchChromiumWithCDP(t, port)
	ctx := testCtx(t)

	var wsEndpoint string
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/json/version/", port))
		if err == nil {
			var v struct {
				WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&v)
			_ = resp.Body.Close()
			wsEndpoint = v.WebSocketDebuggerURL
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	must.NotEmpty(wsEndpoint, "could not resolve ws endpoint")

	bt := globalBrowser.BrowserType()
	for i := 0; i < 2; i++ {
		cdpBrowser, err := bt.ConnectOverCDP(ctx, wsEndpoint)
		must.NoError(err)
		is.Len(cdpBrowser.Contexts(), 1)
		must.NoError(cdpBrowser.Close(ctx))
	}
}

// TestChromiumShouldSendExtraHeadersWithConnectRequest verifies custom headers are sent in the CDP websocket handshake.
// Ref: TestChromium.java#shouldSendExtraHeadersWithConnectRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestChromium.java
func TestChromiumShouldSendExtraHeadersWithConnectRequest(t *testing.T) {
	if globalBTName != "chromium" {
		t.Skip("chromium-specific test")
	}
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	must.NoError(err)
	defer ln.Close()

	reqCh := make(chan *http.Request, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		req, err := http.ReadRequest(bufio.NewReader(conn))
		if err == nil {
			reqCh <- req
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = globalBrowser.BrowserType().ConnectOverCDP(ctx, fmt.Sprintf("ws://localhost:%d/ws", ln.Addr().(*net.TCPAddr).Port),
		&playwright.ConnectOverCDPOptions{Headers: map[string]string{"User-Agent": "Playwright", "foo": "bar"}})

	select {
	case req := <-reqCh:
		is.Equal("Playwright", req.Header.Get("User-Agent"))
		is.Equal("bar", req.Header.Get("foo"))
	case <-time.After(5 * time.Second):
		t.Fatal("no handshake received")
	}
}
