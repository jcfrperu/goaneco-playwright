//go:build e2e

// E2E tests for: TestPopup.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

// TestPopupOpenerIsPage verifies that OnPopup fires and popup.Opener() returns the
// originating page for a window.open() call.
// Ref: TestPopup.java / general OnPopup usage
// Extra: no equivalent method in playwright-java; related class: TestPopup.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupOpenerIsPage(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup := waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, "url => window.open(url)", srv.EmptyPage())
		must.NoError(evalErr)
	})

	must.NotNil(popup)
	opener, err := popup.Opener(ctx)
	must.NoError(err)
	is.Equal(page, opener, "popup.Opener() should return the page that opened it")
}

// TestPopupSharesBrowserContext verifies that a popup opened via window.open()
// belongs to the same BrowserContext as the originating page.
// Ref: TestPopup.java general / BrowserContext.OnPage
// Extra: no equivalent method in playwright-java; related class: TestPopup.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupSharesBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup := waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, "url => window.open(url)", srv.EmptyPage())
		must.NoError(evalErr)
	})

	must.NotNil(popup)
	is.Equal(bc, popup.Context(), "popup should share the parent page's BrowserContext")
}

// TestPopupURL verifies that the popup opens with the correct URL.
// Ref: TestPopup.java general navigation
// Extra: no equivalent method in playwright-java; related class: TestPopup.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupURL(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)
	srv.ServeWithBody("/popup-url.html", "text/html", `<title>Popup URL</title>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	targetURL := srv.Prefix() + "/popup-url.html"
	popup := waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, "url => window.open(url)", targetURL)
		must.NoError(evalErr)
	})

	must.NoError(popup.WaitForLoadState(ctx, "domcontentloaded"))

	popupURL := popup.URL()
	is.Equal(targetURL, popupURL, "popup URL should match the target URL")
}

// TestPopupShouldInheritUserAgentFromBrowserContext verifies that popups inherit the userAgent from the BrowserContext.
// Ref: TestPopup.java#shouldInheritUserAgentFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldInheritUserAgentFromBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	userAgent := "hey"
	bc := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		UserAgent: &userAgent,
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a target=_blank rel=noopener href='/popup-ua.html'>link</a>`))

	var receivedUA string
	srv.ServeWithBody("/popup-ua.html", "text/html", `<!DOCTYPE html><html><head></head><body>
<script>window['initialUserAgent'] = navigator.userAgent;</script>
</body></html>`)

	popup := waitForPopup(t, page, func() {
		must.NoError(page.Locator("a").Click(ctx))
	})

	must.NoError(popup.WaitForLoadState(ctx, "domcontentloaded"))

	uaRaw, err := popup.Evaluate(ctx, "() => window['initialUserAgent']")
	must.NoError(err)
	var ok bool
	receivedUA, ok = uaRaw.(string)
	must.True(ok, "expected string user agent")
	is.Equal("hey", receivedUA)
}

// TestPopupShouldRespectRoutesFromBrowserContext verifies that routes registered on BrowserContext apply to popups
// opened via a link with target="_blank".
// Ref: TestPopup.java#shouldRespectRoutesFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldRespectRoutesFromBrowserContext(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a target=_blank rel=noopener href='`+srv.EmptyPage()+`'>link</a>`))

	var intercepted atomic.Bool
	must.NoError(bc.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		intercepted.Store(true)
		must.NoError(route.Continue(ctx, nil))
	}))

	waitForPopup(t, page, func() {
		must.NoError(page.Locator("a").Click(ctx))
	})

	is.True(intercepted.Load(), "context route should have intercepted the popup request")
}

// TestPopupShouldInheritExtraHeadersFromBrowserContext verifies that extra HTTP headers from BrowserContext
// are sent with requests made from a popup opened via window.open().
// Ref: TestPopup.java#shouldInheritExtraHeadersFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldInheritExtraHeadersFromBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	// Capture the header before any navigation so the route is ready.
	headerCh := make(chan string, 1)
	srv.SetRoute("/dummy-extra.html", func(w http.ResponseWriter, r *http.Request) {
		select {
		case headerCh <- r.Header.Get("foo"):
		default:
		}
		w.WriteHeader(http.StatusOK)
	})

	bc := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		ExtraHTTPHeaders: map[string]string{"foo": "bar"},
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, "url => { window['_popup'] = window.open(url) }", srv.Prefix()+"/dummy-extra.html")
		must.NoError(evalErr)
	})

	select {
	case fooHeader := <-headerCh:
		is.Equal("bar", fooHeader, "popup request should carry extra HTTP header")
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for /dummy-extra.html request with extra header")
	}
}

