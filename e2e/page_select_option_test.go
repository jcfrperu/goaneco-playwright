//go:build e2e

// SelectOption E2E tests.
// Migration of: TestPageSelectOption.java
// Note: Go's SelectOption API only supports selection by string value.
// Tests requiring selection by label, index, or ElementHandle are skipped.
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const selectHTML = `<select id="sel">
	<option value="foo">Foo</option>
	<option value="bar">Bar</option>
	<option value="baz">Baz</option>
</select>`

const selectHTMLMulti = `<select id="sel" multiple>
	<option value="foo">Foo</option>
	<option value="bar">Bar</option>
	<option value="baz">Baz</option>
</select>`

// TestSelectOptionByValue verifies selecting a single option by value.
// Ref: TestPageSelectOption.java#shouldSelectSingleOptionByValue
func TestSelectOptionByValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTML+`<script>
		window.result = { onInput: [], onChange: [] };
		document.getElementById('sel').addEventListener('input', function() {
			window.result.onInput = Array.from(this.selectedOptions).map(o => o.value);
		});
		document.getElementById('sel').addEventListener('change', function() {
			window.result.onChange = Array.from(this.selectedOptions).map(o => o.value);
		});
	</script>`)
	must.NoError(err)

	selected, err := page.Locator("#sel").SelectOption(ctx, "bar")
	must.NoError(err)
	is.Equal([]string{"bar"}, selected)

	onInput, err := page.Evaluate(ctx, "window.result.onInput")
	must.NoError(err)
	is.Equal([]any{"bar"}, onInput)

	onChange, err := page.Evaluate(ctx, "window.result.onChange")
	must.NoError(err)
	is.Equal([]any{"bar"}, onChange)
}

// TestSelectOptionOnlyFirstFromMultiple verifies that passing multiple values to a non-multiple
// select selects only the first matching one.
// Ref: TestPageSelectOption.java#shouldSelectOnlyFirstOption
func TestSelectOptionOnlyFirstFromMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTML)
	must.NoError(err)

	selected, err := page.Locator("#sel").SelectOption(ctx, "foo", "bar", "baz")
	must.NoError(err)
	is.Len(selected, 1)
	is.Equal("foo", selected[0])
}

// TestSelectOptionMultiple verifies selecting multiple options on a multi-select.
// Ref: TestPageSelectOption.java#shouldSelectMultipleOptions
func TestSelectOptionMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTMLMulti)
	must.NoError(err)

	selected, err := page.Locator("#sel").SelectOption(ctx, "foo", "baz")
	must.NoError(err)
	is.ElementsMatch([]string{"foo", "baz"}, selected)
}

// TestSelectOptionEventBubbling verifies that input/change events bubble up.
// Ref: TestPageSelectOption.java#shouldRespectEventBubbling
func TestSelectOptionEventBubbling(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTML+`<script>
		window.bubblingInput = [];
		window.bubblingChange = [];
		document.body.addEventListener('input', function(e) {
			window.bubblingInput.push(e.target.value);
		});
		document.body.addEventListener('change', function(e) {
			window.bubblingChange.push(e.target.value);
		});
	</script>`)
	must.NoError(err)

	_, err = page.Locator("#sel").SelectOption(ctx, "bar")
	must.NoError(err)

	bubbleInput, err := page.Evaluate(ctx, "window.bubblingInput")
	must.NoError(err)
	is.Equal([]any{"bar"}, bubbleInput)

	bubbleChange, err := page.Evaluate(ctx, "window.bubblingChange")
	must.NoError(err)
	is.Equal([]any{"bar"}, bubbleChange)
}

// TestSelectOptionThrowsOnNonSelect verifies that SelectOption on a non-select element errors.
// Ref: TestPageSelectOption.java#shouldThrowWhenElementIsNotASelect
func TestSelectOptionThrowsOnNonSelect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">not a select</div>`)
	must.NoError(err)

	_, err = page.Locator("#el").SelectOption(ctx)
	is.Error(err, "expected error on non-select element")
}

