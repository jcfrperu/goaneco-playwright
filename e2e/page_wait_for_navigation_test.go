//go:build e2e

// E2E tests for: TestPageWaitForNavigation.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

// Ref: TestPageWaitForNavigation.java#shouldWorkWithDOMHistoryBackHistoryForward  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithDOMHistoryBackHistoryForward(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a id=back onclick="javascript:goBack()">back</a>
<a id=forward onclick="javascript:goForward()">forward</a>
<script>
  function goBack() { history.back(); }
  function goForward() { history.forward(); }
  history.pushState({}, '', '/first.html');
  history.pushState({}, '', '/second.html');
</script>`))

	is.Equal(srv.Prefix()+"/second.html", page.URL())

	must.NoError(page.Locator("a#back").Click(ctx))
	err := page.WaitForURL(ctx, srv.Prefix()+"/first.html", 3*time.Second)
	must.NoError(err)
	is.Equal(srv.Prefix()+"/first.html", page.URL())

	must.NoError(page.Locator("a#forward").Click(ctx))
	err = page.WaitForURL(ctx, srv.Prefix()+"/second.html", 3*time.Second)
	must.NoError(err)
	is.Equal(srv.Prefix()+"/second.html", page.URL())
}

// TestPageWaitForNavigationShouldRespectTimeout verifies WaitForURL times out when URL never matches.
// Ref: TestPageWaitForNavigation.java#shouldRespectTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	err := page.WaitForURL(ctx, "http://never-matches.invalid/", 500*time.Millisecond)
	is.Error(err)
	is.Contains(err.Error(), "timeout")
}

// TestPageWaitForNavigationShouldWorkWithClickingOnAnchorLinks verifies URL settles after a click-triggered navigation.
// Ref: TestPageWaitForNavigation.java#shouldWorkWithClickingOnAnchorLinks  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithClickingOnAnchorLinks(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/linked", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<p>linked page</p>`)
	}))
	srv.SetRoute("/link-page", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<a id="link" href="/linked">Go</a>`)
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/link-page"))
	must.NoError(page.Locator("#link").Click(ctx))
	must.NoError(page.WaitForURL(ctx, srv.Prefix()+"/linked"))
	is.Equal(srv.Prefix()+"/linked", page.URL())
}

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestWaitForLoadStateAfterSetContentExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))
	must.NoError(page.WaitForLoadState(ctx, "load"))
}

// TestFrameWaitForLoadState verifies Frame.WaitForLoadState.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestWaitForLoadStateDOMContentLoadedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Content</p>`))
	must.NoError(page.WaitForLoadState(ctx, "domcontentloaded"))
}

// TestWaitForLoadStateLoadEx4 verifies WaitForLoadState load.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestWaitForLoadStateLoadEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Loaded</p>`))
	must.NoError(page.WaitForLoadState(ctx, "load"))
}

// TestWaitForLoadStateNetworkIdleEx4 verifies WaitForLoadState networkidle.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestWaitForLoadStateNetworkIdleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>No network</p>`))
	must.NoError(page.WaitForLoadState(ctx, "networkidle"))
}

// TestWaitForLoadStateDefaultEx4 verifies WaitForLoadState defaults to load.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestWaitForLoadStateDefaultEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Page</p>`))
	must.NoError(page.WaitForLoadState(ctx))
}

// TestWaitForLoadStateAfterSetContentEx5 verifies WaitForLoadState after SetContent.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForNavigation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
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

// Ref: TestPageWaitForNavigation.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWork(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkWithCommit  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithCommit(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkWithClickingOnLinksWhichDoNotCommitNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithClickingOnLinksWhichDoNotCommitNavigation(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkWithHistoryPushState  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithHistoryPushState(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkWithHistoryReplaceState  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithHistoryReplaceState(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkWhenSubframeIssuesWindowStop  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWhenSubframeIssuesWindowStop(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkWithUrlMatch  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithUrlMatch(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkWithUrlMatchForSameDocumentNavigations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkWithUrlMatchForSameDocumentNavigations(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkForCrossProcessNavigations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkForCrossProcessNavigations(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldWorkOnFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldWorkOnFrame(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldFailWhenFrameDetaches  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldFailWhenFrameDetaches(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldThrowOnInvalidUrlMatcherTypeInPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldThrowOnInvalidUrlMatcherTypeInPage(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}

// Ref: TestPageWaitForNavigation.java#shouldThrowOnInvalidUrlMatcherTypeInFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForNavigation.java
func TestPageWaitForNavigationShouldThrowOnInvalidUrlMatcherTypeInFrame(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Frame.WaitForNavigation")
}
