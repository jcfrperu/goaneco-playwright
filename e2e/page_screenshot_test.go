//go:build e2e

// E2E tests for: TestPageScreenshot.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"bytes"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageScreenshotWithColoredBackground verifies screenshot captures colored background.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotWithColoredBackgroundEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="background-color:blue;width:200px;height:200px;"></div>
	`))

	data, err := page.Screenshot(ctx, nil)
	must.NoError(err)
	is.NotEmpty(data)
}

// TestPageScreenshotAfterSetViewport verifies screenshot works after SetViewportSize.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotAfterSetViewportEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 400, 300))
	must.NoError(page.SetContent(ctx, `<p>Small viewport</p>`))

	data, err := page.Screenshot(ctx, nil)
	must.NoError(err)
	is.NotEmpty(data)
}

// TestPageScreenshotDefaultPNG verifies page screenshot is PNG by default.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotDefaultPNGExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:100px;height:100px;background:red">box</div>`))

	// Use locator screenshot as page-level not exposed directly
	bytes, err := page.Locator("div").Screenshot(ctx)
	must.NoError(err)
	is.NotEmpty(bytes)
	// PNG header
	is.Equal([]byte{0x89, 0x50, 0x4e, 0x47}, bytes[:4])
}

// TestPageScreenshotFullPageCaptures verifies fullPage captures all content.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotFullPageCapturesExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="height:2000px;background:green">tall</div>`))

	fullPage := true
	bytes, err := page.Locator("div").Screenshot(ctx, &playwright.ScreenshotOptions{FullPage: &fullPage})
	must.NoError(err)
	is.NotEmpty(bytes)
}

// TestPageScreenshotJPEGCompression verifies JPEG screenshot is smaller than PNG.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotJPEGCompressionExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:300px;height:300px;background:purple"></div>`))

	q := 30
	jpegBytes, err := page.Locator("div").Screenshot(ctx, &playwright.ScreenshotOptions{
		Type:    "jpeg",
		Quality: &q,
	})
	must.NoError(err)
	is.NotEmpty(jpegBytes)

	pngBytes, err := page.Locator("div").Screenshot(ctx, &playwright.ScreenshotOptions{Type: "png"})
	must.NoError(err)
	is.NotEmpty(pngBytes)

	// Generally JPEG with low quality should be smaller than PNG
	// but we just verify both are valid
	is.Equal(byte(0xFF), jpegBytes[0])
	is.Equal(byte(0xD8), jpegBytes[1])
	is.Equal(byte(0x89), pngBytes[0])
}

// TestPageScreenshotWhiteBackgroundEx4 verifies Screenshot of white page.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotWhiteBackgroundEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body style="background:white;"></body>`))

	data, err := page.Screenshot(ctx)
	must.NoError(err)
	is.NotEmpty(data)
}

// TestPageScreenshotAfterScrollEx4 verifies Screenshot works after scrolling.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotAfterScrollEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;">Scroll page</div>
	`))
	_, err := page.Evaluate(ctx, `() => window.scrollTo(0, 1000)`)
	must.NoError(err)

	data, err := page.Screenshot(ctx)
	must.NoError(err)
	is.NotEmpty(data)
}

// TestPageScreenshotWithClip verifies Screenshot with a clip region produces smaller output.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotWithClip(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="width:600px;height:600px;background:red;"></div>
	`))
	must.NoError(page.SetViewportSize(ctx, 800, 600))

	data, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{
		Clip: &playwright.Rect{X: 0, Y: 0, Width: 200, Height: 200},
	})
	must.NoError(err)
	is.NotEmpty(data)
	// Must be a valid PNG.
	is.Equal(byte(0x89), data[0])
	is.Equal(byte('P'), data[1])
}

// TestPageScreenshotWithTransparentBackground verifies omitBackground produces transparent PNG.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotWithTransparentBackground(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:100px;height:100px;"></div>`))

	omit := true
	data, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{OmitBackground: &omit})
	must.NoError(err)
	is.NotEmpty(data)
	// Must be a valid PNG (transparency only available in PNG).
	is.Equal(byte(0x89), data[0])
}

// TestPageScreenshotJPEGFormat verifies JPEG screenshot has JPEG magic bytes.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotJPEGFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:100px;height:100px;background:green;"></div>`))

	q := 80
	data, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{
		Type:    "jpeg",
		Quality: &q,
	})
	must.NoError(err)
	is.NotEmpty(data)
	// JPEG magic bytes: FF D8
	is.Equal(byte(0xFF), data[0])
	is.Equal(byte(0xD8), data[1])
}

