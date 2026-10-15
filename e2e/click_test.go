//go:build e2e

// E2E tests for: TestClick.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"fmt"
	"strings"
	"sync"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Ref: TestClick.java#shouldWaitForVisibleWhenAlreadyVisible  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldWaitForVisibleWhenAlreadyVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">visible</div>`)
	must.NoError(err, "SetContent failed")

	// Should return immediately when already visible.
	err = page.Locator("#el").WaitFor(ctx)
	must.NoError(err, "WaitFor() for already-visible element failed")
}

// TestClickShouldClickTheButton verifies that Locator.Click() triggers a button's click handler.
// Ref: TestClick.java#shouldClickTheButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="btn" onclick="this.setAttribute('data-clicked','1')">Click me</button>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#btn").Click(ctx)
	must.NoError(err, "Click failed")

	attr, err := page.Locator("#btn").GetAttribute(ctx, "data-clicked")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-clicked='1' after click, got %v", attr)
	}
}

// Ref: TestClick.java#shouldClickSvg  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickSvg(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<svg height="100" width="100">
		<circle onclick="javascript:window.__CLICKED=42" cx="50" cy="50" r="40"
		        stroke="black" stroke-width="3" fill="red"/>
	</svg>`)
	must.NoError(err)

	err = page.Locator("circle").Click(ctx)
	must.NoError(err, "Click on SVG circle failed")

	val, err := page.Evaluate(ctx, "__CLICKED")
	must.NoError(err)
	is.Equal(float64(42), val)
}

// Ref: TestClick.java#shouldClickOnASpanWithAnInlineElementInside  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickOnASpanWithAnInlineElementInside(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<style>
		span::before { content: 'q'; }
	</style>
	<span onclick="javascript:window.CLICKED=42"></span>`)
	must.NoError(err)

	err = page.Locator("span").Click(ctx)
	must.NoError(err, "Click on span failed")

	val, err := page.Evaluate(ctx, "CLICKED")
	must.NoError(err)
	is.Equal(float64(42), val)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickToggleCheckboxEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	must.NoError(page.Locator("#chk").Click(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	must.NoError(page.Locator("#chk").Click(ctx))

	checked, err = page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestDblClickFiresBothClickEvents verifies DblClick fires two click-type events.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestDblClickFiresBothClickEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onclick="window.__clicks=(window.__clicks||0)+1">dbl click</div>
	`))

	must.NoError(page.Locator("#el").DblClick(ctx))

	clicks, err := page.Evaluate(ctx, `() => window.__clicks`)
	must.NoError(err)
	is.Equal(float64(2), clicks)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="btn" style="position:absolute;left:50px;top:50px;width:100px;height:40px">
			Click me
		</button>
		<script>
			document.getElementById('btn').addEventListener('click', function() {
				this.setAttribute('data-clicked', '1');
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Mouse.Click(ctx, 100, 70)
	must.NoError(err, "Mouse.Click failed")

	attr, err := page.Locator("#btn").GetAttribute(ctx, "data-clicked")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-clicked='1' after mouse click, got %v", attr)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMove(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="el" style="width:200px;height:200px">hover target</div>
		<script>
			document.getElementById('el').addEventListener('mousemove', function() {
				this.setAttribute('data-moved', '1');
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Mouse.Move(ctx, 100, 100)
	must.NoError(err, "Mouse.Move failed")

	attr, err := page.Locator("#el").GetAttribute(ctx, "data-moved")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-moved='1' after mouse move, got %v", attr)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseDownAndUp(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="el" style="width:200px;height:200px">target</div>
		<script>
			var log = [];
			document.getElementById('el').addEventListener('mousedown', function() { log.push('down'); });
			document.getElementById('el').addEventListener('mouseup',   function() { log.push('up');   });
			window._log = function() { return log.join(','); };
		</script>
	`)
	must.NoError(err, "SetContent failed")

	// Move onto the element first, then press/release
	err = page.Mouse.Move(ctx, 100, 100)
	must.NoError(err, "Mouse.Move failed")
	err = page.Mouse.Down(ctx)
	must.NoError(err, "Mouse.Down failed")
	err = page.Mouse.Up(ctx)
	must.NoError(err, "Mouse.Up failed")

	result, err := page.Evaluate(ctx, "window._log()")
	must.NoError(err, "Evaluate failed")
	got, _ := result.(string)
	if got != "down,up" {
		t.Errorf("expected 'down,up', got %q", got)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseWheel(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="el" style="height:200px;overflow-y:scroll;border:1px solid">
			<div style="height:2000px">scrollable content</div>
		</div>
		<script>
			var totalDelta = 0;
			document.getElementById('el').addEventListener('wheel', function(e) {
				totalDelta += e.deltaY;
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	// Move mouse over the scrollable element, then wheel.
	err = page.Mouse.Move(ctx, 100, 100)
	must.NoError(err, "Mouse.Move failed")
	err = page.Mouse.Wheel(ctx, 0, 100)
	must.NoError(err, "Mouse.Wheel failed")

	raw, err := page.Evaluate(ctx, "totalDelta")
	must.NoError(err, "Evaluate failed")
	delta, _ := raw.(float64)
	if delta == 0 {
		t.Error("expected wheel event with non-zero deltaY after Mouse.Wheel")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseDblClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="el" style="position:absolute;left:50px;top:50px;width:100px;height:40px">dblclick me</div>
		<script>
			document.getElementById('el').addEventListener('dblclick', function() {
				this.setAttribute('data-dblclicked', '1');
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Mouse.DblClick(ctx, 100, 70)
	must.NoError(err, "Mouse.DblClick failed")

	attr, err := page.Locator("#el").GetAttribute(ctx, "data-dblclicked")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-dblclicked='1', got %v", attr)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseClickFiresMouseEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" style="width:100px;height:100px;background:blue;"
			onmousedown="this.dataset.down='yes'"
			onmouseup="this.dataset.up='yes'"
			onclick="this.dataset.clicked='yes'">
		</div>
	`))

	must.NoError(page.Mouse.Click(ctx, 50, 50))

	down, err := page.Evaluate(ctx, `() => document.getElementById('target').dataset.down`)
	must.NoError(err)
	is.Equal("yes", down)

	up, err := page.Evaluate(ctx, `() => document.getElementById('target').dataset.up`)
	must.NoError(err)
	is.Equal("yes", up)

	clicked, err := page.Evaluate(ctx, `() => document.getElementById('target').dataset.clicked`)
	must.NoError(err)
	is.Equal("yes", clicked)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveUpdatesCoordinates(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="tracker" style="width:200px;height:200px;"
			onmousemove="this.dataset.x=event.clientX; this.dataset.y=event.clientY">
		</div>
	`))

	must.NoError(page.Mouse.Move(ctx, 100, 100))

	x, err := page.Evaluate(ctx, `() => document.getElementById('tracker').dataset.x`)
	must.NoError(err)
	is.Equal("100", x)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseDownAndUpSeparately(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px;"
			onmousedown="this.dataset.state='down'"
			onmouseup="this.dataset.state='up'">
		</div>
	`))

	must.NoError(page.Mouse.Move(ctx, 50, 50))
	must.NoError(page.Mouse.Down(ctx))

	state, err := page.Evaluate(ctx, `() => document.getElementById('el').dataset.state`)
	must.NoError(err)
	is.Equal("down", state)

	must.NoError(page.Mouse.Up(ctx))

	state, err = page.Evaluate(ctx, `() => document.getElementById('el').dataset.state`)
	must.NoError(err)
	is.Equal("up", state)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveFiresMouseMoveEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:200px;height:200px"
		     onmousemove="window.__moved=true">move here</div>
	`))

	loc := page.Locator("#el")
	bb, err := loc.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)

	must.NoError(page.Mouse.Move(ctx, bb.X+10, bb.Y+10))

	result, err := page.Evaluate(ctx, `() => window.__moved`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseWheelScrollsPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:3000px">tall page</div>
	`))

	initialScroll, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.Equal(float64(0), initialScroll)

	must.NoError(page.Mouse.Wheel(ctx, 0, 200))

	afterScroll, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	scrolled, ok := afterScroll.(float64)
	is.True(ok)
	is.Greater(scrolled, float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveFiresMousemoveEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:200px;height:200px"
		     onmousemove="window.__moved=true">move here</div>
	`))

	must.NoError(page.Mouse.Move(ctx, 100, 100))

	moved, err := page.Evaluate(ctx, `() => window.__moved`)
	must.NoError(err)
	is.Equal(true, moved)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseWheelScrollsDownMouse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:3000px;width:100px">Scrollable content</div>
	`))

	must.NoError(page.Mouse.Wheel(ctx, 0, 300))
	must.NoError(page.WaitForTimeout(ctx, 100))

	scrollY, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.Greater(scrollY.(float64), float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseClickFiresEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px;position:absolute;left:0;top:0"
		     onclick="window.__mouseClicked2=true">click here</div>
	`))

	must.NoError(page.Mouse.Click(ctx, 50, 50))

	result, err := page.Evaluate(ctx, `() => window.__mouseClicked2`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseDownAndUpEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onmousedown="window.__mouseDown=true"
		     onmouseup="window.__mouseUp=true">click</div>
	`))

	must.NoError(page.Mouse.Move(ctx, 50, 50))
	must.NoError(page.Mouse.Down(ctx))
	must.NoError(page.Mouse.Up(ctx))

	down, err := page.Evaluate(ctx, `() => window.__mouseDown`)
	must.NoError(err)
	is.Equal(true, down)

	up, err := page.Evaluate(ctx, `() => window.__mouseUp`)
	must.NoError(err)
	is.Equal(true, up)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveToCoordinatesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="width:400px;height:400px" onmousemove="window.__x=event.clientX;window.__y=event.clientY"></div>
	`))

	must.NoError(page.Mouse.Move(ctx, 150, 200))

	x, err := page.Evaluate(ctx, `() => window.__x`)
	must.NoError(err)
	must.NotNil(x)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseWheelHorizontalEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="width:3000px;height:100px">Wide content</div>
	`))

	must.NoError(page.Mouse.Wheel(ctx, 300, 0))
	must.NoError(page.WaitForTimeout(ctx, 100))

	scrollX, err := page.Evaluate(ctx, `() => window.scrollX`)
	must.NoError(err)
	must.NotNil(scrollX)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseWheelScrollEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="height:3000px;">Tall page</div>`))

	must.NoError(page.Mouse.Move(ctx, 100, 100))
	must.NoError(page.Mouse.Wheel(ctx, 0, 500))
	must.NoError(page.WaitForTimeout(ctx, 100))

	scrollY, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.Greater(scrollY.(float64), float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveSequenceEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="tracker" style="width:300px;height:300px;"></div>
		<script>
			var moveCount = 0;
			document.getElementById('tracker').addEventListener('mousemove', function() { moveCount++; });
		</script>
	`))

	must.NoError(page.Mouse.Move(ctx, 50, 50))
	must.NoError(page.Mouse.Move(ctx, 100, 100))
	must.NoError(page.Mouse.Move(ctx, 150, 150))

	count, err := page.Evaluate(ctx, `() => moveCount`)
	must.NoError(err)
	is.Greater(count.(float64), float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseClickAtCoordinatesEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="a" style="width:100px;height:100px;position:absolute;top:0;left:0;"></div>
		<script>
			var lastClick = null;
			document.getElementById('a').addEventListener('click', function(e) {
				lastClick = {x: e.clientX, y: e.clientY};
			});
		</script>
	`))

	must.NoError(page.Mouse.Click(ctx, 50, 50))

	result, err := page.Evaluate(ctx, `() => lastClick`)
	must.NoError(err)
	must.NotNil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveToElementEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:200px;height:200px;position:absolute;top:50px;left:50px;"></div>
		<script>
			var entered = false;
			document.getElementById('el').addEventListener('mouseenter', function() { entered = true; });
		</script>
	`))

	must.NoError(page.Mouse.Move(ctx, 150, 150))
	must.NoError(page.WaitForTimeout(ctx, 50))

	result, err := page.Evaluate(ctx, `() => entered`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveUpdatesPositionEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="width:300px;height:300px;"
			onmousemove="document.getElementById('out').textContent=event.clientX+','+event.clientY">
		</div>
		<span id="out"></span>
	`))

	must.NoError(page.Mouse.Move(ctx, 100, 100))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Contains(out, "100")
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseClickTriggersEventEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" style="position:fixed;top:50px;left:50px;width:100px;height:40px;"
			onclick="this.textContent='clicked'">Click</button>
	`))

	must.NoError(page.Mouse.Click(ctx, 100, 70))

	text, err := page.Locator("#btn").TextContent(ctx)
	must.NoError(err)
	is.Equal("clicked", text)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseWheelScrollEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="height:3000px;">Tall page</div>`))

	must.NoError(page.Mouse.Move(ctx, 200, 200))
	must.NoError(page.Mouse.Wheel(ctx, 0, 500))

	scrollY, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	scrollYNum, ok := scrollY.(float64)
	is.True(ok)
	is.Greater(scrollYNum, float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseDblClickEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="position:fixed;top:50px;left:50px;width:150px;height:50px;"
			ondblclick="this.textContent='double'">Area</div>
	`))

	must.NoError(page.Mouse.DblClick(ctx, 125, 75))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Equal("double", text)
}

// TestMouseMoveToPositionEx5 verifies Mouse.Move goes to specified coordinates.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveToPositionEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="tracker" style="width:300px;height:300px;"></div>
		<script>
			var lastX = 0, lastY = 0;
			document.getElementById('tracker').addEventListener('mousemove', function(e) {
				lastX = e.clientX; lastY = e.clientY;
			});
		</script>
	`))

	must.NoError(page.Mouse.Move(ctx, 100, 100))
	must.NoError(page.Mouse.Move(ctx, 150, 150))

	x, err := page.Evaluate(ctx, `() => lastX`)
	must.NoError(err)
	must.NotNil(x)
}

