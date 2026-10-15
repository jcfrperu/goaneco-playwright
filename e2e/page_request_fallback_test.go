//go:build e2e

// E2E tests for: TestPageRequestFallback.java.
package e2e

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageRequestFallbackShouldFallBack verifies Route.Fallback chains handlers in reverse registration order.
// Ref: TestPageRequestFallback.java#shouldFallBack  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldFallBack(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var mu sync.Mutex
	var intercepted []int
	for i := 1; i <= 3; i++ {
		n := i
		must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) {
			mu.Lock()
			intercepted = append(intercepted, n)
			mu.Unlock()
			_ = route.Fallback(ctx)
		}))
	}
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	mu.Lock()
	defer mu.Unlock()
	is.Equal([]int{3, 2, 1}, intercepted)
}

// TestPageRequestFallbackShouldFallBackAfterException verifies falling back after a failed fulfill.
// Ref: TestPageRequestFallback.java#shouldFallBackAfterException  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldFallBackAfterException(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		_ = route.Continue(ctx, nil)
	}))
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		if err := route.Fulfill(ctx, &playwright.RouteFulfillOptions{}); err != nil {
			_ = route.Fallback(ctx)
		}
	}))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
}

// TestPageRequestFallbackShouldChainOnce verifies Times(1) handlers with fallback.
// Ref: TestPageRequestFallback.java#shouldChainOnce  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldChainOnce(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var didFulfill atomic.Bool
	must.NoError(page.Route(ctx, "**/title.html", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Status: a4Ptr(200), Body: a4Ptr("fulfilled one")})
		didFulfill.Store(true)
	}, &playwright.PageRouteOptions{Times: a4Ptr(1)}))
	must.NoError(page.Route(ctx, "**/title.html", func(route *playwright.Route) {
		_ = route.Fallback(ctx)
	}, &playwright.PageRouteOptions{Times: a4Ptr(1)}))
	log := a4Watch(t, page)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/title.html"))
	is.True(didFulfill.Load())
	resp := log.WaitResponse(t, "/title.html", 0)
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("fulfilled one", text)
}

// TestPageRequestFallbackShouldAmendHTTPHeaders verifies Route.Fallback amending headers.
// Ref: TestPageRequestFallback.java#shouldAmendHTTPHeaders  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldAmendHTTPHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var values []string

	must.NoError(page.Route(ctx, "**/sleep.zzz", func(route *playwright.Route) {
		mu.Lock()
		values = append(values, route.Request().Headers()["foo"])
		v, err := route.Request().HeaderValue(ctx, "FOO")
		if err == nil && v != nil {
			values = append(values, *v)
		}
		mu.Unlock()
		_ = route.Resume(ctx)
	}))
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		headers := make(map[string]string)
		for k, v := range route.Request().Headers() {
			headers[k] = v
		}
		headers["FOO"] = "bar"
		_ = route.Fallback(ctx, (&playwright.RouteFallbackOptions{}).SetHeaders(headers))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz')")
	must.NoError(err)
	rec.Wait(t)
	mu.Lock()
	values = append(values, rec.Header.Values("foo")...)
	mu.Unlock()
	is.Equal([]string{"bar", "bar", "bar"}, values)
}

// TestPageRequestFallbackShouldDeleteHeaderWithUndefinedValue verifies Route.Fallback deleting a header.
// Ref: TestPageRequestFallback.java#shouldDeleteHeaderWithUndefinedValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldDeleteHeaderWithUndefinedValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/empty.html"))

	rec := a4RecordRoute(srv, "/something")
	srv.SetRoute("/something", func(w http.ResponseWriter, req *http.Request) {
		b, _ := io.ReadAll(req.Body)
		rec.Method = req.Method
		rec.Body = b
		rec.Header = req.Header.Clone()
		select {
		case <-rec.got:
		default:
			close(rec.got)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("done"))
	})

	var mu sync.Mutex
	var interceptedRequest *playwright.NetworkRequest
	must.NoError(page.Route(ctx, srv.Prefix()+"/something", func(route *playwright.Route) {
		mu.Lock()
		interceptedRequest = route.Request()
		mu.Unlock()
		_ = route.Resume(ctx)
	}))
	must.NoError(page.Route(ctx, srv.Prefix()+"/something", func(route *playwright.Route) {
		headers, err := route.Request().AllHeaders(ctx)
		if err != nil {
			_ = route.Fallback(ctx)
			return
		}
		delete(headers, "foo")
		_ = route.Fallback(ctx, (&playwright.RouteFallbackOptions{}).SetHeaders(headers))
	}))

	text, err := page.Evaluate(ctx, `async url => {
		const data = await fetch(url, {
			headers: {
				foo: 'a',
				bar: 'b',
			}
		});
		return data.text();
	}`, srv.Prefix()+"/something")
	must.NoError(err)
	is.Equal("done", text)

	mu.Lock()
	req := interceptedRequest
	mu.Unlock()
	must.NotNil(req)
	_, hasFoo := req.Headers()["foo"]
	is.False(hasFoo)
	rec.Wait(t)
	is.Empty(rec.Header.Values("foo"))
	is.Equal([]string{"b"}, rec.Header.Values("bar"))
}

