//go:build e2e

// E2E tests for: TestGlobalFetch.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"encoding/base64"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextGet(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/api-get", "application/json", `{"ok":true}`)

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/api-get")
	must.NoError(err, "Get failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	is.True(resp.OK())
	body, err := resp.Text(ctx)
	must.NoError(err, "Text() failed")
	is.Equal(`{"ok":true}`, body)
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextPost(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/api-post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"created":true}`))
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/api-post", &playwright.APIFetchOptions{
		Data: []byte(`{"name":"test"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	must.NoError(err, "Post failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(201, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextFetchStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/not-found", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte("not found"))
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Fetch(ctx, srv.Prefix()+"/not-found")
	must.NoError(err, "Fetch failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(404, resp.Status())
	is.False(resp.OK())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIResponseBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/api-body", "text/plain", "body content")

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/api-body")
	must.NoError(err, "Get failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	body, err := resp.Body(ctx)
	must.NoError(err, "APIResponse.Body() failed")
	is.Equal("body content", string(body))
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIResponseHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/api-hdr", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Api-Test", "api-header-value")
		w.WriteHeader(200)
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/api-hdr")
	must.NoError(err, "Get failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	// Playwright normalises response header names to lowercase.
	is.Equal("api-header-value", resp.Headers()["x-api-test"])
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIResponseJSON(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/api-json", "application/json", `{"key":"value","num":7}`)

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/api-json")
	must.NoError(err, "Get failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	jsonVal, err := resp.JSON(ctx)
	must.NoError(err, "APIResponse.JSON() failed")
	m, ok := jsonVal.(map[string]any)
	is.True(ok, "JSON() returned %T, want map[string]any", jsonVal)
	is.Equal("value", m["key"])
	is.Equal(float64(7), m["num"])
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIResponseURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/api-url", "text/plain", "ok")
	target := srv.Prefix() + "/api-url"

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, target)
	must.NoError(err, "Get failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(target, resp.URL())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextWithBaseURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/health", "text/plain", "ok")

	baseURL := srv.Prefix()
	apiCtx, err := globalPW.NewAPIRequestContext(ctx, &playwright.APIRequestContextOptions{
		BaseURL: &baseURL,
	})
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/health")
	must.NoError(err, "Get /health failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextPut(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/api-put", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(405)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"updated":true}`))
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err)
	defer apiCtx.Dispose(ctx)

	method := "PUT"
	resp, err := apiCtx.Fetch(ctx, srv.Prefix()+"/api-put", &playwright.APIFetchOptions{
		Method: &method,
		Data:   []byte(`{"id":1}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx)

	is.Equal(200, resp.Status())
	is.True(resp.OK())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextDelete(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/api-del", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(405)
			return
		}
		w.WriteHeader(204)
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err)
	defer apiCtx.Dispose(ctx)

	method := "DELETE"
	resp, err := apiCtx.Fetch(ctx, srv.Prefix()+"/api-del", &playwright.APIFetchOptions{
		Method: &method,
	})
	must.NoError(err)
	defer resp.Dispose(ctx)

	is.Equal(204, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextPatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/api-patch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			w.WriteHeader(405)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"patched":true}`))
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err)
	defer apiCtx.Dispose(ctx)

	method := "PATCH"
	resp, err := apiCtx.Fetch(ctx, srv.Prefix()+"/api-patch", &playwright.APIFetchOptions{
		Method:  &method,
		Data:    []byte(`{"field":"value"}`),
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	must.NoError(err)
	defer resp.Dispose(ctx)

	is.Equal(200, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestAPIRequestContextCustomHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedHeader string
	srv.SetRoute("/api-custom-hdr", func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Custom-Header")
		w.WriteHeader(200)
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err)
	defer apiCtx.Dispose(ctx)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/api-custom-hdr", &playwright.APIFetchOptions{
		Headers: map[string]string{"X-Custom-Header": "custom-value"},
	})
	must.NoError(err)
	defer resp.Dispose(ctx)

	is.Equal(200, resp.Status())
	is.Equal("custom-value", receivedHeader)
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchHeadShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Head(ctx, srv.EmptyPage())
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	is.Equal(srv.EmptyPage(), resp.URL())
	is.True(resp.OK())
	// HEAD responses have no body
	body, err := resp.Body(ctx)
	must.NoError(err)
	is.Empty(body)
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldPropagateExtraHttpHeadersWithRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var h1, h2, h3 string
	srv.SetRoute("/a/redirect1", func(w http.ResponseWriter, r *http.Request) {
		h1 = r.Header.Get("My-Secret")
		http.Redirect(w, r, "/b/c/redirect2", http.StatusFound)
	})
	srv.SetRoute("/b/c/redirect2", func(w http.ResponseWriter, r *http.Request) {
		h2 = r.Header.Get("My-Secret")
		http.Redirect(w, r, "/simple.json", http.StatusFound)
	})
	srv.SetRoute("/simple.json", func(w http.ResponseWriter, r *http.Request) {
		h3 = r.Header.Get("My-Secret")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"foo":"bar"}`)
	})

	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		ExtraHTTPHeaders: map[string]string{"My-Secret": "Value"},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/a/redirect1")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("Value", h1, "header on first hop")
	is.Equal("Value", h2, "header on second hop")
	is.Equal("Value", h3, "header on final hop")
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldResolveUrlRelativeToGlobalBaseURLOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		BaseURL: &[]string{srv.Prefix()}[0],
	})

	resp, err := apiCtx.Get(ctx, "/empty.html")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(srv.EmptyPage(), resp.URL())
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnEmptyBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	body, err := resp.Body(ctx)
	must.NoError(err)
	// The empty HTML page has some content; what matters is body() doesn't error
	text, err := resp.Text(ctx)
	must.NoError(err)
	_ = body
	_ = text

	// After explicit dispose, body() must fail
	must.NoError(resp.Dispose(ctx))
	_, err = resp.Body(ctx)
	is.Error(err, "body() after dispose must fail")
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldRetryECONNRESET(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var count int32
	srv.SetRoute("/retry", func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&count, 1)
		if n <= 3 {
			// Drop the connection to simulate ECONNRESET
			if hj, ok := w.(http.Hijacker); ok {
				conn, _, _ := hj.Hijack()
				conn.(*net.TCPConn).SetLinger(0)
				conn.Close()
			}
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		io.WriteString(w, "Hello!")
	})

	apiCtx := newAPICtx(t)

	maxRetries := 3
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/retry", &playwright.APIFetchOptions{
		MaxRetries: &maxRetries,
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("Hello!", text)
	is.EqualValues(4, atomic.LoadInt32(&count))
}

// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportDisposeReason(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)

	// Dispose with a reason
	err := apiCtx.Dispose(ctx, "Test ended.")
	must.NoError(err)

	// Subsequent requests should fail and include the reason
	_, err = apiCtx.Get(ctx, srv.EmptyPage())
	is.Error(err)
	is.Contains(err.Error(), "Test ended.", "error should include dispose reason")
}

// TestGlobalFetchShouldRemoveContentLengthFromRedirectedPostRequests verifies that after
// a POST→redirect the follow-up request is a GET with no Content-Length.
// Extra: no equivalent method in playwright-java; related class: TestGlobalFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldRemoveContentLengthFromRedirectedPostRequests(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var firstMethod string
	var finalMethod, finalContentLength string

	srv.SetRoute("/redirect", func(w http.ResponseWriter, r *http.Request) {
		firstMethod = r.Method
		http.Redirect(w, r, "/final", http.StatusFound)
	})
	srv.SetRoute("/final", func(w http.ResponseWriter, r *http.Request) {
		finalMethod = r.Method
		finalContentLength = r.Header.Get("Content-Length")
		w.WriteHeader(200)
	})

	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/redirect", &playwright.APIFetchOptions{
		Data: []byte("hello"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("POST", firstMethod)
	is.Equal("GET", finalMethod)
	is.Empty(finalContentLength, "Content-Length should be absent on redirected GET")
}

