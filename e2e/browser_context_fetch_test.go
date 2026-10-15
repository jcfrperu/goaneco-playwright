//go:build e2e

// E2E tests for: TestBrowserContextFetch.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
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
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchGetShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/simple.json")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(srv.Prefix()+"/simple.json", resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())

	ct := resp.Headers()["content-type"]
	is.Contains(ct, "application/json")

	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchGetShouldReturnTiming(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/simple.json")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.True(resp.OK())

	timing := resp.Timing()
	must.NotNil(timing)
	is.Greater(timing.StartTime, 0.0, "startTime should be > 0")
	is.GreaterOrEqual(timing.DomainLookupEnd, timing.DomainLookupStart)
	is.GreaterOrEqual(timing.ConnectStart, timing.DomainLookupEnd)
	is.Equal(-1.0, timing.SecureConnectionStart)
	is.GreaterOrEqual(timing.ConnectEnd, timing.ConnectStart)
	is.GreaterOrEqual(timing.RequestStart, timing.ConnectEnd)
	is.GreaterOrEqual(timing.ResponseStart, timing.RequestStart)
	is.GreaterOrEqual(timing.ResponseEnd, timing.ResponseStart)
	is.Less(timing.ResponseEnd, 60_000.0, "responseEnd should be < 60s")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Fetch(ctx, srv.Prefix()+"/simple.json")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(srv.Prefix()+"/simple.json", resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())

	ct := resp.Headers()["content-type"]
	is.Contains(ct, "application/json")

	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldThrowOnNetworkError(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Close the connection without sending a response to simulate socket hang up.
	srv.SetRoute("/test", func(w http.ResponseWriter, r *http.Request) {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
	})

	_, apiCtx := newBCtxWithRequest(t)

	_, err := apiCtx.Get(ctx, srv.Prefix()+"/test")
	is.Error(err, "expected network error")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldAddSessionCookiesToRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedCookie string
	srv.SetRoute("/simple.json", func(w http.ResponseWriter, r *http.Request) {
		receivedCookie = r.Header.Get("Cookie")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"foo\": \"bar\"}\n")
	})

	bCtx, apiCtx := newBCtxWithRequest(t)

	domain := "localhost"
	path := "/"
	expires := -1.0
	httpOnly := false
	secure := false
	sameSite := playwright.SameSiteLax
	err := bCtx.AddCookies(ctx, []playwright.Cookie{
		{
			Name:     "username",
			Value:    "John Doe",
			Domain:   &domain,
			Path:     &path,
			Expires:  &expires,
			HTTPOnly: &httpOnly,
			Secure:   &secure,
			SameSite: sameSite,
		},
	})
	must.NoError(err)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/simple.json")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Contains(receivedCookie, "username=John Doe")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchGetShouldSupportQueryParams(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/empty.html")

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.EmptyPage()+"?p1=foo", &playwright.APIFetchOptions{
		Params: map[string]string{
			"p1":     "v1",
			"param2": "value2",
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	req := <-ch
	rawQuery := req.URL.RawQuery
	is.Contains(rawQuery, "p1=foo")
	is.Contains(rawQuery, "p1=v1")
	is.Contains(rawQuery, "param2=value2")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchGetShouldSupportFailOnStatusCode(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/does-not-exist.html", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.WriteString(w, "Not Found")
	})

	_, apiCtx := newBCtxWithRequest(t)

	failOnStatus := true
	_, err := apiCtx.Get(ctx, srv.Prefix()+"/does-not-exist.html", &playwright.APIFetchOptions{
		FailOnStatusCode: &failOnStatus,
	})
	is.Error(err, "should error on 404 with failOnStatusCode")
	is.True(
		strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found"),
		"error should mention 404, got: %v", err)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldFollowRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)
	srv.ServeWithRedirect("/redirect1", "/redirect2")
	srv.ServeWithRedirect("/redirect2", "/simple.json")

	var receivedCookie string
	origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/simple.json" {
			receivedCookie = r.Header.Get("Cookie")
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, "{\"foo\": \"bar\"}\n")
		}
	})
	srv.SetRoute("/simple.json", origHandler)

	bCtx, apiCtx := newBCtxWithRequest(t)

	domain := "localhost"
	path := "/"
	expires := -1.0
	httpOnly := false
	secure := false
	sameSite := playwright.SameSiteLax
	err := bCtx.AddCookies(ctx, []playwright.Cookie{
		{
			Name:     "username",
			Value:    "John Doe",
			Domain:   &domain,
			Path:     &path,
			Expires:  &expires,
			HTTPOnly: &httpOnly,
			Secure:   &secure,
			SameSite: sameSite,
		},
	})
	must.NoError(err)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/redirect1")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Contains(receivedCookie, "username=John Doe")
	is.Equal(srv.Prefix()+"/simple.json", resp.URL())
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldWorkWithHttpCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/empty.html", "user", "pass", "text/html", "OK")

	_, apiCtx := newBCtxWithRequest(t)

	// Without credentials → 401.
	resp, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck
	is.Equal(401, resp.Status())

	// With credentials (manual Authorization header) → 200.
	creds := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	resp2, err := apiCtx.Get(ctx, srv.EmptyPage(), &playwright.APIFetchOptions{
		Headers: map[string]string{"Authorization": "Basic " + creds},
	})
	must.NoError(err)
	defer resp2.Dispose(ctx) //nolint:errcheck
	is.Equal(200, resp2.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchPostShouldSupportPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)

	var receivedMethod, receivedBody string
	srv.SetRoute("/simple.json", func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		receivedBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"foo\": \"bar\"}\n")
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/simple.json", &playwright.APIFetchOptions{
		Data: []byte("My request"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("POST", receivedMethod)
	is.Equal("My request", receivedBody)
	is.Equal(200, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchDeleteShouldSupportPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedMethod, receivedBody string
	srv.SetRoute("/simple.json", func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		receivedBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"foo\": \"bar\"}\n")
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Delete(ctx, srv.Prefix()+"/simple.json", &playwright.APIFetchOptions{
		Data: []byte("My request"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("DELETE", receivedMethod)
	is.Equal("My request", receivedBody)
	is.Equal(200, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchPatchShouldSupportPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedMethod, receivedBody string
	srv.SetRoute("/simple.json", func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		receivedBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"foo\": \"bar\"}\n")
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Patch(ctx, srv.Prefix()+"/simple.json", &playwright.APIFetchOptions{
		Data: []byte("My request"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("PATCH", receivedMethod)
	is.Equal("My request", receivedBody)
	is.Equal(200, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchPutShouldSupportPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedMethod, receivedBody string
	srv.SetRoute("/simple.json", func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		receivedBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"foo\": \"bar\"}\n")
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Put(ctx, srv.Prefix()+"/simple.json", &playwright.APIFetchOptions{
		Data: []byte("My request"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("PUT", receivedMethod)
	is.Equal("My request", receivedBody)
	is.Equal(200, resp.Status())
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSendContentLength(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedContentLength, receivedContentType string
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		receivedContentLength = r.Header.Get("Content-Length")
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	data := make([]byte, 256)
	for i := range data {
		data[i] = byte(i)
	}
	resp, err := apiCtx.Post(ctx, srv.EmptyPage(), &playwright.APIFetchOptions{
		Data: data,
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("256", receivedContentLength)
	is.Contains(receivedContentType, "application/octet-stream")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportTimeoutOption(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.BlockRequest("/slow")

	_, apiCtx := newBCtxWithRequest(t)

	timeout := 100.0
	_, err := apiCtx.Get(ctx, srv.Prefix()+"/slow", &playwright.APIFetchOptions{
		Timeout: &timeout,
	})
	is.Error(err, "expected timeout error")
	is.True(
		strings.Contains(err.Error(), "Timeout") || strings.Contains(err.Error(), "timeout"),
		"error should mention timeout, got: %v", err)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportATimeoutOf0(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/slow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(100 * time.Millisecond)
		io.WriteString(w, "done")
	})

	_, apiCtx := newBCtxWithRequest(t)

	zero := 0.0
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/slow", &playwright.APIFetchOptions{
		Timeout: &zero,
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("done", text)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldDispose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/simple.json")
	must.NoError(err)

	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)

	must.NoError(resp.Dispose(ctx))

	_, err = resp.Body(ctx)
	is.Error(err, "body() after dispose must fail")
	is.True(
		strings.Contains(err.Error(), "Response has been disposed") || strings.Contains(err.Error(), "disposed"),
		"error should mention disposed, got: %v", err)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldDisposeWhenContextCloses(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveSimpleJSON(srv)

	// Create context with manual lifecycle so we can close it mid-test.
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	}()

	apiCtx := bCtx.Request()
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/simple.json")
	must.NoError(err)

	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)

	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = bCtx.Close(closeCtx)
	must.NoError(err)

	_, err = resp.Body(ctx)
	is.Error(err, "body() after context close must fail")
	is.True(
		strings.Contains(err.Error(), "Response has been disposed") ||
			strings.Contains(err.Error(), "closed") ||
			strings.Contains(err.Error(), "disposed"),
		"error should mention disposal or closure, got: %v", err)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldAcceptBoolAndNumericParams(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/empty.html")

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.EmptyPage(), &playwright.APIFetchOptions{
		Params: map[string]string{
			"str":   "s",
			"num":   "10",
			"bool":  "true",
			"bool2": "false",
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	req := <-ch
	rawQuery := req.URL.RawQuery
	is.Contains(rawQuery, "str=s")
	is.Contains(rawQuery, "num=10")
	is.Contains(rawQuery, "bool=true")
	is.Contains(rawQuery, "bool2=false")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldRetryECONNRESET(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var count int32
	srv.SetRoute("/test", func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&count, 1)
		if n <= 3 {
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

	_, apiCtx := newBCtxWithRequest(t)

	maxRetries := 3
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/test", &playwright.APIFetchOptions{
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

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestAPIRequestContextFetchWithQueryParams(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedParam string
	srv.SetRoute("/api-query", func(w http.ResponseWriter, r *http.Request) {
		receivedParam = r.URL.Query().Get("key")
		w.WriteHeader(200)
	})

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err)
	defer apiCtx.Dispose(ctx)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/api-query?key=testval")
	must.NoError(err)
	defer resp.Dispose(ctx)

	is.Equal(200, resp.Status())
	is.Equal("testval", receivedParam)
}

// TestBrowserContextSessionStorageSetGetEx2 verifies sessionStorage set/get in page.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextFetch.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextSessionStorageSetGetEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => sessionStorage.setItem('sessKey', 'sessValue')`)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => sessionStorage.getItem('sessKey')`)
	must.NoError(err)
	is.Equal("sessValue", result)
}

// Ref: TestBrowserContextFetch.java#contextRequestShouldExportSameStorageStateAsContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchContextRequestShouldExportSameStorageStateAsContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/setcookie.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "a=b")
		w.Header().Add("Set-Cookie", "c=d")
		w.WriteHeader(200)
	})

	bCtx, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/setcookie.html")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	contextState, err := bCtx.StorageState(ctx)
	must.NoError(err)
	is.Len(contextState.Cookies, 2)

	requestState, err := apiCtx.StorageState(ctx)
	must.NoError(err)

	// Both should export the same cookies.
	is.Equal(len(contextState.Cookies), len(requestState.Cookies))
	cNames := cookieNames(contextState.Cookies)
	rNames := cookieNames(requestState.Cookies)
	is.ElementsMatch(cNames, rNames)
}

// Ref: TestBrowserContextFetch.java#shouldAbortRequestsWhenBrowserContextCloses  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldAbortRequestsWhenBrowserContextCloses(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.BlockRequest("/empty.html")

	// Create a separate context with manual lifecycle so we can close it mid-request.
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err)

	// Close the context from a goroutine after a short delay.
	go func() {
		time.Sleep(500 * time.Millisecond)
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	}()

	_, err = bCtx.Request().Get(ctx, srv.EmptyPage())
	is.Error(err, "request should fail when context closes")
	is.True(
		strings.Contains(err.Error(), "Request context disposed") ||
			strings.Contains(err.Error(), "Target page, context or browser has been closed") ||
			strings.Contains(err.Error(), "closed") ||
			strings.Contains(err.Error(), "disposed"),
		"unexpected error: %v", err)
}

// TestBrowserContextFetchShouldThrowOnNetworkErrorAfterRedirect verifies that a network error
// on the destination of a redirect produces an error.
// Ref: TestBrowserContextFetch.java#shouldThrowOnNetworkErrorAfterRedirect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldThrowOnNetworkErrorAfterRedirect(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/redirect", "/fail")
	srv.SetRoute("/fail", func(w http.ResponseWriter, r *http.Request) {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
	})

	_, apiCtx := newBCtxWithRequest(t)

	_, err := apiCtx.Get(ctx, srv.Prefix()+"/redirect")
	is.Error(err, "expected network error after redirect")
}

// TestBrowserContextFetchShouldThrowOnNetworkErrorWhenSendingBody verifies that a connection
// drop while the server receives the request body produces an error.
// Ref: TestBrowserContextFetch.java#shouldThrowOnNetworkErrorWhenSendingBody  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldThrowOnNetworkErrorWhenSendingBody(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/drop", func(w http.ResponseWriter, r *http.Request) {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.(*net.TCPConn).SetLinger(0)
			conn.Close()
		}
	})

	_, apiCtx := newBCtxWithRequest(t)

	_, err := apiCtx.Post(ctx, srv.Prefix()+"/drop", &playwright.APIFetchOptions{
		Data: []byte("my request body"),
	})
	is.Error(err, "expected network error when connection is dropped")
}

// TestBrowserContextFetchShouldThrowOnNetworkErrorWhenSendingBodyAfterRedirect verifies that a
// connection drop at the redirect destination with a POST body produces an error.
// Ref: TestBrowserContextFetch.java#shouldThrowOnNetworkErrorWhenSendingBodyAfterRedirect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldThrowOnNetworkErrorWhenSendingBodyAfterRedirect(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/start", "/drop")
	srv.SetRoute("/drop", func(w http.ResponseWriter, r *http.Request) {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.(*net.TCPConn).SetLinger(0)
			conn.Close()
		}
	})

	_, apiCtx := newBCtxWithRequest(t)

	_, err := apiCtx.Post(ctx, srv.Prefix()+"/start", &playwright.APIFetchOptions{
		Data: []byte("body"),
	})
	is.Error(err, "expected network error after redirect with body")
}

// TestBrowserContextFetchShouldNotAddContextCookieIfCookieHeaderPassedAsAParameter verifies
// that when an explicit Cookie header is provided, session cookies are not added.
// Ref: TestBrowserContextFetch.java#shouldNotAddContextCookieIfCookieHeaderPassedAsAParameter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldNotAddContextCookieIfCookieHeaderPassedAsAParameter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/path")

	bCtx, apiCtx := newBCtxWithRequest(t)

	domain := "localhost"
	path := "/"
	expires := -1.0
	httpOnly := false
	secure := false
	sameSite := playwright.SameSiteLax
	err := bCtx.AddCookies(ctx, []playwright.Cookie{
		{
			Name:     "session",
			Value:    "session-value",
			Domain:   &domain,
			Path:     &path,
			Expires:  &expires,
			HTTPOnly: &httpOnly,
			Secure:   &secure,
			SameSite: sameSite,
		},
	})
	must.NoError(err)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/path", &playwright.APIFetchOptions{
		Headers: map[string]string{"Cookie": "explicit=value"},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	req := <-ch
	cookieHeader := req.Header.Get("Cookie")
	is.Contains(cookieHeader, "explicit=value")
	is.NotContains(cookieHeader, "session=session-value",
		"context cookie should not be added when Cookie header is explicitly provided")
}

// TestBrowserContextFetchShouldReturnErrorWithWrongCredentials verifies that incorrect
// basic-auth credentials result in a 401 response.
// Ref: TestBrowserContextFetch.java#shouldReturnErrorWithWrongCredentials  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldReturnErrorWithWrongCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	_, apiCtx := newBCtxWithRequest(t)

	wrongCreds := base64.StdEncoding.EncodeToString([]byte("user:wrong"))
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected", &playwright.APIFetchOptions{
		Headers: map[string]string{"Authorization": "Basic " + wrongCreds},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status())
}

// TestBrowserContextFetchGetShouldSupportPostData verifies that a GET request can carry
// a body and the server receives it.
// Ref: TestBrowserContextFetch.java#getShouldSupportPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchGetShouldSupportPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedMethod string
	var receivedBody []byte
	srv.SetRoute("/echo", func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/echo", &playwright.APIFetchOptions{
		Data: []byte("my body"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("GET", receivedMethod)
	is.Equal("my body", string(receivedBody))
}

// TestBrowserContextFetchHeadShouldSupportPostData verifies that a HEAD request with a
// body does not produce an error.
// Ref: TestBrowserContextFetch.java#headShouldSupportPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchHeadShouldSupportPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedMethod string
	srv.SetRoute("/echo", func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Head(ctx, srv.Prefix()+"/echo", &playwright.APIFetchOptions{
		Data: []byte("my body"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("HEAD", receivedMethod)
}

// TestBrowserContextFetchShouldAddDefaultHeaders verifies that Playwright adds default
// Accept and Accept-Encoding headers to requests.
// Ref: TestBrowserContextFetch.java#shouldAddDefaultHeaders  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldAddDefaultHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/empty.html")

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	req := <-ch
	is.NotEmpty(req.Header.Get("Accept"), "Accept header should be set by default")
	is.NotEmpty(req.Header.Get("Accept-Encoding"), "Accept-Encoding header should be set by default")
}

// TestBrowserContextFetchShouldAddDefaultHeadersToRedirects verifies that default headers
// (Accept) are propagated to redirect destinations.
// Ref: TestBrowserContextFetch.java#shouldAddDefaultHeadersToRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldAddDefaultHeadersToRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var finalAccept string
	srv.ServeWithRedirect("/redirect", "/final")
	srv.SetRoute("/final", func(w http.ResponseWriter, r *http.Request) {
		finalAccept = r.Header.Get("Accept")
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/redirect")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.NotEmpty(finalAccept, "Accept header should be present after redirect")
}

// TestBrowserContextFetchShouldAllowToOverrideDefaultHeaders verifies that per-request headers
// override default ones (e.g., Accept).
// Ref: TestBrowserContextFetch.java#shouldAllowToOverrideDefaultHeaders  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldAllowToOverrideDefaultHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedAccept string
	srv.SetRoute("/test", func(w http.ResponseWriter, r *http.Request) {
		receivedAccept = r.Header.Get("Accept")
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/test", &playwright.APIFetchOptions{
		Headers: map[string]string{"Accept": "text/html"},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("text/html", receivedAccept)
}

// TestBrowserContextFetchShouldPropagateCustomHeadersWithRedirects verifies that custom
// per-request headers are forwarded after a redirect.
// Ref: TestBrowserContextFetch.java#shouldPropagateCustomHeadersWithRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldPropagateCustomHeadersWithRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var h1, h2 string
	srv.SetRoute("/hop1", func(w http.ResponseWriter, r *http.Request) {
		h1 = r.Header.Get("My-Header")
		http.Redirect(w, r, "/hop2", http.StatusFound)
	})
	srv.SetRoute("/hop2", func(w http.ResponseWriter, r *http.Request) {
		h2 = r.Header.Get("My-Header")
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/hop1", &playwright.APIFetchOptions{
		Headers: map[string]string{"My-Header": "my-value"},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("my-value", h1, "custom header should be on first hop")
	is.Equal("my-value", h2, "custom header should propagate after redirect")
}

// TestBrowserContextFetchShouldPropagateExtraHttpHeadersWithRedirects verifies that context-level
// extra HTTP headers are propagated after a redirect.
// Ref: TestBrowserContextFetch.java#shouldPropagateExtraHttpHeadersWithRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldPropagateExtraHttpHeadersWithRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var h1, h2 string
	srv.SetRoute("/hop1", func(w http.ResponseWriter, r *http.Request) {
		h1 = r.Header.Get("X-Extra")
		http.Redirect(w, r, "/hop2", http.StatusFound)
	})
	srv.SetRoute("/hop2", func(w http.ResponseWriter, r *http.Request) {
		h2 = r.Header.Get("X-Extra")
		w.WriteHeader(200)
	})

	bCtx, apiCtx := newBCtxWithRequest(t)

	err := bCtx.SetExtraHTTPHeaders(ctx, map[string]string{"X-Extra": "extra-value"})
	must.NoError(err)

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/hop1")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal("extra-value", h1, "extra header should be on first hop")
	is.Equal("extra-value", h2, "extra header should propagate after redirect")
}

// TestBrowserContextFetchShouldThrowOnInvalidHeaderValue verifies that a header value
// containing control characters causes an error.
// Ref: TestBrowserContextFetch.java#shouldThrowOnInvalidHeaderValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldThrowOnInvalidHeaderValue(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	_, apiCtx := newBCtxWithRequest(t)

	_, err := apiCtx.Get(ctx, srv.EmptyPage(), &playwright.APIFetchOptions{
		Headers: map[string]string{"Invalid-Header": "value\r\ninjected"},
	})
	is.Error(err, "should error on header value with CRLF")
}

// TestBrowserContextFetchShouldThrowOnNonHttpSProtocol verifies that a non-HTTP/HTTPS URL
// produces an error.
// Ref: TestBrowserContextFetch.java#shouldThrowOnNonHttpSProtocol  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldThrowOnNonHttpSProtocol(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	_, apiCtx := newBCtxWithRequest(t)

	_, err := apiCtx.Get(ctx, "file:///etc/passwd")
	is.Error(err, "should error on non-HTTP/HTTPS protocol")
}

// TestBrowserContextFetchShouldRespectTimeoutAfterRedirects verifies that per-request timeout
// applies even after following a redirect.
// Ref: TestBrowserContextFetch.java#shouldRespectTimeoutAfterRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldRespectTimeoutAfterRedirects(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithRedirect("/redirect", "/slow")
	srv.BlockRequest("/slow")

	_, apiCtx := newBCtxWithRequest(t)

	timeout := 300.0
	_, err := apiCtx.Get(ctx, srv.Prefix()+"/redirect", &playwright.APIFetchOptions{
		Timeout: &timeout,
	})
	is.Error(err, "expected timeout after redirect")
	is.True(
		strings.Contains(strings.ToLower(err.Error()), "timeout"),
		"error should mention timeout, got: %v", err)
}

// TestBrowserContextFetchShouldSupportApplicationXWwwFormUrlencoded verifies that FormData
// is sent as application/x-www-form-urlencoded.
// Ref: TestBrowserContextFetch.java#shouldSupportApplicationXWwwFormUrlencoded  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportApplicationXWwwFormUrlencoded(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedContentType, receivedBody string
	srv.SetRoute("/form", func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		receivedBody = string(b)
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/form", &playwright.APIFetchOptions{
		FormData: []playwright.FormDataField{
			{Name: "firstName", Value: "John"},
			{Name: "lastName", Value: "Doe"},
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Contains(receivedContentType, "application/x-www-form-urlencoded")
	is.Contains(receivedBody, "firstName=John")
	is.Contains(receivedBody, "lastName=Doe")
}

// TestBrowserContextFetchShouldEncodeToApplicationJsonByDefault verifies that JSON bytes
// sent with Content-Type: application/json arrive unmodified at the server.
// Ref: TestBrowserContextFetch.java#shouldEncodeToApplicationJsonByDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldEncodeToApplicationJsonByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedContentType string
	var receivedBody []byte
	srv.SetRoute("/json", func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	payload := `{"key":"value"}`
	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/json", &playwright.APIFetchOptions{
		Headers: map[string]string{"Content-Type": "application/json"},
		Data:    []byte(payload),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Contains(receivedContentType, "application/json")
	is.Equal(payload, string(receivedBody))
}

// TestBrowserContextFetchShouldSupportMultipartFormData verifies that MultipartData fields
// are sent as multipart/form-data and the server can parse them.
// Ref: TestBrowserContextFetch.java#shouldSupportMultipartFormData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportMultipartFormData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	type formCapture struct {
		contentType string
		firstName   string
		lastName    string
	}
	captCh := make(chan formCapture, 1)
	srv.SetRoute("/upload", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		_ = r.ParseMultipartForm(32 << 20)
		captCh <- formCapture{
			contentType: ct,
			firstName:   r.FormValue("firstName"),
			lastName:    r.FormValue("lastName"),
		}
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/upload", &playwright.APIFetchOptions{
		MultipartData: []playwright.MultipartField{
			{Name: "firstName", Value: "John"},
			{Name: "lastName", Value: "Doe"},
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	capt := <-captCh
	is.Contains(capt.contentType, "multipart/form-data")
	is.Equal("John", capt.firstName)
	is.Equal("Doe", capt.lastName)
}

// TestBrowserContextFetchShouldSupportRepeatingNamesInMultipartFormData verifies that
// multiple multipart fields with the same name are all received by the server.
// Ref: TestBrowserContextFetch.java#shouldSupportRepeatingNamesInMultipartFormData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportRepeatingNamesInMultipartFormData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	formCh := make(chan map[string][]string, 1)
	srv.SetRoute("/upload", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(32 << 20)
		values := make(map[string][]string)
		if r.MultipartForm != nil {
			for k, v := range r.MultipartForm.Value {
				values[k] = v
			}
		}
		formCh <- values
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/upload", &playwright.APIFetchOptions{
		MultipartData: []playwright.MultipartField{
			{Name: "field", Value: "value1"},
			{Name: "field", Value: "value2"},
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	form := <-formCh
	is.ElementsMatch([]string{"value1", "value2"}, form["field"],
		"both values for repeated field name should be received")
}

// TestBrowserContextFetchShouldSerializeDataToJsonRegardlessOfContentType verifies that
// raw Data bytes are forwarded unchanged regardless of the Content-Type header.
// Ref: TestBrowserContextFetch.java#shouldSerializeDataToJsonRegardlessOfContentType  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSerializeDataToJsonRegardlessOfContentType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedBody []byte
	srv.SetRoute("/post", func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	payload := `{"foo":"bar"}`
	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/post", &playwright.APIFetchOptions{
		Headers: map[string]string{"Content-Type": "text/plain"},
		Data:    []byte(payload),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(payload, string(receivedBody))
}

// TestBrowserContextFetchShouldNotThrowWhenDataPassedForUnsupportedRequest verifies that
// passing a body to methods that don't conventionally carry one (GET, HEAD) does not error.
// Ref: TestBrowserContextFetch.java#shouldNotThrowWhenDataPassedForUnsupportedRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldNotThrowWhenDataPassedForUnsupportedRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	// HEAD with body should not error.
	resp, err := apiCtx.Head(ctx, srv.Prefix()+"/ok", &playwright.APIFetchOptions{
		Data: []byte("ignored body"),
	})
	must.NoError(err, "HEAD with body should not throw")
	defer resp.Dispose(ctx) //nolint:errcheck
}

// TestBrowserContextFetchShouldWorkWithSetHTTPCredentials verifies that BrowserContext-level
// HTTP credentials are applied to API requests automatically.
// Ref: TestBrowserContextFetch.java#shouldWorkWithSetHTTPCredentials  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldWorkWithSetHTTPCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
}

// TestBrowserContextFetchShouldReturnErrorWithWrongCredentialsContextLevel verifies that
// incorrect BrowserContext-level credentials result in a 401 response.
// Ref: TestBrowserContextFetch.java#shouldReturnErrorWithWrongCredentials (context-level)
func TestBrowserContextFetchShouldReturnErrorWithWrongCredentialsContextLevel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "wrong",
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status())
}

// TestBrowserContextFetchShouldWorkWithSetHTTPCredentialsAndMatchingOrigin verifies that
// context-level credentials are sent when the origin matches exactly.
// Ref: TestBrowserContextFetch.java#shouldWorkWithSetHTTPCredentialsAndMatchingOrigin  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldWorkWithSetHTTPCredentialsAndMatchingOrigin(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	origin := srv.Prefix()
	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &origin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
}

// TestBrowserContextFetchShouldWorkWithSetHTTPCredentialsAndMatchingOriginCaseInsensitive
// verifies that origin matching for context-level credentials is case-insensitive.
// Ref: TestBrowserContextFetch.java#shouldWorkWithSetHTTPCredentialsAndMatchingOriginCaseInsensitive  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldWorkWithSetHTTPCredentialsAndMatchingOriginCaseInsensitive(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	origin := strings.ToUpper(srv.Prefix())
	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
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

// TestBrowserContextFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginScheme verifies
// that context-level credentials are NOT sent when the configured origin has a wrong scheme.
// Ref: TestBrowserContextFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginScheme  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginScheme(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	wrongSchemeOrigin := strings.Replace(srv.Prefix(), "http://", "https://", 1)
	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &wrongSchemeOrigin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status())
}

// TestBrowserContextFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginHostname verifies
// that context-level credentials are NOT sent when the configured origin has a wrong hostname.
// Ref: TestBrowserContextFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginHostname  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginHostname(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "pass", "text/plain", "ok")

	u, err := url.Parse(srv.Prefix())
	must.NoError(err)
	wrongHostOrigin := fmt.Sprintf("%s://wrong-hostname:%s", u.Scheme, u.Port())
	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &wrongHostOrigin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status())
}

// TestBrowserContextFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginPort verifies
// that context-level credentials are NOT sent when the configured origin has a wrong port.
// Ref: TestBrowserContextFetch.java#shouldReturnErrorWithCorrectCredentialsAndWrongOriginPort  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldReturnErrorWithCorrectCredentialsAndWrongOriginPort(t *testing.T) {
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
	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Origin:   &wrongPortOrigin,
		},
	})

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/protected")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(401, resp.Status())
}

// TestBrowserContextFetchShouldSerializeNullValuesInPostData verifies that JSON with null
// values is sent and received correctly.
// Ref: TestBrowserContextFetch.java#shouldSerializeNullValuesInPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSerializeNullValuesInPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	var receivedBody []byte
	srv.SetRoute("/post", func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	})

	_, apiCtx := newBCtxWithRequest(t)

	payload := `{"nullKey":null,"strKey":"value"}`
	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/post", &playwright.APIFetchOptions{
		Headers: map[string]string{"Content-Type": "application/json"},
		Data:    []byte(payload),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(payload, string(receivedBody))
}

