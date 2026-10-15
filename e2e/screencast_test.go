//go:build e2e

// E2E tests for: TestScreencast.java.
package e2e

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestScreencastShouldSaveAsVideo verifies that a recorded video can be saved to a custom path.
// Ref: TestScreencast.java#shouldSaveAsVideo  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastShouldSaveAsVideo(t *testing.T) {
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

// TestScreencastShouldDeleteVideo verifies that deleting a video removes the file.
// Ref: TestScreencast.java#shouldDeleteVideo  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastShouldDeleteVideo(t *testing.T) {
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

// Ref: TestScreencast.java#shouldExposeVideoPath  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastShouldExposeVideoPath(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#shouldWaitForVideoFinishWhenPageIsClosed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastShouldWaitForVideoFinishWhenPageIsClosed(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastStartShouldDeliverFramesViaOnFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastStartShouldDeliverFramesViaOnFrame(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastStartShouldScaleFramesToFitSize  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastStartShouldScaleFramesToFitSize(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#onFrameShouldReceiveViewportSizeAndTimestamp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastOnFrameShouldReceiveViewportSizeAndTimestamp(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastStartShouldRecordVideoWithTheRequestedFps  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastStartShouldRecordVideoWithTheRequestedFps(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastStartShouldThrowOnInvalidFps  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastStartShouldThrowOnInvalidFps(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastStartShouldThrowIfAlreadyStarted  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastStartShouldThrowIfAlreadyStarted(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastStartShouldRecordVideoToPath  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastStartShouldRecordVideoToPath(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastStartReturnsDisposable  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastStartReturnsDisposable(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastShowOverlay  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastShowOverlay(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastShowChapter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastShowChapter(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastHideShowOverlays  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastHideShowOverlays(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastShowAndHideActions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastShowAndHideActions(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastShowActionsShouldAcceptEveryPosition  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastShowActionsShouldAcceptEveryPosition(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastShowActionsShouldAcceptEveryCursor  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastShowActionsShouldAcceptEveryCursor(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}

// Ref: TestScreencast.java#screencastShowActionsShouldRenderStyle  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestScreencast.java
func TestScreencastScreencastShowActionsShouldRenderStyle(t *testing.T) {
	t.Skip("not implemented: Page.Screencast")
}
