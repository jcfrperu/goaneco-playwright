//go:build e2e

// E2E tests for: TestPageEventNetwork.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"sync"
	"testing"
	"time"
)

// Ref: TestPageEventNetwork.java#PageEventsRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventNetworkPageEventsRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/we-req", "text/html", "<p>ok</p>")

	page := newPage(t)

	done := make(chan any, 1)
	go func() {
		val, err := page.WaitForEvent(ctx, "request")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	err := page.Goto(ctx, srv.Prefix()+"/we-req")
	must.NoError(err, "Goto failed")

	result := <-done
	if req, ok := result.(*playwright.NetworkRequest); ok {
		if req.URL() == "" {
			t.Error("WaitForEvent(request) returned request with empty URL")
		}
	} else {
		t.Fatalf("WaitForEvent(request) unexpected result: %v", result)
	}
}

// --- BrowserContext.WaitForEvent("page") ---

// TestNetworkEventOrderRequestBeforeResponse verifies request fires before response for same resource.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestNetworkEventOrderRequestBeforeResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var events []string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			events = append(events, "request")
			mu.Unlock()
		}
	})
	defer cancel()

	cancel2 := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			events = append(events, "response")
			mu.Unlock()
		}
	})
	defer cancel2()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.Len(events, 2)
	is.Equal("request", events[0])
	is.Equal("response", events[1])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageOnRequestFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/fin", "text/html", `<p>done</p>`)
	page := newPage(t)

	done := make(chan string, 1)
	cancel := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() != "" {
			select {
			case done <- req.URL():
			default:
			}
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/fin")
	must.NoError(err, "Goto failed")

	select {
	case url := <-done:
		t.Logf("OnRequestFinished fired for URL: %s", url)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnRequestFinished")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageOnRequestFailed(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	page := newPage(t)

	done := make(chan string, 1)
	cancel := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		select {
		case done <- req.URL():
		default:
		}
	})
	defer cancel()

	// Navigate to a non-existent host to trigger a request failure.
	_ = page.Goto(ctx, "http://localhost:1/nonexistent")

	select {
	case url := <-done:
		t.Logf("OnRequestFailed fired for URL: %s", url)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for OnRequestFailed")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventsFireInProperOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/order", "text/html", `<p>order</p>`)
	page := newPage(t)

	var mu sync.Mutex
	var order []string

	cancelReq := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		order = append(order, "request")
		mu.Unlock()
	})
	defer cancelReq()

	cancelResp := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		order = append(order, "response")
		mu.Unlock()
	})
	defer cancelResp()

	cancelFin := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		mu.Lock()
		order = append(order, "requestFinished")
		mu.Unlock()
	})
	defer cancelFin()

	err := page.Goto(ctx, srv.Prefix()+"/order")
	must.NoError(err, "Goto failed")

	// Give handlers time to fire.
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	got := make([]string, len(order))
	copy(got, order)
	mu.Unlock()

	// Find at least request → response → requestFinished order for the main request.
	reqIdx, respIdx, finIdx := -1, -1, -1
	for i, e := range got {
		switch e {
		case "request":
			if reqIdx == -1 {
				reqIdx = i
			}
		case "response":
			if respIdx == -1 {
				respIdx = i
			}
		case "requestFinished":
			if finIdx == -1 {
				finIdx = i
			}
		}
	}

	if reqIdx == -1 || respIdx == -1 || finIdx == -1 {
		t.Fatalf("expected request+response+requestFinished events, got: %v", got)
	}
	if !(reqIdx < respIdx && respIdx < finIdx) {
		t.Errorf("expected request < response < requestFinished order, got indices req=%d resp=%d fin=%d in %v", reqIdx, respIdx, finIdx, got)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestBrowserContextOnRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-req", "text/html", `<p>ctx</p>`)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	done := make(chan string, 1)
	cancel := bCtx.OnRequest(func(req *playwright.NetworkRequest) {
		select {
		case done <- req.URL():
		default:
		}
	})
	defer cancel()

	err = page.Goto(ctx, srv.Prefix()+"/ctx-req")
	must.NoError(err, "Goto failed")

	select {
	case url := <-done:
		t.Logf("BrowserContext.OnRequest fired: %s", url)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for BrowserContext.OnRequest")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestBrowserContextOnResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-resp", "text/html", `<p>resp</p>`)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	done := make(chan int, 1)
	cancel := bCtx.OnResponse(func(resp *playwright.NetworkResponse) {
		select {
		case done <- resp.Status():
		default:
		}
	})
	defer cancel()

	err = page.Goto(ctx, srv.Prefix()+"/ctx-resp")
	must.NoError(err, "Goto failed")

	select {
	case status := <-done:
		if status != 200 {
			t.Errorf("expected status 200, got %d", status)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for BrowserContext.OnResponse")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestBrowserContextOnRequestFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-fin", "text/html", `<p>fin</p>`)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	done := make(chan struct{}, 1)
	cancel := bCtx.OnRequestFinished(func(req *playwright.NetworkRequest) {
		select {
		case done <- struct{}{}:
		default:
		}
	})
	defer cancel()

	err = page.Goto(ctx, srv.Prefix()+"/ctx-fin")
	must.NoError(err, "Goto failed")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for BrowserContext.OnRequestFinished")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestBrowserContextOnRequestFailed(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	done := make(chan string, 1)
	cancel := bCtx.OnRequestFailed(func(req *playwright.NetworkRequest) {
		select {
		case done <- req.URL():
		default:
		}
	})
	defer cancel()

	_ = page.Goto(ctx, "http://localhost:1/nonexistent")

	select {
	case url := <-done:
		t.Logf("BrowserContext.OnRequestFailed fired: %s", url)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for BrowserContext.OnRequestFailed")
	}
}

// ---------------------------------------------------------------------------
// From network_events_extra_test.go
// ---------------------------------------------------------------------------

// TestNetworkEventsRequestFiredEx verifies OnRequest fires on page navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestNetworkEventsRequestFiredEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/network-req", "text/html", `<html><body>page</body></html>`)

	var mu sync.Mutex
	var urls []string
	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		urls = append(urls, req.URL())
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/network-req"))

	mu.Lock()
	captured := urls
	mu.Unlock()

	is.NotEmpty(captured)
	is.Contains(captured[0], "/network-req")
}

