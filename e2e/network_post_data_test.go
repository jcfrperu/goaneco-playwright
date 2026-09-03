//go:build e2e

// E2E tests for NetworkRequest.PostData.
// Migration of: TestNetworkPostData.java
package e2e

import (
	"strings"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureRequestTo waits for the first POST/PUT request to a URL containing urlSubstr.
func captureRequestTo(t *testing.T, page *playwright.Page, urlSubstr string, trigger func()) *playwright.NetworkRequest {
	t.Helper()
	ch := make(chan *playwright.NetworkRequest, 1)
	var once sync.Once
	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if strings.Contains(req.URL(), urlSubstr) {
			once.Do(func() {
				select {
				case ch <- req:
				default:
				}
			})
		}
	})
	defer cancel()

	trigger()

	select {
	case req := <-ch:
		return req
	case <-time.After(10 * time.Second):
		t.Fatalf("captureRequestTo: timed out waiting for request to URL containing %q", urlSubstr)
		return nil
	}
}

// TestNetworkRequestPostDataUtf8Body verifies PostData() returns the correct UTF-8 encoded body.
// Ref: TestNetworkPostData.java#shouldReturnCorrectPostDataBufferForUtf8Body
func TestNetworkRequestPostDataUtf8Body(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	value := "baẞ"
	req := captureRequestTo(t, page, "title.html", func() {
		_, err := page.Evaluate(ctx, `({url, value}) => {
			const request = new Request(url, {
				method: 'POST',
				body: JSON.stringify(value),
			});
			request.headers.set('content-type', 'application/json;charset=UTF-8');
			return fetch(request);
		}`, map[string]any{"url": srv.Prefix() + "/title.html", "value": value})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	postData := req.PostData()
	is.NotEmpty(postData, "PostData() should not be empty")
	is.Contains(postData, "baẞ", "PostData should contain the unicode value")
}

// TestNetworkRequestPostDataWithoutContentType verifies PostData() works when content-type is empty.
// Ref: TestNetworkPostData.java#shouldReturnPostDataWOContentType
func TestNetworkRequestPostDataWithoutContentType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureRequestTo(t, page, "title.html", func() {
		_, err := page.Evaluate(ctx, `({url}) => {
			const request = new Request(url, {
				method: 'POST',
				body: JSON.stringify({ value: 42 }),
			});
			request.headers.set('content-type', '');
			return fetch(request);
		}`, map[string]any{"url": srv.Prefix() + "/title.html"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	postData := req.PostData()
	is.Contains(postData, "42", "PostData should contain the JSON-encoded value 42")
}

// TestNetworkRequestPostDataForPut verifies PostData() works with PUT requests.
// Ref: TestNetworkPostData.java#shouldReturnPostDataForPUTRequests
func TestNetworkRequestPostDataForPut(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureRequestTo(t, page, "title.html", func() {
		_, err := page.Evaluate(ctx, `({url}) => {
			const request = new Request(url, {
				method: 'PUT',
				body: JSON.stringify({ value: 42 }),
			});
			return fetch(request);
		}`, map[string]any{"url": srv.Prefix() + "/title.html"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	postData := req.PostData()
	is.Contains(postData, "42", "PostData should contain value 42 for PUT request")
}
