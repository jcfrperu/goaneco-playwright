//go:build e2e

// Ref: TestSelectorsText.java

package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSelectorsTextSmoke verifies that the >> chaining operator works with text selectors
// including HTML entities and regex flags.
// Ref: TestSelectorsText.java#shouldWorkSmoke
func TestSelectorsTextSmoke(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Hi&gt;&gt;<span></span></div>`)
	must.NoError(err)

	outer, err := page.EvalOnSelector(ctx, `text="Hi>>">>span`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span></span>", outer)

	outer, err = page.EvalOnSelector(ctx, `text=/Hi\>\>/ >> span`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span></span>", outer)

	err = page.SetContent(ctx, `<div>let's<span>hello</span></div>`)
	must.NoError(err)

	outer, err = page.EvalOnSelector(ctx, `text=/let's/i >> span`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span>hello</span>", outer)

	outer, err = page.EvalOnSelector(ctx, `text=/let\'s/i >> span`, "e => e.outerHTML")
	must.NoError(err)
	is.Equal("<span>hello</span>", outer)
}

// TestSelectorsTextFullNodeMatch verifies that getByText with exact=true matches the full
// concatenated text of a node, distinguishing nodes that share a prefix.
// Ref: TestSelectorsText.java#hasTextAndInternalTextShouldMatchFullNodeTextInStrictMode
func TestSelectorsTextFullNodeMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=div1>hello<span>world</span></div><div id=div2>hello</div>`)
	must.NoError(err)

	exactTrue := true

	id1, err := page.GetByText("helloworld", &playwright.GetByTextOptions{Exact: &exactTrue}).GetAttribute(ctx, "id")
	must.NoError(err)
	must.NotNil(id1)
	is.Equal("div1", *id1)

	id2, err := page.GetByText("hello", &playwright.GetByTextOptions{Exact: &exactTrue}).GetAttribute(ctx, "id")
	must.NoError(err)
	must.NotNil(id2)
	is.Equal("div2", *id2)
}
