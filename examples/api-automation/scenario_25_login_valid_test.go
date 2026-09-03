//go:build e2e

package apiautomation

// Feature: Login with valid credentials
// Scenario: Login returns 200 and a session token in the response body
//
// Feature file: features/25_login_valid.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario25LoginValid(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/user/login?username=test&password=abc123")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)
	is.NotEmpty(body, "login response body should not be empty")
}
