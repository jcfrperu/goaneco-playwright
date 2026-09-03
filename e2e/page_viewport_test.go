//go:build e2e

// Page.ViewportSize E2E tests.
// Migration of: TestPageViewport.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageViewportSizeReturnsCurrentSize verifies ViewportSize returns the set size.
// Ref: TestPageViewport.java#shouldReturnViewportSize
func TestPageViewportSizeReturnsCurrentSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 1024, 768))

	size, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(size)
	is.Equal(1024, size.Width)
	is.Equal(768, size.Height)
}

// TestPageViewportSizeUpdateReflectsChange verifies ViewportSize updates after SetViewportSize.
// Ref: TestPageViewport.java#shouldUpdateViewportSize
func TestPageViewportSizeUpdateReflectsChange(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))
	must.NoError(page.SetViewportSize(ctx, 1280, 720))

	size, err := page.ViewportSize(ctx)
	must.NoError(err)
	is.Equal(1280, size.Width)
	is.Equal(720, size.Height)
}

// TestPageViewportSizeMatchesBrowserWindow verifies viewport size is reflected in window.innerWidth.
// Ref: TestPageViewport.java#shouldMatchWindowInnerWidth
func TestPageViewportSizeMatchesBrowserWindow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 500, 400))

	w, err := page.Evaluate(ctx, `() => window.innerWidth`)
	must.NoError(err)
	is.Equal(float64(500), w)

	h, err := page.Evaluate(ctx, `() => window.innerHeight`)
	must.NoError(err)
	is.Equal(float64(400), h)
}

// TestPageViewportSizeDefaultIs1280x720 verifies default viewport is 1280x720.
// Ref: TestPageViewport.java#shouldHaveDefaultViewport
func TestPageViewportSizeDefaultIs1280x720(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	size, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(size)
	is.Equal(1280, size.Width)
	is.Equal(720, size.Height)
}

// TestSetViewportSizeUpdatesWidth verifies SetViewportSize changes viewport width.
// Ref: TestPageViewport.java#shouldUpdateWidth
func TestSetViewportSizeUpdatesWidth(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Equal(800, vp.Width)
}

// TestSetViewportSizeUpdatesHeight verifies SetViewportSize changes viewport height.
// Ref: TestPageViewport.java#shouldUpdateHeight
func TestSetViewportSizeUpdatesHeight(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 1280, 900))

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Equal(900, vp.Height)
}

// TestSetViewportSizeReflectedInJS verifies viewport is reflected in JS window dimensions.
// Ref: TestPageViewport.java#shouldReflectInJS
func TestSetViewportSizeReflectedInJS(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 640, 480))

	width, err := page.Evaluate(ctx, `() => window.innerWidth`)
	must.NoError(err)
	is.Equal(float64(640), width)
}

// TestSetViewportSizeCanResizeMultipleTimes verifies viewport can be resized multiple times.
// Ref: TestPageViewport.java#shouldResizeMultipleTimes
func TestSetViewportSizeCanResizeMultipleTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 1024, 768))
	must.NoError(page.SetViewportSize(ctx, 800, 600))
	must.NoError(page.SetViewportSize(ctx, 640, 480))

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Equal(640, vp.Width)
	is.Equal(480, vp.Height)
}

// TestViewportSizeDefaultIsNotNil verifies ViewportSize is not nil for new page.
// Ref: TestPageViewport.java#shouldBeNonNilByDefault
func TestViewportSizeDefaultIsNotNil(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Greater(vp.Width, 0)
	is.Greater(vp.Height, 0)
}

// TestSetViewportSizeLargeEx2 verifies SetViewportSize works with large dimensions.
// Ref: TestPageViewport.java#shouldSupportLargeViewport
func TestSetViewportSizeLargeEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 1920, 1080))

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Equal(1920, vp.Width)
	is.Equal(1080, vp.Height)
}

// TestSetViewportSizeSmallEx2 verifies SetViewportSize works with small dimensions.
// Ref: TestPageViewport.java#shouldSupportSmallViewport
func TestSetViewportSizeSmallEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 320, 568))

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Equal(320, vp.Width)
	is.Equal(568, vp.Height)
}

// TestSetViewportSizeReflectsInWindowHeightEx2 verifies window.innerHeight matches viewport.
// Ref: TestPageViewport.java#shouldReflectInWindowHeight
func TestSetViewportSizeReflectsInWindowHeightEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))

	height, err := page.Evaluate(ctx, `() => window.innerHeight`)
	must.NoError(err)
	is.Equal(float64(600), height)
}

// TestSetViewportSizeAfterNavigationEx2 verifies ViewportSize persists after navigation.
// Ref: TestPageViewport.java#shouldPersistAfterNavigation
func TestSetViewportSizeAfterNavigationEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 1024, 768))
	must.NoError(page.SetContent(ctx, `<p>New content</p>`))

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Equal(1024, vp.Width)
}

// TestSetViewportSizeWidthAndHeightBothSetEx2 verifies both width and height are set correctly.
// Ref: TestPageViewport.java#shouldSetBothDimensions
func TestSetViewportSizeWidthAndHeightBothSetEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 1366, 768))

	vp, err := page.ViewportSize(ctx)
	must.NoError(err)
	must.NotNil(vp)
	is.Equal(1366, vp.Width)
	is.Equal(768, vp.Height)
}
