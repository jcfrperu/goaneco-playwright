//go:build e2e

package e2e

import (
	"bytes"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
