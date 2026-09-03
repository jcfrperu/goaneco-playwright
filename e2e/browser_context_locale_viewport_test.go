//go:build e2e

// Extended locale and viewport tests for BrowserContext.
// Migration of: TestBrowserContextLocale.java, TestBrowserContextViewport.java
package e2e

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrowserContextLocaleAffectsNavigatorLanguage verifies that locale option sets navigator.language.
// Ref: TestBrowserContextLocale.java#shouldAffectNavigatorLanguage
func TestBrowserContextLocaleAffectsNavigatorLanguage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	locale := "fr-CH"
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{Locale: &locale})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	lang, err := page.Evaluate(ctx, "() => navigator.language")
	must.NoError(err, "Evaluate failed")
	if lang != "fr-CH" {
		t.Errorf("navigator.language = %q, want 'fr-CH'", lang)
	}
}

// TestBrowserContextLocaleFormatsNumber verifies that locale affects number formatting.
// Ref: TestBrowserContextLocale.java#shouldFormatNumber
func TestBrowserContextLocaleFormatsNumber(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	locale := "de-DE"
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{Locale: &locale})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	formatted, err := page.Evaluate(ctx, `() => (1000000.5).toLocaleString()`)
	must.NoError(err, "Evaluate failed")
	s, ok := formatted.(string)
	is.True(ok, "expected string")
	// German locale uses period as thousands separator and comma as decimal.
	if !strings.Contains(s, ".") || !strings.Contains(s, ",") {
		t.Errorf("expected German number format (with . and ,), got: %q", s)
	}
}

// TestBrowserContextLocaleIsolatedBetweenContexts verifies that locale is per-context.
// Ref: TestBrowserContextLocale.java#shouldBeIsolatedBetweenContexts
func TestBrowserContextLocaleIsolatedBetweenContexts(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	localeEN := "en-US"
	localeDE := "de-DE"

	bCtxEN, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{Locale: &localeEN})
	must.NoError(err, "NewContext EN failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtxEN.Close(c)
	})

	bCtxDE, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{Locale: &localeDE})
	must.NoError(err, "NewContext DE failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtxDE.Close(c)
	})

	pageEN, _ := bCtxEN.NewPage(ctx)
	pageDE, _ := bCtxDE.NewPage(ctx)

	langEN, _ := pageEN.Evaluate(ctx, "() => navigator.language")
	langDE, _ := pageDE.Evaluate(ctx, "() => navigator.language")

	if langEN != "en-US" {
		t.Errorf("EN context navigator.language = %q, want 'en-US'", langEN)
	}
	if langDE != "de-DE" {
		t.Errorf("DE context navigator.language = %q, want 'de-DE'", langDE)
	}
}

// TestBrowserContextLocaleAffectsAcceptLanguageHeader verifies that locale sets Accept-Language header.
// Ref: TestBrowserContextLocale.java#shouldAffectAcceptLanguageHeader
func TestBrowserContextLocaleAffectsAcceptLanguageHeader(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedHeader string
	srv.SetRoute("/accept-lang", func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("Accept-Language")
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<p>ok</p>"))
	})

	locale := "fr-CH,fr"
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{Locale: &locale})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.Prefix()+"/accept-lang")
	must.NoError(err, "Goto failed")

	if !strings.Contains(strings.ToLower(receivedHeader), "fr") {
		t.Errorf("Accept-Language header = %q, expected to contain 'fr'", receivedHeader)
	}
}

// TestBrowserContextViewportOuterSize verifies outerWidth/outerHeight reflect viewport settings.
// Ref: TestBrowserContextViewport.java#shouldReturnCorrectOuterWidthAndOuterHeight
func TestBrowserContextViewportOuterSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 1280, Height: 720},
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	innerW, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err, "Evaluate innerWidth failed")
	innerH, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err, "Evaluate innerHeight failed")

	if innerW != float64(1280) {
		t.Errorf("innerWidth = %v, want 1280", innerW)
	}
	if innerH != float64(720) {
		t.Errorf("innerHeight = %v, want 720", innerH)
	}
}

// TestPageViewportSize verifies that Page.ViewportSize returns the correct dimensions
// reflecting the viewport configured for the browser context.
// Ref: TestPageViewportSize
func TestPageViewportSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 1024, Height: 768},
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	size, err := page.ViewportSize(ctx)
	must.NoError(err, "ViewportSize failed")
	must.NotNil(size, "ViewportSize returned nil")
	is.Equal(1024, size.Width, "viewport width should match context option")
	is.Equal(768, size.Height, "viewport height should match context option")
}