// Ref: TestGlobalFetch.java#shouldDisposeGlobalRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldDisposeGlobalRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/data.json", "application/json", `{"foo":"bar"}`)

	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/data.json")
	must.NoError(err)

	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal(`{"foo":"bar"}`, text)

	// Dispose the context; subsequent body() on the response must fail.
	err = apiCtx.Dispose(ctx)
	must.NoError(err)

	_, err = resp.Body(ctx)
	is.Error(err, "body() after context dispose must fail")
}

// Ref: TestGlobalFetch.java#shouldSupportGlobalUserAgentOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportGlobalUserAgentOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedUA string
	srv.SetRoute("/ua", func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(200)
	})

	ua := "My Custom Agent"
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{UserAgent: &ua})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/ua")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("My Custom Agent", receivedUA)
}

// Ref: TestGlobalFetch.java#shouldSupportGlobalTimeoutOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportGlobalTimeoutOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.BlockRequest("/slow")

	timeout := 100.0
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{Timeout: &timeout})

	_, err := apiCtx.Get(ctx, srv.Prefix()+"/slow")
	is.Error(err, "expected timeout error")
	is.True(
		strings.Contains(err.Error(), "Timeout") || strings.Contains(err.Error(), "timeout"),
		"error should mention timeout, got: %v", err)
}

// Ref: TestGlobalFetch.java#shouldSupportGlobalHttpCredentialsOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportGlobalHttpCredentialsOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "secret")

	// Without credentials → 401
	apiCtx1 := newAPICtx(t)
	resp1, err := apiCtx1.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp1.Dispose(ctx) //nolint:errcheck
	is.Equal(401, resp1.Status())

	// With credentials → 200
	apiCtx2 := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{Username: "user", Password: "pass"},
	})
	resp2, err := apiCtx2.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp2.Dispose(ctx) //nolint:errcheck
	is.Equal(200, resp2.Status())
}

// Ref: TestGlobalFetch.java#shouldReturnServerAddressFromResponse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnServerAddressFromResponse(t *testing.T) {
	t.Skip("serverAddr not returned in APIResponse by driver v1.61.1")
}

// Ref: TestGlobalFetch.java#shouldReturnNullSecurityDetailsForHttpResponse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnNullSecurityDetailsForHttpResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Nil(resp.SecurityDetails(), "securityDetails should be nil for plain HTTP")
}

// Ref: TestGlobalFetch.java#shouldReturnBodyForFailingRequests  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnBodyForFailingRequests(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/not-found", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		io.WriteString(w, "Not found.")
	})

	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/not-found")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(404, resp.Status())
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("Not found.", text)
}

// Ref: TestGlobalFetch.java#shouldJsonStringifyTypeBodyWhenContentTypeIsApplicationJson  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldJsonStringifyTypeBodyWhenContentTypeIsApplicationJson(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedBody []byte
	srv.SetRoute("/json", func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	})

	apiCtx := newAPICtx(t)

	payload := []byte(`{"name":"John","age":30}`)
	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/json", &playwright.APIFetchOptions{
		Headers: map[string]string{"Content-Type": "application/json"},
		Data:    payload,
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(string(payload), string(receivedBody))
}

// Ref: TestGlobalFetch.java#shouldThrowAnErrorWhenMaxRedirectsIsExceeded  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldThrowAnErrorWhenMaxRedirectsIsExceeded(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/r1", "/r2")
	srv.ServeWithRedirect("/r2", "/r3")
	srv.ServeWithRedirect("/r3", "/r4")
	srv.ServeWithBody("/r4", "text/plain", "final")

	apiCtx := newAPICtx(t)

	maxRedirects := 2
	_, err := apiCtx.Get(ctx, srv.Prefix()+"/r1", &playwright.APIFetchOptions{
		MaxRedirects: &maxRedirects,
	})
	is.Error(err, "should error when max redirects exceeded")
	is.True(
		strings.Contains(err.Error(), "redirect") || strings.Contains(err.Error(), "Redirect"),
		"error should mention redirects, got: %v", err)
}

