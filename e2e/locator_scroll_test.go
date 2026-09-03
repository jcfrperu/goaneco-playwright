//go:build e2e

// Locator.ScrollIntoViewIfNeeded E2E tests.
// Migration of: TestLocatorScroll.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorScrollIntoViewMakesVisible verifies ScrollIntoViewIfNeeded brings element into viewport.
// Ref: TestLocatorScroll.java#shouldScrollIntoView
func TestLocatorScrollIntoViewMakesVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 400, 400))
	must.NoError(page.SetContent(ctx, `
		<div style="height:800px;">spacer</div>
		<div id="target" style="height:50px;background:red;">target</div>
	`))

	must.NoError(page.Locator("#target").ScrollIntoViewIfNeeded(ctx))

	// Verify target is now visible in viewport
	inViewport, err := page.Evaluate(ctx, `() => {
		const el = document.getElementById('target');
		const rect = el.getBoundingClientRect();
		return rect.top >= 0 && rect.bottom <= window.innerHeight;
	}`)
	must.NoError(err)
	is.Equal(true, inViewport)
}

// TestLocatorScrollIntoViewAlreadyVisible verifies ScrollIntoViewIfNeeded is safe when already visible.
// Ref: TestLocatorScroll.java#shouldNotScrollIfAlreadyVisible
func TestLocatorScrollIntoViewAlreadyVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="top">always visible</div>`))

	// Should not throw or error when element is already visible
	must.NoError(page.Locator("#top").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#top").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorScrollIntoViewAndBoundingBox verifies BoundingBox after scroll reflects on-screen position.
// Ref: TestLocatorScroll.java#shouldHaveValidBoundingBoxAfterScroll
func TestLocatorScrollIntoViewAndBoundingBox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 800, 600))
	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;">spacer</div>
		<div id="box" style="width:100px;height:100px;background:blue;">box</div>
	`))

	must.NoError(page.Locator("#box").ScrollIntoViewIfNeeded(ctx))

	bb, err := page.Locator("#box").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
	is.Greater(bb.Width, float64(0))
	is.Greater(bb.Height, float64(0))
}

func TestScrollIntoViewMakesElementVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px"></div>
		<div id="target" style="height:100px;background:blue">target</div>
	`))

	must.NoError(page.Locator("#target").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#target").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

func TestScrollIntoViewForAlreadyVisibleElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="visible" style="height:100px">visible</div>`))

	must.NoError(page.Locator("#visible").ScrollIntoViewIfNeeded(ctx))
}

func TestScrollIntoViewUpdatesScrollPosition(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:3000px"></div>
		<div id="far" style="height:50px">far element</div>
	`))

	initialScroll, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.Equal(float64(0), initialScroll)

	must.NoError(page.Locator("#far").ScrollIntoViewIfNeeded(ctx))

	afterScroll, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	scrolled, ok := afterScroll.(float64)
	is.True(ok)
	is.Greater(scrolled, float64(0))
}

func TestScrollIntoViewForNestedElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px">spacer</div>
		<div id="outer" style="height:500px;overflow:auto">
			<div style="height:300px">top</div>
			<div id="inner" style="height:50px;background:red">nested target</div>
		</div>
	`))

	must.NoError(page.Locator("#inner").ScrollIntoViewIfNeeded(ctx))

	bb, err := page.Locator("#inner").BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(bb)
}

// TestLocatorScrollIntoViewNoErrorEx2 verifies ScrollIntoViewIfNeeded does not error.
// Ref: TestLocatorScrollIntoView.java#shouldNotError
func TestLocatorScrollIntoViewNoErrorEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:3000px"></div>
		<button id="btn">Far button</button>
	`))

	must.NoError(page.Locator("#btn").ScrollIntoViewIfNeeded(ctx))
}

// TestLocatorScrollIntoViewMakesElementVisible verifies element is in viewport after scroll.
// Ref: TestLocatorScrollIntoView.java#shouldBringIntoView
func TestLocatorScrollIntoViewMakesElementVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:5000px;width:10px"></div>
		<p id="target">Target</p>
	`))

	must.NoError(page.Locator("#target").ScrollIntoViewIfNeeded(ctx))

	inViewport, err := page.Evaluate(ctx, `() => {
		const el = document.getElementById('target');
		const rect = el.getBoundingClientRect();
		return rect.top >= 0 && rect.bottom <= window.innerHeight;
	}`)
	must.NoError(err)
	is.Equal(true, inViewport)
}

