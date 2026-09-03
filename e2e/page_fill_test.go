//go:build e2e

// E2E tests for Locator.Fill and Locator.Press (Priority 2).
// Ref: TestPageFill.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageFillTextarea verifies that Locator.Fill() sets the value of a textarea.
// Ref: TestPageFill.java#shouldFillTextarea
func TestPageFillTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<textarea id="ta"></textarea>`)
	must.NoError(err, "SetContent failed")

	ta := page.Locator("#ta")
	err = ta.Fill(ctx, "hello playwright")
	must.NoError(err, "Fill failed")

	val, err := ta.InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("hello playwright", val)
}

// TestPageFillInput verifies that Locator.Fill() sets the value of an input.
// Ref: TestPageFill.java#shouldFillInput
func TestPageFillInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" />`))
	inp := page.Locator("#inp")
	must.NoError(inp.Fill(ctx, "some value"))
	val, err := inp.InputValue(ctx)
	must.NoError(err)
	is.Equal("some value", val)
}

// TestPageFillThrowOnUnsupportedInputs verifies that filling unsupported input types returns an error.
// Ref: TestPageFill.java#shouldThrowOnUnsupportedInputs
func TestPageFillThrowOnUnsupportedInputs(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	for _, typ := range []string{"button", "checkbox", "file", "image", "radio", "reset", "submit"} {
		must.NoError(page.SetContent(ctx, `<input type="`+typ+`" id="inp" />`))
		inp := page.Locator("#inp")
		err := inp.Fill(ctx, "")
		must.Errorf(err, "expected error for input type %q", typ)
	}
}

// TestPageFillRangeInput verifies that a range input can be filled with a numeric value.
// Ref: TestPageFill.java#shouldFillRangeInput
func TestPageFillRangeInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=range min=0 max=100 value=50>`))
	must.NoError(page.Locator("input").Fill(ctx, "42"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("42", val)
}

// TestPageFillThrowOnIncorrectRangeValue verifies that filling a range with out-of-range/malformed values errors.
// Ref: TestPageFill.java#shouldThrowOnIncorrectRangeValue
func TestPageFillThrowOnIncorrectRangeValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=range min=0 max=100 value=50>`))
	inp := page.Locator("input")

	for _, bad := range []string{"foo", "200", "15.43"} {
		err := inp.Fill(ctx, bad)
		must.Errorf(err, "expected error for range value %q", bad)
		is.ErrorContains(err, "Malformed value")
	}
}

// TestPageFillDifferentInputTypes verifies that text-like input types (password, search, etc.) can be filled.
// Ref: TestPageFill.java#shouldFillDifferentInputTypes
func TestPageFillDifferentInputTypes(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	for _, typ := range []string{"password", "search", "tel", "text", "url"} {
		must.NoError(page.SetContent(ctx, `<input type="`+typ+`" id="inp" />`))
		expected := "text " + typ
		inp := page.Locator("#inp")
		must.NoErrorf(inp.Fill(ctx, expected), "Fill for type %q failed", typ)
		val, err := inp.InputValue(ctx)
		must.NoError(err)
		must.Equalf(expected, val, "input type=%q", typ)
	}
}

// TestPageFillDateInputAfterClicking verifies that a date input can be filled after clicking.
// Ref: TestPageFill.java#shouldFillDateInputAfterClicking
func TestPageFillDateInputAfterClicking(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=date>`))
	inp := page.Locator("input")
	must.NoError(inp.Click(ctx))
	must.NoError(inp.Fill(ctx, "2020-03-02"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("2020-03-02", val)
}

// TestPageFillThrowOnIncorrectDate verifies that filling a date input with an invalid date errors.
// Ref: TestPageFill.java#shouldThrowOnIncorrectDate (skip for WebKit)
func TestPageFillThrowOnIncorrectDate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName == "webkit" {
		t.Skip("WebKit does not validate date input values")
	}
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=date>`))
	err := page.Locator("input").Fill(ctx, "2020-13-05")
	is.Error(err)
	is.ErrorContains(err, "Malformed value")
}

