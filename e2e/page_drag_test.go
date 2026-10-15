//go:build e2e

// E2E tests for: TestPageDrag.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestLocatorDragTo(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="source" draggable="true"
			style="width:50px;height:50px;background:red;position:absolute;left:0;top:0">
			drag
		</div>
		<div id="target"
			style="width:100px;height:100px;background:blue;position:absolute;left:200px;top:0">
			drop here
		</div>
		<script>
			var source = document.getElementById('source');
			var target = document.getElementById('target');
			source.addEventListener('dragstart', function(e) {
				e.dataTransfer.setData('text/plain', 'dragged');
			});
			target.addEventListener('dragover', function(e) { e.preventDefault(); });
			target.addEventListener('drop', function(e) {
				e.preventDefault();
				target.setAttribute('data-dropped', e.dataTransfer.getData('text/plain'));
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#source").DragTo(ctx, page.Locator("#target"))
	must.NoError(err, "DragTo failed")

	attr, err := page.Locator("#target").GetAttribute(ctx, "data-dropped")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "dragged" {
		t.Errorf("expected data-dropped='dragged', got %v", attr)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragTo(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="src" draggable="true"
			style="width:50px;height:50px;background:green;position:absolute;left:0;top:0">
			src
		</div>
		<div id="tgt"
			style="width:100px;height:100px;background:yellow;position:absolute;left:200px;top:0">
			tgt
		</div>
		<script>
			document.getElementById('src').addEventListener('dragstart', function(e) {
				e.dataTransfer.setData('text/plain', 'ok');
			});
			document.getElementById('tgt').addEventListener('dragover', function(e) { e.preventDefault(); });
			document.getElementById('tgt').addEventListener('drop', function(e) {
				e.preventDefault();
				this.setAttribute('data-got', e.dataTransfer.getData('text/plain'));
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.DragTo(ctx, "#src", "#tgt")
	must.NoError(err, "Page.DragTo failed")

	attr, err := page.Locator("#tgt").GetAttribute(ctx, "data-got")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "ok" {
		t.Errorf("expected data-got='ok', got %v", attr)
	}
}

// TestLocatorDragToNoErrorEx2 verifies DragTo does not error.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestLocatorDragToNoErrorEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="drag" style="width:50px;height:50px;position:absolute;left:10px;top:10px;background:red"></div>
		<div id="drop" style="width:100px;height:100px;position:absolute;left:200px;top:200px;background:blue"></div>
	`))

	must.NoError(page.Locator("#drag").DragTo(ctx, page.Locator("#drop")))
}

// TestLocatorDragToFiresDragEventsEx2 verifies DragTo fires drag events.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestLocatorDragToFiresDragEventsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="drag" draggable="true" style="width:50px;height:50px;position:absolute;left:10px;top:10px"
		     ondragstart="window.__dragStarted=true"></div>
		<div id="drop" style="width:100px;height:100px;position:absolute;left:200px;top:200px"
		     ondragover="event.preventDefault()"
		     ondrop="window.__dropped=true"></div>
	`))

	must.NoError(page.Locator("#drag").DragTo(ctx, page.Locator("#drop")))

	started, err := page.Evaluate(ctx, `() => window.__dragStarted`)
	must.NoError(err)
	is.Equal(true, started)
}

// TestLocatorDragToSameContainerEx2 verifies DragTo within same container.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestLocatorDragToSameContainerEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="container" style="width:500px;height:500px;position:relative">
			<div id="source" style="width:50px;height:50px;position:absolute;left:10px;top:10px;background:green"></div>
			<div id="target" style="width:50px;height:50px;position:absolute;left:400px;top:400px;background:yellow"></div>
		</div>
	`))

	must.NoError(page.Locator("#source").DragTo(ctx, page.Locator("#target")))
}

// TestLocatorDragToFiresDragEvents verifies DragTo fires drag-related events.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestLocatorDragToFiresDragEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="source" style="width:50px;height:50px;background:red;position:absolute;top:10px;left:10px;" draggable="true"
			ondragstart="document.getElementById('log').textContent += 'dragstart;'">source</div>
		<div id="target" style="width:100px;height:100px;background:blue;position:absolute;top:10px;left:200px;"
			ondrop="document.getElementById('log').textContent += 'drop;'"
			ondragover="event.preventDefault()">target</div>
		<div id="log"></div>
	`))

	must.NoError(page.Locator("#source").DragTo(ctx, page.Locator("#target")))

	logText, err := page.Locator("#log").InnerText(ctx)
	must.NoError(err)
	is.Contains(logText, "dragstart")
}

