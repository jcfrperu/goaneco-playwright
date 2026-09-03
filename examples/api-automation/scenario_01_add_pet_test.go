//go:build e2e

package apiautomation

// Feature: Add a new pet
// Scenario: Successfully add a pet with valid data
//
// Feature file: features/01_add_pet.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario01AddPet(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	pet := &models.Pet{
		Name:      "Buddy",
		PhotoURLs: []string{"https://example.com/img.jpg"},
		Status:    "available",
	}

	resp, err := apiCtx.Post(ctx, "/pet", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, pet),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	created := mustUnmarshalPet(t, body)
	is.Greater(created.ID, int64(0))
	is.Equal("Buddy", created.Name)
	is.Equal("available", created.Status)

	t.Cleanup(func() {
		cleanCtx := testCtx(t)
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deletePetIfExists(t, cleanAPICtx, cleanCtx, created.ID)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})
}
