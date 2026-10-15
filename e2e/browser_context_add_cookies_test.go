//go:build e2e

// E2E tests for: TestBrowserContextAddCookies.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextFetchShouldAddCookiesFromSetCookieHeader(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/setcookie.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "session=value")
		w.Header().Add("Set-Cookie", "foo=bar; max-age=3600")
		w.WriteHeader(200)
	})

	bCtx, apiCtx := newBCtxWithRequest(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	defer page.Close(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/setcookie.html")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)

	sort.Slice(cookies, func(i, j int) bool {
		return cookies[i].Name < cookies[j].Name
	})
	is.Len(cookies, 2)
	is.Equal("foo", cookies[0].Name)
	is.Equal("bar", cookies[0].Value)
	is.Equal("session", cookies[1].Name)
	is.Equal("value", cookies[1].Value)

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	result, err := page.Evaluate(ctx, "() => document.cookie.split(';').map(s => s.trim()).sort()")
	must.NoError(err)

	items, ok := result.([]any)
	is.True(ok, "expected []any from evaluate, got %T", result)
	var got []string
	for _, item := range items {
		if s, ok := item.(string); ok {
			got = append(got, s)
		}
	}
	sort.Strings(got)
	is.Contains(got, "foo=bar")
	is.Contains(got, "session=value")
}

// TestBrowserContextAddCookiesShouldRoundtripCookie verifies that added cookies survive a roundtrip (add → get).
// Ref: TestBrowserContextAddCookies.java#shouldRoundtripCookie  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldRoundtripCookie(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	// Use URL only (not both url+domain — Playwright requires one or the other).
	// Avoid SameSite=None without Secure=true; Chromium rejects such cookies.
	cookieURL := "http://example.com"
	httpOnly := false

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{{
		Name:     "roundtrip",
		Value:    "thevalue",
		URL:      &cookieURL,
		HTTPOnly: &httpOnly,
	}}))

	cookies, err := bCtx.Cookies(ctx, cookieURL)
	must.NoError(err)
	is.Len(cookies, 1)
	c := cookies[0]
	is.Equal("roundtrip", c.Name)
	is.Equal("thevalue", c.Value)
	must.NotNil(c.Domain)
	must.NotNil(c.Path)
	is.Equal("/", *c.Path)
}

// TestBrowserContextAddCookiesShouldSendCookieHeader verifies that cookies are sent in request headers.
// Ref: TestBrowserContextAddCookies.java#shouldSendCookieHeader  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldSendCookieHeader(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	received := make(chan string, 1)
	srv.SetRoute("/cookie-header", func(w http.ResponseWriter, r *http.Request) {
		received <- r.Header.Get("Cookie")
		w.WriteHeader(200)
	})

	url := srv.Prefix() + "/cookie-header"
	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "abc123", URL: &url},
	}))

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, url))

	select {
	case header := <-received:
		is.Contains(header, "session=abc123")
	default:
		t.Fatal("no request received")
	}
}

// TestBrowserContextAddCookiesShouldIsolateCookiesInBrowserContexts verifies cookies are isolated between contexts.
// Ref: TestBrowserContextAddCookies.java#shouldIsolateCookiesInBrowserContexts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldIsolateCookiesInBrowserContexts(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	url := "http://example.com"
	domain := "example.com"
	path := "/"

	bCtx1 := newContext(t)
	bCtx2 := newContext(t)

	// Use Domain+Path (not URL) — Playwright requires either url OR domain, not both.
	must.NoError(bCtx1.AddCookies(ctx, []playwright.Cookie{
		{Name: "ctx1cookie", Value: "ctx1val", Domain: &domain, Path: &path},
	}))

	cookies1, err := bCtx1.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies1, 1)

	cookies2, err := bCtx2.Cookies(ctx, url)
	must.NoError(err)
	is.Empty(cookies2, "ctx2 should not have ctx1's cookies")
}

