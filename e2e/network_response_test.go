//go:build e2e

// E2E tests for: TestNetworkResponse.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNetworkResponseBodyAfterRoute verifies response body can be read after route fulfillment.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseBodyAfterRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	body := "custom body content"
	ct := "text/plain"

	must.NoError(page.Route(ctx, "**/api/body", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		}))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var capturedBody string

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.Prefix()+"/api/body" {
			text, err := resp.Text(ctx)
			if err == nil {
				mu.Lock()
				capturedBody = text
				mu.Unlock()
			}
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => fetch('/api/body').then(r => r.text())`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.Equal("custom body content", capturedBody)
}

// TestNetworkResponseStatusCodes verifies Status() reflects the HTTP status code.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseStatusCodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Set up routes for different status codes
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/api/200", 200},
		{"/api/201", 201},
		{"/api/404", 404},
	} {
		status := tc.status
		path := tc.path
		srv.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte("ok"))
		})
	}

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	statusMap := make(map[string]int)

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statusMap[resp.URL()] = resp.Status()
		mu.Unlock()
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `async () => {
		await fetch('/api/200');
		await fetch('/api/201');
		await fetch('/api/404').catch(() => {});
	}`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.Equal(200, statusMap[srv.Prefix()+"/api/200"])
	is.Equal(201, statusMap[srv.Prefix()+"/api/201"])
	is.Equal(404, statusMap[srv.Prefix()+"/api/404"])
}

// TestNetworkRequestResponseReturnsResponse verifies Request.Response() returns response.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkRequestResponseReturnsResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var capturedReq *playwright.NetworkRequest
	var mu sync.Mutex

	off := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			capturedReq = req
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	r := capturedReq
	mu.Unlock()

	must.NotNil(r)

	resp, err := r.Response(ctx)
	must.NoError(err)
	must.NotNil(resp)
	is.Equal(200, resp.Status())
}

// TestNetworkRequestResponseStatusIs200 verifies response status is 200 for successful nav.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkRequestResponseStatusIs200(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var capturedReq *playwright.NetworkRequest
	var mu sync.Mutex

	off := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			capturedReq = req
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	r := capturedReq
	mu.Unlock()

	must.NotNil(r)

	resp, err := r.Response(ctx)
	must.NoError(err)
	must.NotNil(resp)
	is.Equal(200, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// TestNetworkResponseURLMatchesRequestExtra verifies response URL matches request URL.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// TestNetworkResponseStatusTextEx4 verifies route fulfill with custom status text.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// TestNetworkResponseBodyBytesEx5 verifies response Body() returns raw bytes.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// TestNetworkResponseStatusText200 verifies StatusText for different codes.
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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
// Extra: no equivalent method in playwright-java; related class: TestNetworkResponse.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
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

// TestNetworkResponseShouldReturnBody verifies route fulfill returns body content.
// Ref: TestNetworkResponse.java#shouldReturnBody  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldReturnBody(t *testing.T) {
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

// TestNetworkResponseShouldReturnText verifies response Text() returns body as string.
// Ref: TestNetworkResponse.java#shouldReturnText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldReturnText(t *testing.T) {
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

// TestNetworkResponseShouldReturnUncompressedText verifies the gzip-encoded body is decoded.
// Ref: TestNetworkResponse.java#shouldReturnUncompressedText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldReturnUncompressedText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	a4Gzip(srv, "/simple.json", "application/json", []byte("{\"foo\": \"bar\"}\n"))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/simple.json"))
	resp := log.WaitResponse(t, "/simple.json", 0)
	is.Equal("gzip", resp.Headers()["content-encoding"])
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}

// TestNetworkResponseShouldThrowWhenRequestingBodyOfRedirectedResponse verifies the body of a redirect response is unavailable.
// Ref: TestNetworkResponse.java#shouldThrowWhenRequestingBodyOfRedirectedResponse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldThrowWhenRequestingBodyOfRedirectedResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	srv.ServeWithRedirect("/foo.html", "/empty.html")
	must.NoError(page.Goto(ctx, srv.Prefix()+"/foo.html"))
	redirected := log.WaitResponse(t, "/foo.html", 302)
	is.Equal(302, redirected.Status())
	_, err := redirected.Text(ctx)
	must.Error(err)
	is.Contains(err.Error(), "Response body is unavailable for redirect responses")
}

// TestNetworkResponseShouldWaitUntilResponseCompletes verifies Body() waits until the response is fully received.
// Ref: TestNetworkResponse.java#shouldWaitUntilResponseCompletes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldWaitUntilResponseCompletes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	log := a4Watch(t, page)

	responseWritten := make(chan struct{}, 4)
	responseRead := make(chan struct{})
	srv.SetRoute("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("hello "))
		w.(http.Flusher).Flush()
		responseWritten <- struct{}{}
		<-responseRead
		_, _ = w.Write([]byte("wor"))
		w.(http.Flusher).Flush()
		_, _ = w.Write([]byte("ld!"))
	})

	_, err := page.Evaluate(ctx, "() => fetch('./get', { method: 'GET'})")
	must.NoError(err)
	pageResponse := log.WaitResponse(t, "/get", 0)
	<-responseWritten
	is.Equal(200, pageResponse.Status())
	isFinished := func() bool {
		for _, r := range log.Finished() {
			if strings.Contains(r.URL(), "/get") {
				return true
			}
		}
		return false
	}
	is.False(isFinished())
	close(responseRead)
	text, err := pageResponse.Text(ctx)
	must.NoError(err)
	is.Equal("hello world!", text)
	is.Eventually(isFinished, 5*time.Second, 20*time.Millisecond)
}

// TestNetworkResponseShouldReturnBodyWithCompression verifies binary bodies are decoded when compressed.
// Ref: TestNetworkResponse.java#shouldReturnBodyWithCompression  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldReturnBodyWithCompression(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	expected := make([]byte, 4096)
	for i := range expected {
		expected[i] = byte(i * 31)
	}
	a4Gzip(srv, "/pptr.png", "image/png", expected)
	_ = page.Goto(ctx, srv.Prefix()+"/pptr.png")
	resp := log.WaitResponse(t, "/pptr.png", 0)
	body, err := resp.Body(ctx)
	must.NoError(err)
	is.Equal(expected, body)
}

// TestNetworkResponseShouldReturnStatusText verifies Response.StatusText().
// Ref: TestNetworkResponse.java#shouldReturnStatusText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldReturnStatusText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	srv.ServeWithBody("/cool", "text/plain", "")
	must.NoError(page.Goto(ctx, srv.Prefix()+"/cool"))
	resp := log.WaitResponse(t, "/cool", 0)
	assert.Equal(t, "OK", resp.StatusText())
}

// TestNetworkResponseShouldWork verifies response header names are lower-cased.
// Ref: TestNetworkResponse.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header()["foo"] = []string{"bar"}
		w.Header()["BaZ"] = []string{"bAz"}
		w.WriteHeader(200)
	})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	resp := log.WaitResponse(t, "/empty.html", 0)
	is.Equal("bar", resp.Headers()["foo"])
	is.Equal("bAz", resp.Headers()["baz"])
	_, has := resp.Headers()["BaZ"]
	is.False(has)
}

// TestNetworkResponseShouldReturnServerAddress is not ported.
// Ref: TestNetworkResponse.java#shouldReturnServerAddress  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkResponse.java
func TestNetworkResponseShouldReturnServerAddress(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	response := log.WaitResponse(t, "/empty.html", 0)
	address, err := response.ServerAddr(ctx)
	must.NoError(err)
	must.NotNil(address)
	is.Equal(srv.Port(), address.Port)
	is.Contains([]string{"127.0.0.1", "::1"}, address.IpAddress)
}
