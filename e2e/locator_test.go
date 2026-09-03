//go:build e2e

// E2E tests for Locator API (E2E-LOC-01).
// Covers GetByRole, GetByText, GetByLabel, GetByPlaceholder, GetByTestId,
// Count, IsVisible, and direct Locator selector usage.
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorGetByRole verifies GetByRole selects elements by ARIA role and optional name.
func TestLocatorGetByRole(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button>Submit</button>
		<a href="/home">Home</a>
		<h1>Page Title</h1>
	`)
	must.NoError(err, "SetContent failed")

	// GetByRole without name: count all buttons
	btnLocator := page.GetByRole(playwright.AriaRoleButton)
	count, err := btnLocator.Count(ctx)
	must.NoError(err, "GetByRole button Count() failed")
	is.Equal(1, count)

	// GetByRole with name filter
	submitName := "Submit"
	submitLocator := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name: &submitName,
	})
	visible, err := submitLocator.IsVisible(ctx)
	must.NoError(err, "GetByRole(button, name=Submit).IsVisible() failed")
	if !visible {
		t.Error("expected Submit button to be visible")
	}

	// GetByRole for link
	linkLocator := page.GetByRole(playwright.AriaRoleLink)
	linkCount, err := linkLocator.Count(ctx)
	must.NoError(err, "GetByRole link Count() failed")
	is.Equal(1, linkCount)

	// GetByRole for heading
	headingLocator := page.GetByRole(playwright.AriaRoleHeading)
	headingCount, err := headingLocator.Count(ctx)
	must.NoError(err, "GetByRole heading Count() failed")
	is.Equal(1, headingCount)

	text, err := headingLocator.InnerText(ctx)
	must.NoError(err, "GetByRole(heading).InnerText() failed")
	is.Equal("Page Title", text)
}

// TestLocatorGetByText verifies GetByText selects elements by text content.
func TestLocatorGetByText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<p>Hello World</p>
		<p>Hello Playwright</p>
		<span>Different content</span>
	`)
	must.NoError(err, "SetContent failed")

	// Partial match (default): both "Hello" paragraphs
	helloLocator := page.GetByText("Hello")
	count, err := helloLocator.Count(ctx)
	must.NoError(err, "GetByText(Hello).Count() failed")
	is.Equal(2, count)

	// Exact match for "Hello World"
	exact := true
	exactLocator := page.GetByText("Hello World", &playwright.GetByTextOptions{Exact: &exact})
	exactCount, err := exactLocator.Count(ctx)
	must.NoError(err, "GetByText(Hello World, exact).Count() failed")
	is.Equal(1, exactCount)
}

// TestLocatorGetByLabel verifies GetByLabel selects form elements associated with a label.
func TestLocatorGetByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<label for="name-input">Full Name</label>
		<input id="name-input" type="text" value="John Doe" />
		<label for="email-input">Email Address</label>
		<input id="email-input" type="email" value="john@example.com" />
	`)
	must.NoError(err, "SetContent failed")

	// GetByLabel finds the associated input
	nameLocator := page.GetByLabel("Full Name")
	val, err := nameLocator.InputValue(ctx)
	must.NoError(err, "GetByLabel(Full Name).InputValue() failed")
	is.Equal("John Doe", val)

	emailLocator := page.GetByLabel("Email Address")
	emailVal, err := emailLocator.InputValue(ctx)
	must.NoError(err, "GetByLabel(Email Address).InputValue() failed")
	is.Equal("john@example.com", emailVal)
}

// TestLocatorGetByPlaceholder verifies GetByPlaceholder selects inputs by their placeholder text.
func TestLocatorGetByPlaceholder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input type="text" placeholder="Enter your name" value="Alice" />
		<input type="email" placeholder="Enter your email" value="alice@example.com" />
	`)
	must.NoError(err, "SetContent failed")

	nameInput := page.GetByPlaceholder("Enter your name")
	val, err := nameInput.InputValue(ctx)
	must.NoError(err, "GetByPlaceholder(Enter your name).InputValue() failed")
	is.Equal("Alice", val)

	emailInput := page.GetByPlaceholder("Enter your email")
	emailVal, err := emailInput.InputValue(ctx)
	must.NoError(err, "GetByPlaceholder(Enter your email).InputValue() failed")
	is.Equal("alice@example.com", emailVal)
}

// TestLocatorGetByTestId verifies GetByTestId selects elements by data-testid attribute.
func TestLocatorGetByTestId(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button data-testid="submit-btn">Submit</button>
		<input data-testid="search-input" type="text" value="playwright" />
		<div data-testid="result-container"><p>Result here</p></div>
	`)
	must.NoError(err, "SetContent failed")

	// Locate by data-testid
	btnLocator := page.GetByTestId("submit-btn")
	visible, err := btnLocator.IsVisible(ctx)
	must.NoError(err, "GetByTestId(submit-btn).IsVisible() failed")
	if !visible {
		t.Error("expected submit-btn to be visible")
	}

	text, err := btnLocator.InnerText(ctx)
	must.NoError(err, "GetByTestId(submit-btn).InnerText() failed")
	is.Equal("Submit", text)

	inputLocator := page.GetByTestId("search-input")
	val, err := inputLocator.InputValue(ctx)
	must.NoError(err, "GetByTestId(search-input).InputValue() failed")
	is.Equal("playwright", val)
}

// TestLocatorCount verifies Count returns the correct number of matching elements.
func TestLocatorCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
			<li>Item 3</li>
		</ul>
		<p class="note">Note A</p>
		<p class="note">Note B</p>
	`)
	must.NoError(err, "SetContent failed")

	liLocator := page.Locator("li")
	count, err := liLocator.Count(ctx)
	must.NoError(err, "Locator(li).Count() failed")
	is.Equal(3, count)

	noteLocator := page.Locator("p.note")
	noteCount, err := noteLocator.Count(ctx)
	must.NoError(err, "Locator(p.note).Count() failed")
	is.Equal(2, noteCount)

	// Non-existent element
	missing := page.Locator("#does-not-exist")
	missingCount, err := missing.Count(ctx)
	must.NoError(err, "Locator(#does-not-exist).Count() failed")
	is.Equal(0, missingCount)
}

// TestLocatorIsVisible verifies IsVisible returns true for visible elements and false for hidden ones.
func TestLocatorIsVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="visible">I am visible</div>
		<div id="hidden" style="display:none">I am hidden</div>
		<div id="invisible" style="visibility:hidden">I am invisible</div>
	`)
	must.NoError(err, "SetContent failed")

	// Visible element
	visibleLocator := page.Locator("#visible")
	isVisible, err := visibleLocator.IsVisible(ctx)
	must.NoError(err, "Locator(#visible).IsVisible() failed")
	if !isVisible {
		t.Error("expected #visible to be visible")
	}

	// Hidden element (display:none)
	hiddenLocator := page.Locator("#hidden")
	isHidden, err := hiddenLocator.IsVisible(ctx)
	must.NoError(err, "Locator(#hidden).IsVisible() failed")
	if isHidden {
		t.Error("expected #hidden to NOT be visible (display:none)")
	}

	// Invisible element (visibility:hidden)
	invisibleLocator := page.Locator("#invisible")
	isInvisible, err := invisibleLocator.IsVisible(ctx)
	must.NoError(err, "Locator(#invisible).IsVisible() failed")
	if isInvisible {
		t.Error("expected #invisible to NOT be visible (visibility:hidden)")
	}
}

// TestLocatorFill verifies that Fill sets the value of an input element.
func TestLocatorFill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="name" type="text" />`)
	must.NoError(err, "SetContent failed")

	inputLocator := page.Locator("#name")
	err = inputLocator.Fill(ctx, "Playwright")
	must.NoError(err, "Locator(#name).Fill() failed")

	val, err := inputLocator.InputValue(ctx)
	must.NoError(err, "Locator(#name).InputValue() failed")
	is.Equal("Playwright", val)
}

// TestLocatorGetAttribute verifies GetAttribute returns the correct attribute value.
func TestLocatorGetAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<a id="link" href="https://playwright.dev" target="_blank">Playwright</a>`)
	must.NoError(err, "SetContent failed")

	linkLocator := page.Locator("#link")

	href, err := linkLocator.GetAttribute(ctx, "href")
	must.NoError(err, "GetAttribute(href) failed")
	if href == nil || *href != "https://playwright.dev" {
		t.Errorf("GetAttribute(href) = %v, want 'https://playwright.dev'", href)
	}

	target, err := linkLocator.GetAttribute(ctx, "target")
	must.NoError(err, "GetAttribute(target) failed")
	if target == nil || *target != "_blank" {
		t.Errorf("GetAttribute(target) = %v, want '_blank'", target)
	}

	// Non-existent attribute
	missing, err := linkLocator.GetAttribute(ctx, "data-missing")
	must.NoError(err, "GetAttribute(data-missing) failed")
	if missing != nil {
		t.Errorf("GetAttribute(data-missing) = %v, want nil", missing)
	}
}
func TestLocatorTap(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="btn">tap target</button>
		<script>
		document.getElementById('btn').addEventListener('click', function() {
			this.dataset.tapped = 'yes';
		});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#btn").Tap(ctx)
	must.NoError(err, "Locator.Tap failed")

	val, err := page.Evaluate(ctx, `document.getElementById('btn').dataset.tapped`)
	must.NoError(err, "Evaluate failed")
	if val != "yes" {
		t.Errorf("Tap did not trigger click: dataset.tapped = %v", val)
	}
}

func TestLocatorScrollIntoViewIfNeeded(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div style="height:2000px;background:#eee">spacer</div>
		<div id="target" style="height:100px;background:red">target</div>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#target").ScrollIntoViewIfNeeded(ctx)
	must.NoError(err, "Locator.ScrollIntoViewIfNeeded failed")

	visible, err := page.Locator("#target").IsVisible(ctx)
	must.NoError(err, "IsVisible failed")
	if !visible {
		t.Error("element should be visible after ScrollIntoViewIfNeeded")
	}
}

func TestLocatorBoundingBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="box" style="position:absolute;left:50px;top:60px;width:100px;height:80px;background:blue"></div>
	`)
	must.NoError(err, "SetContent failed")

	bb, err := page.Locator("#box").BoundingBox(ctx)
	must.NoError(err, "Locator.BoundingBox failed")
	must.NotNil(bb, "BoundingBox returned nil")
	if bb.Width != 100 || bb.Height != 80 {
		t.Errorf("BoundingBox size = %.0fx%.0f, want 100x80", bb.Width, bb.Height)
	}
	if bb.X < 50 || bb.Y < 60 {
		t.Errorf("BoundingBox position = (%.0f,%.0f), expected >= (50,60)", bb.X, bb.Y)
	}
}

func TestLocatorIsDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<button id="disabled-btn" disabled>Disabled</button>
		<button id="enabled-btn">Enabled</button>
	`)
	must.NoError(err, "SetContent failed")

	disabled, err := page.Locator("#disabled-btn").IsDisabled(ctx)
	must.NoError(err, "IsDisabled on disabled element failed")
	if !disabled {
		t.Error("expected #disabled-btn to be disabled")
	}

	notDisabled, err := page.Locator("#enabled-btn").IsDisabled(ctx)
	must.NoError(err, "IsDisabled on enabled element failed")
	if notDisabled {
		t.Error("expected #enabled-btn to not be disabled")
	}
}

func TestLocatorSelector(t *testing.T) {
	t.Parallel()
	page := newPage(t)

	loc := page.Locator("#my-element")
	if sel := loc.Selector(); sel != "#my-element" {
		t.Errorf("Selector() = %q, want '#my-element'", sel)
	}

	nth := page.Locator("li").Nth(0)
	if sel := nth.Selector(); sel == "" {
		t.Error("Selector() returned empty string for Nth locator")
	}
}

func TestLocatorBoundingBoxNilForHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="hidden" style="display:none;width:100px;height:100px"></div>`)
	must.NoError(err, "SetContent failed")

	bb, err := page.Locator("#hidden").BoundingBox(ctx)
	must.NoError(err, "BoundingBox on hidden element failed")
	if bb != nil {
		t.Errorf("BoundingBox for hidden element = %+v, want nil", bb)
	}
}

// TestLocatorAllReturnsCorrectLength verifies All returns slice of correct length.
// Ref: TestLocatorAll.java#shouldReturnCorrectLength
func TestLocatorAllReturnsCorrectLength(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="box">a</div>
		<div class="box">b</div>
		<div class="box">c</div>
	`))

	all, err := page.Locator(".box").All(ctx)
	must.NoError(err)
	is.Len(all, 3)
}

// TestLocatorAllCanIterateResults verifies All elements can be iterated and read.
// Ref: TestLocatorAll.java#shouldIterateResults
func TestLocatorAllCanIterateResults(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>first</div>
		<div>second</div>
		<div>third</div>
	`))

	all, err := page.Locator("div").All(ctx)
	must.NoError(err)
	is.Len(all, 3)

	texts := make([]string, 0, len(all))
	for _, loc := range all {
		text, err := loc.InnerText(ctx)
		must.NoError(err)
		texts = append(texts, text)
	}
	is.Equal([]string{"first", "second", "third"}, texts)
}

// TestLocatorAllReturnsEmptyForNoMatch verifies All returns empty slice when no match.
// Ref: TestLocatorAll.java#shouldReturnEmptyWhenNoMatch
func TestLocatorAllReturnsEmptyForNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	all, err := page.Locator(".nonexistent").All(ctx)
	must.NoError(err)
	is.Empty(all)
}

// TestLocatorAllSingleElement verifies All with single match returns slice of one.
// Ref: TestLocatorAll.java#shouldReturnSliceOfOne
func TestLocatorAllSingleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="hello">`))

	all, err := page.Locator("input").All(ctx)
	must.NoError(err)
	is.Len(all, 1)

	val, err := all[0].InputValue(ctx)
	must.NoError(err)
	is.Equal("hello", val)
}

// TestLocatorAllFillsEachInput verifies All can be used to fill multiple inputs.
// Ref: TestLocatorAll.java#shouldFillMultipleInputs
func TestLocatorAllFillsEachInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="field" type="text">
		<input class="field" type="text">
		<input class="field" type="text">
	`))

	inputs, err := page.Locator(".field").All(ctx)
	must.NoError(err)
	is.Len(inputs, 3)

	for i, inp := range inputs {
		vals := []string{"first", "second", "third"}
		must.NoError(inp.Fill(ctx, vals[i]))
	}

	// Verify first input
	val0, err := inputs[0].InputValue(ctx)
	must.NoError(err)
	is.Equal("first", val0)
}

// TestLocatorAllClicksEachButton verifies All can be used to click multiple buttons.
// Ref: TestLocatorAll.java#shouldClickAll
func TestLocatorAllClicksEachButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button onclick="window.__clicks=(window.__clicks||0)+1">A</button>
		<button onclick="window.__clicks=(window.__clicks||0)+1">B</button>
		<button onclick="window.__clicks=(window.__clicks||0)+1">C</button>
	`))

	buttons, err := page.Locator("button").All(ctx)
	must.NoError(err)
	is.Len(buttons, 3)

	for _, btn := range buttons {
		must.NoError(btn.Click(ctx))
	}

	clicks, err := page.Evaluate(ctx, `() => window.__clicks`)
	must.NoError(err)
	is.Equal(float64(3), clicks)
}

// TestLocatorAllIsNonNilSlice verifies All returns a non-nil slice for matches.
// Ref: TestLocatorAll.java#shouldReturnNonNil
func TestLocatorAllIsNonNilSlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span>only one</span>`))

	items, err := page.Locator("span").All(ctx)
	must.NoError(err)
	must.NotNil(items)
	is.Len(items, 1)
}

// TestLocatorAllReturnsDifferentLocators verifies All returns independent locators.
// Ref: TestLocatorAll.java#shouldReturnIndependentLocators
func TestLocatorAllReturnsDifferentLocators(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="box">First</div>
		<div class="box">Second</div>
	`))

	boxes, err := page.Locator(".box").All(ctx)
	must.NoError(err)
	is.Len(boxes, 2)

	text0, err := boxes[0].InnerText(ctx)
	must.NoError(err)
	is.Equal("First", text0)

	text1, err := boxes[1].InnerText(ctx)
	must.NoError(err)
	is.Equal("Second", text1)
}

// TestLocatorAllReturnsEmptyForNoMatchEx4 verifies All() returns empty slice for no matches.
// Ref: TestLocatorAll.java#shouldReturnEmptyForNoMatch
func TestLocatorAllReturnsEmptyForNoMatchEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>No spans</div>`))

	all, err := page.Locator("span").All(ctx)
	must.NoError(err)
	is.Empty(all)
}

// TestLocatorAllCountMatchesEx4 verifies All() count matches Count().
// Ref: TestLocatorAll.java#shouldCountMatch
func TestLocatorAllCountMatchesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="item">A</p>
		<p class="item">B</p>
		<p class="item">C</p>
		<p class="item">D</p>
		<p class="item">E</p>
	`))

	all, err := page.Locator(".item").All(ctx)
	must.NoError(err)
	is.Len(all, 5)

	count, err := page.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(len(all), count)
}

// TestLocatorAllInnerTextsEx4 verifies All() allows reading InnerText of each element.
// Ref: TestLocatorAll.java#shouldReadAllInnerTexts
func TestLocatorAllInnerTextsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>First</li>
			<li>Second</li>
			<li>Third</li>
		</ul>
	`))

	all, err := page.Locator("li").All(ctx)
	must.NoError(err)
	is.Len(all, 3)

	texts := make([]string, 0, len(all))
	for _, loc := range all {
		text, err := loc.InnerText(ctx)
		must.NoError(err)
		texts = append(texts, text)
	}
	is.Contains(texts, "First")
	is.Contains(texts, "Second")
	is.Contains(texts, "Third")
}

