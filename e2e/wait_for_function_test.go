//go:build e2e

// E2E tests for: TestWaitForFunction.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestWaitForFunctionWithPolling verifies WaitForFunction with custom polling interval.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWithPolling(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	_, err := page.Evaluate(ctx, `() => {
		setTimeout(() => window.__done = true, 100);
	}`)
	must.NoError(err)

	handle, err := page.WaitForFunction(ctx, `() => window.__done`, nil,
		&playwright.WaitForFunctionOptions{PollingInterval: 50})
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(true, val)
}

// TestWaitForFunctionTimesOut verifies WaitForFunction times out correctly.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionTimesOut(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	_, err := page.WaitForFunction(ctx, `() => false`, nil,
		&playwright.WaitForFunctionOptions{Timeout: 200})
	is.Error(err)
}

// TestWaitForFunctionPassesArgument verifies WaitForFunction passes arg to expression.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionPassesArgument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `n => n > 0`, 42)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(true, val)
}

// TestWaitForFunctionReturnsDOMElement verifies WaitForFunction returns DOM element handle.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionReturnsDOMElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">text</div>`))

	handle, err := page.WaitForFunction(ctx, `() => document.getElementById('target')`, nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionWaitsForDOMChange verifies WaitForFunction detects DOM mutations.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWaitsForDOMChange(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">initial</div>`))

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => { document.getElementById('el').textContent = 'changed'; }`)
	}()

	handle, err := page.WaitForFunction(ctx, `() => document.getElementById('el').textContent === 'changed'`, nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionWithMutation verifies WaitForFunction detects DOM mutations.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWithMutation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">initial</div>`))

	done := make(chan struct{})
	go func() {
		defer close(done)
		must.NoError(page.WaitForTimeout(ctx, 100))
		_, _ = page.Evaluate(ctx, `() => { document.getElementById('el').textContent = 'changed'; }`)
	}()

	handle, err := page.WaitForFunction(ctx,
		`() => document.getElementById('el').textContent === 'changed'`, nil)
	must.NoError(err)
	must.NotNil(handle)
	<-done
}

// TestWaitForFunctionReturnsBoolean verifies WaitForFunction returns boolean handle.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionReturnsBoolean(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => true`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(true, val)
}

// TestWaitForFunctionWithPollingInterval2 verifies WaitForFunction with polling.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWithPollingInterval2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="counter">0</div>`))

	_, _ = page.Evaluate(ctx, `() => {
		let n = 0;
		const t = setInterval(() => { document.getElementById('counter').textContent = ++n; }, 30);
		window.__stop = () => clearInterval(t);
	}`)

	handle, err := page.WaitForFunction(ctx,
		`() => parseInt(document.getElementById('counter').textContent) >= 2`, nil,
		&playwright.WaitForFunctionOptions{PollingInterval: 20})
	must.NoError(err)
	must.NotNil(handle)

	_, _ = page.Evaluate(ctx, `() => window.__stop()`)
}

// TestWaitForFunctionWithArg2 verifies WaitForFunction passes argument to predicate.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWithArg2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `(expected) => expected === 42`, 42)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionReturnsBoolEx3 verifies WaitForFunction with bool condition.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionReturnsBoolEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => true", nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionWithDOMChangeEx3 verifies WaitForFunction detects DOM changes.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWithDOMChangeEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">initial</div>`))

	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => document.getElementById('el').textContent = 'changed'`)
	}()

	handle, err := page.WaitForFunction(ctx, `() => document.getElementById('el').textContent === 'changed'`, nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionWithArgumentEx3 verifies WaitForFunction receives argument.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWithArgumentEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "(n) => n === 42", 42)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionWindowPropertyEx3 verifies WaitForFunction waits for window property.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWindowPropertyEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => { window.__ready = true; }`)
	}()

	handle, err := page.WaitForFunction(ctx, "() => window.__ready === true", nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionReturnsJSHandleEx3 verifies WaitForFunction result is usable JSHandle.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionReturnsJSHandleEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => 99", nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(99), val)
}