// TestBrowserContextFetchShouldSupportHTTPCredentialsSendImmediatelyForNewContext verifies
// that Send="always" on BrowserContext credentials causes Authorization to be sent on the
// very first request without a prior 401 challenge.
// Ref: TestBrowserContextFetch.java#shouldSupportHTTPCredentialsSendImmediatelyForNewContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportHTTPCredentialsSendImmediatelyForNewContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch := testserver.WaitForRequest(srv, "/direct-auth")

	sendAlways := "always"
	origin := srv.Prefix()
	_, apiCtx := newBCtxWithRequest(t, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Send:     &sendAlways,
			Origin:   &origin,
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

// TestBrowserContextFetchShouldNotWorkAfterContextDispose verifies that requests fail after
// the owning BrowserContext has been closed.
// Ref: TestBrowserContextFetch.java#shouldNotWorkAfterContextDispose  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldNotWorkAfterContextDispose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err)

	apiCtx := bCtx.Request()

	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = bCtx.Close(closeCtx)
	must.NoError(err)

	_, err = apiCtx.Get(ctx, srv.EmptyPage())
	is.Error(err, "request should fail after context is closed")
}

// TestBrowserContextFetchShouldOverrideRequestParameters verifies that APIRequestContext.Fetch
// can override method, headers, and body relative to a page navigation request.
// The Java version passes the captured Request object directly; here we navigate to obtain
// the URL and re-issue the request via APIRequestContext with explicit overrides.
// Ref: TestBrowserContextFetch.java#shouldOverrideRequestParameters  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldOverrideRequestParameters(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	type captured struct {
		method string
		foo    string
		body   []byte
	}
	captCh := make(chan captured, 1)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		captCh <- captured{
			method: r.Method,
			foo:    r.Header.Get("foo"),
			body:   b,
		}
		w.WriteHeader(200)
	})

	bCtx, apiCtx := newBCtxWithRequest(t)

	// Navigate to the page first (mimics page.waitForRequest in Java).
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	defer page.Close(ctx)               //nolint:errcheck
	_ = page.Goto(ctx, srv.EmptyPage()) // navigation GET is captured but we only care about the next fetch
	<-captCh                            // drain the navigation request

	// Re-fetch the same URL via APIRequestContext with overrides: POST, custom header, and body.
	method := "POST"
	resp, err := apiCtx.Fetch(ctx, srv.EmptyPage(), &playwright.APIFetchOptions{
		Method:  &method,
		Headers: map[string]string{"foo": "bar"},
		Data:    []byte("data"),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	capt := <-captCh
	is.Equal("POST", capt.method)
	is.Equal("bar", capt.foo)
	is.Equal("data", string(capt.body))
}

// TestBrowserContextFetchShouldSupportMultipartFormDataWithPathValues verifies that a file
// specified as a filesystem path is correctly uploaded as a multipart part, with filename and
// MIME type inferred from the path.
// Ref: TestBrowserContextFetch.java#shouldSupportMultipartFormDataWithPathValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportMultipartFormDataWithPathValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	type fileCapture struct {
		contentType string
		firstName   string
		lastName    string
		fileName    string
		fileMIME    string
		fileContent []byte
	}
	captCh := make(chan fileCapture, 1)
	srv.SetRoute("/upload", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		_ = r.ParseMultipartForm(32 << 20)

		var (
			fileName    string
			fileMIME    string
			fileContent []byte
		)
		if f, fh, err := r.FormFile("file"); err == nil && fh != nil {
			fileName = fh.Filename
			fileMIME = fh.Header.Get("Content-Type")
			fileContent, _ = io.ReadAll(f)
			f.Close() //nolint:errcheck
		}
		captCh <- fileCapture{
			contentType: ct,
			firstName:   r.FormValue("firstName"),
			lastName:    r.FormValue("lastName"),
			fileName:    fileName,
			fileMIME:    fileMIME,
			fileContent: fileContent,
		}
		w.WriteHeader(200)
	})

	// Write a temporary JSON file to upload.
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "simplezip.json")
	must.NoError(os.WriteFile(filePath, []byte(`{"foo":"bar"}`), 0o600))

	_, apiCtx := newBCtxWithRequest(t)

	resp, err := apiCtx.Post(ctx, srv.Prefix()+"/upload", &playwright.APIFetchOptions{
		MultipartData: []playwright.MultipartField{
			{Name: "firstName", Value: "John"},
			{Name: "lastName", Value: "Doe"},
			{Name: "file", FilePath: filePath},
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())

	capt := <-captCh
	is.Contains(capt.contentType, "multipart/form-data")
	is.Equal("John", capt.firstName)
	is.Equal("Doe", capt.lastName)
	is.Equal("simplezip.json", capt.fileName)
	is.Equal(`{"foo":"bar"}`, string(capt.fileContent))
	// MIME type is derived from the .json extension.
	is.Contains(capt.fileMIME, "application/json")
}

// TestBrowserContextFetchShouldSupportHTTPCredentialsSendImmediatelyForBrowserNewPage verifies
// that Browser.NewPage() accepts HttpCredentials with Send="always" so the Authorization header
// is included on the very first request without a prior 401 challenge.
// When the request targets a different origin the header must NOT be sent.
// Ref: TestBrowserContextFetch.java#shouldSupportHTTPCredentialsSendImmediatelyForBrowserNewPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportHTTPCredentialsSendImmediatelyForBrowserNewPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	ch1 := testserver.WaitForRequest(srv, "/auth-check")

	sendAlways := "always"
	origin := strings.ToUpper(srv.Prefix())
	page, err := globalBrowser.NewPage(ctx, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "user",
			Password: "pass",
			Send:     &sendAlways,
			Origin:   &origin,
		},
	})
	must.NoError(err)
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = page.Close(closeCtx)
	}()

	apiCtx := page.Context().Request()

	// First request: same origin → Authorization header must be sent immediately.
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/auth-check")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, resp.Status())

	req1 := <-ch1
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	is.Equal(expected, req1.Header.Get("Authorization"),
		"Authorization header should be sent on first request (same origin)")

	// Second request: cross-process origin → Authorization must NOT be sent.
	var crossAuth string
	srv.SetRoute("/auth-check", func(w http.ResponseWriter, r *http.Request) {
		crossAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
	})

	resp2, err := apiCtx.Get(ctx, srv.CrossProcessPrefix()+"/auth-check")
	must.NoError(err)
	defer resp2.Dispose(ctx) //nolint:errcheck
	is.Equal(200, resp2.Status())
	is.Empty(crossAuth, "Authorization should NOT be sent to a different origin")
}

