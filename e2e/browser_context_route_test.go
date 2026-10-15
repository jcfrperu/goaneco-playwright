//go:build e2e

// E2E tests for: TestBrowserContextRoute.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextFulfillRouteEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	body := "ctx-routed"
	ct := "text/html"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Body:        &body,
			ContentType: &ct,
		}))
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, "http://example.com/"))
	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "ctx-routed")
}

// TestBrowserContextRouteModifyRequestHeaders verifies the context route can modify request headers.
// Ref: TestBrowserContextRoute.java (header modification variant)
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteModifyRequestHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var mu sync.Mutex
	var received string
	srv.SetRoute("/ctx-headers", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		received = r.Header.Get("X-Context-Route")
		mu.Unlock()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bCtx.Route(ctx, "**/ctx-headers", func(r *playwright.Route) {
		must.NoError(r.Continue(ctx, &playwright.RouteContinueOptions{
			Headers: map[string]string{"X-Context-Route": "injected"},
		}))
	}))

	_, err = page.Evaluate(ctx, `() => fetch('/ctx-headers')`)
	must.NoError(err)

	mu.Lock()
	h := received
	mu.Unlock()
	is.Equal("injected", h)
}

// TestBrowserContextRouteMultiplePages verifies context route applies to all pages in the context.
// Ref: TestBrowserContextRoute.java (multi-page variant)
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteMultiplePages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page1, err := bCtx.NewPage(ctx)
	must.NoError(err)
	page2, err := bCtx.NewPage(ctx)
	must.NoError(err)

	intercepted := "intercepted-by-ctx"
	must.NoError(bCtx.Route(ctx, "**/multi-page-route", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &intercepted})
	}))

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	for i, p := range []*playwright.Page{page1, page2} {
		res, err := p.Evaluate(ctx, `() => fetch('/multi-page-route').then(r => r.text())`)
		must.NoError(err)
		is.Equal("intercepted-by-ctx", res, fmt.Sprintf("page%d should be intercepted by context route", i+1))
	}
}

// TestBrowserContextRouteFulfillJSON verifies context route can return a JSON response.
// Ref: TestBrowserContextRoute.java (fulfill with JSON variant)
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteFulfillJSON(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	status := 200
	body := `{"ok":true}`
	contentType := "application/json"
	must.NoError(bCtx.Route(ctx, "**/api/json", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      &status,
			Body:        &body,
			ContentType: &contentType,
		})
	}))

	res, err := page.Evaluate(ctx, `() => fetch('/api/json').then(r => r.json())`)
	must.NoError(err)
	m, ok := res.(map[string]any)
	is.True(ok)
	is.Equal(true, m["ok"])
}

// TestBrowserContextRouteAbortRequest verifies context route can abort a request.
// Ref: TestBrowserContextRoute.java (abort variant)
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteAbortRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bCtx.Route(ctx, "**/api/ctx-abort", func(r *playwright.Route) {
		_ = r.Abort(ctx, "failed")
	}))

	res, err := page.Evaluate(ctx, `async () => {
		try {
			await fetch('/api/ctx-abort');
			return 'ok';
		} catch(e) { return 'aborted'; }
	}`)
	must.NoError(err)
	is.Equal("aborted", res)
}

// TestBrowserContextRouteShouldYieldToPageRoute verifies that a page-level route takes priority
// over a context-level route when both patterns match.
// Ref: TestBrowserContextRoute.java#shouldYieldToPageRoute  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldYieldToPageRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	ctxBody := "context-route"
	pageBody := "page-route"

	// Context route: returns "context-route"
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &ctxBody})
	}))
	// Page route: returns "page-route" — should take priority
	must.NoError(page.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &pageBody})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "page-route", "page-level route should take priority over context route")
}

// TestBrowserContextRouteShouldFallBackToContextRoute verifies that the context route handles requests
// that are not intercepted by any page route.
// Ref: TestBrowserContextRoute.java#shouldFallBackToContextRoute  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldFallBackToContextRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	ctxBody := "from-context-route"

	// Only a context route, no page route
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &ctxBody})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "from-context-route")
}

// TestBrowserContextRouteShouldOverwritePostBodyWithEmptyString verifies that Route.Continue can replace the
// POST body with an empty string, and the server receives an empty body.
// Ref: TestBrowserContextRoute.java#shouldOverwritePostBodyWithEmptyString  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldOverwritePostBodyWithEmptyString(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var mu sync.Mutex
	var receivedBody string
	srv.SetRoute("/echo-body", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		receivedBody = string(b)
		mu.Unlock()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain")
		w.Write(b)
	})

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	// Context route overwrites POST body with empty
	must.NoError(bCtx.Route(ctx, "**/echo-body", func(r *playwright.Route) {
		must.NoError(r.Continue(ctx, &playwright.RouteContinueOptions{
			PostData: []byte(""),
		}))
	}))

	_, err = page.Evaluate(ctx, `() => fetch('/echo-body', {method: 'POST', body: 'original'}).then(r => r.text())`)
	must.NoError(err)

	mu.Lock()
	body := receivedBody
	mu.Unlock()
	is.Equal("", body, "server should receive empty body after overwrite")
}

