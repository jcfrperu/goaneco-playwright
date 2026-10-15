//go:build e2e

package pom

import (
	"context"
	"fmt"
	"strconv"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

const InventoryURL = "https://www.saucedemo.com/inventory.html"

// Sort option values for the product sort dropdown.
const (
	SortNameAZ     = "az"   // Name (A to Z) — default
	SortNameZA     = "za"   // Name (Z to A)
	SortPriceLowH  = "lohi" // Price (low to high)
	SortPriceHighL = "hilo" // Price (high to low)
)

// InventoryPage encapsulates all interactions with the SauceDemo products/inventory page.
type InventoryPage struct {
	page *playwright.Page
}

// NewInventoryPage returns an InventoryPage bound to the given page.
func NewInventoryPage(page *playwright.Page) *InventoryPage {
	return &InventoryPage{page: page}
}

// ProductCount returns the number of products currently displayed.
func (ip *InventoryPage) ProductCount(ctx context.Context) (int, error) {
	return ip.page.Locator(".inventory_item").Count(ctx)
}

// GetProductNames returns the visible names of all inventory items in order.
func (ip *InventoryPage) GetProductNames(ctx context.Context) ([]string, error) {
	return ip.page.Locator(".inventory_item_name").AllInnerTexts(ctx)
}

// GetProductPrices returns the visible prices of all inventory items in order.
func (ip *InventoryPage) GetProductPrices(ctx context.Context) ([]string, error) {
	return ip.page.Locator(".inventory_item_price").AllInnerTexts(ctx)
}

// AddToCart clicks "Add to cart" for the product with the given name.
func (ip *InventoryPage) AddToCart(ctx context.Context, productName string) error {
	item := ip.page.Locator(".inventory_item").Filter(&playwright.LocatorFilterOptions{HasText: &productName})
	if err := item.Locator(".btn_inventory").Click(ctx); err != nil {
		return fmt.Errorf("addToCart(%q): %w", productName, err)
	}
	return nil
}

// CartBadgeCount returns the integer displayed on the cart icon badge (0 if no badge).
func (ip *InventoryPage) CartBadgeCount(ctx context.Context) (int, error) {
	badge := ip.page.Locator(".shopping_cart_badge")
	visible, err := badge.IsVisible(ctx)
	if err != nil {
		return 0, fmt.Errorf("cartBadgeCount visible check: %w", err)
	}
	if !visible {
		return 0, nil
	}
	text, err := badge.InnerText(ctx)
	if err != nil {
		return 0, fmt.Errorf("cartBadgeCount read text: %w", err)
	}
	return strconv.Atoi(text)
}

// GoToCart clicks the shopping cart icon to navigate to the cart page.
func (ip *InventoryPage) GoToCart(ctx context.Context) error {
	if err := ip.page.Locator(".shopping_cart_link").Click(ctx); err != nil {
		return fmt.Errorf("goToCart: %w", err)
	}
	return nil
}

// Sort selects a sort option from the products dropdown.
// Use the SortName* / SortPrice* constants.
func (ip *InventoryPage) Sort(ctx context.Context, value string) error {
	_, err := ip.page.Locator(`[data-test="product_sort_container"]`).SelectOption(ctx, value)
	if err != nil {
		return fmt.Errorf("sort(%q): %w", value, err)
	}
	return nil
}

// OpenBurgerMenu opens the side navigation menu.
func (ip *InventoryPage) OpenBurgerMenu(ctx context.Context) error {
	if err := ip.page.Locator("#react-burger-menu-btn").Click(ctx); err != nil {
		return fmt.Errorf("openBurgerMenu: %w", err)
	}
	return nil
}

// Logout opens the burger menu and clicks the Logout option.
func (ip *InventoryPage) Logout(ctx context.Context) error {
	if err := ip.OpenBurgerMenu(ctx); err != nil {
		return err
	}
	if err := ip.page.Locator("#logout_sidebar_link").Click(ctx); err != nil {
		return fmt.Errorf("logout click: %w", err)
	}
	return nil
}

// ClickProduct navigates to the detail page of the product with the given name.
func (ip *InventoryPage) ClickProduct(ctx context.Context, productName string) error {
	if err := ip.page.Locator(".inventory_item_name").Filter(
		&playwright.LocatorFilterOptions{HasText: &productName},
	).Click(ctx); err != nil {
		return fmt.Errorf("clickProduct(%q): %w", productName, err)
	}
	return nil
}
