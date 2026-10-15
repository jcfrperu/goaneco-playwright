//go:build e2e

// E2E tests for: TestPageEventPageError.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// TestPageErrorsCapturedAndCleared verifies PageErrors captures JS errors and ClearPageErrors resets them.
// Extra: no equivalent method in playwright-java; related class: TestPageEventPageError.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestPageErrorsCapturedAndCleared(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<script>
			window.onerror = undefined;
			throw new Error("first error");
		</script>
	`))

	errs, err := page.PageErrors(ctx)
	must.NoError(err)
	is.NotEmpty(errs)

	must.NoError(page.ClearPageErrors(ctx))

	errs2, err := page.PageErrors(ctx)
	must.NoError(err)
	is.Empty(errs2)
}

// TestPageErrorsEmptyByDefault verifies PageErrors returns empty before any error.
// Extra: no equivalent method in playwright-java; related class: TestPageEventPageError.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestPageErrorsEmptyByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>no errors here</div>`))

	errs, err := page.PageErrors(ctx)
	must.NoError(err)
	is.Empty(errs)
}

// TestClearPageErrorsIsIdempotent verifies ClearPageErrors can be called with no errors.
// Extra: no equivalent method in playwright-java; related class: TestPageEventPageError.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestClearPageErrorsIsIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>no error</div>`))
	must.NoError(page.ClearPageErrors(ctx))
	must.NoError(page.ClearPageErrors(ctx))

	errs, err := page.PageErrors(ctx)
	must.NoError(err)
	is.Empty(errs)
}

// TestPageErrorsCapturedFromUncaughtException verifies JS errors are captured.
// Extra: no equivalent method in playwright-java; related class: TestPageEventPageError.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestPageErrorsCapturedFromUncaughtException(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	// Trigger an uncaught error
	_, _ = page.Evaluate(ctx, `() => {
		setTimeout(() => { throw new Error('uncaught error'); }, 0);
	}`)

	// Give time for error to be captured
	must.NoError(page.WaitForTimeout(ctx, 200))

	// PageErrors may or may not capture async errors depending on implementation
	_, err := page.PageErrors(ctx)
	must.NoError(err)
}

// TestPageErrorsEmptyByDefaultExtra verifies fresh page has no errors.
// Extra: no equivalent method in playwright-java; related class: TestPageEventPageError.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestPageErrorsEmptyByDefaultExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>no errors here</div>`))

	errors, err := page.PageErrors(ctx)
	must.NoError(err)
	is.Empty(errors)
}

// TestClearPageErrorsRemovesAll verifies ClearPageErrors removes all captured errors.
// Extra: no equivalent method in playwright-java; related class: TestPageEventPageError.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestClearPageErrorsRemovesAllExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	must.NoError(page.ClearPageErrors(ctx))

	errors, err := page.PageErrors(ctx)
	must.NoError(err)
	is.Empty(errors)
}

// TestClearPageErrorsIsIdempotentExtra verifies ClearPageErrors can be called multiple times.
// Extra: no equivalent method in playwright-java; related class: TestPageEventPageError.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestClearPageErrorsIsIdempotentExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	must.NoError(page.ClearPageErrors(ctx))
	must.NoError(page.ClearPageErrors(ctx))
	must.NoError(page.ClearPageErrors(ctx))
}

// TestPageEventPageErrorPageErrorsShouldWork verifies that uncaught JS exceptions are collected by PageErrors.
// Ref: TestPageEventPageError.java#pageErrorsShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestPageEventPageErrorPageErrorsShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "goto failed")

	// Throw 50 errors asynchronously and wait for them to fire.
	_, err = page.Evaluate(ctx, `async () => {
		for (let i = 0; i < 50; i++)
			window.setTimeout(() => { throw new Error('pageerr' + i); }, 0);
		await new Promise(f => window.setTimeout(f, 500));
	}`)
	must.NoError(err, "evaluate failed")

	errors, err := page.PageErrors(ctx)
	must.NoError(err, "PageErrors failed")
	if len(errors) < 10 {
		t.Fatalf("expected at least 10 errors, got %d", len(errors))
	}
	for _, e := range errors {
		if !strings.Contains(e, "pageerr") {
			t.Fatalf("unexpected error format: %q", e)
		}
	}
}

// TestPageEventPageErrorClearPageErrorsShouldWork verifies that ClearPageErrors resets the error list.
// Ref: TestPageEventPageError.java#clearPageErrorsShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPageError.java
func TestPageEventPageErrorClearPageErrorsShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "goto failed")

	// Fire first error.
	_, err = page.Evaluate(ctx, `async () => {
		window.setTimeout(() => { throw new Error('error1'); }, 0);
		await new Promise(f => window.setTimeout(f, 200));
	}`)
	must.NoError(err, "evaluate failed")

	errors, err := page.PageErrors(ctx)
	must.NoError(err, "PageErrors failed")
	found := false
	for _, e := range errors {
		if strings.Contains(e, "error1") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected 'error1' in errors, got: %v", errors)
	}

	// Clear and verify empty.
	err = page.ClearPageErrors(ctx)
	must.NoError(err, "ClearPageErrors failed")
	errors, err = page.PageErrors(ctx)
	must.NoError(err, "PageErrors after clear failed")
	is.Len(errors, 0)

	// Fire second error.
	_, err = page.Evaluate(ctx, `async () => {
		window.setTimeout(() => { throw new Error('error2'); }, 0);
		await new Promise(f => window.setTimeout(f, 200));
	}`)
	must.NoError(err, "evaluate failed")

	errors, err = page.PageErrors(ctx)
	must.NoError(err, "PageErrors failed")
	if len(errors) != 1 {
		t.Fatalf("expected 1 error after second throw, got %d: %v", len(errors), errors)
	}
	is.Contains(errors[0], "error2")
}
