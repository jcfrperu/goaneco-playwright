//go:build e2e

package uiautomation

// Feature: Sort by name Z to A
// Scenario: After selecting "Name (Z to A)" the first product is "Test.allTheThings() T-Shirt (Red)"
// and the last product is "Sauce Labs Backpack".
//
// Feature file: features/04_sort_name_z_to_a.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortByNameZToA(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	must.NoError(inv.Sort(ctx, pom.SortNameZA), "sort name Z to A")

	names, err := inv.GetProductNames(ctx)
	must.NoError(err, "get product names after sort")
	is.NotEmpty(names, "product list must not be empty")

	is.Equal("Test.allTheThings() T-Shirt (Red)", names[0], "first should be Z-alphabetically last product")
	is.Equal("Sauce Labs Backpack", names[len(names)-1], "last should be A-alphabetically first product")
}
