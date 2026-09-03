//go:build e2e

// EvalOnSelector and EvalOnSelectorAll E2E tests.
// Migration of: TestEvalOnSelector.java, TestEvalOnSelectorAll.java,
// TestPageEvalOnSelector.java, TestPageQuerySelector.java, TestPageQuerySelectorAll.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvalOnSelectorWithCSSSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section id="testAttribute">43543</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "css=section", "e => e.id")
	must.NoError(err, "EvalOnSelector failed")
	is.Equal("testAttribute", val)
}

func TestEvalOnSelectorWithIDSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section id="testAttribute">43543</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "#testAttribute", "e => e.id")
	must.NoError(err, "EvalOnSelector failed")
	is.Equal("testAttribute", val)
}

func TestEvalOnSelectorAcceptsArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section>hello</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "section", "(e, suffix) => e.textContent + suffix", " world")
	must.NoError(err, "EvalOnSelector failed")
	is.Equal("hello world", val)
}

func TestEvalOnSelectorReturnsComplexValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section>hello</section>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelector(ctx, "section", "e => ({ tag: e.tagName, text: e.textContent })")
	must.NoError(err, "EvalOnSelector failed")
	m, ok := val.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", val, val)
	if m["tag"] != "SECTION" {
		t.Errorf("tag = %v, want 'SECTION'", m["tag"])
	}
	if m["text"] != "hello" {
		t.Errorf("text = %v, want 'hello'", m["text"])
	}
}

func TestEvalOnSelectorThrowsWhenNoElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section>value</section>`)
	must.NoError(err, "SetContent failed")

	_, err = page.EvalOnSelector(ctx, "div", "e => e.id")
	is.Error(err)
}

func TestEvalOnSelectorAllWithCSSSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>A</div><div>B</div><div>C</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelectorAll(ctx, "div", "els => els.map(e => e.textContent)")
	must.NoError(err, "EvalOnSelectorAll failed")
	arr, ok := val.([]any)
	must.Truef(ok, "expected []any, got %T: %v", val, val)
	is.Len(arr, 3)
	if arr[0] != "A" || arr[1] != "B" || arr[2] != "C" {
		t.Errorf("unexpected values: %v", arr)
	}
}

func TestEvalOnSelectorAllReturnsEmptyArrayWhenNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>content</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.EvalOnSelectorAll(ctx, "span", "els => els.length")
	must.NoError(err, "EvalOnSelectorAll failed")
	is.Equal(float64(0), val)
}

func TestEvalOnSelectorReturnsUndefinedForMissingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id="testAttribute">43543</section>`))

	result, err := page.EvalOnSelector(ctx, "css=div", `e => e.id`)
	must.NoError(err)
	is.Nil(result)
}

func TestEvalOnSelectorReturnsInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>hello</section>`))

	result, err := page.EvalOnSelector(ctx, "css=section", `e => e.innerText`)
	must.NoError(err)
	is.Equal("hello", result)
}

func TestEvalOnSelectorWithAttributeSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div data-testid="box">content</div>`))

	result, err := page.EvalOnSelector(ctx, `[data-testid="box"]`, `e => e.textContent`)
	must.NoError(err)
	is.Equal("content", result)
}

func TestEvalOnSelectorWithArgument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class="myDiv">original</div>`))

	result, err := page.EvalOnSelector(ctx, ".myDiv", `(el, value) => el.textContent + value`, " suffix")
	must.NoError(err)
	is.Equal("original suffix", result)
}

func TestEvalOnSelectorReturnsBooleanTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" checked>`))

	result, err := page.EvalOnSelector(ctx, "input", `e => e.checked`)
	must.NoError(err)
	is.Equal(true, result)
}

func TestEvalOnSelectorReturnsBooleanFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox">`))

	result, err := page.EvalOnSelector(ctx, "input", `e => e.checked`)
	must.NoError(err)
	is.Equal(false, result)
}

