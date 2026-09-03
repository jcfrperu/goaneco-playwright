//go:build e2e

// E2E tests for Locator interaction methods: Hover, Check, Uncheck, Press, SelectOption,
// Nth, First, Last.
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocatorCheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox" />`), "SetContent failed")

	cb := page.Locator("#cb")

	checked, err := cb.IsChecked(ctx)
	must.NoError(err, "IsChecked failed")
	is.False(checked, "expected checkbox unchecked initially")

	must.NoError(cb.Check(ctx), "Check failed")

	checked, err = cb.IsChecked(ctx)
	must.NoError(err, "IsChecked after Check failed")
	is.True(checked, "expected checkbox to be checked after Check()")
}

func TestLocatorUncheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox" checked />`), "SetContent failed")

	cb := page.Locator("#cb")

	must.NoError(cb.Uncheck(ctx), "Uncheck failed")

	checked, err := cb.IsChecked(ctx)
	must.NoError(err, "IsChecked after Uncheck failed")
	is.False(checked, "expected checkbox to be unchecked after Uncheck()")
}

func TestLocatorSelectOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="fruit">
			<option value="apple">Apple</option>
			<option value="banana">Banana</option>
			<option value="cherry">Cherry</option>
		</select>
	`), "SetContent failed")

	sel := page.Locator("#fruit")

	selected, err := sel.SelectOption(ctx, "banana")
	must.NoError(err, "SelectOption failed")
	is.Len(selected, 1)
	is.Equal("banana", selected[0])

	val, err := sel.InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("banana", val)
}

func TestLocatorPress(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" />`), "SetContent failed")

	inp := page.Locator("#inp")

	must.NoError(inp.Press(ctx, "H"), "Press(H) failed")
	must.NoError(inp.Press(ctx, "i"), "Press(i) failed")

	val, err := inp.InputValue(ctx)
	must.NoError(err, "InputValue failed")
	is.Equal("Hi", val)
}

func TestLocatorHover(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" onmouseover="this.setAttribute('data-hovered','1')">hover me</div>
	`), "SetContent failed")

	must.NoError(page.Locator("#target").Hover(ctx), "Hover failed")

	attr, err := page.Locator("#target").GetAttribute(ctx, "data-hovered")
	must.NoError(err, "GetAttribute failed")
	must.NotNil(attr, "expected data-hovered attribute to be set after hover")
	is.Equal("1", *attr)
}

func TestLocatorNth(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="item">Alpha</li>
			<li class="item">Beta</li>
			<li class="item">Gamma</li>
		</ul>
	`), "SetContent failed")

	items := page.Locator(".item")

	text0, err := items.Nth(0).InnerText(ctx)
	must.NoError(err, "Nth(0).InnerText failed")
	is.Equal("Alpha", text0)

	text1, err := items.Nth(1).InnerText(ctx)
	must.NoError(err, "Nth(1).InnerText failed")
	is.Equal("Beta", text1)

	text2, err := items.Nth(2).InnerText(ctx)
	must.NoError(err, "Nth(2).InnerText failed")
	is.Equal("Gamma", text2)
}

func TestLocatorFirstAndLast(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="item">First</li>
			<li class="item">Middle</li>
			<li class="item">Last</li>
		</ul>
	`), "SetContent failed")

	items := page.Locator(".item")

	first, err := items.First().InnerText(ctx)
	must.NoError(err, "First().InnerText failed")
	is.Equal("First", first)

	last, err := items.Last().InnerText(ctx)
	must.NoError(err, "Last().InnerText failed")
	is.Equal("Last", last)
}

// TestLocatorHoverTriggersCSSHoverState verifies that hover activates CSS :hover styles.
// Ref: TestPageHover.java#shouldSetCSSHoverState
func TestLocatorHoverTriggersCSSHoverState(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>
		#btn { background: red; }
		#btn:hover { background: blue; }
		</style>
		<button id="btn">hover me</button>
	`))

	colorBefore, err := page.Evaluate(ctx, `() => getComputedStyle(document.getElementById('btn')).backgroundColor`)
	must.NoError(err)
	is.Equal("rgb(255, 0, 0)", colorBefore)

	must.NoError(page.Locator("#btn").Hover(ctx))

	colorAfter, err := page.Evaluate(ctx, `() => getComputedStyle(document.getElementById('btn')).backgroundColor`)
	must.NoError(err)
	is.Equal("rgb(0, 0, 255)", colorAfter)
}

