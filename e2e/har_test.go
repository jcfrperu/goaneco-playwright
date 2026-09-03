//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

// writeHAR writes a single-entry HAR file and returns its path.
func writeHAR(t *testing.T, url, mimeType, body string, status int) string {
	t.Helper()
	content := fmt.Sprintf(`{"log":{"version":"1.2","creator":{"name":"test","version":"1"},"entries":[{`+
		`"startedDateTime":"2024-01-01T00:00:00.000Z",`+
		`"time":0,`+
		`"request":{"method":"GET","url":%q,"headers":[],"queryString":[],"cookies":[],"headersSize":-1,"bodySize":-1},`+
		`"response":{"status":%d,"statusText":"OK",`+
		`"headers":[{"name":"content-type","value":%q}],`+
		`"cookies":[],"content":{"size":%d,"mimeType":%q,"text":%q},`+
		`"redirectURL":"","headersSize":-1,"bodySize":-1},`+
		`"cache":{},"timings":{"send":0,"wait":0,"receive":0}`+
		`}]}}`,
		url, status, mimeType, len(body), mimeType, body)

	dir, err := os.MkdirTemp("", "har-*")
	require.NoError(t, err, "MkdirTemp")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	path := filepath.Join(dir, "test.har")
	err = os.WriteFile(path, []byte(content), 0600)
	require.NoError(t, err, "WriteFile")
	return path
}

func TestRouteFromHAR(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	targetURL := srv.Prefix() + "/har-served"
	harPath := writeHAR(t, targetURL, "text/plain", "served from HAR", 200)

	bCtx := newContext(t)
	err := bCtx.RouteFromHAR(ctx, harPath)
	must.NoError(err, "RouteFromHAR failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, targetURL)
	must.NoError(err, "Goto failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	is.Contains(content, "served from HAR")
}

func TestRouteFromHARPassThrough(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/real-page", "text/plain", "from real server")

	// HAR only has an entry for a different URL — real requests must pass through.
	harPath := writeHAR(t, srv.Prefix()+"/non-existent", "text/plain", "never", 200)

	bCtx := newContext(t)
	err := bCtx.RouteFromHAR(ctx, harPath)
	must.NoError(err, "RouteFromHAR failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.Prefix()+"/real-page")
	must.NoError(err, "Goto failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	if !strings.Contains(content, "from real server") {
		t.Errorf("expected pass-through response 'from real server', got:\n%s", content)
	}
}

func TestRouteFromHARStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	targetURL := srv.Prefix() + "/har-404"
	harPath := writeHAR(t, targetURL, "text/plain", "not found in HAR", 404)

	bCtx := newContext(t)
	err := bCtx.RouteFromHAR(ctx, harPath)
	must.NoError(err, "RouteFromHAR failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// Goto a URL that maps to a 404 in the HAR — should not error (browser loads it).
	_ = page.Goto(ctx, targetURL)

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	if !strings.Contains(content, "not found in HAR") {
		t.Errorf("expected HAR 404 body in content, got:\n%s", content)
	}
}

// Ensure imports are used.
var _ *playwright.RouteFromHAROptions
