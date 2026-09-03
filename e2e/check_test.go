//go:build e2e

// Check/Uncheck E2E tests.
// Migration of: TestCheck.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCheckCheckbox verifies Check() checks an unchecked checkbox.
// Ref: TestCheck.java#shouldCheckTheBox
func TestCheckCheckbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="checkbox" type="checkbox"></input>`))
	must.NoError(page.Locator("input").Check(ctx))

	checked, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	is.Equal(true, checked)
}

// TestCheckDoesNotRecheckCheckedBox verifies Check() on an already-checked box keeps it checked.
// Ref: TestCheck.java#shouldNotCheckTheCheckedBox
func TestCheckDoesNotRecheckCheckedBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="checkbox" type="checkbox" checked></input>`))
	must.NoError(page.Locator("input").Check(ctx))

	checked, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	is.Equal(true, checked)
}

// TestUncheckCheckbox verifies Uncheck() unchecks a checked checkbox.
// Ref: TestCheck.java#shouldUncheckTheBox
func TestUncheckCheckbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="checkbox" type="checkbox" checked></input>`))
	must.NoError(page.Locator("input").Uncheck(ctx))

	checked, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	is.Equal(false, checked)
}

// TestUncheckDoesNotUncheckUncheckedBox verifies Uncheck() on an unchecked box keeps it unchecked.
// Ref: TestCheck.java#shouldNotUncheckTheUncheckedBox
func TestUncheckDoesNotUncheckUncheckedBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="checkbox" type="checkbox"></input>`))
	must.NoError(page.Locator("input").Uncheck(ctx))

	checked, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	is.Equal(false, checked)
}

// TestCheckCheckboxByLabel verifies Check() via a label element.
// Ref: TestCheck.java#shouldCheckTheBoxByLabel
func TestCheckCheckboxByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<label for="checkbox"><input id="checkbox" type="checkbox"></input></label>`))
	must.NoError(page.Locator("label").Check(ctx))

	checked, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	is.Equal(true, checked)
}

// TestCheckCheckboxOutsideLabel verifies Check() via a label that references an input outside of it.
// Ref: TestCheck.java#shouldCheckTheBoxOutsideLabel
func TestCheckCheckboxOutsideLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<label for="checkbox">Text</label><div><input id="checkbox" type="checkbox"></input></div>`))
	must.NoError(page.Locator("label").Check(ctx))

	checked, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	is.Equal(true, checked)
}

// TestCheckCheckboxInsideLabelWithoutId verifies Check() via a label wrapping a span with the input.
// Ref: TestCheck.java#shouldCheckTheBoxInsideLabelWOId
func TestCheckCheckboxInsideLabelWithoutId(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<label>Text<span><input id="checkbox" type="checkbox"></input></span></label>`))
	must.NoError(page.Locator("label").Check(ctx))

	checked, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	is.Equal(true, checked)
}

// TestCheckRadioButton verifies Check() on a radio button.
// Ref: TestCheck.java#shouldCheckRadio
func TestCheckRadioButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="radio">one</input>
<input id="two" type="radio">two</input>
<input type="radio">three</input>`))
	must.NoError(page.Locator("#two").Check(ctx))

	checked, err := page.Evaluate(ctx, "() => window['two'].checked")
	must.NoError(err)
	is.Equal(true, checked)
}

