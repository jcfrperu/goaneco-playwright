//go:build e2e

// E2E tests for: TestBrowserContextProxy.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	"encoding/base64"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"sync"
	"testing"
	"time"
)

// TestBrowserContextOptionsProxy verifies that proxy configuration in BrowserContextOptions
// routes page HTTP traffic through the specified proxy server, while unconfigured contexts bypass it.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextProxy.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextProxy.java
func TestBrowserContextOptionsProxy(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	targetServer := testserver.New(t)

	var proxiedRequests int
	var proxyMu sync.Mutex

	// Start a local HTTP proxy server that forwards requests to target server
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyMu.Lock()
		proxiedRequests++
		proxyMu.Unlock()

		targetURL, _ := url.Parse(targetServer.Prefix())
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	}))
	t.Cleanup(proxyServer.Close)

	// 1. Context with proxy: traffic must go through proxy
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Proxy: &playwright.ProxyOptions{
			Server: proxyServer.URL,
		},
	})
	must.NoError(err, "NewContext with Proxy failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, targetServer.Prefix()+"/title.html")
	must.NoError(err, "Goto failed")

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Woof-Woof", title)

	proxyMu.Lock()
	countAfterProxyContext := proxiedRequests
	proxyMu.Unlock()

	if countAfterProxyContext == 0 {
		t.Error("expected at least one request to pass through the proxy server for proxied context")
	}

	// 2. Context without proxy: traffic must NOT go through proxyServer (OBS-N02 isolation check)
	noProxyCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext without Proxy failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = noProxyCtx.Close(closeCtx)
	})

	noProxyPage, err := noProxyCtx.NewPage(ctx)
	must.NoError(err, "noProxyPage NewPage failed")

	err = noProxyPage.Goto(ctx, targetServer.Prefix()+"/title.html")
	must.NoError(err, "noProxyPage Goto failed")

	proxyMu.Lock()
	countAfterNoProxyContext := proxiedRequests
	proxyMu.Unlock()

	if countAfterNoProxyContext != countAfterProxyContext {
		t.Errorf("no-proxy context routed %d requests through proxy server, want 0",
			countAfterNoProxyContext-countAfterProxyContext)
	}
}

// Ref: TestBrowserContextProxy.java#shouldUseProxy  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextProxy.java
func TestBrowserContextProxyShouldUseProxy(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	a1ServeProxyTarget(srv)
	bCtx := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Proxy: &playwright.ProxyOptions{Server: fmt.Sprintf("localhost:%d", srv.Port())},
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "http://non-existent.com/target.html"))
	title, err := page.Title(ctx)
	must.NoError(err)
	must.Equal("Served by the proxy", title)
}

// Ref: TestBrowserContextProxy.java#shouldUseProxyTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextProxy.java
func TestBrowserContextProxyShouldUseProxyTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	a1ServeProxyTarget(srv)
	bCtx := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Proxy: &playwright.ProxyOptions{Server: fmt.Sprintf("localhost:%d", srv.Port())},
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "http://non-existent.com/target.html"))
	must.NoError(page.Goto(ctx, "http://non-existent-2.com/target.html"))
	title, err := page.Title(ctx)
	must.NoError(err)
	must.Equal("Served by the proxy", title)
}

// Ref: TestBrowserContextProxy.java#shouldUseProxyForSecondPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextProxy.java
func TestBrowserContextProxyShouldUseProxyForSecondPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	a1ServeProxyTarget(srv)
	bCtx := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Proxy: &playwright.ProxyOptions{Server: fmt.Sprintf("localhost:%d", srv.Port())},
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "http://non-existent.com/target.html"))
	title, err := page.Title(ctx)
	must.NoError(err)
	must.Equal("Served by the proxy", title)

	page2, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Goto(ctx, "http://non-existent.com/target.html"))
	title2, err := page2.Title(ctx)
	must.NoError(err)
	must.Equal("Served by the proxy", title2)
}

// Ref: TestBrowserContextProxy.java#shouldWorkWithIPPORTNotion  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextProxy.java
func TestBrowserContextProxyShouldWorkWithIPPORTNotion(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	a1ServeProxyTarget(srv)
	bCtx := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Proxy: &playwright.ProxyOptions{Server: fmt.Sprintf("127.0.0.1:%d", srv.Port())},
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "http://non-existent.com/target.html"))
	title, err := page.Title(ctx)
	must.NoError(err)
	must.Equal("Served by the proxy", title)
}

// Ref: TestBrowserContextProxy.java#shouldAuthenticate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextProxy.java
func TestBrowserContextProxyShouldAuthenticate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/target.html", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Proxy-Authorization")
		if auth == "" {
			w.Header().Set("Proxy-Authenticate", "Basic realm='Access to internal site'")
			w.WriteHeader(407)
			return
		}
		_, _ = fmt.Fprintf(w, "<html><title>%s</title></html>", auth)
	})
	user, pass := "user", "secret"
	bCtx := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Proxy: &playwright.ProxyOptions{Server: fmt.Sprintf("localhost:%d", srv.Port()), Username: &user, Password: &pass},
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "http://non-existent.com/target.html"))
	title, err := page.Title(ctx)
	must.NoError(err)
	must.Equal("Basic "+base64.StdEncoding.EncodeToString([]byte("user:secret")), title)
}

// Ref: TestBrowserContextProxy.java#doesLaunchWithoutAPort  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextProxy.java
func TestBrowserContextProxyDoesLaunchWithoutAPort(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Proxy: &playwright.ProxyOptions{Server: "http://localhost"},
	})
	must.NoError(err)
	must.NoError(bCtx.Close(ctx))
}