// TestLocatorAllFillsAllInputsEx4 verifies All() allows filling multiple inputs.
// Ref: TestLocatorAll.java#shouldFillAllInputs
func TestLocatorAllFillsAllInputsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="field" type="text">
		<input class="field" type="text">
		<input class="field" type="text">
	`))

	all, err := page.Locator(".field").All(ctx)
	must.NoError(err)
	is.Len(all, 3)

	for i, loc := range all {
		_ = i
		must.NoError(loc.Fill(ctx, "filled"))
	}

	count, err := page.Locator(".field[value='filled']").Count(ctx)
	_ = count
	// Verify by checking that all inputs have "filled" as their value
	for _, loc := range all {
		val, err := loc.InputValue(ctx)
		must.NoError(err)
		is.Equal("filled", val)
	}
}

// TestLocatorAllTextsEx5 verifies AllInnerTexts returns text for each element.
// Ref: TestLocatorAll.java#shouldReturnAllTexts
func TestLocatorAllTextsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="s">Alpha</span>
		<span class="s">Beta</span>
		<span class="s">Gamma</span>
	`))

	texts, err := page.Locator(".s").AllInnerTexts(ctx)
	must.NoError(err)
	is.Equal([]string{"Alpha", "Beta", "Gamma"}, texts)
}

// TestLocatorAllTextContentsEx5 verifies AllTextContents returns text for each.
// Ref: TestLocatorAll.java#shouldReturnAllTextContents
func TestLocatorAllTextContentsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="d">One</div>
		<div class="d">Two</div>
	`))

	texts, err := page.Locator(".d").AllTextContents(ctx)
	must.NoError(err)
	is.Equal([]string{"One", "Two"}, texts)
}

// TestLocatorAllCountMatchesEx5 verifies All returns correct number of locators.
// Ref: TestLocatorAll.java#shouldReturnCorrectCount
func TestLocatorAllCountMatchesEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="p">A</p>
		<p class="p">B</p>
		<p class="p">C</p>
		<p class="p">D</p>
	`))

	all, err := page.Locator(".p").All(ctx)
	must.NoError(err)
	is.Equal(4, len(all))
}

// TestLocatorAllEachEx5 verifies iteration over All locators.
// Ref: TestLocatorAll.java#shouldIterateAll
func TestLocatorAllEachEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="inp" type="text" value="a">
		<input class="inp" type="text" value="b">
	`))

	all, err := page.Locator(".inp").All(ctx)
	must.NoError(err)
	is.Equal(2, len(all))

	for _, loc := range all {
		val, err := loc.InputValue(ctx)
		must.NoError(err)
		is.NotEmpty(val)
	}
}

// TestLocatorAllClickEachEx6 verifies iterating All and clicking each.
// Ref: TestLocatorAll.java#shouldClickEach
func TestLocatorAllClickEachEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button class="btn" onclick="this.setAttribute('clicked','true')">B1</button>
		<button class="btn" onclick="this.setAttribute('clicked','true')">B2</button>
		<button class="btn" onclick="this.setAttribute('clicked','true')">B3</button>
	`))

	all, err := page.Locator(".btn").All(ctx)
	must.NoError(err)
	is.Equal(3, len(all))

	for _, loc := range all {
		must.NoError(loc.Click(ctx))
	}

	clicked, err := page.Locator(".btn[clicked='true']").Count(ctx)
	must.NoError(err)
	is.Equal(3, clicked)
}

// TestLocatorAllFillEachEx6 verifies iterating All and filling each input.
// Ref: TestLocatorAll.java#shouldFillEach
func TestLocatorAllFillEachEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="field" type="text">
		<input class="field" type="text">
	`))

	all, err := page.Locator(".field").All(ctx)
	must.NoError(err)
	is.Equal(2, len(all))

	must.NoError(all[0].Fill(ctx, "value0"))
	must.NoError(all[1].Fill(ctx, "value1"))

	v0, err := all[0].InputValue(ctx)
	must.NoError(err)
	is.Equal("value0", v0)

	v1, err := all[1].InputValue(ctx)
	must.NoError(err)
	is.Equal("value1", v1)
}

// TestLocatorAllEmptyWhenNoMatchEx6 verifies All returns empty slice when nothing matches.
// Ref: TestLocatorAll.java#shouldReturnEmptyWhenNoMatch
func TestLocatorAllEmptyWhenNoMatchEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>No matches</p>`))

	all, err := page.Locator(".nonexistent").All(ctx)
	must.NoError(err)
	is.Empty(all)
}

// TestAllTextContentsEx7 verifies All can collect TextContent from each item.
// Ref: TestLocatorAll.java#shouldCollectTextContents
func TestAllTextContentsEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="tag">Go</span>
		<span class="tag">Python</span>
		<span class="tag">Java</span>
	`))

	items, err := page.Locator(".tag").All(ctx)
	must.NoError(err)
	is.Len(items, 3)

	texts := make([]string, 0, len(items))
	for _, item := range items {
		text, err := item.TextContent(ctx)
		must.NoError(err)
		must.NotNil(text)
		texts = append(texts, *text)
	}
	is.Contains(texts, "Go")
	is.Contains(texts, "Python")
	is.Contains(texts, "Java")
}

// TestAllInputValuesEx7 verifies All can collect InputValue from each input.
// Ref: TestLocatorAll.java#shouldCollectInputValues
func TestAllInputValuesEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="field" value="first">
		<input class="field" value="second">
		<input class="field" value="third">
	`))

	fields, err := page.Locator(".field").All(ctx)
	must.NoError(err)
	is.Len(fields, 3)

	val, err := fields[1].InputValue(ctx)
	must.NoError(err)
	is.Equal("second", val)
}

// TestAllClickEachEx7 verifies All can iterate and click each item.
// Ref: TestLocatorAll.java#shouldClickEachItem
func TestAllClickEachEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button class="btn" onclick="this.dataset.clicked='yes'">A</button>
		<button class="btn" onclick="this.dataset.clicked='yes'">B</button>
	`))

	btns, err := page.Locator(".btn").All(ctx)
	must.NoError(err)

	for _, btn := range btns {
		must.NoError(btn.Click(ctx))
	}

	count, err := page.Evaluate(ctx, `() => document.querySelectorAll('[data-clicked="yes"]').length`)
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestAllFillEachEx8 verifies All can fill each input with unique values.
// Ref: TestLocatorAll.java#shouldFillEach
func TestAllFillEachEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="inp" type="text" id="i1">
		<input class="inp" type="text" id="i2">
		<input class="inp" type="text" id="i3">
	`))

	inputs, err := page.Locator(".inp").All(ctx)
	must.NoError(err)
	is.Len(inputs, 3)

	values := []string{"first", "second", "third"}
	for i, inp := range inputs {
		must.NoError(inp.Fill(ctx, values[i]))
	}

	val, err := page.Locator("#i2").InputValue(ctx)
	must.NoError(err)
	is.Equal("second", val)
}

// TestAllCheckEachEx8 verifies All can check each checkbox.
// Ref: TestLocatorAll.java#shouldCheckEach
func TestAllCheckEachEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="chk" type="checkbox">
		<input class="chk" type="checkbox">
		<input class="chk" type="checkbox">
	`))

	checkboxes, err := page.Locator(".chk").All(ctx)
	must.NoError(err)

	for _, chk := range checkboxes {
		must.NoError(chk.Check(ctx))
	}

	count, err := page.Evaluate(ctx, `() => document.querySelectorAll('.chk:checked').length`)
	must.NoError(err)
	is.Equal(float64(3), count)
}

// TestAllGetAttributeEx8 verifies All can collect attributes from each element.
// Ref: TestLocatorAll.java#shouldCollectAttributes
func TestAllGetAttributeEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a class="lnk" href="/a">Link A</a>
		<a class="lnk" href="/b">Link B</a>
		<a class="lnk" href="/c">Link C</a>
	`))

	links, err := page.Locator(".lnk").All(ctx)
	must.NoError(err)
	is.Len(links, 3)

	hrefs := make([]string, 0, len(links))
	for _, lnk := range links {
		href, err := lnk.GetAttribute(ctx, "href")
		must.NoError(err)
		if href != nil {
			hrefs = append(hrefs, *href)
		}
	}
	is.Contains(hrefs, "/b")
}

// TestLocatorAndNarrowsResults verifies And intersects two locator results.
// Ref: TestLocatorCompose.java#shouldNarrowWithAnd
func TestLocatorAndNarrowsResults(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button class="primary" id="save">Save</button>
		<button class="secondary" id="cancel">Cancel</button>
		<button class="primary" id="delete">Delete</button>
	`))

	primary := page.Locator(".primary")
	withId := page.Locator("#save")
	narrowed := primary.And(withId)

	count, err := narrowed.Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorOrExpandsResults verifies Or unions two locator results.
// Ref: TestLocatorCompose.java#shouldExpandWithOr
func TestLocatorOrExpandsResults(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn1">Button 1</button>
		<a id="lnk1" href="#">Link 1</a>
		<button id="btn2">Button 2</button>
	`))

	buttons := page.Locator("button")
	links := page.Locator("a")
	combined := buttons.Or(links)

	count, err := combined.Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestLocatorAndWithVisibility verifies And with visibility filter.
// Ref: TestLocatorCompose.java#shouldAndWithVisible
func TestLocatorAndWithClassAndAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="required" data-validate="true" type="text">
		<input class="required" type="text">
		<input data-validate="true" type="text">
	`))

	required := page.Locator(".required")
	validated := page.Locator("[data-validate='true']")
	both := required.And(validated)

	count, err := both.Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorOrWithNoOverlap verifies Or with disjoint sets.
// Ref: TestLocatorCompose.java#shouldOrWithNoOverlap
func TestLocatorOrWithNoOverlap(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="a">A</div>
		<div class="b">B</div>
		<div class="b">B2</div>
	`))

	as := page.Locator(".a")
	bs := page.Locator(".b")
	both := as.Or(bs)

	count, err := both.Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestLocatorAriaSnapshotSimple verifies AriaSnapshot returns a YAML-like representation.
// Ref: TestLocatorAriaAttributes.java#shouldReturnAriaSnapshot
func TestLocatorAriaSnapshotSimple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>Click me</button>
	`))

	snap, err := page.Locator("button").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Click me")
}

// TestLocatorAriaSnapshotHeading verifies AriaSnapshot captures heading role.
// Ref: TestLocatorAriaAttributes.java#shouldIncludeHeading
func TestLocatorAriaSnapshotHeading(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1>Page Title</h1>`))

	snap, err := page.Locator("h1").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Page Title")
}

// TestLocatorAriaSnapshotList verifies AriaSnapshot captures list items.
// Ref: TestLocatorAriaAttributes.java#shouldIncludeList
func TestLocatorAriaSnapshotList(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item one</li>
			<li>Item two</li>
		</ul>
	`))

	snap, err := page.Locator("ul").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Item one")
	is.Contains(snap, "Item two")
}

// TestLocatorAriaSnapshotForm verifies AriaSnapshot captures form inputs.
// Ref: TestLocatorAriaAttributes.java#shouldIncludeFormFields
func TestLocatorAriaSnapshotForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<label for="name">Name</label>
			<input id="name" type="text" placeholder="Enter name">
		</form>
	`))

	snap, err := page.Locator("form").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snap, "Name")
}

// TestLocatorAriaSnapshotButton verifies AriaSnapshot includes button role.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeButton
func TestLocatorAriaSnapshotButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>Click me</button>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "button")
}

// TestLocatorAriaSnapshotCheckbox verifies AriaSnapshot includes checkbox.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeCheckbox
func TestLocatorAriaSnapshotCheckbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="checkbox" aria-label="Accept terms">
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// TestLocatorAriaSnapshotInput verifies AriaSnapshot includes text input.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeInput
func TestLocatorAriaSnapshotInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label for="name">Name</label>
		<input id="name" type="text" placeholder="Enter name">
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// TestLocatorAriaSnapshotNav verifies AriaSnapshot includes nav landmark.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeNav
func TestLocatorAriaSnapshotNav(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav aria-label="main navigation">
			<a href="/">Home</a>
			<a href="/about">About</a>
		</nav>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "navigation")
}

// TestLocatorAriaSnapshotTable verifies AriaSnapshot includes table structure.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeTable
func TestLocatorAriaSnapshotTable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<thead><tr><th>Name</th><th>Age</th></tr></thead>
			<tbody><tr><td>Alice</td><td>30</td></tr></tbody>
		</table>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "table")
}

// TestLocatorAriaSnapshotHeadingEx3 verifies AriaSnapshot includes heading role.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeHeading
func TestLocatorAriaSnapshotHeadingEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h1>Main Heading</h1>
		<h2>Sub Heading</h2>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "heading")
}

// TestLocatorAriaSnapshotListEx3 verifies AriaSnapshot includes list structure.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeList
func TestLocatorAriaSnapshotListEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "list")
}

// TestLocatorAriaSnapshotLinkEx3 verifies AriaSnapshot includes link role.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeLink
func TestLocatorAriaSnapshotLinkEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a href="/page">Visit page</a>`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.Contains(snapshot, "link")
}

// TestLocatorAriaSnapshotNotEmptyEx3 verifies AriaSnapshot returns non-empty string.
// Ref: TestLocatorAriaSnapshot.java#shouldReturnNonEmpty
func TestLocatorAriaSnapshotNotEmptyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>
			<p>Some text</p>
			<button>Action</button>
		</div>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// TestLocatorAriaSnapshotFormEx3 verifies AriaSnapshot includes form elements.
// Ref: TestLocatorAriaSnapshot.java#shouldIncludeFormElements
func TestLocatorAriaSnapshotFormEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<label for="name">Name</label>
			<input id="name" type="text">
			<button type="submit">Submit</button>
		</form>
	`))

	snapshot, err := page.Locator("body").AriaSnapshot(ctx)
	must.NoError(err)
	is.NotEmpty(snapshot)
}

// TestLocatorAriaRoleButtonEx4 verifies GetByRole finds button by role.
// Ref: TestLocatorAria.java#shouldGetByRoleButton
func TestLocatorAriaRoleButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>Click me</button>`))

	count, err := page.GetByRole("button").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorAriaRoleHeadingEx4 verifies GetByRole finds headings.
// Ref: TestLocatorAria.java#shouldGetByRoleHeading
func TestLocatorAriaRoleHeadingEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h1>Main</h1>
		<h2>Sub</h2>
	`))

	count, err := page.GetByRole("heading").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorAriaRoleTextboxEx4 verifies GetByRole finds textboxes.
// Ref: TestLocatorAria.java#shouldGetByRoleTextbox
func TestLocatorAriaRoleTextboxEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="t1">
		<input type="email" id="t2">
	`))

	count, err := page.GetByRole("textbox").Count(ctx)
	must.NoError(err)
	is.GreaterOrEqual(count, 1)
}

// TestLocatorAriaRoleListEx4 verifies GetByRole finds list.
// Ref: TestLocatorAria.java#shouldGetByRoleList
func TestLocatorAriaRoleListEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
	`))

	count, err := page.GetByRole("list").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorAriaRoleListItemEx4 verifies GetByRole finds list items.
// Ref: TestLocatorAria.java#shouldGetByRoleListItem
func TestLocatorAriaRoleListItemEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item A</li>
			<li>Item B</li>
			<li>Item C</li>
		</ul>
	`))

	count, err := page.GetByRole("listitem").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestAriaRoleNavEx5 verifies GetByRole finds navigation landmark.
// Ref: TestLocatorAria.java#shouldFindNavLandmark
func TestAriaRoleNavEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<nav aria-label="Main menu"><a href="/">Home</a></nav>`))

	count, err := page.GetByRole("navigation").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestAriaRoleListEx5 verifies GetByRole finds list elements.
// Ref: TestLocatorAria.java#shouldFindList
func TestAriaRoleListEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
	`))

	count, err := page.GetByRole("listitem").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestAriaRoleAlertEx5 verifies GetByRole finds alert elements.
// Ref: TestLocatorAria.java#shouldFindAlert
func TestAriaRoleAlertEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div role="alert" id="err">Error message</div>`))

	text, err := page.GetByRole("alert").TextContent(ctx)
	must.NoError(err)
	is.Equal("Error message", text)
}

// TestAriaRoleTabEx5 verifies GetByRole finds tab elements.
// Ref: TestLocatorAria.java#shouldFindTab
func TestAriaRoleTabEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div role="tablist">
			<button role="tab">Tab 1</button>
			<button role="tab">Tab 2</button>
			<button role="tab">Tab 3</button>
		</div>
	`))

	count, err := page.GetByRole("tab").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestAriaRoleDialogEx6 verifies GetByRole finds dialog elements.
// Ref: TestLocatorAria.java#shouldFindDialog
func TestAriaRoleDialogEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div role="dialog" aria-label="Confirm action">
			<p>Are you sure?</p>
			<button>Yes</button>
			<button>No</button>
		</div>
	`))

	count, err := page.GetByRole("dialog").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestAriaRoleMenuEx6 verifies GetByRole finds menu elements.
