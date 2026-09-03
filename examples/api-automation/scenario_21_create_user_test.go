//go:build e2e

package apiautomation

// Feature: Create a user
// Scenario: Create a user and verify the account is retrievable
//
// Feature file: features/21_create_user.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario21CreateUser(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	user := &models.User{
		Username:  "goaneco-create21",
		FirstName: "Go",
		LastName:  "Aneco",
		Email:     "goaneco-create21@example.com",
		Password:  "pass123",
	}

	postResp, err := apiCtx.Post(ctx, "/user", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, user),
	})
	must.NoError(err)
	defer postResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, postResp.Status())

	t.Cleanup(func() {
		cleanCtx := testCtx(t)
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deleteUserIfExists(t, cleanAPICtx, cleanCtx, user.Username)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})

	getResp, err := apiCtx.Get(ctx, "/user/"+user.Username)
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, getResp.Status())

	body, err := getResp.Body(ctx)
	must.NoError(err)
	retrieved := mustUnmarshalUser(t, body)
	is.Equal(user.Username, retrieved.Username)
	is.Equal(user.Email, retrieved.Email)
}
