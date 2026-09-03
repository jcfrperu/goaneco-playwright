//go:build e2e

package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFrameContentReturnsHTML verifies Frame.Content returns the HTML content.
// Ref: TestFrame.java#shouldReturnContent
func TestFrameContentReturnsHTML(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="box">hello</div>`))

	content, err := page.MainFrame().Content(ctx)
	must.NoError(err)
	is.Contains(content, "hello")
	is.Contains(content, "<html")
}

// TestFrameSetContentAndVerifyLocator verifies Frame.SetContent updates the DOM.
// Ref: TestFrame.java#shouldSetContent
func TestFrameSetContentAndVerifyLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.MainFrame().SetContent(ctx, `<p id="msg">frame content</p>`))

	text, err := page.Locator("#msg").InnerText(ctx)
	must.NoError(err)
	is.Equal("frame content", text)
}

// TestFrameEvaluateReturnsValue verifies Frame.Evaluate returns computed values.
// Ref: TestFrame.java#shouldEvaluateExpression
func TestFrameEvaluateReturnsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>42</div>`))

	result, err := page.MainFrame().Evaluate(ctx, `() => 1 + 1`)
	must.NoError(err)
	is.Equal(float64(2), result)
}

// TestFrameURLMatchesPage verifies Frame.URL matches the page URL after navigation.
// Ref: TestFrame.java#shouldReportFrameURL
func TestFrameURLMatchesPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.Equal(srv.EmptyPage(), page.MainFrame().URL())
}

// TestFrameTitleMatchesPageTitle verifies Frame.Title returns the page title.
// Ref: TestFrame.java#shouldReportFrameTitle
func TestFrameTitleMatchesPageTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/title.html"))

	title, err := page.MainFrame().Title(ctx)
	must.NoError(err)
	is.Equal("Woof-Woof", title)
}

// TestFrameQuerySelectorReturnsElement verifies Frame.QuerySelector finds elements.
// Ref: TestFrame.java#shouldQuerySelector
func TestFrameQuerySelectorReturnsElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">found</div>`))

	el, err := page.MainFrame().QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("found", text)
}

// TestFrameQuerySelectorReturnsNilForMissing verifies Frame.QuerySelector returns nil when no match.
// Ref: TestFrame.java#shouldReturnNullForMissingSelector
func TestFrameQuerySelectorReturnsNilForMissing(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>hello</div>`))

	el, err := page.MainFrame().QuerySelector(ctx, "#doesNotExist")
	must.NoError(err)
	is.Nil(el)
}

// TestFrameIsDetachedFalseForMain verifies Frame.IsDetached returns false for main frame.
// Ref: TestFrame.java#mainFrameShouldNotBeDetached
func TestFrameIsDetachedFalseForMain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.False(page.MainFrame().IsDetached())
}

// TestFrameNameForMainIsEmpty verifies Frame.Name returns empty string for main frame.
// Ref: TestFrame.java#mainFrameNameIsEmpty
func TestFrameNameForMainIsEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.Equal("", page.MainFrame().Name())
}

// TestFrameEvaluateReturnsValueExtra2 verifies Frame.Evaluate returns a value.
// Ref: TestFrame.java#shouldEvaluate
func TestFrameEvaluateReturnsValueExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="val">42</div>`))

	frame := page.MainFrame()
	result, err := frame.Evaluate(ctx, `() => parseInt(document.getElementById('val').textContent)`)
	must.NoError(err)
	is.Equal(float64(42), result)
}

// TestFrameSetContentAndGetContent verifies Frame.SetContent and Content round-trip.
// Ref: TestFrame.java#shouldSetAndGetContent
func TestFrameSetContentAndGetContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	frame := page.MainFrame()
	must.NoError(frame.SetContent(ctx, `<p>hello</p>`))

	content, err := frame.Content(ctx)
	must.NoError(err)
	is.Contains(content, "hello")
}

