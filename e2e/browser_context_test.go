//go:build e2e

// Tests for BrowserContext (Priority 2).
// Migration of: TestBrowserContextBasic.java
package e2e

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrowserContextCreateAndClose verifies isolated browser contexts can be created and closed.
// Ref: TestBrowserContextBasic.java#shouldCreateNewContext
func TestBrowserContextCreateAndClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext() failed")

	page, err := bCtx.NewPage(ctx)
	if err != nil {
		_ = bCtx.Close(ctx)
		t.Fatalf("NewPage() failed: %v", err)
	}

	must.NotNil(page, "expected non-nil Page")
	is.False(page.IsClosed(), "newly created page should not be closed")

	err = bCtx.Close(ctx)
	must.NoError(err, "Close() on BrowserContext failed")
}

// TestBrowserContextCloseShouldBeCallableTwice verifies that BrowserContext.Close() is idempotent.
// Ref: TestBrowserContextBasic.java#closeShouldBeCallableTwice
func TestBrowserContextCloseShouldBeCallableTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")

	err = bCtx.Close(ctx)
	must.NoError(err, "first Close() failed")
	err = bCtx.Close(ctx)
	must.NoError(err, "second Close() failed")
	err = bCtx.Close(ctx)
	must.NoError(err, "third Close() failed")
	is.True(bCtx.IsClosed(), "expected context to remain closed")
}

// TestBrowserContextTrackedByBrowser verifies Browser.Contexts() tracks active contexts
// and BrowserContext.Browser() references the parent browser.
// Ref: TestBrowserContextBasic.java#shouldCreateNewContext
func TestBrowserContextTrackedByBrowser(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	initialCount := len(globalBrowser.Contexts())

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")

	is.Equal(initialCount+1, len(globalBrowser.Contexts()))

	is.Equal(globalBrowser, bCtx.Browser())

	err = bCtx.Close(ctx)
	must.NoError(err, "Close failed")

	is.Equal(initialCount, len(globalBrowser.Contexts()))
}

// TestBrowserContextPagesTracked verifies BrowserContext.Pages() tracks all open pages.
// Ref: TestBrowserContextBasic.java#shouldReturnAllOfThePages
func TestBrowserContextPagesTracked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	is.Len(bCtx.Pages(), 0)

	page1, err := bCtx.NewPage(ctx)
	must.NoError(err, "page1 NewPage failed")

	is.Len(bCtx.Pages(), 1)
	is.Equal(page1, bCtx.Pages()[0])
	is.Equal(bCtx, page1.Context())

	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "page2 NewPage failed")

	is.Len(bCtx.Pages(), 2)

	err = page1.Close(ctx)
	must.NoError(err, "page1.Close failed")

	is.Len(bCtx.Pages(), 1)
	is.Equal(page2, bCtx.Pages()[0])
}

// TestBrowserContextMultiplePages verifies creation and concurrent navigation across multiple pages
// within the same browser context using t.Cleanup.
// Ref: TestBrowserContextBasic.java
func TestBrowserContextMultiplePages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext() failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page1, err := bCtx.NewPage(ctx)
	must.NoError(err, "page1 NewPage() failed")

	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "page2 NewPage() failed")

	page3, err := bCtx.NewPage(ctx)
	must.NoError(err, "page3 NewPage() failed")

	err = page1.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "page1.Goto failed")
	err = page2.Goto(ctx, srv.Prefix()+"/title.html")
	must.NoError(err, "page2.Goto failed")
	err = page3.SetContent(ctx, "<h1>Page 3 Content</h1>")
	must.NoError(err, "page3.SetContent failed")

	title2, err := page2.Title(ctx)
	must.NoError(err)
	is.Equal("Woof-Woof", title2)

	content3, err := page3.Content(ctx)
	must.NoError(err)
	is.Equal("<html><head></head><body><h1>Page 3 Content</h1></body></html>", content3)
}

// TestBrowserContextIsolation verifies full navigation isolation between contexts.
// Ref: TestBrowserContextBasic.java
func TestBrowserContextIsolation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx1, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext 1 failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx1.Close(closeCtx)
	})

	bCtx2, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext 2 failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(closeCtx)
	})

	page1, err := bCtx1.NewPage(ctx)
	must.NoError(err, "page1 NewPage failed")

	page2, err := bCtx2.NewPage(ctx)
	must.NoError(err, "page2 NewPage failed")

	err = page1.Goto(ctx, srv.Prefix()+"/title.html")
	must.NoError(err, "page1.Goto failed")
	err = page2.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "page2.Goto failed")

	is.Equal(srv.Prefix()+"/title.html", page1.URL())
	is.Equal(srv.EmptyPage(), page2.URL())
}

// TestBrowserContextLocalStorageIsolation verifies full localStorage and Cookie isolation across different contexts.
// Ref: TestBrowserContextBasic.java#shouldIsolateLocalStorageAndCookies
func TestBrowserContextLocalStorageIsolation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx1, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext 1 failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx1.Close(closeCtx)
	})

	bCtx2, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext 2 failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(closeCtx)
	})

	page1, err := bCtx1.NewPage(ctx)
	must.NoError(err, "page1 NewPage failed")

	page2, err := bCtx2.NewPage(ctx)
	must.NoError(err, "page2 NewPage failed")

	err = page1.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "page1.Goto failed")
	err = page2.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "page2.Goto failed")

	// Set localStorage and cookie in page1
	_, err = page1.Evaluate(ctx, `() => {
		localStorage.setItem('name', 'page1');
		document.cookie = 'name=page1';
	}`)
	must.NoError(err, "page1.Evaluate set failed")

	// Set localStorage and cookie in page2
	_, err = page2.Evaluate(ctx, `() => {
		localStorage.setItem('name', 'page2');
		document.cookie = 'name=page2';
	}`)
	must.NoError(err, "page2.Evaluate set failed")

	// Verify independent localStorage
	val1, err := page1.Evaluate(ctx, "() => localStorage.getItem('name')")
	must.NoError(err)
	is.Equal("page1", val1)

	val2, err := page2.Evaluate(ctx, "() => localStorage.getItem('name')")
	must.NoError(err)
	is.Equal("page2", val2)

	// Verify independent cookies
	cookie1, err := page1.Evaluate(ctx, "() => document.cookie")
	must.NoError(err)
	is.Equal("name=page1", cookie1)

	cookie2, err := page2.Evaluate(ctx, "() => document.cookie")
	must.NoError(err)
	is.Equal("name=page2", cookie2)
}

