//go:build e2e

// Ref: TestQuerySelector.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQuerySelectorWithTextSelector verifies QuerySelector works with text= selector syntax.
// Ref: TestQuerySelector.java#shouldQueryExistingElementWithTextSelector
func TestQuerySelectorWithTextSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "text='test'")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorWithXpathSelector verifies QuerySelector works with xpath= selector syntax.
// Ref: TestQuerySelector.java#shouldQueryExistingElementWithXpathSelector
func TestQuerySelectorWithXpathSelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "xpath=/html/body/section")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorAutoDetectsXpath verifies QuerySelector auto-detects XPath selectors starting with //.
// Ref: TestQuerySelector.java#shouldAutoDetectXpathSelector
func TestQuerySelectorAutoDetectsXpath(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "//html/body/section")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorAutoDetectsXpathWithParenthesis verifies QuerySelector auto-detects XPath selectors starting with (.
// Ref: TestQuerySelector.java#shouldAutoDetectXpathSelectorWithStartingParenthesis
func TestQuerySelectorAutoDetectsXpathWithParenthesis(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "(//section)[1]")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorAutoDetectsText verifies QuerySelector auto-detects text selectors enclosed in quotes.
// Ref: TestQuerySelector.java#shouldAutoDetectTextSelector
func TestQuerySelectorAutoDetectsText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	el, err := page.QuerySelector(ctx, "'test'")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorSupportsSyntax verifies QuerySelector supports the >> combinator syntax.
// Ref: TestQuerySelector.java#shouldSupportSyntax
func TestQuerySelectorSupportsSyntax(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section><div>test</div></section>`))

	el, err := page.QuerySelector(ctx, "css=section >> css=div")
	must.NoError(err)
	must.NotNil(el)
}

// TestQuerySelectorAllXpathFindsElement verifies QuerySelectorAll with xpath= returns a single match.
// Ref: TestQuerySelector.java#xpathShouldQueryExistingElement
func TestQuerySelectorAllXpathFindsElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<section>test</section>`))

	elements, err := page.QuerySelectorAll(ctx, "xpath=/html/body/section")
	must.NoError(err)
	is.Len(elements, 1)
	must.NotNil(elements[0])
}

// TestQuerySelectorAllXpathReturnsEmpty verifies QuerySelectorAll returns empty for a non-existing XPath.
// Ref: TestQuerySelector.java#xpathShouldReturnEmptyArrayForNonExistingElement
func TestQuerySelectorAllXpathReturnsEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	elements, err := page.QuerySelectorAll(ctx, "//html/body/non-existing-element")
	must.NoError(err)
	is.Empty(elements)
}

// TestQuerySelectorAllXpathReturnsMultiple verifies QuerySelectorAll with xpath= returns multiple matches.
// Ref: TestQuerySelector.java#xpathShouldReturnMultipleElements
func TestQuerySelectorAllXpathReturnsMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div><div></div>`))

	elements, err := page.QuerySelectorAll(ctx, "xpath=/html/body/div")
	must.NoError(err)
	is.Len(elements, 2)
}