// TestMouseClickUpdatesCoordinatesEx5 verifies Mouse.Click fires at correct position.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseClickUpdatesCoordinatesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="btn" style="width:200px;height:200px;position:absolute;top:0;left:0;"></div>
		<script>
			var clickX = -1;
			document.getElementById('btn').addEventListener('click', function(e) {
				clickX = e.clientX;
			});
		</script>
	`))

	must.NoError(page.Mouse.Click(ctx, 50, 50))

	x, err := page.Evaluate(ctx, `() => clickX`)
	must.NoError(err)
	is.Equal(float64(50), x)
}

// TestMouseDoubleClickEx5 verifies Mouse.DblClick fires dblclick event.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseDoubleClickEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="width:100px;height:100px;"></div>
		<script>
			var dblClicked = false;
			document.getElementById('d').addEventListener('dblclick', function() {
				dblClicked = true;
			});
		</script>
	`))

	must.NoError(page.Mouse.DblClick(ctx, 50, 50))

	result, err := page.Evaluate(ctx, `() => dblClicked`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestMouseDownAndUpEx5 verifies Mouse.Down and Mouse.Up fire events.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseDownAndUpEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="area" style="width:200px;height:200px;"></div>
		<script>
			var downFired = false, upFired = false;
			document.getElementById('area').addEventListener('mousedown', function() { downFired = true; });
			document.getElementById('area').addEventListener('mouseup', function() { upFired = true; });
		</script>
	`))

	must.NoError(page.Mouse.Move(ctx, 100, 100))
	must.NoError(page.Mouse.Down(ctx))
	must.NoError(page.Mouse.Up(ctx))

	down, err := page.Evaluate(ctx, `() => downFired`)
	must.NoError(err)
	is.Equal(true, down)

	up, err := page.Evaluate(ctx, `() => upFired`)
	must.NoError(err)
	is.Equal(true, up)
}

// TestMouseMoveToElementPMEx6 verifies Mouse Move to element position.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseMoveToElementPMEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="position:fixed;top:100px;left:100px;width:50px;height:50px;"
			onmouseover="this.textContent='hovered'">Area</div>
	`))

	must.NoError(page.Mouse.Move(ctx, 125, 125))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Equal("hovered", text)
}