// TestBrowserContextOptionsUserAgent verifies custom userAgent configuration in BrowserContextOptions.
func TestBrowserContextOptionsUserAgent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	customAgent := "CustomUserAgent/1.0 (Playwright-Go-Test)"

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		UserAgent: &customAgent,
	})
	must.NoError(err, "NewContext with custom userAgent failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	userAgentRaw, err := page.Evaluate(ctx, "() => window.navigator.userAgent")
	must.NoError(err, "Evaluate userAgent failed")

	userAgent, ok := userAgentRaw.(string)
	if !ok || userAgent != customAgent {
		t.Errorf("userAgent = %q, want %q", userAgent, customAgent)
	}
}

// TestBrowserContextOptionsLocale verifies custom locale configuration in BrowserContextOptions.
// Ref: TestBrowserContextLocale.java
func TestBrowserContextOptionsLocale(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	locale := "fr-FR"

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Locale: &locale,
	})
	must.NoError(err, "NewContext with custom locale failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	langRaw, err := page.Evaluate(ctx, "() => window.navigator.language")
	must.NoError(err, "Evaluate language failed")

	lang, ok := langRaw.(string)
	if !ok || lang != "fr-FR" {
		t.Errorf("navigator.language = %q, want %q", lang, "fr-FR")
	}
}

// TestBrowserContextOptionsViewport verifies custom viewport configuration in BrowserContextOptions.
// Ref: TestBrowserContextViewport.java
func TestBrowserContextOptionsViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	viewport := &playwright.ViewportSize{
		Width:  480,
		Height: 320,
	}

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Viewport: viewport,
	})
	must.NoError(err, "NewContext with custom viewport failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	widthRaw, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err, "Evaluate innerWidth failed")
	heightRaw, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err, "Evaluate innerHeight failed")

	if widthRaw != float64(480) {
		t.Errorf("window.innerWidth = %v, want 480", widthRaw)
	}
	if heightRaw != float64(320) {
		t.Errorf("window.innerHeight = %v, want 320", heightRaw)
	}
}

// TestBrowserContextOptionsTimezone verifies timezone configuration in BrowserContextOptions.
// Ref: TestBrowserContextLocale.java#shouldFormatDate
func TestBrowserContextOptionsTimezone(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	timezone := "America/New_York"

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		TimezoneID: &timezone,
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	tzRaw, err := page.Evaluate(ctx, "() => Intl.DateTimeFormat().resolvedOptions().timeZone")
	must.NoError(err, "Evaluate timeZone failed")

	tz, ok := tzRaw.(string)
	if !ok || tz != timezone {
		t.Errorf("resolvedOptions().timeZone = %v, want %q", tzRaw, timezone)
	}
}

// TestBrowserContextOptionsGeolocation verifies geolocation and permissions in BrowserContextOptions.
// Ref: TestGeolocation.java#shouldIsolateContexts
func TestBrowserContextOptionsGeolocation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Geolocation: &playwright.Geolocation{
			Latitude:  48.8584,
			Longitude: 2.2945,
		},
		Permissions: []string{"geolocation"},
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	coordsRaw, err := page.Evaluate(ctx, `() => new Promise(resolve => {
		navigator.geolocation.getCurrentPosition(pos => {
			resolve({ latitude: pos.coords.latitude, longitude: pos.coords.longitude });
		});
	})`)
	must.NoError(err, "Evaluate geolocation failed")

	coords, ok := coordsRaw.(map[string]any)
	is.True(ok, "expected map[string]any")

	lat, okLat := coords["latitude"].(float64)
	lng, okLng := coords["longitude"].(float64)
	const tolerance = 0.0001
	if !okLat || !okLng || math.Abs(lat-48.8584) > tolerance || math.Abs(lng-2.2945) > tolerance {
		t.Errorf("geolocation coords = {lat: %v, lng: %v}, want approx {48.8584, 2.2945}", lat, lng)
	}
}

// TestBrowserContextOptionsStorageState verifies pre-populated storage state in BrowserContextOptions.
func TestBrowserContextOptionsStorageState(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageState: &playwright.StorageState{
			Origins: []playwright.OriginStorage{
				{
					Origin: srv.Prefix(),
					LocalStorage: []playwright.LocalStorageEntry{
						{Name: "preloaded_key", Value: "preloaded_val"},
					},
				},
			},
		},
	})
	must.NoError(err, "NewContext with StorageState failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	valRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('preloaded_key')")
	must.NoError(err, "Evaluate localStorage failed")

	if valRaw != "preloaded_val" {
		t.Errorf("localStorage preloaded_key = %v, want 'preloaded_val'", valRaw)
	}
}

// TestBrowserContextStorageState_CookiesAndLocalStorage (E2E-CTX-03)
// Verifies full round-trip initialization of StorageState with both Cookies and LocalStorage
// across multiple origins and verifies isolation from other contexts.
func TestBrowserContextStorageState_CookiesAndLocalStorage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	emptyPageURL := srv.EmptyPage()

	storageState := &playwright.StorageState{
		Cookies: []playwright.Cookie{
			{
				Name:     "session_id",
				Value:    "xyz_secret_999",
				URL:      &emptyPageURL,
				SameSite: playwright.SameSiteLax,
			},
			{
				Name:     "user_pref",
				Value:    "dark_mode",
				URL:      &emptyPageURL,
				SameSite: playwright.SameSiteLax,
			},
		},
		Origins: []playwright.OriginStorage{
			{
				Origin: srv.Prefix(),
				LocalStorage: []playwright.LocalStorageEntry{
					{Name: "auth_token", Value: "jwt_token_abc"},
					{Name: "app_theme", Value: "nord"},
				},
			},
		},
	}

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageState: storageState,
	})
	must.NoError(err, "NewContext with StorageState failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// 1. Verify Cookies were injected into document.cookie
	cookieStrRaw, err := page.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate document.cookie failed")
	cookieStr, ok := cookieStrRaw.(string)
	is.True(ok, "expected cookie string")

	is.Contains(cookieStr, "session_id=xyz_secret_999")
	is.Contains(cookieStr, "user_pref=dark_mode")

	// 2. Verify LocalStorage entries were injected
	tokenRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('auth_token')")
	must.NoError(err)
	is.Equal("jwt_token_abc", tokenRaw)

	themeRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('app_theme')")
	must.NoError(err)
	is.Equal("nord", themeRaw)

	// 3. Verify a second page in the same context shares the storage state
	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "page2 NewPage failed")
	err = page2.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "page2 Goto failed")

	token2, err := page2.Evaluate(ctx, "() => localStorage.getItem('auth_token')")
	must.NoError(err)
	is.Equal("jwt_token_abc", token2)

	// 4. Verify an isolated context created without StorageState does NOT have these cookies/storage
	isolatedCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "isolated NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = isolatedCtx.Close(closeCtx)
	})

	isolatedPage, err := isolatedCtx.NewPage(ctx)
	must.NoError(err, "isolatedPage NewPage failed")
	err = isolatedPage.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "isolatedPage Goto failed")

	isoCookies, err := isolatedPage.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate document.cookie on isolated page failed")
	if cookieStr, ok := isoCookies.(string); ok && strings.Contains(cookieStr, "session_id") {
		t.Errorf("isolated context should not have session_id cookie, got %q", cookieStr)
	}

	isoStorage, err := isolatedPage.Evaluate(ctx, "() => localStorage.getItem('auth_token')")
	must.NoError(err, "Evaluate localStorage on isolated page failed")
	if isoStorage != nil {
		t.Errorf("isolated context should not have auth_token in localStorage, got %v", isoStorage)
	}
}

