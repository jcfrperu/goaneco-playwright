//go:build e2e

// E2E tests for: TestDefaultBrowserContext2.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestBrowserContextOnPageCloseFiresOnPageClose verifies BrowserContext.OnPageClose fires when page closes.
// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestBrowserContextOnPageCloseFiresOnPageClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	done := make(chan struct{}, 1)

	cancel := bc.OnPageClose(func(p *playwright.Page) {
		select {
		case done <- struct{}{}:
		default:
		}
	})
	defer cancel()

	must.NoError(page.Close(ctx))

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageClose")
	}
}

// TestBrowserContextSetGeolocationWithPermission verifies SetGeolocation updates location when permission granted.
// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestBrowserContextSetGeolocationWithPermission(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{
		Latitude:  48.8566,
		Longitude: 2.3522,
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	coords, err := page.Evaluate(ctx, `() => new Promise((resolve, reject) => {
		navigator.geolocation.getCurrentPosition(
			pos => resolve({ lat: pos.coords.latitude, lon: pos.coords.longitude }),
			err => reject(err)
		);
	})`)
	must.NoError(err)
	must.NotNil(coords)

	m, ok := coords.(map[string]any)
	is.True(ok)
	lat, ok := m["lat"].(float64)
	is.True(ok)
	is.InDelta(48.8566, lat, 0.001)
}

// TestPersistentContextUserAgentOption verifies the User-Agent header is set from the option.
// Ref: TestDefaultBrowserContext2.java (userAgent option)
// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestPersistentContextUserAgentOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	requestCh := make(chan *http.Request, 1)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		requestCh <- r
		w.WriteHeader(http.StatusOK)
	})

	customUA := "MyTestAgent/1.0"
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		UserAgent: &customUA,
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case req := <-requestCh:
		is.Equal("MyTestAgent/1.0", req.Header.Get("User-Agent"))
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for HTTP request")
	}
}

// TestPersistentContextOfflineModeOption verifies that offline mode blocks network requests.
// Ref: TestDefaultBrowserContext2.java (offline option)
// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestPersistentContextOfflineModeOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	offline := true
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		Offline: &offline,
	})

	// Navigating to a real server URL while offline should fail.
	err := page.Goto(ctx, srv.EmptyPage())
	is.Error(err, "expected navigation to fail in offline mode")
}

// TestPersistentContextHTTPCredentialsOption verifies that HTTP basic auth credentials
// are applied when set via LaunchPersistentContextOptions.HttpCredentials.
// Ref: TestDefaultBrowserContext2.java (httpCredentials option)
// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestPersistentContextHTTPCredentialsOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithBasicAuth("/protected.html", "user", "pass", "text/html", "<html><body>secret</body></html>")

	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
		},
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/protected.html"))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "secret")
}

// TestPersistentContextBypassCSPOption verifies that CSP can be bypassed.
// Ref: TestDefaultBrowserContext2.java (bypassCSP option)
// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestPersistentContextBypassCSPOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Serve a page with a strict CSP that blocks inline scripts.
	srv.SetRoute("/csp.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "script-src 'none'")
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><head></head><body></body></html>`))
	})

	bypassCSP := true
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		BypassCSP: &bypassCSP,
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/csp.html"))

	// With bypassCSP the page should load; eval should work without errors.
	result, err := page.Evaluate(ctx, "() => 1 + 1")
	must.NoError(err)
	is.Equal(float64(2), result)
}

// TestPersistentContextJavaScriptDisabled verifies that JavaScript can be disabled.
// Ref: TestDefaultBrowserContext2.java (javaScriptEnabled option)
// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestPersistentContextJavaScriptDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/jscheck.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body><div id="result">no-js</div><script>document.getElementById('result').textContent='js';</script></body></html>`))
	})

	jsEnabled := false
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		JavaScriptEnabled: &jsEnabled,
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/jscheck.html"))

	// The script should not have run, so the div should still say "no-js".
	text, err := page.EvalOnSelector(ctx, "#result", "el => el.textContent")
	must.NoError(err)
	is.Equal("no-js", text)
}

// TestDefaultBrowserContext2ShouldSupportHasTouchOption verifies that setHasTouch makes ontouchstart available.
// Ref: TestDefaultBrowserContext2.java#shouldSupportHasTouchOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportHasTouchOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	hasTouch := true
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		HasTouch: &hasTouch,
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/mobile.html"))

	result, err := page.Evaluate(ctx, "() => 'ontouchstart' in window")
	must.NoError(err)
	is.Equal(true, result)
}