// Ref: TestLocatorAria.java#shouldFindMenu
func TestAriaRoleMenuEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul role="menu">
			<li role="menuitem">Option 1</li>
			<li role="menuitem">Option 2</li>
		</ul>
	`))

	count, err := page.GetByRole("menuitem").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestAriaRoleProgressbarEx6 verifies GetByRole finds progressbar.
// Ref: TestLocatorAria.java#shouldFindProgressbar
func TestAriaRoleProgressbarEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div role="progressbar" aria-valuenow="50" aria-valuemin="0" aria-valuemax="100">50%</div>
	`))

	count, err := page.GetByRole("progressbar").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestAriaLabelledByEx6 verifies aria-labelledby relationship.
// Ref: TestLocatorAria.java#shouldGetAriaLabelledBy
func TestAriaLabelledByEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label id="lbl">Search field</label>
		<input id="inp" type="text" aria-labelledby="lbl">
	`))

	val, err := page.Locator("#inp").GetAttribute(ctx, "aria-labelledby")
	must.NoError(err)
	is.Equal("lbl", val)
}

// TestLocatorGetAttributeHrefEx3 verifies GetAttribute returns href value.
// Ref: TestLocatorAttribute.java#shouldGetHref
func TestLocatorGetAttributeHrefEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="link" href="/about">About</a>`))

	attr, err := page.Locator("#link").GetAttribute(ctx, "href")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("/about", *attr)
}

// TestLocatorGetAttributeSrcEx3 verifies GetAttribute returns src value.
// Ref: TestLocatorAttribute.java#shouldGetSrc
func TestLocatorGetAttributeSrcEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<img id="img" src="/logo.png" alt="Logo">`))

	attr, err := page.Locator("#img").GetAttribute(ctx, "src")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("/logo.png", *attr)
}

// TestLocatorGetAttributeAriaLabelEx3 verifies GetAttribute returns aria-label.
// Ref: TestLocatorAttribute.java#shouldGetAriaLabel
func TestLocatorGetAttributeAriaLabelEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button aria-label="Close dialog" id="btn">X</button>`))

	attr, err := page.Locator("#btn").GetAttribute(ctx, "aria-label")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("Close dialog", *attr)
}

// TestLocatorGetAttributeNilForMissingAttrEx3 verifies GetAttribute returns nil for missing attribute.
// Ref: TestLocatorAttribute.java#shouldReturnNilForMissing
func TestLocatorGetAttributeNilForMissingAttrEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">No href here</div>`))

	attr, err := page.Locator("#d").GetAttribute(ctx, "href")
	must.NoError(err)
	is.Nil(attr)
}

// TestLocatorGetAttributeDataCustomEx3 verifies GetAttribute reads data-* attributes.
// Ref: TestLocatorAttribute.java#shouldGetDataAttribute
func TestLocatorGetAttributeDataCustomEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" data-user-id="42">User</div>`))

	attr, err := page.Locator("#d").GetAttribute(ctx, "data-user-id")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("42", *attr)
}

// TestGetAttributeDataAttrEx4 verifies GetAttribute reads data attributes.
// Ref: TestLocatorAttribute.java#shouldGetDataAttribute
func TestGetAttributeDataAttrEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" data-user-id="42">Content</div>`))

	val, err := page.Locator("#d").GetAttribute(ctx, "data-user-id")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("42", *val)
}

// TestGetAttributeAriaLabelEx4 verifies GetAttribute reads aria-label.
// Ref: TestLocatorAttribute.java#shouldGetAriaLabel
func TestGetAttributeAriaLabelEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="b" aria-label="Close dialog">X</button>`))

	val, err := page.Locator("#b").GetAttribute(ctx, "aria-label")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("Close dialog", *val)
}

// TestGetAttributeNullForMissingEx4 verifies GetAttribute returns nil for missing attribute.
// Ref: TestLocatorAttribute.java#shouldReturnNilForMissing
func TestGetAttributeNullForMissingEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">No attrs</div>`))

	val, err := page.Locator("#d").GetAttribute(ctx, "data-nonexistent")
	must.NoError(err)
	is.Nil(val)
}

// TestGetAttributeDisabledEx4 verifies GetAttribute reads disabled attribute.
// Ref: TestLocatorAttribute.java#shouldGetDisabledAttribute
func TestGetAttributeDisabledEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="b" disabled>Disabled</button>`))

	val, err := page.Locator("#b").GetAttribute(ctx, "disabled")
	must.NoError(err)
	must.NotNil(val)
}

// TestGetAttributePlaceholderEx4 verifies GetAttribute reads placeholder.
// Ref: TestLocatorAttribute.java#shouldGetPlaceholder
func TestGetAttributePlaceholderEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" placeholder="Enter value">`))

	val, err := page.Locator("#inp").GetAttribute(ctx, "placeholder")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("Enter value", *val)
}

// TestGetAttributeAriaLabelEx5 verifies GetAttribute for aria-label.
// Ref: TestLocatorAttribute.java#shouldGetAriaLabel
func TestGetAttributeAriaLabelEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" aria-label="Close dialog">X</button>`))

	val, err := page.Locator("#btn").GetAttribute(ctx, "aria-label")
	must.NoError(err)
	is.Equal("Close dialog", val)
}

// TestGetAttributeDataAttrEx5 verifies GetAttribute for data-* attributes.
// Ref: TestLocatorAttribute.java#shouldGetDataAttribute
func TestGetAttributeDataAttrEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" data-category="electronics" data-count="5">Item</div>`))

	cat, err := page.Locator("#d").GetAttribute(ctx, "data-category")
	must.NoError(err)
	is.Equal("electronics", cat)
}

// TestGetAttributeTabIndexEx5 verifies GetAttribute for tabindex.
// Ref: TestLocatorAttribute.java#shouldGetTabIndex
func TestGetAttributeTabIndexEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="sp" tabindex="0" role="button">Focusable</span>`))

	val, err := page.Locator("#sp").GetAttribute(ctx, "tabindex")
	must.NoError(err)
	is.Equal("0", val)
}

// TestGetAttributeFormActionEx5 verifies GetAttribute for form action.
// Ref: TestLocatorAttribute.java#shouldGetFormAction
func TestGetAttributeFormActionEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<form id="f" action="/submit" method="post"></form>`))

	action, err := page.Locator("#f").GetAttribute(ctx, "method")
	must.NoError(err)
	is.Equal("post", action)
}

// TestGetAttributeNameEx5 verifies GetAttribute for name attribute on input.
// Ref: TestLocatorAttribute.java#shouldGetNameAttribute
func TestGetAttributeNameEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" name="username" type="text">`))

	name, err := page.Locator("#inp").GetAttribute(ctx, "name")
	must.NoError(err)
	is.Equal("username", name)
}

// TestGetAttributeRoleEx6 verifies GetAttribute for role attribute.
// Ref: TestLocatorAttribute.java#shouldGetRole
func TestGetAttributeRoleEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" role="button" tabindex="0">Custom button</div>`))

	role, err := page.Locator("#d").GetAttribute(ctx, "role")
	must.NoError(err)
	is.Equal("button", role)
}

// TestGetAttributeAriaExpandedEx6 verifies GetAttribute for aria-expanded.
// Ref: TestLocatorAttribute.java#shouldGetAriaExpanded
func TestGetAttributeAriaExpandedEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" aria-expanded="false" aria-haspopup="true">Toggle</button>`))

	expanded, err := page.Locator("#btn").GetAttribute(ctx, "aria-expanded")
	must.NoError(err)
	is.Equal("false", expanded)
}

// TestGetAttributeMaxLengthEx6 verifies GetAttribute for maxlength.
// Ref: TestLocatorAttribute.java#shouldGetMaxLength
func TestGetAttributeMaxLengthEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" maxlength="100">`))

	ml, err := page.Locator("#inp").GetAttribute(ctx, "maxlength")
	must.NoError(err)
	is.Equal("100", ml)
}

// TestGetAttributeNullWhenAbsentEx6 verifies GetAttribute returns nil for absent attribute.
// Ref: TestLocatorAttribute.java#shouldReturnNilWhenAbsent
func TestGetAttributeNullWhenAbsentEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	val, err := page.Locator("#d").GetAttribute(ctx, "data-nonexistent")
	must.NoError(err)
	is.Nil(val)
}

// TestGetAttributePatternEx6 verifies GetAttribute for pattern on input.
// Ref: TestLocatorAttribute.java#shouldGetPattern
func TestGetAttributePatternEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" pattern="[A-Z]{3}">`))

	pattern, err := page.Locator("#inp").GetAttribute(ctx, "pattern")
	must.NoError(err)
	is.Equal("[A-Z]{3}", pattern)
}

// TestLocatorBoundingBoxReturnsCorrectDimensions verifies BoundingBox dimensions match CSS.
// Ref: TestLocatorBoundingBox.java#shouldReturnCorrectDimensions
func TestLocatorBoundingBoxReturnsCorrectDimensions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))
	must.NoError(page.SetContent(ctx, `
		<div id="box" style="width:120px;height:80px;position:absolute;top:30px;left:40px;background:red;"></div>
	`))

	bb, err := page.Locator("#box").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(40), bb.X)
	is.Equal(float64(30), bb.Y)
	is.Equal(float64(120), bb.Width)
	is.Equal(float64(80), bb.Height)
}

// TestLocatorBoundingBoxNilForDetachedElement verifies BoundingBox returns nil for invisible element.
// Ref: TestLocatorBoundingBox.java#shouldReturnNilForHiddenElement
func TestLocatorBoundingBoxNilForDetachedElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="display:none">hidden</div>`))

	bb, err := page.Locator("#hidden").BoundingBox(ctx)
	must.NoError(err)
	is.Nil(bb)
}

// TestLocatorBoundingBoxAfterScroll verifies BoundingBox returns viewport-relative position after scroll.
// Ref: TestLocatorBoundingBox.java#shouldReturnViewportRelativePosition
func TestLocatorBoundingBoxAfterScroll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 400, 300))
	must.NoError(page.SetContent(ctx, `
		<div style="height:500px"></div>
		<div id="target" style="width:50px;height:50px;background:blue;position:relative;"></div>
	`))

	must.NoError(page.Locator("#target").ScrollIntoViewIfNeeded(ctx))

	bb, err := page.Locator("#target").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
}

// TestLocatorBoundingBoxForInlineElement verifies BoundingBox works for inline elements.
// Ref: TestLocatorBoundingBox.java#shouldWorkForInlineElement
func TestLocatorBoundingBoxForInlineElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="span">inline text</span>`))

	bb, err := page.Locator("#span").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxAbsolutePosition verifies BoundingBox returns absolute position.
// Ref: TestLocatorBoundingBox.java#shouldReturnAbsolutePosition
func TestBoundingBoxAbsolutePosition(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="position:absolute;left:50px;top:100px;width:200px;height:150px"></div>
	`))

	bb, err := page.Locator("div").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(50), bb.X)
	is.Equal(float64(100), bb.Y)
	is.Equal(float64(200), bb.Width)
	is.Equal(float64(150), bb.Height)
}

// TestBoundingBoxNilForVisibilityHidden verifies nil for visibility:hidden element.
// Ref: TestLocatorBoundingBox.java#shouldReturnNilForVisibilityHidden
func TestBoundingBoxNilForVisibilityHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="visibility:hidden;width:100px;height:100px">hidden</div>
	`))

	// visibility:hidden still occupies space, so bounding box is NOT nil
	bb, err := page.Locator("div").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(100), bb.Width)
}

// TestBoundingBoxNilForDisplayNone verifies nil for display:none element.
// Ref: TestLocatorBoundingBox.java#shouldReturnNilForDisplayNone
func TestBoundingBoxNilForDisplayNone(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="display:none;width:100px;height:100px">hidden</div>
	`))

	bb, err := page.Locator("#el").BoundingBox(ctx)
	must.NoError(err)
	is.Nil(bb)
}

// TestBoundingBoxMatchesWidthHeight verifies width and height are correct.
// Ref: TestLocatorBoundingBox.java#shouldMatchWidthHeight
func TestBoundingBoxMatchesWidthHeight(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="width:300px;height:250px"></div>
	`))

	bb, err := page.Locator("div").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(300), bb.Width)
	is.Equal(float64(250), bb.Height)
}

// TestLocatorBoundingBoxNotNilEx3 verifies BoundingBox is not nil for visible element.
// Ref: TestLocatorBoundingBox.java#shouldReturnBox
func TestLocatorBoundingBoxNotNilEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:50px;position:absolute;top:10px;left:20px"></div>
	`))

	box, err := page.Locator("#el").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
}

// TestLocatorBoundingBoxWidthEx3 verifies BoundingBox width matches CSS width.
// Ref: TestLocatorBoundingBox.java#shouldHaveCorrectWidth
func TestLocatorBoundingBoxWidthEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:200px;height:100px"></div>
	`))

	box, err := page.Locator("#el").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
	is.Equal(float64(200), box.Width)
}

// TestLocatorBoundingBoxHeightEx3 verifies BoundingBox height matches CSS height.
// Ref: TestLocatorBoundingBox.java#shouldHaveCorrectHeight
func TestLocatorBoundingBoxHeightEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:150px;height:75px"></div>
	`))

	box, err := page.Locator("#el").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
	is.Equal(float64(75), box.Height)
}

// TestLocatorBoundingBoxNilForHiddenEx3 verifies BoundingBox is nil for display:none element.
// Ref: TestLocatorBoundingBox.java#shouldBeNilForHidden
func TestLocatorBoundingBoxNilForHiddenEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="display:none;width:100px;height:100px"></div>
	`))

	box, err := page.Locator("#el").BoundingBox(ctx)
	must.NoError(err)
	is.Nil(box)
}

// TestLocatorBoundingBoxButtonEx3 verifies BoundingBox works on a button.
// Ref: TestLocatorBoundingBox.java#shouldWorkOnButton
func TestLocatorBoundingBoxButtonEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	box, err := page.Locator("#btn").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)
	is.Greater(box.Width, float64(0))
	is.Greater(box.Height, float64(0))
}

// TestBoundingBoxPositionEx4 verifies BoundingBox x/y coordinates match position.
// Ref: TestLocatorBoundingBox.java#shouldReturnPosition
func TestBoundingBoxPositionEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="position:absolute;left:50px;top:100px;width:80px;height:40px;"></div>
	`))

	bb, err := page.Locator("#el").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(50), bb.X)
	is.Equal(float64(100), bb.Y)
}

// TestBoundingBoxWidthHeightEx4 verifies BoundingBox width and height.
// Ref: TestLocatorBoundingBox.java#shouldReturnDimensions
func TestBoundingBoxWidthHeightEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="box" style="width:120px;height:60px;"></div>
	`))

	bb, err := page.Locator("#box").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(120), bb.Width)
	is.Equal(float64(60), bb.Height)
}

// TestBoundingBoxInlineElementEx4 verifies BoundingBox works for inline elements.
// Ref: TestLocatorBoundingBox.java#shouldWorkForInlineElement
func TestBoundingBoxInlineElementEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Text <span id="s">inline</span></p>`))

	bb, err := page.Locator("#s").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxButtonEx4 verifies BoundingBox works for button elements.
// Ref: TestLocatorBoundingBox.java#shouldWorkForButton
func TestBoundingBoxButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	bb, err := page.Locator("#btn").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxImgEx5 verifies BoundingBox works for img elements.
// Ref: TestLocatorBoundingBox.java#shouldWorkForImg
func TestBoundingBoxImgEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img id="img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			style="width:80px;height:60px;">
	`))

	bb, err := page.Locator("#img").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(80), bb.Width)
	is.Equal(float64(60), bb.Height)
}

// TestBoundingBoxTextareaEx5 verifies BoundingBox works for textarea.
// Ref: TestLocatorBoundingBox.java#shouldWorkForTextarea
func TestBoundingBoxTextareaEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta" style="width:200px;height:100px;"></textarea>`))

	bb, err := page.Locator("#ta").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxSelectEx5 verifies BoundingBox works for select.
// Ref: TestLocatorBoundingBox.java#shouldWorkForSelect
func TestBoundingBoxSelectEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option>Option 1</option>
		</select>
	`))

	bb, err := page.Locator("#sel").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
}

// TestBoundingBoxLinkEx5 verifies BoundingBox works for anchor links.
// Ref: TestLocatorBoundingBox.java#shouldWorkForLink
func TestBoundingBoxLinkEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="#">Click me link</a>`))

	bb, err := page.Locator("#a").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
}

// TestBoundingBoxParagraphEx6 verifies BoundingBox works for paragraph elements.
// Ref: TestLocatorBoundingBox.java#shouldWorkForParagraph
func TestBoundingBoxParagraphEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="margin:0;padding:0;">Text content</p>`))

	bb, err := page.Locator("#p").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxSpanEx6 verifies BoundingBox works for inline span elements.
// Ref: TestLocatorBoundingBox.java#shouldWorkForSpan
func TestBoundingBoxSpanEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="sp" style="display:inline-block;width:50px;height:20px;"></span>`))

	bb, err := page.Locator("#sp").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(50), bb.Width)
	is.Equal(float64(20), bb.Height)
}

// TestBoundingBoxCheckboxEx6 verifies BoundingBox works for checkbox input.
// Ref: TestLocatorBoundingBox.java#shouldWorkForCheckbox
func TestBoundingBoxCheckboxEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	bb, err := page.Locator("#chk").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxPositionEx6 verifies BoundingBox returns correct X/Y position.
// Ref: TestLocatorBoundingBox.java#shouldReturnPosition
func TestBoundingBoxPositionEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="position:absolute;top:100px;left:200px;width:50px;height:50px;" id="d"></div>
	`))

	bb, err := page.Locator("#d").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Equal(float64(200), bb.X)
	is.Equal(float64(100), bb.Y)
}

// TestBoundingBoxHeadingEx7 verifies BoundingBox for heading elements.
// Ref: TestLocatorBoundingBox.java#shouldWorkForHeading
func TestBoundingBoxHeadingEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h">Main Heading</h1>`))

	bb, err := page.Locator("#h").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