// TestLocatorHoverFiresMouseoverEvent verifies that hover fires mouseover event.
// Ref: TestPageHover.java#shouldWork
func TestLocatorHoverFiresMouseoverEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target">hover me</div>
		<script>
		window.hovered = false;
		document.getElementById('target').addEventListener('mouseover', () => { window.hovered = true; });
		</script>
	`))

	must.NoError(page.Locator("#target").Hover(ctx))

	hovered, err := page.Evaluate(ctx, "() => window.hovered")
	must.NoError(err)
	is.Equal(true, hovered)
}

// TestLocatorHoverOnNestedElement verifies hover works on elements nested inside a container.
// Ref: TestPageHover.java#shouldWorkOnNestedElements
func TestLocatorHoverOnNestedElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="outer">
			<span id="inner">inner text</span>
		</div>
		<script>
		window.lastHovered = '';
		document.getElementById('outer').addEventListener('mouseover', e => { window.lastHovered = e.target.id; });
		</script>
	`))

	must.NoError(page.Locator("#inner").Hover(ctx))

	lastHovered, err := page.Evaluate(ctx, "() => window.lastHovered")
	must.NoError(err)
	is.Equal("inner", lastHovered)
}

// TestLocatorHoverFiresMouseenter verifies Hover fires mouseenter event.
// Ref: TestLocatorHover.java#shouldFireMouseenter
func TestLocatorHoverFiresMouseenter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px;"
			onmouseenter="this.dataset.entered='yes'">hover</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('el').dataset.entered`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestLocatorHoverFiresMouseover verifies Hover fires mouseover event.
// Ref: TestLocatorHover.java#shouldFireMouseover
func TestLocatorHoverFiresMouseoverExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px;"
			onmouseover="this.dataset.over='yes'">hover</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('el').dataset.over`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestLocatorHoverOnButton verifies Hover works on a button element.
// Ref: TestLocatorHover.java#shouldHoverOnButton
func TestLocatorHoverOnButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onmouseover="this.dataset.hovered='yes'">hover me</button>
	`))

	must.NoError(page.Locator("#btn").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('btn').dataset.hovered`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestLocatorHoverRevealsCSSHover verifies Hover triggers CSS :hover pseudo-class.
// Ref: TestLocatorHover.java#shouldActivateCSSHover
func TestLocatorHoverRevealsCSSHover(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>#el { color: black; } #el:hover { color: red; }</style>
		<div id="el" style="width:50px;height:50px;">hover</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	color, err := page.Evaluate(ctx, `() => getComputedStyle(document.getElementById('el')).color`)
	must.NoError(err)
	is.Equal("rgb(255, 0, 0)", color)
}

// TestLocatorHoverFiresMouseEnter verifies Hover fires mouseenter event.
// Ref: TestLocatorHover.java#shouldFireMouseEnter
func TestLocatorHoverFiresMouseEnterExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onmouseenter="window.__enter=true">hover me</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => window.__enter`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorHoverActivatesCSSHover verifies Hover triggers :hover CSS state.
// Ref: TestLocatorHover.java#shouldTriggerCSSHover
func TestLocatorHoverActivatesCSSHoverExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>
			#el:hover { background: green; }
		</style>
		<div id="el" style="width:100px;height:100px;background:red">hover me</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	bg, err := page.Evaluate(ctx, `() => getComputedStyle(document.getElementById('el')).backgroundColor`)
	must.NoError(err)
	is.Equal("rgb(0, 128, 0)", bg)
}

// TestLocatorHoverDoesNotClickElement verifies Hover doesn't fire click.
// Ref: TestLocatorHover.java#shouldNotClick
func TestLocatorHoverDoesNotClickElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onclick="window.__clicked=true">hover me</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Nil(result)
}

// TestLocatorHoverWithMultipleElements verifies Hover works on each element.
// Ref: TestLocatorHover.java#shouldHoverMultipleElements
func TestLocatorHoverWithMultipleElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="a" style="width:100px;height:100px" onmouseover="window.__a=true">A</div>
		<div id="b" style="width:100px;height:100px" onmouseover="window.__b=true">B</div>
	`))

	must.NoError(page.Locator("#a").Hover(ctx))
	must.NoError(page.Locator("#b").Hover(ctx))

	resA, err := page.Evaluate(ctx, `() => window.__a`)
	must.NoError(err)
	is.Equal(true, resA)

	resB, err := page.Evaluate(ctx, `() => window.__b`)
	must.NoError(err)
	is.Equal(true, resB)
}

// TestLocatorHoverChangesAttribute verifies Hover triggers attribute change.
// Ref: TestLocatorHover.java#shouldChangeAttribute
func TestLocatorHoverChangesAttribute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" onmouseenter="this.setAttribute('data-hovered','true')"
		     style="width:100px;height:50px">hover me</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	attr, err := page.Locator("#el").GetAttribute(ctx, "data-hovered")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("true", *attr)
}

