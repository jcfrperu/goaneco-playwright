//go:build e2e

// E2E tests for: TestBrowserContextBaseUrl.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestBrowserContextBaseUrlShouldConstructANewURLWhenABaseURLInBrowserNewPageIsPassedToPageGoto verifies BaseURL in Browser.NewPage resolves relative URLs.
// Ref: TestBrowserContextBaseUrl.java#shouldConstructANewURLWhenABaseURLInBrowserNewPageIsPassedToPageGoto  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseUrlShouldConstructANewURLWhenABaseURLInBrowserNewPageIsPassedToPageGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/empty.html", "text/html", `<title>Empty</title>`)

	baseURL := srv.Prefix()
	page, err := globalBrowser.NewPage(ctx, &playwright.BrowserContextOptions{BaseURL: &baseURL})
	must.NoError(err, "NewPage failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = page.Close(c)
	})

	must.NoError(page.Goto(ctx, "/empty.html"))
	is.Equal(srv.Prefix()+"/empty.html", page.URL())
}

// TestBrowserContextBaseUrlShouldNotConstructANewURLWhenValidURLsArePassed verifies absolute and special URLs
// are not affected by BaseURL.
// Ref: TestBrowserContextBaseUrl.java#shouldNotConstructANewURLWhenValidURLsArePassed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseUrlShouldNotConstructANewURLWhenValidURLsArePassed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/empty.html", "text/html", `<title>Empty</title>`)

	baseURL := "http://microsoft.com"
	page, err := globalBrowser.NewPage(ctx, &playwright.BrowserContextOptions{BaseURL: &baseURL})
	must.NoError(err, "NewPage failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = page.Close(c)
	})

	// Absolute URL should not be affected.
	must.NoError(page.Goto(ctx, srv.Prefix()+"/empty.html"))
	is.Equal(srv.Prefix()+"/empty.html", page.URL())

	// data: URL should not be affected.
	must.NoError(page.Goto(ctx, "data:text/html,Hello world"))
	val, err := page.Evaluate(ctx, "window.location.href")
	must.NoError(err)
	is.Equal("data:text/html,Hello world", val)

	// about:blank should not be affected.
	must.NoError(page.Goto(ctx, "about:blank"))
	val, err = page.Evaluate(ctx, "window.location.href")
	must.NoError(err)
	is.Equal("about:blank", val)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextBaseUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseURLGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/page", "text/html", `<title>Base URL Page</title><p>hello</p>`)

	baseURL := srv.Prefix()
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{BaseURL: &baseURL})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// Navigate using a relative URL — should resolve against BaseURL.
	err = page.Goto(ctx, "/page")
	must.NoError(err, "Goto(/page) with BaseURL failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title failed")
	if title != "Base URL Page" {
		t.Errorf("Title = %q, want 'Base URL Page'", title)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextBaseUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseURLNotAffectsAbsoluteURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/absolute", "text/html", `<title>Absolute</title><p>absolute</p>`)

	baseURL := "http://example.invalid"
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{BaseURL: &baseURL})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// Absolute URL should NOT be affected by BaseURL.
	err = page.Goto(ctx, srv.Prefix()+"/absolute")
	must.NoError(err, "Goto with absolute URL failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title failed")
	if title != "Absolute" {
		t.Errorf("Title = %q, want 'Absolute'", title)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextBaseUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseURLWithTrailingSlash(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/sub", "text/html", `<title>Sub</title><p>sub</p>`)

	baseURL := srv.Prefix() + "/"
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{BaseURL: &baseURL})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, "sub")
	must.NoError(err, "Goto(sub) with trailing-slash BaseURL failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title failed")
	if title != "Sub" {
		t.Errorf("Title = %q, want 'Sub'", title)
	}
}

// TestBrowserContextBaseURLConstructURLsWithoutTrailingSlash verifies URL construction
// when BaseURL in Browser.NewPage has no trailing slash.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextBaseUrl.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseURLConstructURLsWithoutTrailingSlash(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/mypage.html", "text/html", `<title>MyPage</title>`)

	baseURL := srv.Prefix() + "/url-construction"
	page, err := globalBrowser.NewPage(ctx, &playwright.BrowserContextOptions{BaseURL: &baseURL})
	must.NoError(err, "NewPage failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = page.Close(c)
	})

	must.NoError(page.Goto(ctx, "mypage.html"))
	is.Equal(srv.Prefix()+"/mypage.html", page.URL())

	must.NoError(page.Goto(ctx, "./mypage.html"))
	is.Equal(srv.Prefix()+"/mypage.html", page.URL())

	must.NoError(page.Goto(ctx, "/mypage.html"))
	is.Equal(srv.Prefix()+"/mypage.html", page.URL())
}

