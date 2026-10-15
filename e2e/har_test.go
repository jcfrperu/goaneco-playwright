//go:build e2e

// E2E tests for: TestHar.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestHar.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestRouteFromHAR(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	targetURL := srv.Prefix() + "/har-served"
	harPath := writeHAR(t, targetURL, "text/plain", "served from HAR", 200)

	bCtx := newContext(t)
	err := bCtx.RouteFromHAR(ctx, harPath)
	must.NoError(err, "RouteFromHAR failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, targetURL)
	must.NoError(err, "Goto failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	is.Contains(content, "served from HAR")
}

// Extra: no equivalent method in playwright-java; related class: TestHar.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestRouteFromHARPassThrough(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/real-page", "text/plain", "from real server")

	// HAR only has an entry for a different URL — real requests must pass through.
	harPath := writeHAR(t, srv.Prefix()+"/non-existent", "text/plain", "never", 200)

	bCtx := newContext(t)
	err := bCtx.RouteFromHAR(ctx, harPath)
	must.NoError(err, "RouteFromHAR failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.Prefix()+"/real-page")
	must.NoError(err, "Goto failed")

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	if !strings.Contains(content, "from real server") {
		t.Errorf("expected pass-through response 'from real server', got:\n%s", content)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestHar.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestRouteFromHARStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	targetURL := srv.Prefix() + "/har-404"
	harPath := writeHAR(t, targetURL, "text/plain", "not found in HAR", 404)

	bCtx := newContext(t)
	err := bCtx.RouteFromHAR(ctx, harPath)
	must.NoError(err, "RouteFromHAR failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// Goto a URL that maps to a 404 in the HAR — should not error (browser loads it).
	_ = page.Goto(ctx, targetURL)

	content, err := page.Content(ctx)
	must.NoError(err, "Content() failed")
	if !strings.Contains(content, "not found in HAR") {
		t.Errorf("expected HAR 404 body in content, got:\n%s", content)
	}
}

// Ref: TestHar.java#shouldHaveVersionAndCreator  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldHaveVersionAndCreator(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldHaveBrowser  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldHaveBrowser(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldHavePages  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldHavePages(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldHavePagesInPersistentContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldHavePagesInPersistentContext(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldIncludeRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldIncludeRequest(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldIncludeResponse  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldIncludeResponse(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldFilterByGlob  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldFilterByGlob(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldFilterByRegexp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldFilterByRegexp(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldOmitContent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldOmitContent(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldOmitContentLegacy  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldOmitContentLegacy(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}

// Ref: TestHar.java#shouldAttachContent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestHar.java
func TestHarShouldAttachContent(t *testing.T) {
	t.Skip("not implemented: BrowserContextOptions.RecordHar* / HAR recording")
}
