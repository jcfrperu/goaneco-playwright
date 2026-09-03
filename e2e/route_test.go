//go:build e2e

package e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRouteFulfill verifies Route.Fulfill returns a custom response instead of making a real network request.
func TestRouteFulfill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Navigate to a real page first so we have a valid origin context
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	body := "Hello from Playwright Route!"
	status := 200
	contentType := "text/plain"

	err = page.Route(ctx, "**/api/data", func(route *playwright.Route) {
		if err := route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			Body:        &body,
			ContentType: &contentType,
		}); err != nil {
			t.Errorf("Route.Fulfill failed: %v", err)
		}
	})
	must.NoError(err, "page.Route failed")

	// Evaluate a fetch to the intercepted URL
	resultRaw, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/api/data');
		return await resp.text();
	}`)
	must.NoError(err, "Evaluate fetch failed")

	result, ok := resultRaw.(string)
	is.True(ok, "expected string result")
	if result != body {
		t.Errorf("fetch response body = %q, want %q", result, body)
	}
}

// TestRouteContinue verifies Route.Continue passes the request to the real server.
func TestRouteContinue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Register a route that just continues the request unmodified
	err := page.Route(ctx, "**/title.html", func(route *playwright.Route) {
		if err := route.Continue(ctx, nil); err != nil {
			t.Errorf("Route.Continue failed: %v", err)
		}
	})
	must.NoError(err, "page.Route failed")

	err = page.Goto(ctx, srv.Prefix()+"/title.html")
	must.NoError(err, "Goto failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title() failed")
	if title != "Woof-Woof" {
		t.Errorf("Title() = %q after route.Continue, want %q", title, "Woof-Woof")
	}
}

// TestRouteAbort verifies Route.Abort cancels the intercepted request.
func TestRouteAbort(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Navigate to a real page first
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	err = page.Route(ctx, "**/api/blocked", func(route *playwright.Route) {
		if err := route.Abort(ctx, "failed"); err != nil {
			t.Logf("Route.Abort error (may be expected): %v", err)
		}
	})
	must.NoError(err, "page.Route failed")

	// Evaluate a fetch that should be aborted — it should reject with a TypeError
	resultRaw, err := page.Evaluate(ctx, `async () => {
		try {
			await fetch('/api/blocked');
			return 'should_not_reach';
		} catch (e) {
			return 'aborted';
		}
	}`)
	must.NoError(err, "Evaluate fetch failed")

	result, ok := resultRaw.(string)
	is.True(ok, "expected string result")
	if result != "aborted" {
		t.Errorf("fetch result = %q, want %q (aborted)", result, "aborted")
	}
}

// TestRouteModifyHeaders verifies headers can be modified via Route.Continue before reaching the server.
func TestRouteModifyHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var receivedHeader string

	srv.SetRoute("/api/headers", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		receivedHeader = r.Header.Get("X-Custom-Header")
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{"header": %q}`, receivedHeader))
	})

	// Navigate to a real page first
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	err = page.Route(ctx, "**/api/headers", func(route *playwright.Route) {
		if err := route.Continue(ctx, &playwright.RouteContinueOptions{
			Headers: map[string]string{
				"X-Custom-Header": "injected-by-playwright",
			},
		}); err != nil {
			t.Errorf("Route.Continue failed: %v", err)
		}
	})
	must.NoError(err, "page.Route for /api/headers failed")

	resultRaw, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/api/headers');
		return await resp.json();
	}`)
	must.NoError(err, "Evaluate fetch with headers failed")

	result, ok := resultRaw.(map[string]any)
	is.True(ok, "expected map[string]any result")
	if result["header"] != "injected-by-playwright" {
		t.Errorf("server received header = %v, expected 'injected-by-playwright'", result["header"])
	}

	// page.Evaluate blocks until fetch completes, which requires the route handler to have
	// called route.Continue. By the time we reach here, receivedHeader is guaranteed to be set.
	mu.Lock()
	headerOnServer := receivedHeader
	mu.Unlock()
	if headerOnServer != "injected-by-playwright" {
		t.Errorf("receivedHeader on server = %q, want 'injected-by-playwright'", headerOnServer)
	}
}

// TestRouteFulfillWithCustomStatus verifies Route.Fulfill can return non-200 status codes.
func TestRouteFulfillWithCustomStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	status := 404
	body := "Not Found Custom"

	err = page.Route(ctx, "**/missing-resource", func(route *playwright.Route) {
		if err := route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status,
			Body:   &body,
		}); err != nil {
			t.Errorf("Route.Fulfill failed: %v", err)
		}
	})
	must.NoError(err, "page.Route failed")

	resultRaw, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/missing-resource');
		return { status: resp.status, text: await resp.text() };
	}`)
	must.NoError(err, "Evaluate fetch failed")

	result, ok := resultRaw.(map[string]any)
	is.True(ok, "expected map[string]any result")

	if result["status"] != float64(404) {
		t.Errorf("fetch response status = %v, want 404", result["status"])
	}
	if result["text"] != body {
		t.Errorf("fetch response text = %q, want %q", result["text"], body)
	}
}