// TestBrowserContextBaseUrlShouldConstructTheURLsCorrectlyWhenABaseURLWithATrailingSlashInBrowserNewPageIsPassedToPageGoto verifies URL resolution with a trailing-slash base URL.
// Ref: TestBrowserContextBaseUrl.java#shouldConstructTheURLsCorrectlyWhenABaseURLWithATrailingSlashInBrowserNewPageIsPassedToPageGoto  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseUrlShouldConstructTheURLsCorrectlyWhenABaseURLWithATrailingSlashInBrowserNewPageIsPassedToPageGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	for _, p := range []string{"/url-construction/mypage.html", "/mypage.html", "/url-construction/", "/"} {
		srv.ServeWithBody(p, "text/html", "<title>x</title>")
	}
	page := a4NewPage(t, &playwright.BrowserContextOptions{BaseURL: a4Ptr(srv.Prefix() + "/url-construction/")})

	must.NoError(page.Goto(ctx, "mypage.html"))
	is.Equal(srv.Prefix()+"/url-construction/mypage.html", page.URL())
	must.NoError(page.Goto(ctx, "./mypage.html"))
	is.Equal(srv.Prefix()+"/url-construction/mypage.html", page.URL())
	must.NoError(page.Goto(ctx, "/mypage.html"))
	is.Equal(srv.Prefix()+"/mypage.html", page.URL())
	must.NoError(page.Goto(ctx, "."))
	is.Equal(srv.Prefix()+"/url-construction/", page.URL())
	must.NoError(page.Goto(ctx, "/"))
	is.Equal(srv.Prefix()+"/", page.URL())
}

// TestBrowserContextBaseUrlShouldBeAbleToMatchAURLRelativeToItsGivenURLWithUrlMatcher verifies URL matchers are resolved against the base URL.
// Ref: TestBrowserContextBaseUrl.java#shouldBeAbleToMatchAURLRelativeToItsGivenURLWithUrlMatcher  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseUrlShouldBeAbleToMatchAURLRelativeToItsGivenURLWithUrlMatcher(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/foobar/kek/index.html", "text/html", "<title>x</title>")
	srv.ServeWithBody("/kek/index.html", "text/html", "<title>x</title>")
	page := a4NewPage(t, &playwright.BrowserContextOptions{BaseURL: a4Ptr(srv.Prefix() + "/foobar/")})

	must.NoError(page.Goto(ctx, "/kek/index.html"))
	must.NoError(page.WaitForURL(ctx, "/kek/index.html"))
	is.Equal(srv.Prefix()+"/kek/index.html", page.URL())

	must.NoError(page.Route(ctx, "./kek/index.html", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: a4Ptr("base-url-matched-route")})
	}))
	reqCh := make(chan *playwright.NetworkRequest, 1)
	respCh := make(chan *playwright.NetworkResponse, 1)
	offReq := page.OnRequest(func(r *playwright.NetworkRequest) {
		select {
		case reqCh <- r:
		default:
		}
	})
	defer offReq()
	offResp := page.OnResponse(func(r *playwright.NetworkResponse) {
		select {
		case respCh <- r:
		default:
		}
	})
	defer offResp()
	must.NoError(page.Goto(ctx, "./kek/index.html"))

	select {
	case req := <-reqCh:
		is.Equal(srv.Prefix()+"/foobar/kek/index.html", req.URL())
	case <-time.After(5 * time.Second):
		t.Fatal("request not observed")
	}
	select {
	case resp := <-respCh:
		is.Equal(srv.Prefix()+"/foobar/kek/index.html", resp.URL())
		text, err := resp.Text(ctx)
		must.NoError(err)
		is.Equal("base-url-matched-route", text)
	case <-time.After(5 * time.Second):
		t.Fatal("response not observed")
	}
}

// TestBrowserContextBaseUrlShouldConstructANewURLWhenABaseURLInBrowserNewContextIsPassedToPageGoto verifies a relative Goto resolves against the context base URL.
// Ref: TestBrowserContextBaseUrl.java#shouldConstructANewURLWhenABaseURLInBrowserNewContextIsPassedToPageGoto  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseUrlShouldConstructANewURLWhenABaseURLInBrowserNewContextIsPassedToPageGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := a4NewPage(t, &playwright.BrowserContextOptions{BaseURL: a4Ptr(srv.Prefix())})
	must.NoError(page.Goto(ctx, "/empty.html"))
	must.Equal(srv.EmptyPage(), page.URL())
}

// TestBrowserContextBaseUrlShouldConstructTheURLsCorrectlyWhenABaseURLWithoutATrailingSlashInBrowserNewPageIsPassedToPageGoto verifies relative URL resolution against a base URL without trailing slash.
// Ref: TestBrowserContextBaseUrl.java#shouldConstructTheURLsCorrectlyWhenABaseURLWithoutATrailingSlashInBrowserNewPageIsPassedToPageGoto  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextBaseUrl.java
func TestBrowserContextBaseUrlShouldConstructTheURLsCorrectlyWhenABaseURLWithoutATrailingSlashInBrowserNewPageIsPassedToPageGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := a4NewPage(t, &playwright.BrowserContextOptions{BaseURL: a4Ptr(srv.Prefix() + "/url-construction")})
	for _, rel := range []string{"mypage.html", "./mypage.html", "/mypage.html"} {
		must.NoError(page.Goto(ctx, rel), rel)
		must.Equal(srv.Prefix()+"/mypage.html", page.URL(), rel)
	}
}
