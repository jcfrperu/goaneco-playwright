//go:build e2e

package pom

import (
	"context"
	"fmt"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

const CartURL = "https://www.saucedemo.com/cart.html"

// CartPage encapsulates all interactions with the SauceDemo cart page.
type CartPage struct {
	page *playwright.Page
}

// NewCartPage returns a CartPage bound to the given page.
func NewCartPage(page *playwright.Page) *CartPage {
	return &CartPage{page: page}
}

// ItemCount returns the number of line items currently in the cart.
func (cp *CartPage) ItemCount(ctx context.Context) (int, error) {
	return cp.page.Locator(".cart_item").Count(ctx)
}

// GetItemNames returns the names of all items in the cart.
func (cp *CartPage) GetItemNames(ctx context.Context) ([]string, error) {
	return cp.page.Locator(".cart_item .inventory_item_name").AllInnerTexts(ctx)
}

// GetItemPrices returns the prices of all items in the cart.
func (cp *CartPage) GetItemPrices(ctx context.Context) ([]string, error) {
	return cp.page.Locator(".cart_item .inventory_item_price").AllInnerTexts(ctx)
}

// RemoveItem clicks the Remove button for the item with the given product name.
func (cp *CartPage) RemoveItem(ctx context.Context, productName string) error {
	item := cp.page.Locator(".cart_item").Filter(&playwright.LocatorFilterOptions{HasText: &productName})
	if err := item.Locator("button[data-test^='remove']").Click(ctx); err != nil {
		return fmt.Errorf("removeItem(%q): %w", productName, err)
	}
	return nil
}

// Checkout clicks the Checkout button to proceed to the checkout flow.
func (cp *CartPage) Checkout(ctx context.Context) error {
	if err := cp.page.Locator(`[data-test="checkout"]`).Click(ctx); err != nil {
		return fmt.Errorf("cart checkout: %w", err)
	}
	return nil
}

// ContinueShopping clicks the "Continue Shopping" button to return to the inventory.
func (cp *CartPage) ContinueShopping(ctx context.Context) error {
	if err := cp.page.Locator(`[data-test="continue-shopping"]`).Click(ctx); err != nil {
		return fmt.Errorf("cart continueShopping: %w", err)
	}
	return nil
}
