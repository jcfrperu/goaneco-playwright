//go:build e2e

// BrowserContext BaseURL option tests.
// Migration of: TestBrowserContextBaseUrl.java
package e2e

import (
	"context"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

// TestBrowserContextBaseURLNewPageOption verifies BaseURL in Browser.NewPage resolves relative URLs.
// Ref: TestBrowserContextBaseUrl.java#shouldConstructANewURLWhenABaseURLInBrowserNewPageIsPassedToPageGoto
func TestBrowserContextBaseURLNewPageOption(t *testing.T) {
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

// TestBrowserContextBaseURLConstructURLsWithoutTrailingSlash verifies URL construction
// when BaseURL in Browser.NewPage has no trailing slash.
// Ref: TestBrowserContextBaseUrl.java#shouldConstructTheURLsCorrectlyWhenABaseURLWithoutATrailingSlash
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

// TestBrowserContextBaseURLNotConstructNewURLsForValidURLs verifies absolute and special URLs
// are not affected by BaseURL.
// Ref: TestBrowserContextBaseUrl.java#shouldNotConstructANewURLWhenValidURLsArePassed
func TestBrowserContextBaseURLNotConstructNewURLsForValidURLs(t *testing.T) {
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