// TestPageRequestFallbackShouldAmendMethod verifies Route.Fallback amending the HTTP method.
// Ref: TestPageRequestFallback.java#shouldAmendMethod  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldAmendMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	rec := a4RecordRoute(srv, "/sleep.zzz")
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var method string
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		mu.Lock()
		method = route.Request().Method()
		mu.Unlock()
		_ = route.Resume(ctx)
	}))
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Fallback(ctx, (&playwright.RouteFallbackOptions{}).SetMethod("POST"))
	}))

	request, err := page.WaitForRequest(ctx, "**/sleep.zzz", func() error {
		_, evalErr := page.Evaluate(ctx, "() => fetch('/sleep.zzz')")
		return evalErr
	})
	must.NoError(err)
	rec.Wait(t)

	mu.Lock()
	m := method
	mu.Unlock()
	is.Equal("POST", m)
	is.Equal("POST", request.Method())
	is.Equal("POST", rec.Method)
}

// TestPageRequestFallbackShouldOverrideRequestUrl verifies Route.Fallback with a URL override.
// Ref: TestPageRequestFallback.java#shouldOverrideRequestUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldOverrideRequestUrl(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var mu sync.Mutex
	var serverMethod string
	srv.SetRoute("/global-var.html", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		serverMethod = r.Method
		mu.Unlock()
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<script>window['globalVar'] = 123;</script>"))
	})
	var url atomic.Value
	must.NoError(page.Route(ctx, "**/global-var.html", func(route *playwright.Route) {
		url.Store(route.Request().URL())
		_ = route.Continue(ctx, nil)
	}))
	must.NoError(page.Route(ctx, "**/foo", func(route *playwright.Route) {
		_ = route.Fallback(ctx, &playwright.RouteFallbackOptions{URL: a4Ptr(srv.Prefix() + "/global-var.html")})
	}))
	log := a4Watch(t, page)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/foo"))
	resp := log.WaitResponse(t, "/global-var.html", 0)
	is.Equal(srv.Prefix()+"/global-var.html", url.Load())
	is.Equal(srv.Prefix()+"/global-var.html", resp.URL())
	v, err := page.Evaluate(ctx, "() => window['globalVar']")
	must.NoError(err)
	is.EqualValues(123, v)
	mu.Lock()
	defer mu.Unlock()
	is.Equal("GET", serverMethod)
}

// TestPageRequestFallbackShouldAmendPostData verifies Route.Fallback amending request post data.
// Ref: TestPageRequestFallback.java#shouldAmendPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldAmendPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var postData string
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		mu.Lock()
		postData = route.Request().PostData()
		mu.Unlock()
		_ = route.Resume(ctx)
	}))
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Fallback(ctx, (&playwright.RouteFallbackOptions{}).SetPostData("doggo"))
	}))

	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz', { method: 'POST', body: 'birdy' })")
	must.NoError(err)
	rec.Wait(t)

	mu.Lock()
	pd := postData
	mu.Unlock()
	is.Equal("doggo", pd)
	is.Equal("doggo", string(rec.Body))
}

// TestPageRequestFallbackShouldAmendBinaryPostData verifies Route.Fallback amending binary post data.
// Ref: TestPageRequestFallback.java#shouldAmendBinaryPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldAmendBinaryPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	arr := make([]byte, 256)
	for i := 0; i < 256; i++ {
		arr[i] = byte(i)
	}

	var mu sync.Mutex
	var postDataBuffer []byte
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		mu.Lock()
		postDataBuffer = route.Request().PostDataBuffer()
		mu.Unlock()
		_ = route.Resume(ctx)
	}))
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Fallback(ctx, (&playwright.RouteFallbackOptions{}).SetPostData(arr))
	}))

	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz', { method: 'POST', body: 'birdy' })")
	must.NoError(err)
	rec.Wait(t)

	mu.Lock()
	pdb := postDataBuffer
	mu.Unlock()
	is.Equal(arr, pdb)
	is.Equal(arr, rec.Body)
}

// TestPageRequestFallbackShouldWork verifies a route that only falls back lets the request through.
// Ref: TestPageRequestFallback.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) { _ = route.Fallback(ctx) }))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
}

// TestPageRequestFallbackShouldFallBackAsync verifies handlers that fall back after a delay still chain in order.
// Ref: TestPageRequestFallback.java#shouldFallBackAsync  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldFallBackAsync(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var mu sync.Mutex
	var intercepted []int
	for i := 1; i <= 3; i++ {
		n := i
		must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) {
			mu.Lock()
			intercepted = append(intercepted, n)
			mu.Unlock()
			time.Sleep(time.Duration(n) * 50 * time.Millisecond)
			_ = route.Fallback(ctx)
		}))
	}
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	mu.Lock()
	defer mu.Unlock()
	must.Equal([]int{3, 2, 1}, intercepted)
}

// TestPageRequestFallbackShouldNotChainFulfill verifies fulfilling stops the handler chain.
// Ref: TestPageRequestFallback.java#shouldNotChainFulfill  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldNotChainFulfill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	var failed atomic.Bool
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) { failed.Store(true) }))
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Status: a4Ptr(200), Body: a4Ptr("fulfilled")})
	}))
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) { _ = route.Fallback(ctx) }))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	text, err := log.WaitResponse(t, "/empty.html", 0).Text(ctx)
	must.NoError(err)
	must.Equal("fulfilled", text)
	must.False(failed.Load())
}

// TestPageRequestFallbackShouldNotChainAbort verifies aborting stops the handler chain.
// Ref: TestPageRequestFallback.java#shouldNotChainAbort  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFallback.java
func TestPageRequestFallbackShouldNotChainAbort(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var failed atomic.Bool
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) { failed.Store(true) }))
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) { _ = route.Abort(ctx) }))
	must.NoError(page.Route(ctx, "**/empty.html", func(route *playwright.Route) { _ = route.Fallback(ctx) }))
	must.Error(page.Goto(ctx, srv.EmptyPage()))
	must.False(failed.Load())
}
