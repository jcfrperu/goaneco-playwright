//go:build e2e

package e2e

import (
	"regexp"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocatorSetChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox" id="cb"/>`)
	must.NoError(err, "SetContent failed")

	// Check it.
	err = page.Locator("#cb").SetChecked(ctx, true)
	must.NoError(err, "SetChecked(true) failed")
	checked, err := page.Locator("#cb").IsChecked(ctx)
	if err != nil || !checked {
		t.Errorf("expected checkbox to be checked, got %v err=%v", checked, err)
	}

	// Uncheck it.
	err = page.Locator("#cb").SetChecked(ctx, false)
	must.NoError(err, "SetChecked(false) failed")
	checked, err = page.Locator("#cb").IsChecked(ctx)
	if err != nil || checked {
		t.Errorf("expected checkbox to be unchecked, got %v err=%v", checked, err)
	}
}

func TestLocatorWaitForVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" style="display:none">hello</div>`)
	must.NoError(err, "SetContent failed")

	// Make it visible after 100ms.
	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('el').style.display = 'block'`)
	}()

	err = page.Locator("#el").WaitFor(ctx, "visible")
	must.NoError(err, "WaitFor(visible) failed")
}

func TestLocatorWaitForHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">hello</div>`)
	must.NoError(err, "SetContent failed")

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('el').style.display = 'none'`)
	}()

	err = page.Locator("#el").WaitFor(ctx, "hidden")
	must.NoError(err, "WaitFor(hidden) failed")
}

func TestLocatorClearInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="in" value="hello"/>`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#in").Clear(ctx)
	must.NoError(err, "Clear failed")

	val, err := page.Locator("#in").InputValue(ctx)
	must.NoError(err, "InputValue after clear failed")
	is.Equal("", val)
}

func TestLocatorFocusButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn">Click me</button>`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#btn").Focus(ctx)
	must.NoError(err, "Focus failed")

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err, "Evaluate activeElement failed")
	if focused != "btn" {
		t.Errorf("expected activeElement to be 'btn', got %v", focused)
	}
}

func TestLocatorPressSequentially(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="in"/>`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#in").PressSequentially(ctx, "hello")
	must.NoError(err, "PressSequentially failed")

	val, err := page.Locator("#in").InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("hello", val)
}

func TestLocatorWaitForAlreadyVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">visible</div>`)
	must.NoError(err, "SetContent failed")

	// Should return immediately when already visible.
	err = page.Locator("#el").WaitFor(ctx)
	must.NoError(err, "WaitFor() for already-visible element failed")
}

// TestLocatorHasDoesNotEncodeUnicode verifies that Unicode text in locator filters is not encoded
// and appears verbatim in the error message when the element is not found.
// Ref: TestLocatorMisc.java#locatorsHasDoesNotEncodeUnicode
func TestLocatorHasDoesNotEncodeUnicode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	srv := testserver.New(t)
	page := newPage(t)

	ctx := testCtx(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	unicode := "Драматург"
	locators := []*playwright.Locator{
		page.Locator("button").Filter(&playwright.LocatorFilterOptions{HasText: &unicode}),
		page.Locator("button").Filter(&playwright.LocatorFilterOptions{HasTextRegex: regexp.MustCompile(unicode)}),
		page.Locator("button").Filter(&playwright.LocatorFilterOptions{Has: page.Locator("text=" + unicode)}),
	}

	for _, loc := range locators {
		// Use a short context so the click times out quickly (element does not exist).
		shortC := shortCtx(t)
		err := loc.Click(shortC)
		is.Error(err, "click on non-existent unicode element should fail")
		is.Contains(err.Error(), unicode, "error message should contain the unicode text verbatim")
	}
}

// TestLocatorFilterVisible verifies filter(Visible=true/false) narrows by element visibility.
// Ref: TestLocatorMisc.java#shouldSupportFilterVisible
func TestLocatorFilterVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>
  <div class="item" style="display: none">Hidden data0</div>
  <div class="item">visible data1</div>
  <div class="item" style="display: none">Hidden data1</div>
  <div class="item">visible data2</div>
  <div class="item" style="display: none">Hidden data2</div>
  <div class="item">visible data3</div>
