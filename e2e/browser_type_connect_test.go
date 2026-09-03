//go:build e2e

// Tests for BrowserType.ConnectOverCDP.
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConnectOverCDPInvalidEndpoint verifies that ConnectOverCDP returns a meaningful
// error when the endpoint URL is unreachable or invalid.
// Ref: TestBrowserTypeConnect.java (negative endpoint test)
func TestConnectOverCDPInvalidEndpoint(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	bt := globalBrowser.BrowserType()

	_, err := bt.ConnectOverCDP(ctx, "http://localhost:19999/nonexistent")
	is.Error(err)
}
