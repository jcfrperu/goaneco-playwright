//go:build e2e

package pom

import (
	"context"
	"fmt"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

// ProductDetailPage encapsulates all interactions with a SauceDemo product detail page.
type ProductDetailPage struct {
	page *playwright.Page
}

// NewProductDetailPage returns a ProductDetailPage bound to the given page.
func NewProductDetailPage(page *playwright.Page) *ProductDetailPage {
	return &ProductDetailPage{page: page}
}

// GetName returns the product name displayed on the detail page.
func (pd *ProductDetailPage) GetName(ctx context.Context) (string, error) {
	text, err := pd.page.Locator(".inventory_details_name").InnerText(ctx)
	if err != nil {
		return "", fmt.Errorf("productDetail getName: %w", err)
	}
	return text, nil
}

// GetPrice returns the product price displayed on the detail page (e.g. "$29.99").
func (pd *ProductDetailPage) GetPrice(ctx context.Context) (string, error) {
	text, err := pd.page.Locator(".inventory_details_price").InnerText(ctx)
	if err != nil {
		return "", fmt.Errorf("productDetail getPrice: %w", err)
	}
	return text, nil
}

// GetDescription returns the product description text.
func (pd *ProductDetailPage) GetDescription(ctx context.Context) (string, error) {
	text, err := pd.page.Locator(".inventory_details_desc").InnerText(ctx)
	if err != nil {
		return "", fmt.Errorf("productDetail getDescription: %w", err)
	}
	return text, nil
}

// AddToCart clicks the "Add to cart" button on the detail page.
func (pd *ProductDetailPage) AddToCart(ctx context.Context) error {
	if err := pd.page.Locator(".btn_inventory").Click(ctx); err != nil {
		return fmt.Errorf("productDetail addToCart: %w", err)
	}
	return nil
}

// RemoveFromCart clicks the "Remove" button on the detail page.
func (pd *ProductDetailPage) RemoveFromCart(ctx context.Context) error {
	if err := pd.page.Locator(".btn_inventory").Click(ctx); err != nil {
		return fmt.Errorf("productDetail removeFromCart: %w", err)
	}
	return nil
}

// BackToInventory clicks the back button to return to the inventory listing.
func (pd *ProductDetailPage) BackToInventory(ctx context.Context) error {
	if err := pd.page.Locator(`[data-test="back-to-products"]`).Click(ctx); err != nil {
		return fmt.Errorf("productDetail backToInventory: %w", err)
	}
	return nil
}