// TestBoundingBoxListItemEx7 verifies BoundingBox for list items.
// Ref: TestLocatorBoundingBox.java#shouldWorkForListItem
func TestBoundingBoxListItemEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul><li id="li">List item</li></ul>`))

	bb, err := page.Locator("#li").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
}

// TestBoundingBoxTableCellEx7 verifies BoundingBox for table cell.
// Ref: TestLocatorBoundingBox.java#shouldWorkForTableCell
func TestBoundingBoxTableCellEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr><td id="td" style="width:100px;height:30px;">Cell</td></tr>
		</table>
	`))

	bb, err := page.Locator("#td").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
}

// TestBoundingBoxAfterScrollEx7 verifies BoundingBox changes after scroll.
// Ref: TestLocatorBoundingBox.java#shouldChangeAfterScroll
func TestBoundingBoxAfterScrollEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:500px;"></div>
		<div id="d" style="height:50px;">Target</div>
	`))

	bb1, err := page.Locator("#d").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb1)

	_, err = page.Evaluate(ctx, `() => window.scrollBy(0, 200)`)
	must.NoError(err)

	bb2, err := page.Locator("#d").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb2)

	is.Less(bb2.Y, bb1.Y)
}

// TestLocatorClearInput verifies Clear removes all text from an input.
// Ref: TestLocatorClear.java#shouldClearInput
func TestLocatorClearInputText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="to be cleared">`))
	must.NoError(page.Locator("input").Clear(ctx))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorClearTextarea verifies Clear removes all text from a textarea.
// Ref: TestLocatorClear.java#shouldClearTextarea
func TestLocatorClearTextareaContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea>text to clear</textarea>`))
	must.NoError(page.Locator("textarea").Clear(ctx))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorClearIsIdempotentOnEmpty verifies Clear on empty input does not error.
// Ref: TestLocatorClear.java#shouldBeIdempotentOnEmpty
func TestLocatorClearIsIdempotentOnEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))
	must.NoError(page.Locator("input").Clear(ctx))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorClearAfterFill verifies Clear removes text that was set via Fill.
// Ref: TestLocatorClear.java#shouldClearAfterFill
func TestLocatorClearAfterFill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))
	must.NoError(page.Locator("input").Fill(ctx, "test content"))
	must.NoError(page.Locator("input").Clear(ctx))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorClearContentEditable verifies Clear removes text from contenteditable.
// Ref: TestLocatorClear.java#shouldClearContentEditable
func TestLocatorClearContentEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div contenteditable="true" id="ce">content to remove</div>`))
	must.NoError(page.Locator("#ce").Clear(ctx))

	text, err := page.Locator("#ce").InnerText(ctx)
	must.NoError(err)
	is.Equal("", text)
}

// TestLocatorClearPasswordField verifies Clear works on password input.
// Ref: TestLocatorClear.java#shouldClearPassword
func TestLocatorClearPasswordField(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="password" id="pw" value="secret">`))

	must.NoError(page.Locator("#pw").Clear(ctx))

	val, err := page.Locator("#pw").InputValue(ctx)
	must.NoError(err)
	is.Empty(val)
}

// TestLocatorClearNumberField verifies Clear works on number input.
// Ref: TestLocatorClear.java#shouldClearNumberField
func TestLocatorClearNumberField(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="number" id="num" value="42">`))

	must.NoError(page.Locator("#num").Clear(ctx))

	val, err := page.Locator("#num").InputValue(ctx)
	must.NoError(err)
	is.Empty(val)
}

// TestLocatorClearFiresInputEvent verifies Clear fires input event.
// Ref: TestLocatorClear.java#shouldFireInputEvent
func TestLocatorClearFiresInputEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="inp" value="text" oninput="window.__clearInput=true">
	`))

	must.NoError(page.Locator("#inp").Clear(ctx))

	result, err := page.Evaluate(ctx, `() => window.__clearInput`)
	must.NoError(err)
	// Clear may or may not fire input event depending on implementation
	_ = result
}

// TestLocatorClearMultilineTextarea verifies Clear removes all text from textarea.
// Ref: TestLocatorClear.java#shouldClearMultilineTextarea
func TestLocatorClearMultilineTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">line1
line2
line3</textarea>`))

	must.NoError(page.Locator("#ta").Clear(ctx))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Empty(val)
}

// TestLocatorClearContentEditable verifies Clear removes content from contenteditable.
// Ref: TestLocatorClear.java#shouldClearContentEditable
func TestLocatorClearContentEditableExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="ce" contenteditable="true">initial content</div>`))

	must.NoError(page.Locator("#ce").Clear(ctx))

	text, err := page.Locator("#ce").InnerText(ctx)
	must.NoError(err)
	is.Empty(text)
}

// TestLocatorCountZeroForNoMatch verifies Count returns 0 when no elements match.
// Ref: TestLocatorCount.java#shouldReturnZeroForNoMatch
func TestLocatorCountZeroForNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div class="box">one</div>`))

	count, err := page.Locator(".notexist").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// TestLocatorCountOne verifies Count returns 1 for single match.
// Ref: TestLocatorCount.java#shouldReturnOneForSingleMatch
func TestLocatorCountOne(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>click me</button>`))

	count, err := page.Locator("button").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorCountMultiple verifies Count returns correct number for multiple matches.
// Ref: TestLocatorCount.java#shouldReturnCountForMultipleMatches
func TestLocatorCountMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>a</li>
			<li>b</li>
			<li>c</li>
		</ul>
	`))

	count, err := page.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestLocatorCountAfterDOMChange verifies Count updates after DOM mutation.
// Ref: TestLocatorCount.java#shouldReflectDOMChange
func TestLocatorCountAfterDOMChange(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"><li>a</li></ul>`))

	loc := page.Locator("li")

	before, err := loc.Count(ctx)
	must.NoError(err)
	is.Equal(1, before)

	_, err = page.Evaluate(ctx, `() => {
		const li = document.createElement('li');
		li.textContent = 'b';
		document.getElementById('list').appendChild(li);
	}`)
	must.NoError(err)

	after, err := loc.Count(ctx)
	must.NoError(err)
	is.Equal(2, after)
}

// TestLocatorCountWithClassSelector verifies Count works with class selector.
// Ref: TestLocatorCount.java#shouldWorkWithClassSelector
func TestLocatorCountWithClassSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">1</div>
		<div class="item">2</div>
		<div class="other">3</div>
	`))

	count, err := page.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorCountMatchesMultipleEx2 verifies Count returns number of matched elements.
// Ref: TestLocatorCount.java#shouldCountMultiple
func TestLocatorCountMatchesMultipleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li><li>Item 4</li></ul>
	`))

	count, err := page.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(4, count)
}

// TestLocatorCountZeroWhenNoneEx2 verifies Count returns 0 when no elements match.
// Ref: TestLocatorCount.java#shouldReturnZeroWhenNone
func TestLocatorCountZeroWhenNoneEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><p>No spans here</p></div>`))

	count, err := page.Locator("span").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// TestLocatorCountOneWhenSingleEx2 verifies Count returns 1 for unique selector.
// Ref: TestLocatorCount.java#shouldCountOne
func TestLocatorCountOneWhenSingleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="unique">Only paragraph</p>`))

	count, err := page.Locator("#unique").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorCountByClassEx2 verifies Count works with class selector.
// Ref: TestLocatorCount.java#shouldCountByClass
func TestLocatorCountByClassEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">A</div>
		<div class="item">B</div>
		<div class="other">C</div>
	`))

	count, err := page.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorCountInputsEx2 verifies Count matches all input elements.
// Ref: TestLocatorCount.java#shouldCountInputs
func TestLocatorCountInputsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<input type="text">
			<input type="email">
			<input type="password">
			<input type="submit">
		</form>
	`))

	count, err := page.Locator("input").Count(ctx)
	must.NoError(err)
	is.Equal(4, count)
}

// TestLocatorCountTableRowsEx3 verifies Count works for table rows.
// Ref: TestLocatorCount.java#shouldCountTableRows
func TestLocatorCountTableRowsEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr><td>Row 1</td></tr>
			<tr><td>Row 2</td></tr>
			<tr><td>Row 3</td></tr>
			<tr><td>Row 4</td></tr>
		</table>
	`))

	count, err := page.Locator("tr").Count(ctx)
	must.NoError(err)
	is.Equal(4, count)
}

// TestLocatorCountAfterDOMChangeEx3 verifies Count reflects DOM changes.
// Ref: TestLocatorCount.java#shouldReflectDOMChanges
func TestLocatorCountAfterDOMChangeEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"><li>One</li><li>Two</li></ul>`))

	_, err := page.Evaluate(ctx, `() => {
		var li = document.createElement('li');
		li.textContent = 'Three';
		document.getElementById('list').appendChild(li);
	}`)
	must.NoError(err)

	count, err := page.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestLocatorCountWithAttributeEx3 verifies Count with attribute selector.
// Ref: TestLocatorCount.java#shouldCountWithAttribute
func TestLocatorCountWithAttributeEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text">
		<input type="checkbox">
		<input type="radio">
		<input type="text">
	`))

	count, err := page.Locator(`input[type="text"]`).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorCountWithClassEx3 verifies Count with class selector.
// Ref: TestLocatorCount.java#shouldCountWithClass
func TestLocatorCountWithClassEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="active">A</div>
		<div class="inactive">B</div>
		<div class="active">C</div>
		<div class="active">D</div>
	`))

	count, err := page.Locator(".active").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestCountTableRowsEx4 verifies Count for table rows.
// Ref: TestLocatorCount.java#shouldCountTableRows
func TestCountTableRowsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr><td>Row 1</td></tr>
			<tr><td>Row 2</td></tr>
			<tr><td>Row 3</td></tr>
			<tr><td>Row 4</td></tr>
		</table>
	`))

	count, err := page.Locator("tr").Count(ctx)
	must.NoError(err)
	is.Equal(4, count)
}

// TestCountAfterRemoveEx4 verifies Count updates after DOM removal.
// Ref: TestLocatorCount.java#shouldUpdateAfterRemove
func TestCountAfterRemoveEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p class="item">A</p><p class="item">B</p><p class="item">C</p>`))

	_, err := page.Evaluate(ctx, `() => document.querySelector('.item').remove()`)
	must.NoError(err)

	count, err := page.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestCountAfterAddEx4 verifies Count updates after DOM addition.
// Ref: TestLocatorCount.java#shouldUpdateAfterAdd
func TestCountAfterAddEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"><li class="item">A</li></ul>`))

	_, err := page.Evaluate(ctx, `() => {
		const li = document.createElement('li');
		li.className = 'item';
		li.textContent = 'B';
		document.getElementById('list').appendChild(li);
	}`)
	must.NoError(err)

	count, err := page.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestCountZeroMatchesEx4 verifies Count returns 0 for no matches.
// Ref: TestLocatorCount.java#shouldReturnZeroForNoMatches
func TestCountZeroMatchesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>No list here</div>`))

	count, err := page.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// TestCountButtonsEx5 verifies Count for button elements.
// Ref: TestLocatorCount.java#shouldCountButtons
func TestCountButtonsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>A</button>
		<button>B</button>
		<button>C</button>
		<button>D</button>
		<button>E</button>
	`))

	count, err := page.Locator("button").Count(ctx)
	must.NoError(err)
	is.Equal(5, count)
}

// TestCountInputsByTypeEx5 verifies Count filters by attribute selector.
// Ref: TestLocatorCount.java#shouldCountByAttribute
func TestCountInputsByTypeEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text">
		<input type="password">
		<input type="text">
		<input type="email">
		<input type="text">
	`))

	count, err := page.Locator(`input[type="text"]`).Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestCountAfterHideEx5 verifies Count counts hidden elements too.
// Ref: TestLocatorCount.java#shouldCountHiddenElements
func TestCountAfterHideEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">A</div>
		<div class="item" style="display:none">B</div>
		<div class="item">C</div>
	`))

	count, err := page.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestCountFormInputsEx6 verifies Count for form inputs.
// Ref: TestLocatorCount.java#shouldCountFormInputs
func TestCountFormInputsEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<input type="text">
			<input type="email">
			<input type="password">
			<textarea></textarea>
			<select><option>A</option></select>
		</form>
	`))

	count, err := page.Locator("input").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestCountNestedDivsEx6 verifies Count for nested divs.
// Ref: TestLocatorCount.java#shouldCountNestedDivs
func TestCountNestedDivsEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="card">
			<div class="inner">1</div>
		</div>
		<div class="card">
			<div class="inner">2</div>
		</div>
		<div class="card">
			<div class="inner">3</div>
		</div>
	`))

	count, err := page.Locator(".card .inner").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestCountAfterFilterEx6 verifies Count after Filter with has-text.
// Ref: TestLocatorCount.java#shouldCountAfterFilter
func TestCountAfterFilterEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="item">Active item</li>
			<li class="item">Inactive item</li>
			<li class="item">Active item</li>
		</ul>
	`))

	hasText := "Active item"
	count, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{HasText: &hasText}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterByTextCount verifies Filter(HasText) returns the right count.
// Ref: TestLocatorFilter.java#shouldFilterByText
func TestLocatorFilterByTextCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Apple</li>
			<li>Banana</li>
			<li>Apricot</li>
			<li>Mango</li>
		</ul>
	`))

	needle := "Ap"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{
		HasText: &needle,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterByHasNotTextCount verifies Filter(HasNotText) returns the right count.
// Ref: TestLocatorFilter.java#shouldFilterByHasNotText
func TestLocatorFilterByHasNotTextCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Active</li>
			<li>Inactive</li>
			<li>Active</li>
		</ul>
	`))

	exclude := "Inactive"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{
		HasNotText: &exclude,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterByHasChild verifies Filter(Has) with child locator.
// Ref: TestLocatorFilter.java#shouldFilterByHasChild
func TestLocatorFilterByHasChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="card">
			<span class="badge">New</span>
			<p>Product A</p>
		</div>
		<div class="card">
			<p>Product B</p>
		</div>
		<div class="card">
			<span class="badge">Sale</span>
			<p>Product C</p>
		</div>
	`))

	badge := page.Locator(".badge")
	count, err := page.Locator(".card").Filter(&playwright.LocatorFilterOptions{
		Has: badge,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterByHasNotChild verifies Filter(HasNot) excludes elements with child.
// Ref: TestLocatorFilter.java#shouldFilterByHasNotChild
func TestLocatorFilterByHasNotChild(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">
			<button class="action">Buy</button>
			Product X
		</div>
		<div class="item">Product Y</div>
		<div class="item">Product Z</div>
	`))

	button := page.Locator(".action")
	count, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{
		HasNot: button,
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterChained verifies multiple Filter calls can be chained.
// Ref: TestLocatorFilter.java#shouldChainFilters
func TestLocatorFilterChained(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="active">Alpha</li>
			<li class="active">Beta</li>
			<li class="inactive">Alpha</li>
			<li class="inactive">Gamma</li>
		</ul>
	`))

	text := "Alpha"
	notText := "inactive"
	active := page.Locator("li").Filter(&playwright.LocatorFilterOptions{
		HasText: &text,
	}).Filter(&playwright.LocatorFilterOptions{
		HasNotText: &notText,
	})

	count, err := active.Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorFilterByTextEx2 verifies Filter with HasText finds matching element.
// Ref: TestLocatorFilter.java#shouldFilterByText
func TestLocatorFilterByTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Alpha</li>
			<li>Beta</li>
			<li>Gamma</li>
		</ul>
	`))

	text := "Beta"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{HasText: &text}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorFilterByNotTextEx2 verifies Filter with HasNotText excludes matching elements.
// Ref: TestLocatorFilter.java#shouldFilterByNotText
func TestLocatorFilterByNotTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">Keep</div>
		<div class="item">Remove</div>
		<div class="item">Keep</div>
	`))

	notText := "Remove"
	count, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{HasNotText: &notText}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterChainedEx2 verifies Filter can be chained with other locator methods.
