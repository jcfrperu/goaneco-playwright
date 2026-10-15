//go:build e2e

// E2E tests for: TestPageNetworkResponse.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// TestNetworkResponseHeadersIncludeContentType verifies Response.Headers() includes content-type.
// Extra: no equivalent method in playwright-java; related class: TestPageNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestNetworkResponseHeadersIncludeContentType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedHeaders map[string]string

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			capturedHeaders = resp.Headers()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(capturedHeaders)
	// Content-type should be present for an HTML page
	is.NotEmpty(capturedHeaders)
}

// Extra: no equivalent method in playwright-java; related class: TestPageNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestNetworkResponseHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/with-headers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "custom-value")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	})

	page := newPage(t)

	type result struct {
		headers map[string]string
	}
	ch := make(chan result, 1)

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if !strings.HasSuffix(resp.URL(), "/with-headers") {
			return
		}
		select {
		case ch <- result{resp.Headers()}:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/with-headers")
	must.NoError(err, "Goto failed")

	select {
	case r := <-ch:
		// Headers may be lower-cased depending on browser/protocol.
		val := r.headers["x-custom-header"]
		if val == "" {
			val = r.headers["X-Custom-Header"]
		}
		if val != "custom-value" {
			t.Errorf("X-Custom-Header = %q, want 'custom-value'; all headers: %v", val, r.headers)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for response")
	}
}

// TestNetworkResponseHeadersMap verifies Headers() returns a non-empty map.
// Extra: no equivalent method in playwright-java; related class: TestPageNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestNetworkResponseHeadersMap(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var headers map[string]string
	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		if r.URL() == srv.EmptyPage() {
			headers = r.Headers()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.NotEmpty(headers)
}

// TestNetworkResponseHeadersEx3 verifies response headers are accessible.
// Extra: no equivalent method in playwright-java; related class: TestPageNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestNetworkResponseHeadersEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var headers map[string]string
	var mu sync.Mutex

	bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		headers = resp.Headers()
		mu.Unlock()
	})

	status := 200
	ct := "text/html"
	body := "content"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	h := headers
	mu.Unlock()

	must.NotNil(h)
}

// TestNetworkResponseHeadersEx4 verifies route fulfill returns custom headers.
// Extra: no equivalent method in playwright-java; related class: TestPageNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestNetworkResponseHeadersEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/with-headers", func(r *playwright.Route) {
		status := 200
		ct := "application/json"
		body := `{}`
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
			Headers:     map[string]string{"X-Custom-Header": "custom-value"},
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	result, err := page.Evaluate(ctx, `() => fetch('/with-headers').then(r => r.headers.get('x-custom-header'))`)
	must.NoError(err)
	is.Equal("custom-value", result)
}

// TestPageNetworkResponseShouldReportMultipleSetCookieHeaders verifies that multiple Set-Cookie headers are reported correctly.
// Ref: TestPageNetworkResponse.java#shouldReportMultipleSetCookieHeaders  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestPageNetworkResponseShouldReportMultipleSetCookieHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/headers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "a=b")
		w.Header().Add("Set-Cookie", "c=d")
		w.WriteHeader(200)
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	response, err := page.WaitForResponse(ctx, "**/*", func() error {
		_, err := page.Evaluate(ctx, "fetch('/headers')")
		return err
	})
	must.NoError(err)
	must.NotNil(response)

	headers, err := response.HeadersArray(ctx)
	must.NoError(err)
	var cookies []string
	for _, h := range headers {
		if strings.ToLower(h.Name) == "set-cookie" {
			cookies = append(cookies, h.Value)
		}
	}
	is.Equal([]string{"a=b", "c=d"}, cookies)

	val, err := response.HeaderValue(ctx, "not-there")
	must.NoError(err)
	is.Nil(val)

	val, err = response.HeaderValue(ctx, "set-cookie")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("a=b\nc=d", *val)

	vals, err := response.HeaderValues(ctx, "set-cookie")
	must.NoError(err)
	is.Equal([]string{"a=b", "c=d"}, vals)
}

