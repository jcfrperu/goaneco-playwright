//go:build e2e

package apiautomation

// Feature: Get a pet by ID
// Scenario: Retrieve an existing pet by its ID
//
// Feature file: features/02_get_pet_by_id.feature

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario02GetPetByID(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	created := createTestPet(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	pet := mustUnmarshalPet(t, body)
	is.Equal(created.ID, pet.ID)
	is.Equal(created.Name, pet.Name)
	is.Equal(created.Status, pet.Status)
}
