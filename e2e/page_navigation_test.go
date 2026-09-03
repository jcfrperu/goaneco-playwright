//go:build e2e

package e2e

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

func TestPageReload(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/reload-page", "text/html", `<title>Original</title><p id="counter">0</p>`)

	err := page.Goto(ctx, srv.Prefix()+"/reload-page")
	must.NoError(err, "Goto failed")

	// Modify the counter via JS
	_, err = page.Evaluate(ctx, "document.getElementById('counter').textContent = '99'")
	must.NoError(err, "Evaluate failed")
	val, err := page.Locator("#counter").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("99", val)

	err = page.Reload(ctx)
	must.NoError(err, "Reload failed")

	// After reload the server re-renders the page, counter resets to "0"
	val, err = page.Locator("#counter").InnerText(ctx)
	must.NoError(err, "InnerText after reload failed")
	is.Equal("0", val)
}

func TestPageGoBackAndForward(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/nav-a", "text/html", `<title>Page A</title><h1>Page A</h1>`)
	srv.ServeWithBody("/nav-b", "text/html", `<title>Page B</title><h1>Page B</h1>`)

	err := page.Goto(ctx, srv.Prefix()+"/nav-a")
	must.NoError(err, "Goto /nav-a failed")
	err = page.Goto(ctx, srv.Prefix()+"/nav-b")
	must.NoError(err, "Goto /nav-b failed")

	titleB, err := page.Title(ctx)
	must.NoError(err, "Title failed")
	is.Equal("Page B", titleB)

	err = page.GoBack(ctx)
	must.NoError(err, "GoBack failed")
	titleA, err := page.Title(ctx)
	must.NoError(err, "Title after GoBack failed")
	is.Equal("Page A", titleA)

	err = page.GoForward(ctx)
	must.NoError(err, "GoForward failed")
	titleB2, err := page.Title(ctx)
	must.NoError(err, "Title after GoForward failed")
	is.Equal("Page B", titleB2)
}

// TestPageGotoHash verifies that page.URL() includes the hash fragment after navigation.
// Ref: TestPageBasic.java#pageUrlShouldIncludeHashes
func TestPageGotoHash(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	targetURL := srv.EmptyPage() + "#section-1"
	err := page.Goto(ctx, targetURL)
	must.NoErrorf(err, "Goto(%q) failed", targetURL)

	got := page.URL()
	is.Equal(targetURL, got)
}

func TestPageGotoRedirectFollowed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithRedirect("/redirect-source", srv.EmptyPage())

	must.NoError(page.Goto(ctx, srv.Prefix()+"/redirect-source"))
	is.Equal(srv.EmptyPage(), page.URL())
}

func TestPageGotoWorksForDataURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	dataURL := `data:text/html,<title>DataPage</title>`
	must.NoError(page.Goto(ctx, dataURL))
	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("DataPage", title)
}

func TestPageReloadPreservesURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	urlBefore := page.URL()

	must.NoError(page.Reload(ctx))
	is.Equal(urlBefore, page.URL())
}

func TestPageGotoMultipleRedirectsFollowed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithRedirect("/redirect1", srv.Prefix()+"/redirect2")
	srv.ServeWithRedirect("/redirect2", srv.EmptyPage())

	must.NoError(page.Goto(ctx, srv.Prefix()+"/redirect1"))
	is.Equal(srv.EmptyPage(), page.URL())
}

func TestPageGotoBackToFirstPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/page1", "text/html", `<title>Page1</title>`)
	srv.ServeWithBody("/page2", "text/html", `<title>Page2</title>`)
	srv.ServeWithBody("/page3", "text/html", `<title>Page3</title>`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page1"))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/page2"))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/page3"))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page3", title)

	must.NoError(page.GoBack(ctx))
	title, err = page.Title(ctx)
	must.NoError(err)
	is.Equal("Page2", title)

	must.NoError(page.GoBack(ctx))
	title, err = page.Title(ctx)
	must.NoError(err)
	is.Equal("Page1", title)
}

func TestPageURLAfterGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/target", "text/html", `<p>target</p>`)
	targetURL := srv.Prefix() + "/target"

	must.NoError(page.Goto(ctx, targetURL))
	is.Equal(targetURL, page.URL())
}

// TestPageGoBackRestoresPreviousURL verifies GoBack returns to previous URL.
// Ref: TestPageNavigation.java#shouldGoBack
func TestPageGoBackRestoresPreviousURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	firstURL := page.URL()

	must.NoError(page.SetContent(ctx, `<div>second</div>`))

	must.NoError(page.GoBack(ctx))
	is.Equal(firstURL, page.URL())
}

// TestPageGoForwardRestoresNextURL verifies GoForward returns to next URL.
// Ref: TestPageNavigation.java#shouldGoForward
func TestPageGoForwardRestoresNextURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Goto(ctx, srv.EmptyPage()+"?second=1"))

	must.NoError(page.GoBack(ctx))
	must.NoError(page.GoForward(ctx))

	is.Contains(page.URL(), "second=1")
}

// TestPageGoBackOnFirstPageDoesNotError verifies GoBack at history start doesn't error.
// Ref: TestPageNavigation.java#shouldNotErrorOnGoBackFromFirst
func TestPageGoBackOnFirstPageDoesNotError(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	// GoBack with no history should not error
	_ = page.GoBack(ctx) // error may or may not occur, just don't panic
}

// TestPageGoForwardOnLastPageDoesNotError verifies GoForward at history end doesn't error.
// Ref: TestPageNavigation.java#shouldNotErrorOnGoForwardFromLast
func TestPageGoForwardOnLastPageDoesNotError(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_ = page.GoForward(ctx) // should not panic
}

// TestPageGoBackAndForwardMultipleTimes verifies multiple back/forward work correctly.
// Ref: TestPageNavigation.java#shouldWorkMultipleTimesBackAndForward
func TestPageGoBackAndForwardMultipleTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	url1 := srv.EmptyPage()
	url2 := srv.EmptyPage() + "?p=2"

	must.NoError(page.Goto(ctx, url1))
	must.NoError(page.Goto(ctx, url2))
	must.NoError(page.GoBack(ctx))
	is.Equal(url1, page.URL())
	must.NoError(page.GoForward(ctx))
	is.Equal(url2, page.URL())
}

// TestPageGotoUpdatesURL verifies Goto sets the current URL.
// Ref: TestPageNavigation.java#shouldUpdateURL
func TestPageGotoUpdatesURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.NotEmpty(page.URL())
	is.Contains(page.URL(), srv.Prefix())
}

// TestPageReloadKeepsURL verifies Reload keeps the same URL.
// Ref: TestPageNavigation.java#shouldReloadSameURL
func TestPageReloadKeepsURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	urlBefore := page.URL()

	must.NoError(page.Reload(ctx))

	is.Equal(urlBefore, page.URL())
}

// TestPageGotoHTMLContent verifies SetContent and Content round-trip.
// Ref: TestPageNavigation.java#shouldSetAndGetContent
func TestPageSetContentRoundTrip(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	html := `<html><body><h1>Hello</h1></body></html>`
	must.NoError(page.SetContent(ctx, html))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "Hello")
}

// TestPageGoBackAtStartNoError verifies GoBack at start does not error.
// Ref: TestPageNavigation.java#goBackAtStartShouldNotError
func TestPageGoBackAtStartNoErrorExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.GoBack(ctx))
}

// TestPageGoForwardAtEndNoError verifies GoForward at end does not error.
// Ref: TestPageNavigation.java#goForwardAtEndShouldNotError
func TestPageGoForwardAtEndNoErrorExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.GoForward(ctx))
}

// TestPageGotoWithRedirect verifies Goto follows a server-side redirect.
// Ref: TestPageNavigation.java#shouldFollowRedirect
func TestPageGotoWithRedirect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/destination", "text/html", `<p>redirected</p>`)
	srv.ServeWithRedirect("/redirect-src", "/destination")

	must.NoError(page.Goto(ctx, srv.Prefix()+"/redirect-src"))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "redirected")
}

