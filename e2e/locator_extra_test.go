//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func localStrPtrLF3(s string) *string { return &s }

func localStrPtrLF4(s string) *string { return &s }
