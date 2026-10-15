//go:build e2e

package apiautomation

// Feature: Place a completed order
// Scenario: Place an order with complete=true and verify the flag is set in the response
//
// Feature file: features/17_place_order_complete.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario17PlaceOrderComplete(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	order := &models.Order{
		PetID:    1,
		Quantity: 1,
		Status:   "placed",
		Complete: true,
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
	is.Greater(created.ID, int64(0))
	is.True(created.Complete)

	t.Cleanup(func() {
		cleanCtx := testCtx(t)
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deleteOrderIfExists(t, cleanAPICtx, cleanCtx, created.ID)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})
}