// TestBrowserContextRouteShouldNotSwallowExceptionsInRoute verifies that a panic in a context route handler is
// recovered without crashing the process.
// Ref: TestBrowserContextRoute.java#shouldNotSwallowExceptionsInRoute  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldNotSwallowExceptionsInRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	handlerCalled := make(chan struct{}, 1)
	must.NoError(bCtx.Route(ctx, "**/api/panic-ctx", func(r *playwright.Route) {
		select {
		case handlerCalled <- struct{}{}:
		default:
		}
		panic("deliberate context route panic")
	}))

	// The JS AbortController ensures the fetch times out even though the handler panicked.
	_, err = page.Evaluate(ctx, `async () => {
		const ctrl = new AbortController();
		setTimeout(() => ctrl.abort(), 300);
		try {
			await fetch('/api/panic-ctx', { signal: ctrl.signal });
			return 'ok';
		} catch(e) { return 'aborted'; }
	}`)
	must.NoError(err)

	select {
	case <-handlerCalled:
	case <-time.After(5 * time.Second):
		t.Fatal("route handler never invoked")
	}
	// If we reach here, the process survived the panic (recovery works).
}

// TestBrowserContextRouteShouldIntercept verifies request metadata is accessible in the route handler.
// Ref: TestBrowserContextRoute.java#shouldIntercept  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldIntercept(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var intercepted bool
	var reqURL, method, resourceType string
	var isNavigation bool
	var postData string
	var userAgent string

	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		intercepted = true
		req := r.Request()
		reqURL = req.URL()
		method = req.Method()
		resourceType = req.ResourceType()
		isNavigation = req.IsNavigationRequest()
		postData = req.PostData()
		userAgent = req.Headers()["user-agent"]
		_ = r.Continue(ctx, nil)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.True(intercepted, "route handler was not called")
	is.Contains(reqURL, "empty.html")
	is.Equal("GET", method)
	is.Equal("document", resourceType)
	is.True(isNavigation)
	is.Empty(postData)
	is.NotEmpty(userAgent)
}

// TestBrowserContextRouteShouldUnroute verifies that Unroute removes all registered context routes.
// Ref: TestBrowserContextRoute.java#shouldUnroute  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldUnroute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	intercepted := 0
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		intercepted++
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: strPtr("intercepted")})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(1, intercepted, "route should have been called once")

	// Remove all routes.
	must.NoError(bCtx.Unroute(ctx))

	// Next navigation should reach the real server, not the route.
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal(1, intercepted, "route should NOT be called after Unroute")
}

// TestBrowserContextRouteShouldSupportTheTimesParameterWithRouteMatching verifies that a route with Times=N is
// invoked at most N times; subsequent requests pass through normally.
// Ref: TestBrowserContextRoute.java#shouldSupportTheTimesParameterWithRouteMatching  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldSupportTheTimesParameterWithRouteMatching(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	intercepted := 0

	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted++
		mu.Unlock()
		_ = r.Continue(ctx, nil)
	}, &playwright.BrowserContextRouteOptions{Times: localIntPtr(2)}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := intercepted
	mu.Unlock()
	is.Equal(2, got, "handler should be called exactly 2 times")
}

// TestBrowserContextRouteShouldChainFallback verifies that multiple context routes registered for the
// same pattern are called in LIFO order when each calls Fallback.
// Ref: TestBrowserContextRoute.java#shouldChainFallback  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldChainFallback(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var intercepted []int

	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 1)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 2)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 3)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := make([]int, len(intercepted))
	copy(got, intercepted)
	mu.Unlock()
	is.Equal([]int{3, 2, 1}, got)
}

// TestBrowserContextRouteShouldNotChainFulfill verifies that when a handler calls Fulfill, subsequent
// (lower-priority) handlers are NOT called.
// Ref: TestBrowserContextRoute.java#shouldNotChainFulfill  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldNotChainFulfill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	lowestCalled := false

	// lowest-priority handler (registered first) — must NOT be called
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		lowestCalled = true
	}))
	// middle: fulfills the route
	body := "fulfilled"
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))
	// top: falls back to middle
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		_ = r.Fallback(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	// The middle handler fulfilled with body "fulfilled"; verify the content.
	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "fulfilled")
	is.False(lowestCalled, "lowest-priority handler must not be called")
}

// TestBrowserContextRouteShouldNotChainAbort verifies that when a handler aborts the request,
// lower-priority handlers are NOT called.
// Ref: TestBrowserContextRoute.java#shouldNotChainAbort  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldNotChainAbort(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	lowestCalled := false

	// lowest-priority — must NOT be called
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		lowestCalled = true
	}))
	// middle: aborts
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		_ = r.Abort(ctx)
	}))
	// top: falls back to middle
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		_ = r.Fallback(ctx)
	}))

	err = page.Goto(ctx, srv.EmptyPage())
	is.Error(err, "navigation should fail because route was aborted")
	is.False(lowestCalled, "lowest-priority handler must not be called")
}

