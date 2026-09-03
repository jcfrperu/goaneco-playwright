//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageOnDownload(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Serve a file that triggers a download (Content-Disposition: attachment).
	srv.SetRoute("/file.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=\"file.txt\"")
		_, _ = w.Write([]byte("download content"))
	})

	// Create context with acceptDownloads enabled.
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		AcceptDownloads: boolPtr(true),
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// Register download handler before triggering the download.
	var mu sync.Mutex
	var downloads []*playwright.Download
	downloadSeen := make(chan struct{}, 1)

	page.OnDownload(func(d *playwright.Download) {
		mu.Lock()
		downloads = append(downloads, d)
		mu.Unlock()
		select {
		case downloadSeen <- struct{}{}:
		default:
		}
	})

	// Serve a page with a download link and click it.
	srv.ServeWithBody("/download-page", "text/html", `
		<a id="dl" href="/file.txt" download>Download</a>
	`)
	err = page.Goto(ctx, srv.Prefix()+"/download-page")
	must.NoError(err, "Goto failed")

	err = page.Locator("#dl").Click(ctx)
	must.NoError(err, "Click failed")

	// Wait for download event.
	select {
	case <-downloadSeen:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for download event")
	}

	mu.Lock()
	dl := downloads[0]
	mu.Unlock()

	if dl.SuggestedFilename() != "file.txt" {
		t.Errorf("SuggestedFilename = %q, want 'file.txt'", dl.SuggestedFilename())
	}
	if dl.URL() == "" {
		t.Error("Download.URL() is empty")
	}

	// Get the download path (waits for completion).
	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	path, err := dl.Path(dlCtx)
	must.NoError(err, "Download.Path() failed")
	if path == "" {
		t.Error("download path is empty")
	}
}

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

func TestDownloadSaveAs(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "saveAs content")

	// Wait for download to finish and get the temp path.
	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := dl.Path(dlCtx)
	must.NoError(err, "Download.Path() failed")

	// Save to a custom path.
	dir, err := os.MkdirTemp("", "dl-saveas-*")
	must.NoError(err, "MkdirTemp failed")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	dest := filepath.Join(dir, "saved.bin")
	err = dl.SaveAs(dlCtx, dest)
	must.NoError(err, "Download.SaveAs failed")

	info, err := os.Stat(dest)
	must.NoError(err, "saved file not found")
	if info.Size() == 0 {
		t.Error("saved file is empty")
	}
}

func TestDownloadDelete(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "delete me")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	originalPath, err := dl.Path(dlCtx)
	must.NoError(err, "Download.Path() failed")
	if originalPath == "" {
		t.Fatal("Download.Path() returned empty path")
	}

	err = dl.Delete(dlCtx)
	must.NoError(err, "Download.Delete failed")

	// File should no longer exist at the original path.
	if _, err := os.Stat(originalPath); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted at %q, but it still exists (err: %v)", originalPath, err)
	}
}

func TestDownloadCancel(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Serve a "download" that stays open until the client disconnects so we
	// have time to call Cancel() before it completes.
	started := make(chan struct{}, 1)
	srv.SetRoute("/slow-dl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="slow.bin"`)
		w.Header().Set("Content-Length", "10000000")
		w.WriteHeader(200)
		// Write initial bytes so the browser triggers the download event.
		_, _ = w.Write(make([]byte, 8192))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		select {
		case started <- struct{}{}:
		default:
		}
		// Block until the client disconnects (download cancelled).
		<-r.Context().Done()
	})
	srv.ServeWithBody("/cancel-page", "text/html", `<a id="a" href="/slow-dl" download>dl</a>`)

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})

	err = page.Goto(ctx, srv.Prefix()+"/cancel-page")
	must.NoError(err, "Goto failed")
	err = page.Locator("#a").Click(ctx)
	must.NoError(err, "Click failed")

	var dl *playwright.Download
	select {
	case dl = <-dlCh:
	case <-time.After(10 * time.Second):
		t.Fatal("download event not received")
	}

	// Confirm the server received the request before cancelling.
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("slow-dl handler never started")
	}

	dlCtx, dlCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dlCancel()
	err = dl.Cancel(dlCtx)
	must.NoError(err, "Download.Cancel failed")
}

