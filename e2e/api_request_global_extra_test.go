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

// newAPICtx creates a standalone APIRequestContext and registers cleanup.
func newAPICtx(t *testing.T, opts ...*playwright.APIRequestContextOptions) *playwright.APIRequestContext {
	t.Helper()
	ctx := testCtx(t)
	apiCtx, err := globalPW.NewAPIRequestContext(ctx, opts...)
	require.NoError(t, err, "NewAPIRequestContext")
	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = apiCtx.Dispose(cleanCtx)
	})
	return apiCtx
}

// ── TestGlobalFetch equivalents ──────────────────────────────────────────────

const simpleJSONBody = "{\"foo\": \"bar\"}\n"
