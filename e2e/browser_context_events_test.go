//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrowserContextOnRequestFiredForPage verifies BrowserContext.OnRequest fires for page navigations.
// Ref: TestBrowserContextEvents.java#shouldFireOnRequest
func TestBrowserContextOnRequestFiredForPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var urls []string

	cancel := bc.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(urls)
	is.Contains(urls[0], "empty.html")
}

// TestBrowserContextOnResponseFiredForPage verifies BrowserContext.OnResponse fires for page navigations.
// Ref: TestBrowserContextEvents.java#shouldFireOnResponse
func TestBrowserContextOnResponseFiredForPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var statuses []int

	cancel := bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statuses = append(statuses, resp.Status())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(statuses)
	is.Equal(200, statuses[0])
}

// TestBrowserContextOnRequestFinishedFiredForPage verifies BrowserContext.OnRequestFinished fires.
// Ref: TestBrowserContextEvents.java#shouldFireOnRequestFinished
func TestBrowserContextOnRequestFinishedFiredForPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	done := make(chan struct{}, 1)
	var mu sync.Mutex
	var gotURL string

	cancel := bc.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			gotURL = req.URL()
			mu.Unlock()
			select {
			case done <- struct{}{}:
			default:
			}
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnRequestFinished")
	}

	mu.Lock()
	defer mu.Unlock()
	is.Equal(srv.EmptyPage(), gotURL)
}

// TestBrowserContextOnPageCloseFiresOnPageClose verifies BrowserContext.OnPageClose fires when page closes.
// Ref: TestBrowserContextEvents.java#shouldFireOnPageClose
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

// TestBrowserContextOnPageLoadFiresAfterNavigation verifies BrowserContext.OnPageLoad fires.
// Ref: TestBrowserContextEvents.java#shouldFireOnPageLoad
func TestBrowserContextOnPageLoadFiresAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	done := make(chan struct{}, 1)

	cancel := bc.OnPageLoad(func(p *playwright.Page) {
		select {
		case done <- struct{}{}:
		default:
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageLoad")
	}
}

// TestBrowserContextOnRequestFailedFired verifies BrowserContext.OnRequestFailed fires when a request fails.
// Ref: TestBrowserContextEvents.java#shouldFireOnRequestFailed
func TestBrowserContextOnRequestFailedFired(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	// Serve a page that loads a script that we abort via route.
	srv.ServeWithBody("/fail-page", "text/html",
		`<script src="/fail.js"></script>`)
	srv.SetRoute("/fail.js", func(w http.ResponseWriter, r *http.Request) {
		// This handler will never be reached because the route aborts the request.
		w.WriteHeader(200)
	})

	// Abort the script request via route to guarantee requestfailed fires.
	must.NoError(bc.Route(ctx, "**/fail.js", func(route *playwright.Route) {
		_ = route.Abort(ctx, "failed")
	}))

	done := make(chan struct{}, 1)
	var mu sync.Mutex
	var failedURL string

	cancel := bc.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		if failedURL == "" {
			failedURL = req.URL()
		}
		mu.Unlock()
		select {
		case done <- struct{}{}:
		default:
		}
	})
	defer cancel()

	_ = page.Goto(ctx, srv.Prefix()+"/fail-page")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("OnRequestFailed did not fire within 5 seconds")
	}

	mu.Lock()
	u := failedURL
	mu.Unlock()
	is.Contains(u, "/fail.js")
}

// TestBrowserContextOnPageFired verifies BrowserContext.OnPage fires when a new page is created.
// Ref: TestBrowserContextEvents.java#shouldFireOnPage
func TestBrowserContextOnPageFired(t *testing.T) {
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

	done := make(chan *playwright.Page, 1)
	cancel := bc.OnPage(func(p *playwright.Page) {
		select {
		case done <- p:
		default:
		}
	})
	defer cancel()

	// Open a popup — this creates a new page in the same context.
	_, err = page.Evaluate(ctx, "() => window.open('about:blank')")
	must.NoError(err)

	select {
	case newPage := <-done:
		must.NotNil(newPage)
		is.Equal(bc, newPage.Context())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPage event")
	}
}

