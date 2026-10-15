//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func localStatusRB(n int) *int { return &n }

func localStrRB(s string) *string { return &s }

func localStatusPR7(n int) *int { return &n }

func localStrPR7(s string) *string { return &s }
