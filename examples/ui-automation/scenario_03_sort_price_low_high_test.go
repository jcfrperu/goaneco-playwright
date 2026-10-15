//go:build e2e

package uiautomation

// Feature: Sort by price low to high
// Scenario: After selecting "Price (low to high)" the first product listed is
// "Sauce Labs Onesie" at $7.99.
//
// Feature file: features/03_sort_price_low_high.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortByPriceLowToHigh(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	must.NoError(inv.Sort(ctx, pom.SortPriceLowH), "sort price low to high")

	names, err := inv.GetProductNames(ctx)
	must.NoError(err, "get product names after sort")
	is.NotEmpty(names, "product list must not be empty")
	is.Equal("Sauce Labs Onesie", names[0], "first product should be the cheapest")

	prices, err := inv.GetProductPrices(ctx)
	must.NoError(err, "get product prices after sort")
	is.Equal("$7.99", prices[0], "cheapest product price should be $7.99")
}
