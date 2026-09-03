//go:build e2e

package apiautomation

// Feature: Get a non-existent order
// Scenario: Requesting an order that does not exist returns 404
//
// Feature file: features/18_get_order_not_found.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario18GetOrderNotFound(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/store/order/999999")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(404, resp.Status())
	is.False(resp.OK())
}