// TestPopupShouldInheritOfflineFromBrowserContext verifies that popups inherit the offline mode from BrowserContext.
// Ref: TestPopup.java#shouldInheritOfflineFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldInheritOfflineFromBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(bc.SetOffline(ctx, true))

	onlineRaw, err := page.Evaluate(ctx, `url => {
		const win = window.open(url);
		return win.navigator.onLine;
	}`, srv.Prefix()+"/dummy.html")
	must.NoError(err)

	online, ok := onlineRaw.(bool)
	must.True(ok, "expected bool for navigator.onLine")
	is.False(online, "popup in offline context should report navigator.onLine = false")
}

// TestPopupShouldInheritHttpCredentialsFromBrowserContext verifies that HTTP credentials from BrowserContext
// are used for authentication in popups.
// Ref: TestPopup.java#shouldInheritHttpCredentialsFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldInheritHttpCredentialsFromBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/title-auth.html", "user", "pass", "text/html",
		`<!DOCTYPE html><html><head><title>Woof-Woof</title></head><body></body></html>`)

	bc := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
		},
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup := waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, "url => { window['_popup'] = window.open(url) }", srv.Prefix()+"/title-auth.html")
		must.NoError(evalErr)
	})

	must.NoError(popup.WaitForLoadState(ctx, "domcontentloaded"))

	title, err := popup.Title(ctx)
	must.NoError(err)
	is.Equal("Woof-Woof", title, "popup should load authenticated page with inherited credentials")
}

// TestPopupShouldInheritViewportSizeFromBrowserContext verifies that a popup opened via window.open('about:blank')
// inherits the viewport dimensions from the BrowserContext.
// Ref: TestPopup.java#shouldInheritViewportSizeFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldInheritViewportSizeFromBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 400, Height: 500},
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	sizeRaw, err := page.Evaluate(ctx, `() => {
		const win = window.open('about:blank');
		return { width: win.innerWidth, height: win.innerHeight };
	}`)
	must.NoError(err)

	sizeMap, ok := sizeRaw.(map[string]interface{})
	must.True(ok, "expected map for window size")
	is.Equal(float64(400), sizeMap["width"], "popup width should match context viewport")
	is.Equal(float64(500), sizeMap["height"], "popup height should match context viewport")
}

// TestPopupShouldInheritTouchSupportFromBrowserContext verifies that touch support is inherited by popups.
// Ref: TestPopup.java#shouldInheritTouchSupportFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldInheritTouchSupportFromBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	hasTouch := true
	bc := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 400, Height: 500},
		HasTouch: &hasTouch,
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	hasTouchRaw, err := page.Evaluate(ctx, `() => {
		const win = window.open('');
		return 'ontouchstart' in win;
	}`)
	must.NoError(err)

	is.Equal(true, hasTouchRaw, "popup should inherit touch support from browser context")
}

// TestPopupShouldRespectRoutesFromBrowserContextWithWindowOpen verifies that context routes apply to popups opened
// via window.open().
// Ref: TestPopup.java#shouldRespectRoutesFromBrowserContextWithWindowOpen  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldRespectRoutesFromBrowserContextWithWindowOpen(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var intercepted atomic.Bool
	must.NoError(bc.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		intercepted.Store(true)
		must.NoError(route.Continue(ctx, nil))
	}))

	waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, "url => { window['__popup'] = window.open(url) }", srv.EmptyPage())
		must.NoError(evalErr)
	})

	is.True(intercepted.Load(), "context route should intercept popup's request")
}