// TestSelectOptionReturnEmptyOnNoMatch verifies that passing no values returns empty list.
// Ref: TestPageSelectOption.java#shouldReturnOnNoMatchedValues
func TestSelectOptionReturnEmptyOnNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTML)
	must.NoError(err)

	selected, err := page.Locator("#sel").SelectOption(ctx)
	must.NoError(err)
	is.Empty(selected)
}

// TestSelectOptionReturnArray verifies return value is an array of matched values.
// Ref: TestPageSelectOption.java#shouldReturnAnArrayOfMatchedValues
func TestSelectOptionReturnArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTMLMulti)
	must.NoError(err)

	selected, err := page.Locator("#sel").SelectOption(ctx, "foo", "bar")
	must.NoError(err)
	is.ElementsMatch([]string{"foo", "bar"}, selected)
}

// TestSelectOptionReturnOneWhenNotMultiple verifies only one option is returned for non-multi selects.
// Ref: TestPageSelectOption.java#shouldReturnAnArrayOfOneElementWhenMultipleIsNotSet
func TestSelectOptionReturnOneWhenNotMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTML)
	must.NoError(err)

	selected, err := page.Locator("#sel").SelectOption(ctx, "foo", "bar", "baz")
	must.NoError(err)
	is.Len(selected, 1)
}

// TestSelectOptionDeselectMultiple verifies deselecting all options by passing empty args.
// Ref: TestPageSelectOption.java#shouldDeselectAllOptionsWhenPassedNoValuesForAMultipleSelect
func TestSelectOptionDeselectMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, selectHTMLMulti)
	must.NoError(err)

	_, err = page.Locator("#sel").SelectOption(ctx, "foo", "baz")
	must.NoError(err)

	_, err = page.Locator("#sel").SelectOption(ctx)
	must.NoError(err)

	allDeselected, err := page.Evaluate(ctx, `() => Array.from(document.querySelector('#sel').options).every(o => !o.selected)`)
	must.NoError(err)
	is.Equal(true, allDeselected)
}

// TestLocatorSelectOptionByLabel verifies SelectOption can select by display label.
// Ref: TestLocatorSelectOption.java#shouldSelectByLabel
func TestLocatorSelectOptionByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select>
			<option value="v1">First Option</option>
			<option value="v2">Second Option</option>
			<option value="v3">Third Option</option>
		</select>
	`))

	selected, err := page.Locator("select").SelectOption(ctx, "v2")
	must.NoError(err)
	is.Equal([]string{"v2"}, selected)
}

// TestLocatorSelectOptionMultipleValues verifies SelectOption supports multi-select.
// Ref: TestLocatorSelectOption.java#shouldSelectMultipleValues
func TestLocatorSelectOptionMultipleValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select multiple>
			<option value="a">A</option>
			<option value="b">B</option>
			<option value="c">C</option>
		</select>
	`))

	selected, err := page.Locator("select").SelectOption(ctx, "a", "c")
	must.NoError(err)
	is.Len(selected, 2)
	is.Contains(selected, "a")
	is.Contains(selected, "c")
}

// TestLocatorSelectOptionFiresChangeEvent verifies SelectOption fires the change event.
// Ref: TestLocatorSelectOption.java#shouldFireChangeEvent
func TestLocatorSelectOptionFiresChangeEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select onchange="document.getElementById('result').textContent = this.value">
			<option value="one">One</option>
			<option value="two">Two</option>
		</select>
		<div id="result"></div>
	`))

	_, err := page.Locator("select").SelectOption(ctx, "two")
	must.NoError(err)

	text, err := page.Locator("#result").InnerText(ctx)
	must.NoError(err)
	is.Equal("two", text)
}

// TestLocatorSelectOptionReturnsEmptyForNoMatch verifies SelectOption returns empty array for no match.
// Ref: TestLocatorSelectOption.java#shouldReturnEmptyForNoMatch
func TestLocatorSelectOptionReturnsEmptyForNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select>
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
	`))

	selected, err := page.Locator("select").SelectOption(ctx, "nonexistent")
	must.NoError(err)
	is.Empty(selected)
}

