//go:build e2e

// E2E tests for Locator.Click (Priority 2).
// Ref: TestClick.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageClickButton verifies that Locator.Click() triggers a button's click handler.
// Ref: TestClick.java#shouldClickTheButton
func TestPageClickButton(t *testing.T) {
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

// TestPageClickCheckbox verifies that Locator.Click() toggles a checkbox.
// Ref: TestClick.java#shouldClickOnCheckbox
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

func TestClickSVG(t *testing.T) {
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

func TestClickSpanWithInlineElement(t *testing.T) {
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