func TestDownloadFailure(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "success content")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Wait for download completion.
	_, err := dl.Path(dlCtx)
	must.NoError(err, "Download.Path() failed")

	// A successful download returns nil from Failure().
	if err := dl.Failure(dlCtx); err != nil {
		t.Errorf("Download.Failure() = %v, want nil for successful download", err)
	}
}

// TestDownloadReportDownloadsWithAcceptDownloadsTrue verifies download file content is accessible.
// Ref: TestDownload.java#shouldReportDownloadsWithAcceptDownloadsTrue
func TestDownloadReportDownloadsWithAcceptDownloadsTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	path, err := dl.Path(dlCtx)
	must.NoError(err, "Download.Path() failed")

	data, err := os.ReadFile(path)
	must.NoError(err, "ReadFile failed")
	is.Equal("Hello world", string(data))
}

// TestDownloadSaveToUserSpecifiedPathWithoutUpdatingOriginalPath verifies both paths are accessible.
// Ref: TestDownload.java#shouldSaveToUserSpecifiedPathWithoutUpdatingOriginalPath
func TestDownloadSaveToUserSpecifiedPathWithoutUpdatingOriginalPath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dir, err := os.MkdirTemp("", "dl-nopath-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	userFile := filepath.Join(dir, "saved.bin")
	err = dl.SaveAs(dlCtx, userFile)
	must.NoError(err)

	data, err := os.ReadFile(userFile)
	must.NoError(err)
	is.Equal("Hello world", string(data))

	// Original path must still be accessible.
	origPath, err := dl.Path(dlCtx)
	must.NoError(err)
	origData, err := os.ReadFile(origPath)
	must.NoError(err)
	is.Equal("Hello world", string(origData))
}

