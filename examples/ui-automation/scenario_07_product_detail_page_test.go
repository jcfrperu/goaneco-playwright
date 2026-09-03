//go:build e2e

package uiautomation

// Feature: Product detail page
// Scenario: Click a product on the inventory, verify the detail page shows the
// correct name and price, then navigate back to the inventory.
//
// Feature file: features/07_product_detail_page.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductDetailPage(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	must.NoError(inv.ClickProduct(ctx, "Sauce Labs Backpack"), "click product")

	detail := pom.NewProductDetailPage(page)

	name, err := detail.GetName(ctx)
	must.NoError(err, "get product name on detail page")
	is.Equal("Sauce Labs Backpack", name, "product name mismatch on detail page")

	price, err := detail.GetPrice(ctx)
	must.NoError(err, "get product price on detail page")
	is.Equal("$29.99", price, "product price mismatch on detail page")

	desc, err := detail.GetDescription(ctx)
	must.NoError(err, "get product description")
	is.NotEmpty(desc, "product description should not be empty")

	must.NoError(detail.BackToInventory(ctx), "navigate back to inventory")

	count, err := inv.ProductCount(ctx)
	must.NoError(err, "product count after returning to inventory")
	is.Equal(6, count, "inventory should show 6 products after returning")
}
