//go:build e2e

// E2E tests for: TestPageWaitForUrl.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/jcfrperu/goaneco-playwright/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestWaitForURLExactMatch verifies WaitForURL resolves when URL matches exactly.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestWaitForURLExactMatchAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, srv.EmptyPage()))
}

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestWaitForURLWildcardPatternExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, "**/empty*"))
}

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestWaitForURLWithTimeoutSucceeds(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForURL(ctx, srv.EmptyPage(), 5*time.Second))
}

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestWaitForURLCurrentURLEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))
	must.NoError(page.WaitForURL(ctx, "about:blank"))
}

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// TestWFNavLoadStateLoad verifies WaitForLoadState resolves after Goto.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestWFNavLoadStateNetworkIdle(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForLoadState(ctx, "networkidle"))
}

// TestWaitForURLCrossOriginNavigation verifies WaitForURL works after cross-origin navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageWaitForUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
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

// Ref: TestPageWaitForUrl.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	a2Ev(t, ctx, page, "url => window.location.href = url", srv.Prefix()+"/grid.html")
	must.NoError(page.WaitForURL(ctx, "**/grid.html"))
}

// Ref: TestPageWaitForUrl.java#shouldRespectTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	err := page.WaitForURL(ctx, "**/frame.html", 2500*time.Millisecond)
	must.Error(err)
	assert.Contains(t, err.Error(), "Timeout 2500ms exceeded")
}

// Ref: TestPageWaitForUrl.java#shouldWorkWithBothDomcontentloadedAndLoad  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkWithBothDomcontentloadedAndLoad(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/one-style.html"))
	dcl := protocol.LifecycleEventDomcontentloaded
	must.NoError(page.WaitForURL(ctx, "**/one-style.html", &playwright.PageWaitForURLOptions{
		WaitUntil: &dcl,
	}))
	load := protocol.LifecycleEventLoad
	must.NoError(page.WaitForURL(ctx, "**/one-style.html", &playwright.PageWaitForURLOptions{
		WaitUntil: &load,
	}))
}

// Ref: TestPageWaitForUrl.java#shouldWorkWithCommit  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkWithCommit(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, srv, page := a2Setup(t)

	srv.SetRoute("/script.js", func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<title>Hello</title><script src="script.js"></script>`)
	})

	to := 100.0
	_ = page.Goto(ctx, srv.EmptyPage(), &playwright.PageGotoOptions{Timeout: &to})

	commit := protocol.LifecycleEventCommit
	must.NoError(page.WaitForURL(ctx, "**/empty.html", &playwright.PageWaitForURLOptions{
		WaitUntil: &commit,
	}))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Hello", title)
}

// Ref: TestPageWaitForUrl.java#shouldWorkWithCommitAndAboutBlank  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkWithCommitAndAboutBlank(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	commit := protocol.LifecycleEventCommit
	must.NoError(page.WaitForURL(ctx, "about:blank", &playwright.PageWaitForURLOptions{
		WaitUntil: &commit,
	}))
}

// Ref: TestPageWaitForUrl.java#shouldWorkWithClickingOnAnchorLinks  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkWithClickingOnAnchorLinks(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, "<a href='#foobar'>foobar</a>"))
	must.NoError(a2Click(ctx, page, "a"))
	must.NoError(page.WaitForURL(ctx, "**/*#foobar"))
}

// Ref: TestPageWaitForUrl.java#shouldWorkWithHistoryPushState  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkWithHistoryPushState(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a onclick='javascript:pushState()'>SPA</a>
<script>
  function pushState() { history.pushState({}, '', 'wow.html') }
</script>`))
	must.NoError(a2Click(ctx, page, "a"))
	must.NoError(page.WaitForURL(ctx, "**/wow.html"))
	assert.Equal(t, srv.Prefix()+"/wow.html", page.URL())
}

// Ref: TestPageWaitForUrl.java#shouldWorkWithHistoryReplaceState  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkWithHistoryReplaceState(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, ` <a onclick='javascript:replaceState()'>SPA</a>
<script>
  function replaceState() { history.replaceState({}, '', '/replaced.html') }
</script>`))
	must.NoError(a2Click(ctx, page, "a"))
	must.NoError(page.WaitForURL(ctx, "**/replaced.html"))
	assert.Equal(t, srv.Prefix()+"/replaced.html", page.URL())
}

// Ref: TestPageWaitForUrl.java#shouldWorkWithDOMHistoryBackHistoryForward  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkWithDOMHistoryBackHistoryForward(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a id=back onclick='javascript:goBack()'>back</a>
<a id=forward onclick='javascript:goForward()'>forward</a>
<script>
  function goBack() { history.back(); }
  function goForward() { history.forward(); }
  history.pushState({}, '', '/first.html');
  history.pushState({}, '', '/second.html');
</script>`))
	assert.Equal(t, srv.Prefix()+"/second.html", page.URL())
	must.NoError(a2Click(ctx, page, "a#back"))
	must.NoError(page.WaitForURL(ctx, "**/first.html"))
	assert.Equal(t, srv.Prefix()+"/first.html", page.URL())
	must.NoError(a2Click(ctx, page, "a#forward"))
	must.NoError(page.WaitForURL(ctx, "**/second.html"))
	assert.Equal(t, srv.Prefix()+"/second.html", page.URL())
}

// Ref: TestPageWaitForUrl.java#shouldWorkOnFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForUrl.java
func TestPageWaitForUrlShouldWorkOnFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frames/one-frame.html", "text/html",
		`<html><body><iframe src="/frames/frame.html"></iframe></body></html>`)
	srv.ServeWithBody("/frames/frame.html", "text/html",
		`<html><body>frame</body></html>`)
	srv.ServeWithBody("/grid.html", "text/html",
		`<html><body>grid</body></html>`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/frames/one-frame.html"))
	frames := page.Frames()
	must.GreaterOrEqual(len(frames), 2)
	frame := frames[1]

	_, err := frame.Evaluate(ctx, `url => window.location.href = url`, srv.Prefix()+"/grid.html")
	must.NoError(err)

	must.NoError(frame.WaitForURL(ctx, "**/grid.html"))
	is.Contains(frame.URL(), "grid.html")
}
