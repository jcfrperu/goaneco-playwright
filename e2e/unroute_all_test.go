//go:build e2e

// E2E tests for Page.UnrouteAll.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
)

// TestPageUnroute verifies that Unroute() removes all page-level routes, allowing
// subsequent requests to pass through to the real server.
func TestPageUnroute(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/unroute-real", "text/html", "<title>RealPage</title>")
	page := newPage(t)

	// Block all traffic.
	status := 503
	body := "blocked"
	err := page.Route(ctx, "**/*", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Status: &status, Body: &body})
	})
	must.NoError(err, "Route failed")

	err = page.Unroute(ctx)
	must.NoError(err, "Unroute failed")

	err = page.Goto(ctx, srv.Prefix()+"/unroute-real")
	must.NoError(err, "Goto after Unroute failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title failed")
	if title != "RealPage" {
		t.Errorf("Title = %q, want 'RealPage'; route may still be active", title)
	}
}

// TestBrowserContextUnroute verifies that BrowserContext.Unroute() removes context-level routes.
func TestBrowserContextUnroute(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-unroute", "text/html", "<title>CtxReal</title>")
	bCtx := newContext(t)

	// Block all traffic at the context level.
	status := 503
	body := "blocked"
	err := bCtx.Route(ctx, "**/*", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Status: &status, Body: &body})
	})
	must.NoError(err, "Context Route failed")

	err = bCtx.Unroute(ctx)
	must.NoError(err, "Context Unroute failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.Prefix()+"/ctx-unroute")
	must.NoError(err, "Goto after context Unroute failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title failed")
	if title != "CtxReal" {
		t.Errorf("Title = %q, want 'CtxReal'; context route may still be active", title)
	}
}

// TestPageUnrouteAll verifies that UnrouteAll removes all route handlers so subsequent
// requests pass through to the network normally.
func TestPageUnrouteAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/pass", "text/html", `<title>Pass</title>`)

	// Register a route that fulfills all requests with a 404 (so we can detect it)
	status := 404
	notFoundBody := "blocked"
	err := page.Route(ctx, "**/*", func(r *playwright.Route) {
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status,
			Body:   &notFoundBody,
		})
	})
	must.NoError(err, "Route failed")

	// Remove all routes with UnrouteAll
	err = page.UnrouteAll(ctx)
	must.NoError(err, "UnrouteAll failed")

	// After UnrouteAll, navigation should reach the real server
	err = page.Goto(ctx, srv.Prefix()+"/pass")
	must.NoError(err, "Goto after UnrouteAll failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title failed")
	if title != "Pass" {
		t.Errorf("Title = %q, want 'Pass'", title)
	}
}