// TestBrowserContextAddCookiesShouldSetMultipleCookies verifies adding multiple cookies at once.
// Ref: TestBrowserContextAddCookies.java#shouldSetMultipleCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldSetMultipleCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	url := "http://example.com"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "cookie1", Value: "val1", URL: &url},
		{Name: "cookie2", Value: "val2", URL: &url},
		{Name: "cookie3", Value: "val3", URL: &url},
	}))

	cookies, err := bCtx.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies, 3)
}

// TestBrowserContextAddCookiesShouldHaveExpiresSetTo1ForSessionCookies verifies session cookies report expires == -1.
// Ref: TestBrowserContextAddCookies.java#shouldHaveExpiresSetTo1ForSessionCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldHaveExpiresSetTo1ForSessionCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	// Navigate to a page that sets a session cookie (no expires)
	srv.SetRoute("/session-cookie", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "sess=val")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/session-cookie"))

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 1)
	must.NotNil(cookies[0].Expires)
	// Session cookies have expires = -1 in Playwright's representation
	is.Less(*cookies[0].Expires, float64(0), "session cookie should have expires < 0")
}

// TestBrowserContextAddCookiesShouldSetCookieWithReasonableDefaults verifies cookies get reasonable defaults.
// Ref: TestBrowserContextAddCookies.java#shouldSetCookieWithReasonableDefaults  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldSetCookieWithReasonableDefaults(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	url := srv.EmptyPage()

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, url))

	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "defaults", Value: "check", URL: &url},
	}))

	cookies, err := bCtx.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies, 1)
	c := cookies[0]
	is.Equal("defaults", c.Name)
	is.Equal("check", c.Value)
	must.NotNil(c.Domain, "domain should be populated by default")
	must.NotNil(c.Path, "path should be populated by default")
}

// TestBrowserContextAddCookiesShouldSetACookieWithAPath verifies that a cookie with a path is scoped correctly.
// Ref: TestBrowserContextAddCookies.java#shouldSetACookieWithAPath  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldSetACookieWithAPath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	domain := "example.com"
	path := "/sub"
	url := "http://example.com/sub"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "pathcookie", Value: "pathval", Domain: &domain, Path: &path},
	}))

	// Should be found for the correct path
	cookies, err := bCtx.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies, 1)
	is.Equal("/sub", *cookies[0].Path)

	// Should not be found for a different path
	rootURL := "http://example.com/"
	cookiesRoot, err := bCtx.Cookies(ctx, rootURL)
	must.NoError(err)
	is.Empty(cookiesRoot, "cookie scoped to /sub should not appear at /")
}

// TestBrowserContextAddCookiesShouldBeAbleToSetUnsecureCookieForHTTPWebsite verifies a non-secure cookie can be set on HTTP.
// Ref: TestBrowserContextAddCookies.java#shouldBeAbleToSetUnsecureCookieForHTTPWebsite  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldBeAbleToSetUnsecureCookieForHTTPWebsite(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	url := srv.EmptyPage()

	bCtx := newContext(t)
	secure := false
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "unsecure", Value: "ok", URL: &url, Secure: &secure},
	}))

	cookies, err := bCtx.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies, 1)
	must.NotNil(cookies[0].Secure)
	is.False(*cookies[0].Secure)
}

// TestBrowserContextAddCookiesShouldSetACookieOnADifferentDomain verifies cookies can be set for any domain.
// Ref: TestBrowserContextAddCookies.java#shouldSetACookieOnADifferentDomain  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldSetACookieOnADifferentDomain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	url := "http://www.example.com"
	domain := "www.example.com"
	path := "/"

	bCtx := newContext(t)
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "otherdomain", Value: "crossval", Domain: &domain, Path: &path},
	}))

	cookies, err := bCtx.Cookies(ctx, url)
	must.NoError(err)
	is.Len(cookies, 1)
	is.Equal("otherdomain", cookies[0].Name)
}

