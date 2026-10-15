//go:build e2e

// E2E tests for: TestBrowserContextClearCookies.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"regexp"
	"testing"
	"time"
)

// TestBrowserContextClearCookiesShouldIsolateCookiesWhenClearing verifies clearing only affects the target context.
// Ref: TestBrowserContextClearCookies.java#shouldIsolateCookiesWhenClearing  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesShouldIsolateCookiesWhenClearing(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	url := "http://example.com"
	domain := "example.com"
	path := "/"

	bCtx1 := newContext(t)
	bCtx2 := newContext(t)

	must.NoError(bCtx1.AddCookies(ctx, []playwright.Cookie{
		{Name: "ctx1", Value: "v1", Domain: &domain, Path: &path},
	}))
	must.NoError(bCtx2.AddCookies(ctx, []playwright.Cookie{
		{Name: "ctx2", Value: "v2", Domain: &domain, Path: &path},
	}))

	// Clear only ctx1
	must.NoError(bCtx1.ClearCookies(ctx))

	cookies1, err := bCtx1.Cookies(ctx, url)
	must.NoError(err)
	is.Empty(cookies1, "ctx1 cookies should be cleared")

	cookies2, err := bCtx2.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies2, 1, "ctx2 cookies should be unaffected")
}

// TestBrowserContextClearCookiesShouldRemoveCookiesByName verifies cookies are removed by name filter.
// Ref: TestBrowserContextClearCookies.java#shouldRemoveCookiesByName  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesShouldRemoveCookiesByName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	url := "http://example.com"
	domain := "example.com"
	path := "/"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "keep", Value: "k", Domain: &domain, Path: &path},
		{Name: "remove", Value: "r", Domain: &domain, Path: &path},
	}))

	name := "remove"
	must.NoError(bCtx.ClearCookies(ctx, &playwright.ClearCookiesOptions{Name: &name}))

	cookies, err := bCtx.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies, 1)
	is.Equal("keep", cookies[0].Name)
}

// TestBrowserContextClearCookiesShouldRemoveCookiesByDomain verifies cookies are removed by domain filter.
// Ref: TestBrowserContextClearCookies.java#shouldRemoveCookiesByDomain  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesShouldRemoveCookiesByDomain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	url1 := "http://foo.com"
	url2 := "http://bar.com"
	domainFoo := "foo.com"
	domainBar := "bar.com"
	path := "/"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "foocookie", Value: "fv", Domain: &domainFoo, Path: &path},
		{Name: "barcookie", Value: "bv", Domain: &domainBar, Path: &path},
	}))

	must.NoError(bCtx.ClearCookies(ctx, &playwright.ClearCookiesOptions{Domain: &domainFoo}))

	fooC, err := bCtx.Cookies(ctx, url1)
	must.NoError(err)
	is.Empty(fooC, "foo.com cookies should be cleared")

	barC, err := bCtx.Cookies(ctx, url2)
	must.NoError(err)
	is.Len(barC, 1, "bar.com cookies should remain")
}

// TestBrowserContextClearCookiesShouldRemoveCookiesByPath verifies cookies are removed by path filter.
// Ref: TestBrowserContextClearCookies.java#shouldRemoveCookiesByPath  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesShouldRemoveCookiesByPath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	domain := "example.com"
	path1 := "/admin"
	path2 := "/user"
	url1 := "http://example.com/admin"
	url2 := "http://example.com/user"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "admin", Value: "av", Domain: &domain, Path: &path1},
		{Name: "user", Value: "uv", Domain: &domain, Path: &path2},
	}))

	must.NoError(bCtx.ClearCookies(ctx, &playwright.ClearCookiesOptions{Path: &path1}))

	admin, err := bCtx.Cookies(ctx, url1)
	must.NoError(err)
	is.Empty(admin, "/admin cookies should be cleared")

	user, err := bCtx.Cookies(ctx, url2)
	must.NoError(err)
	is.Len(user, 1, "/user cookies should remain")
}

