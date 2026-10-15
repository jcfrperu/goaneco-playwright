//go:build e2e

// E2E tests for: TestElementHandleWaitForElementState.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandleWaitForElementState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleIsDisabledForDisabledButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>go</button>`))

	handle, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(handle)

	disabled, err := handle.IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleWaitForElementState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el" style="display:none">hi</div>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#el")

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('el').style.display = 'block'`)
	}()

	err = el.WaitForElementState(ctx, "visible")
	must.NoError(err, "WaitForElementState(visible) failed")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleWaitForElementState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForAlreadyVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">visible</div>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#el")

	err = el.WaitForElementState(ctx, "visible")
	must.NoError(err, "WaitForElementState(visible) for already-visible element failed")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleWaitForElementState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="el">visible</div>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#el")

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('el').style.display = 'none'`)
	}()

	err = el.WaitForElementState(ctx, "hidden")
	must.NoError(err, "WaitForElementState(hidden) failed")
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleWaitForElementState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForEnabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn" disabled>Click</button>`)
	must.NoError(err, "SetContent failed")

	el, _ := page.QuerySelector(ctx, "#btn")

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `document.getElementById('btn').disabled = false`)
	}()

	err = el.WaitForElementState(ctx, "enabled")
	must.NoError(err, "WaitForElementState(enabled) failed")
}

// Ref: TestElementHandleWaitForElementState.java#shouldWaitForVisible  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldWaitForVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div style='display:none'>content</div>"))
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(div)
	must.NoError(a3GiveItAChanceToResolve(ctx, page))
	_, err = div.Evaluate(ctx, "div => div.style.display = 'block'")
	must.NoError(err)
	must.NoError(div.WaitForElementState(ctx, "visible"))
}

// Ref: TestElementHandleWaitForElementState.java#shouldWaitForAlreadyVisible  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldWaitForAlreadyVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>content</div>"))
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(div)
	must.NoError(div.WaitForElementState(ctx, "visible"))
}

// Ref: TestElementHandleWaitForElementState.java#shouldTimeoutWaitingForVisible  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldTimeoutWaitingForVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div style='display:none'>content</div>"))
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(div)
	timeout := 100.0
	err = div.WaitForElementState(ctx, "visible", &playwright.ElementHandleWaitForElementStateOptions{Timeout: &timeout})
	must.Error(err)
	assert.Contains(t, err.Error(), "Timeout")
}

// Ref: TestElementHandleWaitForElementState.java#shouldThrowWaitingForVisibleWhenDetached  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldThrowWaitingForVisibleWhenDetached(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div style='display:none'>content</div>"))
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(div)
	_, err = div.Evaluate(ctx, "div => div.remove()")
	must.NoError(err)
	err = div.WaitForElementState(ctx, "visible")
	must.Error(err)
	assert.Contains(t, err.Error(), "Element is not attached to the DOM")
}

// Ref: TestElementHandleWaitForElementState.java#shouldWaitForHidden  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldWaitForHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>content</div>"))
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(div)
	must.NoError(a3GiveItAChanceToResolve(ctx, page))
	_, err = div.Evaluate(ctx, "div => div.style.display = 'none'")
	must.NoError(err)
	must.NoError(div.WaitForElementState(ctx, "hidden"))
}

// Ref: TestElementHandleWaitForElementState.java#shouldWaitForAlreadyHidden  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldWaitForAlreadyHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div></div>"))
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(div)
	must.NoError(div.WaitForElementState(ctx, "hidden"))
}

// Ref: TestElementHandleWaitForElementState.java#shouldWaitForHiddenWhenDetached  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldWaitForHiddenWhenDetached(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<div>content</div>"))
	div, err := page.QuerySelector(ctx, "div")
	must.NoError(err)
	must.NotNil(div)
	must.NoError(a3GiveItAChanceToResolve(ctx, page))
	_, err = div.Evaluate(ctx, "div => div.remove()")
	must.NoError(err)
	must.NoError(div.WaitForElementState(ctx, "hidden"))
}

// Ref: TestElementHandleWaitForElementState.java#shouldWaitForEnabledButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldWaitForEnabledButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<button disabled><span>Target</span></button>"))
	span, err := page.QuerySelector(ctx, "text=Target")
	must.NoError(err)
	must.NotNil(span)
	must.NoError(a3GiveItAChanceToResolve(ctx, page))
	_, err = span.Evaluate(ctx, "span => span.parentElement.disabled = false")
	must.NoError(err)
	must.NoError(span.WaitForElementState(ctx, "enabled"))
}

// Ref: TestElementHandleWaitForElementState.java#shouldThrowWaitingForEnabledWhenDetached  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldThrowWaitingForEnabledWhenDetached(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<button disabled>Target</button>"))
	button, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(button)
	_, err = button.Evaluate(ctx, "button => button.remove()")
	must.NoError(err)
	err = button.WaitForElementState(ctx, "enabled")
	must.Error(err)
	assert.Contains(t, err.Error(), "Element is not attached to the DOM")
}

// Ref: TestElementHandleWaitForElementState.java#shouldWaitForDisabledButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleWaitForElementState.java
func TestElementHandleWaitForElementStateShouldWaitForDisabledButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<button><span>Target</span></button>"))
	span, err := page.QuerySelector(ctx, "text=Target")
	must.NoError(err)
	must.NotNil(span)
	must.NoError(a3GiveItAChanceToResolve(ctx, page))
	_, err = span.Evaluate(ctx, "span => span.parentElement.disabled = true")
	must.NoError(err)
	must.NoError(span.WaitForElementState(ctx, "disabled"))
}
