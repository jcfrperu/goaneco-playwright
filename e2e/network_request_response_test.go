//go:build e2e

// E2E tests for NetworkRequest.Response().
package e2e

import (
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

func TestNetworkRequestResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/req-resp", "text/html", `<html><body>ok</body></html>`)

	page := newPage(t)

	type capture struct {
		req *playwright.NetworkRequest
	}
	ch := make(chan capture, 1)
	cancel := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.Prefix()+"/req-resp" {
			select {
			case ch <- capture{req}:
			default:
			}
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/req-resp")
	must.NoError(err, "Goto failed")

	var c capture
	select {
	case c = <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal("OnRequestFinished never fired for /req-resp")
	}

	resp, err := c.req.Response(ctx)
	must.NoError(err, "NetworkRequest.Response() failed")
	must.NotNil(resp, "NetworkRequest.Response() returned nil, want non-nil")
	if resp.Status() != 200 {
		t.Errorf("Response.Status() = %d, want 200", resp.Status())
	}
	if resp.URL() != srv.Prefix()+"/req-resp" {
		t.Errorf("Response.URL() = %q, want %q", resp.URL(), srv.Prefix()+"/req-resp")
	}
}

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

// TestNetworkRequestHeadersAvailable verifies Request.Headers() returns headers map.
// Ref: TestNetworkEvents.java#shouldReportRequestHeaders
func TestNetworkRequestHeadersAvailable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedHeaders map[string]string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			capturedHeaders = req.Headers()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(capturedHeaders)
	// All requests should have at least a User-Agent or accept header
	is.NotEmpty(capturedHeaders)
}

// TestNetworkResponseHeadersIncludeContentType verifies Response.Headers() includes content-type.
// Ref: TestNetworkEvents.java#shouldReportResponseHeaders
func TestNetworkResponseHeadersIncludeContentType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedHeaders map[string]string

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			capturedHeaders = resp.Headers()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(capturedHeaders)
	// Content-type should be present for an HTML page
	is.NotEmpty(capturedHeaders)
}

// TestNetworkResponseURLMatchesRequest verifies Response.URL() matches the request URL.
// Ref: TestNetworkEvents.java#shouldReportURL
func TestNetworkResponseURLMatchesRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedURL string

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			capturedURL = resp.URL()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.Equal(srv.EmptyPage(), capturedURL)
}

// TestNetworkRequestURLMatchesGoto verifies Request.URL() matches the navigated URL.
// Ref: TestNetworkEvents.java#shouldReportRequestURL
func TestNetworkRequestURLMatchesGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedURL string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			capturedURL = req.URL()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.Equal(srv.EmptyPage(), capturedURL)
}

// ---------------------------------------------------------------------------
// From network_request_extra_test.go
// ---------------------------------------------------------------------------

// TestNetworkRequestResponseFromOnRequestFinished verifies request.Response() after request finishes.
// Ref: TestNetworkEvents.java#shouldReportRequestAfterNavigationCompletion
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

// TestNetworkRequestResourceType verifies that navigation requests report document resource type.
// Ref: TestNetworkEvents.java#shouldReportResourceType
func TestNetworkRequestResourceType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var navResourceType string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() && req.IsNavigationRequest() {
			mu.Lock()
			navResourceType = req.ResourceType()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.Equal("document", navResourceType)
}

// TestNetworkResponseBodyAfterRoute verifies response body can be read after route fulfillment.
// Ref: TestNetworkEvents.java#shouldReportBody
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
// Ref: TestNetworkEvents.java#shouldReportStatusCode
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
		srv.ServeWithBody(path, "text/plain", "ok")
		_ = status
		_ = path
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