// Ref: TestBrowserContextFetch.java#shouldWorkWithHttpCredentials  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestShouldWorkWithHttpCredentials(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/guarded", "u", "p", "text/html", "<p>ok</p>")

	statuses := make(chan int, 5)
	page := newPage(t)
	page.OnResponse(func(resp *playwright.NetworkResponse) {
		statuses <- resp.Status()
	})

	_ = page.Goto(ctx, srv.Prefix()+"/guarded")

	select {
	case status := <-statuses:
		if status != 401 {
			t.Errorf("Status = %d, want 401", status)
		}
	case <-time.After(5 * time.Second):
		t.Error("no response received within 5 seconds")
	}
}

// ---------------------------------------------------------------------------
// From browser_context_http_auth_extra_test.go
// ---------------------------------------------------------------------------

// TestBrowserContextFetchPageRequestCookiesShouldReturnBrowserContextCookies verifies the Java behavior.
// Ref: TestBrowserContextFetch.java#pageRequestCookiesShouldReturnBrowserContextCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchPageRequestCookiesShouldReturnBrowserContextCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	emptyPage := srv.EmptyPage()
	domain := "example.com"
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "a", Value: "b", URL: &emptyPage},
		{Name: "c", Value: "d", Domain: &domain, Path: a4Ptr("/")},
	}))

	allCookies, err := page.Request().Cookies(ctx)
	must.NoError(err)
	must.Len(allCookies, 2)
	is.ElementsMatch([]string{"a", "c"}, []string{allCookies[0].Name, allCookies[1].Name})
	urlCookies, err := page.Request().Cookies(ctx, emptyPage)
	must.NoError(err)
	must.Len(urlCookies, 1)
	is.Equal("a", urlCookies[0].Name)

	contextCookies, err := bCtx.Cookies(ctx, emptyPage)
	must.NoError(err)
	is.Equal(contextCookies, urlCookies)
}

