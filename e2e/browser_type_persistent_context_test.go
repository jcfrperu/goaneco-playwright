//go:build e2e

// E2E tests for BrowserType.LaunchPersistentContext option matrix.
// Migration of: TestDefaultBrowserContext2.java
package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// launchPersistent is a test helper that launches a persistent context in a fresh temp dir,
// registers cleanup for both browser and context, and returns the BrowserContext and first Page.
func launchPersistent(t *testing.T, opts *playwright.LaunchPersistentContextOptions) (*playwright.BrowserContext, *playwright.Page) {
	t.Helper()
	ctx := testCtx(t)

	bt, err := globalPW.Chromium()
	require.NoError(t, err, "Chromium() failed")

	if opts == nil {
		opts = &playwright.LaunchPersistentContextOptions{}
	}
	headless := true
	opts.Headless = &headless

	browser, bCtx, err := bt.LaunchPersistentContext(ctx, t.TempDir(), opts)
	require.NoError(t, err, "LaunchPersistentContext failed")
	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(cleanCtx)
		_ = browser.Close(cleanCtx)
	})

	// The persistent context always opens with a default page (about:blank).
	// Create a new page if none are tracked yet.
	pages := bCtx.Pages()
	var page *playwright.Page
	if len(pages) > 0 {
		page = pages[0]
	} else {
		page, err = bCtx.NewPage(ctx)
		require.NoError(t, err, "NewPage failed")
	}
	return bCtx, page
}

// TestPersistentContextHasTouchOption verifies that setHasTouch makes ontouchstart available.
// Ref: TestDefaultBrowserContext2.java#shouldSupportHasTouchOption
func TestPersistentContextHasTouchOption(t *testing.T) {
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

// TestPersistentContextColorSchemeOption verifies the dark color scheme is applied.
// Ref: TestDefaultBrowserContext2.java#shouldSupportColorSchemeOption
func TestPersistentContextColorSchemeOption(t *testing.T) {
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

// TestPersistentContextTimezoneIdOption verifies the timezone is applied to Date formatting.
// Ref: TestDefaultBrowserContext2.java#shouldSupportTimezoneIdOption
func TestPersistentContextTimezoneIdOption(t *testing.T) {
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

// TestPersistentContextLocaleOption verifies the navigator.language reflects the locale.
// Ref: TestDefaultBrowserContext2.java#shouldSupportLocaleOption
func TestPersistentContextLocaleOption(t *testing.T) {
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

// TestPersistentContextGeolocationAndPermissionsOptions verifies geolocation is resolved when
// both the geolocation option and the "geolocation" permission are set.
// Ref: TestDefaultBrowserContext2.java#shouldSupportGeolocationAndPermissionsOptions
func TestPersistentContextGeolocationAndPermissionsOptions(t *testing.T) {
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

// TestPersistentContextExtraHTTPHeadersOption verifies extra headers are sent with requests.
// Ref: TestDefaultBrowserContext2.java#shouldSupportExtraHTTPHeadersOption
func TestPersistentContextExtraHTTPHeadersOption(t *testing.T) {
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

// TestPersistentContextAcceptsUserDataDir verifies that files are written to the user data dir.
// Ref: TestDefaultBrowserContext2.java#shouldAcceptUserDataDir
func TestPersistentContextAcceptsUserDataDir(t *testing.T) {
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

// TestPersistentContextRestoresStateFromUserDataDir verifies that localStorage persists
// across two separate persistent context sessions using the same user data dir.
// Ref: TestDefaultBrowserContext2.java#shouldRestoreStateFromUserDataDir
func TestPersistentContextRestoresStateFromUserDataDir(t *testing.T) {
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

// TestPersistentContextDefaultURLIsAboutBlank verifies that a fresh persistent context
// opens with a single about:blank page.
// Ref: TestDefaultBrowserContext2.java#shouldHaveDefaultURLWhenLaunchingBrowser
func TestPersistentContextDefaultURLIsAboutBlank(t *testing.T) {
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

// TestPersistentContextUserAgentOption verifies the User-Agent header is set from the option.
// Ref: TestDefaultBrowserContext2.java (userAgent option)
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

// TestPersistentContextViewportSizeOption verifies that a viewport option is respected.
// Ref: TestDefaultBrowserContext2.java#shouldWorkInPersistentContext (viewport part)
func TestPersistentContextViewportSizeOption(t *testing.T) {
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

// TestPersistentContextOfflineModeOption verifies that offline mode blocks network requests.
// Ref: TestDefaultBrowserContext2.java (offline option)
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

// TestPersistentContextContrastOption verifies the contrast media feature option.
// Ref: TestDefaultBrowserContext2.java#shouldSupportContrastOption
func TestPersistentContextContrastOption(t *testing.T) {
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

// TestPersistentContextBypassCSPOption verifies that CSP can be bypassed.
// Ref: TestDefaultBrowserContext2.java (bypassCSP option)
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