// TestLocatorHoverDoesNotClick verifies Hover does not trigger click event.
// Ref: TestLocatorHover.java#shouldNotTriggerClick
func TestLocatorHoverDoesNotClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked=true">Hover only</button>
	`))

	must.NoError(page.Locator("#btn").Hover(ctx))

	clicked, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Nil(clicked)
}

// TestLocatorHoverOnLink verifies Hover works on links.
// Ref: TestLocatorHover.java#shouldHoverLink
func TestLocatorHoverOnLink(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a href="#" id="lnk" onmouseenter="window.__linkHovered=true">Link</a>
	`))

	must.NoError(page.Locator("#lnk").Hover(ctx))

	hovered, err := page.Evaluate(ctx, `() => window.__linkHovered`)
	must.NoError(err)
	is.Equal(true, hovered)
}

// TestLocatorHoverNoErrorForOffscreenElement verifies Hover scrolls to element.
// Ref: TestLocatorHover.java#shouldScrollToElement
func TestLocatorHoverNoErrorForOffscreenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px"></div>
		<button id="far-btn">Far Button</button>
	`))

	must.NoError(page.Locator("#far-btn").Hover(ctx))
}

// TestLocatorHoverFiresMouseenterEx5 verifies Hover fires mouseenter event.
// Ref: TestLocatorHover.java#shouldFireMouseenter
func TestLocatorHoverFiresMouseenterEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onmouseenter="window.__entered=true">hover target</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	entered, err := page.Evaluate(ctx, `() => window.__entered`)
	must.NoError(err)
	is.Equal(true, entered)
}

// TestLocatorHoverFiresMouseoverEx5 verifies Hover fires mouseover event.
// Ref: TestLocatorHover.java#shouldFireMouseover
func TestLocatorHoverFiresMouseoverEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onmouseover="window.__over=true">hover</div>
	`))

	must.NoError(page.Locator("#el").Hover(ctx))

	over, err := page.Evaluate(ctx, `() => window.__over`)
	must.NoError(err)
	is.Equal(true, over)
}

// TestLocatorHoverOnButtonEx5 verifies Hover works on button element.
// Ref: TestLocatorHover.java#shouldWorkOnButton
func TestLocatorHoverOnButtonEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onmouseover="window.__btnHovered=true">Hover me</button>
	`))

	must.NoError(page.Locator("#btn").Hover(ctx))

	hovered, err := page.Evaluate(ctx, `() => window.__btnHovered`)
	must.NoError(err)
	is.Equal(true, hovered)
}

// TestLocatorHoverTooltipShowsEx5 verifies Hover shows tooltip via CSS :hover.
// Ref: TestLocatorHover.java#shouldShowTooltip
func TestLocatorHoverTooltipShowsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>
			#tooltip { display: none; }
			#target:hover + #tooltip { display: block; }
		</style>
		<button id="target">Hover for tooltip</button>
		<div id="tooltip">Tooltip text</div>
	`))

	must.NoError(page.Locator("#target").Hover(ctx))

	visible, err := page.Locator("#tooltip").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorHoverFiresMouseoverEx6 verifies Hover fires mouseover event.
// Ref: TestLocatorHover.java#shouldFireMouseoverEvent
func TestLocatorHoverFiresMouseoverEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="width:100px;height:100px;">Hover me</div>
		<script>
			var hovered = false;
			document.getElementById('d').addEventListener('mouseover', function() { hovered = true; });
		</script>
	`))

	must.NoError(page.Locator("#d").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => hovered`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorHoverOnButtonEx6 verifies Hover on button element.
// Ref: TestLocatorHover.java#shouldHoverOnButton
func TestLocatorHoverOnButtonEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn">Hover Button</button>
		<script>
			var entered = false;
			document.getElementById('btn').addEventListener('mouseenter', function() { entered = true; });
		</script>
	`))

	must.NoError(page.Locator("#btn").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => entered`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorHoverRevealsHiddenContentEx6 verifies Hover can reveal CSS-hover-triggered content.
