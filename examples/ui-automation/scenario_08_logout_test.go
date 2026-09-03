//go:build e2e

package uiautomation

// Feature: Logout via burger menu
// Scenario: Open the side navigation menu and click Logout; the browser should
// redirect back to the login page with the login form visible.
//
// Feature file: features/08_logout.feature

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutViaBurgerMenu(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)
	inv := pom.NewInventoryPage(page)

	must.NoError(inv.Logout(ctx), "logout via burger menu")

	lp := pom.NewLoginPage(page)
	visible, err := lp.IsVisible(ctx)
	must.NoError(err, "check login form visibility after logout")
	is.True(visible, "login form should be visible after logout")
}
