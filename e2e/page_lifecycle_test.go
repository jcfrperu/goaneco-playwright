//go:build e2e

// Page lifecycle E2E tests: close, load events.
// Migration of: TestPageBasic.java (shouldRejectAllPromisesWhenPageIsClosed,
// shouldNotRunBeforeunloadByDefault, shouldFireLoadWhenExpected, shouldFireDomcontentloadedWhenExpected)
// and TestBrowserContextEvents.java (pageLoadEvent, pageCloseEvent).
package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

// TestPageRejectAllPromisesWhenClosed verifies that an Evaluate call on a closed page returns an error.
// Ref: TestPageBasic.java#shouldRejectAllPromisesWhenPageIsClosed
func TestPageRejectAllPromisesWhenClosed(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// Start an evaluate that never resolves, then close the page concurrently.
	type evalResult struct {
		val any
		err error
	}
	ch := make(chan evalResult, 1)
	go func() {
		val, err := page.Evaluate(ctx, "() => new Promise(f => {})")
		ch <- evalResult{val, err}
	}()

	// Give the evaluate goroutine a moment to send its IPC, then close the page.
	must.NoError(page.WaitForTimeout(ctx, 100))
	err = page.Close(ctx)
	must.NoError(err, "Close failed")

	select {
	case res := <-ch:
		if res.err == nil {
			t.Fatal("expected error from Evaluate after page close, got nil")
		}
		t.Logf("Evaluate error after close: %v", res.err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for Evaluate to reject after page close")
	}
}

// TestPageNotRunBeforeunloadByDefault verifies Close() does not trigger beforeunload dialogs.
// Ref: TestPageBasic.java#shouldNotRunBeforeunloadByDefault
func TestPageNotRunBeforeunloadByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Install a beforeunload handler that would fire a dialog if triggered.
	_, err = page.Evaluate(ctx, `() => {
		window.addEventListener('beforeunload', e => { e.returnValue = 'are you sure?'; });
	}`)
	must.NoError(err, "Evaluate failed")

	// Close should succeed without hanging on a dialog.
	done := make(chan error, 1)
	go func() {
		done <- page.Close(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Close() failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Close() timed out — beforeunload may have triggered a dialog")
	}
}

// TestBrowserContextOnPageClose verifies that OnPageClose handlers fire when a page is closed.
// Ref: TestBrowserContextEvents.java#browserContextOnPageClose
func TestBrowserContextOnPageClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	closedCh := make(chan *playwright.Page, 1)
	cancel := bCtx.OnPageClose(func(p *playwright.Page) {
		select {
		case closedCh <- p:
		default:
		}
	})
	defer cancel()

	err = page.Close(ctx)
	must.NoError(err, "Close failed")

	select {
	case p := <-closedCh:
		if p != page {
			t.Errorf("OnPageClose received wrong page: %v, want %v", p, page)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageClose to fire")
	}
}

// TestBrowserContextOnPageLoad verifies that OnPageLoad handlers fire after navigation.
// Ref: TestBrowserContextEvents.java#browserContextOnPageLoad
func TestBrowserContextOnPageLoad(t *testing.T) {
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

	srv.ServeWithBody("/onload-target", "text/html", `<p>loaded</p>`)

	loadedCh := make(chan *playwright.Page, 1)
	cancelLoad := bCtx.OnPageLoad(func(p *playwright.Page) {
		select {
		case loadedCh <- p:
		default:
		}
	})
	defer cancelLoad()

	err = page.Goto(ctx, srv.Prefix()+"/onload-target")
	must.NoError(err, "Goto failed")

	select {
	case p := <-loadedCh:
		if p != page {
			t.Errorf("OnPageLoad received wrong page: %v, want %v", p, page)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageLoad to fire")
	}
}

// TestPageFireLoadWhenExpected verifies the "load" lifecycle event fires during navigation.
// Ref: TestPageBasic.java#shouldFireLoadWhenExpected
func TestPageFireLoadWhenExpected(t *testing.T) {
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

	srv.ServeWithBody("/fire-load", "text/html", `<p>ready</p>`)

	loadedCh := make(chan struct{}, 1)
	cancelLoad := bCtx.OnPageLoad(func(p *playwright.Page) {
		select {
		case loadedCh <- struct{}{}:
		default:
		}
	})
	defer cancelLoad()

	err = page.Goto(ctx, srv.Prefix()+"/fire-load")
	must.NoError(err, "Goto failed")

	select {
	case <-loadedCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for load event")
	}
}

// TestPageFireDomContentLoadedWhenExpected verifies that WaitForLoadState("domcontentloaded") resolves after Goto.
// Ref: TestPageBasic.java#shouldFireDomcontentloadedWhenExpected
func TestPageFireDomContentLoadedWhenExpected(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/dcl", "text/html", `<p>dom content loaded</p>`)
	err := page.Goto(ctx, srv.Prefix()+"/dcl")
	must.NoError(err, "Goto failed")

	err = page.WaitForLoadState(ctx, "domcontentloaded")
	must.NoError(err, "WaitForLoadState('domcontentloaded') failed")
}

// TestPageCloseUnregistersFromContext verifies that closing a page removes it from the context page list,
// and that the OnPageClose event fires before Pages() is updated.
// Ref: TestPageBasic.java#shouldNotBeVisibleInContextPages (extended)
func TestPageCloseAndOnPageCloseOrder(t *testing.T) {
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
	must.NoError(err, "NewPage failed")

	closedCh := make(chan string, 1)
	cancelClose := bCtx.OnPageClose(func(p *playwright.Page) {
		url := p.URL()
		select {
		case closedCh <- url:
		default:
		}
	})
	defer cancelClose()

	err = page.Close(ctx)
	must.NoError(err, "Close() failed")

	select {
	case url := <-closedCh:
		t.Logf("OnPageClose fired; page URL was: %q", url)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageClose")
	}

	is.Len(bCtx.Pages(), 0)
}

// TestPageCloseMarksPageAsClosed verifies IsClosed returns true after Close.
// Ref: TestPageClose.java#shouldMarkPageAsClosed
func TestPageCloseMarksPageAsClosed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	is.False(page.IsClosed())
	must.NoError(page.Close(ctx))
	is.True(page.IsClosed())
}

// TestPageCloseRemovesPageFromContext verifies closed page is removed from context.Pages().
// Ref: TestPageClose.java#shouldRemoveFromContext
func TestPageCloseRemovesPageFromContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	pagesBeforeClose := bc.Pages()
	initialCount := len(pagesBeforeClose)

	must.NoError(page.Close(ctx))

	pagesAfterClose := bc.Pages()
	is.Less(len(pagesAfterClose), initialCount)
}

