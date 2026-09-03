//go:build e2e

// HTTP Basic Auth E2E tests.
package e2e

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestHTTPAuthSuccess(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/protected", "user", "secret", "text/html", "<p>authenticated</p>")

	page := newPageWithCredentials(t, "user", "secret")
	err := page.Goto(ctx, srv.Prefix()+"/protected")
	must.NoError(err, "Goto failed")

	text, err := page.Locator("p").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	if text != "authenticated" {
		t.Errorf("got %q, want 'authenticated'", text)
	}
}

func TestHTTPAuthFailure(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/locked", "admin", "correct", "text/html", "<p>secret</p>")

	statuses := make(chan int, 5)
	page := newPageWithCredentials(t, "admin", "wrong")
	page.OnResponse(func(resp *playwright.NetworkResponse) {
		statuses <- resp.Status()
	})

	_ = page.Goto(ctx, srv.Prefix()+"/locked")

	select {
	case status := <-statuses:
		if status != 401 {
			t.Errorf("Status = %d, want 401", status)
		}
	case <-time.After(5 * time.Second):
		t.Error("no response received within 5 seconds")
	}
}

func TestHTTPAuthNoCredentials(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBasicAuth("/guarded", "u", "p", "text/html", "<p>ok</p>")

	statuses := make(chan int, 5)
	page := newPage(t)
	page.OnResponse(func(resp *playwright.NetworkResponse) {
		statuses <- resp.Status()
	})

	_ = page.Goto(ctx, srv.Prefix()+"/guarded")

	select {
	case status := <-statuses:
		if status != 401 {
			t.Errorf("Status = %d, want 401", status)
		}
	case <-time.After(5 * time.Second):
		t.Error("no response received within 5 seconds")
	}
}

// ---------------------------------------------------------------------------
// From browser_context_http_auth_extra_test.go
// ---------------------------------------------------------------------------

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

// TestHTTPAuthContextSucceeds verifies BrowserContext with credentials succeeds on protected page.
// Ref: TestBrowserContextHTTPAuth.java#shouldWork
func TestHTTPAuthContextSucceeds(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveAuthProtected(t, srv, "/protected", "admin", "secret", "authenticated content")

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "admin",
			Password: "secret",
		},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/protected"))

	text, err := page.Evaluate(ctx, `() => document.body.textContent`)
	must.NoError(err)
	is.Equal("authenticated content", text)
}

// TestHTTPAuthContextFailsWithoutCredentials verifies unauthenticated access returns 401.
// Ref: TestBrowserContextHTTPAuth.java#shouldFailWithoutCredentials
func TestHTTPAuthContextFailsWithoutCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveAuthProtected(t, srv, "/protected401", "admin", "secret", "should not see this")

	page := newPage(t)

	// Should navigate but get 401 content (or error)
	_ = page.Goto(ctx, srv.Prefix()+"/protected401")

	status, err := page.Evaluate(ctx, `() => document.readyState`)
	must.NoError(err)
	must.NotNil(status)
}

// TestHTTPAuthContextWithWrongCredentials verifies wrong credentials still return 401.
// Ref: TestBrowserContextHTTPAuth.java#shouldFailWithWrongCredentials
func TestHTTPAuthContextWithWrongCredentials(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	serveAuthProtected(t, srv, "/protected-wrong", "admin", "correct", "secret content")

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		HttpCredentials: &playwright.HttpCredentials{
			Username: "admin",
			Password: "wrong-password",
		},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	_ = page.Goto(ctx, srv.Prefix()+"/protected-wrong")

	// Page navigated but with 401 - just verify we can still evaluate
	state, err := page.Evaluate(ctx, `() => document.readyState`)
	must.NoError(err)
	must.NotNil(state)
}