// TestCheckCheckboxByAriaRole verifies Check() on a div with role=checkbox.
// Ref: TestCheck.java#shouldCheckTheBoxByAriaRole
func TestCheckCheckboxByAriaRole(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div role="checkbox" id="checkbox">CHECKBOX</div>
<script>
  checkbox.addEventListener('click', () => checkbox.setAttribute('aria-checked', 'true'));
</script>`))
	must.NoError(page.Locator("div").Check(ctx))

	ariaChecked, err := page.Evaluate(ctx, "() => window['checkbox'].getAttribute('aria-checked')")
	must.NoError(err)
	is.Equal("true", ariaChecked)
}

// TestLocatorCheckFiresChangeEvent verifies Check fires a change event.
// Ref: TestLocatorCheck.java#shouldFireChangeEvent
func TestLocatorCheckFiresChangeEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="cb" onchange="document.getElementById('result').textContent = this.checked ? 'on' : 'off'">
		<div id="result"></div>
	`))

	must.NoError(page.Locator("#cb").Check(ctx))

	text, err := page.Locator("#result").InnerText(ctx)
	must.NoError(err)
	is.Equal("on", text)
}

// TestLocatorUncheckFiresChangeEvent verifies Uncheck fires a change event.
// Ref: TestLocatorCheck.java#shouldFireChangeEventOnUncheck
func TestLocatorUncheckFiresChangeEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="cb" checked onchange="document.getElementById('result').textContent = this.checked ? 'on' : 'off'">
		<div id="result"></div>
	`))

	must.NoError(page.Locator("#cb").Uncheck(ctx))

	text, err := page.Locator("#result").InnerText(ctx)
	must.NoError(err)
	is.Equal("off", text)
}

// TestLocatorCheckIsCheckedAfterCheck verifies IsChecked is true after Check.
// Ref: TestLocatorCheck.java#shouldBeCheckedAfterCheck
func TestLocatorCheckIsCheckedAfterCheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))
	must.NoError(page.Locator("#cb").Check(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorUncheckIsUncheckedAfterUncheck verifies IsChecked is false after Uncheck.
// Ref: TestLocatorCheck.java#shouldNotBeCheckedAfterUncheck
func TestLocatorUncheckIsUncheckedAfterUncheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))
	must.NoError(page.Locator("#cb").Uncheck(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorCheckByLabel verifies Check can interact with checkbox via associated label.
// Ref: TestLocatorCheck.java#shouldCheckByLabel
func TestLocatorCheckByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="agree">I agree</label>
		<input type="checkbox" id="agree">
	`))

	must.NoError(page.Locator("label").Click(ctx))

	checked, err := page.Locator("#agree").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorCheckIsIdempotent verifies Check on already-checked doesn't error.
