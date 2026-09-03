//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

// waitForPopup subscribes to page.OnPopup, runs trigger(), and returns the
// first popup page received within the test deadline.
func waitForPopup(t *testing.T, page *playwright.Page, trigger func()) *playwright.Page {
	t.Helper()
	popupCh := make(chan *playwright.Page, 1)
	cancel := page.OnPopup(func(p *playwright.Page) {
		select {
		case popupCh <- p:
		default:
		}
	})
	defer cancel()

	trigger()

	select {
	case popup := <-popupCh:
		return popup
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for popup")
		return nil
	}
}

// newContextWithCleanup creates a BrowserContext with the given options and
// registers cleanup. Equivalent to the per-test context pattern used across
// browser_context_test.go.
func newContextWithCleanup(t *testing.T, opts ...*playwright.BrowserContextOptions) *playwright.BrowserContext {
	t.Helper()
	ctx := testCtx(t)
	bc, err := globalBrowser.NewContext(ctx, opts...)
	require.NoError(t, err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(closeCtx)
	})
	return bc
}

// --- TestPageEventPopup tests ---

// TestPopupClickTargetBlank verifies that clicking a link with target="_blank" opens a popup,
// the opener is accessible from the popup, and the popup belongs to the correct frame.
// Ref: TestPageEventPopup.java#shouldWorkWithClickingTarget_blank
func TestPopupClickTargetBlank(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)
	srv.ServeWithBody("/one-style.html", "text/html", `<!DOCTYPE html><html><head><title>Style</title></head><body></body></html>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a target=_blank rel='opener' href='/one-style.html'>yo</a>`))

	popup := waitForPopup(t, page, func() {
		must.NoError(page.Locator("a").Click(ctx))
	})

	must.NotNil(popup)

	pageHasOpener, err := page.Evaluate(ctx, "() => !!window.opener")
	must.NoError(err)
	is.Equal(false, pageHasOpener, "page.opener should be null for the original tab")

	popupHasOpener, err := popup.Evaluate(ctx, "() => !!window.opener")
	must.NoError(err)
	is.Equal(true, popupHasOpener, "popup.opener should be non-null")

	// Verify the popup has a main frame (Go API does not expose Frame.Page())
	is.NotNil(popup.MainFrame(), "popup.MainFrame() should not be nil")
}

// --- TestPopup tests ---

// TestPopupInheritUserAgent verifies that popups inherit the userAgent from the BrowserContext.
// Ref: TestPopup.java#shouldInheritUserAgentFromBrowserContext
func TestPopupInheritUserAgent(t *testing.T) {
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

// TestPopupRespectRoutes verifies that routes registered on BrowserContext apply to popups
// opened via a link with target="_blank".
// Ref: TestPopup.java#shouldRespectRoutesFromBrowserContext
func TestPopupRespectRoutes(t *testing.T) {
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

// TestPopupInheritExtraHTTPHeaders verifies that extra HTTP headers from BrowserContext
// are sent with requests made from a popup opened via window.open().
// Ref: TestPopup.java#shouldInheritExtraHeadersFromBrowserContext
func TestPopupInheritExtraHTTPHeaders(t *testing.T) {
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

// TestPopupInheritOffline verifies that popups inherit the offline mode from BrowserContext.
// Ref: TestPopup.java#shouldInheritOfflineFromBrowserContext
func TestPopupInheritOffline(t *testing.T) {
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

// TestPopupInheritHttpCredentials verifies that HTTP credentials from BrowserContext
// are used for authentication in popups.
// Ref: TestPopup.java#shouldInheritHttpCredentialsFromBrowserContext
func TestPopupInheritHttpCredentials(t *testing.T) {
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

// TestPopupInheritViewportSize verifies that a popup opened via window.open('about:blank')
// inherits the viewport dimensions from the BrowserContext.
// Ref: TestPopup.java#shouldInheritViewportSizeFromBrowserContext
func TestPopupInheritViewportSize(t *testing.T) {
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

// TestPopupInheritHasTouch verifies that touch support is inherited by popups.
// Ref: TestPopup.java#shouldInheritTouchSupportFromBrowserContext
func TestPopupInheritHasTouch(t *testing.T) {
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

// TestPopupRespectRoutesWindowOpen verifies that context routes apply to popups opened
// via window.open().
// Ref: TestPopup.java#shouldRespectRoutesFromBrowserContextWithWindowOpen
func TestPopupRespectRoutesWindowOpen(t *testing.T) {
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

// TestPopupAddInitScriptInProcess verifies that init scripts registered on BrowserContext
// run in popups opened via window.open('about:blank') (in-process popup).
// Ref: TestPopup.java#BrowserContextAddInitScriptShouldApplyToAnInProcessPopup
func TestPopupAddInitScriptInProcess(t *testing.T) {
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

// TestPopupAddInitScriptCrossProcess verifies that init scripts registered on BrowserContext
// run in popups opened cross-process (window.open to a different origin).
// Ref: TestPopup.java#BrowserContextAddInitScriptShouldApplyToACrossProcessPopup
func TestPopupAddInitScriptCrossProcess(t *testing.T) {
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

// TestPopupExposeFunction verifies that functions exposed on BrowserContext are available
// in popups opened via window.open('about:blank').
// Ref: TestPopup.java#shouldExposeFunctionFromBrowserContext
func TestPopupExposeFunction(t *testing.T) {
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

// TestPopupOpenerIsPage verifies that OnPopup fires and popup.Opener() returns the
// originating page for a window.open() call.
// Ref: TestPopup.java / general OnPopup usage
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

// TestPopupOnPopupEvent verifies that page.OnPopup fires with the correct popup Page
// when window.open() is called.
// Ref: TestPageEventPopup.java / TestPopup.java general
func TestPopupOnPopupEvent(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	must := require.New(t)

	srv := testserver.New(t)
	srv.ServeWithBody("/popup-target.html", "text/html", `<title>Popup Target</title>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var receivedPopup *playwright.Page

	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		receivedPopup = p
		mu.Unlock()
	})
	defer cancel()

	_, err := page.Evaluate(ctx, "url => window.open(url)", srv.Prefix()+"/popup-target.html")
	must.NoError(err)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		p := receivedPopup
		mu.Unlock()
		if p != nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	mu.Lock()
	p := receivedPopup
	mu.Unlock()
	must.NotNil(p, "OnPopup should have fired with a non-nil page")
}

// TestPopupUseViewportSizeFromWindowFeatures verifies that window.open with explicit
// width/height features creates a popup with those dimensions, and that the popup
// viewport can subsequently be resized.
// Ref: TestPopup.java#shouldUseViewportSizeFromWindowFeatures
func TestPopupUseViewportSizeFromWindowFeatures(t *testing.T) {
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

// TestPopupWindowOpenerAccess verifies that window.opener is accessible from the popup
// when rel="opener" is set on the link.
// Ref: TestPageEventPopup.java#shouldWorkWithClickingTarget_blank
func TestPopupWindowOpenerAccess(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<button onclick="window.open('about:blank', '_blank')">open</button>`))

	popup := waitForPopup(t, page, func() {
		must.NoError(page.Locator("button").Click(ctx))
	})

	must.NotNil(popup)

	// The popup opened via window.open() should have window.opener set
	popupHasOpener, err := popup.Evaluate(ctx, "() => !!window.opener")
	must.NoError(err)
	is.Equal(true, popupHasOpener, "popup should have window.opener")
}

// TestPopupSharesBrowserContext verifies that a popup opened via window.open()
// belongs to the same BrowserContext as the originating page.
// Ref: TestPopup.java general / BrowserContext.OnPage
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