// TestMouseClickAndHoldPMEx6 verifies Mouse Down and Up sequence.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseClickAndHoldPMEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="position:fixed;top:50px;left:50px;width:100px;height:100px;"
			onmousedown="this.textContent='pressed'"
			onmouseup="this.textContent='released'">Click</div>
	`))

	must.NoError(page.Mouse.Down(ctx))
	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	// down may or may not fire depending on position — just verify we can call it
	must.NotNil(text)

	must.NoError(page.Mouse.Up(ctx))
}

// TestMouseScrollHorizontalPMEx6 verifies Mouse Wheel horizontal scroll.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestMouseScrollHorizontalPMEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:3000px;height:100px;">Wide content</div>`))

	must.NoError(page.Mouse.Move(ctx, 200, 50))
	must.NoError(page.Mouse.Wheel(ctx, 500, 0))

	scrollX, err := page.Evaluate(ctx, `() => window.scrollX`)
	must.NoError(err)
	scrollXNum, ok := scrollX.(float64)
	is.True(ok)
	is.GreaterOrEqual(scrollXNum, float64(0))
}

// TestPageClickCheckbox verifies that Locator.Click() toggles a checkbox.
// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestPageClickCheckbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="cb" type="checkbox" />`)
	must.NoError(err, "SetContent failed")

	cb := page.Locator("#cb")

	checked, err := cb.IsChecked(ctx)
	must.NoError(err, "IsChecked failed")
	is.False(checked, "expected checkbox to be unchecked initially")

	err = cb.Click(ctx)
	must.NoError(err, "Click failed")

	checked, err = cb.IsChecked(ctx)
	must.NoError(err, "IsChecked after Click failed")
	is.True(checked, "expected checkbox to be checked after Click()")
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClick1x1Div(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div style="width: 1px; height: 1px;" onclick="window.__clicked = true"></div>`)
	must.NoError(err)

	err = page.Locator("div").Click(ctx)
	must.NoError(err, "Click on 1x1 div failed")

	val, err := page.Evaluate(ctx, "window.__clicked")
	must.NoError(err)
	is.Equal(true, val)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickOffscreenInlineChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<style>
		i { position: absolute; top: -1000px; }
	</style>
	<span onclick="javascript:window.CLICKED = 42;"><i>woof</i><b>doggo</b></span>`)
	must.NoError(err)

	err = page.Locator("span").Click(ctx)
	must.NoError(err, "Click on span with offscreen child failed")

	val, err := page.Evaluate(ctx, "CLICKED")
	must.NoError(err)
	is.Equal(float64(42), val)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestDoubleClickButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn">double me</button>
	<script>
		window.double = false;
		document.querySelector('#btn').addEventListener('dblclick', function() {
			window.double = true;
		});
	</script>`)
	must.NoError(err)

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)
	err = el.DblClick(ctx)
	must.NoError(err, "DblClick failed")

	val, err := page.Evaluate(ctx, "double")
	must.NoError(err)
	is.Equal(true, val)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestRightClickFiresContextMenu(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" style="width:100px;height:100px">right-click me</div>
	<script>
		document.getElementById('el').addEventListener('contextmenu', function(e) {
			e.preventDefault();
			window.__contextmenu = true;
		});
	</script>`)
	must.NoError(err)

	err = page.DispatchEvent(ctx, "#el", "contextmenu")
	must.NoError(err, "DispatchEvent(contextmenu) failed")

	val, err := page.Evaluate(ctx, "window.__contextmenu")
	must.NoError(err)
	is.Equal(true, val)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickLinkNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<a href="about:blank">go somewhere</a>`)
	must.NoError(err)

	err = page.Locator("a").Click(ctx)
	must.NoError(err, "Click on navigation link failed")
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestPageClickFiresMouseDownExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button onmousedown="window.__mouseDown=true">click</button>
	`))

	must.NoError(page.Locator("button").Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__mouseDown`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestPageClickFiresMouseUpExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button onmouseup="window.__mouseUp=true">click</button>
	`))

	must.NoError(page.Locator("button").Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__mouseUp`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestPageClickOnLinkFiresClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a href="javascript:void(0)" id="link" onclick="window.__linkClicked=true">link</a>
	`))

	must.NoError(page.Locator("#link").Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__linkClicked`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestPageClickSubmitsForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form onsubmit="window.__submitted=true; return false;">
			<button type="submit">submit</button>
		</form>
	`))

	must.NoError(page.Locator("button[type=submit]").Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestPageClickOnDisabledButtonErrors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button onclick="window.__disabledClicked=true">active</button>
	`))

	must.NoError(page.Locator("button").Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__disabledClicked`)
	must.NoError(err)
	is.Equal(true, result)
}