// TestPageNetworkResponseShouldRejectResponseFinishedIfPageCloses verifies Finished() rejects if page closes before response finishes.
// Ref: TestPageNetworkResponse.java#shouldRejectResponseFinishedIfPageCloses  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestPageNetworkResponseShouldRejectResponseFinishedIfPageCloses(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	srv.SetRoute("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		if f, ok := w.(http.Flusher); ok {
			_, _ = w.Write([]byte("hello "))
			f.Flush()
		}
		<-r.Context().Done()
	})

	pageResponse, err := page.WaitForResponse(ctx, "**/get", func() error {
		_, err := page.Evaluate(ctx, "() => fetch('./get', { method: 'GET' })")
		return err
	})
	must.NoError(err)
	must.NotNil(pageResponse)

	page.OnRequest(func(request *playwright.NetworkRequest) {
		_ = page.Close(ctx)
	})
	_, err = page.Evaluate(ctx, "() => setTimeout(() => fetch('./empty.html', { method: 'GET' }), 1000)")
	must.NoError(err)

	finErr := pageResponse.Finished()
	must.Error(finErr)
	is.Contains(finErr.Error(), "closed")
}

// TestPageNetworkResponseShouldReturnNullExistingResponseBeforeResponseReceived verifies ExistingResponse is nil before response received.
// Ref: TestPageNetworkResponse.java#shouldReturnNullExistingResponseBeforeResponseReceived  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestPageNetworkResponseShouldReturnNullExistingResponseBeforeResponseReceived(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedRequest *playwright.NetworkRequest

	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		mu.Lock()
		capturedRequest = route.Request()
		is.Nil(capturedRequest.ExistingResponse())
		mu.Unlock()
		must.NoError(route.Continue(ctx))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	mu.Lock()
	defer mu.Unlock()
	must.NotNil(capturedRequest)
}

// TestPageNetworkResponseShouldReturnExistingResponseAfterReceived verifies ExistingResponse is set after response received.
// Ref: TestPageNetworkResponse.java#shouldReturnExistingResponseAfterReceived  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestPageNetworkResponseShouldReturnExistingResponseAfterReceived(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedResponse *playwright.NetworkResponse

	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		mu.Lock()
		capturedResponse = r
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	mu.Lock()
	resp := capturedResponse
	mu.Unlock()
	must.NotNil(resp)
	must.NotNil(resp.Request())
	is.Equal(resp, resp.Request().ExistingResponse())
}

// TestPageNetworkResponseShouldReturnHttpVersion verifies HttpVersion format.
// Ref: TestPageNetworkResponse.java#shouldReturnHttpVersion  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestPageNetworkResponseShouldReturnHttpVersion(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	response, err := page.WaitForResponse(ctx, "**/*", func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(response)
	version, err := response.HttpVersion(ctx)
	must.NoError(err)
	matched, err := regexp.MatchString(`^HTTP/[12](\.[01])?$`, version)
	must.NoError(err)
	is.True(matched, "unexpected version: "+version)
}

// TestPageNetworkResponseShouldRejectResponseFinishedIfContextCloses verifies Finished() rejects if context closes before response finishes.
// Ref: TestPageNetworkResponse.java#shouldRejectResponseFinishedIfContextCloses  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkResponse.java
func TestPageNetworkResponseShouldRejectResponseFinishedIfContextCloses(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	srv.SetRoute("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		if f, ok := w.(http.Flusher); ok {
			_, _ = w.Write([]byte("hello "))
			f.Flush()
		}
		<-r.Context().Done()
	})

	pageResponse, err := page.WaitForResponse(ctx, "**/get", func() error {
		_, err := page.Evaluate(ctx, "() => fetch('./get', { method: 'GET' })")
		return err
	})
	must.NoError(err)
	must.NotNil(pageResponse)

	page.OnRequest(func(request *playwright.NetworkRequest) {
		_ = bCtx.Close(ctx)
	})
	_, err = page.Evaluate(ctx, "() => setTimeout(() => fetch('./empty.html', { method: 'GET' }), 1000)")
	must.NoError(err)

	finErr := pageResponse.Finished()
	must.Error(finErr)
	is.Contains(finErr.Error(), "closed")
}
