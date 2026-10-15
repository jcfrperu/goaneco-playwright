//go:build e2e

// E2E tests for: TestRequestContinue.java.
package e2e

import (
	"net/http"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequestContinueShouldAmendHTTPHeaders verifies headers can be amended when continuing.
// Ref: TestRequestContinue.java#shouldAmendHTTPHeaders  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldAmendHTTPHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		headers := map[string]string{}
		for k, v := range route.Request().Headers() {
			headers[k] = v
		}
		headers["FOO"] = "bar"
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{Headers: headers})
	}))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz')")
	must.NoError(err)
	rec.Wait(t)
	is.Equal([]string{"bar"}, rec.Header.Values("foo"))
}

// TestRequestContinueShouldAmendMethod verifies the method can be amended when continuing.
// Ref: TestRequestContinue.java#shouldAmendMethod  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldAmendMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{Method: a4Ptr("POST")})
	}))
	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz')")
	must.NoError(err)
	rec.Wait(t)
	is.Equal("POST", rec.Method)
}

// TestRequestContinueShouldOverrideRequestUrl verifies the URL can be overridden when continuing.
// Ref: TestRequestContinue.java#shouldOverrideRequestUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldOverrideRequestUrl(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	srv.ServeWithBody("/global-var.html", "text/html", "<script>window['globalVar'] = 123;</script>")
	rec := a4RecordRoute(srv, "/global-var.html")
	srv.SetRoute("/global-var.html", func(w http.ResponseWriter, r *http.Request) {
		rec.once.Do(func() { rec.Method = r.Method; close(rec.got) })
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<script>window['globalVar'] = 123;</script>"))
	})
	must.NoError(page.Route(ctx, "**/foo", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{URL: a4Ptr(srv.Prefix() + "/global-var.html")})
	}))
	log := a4Watch(t, page)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/foo"))
	resp := log.WaitResponse(t, "/global-var.html", 0)
	is.Equal(srv.Prefix()+"/global-var.html", resp.URL())
	v, err := page.Evaluate(ctx, "window['globalVar']")
	must.NoError(err)
	is.EqualValues(123, v)
	rec.Wait(t)
	is.Equal("GET", rec.Method)
}

// TestRequestContinueShouldOverrideMethodAlongWithUrl verifies method and URL can be overridden together.
// Ref: TestRequestContinue.java#shouldOverrideMethodAlongWithUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldOverrideMethodAlongWithUrl(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	rec := a4RecordRoute(srv, "/empty.html")
	must.NoError(page.Route(ctx, "**/foo", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{URL: a4Ptr(srv.EmptyPage()), Method: a4Ptr("POST")})
	}))
	_ = page.Goto(ctx, srv.Prefix()+"/foo")
	rec.Wait(t)
	is.Equal("POST", rec.Method)
}

// TestRequestContinueShouldAmendMethodOnMainRequest verifies the method of the main request can be amended.
// Ref: TestRequestContinue.java#shouldAmendMethodOnMainRequest  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldAmendMethodOnMainRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	rec := a4RecordRoute(srv, "/empty.html")
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{Method: a4Ptr("POST")})
	}))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	rec.Wait(t)
	is.Equal("POST", rec.Method)
}

// TestRequestContinueShouldAmendPostData verifies post data can be amended.
// Ref: TestRequestContinue.java#shouldAmendPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldAmendPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	data := []byte("doggo")
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{PostData: data})
	}))
	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz', { method: 'POST', body: 'birdy' })")
	must.NoError(err)
	rec.Wait(t)
	is.Equal("POST", rec.Method)
	is.Equal(data, rec.Body)
}

// TestRequestContinueShouldAmendUtf8PostData verifies UTF-8 post data can be amended.
// Ref: TestRequestContinue.java#shouldAmendUtf8PostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldAmendUtf8PostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	data := []byte("пушкин")
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{PostData: data})
	}))
	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz', { method: 'POST', body: 'birdy' })")
	must.NoError(err)
	rec.Wait(t)
	is.Equal("POST", rec.Method)
	is.Equal(data, rec.Body)
}

// TestRequestContinueShouldAmendLongerPostData verifies longer post data can be set.
// Ref: TestRequestContinue.java#shouldAmendLongerPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldAmendLongerPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	data := []byte("doggo-is-longer-than-birdy")
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{PostData: data})
	}))
	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz', { method: 'POST', body: 'birdy' })")
	must.NoError(err)
	rec.Wait(t)
	is.Equal("POST", rec.Method)
	is.Equal(data, rec.Body)
}

// TestRequestContinueShouldAmendBinaryPostData verifies binary post data can be set.
// Ref: TestRequestContinue.java#shouldAmendBinaryPostData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldAmendBinaryPostData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	data := a4Bytes256()
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Continue(ctx, &playwright.RouteContinueOptions{PostData: data})
	}))
	rec := a4RecordRoute(srv, "/sleep.zzz")
	_, err := page.Evaluate(ctx, "() => fetch('/sleep.zzz', { method: 'POST', body: 'birdy' })")
	must.NoError(err)
	rec.Wait(t)
	is.Equal("POST", rec.Method)
	is.Equal(data, rec.Body)
}

// TestRequestContinueShouldWork verifies Route.Continue lets the request through.
// Ref: TestRequestContinue.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestContinue.java
func TestRequestContinueShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) { _ = route.Continue(ctx, nil) }))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
}
