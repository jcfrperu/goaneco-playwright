package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDriver_StartInvalidPath verifies that attempting to start a driver with an invalid
// path returns an error without hanging or leaking resources.
func TestDriver_StartInvalidPath(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	d, err := StartDriver(ctx, filepath.Join(t.TempDir(), "non-existent-cli.js"))
	if err != nil {
		is.Nil(d)
		return
	}

	// If driver started (e.g. node runs and fails immediately), verify Stop completes
	err = d.Stop()
	is.Error(err)
}

// TestDriver_ShouldNotHangWhenProcessExitsUnexpectedly (UNIT-DRV-01)
// Verifies that when the driver process exits unexpectedly (or is killed),
// the driver's watchdog reaps the process, closes channels, and Stop() completes
// within a bounded time without hanging indefinitely.
func TestDriver_ShouldNotHangWhenProcessExitsUnexpectedly(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	// Create a short-lived Node.js script that exits immediately with non-zero exit code
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "exit_early.js")
	err := os.WriteFile(scriptPath, []byte("process.exit(42);\n"), 0600)
	must.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	d, err := StartDriver(ctx, scriptPath)
	must.NoError(err)
	must.NotNil(d)

	// Wait for process to exit and verify Stop returns without hanging
	stopDone := make(chan error, 1)
	go func() {
		stopDone <- d.Stop()
	}()

	select {
	case stopErr := <-stopDone:
		is.Error(stopErr, "expected error from non-zero exit code 42")
	case <-time.After(3 * time.Second):
		t.Fatal("driver.Stop() hung after process exited unexpectedly")
	}

	// Verify that any new SendRequest on the driver's connection fails fast with context error or transport error
	reqCtx, reqCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	t.Cleanup(reqCancel)
	_, reqErr := d.Conn().SendRequest(reqCtx, 1, []byte(`{}`))
	is.Error(reqErr, "SendRequest on dead driver should return error")
}

// TestDriver_KillProcessWhileWaiting (UNIT-DRV-01)
// Verifies that killing the running process explicitly triggers the watchdog
// and unblocks any waiting callers within a bounded time.
func TestDriver_KillProcessWhileWaiting(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	// Create a Node.js script that stays alive until killed
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "sleep_script.js")
	err := os.WriteFile(scriptPath, []byte("setInterval(() => {}, 1000);\n"), 0600)
	must.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	d, err := StartDriver(ctx, scriptPath)
	must.NoError(err)
	must.NotNil(d)

	// Kill the process abruptly
	if d.cmd != nil && d.cmd.Process != nil {
		_ = d.cmd.Process.Kill()
	}

	stopDone := make(chan error, 1)
	go func() {
		stopDone <- d.Stop()
	}()

	select {
	case <-stopDone:
		// Stop completed cleanly after abrupt kill
	case <-time.After(3 * time.Second):
		t.Fatal("driver.Stop() hung after process was killed")
	}
}
