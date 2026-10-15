//go:build e2e

// E2E tests for: TestElementHandleClick.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestElementHandleClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickFiresClick(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="this.dataset.clicked='yes'">click</button>
	`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Click(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('btn').dataset.clicked`)
	must.NoError(err)
	is.Equal("yes", result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleIsHiddenForHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="hidden" style="display:none">hidden</div>`))

	handle, err := page.QuerySelector(ctx, "#hidden")
	must.NoError(err)
	must.NotNil(handle)

	hidden, err := handle.IsHidden(ctx)
	must.NoError(err)
	is.True(hidden)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked=true">Click</button>
	`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Click(ctx))

	clicked, err := page.Evaluate(ctx, `() => window.__clicked`)
	must.NoError(err)
	is.Equal(true, clicked)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.__clicked9=true">Click</button>
	`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NoError(el.Click(ctx))

	result, err := page.Evaluate(ctx, `() => window.__clicked9`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" onclick="this.textContent='clicked'">Click</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.Click(ctx))

	text, err := page.Locator("#btn").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("clicked", *text)
}

// Extra: no equivalent method in playwright-java; related class: TestElementHandleClick.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleDispatchClickEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" onclick="this.textContent='dispatched'">Click</button>`))

	el, err := page.QuerySelector(ctx, "#btn")
	must.NoError(err)
	must.NotNil(el)

	must.NoError(el.DispatchEvent(ctx, "click"))

	text, err := page.Locator("#btn").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("dispatched", *text)
}

// Ref: TestElementHandleClick.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))
	button, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(button)
	must.NoError(button.Click(ctx))
	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Clicked", result)
}

// Ref: TestElementHandleClick.java#shouldWorkWithNodeRemoved  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldWorkWithNodeRemoved(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))
	_, err := page.Evaluate(ctx, "() => delete window['Node']")
	must.NoError(err)
	button, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(button)
	must.NoError(button.Click(ctx))
	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Clicked", result)
}

// Ref: TestElementHandleClick.java#shouldWorkForShadowDOMV1  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldWorkForShadowDOMV1(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/shadow.html"))
	handle, err := page.EvaluateHandle(ctx, "() => window['button']")
	must.NoError(err)
	buttonHandle := handle.AsElement()
	must.NotNil(buttonHandle)
	must.NoError(buttonHandle.Click(ctx))
	clicked, err := page.Evaluate(ctx, "clicked")
	must.NoError(err)
	is.Equal(true, clicked)
}

// Ref: TestElementHandleClick.java#shouldWorkForTextNodes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldWorkForTextNodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	handle, err := page.EvaluateHandle(ctx, "() => document.querySelector('button').firstChild")
	must.NoError(err)
	buttonTextNode := handle.AsElement()
	must.NotNil(buttonTextNode)
	must.NoError(buttonTextNode.Click(ctx))
	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Clicked", result)
}

// Ref: TestElementHandleClick.java#shouldThrowForDetachedNodes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldThrowForDetachedNodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))
	button, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(button)
	_, err = page.Evaluate(ctx, "button => button.remove()", button)
	must.NoError(err)
	err = button.Click(ctx)
	must.Error(err)
	assert.Contains(t, err.Error(), "Element is not attached to the DOM")
}

// Ref: TestElementHandleClick.java#shouldThrowForHiddenNodesWithForce  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldThrowForHiddenNodesWithForce(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	button, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(button)
	_, err = page.Evaluate(ctx, "button => button.style.display = 'none'", button)
	must.NoError(err)
	force := true
	err = button.Click(ctx, &playwright.ElementHandleClickOptions{Force: &force})
	must.Error(err)
	assert.Contains(t, err.Error(), "Element is not visible")
}

// Ref: TestElementHandleClick.java#shouldThrowForRecursivelyHiddenNodesWithForce  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldThrowForRecursivelyHiddenNodesWithForce(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/button.html"))
	button, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(button)
	_, err = page.Evaluate(ctx, "button => button.parentElement.style.display = 'none'", button)
	must.NoError(err)
	force := true
	err = button.Click(ctx, &playwright.ElementHandleClickOptions{Force: &force})
	must.Error(err)
	assert.Contains(t, err.Error(), "Element is not visible")
}

// Ref: TestElementHandleClick.java#shouldThrowForBrElementsWithForce  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldThrowForBrElementsWithForce(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "hello<br>goodbye"))
	br, err := page.QuerySelector(ctx, "br")
	must.NoError(err)
	must.NotNil(br)
	force := true
	err = br.Click(ctx, &playwright.ElementHandleClickOptions{Force: &force})
	must.Error(err)
	assert.Contains(t, err.Error(), "Element is outside of the viewport")
}

// Ref: TestElementHandleClick.java#shouldDoubleClickTheButton  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestElementHandleClick.java
func TestElementHandleClickShouldDoubleClickTheButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))
	_, err := page.Evaluate(ctx, `() => {
  window['double'] = false;
  const button = document.querySelector('button');
  button.addEventListener('dblclick', event => {
    window['double'] = true;
  });
}`)
	must.NoError(err)
	button, err := page.QuerySelector(ctx, "button")
	must.NoError(err)
	must.NotNil(button)
	must.NoError(button.DblClick(ctx))
	dbl, err := page.Evaluate(ctx, "() => window['double']")
	must.NoError(err)
	is.Equal(true, dbl)
	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Clicked", result)
}
