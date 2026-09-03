//go:build e2e

// E2E tests for Page.WaitForEvent, Page.OnPopup, BrowserContext.WaitForEvent, Locator.Screenshot.
package e2e

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

// --- Page.OnPopup ---

func TestPageOnPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/popup-page", "text/html", `<title>Popup</title>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	popupCh := make(chan *playwright.Page, 1)
	cancel := page.OnPopup(func(p *playwright.Page) {
		select {
		case popupCh <- p:
		default:
		}
	})
	defer cancel()

	_, err = page.Evaluate(ctx, `() => window.open('`+srv.Prefix()+`/popup-page')`)
	must.NoError(err, "Evaluate failed")

	select {
	case popup := <-popupCh:
		if popup == nil {
			t.Fatal("OnPopup received nil page")
		}
		opener, err := popup.Opener(ctx)
		must.NoError(err, "popup.Opener() failed")
		if opener != page {
			t.Errorf("popup.Opener() = %v, want original page %v", opener, page)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for popup via OnPopup")
	}
}

// --- Page.WaitForEvent("popup") ---

func TestPageWaitForEventPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/we-popup", "text/html", `<title>WE Popup</title>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	done := make(chan any, 1)
	go func() {
		val, err := page.WaitForEvent(ctx, "popup")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	_, err = page.Evaluate(ctx, `() => window.open('`+srv.Prefix()+`/we-popup')`)
	must.NoError(err, "Evaluate failed")

	result := <-done
	if popup, ok := result.(*playwright.Page); ok {
		if popup == nil {
			t.Fatal("WaitForEvent(popup) returned nil page")
		}
	} else {
		t.Fatalf("WaitForEvent(popup) returned unexpected type/error: %v", result)
	}
}

// --- Page.WaitForEvent("request") ---

func TestPageWaitForEventRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/we-req", "text/html", "<p>ok</p>")

	page := newPage(t)

	done := make(chan any, 1)
	go func() {
		val, err := page.WaitForEvent(ctx, "request")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	err := page.Goto(ctx, srv.Prefix()+"/we-req")
	must.NoError(err, "Goto failed")

	result := <-done
	if req, ok := result.(*playwright.NetworkRequest); ok {
		if req.URL() == "" {
			t.Error("WaitForEvent(request) returned request with empty URL")
		}
	} else {
		t.Fatalf("WaitForEvent(request) unexpected result: %v", result)
	}
}

// --- BrowserContext.WaitForEvent("page") ---

func TestBrowserContextWaitForEventPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-we-page", "text/html", "<p>ctx</p>")

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	done := make(chan any, 1)
	go func() {
		val, err := bCtx.WaitForEvent(ctx, "page")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	// Open a new popup from the page to trigger a "page" event on the context
	_, err = page.Evaluate(ctx, `() => window.open('`+srv.Prefix()+`/ctx-we-page')`)
	must.NoError(err, "Evaluate failed")

	result := <-done
	if _, ok := result.(*playwright.Page); !ok {
		t.Fatalf("BrowserContext.WaitForEvent(page) unexpected result type: %T — %v", result, result)
	}
}

// --- BrowserContext.WaitForEvent("request") ---

func TestBrowserContextWaitForEventRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-we-req", "text/html", "<p>ok</p>")

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	done := make(chan any, 1)
	go func() {
		val, err := bCtx.WaitForEvent(ctx, "request")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	err = page.Goto(ctx, srv.Prefix()+"/ctx-we-req")
	must.NoError(err, "Goto failed")

	result := <-done
	if req, ok := result.(*playwright.NetworkRequest); ok {
		if req.URL() == "" {
			t.Error("BrowserContext.WaitForEvent(request) returned request with empty URL")
		}
	} else {
		t.Fatalf("BrowserContext.WaitForEvent(request) unexpected result: %v", result)
	}
}

// --- Locator.Screenshot ---

func TestLocatorScreenshot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="box" style="width:100px;height:100px;background:red"></div>`)
	must.NoError(err, "SetContent failed")

	data, err := page.Locator("#box").Screenshot(ctx)
	must.NoError(err, "Locator.Screenshot failed")
	if len(data) == 0 {
		t.Fatal("Locator.Screenshot returned empty bytes")
	}
	if !bytes.HasPrefix(data, []byte("\x89PNG")) {
		t.Errorf("expected PNG header, got: %x", data[:min(8, len(data))])
	}
}

func TestLocatorScreenshotJPEG(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="box" style="width:100px;height:100px;background:blue"></div>`)
	must.NoError(err, "SetContent failed")

	q := 80
	data, err := page.Locator("#box").Screenshot(ctx, &playwright.ScreenshotOptions{Type: "jpeg", Quality: &q})
	must.NoError(err, "Locator.Screenshot JPEG failed")
	if len(data) == 0 {
		t.Fatal("Locator.Screenshot JPEG returned empty bytes")
	}
	// JPEG magic: FF D8
	if data[0] != 0xFF || data[1] != 0xD8 {
		t.Errorf("expected JPEG header FF D8, got: %x %x", data[0], data[1])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