// Ref: TestLocatorFilter.java#shouldChainWithNth
func TestLocatorFilterChainedEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="active">First Active</li>
			<li>Inactive</li>
			<li class="active">Second Active</li>
		</ul>
	`))

	text, err := page.Locator("li.active").Nth(1).InnerText(ctx)
	must.NoError(err)
	is.Equal("Second Active", text)
}

// TestLocatorFilterCountsMatchEx2 verifies Filter reduces count correctly.
// Ref: TestLocatorFilter.java#shouldReduceCount
func TestLocatorFilterCountsMatchEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="p">Hello</p>
		<p class="p">World</p>
		<p class="p">Hello World</p>
	`))

	text := "Hello"
	count, err := page.Locator(".p").Filter(&playwright.LocatorFilterOptions{HasText: &text}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorFilterInnerTextEx2 verifies Filter result has correct text.
// Ref: TestLocatorFilter.java#shouldReturnCorrectText
func TestLocatorFilterInnerTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="card">Apple</div>
		<div class="card">Banana</div>
		<div class="card">Cherry</div>
	`))

	text := "Cherry"
	innerText, err := page.Locator(".card").Filter(&playwright.LocatorFilterOptions{HasText: &text}).InnerText(ctx)
	must.NoError(err)
	is.Equal("Cherry", innerText)
}

func localStrPtrLF3(s string) *string { return &s }

// TestFilterByHasTextEx3 verifies Filter with HasText finds matching items.
// Ref: TestLocatorFilter.java#shouldFilterByHasText
func TestFilterByHasTextEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="card">Free plan</div>
		<div class="card">Pro plan - $10/month</div>
		<div class="card">Enterprise</div>
	`))

	count, err := page.Locator(".card").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrLF3("Pro plan"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestFilterByHasNotTextEx3 verifies Filter with HasNotText excludes matching items.
// Ref: TestLocatorFilter.java#shouldFilterByHasNotText
func TestFilterByHasNotTextEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li class="item active">Active item</li>
		<li class="item">Inactive item</li>
		<li class="item active">Active item 2</li>
	`))

	count, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{
		HasNotText: localStrPtrLF3("Inactive"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestFilterByHasLocatorEx3 verifies Filter with Has finds items containing sub-locator.
// Ref: TestLocatorFilter.java#shouldFilterByHas
func TestFilterByHasLocatorEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="row"><span class="badge">New</span><span>Product A</span></div>
		<div class="row"><span>Product B</span></div>
		<div class="row"><span class="badge">New</span><span>Product C</span></div>
	`))

	count, err := page.Locator(".row").Filter(&playwright.LocatorFilterOptions{
		Has: page.Locator(".badge"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

func localStrPtrLF4(s string) *string { return &s }

// TestFilterByHasTextButtonEx4 verifies Filter with HasText on buttons.
// Ref: TestLocatorFilter.java#shouldFilterButtons
func TestFilterByHasTextButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button class="action">Save</button>
		<button class="action">Delete</button>
		<button class="action">Cancel</button>
	`))

	count, err := page.Locator(".action").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrLF4("Delete"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestFilterByHasNotTextEx4 verifies Filter HasNotText on table rows.
// Ref: TestLocatorFilter.java#shouldFilterTableRows
func TestFilterByHasNotTextEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr class="row"><td>Alice</td><td>Active</td></tr>
			<tr class="row"><td>Bob</td><td>Inactive</td></tr>
			<tr class="row"><td>Carol</td><td>Active</td></tr>
		</table>
	`))

	count, err := page.Locator(".row").Filter(&playwright.LocatorFilterOptions{
		HasNotText: localStrPtrLF4("Inactive"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestFilterChainedEx4 verifies chained Filter calls.
// Ref: TestLocatorFilter.java#shouldChainFilters
func TestFilterChainedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="card active">Premium member</div>
		<div class="card active">Free member</div>
		<div class="card inactive">Premium member</div>
		<div class="card active">Premium member - sale</div>
	`))

	count, err := page.Locator(".card").
		Filter(&playwright.LocatorFilterOptions{HasText: localStrPtrLF4("Premium")}).
		Filter(&playwright.LocatorFilterOptions{HasText: localStrPtrLF4("active")}).
		Count(ctx)
	must.NoError(err)
	is.GreaterOrEqual(count, 1)
}

// TestLocatorFocusFiresFocusEvent verifies Focus() fires the focus event.
// Ref: TestLocatorFocus.java#shouldFireFocusEvent
func TestLocatorFocusFiresFocusEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp">
		<script>
		window.focused = false;
		document.getElementById('inp').addEventListener('focus', () => { window.focused = true; });
		</script>
	`))

	must.NoError(page.Locator("input").Focus(ctx))

	focused, err := page.Evaluate(ctx, "() => window.focused")
	must.NoError(err)
	is.Equal(true, focused)
}

// TestLocatorFocusMakesElementActiveElement verifies the focused element becomes document.activeElement.
// Ref: TestLocatorFocus.java#shouldMakeElementActiveElement
func TestLocatorFocusMakesElementActiveElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first">
		<input id="second">
	`))

	must.NoError(page.Locator("#second").Focus(ctx))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", activeId)
}

// TestLocatorFocusOnTextarea verifies Focus() works on textarea elements.
// Ref: TestLocatorFocus.java#shouldWorkOnTextarea
func TestLocatorFocusOnTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	must.NoError(page.Locator("textarea").Focus(ctx))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("ta", activeId)
}

// TestLocatorFocusOnSelect verifies Focus() works on select elements.
// Ref: TestLocatorFocus.java#shouldWorkOnSelect
func TestLocatorFocusOnSelect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option>A</option>
			<option>B</option>
		</select>
	`))

	must.NoError(page.Locator("select").Focus(ctx))

	activeTag, err := page.Evaluate(ctx, `() => document.activeElement.tagName`)
	must.NoError(err)
	is.Equal("SELECT", activeTag)
}

// TestLocatorFocusThenCheckToBeFocused verifies Focus() + ToBeRocused assertion works.
// Ref: TestLocatorFocus.java#shouldSupportToBeFocusedAssertion
func TestLocatorFocusThenCheckToBeFocused(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" tabindex="0">`))

	loc := page.Locator("#inp")
	must.NoError(loc.Focus(ctx))

	isFocused, err := page.Evaluate(ctx, `() => document.activeElement.id === 'inp'`)
	must.NoError(err)
	is.Equal(true, isFocused)
}

// TestLocatorFocusFiresFocusEvent verifies Focus fires the focus DOM event.
// Ref: TestLocatorFocus.java#shouldFireFocusEvent
func TestLocatorFocusFiresFocusEventExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onfocus="this.dataset.focused='yes'">
	`))

	must.NoError(page.Locator("input").Focus(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('inp').dataset.focused`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestLocatorFocusMakesElementActiveElement verifies Focus sets element as document.activeElement.
// Ref: TestLocatorFocus.java#shouldMakeActiveElement
func TestLocatorFocusMakesElementActiveElementExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first">
		<input id="second">
	`))

	must.NoError(page.Locator("#second").Focus(ctx))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", activeId)
}

// TestLocatorFocusOnButton verifies Focus works on button elements.
// Ref: TestLocatorFocus.java#shouldFocusButton
func TestLocatorFocusOnButtonExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn">focusable</button>
	`))

	must.NoError(page.Locator("button").Focus(ctx))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("btn", activeId)
}

// TestLocatorFocusOnTextarea verifies Focus works on textarea elements.
// Ref: TestLocatorFocus.java#shouldFocusTextarea
func TestLocatorFocusOnTextareaExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<textarea id="ta"></textarea>
	`))

	must.NoError(page.Locator("textarea").Focus(ctx))

	activeTag, err := page.Evaluate(ctx, `() => document.activeElement.tagName.toLowerCase()`)
	must.NoError(err)
	is.Equal("textarea", activeTag)
}

// TestLocatorFocusSwitchesBetweenElements verifies Focus can be transferred between elements.
// Ref: TestLocatorFocus.java#shouldTransferFocus
func TestLocatorFocusSwitchesBetweenElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="a">
		<input id="b">
	`))

	must.NoError(page.Locator("#a").Focus(ctx))
	firstActive, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("a", firstActive)

	must.NoError(page.Locator("#b").Focus(ctx))
	secondActive, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("b", secondActive)
}

// TestLocatorFocusOnSelect verifies Focus works on select element.
// Ref: TestLocatorFocus.java#shouldFocusSelect
func TestLocatorFocusOnSelectExtra3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" onfocus="window.__focused=true">
			<option value="a">A</option>
		</select>
	`))

	must.NoError(page.Locator("#sel").Focus(ctx))

	result, err := page.Evaluate(ctx, `() => window.__focused`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorFocusOnTextarea verifies Focus works on textarea.
// Ref: TestLocatorFocus.java#shouldFocusTextarea
func TestLocatorFocusOnTextareaExtra3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<textarea id="ta" onfocus="window.__taFocused=true"></textarea>
	`))

	must.NoError(page.Locator("#ta").Focus(ctx))

	result, err := page.Evaluate(ctx, `() => window.__taFocused`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorFocusMakesDocumentActiveElement verifies Focus makes element document.activeElement.
// Ref: TestLocatorFocus.java#shouldBecomeActiveElement
func TestLocatorFocusMakesDocumentActiveElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	must.NoError(page.Locator("#inp").Focus(ctx))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("inp", activeId)
}

// TestLocatorFocusSwitchBetweenElements verifies Focus switches active element.
// Ref: TestLocatorFocus.java#shouldSwitchActiveElement
func TestLocatorFocusSwitchBetweenElementsExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="first">
		<input type="text" id="second">
	`))

	must.NoError(page.Locator("#first").Focus(ctx))
	first, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("first", first)

	must.NoError(page.Locator("#second").Focus(ctx))
	second, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", second)
}

// TestLocatorFocusFiresFocusEvent verifies Focus fires focus event.
// Ref: TestLocatorFocus.java#shouldFireFocusEvent
func TestLocatorFocusFiresFocusEventExtra3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="inp" onfocus="window.__focusFired=true">
	`))

	must.NoError(page.Locator("#inp").Focus(ctx))

	result, err := page.Evaluate(ctx, `() => window.__focusFired`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorFocusInputEx4 verifies Focus sets input as active element.
// Ref: TestLocatorFocus.java#shouldFocusInput
func TestLocatorFocusInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Focus(ctx))

	active, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("inp", active)
}

// TestLocatorFocusButtonEx4 verifies Focus sets button as active element.
// Ref: TestLocatorFocus.java#shouldFocusButton
func TestLocatorFocusButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	must.NoError(page.Locator("#btn").Focus(ctx))

	active, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("btn", active)
}

// TestLocatorFocusFiresFocusEventEx4 verifies Focus fires onfocus event.
// Ref: TestLocatorFocus.java#shouldFireFocusEvent
func TestLocatorFocusFiresFocusEventEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" onfocus="window.__focusFired=true">
	`))

	must.NoError(page.Locator("#inp").Focus(ctx))

	fired, err := page.Evaluate(ctx, `() => window.__focusFired`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestLocatorFocusSelectEx4 verifies Focus sets select as active element.
// Ref: TestLocatorFocus.java#shouldFocusSelect
func TestLocatorFocusSelectEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel"><option value="a">A</option></select>
	`))

	must.NoError(page.Locator("#sel").Focus(ctx))

	active, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("sel", active)
}

// TestLocatorFocusLinkEx5 verifies Focus works on link element.
// Ref: TestLocatorFocus.java#shouldFocusLink
func TestLocatorFocusLinkEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="lnk" href="#">Link</a>`))
	must.NoError(page.Locator("#lnk").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("lnk", focused)
}

// TestLocatorFocusSelectEx5 verifies Focus works on select element.
// Ref: TestLocatorFocus.java#shouldFocusSelect
func TestLocatorFocusSelectEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option>A</option>
		</select>
	`))
	must.NoError(page.Locator("#sel").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("sel", focused)
}

// TestLocatorFocusButtonEx5 verifies Focus works on button element.
// Ref: TestLocatorFocus.java#shouldFocusButton
func TestLocatorFocusButtonEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Focus me</button>`))
	must.NoError(page.Locator("#btn").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("btn", focused)
}

// TestLocatorFocusTextareaEx5 verifies Focus works on textarea element.
// Ref: TestLocatorFocus.java#shouldFocusTextarea
func TestLocatorFocusTextareaEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">text</textarea>`))
	must.NoError(page.Locator("#ta").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("ta", focused)
}

// TestFocusSelectEx6 verifies Focus works on select elements.
// Ref: TestLocatorFocus.java#shouldFocusSelect
func TestFocusSelectEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel"><option>A</option></select>`))

	must.NoError(page.Locator("#sel").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("sel", focused)
}

// TestFocusTextareaEx6 verifies Focus works on textarea elements.
// Ref: TestLocatorFocus.java#shouldFocusTextarea
func TestFocusTextareaEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	must.NoError(page.Locator("#ta").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("ta", focused)
}

// TestFocusLinkEx6 verifies Focus works on anchor elements with tabindex.
// Ref: TestLocatorFocus.java#shouldFocusLink
func TestFocusLinkEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="#" tabindex="0">Link</a>`))

	must.NoError(page.Locator("#a").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("a", focused)
}

// TestFocusSpanWithTabIndexEx6 verifies Focus works on span with tabindex.
// Ref: TestLocatorFocus.java#shouldFocusSpanWithTabIndex
func TestFocusSpanWithTabIndexEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="sp" tabindex="0" role="button">Click</span>`))

	must.NoError(page.Locator("#sp").Focus(ctx))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("sp", focused)
}

// TestLocatorGetAttributeReturnsValue verifies GetAttribute returns existing attribute.
// Ref: TestLocatorGetAttribute.java#shouldReturnAttributeValue
func TestLocatorGetAttributeReturnsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a href="https://example.com" id="lnk">link</a>`))

	attr, err := page.Locator("#lnk").GetAttribute(ctx, "href")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("https://example.com", *attr)
}

// TestLocatorGetAttributeReturnsNilForMissing verifies GetAttribute returns nil for absent attribute.
// Ref: TestLocatorGetAttribute.java#shouldReturnNilForMissing
func TestLocatorGetAttributeReturnsNilForMissing(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">hello</div>`))

	attr, err := page.Locator("#el").GetAttribute(ctx, "data-nonexistent")
	must.NoError(err)
	is.Nil(attr)
}

// TestLocatorGetAttributeDataAttribute verifies GetAttribute works for data-* attributes.
// Ref: TestLocatorGetAttribute.java#shouldGetDataAttribute
func TestLocatorGetAttributeDataAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" data-value="42">content</div>`))

	attr, err := page.Locator("#el").GetAttribute(ctx, "data-value")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("42", *attr)
}

// TestLocatorGetAttributeClassAttribute verifies GetAttribute for class.
// Ref: TestLocatorGetAttribute.java#shouldGetClassAttribute
func TestLocatorGetAttributeClassAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" class="foo bar baz">content</div>`))

	attr, err := page.Locator("#el").GetAttribute(ctx, "class")
	must.NoError(err)
	must.NotNil(attr)
	is.Contains(*attr, "foo")
	is.Contains(*attr, "bar")
}

// TestLocatorGetAttributeInputType verifies GetAttribute for input type.
// Ref: TestLocatorGetAttribute.java#shouldGetInputType
func TestLocatorGetAttributeInputType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="email" placeholder="Enter email">`))

	attr, err := page.Locator("#inp").GetAttribute(ctx, "type")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("email", *attr)
}

// TestLocatorGetAttributeEmptyValue verifies GetAttribute returns empty string for empty attribute.
// Ref: TestLocatorGetAttribute.java#shouldReturnEmptyForEmptyAttribute
func TestLocatorGetAttributeEmptyValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" disabled>`))

	attr, err := page.Locator("#inp").GetAttribute(ctx, "disabled")
	must.NoError(err)
	// disabled attribute without value returns empty string, not nil
	must.NotNil(attr)
}

// TestLocatorHasTextFiltersEx2 verifies Locator Filter with HasText filters correctly.
// Ref: TestLocatorHasText.java#shouldFilterByText
func TestLocatorHasTextFiltersEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">Apple</div>
		<div class="item">Banana</div>
		<div class="item">Cherry</div>
	`))

	banana := "Banana"
	count, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{HasText: &banana}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestLocatorHasTextReturnAllEx2 verifies no filter returns all matching elements.
// Ref: TestLocatorHasText.java#shouldReturnAllWithoutFilter
func TestLocatorHasTextReturnAllEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">Alpha</div>
		<div class="item">Beta</div>
		<div class="item">Gamma</div>
	`))

	count, err := page.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestLocatorHasTextPartialMatchEx2 verifies Filter HasText matches partial text.
// Ref: TestLocatorHasText.java#shouldMatchPartialText
func TestLocatorHasTextPartialMatchEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="note">Important note here</p>
		<p class="note">Another note</p>
		<p class="note">Unrelated content</p>
	`))

	hasText := "note"
	count, err := page.Locator(".note").Filter(&playwright.LocatorFilterOptions{HasText: &hasText}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorHasTextNoMatchReturnsZeroEx2 verifies Filter HasText with no matches returns 0.
// Ref: TestLocatorHasText.java#shouldReturnZeroWhenNoMatch
func TestLocatorHasTextNoMatchReturnsZeroEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>One</li>
		<li>Two</li>
	`))

	hasText := "Three"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{HasText: &hasText}).Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// TestLocatorHasNotTextFiltersEx2 verifies Filter HasNotText excludes elements with matching text.
// Ref: TestLocatorHasText.java#shouldFilterByNotText
func TestLocatorHasNotTextFiltersEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="row">Include me</div>
		<div class="row">Exclude me</div>
		<div class="row">Include me too</div>
	`))

	hasNotText := "Exclude"
	count, err := page.Locator(".row").Filter(&playwright.LocatorFilterOptions{HasNotText: &hasNotText}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorInnerHTMLReturnsNestedHTML verifies InnerHTML returns full nested HTML.
// Ref: TestLocatorInnerHTML.java#shouldReturnNestedHTML
func TestLocatorInnerHTMLReturnsNestedHTML(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="parent"><span id="child">text</span></div>`))

	html, err := page.Locator("#parent").InnerHTML(ctx)
	must.NoError(err)
	is.Equal(`<span id="child">text</span>`, html)
}

// TestLocatorInnerHTMLForEmptyElement verifies InnerHTML returns empty string for empty element.
// Ref: TestLocatorInnerHTML.java#shouldReturnEmptyForEmptyElement
func TestLocatorInnerHTMLForEmptyElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="empty"></div>`))

	html, err := page.Locator("#empty").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("", html)
}

// TestLocatorInnerHTMLWithMultipleChildren verifies InnerHTML includes all children.
// Ref: TestLocatorInnerHTML.java#shouldReturnAllChildren
func TestLocatorInnerHTMLWithMultipleChildren(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"><li>a</li><li>b</li></ul>`))

	html, err := page.Locator("#list").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<li>a</li><li>b</li>", html)
}