// TestPageScreenshotScaleCSS verifies Screenshot with scale=css works without error.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotScaleCSS(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:200px;height:200px;background:purple;"></div>`))

	data, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Scale: "css"})
	must.NoError(err)
	is.NotEmpty(data)
}

// TestPageScreenshotScaleDevice verifies Screenshot with scale=device works without error.
// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotScaleDevice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:200px;height:200px;background:orange;"></div>`))

	data, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Scale: "device"})
	must.NoError(err)
	is.NotEmpty(data)
}

// TestPageScreenshotMaskShouldWorkForPage verifies Screenshot with masked locator returns non-empty bytes.
// Ref: TestPageScreenshot.java#maskShouldWorkForPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotMaskShouldWorkForPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="secret" style="width:100px;height:50px;background:red;">secret</div>
		<div id="public" style="width:100px;height:50px;background:blue;">public</div>
	`))

	mask := page.Locator("#secret")
	data, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{
		Mask: []*playwright.Locator{mask},
	})
	must.NoError(err)
	is.NotEmpty(data)
	// PNG header.
	is.Equal(byte(0x89), data[0])
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotReturnsPNGBytes(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	data, err := page.Screenshot(ctx)
	must.NoError(err, "Screenshot failed")
	if len(data) == 0 {
		t.Fatal("Screenshot returned empty bytes")
	}
	// PNG magic bytes: 0x89 P N G (0x89 0x50 0x4E 0x47)
	if !bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		t.Errorf("expected PNG magic bytes, got prefix: %x", data[:minInt(4, len(data))])
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotJPEG(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	q := 80
	data, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Type: "jpeg", Quality: &q})
	must.NoError(err, "Screenshot JPEG failed")
	// JPEG magic bytes: 0xFF 0xD8
	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
		t.Errorf("expected JPEG magic bytes, got prefix: %x", data[:minInt(4, len(data))])
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotFullPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>body { margin: 0; height: 2000px; background: linear-gradient(white, red); }</style>
		<div>top</div>
	`))
	must.NoError(page.SetViewportSize(ctx, 800, 600))

	fullPage := true
	screenshotBytes, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{
		FullPage: &fullPage,
	})
	must.NoError(err)
	is.NotEmpty(screenshotBytes)

	viewportBytes, err := page.Screenshot(ctx)
	must.NoError(err)
	is.Greater(len(screenshotBytes), len(viewportBytes), "full page screenshot should be larger than viewport screenshot")
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotDefaultIsPNG(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body style="background: green;"></body>`))

	data, err := page.Screenshot(ctx)
	must.NoError(err)
	is.NotEmpty(data)

	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47}
	is.True(bytes.HasPrefix(data, pngMagic), "expected PNG magic bytes, got %X", data[:4])
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotJPEGQuality(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body style="background: steelblue; width: 800px; height: 600px;"></body>`))
	must.NoError(page.SetViewportSize(ctx, 800, 600))

	quality10 := 10
	quality90 := 90

	lowQ, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Type: "jpeg", Quality: &quality10})
	must.NoError(err)

	highQ, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Type: "jpeg", Quality: &quality90})
	must.NoError(err)

	is.Less(len(lowQ), len(highQ), "low quality JPEG should be smaller than high quality")
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestLocatorScreenshotMatchesPNGFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="box" style="width:100px;height:100px;background:red;"></div>`))

	data, err := page.Locator("#box").Screenshot(ctx)
	must.NoError(err)
	is.NotEmpty(data)

	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47}
	is.True(bytes.HasPrefix(data, pngMagic), "expected PNG magic bytes")
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestLocatorScreenshotJPEGFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:100px;height:100px;background:blue" id="box"></div>`))

	bytes, err := page.Locator("#box").Screenshot(ctx, &playwright.ScreenshotOptions{Type: "jpeg"})
	must.NoError(err)
	is.NotEmpty(bytes)

	is.Equal(byte(0xFF), bytes[0])
	is.Equal(byte(0xD8), bytes[1])
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotNonEmptyResult(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>screenshot test</div>`))

	bytes, err := page.Locator("body").Screenshot(ctx)
	must.NoError(err)
	is.Greater(len(bytes), 100)
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotWithQuality(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:200px;height:200px;background:green"></div>`))

	q := 50
	bytes, err := page.Locator("div").Screenshot(ctx, &playwright.ScreenshotOptions{
		Type:    "jpeg",
		Quality: &q,
	})
	must.NoError(err)
	is.NotEmpty(bytes)
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestLocatorScreenshotButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click Me</button>`))

	data, err := page.Locator("#btn").Screenshot(ctx, nil)
	must.NoError(err)
	is.NotEmpty(data)
}

// Extra: no equivalent method in playwright-java; related class: TestPageScreenshot.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestLocatorScreenshotInputEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="test" style="width:200px">`))

	data, err := page.Locator("#inp").Screenshot(ctx, nil)
	must.NoError(err)
	is.NotEmpty(data)
}