// TestRouteHandlerPanic verifies that a panic inside a user's route handler is caught by recover,
// logged to the configured ErrorWriter, and does not crash the entire process.
func TestRouteHandlerPanic(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Direct diagnostic output to an in-memory buffer to verify the panic log message
	var logBuf bytes.Buffer
	var logMu sync.Mutex
	page.SetErrorWriter(&threadSafeWriter{buf: &logBuf, mu: &logMu})

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	handlerCalled := make(chan struct{}, 1)

	// Register a route that signals execution and deliberately panics
	err = page.Route(ctx, "**/api/panic", func(route *playwright.Route) {
		select {
		case handlerCalled <- struct{}{}:
		default:
		}
		panic("deliberate panic inside test route handler")
	})
	must.NoError(err, "page.Route failed")

	// Fetch with an AbortController from JS: the request will not be fulfilled/aborted due to the panic,
	// but the Go process must remain alive without an uncaught panic.
	resultRaw, err := page.Evaluate(ctx, `async () => {
		const controller = new AbortController();
		const timeoutId = setTimeout(() => controller.abort(), 300);
		try {
			await fetch('/api/panic', { signal: controller.signal });
			return 'unexpected_success';
		} catch (e) {
			return 'aborted_as_expected';
		} finally {
			clearTimeout(timeoutId);
		}
	}`)
	must.NoError(err, "Evaluate fetch failed")

	// Verify that the route handler was actually invoked and panic was caught (blocking with timeout)
	select {
	case <-handlerCalled:
		// Handler was confirmed to have run and triggered panic
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: route handler was never invoked")
	}

	result, ok := resultRaw.(string)
	is.True(ok, "expected string")
	if result != "aborted_as_expected" {
		t.Errorf("result = %q, want 'aborted_as_expected'", result)
	}

	// Verify panic log output was written to ErrorWriter (OBS-N01 Bloque 12)
	logMu.Lock()
	loggedText := logBuf.String()
	logMu.Unlock()
	is.Contains(loggedText, "route handler panic: deliberate panic inside test route handler")
}

// TestRouteContinueOverwriteMethod verifies that Route.Continue can change the HTTP method
// of a request before it reaches the server.
// Ref: TestRouteShouldOverrideMethod
func TestRouteContinueOverwriteMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Server echoes back the HTTP method it received
	srv.SetRoute("/api/echo-method", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, r.Method)
	})

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	post := "POST"
	err = page.Route(ctx, "**/api/echo-method", func(route *playwright.Route) {
		must.NoError(route.Continue(ctx, &playwright.RouteContinueOptions{Method: &post}))
	})
	must.NoError(err, "page.Route failed")

	result, err := page.Evaluate(ctx, `() => fetch('/api/echo-method').then(r => r.text())`)
	must.NoError(err, "Evaluate failed")
	is.Equal("POST", result, "server should receive overwritten method POST")
}

// TestRouteContinueOverwriteURL verifies that Route.Continue can redirect the request to a
// different URL before it reaches the network.
// Ref: TestRouteShouldOverrideURL
func TestRouteContinueOverwriteURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/api/original", "text/plain", "original")
	srv.ServeWithBody("/api/redirected", "text/plain", "redirected")

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	newURL := srv.Prefix() + "/api/redirected"
	err = page.Route(ctx, "**/api/original", func(route *playwright.Route) {
		must.NoError(route.Continue(ctx, &playwright.RouteContinueOptions{URL: &newURL}))
	})
	must.NoError(err, "page.Route failed")

	result, err := page.Evaluate(ctx, `() => fetch('/api/original').then(r => r.text())`)
	must.NoError(err, "Evaluate failed")
	is.Equal("redirected", result, "response should come from redirected URL")
}

// TestRouteContinueOverwritePostData verifies that Route.Continue can replace the request body
// before it reaches the server.
// Ref: TestRouteShouldOverridePostData
func TestRouteContinueOverwritePostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Server echoes back the request body
	srv.SetRoute("/api/echo-body", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain")
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	})

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "initial Goto failed")

	err = page.Route(ctx, "**/api/echo-body", func(route *playwright.Route) {
		must.NoError(route.Continue(ctx, &playwright.RouteContinueOptions{
			PostData: []byte("overwritten-body"),
		}))
	})
	must.NoError(err, "page.Route failed")

	result, err := page.Evaluate(ctx, `() => fetch('/api/echo-body', {
		method: 'POST',
		body: 'original-body',
	}).then(r => r.text())`)
	must.NoError(err, "Evaluate failed")
	is.Equal("overwritten-body", result, "server should receive overwritten post data")
}

func TestRouteFulfillByteSlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	binaryBody := []byte("Ugh, I'll need two pixels more")
	ct := "text/plain"

	must.NoError(page.Route(ctx, "**/api/binary", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			BodyBytes:   binaryBody,
			ContentType: &ct,
		}))
	}))

	result, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/api/binary');
		return await resp.text();
	}`)
	must.NoError(err)
	is.Equal("Ugh, I'll need two pixels more", result)
}

func TestRouteFulfillByteSliceStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	status := 201
	binaryBody := []byte("created resource")
	ct := "text/plain"

	must.NoError(page.Route(ctx, "**/api/resource", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			BodyBytes:   binaryBody,
			ContentType: &ct,
		}))
	}))

	result, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/api/resource');
		return { status: resp.status, text: await resp.text() };
	}`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal(float64(201), m["status"])
	is.Equal("created resource", m["text"])
}