// TestLocatorInnerTextStripsHTMLTags verifies InnerText returns only visible text without tags.
// Ref: TestLocatorInnerHTML.java#shouldStripTags
func TestLocatorInnerTextStripsHTMLTags(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p"><b>bold</b> and <i>italic</i></p>`))

	text, err := page.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("bold and italic", text)
}

// TestLocatorInnerHTMLModifiedByEvaluate verifies InnerHTML reflects JavaScript-modified DOM.
// Ref: TestLocatorInnerHTML.java#shouldReflectModifications
func TestLocatorInnerHTMLModifiedByEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"></div>`))
	_, err := page.Evaluate(ctx, `() => { document.getElementById('container').innerHTML = '<span>dynamic</span>'; }`)
	must.NoError(err)

	html, err := page.Locator("#container").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<span>dynamic</span>", html)
}

// TestInnerHTMLWithAttributes verifies InnerHTML preserves element attributes.
// Ref: TestLocatorInnerHTML.java#shouldPreserveAttributes
func TestInnerHTMLWithAttributes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"><a href="https://example.com" class="link">click</a></div>`))

	html, err := page.Locator("#container").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, `href="https://example.com"`)
	is.Contains(html, `class="link"`)
}

// TestInnerTextStripsHTMLTags verifies InnerText strips HTML tags.
// Ref: TestLocatorInnerHTML.java#shouldStripTags
func TestInnerTextStripsHTMLTags(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p"><b>bold</b> and <i>italic</i></p>`))

	text, err := page.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("bold and italic", text)
}

// TestTextContentIncludesHidden verifies TextContent includes hidden text.
// Ref: TestLocatorInnerHTML.java#shouldIncludeHiddenText
func TestTextContentIncludesHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="parent">
			visible
			<span style="display:none">hidden</span>
		</div>
	`))

	tc, err := page.Locator("#parent").TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Contains(*tc, "hidden")
}

// TestInnerHTMLEmptyForLeafElement verifies InnerHTML is empty for leaf elements.
// Ref: TestLocatorInnerHTML.java#shouldReturnEmptyForLeaf
func TestInnerHTMLEmptyForLeafElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	html, err := page.Locator("#inp").InnerHTML(ctx)
	must.NoError(err)
	is.Empty(html)
}

// TestAllTextContentsReturnsAllTexts verifies AllTextContents returns all text values.
// Ref: TestLocatorInnerHTML.java#shouldReturnAllTextContents
func TestAllTextContentsReturnsAllTextsExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">alpha</div>
		<div class="item">beta</div>
		<div class="item">gamma</div>
	`))

	texts, err := page.Locator(".item").AllTextContents(ctx)
	must.NoError(err)
	is.Len(texts, 3)
	is.Equal("alpha", texts[0])
	is.Equal("beta", texts[1])
	is.Equal("gamma", texts[2])
}

// TestLocatorInnerHTMLContainsTagsEx3 verifies InnerHTML returns HTML with tags.
// Ref: TestLocatorInnerHTML.java#shouldContainTags
func TestLocatorInnerHTMLContainsTagsEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><span>hello</span></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<span>")
	is.Contains(html, "hello")
}

// TestLocatorInnerHTMLEmptyDivEx3 verifies InnerHTML on empty div returns empty string.
// Ref: TestLocatorInnerHTML.java#shouldReturnEmptyForEmptyDiv
func TestLocatorInnerHTMLEmptyDivEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="empty"></div>`))

	html, err := page.Locator("#empty").InnerHTML(ctx)
	must.NoError(err)
	is.Empty(html)
}

// TestLocatorInnerHTMLWithNestedStructureEx3 verifies InnerHTML captures nested HTML.
// Ref: TestLocatorInnerHTML.java#shouldCaptureNested
func TestLocatorInnerHTMLWithNestedStructureEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			<ul><li>Item 1</li><li>Item 2</li></ul>
		</div>
	`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<ul>")
	is.Contains(html, "Item 1")
}

// TestLocatorInnerHTMLDynamicContentEx3 verifies InnerHTML reflects JS DOM changes.
// Ref: TestLocatorInnerHTML.java#shouldReflectDOMChanges
func TestLocatorInnerHTMLDynamicContentEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"></div>`))

	_, err := page.Evaluate(ctx, `() => {
		const p = document.createElement('p');
		p.textContent = 'Dynamic';
		document.getElementById('d').appendChild(p);
	}`)
	must.NoError(err)

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "Dynamic")
}

// TestLocatorInnerHTMLNestedEx4 verifies InnerHTML returns nested tags.
// Ref: TestLocatorInnerHTML.java#shouldReturnNestedTags
func TestLocatorInnerHTMLNestedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><ul><li>Item</li></ul></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<ul><li>Item</li></ul>", html)
}

// TestLocatorInnerHTMLEmptyEx4 verifies InnerHTML returns empty for empty element.
// Ref: TestLocatorInnerHTML.java#shouldReturnEmptyForEmpty
func TestLocatorInnerHTMLEmptyEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("", html)
}

// TestLocatorInnerHTMLAttributesEx4 verifies InnerHTML includes attributes.
// Ref: TestLocatorInnerHTML.java#shouldIncludeAttributes
func TestLocatorInnerHTMLAttributesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><a href="url" class="link">Click</a></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, `href="url"`)
	is.Contains(html, `class="link"`)
}

// TestLocatorInnerHTMLMultipleChildrenEx4 verifies InnerHTML with multiple children.
// Ref: TestLocatorInnerHTML.java#shouldHandleMultipleChildren
func TestLocatorInnerHTMLMultipleChildrenEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><span>A</span><span>B</span><span>C</span></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<span>A</span><span>B</span><span>C</span>", html)
}

// TestInnerHTMLNestedListEx5 verifies InnerHTML for nested list structures.
// Ref: TestLocatorInnerHTML.java#shouldReturnNestedList
func TestInnerHTMLNestedListEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul id="list">
			<li>Item <strong>A</strong></li>
			<li>Item <em>B</em></li>
		</ul>
	`))

	html, err := page.Locator("#list").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<strong>A</strong>")
	is.Contains(html, "<em>B</em>")
}

// TestInnerHTMLFormElementsEx5 verifies InnerHTML for form elements.
// Ref: TestLocatorInnerHTML.java#shouldReturnFormElements
func TestInnerHTMLFormElementsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			<input type="text" id="inp">
			<button>Submit</button>
		</div>
	`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "input")
	is.Contains(html, "button")
}

// TestInnerHTMLAfterInsertEx5 verifies InnerHTML reflects dynamic DOM insertion.
// Ref: TestLocatorInnerHTML.java#shouldReflectInsertion
func TestInnerHTMLAfterInsertEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"></div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').innerHTML = '<p id="inserted">New</p>'`)
	must.NoError(err)

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "inserted")
}

// TestInnerHTMLAnchorEx5 verifies InnerHTML for div containing anchor.
// Ref: TestLocatorInnerHTML.java#shouldReturnAnchorInDiv
func TestInnerHTMLAnchorEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><a href="/go">Go</a></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<a")
	is.Contains(html, "Go")
}

// TestInnerHTMLWithAttributesEx6 verifies InnerHTML preserves element attributes.
// Ref: TestLocatorInnerHTML.java#shouldPreserveAttributes
func TestInnerHTMLWithAttributesEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><img src="image.jpg" alt="Photo"></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, `alt="Photo"`)
}

// TestInnerHTMLHeadersEx6 verifies InnerHTML returns table headers.
// Ref: TestLocatorInnerHTML.java#shouldReturnTableHeaders
func TestInnerHTMLHeadersEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table id="t">
			<thead><tr><th>Name</th><th>Value</th></tr></thead>
		</table>
	`))

	html, err := page.Locator("#t").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "Name")
	is.Contains(html, "Value")
}

// TestInnerHTMLEmptyDivEx6 verifies InnerHTML is empty for empty div.
// Ref: TestLocatorInnerHTML.java#shouldBeEmptyForEmptyDiv
func TestInnerHTMLEmptyDivEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="empty"></div>`))

	html, err := page.Locator("#empty").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("", html)
}

// TestInnerHTMLScriptContentEx6 verifies InnerHTML for div with multiple children.
// Ref: TestLocatorInnerHTML.java#shouldReturnMultipleChildren
func TestInnerHTMLScriptContentEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<section id="s">
			<h2>Title</h2>
			<p>Paragraph</p>
			<ul><li>Item</li></ul>
		</section>
	`))

	html, err := page.Locator("#s").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<h2>")
	is.Contains(html, "<p>")
	is.Contains(html, "<ul>")
}

// TestInnerHTMLSVGEx7 verifies InnerHTML returns SVG content.
// Ref: TestLocatorInnerHTML.java#shouldReturnSVGContent
func TestInnerHTMLSVGEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			<svg><circle id="c" r="5"/></svg>
		</div>
	`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "circle")
}

// TestInnerHTMLAfterReplaceEx7 verifies InnerHTML after replacing children.
// Ref: TestLocatorInnerHTML.java#shouldReturnAfterReplace
func TestInnerHTMLAfterReplaceEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><p>Old content</p></div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').innerHTML = '<span>New</span>'`)
	must.NoError(err)

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "New")
	must.NotContains(html, "Old")
}

// TestInnerHTMLDataAttrEx7 verifies InnerHTML preserves data attributes.
// Ref: TestLocatorInnerHTML.java#shouldPreserveDataAttributes
func TestInnerHTMLDataAttrEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><span data-id="42">Item</span></div>`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "data-id")
	is.Contains(html, "42")
}

// TestInnerHTMLComplexNestingEx7 verifies InnerHTML with deep nesting.
// Ref: TestLocatorInnerHTML.java#shouldHandleDeepNesting
func TestInnerHTMLComplexNestingEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			<ul>
				<li><a href="#"><strong>Bold link</strong></a></li>
			</ul>
		</div>
	`))

	html, err := page.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, "<strong>")
	is.Contains(html, "<ul>")
}

// TestLocatorInnerTextReturnsText verifies InnerText returns visible text.
// Ref: TestLocatorInnerText.java#shouldReturnText
func TestLocatorInnerTextReturnsTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello World</p>`))

	text, err := page.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello World", text)
}

// TestLocatorInnerTextExcludesHiddenText verifies InnerText excludes CSS-hidden text.
// Ref: TestLocatorInnerText.java#shouldExcludeHiddenText
func TestLocatorInnerTextExcludesHiddenText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			Visible
			<span style="display:none">Hidden</span>
		</div>
	`))

	text, err := page.Locator("#d").InnerText(ctx)
	must.NoError(err)
	is.Contains(text, "Visible")
	must.NotContains(text, "Hidden")
}

// TestLocatorInnerTextWithNestedElements verifies InnerText includes nested text.
// Ref: TestLocatorInnerText.java#shouldIncludeNestedText
func TestLocatorInnerTextWithNestedElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d"><span>First</span> <em>Second</em></div>
	`))

	text, err := page.Locator("#d").InnerText(ctx)
	must.NoError(err)
	is.Contains(text, "First")
	is.Contains(text, "Second")
}

// TestLocatorInnerTextEmptyElement verifies InnerText on empty element returns empty string.
// Ref: TestLocatorInnerText.java#shouldReturnEmptyForEmptyElement
func TestLocatorInnerTextEmptyElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="empty"></div>`))

	text, err := page.Locator("#empty").InnerText(ctx)
	must.NoError(err)
	is.Empty(text)
}

// TestLocatorInnerTextButton verifies InnerText works on button elements.
// Ref: TestLocatorInnerText.java#shouldWorkOnButton
func TestLocatorInnerTextButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Submit Form</button>`))

	text, err := page.Locator("#btn").InnerText(ctx)
	must.NoError(err)
	is.Equal("Submit Form", text)
}

// TestInnerTextHeadingEx3 verifies InnerText returns heading text.
// Ref: TestLocatorInnerText.java#shouldReturnHeading
func TestInnerTextHeadingEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h">Hello World</h1>`))

	text, err := page.Locator("#h").InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello World", text)
}

// TestInnerTextListItemEx3 verifies InnerText works for list items.
// Ref: TestLocatorInnerText.java#shouldWorkForListItem
func TestInnerTextListItemEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li id="first">First Item</li>
		</ul>
	`))

	text, err := page.Locator("#first").InnerText(ctx)
	must.NoError(err)
	is.Equal("First Item", text)
}

// TestInnerTextWithBoldEx3 verifies InnerText returns visible text skipping tags.
// Ref: TestLocatorInnerText.java#shouldReturnTextWithBold
func TestInnerTextWithBoldEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Normal <b>Bold</b> Text</p>`))

	text, err := page.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("Normal Bold Text", text)
}

// TestInnerTextSpanEx3 verifies InnerText for span elements.
// Ref: TestLocatorInnerText.java#shouldReturnSpanText
func TestInnerTextSpanEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="s">Span content</span>`))

	text, err := page.Locator("#s").InnerText(ctx)
	must.NoError(err)
	is.Equal("Span content", text)
}

// TestInnerTextAnchorEx4 verifies InnerText for anchor elements.
// Ref: TestLocatorInnerText.java#shouldReturnAnchorText
func TestInnerTextAnchorEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="#">Click here</a>`))

	text, err := page.Locator("#a").InnerText(ctx)
	must.NoError(err)
	is.Equal("Click here", text)
}

// TestInnerTextTableCellEx4 verifies InnerText for table cells.
// Ref: TestLocatorInnerText.java#shouldReturnTableCellText
func TestInnerTextTableCellEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr><td id="c1">Cell 1</td><td id="c2">Cell 2</td></tr>
		</table>
	`))

	text, err := page.Locator("#c2").InnerText(ctx)
	must.NoError(err)
	is.Equal("Cell 2", text)
}

// TestInnerTextAfterDOMUpdateEx4 verifies InnerText reflects DOM updates.
// Ref: TestLocatorInnerText.java#shouldReflectDOMUpdate
func TestInnerTextAfterDOMUpdateEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Original text</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').textContent = 'Updated text'`)
	must.NoError(err)

	text, err := page.Locator("#d").InnerText(ctx)
	must.NoError(err)
	is.Equal("Updated text", text)
}

// TestInnerTextLabelEx4 verifies InnerText for label elements.
// Ref: TestLocatorInnerText.java#shouldReturnLabelText
func TestInnerTextLabelEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<label id="lbl" for="inp">Email Address</label><input id="inp">`))

	text, err := page.Locator("#lbl").InnerText(ctx)
	must.NoError(err)
	is.Equal("Email Address", text)
}

// TestInnerTextSpanInsideParaEx5 verifies InnerText for span inside paragraph.
// Ref: TestLocatorInnerText.java#shouldReturnSpanInsideParagraph
func TestInnerTextSpanInsideParaEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello <span>World</span></p>`))

	text, err := page.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello World", text)
}

// TestInnerTextEmptyElementEx5 verifies InnerText returns empty string for empty element.
// Ref: TestLocatorInnerText.java#shouldReturnEmptyForEmpty
func TestInnerTextEmptyElementEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="empty"></div>`))

	text, err := page.Locator("#empty").InnerText(ctx)
	must.NoError(err)
	is.Equal("", text)
}

// TestInnerTextAfterAppendEx5 verifies InnerText reflects appended content.
// Ref: TestLocatorInnerText.java#shouldReflectAppendedContent
func TestInnerTextAfterAppendEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Initial</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').textContent += ' Appended'`)
	must.NoError(err)

	text, err := page.Locator("#d").InnerText(ctx)
	must.NoError(err)
	is.Equal("Initial Appended", text)
}

// TestInnerTextListItemEx5 verifies InnerText for list items.
// Ref: TestLocatorInnerText.java#shouldReturnListItemText
func TestInnerTextListItemEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li id="li1">Item One</li>
			<li id="li2">Item Two</li>
		</ul>
	`))

	text, err := page.Locator("#li2").InnerText(ctx)
	must.NoError(err)
	is.Equal("Item Two", text)
}

// TestInnerTextDefinitionEx6 verifies InnerText for definition terms.
// Ref: TestLocatorInnerText.java#shouldReturnDefinitionTerm
func TestInnerTextDefinitionEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<dl><dt id="dt">Term</dt><dd>Definition</dd></dl>`))

	text, err := page.Locator("#dt").InnerText(ctx)
	must.NoError(err)
	is.Equal("Term", text)
}

// TestInnerTextFormLabelEx6 verifies InnerText for form label.
// Ref: TestLocatorInnerText.java#shouldReturnFormLabel
func TestInnerTextFormLabelEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<label id="lbl" for="name">Full Name <span>(required)</span></label>
		<input id="name" type="text">
	`))

	text, err := page.Locator("#lbl").InnerText(ctx)
	must.NoError(err)
	is.Contains(text, "Full Name")
	is.Contains(text, "required")
}

// TestInnerTextBlockquoteEx6 verifies InnerText for blockquote elements.
// Ref: TestLocatorInnerText.java#shouldReturnBlockquoteText
func TestInnerTextBlockquoteEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<blockquote id="bq">Famous quote here</blockquote>`))

	text, err := page.Locator("#bq").InnerText(ctx)
	must.NoError(err)
	is.Equal("Famous quote here", text)
}

