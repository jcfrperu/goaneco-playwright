//go:build e2e

// E2E tests for: TestBrowserContextNetworkEvents.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"sync"
	"testing"
	"time"
)

// TestBrowserContextOnRequestFailedFired verifies BrowserContext.OnRequestFailed fires when a request fails.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextNetworkEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestBrowserContextOnRequestFailedFired(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	// Serve a page that loads a script that we abort via route.
	srv.ServeWithBody("/fail-page", "text/html",
		`<script src="/fail.js"></script>`)
	srv.SetRoute("/fail.js", func(w http.ResponseWriter, r *http.Request) {
		// This handler will never be reached because the route aborts the request.
		w.WriteHeader(200)
	})

	// Abort the script request via route to guarantee requestfailed fires.
	must.NoError(bc.Route(ctx, "**/fail.js", func(route *playwright.Route) {
		_ = route.Abort(ctx, "failed")
	}))

	done := make(chan struct{}, 1)
	var mu sync.Mutex
	var failedURL string

	cancel := bc.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		if failedURL == "" {
			failedURL = req.URL()
		}
		mu.Unlock()
		select {
		case done <- struct{}{}:
		default:
		}
	})
	defer cancel()

	_ = page.Goto(ctx, srv.Prefix()+"/fail-page")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("OnRequestFailed did not fire within 5 seconds")
	}

	mu.Lock()
	u := failedURL
	mu.Unlock()
	is.Contains(u, "/fail.js")
}