// TestBrowserContextFetchPageRequestClearCookiesShouldClearBrowserContextCookies verifies the Java behavior.
// Ref: TestBrowserContextFetch.java#pageRequestClearCookiesShouldClearBrowserContextCookies  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchPageRequestClearCookiesShouldClearBrowserContextCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	emptyPage := srv.EmptyPage()
	must.NoError(bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "a", Value: "b", URL: &emptyPage},
		{Name: "c", Value: "d", URL: &emptyPage},
	}))
	name := "a"
	must.NoError(page.Request().ClearCookies(ctx, &playwright.ClearCookiesOptions{Name: &name}))
	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	must.Len(cookies, 1)
	is.Equal("c", cookies[0].Name)

	must.NoError(page.Request().ClearCookies(ctx))
	cookies, err = bCtx.Cookies(ctx)
	must.NoError(err)
	is.Empty(cookies)

	requests := testserver.WaitForRequest(srv, "/empty.html")
	response, err := page.Request().Get(ctx, emptyPage)
	must.NoError(err)
	defer response.Dispose(ctx)
	request := <-requests
	is.Empty(request.Header.Values("Cookie"))
}

// TestBrowserContextFetchShouldSerializeDateAndLocalDateTime verifies the Java behavior.
// Ref: TestBrowserContextFetch.java#shouldSerializeDateAndLocalDateTime  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSerializeDateAndLocalDateTime(t *testing.T) {
	t.Skip("not implemented: Java-specific Date/LocalDateTime serialization and APIRequestContext.Fetch(Request)")
}

