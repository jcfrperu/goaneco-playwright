//go:build e2e

// E2E tests for: TestPageSetExtraHttpHeaders.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestContextSetExtraHTTPHeadersEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	must.NoError(bc.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Context-Header": "context-value",
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<html><body>Content</body></html>`))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "Content")
}

// TestSetExtraHTTPHeadersSentWithRequest verifies extra headers are included in requests.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestSetExtraHTTPHeadersSentWithRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	capturedHeader := make(chan string, 1)

	srv.SetRoute("/check-headers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader <- r.Header.Get("X-Custom-Header")
		w.WriteHeader(200)
	}))

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Custom-Header": "custom-value",
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/check-headers"))

	select {
	case val := <-capturedHeader:
		is.Equal("custom-value", val)
	default:
		t.Fatal("header not captured")
	}
}

// TestSetExtraHTTPHeadersMultipleHeaders verifies multiple extra headers are sent.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestSetExtraHTTPHeadersMultipleHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	captured := make(chan map[string]string, 1)

	srv.SetRoute("/multi-headers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured <- map[string]string{
			"X-Header-One": r.Header.Get("X-Header-One"),
			"X-Header-Two": r.Header.Get("X-Header-Two"),
		}
		w.WriteHeader(200)
	}))

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Header-One": "value-one",
		"X-Header-Two": "value-two",
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/multi-headers"))

	select {
	case headers := <-captured:
		is.Equal("value-one", headers["X-Header-One"])
		is.Equal("value-two", headers["X-Header-Two"])
	default:
		t.Fatal("headers not captured")
	}
}

// TestSetExtraHTTPHeadersOverridePreviousHeaders verifies SetExtraHTTPHeaders replaces previous headers.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestSetExtraHTTPHeadersOverridePreviousHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	capturedHeader := make(chan string, 2)

	srv.SetRoute("/override-check", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader <- r.Header.Get("X-Test")
		w.WriteHeader(200)
	}))

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{"X-Test": "first"}))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/override-check"))
	first := <-capturedHeader
	is.Equal("first", first)

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{"X-Test": "second"}))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/override-check"))
	second := <-capturedHeader
	is.Equal("second", second)
}

// TestSetExtraHTTPHeadersSentToServerEx2 verifies extra headers are sent to server.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestSetExtraHTTPHeadersSentToServerEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var receivedHeader string
	srv.SetRoute("/headers", func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Custom-Token")
		w.WriteHeader(http.StatusOK)
	})

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Custom-Token": "secret-token",
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/headers"))

	is.Equal("secret-token", receivedHeader)
}

// TestSetExtraHTTPHeadersMultipleHeadersEx2 verifies multiple extra headers are sent.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestSetExtraHTTPHeadersMultipleHeadersEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	headers := map[string]string{}
	srv.SetRoute("/multi", func(w http.ResponseWriter, r *http.Request) {
		headers["X-Header-A"] = r.Header.Get("X-Header-A")
		headers["X-Header-B"] = r.Header.Get("X-Header-B")
		w.WriteHeader(http.StatusOK)
	})

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Header-A": "value-a",
		"X-Header-B": "value-b",
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/multi"))

	is.Equal("value-a", headers["X-Header-A"])
	is.Equal("value-b", headers["X-Header-B"])
}

// TestSetExtraHTTPHeadersOverridesOnNextNavigationEx2 verifies headers persist across navigations.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestSetExtraHTTPHeadersOverridesOnNextNavigationEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	count := 0
	srv.SetRoute("/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Persistent") == "yes" {
			count++
		}
		w.WriteHeader(http.StatusOK)
	})

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Persistent": "yes",
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/count"))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/count"))

	is.Equal(2, count)
}

// TestSetExtraHTTPHeadersEmptyMapEx2 verifies setting empty map does not error.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestSetExtraHTTPHeadersEmptyMapEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{}))
}

// TestExtraHeadersOverrideEx3 verifies SetExtraHTTPHeaders override works.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestExtraHeadersOverrideEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var receivedAuth string
	var mu sync.Mutex

	srv.SetRoute("/protected", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		receivedAuth = r.Header.Get("Authorization")
		mu.Unlock()
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<p>Protected</p>`)
	})

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{
		"Authorization": "Bearer token123",
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/protected"))

	mu.Lock()
	auth := receivedAuth
	mu.Unlock()

	is.Equal("Bearer token123", auth)
}

