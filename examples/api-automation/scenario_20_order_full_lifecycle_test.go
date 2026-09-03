//go:build e2e

package apiautomation

// Feature: Order full lifecycle
// Scenario: Place, retrieve, delete and confirm deletion of an order
//
// Feature file: features/20_order_full_lifecycle.feature

import (
	"fmt"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario20OrderFullLifecycle(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	// POST /store/order
	order := &models.Order{
		PetID:    1,
		Quantity: 2,
		Status:   "placed",
	}
	postResp, err := apiCtx.Post(ctx, "/store/order", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, order),
	})
	must.NoError(err)
	defer postResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, postResp.Status())
	postBody, err := postResp.Body(ctx)
	must.NoError(err)
	created := mustUnmarshalOrder(t, postBody)
	is.Greater(created.ID, int64(0))

	// GET /store/order/{id}
	getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/store/order/%d", created.ID))
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, getResp.Status())

	// DELETE /store/order/{id}
	delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/store/order/%d", created.ID))
	must.NoError(err)
	defer delResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, delResp.Status())

	// GET /store/order/{id} — must be 404 now
	confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/store/order/%d", created.ID))
	must.NoError(err)
	defer confirmResp.Dispose(ctx) //nolint:errcheck
	is.Equal(404, confirmResp.Status())
}
