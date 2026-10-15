//go:build e2e

// E2E tests for: TestPageSelectOption.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestElementHandleSelectOptionByValue verifies SelectOption by value on select element.
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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

// TestLocatorSelectOptionByLabel verifies SelectOption can select by display label.
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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

// TestPageSelectOptionShouldSelectSingleOptionByValue verifies SelectOption selects by value.
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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

// TestPageSelectOptionShouldSelectMultipleOptions verifies SelectOption with multiple select.
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSelectOption.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
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

// TestPageSelectOptionShouldSelectSingleOptionByValue verifies selecting a single option by value.
// Ref: TestPageSelectOption.java#shouldSelectSingleOptionByValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectSingleOptionByValue(t *testing.T) {
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

// TestPageSelectOptionShouldSelectOnlyFirstOption verifies that passing multiple values to a non-multiple
// select selects only the first matching one.
// Ref: TestPageSelectOption.java#shouldSelectOnlyFirstOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectOnlyFirstOption(t *testing.T) {
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

// TestPageSelectOptionShouldSelectMultipleOptions verifies selecting multiple options on a multi-select.
// Ref: TestPageSelectOption.java#shouldSelectMultipleOptions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectMultipleOptions(t *testing.T) {
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

// TestPageSelectOptionShouldRespectEventBubbling verifies that input/change events bubble up.
// Ref: TestPageSelectOption.java#shouldRespectEventBubbling  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldRespectEventBubbling(t *testing.T) {
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

// TestPageSelectOptionShouldThrowWhenElementIsNotASelect verifies that SelectOption on a non-select element errors.
// Ref: TestPageSelectOption.java#shouldThrowWhenElementIsNotASelect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldThrowWhenElementIsNotASelect(t *testing.T) {
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

// TestPageSelectOptionShouldReturnOnNoMatchedValues verifies that passing no values returns empty list.
// Ref: TestPageSelectOption.java#shouldReturnOnNoMatchedValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldReturnOnNoMatchedValues(t *testing.T) {
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

// TestPageSelectOptionShouldReturnAnArrayOfMatchedValues verifies return value is an array of matched values.
// Ref: TestPageSelectOption.java#shouldReturnAnArrayOfMatchedValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldReturnAnArrayOfMatchedValues(t *testing.T) {
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

// TestPageSelectOptionShouldReturnAnArrayOfOneElementWhenMultipleIsNotSet verifies only one option is returned for non-multi selects.
// Ref: TestPageSelectOption.java#shouldReturnAnArrayOfOneElementWhenMultipleIsNotSet  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldReturnAnArrayOfOneElementWhenMultipleIsNotSet(t *testing.T) {
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

// TestPageSelectOptionShouldDeselectAllOptionsWhenPassedNoValuesForAMultipleSelect verifies deselecting all options by passing empty args.
// Ref: TestPageSelectOption.java#shouldDeselectAllOptionsWhenPassedNoValuesForAMultipleSelect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldDeselectAllOptionsWhenPassedNoValuesForAMultipleSelect(t *testing.T) {
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

// TestPageSelectOptionShouldFallBackToSelectingByLabel verifies selecting by label when no option matches by value.
// Ref: TestPageSelectOption.java#shouldFallBackToSelectingByLabel  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldFallBackToSelectingByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Locator("select").SelectOption(ctx, "Blue")
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"blue"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"blue"}, onChange)
}

// TestPageSelectOptionShouldSelectSingleOptionByLabel verifies selecting option by label descriptor.
// Ref: TestPageSelectOption.java#shouldSelectSingleOptionByLabel  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectSingleOptionByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Locator("select").SelectOption(ctx, (&playwright.SelectOption{}).SetLabel("Indigo"))
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"indigo"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"indigo"}, onChange)
}

// TestPageSelectOptionShouldSelectSingleOptionByHandle verifies selecting option by ElementHandle.
// Ref: TestPageSelectOption.java#shouldSelectSingleOptionByHandle  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectSingleOptionByHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	handle, err := page.QuerySelector(ctx, "[id=whiteOption]")
	must.NoError(err)
	must.NotNil(handle)
	_, err = page.Locator("select").SelectOption(ctx, handle)
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"white"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"white"}, onChange)
}

