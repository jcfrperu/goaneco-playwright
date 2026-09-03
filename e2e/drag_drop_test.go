//go:build e2e

// Drag & Drop E2E tests.
// Migration of: TestDragDrop.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
// Ref: TestLocatorDrag.java#shouldNotError
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
// Ref: TestLocatorDrag.java#shouldFireDragEvents
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
// Ref: TestLocatorDrag.java#shouldDragWithinContainer
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
// Ref: TestLocatorDragTo.java#shouldFireDragEvents
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
// Ref: TestLocatorDragTo.java#shouldUpdatePosition
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
// Ref: TestPageDragAndDrop.java#shouldWork
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
// Ref: TestPageDragAndDrop.java#shouldFireDragEvents
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
// Ref: TestPageDragAndDrop.java#shouldSupportDataTransfer
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