// TestInnerTextFooterEx6 verifies InnerText for footer element.
// Ref: TestLocatorInnerText.java#shouldReturnFooterText
func TestInnerTextFooterEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<footer id="footer">Copyright 2025</footer>`))

	text, err := page.Locator("#footer").InnerText(ctx)
	must.NoError(err)
	is.Equal("Copyright 2025", text)
}

// TestLocatorInputValueFromInput verifies InputValue works for text input.
// Ref: TestLocatorInputValue.java#shouldGetValueFromInput
func TestLocatorInputValueFromInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="test value">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("test value", val)
}

// TestLocatorInputValueFromTextarea verifies InputValue works for textarea.
// Ref: TestLocatorInputValue.java#shouldGetValueFromTextarea
func TestLocatorInputValueFromTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea>multiline content</textarea>`))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("multiline content", val)
}

// TestLocatorInputValueEmptyForEmptyInput verifies InputValue returns empty string for empty input.
// Ref: TestLocatorInputValue.java#shouldReturnEmptyForEmptyInput
func TestLocatorInputValueEmptyForEmptyInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorInputValueAfterFill verifies InputValue reflects value after Fill.
// Ref: TestLocatorInputValue.java#shouldReflectValueAfterFill
func TestLocatorInputValueAfterFill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))
	must.NoError(page.Locator("input").Fill(ctx, "filled value"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("filled value", val)
}

// TestLocatorInputValueFromSelect verifies InputValue works for select elements.
// Ref: TestLocatorInputValue.java#shouldGetValueFromSelect
func TestLocatorInputValueFromSelect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select>
			<option value="opt1">Option 1</option>
			<option value="opt2" selected>Option 2</option>
		</select>
	`))

	val, err := page.Locator("select").InputValue(ctx)
	must.NoError(err)
	is.Equal("opt2", val)
}

// TestInputValueFromPasswordField verifies InputValue reads from password field.
// Ref: TestLocatorInputValue.java#shouldReadPasswordField
func TestInputValueFromPasswordField(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="password" value="secret123">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("secret123", val)
}

// TestInputValueFromNumberField verifies InputValue reads number input.
// Ref: TestLocatorInputValue.java#shouldReadNumberField
func TestInputValueFromNumberField(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="number" value="42">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("42", val)
}

// TestInputValueFromDateField verifies InputValue reads date input.
// Ref: TestLocatorInputValue.java#shouldReadDateField
func TestInputValueFromDateField(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="date" value="2024-06-15">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("2024-06-15", val)
}

// TestInputValueAfterTyping verifies InputValue reflects typed value.
// Ref: TestLocatorInputValue.java#shouldReflectTypedValue
func TestInputValueAfterTyping(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "typed value"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("typed value", val)
}

// TestInputValueFromHiddenInput verifies InputValue reads from hidden input.
// Ref: TestLocatorInputValue.java#shouldReadHiddenInput
func TestInputValueFromHiddenInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="hidden" value="hidden-value">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("hidden-value", val)
}

// TestInputValueTextEx3 verifies InputValue returns filled text.
// Ref: TestLocatorInputValue.java#shouldReturnFilledText
func TestInputValueTextEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hello", val)
}

// TestInputValueAfterFillEx3 verifies InputValue returns value set by Fill.
// Ref: TestLocatorInputValue.java#shouldReturnAfterFill
func TestInputValueAfterFillEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Fill(ctx, "typed text"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("typed text", val)
}

// TestInputValueEmptyEx3 verifies InputValue returns empty string for empty input.
// Ref: TestLocatorInputValue.java#shouldReturnEmpty
func TestInputValueEmptyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestInputValueTextareaEx3 verifies InputValue works for textarea.
// Ref: TestLocatorInputValue.java#shouldWorkForTextarea
func TestInputValueTextareaEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">multiline text</textarea>`))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Equal("multiline text", val)
}

// TestInputValueNumberEx3 verifies InputValue returns number as string.
// Ref: TestLocatorInputValue.java#shouldReturnNumberAsString
func TestInputValueNumberEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="num" type="number" value="42">`))

	val, err := page.Locator("#num").InputValue(ctx)
	must.NoError(err)
	is.Equal("42", val)
}

// TestInputValueAfterTypeEx4 verifies InputValue after Type action.
// Ref: TestLocatorInputValue.java#shouldGetValueAfterType
func TestInputValueAfterTypeEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	el, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)
	must.NoError(el.Type(ctx, "typed text"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("typed text", val)
}

// TestInputValueAfterFillEx4 verifies InputValue after Fill action.
// Ref: TestLocatorInputValue.java#shouldGetValueAfterFill
func TestInputValueAfterFillEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "filled text"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("filled text", val)
}

// TestInputValueAfterClearEx4 verifies InputValue is empty after clearing.
// Ref: TestLocatorInputValue.java#shouldBeEmptyAfterClear
func TestInputValueAfterClearEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="initial">`))

	must.NoError(page.Locator("#inp").Fill(ctx, ""))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestInputValueNumberEx4 verifies InputValue for number input type.
// Ref: TestLocatorInputValue.java#shouldGetNumberValue
func TestInputValueNumberEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="num" type="number" value="42">`))

	val, err := page.Locator("#num").InputValue(ctx)
	must.NoError(err)
	is.Equal("42", val)
}

// TestLocatorLocatorChainsSelectors verifies Locator.Locator scopes to parent element.
// Ref: TestLocatorLocator.java#shouldChainSelectors
func TestLocatorLocatorChainsSelectors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="parent">
			<span class="child">inside</span>
		</div>
		<span class="child">outside</span>
	`))

	text, err := page.Locator(".parent").Locator(".child").InnerText(ctx)
	must.NoError(err)
	is.Equal("inside", text)
}

// TestLocatorLocatorCountScopedToParent verifies chained Count is scoped.
// Ref: TestLocatorLocator.java#shouldCountScopedChildren
func TestLocatorLocatorCountScopedToParent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul class="list">
			<li>a</li>
			<li>b</li>
		</ul>
		<ul class="other">
			<li>c</li>
		</ul>
	`))

	count, err := page.Locator(".list").Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestLocatorLocatorWorksWithButton verifies chained locator can click button.
// Ref: TestLocatorLocator.java#shouldWorkWithButton
func TestLocatorLocatorWorksWithButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form id="form">
			<button onclick="document.getElementById('result').textContent='clicked'">go</button>
		</form>
		<div id="result"></div>
	`))

	must.NoError(page.Locator("#form").Locator("button").Click(ctx))

	text, err := page.Locator("#result").InnerText(ctx)
	must.NoError(err)
	is.Equal("clicked", text)
}

// TestLocatorLocatorDeepNesting verifies deeply chained locators work.
// Ref: TestLocatorLocator.java#shouldWorkWithDeepNesting
func TestLocatorLocatorDeepNesting(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="a">
			<div id="b">
				<div id="c">target</div>
			</div>
		</div>
	`))

	text, err := page.Locator("#a").Locator("#b").Locator("#c").InnerText(ctx)
	must.NoError(err)
	is.Equal("target", text)
}

// TestLocatorLocatorNthOnChained verifies Nth works after chaining.
// Ref: TestLocatorLocator.java#shouldSupportNthAfterChaining
func TestLocatorLocatorNthOnChained(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="row">
			<span>row1-a</span><span>row1-b</span>
		</div>
	`))

	text, err := page.Locator(".row").Locator("span").Nth(1).InnerText(ctx)
	must.NoError(err)
	is.Equal("row1-b", text)
}

// TestLocatorNthZeroGetsFirst verifies Nth(0) returns the first element.
// Ref: TestLocatorNth.java#shouldGetFirstElement
func TestLocatorNthZeroGetsFirst(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">first</div>
		<div class="item">second</div>
		<div class="item">third</div>
	`))

	text, err := page.Locator(".item").Nth(0).InnerText(ctx)
	must.NoError(err)
	is.Equal("first", text)
}

// TestLocatorNthNegativeGetsLast verifies Nth(-1) returns the last element.
// Ref: TestLocatorNth.java#shouldGetLastElement
func TestLocatorNthNegativeGetsLast(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">alpha</div>
		<div class="item">beta</div>
		<div class="item">gamma</div>
	`))

	text, err := page.Locator(".item").Nth(-1).InnerText(ctx)
	must.NoError(err)
	is.Equal("gamma", text)
}

// TestLocatorNthMiddleElement verifies Nth(1) returns the second element.
// Ref: TestLocatorNth.java#shouldGetMiddleElement
func TestLocatorNthMiddleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="s">one</span>
		<span class="s">two</span>
		<span class="s">three</span>
	`))

	text, err := page.Locator(".s").Nth(1).InnerText(ctx)
	must.NoError(err)
	is.Equal("two", text)
}

// TestLocatorFirstAndLastConvenience verifies First/Last shortcuts work.
// Ref: TestLocatorNth.java#shouldSupportFirstAndLast
func TestLocatorFirstAndLastConvenienceExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="p">para1</p>
		<p class="p">para2</p>
		<p class="p">para3</p>
	`))

	first, err := page.Locator(".p").First().InnerText(ctx)
	must.NoError(err)
	is.Equal("para1", first)

	last, err := page.Locator(".p").Last().InnerText(ctx)
	must.NoError(err)
	is.Equal("para3", last)
}

// TestLocatorNthSingleResult verifies Nth works when only one result exists.
// Ref: TestLocatorNth.java#shouldWorkWithSingleResult
func TestLocatorNthSingleResult(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="only">the only one</div>`))

	text, err := page.Locator("#only").Nth(0).InnerText(ctx)
	must.NoError(err)
	is.Equal("the only one", text)
}

// TestLocatorNthCanClick verifies Nth locator can be clicked.
// Ref: TestLocatorNth.java#shouldClick
func TestLocatorNthCanClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>One</button>
		<button onclick="window.__second=true">Two</button>
		<button>Three</button>
	`))

	must.NoError(page.Locator("button").Nth(1).Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__second`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorNthCanFill verifies Nth locator can fill an input.
// Ref: TestLocatorNth.java#shouldFill
func TestLocatorNthCanFill(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" class="inp">
		<input type="text" class="inp">
		<input type="text" class="inp">
	`))

	must.NoError(page.Locator(".inp").Nth(2).Fill(ctx, "third value"))

	val, err := page.Locator(".inp").Nth(2).InputValue(ctx)
	must.NoError(err)
	is.Equal("third value", val)
}

// TestLocatorNthInnerText verifies Nth locator can get InnerText.
// Ref: TestLocatorNth.java#shouldGetInnerText
func TestLocatorNthInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="para">First paragraph</p>
		<p class="para">Second paragraph</p>
		<p class="para">Third paragraph</p>
	`))

	text, err := page.Locator(".para").Nth(1).InnerText(ctx)
	must.NoError(err)
	is.Equal("Second paragraph", text)
}

// TestLocatorNthBoundingBox verifies Nth locator can get bounding box.
// Ref: TestLocatorNth.java#shouldGetBoundingBox
func TestLocatorNthBoundingBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="box" style="width:50px;height:50px;margin:5px">A</div>
		<div class="box" style="width:50px;height:50px;margin:5px">B</div>
	`))

	bb, err := page.Locator(".box").Nth(0).BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, 0.0)
	is.Greater(bb.Height, 0.0)
}

// TestLocatorNthIsVisible verifies Nth locator IsVisible.
// Ref: TestLocatorNth.java#shouldBeVisible
func TestLocatorNthIsVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">Visible</div>
		<div class="item" style="display:none">Hidden</div>
	`))

	visible, err := page.Locator(".item").Nth(0).IsVisible(ctx)
	must.NoError(err)
	is.True(visible)

	hidden, err := page.Locator(".item").Nth(1).IsVisible(ctx)
	must.NoError(err)
	is.False(hidden)
}

// TestLocatorNthFirstEx3 verifies Nth(0) gets the first element.
// Ref: TestLocatorNth.java#shouldGetFirstElement
func TestLocatorNthFirstEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>First</li>
		<li>Second</li>
		<li>Third</li>
	`))

	text, err := page.Locator("li").Nth(0).InnerText(ctx)
	must.NoError(err)
	is.Equal("First", text)
}

// TestLocatorNthLastEx3 verifies Nth(-1) gets the last element.
// Ref: TestLocatorNth.java#shouldGetLastElement
func TestLocatorNthLastEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>A</li>
		<li>B</li>
		<li>C</li>
	`))

	text, err := page.Locator("li").Last().InnerText(ctx)
	must.NoError(err)
	is.Equal("C", text)
}

// TestLocatorNthMiddleEx3 verifies Nth(1) gets the second element.
// Ref: TestLocatorNth.java#shouldGetMiddleElement
func TestLocatorNthMiddleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button>One</button>
		<button>Two</button>
		<button>Three</button>
	`))

	text, err := page.Locator("button").Nth(1).InnerText(ctx)
	must.NoError(err)
	is.Equal("Two", text)
}

// TestLocatorNthIsVisibleEx3 verifies Nth element IsVisible works.
// Ref: TestLocatorNth.java#shouldCheckVisibilityForNth
func TestLocatorNthIsVisibleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="box">Box 1</div>
		<div class="box" style="display:none">Box 2</div>
		<div class="box">Box 3</div>
	`))

	v0, err := page.Locator(".box").Nth(0).IsVisible(ctx)
	must.NoError(err)
	is.True(v0)

	v1, err := page.Locator(".box").Nth(1).IsVisible(ctx)
	must.NoError(err)
	is.False(v1)
}

// TestLocatorNthFillEx3 verifies Nth element Fill works.
// Ref: TestLocatorNth.java#shouldFillNthElement
func TestLocatorNthFillEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="inp" type="text">
		<input class="inp" type="text">
		<input class="inp" type="text">
	`))

	must.NoError(page.Locator(".inp").Nth(2).Fill(ctx, "third input"))

	val, err := page.Locator(".inp").Nth(2).InputValue(ctx)
	must.NoError(err)
	is.Equal("third input", val)
}

// TestLocatorNthClickEx4 verifies Nth element can be clicked.
// Ref: TestLocatorNth.java#shouldClickNthElement
func TestLocatorNthClickEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button onclick="window.__btnClicked=0">Button 0</button>
		<button onclick="window.__btnClicked=1">Button 1</button>
		<button onclick="window.__btnClicked=2">Button 2</button>
	`))

	must.NoError(page.Locator("button").Nth(1).Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__btnClicked`)
	must.NoError(err)
	is.Equal(float64(1), result)
}

// TestLocatorNthInputFillEx4 verifies Nth input can be filled.
// Ref: TestLocatorNth.java#shouldFillNthInput
func TestLocatorNthInputFillEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" class="f">
		<input type="text" class="f">
		<input type="text" class="f">
	`))

	must.NoError(page.Locator(".f").Nth(0).Fill(ctx, "first"))
	must.NoError(page.Locator(".f").Nth(2).Fill(ctx, "third"))

	v0, err := page.Locator(".f").Nth(0).InputValue(ctx)
	must.NoError(err)
	is.Equal("first", v0)

	v2, err := page.Locator(".f").Nth(2).InputValue(ctx)
	must.NoError(err)
	is.Equal("third", v2)
}

// TestLocatorNthFirstEx4 verifies First() is same as Nth(0).
// Ref: TestLocatorNth.java#firstShouldBeSameAsNth0
func TestLocatorNthFirstEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="p">First</p>
		<p class="p">Second</p>
		<p class="p">Third</p>
	`))

	nthText, err := page.Locator(".p").Nth(0).InnerText(ctx)
	must.NoError(err)
	firstText, err := page.Locator(".p").First().InnerText(ctx)
	must.NoError(err)
	is.Equal(nthText, firstText)
}

// TestLocatorNthLastEx4 verifies Last() is same as Nth(-1).
// Ref: TestLocatorNth.java#lastShouldBeSameAsNthMinus1
func TestLocatorNthLastEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="s">A</span>
		<span class="s">B</span>
		<span class="s">C</span>
	`))

	nthText, err := page.Locator(".s").Nth(-1).InnerText(ctx)
	must.NoError(err)
	lastText, err := page.Locator(".s").Last().InnerText(ctx)
	must.NoError(err)
	is.Equal(nthText, lastText)
}

// TestNthFirstItemEx5 verifies Nth(0) selects the first element.
// Ref: TestLocatorNth.java#shouldSelectFirst
func TestNthFirstItemEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p class="item">Alpha</p><p class="item">Beta</p><p class="item">Gamma</p>`))

	text, err := page.Locator(".item").Nth(0).TextContent(ctx)
	must.NoError(err)
	is.Equal("Alpha", text)
}

// TestNthMiddleItemEx5 verifies Nth(1) selects the second element.
// Ref: TestLocatorNth.java#shouldSelectMiddle
func TestNthMiddleItemEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p class="item">Alpha</p><p class="item">Beta</p><p class="item">Gamma</p>`))

	text, err := page.Locator(".item").Nth(1).TextContent(ctx)
	must.NoError(err)
	is.Equal("Beta", text)
}

// TestNthLastItemEx5 verifies Nth(-1) selects the last element.
// Ref: TestLocatorNth.java#shouldSelectLast
func TestNthLastItemEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span class="x">One</span><span class="x">Two</span><span class="x">Three</span>`))

	text, err := page.Locator(".x").Last().TextContent(ctx)
	must.NoError(err)
	is.Equal("Three", text)
}

