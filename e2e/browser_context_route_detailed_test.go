//go:build e2e

// Detailed BrowserContext route interception tests.
// Migration of: TestBrowserContextRoute.java
package e2e

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func localIntPtr(n int) *int { return &n }

// TestBrowserContextRouteYieldToPageRoute verifies that a page-level route takes priority
// over a context-level route when both patterns match.
// Ref: TestBrowserContextRoute.java#shouldYieldToPageRoute
func TestBrowserContextRouteYieldToPageRoute(t *testing.T) {
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

// TestBrowserContextRouteUsedAsFallback verifies that the context route handles requests
// that are not intercepted by any page route.
// Ref: TestBrowserContextRoute.java#shouldFallBackToContextRoute (simplified)
func TestBrowserContextRouteUsedAsFallback(t *testing.T) {
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

// TestBrowserContextRouteOverwritePostBodyEmpty verifies that Route.Continue can replace the
// POST body with an empty string, and the server receives an empty body.
// Ref: TestBrowserContextRoute.java#shouldOverwritePostBodyWithEmptyString
func TestBrowserContextRouteOverwritePostBodyEmpty(t *testing.T) {
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

// TestBrowserContextRoutePanicRecovery verifies that a panic in a context route handler is
// recovered without crashing the process.
// Ref: TestBrowserContextRoute.java#shouldNotSwallowExceptionsInRoute
func TestBrowserContextRoutePanicRecovery(t *testing.T) {
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

// TestBrowserContextRouteModifyRequestHeaders verifies the context route can modify request headers.
// Ref: TestBrowserContextRoute.java (header modification variant)
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

// TestBrowserContextRouteInterceptRequestProperties verifies request metadata is accessible in the route handler.
// Ref: TestBrowserContextRoute.java#shouldIntercept
func TestBrowserContextRouteInterceptRequestProperties(t *testing.T) {
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

// TestBrowserContextRouteUnrouteAll verifies that Unroute removes all registered context routes.
// Ref: TestBrowserContextRoute.java#shouldUnroute (simplified)
func TestBrowserContextRouteUnrouteAll(t *testing.T) {
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

// TestBrowserContextRouteSupportTheTimesParameter verifies that a route with Times=N is
// invoked at most N times; subsequent requests pass through normally.
// Ref: TestBrowserContextRoute.java#shouldSupportTheTimesParameterWithRouteMatching
func TestBrowserContextRouteSupportTheTimesParameter(t *testing.T) {
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

// TestBrowserContextRouteChainFallback verifies that multiple context routes registered for the
// same pattern are called in LIFO order when each calls Fallback.
// Ref: TestBrowserContextRoute.java#shouldChainFallback
func TestBrowserContextRouteChainFallback(t *testing.T) {
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

// TestBrowserContextRouteNotChainFulfill verifies that when a handler calls Fulfill, subsequent
// (lower-priority) handlers are NOT called.
// Ref: TestBrowserContextRoute.java#shouldNotChainFulfill
func TestBrowserContextRouteNotChainFulfill(t *testing.T) {
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

// TestBrowserContextRouteNotChainAbort verifies that when a handler aborts the request,
// lower-priority handlers are NOT called.
// Ref: TestBrowserContextRoute.java#shouldNotChainAbort
func TestBrowserContextRouteNotChainAbort(t *testing.T) {
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

// TestBrowserContextRouteChainFallbackIntoPage verifies that context routes and page routes form
// a unified LIFO chain: page routes are tried first (most recent first), then context routes.
// Ref: TestBrowserContextRoute.java#shouldChainFallbackIntoPage
func TestBrowserContextRouteChainFallbackIntoPage(t *testing.T) {
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

// TestBrowserContextRouteFallBackAsync verifies that the fallback chain works correctly even when
// handlers introduce delays before calling Fallback.
// Ref: TestBrowserContextRoute.java#shouldFallBackAsync
func TestBrowserContextRouteFallBackAsync(t *testing.T) {
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
