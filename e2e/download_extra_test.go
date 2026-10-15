//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
	"net/http"
	"sync"
	"testing"
	"time"
)

func boolPtr(b bool) *bool { return &b }

// downloadContext creates a BrowserContext with acceptDownloads enabled and registers cleanup.
func downloadContext(t *testing.T) *playwright.BrowserContext {
	t.Helper()
	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		AcceptDownloads: boolPtr(true),
	})
	require.NoError(t, err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	return bCtx
}

// triggerDownload serves a small downloadable file, navigates to a page with a download link,
// clicks it, and returns the Download after the event fires.
func triggerDownload(t *testing.T, bCtx *playwright.BrowserContext, content string) *playwright.Download {
	t.Helper()
	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/dl-file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="download.bin"`)
		_, _ = w.Write([]byte(content))
	})
	srv.ServeWithBody("/dl-page", "text/html", `<a id="a" href="/dl-file" download>dl</a>`)

	page, err := bCtx.NewPage(ctx)
	require.NoError(t, err, "NewPage failed")

	var mu sync.Mutex
	var dl *playwright.Download
	seen := make(chan struct{}, 1)
	page.OnDownload(func(d *playwright.Download) {
		mu.Lock()
		dl = d
		mu.Unlock()
		select {
		case seen <- struct{}{}:
		default:
		}
	})

	err = page.Goto(ctx, srv.Prefix()+"/dl-page")
	require.NoError(t, err, "Goto failed")
	err = page.Locator("#a").Click(ctx)
	require.NoError(t, err, "Click failed")

	select {
	case <-seen:
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for download event")
	}

	mu.Lock()
	d := dl
	mu.Unlock()
	return d
}
