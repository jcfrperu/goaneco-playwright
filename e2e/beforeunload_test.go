//go:build e2e

// E2E tests for page beforeunload behavior.
// Migration of: TestBeforeunload.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
)

// TestBeforeunloadNavigateAwayFromPage verifies that a page with a beforeunload handler
// can be navigated away from without hanging.
// Ref: TestBeforeunload.java#shouldBeAbleToNavigateAwayFromPageWithBeforeunload
func TestBeforeunloadNavigateAwayFromPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// /beforeunload.html is a built-in testserver page that adds a beforeunload handler.
	must.NoError(page.Goto(ctx, srv.Prefix()+"/beforeunload.html"), "Goto beforeunload.html failed")

	// We must interact with the page so that the beforeunload handler fires on navigation.
	must.NoError(page.Locator("body").Click(ctx), "body click failed")

	// Navigate away — should not hang even with the beforeunload dialog.
	must.NoError(page.Goto(ctx, srv.EmptyPage()), "navigate away from beforeunload page failed")
}