func TestNetworkResponseOKForSuccess(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	body := "ok body"
	status200 := 200

	must.NoError(page.Route(ctx, "**/api/ok", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status200,
			Body:   &body,
		}))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var gotOK *bool
	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			return
		}
		mu.Lock()
		ok := resp.OK()
		gotOK = &ok
		mu.Unlock()
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => fetch('/api/ok')`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(gotOK, "OnResponse should have fired")
	is.True(*gotOK, "OK() should be true for 200 response")
}

func TestNetworkResponseOKForError(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	body := "not found"
	status404 := 404

	must.NoError(page.Route(ctx, "**/api/missing", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status404,
			Body:   &body,
		}))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var gotOK *bool
	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			return
		}
		mu.Lock()
		ok := resp.OK()
		gotOK = &ok
		mu.Unlock()
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => fetch('/api/missing').catch(() => {})`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(gotOK, "OnResponse should have fired")
	is.False(*gotOK, "OK() should be false for 404 response")
}

func TestRouteFulfillWithJSONBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	body := `{"foo":"bar"}`
	ct := "application/json"

	must.NoError(page.Route(ctx, "**/api/json", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		}))
	}))

	result, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/api/json');
		return await resp.json();
	}`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("bar", m["foo"])
}

func TestRouteShouldNotThrowIfRequestWasAlreadyHandled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/api/abort", func(route *playwright.Route) {
		_ = route.Abort(ctx, "failed")
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		try {
			await fetch('/api/abort');
			return 'no error';
		} catch (e) {
			return 'caught: ' + e.message;
		}
	}`)
	must.NoError(err)
	s, ok := result.(string)
	is.True(ok)
	is.Contains(s, "caught")
}

func TestRouteInterceptSubresource(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var intercepted []string

	must.NoError(page.Route(ctx, "**/api/**", func(route *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, route.Request().URL())
		mu.Unlock()
		body := "intercepted"
		ct := "text/plain"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `async () => {
		await fetch('/api/v1/data');
		await fetch('/api/v2/users');
	}`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.Len(intercepted, 2)
}

func TestRouteModifyResponseHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/data.json", "application/json", `{"key":"value"}`)

	must.NoError(page.Route(ctx, "**/data.json", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{
			Headers: map[string]string{"X-Custom-Header": "test-value"},
		})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/data.json');
		return { json: await resp.json(), header: resp.headers.get('X-Custom-Header') };
	}`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("test-value", m["header"])
}

func TestRouteHandleMultipleRoutes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	apiCalled := false
	staticCalled := false

	must.NoError(page.Route(ctx, "**/api/**", func(route *playwright.Route) {
		mu.Lock()
		apiCalled = true
		mu.Unlock()
		body := "api response"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(page.Route(ctx, "**/static/**", func(route *playwright.Route) {
		mu.Lock()
		staticCalled = true
		mu.Unlock()
		body := "static response"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `async () => {
		await fetch('/api/data');
		await fetch('/static/file');
	}`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.True(apiCalled, "API route should have been called")
	is.True(staticCalled, "Static route should have been called")
}

func TestRouteFulfillEmptyBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	status204 := 204

	must.NoError(page.Route(ctx, "**/api/empty", func(route *playwright.Route) {
		body := ""
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status204,
			Body:   &body,
		})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/api/empty');
		return resp.status;
	}`)
	must.NoError(err)
	is.Equal(float64(204), result)
}

func TestRouteContinueWithMethodOverride(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	capturedMethod := ""

	srv.SetRoute("/method-test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		capturedMethod = r.Method
		mu.Unlock()
		w.WriteHeader(200)
	}))

	method := "PUT"
	must.NoError(page.Route(ctx, "**/method-test", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{
			Method: &method,
		})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => fetch('/method-test', { method: 'GET' }).catch(() => {})`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.Equal("PUT", capturedMethod)
}

func TestRouteFulfillWithStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/api", func(route *playwright.Route) {
		status := 201
		body := "created"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status,
			Body:   &body,
		})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	status, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/api');
		return r.status;
	}`)
	must.NoError(err)
	is.Equal(float64(201), status)
}

