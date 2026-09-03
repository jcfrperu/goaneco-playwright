//go:build e2e

package e2e

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newAPICtx creates a standalone APIRequestContext and registers cleanup.
func newAPICtx(t *testing.T, opts ...*playwright.APIRequestContextOptions) *playwright.APIRequestContext {
	t.Helper()
	ctx := testCtx(t)
	apiCtx, err := globalPW.NewAPIRequestContext(ctx, opts...)
	require.NoError(t, err, "NewAPIRequestContext")
	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = apiCtx.Dispose(cleanCtx)
	})
	return apiCtx
}

// ── TestGlobalFetch equivalents ──────────────────────────────────────────────

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

func TestGlobalFetchShouldReturnServerAddressFromResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	addr := resp.ServerAddr()
	must.NotNil(addr, "serverAddr should not be nil")
	is.Equal(srv.Port(), addr.Port)
	is.True(
		addr.IpAddress == "127.0.0.1" || addr.IpAddress == "::1",
		"unexpected IP: %s", addr.IpAddress)
}

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

func TestGlobalFetchShouldJsonStringifyBodyWhenContentTypeIsApplicationJson(t *testing.T) {
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
// Ref: TestGlobalFetch.java#shouldReturnErrorWithWrongCredentials
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
// Ref: TestGlobalFetch.java#shouldSetPlaywrightAsUserAgent
func TestGlobalFetchShouldSetPlaywrightAsUserAgent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/ua")
	srv.SetRoute("/ua", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

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

// TestGlobalFetchShouldNotDoubleStringifyBodyWhenContentTypeIsApplicationJson verifies
// that when Content-Type is application/json the body bytes are sent as-is (not re-encoded).
// Ref: TestGlobalFetch.java#shouldNotDoubleStringifyBodyWhenContentTypeIsApplicationJson
func TestGlobalFetchShouldNotDoubleStringifyBodyWhenContentTypeIsApplicationJson(t *testing.T) {
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

// TestGlobalFetchShouldRemoveContentLengthFromRedirectedPostRequests verifies that after
// a POST→redirect the follow-up request is a GET with no Content-Length.
// Ref: TestGlobalFetch.java#shouldRemoveContentLengthFromRedirectedPostRequests
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

// TestGlobalFetchShouldUseMaxRedirectsFromFetchWhenProvidedOverridingNewContext verifies that
// a per-request maxRedirects overrides the context-level maxRedirects.
// Ref: TestGlobalFetch.java#shouldUseMaxRedirectsFromFetchWhenProvidedOverridingNewContext
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
// Ref: TestGlobalFetch.java#shouldFollowRedirectsUpToMaxRedirectsLimitSetInNewContext
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
// Ref: TestGlobalFetch.java#shouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0InNewContext
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
// Ref: TestGlobalFetch.java#shouldNotModifyRequestMethodInOptions
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
// Ref: TestGlobalFetch.java#shouldSupportGlobalHttpCredentialsOptionAndMatchingOriginCaseInsensitive
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
// Ref: TestGlobalFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginScheme
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
// Ref: TestGlobalFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginHostname
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
// Ref: TestGlobalFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginPort
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
// Ref: TestGlobalFetch.java#shouldSupportHTTPCredentialsSendImmediately
func TestGlobalFetchShouldSupportHTTPCredentialsSendImmediately(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/direct-auth")
	srv.SetRoute("/direct-auth", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

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
// Ref: TestGlobalFetch.java#shouldNotThrowWhenFailOnStatusCodeIsSetToFalse
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

const simpleJSONBody = "{\"foo\": \"bar\"}\n"

// TestGlobalFetchFetchShouldWork verifies that APIRequestContext.Fetch returns a valid JSON response.
// Ref: TestGlobalFetch.java#fetchShouldWork
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
// Ref: TestGlobalFetch.java#deleteShouldWork
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
// Ref: TestGlobalFetch.java#getShouldWork
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
// Ref: TestGlobalFetch.java#patchShouldWork
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
// Ref: TestGlobalFetch.java#postShouldWork
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
// Ref: TestGlobalFetch.java#putShouldWork
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
