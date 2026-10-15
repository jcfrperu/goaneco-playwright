//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func localIntPtr(n int) *int { return &n }
