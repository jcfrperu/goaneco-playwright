//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
)

func aiMode() *playwright.LocatorAriaSnapshotOptions {
	m := "ai"
	return &playwright.LocatorAriaSnapshotOptions{Mode: &m}
}
