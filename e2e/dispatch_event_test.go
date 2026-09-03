//go:build e2e

// Page.DispatchEvent and ElementHandle.DispatchEvent E2E tests.
// Migration of: TestDispatchEvent.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageDispatchEventClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="btn">Click me</button>
		<script>
			document.getElementById('btn').addEventListener('click', function() {
				this.setAttribute('data-clicked', '1');
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.DispatchEvent(ctx, "#btn", "click")
	must.NoError(err, "DispatchEvent(click) failed")

	attr, err := page.Locator("#btn").GetAttribute(ctx, "data-clicked")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-clicked='1' after dispatchEvent, got %v", attr)
	}
}

func TestPageDispatchEventCustom(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="el">target</div>
		<script>
			document.getElementById('el').addEventListener('myevent', function(e) {
				this.setAttribute('data-received', '1');
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.DispatchEvent(ctx, "#el", "myevent")
	must.NoError(err, "DispatchEvent(myevent) failed")

	attr, err := page.Locator("#el").GetAttribute(ctx, "data-received")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-received='1', got %v", attr)
	}
}

func TestElementHandleDispatchEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="btn">Click</button>
		<script>
			document.getElementById('btn').addEventListener('click', function() {
				this.setAttribute('data-fired', '1');
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#btn")
	if err != nil || el == nil {
		t.Fatalf("QuerySelector failed: %v", err)
	}

	err = el.DispatchEvent(ctx, "click")
	must.NoError(err, "ElementHandle.DispatchEvent(click) failed")

	attr, err := page.Locator("#btn").GetAttribute(ctx, "data-fired")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-fired='1', got %v", attr)
	}
}

func TestPageDispatchEventWithEventInit(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Dispatch a MouseEvent with clientX set via eventInit and verify event was received.
	err := page.SetContent(ctx, `
		<div id="el">target</div>
		<script>
			document.getElementById('el').addEventListener('mousemove', function(e) {
				this.setAttribute('data-received', '1');
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.DispatchEvent(ctx, "#el", "mousemove", map[string]any{"clientX": 100.0})
	must.NoError(err, "DispatchEvent with eventInit failed")

	attr, err := page.Locator("#el").GetAttribute(ctx, "data-received")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-received='1', got %v", attr)
	}
}

func TestDispatchEventSVGClick(t *testing.T) {
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

	err = page.DispatchEvent(ctx, "circle", "click")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => window['__CLICKED']")
	must.NoError(err)
	is.Equal(float64(42), val)
}

func TestDispatchEventSpanInlineElement(t *testing.T) {
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

	err = page.DispatchEvent(ctx, "span", "click")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => window['CLICKED']")
	must.NoError(err)
	is.Equal(float64(42), val)
}

func TestDispatchEventNotFailWhenBlockedOnHover(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<style>
		container { display: block; position: relative; width: 200px; height: 50px; }
		div, button { position: absolute; left: 0; top: 0; bottom: 0; right: 0; }
		div { pointer-events: none; }
		container:hover div { pointer-events: auto; background: red; }
	</style>
	<container>
		<button onclick="window.clicked=true">Click me</button>
		<div></div>
	</container>`)
	must.NoError(err)

	err = page.DispatchEvent(ctx, "button", "click")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => window['clicked']")
	must.NoError(err)
	is.Equal(true, val)
}

func TestDispatchEventShadowDOM(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
		const div = document.createElement('div');
		div.attachShadow({mode: 'open'});
		document.body.appendChild(div);
	}`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => new Promise(f => setTimeout(f, 100))`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
		const span = document.createElement('span');
		span.textContent = 'Hello from shadow';
		span.addEventListener('click', () => window['clicked'] = true);
		document.querySelector('div').shadowRoot.appendChild(span);
	}`)
	must.NoError(err)

	err = page.DispatchEvent(ctx, "span", "click")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => window['clicked']")
	must.NoError(err)
	is.Equal(true, val)
}

func TestDispatchEventDragDrop(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="source" draggable="true"
		style="width:50px;height:50px;position:absolute;left:0;top:0">source</div>
	<div id="target" style="width:100px;height:100px;position:absolute;left:200px;top:0">target</div>
	<script>
		document.querySelector('#source').addEventListener('dragstart', e => {
			e.dataTransfer.setData('text', 'dragged');
		});
		document.querySelector('#target').addEventListener('dragover', e => e.preventDefault());
		document.querySelector('#target').addEventListener('drop', e => {
			e.preventDefault();
			e.target.appendChild(document.querySelector('#source'));
		});
	</script>`)
	must.NoError(err)

	err = page.Locator("#source").DragTo(ctx, page.Locator("#target"))
	must.NoError(err)

	sourceInTarget, err := page.Evaluate(ctx, `() => {
		const source = document.querySelector('#source');
		const target = document.querySelector('#target');
		return source && source.parentElement === target;
	}`)
	must.NoError(err)
	is.Equal(true, sourceInTarget)
}

