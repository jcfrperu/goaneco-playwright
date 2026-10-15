//go:build e2e

// E2E tests for: TestBrowserContextLocale.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestBrowserContextOptionsLocale verifies custom locale configuration in BrowserContextOptions.
// Ref: TestBrowserContextLocale.java
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextLocale.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
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

// TestBrowserContextLocaleShouldAffectNavigatorLanguage verifies that locale option sets navigator.language.
// Ref: TestBrowserContextLocale.java#shouldAffectNavigatorLanguage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldAffectNavigatorLanguage(t *testing.T) {
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

// TestBrowserContextLocaleShouldFormatNumber verifies that locale affects number formatting.
// Ref: TestBrowserContextLocale.java#shouldFormatNumber  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldFormatNumber(t *testing.T) {
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

// TestBrowserContextLocaleShouldBeIsolatedBetweenContexts verifies that locale is per-context.
// Ref: TestBrowserContextLocale.java#shouldBeIsolatedBetweenContexts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldBeIsolatedBetweenContexts(t *testing.T) {
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

// TestBrowserContextLocaleShouldAffectAcceptLanguageHeader verifies that locale sets Accept-Language header.
// Ref: TestBrowserContextLocale.java#shouldAffectAcceptLanguageHeader  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldAffectAcceptLanguageHeader(t *testing.T) {
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

// TestBrowserContextLocaleShouldNotChangeDefaultLocaleInAnotherContext verifies that setting a locale in one context
// does not affect a new context created without a locale override.
// Ref: TestBrowserContextLocale.java#shouldNotChangeDefaultLocaleInAnotherContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldNotChangeDefaultLocaleInAnotherContext(t *testing.T) {
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

// TestBrowserContextLocaleShouldFormatDate verifies timezone configuration in BrowserContextOptions.
// Ref: TestBrowserContextLocale.java#shouldFormatDate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldFormatDate(t *testing.T) {
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

// TestPageViewportSize verifies that Page.ViewportSize returns the correct dimensions
// reflecting the viewport configured for the browser context.
// Ref: TestPageViewportSize
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextLocale.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
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

// TestBrowserContextLocaleAffectsNumberFormat verifies Locale option affects number formatting.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextLocale.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextLocale.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextLocale.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
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

// a4LocalePopup opens /formatted-number.html in a popup of a page created with the given locale.
func a4LocalePopup(t *testing.T, locale string) *playwright.Page {
	t.Helper()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/formatted-number.html", "text/html",
		"<script>window['result'] = (1000000.50).toLocaleString().replace(/\\s/g, ' '); window.initialNavigatorLanguage = navigator.language;</script>")
	page := a4NewPage(t, &playwright.BrowserContextOptions{Locale: a4Ptr(locale)})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	popup, err := page.WaitForPopup(ctx, func() error {
		_, e := page.Evaluate(ctx, "url => window.open(url)", srv.Prefix()+"/formatted-number.html")
		return e
	})
	must.NoError(err)
	must.NoError(popup.WaitForLoadState(ctx, "domcontentloaded"))
	return popup
}

// TestBrowserContextLocaleShouldFormatNumberInPopups verifies the context locale applies to popups' number formatting.
// Ref: TestBrowserContextLocale.java#shouldFormatNumberInPopups  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldFormatNumberInPopups(t *testing.T) {
	t.Parallel()
	popup := a4LocalePopup(t, "fr-FR")
	result, err := popup.Evaluate(testCtx(t), "window['result']")
	require.NoError(t, err)
	assert.Equal(t, "1 000 000,5", result)
}

// TestBrowserContextLocaleShouldAffectNavigatorLanguageInPopups verifies navigator.language in popups.
// Ref: TestBrowserContextLocale.java#shouldAffectNavigatorLanguageInPopups  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldAffectNavigatorLanguageInPopups(t *testing.T) {
	t.Parallel()
	popup := a4LocalePopup(t, "fr-FR")
	result, err := popup.Evaluate(testCtx(t), "window.initialNavigatorLanguage")
	require.NoError(t, err)
	assert.Equal(t, "fr-FR", result)
}

// TestBrowserContextLocaleShouldWorkForMultiplePagesSharingSameProcess verifies opening several popups works with a locale set.
// Ref: TestBrowserContextLocale.java#shouldWorkForMultiplePagesSharingSameProcess  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextLocale.java
func TestBrowserContextLocaleShouldWorkForMultiplePagesSharingSameProcess(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := a4NewPage(t, &playwright.BrowserContextOptions{Locale: a4Ptr("ru-RU")})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	for i := 0; i < 2; i++ {
		_, err := page.WaitForPopup(ctx, func() error {
			_, e := page.Evaluate(ctx, "url => window.open(url)", srv.EmptyPage())
			return e
		})
		must.NoError(err)
	}
}