// TestPageFillTimeInput verifies that a time input can be filled.
// Ref: TestPageFill.java#shouldFillTimeInput
func TestPageFillTimeInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=time>`))
	must.NoError(page.Locator("input").Fill(ctx, "13:15"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("13:15", val)
}

// TestPageFillThrowOnIncorrectTime verifies that filling a time input with an invalid value errors.
// Ref: TestPageFill.java#shouldThrowOnIncorrectTime (skip for WebKit)
func TestPageFillThrowOnIncorrectTime(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName == "webkit" {
		t.Skip("WebKit does not validate time input values")
	}
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=time>`))
	err := page.Locator("input").Fill(ctx, "25:05")
	is.Error(err)
	is.ErrorContains(err, "Malformed value")
}

// TestPageFillDatetimeLocalInput verifies that a datetime-local input can be filled.
// Ref: TestPageFill.java#shouldFillDatetimeLocalInput
func TestPageFillDatetimeLocalInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=datetime-local>`))
	must.NoError(page.Locator("input").Fill(ctx, "2020-03-02T05:15"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("2020-03-02T05:15", val)
}

// TestPageFillThrowOnIncorrectDatetimeLocal verifies that a datetime-local input rejects invalid values (Chromium only).
// Ref: TestPageFill.java#shouldThrowOnIncorrectDatetimeLocal
func TestPageFillThrowOnIncorrectDatetimeLocal(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("only Chromium validates datetime-local input values")
	}
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type=datetime-local>`))
	err := page.Locator("input").Fill(ctx, "abc")
	is.Error(err)
	is.ErrorContains(err, "Malformed value")
}

// TestPageFillContenteditable verifies that a contenteditable div can be filled.
// Ref: TestPageFill.java#shouldFillContenteditable
func TestPageFillContenteditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div contenteditable="true" id="ed"></div>`))
	must.NoError(page.Locator("#ed").Fill(ctx, "some value"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('div').textContent`)
	must.NoError(err)
	is.Equal("some value", val)
}

// TestPageFillExistingValueAndSelection verifies Fill works when an element has pre-existing value and selection.
// Ref: TestPageFill.java#shouldFillElementsWithExistingValueAndSelection
func TestPageFillExistingValueAndSelection(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" value="value one" />`))
	inp := page.Locator("#inp")
	must.NoError(inp.Fill(ctx, "another value"))
	val, err := inp.InputValue(ctx)
	must.NoError(err)
	is.Equal("another value", val)
}

// TestPageFillThrowOnNonFillableElement verifies that filling a body element errors.
// Ref: TestPageFill.java#shouldThrowWhenElementIsNotAnInputTextareaOrContenteditable
func TestPageFillThrowOnNonFillableElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="notfillable">text</div>`))
	err := page.Locator("#notfillable").Fill(ctx, "")
	is.Error(err, "expected error when filling a non-fillable element")
}

// TestPageFillBody verifies that a contentEditable body can be filled.
// Ref: TestPageFill.java#shouldBeAbleToFillTheBody
func TestPageFillBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body contentEditable='true'></body>`))
	must.NoError(page.Locator("body").Fill(ctx, "some value"))
	val, err := page.Evaluate(ctx, `() => document.body.textContent`)
	must.NoError(err)
	is.Equal("some value", val)
}