// TestExtraHeadersMultipleEx3 verifies multiple extra headers are sent.
// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestExtraHeadersMultipleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var h1, h2 string
	var mu sync.Mutex

	srv.SetRoute("/multi", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		h1 = r.Header.Get("X-Header-1")
		h2 = r.Header.Get("X-Header-2")
		mu.Unlock()
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, `<p>Multi</p>`)
	})

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/multi"))

	mu.Lock()
	v1 := h1
	v2 := h2
	mu.Unlock()

	is.Equal("value1", v1)
	is.Equal("value2", v2)
}

// Extra: no equivalent method in playwright-java; related class: TestPageSetExtraHttpHeaders.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestPageSetExtraHTTPHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	gotHeader := make(chan string, 1)
	srv.SetRoute("/page-headers", func(w http.ResponseWriter, r *http.Request) {
		gotHeader <- r.Header.Get("X-Page-Header")
		w.WriteHeader(200)
	})

	page := newPage(t)
	err := page.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Page-Header": "page-value",
	})
	must.NoError(err, "SetExtraHTTPHeaders failed")

	err = page.Goto(ctx, srv.Prefix()+"/page-headers")
	must.NoError(err, "Goto failed")

	select {
	case val := <-gotHeader:
		if val != "page-value" {
			t.Errorf("X-Page-Header = %q, want 'page-value'", val)
		}
	default:
		t.Error("request never reached the server")
	}
}

// --- Page.SetViewportSize / Page.ViewportSize ---

// Ref: TestPageSetExtraHttpHeaders.java#shouldWorkWithExtraHeadersFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestPageSetExtraHttpHeadersShouldWorkWithExtraHeadersFromBrowserContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	gotHeader := make(chan string, 1)
	srv.SetRoute("/check-headers", func(w http.ResponseWriter, r *http.Request) {
		gotHeader <- r.Header.Get("X-Custom-Header")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = bCtx.SetExtraHTTPHeaders(ctx, map[string]string{
		"X-Custom-Header": "playwright-test",
	})
	must.NoError(err, "SetExtraHTTPHeaders failed")

	err = page.Goto(ctx, srv.Prefix()+"/check-headers")
	must.NoError(err, "Goto failed")

	select {
	case val := <-gotHeader:
		if val != "playwright-test" {
			t.Errorf("X-Custom-Header = %q, want 'playwright-test'", val)
		}
	default:
		t.Error("request never reached the server")
	}
}

// TestPageSetExtraHttpHeadersShouldWork verifies page extra headers are sent.
// Ref: TestPageSetExtraHttpHeaders.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestPageSetExtraHttpHeadersShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{"foo": "bar"}))
	rec := a4RecordRoute(srv, "/empty.html")
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	rec.Wait(t)
	is.Equal([]string{"bar"}, rec.Header.Values("foo"))
	is.Empty(rec.Header.Values("baz"))
}

// TestPageSetExtraHttpHeadersShouldWorkWithRedirects verifies extra headers survive redirects.
// Ref: TestPageSetExtraHttpHeaders.java#shouldWorkWithRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestPageSetExtraHttpHeadersShouldWorkWithRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	srv.ServeWithRedirect("/foo.html", "/empty.html")
	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{"foo": "bar"}))
	rec := a4RecordRoute(srv, "/empty.html")
	must.NoError(page.Goto(ctx, srv.Prefix()+"/foo.html"))
	rec.Wait(t)
	is.Equal([]string{"bar"}, rec.Header.Values("foo"))
}

// TestPageSetExtraHttpHeadersShouldOverrideExtraHeadersFromBrowserContext verifies page headers override context headers.
// Ref: TestPageSetExtraHttpHeaders.java#shouldOverrideExtraHeadersFromBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestPageSetExtraHttpHeadersShouldOverrideExtraHeadersFromBrowserContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := a4NewPage(t, &playwright.BrowserContextOptions{ExtraHTTPHeaders: map[string]string{"fOo": "bAr", "baR": "foO"}})
	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{"Foo": "Bar"}))
	rec := a4RecordRoute(srv, "/empty.html")
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	rec.Wait(t)
	is.Equal([]string{"Bar"}, rec.Header.Values("foo"))
	is.Equal([]string{"foO"}, rec.Header.Values("bar"))
}

// TestPageSetExtraHttpHeadersShouldThrowForNonStringHeaderValues is not ported.
// Ref: TestPageSetExtraHttpHeaders.java#shouldThrowForNonStringHeaderValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetExtraHttpHeaders.java
func TestPageSetExtraHttpHeadersShouldThrowForNonStringHeaderValues(t *testing.T) {
	t.Skip("not implemented: null header values (ExtraHTTPHeaders is map[string]string, a nil value cannot be expressed)")
}
