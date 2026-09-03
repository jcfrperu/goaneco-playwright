//go:build e2e

// Page.Title and Page.Content E2E tests.
// Migration of: TestPageTitle.java / TestPageContent.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageTitleFromSetContent verifies Title returns the title set via SetContent.
// Ref: TestPageTitle.java#shouldReturnTitleFromContent
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
// Ref: TestPageTitle.java#shouldReturnTitleAfterNavigation
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
// Ref: TestPageTitle.java#shouldReturnEmptyTitleWhenNotSet
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
// Ref: TestPageContent.java#shouldReturnContentAfterSetContent
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
// Ref: TestPageContent.java#shouldReturnContentAfterNavigation
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
// Ref: TestPageTitle.java#shouldReturnTitle
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
// Ref: TestPageTitle.java#shouldReturnEmptyWithoutTitle
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
// Ref: TestPageTitle.java#shouldUpdateAfterNavigation
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
// Ref: TestPageTitle.java#shouldHandleSpecialChars
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
// Ref: TestPageTitle.java#shouldReflectDynamicChange
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
// Ref: TestPageTitle.java#shouldHandleLongTitle
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
// Ref: TestPageTitle.java#shouldHandleUnicode
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
// Ref: TestPageTitle.java#shouldHandleSpecialChars
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
// Ref: TestPageTitle.java#shouldHandleNumbers
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
// Ref: TestPageTitle.java#shouldReturnLatestAfterMultipleUpdates
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
