//go:build e2e

package e2e

import (
	"bytes"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElementHandleQuerySelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/el", "text/html", `<div id="target">hello world</div>`)
	err := page.Goto(ctx, srv.Prefix()+"/el")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "#target")
	must.NoError(err, "QuerySelector failed")
	must.NotNil(el, "QuerySelector returned nil for existing element")

	text, err := el.InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("hello world", text)
}

func TestElementHandleQuerySelectorReturnsNilWhenNotFound(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/el", "text/html", `<div>nothing</div>`)
	err := page.Goto(ctx, srv.Prefix()+"/el")
	must.NoError(err, "Goto failed")

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err, "QuerySelector failed")
	is.Nil(el)
}

func TestElementHandleQuerySelectorAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/list", "text/html", `<ul><li>one</li><li>two</li><li>three</li></ul>`)
	err := page.Goto(ctx, srv.Prefix()+"/list")
	must.NoError(err, "Goto failed")

	elements, err := page.QuerySelectorAll(ctx, "li")
	must.NoError(err, "QuerySelectorAll failed")
	is.Len(elements, 3)

	wantTexts := []string{"one", "two", "three"}
	for i, el := range elements {
		text, err := el.InnerText(ctx)
		must.NoErrorf(err, "InnerText[%d] failed", i)
		if text != wantTexts[i] {
			t.Errorf("element[%d] text = %q, want %q", i, text, wantTexts[i])
		}
	}
}

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

func TestElementHandleClickFiresClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="this.dataset.clicked='yes'">click</button>
	`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Click(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('btn').dataset.clicked`)
	must.NoError(err)
	is.Equal("yes", result)
}

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

func TestElementHandleTextContentReturnsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello <span>World</span></p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)
	must.NotNil(el)

	tc, err := el.TextContent(ctx)
	must.NoError(err)
	is.Equal("Hello World", tc)
}

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

func TestElementHandleIsVisibleForVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="visible">I'm visible</div>`))

	el, err := page.QuerySelector(ctx, "#visible")
	must.NoError(err)
	must.NotNil(el)

	visible, err := el.IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

func TestElementHandleQuerySelectorFindsChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="parent"><span id="child">child</span></div>`))

	parent, err := page.QuerySelector(ctx, "#parent")
	must.NoError(err)
	must.NotNil(parent)

	child, err := parent.QuerySelector(ctx, "#child")
	must.NoError(err)
	must.NotNil(child)

	text, err := child.InnerText(ctx)
	must.NoError(err)
	is.Equal("child", text)
}

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

func TestElementHandleInnerHTMLExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"><span>child</span></div>`))

	handle, err := page.QuerySelector(ctx, "#container")
	must.NoError(err)
	must.NotNil(handle)

	html, err := handle.InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<span>child</span>")
}

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

func TestElementHandleIsHiddenForHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="display:none">hidden</div>`))

	handle, err := page.QuerySelector(ctx, "#hidden")
	must.NoError(err)
	must.NotNil(handle)

	hidden, err := handle.IsHidden(ctx)
	must.NoError(err)
	is.True(hidden)
}

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

func TestElementHandleIsDisabledForDisabledButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>go</button>`))

	handle, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(handle)

	disabled, err := handle.IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)
}

func TestElementHandleIsEditableForInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="editable" type="text">
		<input id="readonly" type="text" readonly>
	`))

	editable, err := page.QuerySelector(ctx, "#editable")
	must.NoError(err)
	isEditable, err := editable.IsEditable(ctx)
	must.NoError(err)
	is.True(isEditable)

	readonly, err := page.QuerySelector(ctx, "#readonly")
	must.NoError(err)
	isReadOnly, err := readonly.IsEditable(ctx)
	must.NoError(err)
	is.False(isReadOnly)
}

func TestElementHandleQuerySelectorAllExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			<span>a</span>
			<span>b</span>
			<span>c</span>
		</div>
	`))

	parent, err := page.QuerySelector(ctx, "#parent")
	must.NoError(err)
	must.NotNil(parent)

	children, err := parent.QuerySelectorAll(ctx, "span")
	must.NoError(err)
	is.Len(children, 3)
}

func TestElementHandleSelectOptionExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
			<option value="c">Gamma</option>
		</select>
	`))

	handle, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(handle)

	selected, err := handle.SelectOption(ctx, "b")
	must.NoError(err)
	is.Contains(selected, "b")
}

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

func TestElementHandleTypeInInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.Type(ctx, "hello"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Equal("hello", val)
}

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

func TestElementHandlePressEnterSubmitsForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form onsubmit="window.__submitted=true;return false;">
			<input id="inp" type="text">
			<button type="submit">Go</button>
		</form>
	`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Press(ctx, "Enter"))

	result, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, result)
}

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

func TestElementHandleClickEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked=true">Click</button>
	`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Click(ctx))

	clicked, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Equal(true, clicked)
}

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

func TestElementHandleSelectOptionEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
	`))

	el, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(el)

	_, err = el.SelectOption(ctx, "b")
	must.NoError(err)

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("b", val)
}

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

func TestElementHandleInnerHTMLEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><span>child</span></div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	html, err := el.InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<span>child</span>", html)
}

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

func TestElementHandleClickEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked9=true">Click</button>
	`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NoError(el.Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__clicked9`)
	must.NoError(err)
	is.Equal(true, result)
}

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

func TestElementHandleClickEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" onclick="this.textContent='clicked'">Click</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Click(ctx))

	text, err := page.Locator("#btn").TextContent(ctx)
	must.NoError(err)
	is.Equal("clicked", text)
}

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

func TestElementHandleSelectOptionEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
			<option value="c">Gamma</option>
		</select>
	`))

	el, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(el)

	_, err = el.SelectOption(ctx, "b")
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => document.getElementById('sel').value`)
	must.NoError(err)
	is.Equal("b", val)
}

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

func TestElementHandleDispatchClickEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" onclick="this.textContent='dispatched'">Click</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.DispatchEvent(ctx, "click"))

	text, err := page.Locator("#btn").TextContent(ctx)
	must.NoError(err)
	is.Equal("dispatched", text)
}

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

func TestElementHandleQuerySelectorEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="container">
			<span id="child">Child element</span>
		</div>
	`))

	container, err := page.QuerySelector(ctx, "#container")
	must.NoError(err)
	must.NotNil(container)

	child, err := container.QuerySelector(ctx, "#child")
	must.NoError(err)
	must.NotNil(child)

	text, err := child.TextContent(ctx)
	must.NoError(err)
	is.Equal("Child element", text)
}

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

func TestElementHandleInnerHTMLEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><strong>Bold</strong></div>`))

	el, err := page.QuerySelector(ctx, "#d")
	must.NoError(err)
	must.NotNil(el)

	html, err := el.InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<strong>")
}

func TestElementHandlePressEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Type(ctx, "test"))
	must.NoError(el.Press(ctx, "Backspace"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("tes", val)
}

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

func TestElementHandleSelectOptionEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
			<option value="c">Gamma</option>
		</select>
	`))

	el, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(el)

	_, err = el.SelectOption(ctx, "b")
	must.NoError(err)

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("b", val)
}

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

func TestElementHandleQuerySelectorAllEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul id="list">
			<li>A</li>
			<li>B</li>
			<li>C</li>
		</ul>
	`))

	list, err := page.QuerySelector(ctx, "#list")
	must.NoError(err)
	must.NotNil(list)

	items, err := list.QuerySelectorAll(ctx, "li")
	must.NoError(err)
	is.Len(items, 3)
}

// --- From element_handle_bounding_box_extra_test.go ---

// TestElementHandleBoundingBoxCorrectDimensions verifies BoundingBox returns correct dimensions.
// Ref: TestElementHandleBoundingBox.java#shouldReturnCorrectDimensions
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
// Ref: TestElementHandleBoundingBox.java#shouldReturnNilForDisplayNone
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
// Ref: TestElementHandleBoundingBox.java#shouldUpdateAfterScroll
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
// Ref: TestElementHandleBoundingBox.java#shouldReturnCorrectPosition
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

func TestElementHandleInnerHTML(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el"><b>hello</b></div>`)
	must.NoError(err, "SetContent failed")

	el, err := page.QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el, "QuerySelector failed")

	html, err := el.InnerHTML(ctx)
	must.NoError(err, "InnerHTML failed")
	is.Equal("<b>hello</b>", html)
}

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

func TestElementHandleIsEnabledAndDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="ok">OK</button><button id="no" disabled>No</button>`)
	must.NoError(err, "SetContent failed")

	elOK, _ := page.QuerySelector(ctx, "#ok")
	elNo, _ := page.QuerySelector(ctx, "#no")

	enabled, err := elOK.IsEnabled(ctx)
	if err != nil || !enabled {
		t.Errorf("#ok should be enabled, got %v err=%v", enabled, err)
	}

	disabled, err := elNo.IsEnabled(ctx)
	if err != nil || disabled {
		t.Errorf("#no should be disabled (IsEnabled=false), got %v err=%v", disabled, err)
	}
}

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

func TestElementHandleSelectOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select id="sel"><option value="foo">Foo</option><option value="bar">Bar</option></select>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#sel")

	selected, err := el.SelectOption(ctx, "bar")
	must.NoError(err, "SelectOption failed")
	if len(selected) != 1 || selected[0] != "bar" {
		t.Errorf("SelectOption returned %v, want ['bar']", selected)
	}
}

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

func TestElementHandleWaitForVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" style="display:none">hi</div>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#el")

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('el').style.display = 'block'`)
	}()

	err = el.WaitForElementState(ctx, "visible")
	must.NoError(err, "WaitForElementState(visible) failed")
}

func TestElementHandleWaitForAlreadyVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">visible</div>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#el")

	err = el.WaitForElementState(ctx, "visible")
	must.NoError(err, "WaitForElementState(visible) for already-visible element failed")
}

func TestElementHandleWaitForHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">visible</div>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#el")

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('el').style.display = 'none'`)
	}()

	err = el.WaitForElementState(ctx, "hidden")
	must.NoError(err, "WaitForElementState(hidden) failed")
}

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

func TestElementHandleWaitForEnabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn" disabled>Click</button>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#btn")

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('btn').disabled = false`)
	}()

	err = el.WaitForElementState(ctx, "enabled")
	must.NoError(err, "WaitForElementState(enabled) failed")
}

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

// TestElementHandleBoundingBoxEx verifies ElementHandle BoundingBox returns valid dimensions.
// Ref: TestElementHandle.java#shouldGetBoundingBox
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

// TestElementHandleGetAttributeAbsentEx verifies GetAttribute returns nil for missing attribute.
// Ref: TestElementHandle.java#shouldReturnNilForAbsentAttribute
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
// Ref: TestElementHandle.java#shouldCheckCheckbox
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
// Ref: TestElementHandle.java#shouldUncheckCheckbox
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
// Ref: TestElementHandle.java#shouldFillAndGetValue
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
// Ref: TestElementHandle.java#shouldReturnNullForNonExistingAttribute
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
// Ref: TestElementHandle.java#shouldReturnFalseForHiddenElement
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
// Ref: TestElementHandle.java#shouldReturnTrueForEnabledInput
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
// Ref: TestElementHandle.java#shouldReturnFalseForDisabledInput
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

// TestElementHandleSelectOptionByValue verifies SelectOption by value on select element.
// Ref: TestElementHandle.java#shouldSelectOptionByValue
func TestElementHandleSelectOptionByValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="one">One</option>
			<option value="two">Two</option>
			<option value="three">Three</option>
		</select>
	`))

	el, err := page.QuerySelector(ctx, "#sel")
	must.NoError(err)
	must.NotNil(el)

	selected, err := el.SelectOption(ctx, "two")
	must.NoError(err)
	is.Equal([]string{"two"}, selected)
}

// TestElementHandleGetBoundingBoxForVisibleElement verifies BoundingBox returns non-nil for visible element.
// Ref: TestElementHandle.java#shouldWorkForVisibleElement
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

// TestElementHandleInnerHTMLReturnsContent verifies InnerHTML returns the inner HTML string.
// Ref: TestElementHandle.java#shouldWorkForInnerHTML
func TestElementHandleInnerHTMLReturnsContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"><span>hello</span></div>`))

	el, err := page.QuerySelector(ctx, "#container")
	must.NoError(err)
	must.NotNil(el)

	html, err := el.InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<span>hello</span>", html)
}

// TestElementHandleInnerTextReturnsText verifies InnerText returns visible text.
// Ref: TestElementHandle.java#shouldReturnVisibleText
func TestElementHandleInnerTextReturnsText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello <b>World</b></p>`))

	el, err := page.QuerySelector(ctx, "#p")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello World", text)
}
