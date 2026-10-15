//go:build e2e

package apiautomation

// Feature: Update a user
// Scenario: Update a user's email and verify the change is persisted
//
// Feature file: features/23_update_user.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario23UpdateUser(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	user := createTestUser(t, "updateuser23")

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	updated := &models.User{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     "updated-" + user.Username + "@example.com",
		Password:  user.Password,
	}

	putResp, err := apiCtx.Put(ctx, "/user/"+user.Username, &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, updated),
	})
	must.NoError(err)
	defer putResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, putResp.Status())

	getResp, err := apiCtx.Get(ctx, "/user/"+user.Username)
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, getResp.Status())

	body, err := getResp.Body(ctx)
	must.NoError(err)
	retrieved := mustUnmarshalUser(t, body)
	is.Equal(updated.Email, retrieved.Email)
}
