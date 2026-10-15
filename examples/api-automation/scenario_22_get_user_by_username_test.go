//go:build e2e

package apiautomation

// Feature: Get a user by username
// Scenario: Retrieve an existing user by username and verify fields
//
// Feature file: features/22_get_user_by_username.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario22GetUserByUsername(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	user := createTestUser(t, "getuser22")

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/user/"+user.Username)
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	retrieved := mustUnmarshalUser(t, body)
	is.Equal(user.Username, retrieved.Username)
}