// Ref: TestLocatorCheck.java#shouldBeIdempotent
func TestLocatorCheckIsIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	must.NoError(page.Locator("#cb").Check(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorUncheckIsIdempotent verifies Uncheck on already-unchecked doesn't error.
// Ref: TestLocatorCheck.java#shouldBeIdempotentUncheck
func TestLocatorUncheckIsIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	must.NoError(page.Locator("#cb").Uncheck(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorCheckByLabelAssociation verifies Check works via label association.
// Ref: TestLocatorCheck.java#shouldCheckByLabel
func TestLocatorCheckByLabelAssociation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="cb">Accept terms</label>
		<input type="checkbox" id="cb">
	`))

	must.NoError(page.GetByLabel("Accept terms").Check(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorCheckFiresClickEvent verifies Check fires a click event.
// Ref: TestLocatorCheck.java#shouldFireClickEvent
func TestLocatorCheckFiresClickEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="cb" onclick="window.__clicked=true">
	`))

	must.NoError(page.Locator("#cb").Check(ctx))

	clicked, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Equal(true, clicked)
}

// TestLocatorUncheckFiresClickEvent verifies Uncheck fires a click event.
// Ref: TestLocatorCheck.java#shouldFireClickEventOnUncheck
func TestLocatorUncheckFiresClickEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="cb" checked onclick="window.__unClicked=true">
	`))

	must.NoError(page.Locator("#cb").Uncheck(ctx))

	clicked, err := page.Evaluate(ctx, `() => window.__unClicked`)
	must.NoError(err)
	is.Equal(true, clicked)
}

// TestLocatorCheckSetsCheckedTrue verifies Check sets checkbox to checked.
// Ref: TestLocatorCheck.java#shouldSetChecked
func TestLocatorCheckSetsCheckedTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	must.NoError(page.Locator("#cb").Check(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorUncheckSetsCheckedFalse verifies Uncheck sets checkbox to unchecked.
// Ref: TestLocatorCheck.java#shouldSetUnchecked
func TestLocatorUncheckSetsCheckedFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox" checked>`))

	must.NoError(page.Locator("#cb").Uncheck(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorCheckIdempotent verifies Check on already-checked checkbox does not error.
// Ref: TestLocatorCheck.java#shouldBeIdempotent
func TestLocatorCheckIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox" checked>`))

	must.NoError(page.Locator("#cb").Check(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorUncheckIdempotent verifies Uncheck on already-unchecked checkbox does not error.
// Ref: TestLocatorCheck.java#shouldBeIdempotentUncheck
func TestLocatorUncheckIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	must.NoError(page.Locator("#cb").Uncheck(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorCheckFiresChangeEvent verifies Check fires change event.
// Ref: TestLocatorCheck.java#shouldFireChangeEvent
func TestLocatorCheckFiresChangeEventEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="cb" type="checkbox" onchange="window.__changed=true">
	`))

	must.NoError(page.Locator("#cb").Check(ctx))

	changed, err := page.Evaluate(ctx, `() => window.__changed`)
	must.NoError(err)
	is.Equal(true, changed)
}

// TestLocatorCheckRadioButtonEx4 verifies Check works on radio button.
// Ref: TestLocatorCheck.java#shouldCheckRadio
func TestLocatorCheckRadioButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="radio" id="r1" name="group" value="a">
		<input type="radio" id="r2" name="group" value="b">
	`))

	must.NoError(page.Locator("#r1").Check(ctx))

	checked, err := page.Locator("#r1").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	checked2, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.False(checked2)
}

// TestLocatorCheckFiresInputEventEx4 verifies Check fires input event.
// Ref: TestLocatorCheck.java#shouldFireInputEvent
func TestLocatorCheckFiresInputEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="cb" oninput="window.__inputFired=true">
	`))

	must.NoError(page.Locator("#cb").Check(ctx))

	fired, err := page.Evaluate(ctx, `() => window.__inputFired`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestLocatorIsCheckedAfterClickEx4 verifies IsChecked returns true after Click on checkbox.
// Ref: TestLocatorCheck.java#shouldBeCheckedAfterClick
func TestLocatorIsCheckedAfterClickEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	must.NoError(page.Locator("#cb").Click(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedTrueEx4 verifies SetChecked(true) checks unchecked checkbox.
// Ref: TestLocatorCheck.java#shouldSetCheckedToTrue
func TestLocatorSetCheckedTrueEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, true))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedFalseEx4 verifies SetChecked(false) unchecks checked checkbox.
// Ref: TestLocatorCheck.java#shouldSetCheckedToFalse
func TestLocatorSetCheckedFalseEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, false))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestCheckByIDEx5 verifies Check by ID selector.
// Ref: TestLocatorCheck.java#shouldCheckByID
func TestCheckByIDEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="myCheck" type="checkbox">`))
	must.NoError(page.Locator("#myCheck").Check(ctx))

	checked, err := page.Locator("#myCheck").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestUncheckByIDEx5 verifies Uncheck by ID selector.
// Ref: TestLocatorCheck.java#shouldUncheckByID
func TestUncheckByIDEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="myCheck" type="checkbox" checked>`))
	must.NoError(page.Locator("#myCheck").Uncheck(ctx))

	checked, err := page.Locator("#myCheck").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestCheckFiresInputEventEx5 verifies Check fires input event.
// Ref: TestLocatorCheck.java#shouldFireInputEvent
func TestCheckFiresInputEventEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="chk" type="checkbox">
		<script>
			var fired = false;
			document.getElementById('chk').addEventListener('input', function() { fired = true; });
		</script>
	`))

	must.NoError(page.Locator("#chk").Check(ctx))

	result, err := page.Evaluate(ctx, `() => fired`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestCheckRadioButtonEx5 verifies Check works on radio button.
// Ref: TestLocatorCheck.java#shouldCheckRadioButton
func TestCheckRadioButtonEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="radio" id="r1" name="g" value="one">
		<input type="radio" id="r2" name="g" value="two">
	`))

	must.NoError(page.Locator("#r1").Check(ctx))

	checked, err := page.Locator("#r1").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	checked2, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.False(checked2)
}

// TestCheckTripleCheckEx6 verifies Check is idempotent after multiple calls.
// Ref: TestLocatorCheck.java#shouldBeIdempotent
func TestCheckTripleCheckEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))
	must.NoError(page.Locator("#chk").Check(ctx))
	must.NoError(page.Locator("#chk").Check(ctx))
	must.NoError(page.Locator("#chk").Check(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestUncheckTripleEx6 verifies Uncheck is idempotent after multiple calls.
// Ref: TestLocatorCheck.java#shouldBeIdempotentUncheck
func TestUncheckTripleEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))
	must.NoError(page.Locator("#chk").Uncheck(ctx))
	must.NoError(page.Locator("#chk").Uncheck(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestCheckStateByValueEx6 verifies checkbox state matches value attribute.
// Ref: TestLocatorCheck.java#shouldMatchValueAttribute
func TestCheckStateByValueEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" value="agree">`))
	must.NoError(page.Locator("#chk").Check(ctx))

	val, err := page.Evaluate(ctx, `() => document.getElementById('chk').value`)
	must.NoError(err)
	is.Equal("agree", val)
}

