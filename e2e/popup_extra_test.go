//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// waitForPopup subscribes to page.OnPopup, runs trigger(), and returns the
// first popup page received within the test deadline.
func waitForPopup(t *testing.T, page *playwright.Page, trigger func()) *playwright.Page {
	t.Helper()
	popupCh := make(chan *playwright.Page, 1)
	cancel := page.OnPopup(func(p *playwright.Page) {
		select {
		case popupCh <- p:
		default:
		}
	})
	defer cancel()

	trigger()

	select {
	case popup := <-popupCh:
		return popup
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for popup")
		return nil
	}
}

// newContextWithCleanup creates a BrowserContext with the given options and
// registers cleanup. Equivalent to the per-test context pattern used across
// browser_context_test.go.
func newContextWithCleanup(t *testing.T, opts ...*playwright.BrowserContextOptions) *playwright.BrowserContext {
	t.Helper()
	ctx := testCtx(t)
	bc, err := globalBrowser.NewContext(ctx, opts...)
	require.NoError(t, err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(closeCtx)
	})
	return bc
}

// --- TestPageEventPopup tests ---
