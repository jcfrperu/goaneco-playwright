//go:build e2e

// E2E tests for: TestLocatorConvenience.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorTextContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">hello <b>world</b></div>`)
	must.NoError(err, "SetContent failed")

	tc, err := page.Locator("#el").TextContent(ctx)
	must.NoError(err, "TextContent failed")
	if tc == nil || *tc != "hello world" {
		t.Errorf("TextContent = %v, want 'hello world'", tc)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorInnerHTML(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el"><b>hello</b></div>`)
	must.NoError(err, "SetContent failed")

	html, err := page.Locator("#el").InnerHTML(ctx)
	must.NoError(err, "InnerHTML failed")
	is.Equal("<b>hello</b>", html)
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">hello</div><div id="hidden" style="display:none">hidden</div>`)
	must.NoError(err, "SetContent failed")

	text, err := page.Locator("#el").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	is.Equal("hello", text)
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorInputValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="in" value="initial"/>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator("#in").InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("initial", val)
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="ed"/><input id="ro" readonly/>`)
	must.NoError(err, "SetContent failed")

	editable, err := page.Locator("#ed").IsEditable(ctx)
	must.NoError(err)
	is.True(editable, "expected #ed to be editable")

	readonly, err := page.Locator("#ro").IsEditable(ctx)
	must.NoError(err)
	is.False(readonly, "expected #ro to NOT be editable")
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox" id="ch" checked/><input type="checkbox" id="un"/>`)
	must.NoError(err, "SetContent failed")

	checked, err := page.Locator("#ch").IsChecked(ctx)
	must.NoError(err)
	is.True(checked, "expected #ch to be checked")

	unchecked, err := page.Locator("#un").IsChecked(ctx)
	must.NoError(err)
	is.False(unchecked, "expected #un to NOT be checked")
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorGetAttributeConvenience(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" data-value="42">text</div>`)
	must.NoError(err, "SetContent failed")

	val, err := page.Locator("#el").GetAttribute(ctx, "data-value")
	must.NoError(err, "GetAttribute failed")
	if val == nil || *val != "42" {
		t.Errorf("GetAttribute(data-value) = %v, want '42'", val)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorAllShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul><li>a</li><li>b</li><li>c</li></ul>`)
	must.NoError(err, "SetContent failed")

	items, err := page.Locator("li").All(ctx)
	must.NoError(err, "All failed")
	is.Len(items, 3)

	var texts []string
	for _, loc := range items {
		tc, err := loc.TextContent(ctx)
		if err != nil || tc == nil {
			t.Fatalf("TextContent on nth locator failed: %v", err)
		}
		texts = append(texts, *tc)
	}
	if texts[0] != "a" || texts[1] != "b" || texts[2] != "c" {
		t.Errorf("unexpected texts from All(): %v", texts)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorAllEmptyWhenNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>nothing</div>`)
	must.NoError(err, "SetContent failed")

	items, err := page.Locator("span").All(ctx)
	must.NoError(err, "All failed")
	if len(items) != 0 {
		t.Errorf("expected empty All() for non-matching selector, got %d", len(items))
	}
}

// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorHasLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="parent"><span class="child">yes</span></div><div class="parent">no</div>`)
	must.NoError(err, "SetContent failed")

	loc := page.Locator(".parent").Filter(&playwright.LocatorFilterOptions{
		Has: page.Locator(".child"),
	})
	count, err := loc.Count(ctx)
	must.NoError(err, "Count failed")
	is.Equal(1, count)
}

