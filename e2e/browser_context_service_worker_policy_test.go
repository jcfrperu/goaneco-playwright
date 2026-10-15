//go:build e2e

// E2E tests for: TestBrowserContextServiceWorkerPolicy.java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestBrowserContextServiceWorkerPolicyShouldAllowServiceWorkersByDefault verifies that service workers are allowed by default.
// Ref: TestBrowserContextServiceWorkerPolicy.java#shouldAllowServiceWorkersByDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextServiceWorkerPolicy.java
func TestBrowserContextServiceWorkerPolicyShouldAllowServiceWorkersByDefault(t *testing.T) {
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

// TestBrowserContextServiceWorkerPolicyBlocksServiceWorkerRegistration verifies that service workers can be blocked via context policy.
// Ref: TestBrowserContextServiceWorkerPolicy.java#blocksServiceWorkerRegistration  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextServiceWorkerPolicy.java
func TestBrowserContextServiceWorkerPolicyBlocksServiceWorkerRegistration(t *testing.T) {
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
