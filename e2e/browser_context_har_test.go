//go:build e2e

// E2E tests for: TestBrowserContextHar.java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

// Ref: TestBrowserContextHar.java#shouldContextRouteFromHARMatchingTheMethodAndFollowingRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldContextRouteFromHARMatchingTheMethodAndFollowingRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	must.NoError(bCtx.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har")))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "http://no.playwright/"))
	// HAR contains a redirect for the script that should be followed automatically.
	v, err := page.Evaluate(ctx, "window.value")
	must.NoError(err)
	must.Equal("foo", v)
	// HAR contains a POST for the css file that should not be used.
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgb(255, 0, 0)"))
}

// Ref: TestBrowserContextHar.java#shouldPageRouteFromHARMatchingTheMethodAndFollowingRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldPageRouteFromHARMatchingTheMethodAndFollowingRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har")))
	must.NoError(page.Goto(ctx, "http://no.playwright/"))

	value, err := page.Evaluate(ctx, "window.value")
	must.NoError(err)
	must.Equal("foo", value)
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgb(255, 0, 0)"))
}

// Ref: TestBrowserContextHar.java#fallbackContinueShouldContinueWhenNotFoundInHar  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarFallbackContinueShouldContinueWhenNotFoundInHar(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	must.NoError(bCtx.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har"), &playwright.RouteFromHAROptions{
		NotFound: playwright.HarNotFoundFallback,
	}))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/one-style.html"))
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgb(255, 192, 203)"))
}

// Ref: TestBrowserContextHar.java#byDefaultShouldAbortRequestsNotFoundInHar  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarByDefaultShouldAbortRequestsNotFoundInHar(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	must.NoError(bCtx.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har")))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.Error(page.Goto(ctx, srv.EmptyPage()))
}

// Ref: TestBrowserContextHar.java#fallbackContinueShouldContinueRequestsOnBadHar  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarFallbackContinueShouldContinueRequestsOnBadHar(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	harPath := filepath.Join(t.TempDir(), "test.har")
	must.NoError(os.WriteFile(harPath, []byte(`{"log":{}}`), 0o600))
	bCtx := newContext(t)
	must.NoError(bCtx.RouteFromHAR(ctx, harPath, &playwright.RouteFromHAROptions{
		NotFound: playwright.HarNotFoundFallback,
	}))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/one-style.html"))
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgb(255, 192, 203)"))
}

// Ref: TestBrowserContextHar.java#shouldOnlyHandleRequestsMatchingUrlFilter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldOnlyHandleRequestsMatchingUrlFilter(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	url := "**/*.js"
	// The context route below serves the document; the HAR filter only handles JS requests.
	must.NoError(bCtx.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har"), &playwright.RouteFromHAROptions{URL: &url}))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(bCtx.Route(ctx, "http://no.playwright/", func(route *playwright.Route) {
		status := 200
		ct := "text/html"
		body := "<script src='./script.js'></script><div>hello</div>"
		_ = route.Fulfill(context.Background(), &playwright.RouteFulfillOptions{Status: &status, ContentType: &ct, Body: &body})
	}))
	must.NoError(page.Goto(ctx, "http://no.playwright/"))
	// HAR contains a redirect for the script that should be followed automatically.
	v, err := page.Evaluate(ctx, "window.value")
	must.NoError(err)
	must.Equal("foo", v)
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgba(0, 0, 0, 0)"))
}

// Ref: TestBrowserContextHar.java#shouldOnlyContextRouteFromHARRequestsMatchingUrlFilter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldOnlyContextRouteFromHARRequestsMatchingUrlFilter(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	url := "**/*.js"
	must.NoError(bCtx.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har"), &playwright.RouteFromHAROptions{URL: &url}))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(bCtx.Route(ctx, "http://no.playwright/", func(route *playwright.Route) {
		status := 200
		ct := "text/html"
		body := "<script src='./script.js'></script><div>hello</div>"
		_ = route.Fulfill(context.Background(), &playwright.RouteFulfillOptions{Status: &status, ContentType: &ct, Body: &body})
	}))
	must.NoError(page.Goto(ctx, "http://no.playwright/"))
	v, err := page.Evaluate(ctx, "window.value")
	must.NoError(err)
	must.Equal("foo", v)
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgba(0, 0, 0, 0)"))
}

// Ref: TestBrowserContextHar.java#shouldOnlyPageRouteFromHARRequestsMatchingUrlFilter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldOnlyPageRouteFromHARRequestsMatchingUrlFilter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	url := "**/*.js"
	must.NoError(page.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har"), &playwright.RouteFromHAROptions{URL: &url}))
	must.NoError(bCtx.Route(ctx, "http://no.playwright/", func(route *playwright.Route) {
		status := 200
		contentType := "text/html"
		body := "<script src='./script.js'></script><div>hello</div>"
		_ = route.Fulfill(context.Background(), &playwright.RouteFulfillOptions{
			Status: &status, ContentType: &contentType, Body: &body,
		})
	}))
	must.NoError(page.Goto(ctx, "http://no.playwright/"))

	value, err := page.Evaluate(ctx, "window.value")
	must.NoError(err)
	must.Equal("foo", value)
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgba(0, 0, 0, 0)"))
}

