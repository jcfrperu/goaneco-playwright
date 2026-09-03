//go:build e2e

package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
)

func TestBrowserTypeLaunchPersistentContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	// Create a temporary directory for user data.
	dir, err := os.MkdirTemp("", "playwright-persistent-*")
	must.NoError(err, "MkdirTemp failed")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	bt, err := globalPW.Chromium()
	must.NoError(err, "Chromium() failed")

	headless := true
	browser, bCtx, err := bt.LaunchPersistentContext(ctx, dir, &playwright.LaunchPersistentContextOptions{
		Headless: &headless,
	})
	must.NoError(err, "LaunchPersistentContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
		_ = browser.Close(closeCtx)
	})

	must.NotNil(browser, "LaunchPersistentContext returned nil browser")
	must.NotNil(bCtx, "LaunchPersistentContext returned nil context")

	// Create a page and verify it works.
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, "about:blank")
	must.NoError(err, "Goto about:blank failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title() failed")
	// about:blank title should be empty string.
	if title != "" {
		t.Logf("title = %q (about:blank usually has empty title)", title)
	}
}