// TestLocatorIsEnabledTrueForInputEx2 verifies IsEnabled returns true for normal input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEnabledTrueForInputEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestLocatorIsEnabledFalseForDisabledInputEx2 verifies IsEnabled returns false for disabled input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEnabledFalseForDisabledInputEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" disabled>`))

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestLocatorIsEnabledTrueForButtonEx2 verifies IsEnabled returns true for normal button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEnabledTrueForButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	enabled, err := page.Locator("#btn").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestLocatorIsEnabledFalseForDisabledButtonEx2 verifies IsEnabled returns false for disabled button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEnabledFalseForDisabledButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Click</button>`))

	enabled, err := page.Locator("#btn").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestLocatorIsEnabledSelectEx2 verifies IsEnabled returns true for select.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEnabledSelectEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel"><option value="a">A</option></select>`))

	enabled, err := page.Locator("#sel").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledSelectEx3 verifies IsEnabled for select element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledSelectEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s">
			<option>A</option>
		</select>
	`))

	enabled, err := page.Locator("#s").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledSelectEx3 verifies IsEnabled returns false for disabled select.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledDisabledSelectEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="s" disabled>
			<option>A</option>
		</select>
	`))

	enabled, err := page.Locator("#s").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledTextareaEx3 verifies IsEnabled for textarea.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledTextareaEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">text</textarea>`))

	enabled, err := page.Locator("#ta").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledTextareaEx3 verifies IsEnabled returns false for disabled textarea.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledDisabledTextareaEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta" disabled>text</textarea>`))

	enabled, err := page.Locator("#ta").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledAfterEnableEx3 verifies IsEnabled reflects dynamic enable change.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledAfterEnableEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" disabled>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('inp').disabled = false`)
	must.NoError(err)

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledSelectEx4 verifies IsEnabled for select element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledSelectEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel"><option>One</option></select>`))

	enabled, err := page.Locator("#sel").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledSelectEx4 verifies IsEnabled returns false for disabled select.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledDisabledSelectEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<select id="sel" disabled><option>One</option></select>`))

	enabled, err := page.Locator("#sel").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledTextareaEx4 verifies IsEnabled for textarea element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledTextareaEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">Content</textarea>`))

	enabled, err := page.Locator("#ta").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledAfterJSEnableEx4 verifies IsEnabled reflects JS-enabled element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledAfterJSEnableEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" disabled>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('inp').disabled = false`)
	must.NoError(err)

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledAfterJSEnableFieldsetEx5 verifies IsEnabled after JS enables fieldset child.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledAfterJSEnableFieldsetEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<fieldset disabled>
			<input id="inp" type="text">
		</fieldset>
	`))

	enabled, err := page.Locator("#inp").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestIsEnabledButtonWithTypeEx5 verifies IsEnabled for submit button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledButtonWithTypeEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="submit" type="submit">Submit</button>`))

	enabled, err := page.Locator("#submit").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestIsEnabledDisabledButtonWithTypeEx5 verifies IsEnabled for disabled submit button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsEnabledDisabledButtonWithTypeEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="submit" type="submit" disabled>Submit</button>`))

	enabled, err := page.Locator("#submit").IsEnabled(ctx)
	must.NoError(err)
	is.False(enabled)
}

// TestLocatorConvenienceGetAttributeShouldWork verifies getAttribute works on dom.html elements.
// Ref: TestLocatorConvenience.java#getAttributeShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceGetAttributeShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	val, err := page.Locator("#outer").GetAttribute(ctx, "name")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("value", *val)

	missing, err := page.Locator("#outer").GetAttribute(ctx, "foo")
	must.NoError(err)
	is.Nil(missing, "getAttribute for non-existent attribute should return nil")
}

// TestLocatorConvenienceInnerHTMLShouldWork verifies innerHTML on the #outer div from dom.html.
// Ref: TestLocatorConvenience.java#innerHTMLShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceInnerHTMLShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	html, err := page.Locator("#outer").InnerHTML(ctx)
	must.NoError(err)
	is.Contains(html, `id="inner"`)
	is.Contains(html, "Text,")
}

// TestLocatorConvenienceInnerTextShouldWork verifies innerText on the #inner div from dom.html.
// Ref: TestLocatorConvenience.java#innerTextShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceInnerTextShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	text, err := page.Locator("#inner").InnerText(ctx)
	must.NoError(err)
	is.Equal("Text, more text", text)
}

// TestLocatorConvenienceTextContentShouldWork verifies textContent on the #inner div from dom.html.
// Ref: TestLocatorConvenience.java#textContentShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceTextContentShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	tc, err := page.Locator("#inner").TextContent(ctx)
	must.NoError(err)
	must.NotNil(tc)
	is.Contains(*tc, "Text,")
	is.Contains(*tc, "more text")
}

// TestLocatorConvenienceInputValueShouldWork verifies inputValue on #input and #textarea from dom.html.
// Ref: TestLocatorConvenience.java#inputValueShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceInputValueShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	// Fill textarea and verify InputValue
	must.NoError(page.Locator("#textarea").Fill(ctx, "text value"))
	val, err := page.Locator("#textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("text value", val)

	// Fill input and verify InputValue
	must.NoError(page.Locator("#input").Fill(ctx, "input value"))
	val2, err := page.Locator("#input").InputValue(ctx)
	must.NoError(err)
	is.Equal("input value", val2)
}

// TestLocatorConvenienceInnerTextShouldThrow verifies that InnerText on an SVG element returns an error.
// Ref: TestLocatorConvenience.java#innerTextShouldThrow  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceInnerTextShouldThrow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<svg>text</svg>`))

	_, err := page.Locator("svg").InnerText(ctx)
	is.Error(err, "expected error for InnerText on SVG")
	is.Contains(err.Error(), "HTMLElement")
}

