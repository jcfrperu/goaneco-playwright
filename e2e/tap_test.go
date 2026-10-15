//go:build e2e

// E2E tests for: TestTap.java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestTap.java#shouldSendAllOfTheCorrectEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTap.java
func TestTapShouldSendAllOfTheCorrectEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2TouchPage(t)

	must.NoError(page.SetContent(ctx, "<div id='a' style='background: lightblue; width: 50px; height: 50px'>a</div>\n<div id='b' style='background: pink; width: 50px; height: 50px'>b</div>"))
	must.NoError(page.Locator("#a").Tap(ctx))
	b, err := page.QuerySelector(ctx, "#b")
	must.NoError(err)
	eventsHandle, err := a2TrackTapEvents(ctx, b)
	must.NoError(err)
	must.NoError(page.Locator("#b").Tap(ctx))
	// webkit doesnt send pointerenter or pointerleave or mouseout
	events, err := eventsHandle.JSONValue(ctx)
	must.NoError(err)
	assert.Equal(t, []any{"pointerover", "pointerenter",
		"pointerdown", "touchstart",
		"pointerup", "pointerout",
		"pointerleave", "touchend",
		"mouseover", "mouseenter",
		"mousemove", "mousedown",
		"mouseup", "click"}, events)
}

// Ref: TestTap.java#shouldNotSendMouseEventsTouchstartIsCanceled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTap.java
func TestTapShouldNotSendMouseEventsTouchstartIsCanceled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2TouchPage(t)

	must.NoError(page.SetContent(ctx, "<div style='width: 50px; height: 50px; background: red'>"))
	a2Ev(t, ctx, page, `() => {
    // touchstart is not cancelable unless passive is false
    document.addEventListener('touchstart', t => t.preventDefault(), {passive: false});
  }`)
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	eventsHandle, err := a2TrackTapEvents(ctx, div)
	must.NoError(err)
	must.NoError(page.Locator("div").Tap(ctx))
	events, err := eventsHandle.JSONValue(ctx)
	must.NoError(err)
	assert.Equal(t, []any{"pointerover", "pointerenter",
		"pointerdown", "touchstart",
		"pointerup", "pointerout",
		"pointerleave", "touchend"}, events)
}

// Ref: TestTap.java#shouldNotSendMouseEventsWhenTouchendIsCanceled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTap.java
func TestTapShouldNotSendMouseEventsWhenTouchendIsCanceled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2TouchPage(t)

	must.NoError(page.SetContent(ctx, "<div style='width: 50px; height: 50px; background: red'>"))
	a2Ev(t, ctx, page, "() => document.addEventListener('touchend', t => t.preventDefault())")
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	eventsHandle, err := a2TrackTapEvents(ctx, div)
	must.NoError(err)
	must.NoError(page.Locator("div").Tap(ctx))
	events, err := eventsHandle.JSONValue(ctx)
	must.NoError(err)
	assert.Equal(t, []any{"pointerover", "pointerenter",
		"pointerdown", "touchstart",
		"pointerup", "pointerout",
		"pointerleave", "touchend"}, events)
}

// Ref: TestTap.java#shouldNotWaitForANavigationCausedByATap  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTap.java
func TestTapShouldNotWaitForANavigationCausedByATap(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2TouchPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, "<a href='/intercept-this.html'>link</a>;"))
	received := srv.RequestReceived("/intercept-this.html")
	must.NoError(page.Locator("a").Tap(ctx))
	select {
	case <-received:
	case <-time.After(10 * time.Second):
		t.Fatal("request for /intercept-this.html was not received")
	}
}

// Ref: TestTap.java#shouldWorkWithModifiers  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTap.java
func TestTapShouldWorkWithModifiers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2TouchPage(t)

	must.NoError(page.SetContent(ctx, "hello world"))
	_, err := page.Evaluate(ctx, `() => {
	  window.touchPromise = new Promise(resolve => {
	    document.addEventListener('touchstart', event => {
	      resolve(event.altKey);
	    }, {passive: false});
	  });
	}`)
	must.NoError(err)
	must.NoError(page.Locator("body").Tap(ctx, &playwright.LocatorTapOptions{Modifiers: []string{"Alt"}}))
	altKey, err := page.Evaluate(ctx, "() => window.touchPromise")
	must.NoError(err)
	assert.Equal(t, true, altKey)
}

// Ref: TestTap.java#shouldSendWellFormedTouchPoints  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestTap.java
func TestTapShouldSendWellFormedTouchPoints(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, _, page := a2TouchPage(t)

	_, err := page.Evaluate(ctx, `() => {
	  window.touchStartPromise = new Promise(resolve => {
	    document.addEventListener('touchstart', event => {
	      resolve([...event.touches].map(t => ({
	        identifier: t.identifier,
	        clientX: t.clientX,
	        clientY: t.clientY,
	        pageX: t.pageX,
	        pageY: t.pageY,
	        radiusX: 'radiusX' in t ? t.radiusX : t['webkitRadiusX'],
	        radiusY: 'radiusY' in t ? t.radiusY : t['webkitRadiusY'],
	        rotationAngle: 'rotationAngle' in t ? t.rotationAngle : t['webkitRotationAngle'],
	        force: 'force' in t ? t.force : t['webkitForce'],
	      })));
	    }, false);
	  });
	}`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
	  window.touchEndPromise = new Promise(resolve => {
	    document.addEventListener('touchend', event => {
	      resolve([...event.touches].map(t => ({
	        identifier: t.identifier,
	        clientX: t.clientX,
	        clientY: t.clientY,
	        pageX: t.pageX,
	        pageY: t.pageY,
	        radiusX: 'radiusX' in t ? t.radiusX : t['webkitRadiusX'],
	        radiusY: 'radiusY' in t ? t.radiusY : t['webkitRadiusY'],
	        rotationAngle: 'rotationAngle' in t ? t.rotationAngle : t['webkitRotationAngle'],
	        force: 'force' in t ? t.force : t['webkitForce'],
	      })));
	    }, false);
	  });
	}`)
	must.NoError(err)

	must.NoError(page.Touchscreen().Tap(ctx, 40, 60))

	touchStart, err := page.Evaluate(ctx, "() => window.touchStartPromise")
	must.NoError(err)

	expected := []any{
		map[string]any{
			"clientX":       float64(40),
			"clientY":       float64(60),
			"force":         float64(1),
			"identifier":    float64(0),
			"pageX":         float64(40),
			"pageY":         float64(60),
			"radiusX":       float64(1),
			"radiusY":       float64(1),
			"rotationAngle": float64(0),
		},
	}
	is.Equal(expected, touchStart)

	touchEnd, err := page.Evaluate(ctx, "() => window.touchEndPromise")
	must.NoError(err)
	is.Equal([]any{}, touchEnd)
}
