//go:build e2e

// E2E tests for BrowserContext service worker policy.
// Migration of: TestBrowserContextServiceWorkerPolicy.java
package e2e

import (
	"net/http"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
)

func setupServiceWorkerRoutes(t *testing.T, srv *testserver.Server) {
	t.Helper()
	srv.SetRoute("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.addEventListener('install', () => self.skipWaiting());`))
	})
	srv.ServeWithBody("/sw-page.html", "text/html", `<!DOCTYPE html><html><body>
<script>
  window.registrationPromise = navigator.serviceWorker
    ? navigator.serviceWorker.register('/sw.js')
    : null;
</script>
</body></html>`)
}

// TestServiceWorkersAllowedByDefault verifies that service workers are allowed by default.
// Ref: TestBrowserContextServiceWorkerPolicy.java#shouldAllowServiceWorkersByDefault
func TestServiceWorkersAllowedByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	setupServiceWorkerRoutes(t, srv)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/sw-page.html"))

	result, err := page.Evaluate(ctx, "() => window['registrationPromise'] !== undefined")
	must.NoError(err, "Evaluate failed")
	if v, ok := result.(bool); !ok || !v {
		t.Error("expected registrationPromise to be defined (service workers allowed by default)")
	}
}

// TestServiceWorkersBlocked verifies that service workers can be blocked via context policy.
// Ref: TestBrowserContextServiceWorkerPolicy.java#blocksServiceWorkerRegistration
func TestServiceWorkersBlocked(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	setupServiceWorkerRoutes(t, srv)

	block := "block"
	bCtx := newContextWithCleanup(t, &playwright.BrowserContextOptions{
		ServiceWorkers: &block,
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	must.NoError(page.Goto(ctx, srv.Prefix()+"/sw-page.html"))

	// When service workers are blocked, navigator.serviceWorker is still defined
	// but the registration promise rejects. The promise may also be null if blocked
	// at the browser level before script execution.
	result, err := page.Evaluate(ctx, `() => {
		if (!window.registrationPromise) return 'null-promise';
		return window.registrationPromise
			.then(() => 'registered')
			.catch(() => 'blocked');
	}`)
	must.NoError(err, "Evaluate failed")
	// Result should be "blocked" or "null-promise" depending on browser implementation.
	if v, ok := result.(string); ok {
		if v == "registered" {
			t.Errorf("service worker should not have registered when policy is 'block', got: %q", v)
		}
	}
}
