//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// launchPersistent is a test helper that launches a persistent context in a fresh temp dir,
// registers cleanup for both browser and context, and returns the BrowserContext and first Page.
func launchPersistent(t *testing.T, opts *playwright.LaunchPersistentContextOptions) (*playwright.BrowserContext, *playwright.Page) {
	t.Helper()
	ctx := testCtx(t)

	bt, err := globalPW.Chromium()
	require.NoError(t, err, "Chromium() failed")

	if opts == nil {
		opts = &playwright.LaunchPersistentContextOptions{}
	}
	headless := true
	opts.Headless = &headless

	browser, bCtx, err := bt.LaunchPersistentContext(ctx, t.TempDir(), opts)
	require.NoError(t, err, "LaunchPersistentContext failed")
	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(cleanCtx)
		_ = browser.Close(cleanCtx)
	})

	// The persistent context always opens with a default page (about:blank).
	// Create a new page if none are tracked yet.
	pages := bCtx.Pages()
	var page *playwright.Page
	if len(pages) > 0 {
		page = pages[0]
	} else {
		page, err = bCtx.NewPage(ctx)
		require.NoError(t, err, "NewPage failed")
	}
	return bCtx, page
}