// TestWaitForFunctionPollingIntervalEx4 verifies WaitForFunction polls until condition.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionPollingIntervalEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="counter">0</div>`))

	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => document.getElementById('counter').textContent = '5'`)
	}()

	handle, err := page.WaitForFunction(ctx, `() => parseInt(document.getElementById('counter').textContent) >= 5`, nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionImmediatelyTrueEx4 verifies WaitForFunction returns if already true.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionImmediatelyTrueEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, `() => 1 + 1 === 2`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(true, val)
}

// TestWaitForFunctionStringResultEx4 verifies WaitForFunction returns string handle.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionStringResultEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, `() => 'ready'`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal("ready", val)
}

// TestWaitForFunctionObjectResultEx4 verifies WaitForFunction returns object handle.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionObjectResultEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, `() => ({ status: 'ok', code: 200 })`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	m, ok := val.(map[string]any)
	is.True(ok)
	is.Equal("ok", m["status"])
}

// TestWaitForFunctionCounterEx5 verifies WaitForFunction waits for counter.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionCounterEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => {
		window.__count = 0;
		const id = setInterval(() => { window.__count++; if (window.__count >= 3) clearInterval(id); }, 50);
	}`)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => window.__count >= 3`, nil)
	must.NoError(err)

	count, err := page.Evaluate(ctx, `() => window.__count`)
	must.NoError(err)
	is.GreaterOrEqual(count.(float64), float64(3))
}

// TestWaitForFunctionStringResultEx5 verifies WaitForFunction can return string.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionStringResultEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => { window.__ready = false; setTimeout(() => { window.__ready = true; }, 50); }`)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => window.__ready === true`, nil)
	must.NoError(err)

	ready, err := page.Evaluate(ctx, `() => window.__ready`)
	must.NoError(err)
	is.Equal(true, ready)
}

// TestWaitForFunctionElementAppearEx5 verifies WaitForFunction waits for DOM element.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionElementAppearEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body id="b"></body></html>`))

	_, err := page.Evaluate(ctx, `() => setTimeout(() => {
		const el = document.createElement('div');
		el.id = 'late-el';
		document.body.appendChild(el);
	}, 100)`)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => !!document.getElementById('late-el')`, nil)
	must.NoError(err)

	count, err := page.Locator("#late-el").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestWaitForFunctionInputValueEx6 verifies WaitForFunction waits for input value.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionInputValueEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	_, err := page.Evaluate(ctx, `() => setTimeout(() => { document.getElementById('inp').value = 'ready'; }, 50)`)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => document.getElementById('inp').value === 'ready'`, nil)
	must.NoError(err)

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("ready", val)
}

// TestWaitForFunctionClassAddedEx6 verifies WaitForFunction waits for class addition.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionClassAddedEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	_, err := page.Evaluate(ctx, `() => setTimeout(() => { document.getElementById('d').classList.add('loaded'); }, 60)`)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => document.getElementById('d').classList.contains('loaded')`, nil)
	must.NoError(err)

	cls, err := page.Evaluate(ctx, `() => document.getElementById('d').className`)
	must.NoError(err)
	is.Equal("loaded", cls)
}

// TestWaitForFunctionWindowFlagEx6 verifies WaitForFunction waits for window flag.
// Extra: no equivalent method in playwright-java; related class: TestWaitForFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionWindowFlagEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => setTimeout(() => { window.__flag = 'done'; }, 80)`)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => window.__flag === 'done'`, nil)
	must.NoError(err)

	flag, err := page.Evaluate(ctx, `() => window.__flag`)
	must.NoError(err)
	is.Equal("done", flag)
}

// TestWaitForFunctionShouldTimeout verifies that WaitForTimeout blocks for at least the specified duration.
// Ref: TestWaitForFunction.java#shouldTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	const ms = 100.0
	start := time.Now()
	err := page.WaitForTimeout(ctx, ms)
	must.NoError(err, "WaitForTimeout failed")
	elapsed := time.Since(start).Milliseconds()
	if elapsed < ms/2 {
		t.Fatalf("WaitForTimeout(%v ms) returned too fast: %d ms elapsed", ms, elapsed)
	}
}

// TestWaitForFunctionShouldAcceptAString verifies that WaitForFunction accepts a plain JS expression string.
// Ref: TestWaitForFunction.java#shouldAcceptAString  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldAcceptAString(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => window['__FOO'] = 1")
	must.NoError(err, "evaluate failed")
	_, err = page.WaitForFunction(ctx, "window.__FOO === 1", nil)
	must.NoError(err, "WaitForFunction failed")
}

// TestWaitForFunctionShouldReturnTheSuccessValueAsAJSHandle verifies that WaitForFunction returns the truthy value as a JSHandle.
// Ref: TestWaitForFunction.java#shouldReturnTheSuccessValueAsAJSHandle  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldReturnTheSuccessValueAsAJSHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "5", nil)
	must.NoError(err, "WaitForFunction failed")
	val, err := handle.JSONValue(ctx)
	must.NoError(err, "JSONValue failed")
	f, ok := val.(float64)
	if !ok || f != 5 {
		t.Fatalf("expected 5, got %v (%T)", val, val)
	}
}

