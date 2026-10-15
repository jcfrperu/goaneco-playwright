//go:build e2e

// E2E tests for: TestDialog.java.
package e2e

import (
	"sync/atomic"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestDialog.java#shouldFire  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDialog.java
func TestDialogShouldFire(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	var fired atomic.Bool
	page.OnDialog(func(d *playwright.Dialog) {
		assert.Equal(t, "alert", d.Type())
		assert.Equal(t, "", d.DefaultValue())
		assert.Equal(t, "yo", d.Message())
		fired.Store(true)
		assert.NoError(t, d.Accept(ctx))
	})
	a2Ev(t, ctx, page, "() => alert('yo')")
	assert.True(t, fired.Load())
}

// Ref: TestDialog.java#shouldAllowAcceptingPrompts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDialog.java
func TestDialogShouldAllowAcceptingPrompts(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	page.OnDialog(func(d *playwright.Dialog) {
		assert.Equal(t, "prompt", d.Type())
		assert.Equal(t, "yes.", d.DefaultValue())
		assert.Equal(t, "question?", d.Message())
		assert.NoError(t, d.Accept(ctx, "answer!"))
	})
	assert.Equal(t, "answer!", a2Ev(t, ctx, page, "() => prompt('question?', 'yes.')"))
}

// Ref: TestDialog.java#shouldAcceptPromptOnce  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDialog.java
func TestDialogShouldAcceptPromptOnce(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	var calls atomic.Int32
	var off func()
	off = page.OnDialog(func(d *playwright.Dialog) {
		off() // emulates onceDialog: unsubscribe before the first dialog is handled
		calls.Add(1)
		assert.Equal(t, "prompt", d.Type())
		assert.Equal(t, "question?", d.Message())
		assert.NoError(t, d.Accept(ctx, "answer!"))
	})
	assert.Equal(t, "answer!", a2Ev(t, ctx, page, "prompt('question?')"))
	assert.Nil(t, a2Ev(t, ctx, page, "prompt('question?')"))
	assert.Equal(t, int32(1), calls.Load())
}

// Ref: TestDialog.java#shouldDismissThePrompt  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDialog.java
func TestDialogShouldDismissThePrompt(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	page.OnDialog(func(d *playwright.Dialog) {
		assert.NoError(t, d.Dismiss(ctx))
	})
	assert.Nil(t, a2Ev(t, ctx, page, "() => prompt('question?')"))
}

// Ref: TestDialog.java#shouldAcceptTheConfirmPrompt  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDialog.java
func TestDialogShouldAcceptTheConfirmPrompt(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	page.OnDialog(func(d *playwright.Dialog) {
		assert.NoError(t, d.Accept(ctx))
	})
	assert.Equal(t, true, a2Ev(t, ctx, page, "() => confirm('boolean?')"))
}

// Ref: TestDialog.java#shouldDismissTheConfirmPrompt  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestDialog.java
func TestDialogShouldDismissTheConfirmPrompt(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	page.OnDialog(func(d *playwright.Dialog) {
		assert.NoError(t, d.Dismiss(ctx))
	})
	res, err := page.Evaluate(ctx, "() => confirm('boolean?')")
	must.NoError(err)
	assert.Equal(t, false, res)
}
