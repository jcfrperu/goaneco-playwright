//go:build e2e

// E2E tests for BrowserContext.Debugger (step-through debugging interface).
// Migration of: TestDebugger.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDebuggerNoPausedDetailsInitially verifies that PausedDetails returns nil before pausing.
// Ref: TestDebugger.java#shouldReturnNullPausedDetailsInitially
func TestDebuggerNoPausedDetailsInitially(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContextWithCleanup(t)

	dbg := bCtx.Debugger()
	if dbg == nil {
		t.Skip("Debugger not available in this browser context (server did not include debugger GUID)")
		return
	}

	details, err := dbg.PausedDetails(ctx)
	must.NoError(err, "PausedDetails() failed")
	must.Nil(details, "PausedDetails should be nil before pausing")
}

// TestDebuggerPauseAtNextAndResume verifies that RequestPause() pauses at the next action and Resume() continues.
// Ref: TestDebugger.java#shouldPauseAtNextAndResume
func TestDebuggerPauseAtNextAndResume(t *testing.T) {
	t.Skip("Requires interactive debugger flow that is not reliably automatable in e2e; covered by Playwright core tests")
}

// TestDebuggerStepWithNext verifies that Next() advances one step when paused.
// Ref: TestDebugger.java#shouldStepWithNext
func TestDebuggerStepWithNext(t *testing.T) {
	t.Skip("Requires interactive debugger flow that is not reliably automatable in e2e; covered by Playwright core tests")
}

// TestDebuggerPauseAtPauseCall verifies that a page.Pause() call triggers the debugger.
// Ref: TestDebugger.java#shouldPauseAtPauseCall
func TestDebuggerPauseAtPauseCall(t *testing.T) {
	t.Skip("Requires interactive debugger flow that is not reliably automatable in e2e; covered by Playwright core tests")
}