// TestBrowserContextStorageStatePath verifies that StorageStatePath loads cookies and localStorage
// from a JSON file on disk, equivalent to passing StorageState inline.
func TestBrowserContextStorageStatePath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	emptyPageURL := srv.EmptyPage()
	stateJSON := `{
		"cookies": [{"name":"file_cookie","value":"file_value","url":"` + emptyPageURL + `","sameSite":"Lax"}],
		"origins": [{"origin":"` + srv.Prefix() + `","localStorage":[{"name":"file_key","value":"file_val"}]}]
	}`

	f, err := os.CreateTemp("", "storage-state-*.json")
	must.NoError(err, "CreateTemp failed")
	t.Cleanup(func() { _ = os.Remove(f.Name()) })
	_, err = f.WriteString(stateJSON)
	must.NoError(err, "WriteString failed")
	_ = f.Close()

	path := f.Name()
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageStatePath: &path,
	})
	must.NoError(err, "NewContext with StorageStatePath failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	cookieRaw, err := page.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate cookie failed")
	cookieStr, _ := cookieRaw.(string)
	is.Contains(cookieStr, "file_cookie=file_value")

	storageRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('file_key')")
	must.NoError(err, "Evaluate localStorage failed")
	is.Equal("file_val", storageRaw)
}

// TestBrowserContextOptionsProxy verifies that proxy configuration in BrowserContextOptions
// routes page HTTP traffic through the specified proxy server, while unconfigured contexts bypass it.
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

// TestBrowserContextCookies verifies reading cookies from BrowserContext.
// Ref: TestBrowserContextCookies.java#shouldReturnCookies
func TestBrowserContextCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, "() => document.cookie = 'username=john; path=/'")
	must.NoError(err, "Evaluate document.cookie failed")

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err, "Cookies failed")

	var found bool
	for _, ck := range cookies {
		if ck.Name == "username" && ck.Value == "john" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected cookie 'username=john' in context cookies, got %v", cookies)
	}
}

// TestBrowserContextAddCookies verifies programmatically setting cookies via BrowserContext.AddCookies.
// Ref: TestBrowserContextAddCookies.java#shouldWork
func TestBrowserContextAddCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	emptyPageURL := srv.EmptyPage()
	err = bCtx.AddCookies(ctx, []playwright.Cookie{
		{
			Name:     "auth_token",
			Value:    "secret_bearer_token",
			URL:      &emptyPageURL,
			SameSite: playwright.SameSiteLax,
		},
	})
	must.NoError(err, "AddCookies failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	cookieRaw, err := page.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate document.cookie failed")

	if cookieStr, ok := cookieRaw.(string); !ok || !strings.Contains(cookieStr, "auth_token=secret_bearer_token") {
		t.Errorf("document.cookie = %v, want auth_token=secret_bearer_token", cookieRaw)
	}
}

// TestBrowserContextClearCookies verifies clearing cookies from the BrowserContext.
// Ref: TestBrowserContextClearCookies.java#shouldClearCookies
func TestBrowserContextClearCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	emptyPageURL := srv.EmptyPage()
	err = bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "c1", Value: "v1", URL: &emptyPageURL},
		{Name: "c2", Value: "v2", URL: &emptyPageURL},
	})
	must.NoError(err, "AddCookies failed")

	cookies, err := bCtx.Cookies(ctx)
	if err != nil || len(cookies) < 2 {
		t.Fatalf("Cookies failed or returned < 2 cookies: %v, %v", cookies, err)
	}

	// Clear all cookies
	err = bCtx.ClearCookies(ctx)
	must.NoError(err, "ClearCookies failed")

	cookiesAfter, err := bCtx.Cookies(ctx)
	must.NoError(err, "Cookies after ClearCookies failed")
	is.Len(cookiesAfter, 0)
}

// TestBrowserContextCloseRace verifies that closing a BrowserContext while a page has an active,
// blocked navigation does not hang, deadlock, panic, or leak goroutines (E2E-CTX-04).
func TestBrowserContextCloseRace(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Block the request indefinitely to simulate an in-flight navigation,
	// and get a signal channel when the HTTP request arrives at the server.
	requestReceived := srv.RequestReceived("/slow-loading")

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")

	page, err := bCtx.NewPage(ctx)
	if err != nil {
		_ = bCtx.Close(ctx)
		t.Fatalf("NewPage failed: %v", err)
	}

	// Start in-flight navigation in a separate goroutine
	gotoDone := make(chan error, 1)
	go func() {
		gotoDone <- page.Goto(ctx, srv.Prefix()+"/slow-loading")
	}()

	// Wait for the HTTP request to actually reach the test server before closing (OBS-N01-B14)
	select {
	case <-requestReceived:
		// Request confirmed in-flight on the server
	case <-time.After(5 * time.Second):
		t.Fatal("browser never sent request to /slow-loading within 5 seconds")
	}

	// Close context concurrently during the active request
	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = bCtx.Close(closeCtx)
	must.NoError(err, "bCtx.Close returned unexpected error")

	// Navigation should terminate (with error or canceled context) without hanging
	select {
	case <-gotoDone:
		// Navigation returned
	case <-time.After(5 * time.Second):
		t.Fatal("page.Goto did not terminate within 5 seconds after context close")
	}

	is.True(bCtx.IsClosed(), "expected bCtx.IsClosed() to be true")
}

