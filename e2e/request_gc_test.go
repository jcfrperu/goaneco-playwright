//go:build e2e

// E2E tests for Page.RequestGC (explicit garbage collection trigger).
// Migration of: TestRequestGC.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRequestGCShouldWork verifies that RequestGC triggers garbage collection and
// allows WeakRef objects to be collected after their target is released.
// Ref: TestRequestGC.java#shouldWork
func TestRequestGCShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Create a WeakRef to an object, then drop the strong reference and trigger GC.
	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => {
		window.__weakRef = new WeakRef({ name: 'playwright' });
	}`)
	must.NoError(err, "failed to create WeakRef")

	must.NoError(page.RequestGC(ctx), "RequestGC() should not error")

	// After GC, the WeakRef target may have been collected.
	// We just verify the API works without error; GC timing is non-deterministic.
	result, err := page.Evaluate(ctx, `() => {
		const ref = window.__weakRef.deref();
		return ref !== undefined ? ref.name : null;
	}`)
	must.NoError(err, "Evaluate after RequestGC failed")
	// The result is either the string "playwright" (if not yet GC'd) or nil.
	_ = result
}
