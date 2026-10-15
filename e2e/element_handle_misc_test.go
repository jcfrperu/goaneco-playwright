//go:build e2e

// E2E tests for: TestElementHandleMisc.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"bytes"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/vis", "text/html", `
		<div id="visible">visible</div>
		<div id="hidden" style="display:none">hidden</div>
	`)
	err := page.Goto(ctx, srv.Prefix()+"/vis")
	must.NoError(err, "Goto failed")

	visible, err := page.QuerySelector(ctx, "#visible")
	must.NoError(err)
	must.NotNil(visible, "QuerySelector #visible failed")
	hidden, err := page.QuerySelector(ctx, "#hidden")
	must.NoError(err)
	must.NotNil(hidden, "QuerySelector #hidden failed")

	isVis, err := visible.IsVisible(ctx)
	must.NoError(err, "IsVisible failed")
	is.True(isVis, "expected #visible to be visible")

	isHid, err := hidden.IsHidden(ctx)
	must.NoError(err, "IsHidden failed")
	is.True(isHid, "expected #hidden to be hidden")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleGetAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/attr", "text/html", `<a href="https://example.com" data-id="42">link</a>`)
	err := page.Goto(ctx, srv.Prefix()+"/attr")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "a")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	href, err := el.GetAttribute(ctx, "href")
	must.NoError(err, "GetAttribute href failed")
	must.NotNil(href)
	is.Equal("https://example.com", *href)

	dataID, err := el.GetAttribute(ctx, "data-id")
	must.NoError(err, "GetAttribute data-id failed")
	must.NotNil(dataID)
	is.Equal("42", *dataID)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleInputValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/input", "text/html", `<input id="field" value="initial" />`)
	err := page.Goto(ctx, srv.Prefix()+"/input")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "#field")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	val, err := el.InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("initial", val)

	err = el.Fill(ctx, "updated")
	must.NoError(err, "Fill failed")

	val2, err := el.InputValue(ctx)
	must.NoError(err, "InputValue after Fill failed")
	is.Equal("updated", val2)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestJSHandleEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/js", "text/html", `<div id="target">42</div>`)
	err := page.Goto(ctx, srv.Prefix()+"/js")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	// Use Evaluate on the ElementHandle to get its text via JS
	result, err := el.Evaluate(ctx, "el => el.textContent")
	must.NoError(err, "Evaluate failed")
	is.Equal("42", result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleTap(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="btn" ontouchstart="this.dataset.tapped='yes'">tap me</button>
		<script>
		document.getElementById('btn').addEventListener('click', function() {
			this.dataset.clicked = 'yes';
		});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	err = el.Tap(ctx)
	must.NoError(err, "ElementHandle.Tap failed")

	clicked, err := page.Evaluate(ctx, `document.getElementById('btn').dataset.clicked`)
	must.NoError(err, "Evaluate failed")
	if clicked != "yes" {
		t.Errorf("Tap did not trigger click: dataset.clicked = %v", clicked)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="inp" type="text">`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	err = el.Type(ctx, "hello world")
	must.NoError(err, "ElementHandle.Type failed")

	val, err := el.InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("hello world", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleTypeFiresKeyEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="inp" type="text">
		<div id="log"></div>
		<script>
		document.getElementById('inp').addEventListener('keydown', function(e) {
			document.getElementById('log').textContent += e.key;
		});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	err = el.Type(ctx, "ab")
	must.NoError(err, "ElementHandle.Type failed")

	log, err := page.Locator("#log").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("ab", log)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleScreenshot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="box" style="width:80px;height:80px;background:green"></div>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#box")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	data, err := el.Screenshot(ctx)
	must.NoError(err, "ElementHandle.Screenshot failed")
	if len(data) == 0 {
		t.Fatal("Screenshot returned empty bytes")
	}
	if !bytes.HasPrefix(data, []byte("\x89PNG")) {
		t.Errorf("expected PNG header, got: %x", data[:min(8, len(data))])
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleScreenshotJPEG(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="box" style="width:80px;height:80px;background:navy"></div>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#box")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	q := 75
	data, err := el.Screenshot(ctx, &playwright.ScreenshotOptions{Type: "jpeg", Quality: &q})
	must.NoError(err, "ElementHandle.Screenshot JPEG failed")
	if len(data) < 2 {
		t.Fatal("Screenshot JPEG returned too few bytes")
	}
	if data[0] != 0xFF || data[1] != 0xD8 {
		t.Errorf("expected JPEG header FF D8, got: %02x %02x", data[0], data[1])
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFillInputExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Fill(ctx, "filled text"))

	val, err := el.InputValue(ctx)
	must.NoError(err)
	is.Equal("filled text", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleEvaluateReturnsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" class="test-class">hello</div>`))

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el)

	result, err := el.Evaluate(ctx, `el => el.className`)
	must.NoError(err)
	is.Equal("test-class", result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleGetAttributeExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a href="https://example.com" id="link">link</a>`))

	handle, err := page.QuerySelector(ctx, "#link")
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.GetAttribute(ctx, "href")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("https://example.com", *val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleGetAttributeMissing(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">text</div>`))

	handle, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.GetAttribute(ctx, "nonexistent")
	must.NoError(err)
	is.Nil(val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsEnabledExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="enabled">go</button>
		<button id="disabled" disabled>stop</button>
	`))

	enabled, err := page.QuerySelector(ctx, "#enabled")
	must.NoError(err)
	isEnabled, err := enabled.IsEnabled(ctx)
	must.NoError(err)
	is.True(isEnabled)

	disabled, err := page.QuerySelector(ctx, "#disabled")
	must.NoError(err)
	isDisabled, err := disabled.IsEnabled(ctx)
	must.NoError(err)
	is.False(isDisabled)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsCheckedExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="checked" checked>
		<input type="checkbox" id="unchecked">
	`))

	checkedEl, err := page.QuerySelector(ctx, "#checked")
	must.NoError(err)
	isChecked, err := checkedEl.IsChecked(ctx)
	must.NoError(err)
	is.True(isChecked)

	uncheckedEl, err := page.QuerySelector(ctx, "#unchecked")
	must.NoError(err)
	isNotChecked, err := uncheckedEl.IsChecked(ctx)
	must.NoError(err)
	is.False(isNotChecked)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleInnerTextExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="para">Hello <b>World</b></p>`))

	handle, err := page.QuerySelector(ctx, "#para")
	must.NoError(err)
	must.NotNil(handle)

	text, err := handle.InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello World", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsVisibleForVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="vis">visible</div>`))

	handle, err := page.QuerySelector(ctx, "#vis")
	must.NoError(err)
	must.NotNil(handle)

	visible, err := handle.IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleDispatchEventExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked=true">go</button>
	`))

	handle, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.DispatchEvent(ctx, "click"))

	result, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleSetCheckedExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	handle, err := page.QuerySelector(ctx, "#cb")
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.SetChecked(ctx, true))

	checked, err := handle.IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	must.NoError(handle.SetChecked(ctx, false))

	unchecked, err := handle.IsChecked(ctx)
	must.NoError(err)
	is.False(unchecked)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleScrollIntoViewIfNeededExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px"></div>
		<div id="target" style="height:50px;background:blue">target</div>
	`))

	handle, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.ScrollIntoViewIfNeeded(ctx))

	visible, err := handle.IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleHoverExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" onmouseover="window.__hovered=true" style="width:100px;height:100px">hover me</div>
	`))

	handle, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.Hover(ctx))

	result, err := page.Evaluate(ctx, `() => window.__hovered`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleTextContentExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Hello <span style="display:none">hidden</span> world</div>`))

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Hello")
	is.Contains(text, "hidden")
	is.Contains(text, "world")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleInputValueExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="initial">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	val, err := el.InputValue(ctx)
	must.NoError(err)
	is.Equal("initial", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFillSetsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Fill(ctx, "filled value"))

	val, err := el.InputValue(ctx)
	must.NoError(err)
	is.Equal("filled value", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFocusMakesActive(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("inp", focused)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleEvaluateReadsProperty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Click me</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	disabled, err := el.Evaluate(ctx, "el => el.disabled")
	must.NoError(err)
	is.Equal(true, disabled)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleHoverEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onmouseover="window.__hovered=true">hover</div>
	`))

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Hover(ctx))

	hovered, err := page.Evaluate(ctx, `() => window.__hovered`)
	must.NoError(err)
	is.Equal(true, hovered)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsVisibleEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Visible</p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)
	must.NotNil(el)

	visible, err := el.IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsEnabledEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	enabled, err := el.IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleGetAttributeEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a href="https://example.com" id="lnk">Link</a>`))

	el, err := page.QuerySelector(ctx, "#lnk")
	must.NoError(err)
	must.NotNil(el)

	href, err := el.GetAttribute(ctx, "href")
	must.NoError(err)
	must.NotNil(href)
	is.Equal("https://example.com", *href)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleTextContentEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello <b>World</b></p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.TextContent(ctx)
	must.NoError(err)
	is.Equal("Hello World", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsVisibleEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="visible">Visible</div>
		<div id="hidden" style="display:none">Hidden</div>
	`))

	v, err := page.QuerySelector(ctx, "#visible")
	must.NoError(err)
	visibleResult, err := v.IsVisible(ctx)
	must.NoError(err)
	is.True(visibleResult)

	h, err := page.QuerySelector(ctx, "#hidden")
	must.NoError(err)
	hiddenResult, err := h.IsVisible(ctx)
	must.NoError(err)
	is.False(hiddenResult)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleInnerTextEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello <span style="display:none">Hidden</span>World</p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	must.NotContains(text, "Hidden")
	is.Contains(text, "Hello")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsCheckedEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	el, err := page.QuerySelector(ctx, "#chk")
	must.NoError(err)

	checked, err := el.IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsEnabledEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" disabled>`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)

	enabled, err := el.IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFillEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NoError(el.Fill(ctx, "filled value"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("filled value", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFillEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Fill(ctx, "hello world"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hello world", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleTypeEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Type(ctx, "typed"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("typed", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFocusEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("inp", focused)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsVisibleEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible div</div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	visible, err := el.IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsHiddenEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none">Hidden</div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	hidden, err := el.IsHidden(ctx)
	must.NoError(err)
	is.True(hidden)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsEnabledEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	enabled, err := el.IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsCheckedEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	el, err := page.QuerySelector(ctx, "#chk")
	must.NoError(err)
	must.NotNil(el)

	checked, err := el.IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleGetAttributeEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="/test" target="_blank">Link</a>`))

	el, err := page.QuerySelector(ctx, "#a")
	must.NoError(err)
	must.NotNil(el)

	target, err := el.GetAttribute(ctx, "target")
	must.NoError(err)
	must.NotNil(target)
	is.Equal("_blank", *target)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleHoverEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" onmouseover="document.getElementById('out').textContent='hovered'">Hover me</div>
		<span id="out"></span>
	`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Hover(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("hovered", out)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleScrollIntoViewEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;"></div>
		<button id="far-btn">Far Button</button>
	`))

	el, err := page.QuerySelector(ctx, "#far-btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.ScrollIntoViewIfNeeded(ctx))

	visible, err := el.IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleTextContentEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Paragraph text content</p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.TextContent(ctx)
	must.NoError(err)
	is.Equal("Paragraph text content", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleInnerTextEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h">Page Title</h1>`))

	el, err := page.QuerySelector(ctx, "#h")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("Page Title", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleEvaluateEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" data-val="42">Content</div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	result, err := el.Evaluate(ctx, `el => el.dataset.val`)
	must.NoError(err)
	is.Equal("42", result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsDisabledEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Disabled</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	disabled, err := el.IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">hello</div>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	text, err := el.InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("hello", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleTextContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">hello <b>world</b></div>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	tc, err := el.TextContent(ctx)
	must.NoError(err, "TextContent failed")
	is.Equal("hello world", tc)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="ed"/><input id="ro" readonly/>`)
	must.NoError(err, "SetContent failed")

	elEd, _ := page.QuerySelector(ctx, "#ed")
	elRo, _ := page.QuerySelector(ctx, "#ro")

	editable, err := elEd.IsEditable(ctx)
	if err != nil || !editable {
		t.Errorf("#ed should be editable, got %v err=%v", editable, err)
	}

	readonly, err := elRo.IsEditable(ctx)
	if err != nil || readonly {
		t.Errorf("#ro should NOT be editable, got %v err=%v", readonly, err)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsCheckedAndUnchecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox" id="ch" checked/><input type="checkbox" id="un"/>`)
	must.NoError(err, "SetContent failed")

	elCh, _ := page.QuerySelector(ctx, "#ch")
	elUn, _ := page.QuerySelector(ctx, "#un")

	checked, err := elCh.IsChecked(ctx)
	if err != nil || !checked {
		t.Errorf("#ch should be checked, got %v err=%v", checked, err)
	}

	unchecked, err := elUn.IsChecked(ctx)
	if err != nil || unchecked {
		t.Errorf("#un should NOT be checked, got %v err=%v", unchecked, err)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleHover(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn">Hover me</button>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	err = el.Hover(ctx)
	must.NoError(err, "Hover failed")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFillInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="in"/>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#in")
	err = el.Fill(ctx, "hello")
	must.NoError(err, "Fill failed")

	val, err := el.InputValue(ctx)
	if err != nil || val != "hello" {
		t.Errorf("InputValue = %q, want 'hello' (err=%v)", val, err)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleCheckAndUncheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox" id="cb"/>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#cb")

	err = el.Check(ctx)
	must.NoError(err, "Check failed")
	checked, _ := el.IsChecked(ctx)
	is.True(checked, "expected checkbox to be checked")

	err = el.Uncheck(ctx)
	must.NoError(err, "Uncheck failed")
	checked, _ = el.IsChecked(ctx)
	is.False(checked, "expected checkbox to be unchecked")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleSetChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox" id="cb"/>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#cb")

	err = el.SetChecked(ctx, true)
	must.NoError(err, "SetChecked(true) failed")
	if checked, _ := el.IsChecked(ctx); !checked {
		t.Error("expected checkbox to be checked after SetChecked(true)")
	}

	err = el.SetChecked(ctx, false)
	must.NoError(err, "SetChecked(false) failed")
	if checked, _ := el.IsChecked(ctx); checked {
		t.Error("expected checkbox to be unchecked after SetChecked(false)")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFocusButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn">Click</button>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#btn")
	err = el.Focus(ctx)
	must.NoError(err, "Focus failed")

	focused, _ := page.Evaluate(ctx, "() => document.activeElement.id")
	if focused != "btn" {
		t.Errorf("expected activeElement to be 'btn', got %v", focused)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleDisposeIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">text</div>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#el")

	err = el.Dispose(ctx)
	must.NoError(err, "first Dispose failed")
	// Second dispose should not panic.
	_ = el.Dispose(ctx)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleAsJSHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" data-val="42">hello</div>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	handle := el.AsJSHandle()
	must.NotNil(handle, "AsJSHandle() returned nil")

	// Verify the JSHandle references the element by evaluating a property on it.
	val, err := handle.Evaluate(ctx, "e => e.dataset.val")
	must.NoError(err, "JSHandle.Evaluate failed")
	is.Equal("42", val)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleEvaluateHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul><li>first</li><li>second</li></ul>`)
	must.NoError(err, "SetContent failed")

	ul, err := page.QuerySelector(ctx, "ul")
	must.NoError(err)
	must.NotNil(ul, "QuerySelector failed")

	// EvaluateHandle returns a JSHandle for the first <li>.
	handle, err := ul.EvaluateHandle(ctx, "el => el.querySelector('li')")
	must.NoError(err, "EvaluateHandle failed")
	must.NotNil(handle, "EvaluateHandle returned nil handle")

	// The returned handle should reference the first <li>.
	text, err := handle.Evaluate(ctx, "el => el.textContent")
	must.NoError(err, "Evaluate on returned handle failed")
	is.Equal("first", text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleScrollIntoViewIfNeeded(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div style="height:2000px;background:#eee">spacer</div>
		<div id="target" style="height:100px;background:red">target</div>
	`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	err = el.ScrollIntoViewIfNeeded(ctx)
	must.NoError(err, "ElementHandle.ScrollIntoViewIfNeeded failed")

	visible, err := el.IsVisible(ctx)
	must.NoError(err, "IsVisible failed")
	is.True(visible, "element should be visible after ScrollIntoViewIfNeeded")
}

// --- From element_handle_extra_b_test.go ---

// TestElementHandleGetAttributeAbsentEx verifies GetAttribute returns nil for missing attribute.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleGetAttributeAbsentEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	attr, err := el.GetAttribute(ctx, "data-missing")
	must.NoError(err)
	is.Nil(attr)
}

// TestElementHandleCheckEx verifies ElementHandle Check.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleCheckEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	el, err := page.QuerySelector(ctx, "#cb")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Check(ctx))

	checked, err := el.IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestElementHandleUncheckEx verifies ElementHandle Uncheck.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleUncheckEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox" checked>`))

	el, err := page.QuerySelector(ctx, "#cb")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Uncheck(ctx))

	checked, err := el.IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestElementHandleFillAndGetValueEx verifies ElementHandle Fill then InputValue.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleFillAndGetValueEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Fill(ctx, "test value"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("test value", val)
}

// --- From element_handle_misc_test.go ---

// TestElementHandleGetAttributeReturnsNilWhenMissing verifies GetAttribute returns nil for missing attr.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleGetAttributeReturnsNilWhenMissing(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">content</div>`))

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el)

	attr, err := el.GetAttribute(ctx, "nonexistent-attr")
	must.NoError(err)
	is.Nil(attr, "should return nil for non-existing attribute")
}

