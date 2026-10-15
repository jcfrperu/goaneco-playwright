//go:build e2e

// E2E tests for: TestBrowserContextCookies.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

// TestBrowserContextGetCookiesEx2 verifies Cookies returns added cookies.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextGetCookiesEx2(t *testing.T) {
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
		{Name: "mycookie", Value: "myvalue", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.NotEmpty(cookies)

	found := false
	for _, c := range cookies {
		if c.Name == "mycookie" && c.Value == "myvalue" {
			found = true
		}
	}
	is.True(found)
}

// TestPageCookiesMultipleEx3 verifies multiple cookies can be set.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestPageCookiesMultipleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "c1", Value: "v1", URL: &url},
		{Name: "c2", Value: "v2", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Equal(2, len(cookies))
}

// TestPageCookiesValueEx3 verifies cookie value is preserved.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestPageCookiesValueEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "secret123", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Equal(1, len(cookies))
	is.Equal("secret123", cookies[0].Value)
}

// TestPageCookiesNameEx3 verifies cookie name is preserved.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestPageCookiesNameEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	url := "https://example.com"
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "mytoken", Value: "xyz", URL: &url},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Equal("mytoken", cookies[0].Name)
}

// TestSetCookieWithPathEx4 verifies AddCookies with path attribute.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestSetCookieWithPathEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "abc123", Domain: localStrPtrCk4("localhost"), Path: localStrPtrCk4("/app")},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)

	found := false
	for _, c := range cookies {
		if c.Name == "session" && c.Value == "abc123" {
			found = true
			break
		}
	}
	is.True(found)
	_ = page
}

// TestMultipleCookiesEx4 verifies AddCookies can add multiple cookies at once.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCookies.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestMultipleCookiesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{
		{Name: "auth", Value: "token123", Domain: localStrPtrCk4("localhost"), Path: localStrPtrCk4("/")},
		{Name: "pref", Value: "dark-mode", Domain: localStrPtrCk4("localhost"), Path: localStrPtrCk4("/")},
		{Name: "lang", Value: "en", Domain: localStrPtrCk4("localhost"), Path: localStrPtrCk4("/")},
	}))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.GreaterOrEqual(len(cookies), 3)
	_ = page
}

// TestBrowserContextCookiesShouldGetANonSessionCookie verifies a cookie with MaxAge is reported with expires > 0.
// Ref: TestBrowserContextCookies.java#shouldGetANonSessionCookie  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextCookiesShouldGetANonSessionCookie(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/set-non-session", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "user", Value: "John", MaxAge: 3600})
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/set-non-session"))

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 1)
	is.Equal("user", cookies[0].Name)
	is.Equal("John", cookies[0].Value)
	must.NotNil(cookies[0].Expires)
	is.Greater(*cookies[0].Expires, float64(0), "non-session cookie should have positive expires")
}

// TestBrowserContextCookiesShouldProperlyReportHttpOnlyCookie verifies that HttpOnly cookies are reported correctly.
// Ref: TestBrowserContextCookies.java#shouldProperlyReportHttpOnlyCookie  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextCookiesShouldProperlyReportHttpOnlyCookie(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/http-only", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "secret=value; HttpOnly")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/http-only"))

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 1)
	is.Equal("secret", cookies[0].Name)
	must.NotNil(cookies[0].HTTPOnly)
	is.True(*cookies[0].HTTPOnly, "HttpOnly cookie should report httpOnly=true")
}

// TestBrowserContextCookiesShouldProperlyReportStrictSameSiteCookie verifies SameSite=Strict is reported correctly.
// Ref: TestBrowserContextCookies.java#shouldProperlyReportStrictSameSiteCookie  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextCookiesShouldProperlyReportStrictSameSiteCookie(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/strict-samesite", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "strict=val; SameSite=Strict")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/strict-samesite"))

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 1)
	is.Equal(playwright.SameSiteStrict, cookies[0].SameSite, "expected SameSite=Strict")
}

// TestBrowserContextCookiesShouldProperlyReportLaxSameSiteCookie verifies SameSite=Lax is reported correctly.
// Ref: TestBrowserContextCookies.java#shouldProperlyReportLaxSameSiteCookie  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextCookiesShouldProperlyReportLaxSameSiteCookie(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/lax-samesite", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "lax=val; SameSite=Lax")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/lax-samesite"))

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 1)
	is.Equal(playwright.SameSiteLax, cookies[0].SameSite, "expected SameSite=Lax")
}

// TestBrowserContextCookiesShouldGetMultipleCookies verifies multiple cookies are returned.
// Ref: TestBrowserContextCookies.java#shouldGetMultipleCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextCookiesShouldGetMultipleCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/multi-cookie", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "first=one")
		w.Header().Add("Set-Cookie", "second=two")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/multi-cookie"))

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 2)

	names := map[string]string{}
	for _, c := range cookies {
		names[c.Name] = c.Value
	}
	is.Equal("one", names["first"])
	is.Equal("two", names["second"])
}

// TestBrowserContextCookiesShouldGetCookiesFromMultipleUrls verifies cookies can be retrieved for multiple URLs.
// Ref: TestBrowserContextCookies.java#shouldGetCookiesFromMultipleUrls  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextCookiesShouldGetCookiesFromMultipleUrls(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	parsed, err := url.Parse(srv.Prefix())
	must.NoError(err)
	host := parsed.Hostname() // e.g., "127.0.0.1" (no port — cookies don't include port)

	url1 := srv.Prefix() + "/path1"
	url2 := srv.Prefix() + "/path2"
	path1 := "/path1"
	path2 := "/path2"
	bCtx := newContext(t)

	// Use Domain+Path so cookies are path-scoped to their respective URLs.
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "cookie1", Value: "val1", Domain: &host, Path: &path1},
		{Name: "cookie2", Value: "val2", Domain: &host, Path: &path2},
	}))

	cookies1, err := bCtx.Cookies(ctx, url1)
	must.NoError(err)
	is.Len(cookies1, 1)
	is.Equal("cookie1", cookies1[0].Name)

	cookies2, err := bCtx.Cookies(ctx, url2)
	must.NoError(err)
	is.Len(cookies2, 1)
	is.Equal("cookie2", cookies2[0].Name)

	all, err := bCtx.Cookies(ctx, url1, url2)
	must.NoError(err)
	is.Len(all, 2)
}

// TestBrowserContextCookiesShouldGetACookie verifies that a cookie set via document.cookie is visible via context.cookies().
// Ref: TestBrowserContextCookies.java#shouldGetACookie  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCookies.java
func TestBrowserContextCookiesShouldGetACookie(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	docCookie, err := page.Evaluate(ctx, `() => {
		document.cookie = 'username=John Doe';
		return document.cookie;
	}`)
	must.NoError(err)
	is.Equal("username=John Doe", docCookie)

	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.Len(cookies, 1)
	c := cookies[0]
	is.Equal("username", c.Name)
	is.Equal("John Doe", c.Value)
	must.NotNil(c.Domain)
	must.NotNil(c.Path)
	is.Equal("/", *c.Path)
	must.NotNil(c.Expires)
	is.Less(*c.Expires, float64(0), "session cookie should have expires < 0")
	must.NotNil(c.HTTPOnly)
	is.False(*c.HTTPOnly)
	must.NotNil(c.Secure)
	is.False(*c.Secure)
}
