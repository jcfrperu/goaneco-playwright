//go:build e2e

// Page.WaitForTimeout E2E tests.
// Migration of: TestPageWaitForTimeout.java
package e2e

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageWaitForTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	tests := []struct {
		name  string
		ms    float64
		minMs int64
	}{
		{"zero returns immediately", 0, 0},
		{"short duration 50ms", 50, 30},
		{"standard delay 100ms", 100, 80},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			page := newPage(t)
			start := time.Now()
			must.NoError(page.WaitForTimeout(testCtx(t), tc.ms))
			is.Less(time.Since(start).Milliseconds(), int64(5000), "timed out waiting")
			is.GreaterOrEqual(time.Since(start).Milliseconds(), tc.minMs)
		})
	}
}