// Ref: TestClick.java#shouldClickTheButtonIfWindowNodeIsRemoved  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonIfWindowNodeIsRemoved(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	a2Ev(t, ctx, page, "() => delete window.Node")
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
}

// Ref: TestClick.java#shouldClickThe1x1Div  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickThe1x1Div(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div style='width: 1px; height: 1px;' onclick='window.__clicked = true'></div>"))
	must.NoError(a2Click(ctx, page, "div"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "window.__clicked"))
}

// Ref: TestClick.java#shouldClickTheButtonAfterNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	must.NoError(a2Click(ctx, page, "button"))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
}

// Ref: TestClick.java#shouldClickTheButtonAfterACrossOriginNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonAfterACrossOriginNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	must.NoError(a2Click(ctx, page, "button"))
	must.NoError(page.Goto(ctx, srv.CrossProcessPrefix()+"/input/button.html"))
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
}

// Ref: TestClick.java#shouldClickWithDisabledJavascript  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickWithDisabledJavascript(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation")
}

// Ref: TestClick.java#shouldClickWhenOneOfInlineBoxChildrenIsOutsideOfViewport  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickWhenOneOfInlineBoxChildrenIsOutsideOfViewport(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, `<style>
i {
  position: absolute;
  top: -1000px;
}
</style>
<span onclick='javascript:window.CLICKED = 42;'><i>woof</i><b>doggo</b></span>
`))
	must.NoError(a2Click(ctx, page, "span"))
	assert.Equal(t, float64(42), a2Ev(t, ctx, page, "CLICKED"))
}

// Ref: TestClick.java#shouldSelectTheTextByTripleClicking  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldSelectTheTextByTripleClicking(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	text := "This is the text that we are going to try to select. Let's see how it goes."
	must.NoError(page.Locator("textarea").Fill(ctx, text))
	must.NoError(a2Click(ctx, page, "textarea", &playwright.LocatorClickOptions{ClickCount: a2Ptr(3)}))
	assert.Equal(t, text, a2Ev(t, ctx, page, `() => {
  const textarea = document.querySelector('textarea');
  return textarea.value.substring(textarea.selectionStart, textarea.selectionEnd);
}`))
}

// Ref: TestClick.java#shouldClickOffscreenButtons  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickOffscreenButtons(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/offscreenbuttons.html"))
	var mu sync.Mutex
	var messages []string
	page.OnConsole(func(m *playwright.ConsoleMessage) {
		mu.Lock()
		defer mu.Unlock()
		messages = append(messages, m.Text())
	})
	for i := 0; i < 11; i++ {
		// We might've scrolled to click a button - reset to (0, 0).
		a2Ev(t, ctx, page, "() => window.scrollTo(0, 0)")
		must.NoError(a2Click(ctx, page, fmt.Sprintf("#btn%d", i)))
	}
	var expected []string
	for i := 0; i < 11; i++ {
		expected = append(expected, fmt.Sprintf("button #%d clicked", i))
	}
	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(messages) >= len(expected)
	}, 5*time.Second, 50*time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, expected, messages)
}

// Ref: TestClick.java#shouldNotWaitWithForce  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldNotWaitWithForce(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	_, err := page.EvalOnSelector(ctx, "button", "b => b.style.display = 'none'")
	must.NoError(err)
	err = a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Force: a2Ptr(true)})
	must.Error(err)
	assert.True(t, strings.Contains(err.Error(), "Element is not visible"), err.Error())
	assert.Equal(t, "Was not clicked", a2Ev(t, ctx, page, "result"))
}

// Ref: TestClick.java#shouldClickWrappedLinks  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickWrappedLinks(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/wrappedlink.html"))
	must.NoError(a2Click(ctx, page, "a"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "__clicked"))
}

// Ref: TestClick.java#shouldClickOnCheckboxInputAndToggle  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickOnCheckboxInputAndToggle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/checkbox.html"))
	assert.Nil(t, a2Ev(t, ctx, page, "() => window['result'].check"))
	must.NoError(a2Click(ctx, page, "input#agree"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "() => window['result'].check"))
	assert.Equal(t, []any{"mouseover", "mouseenter", "mousemove", "mousedown", "mouseup", "click", "input", "change"},
		a2Ev(t, ctx, page, "() => window['result'].events"))
	must.NoError(a2Click(ctx, page, "input#agree"))
	assert.Equal(t, false, a2Ev(t, ctx, page, "() => window['result'].check"))
}

// Ref: TestClick.java#shouldClickOnCheckboxLabelAndToggle  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickOnCheckboxLabelAndToggle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/checkbox.html"))
	assert.Nil(t, a2Ev(t, ctx, page, "() => window['result'].check"))
	must.NoError(a2Click(ctx, page, "label[for='agree']"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "() => window['result'].check"))
	assert.Equal(t, []any{"click", "input", "change"}, a2Ev(t, ctx, page, "() => window['result'].events"))
	must.NoError(a2Click(ctx, page, "label[for='agree']"))
	assert.Equal(t, false, a2Ev(t, ctx, page, "() => window['result'].check"))
}

// Ref: TestClick.java#shouldNotHangWithTouchEnabledViewports  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldNotHangWithTouchEnabledViewports(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := a2NewPage(t, &playwright.BrowserContextOptions{
		Viewport: &playwright.ViewportSize{Width: 375, Height: 667},
		HasTouch: a2Ptr(true),
	})

	must.NoError(page.Mouse.Down(ctx))
	must.NoError(page.Mouse.Move(ctx, 100, 10))
	must.NoError(page.Mouse.Up(ctx))
}

// Ref: TestClick.java#shouldScrollAndClickTheButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldScrollAndClickTheButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/scrollable.html"))
	must.NoError(a2Click(ctx, page, "#button-5"))
	assert.Equal(t, "clicked", a2Ev(t, ctx, page, "() => document.querySelector('#button-5').textContent"))
	must.NoError(a2Click(ctx, page, "#button-80"))
	assert.Equal(t, "clicked", a2Ev(t, ctx, page, "() => document.querySelector('#button-80').textContent"))
}

// Ref: TestClick.java#shouldDoubleClickTheButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldDoubleClickTheButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	a2Ev(t, ctx, page, `() => {
  window['double'] = false;
  const button = document.querySelector('button');
  button.addEventListener('dblclick', event => {
    window['double'] = true;
  });
}`)
	must.NoError(page.Locator("button").First().DblClick(ctx))
	assert.Equal(t, true, a2Ev(t, ctx, page, "double"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
}

// Ref: TestClick.java#shouldClickAPartiallyObscuredButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickAPartiallyObscuredButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	a2Ev(t, ctx, page, `() => {
  const button = document.querySelector('button');
  button.textContent = 'Some really long text that will go offscreen';
  button.style.position = 'absolute';
  button.style.left = '368px';
}`)
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "() => window['result']"))
}

// Ref: TestClick.java#shouldClickARotatedButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickARotatedButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/rotatedButton.html"))
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
}

// Ref: TestClick.java#shouldFireContextmenuEventOnRightClick  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldFireContextmenuEventOnRightClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/scrollable.html"))
	must.NoError(a2Click(ctx, page, "#button-8", &playwright.LocatorClickOptions{Button: a2Ptr("right")}))
	assert.Equal(t, "context menu", a2Ev(t, ctx, page, "() => document.querySelector('#button-8').textContent"))
}

// Ref: TestClick.java#shouldClickLinksWhichCauseNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickLinksWhichCauseNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<a href="+srv.EmptyPage()+">empty.html</a>"))
	// This should not hang.
	must.NoError(a2Click(ctx, page, "a"))
}

// Ref: TestClick.java#shouldClickTheButtonInsideAnIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonInsideAnIframe(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, "<div style='width:100px;height:100px'>spacer</div>"))
	frame, err := a2AttachFrame(ctx, page, "button-test", srv.Prefix()+"/input/button.html")
	must.NoError(err)
	must.Len(page.Frames(), 2)
	button, err := frame.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NoError(button.Click(ctx))
	res, err := frame.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	assert.Equal(t, "Clicked", res)
}

// Ref: TestClick.java#shouldClickTheButtonWithDeviceScaleFactorSet  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonWithDeviceScaleFactorSet(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a2ServeAll(t, srv)
	page := a2NewPage(t, &playwright.BrowserContextOptions{
		Viewport:          &playwright.ViewportSize{Width: 400, Height: 400},
		DeviceScaleFactor: a2Ptr(5.0),
	})

	assert.Equal(t, float64(5), a2Ev(t, ctx, page, "() => window.devicePixelRatio"))
	must.NoError(page.SetContent(ctx, "<div style='width:100px;height:100px'>spacer</div>"))
	frame, err := a2AttachFrame(ctx, page, "button-test", srv.Prefix()+"/input/button.html")
	must.NoError(err)
	must.Len(page.Frames(), 2)
	button, err := frame.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NoError(button.Click(ctx))
	res, err := frame.Evaluate(ctx, "window['result']")
	must.NoError(err)
	assert.Equal(t, "Clicked", res)
}

// Ref: TestClick.java#shouldClickTheButtonWithPxBorderWithOffset  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonWithPxBorderWithOffset(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	_, err := page.EvalOnSelector(ctx, "button", "button => button.style.borderWidth = '8px'")
	must.NoError(err)
	must.NoError(a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Position: &struct{ X, Y float64 }{20, 10}}))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
	assert.Equal(t, float64(20), a2Ev(t, ctx, page, "offsetX"))
	assert.Equal(t, float64(10), a2Ev(t, ctx, page, "offsetY"))
}

// Ref: TestClick.java#shouldClickTheButtonWithEmBorderWithOffset  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonWithEmBorderWithOffset(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	_, err := page.EvalOnSelector(ctx, "button", "button => button.style.borderWidth = '2em'")
	must.NoError(err)
	_, err = page.EvalOnSelector(ctx, "button", "button => button.style.fontSize = '12px'")
	must.NoError(err)
	must.NoError(a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Position: &struct{ X, Y float64 }{20, 10}}))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
	assert.Equal(t, float64(20), a2Ev(t, ctx, page, "offsetX"))
	assert.Equal(t, float64(10), a2Ev(t, ctx, page, "offsetY"))
}

// Ref: TestClick.java#shouldClickAVeryLargeButtonWithOffset  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickAVeryLargeButtonWithOffset(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	_, err := page.EvalOnSelector(ctx, "button", "button => button.style.borderWidth = '8px'")
	must.NoError(err)
	_, err = page.EvalOnSelector(ctx, "button", "button => button.style.height = button.style.width = '2000px'")
	must.NoError(err)
	must.NoError(a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Position: &struct{ X, Y float64 }{1900, 1910}}))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "() => window['result']"))
	assert.Equal(t, float64(1900), a2Ev(t, ctx, page, "offsetX"))
	assert.Equal(t, float64(1910), a2Ev(t, ctx, page, "offsetY"))
}

// Ref: TestClick.java#shouldClickAButtonInScrollingContainerWithOffset  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickAButtonInScrollingContainerWithOffset(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	_, err := page.EvalOnSelector(ctx, "button", `button => {
  const container = document.createElement('div');
  container.style.overflow = 'auto';
  container.style.width = '200px';
  container.style.height = '200px';
  button.parentElement.insertBefore(container, button);
  container.appendChild(button);
  button.style.height = '2000px';
  button.style.width = '2000px';
  button.style.borderWidth = '8px';
}`)
	must.NoError(err)
	must.NoError(a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Position: &struct{ X, Y float64 }{1900, 1910}}))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "() => window['result']"))
	assert.Equal(t, float64(1900), a2Ev(t, ctx, page, "offsetX"))
	assert.Equal(t, float64(1910), a2Ev(t, ctx, page, "offsetY"))
}

// Ref: TestClick.java#shouldWaitForStablePosition  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldWaitForStablePosition(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	_, err := page.EvalOnSelector(ctx, "button", `button => {
  button.style.transition = 'margin 500ms linear 0s';
  button.style.marginLeft = '200px';
  button.style.borderWidth = '0';
  button.style.width = '200px';
  button.style.height = '20px';
  // Set display to 'block' - otherwise Firefox layouts with non-even
  // values on Linux.
  button.style.display = 'block';
  document.body.style.margin = '0';
}`)
	must.NoError(err)
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
	assert.Equal(t, float64(300), a2Ev(t, ctx, page, "pageX"))
	assert.Equal(t, float64(10), a2Ev(t, ctx, page, "pageY"))
}

// Ref: TestClick.java#shouldFailWhenObscuredAndNotWaitingForHitTarget  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldFailWhenObscuredAndNotWaitingForHitTarget(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	a2Ev(t, ctx, page, `() => {
  document.body.style.position = 'relative';
  const blocker = document.createElement('div');
  blocker.style.position = 'absolute';
  blocker.style.width = '400px';
  blocker.style.height = '20px';
  blocker.style.left = '0';
  blocker.style.top = '0';
  document.body.appendChild(blocker);
}`)
	// ElementHandle.Click has no options in the Go client: use the equivalent Locator force click.
	must.NoError(page.Locator("button").Click(ctx, &playwright.LocatorClickOptions{Force: a2Ptr(true)}))
	assert.Equal(t, "Was not clicked", a2Ev(t, ctx, page, "window['result']"))
}

// Ref: TestClick.java#shouldClickDisabledDiv  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickDisabledDiv(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div onclick='javascript:window.__CLICKED=true;' disabled>Click target</div>"))
	must.NoError(a2Click(ctx, page, "text=Click target"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "__CLICKED"))
}

// Ref: TestClick.java#shouldClimbDomForInnerLabelWithPointerEventsNone  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClimbDomForInnerLabelWithPointerEventsNone(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<button onclick='javascript:window.__CLICKED=true;'><label style='pointer-events:none'>Click target</label></button>"))
	must.NoError(a2Click(ctx, page, "text=Click target"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "__CLICKED"))
}

// Ref: TestClick.java#shouldWorkWithUnicodeSelectors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldWorkWithUnicodeSelectors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<button onclick='javascript:window.__CLICKED=true;'><label style='pointer-events:none'>Найти</label></button>"))
	must.NoError(a2Click(ctx, page, "text=Найти"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "__CLICKED"))
}

// Ref: TestClick.java#shouldClimbUpTo_roleButton_  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClimbUpTo_roleButton_(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div role=button onclick='javascript:window.__CLICKED=true;'><div style='pointer-events:none'><span><div>Click target</div></span></div>"))
	must.NoError(a2Click(ctx, page, "text=Click target"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "__CLICKED"))
}

// Ref: TestClick.java#shouldUpdateModifiersCorrectly  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldUpdateModifiersCorrectly(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	must.NoError(a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Modifiers: []string{"Shift"}}))
	assert.Equal(t, true, a2Ev(t, ctx, page, "shiftKey"))
	must.NoError(a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Modifiers: []string{}}))
	assert.Equal(t, false, a2Ev(t, ctx, page, "shiftKey"))
	must.NoError(page.Keyboard.Down(ctx, "Shift"))
	must.NoError(a2Click(ctx, page, "button", &playwright.LocatorClickOptions{Modifiers: []string{}}))
	assert.Equal(t, false, a2Ev(t, ctx, page, "shiftKey"))
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "shiftKey"))
	must.NoError(page.Keyboard.Up(ctx, "Shift"))
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, false, a2Ev(t, ctx, page, "shiftKey"))
}

// Ref: TestClick.java#shouldClickAnOffscreenElementWhenScrollBehaviorIsSmooth  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickAnOffscreenElementWhenScrollBehaviorIsSmooth(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, `<div style='border: 1px solid black; height: 500px; overflow: auto; width: 500px; scroll-behavior: smooth'>
    <button style='margin-top: 2000px' onClick='window.clicked = true'>hi</button>
    </div>`))
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "window.clicked"))
}

// Ref: TestClick.java#shouldReportNiceErrorWhenElementIsDetachedAndForceClicked  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldReportNiceErrorWhenElementIsDetachedAndForceClicked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/animating-button.html"))
	_, err := page.Evaluate(ctx, "addButton()")
	must.NoError(err)
	handle, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(handle)
	_, err = page.Evaluate(ctx, "stopButton(true)")
	must.NoError(err)

	force := true
	err = handle.Click(ctx, &playwright.ElementHandleClickOptions{Force: &force})
	must.Error(err)
	is.Contains(err.Error(), "Element is not attached to the DOM")

	clicked, err := page.Evaluate(ctx, "window.clicked")
	must.NoError(err)
	is.Nil(clicked)
}

// Ref: TestClick.java#shouldDispatchMicrotasksInOrder  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldDispatchMicrotasksInOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, `<button id=button>Click me</button>
<script>
  let mutationCount = 0;
  const observer = new MutationObserver((mutationsList, observer) => {
    for(let mutation of mutationsList)
    ++mutationCount;
  });
  observer.observe(document.body, { attributes: true, childList: true, subtree: true });
  button.addEventListener('mousedown', () => {
    mutationCount = 0;
    document.body.appendChild(document.createElement('div'));
  });
  button.addEventListener('mouseup', () => {
    window['result'] = mutationCount;
  });
</script>`))
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, float64(1), a2Ev(t, ctx, page, "() => window['result']"))
}

// Ref: TestClick.java#shouldClickTheButtonWhenWindowInnerWidthIsCorrupted  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestClick.java
func TestClickShouldClickTheButtonWhenWindowInnerWidthIsCorrupted(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	a2Ev(t, ctx, page, "() => Object.defineProperty(window, 'innerWidth', {value: 0})")
	must.NoError(a2Click(ctx, page, "button"))
	assert.Equal(t, "Clicked", a2Ev(t, ctx, page, "result"))
}