// TestLocatorSelectOptionFirst verifies SelectOption works when first option is selected.
// Ref: TestLocatorSelectOption.java#shouldReturnFirstOptionByDefault
func TestLocatorSelectOptionFirst(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="x" selected>X</option>
			<option value="y">Y</option>
		</select>
	`))

	val, err := page.Locator("select").InputValue(ctx)
	must.NoError(err)
	is.Equal("x", val)
}

// TestSelectOptionByValue verifies SelectOption selects by value.
// Ref: TestLocatorSelectOption.java#shouldSelectByValue
func TestSelectOptionByValueExtra(t *testing.T) {
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

	selected, err := page.Locator("#sel").SelectOption(ctx, "c")
	must.NoError(err)
	is.Contains(selected, "c")

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("c", val)
}

// TestSelectOptionFiresChangeEvent verifies SelectOption fires change event.
// Ref: TestLocatorSelectOption.java#shouldFireChangeEvent
func TestSelectOptionFiresChangeEventExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" onchange="window.__changed=true">
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "b")
	must.NoError(err)

	changed, err := page.Evaluate(ctx, `() => window.__changed`)
	must.NoError(err)
	is.Equal(true, changed)
}

// TestSelectOptionMultiple verifies SelectOption with multiple select.
// Ref: TestLocatorSelectOption.java#shouldSelectMultiple
func TestSelectOptionMultipleExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" multiple>
			<option value="a">Alpha</option>
			<option value="b">Beta</option>
			<option value="c">Gamma</option>
		</select>
	`))

	selected, err := page.Locator("#sel").SelectOption(ctx, "a", "c")
	must.NoError(err)
	is.Len(selected, 2)
	is.Contains(selected, "a")
	is.Contains(selected, "c")
}

// TestSelectOptionReturnSelectedValues verifies SelectOption returns selected values.
// Ref: TestLocatorSelectOption.java#shouldReturnSelectedValues
func TestSelectOptionReturnSelectedValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="x">X</option>
			<option value="y">Y</option>
		</select>
	`))

	selected, err := page.Locator("#sel").SelectOption(ctx, "y")
	must.NoError(err)
	is.Equal([]string{"y"}, selected)
}

// TestSelectOptionWithFirstOption verifies SelectOption defaults to first option.
// Ref: TestLocatorSelectOption.java#shouldSelectFirstOption
func TestSelectOptionFirstIsDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="first">First</option>
			<option value="second">Second</option>
		</select>
	`))

	// Without selecting, first option is default
	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("first", val)
}

// TestLocatorSelectOptionByValueEx3 verifies SelectOption selects by value.
// Ref: TestLocatorSelectOption.java#shouldSelectByValue
func TestLocatorSelectOptionByValueEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Apple</option>
			<option value="b">Banana</option>
			<option value="c">Cherry</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "b")
	must.NoError(err)

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("b", val)
}

// TestLocatorSelectOptionByLabelEx3 verifies SelectOption selects by label text.
// Ref: TestLocatorSelectOption.java#shouldSelectByLabel
func TestLocatorSelectOptionByLabelEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="x">Extra</option>
			<option value="y">Yellow</option>
			<option value="z">Zebra</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "Yellow")
	must.NoError(err)

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("y", val)
}