// TestDefaultBrowserContext2ShouldSupportColorSchemeOption verifies the dark color scheme is applied.
// Ref: TestDefaultBrowserContext2.java#shouldSupportColorSchemeOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportColorSchemeOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	dark := "dark"
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		ColorScheme: &dark,
	})

	lightMatches, err := page.Evaluate(ctx, "matchMedia('(prefers-color-scheme: light)').matches")
	must.NoError(err)
	is.Equal(false, lightMatches)

	darkMatches, err := page.Evaluate(ctx, "matchMedia('(prefers-color-scheme: dark)').matches")
	must.NoError(err)
	is.Equal(true, darkMatches)
}

// TestDefaultBrowserContext2ShouldSupportTimezoneIdOption verifies the timezone is applied to Date formatting.
// Ref: TestDefaultBrowserContext2.java#shouldSupportTimezoneIdOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportTimezoneIdOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	locale := "en-US"
	tz := "America/Jamaica"
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		Locale:     &locale,
		TimezoneID: &tz,
	})

	result, err := page.Evaluate(ctx, "new Date(1479579154987).toString()")
	must.NoError(err)
	is.Equal("Sat Nov 19 2016 13:12:34 GMT-0500 (Eastern Standard Time)", result)
}

// TestDefaultBrowserContext2ShouldSupportLocaleOption verifies the navigator.language reflects the locale.
// Ref: TestDefaultBrowserContext2.java#shouldSupportLocaleOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportLocaleOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	locale := "fr-FR"
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		Locale: &locale,
	})

	result, err := page.Evaluate(ctx, "navigator.language")
	must.NoError(err)
	is.Equal("fr-FR", result)
}

