//go:build e2e

package pom

import (
	"context"
	"fmt"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

const CheckoutInfoURL = "https://www.saucedemo.com/checkout-step-one.html"

// CheckoutInfoPage encapsulates step 1 of the checkout flow: customer information.
type CheckoutInfoPage struct {
	page *playwright.Page
}

// NewCheckoutInfoPage returns a CheckoutInfoPage bound to the given page.
func NewCheckoutInfoPage(page *playwright.Page) *CheckoutInfoPage {
	return &CheckoutInfoPage{page: page}
}

// FillInfo fills the First Name, Last Name, and Postal Code fields.
func (ci *CheckoutInfoPage) FillInfo(ctx context.Context, firstName, lastName, postalCode string) error {
	if err := ci.page.Locator("#first-name").Fill(ctx, firstName); err != nil {
		return fmt.Errorf("fillInfo first name: %w", err)
	}
	if err := ci.page.Locator("#last-name").Fill(ctx, lastName); err != nil {
		return fmt.Errorf("fillInfo last name: %w", err)
	}
	if err := ci.page.Locator("#postal-code").Fill(ctx, postalCode); err != nil {
		return fmt.Errorf("fillInfo postal code: %w", err)
	}
	return nil
}

// Continue clicks the Continue button to proceed to the order overview.
func (ci *CheckoutInfoPage) Continue(ctx context.Context) error {
	if err := ci.page.Locator(`[data-test="continue"]`).Click(ctx); err != nil {
		return fmt.Errorf("checkoutInfo continue: %w", err)
	}
	return nil
}

// ErrorMessage returns the validation error text if displayed, or empty string.
func (ci *CheckoutInfoPage) ErrorMessage(ctx context.Context) (string, error) {
	loc := ci.page.Locator("[data-test='error']")
	visible, err := loc.IsVisible(ctx)
	if err != nil || !visible {
		return "", err
	}
	return loc.InnerText(ctx)
}
