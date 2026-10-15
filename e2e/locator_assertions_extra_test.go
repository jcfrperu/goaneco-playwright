//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	"testing"
	"time"
)

// shortCtx returns a context with a 2-second deadline, useful for fail tests.
func shortCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(testCtx(t), 2*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func localStrPtrHT3(s string) *string { return &s }