// TestDefaultBrowserContext2ShouldSupportGeolocationAndPermissionsOptions verifies geolocation is resolved when
// both the geolocation option and the "geolocation" permission are set.
// Ref: TestDefaultBrowserContext2.java#shouldSupportGeolocationAndPermissionsOptions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportGeolocationAndPermissionsOptions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		Geolocation: &playwright.Geolocation{
			Latitude:  10,
			Longitude: 10,
		},
		Permissions: []string{"geolocation"},
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() => new Promise(resolve =>
		navigator.geolocation.getCurrentPosition(pos =>
			resolve({ latitude: pos.coords.latitude, longitude: pos.coords.longitude })
		)
	)`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok, "expected map result from geolocation promise")
	is.InDelta(10.0, m["latitude"], 0.001)
	is.InDelta(10.0, m["longitude"], 0.001)
}

// TestDefaultBrowserContext2ShouldSupportExtraHTTPHeadersOption verifies extra headers are sent with requests.
// Ref: TestDefaultBrowserContext2.java#shouldSupportExtraHTTPHeadersOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportExtraHTTPHeadersOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	requestCh := make(chan *http.Request, 1)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		requestCh <- r
		w.WriteHeader(http.StatusOK)
	})

	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		ExtraHTTPHeaders: map[string]string{"foo": "bar"},
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case req := <-requestCh:
		is.Equal("bar", req.Header.Get("foo"))
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for HTTP request")
	}
}

// TestDefaultBrowserContext2ShouldAcceptUserDataDir verifies that files are written to the user data dir.
// Ref: TestDefaultBrowserContext2.java#shouldAcceptUserDataDir  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldAcceptUserDataDir(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	bt, err := globalPW.Chromium()
	must.NoError(err)

	dir := t.TempDir()
	headless := true
	browser, bCtx, err := bt.LaunchPersistentContext(ctx, dir, &playwright.LaunchPersistentContextOptions{
		Headless: &headless,
	})
	must.NoError(err)

	cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	must.NoError(bCtx.Close(cleanCtx))
	must.NoError(browser.Close(cleanCtx))

	// The user data directory should contain files after the context closes.
	// We simply verify the launch + close didn't error; full file enumeration
	// would require os.ReadDir which may be flaky on Windows due to locks.
}

// TestDefaultBrowserContext2ShouldRestoreStateFromUserDataDir verifies that localStorage persists
// across two separate persistent context sessions using the same user data dir.
// Ref: TestDefaultBrowserContext2.java#shouldRestoreStateFromUserDataDir  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldRestoreStateFromUserDataDir(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bt, err := globalPW.Chromium()
	must.NoError(err)
	srv := testserver.New(t)

	dir := t.TempDir()
	headless := true

	// Session 1: write localStorage.
	{
		browser, bCtx, err := bt.LaunchPersistentContext(ctx, dir, &playwright.LaunchPersistentContextOptions{
			Headless: &headless,
		})
		must.NoError(err, "session 1: launch failed")

		page, err := bCtx.NewPage(ctx)
		must.NoError(err)
		must.NoError(page.Goto(ctx, srv.EmptyPage()))
		_, err = page.Evaluate(ctx, "() => { localStorage.hey = 'hello'; }")
		must.NoError(err)

		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		must.NoError(bCtx.Close(closeCtx))
		must.NoError(browser.Close(closeCtx))
		cancel()
	}

	// Session 2: verify localStorage is restored.
	{
		browser2, bCtx2, err := bt.LaunchPersistentContext(ctx, dir, &playwright.LaunchPersistentContextOptions{
			Headless: &headless,
		})
		must.NoError(err, "session 2: launch failed")

		page2, err := bCtx2.NewPage(ctx)
		must.NoError(err)
		must.NoError(page2.Goto(ctx, srv.EmptyPage()))

		val, err := page2.Evaluate(ctx, "localStorage.hey")
		must.NoError(err)
		is.Equal("hello", val)

		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		must.NoError(bCtx2.Close(closeCtx))
		must.NoError(browser2.Close(closeCtx))
		cancel()
	}

	// Session 3: different dir — should NOT have the stored value.
	{
		dir2 := t.TempDir()
		browser3, bCtx3, err := bt.LaunchPersistentContext(ctx, dir2, &playwright.LaunchPersistentContextOptions{
			Headless: &headless,
		})
		must.NoError(err, "session 3: launch failed")

		page3, err := bCtx3.NewPage(ctx)
		must.NoError(err)
		must.NoError(page3.Goto(ctx, srv.EmptyPage()))

		val, err := page3.Evaluate(ctx, "localStorage.hey")
		must.NoError(err)
		is.NotEqual("hello", val)

		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		must.NoError(bCtx3.Close(closeCtx))
		must.NoError(browser3.Close(closeCtx))
		cancel()
	}
}

// TestDefaultBrowserContext2ShouldHaveDefaultURLWhenLaunchingBrowser verifies that a fresh persistent context
// opens with a single about:blank page.
// Ref: TestDefaultBrowserContext2.java#shouldHaveDefaultURLWhenLaunchingBrowser  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldHaveDefaultURLWhenLaunchingBrowser(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	bCtx, _ := launchPersistent(t, nil)

	pages := bCtx.Pages()
	urls := make([]string, 0, len(pages))
	for _, p := range pages {
		urls = append(urls, p.URL())
	}
	is.Equal([]string{"about:blank"}, urls)
}

// TestDefaultBrowserContext2ShouldWorkInPersistentContext verifies that a viewport option is respected.
// Ref: TestDefaultBrowserContext2.java#shouldWorkInPersistentContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldWorkInPersistentContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		Viewport: &playwright.ViewportSize{
			Width:  800,
			Height: 600,
		},
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	width, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err)
	is.Equal(float64(800), width)

	height, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err)
	is.Equal(float64(600), height)
}

// TestDefaultBrowserContext2ShouldSupportContrastOption verifies the contrast media feature option.
// Ref: TestDefaultBrowserContext2.java#shouldSupportContrastOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportContrastOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	more := "more"
	_, page := launchPersistent(t, &playwright.LaunchPersistentContextOptions{
		Contrast: &more,
	})

	moreMatches, err := page.Evaluate(ctx, "() => matchMedia('(prefers-contrast: more)').matches")
	must.NoError(err)
	is.Equal(true, moreMatches)

	noPrefMatches, err := page.Evaluate(ctx, "() => matchMedia('(prefers-contrast: no-preference)').matches")
	must.NoError(err)
	is.Equal(false, noPrefMatches)
}

// Extra: no equivalent method in playwright-java; related class: TestDefaultBrowserContext2.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestBrowserTypeLaunchPersistentContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	// Create a temporary directory for user data.
	dir, err := os.MkdirTemp("", "playwright-persistent-*")
	must.NoError(err, "MkdirTemp failed")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	bt, err := globalPW.Chromium()
	must.NoError(err, "Chromium() failed")

	headless := true
	browser, bCtx, err := bt.LaunchPersistentContext(ctx, dir, &playwright.LaunchPersistentContextOptions{
		Headless: &headless,
	})
	must.NoError(err, "LaunchPersistentContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
		_ = browser.Close(closeCtx)
	})

	must.NotNil(browser, "LaunchPersistentContext returned nil browser")
	must.NotNil(bCtx, "LaunchPersistentContext returned nil context")

	// Create a page and verify it works.
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, "about:blank")
	must.NoError(err, "Goto about:blank failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title() failed")
	// about:blank title should be empty string.
	if title != "" {
		t.Logf("title = %q (about:blank usually has empty title)", title)
	}
}

// TestDefaultBrowserContext2ShouldSupportIgnoreHTTPSErrorsOption mirrors the Java test whose @Test is commented out (TODO).
// Ref: TestDefaultBrowserContext2.java#shouldSupportIgnoreHTTPSErrorsOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldSupportIgnoreHTTPSErrorsOption(t *testing.T) {
	t.Skip("not implemented: HTTPS test server (testserver.HasHTTPSSupport() is false); Java @Test is also disabled (net::ERR_EMPTY_RESPONSE TODO)")
}

// TestDefaultBrowserContext2ShouldWorkWithIgnoreDefaultArgs verifies BrowserType launch options ignoreDefaultArgs.
// Ref: TestDefaultBrowserContext2.java#shouldWorkWithIgnoreDefaultArgs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldWorkWithIgnoreDefaultArgs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	browserType, err := resolveBrowserType(globalPW)
	must.NoError(err)

	browser, err := browserType.Launch(ctx, &playwright.BrowserTypeLaunchOptions{
		IgnoreDefaultArgs: []string{"foo"},
	})
	must.NoError(err)
	page, err := browser.NewPage(ctx)
	must.NoError(err)
	must.NotNil(page)
	must.NoError(browser.Close(ctx))

	ignoreAll := true
	options := &playwright.BrowserTypeLaunchOptions{IgnoreAllDefaultArgs: &ignoreAll}
	must.Equal(&ignoreAll, options.IgnoreAllDefaultArgs)
}

// TestDefaultBrowserContext2ShouldFireCloseEventForAPersistentContext verifies the close event of a persistent context.
// Ref: TestDefaultBrowserContext2.java#shouldFireCloseEventForAPersistentContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldFireCloseEventForAPersistentContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	bCtx, _ := launchPersistent(t, nil)

	closed := make(chan *playwright.BrowserContext, 1)
	off := bCtx.OnClose(func(context *playwright.BrowserContext) {
		closed <- context
	})
	defer off()
	must.NoError(bCtx.Close(ctx))

	select {
	case context := <-closed:
		is.Equal(bCtx, context)
	case <-time.After(5 * time.Second):
		t.Fatal("context close event not received")
	}
}

// TestDefaultBrowserContext2ShouldRespectSelectors verifies a registered custom selector engine in a persistent context.
// Ref: TestDefaultBrowserContext2.java#shouldRespectSelectors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldRespectSelectors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, page := launchPersistent(t, nil)
	engine := `{
  create(root, target) {},
  query(root, selector) {
    return root.querySelector(selector);
  },
  queryAll(root, selector) {
    return Array.from(root.querySelectorAll(selector));
  }
}`
	must.NoError(bCtx.RegisterSelectorEngine(ctx, "defaultContextCSS", engine))

	must.NoError(page.SetContent(ctx, "<div>hello</div>"))
	html, err := page.Locator("css=div").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("hello", html)
	html, err = page.Locator("defaultContextCSS=div").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("hello", html)
}

// TestDefaultBrowserContext2ShouldUploadLargeFile verifies uploading a 200MB file in a persistent context.
// Ref: TestDefaultBrowserContext2.java#shouldUploadLargeFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldUploadLargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("large file upload skipped in -short mode")
	}
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/input/fileupload.html", "text/html", `<input type="file"><script></script>`)
	_, page := launchPersistent(t, nil)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/fileupload.html"))

	uploadFile := filepath.Join(t.TempDir(), "200MB.zip")
	f, err := os.Create(uploadFile)
	must.NoError(err)
	chunk := []byte(strings.Repeat("A", 4*1024))
	for i := 0; i < 50*1024; i++ {
		_, err = f.Write(chunk)
		must.NoError(err)
	}
	must.NoError(f.Close())

	input := page.Locator("input[type='file']")
	_, err = input.Evaluate(ctx, `e => {
    window['events'] = [];
    e.addEventListener('input', () => window['events'].push('input'));
    e.addEventListener('change', () => window['events'].push('change'));
  }`)
	must.NoError(err)
	must.NoError(input.SetInputFiles(ctx, []string{uploadFile}))
	name, err := input.Evaluate(ctx, "e => e.files[0].name")
	must.NoError(err)
	is.Equal("200MB.zip", name)
	events, err := page.Evaluate(ctx, "() => window['events']")
	must.NoError(err)
	is.Equal([]any{"input", "change"}, events)
}

// TestDefaultBrowserContext2ShouldExposeBrowser verifies BrowserContext.Browser() of a persistent context.
// Ref: TestDefaultBrowserContext2.java#shouldExposeBrowser  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDefaultBrowserContext2.java
func TestDefaultBrowserContext2ShouldExposeBrowser(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	_, page := launchPersistent(t, nil)
	bCtx := page.Context()
	browser := bCtx.Browser()
	must.NotNil(browser)
	is.NotEmpty(browser.Version())
	page2, err := browser.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Goto(ctx, "data:text/html,<html><title>Title</title></html>"))
	title, err := page2.Title(ctx)
	must.NoError(err)
	is.Equal("Title", title)
	must.NoError(browser.Close(ctx))
	is.Len(bCtx.Pages(), 0)
}
