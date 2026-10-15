//go:build e2e

// E2E tests for: TestTracingOverConnect.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConnectOverCDPInvalidEndpoint verifies that ConnectOverCDP returns a meaningful
// error when the endpoint URL is unreachable or invalid.
// Ref: TestBrowserTypeConnect.java (negative endpoint test)
// Extra: no equivalent method in playwright-java; related class: TestTracingOverConnect.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracingOverConnect.java
func TestConnectOverCDPInvalidEndpoint(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	bt := globalBrowser.BrowserType()

	_, err := bt.ConnectOverCDP(ctx, "http://localhost:19999/nonexistent")
	is.Error(err)
}

// Ref: TestTracingOverConnect.java#shouldRecordSourcesForBrowserFromConnect  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTracingOverConnect.java
func TestTracingOverConnectShouldRecordSourcesForBrowserFromConnect(t *testing.T) {
	t.Skip("not implemented: BrowserType.Connect")
}