func TestEvalOnSelectorAllReturnsTextContents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>one</div>
		<div>two</div>
		<div>three</div>
	`))

	result, err := page.EvalOnSelectorAll(ctx, "div", `els => els.map(e => e.textContent.trim())`)
	must.NoError(err)

	arr, ok := result.([]any)
	is.True(ok)
	is.Len(arr, 3)
	is.Equal("one", arr[0])
	is.Equal("two", arr[1])
	is.Equal("three", arr[2])
}

func TestEvalOnSelectorAllCountElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>a</li><li>b</li><li>c</li><li>d</li>
	`))

	result, err := page.EvalOnSelectorAll(ctx, "li", `els => els.length`)
	must.NoError(err)
	is.Equal(float64(4), result)
}

// TestEvalOnSelectorGetsBoundingRect verifies EvalOnSelector can query layout.
// Ref: TestPageEvalOnSelector.java#shouldGetBoundingRect
func TestEvalOnSelectorGetsBoundingRect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="width:100px;height:50px">box</div>`))

	result, err := page.EvalOnSelector(ctx, "div", `el => el.getBoundingClientRect().width`)
	must.NoError(err)
	is.Equal(float64(100), result)
}

// TestEvalOnSelectorSetsAttribute verifies EvalOnSelector can modify DOM.
// Ref: TestPageEvalOnSelector.java#shouldSetAttribute
func TestEvalOnSelectorSetsAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">text</div>`))

	_, err := page.EvalOnSelector(ctx, "#el", `el => el.setAttribute('data-x', 'yes')`)
	must.NoError(err)

	val, err := page.Locator("#el").GetAttribute(ctx, "data-x")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("yes", *val)
}

// TestEvalOnSelectorAllSumsValuesExtra verifies EvalOnSelectorAll can aggregate.
// Ref: TestPageEvalOnSelector.java#shouldSumValuesWithAll
func TestEvalOnSelectorAllSumsValuesExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="number" value="10">
		<input type="number" value="20">
		<input type="number" value="30">
	`))

	result, err := page.EvalOnSelectorAll(ctx, "input", `els => els.reduce((s, el) => s + parseInt(el.value), 0)`)
	must.NoError(err)
	is.Equal(float64(60), result)
}

// TestEvalOnSelectorReturnsNthChildText verifies EvalOnSelector works on nth child.
// Ref: TestPageEvalOnSelector.java#shouldWorkOnNthChild
func TestEvalOnSelectorReturnsNthChildText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>first</li>
			<li>second</li>
			<li>third</li>
		</ul>
	`))

	result, err := page.EvalOnSelector(ctx, "li:nth-child(2)", `el => el.textContent`)
	must.NoError(err)
	is.Equal("second", result)
}

// TestEvalOnSelectorAllMapClassNames verifies EvalOnSelectorAll maps class names.
// Ref: TestPageEvalOnSelector.java#shouldMapClassNames
func TestEvalOnSelectorAllMapClassNames(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="a">1</div>
		<div class="b">2</div>
		<div class="c">3</div>
	`))

	result, err := page.EvalOnSelectorAll(ctx, "div", `els => els.map(el => el.className)`)
	must.NoError(err)
	must.NotNil(result)
	list, ok := result.([]any)
	is.True(ok)
	is.Len(list, 3)
	is.Equal("a", list[0])
}

// TestPageQuerySelectorReturnsNilWhenNotFound verifies QuerySelector returns nil when selector has no match.
// Ref: TestPageQuerySelector.java#shouldReturnNullForNonExistingElement
func TestPageQuerySelectorReturnsNilWhenNotFound(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err)
	is.Nil(el, "QuerySelector should return nil for non-existing element")
}

// TestPageQuerySelectorFindsById verifies QuerySelector works with ID selector.
// Ref: TestPageQuerySelector.java#shouldWorkWithIdSelector
func TestPageQuerySelectorFindsById(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">found</div>`))

	el, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(el, "QuerySelector should find element by id")

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("found", text)
}