// TestElementHandleIsHidden verifies IsVisible returns false for hidden elements.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="display:none">invisible</div>`))

	el, err := page.QuerySelector(ctx, "#hidden")
	must.NoError(err)
	must.NotNil(el)

	visible, err := el.IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestElementHandleIsEnabledOnInput verifies IsEnabled returns true for enabled input.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsEnabledOnInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="active" type="text">`))

	el, err := page.QuerySelector(ctx, "#active")
	must.NoError(err)
	must.NotNil(el)

	enabled, err := el.IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestElementHandleIsDisabledOnInput verifies IsEnabled returns false for disabled input.
// Extra: no equivalent method in playwright-java; related class: TestElementHandleMisc.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleIsDisabledOnInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="disabled" type="text" disabled>`))

	el, err := page.QuerySelector(ctx, "#disabled")
	must.NoError(err)
	must.NotNil(el)

	enabled, err := el.IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// Ref: TestElementHandleMisc.java#shouldFocusAButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldFocusAButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn">Click me</button>`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#btn").Focus(ctx)
	must.NoError(err, "Focus failed")

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err, "Evaluate activeElement failed")
	if focused != "btn" {
		t.Errorf("expected activeElement to be 'btn', got %v", focused)
	}
}

// Ref: TestElementHandleMisc.java#shouldHover  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldHover(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, a3ScrollableHTML))
	button, err := page.QuerySelector(ctx, "#button-6")
	must.NoError(err)
	must.NotNil(button)
	must.NoError(button.Hover(ctx))
	id, err := page.Evaluate(ctx, "() => document.querySelector('button:hover').id")
	must.NoError(err)
	is.Equal("button-6", id)
}

