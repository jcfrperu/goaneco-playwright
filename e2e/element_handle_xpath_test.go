//go:build e2e

// Ref: TestElementHandleQuerySelector.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestElementHandleXpathQuerySelectorAll verifies ElementHandle.QuerySelectorAll with an XPath selector.
// Ref: TestElementHandleQuerySelector.java#xpathShouldQueryExistingElement
func TestElementHandleXpathQuerySelectorAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><div class="second"><div class="inner">A</div></div></body></html>`))

	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)

	second, err := html.QuerySelectorAll(ctx, `xpath=./body/div[contains(@class, 'second')]`)
	must.NoError(err)
	is.Len(second, 1)

	inner, err := second[0].QuerySelectorAll(ctx, `xpath=./div[contains(@class, 'inner')]`)
	must.NoError(err)
	is.Len(inner, 1)

	text, err := inner[0].TextContent(ctx)
	must.NoError(err)
	is.Equal("A", text)
}

// TestElementHandleXpathQuerySelectorAllEmpty verifies QuerySelectorAll with XPath returns empty for no matches.
// Ref: TestElementHandleQuerySelector.java#xpathShouldReturnNullForNonExistingElement
func TestElementHandleXpathQuerySelectorAllEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><div class="second"><div class="inner">B</div></div></body></html>`))

	html, err := page.QuerySelector(ctx, "html")
	must.NoError(err)
	must.NotNil(html)

	elements, err := html.QuerySelectorAll(ctx, `xpath=/div[contains(@class, 'third')]`)
	must.NoError(err)
	is.Empty(elements)
}

// TestElementHandleWorkForAdoptedElements verifies ElementHandle remains valid after cross-document adoption.
// Ref: TestElementHandleQuerySelector.java#shouldWorkForAdoptedElements
func TestElementHandleWorkForAdoptedElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup := waitForPopup(t, page, func() {
		_, err := page.Evaluate(ctx, `url => window['__popup'] = window.open(url)`, srv.EmptyPage())
		must.NoError(err)
	})
	must.NotNil(popup)

	// Create a div with a span in the main page.
	_, err := page.Evaluate(ctx, `() => {
		const div = document.createElement('div');
		document.body.appendChild(div);
		const span = document.createElement('span');
		span.textContent = 'hello';
		div.appendChild(span);
	}`)
	must.NoError(err)

	elHandle, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(elHandle)

	// Span is accessible before adoption.
	span, err := elHandle.QuerySelector(ctx, "span")
	must.NoError(err)
	must.NotNil(span)

	text, err := span.TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text)

	// Adopt the div into the popup's document.
	must.NoError(popup.WaitForLoadState(ctx, "domcontentloaded"))
	_, err = page.Evaluate(ctx, `() => {
		const div = document.querySelector('div');
		window['__popup'].document.body.appendChild(div);
	}`)
	must.NoError(err)

	// ElementHandle still resolves the span after adoption.
	span2, err := elHandle.QuerySelector(ctx, "span")
	must.NoError(err)
	must.NotNil(span2)

	text2, err := span2.TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text2)

	// The span is also queryable from the popup.
	spanInPopup, err := popup.QuerySelector(ctx, "span")
	must.NoError(err)
	must.NotNil(spanInPopup)

	text3, err := spanInPopup.TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text3)
}
