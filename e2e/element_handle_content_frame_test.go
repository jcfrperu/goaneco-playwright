//go:build e2e

// E2E tests for: TestElementHandleContentFrame.java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Ref: TestElementHandleContentFrame.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleContentFrame.java
func TestElementHandleContentFrameShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.EmptyPage())
	elementHandle, err := page.QuerySelector(ctx, "#frame1")
	must.NoError(err)
	must.NotNil(elementHandle)
	frame, err := elementHandle.ContentFrame(ctx)
	must.NoError(err)
	frames := page.Frames()
	must.GreaterOrEqual(len(frames), 2)
	is.Equal(frames[1], frame)
}

// Ref: TestElementHandleContentFrame.java#shouldWorkForCrossProcessIframes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleContentFrame.java
func TestElementHandleContentFrameShouldWorkForCrossProcessIframes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.CrossProcessPrefix()+"/empty.html")
	elementHandle, err := page.QuerySelector(ctx, "#frame1")
	must.NoError(err)
	must.NotNil(elementHandle)
	frame, err := elementHandle.ContentFrame(ctx)
	must.NoError(err)
	frames := page.Frames()
	must.GreaterOrEqual(len(frames), 2)
	is.Equal(frames[1], frame)
}

// Ref: TestElementHandleContentFrame.java#shouldWorkForCrossFrameEvaluations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleContentFrame.java
func TestElementHandleContentFrameShouldWorkForCrossFrameEvaluations(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.EmptyPage())
	frames := page.Frames()
	must.GreaterOrEqual(len(frames), 2)
	frame := frames[1]
	jsHandle, err := frame.EvaluateHandle(ctx, "() => window.top.document.querySelector('#frame1')")
	must.NoError(err)
	elementHandle := jsHandle.AsElement()
	must.NotNil(elementHandle)
	contentFrame, err := elementHandle.ContentFrame(ctx)
	must.NoError(err)
	is.Equal(frame, contentFrame)
}

// Ref: TestElementHandleContentFrame.java#shouldReturnNullForNonIframes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleContentFrame.java
func TestElementHandleContentFrameShouldReturnNullForNonIframes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.EmptyPage())
	frames := page.Frames()
	must.GreaterOrEqual(len(frames), 2)
	frame := frames[1]
	jsHandle, err := frame.EvaluateHandle(ctx, "() => document.body")
	must.NoError(err)
	elementHandle := jsHandle.AsElement()
	must.NotNil(elementHandle)
	contentFrame, err := elementHandle.ContentFrame(ctx)
	must.NoError(err)
	is.Nil(contentFrame)
}

// Ref: TestElementHandleContentFrame.java#shouldReturnNullForDocumentDocumentElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleContentFrame.java
func TestElementHandleContentFrameShouldReturnNullForDocumentDocumentElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.EmptyPage())
	frames := page.Frames()
	must.GreaterOrEqual(len(frames), 2)
	frame := frames[1]
	jsHandle, err := frame.EvaluateHandle(ctx, "() => document.documentElement")
	must.NoError(err)
	elementHandle := jsHandle.AsElement()
	must.NotNil(elementHandle)
	contentFrame, err := elementHandle.ContentFrame(ctx)
	must.NoError(err)
	is.Nil(contentFrame)
}