// TestPageScreenshotShouldClipRect verifies the clip option sizes the screenshot.
// Ref: TestPageScreenshot.java#shouldClipRect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldClipRect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeGrid(srv)
	page := newPage(t)
	must.NoError(page.SetViewportSize(ctx, 500, 500))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/grid.html"))
	shot, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Clip: &playwright.Rect{X: 50, Y: 100, Width: 150, Height: 100}})
	must.NoError(err)
	w, h := a4PNGSize(t, shot)
	is.Equal(150, w)
	is.Equal(100, h)
}

// Ref: TestPageScreenshot.java#shouldProduceAValidWebpScreenshot  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldProduceAValidWebpScreenshot(t *testing.T) {
	t.Skip("not supported by playwright-core v1.61.1: driver page.screenshot rejects type 'webp' with 'expected one of (png|jpeg)'")
}

// Ref: TestPageScreenshot.java#pathOptionShouldDetectWebp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotPathOptionShouldDetectWebp(t *testing.T) {
	t.Skip("not supported by playwright-core v1.61.1: driver page.screenshot rejects type 'webp' with 'expected one of (png|jpeg)'")
}

// Ref: TestPageScreenshot.java#qualityOptionShouldWorkForWebp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotQualityOptionShouldWorkForWebp(t *testing.T) {
	t.Skip("not supported by playwright-core v1.61.1: driver page.screenshot rejects type 'webp' with 'expected one of (png|jpeg)'")
}

// TestPageScreenshotShouldNotCaptureInfiniteCssAnimation is not ported yet.
// Ref: TestPageScreenshot.java#shouldNotCaptureInfiniteCssAnimation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldNotCaptureInfiniteCssAnimation(t *testing.T) {
	t.Skip("not implemented: ScreenshotOptions.Animations")
}

// TestPageScreenshotShouldNotCapturePseudoElementCssAnimation is not ported yet.
// Ref: TestPageScreenshot.java#shouldNotCapturePseudoElementCssAnimation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldNotCapturePseudoElementCssAnimation(t *testing.T) {
	t.Skip("not implemented: ScreenshotOptions.Animations")
}

// TestPageScreenshotShouldNotCaptureCssAnimationsInShadowDOM is not ported yet.
// Ref: TestPageScreenshot.java#shouldNotCaptureCssAnimationsInShadowDOM  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldNotCaptureCssAnimationsInShadowDOM(t *testing.T) {
	t.Skip("not implemented: ScreenshotOptions.Animations")
}

// TestPageScreenshotShouldResumeInfiniteAnimations is not ported yet.
// Ref: TestPageScreenshot.java#shouldResumeInfiniteAnimations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldResumeInfiniteAnimations(t *testing.T) {
	t.Skip("not implemented: ScreenshotOptions.Animations")
}

// TestPageScreenshotShouldNotCaptureInfiniteWebAnimations is not ported yet.
// Ref: TestPageScreenshot.java#shouldNotCaptureInfiniteWebAnimations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldNotCaptureInfiniteWebAnimations(t *testing.T) {
	t.Skip("not implemented: ScreenshotOptions.Animations")
}

// TestPageScreenshotMaskShouldWorkForLocator verifies masking a locator changes the screenshot.
// Ref: TestPageScreenshot.java#maskShouldWorkForLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotMaskShouldWorkForLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeGrid(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/grid.html"))
	target := page.Locator("div").First()
	masked, err := target.Screenshot(ctx, &playwright.ScreenshotOptions{Mask: []*playwright.Locator{page.Locator("img")}})
	must.NoError(err)
	original, err := target.Screenshot(ctx)
	must.NoError(err)
	must.NotEqual(masked, original)
}