// TestBrowserContextRouteShouldChainFallbackIntoPage verifies that context routes and page routes form
// a unified LIFO chain: page routes are tried first (most recent first), then context routes.
// Ref: TestBrowserContextRoute.java#shouldChainFallbackIntoPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldChainFallbackIntoPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var intercepted []int

	// Context routes (called after page routes fall back), registered in order 1→2→3.
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 1)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 2)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 3)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))

	// Page routes (called first), registered in order 4→5→6.
	must.NoError(page.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 4)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))
	must.NoError(page.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 5)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))
	must.NoError(page.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 6)
		mu.Unlock()
		_ = r.Fallback(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := make([]int, len(intercepted))
	copy(got, intercepted)
	mu.Unlock()
	is.Equal([]int{6, 5, 4, 3, 2, 1}, got)
}

// TestBrowserContextRouteShouldFallBackAsync verifies that the fallback chain works correctly even when
// handlers introduce delays before calling Fallback.
// Ref: TestBrowserContextRoute.java#shouldFallBackAsync  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
func TestBrowserContextRouteShouldFallBackAsync(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var intercepted []int

	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 1)
		mu.Unlock()
		time.Sleep(50 * time.Millisecond)
		_ = r.Fallback(ctx)
	}))
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 2)
		mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		_ = r.Fallback(ctx)
	}))
	must.NoError(bCtx.Route(ctx, "**/empty.html", func(r *playwright.Route) {
		mu.Lock()
		intercepted = append(intercepted, 3)
		mu.Unlock()
		time.Sleep(150 * time.Millisecond)
		_ = r.Fallback(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := make([]int, len(intercepted))
	copy(got, intercepted)
	mu.Unlock()
	is.Equal([]int{3, 2, 1}, got)
}

// TestRouteNotChainedWhenFulfilled verifies that once a route handler calls Fulfill,
// subsequent handlers registered for the same pattern are NOT invoked.
// Ref: TestBrowserContextRoute.java#shouldNotChainFulfill
func TestRouteNotChainedWhenFulfilled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	failedHandlerCalled := false

	// The "last" handler (registered first, called last due to LIFO) would set failed=true
	must.NoError(bc.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		mu.Lock()
		failedHandlerCalled = true
		mu.Unlock()
		_ = route.Continue(ctx, nil)
	}))

	// This handler fulfills the response; the previous handler must NOT be called
	fulfilledBody := "fulfilled"
	must.NoError(bc.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &fulfilledBody})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "fulfilled", "fulfill handler should have responded")

	mu.Lock()
	wasCalled := failedHandlerCalled
	mu.Unlock()
	is.False(wasCalled, "handler registered before the fulfilling one must not be called")
}

// TestRouteNotChainedWhenAborted verifies that once a route handler calls Abort,
// subsequent handlers registered for the same pattern are NOT invoked.
// Ref: TestBrowserContextRoute.java#shouldNotChainAbort
func TestRouteNotChainedWhenAborted(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)
	var mu sync.Mutex
	failedHandlerCalled := false

	// The "last" handler would set failedHandlerCalled = true
	must.NoError(bc.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		mu.Lock()
		failedHandlerCalled = true
		mu.Unlock()
		_ = route.Continue(ctx, nil)
	}))

	// This handler aborts; the previous handler must NOT be called
	must.NoError(bc.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	// Navigation is expected to fail because the route is aborted
	_ = page.Goto(ctx, srv.EmptyPage())

	mu.Lock()
	wasCalled := failedHandlerCalled
	mu.Unlock()
	is.False(wasCalled, "handler registered before the aborting one must not be called")
}

// TestBrowserContextRouteUnrouteRemovesSpecificPattern verifies that context.Unroute removes
// the context's route handlers and allows requests to pass through.
// Ref: TestBrowserContextRoute.java#shouldUnroute (simplified)
func TestBrowserContextRouteUnrouteRemovesSpecificPattern(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)
	srv.ServeWithBody("/ctx-unroute-check", "text/plain", "real-response")

	// Register route that intercepts
	interceptBody := "intercepted"
	must.NoError(bc.Route(ctx, "**/ctx-unroute-check", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &interceptBody})
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	r1, err := page.Evaluate(ctx, `() => fetch('/ctx-unroute-check').then(r => r.text())`)
	must.NoError(err)
	is.Equal("intercepted", r1, "before unroute the context handler should intercept")

	// Remove all context routes
	must.NoError(bc.Unroute(ctx))

	r2, err := page.Evaluate(ctx, `() => fetch('/ctx-unroute-check').then(r => r.text())`)
	must.NoError(err)
	is.Equal("real-response", r2, "after unroute requests should reach real server")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
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

// TestBrowserContextRouteInterceptsFetch verifies BrowserContext.Route intercepts fetch requests.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
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

// TestBrowserContextRouteFulfillJSONExtra verifies Route can fulfill with JSON body.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
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

// TestBrowserContextRouteFulfillWithBodyEx3 verifies Route.Fulfill replaces response body.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
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

// TestBrowserContextRouteInterceptsAllPagesEx3 verifies route intercepts all pages in context.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
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

// TestRouteFulfillEmptyBodyEx6 verifies Route Fulfill with empty body.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextRoute.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextRoute.java
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