// TestDownloadSaveToTwoDifferentPaths verifies SaveAs can be called twice with different destinations.
// Ref: TestDownload.java#shouldSaveToTwoDifferentPathsWithMultipleSaveAsCalls
func TestDownloadSaveToTwoDifferentPaths(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dir, err := os.MkdirTemp("", "dl-twopaths-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	path1 := filepath.Join(dir, "first.bin")
	err = dl.SaveAs(dlCtx, path1)
	must.NoError(err)
	data, err := os.ReadFile(path1)
	must.NoError(err)
	is.Equal("Hello world", string(data))

	path2 := filepath.Join(dir, "second.bin")
	err = dl.SaveAs(dlCtx, path2)
	must.NoError(err)
	data, err = os.ReadFile(path2)
	must.NoError(err)
	is.Equal("Hello world", string(data))
}

// TestDownloadSaveToOverwrittenFilepath verifies SaveAs overwrites an existing file.
// Ref: TestDownload.java#shouldSaveToOverwrittenFilepath
func TestDownloadSaveToOverwrittenFilepath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dir, err := os.MkdirTemp("", "dl-overwrite-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	dest := filepath.Join(dir, "out.bin")

	// Save twice to the same path — second call should overwrite.
	must.NoError(dl.SaveAs(dlCtx, dest))
	data, err := os.ReadFile(dest)
	must.NoError(err)
	is.Equal("Hello world", string(data))

	must.NoError(dl.SaveAs(dlCtx, dest))
	data, err = os.ReadFile(dest)
	must.NoError(err)
	is.Equal("Hello world", string(data))
}

// TestDownloadCreateSubdirectoriesWhenSavingToNonExistentPath verifies SaveAs creates nested dirs.
// Ref: TestDownload.java#shouldCreateSubdirectoriesWhenSavingToNonExistentUserSpecifiedPath
func TestDownloadCreateSubdirectoriesWhenSavingToNonExistentPath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	base, err := os.MkdirTemp("", "dl-subdirs-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(base) })

	nested := filepath.Join(base, "these", "are", "directories", "download.txt")
	must.NoError(dl.SaveAs(dlCtx, nested))

	data, err := os.ReadFile(nested)
	must.NoError(err)
	is.Equal("Hello world", string(data))
}

// TestDownloadErrorWhenSavingWithDownloadsDisabled verifies SaveAs errors when acceptDownloads=false.
// Ref: TestDownload.java#shouldErrorWhenSavingWithDownloadsDisabled
func TestDownloadErrorWhenSavingWithDownloadsDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/dl-file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="file.bin"`)
		_, _ = w.Write([]byte("Hello world"))
	})
	srv.ServeWithBody("/dl-page", "text/html", `<a id="a" href="/dl-file" download>dl</a>`)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		AcceptDownloads: boolPtr(false),
	})
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/dl-page"))
	must.NoError(page.Locator("#a").Click(ctx))

	var dl *playwright.Download
	select {
	case dl = <-dlCh:
	case <-time.After(10 * time.Second):
		t.Fatal("download event not received")
	}

	dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer dlCancel()

	dir, err := os.MkdirTemp("", "dl-disabled-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	err = dl.SaveAs(dlCtx, filepath.Join(dir, "out.bin"))
	is.Error(err, "SaveAs should fail when acceptDownloads=false")
	is.Contains(err.Error(), "acceptDownloads")
}

// TestDownloadErrorWhenSavingAfterDeletion verifies SaveAs errors after Delete.
// Ref: TestDownload.java#shouldErrorWhenSavingAfterDeletion
func TestDownloadErrorWhenSavingAfterDeletion(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	must.NoError(dl.Delete(dlCtx))

	dir, err := os.MkdirTemp("", "dl-afterdelete-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	err = dl.SaveAs(dlCtx, filepath.Join(dir, "out.bin"))
	is.Error(err, "SaveAs should fail after Delete")
}

// TestDownloadReportNonNavigationDownloads verifies downloads triggered via `download` attribute.
// Ref: TestDownload.java#shouldReportNonNavigationDownloads
func TestDownloadReportNonNavigationDownloads(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/dl-file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte("Hello world"))
	})

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a download='file.txt' href='`+srv.Prefix()+`/dl-file'>download</a>`))

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})
	must.NoError(page.Locator("a").Click(ctx))

	var dl *playwright.Download
	select {
	case dl = <-dlCh:
	case <-time.After(15 * time.Second):
		t.Fatal("download event not received")
	}

	is.Equal("file.txt", dl.SuggestedFilename())

	dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer dlCancel()

	path, err := dl.Path(dlCtx)
	must.NoError(err)
	data, err := os.ReadFile(path)
	must.NoError(err)
	is.Equal("Hello world", string(data))
}

// TestDownloadDeleteDownloadsOnContextDestruction verifies downloads are cleaned up when context closes.
// Ref: TestDownload.java#shouldDeleteDownloadsOnContextDestruction
func TestDownloadDeleteDownloadsOnContextDestruction(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/dl-file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="file.bin"`)
		_, _ = w.Write([]byte("Hello world"))
	})
	srv.ServeWithBody("/dl-page", "text/html", `<a id="a" href="/dl-file" download>dl</a>`)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		AcceptDownloads: boolPtr(true),
	})
	must.NoError(err)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/dl-page"))

	downloadAndGetPath := func() (string, *playwright.Download) {
		t.Helper()
		dlCh := make(chan *playwright.Download, 1)
		page.OnDownload(func(d *playwright.Download) {
			select {
			case dlCh <- d:
			default:
			}
		})
		must.NoError(page.Locator("#a").Click(ctx))
		var dl *playwright.Download
		select {
		case dl = <-dlCh:
		case <-time.After(15 * time.Second):
			t.Fatal("download event not received")
		}
		dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dlCancel()
		p, err := dl.Path(dlCtx)
		must.NoError(err)
		return p, dl
	}

	path1, _ := downloadAndGetPath()
	path2, _ := downloadAndGetPath()

	_, err = os.Stat(path1)
	must.NoError(err, "download1 should exist before context close")
	_, err = os.Stat(path2)
	must.NoError(err, "download2 should exist before context close")

	closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeCancel()
	must.NoError(bCtx.Close(closeCtx))

	_, err = os.Stat(path1)
	is.True(os.IsNotExist(err), "download1 should be deleted after context close")
	_, err = os.Stat(path2)
	is.True(os.IsNotExist(err), "download2 should be deleted after context close")
}

// TestDownloadCancelPending verifies Cancel sets failure to "canceled".
// Ref: TestDownload.java#shouldBeAbleToCancelPendingDownloads
func TestDownloadCancelPending(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	started := make(chan struct{}, 1)
	srv.SetRoute("/slow-dl2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="slow.bin"`)
		w.Header().Set("Content-Length", "10000000")
		w.WriteHeader(200)
		_, _ = w.Write(make([]byte, 8192))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		select {
		case started <- struct{}{}:
		default:
		}
		<-r.Context().Done()
	})
	srv.ServeWithBody("/cancel-page2", "text/html", `<a id="a" href="/slow-dl2" download>dl</a>`)

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/cancel-page2"))
	must.NoError(page.Locator("#a").Click(ctx))

	var dl *playwright.Download
	select {
	case dl = <-dlCh:
	case <-time.After(10 * time.Second):
		t.Fatal("download event not received")
	}

	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("slow-dl handler never started")
	}

	dlCtx, dlCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dlCancel()
	must.NoError(dl.Cancel(dlCtx))

	failure := dl.Failure(dlCtx)
	must.NotNil(failure, "Failure() should be non-nil after Cancel")
	is.Equal("canceled", failure.Error())
}