// TestLocatorDragToUpdatesPosition verifies DragTo changes element visual position.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestLocatorDragToUpdatesPosition(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))
	must.NoError(page.SetContent(ctx, `
		<div id="source" style="width:50px;height:50px;background:red;position:absolute;top:10px;left:10px;">source</div>
		<div id="target" style="width:100px;height:100px;background:blue;position:absolute;top:100px;left:300px;">target</div>
	`))

	// DragTo should not error even if the actual visual drag depends on browser behavior
	must.NoError(page.Locator("#source").DragTo(ctx, page.Locator("#target")))
}

// TestPageDragToMovesElement verifies page.DragTo moves an element.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragToMovesElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="source" draggable="true"
		     style="width:50px;height:50px;background:red;position:absolute;left:0;top:0">
		</div>
		<div id="target"
		     style="width:100px;height:100px;background:blue;position:absolute;left:200px;top:200px"
		     ondragover="event.preventDefault()"
		     ondrop="this.textContent='dropped'">
		</div>
	`))

	must.NoError(page.DragTo(ctx, "#source", "#target"))

	text, err := page.Locator("#target").InnerText(ctx)
	must.NoError(err)
	is.Equal("dropped", text)
}

// TestPageDragToFiresDragEvents verifies drag events fire during page.DragTo.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragToFiresDragEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="src" draggable="true"
		     style="width:50px;height:50px;background:red;position:absolute;left:0;top:0">
		</div>
		<div id="dst"
		     style="width:100px;height:100px;background:blue;position:absolute;left:200px;top:200px"
		     ondragover="event.preventDefault()"
		     ondrop="window.__dropped=true">
		</div>
		<script>
		  document.getElementById('src').addEventListener('dragstart', () => window.__started = true);
		</script>
	`))

	must.NoError(page.DragTo(ctx, "#src", "#dst"))

	started, err := page.Evaluate(ctx, `() => window.__started`)
	must.NoError(err)
	is.Equal(true, started)

	dropped, err := page.Evaluate(ctx, `() => window.__dropped`)
	must.NoError(err)
	is.Equal(true, dropped)
}

// TestPageDragToFromSourceToTarget verifies page.DragTo transfers data via dataTransfer.
// Extra: no equivalent method in playwright-java; related class: TestPageDrag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragToTransfersData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="src" draggable="true"
		     style="width:50px;height:50px;position:absolute;left:0;top:0">
		</div>
		<div id="dst"
		     style="width:100px;height:100px;position:absolute;left:200px;top:200px"
		     ondragover="event.preventDefault()"
		     ondrop="window.__data = event.dataTransfer.getData('text/plain')">
		</div>
		<script>
		  document.getElementById('src').addEventListener('dragstart', e => {
		    e.dataTransfer.setData('text/plain', 'hello');
		  });
		</script>
	`))

	must.NoError(page.DragTo(ctx, "#src", "#dst"))

	data, err := page.Evaluate(ctx, `() => window.__data`)
	must.NoError(err)
	is.Equal("hello", data)
}

// TestPageDragShouldWorkWithLocators verifies that And() narrows to elements matching both locators.
// Ref: TestPageDrag.java#shouldWorkWithLocators  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldWorkWithLocators(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button class="primary">Save</button>
		<button class="secondary">Cancel</button>
		<button class="primary highlighted">Submit</button>
	`)
	must.NoError(err, "SetContent failed")

	// All primary buttons
	primary := page.Locator("button.primary")
	count, err := primary.Count(ctx)
	if err != nil || count != 2 {
		t.Fatalf("expected 2 primary buttons, got %d (err: %v)", count, err)
	}

	// Primary AND highlighted → only the one with both classes
	highlighted := page.Locator("button.highlighted")
	both := primary.And(highlighted)

	count, err = both.Count(ctx)
	must.NoError(err, "And().Count() failed")
	is.Equal(1, count)

	text, err := both.InnerText(ctx)
	must.NoError(err, "And().InnerText() failed")
	is.Equal("Submit", text)
}

