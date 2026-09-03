//go:build e2e

package apiautomation

// Feature: Update a pet
// Scenario: Update an existing pet's status to "sold"
//
// Feature file: features/03_update_pet.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario03UpdatePet(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	created := createTestPet(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	updated := &models.Pet{
		ID:        created.ID,
		Name:      created.Name,
		PhotoURLs: created.PhotoURLs,
		Status:    "sold",
	}

	resp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, updated),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	result := mustUnmarshalPet(t, body)
	is.Equal(created.ID, result.ID)
	is.Equal("sold", result.Status)
}
