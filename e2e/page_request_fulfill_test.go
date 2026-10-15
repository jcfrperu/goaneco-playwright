//go:build e2e

// E2E tests for: TestPageRequestFulfill.java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPageRequestFulfillShouldFulfillWithGlobalFetchResult fulfills route with global APIRequestContext fetch result.
// Ref: TestPageRequestFulfill.java#shouldFulfillWithGlobalFetchResult  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFulfill.java
func TestPageRequestFulfillShouldFulfillWithGlobalFetchResult(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/simple.json", "application/json", "{\"foo\": \"bar\"}\n")
	request := a4NewAPICtx(t)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		response, err := request.Get(ctx, srv.Prefix()+"/simple.json")
		must.NoError(err)
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{Response: response}))
	}))

	response, err := page.WaitForResponse(ctx, "**/*", func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(response)
	is.Equal("application/json", response.Headers()["content-type"])
	is.Equal(200, response.Status())
	text, err := response.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}

// TestPageRequestFulfillShouldFulfillWithFetchResult fulfills route with page.Request() fetch result.
// Ref: TestPageRequestFulfill.java#shouldFulfillWithFetchResult  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFulfill.java
func TestPageRequestFulfillShouldFulfillWithFetchResult(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/simple.json", "application/json", "{\"foo\": \"bar\"}\n")
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		response, err := page.Request().Get(ctx, srv.Prefix()+"/simple.json")
		must.NoError(err)
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{Response: response}))
	}))

	response, err := page.WaitForResponse(ctx, "**/*", func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(response)
	is.Equal("application/json", response.Headers()["content-type"])
	is.Equal(200, response.Status())
	text, err := response.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}

// TestPageRequestFulfillShouldFulfillWithFetchResultAndOverrides fulfills route with overrides.
// Ref: TestPageRequestFulfill.java#shouldFulfillWithFetchResultAndOverrides  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestFulfill.java
func TestPageRequestFulfillShouldFulfillWithFetchResultAndOverrides(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/simple.json", "application/json", "{\"foo\": \"bar\"}\n")
	status := 201
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		response, err := page.Request().Get(ctx, srv.Prefix()+"/simple.json")
		must.NoError(err)
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Response: response,
			Status:   &status,
			Headers:  map[string]string{"Content-Type": "application/json", "foo": "bar"},
		}))
	}))

	response, err := page.WaitForResponse(ctx, "**/*", func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(response)
	is.Equal("application/json", response.Headers()["content-type"])
	is.Equal(201, response.Status())
	allHeaders, err := response.AllHeaders(ctx)
	must.NoError(err)
	is.Equal("bar", allHeaders["foo"])
	text, err := response.Text(ctx)
	must.NoError(err)
	is.Equal("{\"foo\": \"bar\"}\n", text)
}