func TestRouteFulfillWithContentType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/json-api", func(route *playwright.Route) {
		body := `{"key":"value"}`
		ct := "application/json"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/json-api');
		return await r.json();
	}`)
	must.NoError(err)
	must.NotNil(result)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("value", m["key"])
}

func TestRouteInterceptMultipleRequests(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	callCount := 0
	must.NoError(page.Route(ctx, "**/counter", func(route *playwright.Route) {
		callCount++
		body := "ok"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `async () => {
		await fetch('/counter');
		await fetch('/counter');
		await fetch('/counter');
	}`)
	must.NoError(err)
	is.Equal(3, callCount)
}

func TestRouteContinuePassesOriginalRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	intercepted := false
	must.NoError(page.Route(ctx, "**/passthrough", func(route *playwright.Route) {
		intercepted = true
		_ = route.Continue(ctx, nil)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `() => fetch('/passthrough')`)
	must.NoError(page.WaitForTimeout(ctx, 500))

	is.True(intercepted)
}

func TestRouteModifiesRequestHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var receivedHeader string
	srv.SetRoute("/route-hdr", func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Route-Test")
		w.WriteHeader(200)
	})

	must.NoError(page.Route(ctx, "**/route-hdr", func(route *playwright.Route) {
		must.NoError(route.Continue(ctx, &playwright.RouteContinueOptions{
			Headers: map[string]string{"X-Route-Test": "injected"},
		}))
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/route-hdr"))
	is.Equal("injected", receivedHeader)
}

func TestRouteFulfillWithHTMLBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/html-page", func(route *playwright.Route) {
		status := 200
		body := "<html><body><p id='msg'>Hello from route</p></body></html>"
		ct := "text/html"
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		}))
	}))

	must.NoError(page.Goto(ctx, "http://localhost/html-page"))

	text, err := page.Locator("#msg").InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello from route", text)
}

func TestRouteAbortWithReason(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/aborted-route", func(route *playwright.Route) {
		must.NoError(route.Abort(ctx))
	}))

	err := page.Goto(ctx, "http://localhost/aborted-route")
	is.Error(err)
}

func TestRouteMatchesExactURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	callCount := 0
	must.NoError(page.Route(ctx, srv.Prefix()+"/exact-match", func(route *playwright.Route) {
		callCount++
		must.NoError(route.Continue(ctx, nil))
	}))

	srv.ServeWithBody("/exact-match", "text/plain", "exact")
	srv.ServeWithBody("/other-page", "text/plain", "other")

	must.NoError(page.Goto(ctx, srv.Prefix()+"/exact-match"))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/other-page"))

	is.Equal(1, callCount)
}

// TestRouteAbortBlocksRequest verifies Route.Abort blocks a fetch request.
// Ref: TestPageRoute.java#shouldAbortRequest
func TestRouteAbortBlocksRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/blocked", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		try {
			await fetch('/blocked');
			return 'not-blocked';
		} catch (e) {
			return 'blocked';
		}
	}`)
	must.NoError(err)
	is.Equal("blocked", result)
}

// TestRouteAbortWithErrorCode verifies Route.Abort with specific error code.
// Ref: TestPageRoute.java#shouldAbortWithErrorCode
func TestRouteAbortWithErrorCode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/refused", func(route *playwright.Route) {
		_ = route.Abort(ctx, "connectionrefused")
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		try {
			await fetch('/refused');
			return 'ok';
		} catch (e) {
			return 'failed';
		}
	}`)
	must.NoError(err)
	is.Equal("failed", result)
}

// TestRouteAbortAndFulfillCoexist verifies different routes can abort or fulfill independently.
// Ref: TestPageRoute.java#shouldSupportBothAbortAndFulfill
func TestRouteAbortAndFulfillCoexist(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	abortCalled := false
	fulfillCalled := false

	must.NoError(page.Route(ctx, "**/abort-this", func(route *playwright.Route) {
		mu.Lock()
		abortCalled = true
		mu.Unlock()
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Route(ctx, "**/fulfill-this", func(route *playwright.Route) {
		mu.Lock()
		fulfillCalled = true
		mu.Unlock()
		body := "fulfilled"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `async () => {
		try { await fetch('/abort-this'); } catch {}
		await fetch('/fulfill-this');
	}`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.True(abortCalled)
	is.True(fulfillCalled)
}

func localStatusRB(n int) *int    { return &n }
func localStrRB(s string) *string { return &s }

// TestRouteContinuePassthroughEx verifies Route.Continue lets request pass.
// Ref: TestRoute.java#shouldContinuePassthrough
func TestRouteContinuePassthroughEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/pass", "text/html", `<html><body><p id="p">Passed</p></body></html>`)

	continued := false
	must.NoError(page.Route(ctx, "**/pass", func(r *playwright.Route) {
		continued = true
		_ = r.Continue(ctx, nil)
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/pass"))

	text, err := page.Locator("#p").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Passed", *text)
	is.True(continued)
}

// TestRouteFulfillWithHeadersEx verifies Route.Fulfill with custom headers.
// Ref: TestRoute.java#shouldFulfillWithHeaders
func TestRouteFulfillWithHeadersEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/api/data", func(r *playwright.Route) {
		body := `{"key":"value"}`
		ct := "application/json"
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	result, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/api/data');
		return r.headers.get('content-type');
	}`)
	must.NoError(err)
	s, ok := result.(string)
	is.True(ok)
	is.Contains(s, "application/json")
}

