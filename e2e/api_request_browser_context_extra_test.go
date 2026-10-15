//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// serveSimpleJSON registers a route serving `{"foo": "bar"}\n` as application/json.
func serveSimpleJSON(srv *testserver.Server) {
	srv.ServeWithBody("/simple.json", "application/json", "{\"foo\": \"bar\"}\n")
}

// newBCtxWithRequest creates a BrowserContext that auto-closes on test cleanup
// and returns it along with its bound APIRequestContext.
func newBCtxWithRequest(t *testing.T, opts ...*playwright.BrowserContextOptions) (*playwright.BrowserContext, *playwright.APIRequestContext) {
	t.Helper()
	ctx := testCtx(t)
	var bCtx *playwright.BrowserContext
	var err error
	if len(opts) > 0 && opts[0] != nil {
		bCtx, err = globalBrowser.NewContext(ctx, opts[0])
	} else {
		bCtx, err = globalBrowser.NewContext(ctx)
	}
	require.NoError(t, err, "NewContext")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	return bCtx, bCtx.Request()
}

// ── TestBrowserContextFetch equivalents ─────────────────────────────────────

func cookieNames(cookies []playwright.Cookie) []string {
	names := make([]string, len(cookies))
	for i, c := range cookies {
		names[i] = c.Name
	}
	return names
}