// TestPageQuerySelectorFindsByClass verifies QuerySelector works with class selector.
// Ref: TestPageQuerySelector.java#shouldWorkWithClassSelector
func TestPageQuerySelectorFindsByClass(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">first</div>
		<div class="item">second</div>
	`))

	// QuerySelector returns the FIRST match
	el, err := page.QuerySelector(ctx, ".item")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("first", text)
}

// TestPageQuerySelectorAllReturnsAll verifies QuerySelectorAll returns all matches.
// Ref: TestPageQuerySelector.java#shouldReturnAllMatchingElements
func TestPageQuerySelectorAllReturnsAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="tag">A</span>
		<span class="tag">B</span>
		<span class="tag">C</span>
	`))

	handles, err := page.QuerySelectorAll(ctx, "span.tag")
	must.NoError(err)
	is.Len(handles, 3)

	expected := []string{"A", "B", "C"}
	for i, el := range handles {
		text, err := el.InnerText(ctx)
		must.NoError(err)
		is.Equal(expected[i], text)
	}
}

// TestPageQuerySelectorAllReturnsEmptySlice verifies QuerySelectorAll returns empty slice when no match.
// Ref: TestPageQuerySelector.java#shouldReturnEmptyArrayWhenNoMatch
func TestPageQuerySelectorAllReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	handles, err := page.QuerySelectorAll(ctx, "span.nonexistent")
	must.NoError(err)
	is.Empty(handles)
}

// TestPageQuerySelectorWithAttributeSelector verifies attribute-based CSS selectors.
// Ref: TestPageQuerySelector.java#shouldWorkWithAttributeSelector
func TestPageQuerySelectorWithAttributeSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" data-test="first">
		<input type="password" data-test="second">
	`))

	el, err := page.QuerySelector(ctx, `input[type="password"]`)
	must.NoError(err)
	must.NotNil(el)

	attr, err := el.GetAttribute(ctx, "data-test")
	must.NoError(err)
	is.Equal("second", attr)
}

// TestPageQuerySelectorAllWithNestedSelector verifies QuerySelectorAll with descendant selectors.
// Ref: TestPageQuerySelector.java#shouldWorkWithNestedSelectors
func TestPageQuerySelectorAllWithNestedSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="container">
			<p>one</p>
			<p>two</p>
		</div>
		<p>three</p>
	`))

	handles, err := page.QuerySelectorAll(ctx, ".container p")
	must.NoError(err)
	is.Len(handles, 2)
}

// TestQuerySelectorFindsButtonEx2 verifies QuerySelector finds button by id.
// Ref: TestPageQuerySelector.java#shouldFindButton
func TestQuerySelectorFindsButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Submit</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorReturnsNilForMissingEx2 verifies QuerySelector returns nil for no match.
// Ref: TestPageQuerySelector.java#shouldReturnNilForMissing
func TestQuerySelectorReturnsNilForMissingEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>No button here</div>`))

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err)
	is.Nil(el)
}

// TestQuerySelectorAllFindsMultipleEx2 verifies QuerySelectorAll returns multiple elements.
// Ref: TestPageQuerySelector.java#shouldFindMultiple
func TestQuerySelectorAllFindsMultipleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">A</div>
		<div class="item">B</div>
		<div class="item">C</div>
	`))

	elements, err := page.QuerySelectorAll(ctx, ".item")
	must.NoError(err)
	is.Len(elements, 3)
}

// TestQuerySelectorByCSSClassEx2 verifies QuerySelector finds element by class.
// Ref: TestPageQuerySelector.java#shouldFindByClass
func TestQuerySelectorByCSSClassEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p class="highlight">Highlighted text</p>`))

	el, err := page.QuerySelector(ctx, ".highlight")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorByTagEx2 verifies QuerySelector finds element by tag.
// Ref: TestPageQuerySelector.java#shouldFindByTag
func TestQuerySelectorByTagEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h2>A heading</h2>`))

	el, err := page.QuerySelector(ctx, "h2")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("A heading", text)
}