// Ref: TestElementHandleMisc.java#shouldHoverWhenNodeIsRemoved  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldHoverWhenNodeIsRemoved(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, a3ScrollableHTML))
	_, err := page.Evaluate(ctx, "() => delete window['Node']")
	must.NoError(err)
	button, err := page.QuerySelector(ctx, "#button-6")
	must.NoError(err)
	must.NotNil(button)
	must.NoError(button.Hover(ctx))
	id, err := page.Evaluate(ctx, "() => document.querySelector('button:hover').id")
	must.NoError(err)
	is.Equal("button-6", id)
}

// Ref: TestElementHandleMisc.java#shouldFillInput  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldFillInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, a3InputResultHTML))
	handle, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(handle)
	must.NoError(handle.Fill(ctx, "some value"))
	res, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("some value", res)
}

// Ref: TestElementHandleMisc.java#shouldFillInputWhenNodeIsRemoved  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldFillInputWhenNodeIsRemoved(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, a3InputResultHTML))
	_, err := page.Evaluate(ctx, "() => delete window['Node']")
	must.NoError(err)
	handle, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(handle)
	must.NoError(handle.Fill(ctx, "some value"))
	res, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("some value", res)
}

// Ref: TestElementHandleMisc.java#shouldCheckTheBox  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldCheckTheBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input id='checkbox' type='checkbox'></input>"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Check(ctx))
	res, err := page.Evaluate(ctx, "() => checkbox.checked")
	must.NoError(err)
	is.Equal(true, res)
}

