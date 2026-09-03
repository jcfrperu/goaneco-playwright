//go:build e2e

// E2E tests for Page.AddInitScript, Page.BringToFront, Page.SetExtraHTTPHeaders,
// Page.SetViewportSize, Page.ViewportSize.
package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

// --- Page.AddInitScript ---

func TestPageAddInitScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/init-page", "text/html", `<html><body></body></html>`)

	page := newPage(t)

	err := page.AddInitScript(ctx, `window.__injected = 'hello-init'`)
	must.NoError(err, "AddInitScript failed")

	// AddInitScript runs before page scripts on navigation, not on SetContent.
	err = page.Goto(ctx, srv.Prefix()+"/init-page")
	must.NoError(err, "Goto failed")

	val, err := page.Evaluate(ctx, `window.__injected`)
	must.NoError(err, "Evaluate failed")
	is.Equal("hello-init", val)
}

func TestPageAddInitScriptRunsOnNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/init-check", "text/html", `<html><body></body></html>`)

	page := newPage(t)

	err := page.AddInitScript(ctx, `window.__nav = 'nav-value'`)
	must.NoError(err, "AddInitScript failed")

	err = page.Goto(ctx, srv.Prefix()+"/init-check")
	must.NoError(err, "Goto failed")

	val, err := page.Evaluate(ctx, `window.__nav`)
	must.NoError(err, "Evaluate failed")
	is.Equal("nav-value", val)
}

// --- Page.BringToFront ---

func TestPageBringToFront(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<html><body><p>front</p></body></html>`)
	must.NoError(err, "SetContent failed")

	// BringToFront should not error; behavioural effect (focus) is hard to assert in headless mode.
	err = page.BringToFront(ctx)
	must.NoError(err, "BringToFront failed")

	// Verify page is still functional after BringToFront.
	text, err := page.Locator("p").InnerText(ctx)
	must.NoError(err, "InnerText after BringToFront failed")
	is.Equal("front", text)
}

// --- Page.SetExtraHTTPHeaders ---

func TestPageSetExtraHTTPHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	gotHeader := make(chan string, 1)
	srv.SetRoute("/page-headers", func(w http.ResponseWriter, r *http.Request) {
		gotHeader <- r.Header.Get("X-Page-Header")
		w.WriteHeader(200)
	})

	page := newPage(t)
	err := page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Page-Header": "page-value",
	})
	must.NoError(err, "SetExtraHTTPHeaders failed")

	err = page.Goto(ctx, srv.Prefix()+"/page-headers")
	must.NoError(err, "Goto failed")

	select {
	case val := <-gotHeader:
		if val != "page-value" {
			t.Errorf("X-Page-Header = %q, want 'page-value'", val)
		}
	default:
		t.Error("request never reached the server")
	}
}

// --- Page.SetViewportSize / Page.ViewportSize ---

func TestPageSetViewportSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<html><body></body></html>`)
	must.NoError(err, "SetContent failed")

	err = page.SetViewportSize(ctx, 800, 600)
	must.NoError(err, "SetViewportSize failed")

	size, err := page.ViewportSize(ctx)
	must.NoError(err, "ViewportSize failed")
	if size.Width != 800 || size.Height != 600 {
		t.Errorf("ViewportSize = %dx%d, want 800x600", size.Width, size.Height)
	}
}

func TestPageSetViewportSizeChanges(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<html><body></body></html>`)
	must.NoError(err, "SetContent failed")

	err = page.SetViewportSize(ctx, 1280, 720)
	must.NoError(err, "SetViewportSize 1280x720 failed")

	err = page.SetViewportSize(ctx, 320, 568)
	must.NoError(err, "SetViewportSize 320x568 failed")

	size, err := page.ViewportSize(ctx)
	must.NoError(err, "ViewportSize failed")
	if size.Width != 320 || size.Height != 568 {
		t.Errorf("ViewportSize = %dx%d, want 320x568", size.Width, size.Height)
	}
}

func TestPageURLAfterSetContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>test</div>`))

	url := page.URL()
	is.NotEmpty(url)
}

func TestPageFramesOnlyMain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	frames := page.Frames()
	is.GreaterOrEqual(len(frames), 1)
}

func TestPageMainFrameIsNotNil(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	frame := page.MainFrame()
	must.NotNil(frame)
}

func TestPageFrameByNameNilForUnknown(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	frame := page.Frame("nonexistent")
	is.Nil(frame)
}

func TestPageContextReturnsContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	bc := page.Context()
	must.NotNil(bc)
}

func TestPageURLUpdatesAfterGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	url := page.URL()
	is.Equal(srv.EmptyPage(), url)
}
