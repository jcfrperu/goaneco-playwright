//go:build e2e

package apiautomation

// Feature: Login with invalid credentials
// Scenario: Login with wrong credentials returns a non-200 status
//
// Feature file: features/29_login_invalid.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario29LoginInvalid(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/user/login?username=INVALID_USER_XYZ&password=WRONG_PASS_XYZ")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.False(resp.OK(), "login with invalid credentials should not return 2xx, got %d", resp.Status())
}