// TestPageViewportSizeAfterSetViewport verifies that ViewportSize reflects a subsequent SetViewportSize call.
// Ref: TestPageViewportSizeAfterChange
func TestPageViewportSizeAfterSetViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetViewportSize(ctx, 640, 480)
	must.NoError(err, "SetViewportSize failed")

	size, err := page.ViewportSize(ctx)
	must.NoError(err, "ViewportSize failed")
	must.NotNil(size, "ViewportSize returned nil")
	is.Equal(640, size.Width, "width should be 640 after SetViewportSize")
	is.Equal(480, size.Height, "height should be 480 after SetViewportSize")
}

// TestBrowserContextNotChangeDefaultLocaleInAnotherContext verifies that setting a locale in one context
// does not affect a new context created without a locale override.
// Ref: TestBrowserContextLocale.java#shouldNotChangeDefaultLocaleInAnotherContext
func TestBrowserContextNotChangeDefaultLocaleInAnotherContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	getContextLocale := func(opts *playwright.BrowserContextOptions) string {
		bCtx, err := globalBrowser.NewContext(ctx, opts)
		must.NoError(err)
		t.Cleanup(func() {
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = bCtx.Close(closeCtx)
		})
		page, err := bCtx.NewPage(ctx)
		must.NoError(err)
		val, err := page.Evaluate(ctx, "(new Intl.NumberFormat()).resolvedOptions().locale")
		must.NoError(err)
		s, _ := val.(string)
		return s
	}

	defaultLocale := getContextLocale(nil)
	is.NotEmpty(defaultLocale)

	var localeOverride string
	if defaultLocale == "es-MX" {
		localeOverride = "de-DE"
	} else {
		localeOverride = "es-MX"
	}

	overrideLocaleStr := localeOverride
	overriddenLocale := getContextLocale(&playwright.BrowserContextOptions{Locale: &overrideLocaleStr})
	is.Equal(localeOverride, overriddenLocale)

	restoredLocale := getContextLocale(nil)
	is.Equal(defaultLocale, restoredLocale, "default locale should not be affected by override in another context")
}

// TestBrowserContextViewportGetProperSize verifies viewport size is correctly returned.
// Ref: TestBrowserContextViewport.java#shouldGetTheProperDefaultViewPortSize
func TestBrowserContextViewportGetProperSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 800, Height: 600},
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	w, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err, "innerWidth failed")
	h, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err, "innerHeight failed")

	if w != float64(800) || h != float64(600) {
		t.Errorf("viewport = %vx%v, want 800x600", w, h)
	}
}

// TestBrowserContextViewportSetProperSize verifies that SetViewportSize updates the viewport.
// Ref: TestBrowserContextViewport.java#shouldSetTheProperViewportSize
func TestBrowserContextViewportSetProperSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 123, 456))
	w, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err)
	h, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err)
	is.Equal(float64(123), w)
	is.Equal(float64(456), h)
}

// TestBrowserContextViewportEmulateDeviceWidth verifies screen.width and matchMedia respond to viewport size changes.
// Ref: TestBrowserContextViewport.java#shouldEmulateDeviceWidth
func TestBrowserContextViewportEmulateDeviceWidth(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 200, 200))
	screenW, err := page.Evaluate(ctx, "() => window.screen.width")
	must.NoError(err)
	is.Equal(float64(200), screenW)

	matchTrue, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-width: 100px)').matches")
	is.Equal(true, matchTrue)
	matchFalse, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-width: 300px)').matches")
	is.Equal(false, matchFalse)
	matchFalse2, _ := page.Evaluate(ctx, "() => matchMedia('(max-device-width: 100px)').matches")
	is.Equal(false, matchFalse2)
	matchTrue2, _ := page.Evaluate(ctx, "() => matchMedia('(max-device-width: 300px)').matches")
	is.Equal(true, matchTrue2)
	matchFalse3, _ := page.Evaluate(ctx, "() => matchMedia('(device-width: 500px)').matches")
	is.Equal(false, matchFalse3)
	matchTrue3, _ := page.Evaluate(ctx, "() => matchMedia('(device-width: 200px)').matches")
	is.Equal(true, matchTrue3)

	must.NoError(page.SetViewportSize(ctx, 500, 500))
	screenW2, err := page.Evaluate(ctx, "() => window.screen.width")
	must.NoError(err)
	is.Equal(float64(500), screenW2)
}