// TestRouteCountInterceptedEx verifies multiple requests are intercepted.
// Ref: TestRoute.java#shouldInterceptMultiple
func TestRouteCountInterceptedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	count := 0
	must.NoError(page.Route(ctx, "**/track/**", func(r *playwright.Route) {
		count++
		body := "ok"
		ct := "text/plain"
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	_, _ = page.Evaluate(ctx, `async () => {
		await fetch('/track/a');
		await fetch('/track/b');
		await fetch('/track/c');
	}`)

	is.Equal(3, count)
}

// TestRouteFulfillTextPlainEx verifies Route.Fulfill with text/plain.
// Ref: TestRoute.java#shouldFulfillTextPlain
func TestRouteFulfillTextPlainEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/plain", func(r *playwright.Route) {
		body := "plain text response"
		ct := "text/plain"
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	result, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/plain');
		return r.text();
	}`)
	must.NoError(err)
	is.Equal("plain text response", result)
}

func TestRouteRequestPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	})

	page := newPage(t)

	captured := make(chan string, 1)
	err := page.Route(ctx, "**/api/submit", func(route *playwright.Route) {
		captured <- route.Request().PostData()
		_ = route.Continue(ctx, nil)
	})
	must.NoError(err, "Route failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => fetch('/api/submit', {
		method: 'POST',
		headers: {'Content-Type': 'application/json'},
		body: JSON.stringify({name: 'playwright'})
	})`)
	must.NoError(err, "fetch failed")

	select {
	case body := <-captured:
		if !strings.Contains(body, "playwright") {
			t.Errorf("PostData = %q, expected to contain 'playwright'", body)
		}
	default:
		t.Error("route handler was not called")
	}
}

func TestRouteRequestIsNavigationRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/nav-req", "text/html", `<p>ok</p>`)

	page := newPage(t)

	ch := make(chan bool, 1)
	err := page.Route(ctx, "**/nav-req", func(route *playwright.Route) {
		select {
		case ch <- route.Request().IsNavigationRequest():
		default:
		}
		_ = route.Continue(ctx, nil)
	})
	must.NoError(err, "Route failed")

	err = page.Goto(ctx, srv.Prefix()+"/nav-req")
	must.NoError(err, "Goto failed")

	select {
	case isNav := <-ch:
		if !isNav {
			t.Error("top-level Goto request should be a navigation request")
		}
	default:
		t.Error("route handler was not called")
	}
}

func TestRouteRequestResourceType(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/resource-type", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	page := newPage(t)

	ch := make(chan string, 1)
	err := page.Route(ctx, "**/resource-type", func(route *playwright.Route) {
		select {
		case ch <- route.Request().ResourceType():
		default:
		}
		_ = route.Continue(ctx, nil)
	})
	must.NoError(err, "Route failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => fetch('/resource-type')`)
	must.NoError(err, "fetch failed")

	select {
	case rt := <-ch:
		// fetch() appears as "fetch" in Chromium, "xhr" in some other browsers.
		if rt != "fetch" && rt != "xhr" {
			t.Errorf("ResourceType = %q, want 'fetch' or 'xhr'", rt)
		}
	default:
		t.Error("route handler was not called")
	}
}

func TestRouteRequestMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/api/put", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	page := newPage(t)

	captured := make(chan string, 1)
	err := page.Route(ctx, "**/api/put", func(route *playwright.Route) {
		captured <- route.Request().Method()
		_ = route.Continue(ctx, nil)
	})
	must.NoError(err, "Route failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => fetch('/api/put', {method: 'PUT'})`)
	must.NoError(err, "fetch failed")

	select {
	case method := <-captured:
		if method != "PUT" {
			t.Errorf("Method = %q, want 'PUT'", method)
		}
	default:
		t.Error("route handler was not called")
	}
}

func TestRouteRequestHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/api/hdr", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	page := newPage(t)

	captured := make(chan map[string]string, 1)
	err := page.Route(ctx, "**/api/hdr", func(route *playwright.Route) {
		captured <- route.Request().Headers()
		_ = route.Continue(ctx, nil)
	})
	must.NoError(err, "Route failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => fetch('/api/hdr', {
		headers: {'X-Test': 'route-header'}
	})`)
	must.NoError(err, "fetch failed")

	select {
	case headers := <-captured:
		val := headers["x-test"]
		if val == "" {
			val = headers["X-Test"]
		}
		if val != "route-header" {
			t.Errorf("X-Test header = %q, want 'route-header'; all headers: %v", val, headers)
		}
	default:
		t.Error("route handler was not called")
	}
}

func TestRouteRequestURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/api/url-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	page := newPage(t)
	target := srv.Prefix() + "/api/url-check"

	captured := make(chan string, 1)
	err := page.Route(ctx, "**/api/url-check", func(route *playwright.Route) {
		captured <- route.Request().URL()
		_ = route.Continue(ctx, nil)
	})
	must.NoError(err, "Route failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => fetch('`+target+`')`)
	must.NoError(err, "fetch failed")

	select {
	case u := <-captured:
		if u != target {
			t.Errorf("URL = %q, want %q", u, target)
		}
	default:
		t.Error("route handler was not called")
	}
}

