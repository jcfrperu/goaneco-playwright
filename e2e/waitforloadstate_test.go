//go:build e2e

package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageWaitForLoadState(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/load", "text/html", `<h1>loaded</h1>`)
	err := page.Goto(ctx, srv.Prefix()+"/load")
	must.NoError(err, "Goto failed")

	// After Goto completes, the page should be fully loaded.
	// WaitForLoadState("load") should return immediately.
	err = page.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState('load') failed")
}

func TestPageWaitForLoadStateDOMContentLoaded(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/dom", "text/html", `<p>content</p>`)
	err := page.Goto(ctx, srv.Prefix()+"/dom")
	must.NoError(err, "Goto failed")

	err = page.WaitForLoadState(ctx, "domcontentloaded")
	must.NoError(err, "WaitForLoadState('domcontentloaded') failed")
}

func TestPageWaitForLoadStateDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/default", "text/html", `<p>default</p>`)
	err := page.Goto(ctx, srv.Prefix()+"/default")
	must.NoError(err, "Goto failed")

	// Without explicit state, defaults to "load".
	err = page.WaitForLoadState(ctx)
	must.NoError(err, "WaitForLoadState() default failed")
}

// TestWaitForLoadStateAfterSetContent verifies WaitForLoadState resolves after SetContent.
// Ref: TestPageWaitForLoadState.java#shouldWorkAfterSetContent
func TestWaitForLoadStateAfterSetContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))
	must.NoError(page.WaitForLoadState(ctx, "domcontentloaded"))

	count, err := page.Locator("div").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestWaitForLoadStateNetworkIdle verifies WaitForLoadState with networkidle resolves.
// Ref: TestPageWaitForLoadState.java#shouldWorkWithNetworkIdle
func TestWaitForLoadStateNetworkIdle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx, "networkidle"))

	// Page should be in a stable state
	state, err := page.Evaluate(ctx, `() => document.readyState`)
	must.NoError(err)
	is.Equal("complete", state)
}

// TestWaitForLoadStateLoadResolvesAfterGoto verifies WaitForLoadState with load resolves after navigation.
// Ref: TestPageWaitForLoadState.java#shouldWorkWithLoad
func TestWaitForLoadStateLoadResolvesAfterGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	// Should not block since load is already done after Goto
	must.NoError(page.WaitForLoadState(ctx, "load"))
}

// TestWaitForLoadStateDefaultIsLoad verifies WaitForLoadState with no args uses load state.
// Ref: TestPageWaitForLoadState.java#defaultIsLoad
func TestWaitForLoadStateDefaultIsLoad(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx))
}

// TestWaitForLoadStateNetworkIdleAfterGoto verifies networkidle works after navigation.
// Ref: TestPageWaitForLoadState.java#shouldWaitNetworkIdleAfterGoto
func TestWaitForLoadStateNetworkIdleAfterGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx, "networkidle"))
}

// TestWaitForLoadStateDOMContentLoaded verifies domcontentloaded state.
// Ref: TestPageWaitForLoadState.java#shouldWaitDOMContentLoaded
func TestWaitForLoadStateDOMContentLoaded(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx, "domcontentloaded"))
}

// TestWaitForLoadStateNoArgDefaultsToLoad verifies no arg defaults to load.
// Ref: TestPageWaitForLoadState.java#shouldDefaultToLoad
func TestWaitForLoadStateNoArgDefaultsToLoad(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx))
}

// TestWaitForLoadStateAfterSetContentExtra verifies works after SetContent.
// Ref: TestPageWaitForLoadState.java#shouldWorkAfterSetContent
func TestWaitForLoadStateAfterSetContentExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))
	must.NoError(page.WaitForLoadState(ctx, "load"))
}

// TestFrameWaitForLoadState verifies Frame.WaitForLoadState.
// Ref: TestPageWaitForLoadState.java#shouldWorkOnFrame
func TestFrameWaitForLoadState(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.MainFrame().WaitForLoadState(ctx, "load"))
}