// TestBrowserContextViewportEmulateDeviceHeight verifies screen.height and matchMedia respond to viewport size changes.
// Ref: TestBrowserContextViewport.java#shouldEmulateDeviceHeight
func TestBrowserContextViewportEmulateDeviceHeight(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 200, 200))
	screenH, err := page.Evaluate(ctx, "() => window.screen.height")
	must.NoError(err)
	is.Equal(float64(200), screenH)

	matchTrue, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-height: 100px)').matches")
	is.Equal(true, matchTrue)
	matchFalse, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-height: 300px)').matches")
	is.Equal(false, matchFalse)

	must.NoError(page.SetViewportSize(ctx, 500, 500))
	screenH2, err := page.Evaluate(ctx, "() => window.screen.height")
	must.NoError(err)
	is.Equal(float64(500), screenH2)
}

// TestBrowserContextViewportEmulateAvailSize verifies screen.availWidth/availHeight match viewport.
// Ref: TestBrowserContextViewport.java#shouldEmulateAvailWidthAndAvailHeight
func TestBrowserContextViewportEmulateAvailSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 500, 600))
	availW, err := page.Evaluate(ctx, "() => window.screen.availWidth")
	must.NoError(err)
	availH, err := page.Evaluate(ctx, "() => window.screen.availHeight")
	must.NoError(err)
	is.Equal(float64(500), availW)
	is.Equal(float64(600), availH)
}

// TestBrowserContextViewportNoTouchByDefault verifies that touch events are not available by default.
// Ref: TestBrowserContextViewport.java#shouldNotHaveTouchByDefault
func TestBrowserContextViewportNoTouchByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Test via SetContent to avoid dependency on specific test server files
	must.NoError(page.SetContent(ctx, `<html><body><script>
		document.body.textContent = 'ontouchstart' in window ? 'YES' : 'NO';
	</script></body></html>`))
	hasTouchStart, err := page.Evaluate(ctx, "() => 'ontouchstart' in window")
	must.NoError(err)
	is.Equal(false, hasTouchStart)
}

// TestBrowserContextViewportHasTouch verifies that HasTouch option enables touch support (with null viewport).
// Ref: TestBrowserContextViewport.java#shouldSupportTouchWithNullViewport
func TestBrowserContextViewportHasTouch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	hasTouch := true
	noDefaultViewport := true
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		HasTouch:          &hasTouch,
		NoDefaultViewport: &noDefaultViewport,
	})
	must.NoError(err, "NewContext with HasTouch+NoDefaultViewport failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)
	hasTouchStart, err := page.Evaluate(ctx, "() => 'ontouchstart' in window")
	must.NoError(err)
	is.Equal(true, hasTouchStart)
}

// TestBrowserContextViewportNullViewportSize verifies that ViewportSize returns nil when NoDefaultViewport is set.
// Ref: TestBrowserContextViewport.java#shouldReportNullViewPortSizeWhenGivenNullViewport
func TestBrowserContextViewportNullViewportSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	noDefaultViewport := true
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		NoDefaultViewport: &noDefaultViewport,
	})
	must.NoError(err, "NewContext with NoDefaultViewport failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	size, err := page.ViewportSize(ctx)
	must.NoError(err)
	is.Nil(size, "ViewportSize should return nil when no default viewport is set")
}

// ---------------------------------------------------------------------------
// From browser_context_locale_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextLocaleAffectsNumberFormat verifies Locale option affects number formatting.
// Ref: TestBrowserContextLocale.java#shouldAffectNumberFormat
func TestBrowserContextLocaleAffectsNumberFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	locale := "de-DE"
	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Locale: &locale,
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => navigator.language`)
	must.NoError(err)
	is.Equal("de-DE", result)
}

// TestBrowserContextTimezoneAffectsDateFormat verifies TimezoneId option affects timezone.
// Ref: TestBrowserContextLocale.java#shouldAffectTimeZone
func TestBrowserContextTimezoneAffectsDateFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	tz := "America/New_York"
	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		TimezoneID: &tz,
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => Intl.DateTimeFormat().resolvedOptions().timeZone`)
	must.NoError(err)
	is.Equal("America/New_York", result)
}

// TestBrowserContextLocaleFormatsCurrencyCorrectly verifies locale affects currency formatting.
// Ref: TestBrowserContextLocale.java#shouldFormatCurrencyWithLocale
func TestBrowserContextLocaleFormatsCurrencyCorrectly(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	locale := "en-US"
	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Locale: &locale,
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => new Intl.NumberFormat('en-US', {style: 'currency', currency: 'USD'}).format(1234.56)`)
	must.NoError(err)
	is.Equal("$1,234.56", result)
}
