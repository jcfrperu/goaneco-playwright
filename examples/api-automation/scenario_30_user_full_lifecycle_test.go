//go:build e2e

package apiautomation

// Feature: User full lifecycle
// Scenario: Create, login, retrieve, update, delete and confirm deletion of a user
//
// Feature file: features/30_user_full_lifecycle.feature

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario30UserFullLifecycle(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	username := "goaneco-lifecycle30"
	user := &models.User{
		Username:  username,
		FirstName: "Go",
		LastName:  "Aneco",
		Email:     username + "@example.com",
		Password:  "pass123",
	}

	// POST /user
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
			deleteUserIfExists(t, cleanAPICtx, cleanCtx, username)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})

	// GET /user/login
	loginResp, err := apiCtx.Get(ctx, "/user/login?username="+username+"&password=pass123")
	must.NoError(err)
	defer loginResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, loginResp.Status())

	// GET /user/{username}
	getResp, err := apiCtx.Get(ctx, "/user/"+username)
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, getResp.Status())
	getBody, err := getResp.Body(ctx)
	must.NoError(err)
	retrieved := mustUnmarshalUser(t, getBody)
	is.Equal(username, retrieved.Username)

	// PUT /user/{username}
	retrieved.Email = "updated-" + username + "@example.com"
	putResp, err := apiCtx.Put(ctx, "/user/"+username, &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, retrieved),
	})
	must.NoError(err)
	defer putResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, putResp.Status())

	// DELETE /user/{username}
	delResp, err := apiCtx.Delete(ctx, "/user/"+username)
	must.NoError(err)
	defer delResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, delResp.Status())

	// GET /user/{username} — must be 404
	confirmResp, err := apiCtx.Get(ctx, "/user/"+username)
	must.NoError(err)
	defer confirmResp.Dispose(ctx) //nolint:errcheck
	is.Equal(404, confirmResp.Status())
}
