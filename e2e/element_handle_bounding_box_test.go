//go:build e2e

// E2E tests for: TestElementHandleBoundingBox.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/bb", "text/html", `
		<style>
			body { margin: 0; }
			#box { width: 100px; height: 50px; background: red; }
		</style>
		<div id="box"></div>
	`)
	err := page.Goto(ctx, srv.Prefix()+"/bb")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "#box")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	bb, err := el.BoundingBox(ctx)
	must.NoError(err, "BoundingBox failed")
	must.NotNil(bb, "BoundingBox returned nil")
	if bb.Width != 100 || bb.Height != 50 {
		t.Errorf("BoundingBox = {w=%f, h=%f}, want {w=100, h=50}", bb.Width, bb.Height)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="box" style="width:100px;height:50px;position:absolute;top:10px;left:20px;"></div>`))

	el, err := page.QuerySelector(ctx, "#box")
	must.NoError(err)
	must.NotNil(el)

	bb, err := el.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(100), bb.Width)
	is.Equal(float64(50), bb.Height)
}

// TestElementHandleBoundingBoxCorrectDimensions verifies BoundingBox returns correct dimensions.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxCorrectDimensions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="box" style="width:200px;height:100px;position:absolute;left:10px;top:20px"></div>
	`))

	handle, err := page.QuerySelector(ctx, "#box")
	must.NoError(err)
	must.NotNil(handle)

	bb, err := handle.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(200), bb.Width)
	is.Equal(float64(100), bb.Height)
}

// TestElementHandleBoundingBoxNilForDisplayNone verifies nil for display:none.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxNilForDisplayNone(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="hidden" style="display:none;width:100px;height:100px"></div>
	`))

	handle, err := page.QuerySelector(ctx, "#hidden")
	must.NoError(err)
	must.NotNil(handle)

	bb, err := handle.BoundingBox(ctx)
	must.NoError(err)
	is.Nil(bb)
}

// TestElementHandleBoundingBoxAfterScroll verifies BoundingBox updates after scroll.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxAfterScroll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:1000px"></div>
		<div id="target" style="width:100px;height:50px">target</div>
	`))

	handle, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(handle)

	// Scroll to target
	must.NoError(handle.ScrollIntoViewIfNeeded(ctx))

	bb, err := handle.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(100), bb.Width)
	is.Equal(float64(50), bb.Height)
}

// TestElementHandleBoundingBoxPositionValues verifies X and Y are correct.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxPositionValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="position:absolute;left:30px;top:40px;width:60px;height:70px" id="el"></div>
	`))

	handle, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(handle)

	bb, err := handle.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(30), bb.X)
	is.Equal(float64(40), bb.Y)
}

// --- From element_handle_extended_test.go ---

// TestElementHandleBoundingBoxEx verifies ElementHandle BoundingBox returns valid dimensions.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="width:200px;height:100px;position:fixed;top:10px;left:20px;"></div>
	`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	box, err := el.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
	is.Equal(float64(200), box.Width)
	is.Equal(float64(100), box.Height)
}

// TestElementHandleGetBoundingBoxForVisibleElement verifies BoundingBox returns non-nil for visible element.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleGetBoundingBoxForVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))
	must.NoError(page.SetContent(ctx, `<div style="width:50px;height:50px;position:absolute;top:10px;left:20px;background:red;"></div>`))

	el, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(el)

	bb, err := el.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(20), bb.X)
	is.Equal(float64(10), bb.Y)
	is.Equal(float64(50), bb.Width)
	is.Equal(float64(50), bb.Height)
}

// TestLocatorBoundingBoxReturnsCorrectDimensions verifies BoundingBox dimensions match CSS.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestLocatorBoundingBoxReturnsCorrectDimensions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))
	must.NoError(page.SetContent(ctx, `
		<div id="box" style="width:120px;height:80px;position:absolute;top:30px;left:40px;background:red;"></div>
	`))

	bb, err := page.Locator("#box").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(40), bb.X)
	is.Equal(float64(30), bb.Y)
	is.Equal(float64(120), bb.Width)
	is.Equal(float64(80), bb.Height)
}

// TestLocatorBoundingBoxNilForDetachedElement verifies BoundingBox returns nil for invisible element.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestLocatorBoundingBoxNilForDetachedElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="display:none">hidden</div>`))

	bb, err := page.Locator("#hidden").BoundingBox(ctx)
	must.NoError(err)
	is.Nil(bb)
}