// TestFrameQuerySelectorReturnsElementExtra2 verifies Frame.QuerySelector finds element.
// Ref: TestFrame.java#shouldQuerySelector
func TestFrameQuerySelectorReturnsElementExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">click</button>`))

	frame := page.MainFrame()
	handle, err := frame.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(handle)
}

// TestFrameQuerySelectorReturnsNilForMissingExtra2 verifies QuerySelector returns nil when no match.
// Ref: TestFrame.java#shouldReturnNilForMissing
func TestFrameQuerySelectorReturnsNilForMissingExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>nothing</div>`))

	frame := page.MainFrame()
	handle, err := frame.QuerySelector(ctx, "#nothere")
	must.NoError(err)
	is.Nil(handle)
}

// TestFrameGotoSetsURL verifies Frame.Goto updates the frame URL.
// Ref: TestFrame.java#shouldGoto
func TestFrameGotoSetsURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	frame := page.MainFrame()
	is.Equal(srv.EmptyPage(), frame.URL())
}

// TestFrameTitleMatchesPageTitleExtra2 verifies Frame.Title matches document.title.
// Ref: TestFrame.java#shouldMatchPageTitle
func TestFrameTitleMatchesPageTitleExtra2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<title>MyTitle</title>`))

	title, err := page.MainFrame().Title(ctx)
	must.NoError(err)
	is.Equal("MyTitle", title)
}

// TestFrameURLIsNotEmpty verifies Frame.URL returns a non-empty string.
// Ref: TestFrame.java#shouldReturnURL
func TestFrameURLIsNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.NotEmpty(page.MainFrame().URL())
}

// TestFrameNameForMainIsEmptyEx3 verifies MainFrame name is empty string.
// Ref: TestFrame.java#shouldHaveEmptyNameForMain
func TestFrameNameForMainIsEmptyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.Equal("", page.MainFrame().Name())
}

// TestFrameEvaluateReturnsString verifies Frame.Evaluate can return a string.
// Ref: TestFrame.java#shouldReturnString
func TestFrameEvaluateReturnsString(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	val, err := page.MainFrame().Evaluate(ctx, `() => "frame result"`)
	must.NoError(err)
	is.Equal("frame result", val)
}

// TestFrameLocatorCanFindElement verifies Frame.Locator scopes queries.
// Ref: TestFrame.java#shouldScopeLocator
func TestFrameLocatorCanFindElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="container">
			<p class="text">Frame content</p>
		</div>
	`))

	text, err := page.MainFrame().Locator(".text").InnerText(ctx)
	must.NoError(err)
	is.Equal("Frame content", text)
}

// TestFrameIsDetachedFalseForNewFrame verifies IsDetached is false for new frame.
// Ref: TestFrame.java#shouldReturnFalseForAttachedFrame
func TestFrameIsDetachedFalseForNewFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.False(page.MainFrame().IsDetached())
}

// TestFrameSetContentAndEvaluate verifies Frame.SetContent + Evaluate round-trip.
// Ref: TestFrame.java#shouldSetAndEvaluate
func TestFrameSetContentAndEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.MainFrame().SetContent(ctx, `<div id="msg">hello frame</div>`))

	val, err := page.MainFrame().Evaluate(ctx, `() => document.getElementById('msg').textContent`)
	must.NoError(err)
	is.Equal("hello frame", val)
}

// TestFrameURLMatchesPageURLEx4 verifies Frame URL matches the loaded page URL.
// Ref: TestFrame.java#shouldMatchURL
func TestFrameURLMatchesPageURLEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	frame := page.MainFrame()
	must.NotNil(frame)
	is.Contains(frame.URL(), "empty.html")
}

// TestFrameEvaluateStringEx4 verifies Frame.Evaluate returns a string result.
// Ref: TestFrame.java#shouldEvaluateString
func TestFrameEvaluateStringEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="msg">Hello Frame</div>`))

	frame := page.MainFrame()
	result, err := frame.Evaluate(ctx, `() => document.getElementById('msg').textContent`)
	must.NoError(err)
	is.Equal("Hello Frame", result)
}

