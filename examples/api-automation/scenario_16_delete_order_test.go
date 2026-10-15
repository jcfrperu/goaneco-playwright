//go:build e2e

package apiautomation

// Feature: Delete a store order
// Scenario: Delete an order and verify it is no longer retrievable
//
// Feature file: features/16_delete_order.feature

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario16DeleteOrder(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	pet := createTestPet(t)
	order := createTestOrder(t, pet.ID)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/store/order/%d", order.ID))
	must.NoError(err)
	defer delResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, delResp.Status())

	getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/store/order/%d", order.ID))
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(404, getResp.Status())
}