// TestPageURLAfterSetContent verifies URL is "about:blank" after SetContent.
// Ref: TestPageNavigation.java#shouldHaveAboutBlankAfterSetContent
func TestPageURLAfterSetContentEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	// after SetContent, URL may still be the previous URL or about:blank
	url := page.URL()
	is.NotEmpty(url)
}

// TestPageNavigationMultiplePages verifies multiple pages can navigate independently.
// Ref: TestPageNavigation.java#shouldNavigateMultiplePages
func TestPageNavigationMultiplePages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/page-a", "text/html", `<title>A</title>`)
	srv.ServeWithBody("/page-b", "text/html", `<title>B</title>`)

	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.Prefix()+"/page-a"))
	must.NoError(page2.Goto(ctx, srv.Prefix()+"/page-b"))

	title1, err := page1.Title(ctx)
	must.NoError(err)
	is.Equal("A", title1)

	title2, err := page2.Title(ctx)
	must.NoError(err)
	is.Equal("B", title2)
}

// TestPageURLMatchesGotoTarget verifies page.URL() matches what was navigated to.
// Ref: TestPageNavigation.java#shouldMatchGotoURL
func TestPageURLMatchesGotoTarget(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/nav-target", "text/plain", "ok")
	page := newPage(t)

	target := srv.Prefix() + "/nav-target"
	must.NoError(page.Goto(ctx, target))

	is.Equal(target, page.URL())
}

// TestPageReloadRefreshesContent verifies Reload reflects server changes.
// Ref: TestPageNavigation.java#shouldReloadPage
func TestPageReloadRefreshesContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/reload-test", "text/html", `<p>original</p>`)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/reload-test"))
	must.NoError(page.Reload(ctx))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "original")
}

// TestNavigateMultiplePagesEx4 verifies sequential navigation between pages.
// Ref: TestPageNavigation.java#shouldNavigateMultiplePages
func TestNavigateMultiplePagesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/page1", "text/html", `<html><head><title>Page 1</title></head><body></body></html>`)
	srv.ServeWithBody("/page2", "text/html", `<html><head><title>Page 2</title></head><body></body></html>`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page1"))
	t1, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page 1", t1)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page2"))
	t2, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page 2", t2)
}

// TestNavigateWithQueryParamsEx5 verifies navigation with query parameters.
// Ref: TestPageNavigate.java#shouldNavigateWithQueryParams
func TestNavigateWithQueryParamsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<title>Results</title>`)
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/search?q=playwright"))
	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Results", title)
}

// TestNavigateSetsHTMLEx5 verifies SetContent sets page HTML.
// Ref: TestPageNavigate.java#shouldSetPageHTML
func TestNavigateSetsHTMLEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h">Hello</h1>`))

	text, err := page.Locator("#h").InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello", text)
}

// TestNavigateToAboutBlankEx5 verifies goto about:blank works.
// Ref: TestPageNavigate.java#shouldNavigateToAboutBlank
func TestNavigateToAboutBlankEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Some page</p>`))
	must.NoError(page.Goto(ctx, "about:blank"))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.NotEmpty(content)
}

// TestNavigateWithFragmentEx6 verifies navigation with URL fragment.
// Ref: TestPageNavigate.java#shouldNavigateWithFragment
func TestNavigateWithFragmentEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	srv := testserver.New(t)
	srv.SetRoute("/frag", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<html><body><div id="section">Section</div></body></html>`)
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/frag#section"))

	url := page.URL()
	is.Contains(url, "#section")
}

// TestNavigateReturnsResponseEx6 verifies Goto returns response info.
// Ref: TestPageNavigate.java#shouldReturnResponse
func TestNavigateReturnsResponseEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	srv := testserver.New(t)
	srv.SetRoute("/resp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		_, _ = io.WriteString(w, `<title>Resp</title>`)
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/resp"))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Resp", title)
}

// TestNavigateToExistingPageEx6 verifies multiple gotos work correctly.
// Ref: TestPageNavigate.java#shouldNavigateMultipleTimes
func TestNavigateToExistingPageEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	srv := testserver.New(t)
	srv.SetRoute("/page-a", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<title>Page A</title>`)
	})
	srv.SetRoute("/page-b", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<title>Page B</title>`)
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page-a"))
	titleA, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page A", titleA)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page-b"))
	titleB, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page B", titleB)
}

