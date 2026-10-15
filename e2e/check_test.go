//go:build e2e

// E2E tests for: TestCheck.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestLocatorCheckFiresChangeEvent verifies Check fires a change event.
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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
// Extra: no equivalent method in playwright-java; related class: TestCheck.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
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

// TestCheckShouldCheckTheBox verifies Check() checks an unchecked checkbox.
// Ref: TestCheck.java#shouldCheckTheBox  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldCheckTheBox(t *testing.T) {
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

// TestCheckShouldNotCheckTheCheckedBox verifies Check() on an already-checked box keeps it checked.
// Ref: TestCheck.java#shouldNotCheckTheCheckedBox  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldNotCheckTheCheckedBox(t *testing.T) {
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

// TestCheckShouldUncheckTheBox verifies Uncheck() unchecks a checked checkbox.
// Ref: TestCheck.java#shouldUncheckTheBox  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldUncheckTheBox(t *testing.T) {
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

// TestCheckShouldNotUncheckTheUncheckedBox verifies Uncheck() on an unchecked box keeps it unchecked.
// Ref: TestCheck.java#shouldNotUncheckTheUncheckedBox  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldNotUncheckTheUncheckedBox(t *testing.T) {
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

// TestCheckShouldCheckTheBoxByLabel verifies Check() via a label element.
// Ref: TestCheck.java#shouldCheckTheBoxByLabel  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldCheckTheBoxByLabel(t *testing.T) {
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

// TestCheckShouldCheckTheBoxOutsideLabel verifies Check() via a label that references an input outside of it.
// Ref: TestCheck.java#shouldCheckTheBoxOutsideLabel  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldCheckTheBoxOutsideLabel(t *testing.T) {
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

// TestCheckShouldCheckTheBoxInsideLabelWOId verifies Check() via a label wrapping a span with the input.
// Ref: TestCheck.java#shouldCheckTheBoxInsideLabelWOId  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldCheckTheBoxInsideLabelWOId(t *testing.T) {
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

// TestCheckShouldCheckRadio verifies Check() on a radio button.
// Ref: TestCheck.java#shouldCheckRadio  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldCheckRadio(t *testing.T) {
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

// TestCheckShouldCheckTheBoxByAriaRole verifies Check() on a div with role=checkbox.
// Ref: TestCheck.java#shouldCheckTheBoxByAriaRole  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestCheck.java
func TestCheckShouldCheckTheBoxByAriaRole(t *testing.T) {
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