// TestFrameLocatorCountEx4 verifies Frame.Locator Count works.
// Ref: TestFrame.java#shouldCountElements
func TestFrameLocatorCountEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>A</li>
			<li>B</li>
			<li>C</li>
		</ul>
	`))

	frame := page.MainFrame()
	count, err := frame.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestFrameSetContentAndReadEx4 verifies Frame.SetContent then reading content.
// Ref: TestFrame.java#shouldSetContentAndRead
func TestFrameSetContentAndReadEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	frame := page.MainFrame()
	must.NoError(frame.SetContent(ctx, `<p id="p">Frame content</p>`))

	text, err := frame.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("Frame content", text)
}

// TestFrameOnlyMainFrameForSimplePageEx4 verifies Frames() has only main frame for non-iframe page.
// Ref: TestFrame.java#shouldHaveOnlyMainFrame
func TestFrameOnlyMainFrameForSimplePageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><p>No iframes</p></body></html>`))

	frames := page.Frames()
	is.Len(frames, 1)
}

// TestFrameMainFrameEx5 verifies MainFrame is available on Page.
// Ref: TestFrame.java#shouldReturnMainFrame
func TestFrameMainFrameEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Content</p>`))

	frame := page.MainFrame()
	must.NotNil(frame)
}

// TestFrameEvaluateEx5 verifies Frame can run Evaluate.
// Ref: TestFrame.java#shouldRunEvaluate
func TestFrameEvaluateEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>Content</p>`))

	result, err := page.MainFrame().Evaluate(ctx, `() => 42 * 2`)
	must.NoError(err)
	is.Equal(float64(84), result)
}

// TestFrameTitleEx5 verifies Frame Title returns page title.
// Ref: TestFrame.java#shouldReturnFrameTitle
func TestFrameTitleEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Frame Test</title></head><body></body></html>`))

	title, err := page.MainFrame().Title(ctx)
	must.NoError(err)
	is.Equal("Frame Test", title)
}

// TestFrameURLEx5 verifies Frame URL returns current URL.
// Ref: TestFrame.java#shouldReturnFrameURL
func TestFrameURLEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>test</p>`))

	url := page.MainFrame().URL()
	is.NotEmpty(url)
}

// TestFrameQuerySelectorEx5 verifies Frame QuerySelector works.
// Ref: TestFrame.java#shouldQuerySelector
func TestFrameQuerySelectorEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Found</div>`))

	el, err := page.MainFrame().QuerySelector(ctx, "#el")
	must.NoError(err)
	must.NotNil(el)

	text, err := el.InnerText(ctx)
	must.NoError(err)
	is.Equal("Found", text)
}

// TestFrameEvaluateOnMainEx6 verifies Frame Evaluate can access DOM.
// Ref: TestFrame.java#shouldEvaluateDOMAccess
func TestFrameEvaluateOnMainEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Frame content</div>`))

	frame := page.MainFrame()
	result, err := frame.Evaluate(ctx, `() => document.getElementById('el').textContent`)
	must.NoError(err)
	is.Equal("Frame content", result)
}

// TestFrameLocatorEx6 verifies Frame Locator works like page locator.
// Ref: TestFrame.java#shouldUseLocator
func TestFrameLocatorEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn">Frame Button</button>
	`))

	frame := page.MainFrame()
	text, err := frame.Locator("#btn").InnerText(ctx)
	must.NoError(err)
	is.Equal("Frame Button", text)
}

// TestFrameInnerHTMLEx6 verifies Frame Locator InnerHTML access.
// Ref: TestFrame.java#shouldGetInnerHTML
func TestFrameInnerHTMLEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"><span>Inner</span></div>`))

	frame := page.MainFrame()
	html, err := frame.Locator("#d").InnerHTML(ctx)
	must.NoError(err)
	is.Equal("<span>Inner</span>", html)
}