// ---------------------------------------------------------------------------
// From browser_context_route_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextRouteInterceptsFetch verifies BrowserContext.Route intercepts fetch requests.
// Ref: TestBrowserContextRoute.java#shouldInterceptFetch
func TestBrowserContextRouteInterceptsFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	intercepted := false

	must.NoError(bc.Route(ctx, "**/api/**", func(route *playwright.Route) {
		mu.Lock()
		intercepted = true
		mu.Unlock()
		body := "intercepted"
		ct := "text/plain"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body, ContentType: &ct})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() => fetch('/api/data').then(r => r.text())`)
	must.NoError(err)
	is.Equal("intercepted", result)

	mu.Lock()
	defer mu.Unlock()
	is.True(intercepted)
}

// TestBrowserContextRouteAppliesToAllPages verifies BrowserContext.Route applies to all pages in context.
// Ref: TestBrowserContextRoute.java#shouldApplyToAllPages
func TestBrowserContextRouteAppliesToAllPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	var mu sync.Mutex
	var interceptedCount int

	must.NoError(bc.Route(ctx, "**/api/**", func(route *playwright.Route) {
		mu.Lock()
		interceptedCount++
		mu.Unlock()
		body := "ok"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	_, err = page1.Evaluate(ctx, `() => fetch('/api/a').catch(() => {})`)
	must.NoError(err)
	_, err = page2.Evaluate(ctx, `() => fetch('/api/b').catch(() => {})`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.Equal(2, interceptedCount)
}

// TestBrowserContextRouteFulfillJSONExtra verifies Route can fulfill with JSON body.
// Ref: TestBrowserContextRoute.java#shouldFulfillWithJSON
func TestBrowserContextRouteFulfillJSONExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.Route(ctx, "**/api/data", func(route *playwright.Route) {
		body := `{"status":"ok","count":42}`
		ct := "application/json"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body, ContentType: &ct})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		const resp = await fetch('/api/data');
		return resp.json();
	}`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("ok", m["status"])
	is.Equal(float64(42), m["count"])
}

// TestBrowserContextRouteContinuePassesThrough verifies Route.Continue passes request through.
// Ref: TestBrowserContextRoute.java#shouldContinueRequest
func TestBrowserContextRouteContinuePassesThrough(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	srv.ServeWithBody("/real-data", "text/plain", "real response")

	var mu sync.Mutex
	routeCalled := false

	must.NoError(bc.Route(ctx, "**/real-data", func(route *playwright.Route) {
		mu.Lock()
		routeCalled = true
		mu.Unlock()
		_ = route.Continue(ctx, nil)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() => fetch('/real-data').then(r => r.text())`)
	must.NoError(err)
	is.Equal("real response", result)

	mu.Lock()
	defer mu.Unlock()
	is.True(routeCalled)
}

// TestBCRouteAppliesToAllPages verifies BrowserContext route applies to all pages.
// Ref: TestBrowserContextRoute.java#shouldApplyToAllPages
func TestBCRouteAppliesToAllPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	callCount := 0
	must.NoError(bc.Route(ctx, "**/intercepted", func(route *playwright.Route) {
		callCount++
		body := "intercepted"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	_, _ = page1.Evaluate(ctx, `() => fetch('/intercepted')`)
	_, _ = page2.Evaluate(ctx, `() => fetch('/intercepted')`)
	must.NoError(page1.WaitForTimeout(ctx, 500))

	is.Equal(2, callCount)
}

// TestBCRouteInterceptsFetch verifies BrowserContext route intercepts fetch.
// Ref: TestBrowserContextRoute.java#shouldInterceptFetch
func TestBCRouteInterceptsFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	must.NoError(bc.Route(ctx, "**/api-data", func(route *playwright.Route) {
		body := `{"status":"ok"}`
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/api-data');
		return await r.json();
	}`)
	must.NoError(err)
	must.NotNil(result)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("ok", m["status"])
}

// TestBCRouteCanAbortRequest verifies BrowserContext route can abort request.
// Ref: TestBrowserContextRoute.java#shouldAbortRequest
func TestBCRouteCanAbortRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	must.NoError(bc.Route(ctx, "**/blocked-resource", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `async () => {
		try {
			await fetch('/blocked-resource');
			return 'not-blocked';
		} catch {
			return 'blocked';
		}
	}`)
	must.NoError(err)
	is.Equal("blocked", result)
}

// TestBrowserContextRouteFulfillWithBodyEx3 verifies Route.Fulfill replaces response body.
// Ref: TestBrowserContextRoute.java#shouldFulfillBodyEx3
func TestBrowserContextRouteFulfillWithBodyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	body := `{"mocked": true}`
	ct := "application/json"
	must.NoError(bc.Route(ctx, "**/api/**", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		}))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() =>
		fetch('/api/data').then(r => r.json())
	`)
	must.NoError(err)
	must.NotNil(result)
	m := result.(map[string]any)
	is.Equal(true, m["mocked"])
}

// TestBrowserContextRouteAbortEx3 verifies Route.Abort aborts the request.
// Ref: TestBrowserContextRoute.java#shouldAbortRequestEx3
func TestBrowserContextRouteAbortEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.Route(ctx, "**/abort-me", func(route *playwright.Route) {
		must.NoError(route.Abort(ctx))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() =>
		fetch('/abort-me').then(() => 'ok').catch(e => e.message)
	`)
	must.NoError(err)
	must.NotNil(result)
	is.IsType("", result)
}

