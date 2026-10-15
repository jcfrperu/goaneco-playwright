//go:build e2e

// E2E tests for: TestTracing.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"archive/zip"
	"bytes"
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestTracing.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingStartStop(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/trace-page", "text/html", `<p>tracing test</p>`)

	tr := bCtx.Tracing()
	must.NotNil(tr, "BrowserContext.Tracing() returned nil")

	err := tr.Start(ctx)
	must.NoError(err, "Tracing.Start failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	err = page.Goto(ctx, srv.Prefix()+"/trace-page")
	must.NoError(err, "Goto failed")

	err = tr.Stop(ctx)
	must.NoError(err, "Tracing.Stop failed")
}

// Extra: no equivalent method in playwright-java; related class: TestTracing.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingStopWithoutStart(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	tr := bCtx.Tracing()
	if tr == nil {
		t.Skip("Tracing() returned nil — context may not support tracing")
	}
	// Stopping without start should not error.
	err := tr.Stop(ctx)
	must.NoError(err, "Tracing.Stop (no prior Start) failed")
}

// Extra: no equivalent method in playwright-java; related class: TestTracing.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingStartChunkStopChunkSave(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/trace-chunk", "text/html", `<p>chunk tracing</p>`)

	tr := bCtx.Tracing()
	must.NotNil(tr, "BrowserContext.Tracing() returned nil")

	name := "my-trace"
	err := tr.Start(ctx, &playwright.TracingStartOptions{Name: &name})
	must.NoError(err, "Tracing.Start failed")

	traceName, err := tr.StartChunk(ctx)
	must.NoError(err, "Tracing.StartChunk failed")
	t.Logf("StartChunk returned traceName: %q", traceName)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	err = page.Goto(ctx, srv.Prefix()+"/trace-chunk")
	must.NoError(err, "Goto failed")

	dir, err := os.MkdirTemp("", "trace-*")
	must.NoError(err, "MkdirTemp failed")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	tracePath := filepath.Join(dir, "trace.zip")
	err = tr.StopChunk(ctx, tracePath)
	must.NoError(err, "Tracing.StopChunk failed")

	info, err := os.Stat(tracePath)
	must.NoError(err, "trace file not found")
	if info.Size() == 0 {
		t.Error("trace file is empty")
	}

	// Verify it is a valid ZIP archive.
	data, err := os.ReadFile(tracePath)
	must.NoError(err, "failed to read trace file")
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	must.NoError(err, "trace file is not a valid ZIP")
	if len(zr.File) == 0 {
		t.Error("trace ZIP has no files inside")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestTracing.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingStartChunkDiscard(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	tr := bCtx.Tracing()
	must.NotNil(tr, "BrowserContext.Tracing() returned nil")
	err := tr.Start(ctx)
	must.NoError(err, "Tracing.Start failed")
	_, err = tr.StartChunk(ctx)
	must.NoError(err, "Tracing.StartChunk failed")
	// StopChunk with empty path should discard without error.
	err = tr.StopChunk(ctx, "")
	must.NoError(err, "Tracing.StopChunk (discard) failed")
}

// Extra: no equivalent method in playwright-java; related class: TestTracing.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingGetObject(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	// Create a fresh context to verify Tracing() is non-nil.
	closeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() { _ = bCtx.Close(closeCtx) })

	tr := bCtx.Tracing()
	must.NotNil(tr, "expected non-nil Tracing object")
}

// TestTracingNotFailWhenSourcesSetToFalse verifies tracing doesn't fail with Sources=false.
// Ref: TestTracing.java#shouldNotFailWhenSourcesSetExplicitlyToFalse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldNotFailWhenSourcesSetExplicitlyToFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/trace-nosrc", "text/html", `<p>no sources</p>`)

	tr := bCtx.Tracing()
	must.NotNil(tr)

	sources := false
	err := tr.Start(ctx, &playwright.TracingStartOptions{Sources: &sources})
	must.NoError(err, "Tracing.Start with Sources=false failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/trace-nosrc"))

	dir, err := os.MkdirTemp("", "trace-nosrc-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	tracePath := filepath.Join(dir, "trace.zip")
	must.NoError(tr.StopChunk(ctx, tracePath), "StopChunk with Sources=false failed")

	info, err := os.Stat(tracePath)
	must.NoError(err)
	is.Greater(info.Size(), int64(0))
}

// TestTracingShouldCollectTwoTraces verifies two separate trace chunks can be saved independently.
// Ref: TestTracing.java#shouldCollectTwoTraces  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldCollectTwoTraces(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/trace1", "text/html", `<p>trace1</p>`)
	srv.ServeWithBody("/trace2", "text/html", `<p>trace2</p>`)

	tr := bCtx.Tracing()
	must.NotNil(tr)

	err := tr.Start(ctx)
	must.NoError(err, "Tracing.Start failed")

	dir, err := os.MkdirTemp("", "trace-two-*")
	must.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	// First chunk.
	_, err = tr.StartChunk(ctx)
	must.NoError(err, "StartChunk 1 failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/trace1"))

	trace1Path := filepath.Join(dir, "trace1.zip")
	must.NoError(tr.StopChunk(ctx, trace1Path), "StopChunk 1 failed")

	// Second chunk.
	_, err = tr.StartChunk(ctx)
	must.NoError(err, "StartChunk 2 failed")

	must.NoError(page.Goto(ctx, srv.Prefix()+"/trace2"))

	trace2Path := filepath.Join(dir, "trace2.zip")
	must.NoError(tr.StopChunk(ctx, trace2Path), "StopChunk 2 failed")

	// Verify both ZIP files exist and are non-empty.
	for _, p := range []string{trace1Path, trace2Path} {
		info, err := os.Stat(p)
		must.NoError(err, "trace file not found: %s", p)
		is.Greater(info.Size(), int64(0), "trace file is empty: %s", p)

		data, err := os.ReadFile(p)
		must.NoError(err)
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		must.NoError(err, "not a valid ZIP: %s", p)
		is.Greater(len(zr.File), 0, "ZIP has no files: %s", p)
	}
}