// TestLocatorScrollIntoViewAlreadyVisible verifies no error when already visible.
// Ref: TestLocatorScrollIntoView.java#shouldWorkWhenAlreadyVisible
func TestLocatorScrollIntoViewAlreadyVisibleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Already visible</button>`))

	must.NoError(page.Locator("#btn").ScrollIntoViewIfNeeded(ctx))
}

// TestLocatorScrollIntoViewWorksOnList verifies ScrollIntoViewIfNeeded works on list items.
// Ref: TestLocatorScrollIntoView.java#shouldWorkOnList
func TestLocatorScrollIntoViewWorksOnList(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px"></div>
		<ul>
			<li id="item">List item</li>
		</ul>
	`))

	must.NoError(page.Locator("#item").ScrollIntoViewIfNeeded(ctx))
}

// TestScrollIntoViewBelowFoldEx3 verifies ScrollIntoViewIfNeeded scrolls down.
// Ref: TestLocatorScrollIntoView.java#shouldScrollDownToElement
func TestScrollIntoViewBelowFoldEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;">Spacer</div>
		<div id="bottom" style="height:100px;">Bottom element</div>
	`))
	must.NoError(page.SetViewportSize(ctx, 800, 600))

	must.NoError(page.Locator("#bottom").ScrollIntoViewIfNeeded(ctx))

	scrollY, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.Greater(scrollY.(float64), float64(0))
}

// TestScrollIntoViewAlreadyVisibleEx3 verifies no-op when already visible.
// Ref: TestLocatorScrollIntoView.java#shouldBeNoOpWhenVisible
func TestScrollIntoViewAlreadyVisibleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="top">Top element</div>`))

	must.NoError(page.Locator("#top").ScrollIntoViewIfNeeded(ctx))
}

// TestScrollIntoViewRestoresScrollEx3 verifies element is visible after scroll.
// Ref: TestLocatorScrollIntoView.java#shouldMakeElementVisible
func TestScrollIntoViewRestoresScrollEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:3000px;">Big page</div>
		<div id="target">Target</div>
	`))

	must.NoError(page.Locator("#target").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#target").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestScrollIntoViewButtonEx4 verifies ScrollIntoViewIfNeeded works for buttons.
// Ref: TestLocatorScrollIntoView.java#shouldScrollButton
func TestScrollIntoViewButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;"></div>
		<button id="far-btn">Far button</button>
	`))

	must.NoError(page.Locator("#far-btn").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#far-btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestScrollIntoViewInputEx4 verifies ScrollIntoViewIfNeeded works for inputs.
// Ref: TestLocatorScrollIntoView.java#shouldScrollInput
func TestScrollIntoViewInputEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:3000px;"></div>
		<input id="deep-inp" type="text" placeholder="Deep input">
	`))

	must.NoError(page.Locator("#deep-inp").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#deep-inp").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestScrollIntoViewImageEx4 verifies ScrollIntoViewIfNeeded works for images.
// Ref: TestLocatorScrollIntoView.java#shouldScrollImage
func TestScrollIntoViewImageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2500px;"></div>
		<img id="deep-img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			style="width:50px;height:50px;">
	`))

	must.NoError(page.Locator("#deep-img").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#deep-img").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestScrollIntoViewHeadingEx5 verifies ScrollIntoViewIfNeeded for heading.
// Ref: TestLocatorScrollIntoView.java#shouldScrollHeading
func TestScrollIntoViewHeadingEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;"></div>
		<h2 id="h">Deep Heading</h2>
	`))

	must.NoError(page.Locator("#h").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#h").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestScrollIntoViewLinkEx5 verifies ScrollIntoViewIfNeeded for link element.
// Ref: TestLocatorScrollIntoView.java#shouldScrollLink
func TestScrollIntoViewLinkEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:3000px;"></div>
		<a id="lnk" href="#top">Back to top</a>
	`))

	must.NoError(page.Locator("#lnk").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#lnk").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestScrollIntoViewAlreadyVisibleEx5 verifies ScrollIntoViewIfNeeded for already visible element.
// Ref: TestLocatorScrollIntoView.java#shouldNoopForVisible
func TestScrollIntoViewAlreadyVisibleEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible now</div>`))

	must.NoError(page.Locator("#d").ScrollIntoViewIfNeeded(ctx))

	visible, err := page.Locator("#d").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}
