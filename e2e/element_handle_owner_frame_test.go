//go:build e2e

// E2E tests for: TestElementHandleOwnerFrame.java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Ref: TestElementHandleOwnerFrame.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleOwnerFrame.java
func TestElementHandleOwnerFrameShouldWork(t *testing.T) {
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
	ownerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(frame, ownerFrame)
}

// Ref: TestElementHandleOwnerFrame.java#shouldWorkForCrossProcessIframes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleOwnerFrame.java
func TestElementHandleOwnerFrameShouldWorkForCrossProcessIframes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.CrossProcessPrefix()+"/empty.html")
	frames := page.Frames()
	must.GreaterOrEqual(len(frames), 2)
	frame := frames[1]
	jsHandle, err := frame.EvaluateHandle(ctx, "() => document.body")
	must.NoError(err)
	elementHandle := jsHandle.AsElement()
	must.NotNil(elementHandle)
	ownerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(frame, ownerFrame)
}

// Ref: TestElementHandleOwnerFrame.java#shouldWorkForDocument  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleOwnerFrame.java
func TestElementHandleOwnerFrameShouldWorkForDocument(t *testing.T) {
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
	jsHandle, err := frame.EvaluateHandle(ctx, "document")
	must.NoError(err)
	elementHandle := jsHandle.AsElement()
	must.NotNil(elementHandle)
	ownerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(frame, ownerFrame)
}

// Ref: TestElementHandleOwnerFrame.java#shouldWorkForIframeElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleOwnerFrame.java
func TestElementHandleOwnerFrameShouldWorkForIframeElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.EmptyPage())
	frame := page.MainFrame()
	jsHandle, err := frame.EvaluateHandle(ctx, "() => document.querySelector('#frame1')")
	must.NoError(err)
	elementHandle := jsHandle.AsElement()
	must.NotNil(elementHandle)
	ownerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(frame, ownerFrame)
}

// Ref: TestElementHandleOwnerFrame.java#shouldWorkForCrossFrameEvaluations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleOwnerFrame.java
func TestElementHandleOwnerFrameShouldWorkForCrossFrameEvaluations(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	attachFrame(t, ctx, page, "frame1", srv.EmptyPage())
	frame := page.MainFrame()
	jsHandle, err := frame.EvaluateHandle(ctx, "() => document.querySelector('iframe').contentWindow.document.body")
	must.NoError(err)
	elementHandle := jsHandle.AsElement()
	must.NotNil(elementHandle)
	ownerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	children := frame.ChildFrames()
	must.NotEmpty(children)
	is.Equal(children[0], ownerFrame)
}

// Ref: TestElementHandleOwnerFrame.java#shouldWorkForDetachedElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleOwnerFrame.java
func TestElementHandleOwnerFrameShouldWorkForDetachedElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	divHandle, err := page.EvaluateHandle(ctx, `() => {
		const div = document.createElement('div');
		document.body.appendChild(div);
		return div;
	}`)
	must.NoError(err)
	elementHandle := divHandle.AsElement()
	must.NotNil(elementHandle)
	ownerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(page.MainFrame(), ownerFrame)

	_, err = page.Evaluate(ctx, `() => {
		const div = document.querySelector('div');
		document.body.removeChild(div);
	}`)
	must.NoError(err)
	ownerFrame2, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(page.MainFrame(), ownerFrame2)
}

// Ref: TestElementHandleOwnerFrame.java#shouldWorkForAdoptedElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleOwnerFrame.java
func TestElementHandleOwnerFrameShouldWorkForAdoptedElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "url => window['__popup'] = window.open(url)", srv.EmptyPage())
		return err
	})
	must.NoError(err)
	must.NotNil(popup)

	divHandle, err := page.EvaluateHandle(ctx, `() => {
		const div = document.createElement('div');
		document.body.appendChild(div);
		return div;
	}`)
	must.NoError(err)
	elementHandle := divHandle.AsElement()
	must.NotNil(elementHandle)
	ownerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(page.MainFrame(), ownerFrame)

	must.NoError(popup.WaitForLoadState(ctx, "domcontentloaded"))
	_, err = page.Evaluate(ctx, `() => {
		const div = document.querySelector('div');
		window['__popup'].document.body.appendChild(div);
	}`)
	must.NoError(err)
	adoptedOwnerFrame, err := elementHandle.OwnerFrame(ctx)
	must.NoError(err)
	is.Equal(popup.MainFrame(), adoptedOwnerFrame)
}
