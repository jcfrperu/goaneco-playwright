//go:build e2e

// E2E tests for Page.OnRequestFailed event.
// Migration of: TestPageNetworkEvents.java (requestFailed cases)
package e2e

import (
	"sync"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOnRequestFailedFiredForAbortedRoute verifies OnRequestFailed fires when route aborts.
// Ref: TestPageNetworkEvents.java#shouldFireRequestFailedForAbortedRoute
func TestOnRequestFailedFiredForAbortedRoute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var failedURL string

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		failedURL = req.URL()
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Route(ctx, "**/blocked", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/blocked'); } catch {}
	}`)

	mu.Lock()
	url := failedURL
	mu.Unlock()

	is.Contains(url, "blocked")
}

// TestOnRequestFailedURLNotEmpty verifies OnRequestFailed provides a URL.
// Ref: TestPageNetworkEvents.java#shouldProvideURLInRequestFailed
func TestOnRequestFailedURLNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	urls := make([]string, 0)

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Route(ctx, "**/fail", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/fail'); } catch {}
	}`)

	mu.Lock()
	captured := len(urls)
	mu.Unlock()

	is.Greater(captured, 0)
}

// TestOnRequestFailedOffStopsReceiving verifies unregistering the handler stops events.
// Ref: TestPageNetworkEvents.java#shouldStopAfterUnregister
func TestOnRequestFailedOffStopsReceiving(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	count := 0

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	must.NoError(page.Route(ctx, "**/abort1", func(route *playwright.Route) {
		_ = route.Abort(ctx)
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/abort1'); } catch {}
	}`)

	off() // unregister

	mu.Lock()
	countAfterOff := count
	mu.Unlock()

	// Any subsequent fails after off() should not increment
	_, _ = page.Evaluate(ctx, `async () => {
		try { await fetch('/abort1'); } catch {}
	}`)

	mu.Lock()
	is.Equal(countAfterOff, count)
	mu.Unlock()
}
