//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func strPtr(s string) *string { return &s }

// ---------------------------------------------------------------------------
// From browser_context_storage_state_extra_test.go
// ---------------------------------------------------------------------------