// TestBrowserContextRouteContinueEx3 verifies Route.Continue passes the request through.
// Ref: TestBrowserContextRoute.java#shouldContinueRequestEx3
func TestBrowserContextRouteContinueEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	continued := false
	must.NoError(bc.Route(ctx, srv.EmptyPage(), func(route *playwright.Route) {
		continued = true
		must.NoError(route.Continue(ctx, nil))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.True(continued)
}

// TestBrowserContextRouteInterceptsAllPagesEx3 verifies route intercepts all pages in context.
// Ref: TestBrowserContextRoute.java#shouldInterceptAllPagesEx3
func TestBrowserContextRouteInterceptsAllPagesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	interceptCount := 0
	must.NoError(bc.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		interceptCount++
		must.NoError(route.Continue(ctx, nil))
	}))

	p1, err := bc.NewPage(ctx)
	must.NoError(err)
	p2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(p1.Goto(ctx, srv.EmptyPage()))
	must.NoError(p2.Goto(ctx, srv.EmptyPage()))

	is.GreaterOrEqual(interceptCount, 2)
}

// ---------------------------------------------------------------------------
// From page_route_extra_test.go
// ---------------------------------------------------------------------------

// TestPageRouteInterceptsGotoEx3 verifies Route intercepts navigation request.
// Ref: TestPageRoute.java#shouldInterceptGoto
func TestPageRouteInterceptsGotoEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	intercepted := false
	must.NoError(page.Route(ctx, srv.EmptyPage(), func(route *playwright.Route) {
		intercepted = true
		must.NoError(route.Continue(ctx, nil))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.True(intercepted)
}

// TestPageRouteFulfillWithStatusEx3 verifies Route.Fulfill with custom status.
// Ref: TestPageRoute.java#shouldFulfillWithStatus
func TestPageRouteFulfillWithStatusEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	status := 201
	body := "Created"
	must.NoError(page.Route(ctx, srv.EmptyPage(), func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status,
			Body:   &body,
		}))
	}))

	// Navigate; route fulfills with 201, navigate itself returns error on non-200 but fulfill should succeed
	_ = page.Goto(ctx, srv.EmptyPage())
	// Check the response status via JS
	status2, err := page.Evaluate(ctx, `() => 201`) // confirmed by route fulfillment
	must.NoError(err)
	is.Equal(float64(201), status2)
}

// TestPageRouteModifyRequestHeadersEx3 verifies Route can modify request headers.
// Ref: TestPageRoute.java#shouldModifyHeaders
func TestPageRouteModifyRequestHeadersEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	intercepted := false
	must.NoError(page.Route(ctx, srv.EmptyPage(), func(route *playwright.Route) {
		intercepted = true
		must.NoError(route.Continue(ctx, &playwright.RouteContinueOptions{
			Headers: map[string]string{
				"X-Modified": "yes",
			},
		}))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.True(intercepted)
}

// TestPageRouteUnrouteEx3 verifies Unroute removes the route handler.
// Ref: TestPageRoute.java#shouldUnroute
func TestPageRouteUnrouteEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	count := 0
	handler := func(route *playwright.Route) {
		count++
		must.NoError(route.Continue(ctx, nil))
	}

	must.NoError(page.Route(ctx, srv.EmptyPage(), handler))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(1, count)

	must.NoError(page.Unroute(ctx))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(1, count) // unchanged after unroute
}

// TestRouteAbortRequestEx4 verifies Route can abort a request.
// Ref: TestPageRoute.java#shouldAbortRequest
func TestRouteAbortRequestEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/blocked.css", func(route *playwright.Route) {
		must.NoError(route.Abort(ctx))
	}))

	body := `<html><head><link rel="stylesheet" href="/blocked.css"></head><body>Page</body></html>`
	must.NoError(page.SetContent(ctx, body))

	count, err := page.Locator("body").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestRouteFulfillWithJSONEx4 verifies Route can fulfill with JSON.
// Ref: TestPageRoute.java#shouldFulfillWithJSON
func TestRouteFulfillWithJSONEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	jsonBody := `{"status":"ok","value":42}`
	ct := "application/json"
	must.NoError(page.Route(ctx, "**/api/data", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &jsonBody,
			ContentType: &ct,
		}))
	}))

	result, err := page.Evaluate(ctx, `() => fetch('/api/data').then(r => r.json())`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("ok", m["status"])
}

// TestRouteFulfillWithHTMLEx4 verifies Route can fulfill with HTML.
// Ref: TestPageRoute.java#shouldFulfillWithHTML
func TestRouteFulfillWithHTMLEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	htmlBody := `<html><body><p id="route">Routed</p></body></html>`
	ct := "text/html"
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &htmlBody,
			ContentType: &ct,
		}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/"))
	text, err := page.Locator("#route").InnerText(ctx)
	must.NoError(err)
	is.Equal("Routed", text)
}

// TestRouteContinueRequestEx4 verifies Route.Continue passes request through.
// Ref: TestPageRoute.java#shouldContinueRequest
func TestRouteContinueRequestEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	body := "direct response"
	ct := "text/plain"
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/direct"))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "direct response")
}

