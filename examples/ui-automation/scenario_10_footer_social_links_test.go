//go:build e2e

package uiautomation

// Feature: Footer social media links
// Scenario: The inventory page footer contains visible links to Twitter, Facebook,
// and LinkedIn.
//
// Feature file: features/10_footer_social_links.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFooterSocialLinks(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	page := newPage(t)
	loginStandardUser(t, page)

	ctx := testCtx(t)

	socialLinks := map[string]string{
		"Twitter":  ".social_twitter a",
		"Facebook": ".social_facebook a",
		"LinkedIn": ".social_linkedin a",
	}

	for name, selector := range socialLinks {
		loc := page.Locator(selector)
		visible, err := loc.IsVisible(ctx)
		must.NoError(err, "check %s link visibility", name)
		is.True(visible, "%s link should be visible in footer", name)

		href, err := loc.GetAttribute(ctx, "href")
		must.NoError(err, "get %s link href", name)
		is.NotEmpty(href, "%s link href must not be empty", name)
	}
}