// Ref: TestElementHandleMisc.java#shouldUncheckTheBox  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldUncheckTheBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input id='checkbox' type='checkbox' checked></input>"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.Uncheck(ctx))
	res, err := page.Evaluate(ctx, "() => checkbox.checked")
	must.NoError(err)
	is.Equal(false, res)
}

// Ref: TestElementHandleMisc.java#shouldSelectSingleOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldSelectSingleOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, a3SelectHTML))
	sel, err := page.QuerySelector(ctx, "select")
	must.NoError(err)
	must.NotNil(sel)
	_, err = sel.SelectOption(ctx, "blue")
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"blue"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"blue"}, onChange)
}

// Ref: TestElementHandleMisc.java#shouldCheckTheBoxUsingSetChecked  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldCheckTheBoxUsingSetChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input id='checkbox' type='checkbox'></input>"))
	input, err := page.QuerySelector(ctx, "input")
	must.NoError(err)
	must.NotNil(input)
	must.NoError(input.SetChecked(ctx, true))
	res, err := page.Evaluate(ctx, "() => checkbox.checked")
	must.NoError(err)
	is.Equal(true, res)
	must.NoError(input.SetChecked(ctx, false))
	res, err = page.Evaluate(ctx, "() => checkbox.checked")
	must.NoError(err)
	is.Equal(false, res)
}

// Ref: TestElementHandleMisc.java#shouldFallBackToSelectingByLabel  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldFallBackToSelectingByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, a3SelectHTML))
	sel, err := page.QuerySelector(ctx, "select")
	must.NoError(err)
	must.NotNil(sel)
	_, err = sel.SelectOption(ctx, "Blue")
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"blue"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"blue"}, onChange)
}

// Ref: TestElementHandleMisc.java#shouldAllowDisposingTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleMisc.java
func TestElementHandleMiscShouldAllowDisposingTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section>39</section>"))
	element, err := page.QuerySelector(ctx, "section")
	must.NoError(err)
	must.NotNil(element)
	must.NoError(element.Dispose(ctx))
	must.NoError(element.Dispose(ctx))
}