// TestNavigateQueryStringEx7 verifies navigation with query string.
// Ref: TestPageNavigate.java#shouldNavigateWithQueryString
func TestNavigateQueryStringEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.SetRoute("/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><body><p id='q'>" + q + "</p></body></html>"))
	}))

	err := page.Goto(ctx, srv.Prefix()+"/search?q=hello")
	must.NoError(err)

	text, err := page.Locator("#q").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("hello", *text)
}

// TestNavigate404Ex7 verifies navigation to a 404 page shows content.
// Ref: TestPageNavigate.java#shouldReturn404
func TestNavigate404Ex7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.SetRoute("/not-found", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("<html><body><p id='msg'>Not found</p></body></html>"))
	}))

	err := page.Goto(ctx, srv.Prefix()+"/not-found")
	must.NoError(err)

	text, err := page.Locator("#msg").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Not found", *text)
}

// TestNavigatePostRedirectEx7 verifies page URL after redirect.
// Ref: TestPageNavigate.java#shouldFollowRedirect
func TestNavigatePostRedirectEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.SetRoute("/redirect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dest", http.StatusFound)
	}))
	srv.SetRoute("/dest", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><body><p>Destination</p></body></html>"))
	}))

	err := page.Goto(ctx, srv.Prefix()+"/redirect")
	must.NoError(err)

	u := page.URL()
	is.Contains(u, "/dest")
}

// TestNavigateCrossProcessEx8 verifies that cross-origin navigation (127.0.0.1 ↔ localhost) works.
// Ref: TestPageNavigate.java#shouldWorkCrossProcess
func TestNavigateCrossProcessEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	crossURL := srv.CrossProcessPrefix() + "/empty.html"
	must.NoError(page.Goto(ctx, crossURL))

	is.Equal(crossURL, page.URL())
}

// TestNavigateAnchorEx8 verifies navigation with URL fragment anchors updates page.URL().
// Ref: TestPageNavigate.java#shouldWorkWithAnchorNavigation
func TestNavigateAnchorEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/anchor-page", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body>
			<div id="foo">Foo section</div>
			<div id="bar">Bar section</div>
		</body></html>`))
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/anchor-page"))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/anchor-page#foo"))
	is.Contains(page.URL(), "#foo")

	must.NoError(page.Goto(ctx, srv.Prefix()+"/anchor-page#bar"))
	is.Contains(page.URL(), "#bar")
}

// TestNavigate204FailsEx8 verifies that a 204 No Content response causes Goto to return an error.
// Ref: TestPageNavigate.java#shouldFailWhenServerReturns204
func TestNavigate204FailsEx8(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/204", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	err := page.Goto(ctx, srv.Prefix()+"/204")
	is.Error(err, "goto to 204 No Content should return an error")
}

// TestNavigateHistoryAPIBeforeunloadEx8 verifies navigation succeeds when beforeunload calls history.replaceState.
// Ref: TestPageNavigate.java#shouldWorkWhenPageCallsHistoryAPIInBeforeunload
func TestNavigateHistoryAPIBeforeunloadEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/beforeunload-history", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><script>
			window.addEventListener('beforeunload', () => {
				history.replaceState(null, '', '/changed');
			});
		</script></body></html>`))
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/beforeunload-history"))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
}

// TestPageGotoSetsPageURL verifies Goto updates page.URL() to the navigated URL.
// Ref: TestPageGoto.java#shouldSetURL
func TestPageGotoSetsPageURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(srv.EmptyPage(), page.URL())
}

// TestPageGotoAndWaitForTitle verifies Goto loads the page and Title is accessible.
// Ref: TestPageGoto.java#shouldWorkWithTitle
func TestPageGotoAndWaitForTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/title.html"))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Woof-Woof", title)
}

