//go:build e2e

// Locator.Evaluate and Locator.EvaluateAll E2E tests.
// Migration of: TestLocatorEvaluate.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocatorEvaluateShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="a">hello</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator(".a").Evaluate(ctx, "e => e.textContent")
	must.NoError(err, "Evaluate failed")
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

func TestLocatorEvaluateRetrievesFromSubtree(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="a"><div class="b">hello</div></div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator(".a").Evaluate(ctx, "e => e.querySelector('.b').textContent")
	must.NoError(err, "Evaluate failed")
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

func TestLocatorEvaluateForAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="a">one</div><div class="a">two</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator(".a").EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err, "EvaluateAll failed")
	arr, ok := val.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T: %v", val, val)
	}
	if len(arr) != 2 || arr[0] != "one" || arr[1] != "two" {
		t.Errorf("unexpected result: %v", arr)
	}
}

func TestLocatorEvaluateAllRetrievesFromSubtree(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="a"><div class="b">hello</div><div class="b">world</div></div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator(".b").EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err, "EvaluateAll failed")
	arr, ok := val.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", val)
	}
	if len(arr) != 2 || arr[0] != "hello" || arr[1] != "world" {
		t.Errorf("unexpected result: %v", arr)
	}
}

func TestLocatorEvaluateAllNotThrowOnMissingSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="a">hello</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator(".missing").EvaluateAll(ctx, "els => els.length")
	must.NoError(err, "EvaluateAll on missing selector should not throw")
	if val != float64(0) {
		t.Errorf("expected 0 for missing selector, got %v", val)
	}
}

func TestLocatorEvaluateReturnsBooleanProperty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" checked>`))

	result, err := page.Locator("input").Evaluate(ctx, `el => el.checked`)
	must.NoError(err)
	is.Equal(true, result)
}

func TestLocatorEvaluateReturnsNumericProperty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul><li>a</li><li>b</li><li>c</li></ul>`))

	result, err := page.Locator("ul").Evaluate(ctx, `ul => ul.querySelectorAll('li').length`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

func TestLocatorEvaluatePassesArgument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class="box">content</div>`))

	result, err := page.Locator(".box").Evaluate(ctx, `(el, suffix) => el.textContent + suffix`, " end")
	must.NoError(err)
	is.Equal("content end", result)
}

func TestLocatorEvaluateAllSumsValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">10</div>
		<div class="item">20</div>
		<div class="item">30</div>
	`))

	result, err := page.Locator(".item").EvaluateAll(ctx, `els => els.reduce((sum, el) => sum + parseInt(el.textContent), 0)`)
	must.NoError(err)
	is.Equal(float64(60), result)
}

func TestLocatorEvaluateAllEmptyReturnsEmptyArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>no items here</div>`))

	result, err := page.Locator(".no-match").EvaluateAll(ctx, `els => els.length`)
	must.NoError(err)
	is.Equal(float64(0), result)
}

func TestLocatorEvaluateModifiesElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">original</div>`))

	_, err := page.Locator("#target").Evaluate(ctx, `el => { el.textContent = 'modified'; }`)
	must.NoError(err)

	text, err := page.Locator("#target").InnerText(ctx)
	must.NoError(err)
	is.Equal("modified", text)
}

func TestLocatorEvaluateReturnsTagName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">content</div>`))

	val, err := page.Locator("#el").Evaluate(ctx, "el => el.tagName.toLowerCase()")
	must.NoError(err)
	is.Equal("div", val)
}

func TestLocatorEvaluateComputesStyle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" style="color:red">text</div>`))

	val, err := page.Locator("#el").Evaluate(ctx,
		`el => window.getComputedStyle(el).color`)
	must.NoError(err)
	is.NotEmpty(val)
}

func TestLocatorEvaluateModifiesText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">original</p>`))

	_, err := page.Locator("#p").Evaluate(ctx, `el => { el.textContent = 'modified'; }`)
	must.NoError(err)

	text, err := page.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("modified", text)
}

func TestLocatorEvaluateWithArg(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">0</div>`))

	val, err := page.Locator("#el").Evaluate(ctx, `(el, n) => n * 3`, 7)
	must.NoError(err)
	is.Equal(float64(21), val)
}

func TestLocatorEvaluateAllCountsItems(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>A</li><li>B</li><li>C</li>
		</ul>
	`))

	val, err := page.Locator("li").EvaluateAll(ctx, `els => els.length`)
	must.NoError(err)
	is.Equal(float64(3), val)
}

func TestLocatorEvaluateAllExtractsText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="item">X</span>
		<span class="item">Y</span>
		<span class="item">Z</span>
	`))

	val, err := page.Locator(".item").EvaluateAll(ctx, `els => els.map(e => e.textContent)`)
	must.NoError(err)
	slice, ok := val.([]any)
	is.True(ok)
	is.Len(slice, 3)
	is.Equal("X", slice[0])
}

func TestLocatorEvaluateReturnsTagNameEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.tagName`)
	must.NoError(err)
	is.Equal("DIV", result)
}

func TestLocatorEvaluateComputesStyleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="color:rgb(255,0,0)">Red text</p>`))

	result, err := page.Locator("#p").Evaluate(ctx, `el => window.getComputedStyle(el).color`)
	must.NoError(err)
	is.Contains(result.(string), "255")
}

func TestLocatorEvaluateModifiesDOMEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Original</p>`))

	_, err := page.Locator("#p").Evaluate(ctx, `el => { el.textContent = 'Modified'; }`)
	must.NoError(err)

	text, err := page.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("Modified", text)
}

func TestLocatorEvaluateWithArgumentEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Value</div>`))

	result, err := page.Locator("#d").Evaluate(ctx, `(el, n) => n + 10`, 5)
	must.NoError(err)
	is.Equal(float64(15), result)
}

func TestLocatorEvaluateCountSiblingsEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul><li id="li1">1</li><li>2</li><li>3</li></ul>
	`))

	result, err := page.Locator("#li1").Evaluate(ctx, `el => el.parentElement.children.length`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

func TestLocatorEvaluateClassListEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="foo bar baz">Content</div>`))

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.classList.contains('bar')`)
	must.NoError(err)
	is.Equal(true, result)
}

func TestLocatorEvaluateStyleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="color: red;">Colored</div>`))

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.style.color`)
	must.NoError(err)
	is.Equal("red", result)
}

func TestLocatorEvaluateInnerHTMLEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><b>Bold</b></div>`))

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.innerHTML`)
	must.NoError(err)
	is.Equal("<b>Bold</b>", result)
}

func TestLocatorEvaluateTagNameEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id="s">Content</section>`))

	result, err := page.Locator("#s").Evaluate(ctx, `el => el.tagName.toLowerCase()`)
	must.NoError(err)
	is.Equal("section", result)
}

func TestLocatorEvaluateChildCountEx4(t *testing.T) {
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

	result, err := page.Locator("#list").Evaluate(ctx, `el => el.children.length`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

func TestLocatorEvaluateScrollHeightEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="overflow:auto;height:100px;">
			<div style="height:500px;">Content</div>
		</div>
	`))

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.scrollHeight`)
	must.NoError(err)
	is.Greater(result.(float64), float64(100))
}

func TestLocatorEvaluateClientWidthEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="width:200px;height:50px;"></div>
	`))

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.clientWidth`)
	must.NoError(err)
	is.Equal(float64(200), result)
}

func TestLocatorEvaluateNextSiblingEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="first">First</div>
		<div id="second">Second</div>
	`))

	result, err := page.Locator("#first").Evaluate(ctx, `el => el.nextElementSibling.id`)
	must.NoError(err)
	is.Equal("second", result)
}

func TestLocatorEvaluateParentIDEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			<span id="child">Child</span>
		</div>
	`))

	result, err := page.Locator("#child").Evaluate(ctx, `el => el.parentElement.id`)
	must.NoError(err)
	is.Equal("parent", result)
}

func TestLocatorEvaluateOffsetTopEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:100px;"></div>
		<div id="below" style="height:50px;"></div>
	`))

	result, err := page.Locator("#below").Evaluate(ctx, `el => el.offsetTop`)
	must.NoError(err)
	is.Greater(result.(float64), float64(0))
}

func TestLocatorEvaluateComputedStyleEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="color: red;">Text</div>`))

	color, err := page.Locator("#d").Evaluate(ctx, `el => window.getComputedStyle(el).color`)
	must.NoError(err)
	is.NotEmpty(color)
}

func TestLocatorEvaluateScrollHeightEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="height:50px;overflow:auto;"><p>line1</p><p>line2</p></div>`))

	height, err := page.Locator("#d").Evaluate(ctx, `el => el.scrollHeight`)
	must.NoError(err)
	must.NotNil(height)
}

func TestLocatorEvaluateDatasetEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" data-user="john" data-role="admin">Content</div>`))

	user, err := page.Locator("#d").Evaluate(ctx, `el => el.dataset.user`)
	must.NoError(err)
	is.Equal("john", user)
}

func TestLocatorEvaluateChildCountEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"><li>A</li><li>B</li><li>C</li></ul>`))

	count, err := page.Locator("#list").Evaluate(ctx, `el => el.children.length`)
	must.NoError(err)
	is.Equal(float64(3), count)
}

func TestLocatorEvaluateReturnsBoolEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	checked, err := page.Locator("#chk").Evaluate(ctx, `el => el.checked`)
	must.NoError(err)
	is.Equal(true, checked)
}

func TestLocatorEvaluateGetBoundingClientRectEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="width:100px;height:50px;">Content</div>`))

	rect, err := page.Locator("#d").Evaluate(ctx, `el => { const r = el.getBoundingClientRect(); return {w: r.width, h: r.height}; }`)
	must.NoError(err)
	m, ok := rect.(map[string]interface{})
	is.True(ok)
	is.Equal(float64(100), m["w"])
	is.Equal(float64(50), m["h"])
}

func TestLocatorEvaluateAddClassEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	_, err := page.Locator("#d").Evaluate(ctx, `el => el.classList.add('active')`)
	must.NoError(err)

	cls, err := page.Evaluate(ctx, `() => document.getElementById('d').className`)
	must.NoError(err)
	is.Equal("active", cls)
}

func TestLocatorEvaluateSetStyleEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	_, err := page.Locator("#d").Evaluate(ctx, `el => { el.style.color = 'blue'; }`)
	must.NoError(err)

	color, err := page.Evaluate(ctx, `() => document.getElementById('d').style.color`)
	must.NoError(err)
	is.Equal("blue", color)
}

func TestLocatorEvaluateGetParentIdEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="parent"><span id="child">Child</span></div>`))

	parentId, err := page.Locator("#child").Evaluate(ctx, `el => el.parentElement.id`)
	must.NoError(err)
	is.Equal("parent", parentId)
}

func TestLocatorEvaluateHasAttributeEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" required>`))

	hasRequired, err := page.Locator("#inp").Evaluate(ctx, `el => el.hasAttribute('required')`)
	must.NoError(err)
	is.Equal(true, hasRequired)
}

func TestLocatorEvaluateInnerHTMLEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><b>Bold</b> text</div>`))

	html, err := page.Locator("#d").Evaluate(ctx, `el => el.innerHTML`)
	must.NoError(err)
	is.Contains(html, "<b>Bold</b>")
}

func TestLocatorEvaluateOuterHTMLEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" class="text">Content</p>`))

	html, err := page.Locator("#p").Evaluate(ctx, `el => el.outerHTML`)
	must.NoError(err)
	is.Contains(html, `id="p"`)
	is.Contains(html, "Content")
}

func TestLocatorEvaluateOffsetWidthEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="width:150px;height:50px;">Content</div>`))

	width, err := page.Locator("#d").Evaluate(ctx, `el => el.offsetWidth`)
	must.NoError(err)
	is.Equal(float64(150), width)
}

func TestLocatorEvaluateNodeNameEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	name, err := page.Locator("#btn").Evaluate(ctx, `el => el.nodeName`)
	must.NoError(err)
	is.Equal("BUTTON", name)
}

func TestLocatorEvaluateTypePropertyEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="email">`))

	typ, err := page.Locator("#inp").Evaluate(ctx, `el => el.type`)
	must.NoError(err)
	is.Equal("email", typ)
}

func TestEvaluateFormActionEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<form id="f" action="/submit" method="post"><input type="submit"></form>`))

	result, err := page.Locator("#f").Evaluate(ctx, `el => el.action`)
	must.NoError(err)
	s, ok := result.(string)
	is.True(ok)
	is.Contains(s, "/submit")
}

func TestEvaluateTagNameEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id="s">Section</section>`))

	result, err := page.Locator("#s").Evaluate(ctx, `el => el.tagName`)
	must.NoError(err)
	is.Equal("SECTION", result)
}

func TestEvaluateChildrenCountEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			<span>1</span><span>2</span><span>3</span><span>4</span>
		</div>
	`))

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.children.length`)
	must.NoError(err)
	is.Equal(float64(4), result)
}

func TestEvaluateInputTypeEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="range" min="0" max="100" value="50">`))

	result, err := page.Locator("#inp").Evaluate(ctx, `el => el.type`)
	must.NoError(err)
	is.Equal("range", result)
}

func TestEvaluateScrollTopEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="height:100px;overflow:auto;">
			<div style="height:500px;">Tall content</div>
		</div>
	`))

	_, err := page.Locator("#d").Evaluate(ctx, `el => { el.scrollTop = 50; return el.scrollTop; }`)
	must.NoError(err)

	result, err := page.Locator("#d").Evaluate(ctx, `el => el.scrollTop`)
	must.NoError(err)
	is.Equal(float64(50), result)
}