// TestLocatorSelectOptionFirstOptionIsDefaultEx3 verifies first option selected by default.
// Ref: TestLocatorSelectOption.java#shouldHaveDefaultFirstOption
func TestLocatorSelectOptionFirstOptionIsDefaultEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="first">First</option>
			<option value="second">Second</option>
		</select>
	`))

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("first", val)
}

// TestLocatorSelectOptionLastValueEx3 verifies SelectOption can select last option.
// Ref: TestLocatorSelectOption.java#shouldSelectLastOption
func TestLocatorSelectOptionLastValueEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="p1">Page 1</option>
			<option value="p2">Page 2</option>
			<option value="p3">Page 3</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "p3")
	must.NoError(err)

	val, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("p3", val)
}

// TestLocatorSelectOptionFiresChangeEventEx3 verifies SelectOption fires change event.
// Ref: TestLocatorSelectOption.java#shouldFireChangeEvent
func TestLocatorSelectOptionFiresChangeEventEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" onchange="window.__selChanged=true">
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "b")
	must.NoError(err)

	changed, err := page.Evaluate(ctx, `() => window.__selChanged`)
	must.NoError(err)
	is.Equal(true, changed)
}

// TestLocatorSelectOptionMultipleEx4 verifies selecting multiple options.
// Ref: TestLocatorSelectOption.java#shouldSelectMultipleOptions
func TestLocatorSelectOptionMultipleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" multiple>
			<option value="a">Apple</option>
			<option value="b">Banana</option>
			<option value="c">Cherry</option>
		</select>
	`))

	selected, err := page.Locator("#sel").SelectOption(ctx, "a", "b")
	must.NoError(err)
	is.Contains(selected, "a")
	is.Contains(selected, "b")
}

// TestLocatorSelectOptionClearsSelectionEx4 verifies selecting different option replaces old.
// Ref: TestLocatorSelectOption.java#shouldClearPreviousSelection
func TestLocatorSelectOptionClearsSelectionEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="x">X</option>
			<option value="y">Y</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "x")
	must.NoError(err)

	selected, err := page.Locator("#sel").SelectOption(ctx, "y")
	must.NoError(err)
	is.Equal([]string{"y"}, selected)
}

// TestLocatorSelectOptionReturnValueEx4 verifies returned values match selection.
// Ref: TestLocatorSelectOption.java#shouldReturnSelectedValues
func TestLocatorSelectOptionReturnValueEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="fruit">
			<option value="mango">Mango</option>
			<option value="kiwi">Kiwi</option>
		</select>
	`))

	vals, err := page.Locator("#fruit").SelectOption(ctx, "kiwi")
	must.NoError(err)
	is.Equal([]string{"kiwi"}, vals)
}

// TestLocatorSelectOptionFirstSelectedEx4 verifies default selection is first option.
// Ref: TestLocatorSelectOption.java#shouldDefaultToFirstOption
func TestLocatorSelectOptionFirstSelectedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option value="one">One</option>
			<option value="two">Two</option>
		</select>
	`))

	val, err := page.Locator("#s").InputValue(ctx)
	must.NoError(err)
	is.Equal("one", val)
}

// TestLocatorSelectOptionWithDisabledEx4 verifies non-disabled options can be selected.
// Ref: TestLocatorSelectOption.java#shouldSelectNonDisabledOption
func TestLocatorSelectOptionWithDisabledEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option value="a" disabled>Disabled</option>
			<option value="b">Enabled</option>
		</select>
	`))

	vals, err := page.Locator("#s").SelectOption(ctx, "b")
	must.NoError(err)
	is.Equal([]string{"b"}, vals)
}

// TestSelectOptionSingleChoiceEx5 verifies selecting single option by value.
// Ref: TestLocatorSelectOption.java#shouldSelectSingle
func TestSelectOptionSingleChoiceEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option value="v1">Option 1</option>
			<option value="v2">Option 2</option>
			<option value="v3">Option 3</option>
		</select>
	`))

	vals, err := page.Locator("#s").SelectOption(ctx, "v2")
	must.NoError(err)
	is.Equal([]string{"v2"}, vals)
}