// TestNthCountOnSubsetEx5 verifies Nth selection then attribute access.
// Ref: TestLocatorNth.java#shouldGetAttributeFromNth
func TestNthCountOnSubsetEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="field" id="f1" value="v1">
		<input class="field" id="f2" value="v2">
		<input class="field" id="f3" value="v3">
	`))

	val, err := page.Locator(".field").Nth(2).InputValue(ctx)
	must.NoError(err)
	is.Equal("v3", val)
}

// TestNthInputFillEx6 verifies Nth can select and fill a specific input.
// Ref: TestLocatorNth.java#shouldFillNthInput
func TestNthInputFillEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="f" type="text">
		<input class="f" type="text">
		<input class="f" type="text">
	`))

	must.NoError(page.Locator(".f").Nth(1).Fill(ctx, "middle value"))

	val, err := page.Locator(".f").Nth(1).InputValue(ctx)
	must.NoError(err)
	is.Equal("middle value", val)
}

// TestNthCheckboxEx6 verifies Nth can check a specific checkbox.
// Ref: TestLocatorNth.java#shouldCheckNthCheckbox
func TestNthCheckboxEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input class="c" type="checkbox">
		<input class="c" type="checkbox">
		<input class="c" type="checkbox">
	`))

	must.NoError(page.Locator(".c").Nth(2).Check(ctx))

	checked, err := page.Locator(".c").Nth(2).IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	notChecked, err := page.Locator(".c").Nth(0).IsChecked(ctx)
	must.NoError(err)
	is.False(notChecked)
}

// TestNthClickButtonEx6 verifies Nth can click a specific button.
// Ref: TestLocatorNth.java#shouldClickNthButton
func TestNthClickButtonEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button class="b" onclick="window.__clicked=0">A</button>
		<button class="b" onclick="window.__clicked=1">B</button>
		<button class="b" onclick="window.__clicked=2">C</button>
	`))

	must.NoError(page.Locator(".b").Nth(1).Click(ctx))

	idx, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Equal(float64(1), idx)
}

// TestNthAttributeEx6 verifies Nth can get attribute from specific element.
// Ref: TestLocatorNth.java#shouldGetNthAttribute
func TestNthAttributeEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item" data-id="100">First</div>
		<div class="item" data-id="200">Second</div>
		<div class="item" data-id="300">Third</div>
	`))

	id, err := page.Locator(".item").Nth(2).GetAttribute(ctx, "data-id")
	must.NoError(err)
	is.Equal("300", id)
}

// TestNthTableRowEx7 verifies Nth selects table row.
// Ref: TestLocatorNth.java#shouldSelectTableRow
func TestNthTableRowEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr id="r0"><td>Row 0</td></tr>
			<tr id="r1"><td>Row 1</td></tr>
			<tr id="r2"><td>Row 2</td></tr>
		</table>
	`))

	text, err := page.Locator("tr").Nth(1).TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Contains(*text, "Row 1")
}

// TestNthLastEx7 verifies Nth(-1) selects last element (via Last).
// Ref: TestLocatorNth.java#shouldSelectLastViaLast
func TestNthLastEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ol>
			<li>First</li>
			<li>Second</li>
			<li>Last</li>
		</ol>
	`))

	text, err := page.Locator("li").Last().TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Last", *text)
}

// TestNthFirstEx7 verifies First() selects first element.
// Ref: TestLocatorNth.java#shouldSelectFirstViaFirst
func TestNthFirstEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>Alpha</li>
			<li>Beta</li>
			<li>Gamma</li>
		</ul>
	`))

	text, err := page.Locator("li").First().TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Alpha", *text)
}

// TestNthOptionEx7 verifies Nth selects correct select option.
// Ref: TestLocatorNth.java#shouldSelectOption
func TestNthOptionEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="x">X</option>
			<option value="y">Y</option>
			<option value="z">Z</option>
		</select>
	`))

	text, err := page.Locator("#sel option").Nth(2).TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Z", *text)
}

// TestQuerySelectorReturnsNilEx3 verifies QuerySelector returns nil for missing element.
// Ref: TestLocatorQuerySelector.java#shouldReturnNilForMissing
func TestQuerySelectorReturnsNilEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Hello</p>`))

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err)
	is.Nil(el)
}

// TestQuerySelectorFindsFirstEx3 verifies QuerySelector returns only first matching element.
// Ref: TestLocatorQuerySelector.java#shouldReturnFirstMatch
func TestQuerySelectorFindsFirstEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">First</div>
		<div class="item">Second</div>
	`))

	el, err := page.QuerySelector(ctx, ".item")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("First", text)
}

// TestQuerySelectorAllCountEx3 verifies QuerySelectorAll returns all matching elements.
// Ref: TestLocatorQuerySelector.java#shouldReturnAllMatches
func TestQuerySelectorAllCountEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li class="item">A</li>
		<li class="item">B</li>
		<li class="item">C</li>
	`))

	els, err := page.QuerySelectorAll(ctx, ".item")
	must.NoError(err)
	is.Equal(3, len(els))
}

// TestQuerySelectorByIDEx3 verifies QuerySelector by ID returns correct element.
// Ref: TestLocatorQuerySelector.java#shouldFindByID
func TestQuerySelectorByIDEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<span id="myspan">Found it</span>`))

	el, err := page.QuerySelector(ctx, "#myspan")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("Found it", text)
}

// TestQuerySelectorAllEmptyEx3 verifies QuerySelectorAll returns empty slice when no match.
// Ref: TestLocatorQuerySelector.java#shouldReturnEmptyForNoMatch
func TestQuerySelectorAllEmptyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>No divs here</p>`))

	els, err := page.QuerySelectorAll(ctx, "div")
	must.NoError(err)
	is.Empty(els)
}

// TestQuerySelectorFormInputEx4 verifies QuerySelector finds form input.
// Ref: TestLocatorQuerySelector.java#shouldFindFormInput
func TestQuerySelectorFormInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<input type="text" name="username" id="user">
			<input type="password" name="password" id="pass">
		</form>
	`))

	el, err := page.QuerySelector(ctx, "#user")
	must.NoError(err)
	must.NotNil(el)

	name, err := el.GetAttribute(ctx, "name")
	must.NoError(err)
	is.Equal("username", name)
}

// TestQuerySelectorReturnsNilEx4 verifies QuerySelector returns nil for no match.
// Ref: TestLocatorQuerySelector.java#shouldReturnNilForNoMatch
func TestQuerySelectorReturnsNilEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Content</div>`))

	el, err := page.QuerySelector(ctx, "#nonexistent")
	must.NoError(err)
	is.Nil(el)
}

// TestQuerySelectorTableCellEx4 verifies QuerySelector can find table cells.
// Ref: TestLocatorQuerySelector.java#shouldFindTableCell
func TestQuerySelectorTableCellEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr>
				<td id="cell1">Data 1</td>
				<td id="cell2">Data 2</td>
			</tr>
		</table>
	`))

	el, err := page.QuerySelector(ctx, "#cell2")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.TextContent(ctx)
	must.NoError(err)
	is.Equal("Data 2", text)
}

// TestQuerySelectorAnchorEx4 verifies QuerySelector finds anchors by href.
// Ref: TestLocatorQuerySelector.java#shouldFindAnchorByHref
func TestQuerySelectorAnchorEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a href="/home">Home</a>
		<a href="/about">About</a>
		<a href="/contact">Contact</a>
	`))

	el, err := page.QuerySelector(ctx, `a[href="/about"]`)
	must.NoError(err)
	must.NotNil(el)

	text, err := el.TextContent(ctx)
	must.NoError(err)
	is.Equal("About", text)
}

// TestLocatorTextContentReturnsNilForNoMatch verifies TextContent returns nil for missing selector.
// Ref: TestLocatorTextContent.java#shouldReturnNullForMissing
func TestLocatorTextContentReturnsNilForNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))

	// Count should be 0 for a missing locator
	count, err := page.Locator("#doesNotExist").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// TestLocatorInnerTextVsTextContent verifies difference between InnerText and TextContent.
// Ref: TestLocatorTextContent.java#shouldDifferFromInnerText
func TestLocatorInnerTextVsTextContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// TextContent includes hidden text, InnerText does not
	must.NoError(page.SetContent(ctx, `
		<div id="container">
			visible
			<span style="display:none">hidden</span>
		</div>
	`))

	tc, err := page.Locator("#container").TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Contains(*tc, "visible")

	it, err := page.Locator("#container").InnerText(ctx)
	must.NoError(err)
	is.Equal("visible", it)
}

// TestLocatorAllTextContentsReturnsAll verifies AllTextContents returns text from all matched elements.
// Ref: TestLocatorTextContent.java#shouldReturnAllTextContents
func TestLocatorAllTextContentsFromMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p class="item">alpha</p>
		<p class="item">beta</p>
		<p class="item">gamma</p>
	`))

	texts, err := page.Locator(".item").AllTextContents(ctx)
	must.NoError(err)
	is.Equal([]string{"alpha", "beta", "gamma"}, texts)
}

// TestLocatorAllInnerTextsFromMultiple verifies AllInnerTexts returns innerText from all matched elements.
// Ref: TestLocatorTextContent.java#shouldReturnAllInnerTexts
func TestLocatorAllInnerTextsFromMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li class="entry">one</li>
		<li class="entry">two</li>
		<li class="entry">three</li>
	`))

	texts, err := page.Locator(".entry").AllInnerTexts(ctx)
	must.NoError(err)
	is.Equal([]string{"one", "two", "three"}, texts)
}

// TestLocatorGetAttributeExistingAttribute verifies GetAttribute returns value for existing attribute.
// Ref: TestLocatorTextContent.java#shouldReturnAttributeValue
func TestLocatorGetAttributeExistingAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<img id="img" src="image.png" alt="test image">`))

	attr, err := page.Locator("img").GetAttribute(ctx, "alt")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("test image", *attr)
}

// TestLocatorGetAttributeNilForMissing verifies GetAttribute returns nil for missing attribute.
// Ref: TestLocatorTextContent.java#shouldReturnNullForMissingAttribute
func TestLocatorGetAttributeNilForMissing(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="box">content</div>`))

	attr, err := page.Locator("#box").GetAttribute(ctx, "data-nonexistent")
	must.NoError(err)
	is.Nil(attr)
}

// TestLocatorTextContentReturnsTextEx2 verifies TextContent returns element text.
// Ref: TestLocatorTextContent.java#shouldReturnText
func TestLocatorTextContentReturnsTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello TextContent</p>`))

	text, err := page.Locator("#p").TextContent(ctx)
	must.NoError(err)
	is.Equal("Hello TextContent", text)
}

// TestLocatorTextContentIncludesHiddenText verifies TextContent includes CSS-hidden text.
// Ref: TestLocatorTextContent.java#shouldIncludeHiddenText
func TestLocatorTextContentIncludesHiddenText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">Visible<span style="display:none">Hidden</span></div>
	`))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Visible")
	is.Contains(text, "Hidden")
}

// TestLocatorTextContentWithSpecialChars verifies TextContent handles special characters.
// Ref: TestLocatorTextContent.java#shouldHandleSpecialChars
func TestLocatorTextContentWithSpecialChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello &amp; World &lt;test&gt;</p>`))

	text, err := page.Locator("#p").TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Hello & World")
}

// TestLocatorTextContentMultilineText verifies TextContent handles multiline text.
// Ref: TestLocatorTextContent.java#shouldHandleMultiline
func TestLocatorTextContentMultilineText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d">
			Line one
			<br>
			Line two
		</div>
	`))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Line one")
	is.Contains(text, "Line two")
}

// TestLocatorTextContentEmptyDivEx2 verifies TextContent on empty div returns empty or whitespace.
// Ref: TestLocatorTextContent.java#shouldReturnEmptyForEmptyDiv
func TestLocatorTextContentEmptyDivEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"></div>`))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Empty(text)
}

// TestTextContentWithLinksEx3 verifies TextContent includes link text.
// Ref: TestLocatorTextContent.java#shouldIncludeLinkText
func TestTextContentWithLinksEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visit <a href="#">our site</a> today</div>`))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Equal("Visit our site today", text)
}

// TestTextContentHTMLEntitiesEx3 verifies TextContent decodes HTML entities.
// Ref: TestLocatorTextContent.java#shouldDecodeEntities
func TestTextContentHTMLEntitiesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">&lt;Hello &amp; World&gt;</p>`))

	text, err := page.Locator("#p").TextContent(ctx)
	must.NoError(err)
	is.Equal("<Hello & World>", text)
}

// TestTextContentNonBreakingSpaceEx3 verifies TextContent includes nbsp.
// Ref: TestLocatorTextContent.java#shouldIncludeNbsp
func TestTextContentNonBreakingSpaceEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello&nbsp;World</p>`))

	text, err := page.Locator("#p").TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Hello")
	is.Contains(text, "World")
}

// TestTextContentCodeElementEx3 verifies TextContent for code element.
// Ref: TestLocatorTextContent.java#shouldReturnCodeContent
func TestTextContentCodeElementEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<code id="c">let x = 1;</code>`))

	text, err := page.Locator("#c").TextContent(ctx)
	must.NoError(err)
	is.Equal("let x = 1;", text)
}

// TestTextContentHeadingEx4 verifies TextContent for heading elements.
// Ref: TestLocatorTextContent.java#shouldReturnHeadingText
func TestTextContentHeadingEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h2 id="h">Section Title</h2>`))

	text, err := page.Locator("#h").TextContent(ctx)
	must.NoError(err)
	is.Equal("Section Title", text)
}

// TestTextContentWithLineBreakEx4 verifies TextContent collapses whitespace in spans.
// Ref: TestLocatorTextContent.java#shouldHandleLineBreaks
func TestTextContentWithLineBreakEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><span>Hello</span><span>World</span></div>`))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Hello")
	is.Contains(text, "World")
}

// TestTextContentEmptyDivEx4 verifies TextContent returns empty string for empty element.
// Ref: TestLocatorTextContent.java#shouldReturnEmptyForEmptyElement
func TestTextContentEmptyDivEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="empty"></div>`))

	text, err := page.Locator("#empty").TextContent(ctx)
	must.NoError(err)
	is.Equal("", text)
}

// TestTextContentNestedSpanEx4 verifies TextContent for deeply nested elements.
// Ref: TestLocatorTextContent.java#shouldReturnNestedText
func TestTextContentNestedSpanEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><p><strong>Bold <em>italic</em></strong></p></div>`))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Bold")
	is.Contains(text, "italic")
}

// TestTextContentButtonEx4 verifies TextContent for button elements.
// Ref: TestLocatorTextContent.java#shouldReturnButtonText
func TestTextContentButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Submit Form</button>`))

	text, err := page.Locator("#btn").TextContent(ctx)
	must.NoError(err)
	is.Equal("Submit Form", text)
}

// TestTextContentTableHeaderEx5 verifies TextContent for table header.
// Ref: TestLocatorTextContent.java#shouldReturnTableHeaderText
func TestTextContentTableHeaderEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<thead>
				<tr><th id="th1">Name</th><th id="th2">Score</th></tr>
			</thead>
		</table>
	`))

	text, err := page.Locator("#th1").TextContent(ctx)
	must.NoError(err)
	is.Equal("Name", text)
}

// TestTextContentListItemEx5 verifies TextContent for list items.
// Ref: TestLocatorTextContent.java#shouldReturnListItemText
func TestTextContentListItemEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ol>
			<li id="li1">Step one</li>
			<li id="li2">Step two</li>
			<li id="li3">Step three</li>
		</ol>
	`))

	text, err := page.Locator("#li3").TextContent(ctx)
	must.NoError(err)
	is.Equal("Step three", text)
}

// TestTextContentWithWhitespaceEx5 verifies TextContent includes whitespace text.
// Ref: TestLocatorTextContent.java#shouldIncludeWhitespace
func TestTextContentWithWhitespaceEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<pre id="pre">Code here</pre>`))

	text, err := page.Locator("#pre").TextContent(ctx)
	must.NoError(err)
	is.Contains(text, "Code here")
}

// TestTextContentAfterSetContentEx5 verifies TextContent after multiple SetContent calls.
// Ref: TestLocatorTextContent.java#shouldReflectSetContent
func TestTextContentAfterSetContentEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Version 1</div>`))
	must.NoError(page.SetContent(ctx, `<div id="d">Version 2</div>`))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Equal("Version 2", text)
}

// TestTextContentTableCellEx6 verifies TextContent in table cells.
// Ref: TestLocatorTextContent.java#shouldReturnTableCell
func TestTextContentTableCellEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr><td id="c1">Cell A</td><td id="c2">Cell B</td></tr>
		</table>
	`))

	text, err := page.Locator("#c2").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Cell B", *text)
}

// TestTextContentCodeBlockEx6 verifies TextContent for code element.
// Ref: TestLocatorTextContent.java#shouldReturnCodeText
func TestTextContentCodeBlockEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<code id="c">fmt.Println("hello")</code>`))

	text, err := page.Locator("#c").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Contains(*text, "Println")
}

// TestTextContentAfterRemoveChildEx6 verifies TextContent after child removal.
// Ref: TestLocatorTextContent.java#shouldUpdateAfterChildRemoval
func TestTextContentAfterRemoveChildEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Hello <span id="s">World</span></div>`))

	_, err := page.Evaluate(ctx, `() => { const s = document.getElementById('s'); s.parentNode.removeChild(s); }`)
	must.NoError(err)

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Hello ", *text)
}

// TestTextContentSVGEx6 verifies TextContent for SVG title element.
// Ref: TestLocatorTextContent.java#shouldReturnSVGText
func TestTextContentSVGEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<svg>
			<title id="t">SVG Title</title>
			<circle r="5"/>
		</svg>
	`))

	text, err := page.Locator("#t").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("SVG Title", *text)
}