// TestPopupBrowserContextAddInitScriptShouldApplyToAnInProcessPopup verifies that init scripts registered on BrowserContext
// run in popups opened via window.open('about:blank') (in-process popup).
// Ref: TestPopup.java#BrowserContextAddInitScriptShouldApplyToAnInProcessPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupBrowserContextAddInitScriptShouldApplyToAnInProcessPopup(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t)
	must.NoError(bc.AddInitScript(ctx, "window['injected'] = 123"))

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	injectedRaw, err := page.Evaluate(ctx, `() => {
		const win = window.open('about:blank');
		return win['injected'];
	}`)
	must.NoError(err)
	is.Equal(float64(123), injectedRaw, "in-process popup should have init script applied")
}

// TestPopupBrowserContextAddInitScriptShouldApplyToACrossProcessPopup verifies that init scripts registered on BrowserContext
// run in popups opened cross-process (window.open to a different origin).
// Ref: TestPopup.java#BrowserContextAddInitScriptShouldApplyToACrossProcessPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupBrowserContextAddInitScriptShouldApplyToACrossProcessPopup(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)
	srv.ServeWithBody("/title-xp.html", "text/html", `<!DOCTYPE html><html><head><title>Cross</title></head><body></body></html>`)

	// Use cross-process URL to force a new renderer process
	crossURL := srv.CrossProcessPrefix() + "/title-xp.html"

	bc := newContextWithCleanup(t)
	must.NoError(bc.AddInitScript(ctx, "(() => { window['injected'] = 123 })()"))

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup := waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, "url => window.open(url)", crossURL)
		must.NoError(evalErr)
	})

	injectedRaw, err := popup.Evaluate(ctx, "injected")
	must.NoError(err)
	is.Equal(float64(123), injectedRaw, "cross-process popup should have init script injected")

	// Verify init script survives reload
	must.NoError(popup.Reload(ctx))

	injectedAfterReload, err := popup.Evaluate(ctx, "injected")
	must.NoError(err)
	is.Equal(float64(123), injectedAfterReload, "init script should survive popup reload")
}

// TestPopupShouldExposeFunctionFromBrowserContext verifies that functions exposed on BrowserContext are available
// in popups opened via window.open('about:blank').
// Ref: TestPopup.java#shouldExposeFunctionFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldExposeFunctionFromBrowserContext(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t)
	must.NoError(bc.ExposeBinding(ctx, "add", func(args ...any) any {
		if len(args) < 2 {
			return 0
		}
		a, _ := args[0].(float64)
		b, _ := args[1].(float64)
		return a + b
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	addedRaw, err := page.Evaluate(ctx, `async () => {
		const win = window.open('about:blank');
		return win['add'](9, 4);
	}`)
	must.NoError(err)
	is.Equal(float64(13), addedRaw, "exposed function should be callable from popup")
}

// TestPopupShouldUseViewportSizeFromWindowFeatures verifies that window.open with explicit
// width/height features creates a popup with those dimensions, and that the popup
// viewport can subsequently be resized.
// Ref: TestPopup.java#shouldUseViewportSizeFromWindowFeatures  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPopup.java
func TestPopupShouldUseViewportSizeFromWindowFeatures(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	bc := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 700, Height: 700},
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup := waitForPopup(t, page, func() {
		_, evalErr := page.Evaluate(ctx, `async () => {
			const win = window.open(
				window.location.href,
				'Title',
				'toolbar=no,location=no,directories=no,status=no,menubar=no,scrollbars=yes,resizable=yes,width=600,height=300,top=0,left=0'
			);
			await new Promise(resolve => {
				const interval = setInterval(() => {
					if (win.innerWidth === 600 && win.innerHeight === 300) {
						clearInterval(interval);
						resolve();
					}
				}, 10);
			});
			return { width: win.innerWidth, height: win.innerHeight };
		}`)
		must.NoError(evalErr)
	})

	must.NoError(popup.SetViewportSize(ctx, 500, 400))
	must.NoError(popup.WaitForLoadState(ctx))

	resizedRaw, err := popup.Evaluate(ctx, "() => ({ width: window.innerWidth, height: window.innerHeight })")
	must.NoError(err)

	resizedMap, ok := resizedRaw.(map[string]interface{})
	must.True(ok, "expected map for resized window size")
	is.Equal(float64(500), resizedMap["width"], "popup width after SetViewportSize")
	is.Equal(float64(400), resizedMap["height"], "popup height after SetViewportSize")
}