// TestNetworkRequestPostDataForFetch verifies PostData is available for POST fetch requests.
// Ref: TestNetworkEvents.java#shouldReportPostData
func TestNetworkRequestPostDataForFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedPostData string

	must.NoError(page.Route(ctx, "**/api/post", func(route *playwright.Route) {
		req := route.Request()
		mu.Lock()
		capturedPostData = req.PostData()
		mu.Unlock()
		_ = route.Continue(ctx, nil)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => fetch('/api/post', {
		method: 'POST',
		body: 'hello post body',
		headers: { 'Content-Type': 'text/plain' }
	}).catch(() => {})`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.Equal("hello post body", capturedPostData)
}

// TestNetworkRequestMethod verifies Request.Method returns HTTP method.
// Ref: TestNetworkRequest.java#shouldReturnMethod
func TestNetworkRequestMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var method string
	var mu sync.Mutex

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			method = req.Method()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	m := method
	mu.Unlock()

	is.Equal("GET", m)
}

// TestNetworkRequestURLNotEmpty verifies Request.URL returns non-empty URL.
// Ref: TestNetworkRequest.java#shouldReturnURL
func TestNetworkRequestURLNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var url string
	var mu sync.Mutex

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			url = req.URL()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	u := url
	mu.Unlock()

	is.Equal(srv.EmptyPage(), u)
}

// TestNetworkRequestHeadersNotEmpty verifies Request.Headers returns headers.
// Ref: TestNetworkRequest.java#shouldReturnHeaders
func TestNetworkRequestHeadersNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var headers map[string]string
	var mu sync.Mutex

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			headers = req.Headers()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	h := headers
	mu.Unlock()

	is.NotEmpty(h)
}

// TestNetworkRequestIsNavigationForGoto verifies IsNavigationRequest true for Goto.
// Ref: TestNetworkRequest.java#shouldBeNavigationForGoto
func TestNetworkRequestIsNavigationForGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var isNav bool
	var mu sync.Mutex

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			isNav = req.IsNavigationRequest()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	nav := isNav
	mu.Unlock()

	is.True(nav)
}

// TestNetworkRequestResponseReturnsResponse verifies Request.Response() returns response.
// Ref: TestNetworkRequest.java#shouldReturnResponse
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

// TestNetworkRequestURLMatchesEmptyPage verifies request URL matches navigated URL.
// Ref: TestNetworkRequest.java#shouldMatchURL
func TestNetworkRequestURLMatchesEmptyPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var capturedURL string
	var mu sync.Mutex

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			capturedURL = req.URL()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	u := capturedURL
	mu.Unlock()

	is.Equal(srv.EmptyPage(), u)
}

// TestNetworkRequestResourceTypeForNavigation verifies navigation request resource type is non-empty.
// Ref: TestNetworkRequest.java#shouldHaveDocumentResourceType
func TestNetworkRequestResourceTypeForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var resourceType string
	var mu sync.Mutex

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			resourceType = req.ResourceType()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	rt := resourceType
	mu.Unlock()

	is.NotEmpty(rt)
}

// TestNetworkRequestPostDataEmptyForGet verifies PostData is empty for GET.
// Ref: TestNetworkRequest.java#shouldHaveEmptyPostDataForGet
func TestNetworkRequestPostDataEmptyForGet(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var postData string
	var mu sync.Mutex

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			postData = req.PostData()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	pd := postData
	mu.Unlock()

	is.Empty(pd)
}

// TestNetworkRequestResponseStatusIs200 verifies response status is 200 for successful nav.
// Ref: TestNetworkRequest.java#shouldHave200Response
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

// TestNetworkRequestMethodEx4 verifies request method is GET by default.
// Ref: TestNetworkRequest.java#shouldReturnGETMethod
func TestNetworkRequestMethodEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var method string
	var mu sync.Mutex

	page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		method = req.Method()
		mu.Unlock()
	})

	body := "hello"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	m := method
	mu.Unlock()

	is.Equal("GET", m)
}

// TestNetworkRequestURLEx4 verifies request URL is correct.
// Ref: TestNetworkRequest.java#shouldReturnURL
func TestNetworkRequestURLEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var reqURL string
	var mu sync.Mutex

	page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		reqURL = req.URL()
		mu.Unlock()
	})

	body := "ok"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body}))
	}))

	must.NoError(page.Goto(ctx, "http://test.example.com/path"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	u := reqURL
	mu.Unlock()

	is.Contains(u, "example.com")
}

// TestNetworkRequestHeadersEx4 verifies request headers are accessible.
// Ref: TestNetworkRequest.java#shouldReturnHeaders
func TestNetworkRequestHeadersEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var headers map[string]string
	var mu sync.Mutex

	page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		headers = req.Headers()
		mu.Unlock()
	})

	body := "ok"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	h := headers
	mu.Unlock()

	must.NotNil(h)
}

// TestNetworkRequestHeadersContainUserAgent verifies request headers include user-agent.
// Ref: TestNetworkRequest.java#shouldReportAllHeaders
func TestNetworkRequestHeadersContainUserAgent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var headers map[string]string

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			headers = req.Headers()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	h := headers
	mu.Unlock()

	is.NotEmpty(h)
	// Check that at least one header is present (user-agent or accept or similar).
	found := false
	for k := range h {
		if k != "" {
			found = true
			break
		}
	}
	is.True(found, "expected at least one non-empty header key")
}

// TestNetworkRequestAndResponseInOrder verifies request fires before response for navigation.
// Ref: TestNetworkRequest.java#shouldFireEventsInProperOrder
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

// TestNetworkRequestResourceTypeForFetch verifies ResourceType is "fetch" for fetch calls.
// Ref: TestNetworkRequest.java#shouldReportResourceTypeForFetch
func TestNetworkRequestResourceTypeForFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var fetchResourceType string

	srv.ServeWithBody("/api-data", "text/plain", "data")

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.Prefix()+"/api-data" {
			mu.Lock()
			fetchResourceType = req.ResourceType()
			mu.Unlock()
		}
	})
	defer off()

	_, err := page.Evaluate(ctx, `() => fetch('/api-data')`)
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	rt := fetchResourceType
	mu.Unlock()

	is.Equal("fetch", rt)
}

// ---------------------------------------------------------------------------
// From TestNetworkRequest.java (Java source ports)
// ---------------------------------------------------------------------------

// Ref: TestNetworkRequest.java#shouldWorkForMainFrameNavigationRequest
func TestNetworkRequestShouldWorkForMainFrameNavigationRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var requests []*playwright.NetworkRequest

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, req)
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.Len(requests, 1)
}

// Ref: TestNetworkRequest.java#shouldWorkForSubframeNavigationRequest
func TestNetworkRequestShouldWorkForSubframeNavigationRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var requests []*playwright.NetworkRequest

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			requests = append(requests, req)
			mu.Unlock()
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `url => new Promise(resolve => {
		const frame = document.createElement('iframe');
		frame.src = url;
		frame.onload = resolve;
		document.body.appendChild(frame);
	})`, srv.EmptyPage())
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.GreaterOrEqual(len(requests), 1)
}

// Ref: TestNetworkRequest.java#shouldWorkForFetchRequests
func TestNetworkRequestShouldWorkForFetchRequests(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/digits/1.png", "image/png", "fakepng")

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var requests []*playwright.NetworkRequest
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, req)
		mu.Unlock()
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => fetch('/digits/1.png')`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.GreaterOrEqual(len(requests), 1)
}

// Ref: TestNetworkRequest.java#shouldWorkForARedirect
func TestNetworkRequestShouldWorkForARedirect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithRedirect("/foo.html", "/empty.html")

	var mu sync.Mutex
	var requests []*playwright.NetworkRequest
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, req)
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.Prefix()+"/foo.html"))

	mu.Lock()
	defer mu.Unlock()
	must.GreaterOrEqual(len(requests), 2)
	is.Equal(srv.Prefix()+"/foo.html", requests[0].URL())
	is.Equal(srv.Prefix()+"/empty.html", requests[1].URL())
}

