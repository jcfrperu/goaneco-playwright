//go:build e2e

package apiautomation

// Feature: Get an order by ID
// Scenario: Retrieve an existing order by its ID and verify fields match
//
// Feature file: features/15_get_order_by_id.feature

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario15GetOrderByID(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	pet := createTestPet(t)
	order := createTestOrder(t, pet.ID)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, fmt.Sprintf("/store/order/%d", order.ID))
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	retrieved := mustUnmarshalOrder(t, body)
	is.Equal(order.ID, retrieved.ID)
	is.Equal(pet.ID, retrieved.PetID)
	is.Equal(1, retrieved.Quantity)
}
