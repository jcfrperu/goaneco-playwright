//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func roleBoolPtr(b bool) *bool { return &b }

func roleStrPtr(s string) *string { return &s }

func roleIntPtr(i int) *int { return &i }

func localBoolPtrST4(b bool) *bool { return &b }

func localBoolPtrST5(b bool) *bool { return &b }
