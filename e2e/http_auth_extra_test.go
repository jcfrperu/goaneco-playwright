//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	"encoding/base64"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func newPageWithCredentials(t *testing.T, username, password string) *playwright.Page {
	t.Helper()
	ctx := testCtx(t)
	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: username,
			Password: password,
		},
	})
	require.NoError(t, err, "NewContext with credentials failed")
	page, err := bc.NewPage(ctx)
	require.NoError(t, err, "NewPage failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bc.Close(closeCtx)
	})
	return page
}

func serveAuthProtected(t *testing.T, srv *testserver.Server, path, user, pass, body string) {
	t.Helper()
	srv.SetRoute(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expected := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
		if authHeader != expected {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(body))
	}))
}
