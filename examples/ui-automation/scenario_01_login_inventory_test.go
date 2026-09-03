//go:build e2e

package uiautomation

// Feature: Login and Inventory
// Scenario: Verify that after a successful login the inventory page shows 6 products
// and the sort dropdown is visible.
//
// Feature file: features/01_login_inventory.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginAndInventory(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	count, err := inv.ProductCount(ctx)
	must.NoError(err, "count products")
	is.Equal(6, count, "inventory should list 6 products")

	names, err := inv.GetProductNames(ctx)
	must.NoError(err, "get product names")
	is.Len(names, 6, "should have 6 product names")

	// Sort dropdown must be present.
	visible, err := page.Locator(`[data-test="product_sort_container"]`).IsVisible(ctx)
	must.NoError(err, "sort dropdown visibility check")
	is.True(visible, "sort dropdown should be visible on inventory page")
}
