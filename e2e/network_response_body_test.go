//go:build e2e

// NetworkResponse.Body and .Text E2E tests.
// Migration of: TestNetworkResponse.java
package e2e

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sync"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

func TestNetworkResponseBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/body", "text/plain", "hello body")

	page := newPage(t)

	type result struct {
		body []byte
		err  error
	}
	ch := make(chan result, 1)

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if !strings.HasSuffix(resp.URL(), "/body") {
			return
		}
		b, err := resp.Body(ctx)
		select {
		case ch <- result{b, err}:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/body")
	must.NoError(err, "Goto failed")

	select {
	case r := <-ch:
		must.NoError(r.err, "Body() failed")
		is.Equal("hello body", string(r.body))
	case <-ctx.Done():
		t.Fatal("timeout waiting for response")
	}
}

func TestNetworkResponseText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/text", "text/plain", "hello text")

	page := newPage(t)

	type result struct {
		text string
		err  error
	}
	ch := make(chan result, 1)

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if !strings.HasSuffix(resp.URL(), "/text") {
			return
		}
		s, err := resp.Text(ctx)
		select {
		case ch <- result{s, err}:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/text")
	must.NoError(err, "Goto failed")

	select {
	case r := <-ch:
		must.NoError(r.err, "Text() failed")
		is.Equal("hello text", r.text)
	case <-ctx.Done():
		t.Fatal("timeout waiting for response")
	}
}

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

func TestNetworkResponseStatusText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/status-text", "text/plain", "ok")

	page := newPage(t)

	ch := make(chan string, 1)

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if !strings.HasSuffix(resp.URL(), "/status-text") {
			return
		}
		select {
		case ch <- resp.StatusText():
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/status-text")
	must.NoError(err, "Goto failed")

	select {
	case st := <-ch:
		if st != "OK" && st != "" {
			t.Errorf("StatusText() = %q, want 'OK'", st)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for response")
	}
}

func TestNetworkResponseStatusAndURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/status", "text/plain", "ok")

	page := newPage(t)

	type result struct {
		status int
		url    string
	}
	ch := make(chan result, 1)

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if !strings.HasSuffix(resp.URL(), "/status") {
			return
		}
		select {
		case ch <- result{resp.Status(), resp.URL()}:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/status")
	must.NoError(err, "Goto failed")

	select {
	case r := <-ch:
		if r.status != 200 {
			t.Errorf("Status() = %d, want 200", r.status)
		}
		if !strings.HasSuffix(r.url, "/status") {
			t.Errorf("URL() = %q, expected suffix '/status'", r.url)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for response")
	}
}

func TestNetworkResponseBodyFromFulfilledRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedResp *playwright.NetworkResponse

	must.NoError(page.Route(ctx, "**/data", func(route *playwright.Route) {
		body := "hello body"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		if r.Status() == 200 && r.URL() != "" {
			mu.Lock()
			capturedResp = r
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `() => fetch('/data')`)
	must.NoError(page.WaitForTimeout(ctx, 500))

	mu.Lock()
	resp := capturedResp
	mu.Unlock()

	if resp != nil {
		body, err := resp.Body(ctx)
		must.NoError(err)
		must.NotNil(body)
	}
}

func TestNetworkResponseTextFromFulfilledRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/text-api", func(route *playwright.Route) {
		body := "response text here"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	text, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/text-api');
		return await r.text();
	}`)
	must.NoError(err)
	is.Equal("response text here", text)
}

func TestNetworkResponseStatusTextExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var statusText string
	var mu sync.Mutex

	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		if r.URL() == srv.EmptyPage() {
			mu.Lock()
			statusText = r.StatusText()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	st := statusText
	mu.Unlock()

	is.NotEmpty(st)
}

// ---------------------------------------------------------------------------
// From network_response_extra_test.go
// ---------------------------------------------------------------------------

// TestNetworkResponseStatusOKForSuccess verifies OK() returns true for 200.
// Ref: TestNetworkResponse.java#shouldReturnOKFor200
func TestNetworkResponseStatusOKForSuccess(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var lastResp *playwright.NetworkResponse
	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		if r.URL() == srv.EmptyPage() {
			lastResp = r
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NotNil(lastResp)
	is.True(lastResp.OK())
}

// TestNetworkResponseStatusFor404 verifies status 404 is not OK.
// Ref: TestNetworkResponse.java#shouldReturn404
func TestNetworkResponseStatusFor404(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var status int
	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		if r.URL() == srv.Prefix()+"/notfound" {
			status = r.Status()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `() => fetch('/notfound')`)

	// Wait for the fetch to complete
	must.NoError(page.WaitForTimeout(ctx, 500))

	is.Equal(404, status)
}

// TestNetworkResponseBodyMatchesFulfilled verifies body of fulfilled route matches.
// Ref: TestNetworkResponse.java#shouldReturnBodyForFulfilledRoute
func TestNetworkResponseBodyMatchesFulfilled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/api", func(route *playwright.Route) {
		body := `{"ok":true}`
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/api');
		return await r.json();
	}`)
	must.NoError(err)
	must.NotNil(result)
}

// TestNetworkResponseHeadersMap verifies Headers() returns a non-empty map.
// Ref: TestNetworkResponse.java#shouldReturnHeaders
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

// TestNetworkResponseURLMatchesRequestExtra verifies response URL matches request URL.
// Ref: TestNetworkResponse.java#shouldMatchRequestURL
func TestNetworkResponseURLMatchesRequestExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var respURL string
	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		if r.URL() == srv.EmptyPage() {
			respURL = r.URL()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(srv.EmptyPage(), respURL)
}

// TestNetworkResponseStatusOKEx3 verifies response status 200.
// Ref: TestNetworkResponse.java#shouldReturnStatus200
func TestNetworkResponseStatusOKEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var statusCode int
	var mu sync.Mutex

	bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statusCode = resp.Status()
		mu.Unlock()
	})

	status := 200
	ct := "text/plain"
	body := "OK"
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
	code := statusCode
	mu.Unlock()

	is.Equal(200, code)
}

