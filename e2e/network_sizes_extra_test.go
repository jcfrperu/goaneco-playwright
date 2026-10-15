//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"strings"
	"sync"
	"testing"
	"time"
)

// captureFinishedRequestTo waits for the first completed request to a URL containing urlSubstr.
func captureFinishedRequestTo(t *testing.T, page *playwright.Page, urlSubstr string, trigger func()) *playwright.NetworkRequest {
	t.Helper()
	ch := make(chan *playwright.NetworkRequest, 1)
	var once sync.Once
	cancel := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if strings.Contains(req.URL(), urlSubstr) {
			once.Do(func() {
				select {
				case ch <- req:
				default:
				}
			})
		}
	})
	defer cancel()

	trigger()

	select {
	case req := <-ch:
		return req
	case <-time.After(10 * time.Second):
		t.Fatalf("captureFinishedRequestTo: timed out waiting for request to URL containing %q", urlSubstr)
		return nil
	}
}