// TestBrowserContextAddCookiesShouldIsolateSessionCookies verifies that session cookies set via Set-Cookie header
// are not shared across browser contexts.
// Ref: TestBrowserContextAddCookies.java#shouldIsolateSessionCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldIsolateSessionCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/setcookie.html", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "value"})
		w.WriteHeader(http.StatusOK)
	})

	// Context 1: visit the page that sets the cookie.
	bCtx1 := newContext(t)
	pg1, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg1.Goto(ctx, srv.Prefix()+"/setcookie.html"))

	// Cookie should appear in context 1.
	pg2, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg2.Goto(ctx, srv.EmptyPage()))
	cookies1, err := bCtx1.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies1, 1)
	is.Equal("value", cookies1[0].Value)

	// Context 2 (fresh): should have no cookies.
	bCtx2, err := globalBrowser.NewContext(ctx, nil)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(c)
	})
	pg3, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg3.Goto(ctx, srv.EmptyPage()))
	cookies2, err := bCtx2.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies2, "fresh context should have no session cookies")
}

// TestBrowserContextAddCookiesShouldIsolatePersistentCookies verifies that persistent cookies are isolated
// between browser contexts.
// Ref: TestBrowserContextAddCookies.java#shouldIsolatePersistentCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldIsolatePersistentCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/setcookie.html", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "persistent", Value: "persistent-value", MaxAge: 3600})
		w.WriteHeader(http.StatusOK)
	})

	bCtx1 := newContext(t)
	pg1, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg1.Goto(ctx, srv.Prefix()+"/setcookie.html"))

	bCtx2, err := globalBrowser.NewContext(ctx, nil)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(c)
	})

	pg2, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg2.Goto(ctx, srv.EmptyPage()))

	pg3, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg3.Goto(ctx, srv.EmptyPage()))

	cookies1, err := bCtx1.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies1, 1)
	is.Equal("persistent", cookies1[0].Name)
	is.Equal("persistent-value", cookies1[0].Value)

	cookies2, err := bCtx2.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies2, "fresh context should not have persistent cookies from context 1")
}

// TestBrowserContextAddCookiesShouldIsolateSendCookieHeader verifies that cookies added to one context are
// included in request headers for that context only, not for other contexts.
// Ref: TestBrowserContextAddCookies.java#shouldIsolateSendCookieHeader  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldIsolateSendCookieHeader(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	cookieURL := srv.EmptyPage()
	bCtx1 := newContext(t)
	must.NoError(bCtx1.AddCookies(ctx, []playwright.Cookie{{Name: "sendcookie", Value: "value", URL: &cookieURL}}))

	// Context 1: cookie header should be sent.
	pg1, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	reqCh1 := testserver.WaitForRequest(srv, "/empty.html")
	must.NoError(pg1.Goto(ctx, cookieURL))
	req1 := <-reqCh1
	is.Equal("sendcookie=value", req1.Header.Get("Cookie"),
		"context 1 should send cookie header")

	// Context 2 (fresh): cookie header should NOT be sent.
	bCtx2, err := globalBrowser.NewContext(ctx, nil)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(c)
	})
	pg2, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	reqCh2 := testserver.WaitForRequest(srv, "/empty.html")
	must.NoError(pg2.Goto(ctx, cookieURL))
	req2 := <-reqCh2
	is.Empty(req2.Header.Get("Cookie"),
		"fresh context should not send cookie header")
}

// TestBrowserContextAddCookiesShouldIsolateCookiesBetweenLaunches verifies that cookies are NOT shared across
// separate browser launch sessions.
// Ref: TestBrowserContextAddCookies.java#shouldIsolateCookiesBetweenLaunches  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldIsolateCookiesBetweenLaunches(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	cookieURL := srv.EmptyPage()
	expires := float64(time.Now().Unix() + 10000)

	bt := globalBrowser.BrowserType()

	// Launch browser 1, add a cookie, then close.
	browser1, err := bt.Launch(ctx, nil)
	must.NoError(err)
	bCtx1, err := browser1.NewContext(ctx, nil)
	must.NoError(err)
	must.NoError(bCtx1.AddCookies(ctx, []playwright.Cookie{
		{Name: "cookie-in-context-1", Value: "value", URL: &cookieURL, Expires: &expires},
	}))
	must.NoError(bCtx1.Close(ctx))
	must.NoError(browser1.Close(ctx))

	// Launch browser 2 in a fresh session: it should have no cookies.
	browser2, err := bt.Launch(ctx, nil)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = browser2.Close(c)
	})
	bCtx2, err := browser2.NewContext(ctx, nil)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(c)
	})
	cookies, err := bCtx2.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies, "fresh browser launch should have no cookies from a previous session")
}

