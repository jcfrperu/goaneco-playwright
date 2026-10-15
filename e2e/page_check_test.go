//go:build e2e

// E2E tests for: TestPageCheck.java.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestPageCheck.java#shouldCheckTheLabelWithPosition  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageCheck.java
func TestPageCheckShouldCheckTheLabelWithPosition(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	content := "<input id='checkbox' type='checkbox' style='width: 5px; height: 5px;'>\n" +
		"    <label for='checkbox'>\n" +
		"      <a href=" + srv.EmptyPage() + ">I am a long link that goes away so that nothing good will happen if you click on me</a>\n" +
		"      Click me\n" +
		"    </label>"
	must.NoError(page.SetContent(ctx, content))

	loc := page.Locator("text=Click me")
	box, err := loc.BoundingBox(ctx)
	must.NoError(err)
	must.NotNil(box)

	must.NoError(page.Check(ctx, "text=Click me", &playwright.PageCheckOptions{
		Position: &playwright.Point{X: box.Width - 10, Y: 2},
	}))

	checked, err := page.EvalOnSelector(ctx, "input", "input => input.checked")
	must.NoError(err)
	is.Equal(true, checked)
}

// Ref: TestPageCheck.java#trialRunShouldNotCheck  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageCheck.java
func TestPageCheckTrialRunShouldNotCheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	trial := true

	must.NoError(page.SetContent(ctx, "<input id='checkbox' type='checkbox' />"))
	must.NoError(page.Check(ctx, "#checkbox", &playwright.PageCheckOptions{Trial: &trial}))
	res, err := page.Evaluate(ctx, "() => document.getElementById('checkbox').checked")
	must.NoError(err)
	assert.Equal(t, false, res)
}

// Ref: TestPageCheck.java#trialRunShouldNotUncheck  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageCheck.java
func TestPageCheckTrialRunShouldNotUncheck(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	trial := true

	must.NoError(page.SetContent(ctx, "<input id='checkbox' type='checkbox' checked />"))
	must.NoError(page.Uncheck(ctx, "#checkbox", &playwright.PageCheckOptions{Trial: &trial}))
	res, err := page.Evaluate(ctx, "() => document.getElementById('checkbox').checked")
	must.NoError(err)
	assert.Equal(t, true, res)
}

// Ref: TestPageCheck.java#shouldCheckTheBoxUsingSetChecked  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageCheck.java
func TestPageCheckShouldCheckTheBoxUsingSetChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<input id='checkbox' type='checkbox'></input>"))
	// Page.SetChecked is not available; the equivalent Locator.SetChecked is used.
	must.NoError(page.Locator("input").SetChecked(ctx, true))
	res, err := page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	assert.Equal(t, true, res)
	must.NoError(page.Locator("input").SetChecked(ctx, false))
	res, err = page.Evaluate(ctx, "() => window['checkbox'].checked")
	must.NoError(err)
	assert.Equal(t, false, res)
}
