//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"net/http"
	"testing"
)

func setupServiceWorkerRoutes(t *testing.T, srv *testserver.Server) {
	t.Helper()
	srv.SetRoute("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.addEventListener('install', () => self.skipWaiting());`))
	})
	srv.ServeWithBody("/sw-page.html", "text/html", `<!DOCTYPE html><html><body>
<script>
  window.registrationPromise = navigator.serviceWorker
    ? navigator.serviceWorker.register('/sw.js')
    : null;
</script>
</body></html>`)
}
