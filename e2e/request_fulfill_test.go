//go:build e2e

// E2E tests for: TestRequestFulfill.java.
package e2e

import (
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequestFulfillShouldWorkWithStatusCode422 verifies fulfilling with status 422.
// Ref: TestRequestFulfill.java#shouldWorkWithStatusCode422  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestFulfill.java
func TestRequestFulfillShouldWorkWithStatusCode422(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Status: a4Ptr(422), Body: a4Ptr("Yo, page!")})
	}))
	_ = page.Goto(ctx, srv.EmptyPage())
	resp := log.WaitResponse(t, "/empty.html", 0)
	is.Equal(422, resp.Status())
	is.Equal("Unprocessable Entity", resp.StatusText())
	v, err := page.Evaluate(ctx, "document.body.textContent")
	must.NoError(err)
	is.Equal("Yo, page!", v)
}

// TestRequestFulfillFulfillShouldThrowIfHandledTwice verifies fulfilling twice fails.
// Ref: TestRequestFulfill.java#fulfillShouldThrowIfHandledTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestFulfill.java
func TestRequestFulfillFulfillShouldThrowIfHandledTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	errCh := make(chan error, 1)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{})
		errCh <- route.Fulfill(ctx, &playwright.RouteFulfillOptions{})
	}))
	_ = page.Goto(ctx, srv.EmptyPage())
	select {
	case err := <-errCh:
		must.Error(err)
		is.Contains(err.Error(), "Route is already handled!")
	case <-time.After(10 * time.Second):
		t.Fatal("route handler did not run")
	}
}

// TestRequestFulfillAbortShouldThrowIfHandledTwice verifies aborting twice fails.
// Ref: TestRequestFulfill.java#abortShouldThrowIfHandledTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestFulfill.java
func TestRequestFulfillAbortShouldThrowIfHandledTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	errCh := make(chan error, 1)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Abort(ctx)
		errCh <- route.Abort(ctx)
	}))
	_ = page.Goto(ctx, srv.EmptyPage())
	select {
	case err := <-errCh:
		must.Error(err)
		is.Contains(err.Error(), "Route is already handled!")
	case <-time.After(10 * time.Second):
		t.Fatal("route handler did not run")
	}
}

// TestRequestFulfillResumeShouldThrowIfHandledTwice verifies continuing twice fails.
// Ref: TestRequestFulfill.java#resumeShouldThrowIfHandledTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestFulfill.java
func TestRequestFulfillResumeShouldThrowIfHandledTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	errCh := make(chan error, 1)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Continue(ctx, nil)
		errCh <- route.Continue(ctx, nil)
	}))
	_ = page.Goto(ctx, srv.EmptyPage())
	select {
	case err := <-errCh:
		must.Error(err)
		is.Contains(err.Error(), "Route is already handled!")
	case <-time.After(10 * time.Second):
		t.Fatal("route handler did not run")
	}
}

// TestRequestFulfillShouldWork verifies Route.Fulfill with status, headers and body.
// Ref: TestRequestFulfill.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestRequestFulfill.java
func TestRequestFulfillShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	log := a4Watch(t, page)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status:      a4Ptr(201),
			ContentType: a4Ptr("text/html"),
			Headers:     map[string]string{"foo": "bar"},
			Body:        a4Ptr("Yo, page!"),
		})
	}))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	resp := log.WaitResponse(t, "/empty.html", 0)
	is.Equal(201, resp.Status())
	is.Equal("bar", resp.Headers()["foo"])
	v, err := page.Evaluate(ctx, "() => document.body.textContent")
	must.NoError(err)
	is.Equal("Yo, page!", v)
}