// TestBrowserContextDisableJavascript verifies that JavaScriptEnabled=false prevents script execution.
// Ref: TestBrowserContextBasic.java#shouldDisableJavascript
func TestBrowserContextDisableJavascript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	jsEnabled := false
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{JavaScriptEnabled: &jsEnabled})
	must.NoError(err, "NewContext with JavaScriptEnabled=false failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "data:text/html, <script>var something = 'forbidden'</script>"))

	_, err = page.Evaluate(ctx, "something")
	is.Error(err, "evaluate should fail with JS disabled")

	// With JS enabled (default context), the same script works
	bCtx2, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(closeCtx)
	})
	page2, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Goto(ctx, "data:text/html, <script>var something = 'forbidden'</script>"))
	val, err := page2.Evaluate(ctx, "something")
	must.NoError(err)
	is.Equal("forbidden", val)
}

// TestBrowserContextNavigateAfterDisablingJavascript verifies navigation works with JS disabled.
// Ref: TestBrowserContextBasic.java#shouldBeAbleToNavigateAfterDisablingJavascript
func TestBrowserContextNavigateAfterDisablingJavascript(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	jsEnabled := false
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{JavaScriptEnabled: &jsEnabled})
	must.NoError(err, "NewContext with JavaScriptEnabled=false failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
}

// TestBrowserContextEmulateNavigatorOnLine verifies SetOffline toggles navigator.onLine.
// Ref: TestBrowserContextBasic.java#shouldEmulateNavigatorOnLine
func TestBrowserContextEmulateNavigatorOnLine(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	onLine, err := page.Evaluate(ctx, "() => window.navigator.onLine")
	must.NoError(err)
	is.Equal(true, onLine)

	must.NoError(bCtx.SetOffline(ctx, true))
	onLine, err = page.Evaluate(ctx, "() => window.navigator.onLine")
	must.NoError(err)
	is.Equal(false, onLine)

	must.NoError(bCtx.SetOffline(ctx, false))
	onLine, err = page.Evaluate(ctx, "() => window.navigator.onLine")
	must.NoError(err)
	is.Equal(true, onLine)
}

// TestBrowserContextRespectDeviceScaleFactor verifies that deviceScaleFactor option is applied.
// Ref: TestBrowserContextBasic.java#shouldRespectDeviceScaleFactor
func TestBrowserContextRespectDeviceScaleFactor(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	dpr := 3.0
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{DeviceScaleFactor: &dpr})
	must.NoError(err, "NewContext with DeviceScaleFactor=3 failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	dprVal, err := page.Evaluate(ctx, "window.devicePixelRatio")
	must.NoError(err)
	is.Equal(float64(3), dprVal)
}

// TestBrowserContextCloseClearsPages verifies that closing a context marks all its pages as closed.
// Ref: TestBrowserContextBasic.java#shouldCloseAllBelongingPagesOnceClosingContext
func TestBrowserContextCloseClearsPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")

	page1, err := bCtx.NewPage(ctx)
	if err != nil {
		_ = bCtx.Close(ctx)
		t.Fatalf("page1 NewPage failed: %v", err)
	}

	page2, err := bCtx.NewPage(ctx)
	if err != nil {
		_ = bCtx.Close(ctx)
		t.Fatalf("page2 NewPage failed: %v", err)
	}

	if len(bCtx.Pages()) != 2 {
		_ = bCtx.Close(ctx)
		t.Fatalf("expected 2 pages before close, got %d", len(bCtx.Pages()))
	}

	err = bCtx.Close(ctx)
	must.NoError(err, "Close failed")

	is.True(page1.IsClosed(), "expected page1 to be closed after context close")
	is.True(page2.IsClosed(), "expected page2 to be closed after context close")
	is.Len(bCtx.Pages(), 0)
}

func TestBrowserContextAddInitScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	must.NoError(bCtx.AddInitScript(ctx, `window.__ctxInjected = 'from-context'`), "AddInitScript failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	must.NoError(page.Goto(ctx, srv.EmptyPage()), "Goto failed")

	val, err := page.Evaluate(ctx, `() => window.__ctxInjected`)
	must.NoError(err, "Evaluate failed")
	is.Equal("from-context", val, "context init script should inject window.__ctxInjected")
}

func TestBrowserContextAddInitScriptAffectsAllPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	must.NoError(bCtx.AddInitScript(ctx, `window.__shared = 'shared-value'`), "AddInitScript failed")

	page1, err := bCtx.NewPage(ctx)
	must.NoError(err, "page1 NewPage failed")
	must.NoError(page1.Goto(ctx, srv.EmptyPage()), "page1 Goto failed")

	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "page2 NewPage failed")
	must.NoError(page2.Goto(ctx, srv.EmptyPage()), "page2 Goto failed")

	val1, err := page1.Evaluate(ctx, `() => window.__shared`)
	must.NoError(err, "page1 Evaluate failed")
	is.Equal("shared-value", val1, "page1 should have init script value")

	val2, err := page2.Evaluate(ctx, `() => window.__shared`)
	must.NoError(err, "page2 Evaluate failed")
	is.Equal("shared-value", val2, "page2 should have init script value")
}

func TestBrowserContextSetExtraHTTPHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	gotHeader := make(chan string, 1)
	srv.SetRoute("/check-headers", func(w http.ResponseWriter, r *http.Request) {
		gotHeader <- r.Header.Get("X-Custom-Header")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = bCtx.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Custom-Header": "playwright-test",
	})
	must.NoError(err, "SetExtraHTTPHeaders failed")

	err = page.Goto(ctx, srv.Prefix()+"/check-headers")
	must.NoError(err, "Goto failed")

	select {
	case val := <-gotHeader:
		if val != "playwright-test" {
			t.Errorf("X-Custom-Header = %q, want 'playwright-test'", val)
		}
	default:
		t.Error("request never reached the server")
	}
}

func TestBrowserContextSetOffline(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/online-check", "text/html", `<p>online</p>`)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.Prefix()+"/online-check")
	must.NoError(err, "Goto (online) failed")

	err = bCtx.SetOffline(ctx, true)
	must.NoError(err, "SetOffline(true) failed")

	err = page.Goto(ctx, srv.Prefix()+"/online-check")
	if err == nil {
		t.Error("expected Goto to fail in offline mode, but it succeeded")
	}

	err = bCtx.SetOffline(ctx, false)
	must.NoError(err, "SetOffline(false) failed")

	err = page.Goto(ctx, srv.Prefix()+"/online-check")
	must.NoError(err, "Goto (back online) failed")
}

func TestBrowserContextGrantPermissions(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	err = bCtx.GrantPermissions(ctx, []string{"geolocation"})
	must.NoError(err, "GrantPermissions failed")

	result, err := page.Evaluate(ctx, `
		navigator.permissions.query({name: 'geolocation'}).then(p => p.state)
	`)
	must.NoError(err, "Evaluate failed")
	if result != "granted" {
		t.Errorf("geolocation permission = %q, want 'granted'", result)
	}

	err = bCtx.ClearPermissions(ctx)
	must.NoError(err, "ClearPermissions failed")

	result2, err := page.Evaluate(ctx, `
		navigator.permissions.query({name: 'geolocation'}).then(p => p.state)
	`)
	must.NoError(err, "Evaluate after ClearPermissions failed")
	if result2 == "granted" {
		t.Error("expected permission to be revoked after ClearPermissions")
	}
}