// TestBrowserContextClearCookiesShouldRemoveCookiesByNameAndDomain verifies cookies removed when both name and domain match.
// Ref: TestBrowserContextClearCookies.java#shouldRemoveCookiesByNameAndDomain  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesShouldRemoveCookiesByNameAndDomain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	domainA := "aaa.com"
	domainB := "bbb.com"
	path := "/"
	urlA := "http://aaa.com"
	urlB := "http://bbb.com"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "target", Value: "v1", Domain: &domainA, Path: &path},
		{Name: "target", Value: "v2", Domain: &domainB, Path: &path},
	}))

	targetName := "target"
	must.NoError(bCtx.ClearCookies(ctx, &playwright.ClearCookiesOptions{
		Name:   &targetName,
		Domain: &domainA,
	}))

	cA, err := bCtx.Cookies(ctx, urlA)
	must.NoError(err)
	is.Empty(cA, "aaa.com target cookie should be cleared")

	cB, err := bCtx.Cookies(ctx, urlB)
	must.NoError(err)
	is.Len(cB, 1, "bbb.com target cookie should remain")
}

// TestBrowserContextClearCookiesShouldRemoveCookiesByNameRegex verifies cookies are removed by name regex filter.
// Ref: TestBrowserContextClearCookies.java#shouldRemoveCookiesByNameRegex  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesShouldRemoveCookiesByNameRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	parsed, err := url.Parse(srv.Prefix())
	must.NoError(err)
	host := parsed.Hostname()
	path := "/"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "cookie1", Value: "1", Domain: &host, Path: &path},
		{Name: "cookie2", Value: "2", Domain: &host, Path: &path},
	}))

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()))

	// Both cookies initially visible
	docCookie, err := page.Evaluate(ctx, "document.cookie")
	must.NoError(err)
	is.Contains(docCookie.(string), "cookie1=1")
	is.Contains(docCookie.(string), "cookie2=2")

	// Clear only cookies whose name matches the regex coo.*1
	must.NoError(bCtx.ClearCookies(ctx, &playwright.ClearCookiesOptions{
		NameRegex: regexp.MustCompile("coo.*1"),
	}))

	// Only cookie2 should remain
	docCookie2, err := page.Evaluate(ctx, "document.cookie")
	must.NoError(err)
	is.Equal("cookie2=2", docCookie2)
}

// TestBrowserContextClearCookiesShouldClearCookies verifies clearing cookies from the BrowserContext.
// Ref: TestBrowserContextClearCookies.java#shouldClearCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesShouldClearCookies(t *testing.T) {
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

	emptyPageURL := srv.EmptyPage()
	err = bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "c1", Value: "v1", URL: &emptyPageURL},
		{Name: "c2", Value: "v2", URL: &emptyPageURL},
	})
	must.NoError(err, "AddCookies failed")

	cookies, err := bCtx.Cookies(ctx)
	if err != nil || len(cookies) < 2 {
		t.Fatalf("Cookies failed or returned < 2 cookies: %v, %v", cookies, err)
	}

	// Clear all cookies
	err = bCtx.ClearCookies(ctx)
	must.NoError(err, "ClearCookies failed")

	cookiesAfter, err := bCtx.Cookies(ctx)
	must.NoError(err, "Cookies after ClearCookies failed")
	is.Len(cookiesAfter, 0)
}

// TestBrowserContextClearCookiesEx2 verifies ClearCookies removes all cookies.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextClearCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	url := srv.EmptyPage()
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "toclear", Value: "yes", URL: &url},
	}))

	must.NoError(bc.ClearCookies(ctx))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)
}

// TestPageCookiesClearEx3 verifies clearing removes all cookies.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextClearCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestPageCookiesClearEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "a", Value: "1", URL: &url},
		{Name: "b", Value: "2", URL: &url},
		{Name: "c", Value: "3", URL: &url},
	}))

	must.NoError(bc.ClearCookies(ctx))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)
}

// TestClearCookiesEx4 verifies ClearCookies removes all cookies.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextClearCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestClearCookiesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "c1", Value: "v1", Domain: localStrPtrCk4("localhost"), Path: localStrPtrCk4("/")},
		{Name: "c2", Value: "v2", Domain: localStrPtrCk4("localhost"), Path: localStrPtrCk4("/")},
	}))

	must.NoError(bc.ClearCookies(ctx))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)
	_ = page
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextClearCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextClearCookies.java
func TestBrowserContextClearCookiesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "tok", Value: "val", URL: &url},
	}))

	must.NoError(bc.ClearCookies(ctx))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)
}