// TestPageFillFixedPositionInput verifies that a fixed-position input can be filled.
// Ref: TestPageFill.java#shouldFillFixedPositionInput
func TestPageFillFixedPositionInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input style='position: fixed;' />`))
	must.NoError(page.Locator("input").Fill(ctx, "some value"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("some value", val)
}

// TestPageFillWhenFocusInWrongFrame verifies fill works even when focus is in a different iframe.
// Ref: TestPageFill.java#shouldBeAbleToFillWhenFocusIsInTheWrongFrame
func TestPageFillWhenFocusInWrongFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div contentEditable='true' id='ed'></div><iframe></iframe>`))
	// Focus the iframe first (wrong frame)
	must.NoError(page.Locator("iframe").Focus(ctx))
	// Fill should still work on the div
	must.NoError(page.Locator("#ed").Fill(ctx, "some value"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('#ed').textContent`)
	must.NoError(err)
	is.Equal("some value", val)
}

// TestPageFillInputTypeNumber verifies that a number input can be filled with an integer.
// Ref: TestPageFill.java#shouldBeAbleToFillTheInputTypeNumber
func TestPageFillInputTypeNumber(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id='input' type='number'>`))
	must.NoError(page.Locator("input").Fill(ctx, "42"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("42", val)
}

// TestPageFillExponentIntoNumberInput verifies that an exponent notation value can be filled into a number input.
// Ref: TestPageFill.java#shouldBeAbleToFillExponentIntoTheInputTypeNumber
func TestPageFillExponentIntoNumberInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id='input' type='number'>`))
	must.NoError(page.Locator("input").Fill(ctx, "-10e5"))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("-10e5", val)
}

// TestPageFillNumberInputWithEmptyString verifies that a number input can be cleared.
// Ref: TestPageFill.java#shouldBeAbleToFillInputTypeNumberWithEmptyString
func TestPageFillNumberInputWithEmptyString(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id='input' type='number' value='123'>`))
	must.NoError(page.Locator("input").Fill(ctx, ""))
	val, err := page.Evaluate(ctx, `() => document.querySelector('input').value`)
	must.NoError(err)
	is.Equal("", val)
}

// TestPageFillTextIntoNumberInputShouldFail verifies that filling text into a number input errors.
// Ref: TestPageFill.java#shouldNotBeAbleToFillTextIntoTheInputTypeNumber
func TestPageFillTextIntoNumberInputShouldFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id='input' type='number'>`))
	err := page.Locator("input").Fill(ctx, "abc")
	is.Error(err)
	is.ErrorContains(err, "Cannot type text into input[type=number]")
}

// TestPageFillClear verifies that filling with an empty string clears the input.
// Ref: TestPageFill.java#shouldBeAbleToClear
func TestPageFillClear(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" />`))
	inp := page.Locator("#inp")
	must.NoError(inp.Fill(ctx, "some value"))
	val, err := inp.InputValue(ctx)
	must.NoError(err)
	is.Equal("some value", val)
	must.NoError(inp.Fill(ctx, ""))
	val, err = inp.InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestPageFillInputValue verifies that InputValue returns the correct value after filling.
// Ref: TestPageFill.java#inputValueShouldWork
func TestPageFillInputValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" />`))
	inp := page.Locator("#inp")
	must.NoError(inp.Fill(ctx, "my-text-content"))
	val, err := inp.InputValue(ctx)
	must.NoError(err)
	is.Equal("my-text-content", val)
	must.NoError(inp.Fill(ctx, ""))
	val, err = inp.InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestPagePressKey verifies that Locator.Press() sends a key to an input element.
func TestPagePressKey(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="inp" type="text" />
		<div id="out"></div>
		<script>
			document.getElementById('inp').addEventListener('keydown', function(e) {
				if (e.key === 'Enter') {
					document.getElementById('out').textContent = 'submitted';
				}
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	inp := page.Locator("#inp")
	err = inp.Press(ctx, "Enter")
	must.NoError(err, "Press(Enter) failed")

	text, err := page.Locator("#out").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("submitted", text)
}

// TestLocatorFillClearsExistingText verifies Fill replaces existing text in input.
// Ref: TestLocatorFill.java#shouldClearAndFill
func TestLocatorFillClearsExistingText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="old value">`))
	must.NoError(page.Locator("input").Fill(ctx, "new value"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("new value", val)
}

// TestLocatorFillEmptyStringClearsInput verifies Fill with empty string clears input.
// Ref: TestLocatorFill.java#shouldFillWithEmptyString
func TestLocatorFillEmptyStringClearsInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="some text">`))
	must.NoError(page.Locator("input").Fill(ctx, ""))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorFillTextarea verifies Fill works on textarea elements.
// Ref: TestLocatorFill.java#shouldFillInTextarea
func TestLocatorFillTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))
	must.NoError(page.Locator("textarea").Fill(ctx, "multi\nline\ncontent"))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("multi\nline\ncontent", val)
}

