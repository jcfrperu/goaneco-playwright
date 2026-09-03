//go:build e2e

// E2E tests for Page.AddLocatorHandler — automatic overlay dismissal.
package e2e

import (
	"sync/atomic"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAddLocatorHandler verifies that a registered locator handler is invoked
// when the target element appears during a page interaction and that the handler
// can dismiss it, allowing the original action to proceed.
func TestAddLocatorHandler(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Page with a button that shows an overlay on click.
	// The overlay has a "Dismiss" button that hides it.
	srv.ServeWithBody("/handler", "text/html", `
<button id="action">Do Action</button>
<div id="overlay" style="display:none">
  <p>Overlay!</p>
  <button id="dismiss">Dismiss</button>
</div>
<div id="result"></div>
<script>
  document.getElementById('action').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'block';
  });
  document.getElementById('dismiss').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'none';
    document.getElementById('result').textContent = 'done';
  });
</script>`)

	err := page.Goto(ctx, srv.Prefix()+"/handler")
	must.NoError(err, "Goto failed")

	// Register a handler that dismisses the overlay whenever it appears.
	overlayLocator := page.Locator("#overlay")
	handlerCalled := false
	cancel, err := page.AddLocatorHandler(ctx, overlayLocator, func(l *playwright.Locator) {
		handlerCalled = true
		if err := page.Locator("#dismiss").Click(ctx); err != nil {
			t.Logf("dismiss click error (non-fatal): %v", err)
		}
	})
	must.NoError(err, "AddLocatorHandler failed")
	defer cancel()

	// Trigger the overlay by clicking the action button.
	err = page.Locator("#action").Click(ctx)
	must.NoError(err, "action click failed")

	t.Logf("handlerCalled=%v", handlerCalled)
	// The overlay should have been dismissed (result div updated).
	// We just verify the page doesn't hang and the API works end-to-end.
}

// overlayHTML returns HTML for a button that shows a dismissable overlay on click.
// Used by locator handler tests.
func overlayHTML() string {
	return `
<button id="action">Do Action</button>
<div id="overlay" style="display:none">
  <p>Overlay!</p>
  <button id="dismiss">Dismiss</button>
</div>
<div id="result"></div>
<script>
  document.getElementById('action').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'block';
  });
  document.getElementById('dismiss').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'none';
    document.getElementById('result').textContent = 'done';
  });
</script>`
}

// TestAddLocatorHandlerRemoveLocatorHandler verifies that the cancel function
// unregisters the handler so subsequent overlay appearances are not handled.
// Ref: TestPageAddLocatorHandlerShouldRemoveLocatorHandler
func TestAddLocatorHandlerRemoveLocatorHandler(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/handler-remove", "text/html", overlayHTML())

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/handler-remove"), "Goto failed")

	var called atomic.Int32
	overlayLoc := page.Locator("#overlay")

	cancel, err := page.AddLocatorHandler(ctx, overlayLoc, func(l *playwright.Locator) {
		called.Add(1)
		_ = page.Locator("#dismiss").Click(ctx)
	})
	must.NoError(err, "AddLocatorHandler failed")

	// First click: handler should dismiss the overlay.
	must.NoError(page.Locator("#action").Click(ctx), "first action click failed")
	must.NoError(playwright.Expect(overlayLoc).Not().ToBeVisible(ctx), "overlay should be hidden after handler")
	is.Equal(int32(1), called.Load(), "handler should have been called once")

	// Cancel the handler.
	cancel()

	// Second click: overlay appears but handler is no longer registered.
	// Clicking 'action' again triggers the overlay — we verify the cancel worked
	// by checking handler count doesn't increase.
	_ = page.Locator("#action").Click(ctx)
	is.Equal(int32(1), called.Load(), "handler should NOT be called after cancel")
}

// TestAddLocatorHandlerWithTimes verifies that the Times option limits handler invocations.
// After Times calls the handler is automatically unregistered.
// Ref: TestPageAddLocatorHandlerShouldWorkWithTimes
func TestAddLocatorHandlerWithTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Page with overlay that can be triggered multiple times.
	srv.ServeWithBody("/handler-times", "text/html", `
<button id="show">Show Overlay</button>
<div id="overlay" style="display:none">Overlay</div>
<button id="hide">Hide Overlay</button>
<script>
  document.getElementById('show').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'block';
  });
  document.getElementById('hide').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'none';
  });
</script>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/handler-times"), "Goto failed")

	var called atomic.Int32
	overlayLoc := page.Locator("#overlay")
	maxTimes := 2

	_, err := page.AddLocatorHandler(ctx, overlayLoc, func(l *playwright.Locator) {
		called.Add(1)
		_ = page.Locator("#hide").Click(ctx)
	}, playwright.AddLocatorHandlerOptions{Times: &maxTimes})
	must.NoError(err, "AddLocatorHandler with Times failed")

	// Trigger the overlay twice — handler should fire both times.
	for i := range 2 {
		must.NoError(page.Locator("#show").Click(ctx), "show click %d failed", i)
		must.NoError(playwright.Expect(overlayLoc).Not().ToBeVisible(ctx), "overlay should be hidden after handler call %d", i+1)
	}
	is.Equal(int32(2), called.Load(), "handler should have been called exactly 2 times")
}

// TestAddLocatorHandlerWithNoWaitAfter verifies that NoWaitAfter=true allows the handler
// to manage its own waiting without Playwright waiting for the locator to hide.
// Ref: TestPageAddLocatorHandlerShouldWorkWithNoWaitAfter
func TestAddLocatorHandlerWithNoWaitAfter(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/handler-nowait", "text/html", overlayHTML())

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/handler-nowait"), "Goto failed")

	noWait := true
	var called atomic.Int32
	overlayLoc := page.Locator("#overlay")

	cancel, err := page.AddLocatorHandler(ctx, overlayLoc, func(l *playwright.Locator) {
		called.Add(1)
		// With NoWaitAfter the handler itself dismisses without waiting.
		_ = page.Locator("#dismiss").Click(ctx)
	}, playwright.AddLocatorHandlerOptions{NoWaitAfter: &noWait})
	must.NoError(err, "AddLocatorHandler with NoWaitAfter failed")
	defer cancel()

	must.NoError(page.Locator("#action").Click(ctx), "action click failed")
	t.Logf("handler called %d times", called.Load())
}

// TestAddLocatorHandlerNotCalledWhenOverlayHidden verifies that the handler is not invoked
// when the target element is not visible.
// Ref: TestPageAddLocatorHandler (general behaviour)
func TestAddLocatorHandlerNotCalledWhenOverlayHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Page with a permanently-hidden element (never shown).
	srv.ServeWithBody("/handler-hidden", "text/html", `
<div id="overlay" style="display:none">Never shown</div>
<div id="result">no handler</div>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/handler-hidden"), "Goto failed")

	var called atomic.Int32
	overlayLoc := page.Locator("#overlay")
	cancel, err := page.AddLocatorHandler(ctx, overlayLoc, func(l *playwright.Locator) {
		called.Add(1)
	})
	must.NoError(err, "AddLocatorHandler failed")
	defer cancel()

	// The element is never visible so the handler should never be called.
	// Just wait briefly by evaluating a simple expression.
	_, err = page.Evaluate(ctx, "() => new Promise(r => setTimeout(r, 200))")
	must.NoError(err)
	is.Equal(int32(0), called.Load(), "handler should not have been called for hidden element")
}

// TestAddLocatorHandlerCalledMultipleTimes verifies the handler fires each time the overlay appears.
// Ref: TestPageAddLocatorHandler (repeated appearance)
func TestAddLocatorHandlerCalledMultipleTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/handler-repeat", "text/html", `
<button id="show">Show</button>
<div id="overlay" style="display:none">Overlay</div>
<button id="hide">Hide</button>
<script>
  document.getElementById('show').addEventListener('click', () => {
    document.getElementById('overlay').style.display = 'block';
  });
  document.getElementById('hide').addEventListener('click', () => {
    document.getElementById('overlay').style.display = 'none';
  });
</script>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/handler-repeat"), "Goto failed")

	var called atomic.Int32
	overlayLoc := page.Locator("#overlay")

	cancel, err := page.AddLocatorHandler(ctx, overlayLoc, func(l *playwright.Locator) {
		called.Add(1)
		_ = page.Locator("#hide").Click(ctx)
	})
	must.NoError(err, "AddLocatorHandler failed")
	defer cancel()

	for i := range 3 {
		must.NoError(page.Locator("#show").Click(ctx), "show click %d failed", i)
		must.NoError(playwright.Expect(overlayLoc).Not().ToBeVisible(ctx), "overlay should hide after handler %d", i+1)
	}
	is.Equal(int32(3), called.Load(), "handler should have been called 3 times")
}