// TestWaitForLoadStateWithHTMLPageEx3 verifies WaitForLoadState with full HTML page.
// Ref: TestPageWaitForLoadState.java#shouldWorkWithHTMLPage
func TestWaitForLoadStateWithHTMLPageEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/page", "text/html", `<html><body><h1>Page</h1></body></html>`)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/page"))
	must.NoError(page.WaitForLoadState(ctx, "load"))
}

// TestWaitForLoadStateAfterReloadEx3 verifies WaitForLoadState works after page reload.
// Ref: TestPageWaitForLoadState.java#shouldWorkAfterReload
func TestWaitForLoadStateAfterReloadEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Reload(ctx))
	must.NoError(page.WaitForLoadState(ctx, "domcontentloaded"))
}

// TestWaitForLoadStateMultipleCallsEx3 verifies multiple WaitForLoadState calls do not error.
// Ref: TestPageWaitForLoadState.java#shouldAllowMultipleCalls
func TestWaitForLoadStateMultipleCallsEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx, "load"))
	must.NoError(page.WaitForLoadState(ctx, "domcontentloaded"))
	must.NoError(page.WaitForLoadState(ctx))
}

// TestWaitForLoadStateScriptInjectionEx3 verifies load state after script injection via SetContent.
// Ref: TestPageWaitForLoadState.java#shouldWorkWithScriptContent
func TestWaitForLoadStateScriptInjectionEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<html><head><script>window.__loaded = true;</script></head><body></body></html>
	`))
	must.NoError(page.WaitForLoadState(ctx, "load"))

	result, err := page.Evaluate(ctx, `() => window.__loaded`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestWaitForLoadStateDOMContentLoadedEx4 verifies WaitForLoadState domcontentloaded.
// Ref: TestPageWaitForLoadState.java#shouldWaitForDOMContentLoaded
func TestWaitForLoadStateDOMContentLoadedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Content</p>`))
	must.NoError(page.WaitForLoadState(ctx, "domcontentloaded"))
}

// TestWaitForLoadStateLoadEx4 verifies WaitForLoadState load.
// Ref: TestPageWaitForLoadState.java#shouldWaitForLoad
func TestWaitForLoadStateLoadEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Loaded</p>`))
	must.NoError(page.WaitForLoadState(ctx, "load"))
}

// TestWaitForLoadStateNetworkIdleEx4 verifies WaitForLoadState networkidle.
// Ref: TestPageWaitForLoadState.java#shouldWaitForNetworkIdle
func TestWaitForLoadStateNetworkIdleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>No network</p>`))
	must.NoError(page.WaitForLoadState(ctx, "networkidle"))
}

// TestWaitForLoadStateDefaultEx4 verifies WaitForLoadState defaults to load.
// Ref: TestPageWaitForLoadState.java#shouldDefaultToLoad
func TestWaitForLoadStateDefaultEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Page</p>`))
	must.NoError(page.WaitForLoadState(ctx))
}

// TestWaitForLoadStateAfterSetContentEx5 verifies WaitForLoadState after SetContent.
// Ref: TestPageWaitForLoadState.java#shouldWaitAfterSetContent
func TestWaitForLoadStateAfterSetContentEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Hello</div>`))
	must.NoError(page.WaitForLoadState(ctx))

	text, err := page.Locator("div").TextContent(ctx)
	must.NoError(err)
	is.Equal("Hello", text)
}

// TestWaitForLoadStateDomContentLoadedEx5 verifies WaitForLoadState with domcontentloaded.
// Ref: TestPageWaitForLoadState.java#shouldWaitForDOMContentLoaded
func TestWaitForLoadStateDomContentLoadedEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><p>Content</p></body></html>`))
	must.NoError(page.WaitForLoadState(ctx, "domcontentloaded"))

	count, err := page.Locator("p").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestWaitForLoadStateNetworkIdleEx5 verifies WaitForLoadState with networkidle.
// Ref: TestPageWaitForLoadState.java#shouldWaitForNetworkIdle
func TestWaitForLoadStateNetworkIdleEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body>Static</body></html>`))
	must.NoError(page.WaitForLoadState(ctx, "networkidle"))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "Static")
}