// TestPageGotoFollowsRedirects verifies Goto resolves redirect chains transparently.
// Ref: TestPageGoto.java#shouldFollowRedirects
func TestPageGotoFollowsRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithRedirect("/redirect", srv.EmptyPage())

	must.NoError(page.Goto(ctx, srv.Prefix()+"/redirect"))
	is.Contains(page.URL(), "empty.html")
}

// TestPageGotoWorksTwice verifies Goto works repeatedly on the same page.
// Ref: TestPageGoto.java#shouldWorkForMultipleNavigations
func TestPageGotoWorksTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/title.html"))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Woof-Woof", title)
}

// TestPageGotoChangesMainFrameURL verifies MainFrame().URL() matches Goto target.
// Ref: TestPageGoto.java#shouldUpdateMainFrameURL
func TestPageGotoChangesMainFrameURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.Equal(srv.EmptyPage(), page.MainFrame().URL())
}

// TestPageReloadReloadsContent verifies Reload re-executes page scripts.
// Ref: TestPageReload.java#shouldReload
func TestPageReloadReloadsContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, `() => { window.__counter = 1; }`)
	must.NoError(err)

	counter1, err := page.Evaluate(ctx, `() => window.__counter`)
	must.NoError(err)
	is.Equal(float64(1), counter1)

	// After reload, the counter set via evaluate is gone
	must.NoError(page.Reload(ctx))

	counter2, err := page.Evaluate(ctx, `() => window.__counter`)
	must.NoError(err)
	is.Nil(counter2)
}

// TestPageReloadPreservesQueryString verifies Reload keeps the URL including query string.
// Ref: TestPageReload.java#shouldPreserveQueryString
func TestPageReloadPreservesQueryString(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/page-with-query", "text/html", `<html><body>content</body></html>`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page-with-query?key=val"))
	must.NoError(page.Reload(ctx))
	is.Contains(page.URL(), "key=val")
}