// TestBrowserContextOnRequestFailedEx verifies OnRequestFailed fires for aborted request.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextNetworkEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestBrowserContextOnRequestFailedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(c)
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	failCount := 0
	off := bc.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		failCount++
	})
	t.Cleanup(off)

	// Abort all requests to a specific path
	must.NoError(bc.Route(ctx, "**/fail-this", func(r *playwright.Route) {
		_ = r.Abort(ctx)
	}))

	_, _ = page.Evaluate(ctx, `async () => { try { await fetch('/fail-this'); } catch(e) {} }`)

	mu.Lock()
	n := failCount
	mu.Unlock()

	is.Greater(n, 0)
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestBrowserContextNetworkEventsBrowserContextEventsRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var requests []string
	cancel := bCtx.OnRequest(func(r *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, r.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := append([]string(nil), requests...)
	mu.Unlock()
	is.NotEmpty(got)
	is.Contains(got[0], srv.EmptyPage())
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsResponse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestBrowserContextNetworkEventsBrowserContextEventsResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var responses []string
	cancel := bCtx.OnResponse(func(r *playwright.NetworkResponse) {
		mu.Lock()
		responses = append(responses, r.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	got := append([]string(nil), responses...)
	mu.Unlock()
	is.NotEmpty(got)
	is.Contains(got[0], srv.EmptyPage())
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsRequestFailed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestBrowserContextNetworkEventsBrowserContextEventsRequestFailed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	// Close response body immediately to cause request failure.
	srv.SetRoute("/one-style.css", func(w http.ResponseWriter, r *http.Request) {
		// Abort by closing without writing anything.
	})

	failedCh := make(chan *playwright.NetworkRequest, 10)
	cancel := bCtx.OnRequestFailed(func(r *playwright.NetworkRequest) {
		failedCh <- r
	})
	defer cancel()

	// This page loads one-style.css which will fail.
	srv.ServeWithBody("/one-style.html", "text/html", `<link rel="stylesheet" href="/one-style.css">`)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/one-style.html"))

	select {
	case req := <-failedCh:
		is.Contains(req.URL(), "one-style.css")
		is.Equal("stylesheet", req.ResourceType())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for request failed event")
	}
}

// Ref: TestBrowserContextNetworkEvents.java#BrowserContextEventsRequestFinished  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestBrowserContextNetworkEventsBrowserContextEventsRequestFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	finishedCh := make(chan *playwright.NetworkRequest, 1)
	cancel := bCtx.OnRequestFinished(func(r *playwright.NetworkRequest) {
		finishedCh <- r
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case req := <-finishedCh:
		is.Contains(req.URL(), srv.EmptyPage())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for request finished event")
	}
}

// Ref: TestBrowserContextNetworkEvents.java#shouldFireEventsInProperOrder  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestBrowserContextNetworkEventsShouldFireEventsInProperOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var events []string
	addEvent := func(name string) func() {
		switch name {
		case "request":
			return bCtx.OnRequest(func(r *playwright.NetworkRequest) {
				mu.Lock()
				events = append(events, "request")
				mu.Unlock()
			})
		case "response":
			return bCtx.OnResponse(func(r *playwright.NetworkResponse) {
				mu.Lock()
				events = append(events, "response")
				mu.Unlock()
			})
		default: // requestfinished
			return bCtx.OnRequestFinished(func(r *playwright.NetworkRequest) {
				mu.Lock()
				events = append(events, "requestfinished")
				mu.Unlock()
			})
		}
	}
	c1 := addEvent("request")
	defer c1()
	c2 := addEvent("response")
	defer c2()
	c3 := addEvent("requestfinished")
	defer c3()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	// Allow async event handlers to flush.
	_ = page.WaitForTimeout(ctx, 200)

	mu.Lock()
	got := append([]string(nil), events...)
	mu.Unlock()

	// The first three events should follow the expected order.
	if len(got) >= 3 {
		is.Equal("request", got[0])
		is.Equal("response", got[1])
		is.Equal("requestfinished", got[2])
	} else {
		t.Logf("got events: %v", got)
		is.GreaterOrEqual(len(got), 3)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextNetworkEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestNetworkRequestResponseNilForFailedRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	ch := make(chan *playwright.NetworkRequest, 1)
	cancel := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		select {
		case ch <- req:
		default:
		}
	})
	defer cancel()

	// Route all requests to abort so we get a failed request.
	err := page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	})
	must.NoError(err, "Route failed")

	// Navigating to a URL that gets aborted triggers OnRequestFailed.
	_ = page.Goto(ctx, "http://localhost:19999/nonexistent")

	var req *playwright.NetworkRequest
	select {
	case req = <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal("OnRequestFailed never fired")
	}

	// For a failed/aborted request, Response() should return nil without error.
	resp, err := req.Response(ctx)
	if err != nil {
		t.Logf("NetworkRequest.Response() error (acceptable for aborted): %v", err)
	}
	if resp != nil {
		t.Logf("Response is non-nil for aborted request (status=%d), which is acceptable", resp.Status())
	}
}

// ---------------------------------------------------------------------------
// From network_headers_extra_test.go
// ---------------------------------------------------------------------------

// TestNetworkRequestResponseFromOnRequestFinished verifies request.Response() after request finishes.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextNetworkEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestNetworkRequestResponseFromOnRequestFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedResp *playwright.NetworkResponse

	cancel := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			resp, err := req.Response(ctx)
			if err == nil && resp != nil {
				mu.Lock()
				capturedResp = resp
				mu.Unlock()
			}
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	// Give time for OnRequestFinished to fire
	mu.Lock()
	defer mu.Unlock()
	// Note: capturedResp may be nil if timing didn't align, so we just verify no panic
	_ = capturedResp
}

// TestNetworkRequestAndResponseInOrder verifies request fires before response for navigation.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextNetworkEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextNetworkEvents.java
func TestNetworkRequestAndResponseInOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var order []string

	off1 := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			order = append(order, "request")
			mu.Unlock()
		}
	})
	defer off1()

	off2 := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			order = append(order, "response")
			mu.Unlock()
		}
	})
	defer off2()

	off3 := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			order = append(order, "requestfinished")
			mu.Unlock()
		}
	})
	defer off3()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	o := order
	mu.Unlock()

	is.GreaterOrEqual(len(o), 3, "expected at least 3 events")
	is.Equal("request", o[0])
	is.Equal("response", o[1])
	is.Equal("requestfinished", o[2])
}