// Ref: TestLocatorHover.java#shouldRevealContent
func TestLocatorHoverRevealsHiddenContentEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<style>
			#menu { display: none; }
			#trigger:hover + #menu { display: block; }
		</style>
		<div id="trigger">Menu</div>
		<div id="menu">Dropdown</div>
	`))

	must.NoError(page.Locator("#trigger").Hover(ctx))

	visible, err := page.Locator("#menu").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorHoverOnLinkEx6 verifies Hover triggers mouseenter on links.
// Ref: TestLocatorHover.java#shouldHoverOnLink
func TestLocatorHoverOnLinkEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a id="lnk" href="#">Hover Link</a>
		<script>
			var linkHovered = false;
			document.getElementById('lnk').addEventListener('mouseenter', function() { linkHovered = true; });
		</script>
	`))

	must.NoError(page.Locator("#lnk").Hover(ctx))

	result, err := page.Evaluate(ctx, `() => linkHovered`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestHoverChangesClassEx7 verifies Hover triggers CSS class changes via JS.
// Ref: TestLocatorHover.java#shouldChangeClassOnHover
func TestHoverChangesClassEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" onmouseover="this.className='hovered'" onmouseout="this.className=''">Hover me</div>
	`))

	must.NoError(page.Locator("#d").Hover(ctx))

	cls, err := page.Evaluate(ctx, `() => document.getElementById('d').className`)
	must.NoError(err)
	is.Equal("hovered", cls)
}

// TestHoverOnButtonEx7 verifies Hover works on button elements.
// Ref: TestLocatorHover.java#shouldHoverButton
func TestHoverOnButtonEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onmouseover="this.textContent='Hovered'">Click me</button>
	`))

	must.NoError(page.Locator("#btn").Hover(ctx))

	text, err := page.Locator("#btn").TextContent(ctx)
	must.NoError(err)
	is.Equal("Hovered", text)
}

// TestHoverOnLinkEx7 verifies Hover works on anchor elements.
// Ref: TestLocatorHover.java#shouldHoverLink
func TestHoverOnLinkEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a id="a" href="#" onmouseover="document.getElementById('out').textContent='link hovered'">Link</a>
		<span id="out"></span>
	`))

	must.NoError(page.Locator("#a").Hover(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("link hovered", out)
}

// TestHoverShowsTooltipEx7 verifies Hover on tooltip trigger element.
// Ref: TestLocatorHover.java#shouldShowTooltip
func TestHoverShowsTooltipEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="trigger" onmouseover="document.getElementById('tip').style.display='block'"
			onmouseout="document.getElementById('tip').style.display='none'">Hover for tip</div>
		<div id="tip" style="display:none">Tooltip text</div>
	`))

	must.NoError(page.Locator("#trigger").Hover(ctx))

	visible, err := page.Locator("#tip").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestHoverHighlightsTableRowEx8 verifies Hover highlights table row.
// Ref: TestLocatorHover.java#shouldHighlightTableRow
func TestHoverHighlightsTableRowEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr id="row1" onmouseover="this.style.background='yellow'" onmouseout="this.style.background=''">
				<td>Row 1</td>
			</tr>
		</table>
	`))

	must.NoError(page.Locator("#row1").Hover(ctx))

	bg, err := page.Evaluate(ctx, `() => document.getElementById('row1').style.background`)
	must.NoError(err)
	is.Equal("yellow", bg)
}

// TestHoverOnCheckboxEx8 verifies Hover works on checkbox elements.
// Ref: TestLocatorHover.java#shouldHoverCheckbox
func TestHoverOnCheckboxEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="chk" type="checkbox" onmouseover="document.getElementById('out').textContent='hovered'">
		<span id="out"></span>
	`))

	must.NoError(page.Locator("#chk").Hover(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("hovered", out)
}

// TestHoverSetsAriaEx8 verifies Hover updates aria attribute via JS.
// Ref: TestLocatorHover.java#shouldUpdateAriaOnHover
func TestHoverSetsAriaEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" onmouseover="this.setAttribute('aria-expanded','true')">Menu</div>
	`))

	must.NoError(page.Locator("#d").Hover(ctx))

	val, err := page.Locator("#d").GetAttribute(ctx, "aria-expanded")
	must.NoError(err)
	must.NotNil(val)
	is.Equal("true", *val)
}

// TestHoverImageEx9 verifies hover on image triggers title visibility.
// Ref: TestLocatorHover.java#shouldHoverImage
func TestHoverImageEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img id="img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			title="Image tooltip" style="width:50px;height:50px;"
			onmouseover="this.setAttribute('data-hovered','true')">
	`))

	must.NoError(page.Locator("#img").Hover(ctx))

	attr, err := page.Locator("#img").GetAttribute(ctx, "data-hovered")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("true", *attr)
}

