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

// videoSetup creates a BrowserContext with RecordVideo, navigates to a page,
// closes it to flush the recording, and returns the Video (or skips the test if unavailable).
func videoSetup(t *testing.T, dir string) *playwright.Video {
	t.Helper()
	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		RecordVideo: &playwright.RecordVideoOptions{Dir: dir},
	})
	require.NoError(t, err, "NewContext with RecordVideo failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	require.NoError(t, err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	require.NoError(t, err, "Goto failed")

	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = page.Close(closeCtx)
	require.NoError(t, err, "page.Close failed")

	require.NoError(t, page.WaitForTimeout(testCtx(t), 500))

	v := page.Video()
	if v == nil {
		t.Skip("video not available in this environment")
	}
	return v
}
