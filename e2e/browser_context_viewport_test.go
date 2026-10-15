//go:build e2e

// E2E tests for: TestBrowserContextViewport.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestPageViewportSizeAfterSetViewport verifies that ViewportSize reflects a subsequent SetViewportSize call.
// Ref: TestPageViewportSizeAfterChange
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestPageViewportSizeAfterSetViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetViewportSize(ctx, 640, 480)
	must.NoError(err, "SetViewportSize failed")

	size, err := page.ViewportSize(ctx)
	must.NoError(err, "ViewportSize failed")
	must.NotNil(size, "ViewportSize returned nil")
	is.Equal(640, size.Width, "width should be 640 after SetViewportSize")
	is.Equal(480, size.Height, "height should be 480 after SetViewportSize")
}

// TestBrowserContextViewportShouldReturnCorrectOuterWidthAndOuterHeight verifies outerWidth/outerHeight reflect viewport settings.
// Ref: TestBrowserContextViewport.java#shouldReturnCorrectOuterWidthAndOuterHeight  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldReturnCorrectOuterWidthAndOuterHeight(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 1280, Height: 720},
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	innerW, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err, "Evaluate innerWidth failed")
	innerH, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err, "Evaluate innerHeight failed")

	if innerW != float64(1280) {
		t.Errorf("innerWidth = %v, want 1280", innerW)
	}
	if innerH != float64(720) {
		t.Errorf("innerHeight = %v, want 720", innerH)
	}
}

// TestBrowserContextViewportShouldGetTheProperDefaultViewPortSize verifies viewport size is correctly returned.
// Ref: TestBrowserContextViewport.java#shouldGetTheProperDefaultViewPortSize  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldGetTheProperDefaultViewPortSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 800, Height: 600},
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	w, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err, "innerWidth failed")
	h, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err, "innerHeight failed")

	if w != float64(800) || h != float64(600) {
		t.Errorf("viewport = %vx%v, want 800x600", w, h)
	}
}

// TestBrowserContextViewportShouldSetTheProperViewportSize verifies that SetViewportSize updates the viewport.
// Ref: TestBrowserContextViewport.java#shouldSetTheProperViewportSize  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldSetTheProperViewportSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 123, 456))
	w, err := page.Evaluate(ctx, "() => window.innerWidth")
	must.NoError(err)
	h, err := page.Evaluate(ctx, "() => window.innerHeight")
	must.NoError(err)
	is.Equal(float64(123), w)
	is.Equal(float64(456), h)
}

// TestBrowserContextViewportShouldEmulateDeviceWidth verifies screen.width and matchMedia respond to viewport size changes.
// Ref: TestBrowserContextViewport.java#shouldEmulateDeviceWidth  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldEmulateDeviceWidth(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 200, 200))
	screenW, err := page.Evaluate(ctx, "() => window.screen.width")
	must.NoError(err)
	is.Equal(float64(200), screenW)

	matchTrue, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-width: 100px)').matches")
	is.Equal(true, matchTrue)
	matchFalse, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-width: 300px)').matches")
	is.Equal(false, matchFalse)
	matchFalse2, _ := page.Evaluate(ctx, "() => matchMedia('(max-device-width: 100px)').matches")
	is.Equal(false, matchFalse2)
	matchTrue2, _ := page.Evaluate(ctx, "() => matchMedia('(max-device-width: 300px)').matches")
	is.Equal(true, matchTrue2)
	matchFalse3, _ := page.Evaluate(ctx, "() => matchMedia('(device-width: 500px)').matches")
	is.Equal(false, matchFalse3)
	matchTrue3, _ := page.Evaluate(ctx, "() => matchMedia('(device-width: 200px)').matches")
	is.Equal(true, matchTrue3)

	must.NoError(page.SetViewportSize(ctx, 500, 500))
	screenW2, err := page.Evaluate(ctx, "() => window.screen.width")
	must.NoError(err)
	is.Equal(float64(500), screenW2)
}

// TestBrowserContextViewportShouldEmulateDeviceHeight verifies screen.height and matchMedia respond to viewport size changes.
// Ref: TestBrowserContextViewport.java#shouldEmulateDeviceHeight  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldEmulateDeviceHeight(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 200, 200))
	screenH, err := page.Evaluate(ctx, "() => window.screen.height")
	must.NoError(err)
	is.Equal(float64(200), screenH)

	matchTrue, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-height: 100px)').matches")
	is.Equal(true, matchTrue)
	matchFalse, _ := page.Evaluate(ctx, "() => matchMedia('(min-device-height: 300px)').matches")
	is.Equal(false, matchFalse)

	must.NoError(page.SetViewportSize(ctx, 500, 500))
	screenH2, err := page.Evaluate(ctx, "() => window.screen.height")
	must.NoError(err)
	is.Equal(float64(500), screenH2)
}

// TestBrowserContextViewportShouldEmulateAvailWidthAndAvailHeight verifies screen.availWidth/availHeight match viewport.
// Ref: TestBrowserContextViewport.java#shouldEmulateAvailWidthAndAvailHeight  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldEmulateAvailWidthAndAvailHeight(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 500, 600))
	availW, err := page.Evaluate(ctx, "() => window.screen.availWidth")
	must.NoError(err)
	availH, err := page.Evaluate(ctx, "() => window.screen.availHeight")
	must.NoError(err)
	is.Equal(float64(500), availW)
	is.Equal(float64(600), availH)
}

// TestBrowserContextViewportShouldNotHaveTouchByDefault verifies that touch events are not available by default.
// Ref: TestBrowserContextViewport.java#shouldNotHaveTouchByDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldNotHaveTouchByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Test via SetContent to avoid dependency on specific test server files
	must.NoError(page.SetContent(ctx, `<html><body><script>
		document.body.textContent = 'ontouchstart' in window ? 'YES' : 'NO';
	</script></body></html>`))
	hasTouchStart, err := page.Evaluate(ctx, "() => 'ontouchstart' in window")
	must.NoError(err)
	is.Equal(false, hasTouchStart)
}

// TestBrowserContextViewportShouldSupportTouchWithNullViewport verifies that HasTouch option enables touch support (with null viewport).
// Ref: TestBrowserContextViewport.java#shouldSupportTouchWithNullViewport  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldSupportTouchWithNullViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	hasTouch := true
	noDefaultViewport := true
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		HasTouch:          &hasTouch,
		NoDefaultViewport: &noDefaultViewport,
	})
	must.NoError(err, "NewContext with HasTouch+NoDefaultViewport failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)
	hasTouchStart, err := page.Evaluate(ctx, "() => 'ontouchstart' in window")
	must.NoError(err)
	is.Equal(true, hasTouchStart)
}

// TestBrowserContextViewportShouldReportNullViewPortSizeWhenGivenNullViewport verifies that ViewportSize returns nil when NoDefaultViewport is set.
// Ref: TestBrowserContextViewport.java#shouldReportNullViewPortSizeWhenGivenNullViewport  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
func TestBrowserContextViewportShouldReportNullViewPortSizeWhenGivenNullViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	noDefaultViewport := true
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		NoDefaultViewport: &noDefaultViewport,
	})
	must.NoError(err, "NewContext with NoDefaultViewport failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	size, err := page.ViewportSize(ctx)
	must.NoError(err)
	is.Nil(size, "ViewportSize should return nil when no default viewport is set")
}

// ---------------------------------------------------------------------------
// From browser_context_locale_extra_test.go
// ---------------------------------------------------------------------------

// TestPageViewportSizeReturnsCurrentSize verifies ViewportSize returns the set size.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextViewport.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextViewport.java
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