// TestLocatorFillWithSpecialChars verifies Fill handles special characters correctly.
// Ref: TestLocatorFill.java#shouldFillWithSpecialCharacters
func TestLocatorFillWithSpecialChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	specialText := "Hello! @#$%^&*() test"
	must.NoError(page.Locator("input").Fill(ctx, specialText))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal(specialText, val)
}

// TestLocatorFillPasswordInput verifies Fill works on password inputs.
// Ref: TestLocatorFill.java#shouldFillPasswordInput
func TestLocatorFillPasswordInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="password" id="pwd">`))
	must.NoError(page.Locator("input[type=password]").Fill(ctx, "s3cr3t"))

	val, err := page.Locator("input[type=password]").InputValue(ctx)
	must.NoError(err)
	is.Equal("s3cr3t", val)
}

// TestLocatorFillContentEditableDiv verifies Fill works on contenteditable elements.
// Ref: TestLocatorFill.java#shouldFillContentEditable
func TestLocatorFillContentEditableDiv(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div contenteditable="true" id="ce">old text</div>`))
	must.NoError(page.Locator("#ce").Fill(ctx, "new text"))

	text, err := page.Locator("#ce").InnerText(ctx)
	must.NoError(err)
	is.Equal("new text", text)
}

// TestLocatorFillDateInput verifies Fill works on date input.
// Ref: TestLocatorFill.java#shouldFillDateInput
func TestLocatorFillDateInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="date" id="date">`))

	must.NoError(page.Locator("#date").Fill(ctx, "2024-01-15"))

	val, err := page.Locator("#date").InputValue(ctx)
	must.NoError(err)
	is.Equal("2024-01-15", val)
}

// TestLocatorFillNumberInput verifies Fill works on number input.
// Ref: TestLocatorFill.java#shouldFillNumberInput
func TestLocatorFillNumberInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="number" id="num">`))

	must.NoError(page.Locator("#num").Fill(ctx, "42"))

	val, err := page.Locator("#num").InputValue(ctx)
	must.NoError(err)
	is.Equal("42", val)
}

// TestLocatorFillEmailInput verifies Fill works on email input.
// Ref: TestLocatorFill.java#shouldFillEmailInput
func TestLocatorFillEmailInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="email" id="email">`))

	must.NoError(page.Locator("#email").Fill(ctx, "user@example.com"))

	val, err := page.Locator("#email").InputValue(ctx)
	must.NoError(err)
	is.Equal("user@example.com", val)
}

// TestLocatorFillWithUnicodeChars verifies Fill works with Unicode characters.
// Ref: TestLocatorFill.java#shouldFillWithUnicode
func TestLocatorFillWithUnicodeChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "日本語テスト"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("日本語テスト", val)
}

// TestLocatorFillFiresInputEvent verifies Fill fires input event.
// Ref: TestLocatorFill.java#shouldFireInputEvent
func TestLocatorFillFiresInputEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="inp" oninput="window.__inputFired=true">
	`))

	must.NoError(page.Locator("#inp").Fill(ctx, "text"))

	fired, err := page.Evaluate(ctx, `() => window.__inputFired`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestLocatorFillFiresChangeEvent verifies Fill fires change event on blur.
// Ref: TestLocatorFill.java#shouldFireChangeEventOnBlur
func TestLocatorFillFiresChangeEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="inp" onchange="window.__changeFired=true">
		<button id="btn">blur</button>
	`))

	must.NoError(page.Locator("#inp").Fill(ctx, "text"))
	must.NoError(page.Locator("#btn").Click(ctx)) // trigger blur

	fired, err := page.Evaluate(ctx, `() => window.__changeFired`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestLocatorFillClearsExistingContent verifies Fill replaces existing content.
// Ref: TestLocatorFill.java#shouldClearAndFill
func TestLocatorFillClearsExistingContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="old value">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "new value"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("new value", val)
}

// TestLocatorFillWithEmptyString verifies Fill with empty string clears input.
// Ref: TestLocatorFill.java#shouldFillEmptyString
func TestLocatorFillWithEmptyString(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, ""))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Empty(val)
}

// TestLocatorFillTextarea verifies Fill works on textarea.
// Ref: TestLocatorFill.java#shouldFillTextarea
func TestLocatorFillTextareaEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	must.NoError(page.Locator("#ta").Fill(ctx, "multiline\ncontent"))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "multiline")
	is.Contains(val, "content")
}

// TestLocatorFillWithNumbers verifies Fill with numeric string.
// Ref: TestLocatorFill.java#shouldFillNumbers
func TestLocatorFillWithNumbers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "12345"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("12345", val)
}

// TestLocatorFillLongString verifies Fill works with long strings.
// Ref: TestLocatorFill.java#shouldFillLongString
func TestLocatorFillLongString(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	longText := "This is a very long string that contains many characters to test the fill functionality thoroughly."
	must.NoError(page.Locator("#ta").Fill(ctx, longText))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Equal(longText, val)
}

// TestLocatorFillPasswordInputEx4 verifies Fill works on password input.
// Ref: TestLocatorFill.java#shouldFillPasswordInput
func TestLocatorFillPasswordInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="pwd" type="password">`))

	must.NoError(page.Locator("#pwd").Fill(ctx, "secret123"))

	val, err := page.Locator("#pwd").InputValue(ctx)
	must.NoError(err)
	is.Equal("secret123", val)
}

