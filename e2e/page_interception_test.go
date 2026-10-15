//go:build e2e

// E2E tests for: TestPageInterception.java.
package e2e

import (
	"io"
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

// TestPageInterceptionShouldWorkWithNavigationSmoke verifies IsNavigationRequest for intercepted requests.
// Ref: TestPageInterception.java#shouldWorkWithNavigationSmoke  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldWorkWithNavigationSmoke(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	srv.ServeWithBody("/frames/one-frame.html", "text/html", `<link rel="stylesheet" href="./style.css"><script src="./script.js"></script><iframe src="./frame.html"></iframe>`)
	srv.ServeWithBody("/frames/frame.html", "text/html", "<div>frame</div>")
	srv.ServeWithBody("/frames/script.js", "text/javascript", "window.scriptLoaded = true;")
	srv.ServeWithBody("/frames/style.css", "text/css", "body { color: red; }")
	srv.ServeWithRedirect("/rrredirect", "/frames/one-frame.html")

	var mu sync.Mutex
	requests := map[string]*playwright.NetworkRequest{}
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		parts := strings.Split(route.Request().URL(), "/")
		mu.Lock()
		requests[parts[len(parts)-1]] = route.Request()
		mu.Unlock()
		_ = route.Continue(ctx, nil)
	}))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/rrredirect"))
	is.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return requests["style.css"] != nil && requests["script.js"] != nil && requests["frame.html"] != nil
	}, 5*time.Second, 20*time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	must.NotNil(requests["rrredirect"])
	is.True(requests["rrredirect"].IsNavigationRequest())
	is.True(requests["frame.html"].IsNavigationRequest())
	is.False(requests["script.js"].IsNavigationRequest())
	is.False(requests["style.css"].IsNavigationRequest())
}

// TestPageInterceptionShouldInterceptAfterAServiceWorker is not ported yet.
// Ref: TestPageInterception.java#shouldInterceptAfterAServiceWorker  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldInterceptAfterAServiceWorker(t *testing.T) {
	t.Skip("not implemented: service worker test resources (/serviceworkers/fetchdummy) are not served by testserver")
}

// TestPageInterceptionShouldFulfillInterceptedResponseUsingAlias verifies Route.Fetch and fulfill with response.
// Ref: TestPageInterception.java#shouldFulfillInterceptedResponseUsingAlias  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldFulfillInterceptedResponseUsingAlias(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		response, err := route.Fetch(ctx)
		must.NoError(err)
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{Response: response}))
	}))

	response, err := page.WaitForResponse(ctx, "**/*", func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(response)
	is.Equal(200, response.Status())
	is.Contains(response.Headers()["content-type"], "text/html")
}

// TestPageInterceptionShouldInterceptWithUrlOverride verifies Route.Fetch with URL override.
// Ref: TestPageInterception.java#shouldInterceptWithUrlOverride  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldInterceptWithUrlOverride(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/one-style.html", "text/html", `<link rel="stylesheet" href="./one-style.css">`)
	targetURL := srv.Prefix() + "/one-style.html"
	must.NoError(page.Route(ctx, "**/*.html", func(route *playwright.Route) {
		response, err := route.Fetch(ctx, &playwright.RouteFetchOptions{URL: &targetURL})
		must.NoError(err)
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{Response: response}))
	}))

	response, err := page.WaitForResponse(ctx, "**/*", func() error {
		return page.Goto(ctx, srv.EmptyPage())
	})
	must.NoError(err)
	must.NotNil(response)
	is.Equal(200, response.Status())
	text, err := response.Text(ctx)
	must.NoError(err)
	is.Contains(text, "one-style.css")
}

// TestPageInterceptionShouldInterceptWithPostDataOverride verifies Route.Fetch with PostData override.
// Ref: TestPageInterception.java#shouldInterceptWithPostDataOverride  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldInterceptWithPostDataOverride(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var capturedBody string
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		capturedBody = string(body)
		mu.Unlock()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<html></html>"))
	})

	must.NoError(page.Route(ctx, "**/*.html", func(route *playwright.Route) {
		response, err := route.Fetch(ctx, &playwright.RouteFetchOptions{PostData: `{ "foo": "bar" }`})
		must.NoError(err)
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{Response: response}))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	mu.Lock()
	defer mu.Unlock()
	is.Equal(`{ "foo": "bar" }`, capturedBody)
}

// TestPageInterceptionShouldWorkWithGlob is not ported yet.
// Ref: TestPageInterception.java#shouldWorkWithGlob  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldWorkWithGlob(t *testing.T) {
	t.Skip("not implemented: Java-internal Utils.urlMatches glob matcher (no Go counterpart)")
}

// TestPageInterceptionShouldSupportTimeoutOptionInRouteFetch verifies Timeout option in Route.Fetch.
// Ref: TestPageInterception.java#shouldSupportTimeoutOptionInRouteFetch  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldSupportTimeoutOptionInRouteFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.SetRoute("/slow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		if f, ok := w.(http.Flusher); ok {
			_, _ = w.Write([]byte("start"))
			f.Flush()
		}
		<-r.Context().Done()
	})

	to := float64(1000)
	var routeFetchErr error
	var wg sync.WaitGroup
	var slowRoute *playwright.Route
	wg.Add(1)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		slowRoute = route
		_, routeFetchErr = route.Fetch(ctx, &playwright.RouteFetchOptions{Timeout: &to})
		wg.Done()
	}))

	navTo := float64(2000)
	navErr := page.Goto(ctx, srv.Prefix()+"/slow", &playwright.PageGotoOptions{Timeout: &navTo})
	must.Error(navErr)
	is.Contains(navErr.Error(), "Timeout 2000ms exceeded")
	wg.Wait()
	must.Error(routeFetchErr)
	is.Contains(routeFetchErr.Error(), "Timeout 1000ms exceeded")

	if slowRoute != nil {
		_ = slowRoute.Abort(ctx)
	}
}

// TestPageInterceptionShouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0InRouteFetch verifies MaxRedirects=0 in Route.Fetch.
// Ref: TestPageInterception.java#shouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0InRouteFetch  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageInterception.java
func TestPageInterceptionShouldNotFollowRedirectsWhenMaxRedirectsIsSetTo0InRouteFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithRedirect("/foo", "/empty.html")
	maxRedir := 0
	bodyText := "hello"
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		response, err := route.Fetch(ctx, &playwright.RouteFetchOptions{MaxRedirects: &maxRedir})
		must.NoError(err)
		is.Equal("/empty.html", response.Headers()["location"])
		is.Equal(302, response.Status())
		must.NoError(route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &bodyText}))
	}))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/foo"))
}
