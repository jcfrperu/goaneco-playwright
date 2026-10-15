//go:build e2e

// E2E tests for: TestLocatorElementHandle.java.
package e2e

import (
	"fmt"
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestLocatorElementHandle.java#shouldQueryExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorElementHandle.java
func TestLocatorElementHandleShouldQueryExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/playground.html"))
	must.NoError(page.SetContent(ctx, "<html><body><div class='second'><div class='inner'>A</div></div></body></html>"))
	html := page.Locator("html")
	second := html.Locator(".second")
	inner := second.Locator(".inner")
	eh, err := inner.ElementHandle(ctx)
	must.NoError(err)
	content, err := page.Evaluate(ctx, "e => e.textContent", eh)
	must.NoError(err)
	is.Equal("A", content)
}

// Ref: TestLocatorElementHandle.java#shouldQueryExistingElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorElementHandle.java
func TestLocatorElementHandleShouldQueryExistingElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<html><body><div>A</div><br/><div>B</div></body></html>"))
	html := page.Locator("html")
	elements, err := html.Locator("div").ElementHandles(ctx)
	must.NoError(err)
	must.Len(elements, 2)
	texts := make([]string, len(elements))
	for i, el := range elements {
		txt, err := page.Evaluate(ctx, "e => e.textContent", el)
		must.NoError(err)
		texts[i] = fmt.Sprintf("%v", txt)
	}
	is.Equal([]string{"A", "B"}, texts)
}

// Ref: TestLocatorElementHandle.java#shouldReturnEmptyArrayForNonExistingElements  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorElementHandle.java
func TestLocatorElementHandleShouldReturnEmptyArrayForNonExistingElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<html><body><span>A</span><br/><span>B</span></body></html>"))
	html := page.Locator("html")
	elements, err := html.Locator("div").ElementHandles(ctx)
	must.NoError(err)
	is.Len(elements, 0)
}

// Ref: TestLocatorElementHandle.java#xpathShouldQueryExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorElementHandle.java
func TestLocatorElementHandleXpathShouldQueryExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/playground.html"))
	must.NoError(page.SetContent(ctx, "<html><body><div class='second'><div class='inner'>A</div></div></body></html>"))
	html := page.Locator("html")
	second := html.Locator("xpath=./body/div[contains(@class, 'second')]")
	inner := second.Locator("xpath=./div[contains(@class, 'inner')]")
	eh, err := inner.ElementHandle(ctx)
	must.NoError(err)
	content, err := page.Evaluate(ctx, "e => e.textContent", eh)
	must.NoError(err)
	is.Equal("A", content)
}

// Ref: TestLocatorElementHandle.java#xpathShouldReturnNullForNonExistingElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorElementHandle.java
func TestLocatorElementHandleXpathShouldReturnNullForNonExistingElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<html><body><div class='second'><div class='inner'>B</div></div></body></html>"))
	html := page.Locator("html")
	second, err := html.Locator("xpath=/div[contains(@class, 'third')]").ElementHandles(ctx)
	must.NoError(err)
	is.Empty(second)
}
