//go:build e2e

package uiautomation

// Feature: Add items to cart
// Scenario: Add 3 products and verify their names and prices appear in the cart.
//
// Feature file: features/02_add_items_to_cart.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddItemsToCart(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	products := []string{
		"Sauce Labs Backpack",
		"Sauce Labs Bike Light",
		"Sauce Labs Bolt T-Shirt",
	}

	for _, name := range products {
		must.NoError(inv.AddToCart(ctx, name), "add %q to cart", name)
	}

	badge, err := inv.CartBadgeCount(ctx)
	must.NoError(err, "read cart badge")
	is.Equal(3, badge, "cart badge should show 3")

	must.NoError(inv.GoToCart(ctx), "navigate to cart")

	cart := pom.NewCartPage(page)

	count, err := cart.ItemCount(ctx)
	must.NoError(err, "cart item count")
	is.Equal(3, count, "cart should have 3 items")

	names, err := cart.GetItemNames(ctx)
	must.NoError(err, "get cart item names")
	is.ElementsMatch(products, names, "cart names should match added products")

	prices, err := cart.GetItemPrices(ctx)
	must.NoError(err, "get cart item prices")
	is.Len(prices, 3, "cart should show 3 prices")
	for _, p := range prices {
		is.NotEmpty(p, "each price should be non-empty")
	}
}
