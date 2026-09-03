//go:build e2e

// Page and BrowserContext network event tests.
// Migration of: TestPageEventNetwork.java, TestBrowserContextNetworkEvents.java
package e2e

import (
	"context"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
// Ref: TestNetworkEvents.java#shouldFireRequest
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
// Ref: TestNetworkEvents.java#shouldFireResponse
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
// Ref: TestNetworkEvents.java#shouldFireRequestFinished
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
// Ref: TestNetworkEvents.java#shouldReportRequestMethod
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
