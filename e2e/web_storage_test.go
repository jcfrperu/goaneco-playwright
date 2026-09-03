//go:build e2e

// E2E tests for web storage (localStorage and sessionStorage) via page.Evaluate.
// Ref: TestWebStorage.java
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocalStorageEmptyOnFreshOrigin verifies that localStorage has no items on a fresh origin.
// Ref: TestWebStorage.java#localStorageItemsReturnsEmptyListOnFreshOrigin
func TestLocalStorageEmptyOnFreshOrigin(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	length, err := page.Evaluate(ctx, "() => localStorage.length")
	must.NoError(err)
	is.Equal(float64(0), length)
}

// TestLocalStorageGetItemReturnsNullForMissingKey verifies getItem returns null for an absent key.
// Ref: TestWebStorage.java#localStorageGetItemReturnsNullForMissingKey
func TestLocalStorageGetItemReturnsNullForMissingKey(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	val, err := page.Evaluate(ctx, "() => localStorage.getItem('absent')")
	must.NoError(err)
	is.Nil(val)
}

// TestLocalStorageSetItemPersistsAndSurfaces verifies setItem persists values accessible via
// getItem and is also visible to JS running in the page.
// Ref: TestWebStorage.java#localStorageSetItemPersistsAndSurfacesInItemsAndGetItem
func TestLocalStorageSetItemPersistsAndSurfaces(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	is := assert.New(t)
	must := require.New(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, "() => localStorage.setItem('alpha', '1')")
	must.NoError(err)
	_, err = page.Evaluate(ctx, "() => localStorage.setItem('beta', '2')")
	must.NoError(err)

	alpha, err := page.Evaluate(ctx, "() => localStorage.getItem('alpha')")
	must.NoError(err)
	is.Equal("1", alpha)

	beta, err := page.Evaluate(ctx, "() => localStorage.getItem('beta')")
	must.NoError(err)
	is.Equal("2", beta)

	// Verify length reflects both items
	length, err := page.Evaluate(ctx, "() => localStorage.length")
	must.NoError(err)
	is.Equal(float64(2), length)
}

// TestLocalStorageSetItemOverwritesExistingValue verifies that setItem replaces an existing value.
// Ref: TestWebStorage.java#localStorageSetItemOverwritesExistingValue
func TestLocalStorageSetItemOverwritesExistingValue(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must := require.New(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, "() => localStorage.setItem('k', 'first')")
	must.NoError(err)
	_, err = page.Evaluate(ctx, "() => localStorage.setItem('k', 'second')")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => localStorage.getItem('k')")
	must.NoError(err)
	is.Equal("second", val)
}

// TestLocalStorageRemoveItemRemovesSingleItem verifies removeItem deletes only the target key.
// Ref: TestWebStorage.java#localStorageRemoveItemRemovesSingleItem
func TestLocalStorageRemoveItemRemovesSingleItem(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	is := assert.New(t)
	must := require.New(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, "() => localStorage.setItem('a', '1')")
	must.NoError(err)
	_, err = page.Evaluate(ctx, "() => localStorage.setItem('b', '2')")
	must.NoError(err)

	_, err = page.Evaluate(ctx, "() => localStorage.removeItem('a')")
	must.NoError(err)

	removed, err := page.Evaluate(ctx, "() => localStorage.getItem('a')")
	must.NoError(err)
	is.Nil(removed)

	remaining, err := page.Evaluate(ctx, "() => localStorage.getItem('b')")
	must.NoError(err)
	is.Equal("2", remaining)

	length, err := page.Evaluate(ctx, "() => localStorage.length")
	must.NoError(err)
	is.Equal(float64(1), length)
}