// Ref: TestGlobalFetch.java#shouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/start", "/final")
	srv.ServeWithBody("/final", "text/plain", "reached")

	apiCtx := newAPICtx(t)

	maxRedirects := 0
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/start", &playwright.APIFetchOptions{
		MaxRedirects: &maxRedirects,
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	// Should get the 302 redirect response, not the final destination
	is.Equal(302, resp.Status())
	loc := resp.Headers()["location"]
	if loc == "" {
		loc = resp.Headers()["Location"]
	}
	is.Contains(loc, "/final")
}

// Ref: TestGlobalFetch.java#shouldThrowAnErrorWhenMaxRedirectsIsLessThan0  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldThrowAnErrorWhenMaxRedirectsIsLessThan0(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ok", "text/plain", "ok")

	apiCtx := newAPICtx(t)

	maxRedirects := -1
	_, err := apiCtx.Get(ctx, srv.Prefix()+"/ok", &playwright.APIFetchOptions{
		MaxRedirects: &maxRedirects,
	})
	is.Error(err, "negative maxRedirects should error")
}

// Ref: TestGlobalFetch.java#shouldThrowWhenFailOnStatusCodeIsSetToTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldThrowWhenFailOnStatusCodeIsSetToTrue(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/404", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		io.WriteString(w, "Not found.")
	})

	failOnStatusCode := true
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		FailOnStatusCode: &failOnStatusCode,
	})

	_, err := apiCtx.Get(ctx, srv.Prefix()+"/404")
	is.Error(err, "should error on 404 when failOnStatusCode=true")
	is.True(
		strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found"),
		"error should mention 404, got: %v", err)
}

// Ref: TestGlobalFetch.java#shouldSupportGlobalHttpCredentialsOptionAndMatchingOrigin  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportGlobalHttpCredentialsOptionAndMatchingOrigin(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	origin := srv.Prefix()
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &origin,
		},
	})

	// Same origin → credentials sent → 200
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, resp.Status())

	// Cross-origin → credentials NOT sent → 401
	crossOrigin := srv.CrossProcessPrefix()
	if crossOrigin != srv.Prefix() {
		srv2 := testserver.New(t)
		srv2.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")
		resp2, err := apiCtx.Get(ctx, srv2.CrossProcessPrefix()+"/protected")
		if err == nil {
			defer resp2.Dispose(ctx) //nolint:errcheck
			is.Equal(401, resp2.Status())
		}
	}
}

// TestGlobalFetchShouldReturnErrorWithWrongCredentials verifies that wrong basic-auth
// credentials result in a 401 response (not an error).
// Ref: TestGlobalFetch.java#shouldReturnErrorWithWrongCredentials  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnErrorWithWrongCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{Username: "user", Password: "wrong"},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status())
}

// TestGlobalFetchShouldSetPlaywrightAsUserAgent verifies the default User-Agent
// includes "playwright" (case-insensitive).
// Ref: TestGlobalFetch.java#shouldSetPlaywrightAsUserAgent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSetPlaywrightAsUserAgent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/ua")

	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/ua")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	req := <-ch
	ua := req.Header.Get("User-Agent")
	is.True(
		strings.Contains(strings.ToLower(ua), "playwright"),
		"User-Agent should contain 'playwright', got: %s", ua)
}

// TestGlobalFetchShouldNotDoubleStringifyTypeBodyWhenContentTypeIsApplicationJson verifies
// that when Content-Type is application/json the body bytes are sent as-is (not re-encoded).
// Ref: TestGlobalFetch.java#shouldNotDoubleStringifyTypeBodyWhenContentTypeIsApplicationJson  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldNotDoubleStringifyTypeBodyWhenContentTypeIsApplicationJson(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedBody []byte
	srv.SetRoute("/json", func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	})

	apiCtx := newAPICtx(t)

	payload := `{"foo":"bar"}`
	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/json", &playwright.APIFetchOptions{
		Headers: map[string]string{"Content-Type": "application/json"},
		Data:    []byte(payload),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(payload, string(receivedBody))
}

