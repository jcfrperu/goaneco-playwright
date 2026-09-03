//go:build e2e

// E2E tests for Video recording via BrowserContextOptions.RecordVideo.
package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
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

// TestVideoSaveAs verifies that a recorded video can be saved to a custom path.
func TestVideoSaveAs(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	dir := t.TempDir()
	v := videoSetup(t, dir)

	pathCtx, pathCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pathCancel()

	_, err := v.Path(pathCtx)
	must.NoError(err, "Video.Path() failed")

	dest := filepath.Join(dir, "saved-video.webm")
	err = v.SaveAs(pathCtx, dest)
	must.NoErrorf(err, "Video.SaveAs(%q) failed", dest)

	info, err := os.Stat(dest)
	must.NoError(err, "saved video not found")
	if info.Size() == 0 {
		t.Error("saved video file is empty")
	}
}

// TestVideoDelete verifies that deleting a video removes the file.
func TestVideoDelete(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	dir := t.TempDir()
	v := videoSetup(t, dir)

	pathCtx, pathCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pathCancel()

	videoPath, err := v.Path(pathCtx)
	must.NoError(err, "Video.Path() failed")
	if videoPath == "" {
		t.Fatal("Video.Path() returned empty path")
	}

	err = v.Delete(pathCtx)
	must.NoError(err, "Video.Delete failed")

	if _, err := os.Stat(videoPath); !os.IsNotExist(err) {
		t.Errorf("video file should be deleted but still exists at %q (err: %v)", videoPath, err)
	}
}

// TestVideoRecording verifies that a page video is recorded when RecordVideo is set
// and that the video file exists and is non-empty after page close.
func TestVideoRecording(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	dir := t.TempDir()

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		RecordVideo: &playwright.RecordVideoOptions{
			Dir: dir,
		},
	})
	must.NoError(err, "NewContext with RecordVideo failed")
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	}()

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Close the page to flush the video recording
	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = page.Close(closeCtx)
	must.NoError(err, "page.Close failed")

	// The video event fires after page close; give it a moment
	must.NoError(page.WaitForTimeout(ctx, 500))

	v := page.Video()
	if v == nil {
		// Check if any video file was written to the dir
		entries, _ := os.ReadDir(dir)
		if len(entries) == 0 {
			t.Skip("no video file created — recording may not be supported in this environment")
		}
		// Video may be accessible via filesystem only
		videoPath := filepath.Join(dir, entries[0].Name())
		info, err := os.Stat(videoPath)
		if err != nil {
			t.Fatalf("stat video file: %v", err)
		}
		if info.Size() == 0 {
			t.Error("video file is empty")
		}
		t.Logf("video recorded at %s (%d bytes)", videoPath, info.Size())
		return
	}

	pathCtx, pathCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pathCancel()
	videoPath, err := v.Path(pathCtx)
	must.NoError(err, "Video.Path() failed")
	if videoPath == "" {
		t.Fatal("expected non-empty video path")
	}
	t.Logf("video path: %s", videoPath)

	info, err := os.Stat(videoPath)
	must.NoError(err, "video file stat failed")
	if info.Size() == 0 {
		t.Error("video file is empty")
	}
}