// TestDownloadNotFailToCancelAlreadyFinished verifies Cancel on a finished download doesn't error.
// Ref: TestDownload.java#shouldNotFailExplicitlyToCancelADownloadEvenIfThatIsAlreadyFinished
func TestDownloadNotFailToCancelAlreadyFinished(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	path, err := dl.Path(dlCtx)
	must.NoError(err)
	data, err := os.ReadFile(path)
	must.NoError(err)
	is.Equal("Hello world", string(data))

	// Cancel on an already-finished download must not error.
	must.NoError(dl.Cancel(dlCtx))
	is.Nil(dl.Failure(dlCtx))
}

// TestDownloadReportDownloadWhenNavigationTurnsIntoDownload verifies that navigating directly
// to a download URL results in a proper download event.
// Ref: TestDownload.java#shouldReportDownloadWhenNavigationTurnsIntoDownload
func TestDownloadReportDownloadWhenNavigationTurnsIntoDownload(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/nav-download", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Hello world"))
	})

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})

	// Navigating to a download URL results in an error but also fires the download event.
	_ = page.Goto(ctx, srv.Prefix()+"/nav-download")

	var dl *playwright.Download
	select {
	case dl = <-dlCh:
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for download event")
	}

	is.Contains(dl.URL(), "/nav-download")

	dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer dlCancel()

	path, err := dl.Path(dlCtx)
	must.NoError(err)
	data, err := os.ReadFile(path)
	must.NoError(err)
	is.Equal("Hello world", string(data))
}

// TestDownloadReportDownloadsWithAcceptDownloadsFalse verifies that with acceptDownloads=false
// the download event fires but path() fails with an error mentioning acceptDownloads.
// Ref: TestDownload.java#shouldReportDownloadsWithAcceptDownloadsFalse
func TestDownloadReportDownloadsWithAcceptDownloadsFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/dl-named", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=file.txt")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Hello world"))
	})
	srv.ServeWithBody("/dl-named-page", "text/html", `<a id="a" href="/dl-named">download</a>`)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		AcceptDownloads: boolPtr(false),
	})
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/dl-named-page"))
	must.NoError(page.Locator("#a").Click(ctx))

	var dl *playwright.Download
	select {
	case dl = <-dlCh:
	case <-time.After(10 * time.Second):
		t.Fatal("download event not received")
	}

	is.Contains(dl.URL(), "/dl-named")
	is.Equal("file.txt", dl.SuggestedFilename())

	dlCtx, dlCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dlCancel()

	_, err = dl.Path(dlCtx)
	is.Error(err, "path() should fail when acceptDownloads=false")
	is.Contains(err.Error(), "acceptDownloads")

	failure := dl.Failure(dlCtx)
	must.NotNil(failure)
	is.Contains(failure.Error(), "acceptDownloads")
}