// TestGlobalFetchShouldUseMaxRedirectsFromFetchWhenProvidedOverridingNewContext verifies that
// a per-request maxRedirects overrides the context-level maxRedirects.
// Ref: TestGlobalFetch.java#shouldUseMaxRedirectsFromFetchWhenProvidedOverridingNewContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldUseMaxRedirectsFromFetchWhenProvidedOverridingNewContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/r1", "/r2")
	srv.ServeWithRedirect("/r2", "/r3")
	srv.ServeWithRedirect("/r3", "/r4")
	srv.ServeWithBody("/r4", "text/plain", "final")

	contextMax := 1
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{MaxRedirects: &contextMax})

	fetchMax := 10
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/r1", &playwright.APIFetchOptions{
		MaxRedirects: &fetchMax,
	})
	must.NoError(err, "fetch maxRedirects=10 should override context maxRedirects=1")
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("final", text)
}

// TestGlobalFetchShouldFollowRedirectsUpToMaxRedirectsLimitSetInNewContext verifies that
// context-level maxRedirects is enforced when no per-request override is given.
// Ref: TestGlobalFetch.java#shouldFollowRedirectsUpToMaxRedirectsLimitSetInNewContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldFollowRedirectsUpToMaxRedirectsLimitSetInNewContext(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/r1", "/r2")
	srv.ServeWithRedirect("/r2", "/r3")
	srv.ServeWithRedirect("/r3", "/r4")
	srv.ServeWithBody("/r4", "text/plain", "final")

	contextMax := 2
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{MaxRedirects: &contextMax})

	_, err := apiCtx.Get(ctx, srv.Prefix()+"/r1")
	is.Error(err, "should error when redirect chain exceeds context maxRedirects")
	is.True(
		strings.Contains(strings.ToLower(err.Error()), "redirect"),
		"error should mention redirects, got: %v", err)
}

// TestGlobalFetchShouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0InNewContext verifies that
// context-level maxRedirects=0 returns the redirect response without following it.
// Ref: TestGlobalFetch.java#shouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0InNewContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0InNewContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/start", "/final")
	srv.ServeWithBody("/final", "text/plain", "reached")

	contextMax := 0
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{MaxRedirects: &contextMax})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/start")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(302, resp.Status())
}

// TestGlobalFetchShouldNotModifyRequestMethodInOptions verifies that calling Get/Post with a
// shared options struct does not mutate the struct's Method field.
// Ref: TestGlobalFetch.java#shouldNotModifyRequestMethodInOptions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldNotModifyRequestMethodInOptions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var methods []string
	srv.SetRoute("/echo", func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		w.WriteHeader(200)
	})

	apiCtx := newAPICtx(t)

	opts := &playwright.APIFetchOptions{
		Data: []byte("body"),
	}

	resp1, err := apiCtx.Get(ctx, srv.Prefix()+"/echo", opts)
	must.NoError(err)
	defer resp1.Dispose(ctx) //nolint:errcheck

	resp2, err := apiCtx.Post(ctx, srv.Prefix()+"/echo", opts)
	must.NoError(err)
	defer resp2.Dispose(ctx) //nolint:errcheck

	is.Nil(opts.Method, "original opts.Method must not be modified by Get or Post")
	is.Equal([]string{"GET", "POST"}, methods)
}

// TestGlobalFetchShouldSupportGlobalHttpCredentialsOptionAndMatchingOriginCaseInsensitive
// verifies that origin matching is case-insensitive.
// Ref: TestGlobalFetch.java#shouldSupportGlobalHttpCredentialsOptionAndMatchingOriginCaseInsensitive  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportGlobalHttpCredentialsOptionAndMatchingOriginCaseInsensitive(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	origin := strings.ToUpper(srv.Prefix())
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &origin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status(), "credentials should be sent with upper-case origin")
}

// TestGlobalFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginScheme verifies that
// credentials are NOT sent when the configured origin has a different scheme.
// Ref: TestGlobalFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginScheme  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginScheme(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	wrongSchemeOrigin := strings.Replace(srv.Prefix(), "http://", "https://", 1)
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &wrongSchemeOrigin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status(), "credentials should NOT be sent when origin scheme differs")
}

// TestGlobalFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginHostname verifies that
// credentials are NOT sent when the configured origin has a different hostname.
// Ref: TestGlobalFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginHostname  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginHostname(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	u, err := url.Parse(srv.Prefix())
	must.NoError(err)
	wrongHostOrigin := fmt.Sprintf("%s://wrong-hostname:%s", u.Scheme, u.Port())
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &wrongHostOrigin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status(), "credentials should NOT be sent when origin hostname differs")
}

// TestGlobalFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginPort verifies that
// credentials are NOT sent when the configured origin has a different port.
// Ref: TestGlobalFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginPort  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginPort(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	u, err := url.Parse(srv.Prefix())
	must.NoError(err)
	wrongPort := "1234"
	if u.Port() == "1234" {
		wrongPort = "5678"
	}
	wrongPortOrigin := fmt.Sprintf("%s://%s:%s", u.Scheme, u.Hostname(), wrongPort)
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &wrongPortOrigin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status(), "credentials should NOT be sent when origin port differs")
}

// TestGlobalFetchShouldSupportHTTPCredentialsSendImmediately verifies that Send="always"
// causes credentials to be included on the very first request without a prior 401 challenge.
// Ref: TestGlobalFetch.java#shouldSupportHTTPCredentialsSendImmediately  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportHTTPCredentialsSendImmediately(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/direct-auth")

	sendAlways := "always"
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Send:     &sendAlways,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/direct-auth")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	req := <-ch
	authHeader := req.Header.Get("Authorization")
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	is.Equal(expected, authHeader, "Authorization header should be sent on first request")
}

// TestGlobalFetchShouldNotThrowWhenFailOnStatusCodeIsSetToFalse verifies that a 4xx response
// does not produce an error when failOnStatusCode is explicitly false.
// Ref: TestGlobalFetch.java#shouldNotThrowWhenFailOnStatusCodeIsSetToFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldNotThrowWhenFailOnStatusCodeIsSetToFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/404", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.WriteString(w, "Not found") //nolint:errcheck
	})

	failOnStatusCode := false
	apiCtx := newAPICtx(t, &playwright.APIRequestContextOptions{
		FailOnStatusCode: &failOnStatusCode,
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/404")
	must.NoError(err, "should not throw when failOnStatusCode=false")
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(404, resp.Status())
}

// TestGlobalFetchFetchShouldWork verifies that APIRequestContext.Fetch returns a valid JSON response.
// Ref: TestGlobalFetch.java#fetchShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchFetchShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)
	url := srv.Prefix() + "/simple.json"
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Fetch(ctx, url)
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(url, resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Contains(resp.Headers()["content-type"], "application/json")
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal(simpleJSONBody, text)
}

// TestGlobalFetchDeleteShouldWork verifies that APIRequestContext.Delete returns a valid JSON response.
// Ref: TestGlobalFetch.java#deleteShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchDeleteShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)
	url := srv.Prefix() + "/simple.json"
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Delete(ctx, url)
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(url, resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Contains(resp.Headers()["content-type"], "application/json")
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal(simpleJSONBody, text)
}

// TestGlobalFetchGetShouldWork verifies that APIRequestContext.Get returns a valid JSON response.
// Ref: TestGlobalFetch.java#getShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchGetShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)
	url := srv.Prefix() + "/simple.json"
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, url)
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(url, resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Contains(resp.Headers()["content-type"], "application/json")
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal(simpleJSONBody, text)
}