// TestPageSelectOptionShouldSelectSingleOptionByIndex verifies selecting option by 0-based index descriptor.
// Ref: TestPageSelectOption.java#shouldSelectSingleOptionByIndex  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectSingleOptionByIndex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Locator("select").SelectOption(ctx, (&playwright.SelectOption{}).SetIndex(2))
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"brown"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"brown"}, onChange)
}

// TestPageSelectOptionShouldSelectSingleOptionByMultipleAttributes verifies selecting option matching value and label.
// Ref: TestPageSelectOption.java#shouldSelectSingleOptionByMultipleAttributes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectSingleOptionByMultipleAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Locator("select").SelectOption(ctx, (&playwright.SelectOption{}).SetValue("green").SetLabel("Green"))
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"green"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"green"}, onChange)
}

// TestPageSelectOptionShouldNotSelectSingleOptionWhenSomeAttributesDoNotMatch verifies error when attributes conflict.
// Ref: TestPageSelectOption.java#shouldNotSelectSingleOptionWhenSomeAttributesDoNotMatch  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldNotSelectSingleOptionWhenSomeAttributesDoNotMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.EvalOnSelector(ctx, "select", "s => s.value = undefined")
	must.NoError(err)
	_, err = page.Locator("select").SelectOption(ctx, (&playwright.SelectOption{}).SetValue("green").SetLabel("Brown"), (&playwright.SelectOptionOptions{}).SetTimeout(300))
	must.Error(err)
	is.Contains(err.Error(), "Timeout")
	val, err := page.Evaluate(ctx, "() => document.querySelector('select').value")
	must.NoError(err)
	is.Equal("", val)
}

// TestPageSelectOptionShouldNotThrowWhenSelectCausesNavigation verifies selecting an option that navigates does not fail.
// Ref: TestPageSelectOption.java#shouldNotThrowWhenSelectCausesNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldNotThrowWhenSelectCausesNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.EvalOnSelector(ctx, "select", "select => select.addEventListener('input', () => window.location.href = '/empty.html')")
	must.NoError(err)
	_, err = page.Locator("select").SelectOption(ctx, "blue")
	must.NoError(err)
	must.NoError(page.WaitForURL(ctx, "**/empty.html"))
	is.Contains(page.URL(), "empty.html")
}

// TestPageSelectOptionShouldNotAllowNullItems is disabled in Java (// @Test) and null items in an array cannot be expressed with Go typed arguments.
// Ref: TestPageSelectOption.java#shouldNotAllowNullItems  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldNotAllowNullItems(t *testing.T) {
	t.Skip("disabled in playwright-java (// @Test) and null items in an array cannot be expressed with Go typed arguments")
}

// TestPageSelectOptionShouldUnselectWithNull verifies selecting no values clears a multiple select.
// Ref: TestPageSelectOption.java#shouldUnselectWithNull  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldUnselectWithNull(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Evaluate(ctx, "() => window['makeMultiple']()")
	must.NoError(err)
	result, err := page.Locator("select").SelectOption(ctx, "blue", "black", "magenta")
	must.NoError(err)
	is.ElementsMatch([]string{"blue", "black", "magenta"}, result)
	_, err = page.Locator("select").SelectOption(ctx)
	must.NoError(err)
	all, err := page.EvalOnSelector(ctx, "select", "select => Array.from(select.options).every(option => !option.selected)")
	must.NoError(err)
	is.Equal(true, all)
}

// TestPageSelectOptionShouldDeselectAllOptionsWhenPassedNoValuesForASelectWithoutMultiple verifies selecting no values deselects all options.
// Ref: TestPageSelectOption.java#shouldDeselectAllOptionsWhenPassedNoValuesForASelectWithoutMultiple  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldDeselectAllOptionsWhenPassedNoValuesForASelectWithoutMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Locator("select").SelectOption(ctx, "blue", "black", "magenta")
	must.NoError(err)
	_, err = page.Locator("select").SelectOption(ctx)
	must.NoError(err)
	all, err := page.EvalOnSelector(ctx, "select", "select => Array.from(select.options).every(option => !option.selected)")
	must.NoError(err)
	is.Equal(true, all)
}