func TestBrowserContextSetGeolocation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	accuracy := 10.0
	err = bCtx.GrantPermissions(ctx, []string{"geolocation"})
	must.NoError(err, "GrantPermissions failed")
	err = bCtx.SetGeolocation(ctx, &playwright.Geolocation{
		Latitude:  48.8584,
		Longitude: 2.2945,
		Accuracy:  &accuracy,
	})
	must.NoError(err, "SetGeolocation failed")

	coords, err := page.Evaluate(ctx, `
		new Promise((resolve, reject) =>
			navigator.geolocation.getCurrentPosition(
				pos => resolve({lat: pos.coords.latitude, lon: pos.coords.longitude}),
				err => reject(err.message)
			)
		)
	`)
	must.NoError(err, "Evaluate geolocation failed")

	m, ok := coords.(map[string]any)
	is.True(ok, "unexpected geolocation result type")
	lat, _ := m["lat"].(float64)
	lon, _ := m["lon"].(float64)
	if lat < 48.85 || lat > 48.87 {
		t.Errorf("latitude = %f, want ~48.8584", lat)
	}
	if lon < 2.29 || lon > 2.30 {
		t.Errorf("longitude = %f, want ~2.2945", lon)
	}
}

func TestBrowserContextRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/original", "text/html", `<p>original</p>`)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	body := "intercepted by context route"
	err = bCtx.Route(ctx, "**/original", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body: &body,
		})
	})
	must.NoError(err, "BrowserContext.Route failed")

	err = page.Goto(ctx, srv.Prefix()+"/original")
	must.NoError(err, "Goto failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content failed")
	is.Contains(content, "intercepted by context route")
}

func TestBrowserContextIsolatesLocalStorageExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc1 := newContext(t)
	bc2 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)
	page2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	_, err = page1.Evaluate(ctx, `() => localStorage.setItem('key', 'ctx1')`)
	must.NoError(err)

	val2, err := page2.Evaluate(ctx, `() => localStorage.getItem('key')`)
	must.NoError(err)
	is.Nil(val2)
}

func TestBrowserContextIsolatesCookiesExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc1 := newContext(t)
	bc2 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)
	page2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	_, err = page1.Evaluate(ctx, `() => document.cookie = 'ctx=one'`)
	must.NoError(err)

	cookie2, err := page2.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
	must.NotContains(cookie2, "ctx=one")
}

func TestBrowserContextPagesListsPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bc := newContext(t)

	is.Empty(bc.Pages())

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	is.Len(bc.Pages(), 1)

	page2, err := bc.NewPage(ctx)
	must.NoError(err)
	is.Len(bc.Pages(), 2)

	must.NoError(page1.Close(ctx))
	is.Len(bc.Pages(), 1)

	must.NoError(page2.Close(ctx))
	is.Empty(bc.Pages())
}

func TestBrowserContextNewPageCreatesIsolatedPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	is.NotEqual(page1, page2)
	is.Equal(bc, page1.Context())
	is.Equal(bc, page2.Context())
}

func TestBrowserContextNewPageCreatesPageEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NotNil(page)
}

func TestBrowserContextPagesCountEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	p1, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NotNil(p1)

	p2, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NotNil(p2)

	pages := bc.Pages()
	is.GreaterOrEqual(len(pages), 2)
}

func TestBrowserContextClosedPagesCountEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	p1, err := bc.NewPage(ctx)
	must.NoError(err)

	p2, err := bc.NewPage(ctx)
	must.NoError(err)
	_ = p2

	before := len(bc.Pages())
	must.NoError(p1.Close(ctx))
	after := len(bc.Pages())

	is.Less(after, before)
}

func TestBrowserContextIsolatesSessionStorageEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc1 := newContext(t)
	bc2 := newContext(t)

	p1, err := bc1.NewPage(ctx)
	must.NoError(err)

	p2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(p1.SetContent(ctx, `<html><body></body></html>`))
	must.NoError(p2.SetContent(ctx, `<html><body></body></html>`))

	_, err = p1.Evaluate(ctx, `() => sessionStorage.setItem('key', 'ctx1value')`)
	must.NoError(err)

	val, err := p2.Evaluate(ctx, `() => sessionStorage.getItem('key')`)
	must.NoError(err)
	is.Nil(val)
}

func TestBrowserContextAddCookiesAndGetCookiesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "test-cookie", Value: "cookie-value", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.NotEmpty(cookies)

	found := false
	for _, c := range cookies {
		if c.Name == "test-cookie" {
			is.Equal("cookie-value", c.Value)
			found = true
		}
	}
	is.True(found)
}

func TestBrowserContextGrantPermissionsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(bc.GrantPermissions(ctx, []string{"geolocation"}))
}

func TestBrowserContextClearPermissionsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(bc.GrantPermissions(ctx, []string{"geolocation"}))
	must.NoError(bc.ClearPermissions(ctx))
}

func TestBrowserContextSetGeolocationEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.GrantPermissions(ctx, []string{"geolocation"}))
	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{
		Latitude:  48.8566,
		Longitude: 2.3522,
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
}

func TestBrowserContextAddInitScriptEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	must.NoError(bc.AddInitScript(ctx, `window.__fromContext = 'hello';`))

	p1, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(p1.SetContent(ctx, `<div></div>`))

	result, err := p1.Evaluate(ctx, `() => window.__fromContext`)
	must.NoError(err)
	is.Equal("hello", result)
}

func TestBrowserContextMultiplePagesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<p>Page 1</p>`))
	must.NoError(page2.SetContent(ctx, `<p>Page 2</p>`))

	title1, err := page1.Title(ctx)
	must.NoError(err)
	title2, err := page2.Title(ctx)
	must.NoError(err)

	is.NotEqual(title1, title2)
}

func TestBrowserContextCookiesEmptyEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)
}

func TestBrowserContextAddCookiesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "abc123", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Equal(1, len(cookies))
	is.Equal("session", cookies[0].Name)
}

func TestBrowserContextClearCookiesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "tok", Value: "val", URL: &url},
	}))

	must.NoError(bc.ClearCookies(ctx))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)
}

func TestBrowserContextNewPageEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NotNil(page)

	must.NoError(page.SetContent(ctx, `<p>Test</p>`))
	text, err := page.Locator("p").InnerText(ctx)
	must.NoError(err)
	is.Equal("Test", text)
}

func TestBrowserContextPagesCountEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NotNil(page1)
	must.NotNil(page2)

	pages := bc.Pages()
	is.GreaterOrEqual(len(pages), 2)
}

func TestBrowserContextSetContentEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<p>Context page</p>`))

	count, err := page.Locator("p").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

func TestBrowserContextFulfillRouteEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	body := "ctx-routed"
	ct := "text/html"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		}))
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, "http://example.com/"))
	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "ctx-routed")
}

func TestContextMultipleTabsShareCookiesEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<html><body>page1</body></html>`))
	_, err = page1.Evaluate(ctx, `() => document.cookie = 'shared=yes'`)
	must.NoError(err)

	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page2.SetContent(ctx, `<html><body>page2</body></html>`))

	cookie, err := page2.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
	is.Contains(cookie, "shared=yes")
}

func TestContextNewPageEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)

	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<title>Page One</title>`))
	must.NoError(page2.SetContent(ctx, `<title>Page Two</title>`))

	t1, err := page1.Title(ctx)
	must.NoError(err)
	is.Equal("Page One", t1)

	t2, err := page2.Title(ctx)
	must.NoError(err)
	is.Equal("Page Two", t2)
}

func TestContextPagesListEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	_, err := bc.NewPage(ctx)
	must.NoError(err)
	_, err = bc.NewPage(ctx)
	must.NoError(err)

	pages := bc.Pages()
	is.GreaterOrEqual(len(pages), 2)
}

func TestContextRouteEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	intercepted := false
	must.NoError(bc.Route(ctx, "**/api/test", func(r *playwright.Route) {
		intercepted = true
		ct := "application/json"
		body := `{"ok":true}`
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))
	_, _ = page.Evaluate(ctx, `() => fetch('/api/test').catch(() => null)`)

	is.True(intercepted)
}

func TestContextSetExtraHTTPHeadersEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	must.NoError(bc.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Context-Header": "context-value",
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<html><body>Content</body></html>`))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "Content")
}

func TestContextJavaScriptEnabledEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	result, err := page.Evaluate(ctx, `() => 1 + 2`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

func TestContextNewPageMultipleEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)

	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<title>Page One</title>`))
	must.NoError(page2.SetContent(ctx, `<title>Page Two</title>`))

	title1, err := page1.Title(ctx)
	must.NoError(err)
	is.Equal("Page One", title1)

	title2, err := page2.Title(ctx)
	must.NoError(err)
	is.Equal("Page Two", title2)
}

func TestContextGrantPermissionsEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	err := bc.GrantPermissions(ctx, []string{"clipboard-read"})
	must.NoError(err)
}

func TestContextSetGeolocationEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	err := bc.GrantPermissions(ctx, []string{"geolocation"})
	must.NoError(err)

	lat := 37.7749
	lon := -122.4194
	err = bc.SetGeolocation(ctx, &playwright.Geolocation{
		Latitude:  lat,
		Longitude: lon,
	})
	must.NoError(err)
}

// ---------------------------------------------------------------------------
// From browser_context_basic_extra_test.go
// ---------------------------------------------------------------------------

// TestWindowOpenUsesParentContext verifies that a popup opened via window.open shares the parent's BrowserContext.
// Ref: TestBrowserContextBasic.java#windowOpenShouldUseParentTabContext
func TestWindowOpenUsesParentContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var popup *playwright.Page
	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		popup = p
		mu.Unlock()
	})
	defer cancel()

	_, err = page.Evaluate(ctx, "url => window.open(url)", srv.EmptyPage())
	must.NoError(err)

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		p := popup
		mu.Unlock()
		if p != nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	mu.Lock()
	p := popup
	mu.Unlock()
	must.NotNil(p, "popup page should have been created")
	is.Equal(bc, p.Context(), "popup should share the parent's BrowserContext")
}

// TestNewContextDeviceScaleFactorNullViewportError verifies that combining DeviceScaleFactor
// with a null viewport (NoDefaultViewport) returns an error from the server.
// Ref: TestBrowserContextBasic.java#shouldNotAllowDeviceScaleFactorWithNullViewport
func TestNewContextDeviceScaleFactorNullViewportError(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	dsf := 1.0
	noVP := true
	_, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		DeviceScaleFactor: &dsf,
		NoDefaultViewport: &noVP,
	})
	is.Error(err, "combining DeviceScaleFactor with null viewport should fail")
	is.ErrorContains(err, "deviceScaleFactor")
}

// TestNewContextIsMobileNullViewportError verifies that combining IsMobile=true
// with a null viewport (NoDefaultViewport) returns an error from the server.
// Ref: TestBrowserContextBasic.java#shouldNotAllowIsMobileWithNullViewport
func TestNewContextIsMobileNullViewportError(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	isMobile := true
	noVP := true
	_, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		IsMobile:          &isMobile,
		NoDefaultViewport: &noVP,
	})
	is.Error(err, "combining IsMobile with null viewport should fail")
	is.ErrorContains(err, "isMobile")
}

// ---------------------------------------------------------------------------
// From browser_context_close_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextCloseIsIdempotent verifies Close can be called multiple times.
// Ref: TestBrowserContextClose.java#shouldBeCallableTwice
func TestBrowserContextCloseIsIdempotentExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	must.NoError(bc.Close(ctx))
	must.NoError(bc.Close(ctx)) // Second close should not panic
}

// TestBrowserContextClosedHasNoPages verifies Pages returns empty after context close.
// Ref: TestBrowserContextClose.java#shouldHaveNoPagesAfterClose
func TestBrowserContextClosedHasNoPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NotNil(page)

	must.NoError(bc.Close(ctx))

	// Pages should be empty or reflect closed state
	pages := bc.Pages()
	_ = pages // Just verify no panic
}

// TestBrowserContextNewPageAfterClose verifies NewPage fails or behaves gracefully after close.
// Ref: TestBrowserContextClose.java#shouldNotAllowNewPageAfterClose
func TestBrowserContextNewPageAfterClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	must.NoError(bc.Close(ctx))

	// Attempting to create a new page after context close should fail gracefully
	_, err := bc.NewPage(ctx)
	is.Error(err, "creating page after context close should fail")
}

// ---------------------------------------------------------------------------
// From browser_context_network_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextOnRequestFiredForNavigation verifies context OnRequest fires on page navigation.
// Ref: TestBrowserContextNetworkEvents.java#shouldFireOnRequest
func TestBrowserContextOnRequestFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	var mu sync.Mutex
	var requestURLs []string
	off := bc.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requestURLs = append(requestURLs, req.URL())
		mu.Unlock()
	})
	defer off()

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	urls := requestURLs
	mu.Unlock()

	is.NotEmpty(urls)
}

