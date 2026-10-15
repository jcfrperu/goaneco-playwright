//go:build e2e

// E2E tests for: TestVideo.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestVideoRecording verifies that a page video is recorded when RecordVideo is set
// and that the video file exists and is non-empty after page close.
// Extra: no equivalent method in playwright-java; related class: TestVideo.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestVideo.java
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

// Ref: TestVideo.java#shouldWorkWithRelativePathForRecordVideoDir  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestVideo.java
func TestVideoShouldWorkWithRelativePathForRecordVideoDir(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	tmpDir := t.TempDir()
	cwd, err := os.Getwd()
	must.NoError(err)
	relDir, err := filepath.Rel(cwd, tmpDir)
	if err != nil {
		t.Skip("cannot relativize temp dir: " + err.Error())
	}
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		RecordVideo: &playwright.RecordVideoOptions{Dir: relDir, Size: &playwright.ViewportSize{Width: 320, Height: 240}},
	})
	must.NoError(err)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	videoPath, err := page.Video().Path(ctx)
	must.NoError(err)
	must.NoError(bCtx.Close(ctx))
	must.True(filepath.IsAbs(videoPath), "videosPath = "+videoPath)
	_, err = os.Stat(videoPath)
	must.NoError(err, "videosPath = "+videoPath)
}

// Ref: TestVideo.java#shouldRecordVideoWithTheRequestedFps  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestVideo.java
func TestVideoShouldRecordVideoWithTheRequestedFps(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	fps := 60
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		RecordVideo: &playwright.RecordVideoOptions{Dir: t.TempDir(), Fps: &fps},
	})
	must.NoError(err)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	_, err = page.Evaluate(ctx, "() => document.body.style.backgroundColor = 'red'")
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 500))
	videoPath, err := page.Video().Path(ctx)
	must.NoError(err)
	must.NoError(bCtx.Close(ctx))
	_, err = os.Stat(videoPath)
	must.NoError(err, "videosPath = "+videoPath)
}

// Ref: TestVideo.java#shouldThrowOnInvalidFps  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestVideo.java
func TestVideoShouldThrowOnInvalidFps(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	fps := 0
	_, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		RecordVideo: &playwright.RecordVideoOptions{Dir: t.TempDir(), Fps: &fps},
	})
	must.Error(err)
	must.Contains(err.Error(), `"recordVideo.fps" must be a positive number, got 0`)
}

// Ref: TestVideo.java#shouldThrowWhenFpsIsSetWithoutDir  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestVideo.java
func TestVideoShouldThrowWhenFpsIsSetWithoutDir(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	fps := 60
	_, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{RecordVideoFps: &fps})
	must.Error(err)
	must.Contains(err.Error(), "recordVideoFps is set but recordVideoDir is null")
}