// TestLocatorFillEmailInputEx4 verifies Fill works on email input.
// Ref: TestLocatorFill.java#shouldFillEmailInput
func TestLocatorFillEmailInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="email" type="email">`))

	must.NoError(page.Locator("#email").Fill(ctx, "user@example.com"))

	val, err := page.Locator("#email").InputValue(ctx)
	must.NoError(err)
	is.Equal("user@example.com", val)
}

// TestLocatorFillSearchInputEx4 verifies Fill works on search input.
// Ref: TestLocatorFill.java#shouldFillSearchInput
func TestLocatorFillSearchInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="search" type="search">`))

	must.NoError(page.Locator("#search").Fill(ctx, "my search query"))

	val, err := page.Locator("#search").InputValue(ctx)
	must.NoError(err)
	is.Equal("my search query", val)
}

// TestLocatorFillTelInputEx4 verifies Fill works on tel input.
// Ref: TestLocatorFill.java#shouldFillTelInput
func TestLocatorFillTelInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="tel" type="tel">`))

	must.NoError(page.Locator("#tel").Fill(ctx, "+1-555-555-5555"))

	val, err := page.Locator("#tel").InputValue(ctx)
	must.NoError(err)
	is.Equal("+1-555-555-5555", val)
}

// TestLocatorFillURLInputEx4 verifies Fill works on url input.
// Ref: TestLocatorFill.java#shouldFillURLInput
func TestLocatorFillURLInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="url" type="url">`))

	must.NoError(page.Locator("#url").Fill(ctx, "https://example.com"))

	val, err := page.Locator("#url").InputValue(ctx)
	must.NoError(err)
	is.Equal("https://example.com", val)
}

// TestLocatorFillClearsExistingEx5 verifies Fill replaces existing value.
// Ref: TestLocatorFill.java#shouldClearAndFill
func TestLocatorFillClearsExistingEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="old value">`))
	must.NoError(page.Locator("#inp").Fill(ctx, "new value"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("new value", val)
}

// TestLocatorFillEmptyStringEx5 verifies Fill with empty string clears input.
// Ref: TestLocatorFill.java#shouldFillWithEmptyString
func TestLocatorFillEmptyStringEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="something">`))
	must.NoError(page.Locator("#inp").Fill(ctx, ""))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorFillUnicodeEx5 verifies Fill handles unicode characters.
// Ref: TestLocatorFill.java#shouldHandleUnicode
func TestLocatorFillUnicodeEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Fill(ctx, "こんにちは"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("こんにちは", val)
}

// TestLocatorFillFiresInputEventEx5 verifies Fill triggers input event.
// Ref: TestLocatorFill.java#shouldFireInputEvent
func TestLocatorFillFiresInputEventEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text">
		<script>
			var fired = false;
			document.getElementById('inp').addEventListener('input', function() { fired = true; });
		</script>
	`))
	must.NoError(page.Locator("#inp").Fill(ctx, "text"))

	result, err := page.Evaluate(ctx, `() => fired`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorFillContenteditableEx5 verifies Fill works on contenteditable elements.
// Ref: TestLocatorFill.java#shouldFillContenteditable
func TestLocatorFillContenteditableEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="editor" contenteditable="true"></div>`))
	must.NoError(page.Locator("#editor").Fill(ctx, "editable content"))

	text, err := page.Locator("#editor").InnerText(ctx)
	must.NoError(err)
	is.Equal("editable content", text)
}