// TestLocatorDispatchClickEx2 verifies page DispatchEvent fires click via locator-style selector.
// Ref: TestLocatorDispatchEvent.java#shouldFireClick
func TestLocatorDispatchClickEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__dispatched=true">Click</button>
	`))

	must.NoError(page.DispatchEvent(ctx, "#btn", "click"))

	dispatched, err := page.Evaluate(ctx, `() => window.__dispatched`)
	must.NoError(err)
	is.Equal(true, dispatched)
}

// TestLocatorDispatchMouseoverEx2 verifies DispatchEvent fires mouseover.
// Ref: TestLocatorDispatchEvent.java#shouldFireMouseover
func TestLocatorDispatchMouseoverEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" onmouseover="window.__mouseOver=true">hover</div>
	`))

	must.NoError(page.DispatchEvent(ctx, "#el", "mouseover"))

	mo, err := page.Evaluate(ctx, `() => window.__mouseOver`)
	must.NoError(err)
	is.Equal(true, mo)
}

// TestLocatorDispatchFocusEx2 verifies DispatchEvent fires focus on input.
// Ref: TestLocatorDispatchEvent.java#shouldFireFocus
func TestLocatorDispatchFocusEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onfocus="window.__focused=true">
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "focus"))

	focused, err := page.Evaluate(ctx, `() => window.__focused`)
	must.NoError(err)
	is.Equal(true, focused)
}

// TestLocatorDispatchCustomEventEx2 verifies DispatchEvent fires custom events.
// Ref: TestLocatorDispatchEvent.java#shouldFireCustomEvent
func TestLocatorDispatchCustomEventEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el"></div>
		<script>
			document.getElementById('el').addEventListener('custom', () => {
				window.__customEvent = true;
			});
		</script>
	`))

	must.NoError(page.DispatchEvent(ctx, "#el", "custom"))

	fired, err := page.Evaluate(ctx, `() => window.__customEvent`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestPageDispatchEventFocus verifies DispatchEvent can fire focus events.
// Ref: TestPageDispatchEvent.java#shouldFireFocusEvent
func TestPageDispatchEventFocus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onfocus="this.dataset.focused='yes'">
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "focus"))

	result, err := page.Evaluate(ctx, `() => document.getElementById('inp').dataset.focused`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestPageDispatchEventMouseover verifies DispatchEvent can fire mouseover events.
// Ref: TestPageDispatchEvent.java#shouldFireMouseoverEvent
func TestPageDispatchEventMouseover(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" onmouseover="this.dataset.hovered='yes'">hover me</div>
	`))

	must.NoError(page.DispatchEvent(ctx, "#el", "mouseover"))

	result, err := page.Evaluate(ctx, `() => document.getElementById('el').dataset.hovered`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestPageDispatchEventKeydown verifies DispatchEvent can fire keydown events.
// Ref: TestPageDispatchEvent.java#shouldFireKeydownEvent
func TestPageDispatchEventKeydown(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onkeydown="document.getElementById('result').textContent = event.key">
		<div id="result"></div>
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "keydown", map[string]any{"key": "Enter"}))

	text, err := page.Locator("#result").InnerText(ctx)
	must.NoError(err)
	is.Equal("Enter", text)
}

// TestPageDispatchEventViaSelector verifies DispatchEvent works via page-level selector.
// Ref: TestPageDispatchEvent.java#shouldWorkWithSelector
func TestPageDispatchEventViaSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="this.dataset.pressed='yes'">click</button>
	`))

	must.NoError(page.DispatchEvent(ctx, "#btn", "click"))

	result, err := page.Evaluate(ctx, `() => document.getElementById('btn').dataset.pressed`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestPageDispatchEventClickExtra2 verifies DispatchEvent fires click event.
// Ref: TestPageDispatchEvent.java#shouldFireClick
func TestPageDispatchEventClickExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked=true">Click</button>
	`))

	must.NoError(page.DispatchEvent(ctx, "#btn", "click"))

	clicked, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Equal(true, clicked)
}