// TestFrameInnerTextEx6 verifies Frame Locator InnerText access.
// Ref: TestFrame.java#shouldGetInnerText
func TestFrameInnerTextEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Frame text</p>`))

	frame := page.MainFrame()
	text, err := frame.Locator("#p").InnerText(ctx)
	must.NoError(err)
	is.Equal("Frame text", text)
}

// TestFrameInputValueEx6 verifies Frame Locator InputValue access.
// Ref: TestFrame.java#shouldGetInputValue
func TestFrameInputValueEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="frame-val">`))

	frame := page.MainFrame()
	val, err := frame.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("frame-val", val)
}

// TestFrameIsDetachedEx7 verifies MainFrame is not detached.
// Ref: TestFrame.java#mainFrameIsNotDetached
func TestFrameIsDetachedEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Content</div>`))

	frame := page.MainFrame()
	is.False(frame.IsDetached())
}

// TestFrameNameEx7 verifies MainFrame name is empty string.
// Ref: TestFrame.java#mainFrameNameIsEmpty
func TestFrameNameEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Content</div>`))

	frame := page.MainFrame()
	is.Equal("", frame.Name())
}

// TestFrameQuerySelectorEx7 verifies Frame QuerySelector works.
// Ref: TestFrame.java#shouldQuerySelector
func TestFrameQuerySelectorEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Press</button>`))

	frame := page.MainFrame()
	el, err := frame.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)
}

// TestFrameLocatorCountEx7 verifies Frame Locator can count elements.
// Ref: TestFrame.java#shouldCountViaLocator
func TestFrameLocatorCountEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul><li>A</li><li>B</li><li>C</li></ul>`))

	frame := page.MainFrame()
	count, err := frame.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestFrameEvaluateListLengthEx7 verifies Frame Evaluate returns list length.
// Ref: TestFrame.java#shouldEvaluateListLength
func TestFrameEvaluateListLengthEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul><li>A</li><li>B</li></ul>`))

	frame := page.MainFrame()
	count, err := frame.Evaluate(ctx, `() => document.querySelectorAll('li').length`)
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestFrameURLEx8 verifies MainFrame URL matches page URL.
// Ref: TestFrame.java#shouldReturnURL
func TestFrameURLEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Content</div>`))

	frame := page.MainFrame()
	url := frame.URL()
	is.NotEmpty(url)
}

// TestFrameSetContentEx8 verifies Frame SetContent works.
// Ref: TestFrame.java#shouldSetContent
func TestFrameSetContentEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	frame := page.MainFrame()
	must.NoError(frame.SetContent(ctx, `<p id="p">Set via frame</p>`))

	text, err := frame.Locator("#p").TextContent(ctx)
	must.NoError(err)
	is.Equal("Set via frame", text)
}

// TestFrameContentEx8 verifies Frame Content returns HTML.
// Ref: TestFrame.java#shouldReturnContent
func TestFrameContentEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="unique-marker">Marker</div>`))

	frame := page.MainFrame()
	content, err := frame.Content(ctx)
	must.NoError(err)
	is.Contains(content, "unique-marker")
}

// TestFrameTitleEx8 verifies Frame Title returns page title.
// Ref: TestFrame.java#shouldReturnTitle
func TestFrameTitleEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Frame Title Test</title></head><body></body></html>`))

	frame := page.MainFrame()
	title, err := frame.Title(ctx)
	must.NoError(err)
	is.Equal("Frame Title Test", title)
}

