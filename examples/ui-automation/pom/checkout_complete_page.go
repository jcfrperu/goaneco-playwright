//go:build e2e

package pom

import (
	"context"
	"fmt"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

const CheckoutCompleteURL = "https://www.saucedemo.com/checkout-complete.html"

// CheckoutCompletePage encapsulates the order confirmation page shown after a successful checkout.
type CheckoutCompletePage struct {
	page *playwright.Page
}

// NewCheckoutCompletePage returns a CheckoutCompletePage bound to the given page.
func NewCheckoutCompletePage(page *playwright.Page) *CheckoutCompletePage {
	return &CheckoutCompletePage{page: page}
}

// GetHeader returns the main confirmation heading (e.g. "Thank you for your order!").
func (cc *CheckoutCompletePage) GetHeader(ctx context.Context) (string, error) {
	text, err := cc.page.Locator(".complete-header").InnerText(ctx)
	if err != nil {
		return "", fmt.Errorf("checkoutComplete getHeader: %w", err)
	}
	return text, nil
}

// GetSubHeader returns the secondary confirmation text below the heading.
func (cc *CheckoutCompletePage) GetSubHeader(ctx context.Context) (string, error) {
	text, err := cc.page.Locator(".complete-text").InnerText(ctx)
	if err != nil {
		return "", fmt.Errorf("checkoutComplete getSubHeader: %w", err)
	}
	return text, nil
}

// BackHome clicks the "Back Home" button to return to the inventory page.
func (cc *CheckoutCompletePage) BackHome(ctx context.Context) error {
	if err := cc.page.Locator(`[data-test="back-to-products"]`).Click(ctx); err != nil {
		return fmt.Errorf("checkoutComplete backHome: %w", err)
	}
	return nil
}