// TestCheckFocusedCheckboxEx6 verifies Check focuses and checks.
// Ref: TestLocatorCheck.java#shouldFocusAndCheck
func TestCheckFocusedCheckboxEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))
	must.NoError(page.Locator("#chk").Check(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("chk", focused)
}

// TestCheckTriggersChangeEventEx7 verifies Check fires change event.
// Ref: TestLocatorCheck.java#shouldTriggerChangeEvent
func TestCheckTriggersChangeEventEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="chk" type="checkbox" onchange="document.getElementById('out').textContent='changed'">
		<span id="out"></span>
	`))

	must.NoError(page.Locator("#chk").Check(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("changed", out)
}

// TestUncheckTriggersChangeEventEx7 verifies Uncheck fires change event.
// Ref: TestLocatorCheck.java#shouldTriggerChangeEventOnUncheck
func TestUncheckTriggersChangeEventEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="chk" type="checkbox" checked onchange="document.getElementById('out').textContent='unchecked'">
		<span id="out"></span>
	`))

	must.NoError(page.Locator("#chk").Uncheck(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("unchecked", out)
}

// TestCheckInGroupEx7 verifies Check works within a group of checkboxes.
// Ref: TestLocatorCheck.java#shouldCheckInGroup
func TestCheckInGroupEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="opt" type="checkbox" value="a">
		<input class="opt" type="checkbox" value="b">
		<input class="opt" type="checkbox" value="c">
	`))

	opts, err := page.Locator(".opt").All(ctx)
	must.NoError(err)

	must.NoError(opts[1].Check(ctx))

	val, err := page.Evaluate(ctx, `() => [...document.querySelectorAll('.opt:checked')].map(e => e.value)`)
	must.NoError(err)
	arr, ok := val.([]interface{})
	is.True(ok)
	is.Len(arr, 1)
	is.Equal("b", arr[0])
}

// TestCheckLabelAssociatedEx7 verifies Check works when clicking associated label.
// Ref: TestLocatorCheck.java#shouldCheckViaLabel
func TestCheckLabelAssociatedEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="chk" type="checkbox">
		<label for="chk">Accept terms</label>
	`))

	must.NoError(page.Locator("label[for='chk']").Click(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestIsCheckedInitiallyFalseEx2 verifies unchecked checkbox returns false.
// Ref: TestLocatorIsChecked.java#shouldReturnFalseForUnchecked
func TestIsCheckedInitiallyFalseEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestIsCheckedInitiallyTrueEx2 verifies checked checkbox returns true.
// Ref: TestLocatorIsChecked.java#shouldReturnTrueForChecked
func TestIsCheckedInitiallyTrueEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestIsCheckedAfterCheckEx2 verifies IsChecked returns true after Check.
// Ref: TestLocatorIsChecked.java#shouldReturnTrueAfterCheck
func TestIsCheckedAfterCheckEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))
	must.NoError(page.Locator("#chk").Check(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestIsCheckedAfterUncheckEx2 verifies IsChecked returns false after Uncheck.
// Ref: TestLocatorIsChecked.java#shouldReturnFalseAfterUncheck
func TestIsCheckedAfterUncheckEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))
	must.NoError(page.Locator("#chk").Uncheck(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestIsCheckedRadioEx2 verifies IsChecked works for radio buttons.
// Ref: TestLocatorIsChecked.java#shouldWorkForRadio
func TestIsCheckedRadioEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="radio" id="r1" name="g" value="a" checked>
		<input type="radio" id="r2" name="g" value="b">
	`))

	r1, err := page.Locator("#r1").IsChecked(ctx)
	must.NoError(err)
	is.True(r1)

	r2, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.False(r2)
}

// TestIsCheckedAfterJSEx3 verifies IsChecked reflects JS-toggled state.
// Ref: TestLocatorIsChecked.java#shouldReflectJSToggle
func TestIsCheckedAfterJSEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('chk').checked = true`)
	must.NoError(err)

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestIsCheckedRadioEx3 verifies IsChecked for radio buttons.
// Ref: TestLocatorIsChecked.java#shouldWorkForRadioButton
func TestIsCheckedRadioEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="r1" type="radio" name="grp" value="a" checked>
		<input id="r2" type="radio" name="grp" value="b">
	`))

	r1, err := page.Locator("#r1").IsChecked(ctx)
	must.NoError(err)
	is.True(r1)

	r2, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.False(r2)
}

// TestIsCheckedAfterUncheckEx3 verifies IsChecked returns false after Uncheck.
// Ref: TestLocatorIsChecked.java#shouldReturnFalseAfterUncheck
func TestIsCheckedAfterUncheckEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))
	must.NoError(page.Locator("#chk").Uncheck(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestIsCheckedGroupEx3 verifies IsChecked on multiple checkboxes in a group.
// Ref: TestLocatorIsChecked.java#shouldCheckGroup
func TestIsCheckedGroupEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="opt" type="checkbox" id="o1" checked>
		<input class="opt" type="checkbox" id="o2">
		<input class="opt" type="checkbox" id="o3" checked>
	`))

	count, err := page.Evaluate(ctx, `() => document.querySelectorAll('.opt:checked').length`)
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestLocatorSetCheckedToTrue verifies SetChecked(true) checks an unchecked checkbox.
// Ref: TestLocatorMisc.java#shouldSetCheckedToTrue
func TestLocatorSetCheckedToTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, true))
	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedToFalse verifies SetChecked(false) unchecks a checked checkbox.