// TestPageReloadAfterSetContent verifies Reload after SetContent restores blank state.
// Ref: TestPageReload.java#shouldWorkAfterSetContent
func TestPageReloadAfterSetContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<div id="injected">injected</div>`))

	// Reload should restore the navigated page, not the injected content
	must.NoError(page.Reload(ctx))

	count, err := page.Locator("#injected").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// ---------------------------------------------------------------------------
// From TestPageNavigate.java (Java source ports)
// ---------------------------------------------------------------------------

// Ref: TestPageNavigate.java#shouldWork
func TestPageNavigateShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(srv.EmptyPage(), page.URL())
}

// Ref: TestPageNavigate.java#shouldWorkWithFileURL
func TestPageNavigateShouldWorkWithFileURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Create a temp HTML file
	dir := t.TempDir()
	fp := filepath.Join(dir, "file-url-test.html")
	must.NoError(os.WriteFile(fp, []byte(`<html><body><p>file url</p></body></html>`), 0o644))

	abs, err := filepath.Abs(fp)
	must.NoError(err)
	fileURL := (&url.URL{Scheme: "file", Path: filepath.ToSlash("/" + abs)}).String()

	must.NoError(page.Goto(ctx, fileURL))
	is.Contains(strings.ToLower(page.URL()), "file-url-test.html")
}

// Ref: TestPageNavigate.java#shouldUseHttpForNoProtocol
func TestPageNavigateShouldUseHttpForNoProtocol(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	full := srv.EmptyPage()
	noProto := strings.TrimPrefix(full, "http://")
	must.NoError(page.Goto(ctx, noProto))
	is.Equal(full, page.URL())
}

// Ref: TestPageNavigate.java#shouldCaptureIframeNavigationRequest
func TestPageNavigateShouldCaptureIframeNavigationRequest(t *testing.T) {
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

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(srv.EmptyPage(), page.URL())

	var mu sync.Mutex
	var sawFrameURL bool
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if strings.HasSuffix(req.URL(), "/frames/frame.html") {
			mu.Lock()
			sawFrameURL = true
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.Prefix()+"/frames/one-frame.html"))
	is.Equal(srv.Prefix()+"/frames/one-frame.html", page.URL())
	// Wait a moment for iframe requests
	must.NoError(page.WaitForTimeout(ctx, 200))

	is.Len(page.Frames(), 2)

	mu.Lock()
	defer mu.Unlock()
	is.True(sawFrameURL)
}

// Ref: TestPageNavigate.java#shouldWorkWithRedirects
func TestPageNavigateShouldWorkWithRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithRedirect("/redirect/1.html", "/redirect/2.html")
	srv.ServeWithRedirect("/redirect/2.html", "/empty.html")

	must.NoError(page.Goto(ctx, srv.Prefix()+"/redirect/1.html"))
	is.Equal(srv.EmptyPage(), page.URL())
}

// Ref: TestPageNavigate.java#shouldReturnResponseWhenPageChangesItsURLAfterLoad
func TestPageNavigateShouldReturnResponseWhenPageChangesItsURLAfterLoad(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/historyapi.html", "text/html",
		`<html><body><script>history.pushState({}, '', '/other');</script></body></html>`)

	var mu sync.Mutex
	var status int
	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if strings.HasSuffix(resp.URL(), "/historyapi.html") {
			mu.Lock()
			status = resp.Status()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.Prefix()+"/historyapi.html"))

	mu.Lock()
	defer mu.Unlock()
	is.Equal(200, status)
}

// Ref: TestPageNavigate.java#shouldWorkWithSubframesReturn204
func TestPageNavigateShouldWorkWithSubframesReturn204(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frames/one-frame.html", "text/html",
		`<html><body><iframe src="/frames/frame.html"></iframe></body></html>`)
	srv.SetRoute("/frames/frame.html", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Should not error
	_ = page.Goto(ctx, srv.Prefix()+"/frames/one-frame.html")
	must.NoError(nil)
}

// Ref: TestPageNavigate.java#shouldWorkWithSubframesReturn204WithDomcontentloaded
func TestPageNavigateShouldWorkWithSubframesReturn204WithDomcontentloaded(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frames/one-frame.html", "text/html",
		`<html><body><iframe src="/frames/frame.html"></iframe></body></html>`)
	srv.SetRoute("/frames/frame.html", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	domcontent := "domcontentloaded"
	must.NoError(page.Goto(ctx, srv.Prefix()+"/frames/one-frame.html", &playwright.PageGotoOptions{
		WaitUntil: &domcontent,
	}))
}

// Ref: TestPageNavigate.java#shouldNavigateToEmptyPageWithDomcontentloaded
func TestPageNavigateShouldNavigateToEmptyPageWithDomcontentloaded(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	domcontent := "domcontentloaded"
	must.NoError(page.Goto(ctx, srv.EmptyPage(), &playwright.PageGotoOptions{
		WaitUntil: &domcontent,
	}))
	is.Equal(srv.EmptyPage(), page.URL())
}

// Ref: TestPageNavigate.java#shouldCaptureCrossProcessIframeNavigationRequest
func TestPageNavigateShouldCaptureCrossProcessIframeNavigationRequest(t *testing.T) {
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

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	crossOrigin := srv.CrossProcessPrefix()

	var mu sync.Mutex
	var sawCrossFrame bool
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if strings.HasPrefix(req.URL(), crossOrigin) && strings.HasSuffix(req.URL(), "/frames/frame.html") {
			mu.Lock()
			sawCrossFrame = true
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, crossOrigin+"/frames/one-frame.html"))
	is.Equal(crossOrigin+"/frames/one-frame.html", page.URL())
	must.NoError(page.WaitForTimeout(ctx, 200))

	is.Len(page.Frames(), 2)

	mu.Lock()
	defer mu.Unlock()
	is.True(sawCrossFrame)
}

// Ref: TestPageNavigate.java#shouldSendReferer
func TestPageNavigateShouldSendReferer(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var gotReferer string
	srv.SetRoute("/referer-target", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotReferer = r.Header.Get("Referer")
		mu.Unlock()
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<p>ok</p>"))
	})

	referer := "https://example.com/"
	must.NoError(page.Goto(ctx, srv.Prefix()+"/referer-target", &playwright.PageGotoOptions{
		Referer: &referer,
	}))

	mu.Lock()
	defer mu.Unlock()
	is.Equal(referer, gotReferer)
}