// TestQuerySelectorAllReturnsMultiple verifies QuerySelectorAll returns all matching elements.
// Ref: TestPageQuerySelectorAll.java#shouldReturnMultiple
func TestQuerySelectorAllReturnsMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">a</div>
		<div class="item">b</div>
		<div class="item">c</div>
	`))

	handles, err := page.QuerySelectorAll(ctx, ".item")
	must.NoError(err)
	is.Len(handles, 3)
}

// TestQuerySelectorAllReturnsEmpty verifies QuerySelectorAll returns empty for no match.
// Ref: TestPageQuerySelectorAll.java#shouldReturnEmptyForNoMatch
func TestQuerySelectorAllReturnsEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))

	handles, err := page.QuerySelectorAll(ctx, ".nonexistent")
	must.NoError(err)
	is.Empty(handles)
}

// TestQuerySelectorAllTextContent verifies TextContent from each QuerySelectorAll result.
// Ref: TestPageQuerySelectorAll.java#shouldGetTextContent
func TestQuerySelectorAllTextContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>first</p>
		<p>second</p>
	`))

	handles, err := page.QuerySelectorAll(ctx, "p")
	must.NoError(err)
	is.Len(handles, 2)

	text0, err := handles[0].TextContent(ctx)
	must.NoError(err)
	is.Equal("first", text0)

	text1, err := handles[1].TextContent(ctx)
	must.NoError(err)
	is.Equal("second", text1)
}

// TestQuerySelectorAllInputValues verifies InputValue from QuerySelectorAll inputs.
// Ref: TestPageQuerySelectorAll.java#shouldGetInputValues
func TestQuerySelectorAllInputValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input value="one">
		<input value="two">
	`))

	handles, err := page.QuerySelectorAll(ctx, "input")
	must.NoError(err)
	is.Len(handles, 2)

	v0, err := handles[0].InputValue(ctx)
	must.NoError(err)
	is.Equal("one", v0)

	v1, err := handles[1].InputValue(ctx)
	must.NoError(err)
	is.Equal("two", v1)
}

// TestQuerySelectorAllSingleMatch verifies QuerySelectorAll with single match.
// Ref: TestPageQuerySelectorAll.java#shouldWorkWithSingleMatch
func TestQuerySelectorAllSingleMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="title">hello</h1>`))

	handles, err := page.QuerySelectorAll(ctx, "h1")
	must.NoError(err)
	is.Len(handles, 1)

	text, err := handles[0].InnerText(ctx)
	must.NoError(err)
	is.Equal("hello", text)
}

// deepShadowHTML is the content served for shadow-DOM selector tests.
const deepShadowHTML = `<script>
window.addEventListener('DOMContentLoaded', () => {
  const outer = document.createElement('section');
  document.body.appendChild(outer);

  const root1 = document.createElement('div');
  root1.setAttribute('id', 'root1');
  outer.appendChild(root1);
  const shadowRoot1 = root1.attachShadow({mode: 'open'});
  const span1 = document.createElement('span');
  span1.setAttribute('data-testid', 'foo');
  span1.textContent = 'Hello from root1';
  shadowRoot1.appendChild(span1);

  const root2 = document.createElement('div');
  shadowRoot1.appendChild(root2);
  const shadowRoot2 = root2.attachShadow({mode: 'open'});
  const span2 = document.createElement('span');
  span2.setAttribute('data-testid', 'foo');
  span2.setAttribute('id', 'target');
  span2.textContent = 'Hello from root2';
  shadowRoot2.appendChild(span2);

  const root3 = document.createElement('div');
  shadowRoot1.appendChild(root3);
  const shadowRoot3 = root3.attachShadow({mode: 'open'});
  const span3 = document.createElement('span');
  span3.setAttribute('data-testid', 'foo');
  span3.textContent = 'Hello from root3';
  shadowRoot3.appendChild(span3);
  const span4 = document.createElement('span');
  span4.textContent = 'Hello from root3 #2';
  span4.setAttribute('attr', 'value space');
  shadowRoot3.appendChild(span4);
});
</script>`