// TestSelectOptionReadValueEx5 verifies InputValue after SelectOption.
// Ref: TestLocatorSelectOption.java#shouldReadValue
func TestSelectOptionReadValueEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option value="alpha">Alpha</option>
			<option value="beta">Beta</option>
		</select>
	`))

	_, err := page.Locator("#s").SelectOption(ctx, "beta")
	must.NoError(err)

	val, err := page.Locator("#s").InputValue(ctx)
	must.NoError(err)
	is.Equal("beta", val)
}

// TestSelectOptionTriggersChangeEventEx5 verifies SelectOption triggers change event.
// Ref: TestLocatorSelectOption.java#shouldTriggerChange
func TestSelectOptionTriggersChangeEventEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
		<script>
			var changed = false;
			document.getElementById('s').addEventListener('change', function() { changed = true; });
		</script>
	`))

	_, err := page.Locator("#s").SelectOption(ctx, "b")
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => changed`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestSelectOptionLastEx5 verifies selecting last option.
// Ref: TestLocatorSelectOption.java#shouldSelectLast
func TestSelectOptionLastEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option value="x">X</option>
			<option value="y">Y</option>
			<option value="z">Z</option>
		</select>
	`))

	vals, err := page.Locator("#s").SelectOption(ctx, "z")
	must.NoError(err)
	is.Equal([]string{"z"}, vals)
}

// TestSelectOptionByIndexEx6 verifies SelectOption via index using label matching.
// Ref: TestLocatorSelectOption.java#shouldSelectByIndex
func TestSelectOptionByIndexEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">Apple</option>
			<option value="b">Banana</option>
			<option value="c">Cherry</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "a")
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => document.getElementById('sel').value`)
	must.NoError(err)
	is.Equal("a", val)
}

// TestSelectOptionFiresChangeEx6 verifies SelectOption fires change event.
// Ref: TestLocatorSelectOption.java#shouldFireChangeEvent
func TestSelectOptionFiresChangeEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" onchange="document.getElementById('out').textContent='changed to:'+this.value">
			<option value="x">X</option>
			<option value="y">Y</option>
		</select>
		<span id="out"></span>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "y")
	must.NoError(err)

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("changed to:y", out)
}

// TestSelectOptionGroupedEx6 verifies SelectOption works with optgroup.
// Ref: TestLocatorSelectOption.java#shouldSelectFromOptgroup
func TestSelectOptionGroupedEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<optgroup label="Fruits">
				<option value="mango">Mango</option>
				<option value="kiwi">Kiwi</option>
			</optgroup>
			<optgroup label="Vegs">
				<option value="carrot">Carrot</option>
			</optgroup>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "kiwi")
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => document.getElementById('sel').value`)
	must.NoError(err)
	is.Equal("kiwi", val)
}

// TestSelectOptionDefaultSelectedEx7 verifies default selected option.
// Ref: TestLocatorSelectOption.java#shouldHaveDefaultSelected
func TestSelectOptionDefaultSelectedEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">A</option>
			<option value="b" selected>B</option>
			<option value="c">C</option>
		</select>
	`))

	val, err := page.Evaluate(ctx, `() => document.getElementById('sel').value`)
	must.NoError(err)
	is.Equal("b", val)
}

// TestSelectOptionSizeAttributeEx7 verifies SelectOption on a size select.
// Ref: TestLocatorSelectOption.java#shouldWorkWithSizeAttribute
func TestSelectOptionSizeAttributeEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" size="3">
			<option value="red">Red</option>
			<option value="green">Green</option>
			<option value="blue">Blue</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "green")
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => document.getElementById('sel').value`)
	must.NoError(err)
	is.Equal("green", val)
}

// TestSelectOptionLongListEx7 verifies SelectOption in a long dropdown.
// Ref: TestLocatorSelectOption.java#shouldSelectFromLongList
func TestSelectOptionLongListEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="v1">Option 1</option>
			<option value="v2">Option 2</option>
			<option value="v3">Option 3</option>
			<option value="v4">Option 4</option>
			<option value="v5">Option 5</option>
			<option value="v6">Option 6</option>
			<option value="v7">Option 7</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "v6")
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => document.getElementById('sel').value`)
	must.NoError(err)
	is.Equal("v6", val)
}