// TestLocatorConvenienceIsCheckedShouldWork verifies isChecked using the #check element in dom.html.
// Ref: TestLocatorConvenience.java#isCheckedShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceIsCheckedShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	checked, err := page.Locator("#check").IsChecked(ctx)
	must.NoError(err)
	is.True(checked, "#check input should be checked by default")

	// Uncheck it via evaluate and verify again
	_, err = page.Locator("#check").Evaluate(ctx, "input => input.checked = false")
	must.NoError(err)

	checked2, err := page.Locator("#check").IsChecked(ctx)
	must.NoError(err)
	is.False(checked2, "#check input should now be unchecked")
}

// TestLocatorConvenienceIsEditableShouldWork verifies isEditable on #input and #check from dom.html.
// Ref: TestLocatorConvenience.java#isEditableShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceIsEditableShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.DomPage()))

	// #textarea with readOnly = true should not be editable
	_, err := page.EvalOnSelector(ctx, "#textarea", "t => t.readOnly = true")
	must.NoError(err)

	editable, err := page.Locator("#input").IsEditable(ctx)
	must.NoError(err)
	is.True(editable, "#input should be editable")

	editable2, err := page.Locator("#textarea").IsEditable(ctx)
	must.NoError(err)
	is.False(editable2, "#textarea (readOnly) should not be editable")
}

// Ref: TestLocatorConvenience.java#allTextContentsShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceAllTextContentsShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<ul><li>one</li><li>two</li><li>three</li></ul>`)
	must.NoError(err, "SetContent failed")

	texts, err := page.Locator("li").AllTextContents(ctx)
	must.NoError(err, "AllTextContents failed")
	is.Len(texts, 3)
	if texts[0] != "one" || texts[1] != "two" || texts[2] != "three" {
		t.Errorf("unexpected texts: %v", texts)
	}
}

// Ref: TestLocatorConvenience.java#allInnerTextsShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceAllInnerTextsShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="item">alpha</div><div class="item">beta</div>`)
	must.NoError(err, "SetContent failed")

	texts, err := page.Locator(".item").AllInnerTexts(ctx)
	must.NoError(err, "AllInnerTexts failed")
	is.Len(texts, 2)
	if texts[0] != "alpha" || texts[1] != "beta" {
		t.Errorf("unexpected texts: %v", texts)
	}
}

// TestLocatorInnerHTMLReturnsNestedHTML verifies InnerHTML returns full nested HTML.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
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

// TestLocatorInnerHTMLAttributesEx4 verifies InnerHTML includes attributes.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
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

// TestLocatorFirstAndLastConvenience verifies First/Last shortcuts work.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
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

// TestLocatorInnerTextVsTextContent verifies difference between InnerText and TextContent.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
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

// TestLocatorIsVisibleFalseForDisplayNoneEx2 verifies IsVisible returns false for display:none.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsVisibleFalseForDisplayNoneEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="display:none">Hidden</p>`))

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsVisibleFalseForVisibilityHiddenEx2 verifies IsVisible returns false for visibility:hidden.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsVisibleFalseForVisibilityHiddenEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="visibility:hidden">Hidden</p>`))

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsVisibleTrueAfterShowEx2 verifies IsVisible returns true after making element visible.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsVisibleTrueAfterShowEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="display:none">Initially hidden</p>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('p').style.display = 'block'`)
	must.NoError(err)

	visible, err := page.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleButton verifies IsVisible returns true for button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsVisibleButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click me</button>`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleAfterShowEx3 verifies IsVisible returns true after display change.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleAfterShowEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none;">Hidden</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').style.display = ''`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleAfterHideEx3 verifies IsVisible returns false after hiding.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleAfterHideEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').style.visibility = 'hidden'`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleButtonEx3 verifies IsVisible for button elements.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleButtonEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click me</button>`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleOpacityZeroEx3 verifies IsVisible for opacity:0 element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleOpacityZeroEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="opacity:0;">Invisible</div>`))

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleAfterShowEx4 verifies IsVisible reflects JS-shown element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleAfterShowEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none">Hidden</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').style.display = 'block'`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleButtonEx4 verifies IsVisible for visible button.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleHiddenByOpacityEx4 verifies IsVisible for opacity-zero element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleHiddenByOpacityEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="opacity:0">Invisible</div>`))

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleLinkEx4 verifies IsVisible for anchor element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleLinkEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="#">Link</a>`))

	visible, err := page.Locator("#a").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleDisplayNoneEx verifies IsVisible returns false for display:none.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleDisplayNoneEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none">Hidden</div>`))

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestIsVisibleAfterClassRemoveEx verifies IsVisible after removing hidden class.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleAfterClassRemoveEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>.hidden { display:none; }</style>
		<div id="d" class="hidden">Content</div>
	`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').classList.remove('hidden')`)
	must.NoError(err)

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleSelectElementEx verifies IsVisible returns true for select.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleSelectElementEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">A</option>
		</select>
	`))

	visible, err := page.Locator("#sel").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleImageEx verifies IsVisible returns true for img element.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleImageEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img id="img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			style="width:20px;height:20px;">
	`))

	visible, err := page.Locator("#img").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestIsVisibleButtonInHiddenContainerEx verifies IsVisible for button inside hidden div.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestIsVisibleButtonInHiddenContainerEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="display:none">
			<button id="btn">Hidden button</button>
		</div>
	`))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorIsEditableFalseForReadOnly verifies IsEditable returns false for readonly input.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEditableFalseForReadOnly(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" readonly>`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.False(editable)
}

