//go:build e2e

// E2E tests for: TestAssertThatIsInViewport.java.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestAssertThatIsInViewport.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestAssertThatIsInViewport.java
func TestAssertThatIsInViewportShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div id=big style=\"height: 10000px;\"></div>\n<div id=small>foo</div>"))
	must.NoError(playwright.Expect(page.Locator("#big")).ToBeInViewport(ctx))
	must.NoError(playwright.Expect(page.Locator("#small")).Not().ToBeInViewport(ctx))
	must.NoError(page.Locator("#small").ScrollIntoViewIfNeeded(ctx))
	must.NoError(playwright.Expect(page.Locator("#small")).ToBeInViewport(ctx))
	// the ratio=1 variant of the Java test is covered by TestAssertThatIsInViewportShouldRespectRatioOption (ratio option not available)
}

// Ref: TestAssertThatIsInViewport.java#shouldRespectRatioOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestAssertThatIsInViewport.java
func TestAssertThatIsInViewportShouldRespectRatioOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<style>body, div, html { padding: 0; margin: 0; }</style><div id=big style='height: 400vh;'></div>"))
	must.NoError(playwright.Expect(page.Locator("div")).ToBeInViewport(ctx))
	for _, ratio := range []float64{0.1, 0.2, 0.24, 0.25} {
		must.NoError(playwright.Expect(page.Locator("div")).ToBeInViewport(ctx, &playwright.ToBeInViewportOptions{Ratio: &ratio}))
	}
	for _, ratio := range []float64{0.26, 0.3, 0.7, 0.8} {
		must.NoError(playwright.Expect(page.Locator("div")).Not().ToBeInViewport(ctx, &playwright.ToBeInViewportOptions{Ratio: &ratio}))
	}
}

// Ref: TestAssertThatIsInViewport.java#shouldHaveGoodStack  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestAssertThatIsInViewport.java
func TestAssertThatIsInViewportShouldHaveGoodStack(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := playwright.Expect(page.Locator("body")).Not().ToBeInViewport(ctx)
	must.Error(err)
	assert.Contains(t, err.Error(), "not.to.be.in.viewport")
}

// Ref: TestAssertThatIsInViewport.java#shouldReportIntersectionEvenIfFullyCoveredByOtherElement  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestAssertThatIsInViewport.java
func TestAssertThatIsInViewportShouldReportIntersectionEvenIfFullyCoveredByOtherElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<h1>hello</h1>\n<div style=\"position: relative; height: 10000px; top: -5000px;></div>"))
	must.NoError(playwright.Expect(page.Locator("h1")).ToBeInViewport(ctx))
}