// Ref: TestTracing.java#shouldCollectTrace1  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldCollectTrace1(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldWorkWithMultipleChunks  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldWorkWithMultipleChunks(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldCollectSources  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldCollectSources(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldRespectTracesDirAndName  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldRespectTracesDirAndName(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldShowWaitForLoadState  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldShowWaitForLoadState(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldRecoverTracingAfterAFailedStop  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldRecoverTracingAfterAFailedStop(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldCollectActionScreenshots  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldCollectActionScreenshots(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldCollectAriaSnapshots  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldCollectAriaSnapshots(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldTraceVariousAPIs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldTraceVariousAPIs(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#shouldNotRecordNetworkActions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldNotRecordNetworkActions(t *testing.T) {
	t.Skip("not implemented: TraceViewerPage helper / Tracing.Stop path option / Tracing.Group")
}

// Ref: TestTracing.java#canCallTracingGroupGroupEndAtAnyTimeAndAutoClose  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingCanCallTracingGroupGroupEndAtAnyTimeAndAutoClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	tracing := bCtx.Tracing()
	must.NotNil(tracing)

	must.NoError(tracing.Group(ctx, "ignored"))
	must.NoError(tracing.GroupEnd(ctx))
	must.NoError(tracing.Group(ctx, "ignored2"))
	must.NoError(tracing.Start(ctx))
	must.NoError(tracing.Group(ctx, "actual"))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	srv := testserver.New(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	tracePath := filepath.Join(t.TempDir(), "trace.zip")
	must.NoError(tracing.StopChunk(ctx, tracePath))
	must.NoError(tracing.Group(ctx, "ignored3"))
	must.NoError(tracing.GroupEnd(ctx))
	must.NoError(tracing.GroupEnd(ctx))
	must.NoError(tracing.GroupEnd(ctx))
	must.NoError(tracing.Stop(ctx))

	trace := tracingArchiveTraceText(t, tracePath)
	must.Contains(trace, "actual")
	must.Contains(trace, "Navigate")
}

// Ref: TestTracing.java#traceGroupGroupEnd  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingTraceGroupGroupEnd(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	tracing := bCtx.Tracing()
	must.NotNil(tracing)
	must.NoError(tracing.Start(ctx))
	must.NoError(tracing.Group(ctx, "outer group"))

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "data:text/html,<!DOCTYPE html><body><div>Hello world</div></body>"))
	line, column := 17, 1
	must.NoError(tracing.Group(ctx, "inner group 1", &playwright.TracingGroupOptions{
		Location: &playwright.TracingGroupLocation{File: "foo.java", Line: &line, Column: &column},
	}))
	must.NoError(page.Locator("body").Click(ctx))
	must.NoError(tracing.GroupEnd(ctx))
	must.NoError(tracing.Group(ctx, "inner group 2"))
	visible, err := page.Locator("text=Hello").IsVisible(ctx)
	must.NoError(err)
	must.True(visible)
	must.NoError(tracing.GroupEnd(ctx))
	must.NoError(tracing.GroupEnd(ctx))

	tracePath := filepath.Join(t.TempDir(), "trace.zip")
	must.NoError(tracing.StopChunk(ctx, tracePath))
	must.NoError(tracing.Stop(ctx))

	trace := tracingArchiveTraceText(t, tracePath)
	for _, expected := range []string{"outer group", "inner group 1", "inner group 2", "Click", "foo.java"} {
		must.Contains(trace, expected)
	}
}

func tracingArchiveTraceText(t *testing.T, path string) string {
	t.Helper()
	must := require.New(t)
	data, err := os.ReadFile(path)
	must.NoError(err)
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	must.NoError(err)
	var contents strings.Builder
	for _, file := range archive.File {
		if !strings.HasSuffix(file.Name, ".trace") {
			continue
		}
		reader, err := file.Open()
		must.NoError(err)
		fileData, err := io.ReadAll(reader)
		must.NoError(err)
		must.NoError(reader.Close())
		contents.Write(fileData)
	}
	return contents.String()
}

// Ref: TestTracing.java#shouldRecordHarWithStartHarStopHar  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingShouldRecordHarWithStartHarStopHar(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	tracing := bCtx.Tracing()
	must.NotNil(tracing)
	harPath := filepath.Join(t.TempDir(), "tracing.har")
	mode := "minimal"
	must.NoError(tracing.StartHar(ctx, harPath, &playwright.TracingStartHarOptions{Mode: &mode}))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/one-style.html"))
	must.NoError(tracing.StopHar(ctx))
	content, err := os.ReadFile(harPath)
	must.NoError(err)
	must.Contains(string(content), `"log"`)
	must.Contains(string(content), "/one-style.html")
}

// Ref: TestTracing.java#startShouldReturnADisposableThatDiscardsTheTrace  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingStartShouldReturnADisposableThatDiscardsTheTrace(t *testing.T) {
	t.Skip("not implemented: Tracing.Start returning Disposable")
}

// Ref: TestTracing.java#startChunkShouldReturnADisposableThatDiscardsTheChunk  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracing.java
func TestTracingStartChunkShouldReturnADisposableThatDiscardsTheChunk(t *testing.T) {
	t.Skip("not implemented: Tracing.StartChunk returning Disposable")
}