// TestEvalOnSelectorWithDataTestSelector verifies the data-test= attribute selector engine.
// Ref: TestEvalOnSelector.java#shouldWorkWithDataTestSelector
func TestEvalOnSelectorWithDataTestSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section data-test=foo id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "data-test=foo", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorWithDataTestidSelector verifies the data-testid= attribute selector engine.
// Ref: TestEvalOnSelector.java#shouldWorkWithDataTestidSelector
func TestEvalOnSelectorWithDataTestidSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section data-testid=foo id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "data-testid=foo", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorWithDataTestIdSelector verifies the data-test-id= attribute selector engine.
// Ref: TestEvalOnSelector.java#shouldWorkWithDataTestIdSelector
func TestEvalOnSelectorWithDataTestIdSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section data-test-id=foo id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "data-test-id=foo", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorWithTextSelector1 verifies text= selector with a single-quoted value.
// Ref: TestEvalOnSelector.java#shouldWorkWithTextSelector1
func TestEvalOnSelectorWithTextSelector1(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "text='43543'", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorWithXpathSelector verifies the xpath= selector engine prefix.
// Ref: TestEvalOnSelector.java#shouldWorkWithXpathSelector
func TestEvalOnSelectorWithXpathSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "xpath=/html/body/section", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorWithTextSelector2 verifies text= selector without quotes does a substring match.
// Ref: TestEvalOnSelector.java#shouldWorkWithTextSelector2
func TestEvalOnSelectorWithTextSelector2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "text=43543", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorAutoDetectCssSelectorWithAttributes verifies CSS attribute selectors are auto-detected without an engine prefix.
// Ref: TestEvalOnSelector.java#shouldAutoDetectCssSelectorWithAttributes
func TestEvalOnSelectorAutoDetectCssSelectorWithAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section id='testAttribute'>43543</section>`))

	val, err := page.EvalOnSelector(ctx, "section[id='testAttribute']", "e => e.id")
	must.NoError(err)
	is.Equal("testAttribute", val)
}

// TestEvalOnSelectorAutoDetectNestedSelectors verifies >> chains with auto-detected CSS and quoted text selectors.
// Ref: TestEvalOnSelector.java#shouldAutoDetectNestedSelectors
func TestEvalOnSelectorAutoDetectNestedSelectors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div foo=bar><section>43543<span>Hello<div id=target></div></span></section></div>`))

	val, err := page.EvalOnSelector(ctx, "div[foo=bar] > section >> 'Hello' >> div", "e => e.id")
	must.NoError(err)
	is.Equal("target", val)
}

// TestEvalOnSelectorAcceptsElementHandlesAsArguments verifies an ElementHandle can be passed as an extra argument.
// Ref: TestEvalOnSelector.java#shouldAcceptElementHandlesAsArguments
func TestEvalOnSelectorAcceptsElementHandlesAsArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>hello</section><div> world</div>`))

	divHandle, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(divHandle)

	val, err := page.EvalOnSelector(ctx, "section", "(e, div) => e.textContent + div.textContent", divHandle)
	must.NoError(err)
	is.Equal("hello world", val)
}

// TestEvalOnSelectorSupportsSyntax verifies the css=A >> css=B chained selector syntax.
// Ref: TestEvalOnSelector.java#shouldSupportSyntax
func TestEvalOnSelectorSupportsSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div>hello</div></section>`))

	val, err := page.EvalOnSelector(ctx, "css=section >> css=div", "(e, suffix) => e.textContent + suffix", " world!")
	must.NoError(err)
	is.Equal("hello world!", val)
}

// TestEvalOnSelectorSupportsSyntaxWithDifferentEngines verifies mixed xpath/css/text engines chained with >>.
// Ref: TestEvalOnSelector.java#shouldSupportSyntaxWithDifferentEngines
func TestEvalOnSelectorSupportsSyntaxWithDifferentEngines(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span>hello</span></div></section>`))

	val, err := page.EvalOnSelector(ctx, "xpath=/html/body/section >> css=div >> text='hello'", "(e, suffix) => e.textContent + suffix", " world!")
	must.NoError(err)
	is.Equal("hello world!", val)
}

// TestEvalOnSelectorSupportsSpacesWithSyntax verifies extra whitespace around >> and engine= prefixes is ignored.
// Ref: TestEvalOnSelector.java#shouldSupportSpacesWithSyntax
func TestEvalOnSelectorSupportsSpacesWithSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/deep-shadow.html", "text/html", deepShadowHTML)
	err := page.Goto(ctx, srv.Prefix()+"/deep-shadow.html")
	must.NoError(err)

	val, err := page.EvalOnSelector(ctx, " css = div >>css=div>>css   = span  ", "e => e.textContent")
	must.NoError(err)
	is.Equal("Hello from root2", val)
}

// TestEvalOnSelectorContinuesPastFirstFailure verifies >> chaining scans all siblings, not just the first.
// Ref: TestEvalOnSelector.java#shouldNotStopAtFirstFailureWithSyntax
func TestEvalOnSelectorContinuesPastFirstFailure(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>Next</span><button>Previous</button><button>Next</button></div>`))

	val, err := page.EvalOnSelector(ctx, "button >> 'Next'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<button>Next</button>", val)
}