// Ref: TestLocatorMisc.java#shouldSetCheckedToFalse
func TestLocatorSetCheckedToFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, false))
	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorSetCheckedTrueIsIdempotent verifies SetChecked(true) on already-checked is no-op.
// Ref: TestLocatorMisc.java#shouldSetCheckedTrueIdempotent
func TestLocatorSetCheckedTrueIsIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	// Should succeed without error
	must.NoError(page.Locator("#cb").SetChecked(ctx, true))
	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedFalseIsIdempotent verifies SetChecked(false) on unchecked is no-op.
// Ref: TestLocatorMisc.java#shouldSetCheckedFalseIdempotent
func TestLocatorSetCheckedFalseIsIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	// Should succeed without error
	must.NoError(page.Locator("#cb").SetChecked(ctx, false))
	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorSetCheckedOnRadio verifies SetChecked(true) works on radio buttons.
// Ref: TestLocatorMisc.java#shouldWorkOnRadioButton
func TestLocatorSetCheckedOnRadio(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="radio" name="r" id="r1" value="one">
		<input type="radio" name="r" id="r2" value="two">
	`))

	must.NoError(page.Locator("#r2").SetChecked(ctx, true))
	checked, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	// r1 should be unchecked since radio group only allows one
	r1Checked, err := page.Locator("#r1").IsChecked(ctx)
	must.NoError(err)
	is.False(r1Checked)
}

// TestLocatorSetCheckedToTrue verifies SetChecked(true) checks the checkbox.
// Ref: TestLocatorSetChecked.java#shouldCheckBox
func TestLocatorSetCheckedToTrueEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, true))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedToFalse verifies SetChecked(false) unchecks the checkbox.
// Ref: TestLocatorSetChecked.java#shouldUncheckBox
func TestLocatorSetCheckedToFalseEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, false))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorSetCheckedIdempotentTrue verifies SetChecked(true) is idempotent.
// Ref: TestLocatorSetChecked.java#shouldBeIdempotentForTrue
func TestLocatorSetCheckedIdempotentTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, true))
	must.NoError(page.Locator("#cb").SetChecked(ctx, true))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedFiresChangeEvent verifies SetChecked fires change event.
// Ref: TestLocatorSetChecked.java#shouldFireChangeEvent
func TestLocatorSetCheckedFiresChangeEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" id="cb"
		       onchange="window.__changed=(window.__changed||0)+1">
	`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, true))

	count, err := page.Evaluate(ctx, `() => window.__changed`)
	must.NoError(err)
	is.Equal(float64(1), count)
}

