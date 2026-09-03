//go:build e2e

package apiautomation

// Feature: Place a store order
// Scenario: Place an order and verify it has a valid ID and "placed" status
//
// Feature file: features/14_place_order.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario14PlaceOrder(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	order := &models.Order{
		PetID:    1,
		Quantity: 1,
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
	is.Greater(created.ID, int64(0))
	is.Equal("placed", created.Status)

	t.Cleanup(func() {
		cleanCtx := testCtx(t)
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deleteOrderIfExists(t, cleanAPICtx, cleanCtx, created.ID)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})
}