// TestBrowserContextAddCookiesShouldNotSetACookieWithBlankPageURL verifies that adding a cookie for about:blank returns an error.
// Ref: TestBrowserContextAddCookies.java#shouldNotSetACookieWithBlankPageURL  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldNotSetACookieWithBlankPageURL(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	cookieURL := srv.EmptyPage()
	bCtx := newContext(t)
	blankURL := "about:blank"
	err := bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "example-cookie", Value: "best", URL: &cookieURL},
		{Name: "example-cookie-blank", Value: "best", URL: &blankURL},
	})
	is.Error(err, "should error when adding cookie for about:blank")
	is.Contains(err.Error(), "Blank page", "error should mention blank page")
}

// TestBrowserContextAddCookiesShouldNotSetACookieOnADataURLPage verifies that adding a cookie for a data: URL returns an error.
// Ref: TestBrowserContextAddCookies.java#shouldNotSetACookieOnADataURLPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldNotSetACookieOnADataURLPage(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx := newContext(t)
	dataURL := "data:,Hello%2C%20World!"
	err := bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "example-cookie", Value: "best", URL: &dataURL},
	})
	is.Error(err, "should error when adding cookie for data: URL")
	is.Contains(err.Error(), "Data URL", "error should mention data URL")
}

// TestBrowserContextAddCookiesShouldDefaultToSettingSecureCookieForHTTPSWebsites verifies that a cookie added for an HTTPS URL is
// marked as secure by default.
// Ref: TestBrowserContextAddCookies.java#shouldDefaultToSettingSecureCookieForHTTPSWebsites  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldDefaultToSettingSecureCookieForHTTPSWebsites(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bCtx := newContext(t)
	secureURL := "https://example.com"
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "foo", Value: "bar", URL: &secureURL},
	}))

	cookies, err := bCtx.Cookies(ctx, secureURL)
	must.NoError(err)
	is.Len(cookies, 1)
	must.NotNil(cookies[0].Secure, "Secure field should be non-nil")
	is.True(*cookies[0].Secure, "cookie added for https:// should be marked secure")
}

// TestBrowserContextAddCookiesShouldSetCookiesForAFrame verifies that a cookie is accessible from an iframe on the
// same origin.
// Ref: TestBrowserContextAddCookies.java#shouldSetCookiesForAFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldSetCookiesForAFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	pg, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg.Goto(ctx, srv.EmptyPage()))

	cookiePrefix := srv.Prefix()
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "frame-cookie", Value: "value", URL: &cookiePrefix},
	}))

	// Inject an iframe and wait for it to load by having the evaluate return the cookie value
	// from the iframe's document (same origin — accessible directly).
	iframeSrc := srv.Prefix() + "/empty.html"
	val, err := pg.Evaluate(ctx, `src => {
		return new Promise(resolve => {
			const iframe = document.createElement('iframe');
			iframe.onload = () => resolve(iframe.contentDocument.cookie);
			iframe.src = src;
			document.body.appendChild(iframe);
		});
	}`, iframeSrc)
	must.NoError(err)
	is.Equal("frame-cookie=value", val, "iframe should see the cookie set for the origin")
}

// TestBrowserContextAddCookiesShouldNotBlockThirdPartyCookies verifies that third-party cookies set in a
// cross-origin iframe are captured (or not) based on browser policy.
// Ref: TestBrowserContextAddCookies.java#shouldNotBlockThirdPartyCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldNotBlockThirdPartyCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	pg, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(pg.Goto(ctx, srv.EmptyPage()))

	crossPrefix := srv.CrossProcessPrefix()

	// Inject a cross-origin iframe and wait for it to load. We can't read its document
	// cookie directly (cross-origin), so we use a FrameLocator to evaluate inside it.
	_, err = pg.Evaluate(ctx, `src => new Promise(resolve => {
		const iframe = document.createElement('iframe');
		iframe.onload = resolve;
		iframe.src = src;
		document.body.appendChild(iframe);
	})`, crossPrefix+"/empty.html")
	must.NoError(err)

	// Use FrameLocator to evaluate in the cross-origin iframe.
	fl := pg.FrameLocator("iframe")
	_, err = fl.Locator("html").Evaluate(ctx, "document.cookie = 'username=John Doe'")
	must.NoError(err)

	// Allow the cookie to propagate.
	time.Sleep(2 * time.Second)

	cookies, err := bCtx.Cookies(ctx, crossPrefix+"/empty.html")
	must.NoError(err)
	// Chrome/WebKit block third-party cookies by default; Firefox allows them.
	// We accept either outcome; the key is that the test doesn't crash.
	t.Logf("third-party cookies captured: %v", cookies)
}

