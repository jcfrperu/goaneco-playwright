//go:build e2e

// E2E tests for: TestPageSetContent.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageSetContentShouldWork verifies that setContent + content roundtrip works properly.
// Ref: TestPageSetContent.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	expectedOutput := "<html><head></head><body><div>hello</div></body></html>"

	err := page.SetContent(ctx, "<div>hello</div>")
	must.NoError(err, "SetContent() failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	is.Equal(expectedOutput, content)
}

// TestPageSetContentShouldWorkWithDoctype verifies setContent with a full DOCTYPE declaration.
// Ref: TestPageSetContent.java#shouldWorkWithDoctype  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithDoctype(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	doctype := "<!DOCTYPE html>"
	err := page.SetContent(ctx, doctype+"<div>hello</div>")
	must.NoError(err, "SetContent() failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")

	expectedOutput := doctype + "<html><head></head><body><div>hello</div></body></html>"
	is.Equal(expectedOutput, content)
}

// TestPageSetContentShouldWorkWithEmojis verifies that HTML content containing emojis is handled properly.
// Ref: TestPageSetContent.java#shouldWorkWithEmojis  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithEmojis(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div>🐥</div>")
	must.NoError(err, "SetContent() failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	is.Contains(content, "🐥")
}

// TestPageSetContentShouldWorkWithAccents verifies that HTML content containing accented characters is handled properly.
// Ref: TestPageSetContent.java#shouldWorkWithAccents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithAccents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div>aberración</div>")
	must.NoError(err, "SetContent() failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	is.Contains(content, "aberración")
}

// TestPageSetContentShouldWorkFastEnough verifies rapid successive setContent invocations without memory/performance issues.
// Ref: TestPageSetContent.java#shouldWorkFastEnough  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkFastEnough(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	for i := range 20 {
		err := page.SetContent(ctx, "<div>yo</div>")
		must.NoErrorf(err, "SetContent() failed at iteration %d", i)
	}
}

// TestPageSetContentShouldWorkWithNewline verifies newline characters within elements are preserved by evaluating textContent.
// Ref: TestPageSetContent.java#shouldWorkWithNewline  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithNewline(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div>\n</div>")
	must.NoError(err, "SetContent() failed")

	// Strong verification by evaluating DOM element textContent
	text, err := page.Evaluate(ctx, "() => document.querySelector('div').textContent")
	must.NoError(err, "Evaluate failed")
	is.Equal("\n", text)
}

// TestPageSetContentShouldWorkWithHTML4Doctype verifies SetContent with an HTML4 DOCTYPE declaration.
// Ref: TestPageSetContent.java#shouldWorkWithHTML4Doctype  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithHTML4Doctype(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	doctype := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`
	must.NoError(page.SetContent(ctx, doctype+"<div>hello</div>"))

	content, err := page.Content(ctx)
	must.NoError(err)

	expected := doctype + "<html><head></head><body><div>hello</div></body></html>"
	is.Equal(expected, content)
}

// TestPageSetContentShouldWorkWithTrickyContent verifies SetContent handles tricky content (control characters).
// Ref: TestPageSetContent.java#shouldWorkWithTrickyContent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithTrickyContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>hello world</div>"+"\x7F"))
	val, err := page.EvalOnSelector(ctx, "div", "div => div.textContent")
	must.NoError(err)
	is.Equal("hello world", val)
}

// TestPageTitleFromSetContent verifies Title returns the title set via SetContent.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleFromSetContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>My Page</title></head><body></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("My Page", title)
}

// TestPageTitleFromNavigation verifies Title returns the title after navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleFromNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/title.html"))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Woof-Woof", title)
}

// TestPageTitleEmptyWhenNotSet verifies Title returns empty string when no <title> element.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleEmptyWhenNotSet(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head></head><body>no title</body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("", title)
}

// TestPageContentContainsSetHTML verifies Content returns the full page HTML after SetContent.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageContentContainsSetHTML(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="unique-marker">marker</div>`))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "unique-marker")
	is.Contains(content, "marker")
}

// TestPageContentAfterNavigation verifies Content works after navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageContentAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Contains(content, "<html")
	is.Contains(content, "</html>")
}