</div>`)
	must.NoError(err)

	// Nth(1) among visible items should be "visible data2"
	must.NoError(playwright.Expect(
		page.Locator(".item").Filter(&playwright.LocatorFilterOptions{Visible: boolPtr(true)}).Nth(1),
	).ToHaveText(ctx, "visible data2"))

	// Filter visible + has text "data3" → "visible data3"
	data3 := "data3"
	must.NoError(playwright.Expect(
		page.Locator(".item").
			Filter(&playwright.LocatorFilterOptions{Visible: boolPtr(true)}).
			Filter(&playwright.LocatorFilterOptions{HasText: &data3}),
	).ToHaveText(ctx, "visible data3"))

	// Filter not-visible + has text "data1" → "Hidden data1"
	data1 := "data1"
	must.NoError(playwright.Expect(
		page.Locator(".item").
			Filter(&playwright.LocatorFilterOptions{Visible: boolPtr(false)}).
			Filter(&playwright.LocatorFilterOptions{HasText: &data1}),
	).ToHaveText(ctx, "Hidden data1"))
}

// TestLocatorWaitForFunctionAttributeAppears verifies waitForFunction fires once an attribute is added.
// Ref: TestLocatorMisc.java#waitForFunctionShouldWaitForAnAttributeToAppear
func TestLocatorWaitForFunctionAttributeAppears(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id=toggle>Menu</button>`))

	_, err := page.Evaluate(ctx, `() => setTimeout(() =>
		document.querySelector('#toggle').setAttribute('aria-expanded', 'true'), 500)`)
	must.NoError(err)

	must.NoError(
		page.Locator("#toggle").WaitForFunction(ctx, "element => element.hasAttribute('aria-expanded')", nil))
}

// TestLocatorWaitForFunctionTruthy verifies waitForFunction passes immediately when already truthy.
// Ref: TestLocatorMisc.java#waitForFunctionShouldReturnImmediatelyWhenAlreadyTruthy
func TestLocatorWaitForFunctionTruthy(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id=target>yes</div>`))

	must.NoError(
		page.Locator("#target").WaitForFunction(ctx, "element => element.textContent === 'yes'", nil))
}

// TestLocatorWaitForFunctionElementHandleArg verifies waitForFunction accepts an extra ElementHandle argument.
// Ref: TestLocatorMisc.java#waitForFunctionShouldAcceptElementHandleArguments
func TestLocatorWaitForFunctionElementHandleArg(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id=a></div><div id=b>value</div>`))

	handle, err := page.QuerySelector(ctx, "#b")
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(
		page.Locator("#a").WaitForFunction(ctx, "(element, other) => other.textContent === 'value'", handle))
}

// TestLocatorWaitForFunctionPredicateThrows verifies that a throwing predicate propagates an error.
// Ref: TestLocatorMisc.java#waitForFunctionShouldThrowWhenPredicateThrows
func TestLocatorWaitForFunctionPredicateThrows(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id=target>no</div>`))

	err := page.Locator("#target").WaitForFunction(ctx, "() => { throw new Error('oh my'); }", nil)
	is.Error(err)
	is.ErrorContains(err, "oh my")
}

// TestLocatorWaitForFunctionStrictModeViolation verifies that a locator matching multiple elements fails.
// Ref: TestLocatorMisc.java#waitForFunctionShouldThrowOnStrictModeViolation
func TestLocatorWaitForFunctionStrictModeViolation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class=x>1</div><div class=x>2</div>`))

	err := page.Locator("div.x").WaitForFunction(ctx, "() => true", nil)
	is.Error(err)
	is.ErrorContains(err, "strict mode violation")
}