// Ref: TestBrowserContextHar.java#shouldSupportRegexFilter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldSupportRegexFilter(t *testing.T) {
	t.Skip("not implemented: RouteFromHAROptions.URL as regexp (only glob string)")
}

// Ref: TestBrowserContextHar.java#newPageShouldFulfillFromHarMatchingTheMethodAndFollowingRedirects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarNewPageShouldFulfillFromHarMatchingTheMethodAndFollowingRedirects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.RouteFromHAR(ctx, a1ResPath(t, "har-fulfill.har")))
	must.NoError(page.Goto(ctx, "http://no.playwright/"))

	value, err := page.Evaluate(ctx, "window.value")
	must.NoError(err)
	must.Equal("foo", value)
	must.NoError(playwright.Expect(page.Locator("body")).ToHaveCSS(ctx, "background-color", "rgb(255, 0, 0)"))
}

// Ref: TestBrowserContextHar.java#shouldChangeDocumentURLAfterRedirectedNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldChangeDocumentURLAfterRedirectedNavigation(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation / Page.Goto returning *Response")
}

// Ref: TestBrowserContextHar.java#shouldChangeDocumentURLAfterRedirectedNavigationOnClick  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldChangeDocumentURLAfterRedirectedNavigationOnClick(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation returning *Response; RouteFromHAROptions.URL as regexp")
}

// Ref: TestBrowserContextHar.java#shouldGoBackToRedirectedNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldGoBackToRedirectedNavigation(t *testing.T) {
	t.Skip("not implemented: Page.GoBack returning *Response; RouteFromHAROptions.URL as regexp")
}

// Ref: TestBrowserContextHar.java#shouldReloadRedirectedNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldReloadRedirectedNavigation(t *testing.T) {
	t.Skip("not implemented: Page.Reload returning *Response; RouteFromHAROptions.URL as regexp")
}

// Ref: TestBrowserContextHar.java#shouldFulfillFromHarWithContentInAFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldFulfillFromHarWithContentInAFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	must.NoError(bCtx.RouteFromHAR(ctx, a1ResPath(t, "har-sha1.har")))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, "http://no.playwright/"))
	content, err := page.Content(ctx)
	must.NoError(err)
	must.Equal("<html><head></head><body>Hello, world</body></html>", content)
}

// Ref: TestBrowserContextHar.java#shouldRoundTripHarZip  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldRoundTripHarZip(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / RouteFromHAR zip + NotFound")
}

// Ref: TestBrowserContextHar.java#shouldProduceExtractedZip  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldProduceExtractedZip(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* (path, mode, content)")
}

// Ref: TestBrowserContextHar.java#shouldRoundTripExtractedHarZip  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldRoundTripExtractedHarZip(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar*")
}

// Ref: TestBrowserContextHar.java#shouldRoundTripHarWithPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldRoundTripHarWithPostData(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar*")
}

// Ref: TestBrowserContextHar.java#shouldDisambiguateByHeader  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldDisambiguateByHeader(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar*")
}

// Ref: TestBrowserContextHar.java#shouldUpdateHarZipForContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldUpdateHarZipForContext(t *testing.T) {
	t.Skip("not implemented: RouteFromHAROptions.Update / NotFound")
}

// Ref: TestBrowserContextHar.java#shouldUpdateHarZipForPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldUpdateHarZipForPage(t *testing.T) {
	t.Skip("not implemented: Page.RouteFromHAR (Update / NotFound)")
}

// Ref: TestBrowserContextHar.java#shouldUpdateHarZipForPageWithDifferentOptions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldUpdateHarZipForPageWithDifferentOptions(t *testing.T) {
	t.Skip("not implemented: Page.RouteFromHAR (Update, UpdateContent, UpdateMode, NotFound)")
}

// Ref: TestBrowserContextHar.java#shouldUpdateExtractedHarZipForPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldUpdateExtractedHarZipForPage(t *testing.T) {
	t.Skip("not implemented: Page.RouteFromHAR (Update / NotFound)")
}

// Ref: TestBrowserContextHar.java#shouldIgnoreAbortedRequests  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldIgnoreAbortedRequests(t *testing.T) {
	t.Skip("not implemented: RouteFromHAROptions.Update")
}

// Ref: TestBrowserContextHar.java#shouldNotInterceptAPIRequestContextRequestsByDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextHar.java
func TestBrowserContextHarShouldNotInterceptAPIRequestContextRequestsByDefault(t *testing.T) {
	t.Skip("not implemented: RouteFromHAROptions.Update / NotFound; Page.Request()")
}