// TestLocatorBoundingBoxForInlineElement verifies BoundingBox works for inline elements.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestLocatorBoundingBoxForInlineElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="span">inline text</span>`))

	bb, err := page.Locator("#span").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxInlineElementEx4 verifies BoundingBox works for inline elements.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleBoundingBox.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestBoundingBoxInlineElementEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Text <span id="s">inline</span></p>`))

	bb, err := page.Locator("#s").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// Ref: TestElementHandleBoundingBox.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 500, 500))
	must.NoError(page.SetContent(ctx, a3GridHTML))
	elementHandle, err := page.QuerySelector(ctx, ".box:nth-of-type(13)")
	must.NoError(err)
	must.NotNil(elementHandle)
	box, err := elementHandle.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
	is.Equal(float64(100), box.X)
	is.Equal(float64(50), box.Y)
	is.Equal(float64(50), box.Width)
	is.Equal(float64(50), box.Height)
}

// Ref: TestElementHandleBoundingBox.java#shouldHandleNestedFrames  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxShouldHandleNestedFrames(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frames/frame.html", "text/html", `<div style="height:18px;width:100%">Hi, I'm frame</div>`)
	srv.ServeWithBody("/frames/nested-frames.html", "text/html", `<style>body { margin: 0; }
iframe { position: absolute; border: 0; left: 16px; top: 216px; width: 284px; height: 100px; }</style>
<iframe src="./frame.html" name="dos"></iframe>`)

	must.NoError(page.SetViewportSize(ctx, 616, 500))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/frames/nested-frames.html"))
	nestedFrame := page.Frame("dos")
	must.NotNil(nestedFrame)
	elementHandle, err := nestedFrame.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(elementHandle)
	box, err := elementHandle.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
	is.Equal(float64(24), box.X)
	is.Equal(float64(224), box.Y)
	is.Equal(float64(268), box.Width)
	is.Equal(float64(18), box.Height)
}

// Ref: TestElementHandleBoundingBox.java#shouldReturnNullForInvisibleElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxShouldReturnNullForInvisibleElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div style='display:none'>hi</div>"))
	element, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(element)
	box, err := element.BoundingBox(ctx)
	must.NoError(err)
	must.Nil(box)
}

// Ref: TestElementHandleBoundingBox.java#shouldForceALayout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxShouldForceALayout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 500, 500))
	must.NoError(page.SetContent(ctx, "<div style='width: 100px; height: 100px'>hello</div>"))
	elementHandle, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(elementHandle)
	_, err = page.Evaluate(ctx, "element => element.style.height = '200px'", elementHandle)
	must.NoError(err)
	box, err := elementHandle.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
	is.Equal(float64(8), box.X)
	is.Equal(float64(8), box.Y)
	is.Equal(float64(100), box.Width)
	is.Equal(float64(200), box.Height)
}

// Ref: TestElementHandleBoundingBox.java#shouldWorkWithSVGNodes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxShouldWorkWithSVGNodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<svg xmlns='http://www.w3.org/2000/svg' width='500' height='500'>\n"+
		"<rect id='theRect' x='30' y='50' width='200' height='300'></rect>\n"+
		"</svg>"))
	element, err := page.QuerySelector(ctx, "#therect")
	must.NoError(err)
	must.NotNil(element)
	pwBoundingBox, err := element.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(pwBoundingBox)
	res, err := page.Evaluate(ctx, `e => {
  const rect = e.getBoundingClientRect();
  return { x: rect.x, y: rect.y, width: rect.width, height: rect.height };
}`, element)
	must.NoError(err)
	web, ok := res.(map[string]any)
	must.True(ok)
	is.Equal(a3Num(web["x"]), pwBoundingBox.X)
	is.Equal(a3Num(web["y"]), pwBoundingBox.Y)
	is.Equal(a3Num(web["width"]), pwBoundingBox.Width)
	is.Equal(a3Num(web["height"]), pwBoundingBox.Height)
}

// Ref: TestElementHandleBoundingBox.java#shouldWorkWhenInlineBoxChildIsOutsideOfViewport  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleBoundingBox.java
func TestElementHandleBoundingBoxShouldWorkWhenInlineBoxChildIsOutsideOfViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<style>
i {
  position: absolute;
  top: -1000px;
}
body {
  margin: 0;
  font-size: 12px;
}
</style>
<span><i>woof</i><b>doggo</b></span>`))
	handle, err := page.QuerySelector(ctx, "span")
	must.NoError(err)
	must.NotNil(handle)
	pwBoundingBox, err := handle.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(pwBoundingBox)
	res, err := handle.Evaluate(ctx, `e => {
  const rect = e.getBoundingClientRect();
  return { x: rect.x, y: rect.y, width: rect.width, height: rect.height };
}`)
	must.NoError(err)
	web, ok := res.(map[string]any)
	must.True(ok)
	is.Equal(a3Round(a3Num(web["x"])), a3Round(pwBoundingBox.X))
	is.Equal(a3Round(a3Num(web["y"])), a3Round(pwBoundingBox.Y))
	is.Equal(a3Round(a3Num(web["width"])), a3Round(pwBoundingBox.Width))
	is.Equal(a3Round(a3Num(web["height"])), a3Round(pwBoundingBox.Height))
}
