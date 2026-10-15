//go:build e2e

// E2E tests for: TestBrowserTypeConnect.java.
package e2e

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Ref: TestBrowserTypeConnect.java#shouldSaveDownload  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldSaveDownload(t *testing.T) {
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

// Ref: TestBrowserTypeConnect.java#shouldBeAbleToReconnectToABrowser  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldBeAbleToReconnectToABrowser(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldSupportSlowMo  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldSupportSlowMo(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldBeAbleToConnectTwoBrowsersAtTheSameTime  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldBeAbleToConnectTwoBrowsersAtTheSameTime(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldSendExtraHeadersWithConnectRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldSendExtraHeadersWithConnectRequest(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#disconnectedEventShouldBeEmittedWhenBrowserIsClosedOrServerIsClosed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectDisconnectedEventShouldBeEmittedWhenBrowserIsClosedOrServerIsClosed(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#disconnectedEventShouldHaveBrowserAsArgument  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectDisconnectedEventShouldHaveBrowserAsArgument(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldSetTheBrowserConnectedState  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldSetTheBrowserConnectedState(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldThrowWhenUsedAfterIsConnectedReturnsFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldThrowWhenUsedAfterIsConnectedReturnsFalse(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldThrowWhenCallingWaitForNavigationAfterDisconnect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldThrowWhenCallingWaitForNavigationAfterDisconnect(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldRejectNavigationWhenBrowserCloses  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldRejectNavigationWhenBrowserCloses(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldEmitCloseEventsOnPagesAndContexts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldEmitCloseEventsOnPagesAndContexts(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldRespectSelectors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldRespectSelectors(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldNotThrowOnCloseAfterDisconnect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldNotThrowOnCloseAfterDisconnect(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldSaveAsVideosFromRemoteBrowser  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldSaveAsVideosFromRemoteBrowser(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldErrorWhenSavingDownloadAfterDeletion  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldErrorWhenSavingDownloadAfterDeletion(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldSupportTracingOverWebSocket  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldSupportTracingOverWebSocket(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldRecordTraceWithSources  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldRecordTraceWithSources(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldFulfillWithGlobalFetchResult  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldFulfillWithGlobalFetchResult(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#shouldUploadLargeFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectShouldUploadLargeFile(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}

// Ref: TestBrowserTypeConnect.java#setInputFilesShouldPreserveLastModifiedTimestamp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserTypeConnect.java
func TestBrowserTypeConnectSetInputFilesShouldPreserveLastModifiedTimestamp(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect (only ConnectOverCDP exists)")
}
