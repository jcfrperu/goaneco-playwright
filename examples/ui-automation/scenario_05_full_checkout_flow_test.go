//go:build e2e

package uiautomation

// Feature: Full checkout flow
// Scenario: Add one item, proceed through checkout info and overview,
// complete the purchase, and verify the confirmation message.
//
// Feature file: features/05_full_checkout_flow.feature

import (
	"strings"
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullCheckoutFlow(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	must.NoError(inv.AddToCart(ctx, "Sauce Labs Backpack"), "add backpack to cart")
	must.NoError(inv.GoToCart(ctx), "go to cart")

	cart := pom.NewCartPage(page)
	must.NoError(cart.Checkout(ctx), "click checkout")

	info := pom.NewCheckoutInfoPage(page)
	must.NoError(info.FillInfo(ctx, "John", "Doe", "10001"), "fill checkout info")
	must.NoError(info.Continue(ctx), "continue to overview")

	overview := pom.NewCheckoutOverviewPage(page)

	count, err := overview.GetItemCount(ctx)
	must.NoError(err, "overview item count")
	is.Equal(1, count, "overview should list 1 item")

	subtotal, err := overview.GetSubtotal(ctx)
	must.NoError(err, "get subtotal")
	is.True(strings.HasPrefix(subtotal, "Item total:"), "subtotal should start with 'Item total:'")

	total, err := overview.GetTotal(ctx)
	must.NoError(err, "get total")
	is.True(strings.HasPrefix(total, "Total:"), "total should start with 'Total:'")

	must.NoError(overview.Finish(ctx), "click finish")

	complete := pom.NewCheckoutCompletePage(page)
	header, err := complete.GetHeader(ctx)
	must.NoError(err, "get confirmation header")
	is.Equal("Thank you for your order!", header, "confirmation header mismatch")
}