// TestGlobalFetchPatchShouldWork verifies that APIRequestContext.Patch returns a valid JSON response.
// Ref: TestGlobalFetch.java#patchShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchPatchShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)
	url := srv.Prefix() + "/simple.json"
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Patch(ctx, url)
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(url, resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Contains(resp.Headers()["content-type"], "application/json")
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal(simpleJSONBody, text)
}

// TestGlobalFetchPostShouldWork verifies that APIRequestContext.Post returns a valid JSON response.
// Ref: TestGlobalFetch.java#postShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchPostShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)
	url := srv.Prefix() + "/simple.json"
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Post(ctx, url)
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(url, resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Contains(resp.Headers()["content-type"], "application/json")
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal(simpleJSONBody, text)
}

// TestGlobalFetchPutShouldWork verifies that APIRequestContext.Put returns a valid JSON response.
// Ref: TestGlobalFetch.java#putShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchPutShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)
	url := srv.Prefix() + "/simple.json"
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Put(ctx, url)
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(url, resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Contains(resp.Headers()["content-type"], "application/json")
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal(simpleJSONBody, text)
}

// Ref: TestGlobalFetch.java#clearCookiesShouldRemoveAllCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchClearCookiesShouldRemoveAllCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	url := srv.EmptyPage()
	must.NoError(bc.AddCookies(ctx, []playwright.Cookie{{
		Name:  "todelete",
		Value: "yes",
		URL:   &url,
	}}))

	must.NoError(bc.ClearCookies(ctx))

	cookies, err := bc.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)
}

// TestGlobalFetchShouldHaveJavaInDefaultUesrAgent verifies the default user agent of the global request context.
// Ref: TestGlobalFetch.java#shouldHaveJavaInDefaultUesrAgent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldHaveJavaInDefaultUesrAgent(t *testing.T) {
	t.Skip("not implemented: Java-specific default user agent (java/<version>) has no Go equivalent")
}

// TestGlobalFetchHeadShouldWorkA4 verifies APIRequestContext.Head.
// Ref: TestGlobalFetch.java#headShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchHeadShouldWorkA4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/empty.html", "text/html", "")
	apiCtx := a4NewAPICtx(t)

	resp, err := apiCtx.Head(ctx, srv.EmptyPage())
	must.NoError(err)
	is.Equal(srv.EmptyPage(), resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Equal("text/html", resp.Headers()["content-type"])
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("", text)
}

// TestGlobalFetchAddCookiesShouldAddCookiesToTheCookieJar verifies APIRequestContext.addCookies.
// Ref: TestGlobalFetch.java#addCookiesShouldAddCookiesToTheCookieJar  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchAddCookiesShouldAddCookiesToTheCookieJar(t *testing.T) {
	t.Skip("APIRequestContext.addCookies not supported by driver v1.61.1")
}

// TestGlobalFetchAddCookiesShouldValidateCookies verifies cookie validation in APIRequestContext.addCookies.
// Ref: TestGlobalFetch.java#addCookiesShouldValidateCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchAddCookiesShouldValidateCookies(t *testing.T) {
	t.Skip("APIRequestContext.addCookies not supported by driver v1.61.1")
}

// TestGlobalFetchCookiesShouldReturnCookiesFilteredByUrls verifies APIRequestContext.cookies(urls).
// Ref: TestGlobalFetch.java#cookiesShouldReturnCookiesFilteredByUrls  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchCookiesShouldReturnCookiesFilteredByUrls(t *testing.T) {
	t.Skip("APIRequestContext.cookies not supported by driver v1.61.1")
}

// TestGlobalFetchClearCookiesShouldFilterByNameDomainAndPath verifies APIRequestContext.clearCookies filters.
// Ref: TestGlobalFetch.java#clearCookiesShouldFilterByNameDomainAndPath  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchClearCookiesShouldFilterByNameDomainAndPath(t *testing.T) {
	t.Skip("APIRequestContext.clearCookies not supported by driver v1.61.1")
}