// TestBrowserContextFetchShouldSupportOffsetDateTimeInData verifies the Java behavior.
// Ref: TestBrowserContextFetch.java#shouldSupportOffsetDateTimeInData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportOffsetDateTimeInData(t *testing.T) {
	t.Skip("not implemented: Java-specific OffsetDateTime serialization of request data")
}

// TestBrowserContextFetchShouldSupportLocalDateInData verifies the Java behavior.
// Ref: TestBrowserContextFetch.java#shouldSupportLocalDateInData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportLocalDateInData(t *testing.T) {
	t.Skip("not implemented: Java-specific LocalDate serialization of request data")
}

// TestBrowserContextFetchShouldSupportMultipleHttpCredentials verifies the Java behavior.
// Ref: TestBrowserContextFetch.java#shouldSupportMultipleHttpCredentials  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchShouldSupportMultipleHttpCredentials(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.HttpCredentials as a list")
}

// TestBrowserContextFetchFetchShouldWork verifies BrowserContext.Request().Fetch.
// Ref: TestBrowserContextFetch.java#fetchShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchFetchShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/simple.json", "application/json", "{\"foo\": \"bar\"}\n")
	bCtx := newContext(t)
	resp, err := bCtx.Request().Fetch(ctx, srv.Prefix()+"/simple.json")
	must.NoError(err)
	is.Equal(srv.Prefix()+"/simple.json", resp.URL())
	is.Equal(200, resp.Status())
	is.Equal("OK", resp.StatusText())
	is.True(resp.OK())
	is.Equal("application/json", resp.Headers()["content-type"])
	text, err := resp.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}