// TestBrowserContextOnRequestFiredEx verifies OnRequest fires for navigation.
// Ref: TestBrowserContextEvents.java#shouldFireOnRequest
func TestBrowserContextOnRequestFiredEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/req-page", "text/html", `<html><body>Request page</body></html>`)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(c)
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var captured string
	off := bc.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		if captured == "" {
			captured = req.URL()
		}
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/req-page"))

	mu.Lock()
	u := captured
	mu.Unlock()

	is.NotEmpty(u)
	is.Contains(u, "/req-page")
}

// TestBrowserContextOnRequestFailedEx verifies OnRequestFailed fires for aborted request.
// Ref: TestBrowserContextEvents.java#shouldFireOnRequestFailed
func TestBrowserContextOnRequestFailedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(c)
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	failCount := 0
	off := bc.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		failCount++
	})
	t.Cleanup(off)

	// Abort all requests to a specific path
	must.NoError(bc.Route(ctx, "**/fail-this", func(r *playwright.Route) {
		_ = r.Abort(ctx)
	}))

	_, _ = page.Evaluate(ctx, `async () => { try { await fetch('/fail-this'); } catch(e) {} }`)

	mu.Lock()
	n := failCount
	mu.Unlock()

	is.Greater(n, 0)
}

// TestBrowserContextOnResponseFiredEx verifies OnResponse fires after request.
// Ref: TestBrowserContextEvents.java#shouldFireOnResponseForNav
func TestBrowserContextOnResponseFiredEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/resp-page", "text/html", `<html><body>Response page</body></html>`)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(c)
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var respURL string
	off := bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		defer mu.Unlock()
		if respURL == "" {
			respURL = resp.URL()
		}
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/resp-page"))

	mu.Lock()
	u := respURL
	mu.Unlock()

	is.NotEmpty(u)
	is.Contains(u, "/resp-page")
}

// TestBrowserContextConsoleEventInPopup verifies that console messages from a popup page
// are captured by a context-level console handler.
// Ref: TestBrowserContextEvents.java#consoleEventShouldWorkInPopup
func TestBrowserContextConsoleEventInPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	msgCh := make(chan *playwright.ConsoleMessage, 1)
	cancelConsole := bc.OnConsole(func(m *playwright.ConsoleMessage) {
		if m.Text() == "hello" {
			select {
			case msgCh <- m:
			default:
			}
		}
	})
	defer cancelConsole()

	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "const win = window.open(''); win.console.log('hello');")
		return err
	})
	must.NoError(err)
	must.NotNil(popup)

	select {
	case msg := <-msgCh:
		is.Equal("hello", msg.Text())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for console message from popup")
	}
}

// TestBrowserContextConsoleEventInPopup2 verifies console messages from a javascript: URL popup.
// Ref: TestBrowserContextEvents.java#consoleEventShouldWorkInPopup2
func TestBrowserContextConsoleEventInPopup2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	msgCh := make(chan *playwright.ConsoleMessage, 1)
	cancelConsole := bc.OnConsole(func(m *playwright.ConsoleMessage) {
		if m.Type() == "log" && m.Text() == "hello" {
			select {
			case msgCh <- m:
			default:
			}
		}
	})
	defer cancelConsole()

	popupCh := make(chan *playwright.Page, 1)
	cancelPage := bc.OnPage(func(pg *playwright.Page) {
		select {
		case popupCh <- pg:
		default:
		}
	})
	defer cancelPage()

	_, err = page.Evaluate(ctx, `async () => {
		const win = window.open('javascript:console.log("hello")');
		await new Promise(f => setTimeout(f, 0));
		win.close();
	}`)
	must.NoError(err)

	select {
	case <-popupCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for popup page")
	}

	select {
	case msg := <-msgCh:
		is.Equal("hello", msg.Text())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for console message from popup2")
	}
}

