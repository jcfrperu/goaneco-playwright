//go:build e2e

// E2E tests for: TestWheel.java.
package e2e

import (
	"context"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func a2ListenForWheelEvents(t *testing.T, ctx context.Context, page *playwright.Page, selector string) {
	t.Helper()
	a2Ev(t, ctx, page, `selector => {
    document.querySelector(selector).addEventListener('wheel', e => {
      window['lastEvent'] = {
        deltaX: e.deltaX,
        deltaY: e.deltaY,
        clientX: e.clientX,
        clientY: e.clientY,
        deltaMode: e.deltaMode,
        ctrlKey: e.ctrlKey,
        shiftKey: e.shiftKey,
        altKey: e.altKey,
        metaKey: e.metaKey,
      };
    }, { passive: false });
  }`, selector)
}

func a2ExpectWheelEvent(t *testing.T, ctx context.Context, page *playwright.Page, deltaX, deltaY float64, shift bool) {
	t.Helper()
	expected := map[string]any{
		"deltaX": deltaX, "deltaY": deltaY, "clientX": float64(50), "clientY": float64(60), "deltaMode": float64(0),
		"ctrlKey": false, "shiftKey": shift, "altKey": false, "metaKey": false,
	}
	assert.Equal(t, expected, a2Ev(t, ctx, page, "window.lastEvent"))
}

func a2WaitForFunction(t *testing.T, ctx context.Context, page *playwright.Page, expr string) {
	t.Helper()
	if _, err := page.WaitForFunction(ctx, expr, nil); err != nil {
		t.Fatalf("WaitForFunction(%q): %v", expr, err)
	}
}

// Ref: TestWheel.java#shouldDispatchWheelEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWheel.java
func TestWheelShouldDispatchWheelEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div style='width: 5000px; height: 5000px;'></div>"))
	must.NoError(page.Mouse.Move(ctx, 50, 60))
	a2ListenForWheelEvents(t, ctx, page, "div")
	must.NoError(page.Mouse.Wheel(ctx, 0, 100))
	a2WaitForFunction(t, ctx, page, "window.scrollY === 100")
	a2ExpectWheelEvent(t, ctx, page, 0, 100, false)
}

// Ref: TestWheel.java#shouldScrollWhenNobodyIsListening  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWheel.java
func TestWheelShouldScrollWhenNobodyIsListening(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/scrollable.html"))
	must.NoError(page.Mouse.Move(ctx, 50, 60))
	must.NoError(page.Mouse.Wheel(ctx, 0, 100))
	a2WaitForFunction(t, ctx, page, "window.scrollY === 100")
}

// Ref: TestWheel.java#shouldSetTheModifiers  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWheel.java
func TestWheelShouldSetTheModifiers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div style='width: 5000px; height: 5000px;'></div>"))
	must.NoError(page.Mouse.Move(ctx, 50, 60))
	a2ListenForWheelEvents(t, ctx, page, "div")
	must.NoError(page.Keyboard.Down(ctx, "Shift"))
	must.NoError(page.Mouse.Wheel(ctx, 0, 100))
	a2ExpectWheelEvent(t, ctx, page, 0, 100, true)
}

// Ref: TestWheel.java#shouldScrollHorizontally  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWheel.java
func TestWheelShouldScrollHorizontally(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div style='width: 5000px; height: 5000px;'></div>"))
	must.NoError(page.Mouse.Move(ctx, 50, 60))
	a2ListenForWheelEvents(t, ctx, page, "div")
	must.NoError(page.Mouse.Wheel(ctx, 100, 0))
	a2WaitForFunction(t, ctx, page, "window.scrollX === 100")
	a2ExpectWheelEvent(t, ctx, page, 100, 0, false)
}

// Ref: TestWheel.java#shouldWorkWhenTheEventIsCanceled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWheel.java
func TestWheelShouldWorkWhenTheEventIsCanceled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div style='width: 5000px; height: 5000px;'></div>"))
	must.NoError(page.Mouse.Move(ctx, 50, 60))
	a2ListenForWheelEvents(t, ctx, page, "div")
	a2Ev(t, ctx, page, `() => {
    document.querySelector('div').addEventListener('wheel', e => e.preventDefault());
  }`)
	must.NoError(page.Mouse.Wheel(ctx, 0, 100))
	// give the page a chance to scroll
	a2WaitForFunction(t, ctx, page, "!!window['lastEvent']")
	a2ExpectWheelEvent(t, ctx, page, 0, 100, false)
	// ensure that it did not.
	assert.Equal(t, float64(0), a2Ev(t, ctx, page, "window.scrollY"))
}
