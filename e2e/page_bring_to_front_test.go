//go:build e2e

// Page.BringToFront E2E tests.
// Migration of: TestPageBringToFront.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBringToFrontDoesNotThrow verifies BringToFront can be called without error.
// Ref: TestPageBringToFront.java#shouldWork
func TestBringToFrontDoesNotThrow(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>page content</div>`))
	must.NoError(page.BringToFront(ctx))
}

// TestBringToFrontWithMultiplePages verifies BringToFront works when multiple pages exist.
// Ref: TestPageBringToFront.java#shouldWorkWithMultiplePages
func TestBringToFrontWithMultiplePages(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<div>page 1</div>`))
	must.NoError(page2.SetContent(ctx, `<div>page 2</div>`))

	// BringToFront should work without error on each page
	must.NoError(page1.BringToFront(ctx))
	must.NoError(page2.BringToFront(ctx))
	must.NoError(page1.BringToFront(ctx))
}

// TestBringToFrontAfterNavigate verifies BringToFront works after navigation.
// Ref: TestPageBringToFront.java#shouldWorkAfterNavigation
func TestBringToFrontAfterNavigate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>initial</div>`))
	must.NoError(page.BringToFront(ctx))
	must.NoError(page.SetContent(ctx, `<div>after navigation</div>`))
	must.NoError(page.BringToFront(ctx))

	text, err := page.Locator("div").InnerText(ctx)
	must.NoError(err)
	is.Equal("after navigation", text)
}