// TestPageTitleReturnsHTMLTitleEx2 verifies Title returns HTML title tag content.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleReturnsHTMLTitleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>My Page</title></head><body></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("My Page", title)
}

// TestPageTitleEmptyWhenNoTitleTagEx2 verifies Title returns empty string when no title tag.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleEmptyWhenNoTitleTagEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head></head><body><p>No title</p></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Empty(title)
}

// TestPageTitleAfterNavigationEx2 verifies Title updates after navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleAfterNavigationEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/page1", "text/html", `<html><head><title>Page One</title></head><body></body></html>`)
	srv.ServeWithBody("/page2", "text/html", `<html><head><title>Page Two</title></head><body></body></html>`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page1"))
	t1, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page One", t1)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/page2"))
	t2, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page Two", t2)
}

// TestPageTitleWithSpecialCharsEx2 verifies Title handles special characters.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleWithSpecialCharsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Hello &amp; World</title></head><body></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Hello & World", title)
}

// TestPageTitleUpdatesViaEvaluateEx3 verifies Title reflects dynamic changes.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleUpdatesViaEvaluateEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Original</title></head><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => document.title = 'Updated'`)
	must.NoError(err)

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Updated", title)
}

// TestPageTitleLongTitleEx3 verifies Title handles long titles.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleLongTitleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	longTitle := "A very long title that goes on and on and on and on and on and on and on"
	must.NoError(page.SetContent(ctx, `<html><head><title>`+longTitle+`</title></head><body></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal(longTitle, title)
}

// TestPageTitleUnicode3Ex verifies Title handles unicode characters.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleUnicode3Ex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>日本語タイトル</title></head><body></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("日本語タイトル", title)
}

// TestPageTitleWithSpecialCharsEx4 verifies Title handles special HTML characters.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleWithSpecialCharsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Hello &amp; World</title></head><body></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Hello & World", title)
}

// TestPageTitleWithNumbersEx4 verifies Title handles numeric values.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleWithNumbersEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Page 42</title></head><body></body></html>`))

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Page 42", title)
}

// TestPageTitleMultipleUpdatesEx4 verifies Title returns latest value after multiple updates.
// Extra: no equivalent method in playwright-java; related class: TestPageSetContent.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageTitleMultipleUpdatesEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>First</title></head><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => document.title = 'Second'`)
	must.NoError(err)
	_, err = page.Evaluate(ctx, `() => document.title = 'Third'`)
	must.NoError(err)

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Third", title)
}

// Ref: TestPageSetContent.java#shouldIncludeShadowRoots  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldIncludeShadowRoots(t *testing.T) {
	t.Skip("not supported by playwright-core v1.61.1: driver frame.content ignores params and does not support includeShadow")
}

// Ref: TestPageSetContent.java#shouldWorkWithDomcontentloaded  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithDomcontentloaded(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	wu := "domcontentloaded"
	must.NoError(page.SetContent(ctx, "<div>hello</div>", &playwright.SetContentOptions{WaitUntil: &wu}))

	result, err := page.Content(ctx)
	must.NoError(err)
	is.Equal("<html><head></head><body><div>hello</div></body></html>", result)
}

// Ref: TestPageSetContent.java#shouldWorkWithCommit  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldWorkWithCommit(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	wu := "commit"
	must.NoError(page.SetContent(ctx, "<div>hello</div>", &playwright.SetContentOptions{WaitUntil: &wu}))

	result, err := page.Content(ctx)
	must.NoError(err)
	is.Equal("<html><head></head><body><div>hello</div></body></html>", result)
}

// Ref: TestPageSetContent.java#shouldRespectTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.BlockRequest("/img.png")
	timeout := float64(100)
	err := page.SetContent(ctx, "<img src='"+srv.Prefix()+"/img.png'></img>", &playwright.SetContentOptions{Timeout: &timeout})
	must.Error(err)
	is.Contains(err.Error(), "Timeout 100ms exceeded.")
}

// Ref: TestPageSetContent.java#shouldRespectDefaultNavigationTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetContent.java
func TestPageSetContentShouldRespectDefaultNavigationTimeout(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.SetDefaultNavigationTimeout(ctx, 100))
	srv.BlockRequest("/img.png")
	err := page.SetContent(ctx, "<img src='"+srv.Prefix()+"/img.png'></img>")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Timeout 100ms exceeded.")
}
