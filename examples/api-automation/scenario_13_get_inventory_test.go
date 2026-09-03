//go:build e2e

package apiautomation

// Feature: Get store inventory
// Scenario: Retrieve the inventory and verify it is a map of status counts
//
// Feature file: features/13_get_inventory.feature

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario13GetInventory(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/store/inventory")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	var inventory map[string]int
	must.NoError(json.Unmarshal(body, &inventory), "inventory must be a map[string]int")
	is.NotEmpty(inventory, "inventory should not be empty")
}
