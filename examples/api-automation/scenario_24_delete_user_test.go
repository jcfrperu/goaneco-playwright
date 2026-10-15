//go:build e2e

package apiautomation

// Feature: Delete a user
// Scenario: Delete a user and verify the account is no longer retrievable
//
// Feature file: features/24_delete_user.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario24DeleteUser(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	user := createTestUser(t, "deleteuser24")

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	delResp, err := apiCtx.Delete(ctx, "/user/"+user.Username)
	must.NoError(err)
	defer delResp.Dispose(ctx) //nolint:errcheck
	is.Equal(200, delResp.Status())

	getResp, err := apiCtx.Get(ctx, "/user/"+user.Username)
	must.NoError(err)
	defer getResp.Dispose(ctx) //nolint:errcheck
	is.Equal(404, getResp.Status())
}
