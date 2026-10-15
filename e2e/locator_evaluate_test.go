//go:build e2e

// E2E tests for: TestLocatorEvaluate.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
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

// Ref: TestLocatorEvaluate.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
func TestLocatorEvaluateShouldWorkA3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<html><body><div class='tweet'><div class='like'>100</div><div class='retweets'>10</div></div></body></html>"))
	content, err := page.Locator(".tweet .like").Evaluate(ctx, "node => node.innerText")
	must.NoError(err)
	assert.Equal(t, "100", content)
}

// Ref: TestLocatorEvaluate.java#shouldRetrieveContentFromSubtree  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
func TestLocatorEvaluateShouldRetrieveContentFromSubtree(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div class='a'>not-a-child-div</div><div id='myId'><div class='a'>a-child-div</div></div>"))
	content, err := page.Locator("#myId .a").Evaluate(ctx, "node => node.innerText")
	must.NoError(err)
	assert.Equal(t, "a-child-div", content)
}

// Ref: TestLocatorEvaluate.java#shouldWorkForAll  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
func TestLocatorEvaluateShouldWorkForAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<html><body><div class='tweet'><div class='like'>100</div><div class='like'>10</div></div></body></html>"))
	content, err := page.Locator(".tweet .like").EvaluateAll(ctx, "nodes => nodes.map(n => n.innerText)")
	must.NoError(err)
	assert.Equal(t, []any{"100", "10"}, content)
}

// Ref: TestLocatorEvaluate.java#shouldRetrieveContentFromSubtreeForAll  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
func TestLocatorEvaluateShouldRetrieveContentFromSubtreeForAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div class='a'>not-a-child-div</div><div id='myId'><div class='a'>a1-child-div</div><div class='a'>a2-child-div</div></div>"))
	content, err := page.Locator("#myId .a").EvaluateAll(ctx, "nodes => nodes.map(n => n.innerText)")
	must.NoError(err)
	assert.Equal(t, []any{"a1-child-div", "a2-child-div"}, content)
}

// Ref: TestLocatorEvaluate.java#shouldNotThrowInCaseOfMissingSelectorForAll  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorEvaluate.java
func TestLocatorEvaluateShouldNotThrowInCaseOfMissingSelectorForAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div class='a'>not-a-child-div</div><div id='myId'></div>"))
	nodesLength, err := page.Locator("#myId .a").EvaluateAll(ctx, "nodes => nodes.length")
	must.NoError(err)
	assert.EqualValues(t, 0, nodesLength)
}
