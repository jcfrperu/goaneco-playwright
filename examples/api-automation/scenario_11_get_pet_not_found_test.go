//go:build e2e

package apiautomation

// Feature: Get a non-existent pet
// Scenario: Requesting a pet that does not exist returns 404
//
// Feature file: features/11_get_pet_not_found.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario11GetPetNotFound(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/pet/999999999")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(404, resp.Status())
	is.False(resp.OK())
}