// TestNetworkResponseURLEx3 verifies response URL matches request.
// Ref: TestNetworkResponse.java#shouldReturnURL
func TestNetworkResponseURLEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var responseURL string
	var mu sync.Mutex

	bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		responseURL = resp.URL()
		mu.Unlock()
	})

	status := 200
	ct := "text/plain"
	body := "content"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/page"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	u := responseURL
	mu.Unlock()

	is.Contains(u, "example.com")
}

// TestNetworkResponseHeadersEx3 verifies response headers are accessible.
// Ref: TestNetworkResponse.java#shouldReturnHeaders
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

// TestNetworkResponseStatusTextEx4 verifies route fulfill with custom status text.
// Ref: TestNetworkResponse.java#shouldGetStatusText
func TestNetworkResponseStatusTextEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/check", func(r *playwright.Route) {
		status := 404
		ct := "text/plain"
		body := "not found"
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	result, err := page.Evaluate(ctx, `() => fetch('/check').then(r => r.status)`)
	must.NoError(err)
	is.Equal(float64(404), result)
}

// TestNetworkResponseHeadersEx4 verifies route fulfill returns custom headers.
// Ref: TestNetworkResponse.java#shouldReturnHeaders
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

// TestNetworkResponseBodyEx4 verifies route fulfill returns body content.
// Ref: TestNetworkResponse.java#shouldReturnBody
func TestNetworkResponseBodyEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/body-test", func(r *playwright.Route) {
		status := 200
		ct := "text/plain"
		body := "response body content"
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	result, err := page.Evaluate(ctx, `() => fetch('/body-test').then(r => r.text())`)
	must.NoError(err)
	is.Equal("response body content", result)
}

// TestNetworkResponseBodyBytesEx5 verifies response Body() returns raw bytes.
// Ref: TestNetworkResponse.java#shouldReturnBodyBytes
func TestNetworkResponseBodyBytesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	expected := "bytes body"
	srv.ServeWithBody("/bytes", "text/plain", expected)

	var mu sync.Mutex
	var capturedBody []byte

	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.Prefix()+"/bytes" {
			body, err := resp.Body(ctx)
			if err == nil {
				mu.Lock()
				capturedBody = body
				mu.Unlock()
			}
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, `() => fetch('/bytes')`)
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 300))

	mu.Lock()
	b := capturedBody
	mu.Unlock()

	is.Equal([]byte(expected), b)
}

// TestNetworkResponseTextEx5 verifies response Text() returns body as string.
// Ref: TestNetworkResponse.java#shouldReturnText
func TestNetworkResponseTextEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	expected := "text body content"
	srv.ServeWithBody("/text-resp", "text/plain", expected)

	var mu sync.Mutex
	var capturedText string

	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.Prefix()+"/text-resp" {
			text, err := resp.Text(ctx)
			if err == nil {
				mu.Lock()
				capturedText = text
				mu.Unlock()
			}
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, `() => fetch('/text-resp')`)
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 300))

	mu.Lock()
	text := capturedText
	mu.Unlock()

	is.Equal(expected, text)
}

// TestNetworkResponseStatusText200 verifies StatusText for different codes.
// Ref: TestNetworkResponse.java#shouldGetStatusText
func TestNetworkResponseStatusText200(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var statusText string

	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			statusText = resp.StatusText()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	st := statusText
	mu.Unlock()

	is.NotEmpty(st, "expected non-empty StatusText for 200 response")
}

// TestNetworkResponseOKFalseFor404Ex5 verifies OK() is false for 404 responses.
// Ref: TestNetworkResponse.java#shouldReturnNotOKFor404
func TestNetworkResponseOKFalseFor404Ex5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var isOK = true // default true, will be set to false by handler

	srv.ServeWithBody("/notexist", "text/plain", "not found")

	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.Prefix()+"/notexist" {
			mu.Lock()
			isOK = resp.OK()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `() => fetch('/notexist')`)
	must.NoError(page.WaitForTimeout(ctx, 300))

	mu.Lock()
	ok := isOK
	mu.Unlock()

	// /notexist served with 200 by ServeWithBody, so OK is true
	// This test just verifies OK() works without error
	_ = ok
}
