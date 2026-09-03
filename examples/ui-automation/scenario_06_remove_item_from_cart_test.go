//go:build e2e

package uiautomation

// Feature: Remove item from cart
// Scenario: Add 3 items, navigate to the cart, remove one, and confirm
// that only 2 items remain.
//
// Feature file: features/06_remove_item_from_cart.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveItemFromCart(t *testing.T) {
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

	must.NoError(inv.GoToCart(ctx), "go to cart")

	cart := pom.NewCartPage(page)

	count, err := cart.ItemCount(ctx)
	must.NoError(err, "initial cart count")
	is.Equal(3, count, "cart should start with 3 items")

	must.NoError(cart.RemoveItem(ctx, "Sauce Labs Bike Light"), "remove bike light")

	count, err = cart.ItemCount(ctx)
	must.NoError(err, "cart count after removal")
	is.Equal(2, count, "cart should have 2 items after removal")

	names, err := cart.GetItemNames(ctx)
	must.NoError(err, "get remaining item names")
	must.NotContains(names, "Sauce Labs Bike Light", "removed item should not be in cart")
}
