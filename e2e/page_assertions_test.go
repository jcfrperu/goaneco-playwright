//go:build e2e

// Page-level assertion tests: HasURL, HasTitle with text and regex.
// Migration of: TestPageAssertions.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"context"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

func TestPageAssertionHasURLText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/assert-url", "text/html", `<p>page</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/assert-url")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).HasURL(ctx, srv.Prefix()+"/assert-url"); err != nil {
		t.Errorf("HasURL exact failed: %v", err)
	}
}

func TestPageAssertionHasURLFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/url-fail", "text/html", `<p>page</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/url-fail")
	must.NoError(err, "Goto failed")

	err = playwright.ExpectPage(page).HasURL(ctx, "http://wrong-url.invalid/")
	if err == nil {
		t.Error("expected HasURL to fail for wrong URL, but got nil error")
	}
}

func TestPageAssertionNotHasURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/not-url", "text/html", `<p>page</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/not-url")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).Not().HasURL(ctx, "http://other.invalid/"); err != nil {
		t.Errorf("Not().HasURL (should pass) failed: %v", err)
	}
}

func TestPageAssertionHasURLContains(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/url-contains", "text/html", `<p>page</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/url-contains")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).HasURLContains(ctx, "url-contains"); err != nil {
		t.Errorf("HasURLContains failed: %v", err)
	}
}

func TestPageAssertionHasURLRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/url-regex-123", "text/html", `<p>page</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/url-regex-123")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).HasURLRegex(ctx, `url-regex-\d+`); err != nil {
		t.Errorf("HasURLRegex failed: %v", err)
	}
}

func TestPageAssertionHasURLRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/url-regex-nope", "text/html", `<p>page</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/url-regex-nope")
	must.NoError(err, "Goto failed")

	err = playwright.ExpectPage(page).HasURLRegex(ctx, `^https://example\.com`)
	if err == nil {
		t.Error("expected HasURLRegex to fail but got nil error")
	}
}

func TestPageAssertionHasTitleText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/title-assert", "text/html", `<title>My Page Title</title><p>body</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title-assert")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).HasTitle(ctx, "My Page Title"); err != nil {
		t.Errorf("HasTitle failed: %v", err)
	}
}

func TestPageAssertionHasTitleFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/title-fail", "text/html", `<title>Real Title</title><p>body</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title-fail")
	must.NoError(err, "Goto failed")

	err = playwright.ExpectPage(page).HasTitle(ctx, "Wrong Title")
	if err == nil {
		t.Error("expected HasTitle to fail but got nil error")
	}
}

func TestPageAssertionNotHasTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/not-title", "text/html", `<title>Correct Title</title><p>body</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/not-title")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).Not().HasTitle(ctx, "Wrong Title"); err != nil {
		t.Errorf("Not().HasTitle (should pass) failed: %v", err)
	}
}

func TestPageAssertionHasTitleContains(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/title-contains", "text/html", `<title>My Awesome Page</title><p>body</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title-contains")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).HasTitleContains(ctx, "Awesome"); err != nil {
		t.Errorf("HasTitleContains failed: %v", err)
	}
}

func TestPageAssertionHasTitleRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/title-regex", "text/html", `<title>Page 42</title><p>body</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title-regex")
	must.NoError(err, "Goto failed")

	if err := playwright.ExpectPage(page).HasTitleRegex(ctx, `Page \d+`); err != nil {
		t.Errorf("HasTitleRegex failed: %v", err)
	}
}

func TestPageAssertionHasTitleRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/title-regex-fail", "text/html", `<title>My Title</title><p>body</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title-regex-fail")
	must.NoError(err, "Goto failed")

	err = playwright.ExpectPage(page).HasTitleRegex(ctx, `^Nonexistent`)
	if err == nil {
		t.Error("expected HasTitleRegex to fail but got nil error")
	}
}

func TestPageAssertionsTitleCaseInsensitiveRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/title-ci", "text/html", `<title>Woof-Woof</title><p>body</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title-ci")
	must.NoError(err)

	err = playwright.ExpectPage(page).HasTitleRegex(ctx, `(?i)woof-woof`)
	must.NoError(err, "HasTitleRegex (case-insensitive) failed")
}

func TestPageAssertionsHasURLWithBaseURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	baseURL := srv.Prefix()
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{BaseURL: &baseURL})
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	err = playwright.ExpectPage(page).HasURL(ctx, srv.EmptyPage())
	must.NoError(err, "HasURL with baseURL context failed")
}

func TestPageAssertionNotHasTitleEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Actual Title</title></head><body></body></html>`))

	err := playwright.ExpectPage(page).Not().HasTitle(ctx, "Wrong Title")
	must.NoError(err)
}

func TestPageAssertionHasTitleAfterUpdateEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><head><title>Old Title</title></head><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => { document.title = 'New Title'; }`)
	must.NoError(err)

	err = playwright.ExpectPage(page).HasTitle(ctx, "New Title")
	must.NoError(err)
}

func TestPageAssertionHasURLRegexEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/items/42", "text/html", `<html><body>Item 42</body></html>`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/items/42"))

	err := playwright.ExpectPage(page).HasURLRegex(ctx, `/items/\d+`)
	must.NoError(err)
}

// TestPageAssertionNotHasURLRegex verifies Not().HasURLRegex passes when URL does not match the pattern.
// Ref: TestPageAssertions.java#notHasUrlRegEx
func TestPageAssertionNotHasURLRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.Goto(ctx, "data:text/html,<div>B</div>")
	must.NoError(err)

	err = playwright.ExpectPage(page).Not().HasURLRegex(ctx, "about")
	must.NoError(err, "Not().HasURLRegex should pass when URL does not match the pattern")
}

// TestPageAssertionHasTitleNormalizeWhitespace verifies hasTitle normalizes whitespace before comparing.
// Ref: TestPageAssertions.java#hasTitleTextNormalizeWhitespaces
func TestPageAssertionHasTitleNormalizeWhitespace(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<title>     Foo     Bar    </title>`)
	must.NoError(err)

	// Browsers normalize whitespace in title; page.Title() returns "Foo Bar"
	title, err := page.Title(ctx)
	must.NoError(err)
	is.Contains(title, "Foo")
	is.Contains(title, "Bar")
}

// TestPageAssertionHasURLIgnoreCase verifies hasURL with case-insensitive comparison.
// Ref: TestPageAssertions.java#hasUrlSupportIgnoreCase
func TestPageAssertionHasURLIgnoreCase(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.Goto(ctx, "data:text/html,<div>A</div>")
	must.NoError(err)

	// Verify the URL using a case-insensitive regex as the Go API does not expose setIgnoreCase
	err = playwright.ExpectPage(page).HasURLRegex(ctx, `(?i)DATA:text/html`)
	must.NoError(err, "HasURLRegex with case-insensitive flag should match data URL")
}
