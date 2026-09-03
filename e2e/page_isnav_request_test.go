//go:build e2e

// Additional page navigation request E2E tests.
// Migration of: TestNetworkEvents.java (navigation request cases)
package e2e

import (
	"sync"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsNavigationRequestTrueForGoto verifies navigation requests from Goto are marked as nav requests.
// Ref: TestNetworkEvents.java#shouldMarkNavigationRequestsAsNavigation
func TestIsNavigationRequestTrueForGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	navCount := 0

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			navCount++
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.GreaterOrEqual(navCount, 1)
}

// TestIsNavigationRequestFalseForFetch verifies fetch requests are NOT navigation requests.
// Ref: TestNetworkEvents.java#shouldNotMarkFetchAsNavigation
func TestIsNavigationRequestFalseForFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	srv.ServeWithBody("/api/data", "application/json", `{}`)

	var mu sync.Mutex
	var fetchIsNav *bool

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.Prefix()+"/api/data" {
			mu.Lock()
			v := req.IsNavigationRequest()
			fetchIsNav = &v
			mu.Unlock()
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => fetch('/api/data').catch(() => {})`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	if fetchIsNav != nil {
		is.False(*fetchIsNav)
	}
}

// TestRequestMethodIsGETForNavigation verifies navigation requests use GET method.
// Ref: TestNetworkEvents.java#shouldReportMethodForNavigation
func TestRequestMethodIsGETForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var method string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			method = req.Method()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.Equal("GET", method)
}

// TestNetworkEventOrderRequestBeforeResponse verifies request fires before response for same resource.
// Ref: TestNetworkEvents.java#shouldFireEventsInProperOrder
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
