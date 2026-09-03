//go:build e2e

// E2E tests for NetworkRequest.Sizes (request/response byte-level sizes).
// Migration of: TestPageNetworkSizes.java
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

// captureFinishedRequestTo waits for the first completed request to a URL containing urlSubstr.
func captureFinishedRequestTo(t *testing.T, page *playwright.Page, urlSubstr string, trigger func()) *playwright.NetworkRequest {
	t.Helper()
	ch := make(chan *playwright.NetworkRequest, 1)
	var once sync.Once
	cancel := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
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
		t.Fatalf("captureFinishedRequestTo: timed out waiting for request to URL containing %q", urlSubstr)
		return nil
	}
}

// TestNetworkSizesRequestBodyAndHeaders verifies request body and header sizes for a POST request.
// Ref: TestPageNetworkSizes.java#shouldSetBodySizeAndHeadersSize
func TestNetworkSizesRequestBodyAndHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/upload", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureFinishedRequestTo(t, page, "/upload", func() {
		_, err := page.Evaluate(ctx, `({url}) => fetch(url, {
			method: 'POST',
			body: 'hello',
			headers: {'Content-Type': 'text/plain'},
		})`, map[string]any{"url": srv.Prefix() + "/upload"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	sizes, err := req.Sizes(ctx)
	must.NoError(err, "Sizes() failed")
	must.NotNil(sizes, "Sizes() returned nil")

	is.Greater(sizes.RequestHeadersSize, 0, "requestHeadersSize should be > 0")
	is.Greater(sizes.RequestBodySize, 0, "requestBodySize should be > 0 for POST with body")
}

// TestNetworkSizesRequestBodyIsZeroWhenNoBody verifies body size is 0 when no body is sent.
// Ref: TestPageNetworkSizes.java#shouldSetBodySizeTo0IfThereWasNoBody
func TestNetworkSizesRequestBodyIsZeroWhenNoBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/get", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureFinishedRequestTo(t, page, "/get", func() {
		_, err := page.Evaluate(ctx, `({url}) => fetch(url)`,
			map[string]any{"url": srv.Prefix() + "/get"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	sizes, err := req.Sizes(ctx)
	must.NoError(err, "Sizes() failed")
	must.NotNil(sizes, "Sizes() returned nil")

	is.Equal(0, sizes.RequestBodySize, "requestBodySize should be 0 for GET request")
	is.Greater(sizes.RequestHeadersSize, 0, "requestHeadersSize should be > 0")
}

// TestNetworkSizesResponseBodyHeadersAndTransfer verifies response body and header sizes.
// Ref: TestPageNetworkSizes.java#shouldSetBodySizeHeadersSizeAndTransferSize
func TestNetworkSizesResponseBodyHeadersAndTransfer(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/data", "application/json", `{"key":"value"}`)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureFinishedRequestTo(t, page, "/data", func() {
		_, err := page.Evaluate(ctx, `({url}) => fetch(url)`,
			map[string]any{"url": srv.Prefix() + "/data"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	sizes, err := req.Sizes(ctx)
	must.NoError(err, "Sizes() failed")
	must.NotNil(sizes, "Sizes() returned nil")

	is.Greater(sizes.ResponseHeadersSize, 0, "responseHeadersSize should be > 0")
	is.Greater(sizes.ResponseBodySize, 0, "responseBodySize should be > 0")
}

// TestNetworkSizesResponseBodyIsZeroWhenNoBody verifies response body size is 0 for 204.
// Ref: TestPageNetworkSizes.java#shouldSetBodySizeTo0WhenThereWasNoResponseBody
func TestNetworkSizesResponseBodyIsZeroWhenNoBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/no-body", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureFinishedRequestTo(t, page, "/no-body", func() {
		_, err := page.Evaluate(ctx, `({url}) => fetch(url)`,
			map[string]any{"url": srv.Prefix() + "/no-body"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	sizes, err := req.Sizes(ctx)
	must.NoError(err, "Sizes() failed")
	must.NotNil(sizes, "Sizes() returned nil")

	is.Equal(0, sizes.ResponseBodySize, "responseBodySize should be 0 for 204 No Content")
}

// TestNetworkSizesCorrectResponseBodySize verifies the response body size matches the actual content length.
// Ref: TestPageNetworkSizes.java#shouldHaveTheCorrectResponseBodySize
func TestNetworkSizesCorrectResponseBodySize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	body := "hello playwright"
	srv.ServeWithBody("/file", "text/plain", body)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureFinishedRequestTo(t, page, "/file", func() {
		_, err := page.Evaluate(ctx, `({url}) => fetch(url)`,
			map[string]any{"url": srv.Prefix() + "/file"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	sizes, err := req.Sizes(ctx)
	must.NoError(err, "Sizes() failed")
	must.NotNil(sizes, "Sizes() returned nil")

	is.Equal(len(body), sizes.ResponseBodySize, "responseBodySize should equal the content length")
}