// TestLocalStorageClearEmptiesStorage verifies clear() removes all localStorage entries.
// Ref: TestWebStorage.java#localStorageClearEmptiesStorage
func TestLocalStorageClearEmptiesStorage(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must := require.New(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := page.Evaluate(ctx, "() => localStorage.setItem('a', '1')")
	must.NoError(err)
	_, err = page.Evaluate(ctx, "() => localStorage.setItem('b', '2')")
	must.NoError(err)

	_, err = page.Evaluate(ctx, "() => localStorage.clear()")
	must.NoError(err)

	length, err := page.Evaluate(ctx, "() => localStorage.length")
	must.NoError(err)
	is.Equal(float64(0), length)
}

// TestSessionStorageRoundTrip verifies the full sessionStorage lifecycle:
// empty initially, set/get items, removeItem, clear.
// Ref: TestWebStorage.java#sessionStorageRoundTrip
func TestSessionStorageRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	is := assert.New(t)
	must := require.New(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	// Initially empty
	length, err := page.Evaluate(ctx, "() => sessionStorage.length")
	must.NoError(err)
	is.Equal(float64(0), length)

	// Set two items
	_, err = page.Evaluate(ctx, "() => sessionStorage.setItem('s1', 'v1')")
	must.NoError(err)
	_, err = page.Evaluate(ctx, "() => sessionStorage.setItem('s2', 'v2')")
	must.NoError(err)

	v1, err := page.Evaluate(ctx, "() => sessionStorage.getItem('s1')")
	must.NoError(err)
	is.Equal("v1", v1)

	v2, err := page.Evaluate(ctx, "() => sessionStorage.getItem('s2')")
	must.NoError(err)
	is.Equal("v2", v2)

	// Remove one item
	_, err = page.Evaluate(ctx, "() => sessionStorage.removeItem('s1')")
	must.NoError(err)

	removed, err := page.Evaluate(ctx, "() => sessionStorage.getItem('s1')")
	must.NoError(err)
	is.Nil(removed)

	still, err := page.Evaluate(ctx, "() => sessionStorage.getItem('s2')")
	must.NoError(err)
	is.Equal("v2", still)

	// Clear all
	_, err = page.Evaluate(ctx, "() => sessionStorage.clear()")
	must.NoError(err)

	finalLen, err := page.Evaluate(ctx, "() => sessionStorage.length")
	must.NoError(err)
	is.Equal(float64(0), finalLen)
}

// TestLocalStorageAndSessionStorageAreIndependent verifies that localStorage and sessionStorage
// hold independent values even when the same key is used, and clearing one does not affect the other.
// Ref: TestWebStorage.java#localStorageAndSessionStorageAreIndependent
func TestLocalStorageAndSessionStorageAreIndependent(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	is := assert.New(t)
	must := require.New(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, "() => localStorage.setItem('shared', 'local')")
	must.NoError(err)
	_, err = page.Evaluate(ctx, "() => sessionStorage.setItem('shared', 'session')")
	must.NoError(err)

	localVal, err := page.Evaluate(ctx, "() => localStorage.getItem('shared')")
	must.NoError(err)
	is.Equal("local", localVal)

	sessionVal, err := page.Evaluate(ctx, "() => sessionStorage.getItem('shared')")
	must.NoError(err)
	is.Equal("session", sessionVal)

	// Clear localStorage; sessionStorage must be unaffected
	_, err = page.Evaluate(ctx, "() => localStorage.clear()")
	must.NoError(err)

	localLen, err := page.Evaluate(ctx, "() => localStorage.length")
	must.NoError(err)
	is.Equal(float64(0), localLen)

	sessionStill, err := page.Evaluate(ctx, "() => sessionStorage.getItem('shared')")
	must.NoError(err)
	is.Equal("session", sessionStill)
}

// TestStorageScopedToCurrentOrigin verifies that localStorage is origin-scoped:
// a value set on one origin is not visible on a different origin, and each origin
// maintains its own independent storage.
// Ref: TestWebStorage.java#storageMethodsAreScopedToTheCurrentOrigin
func TestStorageScopedToCurrentOrigin(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	is := assert.New(t)
	must := require.New(t)

	// Navigate to origin-1, set a value
	must.NoError(page.Goto(ctx, srv.Prefix()+"/empty.html"))
	_, err := page.Evaluate(ctx, "() => localStorage.setItem('k', 'origin-1')")
	must.NoError(err)

	// Navigate to cross-process origin; storage should be empty there
	must.NoError(page.Goto(ctx, srv.CrossProcessPrefix()+"/empty.html"))
	crossLen, err := page.Evaluate(ctx, "() => localStorage.length")
	must.NoError(err)
	is.Equal(float64(0), crossLen)

	// Set a different value on the cross-process origin
	_, err = page.Evaluate(ctx, "() => localStorage.setItem('k', 'origin-2')")
	must.NoError(err)

	// Navigate back to origin-1; original value must still be present
	must.NoError(page.Goto(ctx, srv.Prefix()+"/empty.html"))
	val, err := page.Evaluate(ctx, "() => localStorage.getItem('k')")
	must.NoError(err)
	is.Equal("origin-1", val)
}

// TestLocalStorageNotSharedBetweenContexts verifies that localStorage is isolated between
// different BrowserContext instances even when navigating to the same origin.
func TestLocalStorageNotSharedBetweenContexts(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	is := assert.New(t)
	must := require.New(t)

	bCtx1 := newContext(t)
	page1, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	_, err = page1.Evaluate(ctx, "() => localStorage.setItem('ctx', 'ctx1-value')")
	must.NoError(err)

	bCtx2 := newContext(t)
	page2, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	// Context 2 must not see context 1's localStorage entry
	val, err := page2.Evaluate(ctx, "() => localStorage.getItem('ctx')")
	must.NoError(err)
	is.Nil(val, "localStorage must be isolated between browser contexts")
}

// TestLocalStorageSnapshotViaStorageState verifies that localStorage entries set in a page are
// captured by BrowserContext.StorageState and reflected in the returned Origins list.
func TestLocalStorageSnapshotViaStorageState(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	is := assert.New(t)
	must := require.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => localStorage.setItem('snap-key', 'snap-value')")
	must.NoError(err)

	state, err := bCtx.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)

	found := false
	for _, origin := range state.Origins {
		for _, entry := range origin.LocalStorage {
			if entry.Name == "snap-key" && entry.Value == "snap-value" {
				found = true
			}
		}
	}
	is.True(found, "StorageState should capture the localStorage entry 'snap-key=snap-value'; got: %+v", state.Origins)
}

// TestLocalStorageRestoreViaStorageState verifies that a BrowserContext created with a pre-populated
// StorageState exposes the restored localStorage values to pages navigating to that origin.
func TestLocalStorageRestoreViaStorageState(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	is := assert.New(t)
	must := require.New(t)

	// Step 1: capture a storage state containing a localStorage entry.
	bCtx1 := newContext(t)
	page1, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	_, err = page1.Evaluate(ctx, "() => localStorage.setItem('restored-key', 'restored-value')")
	must.NoError(err)

	state, err := bCtx1.StorageState(ctx)
	must.NoError(err)
	must.NotEmpty(state.Origins, "expected at least one origin in StorageState after setItem")

	// Step 2: create a new context seeded with the captured state and verify restoration.
	bCtx2, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageState: state,
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bCtx2.Close(testCtx(t)) })

	page2, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	val, err := page2.Evaluate(ctx, "() => localStorage.getItem('restored-key')")
	must.NoError(err)
	is.Equal("restored-value", val, "localStorage should be restored from StorageState")
}
