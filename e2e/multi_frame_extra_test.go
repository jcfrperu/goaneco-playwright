//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"testing"
	"time"
)

// waitForFrame polls page.Frame(name) until it appears or the deadline is reached.
func waitForFrame(t *testing.T, page *playwright.Page, name string) *playwright.Frame {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if f := page.Frame(name); f != nil {
			return f
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("frame %q never attached within 5 seconds", name)
	return nil
}
