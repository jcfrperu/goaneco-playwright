//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
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