// TestPageSelectOptionShouldWorkWhenReDefiningTopLevelEventClass verifies selectOption works when window.Event is redefined.
// Ref: TestPageSelectOption.java#shouldWorkWhenReDefiningTopLevelEventClass  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldWorkWhenReDefiningTopLevelEventClass(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Evaluate(ctx, "() => window.Event = null")
	must.NoError(err)
	_, err = page.Locator("select").SelectOption(ctx, "blue")
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"blue"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"blue"}, onChange)
}

// TestPageSelectOptionShouldWaitForOptionToBeEnabled verifies selectOption waits until the select is enabled.
// Ref: TestPageSelectOption.java#shouldWaitForOptionToBeEnabled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldWaitForOptionToBeEnabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<select disabled>
  <option>one</option>
  <option>two</option>
</select>
<script>
function hydrate() {
  const select = document.querySelector('select');
  select.removeAttribute('disabled');
  select.addEventListener('change', () => {
    window['result'] = select.value;
  });
}
</script>`))
	_, err := page.Evaluate(ctx, "() => setTimeout(hydrate, 1000)")
	must.NoError(err)
	_, err = page.Locator("select").SelectOption(ctx, "two")
	must.NoError(err)
	v, err := page.Evaluate(ctx, "window['result']")
	must.NoError(err)
	is.Equal("two", v)
	val, err := page.Locator("select").InputValue(ctx)
	must.NoError(err)
	is.Equal("two", val)
}

// TestPageSelectOptionShouldWaitForSelectToBeSwapped verifies selectOption retries when the select element is replaced.
// Ref: TestPageSelectOption.java#shouldWaitForSelectToBeSwapped  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldWaitForSelectToBeSwapped(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<select disabled>
  <option>one</option>
  <option>two</option>
</select>
<script>
function hydrate() {
  const select = document.querySelector('select');
  select.remove();

  const newSelect = document.createElement('select');
  const option1 = document.createElement('option');
  option1.textContent = 'one';
  newSelect.appendChild(option1);
  const option2 = document.createElement('option');
  option2.textContent = 'two';
  newSelect.appendChild(option2);

  document.body.appendChild(newSelect);

  newSelect.addEventListener('change', () => {
    window['result'] = newSelect.value;
  });
}
</script>`))
	_, err := page.Evaluate(ctx, "() => setTimeout(window.hydrate, 1000)")
	must.NoError(err)
	_, err = page.Locator("select").SelectOption(ctx, "two")
	must.NoError(err)
	val, err := page.Locator("select").InputValue(ctx)
	must.NoError(err)
	is.Equal("two", val)
	v, err := page.Evaluate(ctx, "window['result']")
	must.NoError(err)
	is.Equal("two", v)
}

// TestPageSelectOptionShouldSelectSingleOption verifies selecting an option by value.
// Ref: TestPageSelectOption.java#shouldSelectSingleOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectSingleOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Locator("select").SelectOption(ctx, "blue")
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"blue"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"blue"}, onChange)
}

// TestPageSelectOptionShouldSelectMultipleOptionsWithAttributes verifies selecting multiple options with attributes.
// Ref: TestPageSelectOption.java#shouldSelectMultipleOptionsWithAttributes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldSelectMultipleOptionsWithAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	_, err := page.Evaluate(ctx, "() => window['makeMultiple']()")
	must.NoError(err)
	_, err = page.Locator("select").SelectOption(ctx, []playwright.SelectOption{
		*(&playwright.SelectOption{}).SetValue("blue"),
		*(&playwright.SelectOption{}).SetLabel("Green"),
		*(&playwright.SelectOption{}).SetIndex(4),
	})
	must.NoError(err)
	onInput, err := page.Evaluate(ctx, "() => window['result'].onInput")
	must.NoError(err)
	is.Equal([]any{"blue", "gray", "green"}, onInput)
	onChange, err := page.Evaluate(ctx, "() => window['result'].onChange")
	must.NoError(err)
	is.Equal([]any{"blue", "gray", "green"}, onChange)
}

// TestPageSelectOptionShouldReturnOnNoValues verifies selecting no values returns an empty list.
// Ref: TestPageSelectOption.java#shouldReturnOnNoValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSelectOption.java
func TestPageSelectOptionShouldReturnOnNoValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	a4ServeSelect(srv)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/select.html"))
	result, err := page.Locator("select").SelectOption(ctx)
	must.NoError(err)
	must.Empty(result)
}
