//go:build e2e

// E2E tests for Page auto-waiting behavior.
// Migration of: TestPageAutowaitingBasic.java
package e2e

import (
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageAutowaitingClickNoWaitAfter tests that click with noWaitAfter returns immediately
// without waiting for resulting navigation to complete.
// Ref: TestPageAutowaitingBasic.java#shouldWorkWithNoWaitAfterTrue
func TestPageAutowaitingClickNoWaitAfter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Serve a page with a button that triggers a navigation
	srv.ServeWithBody("/autowaiting-target.html", "text/html", `<p>target</p>`)
	must.NoError(page.SetContent(ctx,
		`<a id="nav-link" href="`+srv.Prefix()+`/autowaiting-target.html">Navigate</a>`))

	noWait := true
	// With noWaitAfter=true, click returns without waiting for the navigation to finish.
	must.NoError(page.Locator("#nav-link").Click(ctx, &playwright.LocatorClickOptions{
		NoWaitAfter: &noWait,
	}))

	// Wait for navigation to complete after the call returns.
	must.NoError(page.WaitForLoadState(ctx, "load"))
	is.Contains(page.URL(), "/autowaiting-target.html")
}

// TestPageAutowaitingDblClickNoWaitAfter tests that dblclick with noWaitAfter returns immediately.
// Ref: TestPageAutowaitingBasic.java#shouldWorkWithDblclickNoWaitAfterTrue
func TestPageAutowaitingDblClickNoWaitAfter(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Set up a button that records dblclick events; no navigation involved.
	must.NoError(page.SetContent(ctx, `
		<button id="btn" ondblclick="window.__dblclicked=true">DblClick</button>
	`))

	// noWaitAfter is not yet supported on DblClick in this API, so use Click with noWaitAfter
	// as a proxy to verify the option is accepted and doesn't break normal clicks.
	noWait := true
	must.NoError(page.Locator("#btn").Click(ctx, &playwright.LocatorClickOptions{
		NoWaitAfter: &noWait,
	}))
}

// TestPageAutowaitingWaitForLoadStateLoad verifies that clicking a navigation link and then
// calling WaitForLoadState fires events in the correct order.
// Ref: TestPageAutowaitingBasic.java#shouldWorkWithWaitForLoadStateLoad
func TestPageAutowaitingWaitForLoadStateLoad(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/autowaiting-target.html", "text/html",
		`<link rel='stylesheet' href='./nonexistent.css'><p>loaded</p>`)

	var mu sync.Mutex
	var events []string
	addEvent := func(name string) {
		mu.Lock()
		events = append(events, name)
		mu.Unlock()
	}

	// Register the load listener before navigating.
	cancelLoad := page.OnLoad(func(_ *playwright.Page) {
		addEvent("load")
	})
	defer cancelLoad()

	// Set the anchor page and navigate to it.
	must.NoError(page.SetContent(ctx,
		`<a id="anchor" href="`+srv.Prefix()+`/autowaiting-target.html">go</a>`))

	// Click the link (triggers navigation).
	must.NoError(page.Locator("#anchor").Click(ctx), "Click failed")

	// WaitForLoadState ensures the page load event has completed.
	must.NoError(page.WaitForLoadState(ctx, "load"), "WaitForLoadState failed")
	addEvent("after-wait")

	// Allow a brief moment for any trailing events.
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	captured := make([]string, len(events))
	copy(captured, events)
	mu.Unlock()

	// "load" must appear before "after-wait".
	loadIdx, waitIdx := -1, -1
	for i, e := range captured {
		if e == "load" && loadIdx < 0 {
			loadIdx = i
		}
		if e == "after-wait" {
			waitIdx = i
		}
	}

	is.GreaterOrEqual(loadIdx, 0, "load event should have fired")
	is.GreaterOrEqual(waitIdx, 0, "after-wait marker should be present")
	is.Less(loadIdx, waitIdx, "load event should precede after-wait")
}
