//go:build e2e

package apiautomation

// Feature: Update a pet using form data
// Scenario: Update a pet's name and status via POST with form-encoded data
//
// Feature file: features/09_update_pet_with_form.feature

import (
	"fmt"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario09UpdatePetWithForm(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	created := createTestPet(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Post(ctx, fmt.Sprintf("/pet/%d", created.ID), &playwright.APIFetchOptions{
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		FormData: []playwright.FormDataField{
			{Name: "name", Value: "UpdatedName"},
			{Name: "status", Value: "pending"},
		},
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
}
