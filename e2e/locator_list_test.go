//go:build e2e

// E2E tests for Locator.All — iterating over multiple matched elements.
// Migration of: TestLocatorList.java
//
// NOTE: TestLocatorAllShouldWork (the direct Java port) already exists in
// locator_convenience_test.go. This file adds complementary coverage.
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorAllWithDivParagraphs verifies Locator.All on nested div > p selectors.
// Ref: TestLocatorList.java#locatorAllShouldWork (div >> p variant)
func TestLocatorAllWithDivParagraphs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><p>A</p><p>B</p><p>C</p></div>`))

	locators, err := page.Locator("div >> p").All(ctx)
	must.NoError(err, "All() failed")
	must.Len(locators, 3, "expected 3 <p> elements")

	var texts []string
	for _, loc := range locators {
		text, err := loc.TextContent(ctx)
		must.NoError(err, "TextContent failed")
		must.NotNil(text, "TextContent returned nil")
		texts = append(texts, *text)
	}

	is.Equal([]string{"A", "B", "C"}, texts)
}

// TestLocatorAllReturnsEmptyWhenNoMatch verifies Locator.All returns empty slice when nothing matches.
func TestLocatorAllReturnsEmptyWhenNoMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><p>Only paragraph</p></div>`))

	locators, err := page.Locator("span").All(ctx)
	must.NoError(err, "All() failed")
	is.Empty(locators, "expected empty slice when no <span> elements exist")
}

// TestLocatorAllSingleMatch verifies Locator.All with exactly one matching element.
func TestLocatorAllSingleMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul><li>only item</li></ul>`))

	locators, err := page.Locator("li").All(ctx)
	must.NoError(err, "All() failed")
	must.Len(locators, 1, "expected exactly one <li>")

	text, err := locators[0].TextContent(ctx)
	must.NoError(err, "TextContent failed")
	must.NotNil(text)
	is.Equal("only item", *text)
}
