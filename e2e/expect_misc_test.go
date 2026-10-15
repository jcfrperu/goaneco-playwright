//go:build e2e

// E2E tests for: TestExpectMisc.java.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestExpectMisc.java#strictModeViolationErrorFormat  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestExpectMisc.java
func TestExpectMiscStrictModeViolationErrorFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, " <div>hello</div><div>hi</div>"))
	err := playwright.Expect(page.Locator("div")).ToBeVisible(ctx)
	must.Error(err)
	assert.Contains(t, err.Error(), "to.be.visible")
	assert.Contains(t, err.Error(), "strict mode violation: locator(\"div\") resolved to 2 elements:")
}

// Ref: TestExpectMisc.java#invalidSelectorErrorFormat  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestExpectMisc.java
func TestExpectMiscInvalidSelectorErrorFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>hello</div><div>hi</div>"))
	err := playwright.Expect(page.Locator("##")).ToBeVisible(ctx)
	must.Error(err)
	assert.Contains(t, err.Error(), "to.be.visible")
	assert.Contains(t, err.Error(), "Unexpected token \"#\" while parsing css selector \"##\".")
}
