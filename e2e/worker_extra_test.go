//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
)

// waitForWorker registers an OnWorker handler and returns a channel that receives the first Worker.
func waitForWorker(page *playwright.Page) <-chan *playwright.Worker {
	ch := make(chan *playwright.Worker, 1)
	page.OnWorker(func(w *playwright.Worker) {
		select {
		case ch <- w:
		default:
		}
	})
	return ch
}
