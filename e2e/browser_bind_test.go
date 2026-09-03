//go:build e2e

// E2E tests for Browser.Bind / Browser.Unbind (expose browser via a named server endpoint).
// Migration of: TestBrowserBind.java
//
// NOTE: Browser.Bind() and Browser.Unbind() are not yet implemented in the Go API.
// Tests are skipped until the feature is available.
package e2e

import "testing"

// TestBrowserBindAndUnbind verifies that a browser can be bound to an endpoint and unbound.
// Ref: TestBrowserBind.java#shouldBindAndUnbindBrowser
func TestBrowserBindAndUnbind(t *testing.T) {
	t.Skip("Browser.Bind() / Browser.Unbind() not yet implemented in Go API")
}

// TestBrowserBindWithCustomTitleAndOptions verifies bind with custom title and host/port options.
// Ref: TestBrowserBind.java#shouldBindWithCustomTitleAndOptions
func TestBrowserBindWithCustomTitleAndOptions(t *testing.T) {
	t.Skip("Browser.Bind() / Browser.Unbind() not yet implemented in Go API")
}
