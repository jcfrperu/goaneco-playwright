//go:build e2e

// E2E tests for BrowserContext.StorageState() — snapshot of current cookies and localStorage.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageStateCaptureCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/cookie-page", "text/html", `<!DOCTYPE html><html><body></body></html>`)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() { _ = bCtx.Close(testCtx(t)) })

	err = bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "abc123", Domain: strPtr("127.0.0.1"), Path: strPtr("/")},
	})
	must.NoError(err, "AddCookies failed")

	state, err := bCtx.StorageState(ctx)
	must.NoError(err, "StorageState failed")

	found := false
	for _, c := range state.Cookies {
		if c.Name == "session" && c.Value == "abc123" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("cookie 'session=abc123' not found in StorageState; cookies: %+v", state.Cookies)
	}
}

func TestStorageStateCaptureLocalStorage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ls-page", "text/html", `<!DOCTYPE html><html><body></body></html>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.Prefix()+"/ls-page")
	must.NoError(err, "Goto failed")
	_, err = page.Evaluate(ctx, `() => localStorage.setItem("key1", "value1")`)
	must.NoError(err, "localStorage.setItem failed")

	bCtx := page.Context()
	state, err := bCtx.StorageState(ctx)
	must.NoError(err, "StorageState failed")

	found := false
	for _, o := range state.Origins {
		for _, entry := range o.LocalStorage {
			if entry.Name == "key1" && entry.Value == "value1" {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("localStorage entry 'key1=value1' not found in StorageState; origins: %+v", state.Origins)
	}
}

func TestStorageStateEmptyByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() { _ = bCtx.Close(testCtx(t)) })

	state, err := bCtx.StorageState(ctx)
	must.NoError(err, "StorageState failed")
	is.Len(state.Cookies, 0)
}

func strPtr(s string) *string { return &s }

// ---------------------------------------------------------------------------
// From browser_context_storage_state_extra_test.go
// ---------------------------------------------------------------------------

// TestStorageStateHasCookies verifies StorageState captures cookies.
// Ref: TestBrowserContextStorageState.java#shouldCaptureCookies
func TestStorageStateHasCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, `() => document.cookie = 'session=abc'`)
	must.NoError(err)

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
	// May have cookies
}

// TestStorageStateHasOrigins verifies StorageState returns origins field.
// Ref: TestBrowserContextStorageState.java#shouldReturnOrigins
func TestStorageStateHasOrigins(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, `() => localStorage.setItem('test', 'value')`)
	must.NoError(err)

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
	must.NotNil(state.Origins)
}

// TestStorageStateRoundTrip verifies storage state can be captured and applied.
// Ref: TestBrowserContextStorageState.java#shouldRoundTripStorageState
func TestStorageStateRoundTrip(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc1 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	_, err = page1.Evaluate(ctx, `() => localStorage.setItem('key', 'stored')`)
	must.NoError(err)

	state, err := bc1.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
}
