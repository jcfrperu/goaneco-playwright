//go:build e2e

package apiautomation

// Feature: Create users with array
// Scenario: Create multiple users via POST /user/createWithArray and verify they are retrievable
//
// Feature file: features/27_create_users_with_array.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario27CreateUsersWithArray(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	users := []models.User{
		{Username: "goaneco-array27a", Email: "array27a@example.com", Password: "pass123"},
		{Username: "goaneco-array27b", Email: "array27b@example.com", Password: "pass123"},
	}

	resp, err := apiCtx.Post(ctx, "/user/createWithArray", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, users),
	})
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())

	t.Cleanup(func() {
		cleanCtx := testCtx(t)
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			for _, u := range users {
				deleteUserIfExists(t, cleanAPICtx, cleanCtx, u.Username)
			}
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})

	for _, u := range users {
		getResp, err := apiCtx.Get(ctx, "/user/"+u.Username)
		must.NoError(err)
		defer getResp.Dispose(ctx) //nolint:errcheck
		is.Equal(200, getResp.Status(), "user %s should be retrievable", u.Username)
	}
}