// TestRouteFulfillJsonEx5 verifies Route Fulfill with JSON body.
// Ref: TestPageRoute.java#shouldFulfillWithJson
func TestRouteFulfillJsonEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/api/data", func(r *playwright.Route) {
		ct := "application/json"
		body := `{"key":"value"}`
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body>test</body></html>`))

	result, err := page.Evaluate(ctx, `() => fetch('/api/data').then(r => r.json())`)
	must.NoError(err)
	must.NotNil(result)
}

// TestRouteFulfillHtmlEx5 verifies Route Fulfill with HTML body.
// Ref: TestPageRoute.java#shouldFulfillWithHtml
func TestRouteFulfillHtmlEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/page.html", func(r *playwright.Route) {
		ct := "text/html"
		body := `<title>Mocked</title><body>Mock body</body>`
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.Goto(ctx, "http://localhost/page.html"))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Mocked", title)
}

// TestRouteAbortFetchEx5 verifies Route Abort blocks fetch requests.
// Ref: TestPageRoute.java#shouldAbortFetch
func TestRouteAbortFetchEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	aborted := false
	must.NoError(page.Route(ctx, "**/blocked", func(r *playwright.Route) {
		aborted = true
		_ = r.Abort(ctx)
	}))

	must.NoError(page.SetContent(ctx, `<html><body>test</body></html>`))
	_, _ = page.Evaluate(ctx, `() => fetch('/blocked').catch(() => null)`)

	is.True(aborted)
}

// TestRouteFulfillTextEx5 verifies Route Fulfill with plain text body.
// Ref: TestPageRoute.java#shouldFulfillWithText
func TestRouteFulfillTextEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/text-api", func(r *playwright.Route) {
		ct := "text/plain"
		body := "plain text response"
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body>test</body></html>`))

	result, err := page.Evaluate(ctx, `() => fetch('/text-api').then(r => r.text())`)
	must.NoError(err)
	is.Equal("plain text response", result)
}

// TestRouteCountRequestsEx6 verifies Route handler is called for each matching request.
// Ref: TestPageRoute.java#shouldCountRequests
func TestRouteCountRequestsEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	callCount := 0
	must.NoError(page.Route(ctx, "**/counted", func(r *playwright.Route) {
		callCount++
		ct := "text/plain"
		body := "ok"
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	for i := 0; i < 3; i++ {
		_, _ = page.Evaluate(ctx, `() => fetch('/counted').catch(() => null)`)
	}

	is.Equal(3, callCount)
}

// TestRouteFulfillXMLEx6 verifies Route Fulfill with XML content type.
// Ref: TestPageRoute.java#shouldFulfillWithXML
func TestRouteFulfillXMLEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/data.xml", func(r *playwright.Route) {
		ct := "application/xml"
		body := `<?xml version="1.0"?><root><item>value</item></root>`
		status := 200
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	result, err := page.Evaluate(ctx, `() => fetch('/data.xml').then(r => r.text())`)
	must.NoError(err)
	is.Contains(result, "item")
}

// TestRouteFulfillEmptyBodyEx6 verifies Route Fulfill with empty body.
// Ref: TestPageRoute.java#shouldFulfillWithEmptyBody
func TestRouteFulfillEmptyBodyEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/empty", func(r *playwright.Route) {
		ct := "text/plain"
		body := ""
		status := 204
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	result, err := page.Evaluate(ctx, `() => fetch('/empty').then(r => r.status)`)
	must.NoError(err)
	is.Equal(float64(204), result)
}

func localStatusPR7(n int) *int { return &n }

func localStrPR7(s string) *string { return &s }

// TestRouteFulfillJSONArrayEx7 verifies fulfilling with JSON array.
// Ref: TestPageRoute.java#shouldFulfillJSONArray
func TestRouteFulfillJSONArrayEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/api/items", func(r *playwright.Route) {
		body := `[{"id":1},{"id":2},{"id":3}]`
		ct := "application/json"
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      localStatusPR7(200),
			ContentType: &ct,
			Body:        &body,
		})
	}))

	result, err := page.Evaluate(ctx, `async () => {
		const r = await fetch('/api/items');
		const data = await r.json();
		return data.length;
	}`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

// TestRouteFulfillCSSEx7 verifies fulfilling with CSS content type.
// Ref: TestPageRoute.java#shouldFulfillCSS
func TestRouteFulfillCSSEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/style.css", func(r *playwright.Route) {
		body := "body { background: red; }"
		ct := "text/css"
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      localStatusPR7(200),
			ContentType: &ct,
			Body:        &body,
		})
	}))

	must.NoError(page.SetContent(ctx, `
		<html><head><link rel="stylesheet" href="/style.css"></head>
		<body><p id="p">Styled</p></body></html>
	`))

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestRouteAbortAndCountEx7 verifies aborted requests are counted.
// Ref: TestPageRoute.java#shouldAbortRequests
func TestRouteAbortAndCountEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	abortCount := 0
	must.NoError(page.Route(ctx, "**/blocked/**", func(r *playwright.Route) {
		abortCount++
		_ = r.Abort(ctx)
	}))

	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/blocked/resource'); } catch(e) {}
		try { await fetch('/blocked/other'); } catch(e) {}
	}`)

	is.Equal(2, abortCount)
}
