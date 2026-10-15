//go:build e2e

package pom

import (
	"context"
	"fmt"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

const CheckoutOverviewURL = "https://www.saucedemo.com/checkout-step-two.html"

// CheckoutOverviewPage encapsulates step 2 of the checkout flow: order summary.
type CheckoutOverviewPage struct {
	page *playwright.Page
}

// NewCheckoutOverviewPage returns a CheckoutOverviewPage bound to the given page.
func NewCheckoutOverviewPage(page *playwright.Page) *CheckoutOverviewPage {
	return &CheckoutOverviewPage{page: page}
}

// GetItemNames returns the names of all items in the order summary.
func (co *CheckoutOverviewPage) GetItemNames(ctx context.Context) ([]string, error) {
	return co.page.Locator(".cart_item .inventory_item_name").AllInnerTexts(ctx)
}

// GetItemCount returns the number of items in the order summary.
func (co *CheckoutOverviewPage) GetItemCount(ctx context.Context) (int, error) {
	return co.page.Locator(".cart_item").Count(ctx)
}

// GetTotal returns the full "Total:" label text (e.g. "Total: $32.39").
func (co *CheckoutOverviewPage) GetTotal(ctx context.Context) (string, error) {
	text, err := co.page.Locator(".summary_total_label").InnerText(ctx)
	if err != nil {
		return "", fmt.Errorf("checkoutOverview getTotal: %w", err)
	}
	return text, nil
}

// GetSubtotal returns the item subtotal text (e.g. "Item total: $29.99").
func (co *CheckoutOverviewPage) GetSubtotal(ctx context.Context) (string, error) {
	text, err := co.page.Locator(".summary_subtotal_label").InnerText(ctx)
	if err != nil {
		return "", fmt.Errorf("checkoutOverview getSubtotal: %w", err)
	}
	return text, nil
}

// Finish clicks the Finish button to complete the purchase.
func (co *CheckoutOverviewPage) Finish(ctx context.Context) error {
	if err := co.page.Locator(`[data-test="finish"]`).Click(ctx); err != nil {
		return fmt.Errorf("checkoutOverview finish: %w", err)
	}
	return nil
}

// Cancel returns to the cart without completing the purchase.
func (co *CheckoutOverviewPage) Cancel(ctx context.Context) error {
	if err := co.page.Locator(`[data-test="cancel"]`).Click(ctx); err != nil {
		return fmt.Errorf("checkoutOverview cancel: %w", err)
	}
	return nil
}