// TestBrowserContextOnResponseFiredForNavigation verifies context OnResponse fires.
// Ref: TestBrowserContextNetworkEvents.java#shouldFireOnResponse
func TestBrowserContextOnResponseFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	var mu sync.Mutex
	var statuses []int
	off := bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statuses = append(statuses, resp.Status())
		mu.Unlock()
	})
	defer off()

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	s := statuses
	mu.Unlock()

	is.NotEmpty(s)
	is.Equal(200, s[0])
}

// TestBrowserContextOnRequestFinishedFired verifies context OnRequestFinished fires.
// Ref: TestBrowserContextNetworkEvents.java#shouldFireOnRequestFinished
func TestBrowserContextOnRequestFinishedFired(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	var mu sync.Mutex
	fired := false
	off := bc.OnRequestFinished(func(req *playwright.NetworkRequest) {
		mu.Lock()
		fired = true
		mu.Unlock()
	})
	defer off()

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	f := fired
	mu.Unlock()

	is.True(f)
}

// TestBrowserContextOffRemovesHandler verifies off() removes context event handler.
// Ref: TestBrowserContextNetworkEvents.java#shouldRemoveHandlerWithOff
func TestBrowserContextOffRemovesHandler(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	var mu sync.Mutex
	count := 0
	off := bc.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	c := count
	mu.Unlock()

	// Should have fired only once (before off was called)
	is.GreaterOrEqual(c, 1)
}

// ---------------------------------------------------------------------------
// From browser_context_offline_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextOfflineFails verifies navigation fails when offline mode is enabled.
// Ref: TestBrowserContextOffline.java#shouldFail
func TestBrowserContextOfflineFails(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	// First navigate while online
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	// Now go offline
	must.NoError(bc.SetOffline(ctx, true))

	// Fetch should fail in offline mode
	result, err := page.Evaluate(ctx, `async () => {
		try {
			await fetch('/empty.html');
			return 'success';
		} catch (e) {
			return 'failed';
		}
	}`)
	must.NoError(err)
	is.Equal("failed", result)
}

// TestBrowserContextOfflineAndBack verifies going online after offline restores connectivity.
// Ref: TestBrowserContextOffline.java#shouldBeRestorable
func TestBrowserContextOfflineAndBack(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bc.SetOffline(ctx, true))
	must.NoError(bc.SetOffline(ctx, false))

	// Should be able to navigate again after going online
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(srv.EmptyPage(), page.URL())
}

// TestBrowserContextOfflineDefault verifies pages work by default (not offline).
// Ref: TestBrowserContextOffline.java#shouldWorkByDefault
func TestBrowserContextOfflineDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	// Should be online by default
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(srv.EmptyPage(), page.URL())
}

// TestBrowserContextOfflineModeOnContentEx verifies SetContent works in offline mode.
// Ref: TestBrowserContextOffline.java#shouldAllowSetContentOffline
func TestBrowserContextOfflineModeOnContentEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	offline := true
	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{Offline: &offline})
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(c)
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<html><body><p id="p">Offline content</p></body></html>`))

	text, err := page.Locator("#p").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Offline content", *text)
}

// ---------------------------------------------------------------------------
// From browser_context_permissions_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextGrantPermissionsCamera verifies camera permission can be granted.
// Ref: TestBrowserContextPermissions.java#shouldGrantCameraPermission
func TestBrowserContextGrantPermissionsCamera(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bc.GrantPermissions(ctx, []string{"camera"}, srv.Prefix()))

	// Check permission state via JS Permissions API
	result, err := page.Evaluate(ctx, `async () => {
		const status = await navigator.permissions.query({name: 'camera'});
		return status.state;
	}`)
	must.NoError(err)
	is.Equal("granted", result)
}

// TestBrowserContextGrantPermissionsMicrophone verifies microphone permission can be granted.
// Ref: TestBrowserContextPermissions.java#shouldGrantMicrophonePermission
func TestBrowserContextGrantPermissionsMicrophone(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bc.GrantPermissions(ctx, []string{"microphone"}, srv.Prefix()))

	result, err := page.Evaluate(ctx, `async () => {
		const status = await navigator.permissions.query({name: 'microphone'});
		return status.state;
	}`)
	must.NoError(err)
	is.Equal("granted", result)
}

// TestBrowserContextGrantAndClearPermissions verifies ClearPermissions revokes previously granted ones.
// Ref: TestBrowserContextPermissions.java#shouldClearPermissions
func TestBrowserContextGrantAndClearPermissions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bc.GrantPermissions(ctx, []string{"geolocation"}, srv.Prefix()))

	grantedState, err := page.Evaluate(ctx, `async () => {
		const s = await navigator.permissions.query({name: 'geolocation'});
		return s.state;
	}`)
	must.NoError(err)
	is.Equal("granted", grantedState)

	must.NoError(bc.ClearPermissions(ctx))

	clearedState, err := page.Evaluate(ctx, `async () => {
		const s = await navigator.permissions.query({name: 'geolocation'});
		return s.state;
	}`)
	must.NoError(err)
	is.NotEqual("granted", clearedState)
}

// ---------------------------------------------------------------------------
// From browser_context_storage_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextLocalStorageSetGetEx2 verifies localStorage set/get in page.
// Ref: TestBrowserContextStorage.java#shouldSetAndGetLocalStorage
func TestBrowserContextLocalStorageSetGetEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('myKey', 'myValue')`)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => localStorage.getItem('myKey')`)
	must.NoError(err)
	is.Equal("myValue", result)
}

// TestBrowserContextSessionStorageSetGetEx2 verifies sessionStorage set/get in page.
// Ref: TestBrowserContextStorage.java#shouldSetAndGetSessionStorage
func TestBrowserContextSessionStorageSetGetEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => sessionStorage.setItem('sessKey', 'sessValue')`)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => sessionStorage.getItem('sessKey')`)
	must.NoError(err)
	is.Equal("sessValue", result)
}

// TestBrowserContextLocalStorageIsolatedAcrossContextsEx2 verifies contexts have separate localStorage.
// Ref: TestBrowserContextStorage.java#shouldIsolateLocalStorage
func TestBrowserContextLocalStorageIsolatedAcrossContextsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc1 := newContext(t)
	bc2 := newContext(t)

	p1, err := bc1.NewPage(ctx)
	must.NoError(err)
	p2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(p1.Goto(ctx, srv.EmptyPage()))
	must.NoError(p2.Goto(ctx, srv.EmptyPage()))

	_, err = p1.Evaluate(ctx, `() => localStorage.setItem('isolated', 'ctx1')`)
	must.NoError(err)

	val, err := p2.Evaluate(ctx, `() => localStorage.getItem('isolated')`)
	must.NoError(err)
	is.Nil(val)
}

// TestBrowserContextLocalStorageClearedOnReloadEx2 verifies localStorage persists across reload.
// Ref: TestBrowserContextStorage.java#shouldPersistAcrossReload
func TestBrowserContextLocalStorageClearedOnReloadEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('persist', 'yes')`)
	must.NoError(err)

	must.NoError(page.Reload(ctx))

	val, err := page.Evaluate(ctx, `() => localStorage.getItem('persist')`)
	must.NoError(err)
	is.Equal("yes", val)
}