// Ref: TestPageDrag.java#shouldWorkIfTheDragIsCanceled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldWorkIfTheDragIsCanceled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/drag-n-drop.html"))
	a2Ev(t, ctx, page, `() => {
      document.body.addEventListener('dragstart', event => {
        event.preventDefault();
      }, false);
    }`)
	must.NoError(page.Locator("#source").Hover(ctx))
	must.NoError(page.Mouse.Down(ctx))
	must.NoError(page.Locator("#target").Hover(ctx))
	must.NoError(page.Mouse.Up(ctx))
	res, err := page.EvalOnSelector(ctx, "#target", "target => target.contains(document.querySelector('#source'))")
	must.NoError(err)
	assert.Equal(t, false, res)
}

// Ref: TestPageDrag.java#shouldWorkIfTheDragEventIsCapturedButNotCanceled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldWorkIfTheDragEventIsCapturedButNotCanceled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/drag-n-drop.html"))
	a2Ev(t, ctx, page, `() => {
      document.body.addEventListener('dragstart', event => {
        event.stopImmediatePropagation();
      }, false);
    }`)
	must.NoError(page.Locator("#source").Hover(ctx))
	must.NoError(page.Mouse.Down(ctx))
	must.NoError(page.Locator("#target").Hover(ctx))
	must.NoError(page.Mouse.Up(ctx))
	res, err := page.EvalOnSelector(ctx, "#target", "target => target.contains(document.querySelector('#source'))")
	must.NoError(err)
	assert.Equal(t, true, res)
}

// Ref: TestPageDrag.java#shouldBeAbleToDragTheMouseInAFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldBeAbleToDragTheMouseInAFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/frames/one-frame.html"))
	html, err := page.Frames()[1].QuerySelector(ctx, "html")
	must.NoError(err)
	eventsHandle, err := html.EvaluateHandle(ctx, `target => {
  const events = [];
  for (const event of [
    'mousedown', 'mousemove', 'mouseup',
    'dragstart', 'dragend', 'dragover', 'dragenter', 'dragleave', 'dragexit',
    'drop'
  ])
  target.addEventListener(event, () => events.push(event), false);
  return events;
}`)
	must.NoError(err)
	must.NoError(page.Mouse.Move(ctx, 30, 30))
	must.NoError(page.Mouse.Down(ctx))
	must.NoError(page.Mouse.Move(ctx, 60, 60))
	must.NoError(page.Mouse.Up(ctx))
	events, err := eventsHandle.JSONValue(ctx)
	must.NoError(err)
	assert.Equal(t, []any{"mousemove", "mousedown", "mousemove", "mouseup"}, events)
}

// Ref: TestPageDrag.java#shouldWorkWithTheHelperMethod  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldWorkWithTheHelperMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/drag-n-drop.html"))
	must.NoError(page.DragTo(ctx, "#source", "#target"))
	res, err := page.EvalOnSelector(ctx, "#target", "target => target.contains(document.querySelector('#source'))")
	must.NoError(err)
	assert.Equal(t, true, res)
}

// Ref: TestPageDrag.java#shouldAllowSpecifyingThePosition  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldAllowSpecifyingThePosition(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div style='width:100px;height:100px;background:red;' id='red'>\n</div>\n<div style='width:100px;height:100px;background:blue;' id='blue'>\n</div>"))
	eventsHandle, err := page.EvaluateHandle(ctx, `() => {
      const events = [];
      document.getElementById('red').addEventListener('mousedown', event => {
        events.push({
          type: 'mousedown',
          x: event.offsetX,
          y: event.offsetY,
        });
      });
      document.getElementById('blue').addEventListener('mouseup', event => {
        events.push({
          type: 'mouseup',
          x: event.offsetX,
          y: event.offsetY,
        });
      });
      return events;
    }`)
	must.NoError(err)
	must.NoError(page.DragTo(ctx, "#red", "#blue", &playwright.DragAndDropOptions{
		SourcePosition: &playwright.Position{X: 34, Y: 7},
		TargetPosition: &playwright.Position{X: 10, Y: 20},
	}))
	jsonVal, err := eventsHandle.JSONValue(ctx)
	must.NoError(err)
	expected := []any{
		map[string]any{"type": "mousedown", "x": float64(34), "y": float64(7)},
		map[string]any{"type": "mouseup", "x": float64(10), "y": float64(20)},
	}
	is.Equal(expected, jsonVal)
}

