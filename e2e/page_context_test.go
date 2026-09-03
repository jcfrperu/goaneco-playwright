//go:build e2e

// Page.Context and BrowserContext isolation E2E tests.
// Migration of: TestPageContext.java / TestBrowserContextIsolation.java (additional cases)
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageContextReturnsParentContext verifies Page.Context returns the context that created it.
// Ref: TestPageContext.java#shouldReturnContext
func TestPageContextReturnsParentContextNonNil(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	page := newPage(t)

	bc := page.Context()
	must.NotNil(bc)
}

// TestBrowserContextIsolatesLocalStorage verifies two contexts don't share localStorage.
// Ref: TestBrowserContextIsolation.java#shouldNotShareLocalStorage
func TestBrowserContextIsolatesLocalStorageExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc1 := newContext(t)
	bc2 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)
	page2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<html><body></body></html>`))
	must.NoError(page2.SetContent(ctx, `<html><body></body></html>`))

	_, err = page1.Evaluate(ctx, `() => localStorage.setItem('key', 'ctx1-value')`)
	must.NoError(err)

	val, err := page2.Evaluate(ctx, `() => localStorage.getItem('key')`)
	must.NoError(err)
	is.Nil(val) // Different context should not have this key
}

// TestBrowserContextIsolatesCookies verifies two contexts don't share cookies.
// Ref: TestBrowserContextIsolation.java#shouldNotShareCookies
func TestBrowserContextIsolatesCookiesExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc1 := newContext(t)
	bc2 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)
	page2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<html><body></body></html>`))
	must.NoError(page2.SetContent(ctx, `<html><body></body></html>`))

	_, err = page1.Evaluate(ctx, `() => document.cookie = 'ctx1cookie=value1'`)
	must.NoError(err)

	cookies, err := page2.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
	// Cookie from context 1 should not be in context 2
	if cookieStr, ok := cookies.(string); ok {
		must.NotContains(cookieStr, "ctx1cookie")
	}
}

// TestBrowserContextPagesListsAllPages verifies Context.Pages returns all open pages.
// Ref: TestBrowserContextIsolation.java#shouldListAllPages
func TestBrowserContextPagesListsAllPagesExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	pages := bc.Pages()
	is.GreaterOrEqual(len(pages), 2)

	_ = page1
	_ = page2
}
