//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func localStrPtrCk4(s string) *string { return &s }
