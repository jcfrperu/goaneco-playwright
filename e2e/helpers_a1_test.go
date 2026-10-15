//go:build e2e

// Shared helpers for ported playwright-java tests (agent 1).
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"net/http"
	"path/filepath"
	"testing"
)

// a1ResPath returns the absolute path of a playwright-java test resource file.
func a1ResPath(t testing.TB, name string) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "playwright-reference-projects", "playwright-java", "playwright", "src", "test", "resources", name))
	if err != nil {
		t.Fatalf("a1ResPath(%q): %v", name, err)
	}
	return p
}

// a1ServeProxyTarget makes the test server answer /target.html like an HTTP proxy would.
func a1ServeProxyTarget(srv *testserver.Server) {
	srv.SetRoute("/target.html", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html><title>Served by the proxy</title></html>"))
	})
}
