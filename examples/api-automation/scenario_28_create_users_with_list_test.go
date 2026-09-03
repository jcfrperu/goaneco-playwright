//go:build e2e

package apiautomation

// Feature: Create users with list
// Scenario: Create multiple users via POST /user/createWithList and verify they are retrievable
//
// Feature file: features/28_create_users_with_list.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario28CreateUsersWithList(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	users := []models.User{
		{Username: "goaneco-list28a", Email: "list28a@example.com", Password: "pass123"},
		{Username: "goaneco-list28b", Email: "list28b@example.com", Password: "pass123"},
	}

	resp, err := apiCtx.Post(ctx, "/user/createWithList", &playwright.APIFetchOptions{
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