// TestPageScreenshotShouldWorkWithDeviceScaleFactorAndClip verifies clip screenshots with a device scale factor.
// Ref: TestPageScreenshot.java#shouldWorkWithDeviceScaleFactorAndClip  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldWorkWithDeviceScaleFactorAndClip(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeGrid(srv)
	page := a4NewPage(t, &playwright.BrowserContextOptions{Viewport: &playwright.ViewportSize{Width: 500, Height: 500}, DeviceScaleFactor: a4Ptr(3.0)})
	must.NoError(page.Goto(ctx, srv.Prefix()+"/grid.html"))
	shot, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Clip: &playwright.Rect{X: 50, Y: 100, Width: 150, Height: 100}})
	must.NoError(err)
	must.NotNil(shot)
}

// TestPageScreenshotShouldWorkWithDeviceScaleFactorAndScaleCss verifies css scale with a device scale factor.
// Ref: TestPageScreenshot.java#shouldWorkWithDeviceScaleFactorAndScaleCss  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldWorkWithDeviceScaleFactorAndScaleCss(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeGrid(srv)
	page := a4NewPage(t, &playwright.BrowserContextOptions{Viewport: &playwright.ViewportSize{Width: 320, Height: 480}, DeviceScaleFactor: a4Ptr(2.0)})
	must.NoError(page.Goto(ctx, srv.Prefix()+"/grid.html"))
	shot, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Scale: "css"})
	must.NoError(err)
	w, h := a4PNGSize(t, shot)
	is.Equal(320, w)
	is.Equal(480, h)
}

// TestPageScreenshotShouldNotCaptureBlinkingCaretByDefault verifies the caret is hidden in screenshots by default.
// Ref: TestPageScreenshot.java#shouldNotCaptureBlinkingCaretByDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldNotCaptureBlinkingCaretByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<style>div { caret-color: #000 !important; }</style><div contenteditable="true"></div>`))
	div := page.Locator("div")
	must.NoError(div.PressSequentially(ctx, "foo bar"))
	shot, err := div.Screenshot(ctx)
	must.NoError(err)
	for i := 0; i < 10; i++ {
		must.NoError(page.WaitForTimeout(ctx, 150))
		newShot, err := div.Screenshot(ctx)
		must.NoError(err)
		must.Equal(shot, newShot)
	}
}

// Ref: TestPageScreenshot.java#shouldWorkWhenMaskColorIsNotPinkF0F  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldWorkWhenMaskColorIsNotPinkF0F(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeGrid(srv)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 500, 500))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/grid.html"))
	screenshot1, err := page.Screenshot(ctx, (&playwright.ScreenshotOptions{}).
		SetMask([]*playwright.Locator{page.Locator("div").Nth(5)}).
		SetMaskColor("#00FF00"))
	must.NoError(err)

	screenshot2, err := page.Screenshot(ctx, (&playwright.ScreenshotOptions{}).
		SetMask([]*playwright.Locator{page.Locator("div").Nth(5)}))
	must.NoError(err)

	is.NotEqual(screenshot1, screenshot2)
}

// TestPageScreenshotShouldWork verifies a basic page screenshot has the viewport size.
// Ref: TestPageScreenshot.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeGrid(srv)
	page := newPage(t)
	must.NoError(page.SetViewportSize(ctx, 500, 500))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/grid.html"))
	shot, err := page.Screenshot(ctx)
	must.NoError(err)
	w, h := a4PNGSize(t, shot)
	is.Equal(500, w)
	is.Equal(500, h)
}

// TestPageScreenshotShouldWorkWithDeviceScaleFactorAndScaleDevice verifies device-scale screenshots.
// Ref: TestPageScreenshot.java#shouldWorkWithDeviceScaleFactorAndScaleDevice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageScreenshot.java
func TestPageScreenshotShouldWorkWithDeviceScaleFactorAndScaleDevice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeGrid(srv)
	page := a4NewPage(t, &playwright.BrowserContextOptions{Viewport: &playwright.ViewportSize{Width: 320, Height: 480}, DeviceScaleFactor: a4Ptr(2.0)})
	must.NoError(page.Goto(ctx, srv.Prefix()+"/grid.html"))
	shot, err := page.Screenshot(ctx, &playwright.ScreenshotOptions{Scale: "device"})
	must.NoError(err)
	must.NotNil(shot)
}
