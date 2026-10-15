//go:build e2e

// E2E tests for: TestPageLocatorConvenience.java.
package e2e

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Ref: TestPageLocatorConvenience.java#shouldReturnPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageLocatorConvenience.java
func TestPageLocatorConvenienceShouldReturnPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	bc := page.Context()
	must.NotNil(bc)
}
