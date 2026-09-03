//go:build e2e

// Ref: TestSelectorsMisc.java

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCSSLayoutPseudoSelectors verifies :right-of, :left-of, :above, :below, and :near CSS pseudo-selectors.
// Ref: TestSelectorsMisc.java#shouldWorkWithLayoutSelectors
func TestCSSLayoutPseudoSelectors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<container style='width: 500px; height: 500px; position: relative;'></container>`)
	must.NoError(err)

	setupScript := `(container, boxes) => {
    for (let i = 0; i < boxes.length; i++) {
      const div = document.createElement('div');
      div.style.position = 'absolute';
      div.style.overflow = 'hidden';
      div.style.boxSizing = 'border-box';
      div.style.border = '1px solid black';
      div.id = 'id' + i;
      div.textContent = 'id' + i;
      const box = boxes[i];
      div.style.left = box[0] + 'px';
      // Note: top is a flipped y coordinate
      div.style.top = (250 - box[1] - box[3]) + 'px';
      div.style.width = box[2] + 'px';
      div.style.height = box[3] + 'px';
      container.appendChild(div);
      const span = document.createElement('span');
      span.textContent = '' + i;
      div.appendChild(span);
    }
  }`

	boxes := [][]any{
		{0, 0, 150, 150},
		{100, 200, 50, 50},
		{200, 200, 50, 50},
		{100, 150, 50, 50},
		{201, 150, 50, 50},
		{200, 100, 50, 50},
		{50, 50, 50, 50},
		{150, 50, 50, 50},
		{150, -51, 50, 50},
		{201, -101, 50, 50},
	}

	_, err = page.EvalOnSelector(ctx, "container", setupScript, boxes)
	must.NoError(err)

	// :right-of assertions
	val, err := page.EvalOnSelector(ctx, "div:right-of(#id6)", "e => e.id")
	must.NoError(err)
	is.Equal("id7", val)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id1)", "e => e.id")
	must.NoError(err)
	is.Equal("id2", val)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id3)", "e => e.id")
	must.NoError(err)
	is.Equal("id4", val)

	el, err := page.QuerySelector(ctx, "div:right-of(#id4)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id0)", "e => e.id")
	must.NoError(err)
	is.Equal("id7", val)

	val, err = page.EvalOnSelector(ctx, "div:right-of(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id9", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:right-of(#id3)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id4,id2,id5,id7,id8,id9", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:right-of(#id3, 50)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id2,id5,id7,id8", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:right-of(#id3, 49)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id7,id8", val)

	// :left-of assertions
	val, err = page.EvalOnSelector(ctx, "div:left-of(#id2)", "e => e.id")
	must.NoError(err)
	is.Equal("id1", val)

	el, err = page.QuerySelector(ctx, "div:left-of(#id0)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelector(ctx, "div:left-of(#id5)", "e => e.id")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelector(ctx, "div:left-of(#id9)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	val, err = page.EvalOnSelector(ctx, "div:left-of(#id4)", "e => e.id")
	must.NoError(err)
	is.Equal("id3", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:left-of(#id5)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id7,id3,id1,id6,id8", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:left-of(#id5, 3)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id7,id8", val)

	// :above assertions
	val, err = page.EvalOnSelector(ctx, "div:above(#id0)", "e => e.id")
	must.NoError(err)
	is.Equal("id3", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id5)", "e => e.id")
	must.NoError(err)
	is.Equal("id4", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id7)", "e => e.id")
	must.NoError(err)
	is.Equal("id5", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelector(ctx, "div:above(#id9)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	el, err = page.QuerySelector(ctx, "div:above(#id2)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelectorAll(ctx, "div:above(#id5)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id4,id2,id3,id1", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:above(#id5, 20)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id4,id3", val)

	// :below assertions
	val, err = page.EvalOnSelector(ctx, "div:below(#id4)", "e => e.id")
	must.NoError(err)
	is.Equal("id5", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id3)", "e => e.id")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id2)", "e => e.id")
	must.NoError(err)
	is.Equal("id4", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id6)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id7)", "e => e.id")
	must.NoError(err)
	is.Equal("id8", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id9", val)

	el, err = page.QuerySelector(ctx, "div:below(#id9)")
	must.NoError(err)
	is.Nil(el)

	val, err = page.EvalOnSelectorAll(ctx, "div:below(#id3)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id5,id6,id7,id8,id9", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:below(#id3, 105)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id5,id6,id7", val)

	// :near assertions
	val, err = page.EvalOnSelector(ctx, "div:near(#id0)", "e => e.id")
	must.NoError(err)
	is.Equal("id3", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id7)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id5,id3,id6", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id0)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id3,id6,id7,id8,id1,id5", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id6)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0,id3,id7", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id6, 10)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id0", val)

	val, err = page.EvalOnSelectorAll(ctx, "div:near(#id0, 100)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id3,id6,id7,id8,id1,id5,id4,id2", val)

	// Combined layout
	val, err = page.EvalOnSelectorAll(ctx, "div:below(#id5):above(#id8)", "els => els.map(e => e.id).join(',')")
	must.NoError(err)
	is.Equal("id7,id6", val)

	val, err = page.EvalOnSelector(ctx, "div:below(#id5):above(#id8)", "e => e.id")
	must.NoError(err)
	is.Equal("id7", val)

	// Error cases
	_, nearErr := page.QuerySelector(ctx, ":near(50)")
	is.Error(nearErr)
	is.ErrorContains(nearErr, `"near" engine expects a selector list`)

	_, leftOfErr := page.QuerySelector(ctx, "div >> left-of=abc")
	is.Error(leftOfErr)
	is.ErrorContains(leftOfErr, "Malformed selector: left-of=abc")
}

// TestSelectorsMiscInternalHasNot verifies internal:has-not selector filters sections without matching descendants.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalHasNot
func TestSelectorsMiscInternalHasNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<section><span></span><div></div></section><section><br></section>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `section >> internal:has-not="span"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="span, div, br"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="br"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="span, div"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), val)

	val, err = page.EvalOnSelectorAll(ctx, `section >> internal:has-not="article"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), val)
}

// TestSelectorsMiscInternalAnd verifies internal:and selector intersects element sets.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalAnd
func TestSelectorsMiscInternalAnd(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class=foo>hello</div><div class=bar>world</div><span class=foo>hello2</span><span class=bar>world2</span>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `div >> internal:and="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	arr, ok := val.([]any)
	is.True(ok)
	is.Empty(arr)

	val, err = page.EvalOnSelectorAll(ctx, `div >> internal:and=".foo"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `div >> internal:and=".bar"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"world"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `span >> internal:and="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello2", "world2"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `.foo >> internal:and="div"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `.bar >> internal:and="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"world2"}, val)
}

// TestSelectorsMiscInternalOr verifies internal:or selector unions element sets.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalOr
func TestSelectorsMiscInternalOr(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>hello</div><span>world</span>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `div >> internal:or="span"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello", "world"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `span >> internal:or="div"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"hello", "world"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `article >> internal:or="something"`, "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), val)

	text, err := page.Locator(`article >> internal:or="div"`).TextContent(ctx)
	must.NoError(err)
	is.Equal("hello", text)

	text, err = page.Locator(`article >> internal:or="span"`).TextContent(ctx)
	must.NoError(err)
	is.Equal("world", text)
}

// TestSelectorsMiscInternalChain verifies internal:chain selector scopes results within a parent match.
// Ref: TestSelectorsMisc.java#shouldWorkWithInternalChain
func TestSelectorsMiscInternalChain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>one <span>two</span> <button>three</button> </div><span>four</span><button>five</button>`)
	must.NoError(err)

	val, err := page.EvalOnSelectorAll(ctx, `div >> internal:chain="button"`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"three"}, val)

	val, err = page.EvalOnSelectorAll(ctx, `div >> internal:chain="span >> internal:or=\"button\""`, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"two", "three"}, val)
}
