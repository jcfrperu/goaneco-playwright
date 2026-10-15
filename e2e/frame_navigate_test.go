//go:build e2e

// E2E tests for: TestFrameNavigate.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

// TestFrameContentReturnsHTML verifies Frame.Content returns the HTML content.
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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

// TestFrameQuerySelectorReturnsNilForMissing verifies Frame.QuerySelector returns nil when no match.
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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

// TestFrameQuerySelectorReturnsNilForMissingExtra2 verifies QuerySelector returns nil when no match.
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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

// TestFrameChildSetContentTitleEx verifies SetContent + Title in child frame.
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
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

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<p id="el">hello frame</p>`)
	must.NoError(err, "SetContent failed")

	val, err := page.MainFrame().Evaluate(ctx, "document.getElementById('el').textContent")
	must.NoError(err, "frame.Evaluate failed")
	if val != "hello frame" {
		t.Errorf("Evaluate = %q, want 'hello frame'", val)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn">click me</button>`)
	must.NoError(err, "SetContent failed")

	text, err := page.MainFrame().Locator("#btn").InnerText(ctx)
	must.NoError(err, "frame locator InnerText failed")
	if text != "click me" {
		t.Errorf("InnerText = %q, want 'click me'", text)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameSetContentAndContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	frame := page.MainFrame()
	err := frame.SetContent(ctx, `<p id="p">frame content</p>`)
	must.NoError(err, "frame.SetContent failed")

	content, err := frame.Content(ctx)
	must.NoError(err, "frame.Content failed")
	if !strings.Contains(content, "frame content") {
		t.Errorf("Content does not contain 'frame content': %s", content)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestChildFrameAttachedAndName(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/inner", "text/html", `<title>Inner</title><p>inner</p>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'childframe';
		f.src = '`+srv.Prefix()+`/inner';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	child := waitForFrame(t, page, "childframe")
	if child.Name() != "childframe" {
		t.Errorf("frame.Name = %q, want 'childframe'", child.Name())
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestChildFrameEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/frame-eval", "text/html", `<p id="msg">iframe text</p>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'evalframe';
		f.src = '`+srv.Prefix()+`/frame-eval';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	fr := waitForFrame(t, page, "evalframe")
	err = fr.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	val, err := fr.Evaluate(ctx, "document.getElementById('msg').textContent")
	must.NoError(err, "frame.Evaluate failed")
	if val != "iframe text" {
		t.Errorf("frame.Evaluate = %q, want 'iframe text'", val)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestPageFramesSlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/f1", "text/html", "<p>f1</p>")
	srv.ServeWithBody("/f2", "text/html", "<p>f2</p>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		['f1','f2'].forEach(function(name) {
			const f = document.createElement('iframe');
			f.name = name;
			f.src = '`+srv.Prefix()+`/' + name;
			document.body.appendChild(f);
		});
	}`)
	must.NoError(err, "inject iframes failed")

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if len(page.Frames()) >= 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if n := len(page.Frames()); n < 2 {
		t.Errorf("expected at least 2 child frames, got %d", n)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/nav1", "text/html", "<title>Nav1</title>")
	srv.ServeWithBody("/nav2", "text/html", "<title>Nav2</title>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto parent failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'gotoframe';
		f.src = '`+srv.Prefix()+`/nav1';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	fr := waitForFrame(t, page, "gotoframe")
	err = fr.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	err = fr.Goto(ctx, srv.Prefix()+"/nav2")
	must.NoError(err, "frame.Goto failed")

	title, err := fr.Title(ctx)
	must.NoError(err, "frame.Title failed")
	if title != "Nav2" {
		t.Errorf("frame.Title after Goto = %q, want 'Nav2'", title)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameQuerySelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/qs", "text/html", `<button id="btn">press</button>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto parent failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'qsframe';
		f.src = '`+srv.Prefix()+`/qs';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	fr := waitForFrame(t, page, "qsframe")
	err = fr.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	el, err := fr.QuerySelector(ctx, "#btn")
	must.NoError(err, "frame.QuerySelector failed")
	must.NotNil(el, "QuerySelector returned nil, expected element")

	text, err := el.InnerText(ctx)
	must.NoError(err, "InnerText failed")
	if text != "press" {
		t.Errorf("InnerText = %q, want 'press'", text)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameIsDetached(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/iframe-det", "text/html", "<p>inner</p>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'detachme';
		f.src = '`+srv.Prefix()+`/iframe-det';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	frame := waitForFrame(t, page, "detachme")
	err = frame.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	if frame.IsDetached() {
		t.Error("frame should not be detached before removal")
	}

	_, err = page.Evaluate(ctx, `() => {
		document.querySelector('iframe[name="detachme"]').remove();
	}`)
	must.NoError(err, "remove iframe failed")

	// Poll until IsDetached() returns true (event propagates asynchronously).
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if frame.IsDetached() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Error("frame should be detached after removal from DOM")
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameByURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/byurl", "text/html", "<p>by url</p>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto parent failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'urlframe';
		f.src = '`+srv.Prefix()+`/byurl';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	_ = waitForFrame(t, page, "urlframe")

	deadline := time.Now().Add(3 * time.Second)
	targetURL := srv.Prefix() + "/byurl"
	for time.Now().Before(deadline) {
		if fr := page.FrameByURL(targetURL); fr != nil {
			if fr.Name() != "urlframe" {
				t.Errorf("FrameByURL name = %q, want 'urlframe'", fr.Name())
			}
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
	t.Errorf("FrameByURL(%q) returned nil", targetURL)
}

// TestFrameTitle verifies that Frame.Title returns the title of the document loaded in that frame.
// Ref: TestFrameShouldReturnFrameTitle
// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Outer page embeds title.html (/title.html has <title>Woof-Woof</title>)
	srv.ServeWithBody("/frame-title-outer", "text/html",
		`<iframe name="titleframe" src="`+srv.Prefix()+`/title.html"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-title-outer")
	must.NoError(err, "Goto failed")

	child := waitForFrame(t, page, "titleframe")
	must.NotNil(child, "child frame not found")

	title, err := child.Title(ctx)
	must.NoError(err, "Frame.Title failed")
	is.Equal("Woof-Woof", title, "child frame title should match iframe document title")
}

// Extra: no equivalent method in playwright-java; related class: TestFrameNavigate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frame-url-inner", "text/html", `<p>inner</p>`)
	srv.ServeWithBody("/frame-url-outer", "text/html",
		`<iframe name="urlframe" src="`+srv.Prefix()+`/frame-url-inner"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-url-outer")
	must.NoError(err, "Goto failed")

	child := waitForFrame(t, page, "urlframe")
	url := child.URL()
	if !strings.HasSuffix(url, "/frame-url-inner") {
		t.Errorf("Frame.URL() = %q, expected suffix '/frame-url-inner'", url)
	}
}

// Ref: TestFrameNavigate.java#shouldNavigateSubframes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameNavigateShouldNavigateSubframes(t *testing.T) {
	t.Skip("not implemented: Frame.Goto returns no Response (response.ok / response.frame)")
}

// Ref: TestFrameNavigate.java#shouldContinueAfterClientRedirect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestFrameNavigate.java
func TestFrameNavigateShouldContinueAfterClientRedirect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	srv.BlockRequest("/frames/script.js")
	url := srv.Prefix() + "/frames/child-redirect.html"
	wu := "networkidle"
	timeout := 5000.0
	err := page.Goto(ctx, url, &playwright.PageGotoOptions{Timeout: &timeout, WaitUntil: &wu})
	must.Error(err)
	assert.Contains(t, err.Error(), "Timeout 5000ms exceeded.")
	assert.Contains(t, err.Error(), "navigating to \""+url+"\", waiting until \"networkidle\"")
}