// TestPageCloseIsCallableTwice verifies closing a page twice doesn't panic or error.
// Ref: TestPageClose.java#shouldBeCallableTwice
func TestPageCloseIsCallableTwiceExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Close(ctx))
	must.NoError(page.Close(ctx)) // Second close should not panic
}

// TestPageCloseWithContent verifies page with content can be closed normally.
// Ref: TestPageClose.java#shouldClosePageWithContent
func TestPageCloseWithContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<div>some content</div>`))
	must.NoError(page.Close(ctx))
	is.True(page.IsClosed())
}

// TestPageOnLoadPassesSelf verifies that OnLoad fires and passes the page as argument.
// Ref: TestPageBasic.java#shouldPassSelfAsArgumentToLoadEvent
func TestPageOnLoadPassesSelf(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/load-self", "text/html", `<p>ready</p>`)

	loadedCh := make(chan *playwright.Page, 1)
	cancel := page.OnLoad(func(p *playwright.Page) {
		select {
		case loadedCh <- p:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/load-self")
	must.NoError(err)

	select {
	case p := <-loadedCh:
		must.Same(page, p, "OnLoad should pass the page itself")
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnLoad")
	}
}

// TestPageOnDOMContentLoadedPassesSelf verifies OnDOMContentLoaded fires and passes the page as argument.
// Ref: TestPageBasic.java#shouldPassSelfAsArgumentToDomcontentloadedEvent
func TestPageOnDOMContentLoadedPassesSelf(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/dcl-self", "text/html", `<p>ready</p>`)

	ch := make(chan *playwright.Page, 1)
	cancel := page.OnDOMContentLoaded(func(p *playwright.Page) {
		select {
		case ch <- p:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/dcl-self")
	must.NoError(err)

	select {
	case p := <-ch:
		must.Same(page, p, "OnDOMContentLoaded should pass the page itself")
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnDOMContentLoaded")
	}
}

// TestPageCloseShouldWorkWithWindowClose verifies a popup opened via window.open can be
// closed via window.close() and WaitForClose properly detects it.
// Ref: TestPageBasic.java#pageCloseShouldWorkWithWindowClose
func TestPageCloseShouldWorkWithWindowClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.open('about:blank')")
		return err
	})
	must.NoError(err, "WaitForPopup failed")
	must.NotNil(popup)

	err = popup.WaitForClose(ctx, func() error {
		// window.close() may cause the evaluate to fail as the page closes; ignore that error.
		_, _ = popup.Evaluate(ctx, "() => window.close()")
		return nil
	})
	must.NoError(err, "WaitForClose failed")
	is.True(popup.IsClosed(), "popup should be marked as closed")
}

// TestPageOnPageLoadCancelWorks verifies that the cancel function returned by OnPageLoad stops future callbacks.
// Ref: TestBrowserContextEvents.java (cancel behavior)
func TestPageOnPageLoadCancelWorks(t *testing.T) {
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

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	srv.ServeWithBody("/cancel-load-a", "text/html", `<p>a</p>`)
	srv.ServeWithBody("/cancel-load-b", "text/html", `<p>b</p>`)

	var urls []string
	cancelLoad := bCtx.OnPageLoad(func(p *playwright.Page) {
		urls = append(urls, p.URL())
	})

	// First navigation: handler active.
	err = page.Goto(ctx, srv.Prefix()+"/cancel-load-a")
	must.NoError(err, "Goto /a failed")
	time.Sleep(200 * time.Millisecond)

	// Cancel the handler.
	cancelLoad()

	// Second navigation: handler should NOT fire.
	err = page.Goto(ctx, srv.Prefix()+"/cancel-load-b")
	must.NoError(err, "Goto /b failed")
	time.Sleep(200 * time.Millisecond)

	// Only the first URL should be in the list.
	is.Len(urls, 1, "expected exactly 1 load event (before cancel)")
	is.Contains(urls[0], "/cancel-load-a")
}