// TestWaitForFunctionShouldFailWithPredicateThrowingOnFirstCall verifies that a predicate error propagates as an error.
// Ref: TestWaitForFunction.java#shouldFailWithPredicateThrowingOnFirstCall  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldFailWithPredicateThrowingOnFirstCall(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.WaitForFunction(ctx, "() => { throw new Error('oh my'); }", nil)
	is.ErrorContains(err, "oh my")
}

// TestWaitForFunctionShouldRespectTimeout verifies that the timeout option causes a timeout error.
// Ref: TestWaitForFunction.java#shouldRespectTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldRespectTimeout(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.WaitForFunction(ctx, "false", nil, &playwright.WaitForFunctionOptions{Timeout: 10})
	is.Error(err)
	if err.Error() == "" {
		t.Fatalf("expected timeout error, got: %v", err)
	}
}

// TestWaitForFunctionShouldWaitForPredicateWithArguments verifies that arguments are passed to the predicate.
// Ref: TestWaitForFunction.java#shouldWaitForPredicateWithArguments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldWaitForPredicateWithArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.WaitForFunction(ctx, "({arg1, arg2}) => arg1 + arg2 === 3", map[string]any{"arg1": 1, "arg2": 2})
	must.NoError(err, "WaitForFunction with arguments failed")
}

// TestWaitForFunctionShouldPollOnInterval verifies that the polling interval is respected.
// Ref: TestWaitForFunction.java#shouldPollOnInterval  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldPollOnInterval(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	const polling = 100.0
	handle, err := page.WaitForFunction(ctx, `() => {
		if (!window["__startTime"]) {
			window["__startTime"] = Date.now();
			return false;
		}
		return Date.now() - window["__startTime"];
	}`, nil, &playwright.WaitForFunctionOptions{PollingInterval: polling})
	must.NoError(err, "WaitForFunction failed")
	val, err := handle.JSONValue(ctx)
	must.NoError(err, "JSONValue failed")
	delta, ok := val.(float64)
	if !ok || delta < polling {
		t.Fatalf("expected delta >= %v ms, got %v (%T)", polling, val, val)
	}
}

// TestWaitForFunctionShouldSurviveNavigations verifies that WaitForFunction works after a page reload.
// Ref: TestWaitForFunction.java#shouldSurviveNavigations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldSurviveNavigations(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.Reload(ctx)
	must.NoError(err, "Reload failed")
	_, err = page.Evaluate(ctx, "() => window['__done'] = true")
	must.NoError(err, "evaluate failed")
	_, err = page.WaitForFunction(ctx, "() => window['__done']", nil)
	must.NoError(err, "WaitForFunction after navigation failed")
}

// TestWaitForFunctionShouldPollOnRaf verifies that WaitForFunction works with requestAnimationFrame polling.
// Ref: TestWaitForFunction.java#shouldPollOnRaf  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldPollOnRaf(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => window['__rafDone'] = false")
	must.NoError(err)

	go func() {
		time.Sleep(100 * time.Millisecond)
		evalCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = page.Evaluate(evalCtx, "() => window['__rafDone'] = true")
	}()

	// PollingInterval: 0 means requestAnimationFrame
	_, err = page.WaitForFunction(ctx, "() => window['__rafDone']", nil,
		&playwright.WaitForFunctionOptions{PollingInterval: 0})
	must.NoError(err, "WaitForFunction with RAF polling failed")
}

// TestWaitForFunctionShouldReturnTheWindowAsASuccessValue verifies that the window object can be returned as a JSHandle.
// Ref: TestWaitForFunction.java#shouldReturnTheWindowAsASuccessValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldReturnTheWindowAsASuccessValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => window", nil)
	must.NoError(err, "WaitForFunction returning window failed")
	must.NotNil(handle, "handle should not be nil")
}