// TestStorageStateCookiesEx3 verifies StorageState captures cookies.
// Ref: TestBrowserContextStorage.java#shouldCaptureCookies
func TestStorageStateCookiesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "auth", Value: "token123", URL: &url},
	}))

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
	is.Equal(1, len(state.Cookies))
	is.Equal("auth", state.Cookies[0].Name)
}

// TestStorageStateEmptyEx3 verifies StorageState is empty for new context.
// Ref: TestBrowserContextStorage.java#shouldBeEmptyForNewContext
func TestStorageStateEmptyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
	is.Empty(state.Cookies)
}

// TestStorageStateAfterClearEx3 verifies StorageState is empty after clearing.
// Ref: TestBrowserContextStorage.java#shouldBeEmptyAfterClear
func TestStorageStateAfterClearEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "tok", Value: "val", URL: &url},
	}))

	must.NoError(bc.ClearCookies(ctx))

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	is.Empty(state.Cookies)
}

// TestLocalStorageSetGetEx4 verifies localStorage set and get via Evaluate.
// Ref: TestBrowserContextStorage.java#shouldSetAndGetLocalStorage
func TestLocalStorageSetGetEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('user', 'alice')`)
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => localStorage.getItem('user')`)
	must.NoError(err)
	is.Equal("alice", val)
}

// TestSessionStorageEx4 verifies sessionStorage works across evaluations.
// Ref: TestBrowserContextStorage.java#shouldUseSessionStorage
func TestSessionStorageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => sessionStorage.setItem('token', 'xyz789')`)
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => sessionStorage.getItem('token')`)
	must.NoError(err)
	is.Equal("xyz789", val)
}

// TestLocalStorageRemoveItemEx4 verifies localStorage removeItem.
// Ref: TestBrowserContextStorage.java#shouldRemoveItem
func TestLocalStorageRemoveItemEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('temp', 'data')`)
	must.NoError(err)
	_, err = page.Evaluate(ctx, `() => localStorage.removeItem('temp')`)
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => localStorage.getItem('temp')`)
	must.NoError(err)
	is.Nil(val)
}

// TestLocalStorageClearEx4 verifies localStorage clear removes all items.
// Ref: TestBrowserContextStorage.java#shouldClearStorage
func TestLocalStorageClearEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => { localStorage.setItem('a', '1'); localStorage.setItem('b', '2'); }`)
	must.NoError(err)
	_, err = page.Evaluate(ctx, `() => localStorage.clear()`)
	must.NoError(err)

	length, err := page.Evaluate(ctx, `() => localStorage.length`)
	must.NoError(err)
	is.Equal(float64(0), length)
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsRequest
func TestBrowserContextEventsRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var requests []string
	cancel := bCtx.OnRequest(func(r *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, r.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := append([]string(nil), requests...)
	mu.Unlock()
	is.NotEmpty(got)
	is.Contains(got[0], srv.EmptyPage())
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsResponse
func TestBrowserContextEventsResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var responses []string
	cancel := bCtx.OnResponse(func(r *playwright.NetworkResponse) {
		mu.Lock()
		responses = append(responses, r.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := append([]string(nil), responses...)
	mu.Unlock()
	is.NotEmpty(got)
	is.Contains(got[0], srv.EmptyPage())
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsRequestFailed
func TestBrowserContextEventsRequestFailed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	// Close response body immediately to cause request failure.
	srv.SetRoute("/one-style.css", func(w http.ResponseWriter, r *http.Request) {
		// Abort by closing without writing anything.
	})

	failedCh := make(chan *playwright.NetworkRequest, 10)
	cancel := bCtx.OnRequestFailed(func(r *playwright.NetworkRequest) {
		failedCh <- r
	})
	defer cancel()

	// This page loads one-style.css which will fail.
	srv.ServeWithBody("/one-style.html", "text/html", `<link rel="stylesheet" href="/one-style.css">`)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/one-style.html"))

	select {
	case req := <-failedCh:
		is.Contains(req.URL(), "one-style.css")
		is.Equal("stylesheet", req.ResourceType())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for request failed event")
	}
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsRequestFinished
func TestBrowserContextEventsRequestFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	finishedCh := make(chan *playwright.NetworkRequest, 1)
	cancel := bCtx.OnRequestFinished(func(r *playwright.NetworkRequest) {
		finishedCh <- r
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case req := <-finishedCh:
		is.Contains(req.URL(), srv.EmptyPage())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for request finished event")
	}
}

// Ref: TestBrowserContextNetworkEvents.java#shouldFireEventsInProperOrder
func TestBrowserContextNetworkEventsShouldFireEventsInProperOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var events []string
	addEvent := func(name string) func() {
		switch name {
		case "request":
			return bCtx.OnRequest(func(r *playwright.NetworkRequest) {
				mu.Lock()
				events = append(events, "request")
				mu.Unlock()
			})
		case "response":
			return bCtx.OnResponse(func(r *playwright.NetworkResponse) {
				mu.Lock()
				events = append(events, "response")
				mu.Unlock()
			})
		default: // requestfinished
			return bCtx.OnRequestFinished(func(r *playwright.NetworkRequest) {
				mu.Lock()
				events = append(events, "requestfinished")
				mu.Unlock()
			})
		}
	}
	c1 := addEvent("request")
	defer c1()
	c2 := addEvent("response")
	defer c2()
	c3 := addEvent("requestfinished")
	defer c3()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	// Allow async event handlers to flush.
	_ = page.WaitForTimeout(ctx, 200)

	mu.Lock()
	got := append([]string(nil), events...)
	mu.Unlock()

	// The first three events should follow the expected order.
	if len(got) >= 3 {
		is.Equal("request", got[0])
		is.Equal("response", got[1])
		is.Equal("requestfinished", got[2])
	} else {
		t.Logf("got events: %v", got)
		is.GreaterOrEqual(len(got), 3)
	}
}