// TestBrowserContextFetchPageRequestAddCookiesShouldAddCookiesToTheBrowserContext adds cookies through the page request context.
// Ref: TestBrowserContextFetch.java#pageRequestAddCookiesShouldAddCookiesToTheBrowserContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextFetch.java
func TestBrowserContextFetchPageRequestAddCookiesShouldAddCookiesToTheBrowserContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	emptyPage := srv.EmptyPage()
	domain := strings.TrimPrefix(strings.TrimPrefix(srv.Prefix(), "http://"), "https://")
	domain = strings.Split(domain, ":")[0]
	expires := float64(time.Now().Add(time.Hour).Unix())
	must.NoError(page.Request().AddCookies(ctx, []playwright.Cookie{
		{Name: "a", Value: "b", URL: &emptyPage},
		{Name: "c", Value: "d", Domain: &domain, Path: a4Ptr("/"), Expires: &expires},
	}))
	cookies, err := bCtx.Cookies(ctx)
	must.NoError(err)
	is.ElementsMatch([]string{"a=b", "c=d"}, []string{
		cookies[0].Name + "=" + cookies[0].Value,
		cookies[1].Name + "=" + cookies[1].Value,
	})

	requests := testserver.WaitForRequest(srv, "/empty.html")
	response, err := page.Request().Get(ctx, emptyPage)
	must.NoError(err)
	defer response.Dispose(ctx)
	request := <-requests
	is.ElementsMatch([]string{"a=b", "c=d"}, strings.Split(strings.Join(request.Header.Values("Cookie"), "; "), "; "))

	must.NoError(page.Goto(ctx, emptyPage))
	documentCookies, err := page.Evaluate(ctx, "document.cookie")
	must.NoError(err)
	must.IsType("", documentCookies)
	is.ElementsMatch([]string{"a=b", "c=d"}, strings.Split(documentCookies.(string), "; "))
}