// TestEvalOnSelectorSupportsCapture verifies the * capture modifier returns the captured element rather than the final match.
// Ref: TestEvalOnSelector.java#shouldSupportCapture
func TestEvalOnSelectorSupportsCapture(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span>a</span></div></section><section><div><span>b</span></div></section>`))

	val, err := page.EvalOnSelector(ctx, "*css=div >> 'b'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<div><span>b</span></div>", val)

	val, err = page.EvalOnSelector(ctx, "section >> *css=div >> 'b'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<div><span>b</span></div>", val)

	val, err = page.EvalOnSelector(ctx, "css=div >> *text='b'", "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span>b</span>", val)
}

// TestEvalOnSelectorThrowsOnMultipleCaptures verifies that using * on more than one part of a chain returns an error.
// Ref: TestEvalOnSelector.java#shouldThrowOnMultipleCaptures
func TestEvalOnSelectorThrowsOnMultipleCaptures(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span></span></div></section>`))

	_, err := page.EvalOnSelector(ctx, "*css=div >> *css=span", "e => e.outerHTML")
	is.Error(err)
	is.Contains(err.Error(), "Only one of the selectors can capture using * modifier")
}

// TestEvalOnSelectorThrowsOnMalformedCapture verifies that *=engine syntax without a valid engine name returns an error.
// Ref: TestEvalOnSelector.java#shouldThrowOnMalformedCapture
func TestEvalOnSelectorThrowsOnMalformedCapture(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div></div></section>`))

	_, err := page.EvalOnSelector(ctx, "*=div", "e => e.outerHTML")
	is.Error(err)
	is.Contains(err.Error(), "Unknown engine")
}

// TestEvalOnSelectorAllWithTextSelector verifies text= engine matches all elements with the given text.
// Ref: TestEvalOnSelectorAll.java#shouldWorkWithTextSelector
func TestEvalOnSelectorAllWithTextSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div><div>beautiful</div><div>beautiful</div><div>world!</div>`))

	val, err := page.EvalOnSelectorAll(ctx, "text='beautiful'", "divs => divs.length")
	must.NoError(err)
	is.Equal(float64(2), val)
}

// TestEvalOnSelectorAllWithXpathSelector verifies xpath= engine returns all matching elements.
// Ref: TestEvalOnSelectorAll.java#shouldWorkWithXpathSelector
func TestEvalOnSelectorAllWithXpathSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div><div>beautiful</div><div>world!</div>`))

	val, err := page.EvalOnSelectorAll(ctx, "xpath=/html/body/div", "divs => divs.length")
	must.NoError(err)
	is.Equal(float64(3), val)
}

// TestEvalOnSelectorAllSupportsSyntax verifies css=A >> css=B chaining collects all descendant matches.
// Ref: TestEvalOnSelectorAll.java#shouldSupportSyntax
func TestEvalOnSelectorAllSupportsSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>hello</span></div><div>beautiful</div><div><span>wo</span><span>rld!</span></div><span>Not this one</span>`))

	val, err := page.EvalOnSelectorAll(ctx, "css=div >> css=span", "spans => spans.length")
	must.NoError(err)
	is.Equal(float64(3), val)
}

// TestEvalOnSelectorAllSupportsCapture verifies the * capture modifier works with EvalOnSelectorAll.
// Ref: TestEvalOnSelectorAll.java#shouldSupportCapture
func TestEvalOnSelectorAllSupportsCapture(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div><span>a</span></div></section><section><div><span>b</span></div></section>`))

	val, err := page.EvalOnSelectorAll(ctx, "*css=div >> 'b'", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, "section >> *css=div >> 'b'", "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, "section >> *", "els => els.length")
	must.NoError(err)
	is.Equal(float64(4), val)
}