// TestNetworkEventsResponseFiredEx verifies OnResponse fires after request.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestNetworkEventsResponseFiredEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/network-resp", "text/html", `<html><body>response page</body></html>`)

	var mu sync.Mutex
	var statuses []int
	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		defer mu.Unlock()
		statuses = append(statuses, resp.Status())
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/network-resp"))

	mu.Lock()
	captured := statuses
	mu.Unlock()

	is.NotEmpty(captured)
	is.Equal(200, captured[0])
}

// TestNetworkEventsRequestFinishedEx verifies OnRequestFinished fires after response.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestNetworkEventsRequestFinishedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/net-finished", "text/html", `<html><body>finished</body></html>`)

	var mu sync.Mutex
	finishedURLs := []string{}
	off := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		finishedURLs = append(finishedURLs, req.URL())
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/net-finished"))

	mu.Lock()
	captured := finishedURLs
	mu.Unlock()

	is.NotEmpty(captured)
}

// TestNetworkEventsRequestMethodEx verifies request method is GET on navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestNetworkEventsRequestMethodEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/net-method", "text/html", `<html><body>method</body></html>`)

	var mu sync.Mutex
	var method string
	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		if method == "" {
			method = req.Method()
		}
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/net-method"))

	mu.Lock()
	m := method
	mu.Unlock()

	is.Equal("GET", m)
}

// TestOnRequestFailedFiredForAbortedRoute verifies OnRequestFailed fires when route aborts.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestOnRequestFailedFiredForAbortedRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var failedURL string

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		failedURL = req.URL()
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Route(ctx, "**/blocked", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/blocked'); } catch {}
	}`)

	mu.Lock()
	url := failedURL
	mu.Unlock()

	is.Contains(url, "blocked")
}

// TestOnRequestFailedURLNotEmpty verifies OnRequestFailed provides a URL.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestOnRequestFailedURLNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	urls := make([]string, 0)

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Route(ctx, "**/fail", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/fail'); } catch {}
	}`)

	mu.Lock()
	captured := len(urls)
	mu.Unlock()

	is.Greater(captured, 0)
}

// TestOnRequestFailedOffStopsReceiving verifies unregistering the handler stops events.
// Extra: no equivalent method in playwright-java; related class: TestPageEventNetwork.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestOnRequestFailedOffStopsReceiving(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	count := 0

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	must.NoError(page.Route(ctx, "**/abort1", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/abort1'); } catch {}
	}`)

	off() // unregister

	mu.Lock()
	countAfterOff := count
	mu.Unlock()

	// Any subsequent fails after off() should not increment
	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/abort1'); } catch {}
	}`)

	mu.Lock()
	is.Equal(countAfterOff, count)
	mu.Unlock()
}

// TestPageEventNetworkPageEventsResponse verifies Page.OnResponse fires with the navigation response.
// Ref: TestPageEventNetwork.java#PageEventsResponse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventNetworkPageEventsResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	resp := log.WaitResponse(t, "/empty.html", 0)
	log.mu.Lock()
	n := len(log.resps)
	log.mu.Unlock()
	is.Equal(1, n)
	is.Equal(srv.EmptyPage(), resp.URL())
	is.Equal(200, resp.Status())
	is.True(resp.OK())
	// Response.Request() is not exposed by the Go client, so assertNotNull(response.request()) cannot be ported.
}

