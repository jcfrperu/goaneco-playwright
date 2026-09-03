//go:build e2e

// E2E tests for page network event listeners and URL waiting.
package e2e

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

func TestPageOnRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/track-me", "text/html", `<p>hello</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	var mu sync.Mutex
	var seen []*playwright.NetworkRequest
	done := make(chan struct{})
	var once sync.Once

	page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		seen = append(seen, req)
		mu.Unlock()
		once.Do(func() { close(done) })
	})

	err = page.Goto(ctx, srv.Prefix()+"/track-me")
	must.NoError(err, "Goto /track-me failed")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnRequest to fire")
	}

	mu.Lock()
	reqs := seen
	mu.Unlock()

	found := false
	for _, r := range reqs {
		if r.URL() == srv.Prefix()+"/track-me" {
			found = true
			if r.Method() != "GET" {
				t.Errorf("request method = %q, want GET", r.Method())
			}
			break
		}
	}
	if !found {
		t.Errorf("did not receive OnRequest for /track-me; got %d requests", len(reqs))
	}
}

func TestPageOnResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/track-response", "text/html", `<p>response</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	var mu sync.Mutex
	var seen []*playwright.NetworkResponse
	done := make(chan struct{})
	var once sync.Once

	page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		seen = append(seen, resp)
		mu.Unlock()
		once.Do(func() { close(done) })
	})

	err = page.Goto(ctx, srv.Prefix()+"/track-response")
	must.NoError(err, "Goto /track-response failed")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnResponse to fire")
	}

	mu.Lock()
	resps := seen
	mu.Unlock()

	found := false
	for _, r := range resps {
		if r.URL() == srv.Prefix()+"/track-response" {
			found = true
			if r.Status() != 200 {
				t.Errorf("response status = %d, want 200", r.Status())
			}
			break
		}
	}
	if !found {
		t.Errorf("did not receive OnResponse for /track-response; got %d responses", len(resps))
	}
}

func TestPageWaitForURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/start", "text/html", `<p>start</p>`)
	srv.ServeWithBody("/destination", "text/html", `<p>destination</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/start")
	must.NoError(err, "Goto /start failed")

	// Trigger client-side navigation in background
	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `window.location.href = '`+srv.Prefix()+`/destination'`)
	}()

	err = page.WaitForURL(ctx, srv.Prefix()+"/destination", 5*time.Second)
	must.NoError(err, "WaitForURL failed")

	if page.URL() != srv.Prefix()+"/destination" {
		t.Errorf("URL = %q, want /destination", page.URL())
	}
}

func TestPageWaitForURLAlreadyMatches(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// WaitForURL with a wildcard pattern matching the current URL should return immediately
	err = page.WaitForURL(ctx, srv.Prefix()+"*", time.Second)
	must.NoError(err, "WaitForURL (already matches) failed")
}

func TestPageOnRequestURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var urls []string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(urls)
	is.Contains(urls, srv.EmptyPage())
}

func TestPageOnResponseURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var respURLs []string

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		respURLs = append(respURLs, resp.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(respURLs)
	is.Contains(respURLs, srv.EmptyPage())
}

func TestPageOnRequestMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var methods []string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		methods = append(methods, req.Method())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(methods)
	is.Contains(methods, "GET")
}

func TestPageOnResponseStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var statuses []int

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statuses = append(statuses, resp.Status())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(statuses)
	is.Contains(statuses, 200)
}

func TestPageOnResponseOKForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var emptyPageOK *bool

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			ok := resp.OK()
			emptyPageOK = &ok
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(emptyPageOK, "OnResponse should have fired for empty page")
	is.True(*emptyPageOK, "empty page response should be OK")
}

func TestPageOnRequestHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var headers map[string]string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			headers = req.Headers()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(headers, "headers should have been captured")

	_, hasUA := headers["user-agent"]
	is.True(hasUA, "request headers should include user-agent")
}

func TestPageOnRequestIsNavigationRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var navRequestFound bool

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() && req.IsNavigationRequest() {
			mu.Lock()
			navRequestFound = true
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.True(navRequestFound, "should have found a navigation request to empty page")
}

func TestPageOnRequestFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	requests := make([]*playwright.NetworkRequest, 0)

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, req)
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	count := len(requests)
	mu.Unlock()

	is.Greater(count, 0)
}

func TestPageOnResponseFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	responses := make([]*playwright.NetworkResponse, 0)

	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		responses = append(responses, resp)
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	count := len(responses)
	mu.Unlock()

	is.Greater(count, 0)
}

func TestPageOnRequestFinishedFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	finished := make([]*playwright.NetworkRequest, 0)

	off := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		mu.Lock()
		finished = append(finished, req)
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	count := len(finished)
	mu.Unlock()

	is.Greater(count, 0)
}

func TestPageOffRemovesRequestHandler(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	count := 0

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	countBeforeOff := count
	mu.Unlock()

	off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	countAfterOff := count
	mu.Unlock()

	is.Equal(countBeforeOff, countAfterOff)
}

func TestPageOnRequestFiredForFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/fetch-api", "application/json", `{"ok":true}`)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var urls []string
	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer off()

	_, err := page.Evaluate(ctx, `() => fetch('/fetch-api')`)
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	u := urls
	mu.Unlock()

	is.NotEmpty(u)
}

func TestPageOnResponseFiredForFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/resp-api", "application/json", `{"ok":true}`)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var statuses []int
	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statuses = append(statuses, resp.Status())
		mu.Unlock()
	})
	defer off()

	_, err := page.Evaluate(ctx, `() => fetch('/resp-api')`)
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	s := statuses
	mu.Unlock()

	is.NotEmpty(s)
}

func TestPageOnRequestMethodIsGet(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var method string
	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if method == "" {
			method = req.Method()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal("GET", method)
}

func TestPageOnRequestFinishedFiredEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var finished *playwright.NetworkRequest
	var mu sync.Mutex

	off := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			finished = req
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	f := finished
	mu.Unlock()

	must.NotNil(f)
	is.Equal("GET", f.Method())
}

func TestPageOnRequestFailedFiredOnAbortEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var failed *playwright.NetworkRequest
	var mu sync.Mutex

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		failed = req
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Route(ctx, "**/blocked", func(route *playwright.Route) {
		must.NoError(route.Abort(ctx))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `() => fetch('/blocked').catch(() => {})`)
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	f := failed
	mu.Unlock()

	must.NotNil(f)
}

func TestPageOnConsoleMessageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var received string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "log" {
			mu.Lock()
			received = msg.Text()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.log('hello from console')</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	r := received
	mu.Unlock()

	is.Contains(r, "hello from console")
}

// TestOnPopupReceivesPageReference verifies OnPopup callback receives the opened page.
// Ref: TestPagePopup.java#shouldReceivePageReference
func TestOnPopupReceivesPageReference(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var popup *playwright.Page

	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		popup = p
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.SetContent(ctx, `<button onclick="window.open('about:blank')">open</button>`))
	must.NoError(page.Locator("button").Click(ctx))

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		p := popup
		mu.Unlock()
		if p != nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	mu.Lock()
	p := popup
	mu.Unlock()
	must.NotNil(p, "popup page should not be nil")
}

// TestOnPopupIsNotFiredForSameTab verifies OnPopup is NOT fired for same-tab navigation.
// Ref: TestPagePopup.java#shouldNotFireForSameTab
func TestOnPopupIsNotFiredForSameTab(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	popupFired := false

	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		popupFired = true
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.SetContent(ctx, `<a href="about:blank">same tab link</a>`))
	must.NoError(page.Locator("a").Click(ctx))

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	is.False(popupFired)
}

// TestOnPopupCanNavigate verifies the popup page can be navigated after opening.
// Ref: TestPagePopup.java#shouldNavigatePopup
func TestOnPopupCanNavigate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	done := make(chan *playwright.Page, 1)

	cancel := page.OnPopup(func(p *playwright.Page) {
		select {
		case done <- p:
		default:
		}
	})
	defer cancel()

	must.NoError(page.SetContent(ctx, `<button onclick="window.open('about:blank', '_blank')">open</button>`))
	must.NoError(page.Locator("button").Click(ctx))

	var popup *playwright.Page
	select {
	case popup = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for popup")
	}

	must.NotNil(popup)
	must.NoError(popup.SetContent(ctx, `<div>popup content</div>`))

	text, err := popup.Locator("div").InnerText(ctx)
	must.NoError(err)
	is.Equal("popup content", text)
}

func TestPageOnRequestEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var requestURL string
	var mu sync.Mutex

	page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requestURL = req.URL()
		mu.Unlock()
	})

	body := "hello"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/test"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	u := requestURL
	mu.Unlock()

	is.Contains(u, "example.com")
}

func TestPageOnResponseEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var responseStatus int
	var mu sync.Mutex

	page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		responseStatus = resp.Status()
		mu.Unlock()
	})

	status := 200
	body := "response body"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status,
			Body:   &body,
		}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	s := responseStatus
	mu.Unlock()

	is.Equal(200, s)
}

func TestPageOnRequestFinishedEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var finishedURL string
	var mu sync.Mutex

	page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		mu.Lock()
		finishedURL = req.URL()
		mu.Unlock()
	})

	body := "done"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/finish"))
	must.NoError(page.WaitForTimeout(ctx, 150))

	mu.Lock()
	u := finishedURL
	mu.Unlock()

	is.Contains(u, "example.com")
}