// Ref: TestPageDrag.java#shouldDragAndDropWithTweenedMouseMovement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldDragAndDropWithTweenedMouseMovement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body style='margin: 0; padding: 0;'>
  <div style='width:100px;height:100px;background:red;' id='red'></div>
  <div style='width:300px;height:100px;background:blue;' id='blue'></div>
</body>`))

	eventsHandle, err := page.EvaluateHandle(ctx, `() => {
  const events = [];
  document.addEventListener('mousedown', event => {
    events.push({ type: 'mousedown', x: event.pageX, y: event.pageY });
  });
  document.addEventListener('mouseup', event => {
    events.push({ type: 'mouseup', x: event.pageX, y: event.pageY });
  });
  document.addEventListener('mousemove', event => {
    events.push({ type: 'mousemove', x: event.pageX, y: event.pageY });
  });
  return events;
}`)
	must.NoError(err)

	steps := 4
	must.NoError(page.DragTo(ctx, "#red", "#blue", &playwright.DragAndDropOptions{Steps: &steps}))

	jsonVal, err := eventsHandle.JSONValue(ctx)
	must.NoError(err)
	expected := []any{
		map[string]any{"type": "mousemove", "x": float64(50), "y": float64(50)},
		map[string]any{"type": "mousedown", "x": float64(50), "y": float64(50)},
		map[string]any{"type": "mousemove", "x": float64(75), "y": float64(75)},
		map[string]any{"type": "mousemove", "x": float64(100), "y": float64(100)},
		map[string]any{"type": "mousemove", "x": float64(125), "y": float64(125)},
		map[string]any{"type": "mousemove", "x": float64(150), "y": float64(150)},
		map[string]any{"type": "mouseup", "x": float64(150), "y": float64(150)},
	}
	is.Equal(expected, jsonVal)
}

// Ref: TestPageDrag.java#shouldDragToWithTweenedMouseMovement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrag.java
func TestPageDragShouldDragToWithTweenedMouseMovement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body style='margin: 0; padding: 0;'>
  <div style='width:100px;height:100px;background:red;' id='red'></div>
  <div style='width:300px;height:100px;background:blue;' id='blue'></div>
</body>`))

	eventsHandle, err := page.EvaluateHandle(ctx, `() => {
  const events = [];
  document.addEventListener('mousedown', event => {
    events.push({ type: 'mousedown', x: event.pageX, y: event.pageY });
  });
  document.addEventListener('mouseup', event => {
    events.push({ type: 'mouseup', x: event.pageX, y: event.pageY });
  });
  document.addEventListener('mousemove', event => {
    events.push({ type: 'mousemove', x: event.pageX, y: event.pageY });
  });
  return events;
}`)
	must.NoError(err)

	steps := 4
	must.NoError(page.Locator("#red").DragTo(ctx, page.Locator("#blue"), &playwright.LocatorDragToOptions{Steps: &steps}))

	jsonVal, err := eventsHandle.JSONValue(ctx)
	must.NoError(err)
	expected := []any{
		map[string]any{"type": "mousemove", "x": float64(50), "y": float64(50)},
		map[string]any{"type": "mousedown", "x": float64(50), "y": float64(50)},
		map[string]any{"type": "mousemove", "x": float64(75), "y": float64(75)},
		map[string]any{"type": "mousemove", "x": float64(100), "y": float64(100)},
		map[string]any{"type": "mousemove", "x": float64(125), "y": float64(125)},
		map[string]any{"type": "mousemove", "x": float64(150), "y": float64(150)},
		map[string]any{"type": "mouseup", "x": float64(150), "y": float64(150)},
	}
	is.Equal(expected, jsonVal)
}
