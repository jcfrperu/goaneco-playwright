//go:build e2e

// Page.WaitForURL E2E tests.
// Migration of: TestPageWaitForURL.java (additional cases)
package e2e

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWaitForURLExactMatch verifies WaitForURL resolves when URL matches exactly.
// Ref: TestPageWaitForURL.java#shouldWaitForExactURL
func TestWaitForURLExactMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, srv.EmptyPage()))

	is.Equal(srv.EmptyPage(), page.URL())
}

// TestWaitForURLWithWildcard verifies WaitForURL with glob pattern.
// Ref: TestPageWaitForURL.java#shouldWaitForURLWithWildcard
func TestWaitForURLWithWildcard(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, "**/empty.html"))

	is.Contains(page.URL(), "empty.html")
}

// TestWaitForURLAlreadyAtURL verifies WaitForURL resolves immediately if already at URL.
// Ref: TestPageWaitForURL.java#shouldResolveImmediatelyIfAlreadyAtURL
func TestWaitForURLAlreadyAtURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	// Already at the URL, should resolve immediately
	must.NoError(page.WaitForURL(ctx, srv.EmptyPage()))
}

// TestWaitForURLAfterHashNavigation verifies WaitForURL works for hash navigation.
// Ref: TestPageWaitForURL.java#shouldWorkWithHashNavigation
func TestWaitForURLAfterHashNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	// Navigate to hash
	_, err := page.Evaluate(ctx, `() => { window.location.hash = '#section'; }`)
	must.NoError(err)

	must.NoError(page.WaitForURL(ctx, "*#section"))
	is.Contains(page.URL(), "#section")
}

func TestWaitForURLRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	err := page.WaitForURL(ctx, srv.Prefix()+"/frame.html", 500*time.Millisecond)
	is.Error(err, "expected timeout error from WaitForURL")
}

func TestWaitForURLWithAnchorLinks(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a href="#foobar">foobar</a>`))

	must.NoError(page.Locator("a").Click(ctx))

	err := page.WaitForURL(ctx, srv.EmptyPage()+"*", 3*time.Second)
	must.NoError(err)
	is.Contains(page.URL(), "#foobar")
}

func TestWaitForURLWithHistoryPushState(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a onclick="javascript:pushState()">SPA</a>
<script>
  function pushState() { history.pushState({}, '', 'wow.html') }
</script>`))

	must.NoError(page.Locator("a").Click(ctx))

	err := page.WaitForURL(ctx, srv.Prefix()+"/wow.html", 3*time.Second)
	must.NoError(err)
	is.Equal(srv.Prefix()+"/wow.html", page.URL())
}

func TestWaitForURLWithHistoryReplaceState(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a onclick="javascript:replaceState()">SPA</a>
<script>
  function replaceState() { history.replaceState({}, '', '/replaced.html') }
</script>`))

	must.NoError(page.Locator("a").Click(ctx))

	err := page.WaitForURL(ctx, srv.Prefix()+"/replaced.html", 3*time.Second)
	must.NoError(err)
	is.Equal(srv.Prefix()+"/replaced.html", page.URL())
}

func TestWaitForURLWithDOMHistoryBackForward(t *testing.T) {
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

func TestWaitForURLExactMatchAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, srv.EmptyPage()))
}

func TestWaitForURLWildcardPatternExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, "**/empty*"))
}

func TestWaitForURLWithTimeoutSucceeds(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, srv.EmptyPage(), 5*time.Second))
}

func TestWaitForURLAfterHashChangeExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, `() => window.location.hash = '#section'`)
	must.NoError(err)

	must.NoError(page.WaitForURL(ctx, "*#section"))
}

func TestWaitForURLCurrentURLEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))
	must.NoError(page.WaitForURL(ctx, "about:blank"))
}

func TestWaitForURLHistoryPushStateEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))
	_, err := page.Evaluate(ctx, `() => history.pushState({}, '', '/pushed-url')`)
	must.NoError(err)
	must.NoError(page.WaitForURL(ctx, "**/pushed-url"))
}

func TestWaitForURLWildcardSuffixEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))
	_, err := page.Evaluate(ctx, `() => history.pushState({}, '', '/path/to/page.html')`)
	must.NoError(err)
	must.NoError(page.WaitForURL(ctx, "**/*.html"))
}

func TestWaitForURLHashChangeEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	srv := testserver.New(t)
	srv.SetRoute("/hash-page", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<html><body><a href="#section1" id="link">Jump</a></body></html>`)
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/hash-page"))
	must.NoError(page.Locator("#link").Click(ctx))

	must.NoError(page.WaitForURL(ctx, "**/hash-page#section1"))
}

func TestWaitForURLQueryParamEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	srv := testserver.New(t)
	srv.SetRoute("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<title>Search</title>`)
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/search?q=test"))

	url := page.URL()
	is.Contains(url, "q=test")
}

func TestWaitForURLAfterGoToEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	srv := testserver.New(t)
	srv.SetRoute("/target-page", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<title>Target</title>`)
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/target-page"))
	must.NoError(page.WaitForURL(ctx, "**/target-page"))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Target", title)
}

// TestWaitForURLWithPrefixPattern verifies WaitForURL with wildcard prefix.
// Ref: TestPageWaitForNavigation.java#shouldWorkWithURLPattern
func TestWaitForURLWithPrefixPattern(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, srv.Prefix()+"*"))
}

// TestWaitForURLAfterNavigation verifies WaitForURL resolves after sequential navigation.
// Ref: TestPageWaitForNavigation.java#shouldWorkWithNavigationToNewPage
func TestWaitForURLAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/target", "text/html", `<p>target</p>`)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/target"))
	must.NoError(page.WaitForURL(ctx, srv.Prefix()+"/target"))
}

// TestWaitForURLTimeoutOnMismatch verifies WaitForURL times out when URL never matches.
// Ref: TestPageWaitForNavigation.java#shouldRespectTimeout
func TestWaitForURLTimeoutOnMismatch(t *testing.T) {
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

// TestWFNavLoadStateLoad verifies WaitForLoadState resolves after Goto.
// Ref: TestPageWaitForNavigation.java#shouldWorkWithLoad
func TestWFNavLoadStateLoad(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx))
}

// TestWFNavLoadStateNetworkIdle verifies WaitForLoadState with networkidle.
// Ref: TestPageWaitForNavigation.java#shouldWorkWithNetworkIdle
func TestWFNavLoadStateNetworkIdle(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx, "networkidle"))
}

// TestWaitForURLAfterClickNavigation verifies URL settles after a click-triggered navigation.
// Ref: TestPageWaitForNavigation.java#shouldWorkWithClickingOnAnchorLinks
func TestWaitForURLAfterClickNavigation(t *testing.T) {
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

// TestWaitForURLCrossOriginNavigation verifies WaitForURL works after cross-origin navigation.
// Ref: TestPageWaitForNavigation.java#shouldWorkWithCrossProcessNavigation
func TestWaitForURLCrossOriginNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	crossURL := srv.CrossProcessPrefix() + "/empty.html"

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Goto(ctx, crossURL))
	must.NoError(page.WaitForURL(ctx, crossURL))
}

// TestWFNavLoadStateIdempotent verifies WaitForLoadState returns immediately when already at state.
// Ref: TestPageWaitForNavigation.java#shouldReturnImmediatelyIfAlreadyAtState
func TestWFNavLoadStateIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx, "load"))
	// Calling again on an already-loaded page should return immediately.
	must.NoError(page.WaitForLoadState(ctx, "load"))
}