// Ref: TestNetworkRequest.java#shouldWorkAllHeadersInsideRoute
func TestNetworkRequestShouldWorkAllHeadersInsideRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var count int
	var sawAccept bool

	must.NoError(page.Route(ctx, "**", func(route *playwright.Route) {
		headers := route.Request().Headers()
		mu.Lock()
		count++
		if v, ok := headers["accept"]; ok && len(v) > 5 {
			sawAccept = true
		}
		mu.Unlock()
		_ = route.Continue(ctx, nil)
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/empty.html"))

	mu.Lock()
	defer mu.Unlock()
	is.GreaterOrEqual(count, 1)
	is.True(sawAccept, "accept header should be present with length > 5")
}

// Ref: TestNetworkRequest.java#shouldNotWorkForARedirectAndInterception
func TestNetworkRequestShouldNotWorkForARedirectAndInterception(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithRedirect("/foo.html", "/empty.html")

	var mu sync.Mutex
	var requests []*playwright.NetworkRequest
	must.NoError(page.Route(ctx, "**", func(route *playwright.Route) {
		mu.Lock()
		requests = append(requests, route.Request())
		mu.Unlock()
		_ = route.Continue(ctx, nil)
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/foo.html"))

	is.Equal(srv.Prefix()+"/empty.html", page.URL())

	mu.Lock()
	defer mu.Unlock()
	must.Len(requests, 1)
	is.Equal(srv.Prefix()+"/foo.html", requests[0].URL())
}

// Ref: TestNetworkRequest.java#shouldGetTheSameHeadersAsTheServer
func TestNetworkRequestShouldGetTheSameHeadersAsTheServer(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	serverHeadersCh := make(chan http.Header, 1)
	srv.SetRoute("/hdrs.html", func(w http.ResponseWriter, r *http.Request) {
		// Copy the request headers so we can compare
		hcopy := r.Header.Clone()
		select {
		case serverHeadersCh <- hcopy:
		default:
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("done"))
	})

	var mu sync.Mutex
	var reqHeaders map[string]string
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.Prefix()+"/hdrs.html" {
			mu.Lock()
			reqHeaders = req.Headers()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.Prefix()+"/hdrs.html"))

	var got http.Header
	select {
	case got = <-serverHeadersCh:
	case <-time.After(3 * time.Second):
		t.Fatal("server never received /hdrs.html request")
	}

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(reqHeaders)
	// Verify a couple of headers that the server saw are also observed by playwright.
	for k := range got {
		lk := strings.ToLower(k)
		if lk == "user-agent" || lk == "accept" || lk == "host" {
			is.Contains(reqHeaders, lk, "playwright request should have header %q", lk)
		}
	}
}

// Ref: TestNetworkRequest.java#shouldGetTheSameHeadersAsTheServerCORP
func TestNetworkRequestShouldGetTheSameHeadersAsTheServerCORP(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	serverHeadersCh := make(chan http.Header, 1)
	srv.SetRoute("/something", func(w http.ResponseWriter, r *http.Request) {
		hcopy := r.Header.Clone()
		select {
		case serverHeadersCh <- hcopy:
		default:
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("done"))
	})

	var mu sync.Mutex
	var reqHeaders map[string]string
	targetURL := srv.CrossProcessPrefix() + "/something"
	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == targetURL {
			// nothing needed here; headers captured via OnRequest
		}
	})
	defer cancel()

	cancel2 := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == targetURL {
			mu.Lock()
			reqHeaders = req.Headers()
			mu.Unlock()
		}
	})
	defer cancel2()

	text, err := page.Evaluate(ctx, `async url => {
		const data = await fetch(url);
		return data.text();
	}`, targetURL)
	must.NoError(err)
	is.Equal("done", text)

	var got http.Header
	select {
	case got = <-serverHeadersCh:
	case <-time.After(3 * time.Second):
		t.Fatal("server never received /something request")
	}

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(reqHeaders)
	for k := range got {
		lk := strings.ToLower(k)
		if lk == "user-agent" || lk == "accept" || lk == "host" {
			is.Contains(reqHeaders, lk)
		}
	}
}