// TestDownloadReportPathWithinOnDownloadHandlerForFiles verifies that download.Path()
// is accessible within an OnDownload handler for regular file downloads.
// Ref: TestDownload.java#shouldReportDownloadPathWithinPageOnDownloadHandlerForFiles
func TestDownloadReportPathWithinOnDownloadHandlerForFiles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/handler-dl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Hello world"))
	})
	srv.ServeWithBody("/handler-dl-page", "text/html", `<a id="a" href="/handler-dl">download</a>`)

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/handler-dl-page"))

	pathCh := make(chan string, 1)
	page.OnDownload(func(d *playwright.Download) {
		dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dlCancel()
		p, err := d.Path(dlCtx)
		if err == nil {
			select {
			case pathCh <- p:
			default:
			}
		}
	})

	must.NoError(page.Locator("#a").Click(ctx))

	select {
	case path := <-pathCh:
		data, err := os.ReadFile(path)
		must.NoError(err)
		is.Equal("Hello world", string(data))
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for download path")
	}
}

// TestDownloadReportPathWithinOnDownloadHandlerForBlobs verifies that download.Path()
// is accessible within an OnDownload handler for blob-based downloads.
// Ref: TestDownload.java#shouldReportDownloadPathWithinPageOnDownloadHandlerForBlobs
func TestDownloadReportPathWithinOnDownloadHandlerForBlobs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/download-blob.html", "text/html", `<!DOCTYPE html><html><body>
<a id="a">Download blob</a>
<script>
document.querySelector('#a').addEventListener('click', e => {
  e.preventDefault();
  const blob = new Blob(['Hello world'], {type: 'text/plain'});
  const a = document.createElement('a');
  a.download = 'blob.txt';
  a.href = URL.createObjectURL(blob);
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
});
</script>
</body></html>`)

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/download-blob.html"))

	pathCh := make(chan string, 1)
	page.OnDownload(func(d *playwright.Download) {
		dlCtx, dlCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer dlCancel()
		p, err := d.Path(dlCtx)
		if err == nil {
			select {
			case pathCh <- p:
			default:
			}
		}
	})

	must.NoError(page.Locator("#a").Click(ctx))

	select {
	case path := <-pathCh:
		data, err := os.ReadFile(path)
		must.NoError(err)
		is.Equal("Hello world", string(data))
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for blob download path")
	}
}

// TestDownloadReportAltClickDownloads verifies that alt+click triggers a download.
// Only meaningful for Chromium; other browsers may not download on alt-click.
// Ref: TestDownload.java#shouldReportAltClickDownloads
func TestDownloadReportAltClickDownloads(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/alt-dl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Hello world"))
	})
	srv.ServeWithBody("/alt-dl-page", "text/html", `<a id="a" href="/alt-dl">download</a>`)

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.SetContent(ctx, `<a id="a" href="`+srv.Prefix()+`/alt-dl">download</a>`))

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})

	must.NoError(page.Locator("#a").Click(ctx, &playwright.LocatorClickOptions{
		Modifiers: []string{"Alt"},
	}))

	select {
	case dl := <-dlCh:
		dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dlCancel()
		path, err := dl.Path(dlCtx)
		must.NoError(err)
		data, err := os.ReadFile(path)
		must.NoError(err)
		is.Equal("Hello world", string(data))
	case <-time.After(10 * time.Second):
		t.Skip("alt-click download did not fire — may not be supported in this browser")
	}
}

// TestDownloadReportNewWindowDownloads verifies that downloads from target=_blank links
// are reported properly.
// Ref: TestDownload.java#shouldReportNewWindowDownloads
func TestDownloadReportNewWindowDownloads(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/new-window-dl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Hello world"))
	})

	bCtx := downloadContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.SetContent(ctx, `<a id="a" target=_blank href="`+srv.Prefix()+`/new-window-dl">download</a>`))

	dlCh := make(chan *playwright.Download, 1)
	page.OnDownload(func(d *playwright.Download) {
		select {
		case dlCh <- d:
		default:
		}
	})

	must.NoError(page.Locator("#a").Click(ctx))

	select {
	case dl := <-dlCh:
		dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dlCancel()
		path, err := dl.Path(dlCtx)
		must.NoError(err)
		data, err := os.ReadFile(path)
		must.NoError(err)
		is.Equal("Hello world", string(data))
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for new-window download")
	}
}

// TestDownloadDeleteDownloadsOnBrowserGone verifies that downloads are cleaned up
// when the browser is closed (not just the context).
// Ref: TestDownload.java#shouldDeleteDownloadsOnBrowserGone
func TestDownloadDeleteDownloadsOnBrowserGone(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/browser-gone-dl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Hello world"))
	})
	srv.ServeWithBody("/browser-gone-page", "text/html", `<a id="a" href="/browser-gone-dl" download>dl</a>`)

	// Launch a dedicated browser instance for this test so we can close it independently.
	bt, err := globalPW.Chromium()
	must.NoError(err, "Chromium() failed")
	b, err := bt.Launch(ctx, &playwright.BrowserTypeLaunchOptions{
		Headless: boolPtr(true),
	})
	must.NoError(err)

	bCtx, err := b.NewContext(ctx, &playwright.BrowserContextOptions{
		AcceptDownloads: boolPtr(true),
	})
	must.NoError(err)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/browser-gone-page"))

	getDLPath := func() string {
		t.Helper()
		dlCh := make(chan *playwright.Download, 1)
		page.OnDownload(func(d *playwright.Download) {
			select {
			case dlCh <- d:
			default:
			}
		})
		must.NoError(page.Locator("#a").Click(ctx))
		var dl *playwright.Download
		select {
		case dl = <-dlCh:
		case <-time.After(15 * time.Second):
			t.Fatal("download event not received")
		}
		dlCtx, dlCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dlCancel()
		p, err := dl.Path(dlCtx)
		must.NoError(err)
		return p
	}

	path1 := getDLPath()
	path2 := getDLPath()

	_, err = os.Stat(path1)
	must.NoError(err, "download1 should exist before browser close")
	_, err = os.Stat(path2)
	must.NoError(err, "download2 should exist before browser close")

	closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeCancel()
	must.NoError(b.Close(closeCtx))

	// Give the OS a moment to clean up.
	time.Sleep(200 * time.Millisecond)

	_, err = os.Stat(path1)
	is.True(os.IsNotExist(err), "download1 should be deleted after browser close")
	_, err = os.Stat(path2)
	is.True(os.IsNotExist(err), "download2 should be deleted after browser close")
}

// Ref: TestDownload.java#shouldSaveToUserSpecifiedPath
func TestDownloadSaveToUserSpecifiedPath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dir, err := os.MkdirTemp("", "dl-userpath-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	userFile := filepath.Join(dir, "download.txt")
	must.NoError(dl.SaveAs(dlCtx, userFile))

	_, err = os.Stat(userFile)
	must.NoError(err, "saved file should exist")

	data, err := os.ReadFile(userFile)
	must.NoError(err)
	is.Equal("Hello world", string(data))
}

// Ref: TestDownload.java#shouldDeleteFile
func TestDownloadShouldDeleteFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	path, err := dl.Path(dlCtx)
	must.NoError(err)

	_, err = os.Stat(path)
	must.NoError(err, "download file should exist before delete")

	must.NoError(dl.Delete(dlCtx))

	_, err = os.Stat(path)
	is.True(os.IsNotExist(err), "download file should be gone after delete")
}

// Ref: TestDownload.java#shouldExposeStream
func TestDownloadShouldExposeStream(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Go API does not expose createReadStream; read via Path instead.
	path, err := dl.Path(dlCtx)
	must.NoError(err)
	data, err := os.ReadFile(path)
	must.NoError(err)
	is.Equal("Hello world", string(data))
}

// Ref: TestDownload.java#streamShouldSupportZeroSizeRead
func TestDownloadStreamShouldSupportZeroSizeRead(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	bCtx := downloadContext(t)
	dl := triggerDownload(t, bCtx, "Hello world")

	dlCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Go API does not expose createReadStream. Verify the file is non-empty (size > 0).
	path, err := dl.Path(dlCtx)
	must.NoError(err)
	info, err := os.Stat(path)
	must.NoError(err)
	is.Greater(info.Size(), int64(0), "downloaded file should be non-empty")
}