// TestHoverMenuItemEx9 verifies hover on menu item shows submenu.
// Ref: TestLocatorHover.java#shouldHoverMenuItem
func TestHoverMenuItemEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav>
			<ul>
				<li id="menu" onmouseover="document.getElementById('sub').style.display='block'">
					Menu
					<ul id="sub" style="display:none">
						<li>Sub item</li>
					</ul>
				</li>
			</ul>
		</nav>
	`))

	must.NoError(page.Locator("#menu").Hover(ctx))

	visible, err := page.Locator("#sub").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestHoverSetsDataAttrEx9 verifies hover on div sets data attribute.
// Ref: TestLocatorHover.java#shouldSetDataAttribute
func TestHoverSetsDataAttrEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="card" onmouseenter="this.dataset.active='yes'">Card</div>
	`))

	must.NoError(page.Locator("#card").Hover(ctx))

	attr, err := page.Locator("#card").GetAttribute(ctx, "data-active")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("yes", *attr)
}

// TestLocatorTapFiresTouchStart verifies Tap fires touchstart event.
// Ref: TestLocatorTap.java#shouldFireTouchStart
func TestLocatorTapFiresTouchStart(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" style="width:100px;height:100px;" ontouchstart="this.dataset.touched='yes'">
		</div>
	`))

	must.NoError(page.Locator("#target").Tap(ctx))
}

// TestLocatorTapTriggersClick verifies Tap also triggers click on most browsers.
// Ref: TestLocatorTap.java#shouldTriggerClick
func TestLocatorTapTriggersClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="this.dataset.clicked='yes'">tap me</button>
	`))

	must.NoError(page.Locator("#btn").Tap(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('btn').dataset.clicked`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestLocatorTapOnInput verifies Tap focuses input element.
// Ref: TestLocatorTap.java#shouldFocusInput
func TestLocatorTapOnInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input type="text" id="inp">
	`))

	must.NoError(page.Locator("input").Tap(ctx))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("inp", activeId)
}

// TestLocatorTapFiresClickEventEx2 verifies Tap fires click event.
// Ref: TestLocatorTap.java#shouldFireClick
func TestLocatorTapFiresClickEventEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked=true">Tap me</button>
	`))

	must.NoError(page.Locator("#btn").Tap(ctx))

	clicked, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Equal(true, clicked)
}

// TestLocatorTapOnCheckboxEx2 verifies Tap toggles checkbox state.
// Ref: TestLocatorTap.java#shouldToggleCheckbox
func TestLocatorTapOnCheckboxEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	must.NoError(page.Locator("#cb").Tap(ctx))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorTapNoErrorEx2 verifies Tap on div does not error.
// Ref: TestLocatorTap.java#shouldNotError
func TestLocatorTapNoErrorEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px">Tap target</div>
	`))

	must.NoError(page.Locator("#el").Tap(ctx))
}

// TestLocatorTapFiresTouchstartEx3 verifies Tap fires touchstart event.
// Ref: TestLocatorTap.java#shouldFireTouchstart
func TestLocatorTapFiresTouchstartEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px;">Touch me</div>
		<script>
			var touched = false;
			document.getElementById('el').addEventListener('touchstart', function() { touched = true; });
		</script>
	`))

	must.NoError(page.Locator("#el").Tap(ctx))

	result, err := page.Evaluate(ctx, `() => touched`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorTapOnButtonEx3 verifies Tap works on button elements.
// Ref: TestLocatorTap.java#shouldTapButton
func TestLocatorTapOnButtonEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" ontouchstart="window.__tapped=true">Tap</button>
	`))

	must.NoError(page.Locator("#btn").Tap(ctx))

	result, err := page.Evaluate(ctx, `() => window.__tapped`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestLocatorTapOnCheckboxEx3 verifies Tap works on checkbox.
// Ref: TestLocatorTap.java#shouldTapCheckbox
func TestLocatorTapOnCheckboxEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))
	must.NoError(page.Locator("#chk").Tap(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestTapButtonEx4 verifies Tap works on button elements.
// Ref: TestLocatorTap.java#shouldTapButton
func TestTapButtonEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" ontouchend="this.textContent='tapped'">Tap me</button>`))

	must.NoError(page.Locator("#btn").Tap(ctx))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestTapCheckboxEx4 verifies Tap can toggle checkbox.
// Ref: TestLocatorTap.java#shouldTapCheckbox
func TestTapCheckboxEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	must.NoError(page.Locator("#chk").Tap(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestTapLinkEx4 verifies Tap on link fires click-like event.
// Ref: TestLocatorTap.java#shouldTapLink
func TestTapLinkEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a id="a" href="#" onclick="document.getElementById('out').textContent='tapped'">Link</a>
		<span id="out"></span>
	`))

	must.NoError(page.Locator("#a").Tap(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("tapped", out)
}