// TestBrowserContextConsoleEventInImmediatelyClosedPopup verifies console messages
// from a popup that opens and immediately closes.
// Ref: TestBrowserContextEvents.java#consoleEventShouldWorkInImmediatelyClosedPopup
func TestBrowserContextConsoleEventInImmediatelyClosedPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	msgCh := make(chan *playwright.ConsoleMessage, 1)
	cancelConsole := bc.OnConsole(func(m *playwright.ConsoleMessage) {
		if m.Text() == "hello" {
			select {
			case msgCh <- m:
			default:
			}
		}
	})
	defer cancelConsole()

	_, err = page.Evaluate(ctx, `async () => {
		const win = window.open();
		win.console.log('hello');
		win.close();
	}`)
	must.NoError(err)

	select {
	case msg := <-msgCh:
		is.Equal("hello", msg.Text())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for console message from immediately-closed popup")
	}
}

// TestBrowserContextDialogEventInPopup verifies that dialog events from a popup page
// are captured and handled by a page-level dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkInPopup
func TestBrowserContextDialogEventInPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := page.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	resultCh := make(chan any, 1)
	popup, err := page.WaitForPopup(ctx, func() error {
		go func() {
			result, _ := page.Evaluate(ctx, "() => { const win = window.open(''); return win.prompt('hey?'); }")
			select {
			case resultCh <- result:
			default:
			}
		}()
		return nil
	})
	must.NoError(err)
	must.NotNil(popup)

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from popup")
	}

	select {
	case result := <-resultCh:
		is.Equal("hello", result)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for evaluate result")
	}
}

// TestBrowserContextDialogEventInPopup2 verifies that a dialog opened via javascript: URL
// is captured by a context-level dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkInPopup2
func TestBrowserContextDialogEventInPopup2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := page.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	_, err = page.Evaluate(ctx, `window.open('javascript:prompt("hey?")');`)
	must.NoError(err)

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from javascript: URL popup")
	}
}

// TestBrowserContextDialogEventInImmediatelyClosedPopup verifies that a dialog from a popup
// that immediately closes is captured by a page-level dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkInImmdiatelyClosedPopup
func TestBrowserContextDialogEventInImmediatelyClosedPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := page.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	resultCh := make(chan any, 1)
	popup, err := page.WaitForPopup(ctx, func() error {
		go func() {
			result, _ := page.Evaluate(ctx, `async () => {
				const win = window.open();
				const result = win.prompt('hey?');
				win.close();
				return result;
			}`)
			select {
			case resultCh <- result:
			default:
			}
		}()
		return nil
	})
	must.NoError(err)
	must.NotNil(popup)

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from immediately-closed popup")
	}

	select {
	case result := <-resultCh:
		is.Equal("hello", result)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for evaluate result")
	}
}

// TestBrowserContextDialogEventWithInlineScriptTag verifies that a dialog triggered by an
// inline <script> in a popup loaded via link click is captured by a context dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkWithInlineScriptTag
func TestBrowserContextDialogEventWithInlineScriptTag(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	srv.SetRoute("/popup.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`<script>window.result = prompt('hey?')</script>`))
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, "<a href='"+srv.Prefix()+"/popup.html' target=_blank>Click me</a>"))

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := page.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	popupCh := make(chan *playwright.Page, 1)
	cancelPage := bc.OnPage(func(pg *playwright.Page) {
		select {
		case popupCh <- pg:
		default:
		}
	})
	defer cancelPage()

	must.NoError(page.Locator("a").Click(ctx))

	var popup *playwright.Page
	select {
	case popup = <-popupCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for popup page")
	}

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from inline script")
	}

	// Wait for window.result to be set to "hello" in popup
	is.Eventually(func() bool {
		result, err := popup.Evaluate(ctx, "window.result")
		return err == nil && result == "hello"
	}, 5*time.Second, 100*time.Millisecond, "timed out waiting for window.result to be 'hello'")
}