// TestGlobalFetchShouldRemoveContentLengthFromReidrectedPostRequests verifies content-length is dropped on redirected POSTs.
// Ref: TestGlobalFetch.java#shouldRemoveContentLengthFromReidrectedPostRequests  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldRemoveContentLengthFromReidrectedPostRequests(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	var mu sync.Mutex
	var cl1 string
	var cl2 []string
	srv.SetRoute("/redirect", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		cl1 = r.Header.Get("Content-Length")
		mu.Unlock()
		http.Redirect(w, r, "/empty.html", http.StatusFound)
	})
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		cl2 = r.Header.Values("Content-Length")
		mu.Unlock()
		w.Header().Set("Content-Type", "text/html")
	})
	apiCtx := a4NewAPICtx(t)

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/redirect", &playwright.APIFetchOptions{Data: []byte(`{"foo":"bar"}`)})
	must.NoError(err)
	is.Equal(200, resp.Status())
	mu.Lock()
	defer mu.Unlock()
	is.Equal("13", cl1)
	is.Empty(cl2)
}

// TestGlobalFetchShouldSupportMultipleHttpCredentials verifies a list of http credentials in a global request context.
// Ref: TestGlobalFetch.java#shouldSupportMultipleHttpCredentials  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportMultipleHttpCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "user1" || password != "pass1" {
			w.Header().Set("WWW-Authenticate", `Basic realm="test"`)
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
	prefix := srv.Prefix()
	crossPrefix := srv.CrossProcessPrefix()
	apiCtx := a4NewAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentialsList: []playwright.HttpCredentials{
			{Username: "user1", Password: "pass1", Origin: &prefix},
			{Username: "user2", Password: "pass2", Origin: &crossPrefix},
		},
	})
	response, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	is.Equal(200, response.Status())
	response, err = apiCtx.Get(ctx, crossPrefix+"/empty.html")
	must.NoError(err)
	is.Equal(401, response.Status())
}

// TestGlobalFetchShouldSupportHTTPCredentialsSendWithMultipleHttpCredentials verifies per-origin "send" of http credentials.
// Ref: TestGlobalFetch.java#shouldSupportHTTPCredentialsSendWithMultipleHttpCredentials  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldSupportHTTPCredentialsSendWithMultipleHttpCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	authorization := make(chan string, 2)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		authorization <- r.Header.Get("Authorization")
	})
	prefix := srv.Prefix()
	crossPrefix := srv.CrossProcessPrefix()
	always := "always"
	unauthorized := "unauthorized"
	apiCtx := a4NewAPICtx(t, &playwright.APIRequestContextOptions{
		HttpCredentialsList: []playwright.HttpCredentials{
			{Username: "user1", Password: "pass1", Origin: &prefix, Send: &always},
			{Username: "user2", Password: "pass2", Origin: &crossPrefix, Send: &unauthorized},
		},
	})

	response, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	is.Equal(200, response.Status())
	is.Equal("Basic "+base64.StdEncoding.EncodeToString([]byte("user1:pass1")), <-authorization)

	response, err = apiCtx.Get(ctx, crossPrefix+"/empty.html")
	must.NoError(err)
	is.Equal(200, response.Status())
	is.Empty(<-authorization)
}

// TestGlobalFetchShouldResolveUrlRelativeToGobalBaseURLOption verifies a relative URL resolves against the global base URL.
// Ref: TestGlobalFetch.java#shouldResolveUrlRelativeToGobalBaseURLOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGlobalFetch.java
func TestGlobalFetchShouldResolveUrlRelativeToGobalBaseURLOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	request := a4NewAPICtx(t, &playwright.APIRequestContextOptions{BaseURL: a4Ptr(srv.Prefix())})
	resp, err := request.Get(ctx, "/empty.html")
	must.NoError(err)
	must.Equal(srv.EmptyPage(), resp.URL())
}

func sortedCookieValues(cookies []playwright.Cookie) []string {
	values := make([]string, len(cookies))
	for i, cookie := range cookies {
		values[i] = cookie.Value
	}
	sort.Strings(values)
	return values
}