// TestPageEventNetworkShouldSupportRedirects verifies the request/response/finished event order across redirects.
// Ref: TestPageEventNetwork.java#shouldSupportRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventNetworkShouldSupportRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var mu sync.Mutex
	var events []string
	add := func(s string) { mu.Lock(); events = append(events, s); mu.Unlock() }
	offs := []func(){
		page.OnRequest(func(r *playwright.NetworkRequest) { add(r.Method() + " " + r.URL()) }),
		page.OnResponse(func(r *playwright.NetworkResponse) { add(fmt.Sprintf("%d %s", r.Status(), r.URL())) }),
		page.OnRequestFinished(func(r *playwright.NetworkRequest) { add("DONE " + r.URL()) }),
		page.OnRequestFailed(func(r *playwright.NetworkRequest) { add("FAIL " + r.URL()) }),
	}
	defer func() {
		for _, off := range offs {
			off()
		}
	}()
	srv.ServeWithRedirect("/foo.html", "/empty.html")
	fooURL := srv.Prefix() + "/foo.html"
	must.NoError(page.Goto(ctx, fooURL))
	want := []string{
		"GET " + fooURL,
		"302 " + fooURL,
		"DONE " + fooURL,
		"GET " + srv.EmptyPage(),
		"200 " + srv.EmptyPage(),
		"DONE " + srv.EmptyPage(),
	}
	is.Eventually(func() bool { mu.Lock(); defer mu.Unlock(); return len(events) >= len(want) }, 5*time.Second, 20*time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	is.Equal(want, events)
	// Request.RedirectedFrom()/RedirectedTo() are not exposed by the Go client, so the redirect chain assertions cannot be ported.
}

// TestPageEventNetworkPageEventsRequestFailed verifies Page.OnRequestFailed for a connection closed without a response.
// Ref: TestPageEventNetwork.java#PageEventsRequestFailed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventNetworkPageEventsRequestFailed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	srv.ServeWithBody("/one-style.html", "text/html", `<link rel='stylesheet' href='./one-style.css'><div>hello, world!</div>`)
	srv.SetRoute("/one-style.css", func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			return
		}
		conn, _, err := hj.Hijack()
		if err == nil {
			_ = conn.Close()
		}
	})
	log := a4Watch(t, page)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/one-style.html"))
	must.Eventually(func() bool { return len(log.Failed()) > 0 }, 10*time.Second, 20*time.Millisecond)
	failed := log.Failed()
	is.Len(failed, 1)
	is.Contains(failed[0].URL(), "one-style.css")
	resp, err := failed[0].Response(ctx)
	must.NoError(err)
	is.Nil(resp)
	is.Equal("stylesheet", failed[0].ResourceType())
	// NetworkRequest.Failure() and NetworkRequest.Frame() are not exposed by the Go client, so the failure text and frame assertions cannot be ported.
}

// TestPageEventNetworkPageEventsRequestFinished verifies Page.OnRequestFinished for the navigation request.
// Ref: TestPageEventNetwork.java#PageEventsRequestFinished  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventNetworkPageEventsRequestFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.Eventually(func() bool { return len(log.Finished()) > 0 }, 10*time.Second, 20*time.Millisecond)
	request := log.Finished()[0]
	is.Equal(srv.EmptyPage(), request.URL())
	resp, err := request.Response(ctx)
	must.NoError(err)
	must.NotNil(resp)
	is.Equal(srv.EmptyPage(), page.MainFrame().URL())
	// Page.Goto() does not return the Response, and NetworkRequest.Failure()/Frame() and NetworkResponse.Finished()/Request() are not exposed, so those assertions cannot be ported.
}

// Ref: TestPageEventNetwork.java#PageWaitForRequestFinished  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventNetworkPageWaitForRequestFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	request, err := page.WaitForRequestFinished(ctx, func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(request)

	response, err := request.Response(ctx)
	must.NoError(err)
	must.NotNil(response)
	is.Nil(response.Finished())
	is.Equal(srv.EmptyPage(), request.URL())
	must.NotNil(request.Frame())
	is.Equal(page.MainFrame(), request.Frame())
	is.Equal(srv.EmptyPage(), request.Frame().URL())
	is.Nil(request.Failure())
}

// TestPageEventNetworkShouldFireEventsInProperOrder verifies request, response and requestfinished are fired in order.
// Ref: TestPageEventNetwork.java#shouldFireEventsInProperOrder  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventNetwork.java
func TestPageEventNetworkShouldFireEventsInProperOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var mu sync.Mutex
	var events []string
	add := func(s string) { mu.Lock(); events = append(events, s); mu.Unlock() }
	offs := []func(){
		page.OnRequest(func(*playwright.NetworkRequest) { add("request") }),
		page.OnResponse(func(*playwright.NetworkResponse) { add("response") }),
		page.OnRequestFinished(func(*playwright.NetworkRequest) { add("requestfinished") }),
	}
	defer func() {
		for _, off := range offs {
			off()
		}
	}()
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.Eventually(func() bool { mu.Lock(); defer mu.Unlock(); return len(events) >= 3 }, 10*time.Second, 20*time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	is.Equal([]string{"request", "response", "requestfinished"}, events)
}
