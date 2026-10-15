//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