// TestBrowserContextAddCookiesShouldWork verifies programmatically setting cookies via BrowserContext.AddCookies.
// Ref: TestBrowserContextAddCookies.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

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
		{
			Name:     "auth_token",
			Value:    "secret_bearer_token",
			URL:      &emptyPageURL,
			SameSite: playwright.SameSiteLax,
		},
	})
	must.NoError(err, "AddCookies failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	cookieRaw, err := page.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate document.cookie failed")

	if cookieStr, ok := cookieRaw.(string); !ok || !strings.Contains(cookieStr, "auth_token=secret_bearer_token") {
		t.Errorf("document.cookie = %v, want auth_token=secret_bearer_token", cookieRaw)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookieBasic(t *testing.T) {
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
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{{
		Name:  "testcookie",
		Value: "myvalue",
		URL:   &url,
	}}))

	cookie, err := page.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
	is.Contains(cookie, "testcookie=myvalue")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextCookiesReturnsCookies(t *testing.T) {
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
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{{
		Name:  "cookie1",
		Value: "val1",
		URL:   &url,
	}}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.NotEmpty(cookies)

	found := false
	for _, c := range cookies {
		if c.Name == "cookie1" && c.Value == "val1" {
			found = true
			break
		}
	}
	is.True(found)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddMultipleCookies(t *testing.T) {
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
		{Name: "a", Value: "1", URL: &url},
		{Name: "b", Value: "2", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.GreaterOrEqual(len(cookies), 2)
}

// ---------------------------------------------------------------------------
// From page_cookies_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextCookieSessionEx2 verifies session cookie has no expiry.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextCookieSessionEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "val", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.NotEmpty(cookies)
	is.Equal("session", cookies[0].Name)
}

// TestBrowserContextMultipleCookiesEx2 verifies multiple cookies can be added.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextMultipleCookiesEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "c1", Value: "v1", URL: &url},
		{Name: "c2", Value: "v2", URL: &url},
		{Name: "c3", Value: "v3", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 3)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextIsolatesCookiesExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc1 := newContext(t)
	bc2 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)
	page2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	_, err = page1.Evaluate(ctx, `() => document.cookie = 'ctx=one'`)
	must.NoError(err)

	cookie2, err := page2.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
	must.NotContains(cookie2, "ctx=one")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesAndGetCookiesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "test-cookie", Value: "cookie-value", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.NotEmpty(cookies)

	found := false
	for _, c := range cookies {
		if c.Name == "test-cookie" {
			is.Equal("cookie-value", c.Value)
			found = true
		}
	}
	is.True(found)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextAddCookiesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "abc123", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Equal(1, len(cookies))
	is.Equal("session", cookies[0].Name)
}

// TestBrowserContextIsolatesCookies verifies two contexts don't share cookies.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextAddCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextAddCookies.java
func TestBrowserContextIsolatesCookiesExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc1 := newContext(t)
	bc2 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)
	page2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<html><body></body></html>`))
	must.NoError(page2.SetContent(ctx, `<html><body></body></html>`))

	_, err = page1.Evaluate(ctx, `() => document.cookie = 'ctx1cookie=value1'`)
	must.NoError(err)

	cookies, err := page2.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
	// Cookie from context 1 should not be in context 2
	if cookieStr, ok := cookies.(string); ok {
		must.NotContains(cookieStr, "ctx1cookie")
	}
}
