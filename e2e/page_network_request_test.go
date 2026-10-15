//go:build e2e

// E2E tests for: TestPageNetworkRequest.java.
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

// TestNetworkRequestURLNotEmpty verifies Request.URL returns non-empty URL.
// Extra: no equivalent method in playwright-java; related class: TestPageNetworkRequest.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
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

// TestNetworkRequestURLMatchesEmptyPage verifies request URL matches navigated URL.
// Extra: no equivalent method in playwright-java; related class: TestPageNetworkRequest.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
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

// Ref: TestPageNetworkRequest.java#shouldReportRawHeaders  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
func TestPageNetworkRequestShouldReportRawHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	done := make(chan struct{}, 1)
	srv.SetRoute("/headers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		select {
		case done <- struct{}{}:
		default:
		}
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	request, err := page.WaitForRequest(ctx, "**/headers", func() error {
		_, err := page.Evaluate(ctx, `() => fetch('/headers', {
			headers: [
				['header-a', 'value-a'],
				['header-b', 'value-b'],
				['header-a', 'value-a-1'],
				['header-a', 'value-a-2'],
			]
		})`)
		return err
	})
	must.NoError(err)
	must.NotNil(request)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for route")
	}

	rawHeaders, err := request.HeadersArray(ctx)
	must.NoError(err)
	must.NotEmpty(rawHeaders)

	valA, err := request.HeaderValue(ctx, "header-a")
	must.NoError(err)
	must.NotNil(valA)
	is.Equal("value-a, value-a-1, value-a-2", *valA)

	valMissing, err := request.HeaderValue(ctx, "not-there")
	must.NoError(err)
	is.Nil(valMissing)
}

// Ref: TestPageNetworkRequest.java#shouldReportAllCookiesInOneHeader  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
func TestPageNetworkRequestShouldReportAllCookiesInOneHeader(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, `() => {
		document.cookie = 'myCookie=myValue';
		document.cookie = 'myOtherCookie=myOtherValue';
	}`)
	must.NoError(err)

	request, err := page.WaitForRequest(ctx, srv.EmptyPage(), func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(request)

	allHeaders, err := request.AllHeaders(ctx)
	must.NoError(err)
	cookie := allHeaders["cookie"]
	is.Equal("myCookie=myValue; myOtherCookie=myOtherValue", cookie)
}

// TestPageNetworkRequestShouldReportPostDataFor403Request verifies post data of a request answered with 403.
// Ref: TestPageNetworkRequest.java#shouldReportPostDataFor403Request  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
func TestPageNetworkRequestShouldReportPostDataFor403Request(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	var mu sync.Mutex
	var serverBody string
	srv.SetRoute("/upload", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		serverBody = string(b)
		mu.Unlock()
		w.WriteHeader(403)
	})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	log := a4Watch(t, page)
	_, err := page.Evaluate(ctx, "() => fetch('/upload', { method: 'POST', body: 'test'})")
	must.NoError(err)
	req := log.WaitRequest(t, "/upload")
	is.Eventually(func() bool { mu.Lock(); defer mu.Unlock(); return serverBody != "" }, 5*time.Second, 20*time.Millisecond)
	mu.Lock()
	is.Equal("test", serverBody)
	mu.Unlock()
	is.Equal("test", req.PostData())
	is.Equal("POST", req.Method())
	resp := log.WaitResponse(t, "/upload", 0)
	is.Equal(403, resp.Status())
}

// Ref: TestPageNetworkRequest.java#shouldNotAllowToAccessFrameOnPopupMainRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
func TestPageNetworkRequestShouldNotAllowToAccessFrameOnPopupMainRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<a target=_blank href='"+srv.EmptyPage()+"'>click me</a>"))

	var mu sync.Mutex
	var capturedReq *playwright.NetworkRequest
	var panicMsg any

	off := page.Context().OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		capturedReq = req
		defer func() {
			panicMsg = recover()
		}()
		_ = req.Frame()
	})
	defer off()

	must.NoError(page.Locator("text=click me").Click(ctx))

	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return capturedReq != nil
	}, 10*time.Second, 50*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	is.True(capturedReq.IsNavigationRequest())
	must.NotNil(panicMsg)
	is.Contains(fmt.Sprint(panicMsg), "Frame for this navigation request is not available")
}

// TestPageNetworkRequestShouldParseDataIfContentTypeIsApplicationFormUrlEncoded verifies PostData of a form submission.
// Ref: TestPageNetworkRequest.java#shouldParseDataIfContentTypeIsApplicationFormUrlEncoded  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
func TestPageNetworkRequestShouldParseDataIfContentTypeIsApplicationFormUrlEncoded(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	srv.SetRoute("/post", func(w http.ResponseWriter, r *http.Request) {})
	log := a4Watch(t, page)
	must.NoError(page.SetContent(ctx, "<form method='POST' action='/post'><input type='text' name='foo' value='bar'><input type='number' name='baz' value='123'><input type='submit'></form>"))
	must.NoError(page.Locator("input[type=submit]").Click(ctx))
	req := log.WaitRequest(t, "/post")
	is.Equal("foo=bar&baz=123", req.PostData())
}

// TestPageNetworkRequestShouldParseDataIfContentTypeIsApplicationFormUrlEncodedCharsetUTF8 verifies PostData with a charset content type.
// Ref: TestPageNetworkRequest.java#shouldParseDataIfContentTypeIsApplicationFormUrlEncodedCharsetUTF8  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageNetworkRequest.java
func TestPageNetworkRequestShouldParseDataIfContentTypeIsApplicationFormUrlEncodedCharsetUTF8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	srv.SetRoute("/post", func(w http.ResponseWriter, r *http.Request) {})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	log := a4Watch(t, page)
	_, err := page.Evaluate(ctx, `() => fetch('./post', {
  method: 'POST',
  headers: {
    'content-type': 'application/x-www-form-urlencoded; charset=UTF-8',
  },
  body: 'foo=bar&baz=123'
})`)
	must.NoError(err)
	req := log.WaitRequest(t, "/post")
	is.Equal("foo=bar&baz=123", req.PostData())
}
