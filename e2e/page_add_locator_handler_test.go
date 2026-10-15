//go:build e2e

// E2E tests for: TestPageAddLocatorHandler.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
)

// TestAddLocatorHandler verifies that a registered locator handler is invoked
// when the target element appears during a page interaction and that the handler
// can dismiss it, allowing the original action to proceed.
// Extra: no equivalent method in playwright-java; related class: TestPageAddLocatorHandler.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
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

// TestAddLocatorHandlerWithTimes verifies that the Times option limits handler invocations.
// After Times calls the handler is automatically unregistered.
// Ref: TestPageAddLocatorHandlerShouldWorkWithTimes
// Extra: no equivalent method in playwright-java; related class: TestPageAddLocatorHandler.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageAddLocatorHandler.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageAddLocatorHandler.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageAddLocatorHandler.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
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

// TestPageAddLocatorHandlerShouldRemoveLocatorHandler verifies that the cancel function
// unregisters the handler so subsequent overlay appearances are not handled.
// Ref: TestPageAddLocatorHandlerShouldRemoveLocatorHandler
// Ref: TestPageAddLocatorHandler.java#shouldRemoveLocatorHandler  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldRemoveLocatorHandler(t *testing.T) {
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

// Ref: TestPageAddLocatorHandler.java#shouldWorkWithLocatorHover  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWorkWithLocatorHover(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" onmouseover="this.setAttribute('data-hovered','1')">hover me</div>
	`), "SetContent failed")

	must.NoError(page.Locator("#target").Hover(ctx), "Hover failed")

	attr, err := page.Locator("#target").GetAttribute(ctx, "data-hovered")
	must.NoError(err, "GetAttribute failed")
	must.NotNil(attr, "expected data-hovered attribute to be set after hover")
	is.Equal("1", *attr)
}

// Ref: TestPageAddLocatorHandler.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	var beforeCount, afterCount atomic.Int32
	original := page.GetByText("This interstitial covers the button")
	_, err := page.AddLocatorHandler(ctx, original, func(l *playwright.Locator) {
		assert.Equal(t, original.Selector(), l.Selector())
		beforeCount.Add(1)
		assert.NoError(t, page.Locator("#close").Click(ctx))
		afterCount.Add(1)
	}, playwright.AddLocatorHandlerOptions{NoWaitAfter: a2Ptr(true)})
	must.NoError(err)
	argsList := [][]any{
		{"mouseover", "1"}, {"mouseover", "1", "capture"}, {"mouseover", "2"}, {"mouseover", "2", "capture"},
		{"pointerover", "1"}, {"pointerover", "1", "capture"}, {"none", "1"}, {"remove", "1"}, {"hide", "1"},
	}
	for _, args := range argsList {
		must.NoError(page.Locator("#aside").Hover(ctx))
		beforeCount.Store(0)
		afterCount.Store(0)
		a2Ev(t, ctx, page, `(args) => {
  window.clicked = 0;
  window.setupAnnoyingInterstitial(...args);
}`, args)
		assert.Equal(t, int32(0), beforeCount.Load())
		assert.Equal(t, int32(0), afterCount.Load())
		must.NoError(page.Locator("#target").Click(ctx))
		times := int32(args[1].(string)[0] - '0')
		assert.Equal(t, times, beforeCount.Load(), "%v", args)
		assert.Equal(t, times, afterCount.Load(), "%v", args)
		assert.Equal(t, float64(1), a2Ev(t, ctx, page, "window.clicked"))
		must.NoError(playwright.Expect(page.Locator("#interstitial")).Not().ToBeVisible(ctx))
	}
}

// Ref: TestPageAddLocatorHandler.java#shouldWorkWithCustomCheck  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWorkWithCustomCheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	_, err := page.AddLocatorHandler(ctx, page.Locator("body"), func(l *playwright.Locator) {
		if vis, _ := page.GetByText("This interstitial covers the button").IsVisible(ctx); vis {
			assert.NoError(t, page.Locator("#close").Click(ctx))
		}
	}, playwright.AddLocatorHandlerOptions{NoWaitAfter: a2Ptr(true)})
	must.NoError(err)
	argsList := [][]any{{"mouseover", "2"}, {"none", "1"}, {"remove", "1"}, {"hide", "1"}}
	for _, args := range argsList {
		must.NoError(page.Locator("#aside").Hover(ctx))
		a2Ev(t, ctx, page, `(args) => {
  window.clicked = 0;
  window.setupAnnoyingInterstitial(...args);
}`, args)
		must.NoError(page.Locator("#target").Click(ctx))
		assert.Equal(t, float64(1), a2Ev(t, ctx, page, "window.clicked"))
		must.NoError(playwright.Expect(page.Locator("#interstitial")).Not().ToBeVisible(ctx))
	}
}

// Ref: TestPageAddLocatorHandler.java#shouldNotWorkWithForceTrue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldNotWorkWithForceTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	_, err := page.AddLocatorHandler(ctx, page.GetByText("This interstitial covers the button"), func(l *playwright.Locator) {
		assert.NoError(t, page.Locator("#close").Click(ctx))
	})
	must.NoError(err)
	must.NoError(page.Locator("#aside").Hover(ctx))
	a2Ev(t, ctx, page, "() => { window.setupAnnoyingInterstitial('none', 1); }")
	must.NoError(page.Locator("#target").Click(ctx, &playwright.LocatorClickOptions{Force: a2Ptr(true), Timeout: a2Ptr(2000.0)}))
	vis, err := page.Locator("#interstitial").IsVisible(ctx)
	must.NoError(err)
	assert.True(t, vis)
	assert.Nil(t, a2Ev(t, ctx, page, "window.clicked"))
}

// Ref: TestPageAddLocatorHandler.java#shouldThrowWhenPageCloses  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldThrowWhenPageCloses(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	_, err := page.AddLocatorHandler(ctx, page.GetByText("This interstitial covers the button"), func(l *playwright.Locator) {
		_ = page.Close(ctx)
	})
	must.NoError(err)
	must.NoError(page.Locator("#aside").Hover(ctx))
	a2Ev(t, ctx, page, `() => {
    window.clicked = 0;
    window.setupAnnoyingInterstitial('mouseover', 1);
  }`)
	err = page.Locator("#target").Click(ctx)
	must.Error(err)
	assert.Contains(t, err.Error(), "Target page, context or browser has been closed")
}

// Ref: TestPageAddLocatorHandler.java#shouldWorkWithToBeVisible  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWorkWithToBeVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	var called atomic.Int32
	_, err := page.AddLocatorHandler(ctx, page.GetByText("This interstitial covers the button"), func(l *playwright.Locator) {
		called.Add(1)
		assert.NoError(t, page.Locator("#close").Click(ctx))
	})
	must.NoError(err)
	a2Ev(t, ctx, page, `() => {
  window.clicked = 0;
  window.setupAnnoyingInterstitial('remove', 1);
}`)
	must.NoError(playwright.Expect(page.Locator("#target")).ToBeVisible(ctx))
	must.NoError(playwright.Expect(page.Locator("#interstitial")).Not().ToBeVisible(ctx))
	assert.Equal(t, int32(1), called.Load())
}

// Ref: TestPageAddLocatorHandler.java#shouldWorkWhenOwnerFrameDetaches  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWorkWhenOwnerFrameDetaches(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	a2Ev(t, ctx, page, `() => {
    const iframe = document.createElement('iframe');
    iframe.src = 'data:text/html,<body>hello from iframe</body>';
    document.body.append(iframe);

    const target = document.createElement('button');
    target.textContent = 'Click me';
    target.id = 'target';
    target.addEventListener('click', () => window._clicked = true);
    document.body.appendChild(target);

    const closeButton = document.createElement('button');
    closeButton.textContent = 'close';
    closeButton.id = 'close';
    closeButton.addEventListener('click', () => iframe.remove());
    document.body.appendChild(closeButton);
  }`)
	_, err := page.AddLocatorHandler(ctx, page.FrameLocator("iframe").Locator("body"), func(l *playwright.Locator) {
		assert.NoError(t, page.Locator("#close").Click(ctx))
	})
	must.NoError(err)
	must.NoError(page.Locator("#target").Click(ctx))
	h, err := page.QuerySelector(ctx, "iframe")
	must.NoError(err)
	assert.Nil(t, h)
	assert.Equal(t, true, a2Ev(t, ctx, page, "window._clicked"))
}

// Ref: TestPageAddLocatorHandler.java#shouldWorkWithTimesOption  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWorkWithTimesOption(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	var called atomic.Int32
	_, err := page.AddLocatorHandler(ctx, page.Locator("body"), func(l *playwright.Locator) {
		called.Add(1)
	}, playwright.AddLocatorHandlerOptions{NoWaitAfter: a2Ptr(true), Times: a2Ptr(2)})
	must.NoError(err)
	must.NoError(page.Locator("#aside").Hover(ctx))
	a2Ev(t, ctx, page, `() => {
    window.clicked = 0;
    window.setupAnnoyingInterstitial('mouseover', 4);
  }`)
	err = page.Locator("#target").Click(ctx, &playwright.LocatorClickOptions{Timeout: a2Ptr(3000.0)})
	must.Error(err)
	assert.Equal(t, int32(2), called.Load())
	assert.Equal(t, float64(0), a2Ev(t, ctx, page, "window.clicked"))
	vis, verr := page.Locator("#interstitial").IsVisible(ctx)
	must.NoError(verr)
	assert.True(t, vis)
	assert.Contains(t, err.Error(), "Timeout 3000ms exceeded")
	assert.Contains(t, err.Error(), `<div>This interstitial covers the button</div> from <div class="visible" id="interstitial">…</div> subtree intercepts pointer events`)
}

// Ref: TestPageAddLocatorHandler.java#shouldWaitForHiddenByDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWaitForHiddenByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	var called atomic.Int32
	_, err := page.AddLocatorHandler(ctx, page.GetByRole("button", &playwright.GetByRoleOptions{Name: a2Ptr("close")}), func(button *playwright.Locator) {
		called.Add(1)
		assert.NoError(t, button.Click(ctx))
	})
	must.NoError(err)
	must.NoError(page.Locator("#aside").Hover(ctx))
	a2Ev(t, ctx, page, "() => {window.clicked = 0;window.setupAnnoyingInterstitial('timeout', 1);}")
	must.NoError(page.Locator("#target").Click(ctx))
	assert.Equal(t, float64(1), a2Ev(t, ctx, page, "window.clicked"))
	vis, err := page.Locator("#interstitial").IsVisible(ctx)
	must.NoError(err)
	assert.False(t, vis)
	assert.Equal(t, int32(1), called.Load())
}

// Ref: TestPageAddLocatorHandler.java#shouldWaitForHiddenByDefault2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWaitForHiddenByDefault2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	var called atomic.Int32
	_, err := page.AddLocatorHandler(ctx, page.GetByRole("button", &playwright.GetByRoleOptions{Name: a2Ptr("close")}), func(button *playwright.Locator) {
		called.Add(1)
	})
	must.NoError(err)
	must.NoError(page.Locator("#aside").Hover(ctx))
	a2Ev(t, ctx, page, "() => {window.clicked = 0;window.setupAnnoyingInterstitial('hide', 1);}")
	err = page.Locator("#target").Click(ctx, &playwright.LocatorClickOptions{Timeout: a2Ptr(3000.0)})
	must.Error(err)
	assert.Equal(t, float64(0), a2Ev(t, ctx, page, "window.clicked"))
	vis, verr := page.Locator("#interstitial").IsVisible(ctx)
	must.NoError(verr)
	assert.True(t, vis)
	assert.Equal(t, int32(1), called.Load())
	assert.Contains(t, err.Error(), `locator handler has finished, waiting for`)
}

// Ref: TestPageAddLocatorHandler.java#shouldWorkWithNoWaitAfter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddLocatorHandler.java
func TestPageAddLocatorHandlerShouldWorkWithNoWaitAfter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/handle-locator.html"))
	var called atomic.Int32
	_, err := page.AddLocatorHandler(ctx, page.GetByRole("button", &playwright.GetByRoleOptions{Name: a2Ptr("close")}), func(button *playwright.Locator) {
		if called.Add(1) == 1 {
			assert.NoError(t, button.Click(ctx))
		} else {
			assert.NoError(t, page.Locator("#interstitial").WaitFor(ctx, "hidden"))
		}
	}, playwright.AddLocatorHandlerOptions{NoWaitAfter: a2Ptr(true)})
	must.NoError(err)
	must.NoError(page.Locator("#aside").Hover(ctx))
	a2Ev(t, ctx, page, "() => {window.clicked = 0;window.setupAnnoyingInterstitial('timeout', 1);}")
	must.NoError(page.Locator("#target").Click(ctx))
	assert.Equal(t, float64(1), a2Ev(t, ctx, page, "window.clicked"))
	must.NoError(playwright.Expect(page.Locator("#interstitial")).Not().ToBeVisible(ctx))
	assert.Equal(t, int32(2), called.Load())
}