// TestLocatorSetCheckedWithRadioButton verifies SetChecked works on radio button.
// Ref: TestLocatorSetChecked.java#shouldWorkWithRadio
func TestLocatorSetCheckedWithRadioButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="radio" name="g" id="r1" value="a">
		<input type="radio" name="g" id="r2" value="b">
	`))

	must.NoError(page.Locator("#r2").SetChecked(ctx, true))

	checked, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	notChecked, err := page.Locator("#r1").IsChecked(ctx)
	must.NoError(err)
	is.False(notChecked)
}

// TestLocatorSetCheckedToTrueEx3 verifies SetChecked(true) checks checkbox.
// Ref: TestLocatorSetChecked.java#shouldCheckEx3
func TestLocatorSetCheckedToTrueEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, true))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedToFalseEx3 verifies SetChecked(false) unchecks checkbox.
// Ref: TestLocatorSetChecked.java#shouldUncheckEx3
func TestLocatorSetCheckedToFalseEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox" checked>`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, false))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorSetCheckedToTrueOnRadioEx3 verifies SetChecked(true) checks radio button.
// Ref: TestLocatorSetChecked.java#shouldCheckRadioEx3
func TestLocatorSetCheckedToTrueOnRadioEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="radio" id="r1" name="g">
		<input type="radio" id="r2" name="g">
	`))

	must.NoError(page.Locator("#r2").SetChecked(ctx, true))

	checked, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedIdempotentTrueEx3 verifies SetChecked(true) on already checked is idempotent.
// Ref: TestLocatorSetChecked.java#shouldBeIdempotentTrueEx3
func TestLocatorSetCheckedIdempotentTrueEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox" checked>`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, true))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorSetCheckedIdempotentFalseEx3 verifies SetChecked(false) on already unchecked is idempotent.
// Ref: TestLocatorSetChecked.java#shouldBeIdempotentFalseEx3
func TestLocatorSetCheckedIdempotentFalseEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	must.NoError(page.Locator("#cb").SetChecked(ctx, false))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestSetCheckedUnchecksWhenTrueEx4 verifies SetChecked(true) checks unchecked box.
// Ref: TestLocatorSetChecked.java#shouldCheckUnchecked
func TestSetCheckedUnchecksWhenTrueEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))
	must.NoError(page.Locator("#chk").SetChecked(ctx, true))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestSetCheckedFalseOnCheckedEx4 verifies SetChecked(false) unchecks checked box.
