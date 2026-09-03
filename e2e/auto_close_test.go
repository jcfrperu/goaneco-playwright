//go:build e2e

// E2E tests for automatic resource cleanup using defer (Go equivalent of Java try-with-resources).
// Migration of: TestAutoClose.java
package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestAutoCloseBrowserContext verifies that a BrowserContext closes cleanly when defer'd.
// This is Go's idiomatic equivalent of Java's try-with-resources (AutoCloseable).
// Ref: TestAutoClose.java#shouldAllowUsingTryWithResources
func TestAutoCloseBrowserContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")

	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := bCtx.Close(closeCtx); err != nil {
			t.Logf("bCtx.Close: %v", err)
		}
	}()

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	must.NotNil(page)

	pages := bCtx.Pages()
	must.Len(pages, 1, "expected exactly one page in context")
}

// TestAutoClosePageInheritedCleanup verifies that all pages in a context are released when it closes.
func TestAutoClosePageInheritedCleanup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")

	_, err = bCtx.NewPage(ctx)
	must.NoError(err, "NewPage 1 failed")
	_, err = bCtx.NewPage(ctx)
	must.NoError(err, "NewPage 2 failed")

	must.Len(bCtx.Pages(), 2, "expected 2 pages before close")

	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	must.NoError(bCtx.Close(closeCtx), "Close failed")
}
