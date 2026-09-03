//go:build e2e

// Page.SetExtraHTTPHeaders E2E tests.
// Migration of: TestPageExtraHeaders.java, TestPageSetExtraHTTPHeaders.java
package e2e

import (
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetExtraHTTPHeadersSentWithRequest verifies extra headers are included in requests.
// Ref: TestPageExtraHeaders.java#shouldWork
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
// Ref: TestPageExtraHeaders.java#shouldWorkWithMultipleHeaders
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
// Ref: TestPageExtraHeaders.java#shouldOverridePreviousHeaders
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
// Ref: TestPageSetExtraHTTPHeaders.java#shouldSendHeaders
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
// Ref: TestPageSetExtraHTTPHeaders.java#shouldSendMultipleHeaders
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
// Ref: TestPageSetExtraHTTPHeaders.java#shouldPersistAcrossNavigation
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
// Ref: TestPageSetExtraHTTPHeaders.java#shouldAllowEmptyMap
func TestSetExtraHTTPHeadersEmptyMapEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetExtraHTTPHeaders(ctx, map[string]string{}))
}

// TestExtraHeadersOverrideEx3 verifies SetExtraHTTPHeaders override works.
// Ref: TestPageExtraHeaders.java#shouldOverrideHeader
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
// Ref: TestPageExtraHeaders.java#shouldSendMultipleHeaders
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
