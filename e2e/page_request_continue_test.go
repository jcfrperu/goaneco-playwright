//go:build e2e

// E2E tests for: TestPageRequestContinue.java.
package e2e

import (
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageRequestContinueShouldDeleteHeaderWithUndefinedValue verifies removing a header while continuing a request.
// Ref: TestPageRequestContinue.java#shouldDeleteHeaderWithUndefinedValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestContinue.java
func TestPageRequestContinueShouldDeleteHeaderWithUndefinedValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	var mu sync.Mutex
	var serverHeaders http.Header
	srv.SetRoute("/something", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		serverHeaders = r.Header.Clone()
		mu.Unlock()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_, _ = w.Write([]byte("done"))
	})
	var interceptedHeaders map[string]string
	must.NoError(page.Route(ctx, srv.Prefix()+"/something", func(route *playwright.Route) {
		headers := map[string]string{}
		for k, v := range route.Request().Headers() {
			headers[k] = v
		}
		mu.Lock()
		interceptedHeaders = route.Request().Headers()
		mu.Unlock()
		delete(headers, "foo")
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{Headers: headers})
	}))
	text, err := page.Evaluate(ctx, `async url => {
      const data = await fetch(url, {
        headers: {
          foo: 'a',
          bar: 'b',
        }
      });
      return data.text();
    }`, srv.Prefix()+"/something")
	must.NoError(err)
	is.Equal("done", text)
	mu.Lock()
	defer mu.Unlock()
	_, hasFoo := interceptedHeaders["foo"]
	is.True(hasFoo, "intercepted request is expected to carry the original foo header before it is removed")
	is.Empty(serverHeaders.Values("foo"))
	is.Equal([]string{"b"}, serverHeaders.Values("bar"))
}

// TestPageRequestContinueShouldNotThrowWhenContinuingAfterPageIsClosed verifies Route.Continue after the page is closed does not fail.
// Ref: TestPageRequestContinue.java#shouldNotThrowWhenContinuingAfterPageIsClosed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestContinue.java
func TestPageRequestContinueShouldNotThrowWhenContinuingAfterPageIsClosed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	done := make(chan struct{})
	var once sync.Once
	var contErr error
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = page.Close(ctx)
		contErr = route.Continue(ctx, nil)
		once.Do(func() { close(done) })
	}))
	err := page.Goto(ctx, srv.EmptyPage())
	must.Error(err)
	is.True(strings.Contains(err.Error(), "Target page, context or browser has been closed") || strings.Contains(err.Error(), "frame was detached"), err.Error())
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("route handler did not complete")
	}
	is.NoError(contErr)
}

// TestPageRequestContinueRedirectAfterContinueShouldBeAbleToDeleteCookie verifies a cookie deleted by a redirected continued request is gone.
// Ref: TestPageRequestContinue.java#redirectAfterContinueShouldBeAbleToDeleteCookie  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestContinue.java
func TestPageRequestContinueRedirectAfterContinueShouldBeAbleToDeleteCookie(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	srv.SetRoute("/set-cookie", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "foo=bar;")
	})
	must.NoError(page.Goto(ctx, srv.Prefix()+"/set-cookie"))
	cookie, err := page.Evaluate(ctx, "() => document.cookie")
	must.NoError(err)
	is.Equal("foo=bar", cookie)

	srv.SetRoute("/delete-cookie", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "foo=bar; expires=Thu, 01 Jan 1970 00:00:00 GMT")
	})
	srv.ServeWithRedirect("/redirect", "/delete-cookie")
	must.NoError(page.Route(ctx, "**/redirect", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{Headers: route.Request().Headers()})
	}))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/redirect"))

	var mu sync.Mutex
	var cookieHdr []string
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		cookieHdr = r.Header.Values("Cookie")
		mu.Unlock()
		w.Header().Set("Content-Type", "text/html")
	})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	mu.Lock()
	defer mu.Unlock()
	is.Empty(cookieHdr)
}
