//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

func localBoolPtrGB6(b bool) *bool { return &b }

func localStringPtrGB6(s string) *string { return &s }

func localBoolPtrGB7(b bool) *bool { return &b }

func localStringPtrGB7(s string) *string { return &s }

func localBoolPtrGB8(b bool) *bool { return &b }

func localStringPtrGB8(s string) *string { return &s }

func localStringPtrGB9(s string) *string { return &s }