// TestPageDispatchEventCustomExtra2 verifies DispatchEvent can fire custom events.
// Ref: TestPageDispatchEvent.java#shouldFireCustomEvent
func TestPageDispatchEventCustomExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el"></div>
		<script>
			document.getElementById('el').addEventListener('my-event', () => {
				window.__customFired = true;
			});
		</script>
	`))

	must.NoError(page.DispatchEvent(ctx, "#el", "my-event"))

	fired, err := page.Evaluate(ctx, `() => window.__customFired`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestPageDispatchEventChangeExtra2 verifies DispatchEvent fires change event on input.
// Ref: TestPageDispatchEvent.java#shouldFireChangeEvent
func TestPageDispatchEventChangeExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onchange="window.__changed=true">
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "change"))

	changed, err := page.Evaluate(ctx, `() => window.__changed`)
	must.NoError(err)
	is.Equal(true, changed)
}

// TestPageDispatchEventSubmitExtra2 verifies DispatchEvent fires submit on form.
// Ref: TestPageDispatchEvent.java#shouldFireSubmitEvent
func TestPageDispatchEventSubmitExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form id="frm" onsubmit="window.__submitted=true; return false;">
			<input type="submit" value="Submit">
		</form>
	`))

	must.NoError(page.DispatchEvent(ctx, "#frm", "submit"))

	submitted, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, submitted)
}

// TestPageDispatchEventDblClickExtra2 verifies DispatchEvent fires dblclick event.
// Ref: TestPageDispatchEvent.java#shouldFireDblClick
func TestPageDispatchEventDblClickExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" ondblclick="window.__dblClicked=true">double click</div>
	`))

	must.NoError(page.DispatchEvent(ctx, "#el", "dblclick"))

	dbl, err := page.Evaluate(ctx, `() => window.__dblClicked`)
	must.NoError(err)
	is.Equal(true, dbl)
}

// TestDispatchEventFocusEx3 verifies dispatching focus event.
// Ref: TestPageDispatchEvent.java#shouldDispatchFocus
func TestDispatchEventFocusEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text">
		<script>
			var focused = false;
			document.getElementById('inp').addEventListener('focus', function() { focused = true; });
		</script>
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "focus"))

	result, err := page.Evaluate(ctx, `() => focused`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestDispatchEventBlurEx3 verifies dispatching blur event.
// Ref: TestPageDispatchEvent.java#shouldDispatchBlur
func TestDispatchEventBlurEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text">
		<script>
			var blurred = false;
			document.getElementById('inp').addEventListener('blur', function() { blurred = true; });
		</script>
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "focus"))
	must.NoError(page.DispatchEvent(ctx, "#inp", "blur"))

	result, err := page.Evaluate(ctx, `() => blurred`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestDispatchEventMouseOverEx3 verifies dispatching mouseover event.
// Ref: TestPageDispatchEvent.java#shouldDispatchMouseover
func TestDispatchEventMouseOverEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px;"></div>
		<script>
			var overFired = false;
			document.getElementById('el').addEventListener('mouseover', function() { overFired = true; });
		</script>
	`))

	must.NoError(page.DispatchEvent(ctx, "#el", "mouseover"))

	result, err := page.Evaluate(ctx, `() => overFired`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestDispatchEventKeyPressEx3 verifies dispatching keypress event.
// Ref: TestPageDispatchEvent.java#shouldDispatchKeypress
func TestDispatchEventKeyPressEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text">
		<script>
			var keyFired = false;
			document.getElementById('inp').addEventListener('keydown', function() { keyFired = true; });
		</script>
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "keydown"))

	result, err := page.Evaluate(ctx, `() => keyFired`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestDispatchFocusEventEx4 verifies DispatchEvent fires focus event.
// Ref: TestPageDispatchEvent.java#shouldDispatchFocusEvent
func TestDispatchFocusEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onfocus="document.getElementById('out').textContent='focused'">
		<span id="out"></span>
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "focus"))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("focused", out)
}

// TestDispatchBlurEventEx4 verifies DispatchEvent fires blur event.
// Ref: TestPageDispatchEvent.java#shouldDispatchBlurEvent
func TestDispatchBlurEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onblur="document.getElementById('out').textContent='blurred'">
		<span id="out"></span>
	`))

	must.NoError(page.DispatchEvent(ctx, "#inp", "focus"))
	must.NoError(page.DispatchEvent(ctx, "#inp", "blur"))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("blurred", out)
}

// TestDispatchKeydownEventEx4 verifies DispatchEvent fires keydown event.
// Ref: TestPageDispatchEvent.java#shouldDispatchKeydownEvent
func TestDispatchKeydownEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" tabindex="0" onkeydown="document.getElementById('out').textContent='key:'+event.key"></div>
		<span id="out"></span>
	`))

	must.NoError(page.DispatchEvent(ctx, "#d", "keydown", map[string]interface{}{"key": "Enter"}))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("key:Enter", out)
}

// TestDispatchChangeEventEx4 verifies DispatchEvent fires change event on select.
// Ref: TestPageDispatchEvent.java#shouldDispatchChangeOnSelect
func TestDispatchChangeEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" onchange="document.getElementById('out').textContent='changed'">
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
		<span id="out"></span>
	`))

	must.NoError(page.DispatchEvent(ctx, "#sel", "change"))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("changed", out)
}
