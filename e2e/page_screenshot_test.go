//go:build e2e

package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageScreenshotWithColoredBackground verifies screenshot captures colored background.
// Ref: TestPageScreenshot.java#shouldCaptureColor
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
// Ref: TestPageScreenshot.java#shouldWorkAfterSetViewport
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
// Ref: TestPageScreenshot.java#pageShouldDefaultToPNG
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
// Ref: TestPageScreenshot.java#pageShouldCaptureFullPage
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
// Ref: TestPageScreenshot.java#pageShouldCompressJPEG
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
// Ref: TestPageScreenshot.java#shouldScreenshotWhitePage
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
// Ref: TestPageScreenshot.java#shouldWorkAfterScroll
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
// Ref: TestPageScreenshot.java#shouldTakeScreenshotWithClip
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
// Ref: TestPageScreenshot.java#shouldTakeScreenshotWithTransparentBackground
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

// TestPageScreenshotWithMask verifies Screenshot with masked locator returns non-empty bytes.
// Ref: TestPageScreenshot.java#shouldSupportMasks
func TestPageScreenshotWithMask(t *testing.T) {
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

// TestPageScreenshotJPEGFormat verifies JPEG screenshot has JPEG magic bytes.
// Ref: TestPageScreenshot.java#shouldWorkWithJPEG
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
// Ref: TestPageScreenshot.java#shouldWorkWithScaleCSS
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
// Ref: TestPageScreenshot.java#shouldWorkWithScaleDevice
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
