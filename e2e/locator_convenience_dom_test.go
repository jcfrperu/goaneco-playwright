//go:build e2e

// Locator convenience tests using dom.html.
// Migration of: TestLocatorConvenience.java (getAttribute, inputValue, innerHTML,
// innerText, textContent, isChecked, isEditable, isVisible tests using dom.html page)
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorGetAttributeWithDomPage verifies getAttribute works on dom.html elements.
// Ref: TestLocatorConvenience.java#getAttributeShouldWork
func TestLocatorGetAttributeWithDomPage(t *testing.T) {
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

// TestLocatorInnerHTMLWithDomPage verifies innerHTML on the #outer div from dom.html.
// Ref: TestLocatorConvenience.java#innerHTMLShouldWork
func TestLocatorInnerHTMLWithDomPage(t *testing.T) {
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

// TestLocatorInnerTextWithDomPage verifies innerText on the #inner div from dom.html.
// Ref: TestLocatorConvenience.java#innerTextShouldWork
func TestLocatorInnerTextWithDomPage(t *testing.T) {
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

// TestLocatorTextContentWithDomPage verifies textContent on the #inner div from dom.html.
// Ref: TestLocatorConvenience.java#textContentShouldWork
func TestLocatorTextContentWithDomPage(t *testing.T) {
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

// TestLocatorInputValueWithDomPage verifies inputValue on #input and #textarea from dom.html.
// Ref: TestLocatorConvenience.java#inputValueShouldWork
func TestLocatorInputValueWithDomPage(t *testing.T) {
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

// TestLocatorInnerTextShouldThrowOnSVG verifies that InnerText on an SVG element returns an error.
// Ref: TestLocatorConvenience.java#innerTextShouldThrow
func TestLocatorInnerTextShouldThrowOnSVG(t *testing.T) {
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

// TestLocatorIsCheckedShouldWorkWithDomPage verifies isChecked using the #check element in dom.html.
// Ref: TestLocatorConvenience.java#isCheckedShouldWork
func TestLocatorIsCheckedShouldWorkWithDomPage(t *testing.T) {
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

// TestLocatorIsEditableShouldWorkWithDomPage verifies isEditable on #input and #check from dom.html.
// Ref: TestLocatorConvenience.java#isEditableShouldWork
func TestLocatorIsEditableShouldWorkWithDomPage(t *testing.T) {
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