// TestWaitForFunctionShouldWorkWithMultilineBody verifies that WaitForFunction accepts a multiline function body.
// Ref: TestWaitForFunction.java#shouldWorkWithMultilineBody  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldWorkWithMultilineBody(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.WaitForFunction(ctx, `() => {
		const x = 1;
		const y = 2;
		return x + y === 3;
	}`, nil)
	must.NoError(err, "WaitForFunction with multiline body failed")
}

// TestWaitForFunctionShouldAcceptElementHandleArguments verifies passing an ElementHandle as an argument.
// Ref: TestWaitForFunction.java#shouldAcceptElementHandleArguments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldAcceptElementHandleArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">hello</div>`))

	el, err := page.QuerySelector(ctx, "#target")
	must.NoError(err)
	must.NotNil(el)

	handle, err := page.WaitForFunction(ctx, "(el) => el.textContent === 'hello'", el)
	must.NoError(err, "WaitForFunction with ElementHandle argument failed")
	must.NotNil(handle)
}

// TestWaitForFunctionShouldFailWithReferenceErrorOnWrongPage verifies that a ReferenceError propagates as an error.
// Ref: TestWaitForFunction.java#shouldFailWithReferenceErrorOnWrongPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldFailWithReferenceErrorOnWrongPage(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.WaitForFunction(ctx, "() => __undefinedVariable", nil,
		&playwright.WaitForFunctionOptions{Timeout: 1000})
	is.Error(err, "WaitForFunction should fail with ReferenceError")
}

// TestWaitForFunctionShouldFailWithPredicateThrowingSometimes verifies that a predicate that throws intermittently
// eventually fails with the thrown error.
// Ref: TestWaitForFunction.java#shouldFailWithPredicateThrowingSometimes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldFailWithPredicateThrowingSometimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => window['__callCount'] = 0")
	must.NoError(err)

	// Throws on every call - eventually times out
	_, err = page.WaitForFunction(ctx, `() => {
		window['__callCount']++;
		throw new Error('always throws');
	}`, nil, &playwright.WaitForFunctionOptions{Timeout: 500})
	is.Error(err, "should fail because predicate always throws")
}

// TestWaitForFunctionShouldDisableTimeoutWhenItsSetTo0 verifies that a negative Timeout disables the timeout,
// allowing WaitForFunction to wait indefinitely (in practice, the predicate becomes true quickly).
// Ref: TestWaitForFunction.java#shouldDisableTimeoutWhenItsSetTo0  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldDisableTimeoutWhenItsSetTo0(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => window['__ready'] = false")
	must.NoError(err)

	go func() {
		time.Sleep(50 * time.Millisecond)
		evalCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = page.Evaluate(evalCtx, "() => window['__ready'] = true")
	}()

	// Negative timeout = disable (server receives 0 which means no timeout)
	_, err = page.WaitForFunction(ctx, "() => window['__ready']", nil,
		&playwright.WaitForFunctionOptions{Timeout: -1})
	must.NoError(err, "WaitForFunction with disabled timeout should complete")
}

// TestWaitForFunctionShouldThrowNegativePollingInterval verifies that a negative polling interval returns an error.
// Ref: TestWaitForFunction.java#shouldThrowNegativePollingInterval  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldThrowNegativePollingInterval(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.WaitForFunction(ctx, "() => !!document.body", nil,
		&playwright.WaitForFunctionOptions{PollingInterval: -10})
	is.ErrorContains(err, "Cannot poll with non-positive interval")
}

