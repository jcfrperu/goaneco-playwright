//go:build e2e

package apiautomation

// Feature: Pet full lifecycle
// Scenario: Create, retrieve, update, delete and confirm deletion of a pet
//
// Feature file: features/12_pet_full_lifecycle.feature

import (
	"fmt"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario12PetFullLifecycle(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	// POST /pet
	pet := &models.Pet{
		Name:      "LifecyclePet",
		PhotoURLs: []string{"https://example.com/lifecycle.jpg"},
		Status:    "available",
	}
	postResp, err := apiCtx.Post(ctx, "/pet", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, pet),
	})
	must.NoError(err)
	defer postResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, postResp.Status())
	postBody, err := postResp.Body(ctx)
	must.NoError(err)
	created := mustUnmarshalPet(t, postBody)
	is.Greater(created.ID, int64(0))

	// GET /pet/{id}
	getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, getResp.Status())

	// PUT /pet — update status to "sold"
	created.Status = "sold"
	putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, created),
	})
	must.NoError(err)
	defer putResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, putResp.Status())

	// DELETE /pet/{id}
	delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
	must.NoError(err)
	defer delResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, delResp.Status())

	// GET /pet/{id} — must be 404 now
	confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
	must.NoError(err)
	defer confirmResp.Dispose(ctx) //nolint:errcheck
	is.Equal(404, confirmResp.Status())
}