// TestLocatorIsVisibleAfterStyleChange verifies IsVisible reflects style changes.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsVisibleAfterStyleChange(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" style="display:none">content</div>`))

	hidden, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.False(hidden)

	_, err = page.Evaluate(ctx, `() => document.getElementById('el').style.display = 'block'`)
	must.NoError(err)

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsVisibleForOpacityZero verifies opacity:0 element is visible.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsVisibleForOpacityZero(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// opacity:0 elements are visible (take up space)
	must.NoError(page.SetContent(ctx, `<div id="el" style="opacity:0;width:100px;height:100px">text</div>`))

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorIsEnabledForNonFormElement verifies IsEnabled for div is true.
// Extra: no equivalent method in playwright-java; related class: TestLocatorConvenience.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorIsEnabledForNonFormElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">content</div>`))

	enabled, err := page.Locator("#el").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// Ref: TestLocatorConvenience.java#shouldHaveANicePreview  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceShouldHaveANicePreview(t *testing.T) {
	t.Skip("not implemented: Locator.toString / JSHandle.toString (preview)")
}

// Ref: TestLocatorConvenience.java#isVisibleAndIsHiddenShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceIsVisibleAndIsHiddenShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>Hi</div><span></span>"))

	// Locator.IsHidden is not available; hidden state is checked through the ToBeHidden assertion.
	div := page.Locator("div")
	v, err := div.IsVisible(ctx)
	must.NoError(err)
	is.True(v)
	must.NoError(playwright.Expect(div).Not().ToBeHidden(ctx))

	span := page.Locator("span")
	v, err = span.IsVisible(ctx)
	must.NoError(err)
	is.False(v)
	must.NoError(playwright.Expect(span).ToBeHidden(ctx))

	v, err = page.Locator("no-such-element").IsVisible(ctx)
	must.NoError(err)
	is.False(v)
	must.NoError(playwright.Expect(page.Locator("no-such-element")).ToBeHidden(ctx))
}

// Ref: TestLocatorConvenience.java#isEnabledAndIsDisabledShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceIsEnabledAndIsDisabledShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<button disabled>button1</button>\n<button>button2</button>\n<div>div</div>"))
	check := func(selector string, enabled bool) {
		l := page.Locator(selector)
		e, err := l.IsEnabled(ctx)
		must.NoError(err, selector)
		is.Equal(enabled, e, selector)
		d, err := l.IsDisabled(ctx)
		must.NoError(err, selector)
		is.Equal(!enabled, d, selector)
	}
	check("div", true)
	check(":text('button1')", false)
	check(":text('button2')", true)
}

// Ref: TestLocatorConvenience.java#descriptionShouldReturnNullForUndescribedLocators  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceDescriptionShouldReturnNullForUndescribedLocators(t *testing.T) {
	t.Skip("not implemented: Locator.Description")
}

// Ref: TestLocatorConvenience.java#descriptionShouldReturnSimpleDescription  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceDescriptionShouldReturnSimpleDescription(t *testing.T) {
	t.Skip("not implemented: Locator.Describe, Locator.Description")
}

// Ref: TestLocatorConvenience.java#descriptionShouldHandleSpecialCharacters  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceDescriptionShouldHandleSpecialCharacters(t *testing.T) {
	t.Skip("not implemented: Locator.Describe, Locator.Description")
}

// Ref: TestLocatorConvenience.java#descriptionShouldWorkWithChainedLocators  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceDescriptionShouldWorkWithChainedLocators(t *testing.T) {
	t.Skip("not implemented: Locator.Describe, Locator.Description")
}

// Ref: TestLocatorConvenience.java#descriptionShouldReturnLastDescriptionForMultipleCalls  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceDescriptionShouldReturnLastDescriptionForMultipleCalls(t *testing.T) {
	t.Skip("not implemented: Locator.Describe, Locator.Description")
}

// Ref: TestLocatorConvenience.java#toStringShouldReturnFormattedLocator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceToStringShouldReturnFormattedLocator(t *testing.T) {
	t.Skip("not implemented: Locator.String (toString)")
}

// Ref: TestLocatorConvenience.java#toStringShouldPreferDescription  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorConvenience.java
func TestLocatorConvenienceToStringShouldPreferDescription(t *testing.T) {
	t.Skip("not implemented: Locator.Describe, Locator.String (toString)")
}
