//go:build e2e

package apiautomation

// Feature: Place an order with a specific quantity
// Scenario: Place an order with quantity=5 and verify the quantity in the response
//
// Feature file: features/19_place_order_with_quantity.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario19PlaceOrderWithQuantity(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	order := &models.Order{
		PetID:    1,
		Quantity: 5,
		Status:   "placed",
	}

	resp, err := apiCtx.Post(ctx, "/store/order", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, order),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	created := mustUnmarshalOrder(t, body)
	is.Equal(5, created.Quantity)

	t.Cleanup(func() {
		cleanCtx := testCtx(t)
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deleteOrderIfExists(t, cleanAPICtx, cleanCtx, created.ID)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})
}
