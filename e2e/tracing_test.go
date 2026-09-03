//go:build e2e

package e2e

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

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

// TestTracingCollectTwoTraces verifies two separate trace chunks can be saved independently.
// Ref: TestTracing.java#shouldCollectTwoTraces
func TestTracingCollectTwoTraces(t *testing.T) {
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

// TestTracingNotFailWhenSourcesSetToFalse verifies tracing doesn't fail with Sources=false.
// Ref: TestTracing.java#shouldNotFailWhenSourcesSetToFalse
func TestTracingNotFailWhenSourcesSetToFalse(t *testing.T) {
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