// TestLocatorFillDateInputEx6 verifies Fill works on date input.
// Ref: TestLocatorFill.java#shouldFillDateInput
func TestLocatorFillDateInputEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="d" type="date">`))
	must.NoError(page.Locator("#d").Fill(ctx, "2024-01-15"))

	val, err := page.Locator("#d").InputValue(ctx)
	must.NoError(err)
	is.Equal("2024-01-15", val)
}

// TestLocatorFillLargeTextEx6 verifies Fill handles large text.
// Ref: TestLocatorFill.java#shouldHandleLargeText
func TestLocatorFillLargeTextEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	largeText := ""
	for i := 0; i < 100; i++ {
		largeText += "abcdefghij"
	}

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))
	must.NoError(page.Locator("#ta").Fill(ctx, largeText))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Equal(largeText, val)
}

// TestLocatorFillWithNewlinesEx6 verifies Fill handles newlines in textarea.
// Ref: TestLocatorFill.java#shouldHandleNewlines
func TestLocatorFillWithNewlinesEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))
	must.NoError(page.Locator("#ta").Fill(ctx, "line1\nline2\nline3"))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "line1")
	is.Contains(val, "line2")
}

// TestLocatorFillEmojiEx6 verifies Fill handles emoji characters.
// Ref: TestLocatorFill.java#shouldHandleEmoji
func TestLocatorFillEmojiEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Fill(ctx, "Hello 🌍"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "Hello")
}

// TestFillSearchInputEx7 verifies Fill works for search inputs.
// Ref: TestLocatorFill.java#shouldFillSearchInput
func TestFillSearchInputEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="search" type="search">`))

	must.NoError(page.Locator("#search").Fill(ctx, "playwright"))

	val, err := page.Locator("#search").InputValue(ctx)
	must.NoError(err)
	is.Equal("playwright", val)
}

// TestFillEmailInputEx7 verifies Fill works for email inputs.
// Ref: TestLocatorFill.java#shouldFillEmailInput
func TestFillEmailInputEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="email" type="email">`))

	must.NoError(page.Locator("#email").Fill(ctx, "test@example.com"))

	val, err := page.Locator("#email").InputValue(ctx)
	must.NoError(err)
	is.Equal("test@example.com", val)
}

// TestFillUrlInputEx7 verifies Fill works for URL inputs.
// Ref: TestLocatorFill.java#shouldFillUrlInput
func TestFillUrlInputEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="url" type="url">`))

	must.NoError(page.Locator("#url").Fill(ctx, "https://example.com"))

	val, err := page.Locator("#url").InputValue(ctx)
	must.NoError(err)
	is.Equal("https://example.com", val)
}

// TestFillPhoneInputEx7 verifies Fill works for tel inputs.
// Ref: TestLocatorFill.java#shouldFillTelInput
func TestFillPhoneInputEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="tel" type="tel">`))

	must.NoError(page.Locator("#tel").Fill(ctx, "+1-555-0100"))

	val, err := page.Locator("#tel").InputValue(ctx)
	must.NoError(err)
	is.Equal("+1-555-0100", val)
}

// TestFillDateInputEx8 verifies Fill works for date inputs.
// Ref: TestLocatorFill.java#shouldFillDateInput
func TestFillDateInputEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="dt" type="date">`))

	must.NoError(page.Locator("#dt").Fill(ctx, "2025-01-15"))

	val, err := page.Locator("#dt").InputValue(ctx)
	must.NoError(err)
	is.Equal("2025-01-15", val)
}

// TestFillTextareaLargeEx8 verifies Fill works for large text in textarea.
// Ref: TestLocatorFill.java#shouldFillLargeText
func TestFillTextareaLargeEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta" rows="10" cols="50"></textarea>`))

	largeText := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	must.NoError(page.Locator("#ta").Fill(ctx, largeText))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Equal(largeText, val)
}

