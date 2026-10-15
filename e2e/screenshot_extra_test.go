//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
