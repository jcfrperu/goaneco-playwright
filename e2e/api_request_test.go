//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestAPIResponseStatusText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/api-st", "text/plain", "ok")

	apiCtx, err := globalPW.NewAPIRequestContext(ctx)
	must.NoError(err, "NewAPIRequestContext failed")
	defer apiCtx.Dispose(ctx) //nolint:errcheck

	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/api-st")
	must.NoError(err, "Get failed")
	defer resp.Dispose(ctx) //nolint:errcheck

	// HTTP/1.1 returns "OK"; HTTP/2 may return an empty string.
	st := resp.StatusText()
	is.True(st == "OK" || st == "", "StatusText() = %q, want 'OK' or ''", st)
}

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