// TestFillClearsExistingEx8 verifies Fill clears existing content before typing.
// Ref: TestLocatorFill.java#shouldClearBeforeFill
func TestFillClearsExistingEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="old value">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "new value"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("new value", val)
}

// TestFillHiddenInputEx8 verifies Fill works on hidden-to-user but accessible input.
// Ref: TestLocatorFill.java#shouldFillAccessibleInput
func TestFillHiddenInputEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" style="opacity:0.1">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "hidden but accessible"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hidden but accessible", val)
}

// TestFillColorInputEx9 verifies Fill works for color inputs.
// Ref: TestLocatorFill.java#shouldFillColorInput
func TestFillColorInputEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="c" type="color">`))

	must.NoError(page.Locator("#c").Fill(ctx, "#ff0000"))

	val, err := page.Locator("#c").InputValue(ctx)
	must.NoError(err)
	is.Equal("#ff0000", val)
}

// TestFillNumberInputEx9 verifies Fill works for number inputs.
// Ref: TestLocatorFill.java#shouldFillNumberInput
func TestFillNumberInputEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="num" type="number" min="1" max="100">`))

	must.NoError(page.Locator("#num").Fill(ctx, "42"))

	val, err := page.Locator("#num").InputValue(ctx)
	must.NoError(err)
	is.Equal("42", val)
}

// TestFillTimeInputEx9 verifies Fill works for time inputs.
// Ref: TestLocatorFill.java#shouldFillTimeInput
func TestFillTimeInputEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="t" type="time">`))

	must.NoError(page.Locator("#t").Fill(ctx, "14:30"))

	val, err := page.Locator("#t").InputValue(ctx)
	must.NoError(err)
	is.Equal("14:30", val)
}

// TestFillMultilineTextareaEx9 verifies Fill works with unicode content.
// Ref: TestLocatorFill.java#shouldFillUnicode
func TestFillMultilineTextareaEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	unicodeText := "Hello 世界 🌍"
	must.NoError(page.Locator("#ta").Fill(ctx, unicodeText))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Equal(unicodeText, val)
}

// TestFillRangeInputEx10 verifies Fill sets value on range input.
// Ref: TestLocatorFill.java#shouldFillRange
func TestFillRangeInputEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="r" type="range" min="0" max="100" value="0">`))

	must.NoError(page.Locator("#r").Fill(ctx, "75"))

	val, err := page.Locator("#r").InputValue(ctx)
	must.NoError(err)
	is.Equal("75", val)
}

// TestFillMultiLineEx10 verifies Fill with multi-line text in textarea.
// Ref: TestLocatorFill.java#shouldFillMultiLine
func TestFillMultiLineEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta" rows="5"></textarea>`))

	must.NoError(page.Locator("#ta").Fill(ctx, "Line 1\nLine 2\nLine 3"))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "Line 1")
	is.Contains(val, "Line 3")
}

// TestFillPasswordInputEx10 verifies Fill works on password input.
// Ref: TestLocatorFill.java#shouldFillPassword
func TestFillPasswordInputEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="pw" type="password">`))

	must.NoError(page.Locator("#pw").Fill(ctx, "secret123"))

	val, err := page.Locator("#pw").InputValue(ctx)
	must.NoError(err)
	is.Equal("secret123", val)
}

// TestFillHiddenInputEx10 verifies Fill sets value on hidden input via JS.
// Ref: TestLocatorFill.java#shouldFillHiddenViaJS
func TestFillHiddenInputEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="h" type="text" style="opacity:0;position:absolute;">`))

	must.NoError(page.Locator("#h").Fill(ctx, "hidden value"))

	val, err := page.Locator("#h").InputValue(ctx)
	must.NoError(err)
	is.Equal("hidden value", val)
}

// Ref: TestPageFill.java#shouldBeAbleToClearUsingFill
func TestShouldBeAbleToClearUsingFill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" />`))
	must.NoError(page.Locator("#inp").Fill(ctx, "some value"))
	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("some value", val)

	must.NoError(page.Locator("#inp").Fill(ctx, ""))
	val, err = page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}