// Ref: TestLocatorSetChecked.java#shouldUncheckChecked
func TestSetCheckedFalseOnCheckedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))
	must.NoError(page.Locator("#chk").SetChecked(ctx, false))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestSetCheckedIdempotentTrueEx4 verifies SetChecked(true) on already-checked is no-op.
// Ref: TestLocatorSetChecked.java#shouldBeIdempotentTrue
func TestSetCheckedIdempotentTrueEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))
	must.NoError(page.Locator("#chk").SetChecked(ctx, true))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestSetCheckedOnRadioEx4 verifies SetChecked works on radio buttons.
// Ref: TestLocatorSetChecked.java#shouldWorkOnRadio
func TestSetCheckedOnRadioEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="r1" type="radio" name="group" value="a">
		<input id="r2" type="radio" name="group" value="b">
	`))

	must.NoError(page.Locator("#r2").SetChecked(ctx, true))

	checked, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestSetCheckedFiresChangeEventEx4 verifies SetChecked fires change event.
// Ref: TestLocatorSetChecked.java#shouldFireChangeEvent
func TestSetCheckedFiresChangeEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="chk" type="checkbox">
		<script>
			var changed = false;
			document.getElementById('chk').addEventListener('change', function() { changed = true; });
		</script>
	`))

	must.NoError(page.Locator("#chk").SetChecked(ctx, true))

	result, err := page.Evaluate(ctx, `() => changed`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestSetCheckedGroupEx5 verifies SetChecked on multiple checkboxes.
// Ref: TestLocatorSetChecked.java#shouldSetMultipleCheckboxes
func TestSetCheckedGroupEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="opt" type="checkbox">
		<input class="opt" type="checkbox">
		<input class="opt" type="checkbox">
	`))

	opts, err := page.Locator(".opt").All(ctx)
	must.NoError(err)
	for _, opt := range opts {
		must.NoError(opt.SetChecked(ctx, true))
	}

	count, err := page.Evaluate(ctx, `() => document.querySelectorAll('.opt:checked').length`)
	must.NoError(err)
	is.Equal(float64(3), count)
}

// TestSetCheckedFalseWhenAlreadyUncheckedEx5 verifies SetChecked(false) on unchecked is idempotent.
// Ref: TestLocatorSetChecked.java#shouldBeIdempotentFalse
func TestSetCheckedFalseWhenAlreadyUncheckedEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="c" type="checkbox">`))

	must.NoError(page.Locator("#c").SetChecked(ctx, false))

	checked, err := page.Locator("#c").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestSetCheckedTrueWhenAlreadyCheckedEx5 verifies SetChecked(true) on checked is idempotent.
// Ref: TestLocatorSetChecked.java#shouldBeIdempotentTrue
func TestSetCheckedTrueWhenAlreadyCheckedEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="c" type="checkbox" checked>`))

	must.NoError(page.Locator("#c").SetChecked(ctx, true))

	checked, err := page.Locator("#c").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestSetCheckedRadioButtonEx5 verifies SetChecked on radio buttons.
// Ref: TestLocatorSetChecked.java#shouldSetRadioButton
func TestSetCheckedRadioButtonEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="r1" type="radio" name="grp" value="a">
		<input id="r2" type="radio" name="grp" value="b">
	`))

	must.NoError(page.Locator("#r2").SetChecked(ctx, true))

	checked, err := page.Locator("#r2").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}
