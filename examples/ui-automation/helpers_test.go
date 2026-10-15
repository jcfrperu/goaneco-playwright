//go:build e2e

package uiautomation

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/require"
)

// loginStandardUser navigates to the SauceDemo login page and authenticates
// with the standard test credentials, leaving the page on the inventory URL.
func loginStandardUser(t *testing.T, page *playwright.Page) {
	t.Helper()
	ctx := testCtx(t)
	lp := pom.NewLoginPage(page)
	require.NoError(t, lp.Navigate(ctx), "navigate to login page")
	require.NoError(t, lp.Login(ctx, pom.StandardUsername, pom.StandardPassword), "login with standard_user")
}