// TestLocatorWaitForFunctionRespectTimeout verifies that the timeout option causes a timeout error.
// Ref: TestLocatorMisc.java#waitForFunctionShouldRespectTimeout
func TestLocatorWaitForFunctionRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id=target>no</div>`))

	err := page.Locator("#target").WaitForFunction(ctx,
		"element => element.textContent === 'yes'", nil,
		&playwright.LocatorWaitForFunctionOptions{Timeout: 500})
	is.Error(err)
	is.ErrorContains(err, "Timeout 500ms exceeded")
}

// Ref: TestLocatorMisc.java#shouldCheckTheBoxUsingSetChecked
func TestLocatorMiscShouldCheckTheBoxUsingSetChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id='checkbox' type='checkbox'></input>`))
	input := page.Locator("input")
	must.NoError(input.SetChecked(ctx, true))
	v, err := page.Evaluate(ctx, "checkbox.checked")
	must.NoError(err)
	is.Equal(true, v)
	must.NoError(input.SetChecked(ctx, false))
	v, err = page.Evaluate(ctx, "checkbox.checked")
	must.NoError(err)
	is.Equal(false, v)
}

// Ref: TestLocatorMisc.java#shouldWaitFor
func TestLocatorMiscShouldWaitFor(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))
	locator := page.Locator("span")
	_, err := page.EvalOnSelector(ctx, "div", "div => setTimeout(() => div.innerHTML = '<span>target</span>', 500)")
	must.NoError(err)
	must.NoError(locator.WaitFor(ctx))
	tc, err := locator.TextContent(ctx)
	must.NoError(err)
	is.Contains(*tc, "target")
}

// Ref: TestLocatorMisc.java#shouldWaitForHidden
func TestLocatorMiscShouldWaitForHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>target</span></div>`))
	locator := page.Locator("span")
	_, err := page.EvalOnSelector(ctx, "div", "div => setTimeout(() => div.innerHTML = '', 500)")
	must.NoError(err)
	must.NoError(locator.WaitFor(ctx, "hidden"))
}

// Ref: TestLocatorMisc.java#shouldClearInput
func TestLocatorMiscShouldClearInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	handle := page.Locator("input")
	must.NoError(handle.Fill(ctx, "some value"))
	v, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("some value", v)
	must.NoError(handle.Clear(ctx))
	v, err = page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("", v)
}

// Ref: TestLocatorMisc.java#shouldFocusAndBlurAButton
func TestLocatorMiscShouldFocusAndBlurAButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	button := page.Locator("button")

	active, err := button.Evaluate(ctx, "button => document.activeElement === button")
	must.NoError(err)
	is.Equal(false, active)

	must.NoError(button.Focus(ctx))

	active, err = button.Evaluate(ctx, "button => document.activeElement === button")
	must.NoError(err)
	is.Equal(true, active)

	// Blur via evaluate since Locator.Blur is not yet implemented.
	_, err = button.Evaluate(ctx, "button => button.blur()")
	must.NoError(err)

	active, err = button.Evaluate(ctx, "button => document.activeElement === button")
	must.NoError(err)
	is.Equal(false, active)
}

// Ref: TestLocatorMisc.java#LocatorLocatorAndFrameLocatorLocatorShouldAcceptLocator
func TestLocatorMiscLocatorLocatorAndFrameLocatorLocatorShouldAcceptLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><input value=outer></div>`+
		`<iframe srcdoc="<div><input value=inner></div>"></iframe>`))

	inputLocator := page.Locator("input")
	v, err := inputLocator.InputValue(ctx)
	must.NoError(err)
	is.Equal("outer", v)

	v, err = page.Locator("div").Locator(inputLocator.Selector()).InputValue(ctx)
	must.NoError(err)
	is.Equal("outer", v)

	v, err = page.FrameLocator("iframe").Locator(inputLocator.Selector()).InputValue(ctx)
	must.NoError(err)
	is.Equal("inner", v)

	v, err = page.FrameLocator("iframe").Locator("div").Locator(inputLocator.Selector()).InputValue(ctx)
	must.NoError(err)
	is.Equal("inner", v)
}

// Ref: TestLocatorMisc.java#shouldPressSequentially
func TestLocatorMiscShouldPressSequentially(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type='text' />`))
	must.NoError(page.Locator("input").PressSequentially(ctx, "hello"))
	v, err := page.EvalOnSelector(ctx, "input", "input => input.value")
	must.NoError(err)
	is.Equal("hello", v)
}