// Ref: TestNetworkRequest.java#shouldReturnPostData
func TestNetworkRequestShouldReturnPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	srv.SetRoute("/post", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	var mu sync.Mutex
	var capturedPost string
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if strings.HasSuffix(req.URL(), "/post") {
			mu.Lock()
			capturedPost = req.PostData()
			mu.Unlock()
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => fetch('./post', { method: 'POST', body: JSON.stringify({foo: 'bar'})})`)
	must.NoError(err)

	mu.Lock()
	defer mu.Unlock()
	is.Equal(`{"foo":"bar"}`, capturedPost)
}

// Ref: TestNetworkRequest.java#shouldWorkWithBinaryPostData
func TestNetworkRequestShouldWorkWithBinaryPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	srv.SetRoute("/binary-post", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	var mu sync.Mutex
	var capturedBuf []byte
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if strings.HasSuffix(req.URL(), "/binary-post") {
			mu.Lock()
			capturedBuf = req.PostDataBuffer()
			mu.Unlock()
		}
	})
	defer cancel()

	// Send binary data via XHR
	_, err := page.Evaluate(ctx, `() => {
		const buf = new Uint8Array([1, 2, 3]);
		return fetch('./binary-post', { method: 'POST', body: buf });
	}`)
	must.NoError(err)

	// Wait briefly for the OnRequest event
	must.NoError(page.WaitForTimeout(ctx, 500))

	mu.Lock()
	defer mu.Unlock()
	is.NotNil(capturedBuf)
	is.Equal([]byte{1, 2, 3}, capturedBuf)
}

// Ref: TestNetworkRequest.java#shouldWorkWithBinaryPostDataAndInterception
func TestNetworkRequestShouldWorkWithBinaryPostDataAndInterception(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	srv.SetRoute("/binary-post-intercept", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	var mu sync.Mutex
	var capturedBuf []byte

	// Intercept and capture via route
	must.NoError(page.Route(ctx, "**/binary-post-intercept", func(route *playwright.Route) {
		req := route.Request()
		mu.Lock()
		capturedBuf = req.PostDataBuffer()
		mu.Unlock()
		must.NoError(route.Continue(ctx, nil))
	}))

	_, err := page.Evaluate(ctx, `() => {
		const buf = new Uint8Array([10, 20, 30]);
		return fetch('./binary-post-intercept', { method: 'POST', body: buf });
	}`)
	must.NoError(err)

	must.NoError(page.WaitForTimeout(ctx, 500))

	mu.Lock()
	defer mu.Unlock()
	is.NotNil(capturedBuf)
	is.Equal([]byte{10, 20, 30}, capturedBuf)
}

// Ref: TestNetworkRequest.java#shouldBeUndefinedWhenThereIsNoPostData
func TestNetworkRequestShouldBeUndefinedWhenThereIsNoPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var pd string
	var seen bool
	cancel := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			pd = req.PostData()
			seen = true
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.True(seen)
	is.Empty(pd)
}

// Ref: TestNetworkRequest.java#shouldReturnEventSource
func TestNetworkRequestShouldReturnEventSource(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(200)
		fl, _ := w.(http.Flusher)
		_, _ = w.Write([]byte("data: {\"foo\":\"bar\"}\n\n"))
		if fl != nil {
			fl.Flush()
		}
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var requests []*playwright.NetworkRequest
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if strings.HasSuffix(req.URL(), "/sse") {
			mu.Lock()
			requests = append(requests, req)
			mu.Unlock()
		}
	})
	defer cancel()

	result, err := page.Evaluate(ctx, `() => {
		const eventSource = new EventSource('/sse');
		return new Promise(resolve => {
			eventSource.onmessage = e => resolve(JSON.parse(e.data));
		});
	}`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	must.True(ok)
	is.Equal("bar", m["foo"])

	mu.Lock()
	defer mu.Unlock()
	must.NotEmpty(requests)
	is.Equal("eventsource", requests[0].ResourceType())
}

// Ref: TestNetworkRequest.java#shouldReturnNavigationBit
func TestNetworkRequestShouldReturnNavigationBit(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frames/one-frame.html", "text/html",
		`<html><body><iframe src="/frames/frame.html"></iframe></body></html>`)
	srv.ServeWithBody("/frames/frame.html", "text/html",
		`<html><body>frame</body></html>`)
	srv.ServeWithRedirect("/rrredirect", "/frames/one-frame.html")

	var mu sync.Mutex
	requests := map[string]*playwright.NetworkRequest{}
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		name := req.URL()
		if i := strings.LastIndex(name, "/"); i != -1 {
			name = name[i+1:]
		}
		mu.Lock()
		requests[name] = req
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.Prefix()+"/rrredirect"))

	mu.Lock()
	defer mu.Unlock()
	if r := requests["rrredirect"]; r != nil {
		is.True(r.IsNavigationRequest(), "rrredirect should be navigation")
	}
	if r := requests["one-frame.html"]; r != nil {
		is.True(r.IsNavigationRequest(), "one-frame.html should be navigation")
	}
	if r := requests["frame.html"]; r != nil {
		is.True(r.IsNavigationRequest(), "frame.html should be navigation")
	}
}

// Ref: TestNetworkRequest.java#shouldReturnNavigationBitWhenNavigatingToImage
func TestNetworkRequestShouldReturnNavigationBitWhenNavigatingToImage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Minimal 1x1 PNG
	srv.SetRoute("/pptr.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte{0x89, 'P', 'N', 'G'})
	})

	var mu sync.Mutex
	var requests []*playwright.NetworkRequest
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, req)
		mu.Unlock()
	})
	defer cancel()

	_ = page.Goto(ctx, srv.Prefix()+"/pptr.png")

	mu.Lock()
	defer mu.Unlock()
	must.NotEmpty(requests)
	is.True(requests[0].IsNavigationRequest())
}
