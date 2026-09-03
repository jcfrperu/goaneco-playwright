//go:build e2e

package uiautomation

// Feature: Add most expensive item
// Scenario: Sort by price high to low, add the first (most expensive) product to
// the cart, and verify it appears in the cart with the expected price.
//
// Feature file: features/09_add_most_expensive_item.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMostExpensiveItem(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	must.NoError(inv.Sort(ctx, pom.SortPriceHighL), "sort price high to low")

	names, err := inv.GetProductNames(ctx)
	must.NoError(err, "get names after sort")
	is.NotEmpty(names, "product list must not be empty")

	prices, err := inv.GetProductPrices(ctx)
	must.NoError(err, "get prices after sort")
	is.NotEmpty(prices, "price list must not be empty")

	mostExpensiveName := names[0]
	mostExpensivePrice := prices[0]

	is.Equal("Sauce Labs Fleece Jacket", mostExpensiveName, "most expensive product name mismatch")
	is.Equal("$49.99", mostExpensivePrice, "most expensive product price mismatch")

	must.NoError(inv.AddToCart(ctx, mostExpensiveName), "add most expensive item to cart")
	must.NoError(inv.GoToCart(ctx), "go to cart")

	cart := pom.NewCartPage(page)

	cartNames, err := cart.GetItemNames(ctx)
	must.NoError(err, "get cart item names")
	is.Contains(cartNames, mostExpensiveName, "most expensive item should be in cart")

	cartPrices, err := cart.GetItemPrices(ctx)
	must.NoError(err, "get cart item prices")
	is.Contains(cartPrices, mostExpensivePrice, "most expensive price should be in cart")
}