// TestFrameLocatorCountEx8 verifies Frame Locator Count finds elements.
// Ref: TestFrame.java#shouldCountViaLocator
func TestFrameLocatorCountEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">A</div>
		<div class="item">B</div>
		<div class="item">C</div>
		<div class="item">D</div>
	`))

	frame := page.MainFrame()
	count, err := frame.Locator(".item").Count(ctx)
	must.NoError(err)
	is.Equal(4, count)
}

// TestFrameEvaluateArithmetic9 verifies frame Evaluate with arithmetic.
// Ref: TestFrame.java#shouldEvaluateArithmetic
func TestFrameEvaluateArithmetic9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	result, err := frame.Evaluate(ctx, `() => 7 * 6`)
	must.NoError(err)
	is.Equal(float64(42), result)
}

// TestFrameLocatorFill9 verifies frame Locator fill works.
// Ref: TestFrame.java#shouldFillInput
func TestFrameLocatorFill9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	must.NoError(frame.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(frame.Locator("#inp").Fill(ctx, "frame input"))

	val, err := frame.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("frame input", val)
}

// TestFrameLocatorIsVisible9 verifies frame Locator IsVisible.
// Ref: TestFrame.java#shouldCheckVisibility
func TestFrameLocatorIsVisible9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	must.NoError(frame.SetContent(ctx, `<p id="p" style="display:none">Hidden</p>`))

	visible, err := frame.Locator("#p").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestFrameEvaluateWindowVar9 verifies frame Evaluate can access window vars.
// Ref: TestFrame.java#shouldAccessWindowVariable
func TestFrameEvaluateWindowVar9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	_, err := frame.Evaluate(ctx, `() => { window.myVar = 'frame_val'; }`)
	must.NoError(err)

	result, err := frame.Evaluate(ctx, `() => window.myVar`)
	must.NoError(err)
	is.Equal("frame_val", result)
}

// TestFrameLocatorTextContent9 verifies frame Locator TextContent.
// Ref: TestFrame.java#shouldGetTextContent
func TestFrameLocatorTextContent9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	must.NoError(frame.SetContent(ctx, `<h1 id="h">Frame Heading</h1>`))

	text, err := frame.Locator("#h").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Frame Heading", *text)
}

// TestFrameMainIsNotDetachedEx verifies main frame is not detached.
// Ref: TestFrame.java#shouldNotBeDetachedForMainFrame
func TestFrameMainIsNotDetachedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body>Main</body></html>`))

	frames := page.Frames()
	is.NotEmpty(frames)
	is.False(frames[0].IsDetached())
}

// TestFrameHasOneFrameByDefaultEx verifies single frame on plain page.
// Ref: TestFrame.java#shouldHaveOneFrameForPlainPage
func TestFrameHasOneFrameByDefaultEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body>No iframes</body></html>`))

	is.Len(page.Frames(), 1)
}

// TestFrameMainFrameContentEx verifies frames[0] Content matches page SetContent.
// Ref: TestFrame.java#shouldGetMainFrameContent
func TestFrameMainFrameContentEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body><p id="p">Main content</p></body></html>`))

	content, err := page.MainFrame().Content(ctx)
	must.NoError(err)
	is.Contains(content, "Main content")
}

// TestFrameChildFrameLocatorIsCheckedEx verifies child frame locator IsChecked.
// Ref: TestFrame.java#shouldCheckChildFrameLocator
func TestFrameChildFrameLocatorIsCheckedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	must.NoError(frame.SetContent(ctx, `<input id="cb" type="checkbox" checked>`))

	checked, err := frame.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestFrameChildFrameLocatorClickEx verifies click inside iframe.
// Ref: TestFrame.java#shouldClickInChildFrame
func TestFrameChildFrameLocatorClickEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	must.NoError(frame.SetContent(ctx, `
		<button id="btn" onclick="this.dataset.clicked='yes'">Click me</button>
	`))

	must.NoError(frame.Locator("#btn").Click(ctx))

	attr, err := frame.Locator("#btn").GetAttribute(ctx, "data-clicked")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("yes", *attr)
}

// TestFrameChildSetContentTitleEx verifies SetContent + Title in child frame.
// Ref: TestFrame.java#shouldSetAndGetTitle
func TestFrameChildSetContentTitleEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe id="f" src="about:blank"></iframe>`))

	frames := page.Frames()
	is.Greater(len(frames), 1)
	frame := frames[1]

	must.NoError(frame.SetContent(ctx, `<html><head><title>Frame Title</title></head><body></body></html>`))

	title, err := frame.Title(ctx)
	must.NoError(err)
	is.Equal("Frame Title", title)
}