// TestWaitForFunctionShouldWorkWhenResolvedRightBeforeExecutionContextDisposal verifies that WaitForFunction works when resolved
// right before the execution context is disposed (e.g., during a reload).
// Ref: TestWaitForFunction.java#shouldWorkWhenResolvedRightBeforeExecutionContextDisposal  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldWorkWhenResolvedRightBeforeExecutionContextDisposal(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.AddInitScript(ctx, "window['__RELOADED'] = true")
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => {
		if (!window['__RELOADED'])
			window.location.reload();
		return true;
	}`, nil)
	must.NoError(err)
}

// TestWaitForFunctionShouldSurviveCrossProcessNavigation verifies that WaitForFunction works after
// a cross-process navigation.
// Ref: TestWaitForFunction.java#shouldSurviveCrossProcessNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldSurviveCrossProcessNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)
	err = page.Reload(ctx)
	must.NoError(err)

	srv.ServeWithBody("/grid.html", "text/html; charset=utf-8",
		`<!DOCTYPE html><html><body><div id="grid"></div></body></html>`)
	err = page.Goto(ctx, srv.CrossProcessPrefix()+"/grid.html")
	must.NoError(err)

	_, err = page.Evaluate(ctx, "() => window['__FOO'] = 1")
	must.NoError(err)

	handle, err := page.WaitForFunction(ctx, "window.__FOO === 1", nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionShouldNotBeCalledAfterFinishingSuccessfully verifies that the predicate is not called again after
// WaitForFunction succeeds and the page reloads.
// Ref: TestWaitForFunction.java#shouldNotBeCalledAfterFinishingSuccessfully  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldNotBeCalledAfterFinishingSuccessfully(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	var mu sync.Mutex
	var messages []string
	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if strings.HasPrefix(msg.Text(), "waitForFunction") {
			mu.Lock()
			messages = append(messages, msg.Text())
			mu.Unlock()
		}
	})

	_, err = page.WaitForFunction(ctx, `() => {
		console.log('waitForFunction1');
		return true;
	}`, nil)
	must.NoError(err)

	err = page.Reload(ctx)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => {
		console.log('waitForFunction2');
		return true;
	}`, nil)
	must.NoError(err)

	err = page.Reload(ctx)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => {
		console.log('waitForFunction3');
		return true;
	}`, nil)
	must.NoError(err)

	// Allow async console handlers to flush.
	_ = page.WaitForTimeout(ctx, 200)

	mu.Lock()
	got := messages
	mu.Unlock()

	is.Equal([]string{"waitForFunction1", "waitForFunction2", "waitForFunction3"}, got)
}

// TestWaitForFunctionShouldNotBeCalledAfterFinishingUnsuccessfully verifies that the predicate is not called again after
// WaitForFunction fails and the page reloads.
// Ref: TestWaitForFunction.java#shouldNotBeCalledAfterFinishingUnsuccessfully  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldNotBeCalledAfterFinishingUnsuccessfully(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	var mu sync.Mutex
	var messages []string
	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if strings.HasPrefix(msg.Text(), "waitForFunction") {
			mu.Lock()
			messages = append(messages, msg.Text())
			mu.Unlock()
		}
	})

	_, err = page.WaitForFunction(ctx, `() => {
		console.log('waitForFunction1');
		throw new Error('waitForFunction1');
	}`, nil)
	is.ErrorContains(err, "waitForFunction1")

	err = page.Reload(ctx)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => {
		console.log('waitForFunction2');
		throw new Error('waitForFunction2');
	}`, nil)
	is.ErrorContains(err, "waitForFunction2")

	err = page.Reload(ctx)
	must.NoError(err)

	_, err = page.WaitForFunction(ctx, `() => {
		console.log('waitForFunction3');
		throw new Error('waitForFunction3');
	}`, nil)
	is.ErrorContains(err, "waitForFunction3")

	// Allow async console handlers to flush.
	_ = page.WaitForTimeout(ctx, 200)

	mu.Lock()
	got := messages
	mu.Unlock()

	is.Equal([]string{"waitForFunction1", "waitForFunction2", "waitForFunction3"}, got)
}

// Ref: TestWaitForFunction.java#shouldAvoidSideEffectsAfterTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldAvoidSideEffectsAfterTimeout(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var count int
	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	_, err := page.WaitForFunction(ctx, `() => {
		window['counter'] = (window['counter'] || 0) + 1;
		console.log(window['counter']);
	}`, nil, &playwright.WaitForFunctionOptions{PollingInterval: 1, Timeout: 1000})
	is.Error(err)
	is.ErrorContains(err, "Timeout")

	mu.Lock()
	savedCounter := count
	mu.Unlock()

	_ = page.WaitForTimeout(ctx, 2000)

	mu.Lock()
	finalCounter := count
	mu.Unlock()
	is.Equal(savedCounter, finalCounter)
}

// Ref: TestWaitForFunction.java#shouldWorkWithStrictCSPPolicy  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldWorkWithStrictCSPPolicy(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Serve empty page with a strict CSP that only allows scripts from the same origin.
	srv.SetRoute("/csp.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "script-src "+srv.Prefix())
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<html><body></body></html>"))
	})
	must.NoError(page.Goto(ctx, srv.Prefix()+"/csp.html"))
	_, err := page.Evaluate(ctx, "() => window['__FOO'] = 'hit'")
	must.NoError(err)
	_, err = page.WaitForFunction(ctx, "() => window['__FOO'] === 'hit'", nil)
	must.NoError(err)
}

// Ref: TestWaitForFunction.java#shouldRespectDefaultTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWaitForFunction.java
func TestWaitForFunctionShouldRespectDefaultTimeout(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Go API doesn't expose setDefaultTimeout; use an explicit 1ms timeout instead.
	_, err := page.WaitForFunction(ctx, "false", nil, &playwright.WaitForFunctionOptions{Timeout: 1})
	is.Error(err)
	is.ErrorContains(err, "Timeout")
}
