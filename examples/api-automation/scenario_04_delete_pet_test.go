//go:build e2e

package apiautomation

// Feature: Delete a pet
// Scenario: Delete a pet and verify it is no longer retrievable
//
// Feature file: features/04_delete_pet.feature

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario04DeletePet(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	created := createTestPet(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
	must.NoError(err)
	defer delResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, delResp.Status())

	getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(404, getResp.Status())
}
