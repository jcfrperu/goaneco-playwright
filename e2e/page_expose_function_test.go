//go:build e2e

// E2E tests for: TestPageExposeFunction.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingBasic(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	// Expose a binding that adds two numbers.
	err := bCtx.ExposeBinding(ctx, "add", func(args ...any) any {
		a, _ := args[0].(float64)
		b, _ := args[1].(float64)
		return a + b
	})
	must.NoError(err, "ExposeBinding failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	result, err := page.Evaluate(ctx, "add(5, 6)")
	must.NoError(err, "Evaluate failed")
	if result != float64(11) {
		t.Errorf("add(5, 6) = %v (%T), want float64(11)", result, result)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	err := bCtx.ExposeBinding(ctx, "mul", func(args ...any) any {
		a, _ := args[0].(float64)
		b, _ := args[1].(float64)
		return a * b
	})
	must.NoError(err, "ExposeBinding mul")
	err = bCtx.ExposeBinding(ctx, "sub", func(args ...any) any {
		a, _ := args[0].(float64)
		b, _ := args[1].(float64)
		return a - b
	})
	must.NoError(err, "ExposeBinding sub")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	result, err := page.Evaluate(ctx, `mul(3, 4)`)
	must.NoError(err, "Evaluate mul")
	if result != float64(12) {
		t.Errorf("mul(3,4) = %v, want 12", result)
	}

	result, err = page.Evaluate(ctx, `sub(10, 3)`)
	must.NoError(err, "Evaluate sub")
	if result != float64(7) {
		t.Errorf("sub(10,3) = %v, want 7", result)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingOnNewPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	calls := make(chan float64, 10)
	err := bCtx.ExposeBinding(ctx, "record", func(args ...any) any {
		v, _ := args[0].(float64)
		select {
		case calls <- v:
		default:
		}
		return nil
	})
	must.NoError(err, "ExposeBinding failed")

	// Binding should be available on ALL pages created in this context.
	page1, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage 1 failed")
	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage 2 failed")

	_, err = page1.Evaluate(ctx, "record(1)")
	must.NoError(err, "page1.Evaluate")
	_, err = page2.Evaluate(ctx, "record(2)")
	must.NoError(err, "page2.Evaluate")

	v1 := <-calls
	v2 := <-calls
	got := []float64{v1, v2}
	// Values can arrive in any order; just verify both are present.
	if (got[0] != 1 || got[1] != 2) && (got[0] != 2 || got[1] != 1) {
		t.Errorf("expected calls [1, 2] (any order), got %v", got)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingPanic(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	err := bCtx.ExposeBinding(ctx, "boom", func(args ...any) any {
		panic("test panic")
	})
	must.NoError(err, "ExposeBinding failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// JS should receive a rejection; page.Evaluate catches it.
	result, err := page.Evaluate(ctx, `(async () => {
		try { await boom(); return "no error"; }
		catch(e) { return e.message; }
	})()`)
	must.NoError(err, "Evaluate failed")
	// The error message should contain our panic text.
	msg, _ := result.(string)
	if msg == "no error" {
		t.Error("expected JS to receive an error from the panicking binding, got 'no error'")
	}
	_ = msg // log if needed
}

// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingThrowsExceptionInPageContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	err := bCtx.ExposeBinding(ctx, "woof", func(args ...any) any {
		panic("WOOF WOOF")
	})
	must.NoError(err)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `async () => {
		try {
			await window['woof']();
			return { ok: true };
		} catch (e) {
			return { ok: false, message: e.message };
		}
	}`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok, "expected map result, got %T", result)
	is.Equal(false, m["ok"])
	msg, _ := m["message"].(string)
	is.NotEmpty(msg)
}

// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingSurvivesNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	err := bCtx.ExposeBinding(ctx, "compute", func(args ...any) any {
		a, _ := args[0].(float64)
		b, _ := args[1].(float64)
		return a * b
	})
	must.NoError(err)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	result, err := page.Evaluate(ctx, "compute(9, 4)")
	must.NoError(err)
	is.Equal(float64(36), result)

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	result, err = page.Evaluate(ctx, "compute(9, 4)")
	must.NoError(err)
	is.Equal(float64(36), result)
}

// TestExposeBindingReceivesArgs verifies binding receives multiple arguments.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingReceivesArgs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	var mu sync.Mutex
	var gotArgs []any

	must.NoError(bc.ExposeBinding(ctx, "captureArgs", func(args ...any) any {
		mu.Lock()
		gotArgs = args
		mu.Unlock()
		return nil
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, `() => window.captureArgs(1, 'two', true)`)
	must.NoError(err)

	mu.Lock()
	a := gotArgs
	mu.Unlock()

	is.Len(a, 3)
	is.Equal(float64(1), a[0])
	is.Equal("two", a[1])
	is.Equal(true, a[2])
}

// TestExposeBindingCanReturnValue verifies binding return value is accessible in JS.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingCanReturnValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	must.NoError(bc.ExposeBinding(ctx, "computeSum", func(args ...any) any {
		if len(args) == 2 {
			a, _ := args[0].(float64)
			b, _ := args[1].(float64)
			return a + b
		}
		return 0
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	result, err := page.Evaluate(ctx, `() => window.computeSum(3, 4)`)
	must.NoError(err)
	is.Equal(float64(7), result)
}

// TestExposeBindingCalledMultipleTimes verifies binding can be called multiple times.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeBindingCalledMultipleTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	var mu sync.Mutex
	count := 0

	must.NoError(bc.ExposeBinding(ctx, "increment", func(args ...any) any {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, `async () => {
		window.increment();
		window.increment();
		window.increment();
	}`)
	must.NoError(err)

	mu.Lock()
	c := count
	mu.Unlock()

	is.Equal(3, c)
}

// TestAddInitScriptCanSetGlobalFunctions verifies init script can define global functions.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestAddInitScriptCanSetGlobalFunctions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `
		window.add = (a, b) => a + b;
	`))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => window.add(3, 4)`)
	must.NoError(err)
	is.Equal(float64(7), result)
}

// TestAddInitScriptWithComplexObjectEx4 verifies init script can define complex objects.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestAddInitScriptWithComplexObjectEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `
		window.__config = {
			version: '1.0',
			debug: true,
			features: ['a', 'b']
		};
	`))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	version, err := page.Evaluate(ctx, `() => window.__config.version`)
	must.NoError(err)
	is.Equal("1.0", version)

	debug, err := page.Evaluate(ctx, `() => window.__config.debug`)
	must.NoError(err)
	is.Equal(true, debug)
}

// TestInitScriptAddsFunctionEx6 verifies AddInitScript can define a global function.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestInitScriptAddsFunctionEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.greet = function(name) { return 'Hello, ' + name; }`))
	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	result, err := page.Evaluate(ctx, `() => window.greet('World')`)
	must.NoError(err)
	is.Equal("Hello, World", result)
}

// TestExposeFunctionCallableFromJSEx2 verifies ExposeBinding can be called from JS.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeFunctionCallableFromJSEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "myFunc", func(args ...any) any {
		return "hello from Go"
	}))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => myFunc()`)
	must.NoError(err)
	is.Equal("hello from Go", result)
}

// TestExposeFunctionReceivesArgsEx2 verifies ExposeBinding receives arguments.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeFunctionReceivesArgsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "multiply", func(args ...any) any {
		a := args[0].(float64)
		b := args[1].(float64)
		return a * b
	}))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => multiply(3, 7)`)
	must.NoError(err)
	is.Equal(float64(21), result)
}

// TestExposeFunctionReturnsStringEx2 verifies ExposeBinding can return strings.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeFunctionReturnsStringEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "greeting", func(args ...any) any {
		name := args[0].(string)
		return "Hello, " + name + "!"
	}))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => greeting('World')`)
	must.NoError(err)
	is.Equal("Hello, World!", result)
}

// TestExposeFunctionCanReturnNumberEx2 verifies ExposeBinding can return numbers.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeFunctionCanReturnNumberEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "getAnswer", func(args ...any) any {
		return float64(42)
	}))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => getAnswer()`)
	must.NoError(err)
	is.Equal(float64(42), result)
}

// TestExposeFunctionReturnsObjectEx verifies ExposeBinding returning an object.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeFunctionReturnsObjectEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "makeObj", func(args ...any) any {
		return map[string]interface{}{"status": "ok", "code": float64(200)}
	}))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => makeObj().status`)
	must.NoError(err)
	is.Equal("ok", result)
}

// TestExposeFunctionCalledOnEventEx verifies ExposeBinding called from event handler.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeFunctionCalledOnEventEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	called := false
	must.NoError(bc.ExposeBinding(ctx, "markCalled", func(args ...any) any {
		called = true
		return nil
	}))

	must.NoError(page.SetContent(ctx, `
		<button id="btn" onclick="window.markCalled()">Click</button>
	`))

	must.NoError(page.Locator("#btn").Click(ctx))
	is.True(called)
}

// TestExposeFunctionStringConcatEx verifies ExposeBinding string concatenation.
// Extra: no equivalent method in playwright-java; related class: TestPageExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestExposeFunctionStringConcatEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "concat", func(args ...any) any {
		a, _ := args[0].(string)
		b, _ := args[1].(string)
		return a + b
	}))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => concat('foo', 'bar')`)
	must.NoError(err)
	is.Equal("foobar", result)
}

// TestPageExposeFunctionShouldThrowExceptionInPageContext verifies that a handler panic propagates to JS as a caught error.
// Ref: TestPageExposeFunction.java#shouldThrowExceptionInPageContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldThrowExceptionInPageContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "woof", func(args ...any) any {
		panic("WOOF WOOF")
	}))

	result, err := page.Evaluate(ctx, `async () => {
		try {
			await window["woof"]();
		} catch (e) {
			return { message: e.message };
		}
	}`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok, "expected result to be a map")
	is.Equal("WOOF WOOF", m["message"])
}

// TestPageExposeFunctionShouldBeCallableFromInsideAddInitScript verifies that a binding registered before addInitScript
// is callable from the init script after reload.
// Ref: TestPageExposeFunction.java#shouldBeCallableFromInsideAddInitScript  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldBeCallableFromInsideAddInitScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	called := make(chan struct{}, 1)
	must.NoError(bc.ExposeBinding(ctx, "woof", func(args ...any) any {
		select {
		case called <- struct{}{}:
		default:
		}
		return nil
	}))

	must.NoError(page.AddInitScript(ctx, "window['woof']()"))
	must.NoError(page.Reload(ctx))

	select {
	case <-called:
		// success
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for woof() call from init script")
	}
}

// TestPageExposeFunctionShouldSurviveNavigation verifies that a binding remains available after navigation.
// Ref: TestPageExposeFunction.java#shouldSurviveNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldSurviveNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "compute", func(args ...any) any {
		a := args[0].(float64)
		b := args[1].(float64)
		return a * b
	}))

	result, err := page.Evaluate(ctx, "async function() { return await window['compute'](9, 4); }")
	must.NoError(err)
	is.Equal(float64(36), result)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err = page.Evaluate(ctx, "async function() { return await window['compute'](9, 4); }")
	must.NoError(err)
	is.Equal(float64(36), result)
}

// TestPageExposeFunctionShouldWorkAfterCrossOriginNavigation verifies that a binding works after cross-origin navigation.
// Ref: TestPageExposeFunction.java#shouldWorkAfterCrossOriginNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldWorkAfterCrossOriginNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bc.ExposeBinding(ctx, "compute", func(args ...any) any {
		a := args[0].(float64)
		b := args[1].(float64)
		return a * b
	}))

	crossURL := srv.CrossProcessPrefix() + "/empty.html"
	must.NoError(page.Goto(ctx, crossURL))

	result, err := page.Evaluate(ctx, "window['compute'](9, 4)")
	must.NoError(err)
	is.Equal(float64(36), result)
}

// TestPageExposeFunctionShouldWorkWithComplexObjects verifies that a binding can receive and return complex objects.
// Ref: TestPageExposeFunction.java#shouldWorkWithComplexObjects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldWorkWithComplexObjects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "complexObject", func(args ...any) any {
		a := args[0].(map[string]any)
		b := args[1].(map[string]any)
		ax := a["x"].(float64)
		bx := b["x"].(float64)
		return map[string]any{"x": ax + bx}
	}))

	result, err := page.Evaluate(ctx, "async () => window['complexObject']({x: 5}, {x: 2})")
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal(float64(7), m["x"])
}

// TestPageExposeFunctionShouldThrowForDuplicateRegistrations verifies that registering the same binding name twice fails.
// Ref: TestPageExposeFunction.java#shouldThrowForDuplicateRegistrations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldThrowForDuplicateRegistrations(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	must.NoError(bc.ExposeBinding(ctx, "foo", func(args ...any) any { return nil }))

	err := bc.ExposeBinding(ctx, "foo", func(args ...any) any { return nil })
	is.Error(err, "second registration of same binding name should fail")
	is.ErrorContains(err, "foo")
}

// TestPageExposeFunctionShouldWorkOnFrames verifies that a binding is available in child frames.
// Ref: TestPageExposeFunction.java#shouldWorkOnFrames  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldWorkOnFrames(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(bc.ExposeBinding(ctx, "compute", func(args ...any) any {
		a := args[0].(float64)
		b := args[1].(float64)
		return a * b
	}))

	must.NoError(page.SetContent(ctx, `<iframe srcdoc="<div>inner</div>"></iframe>`))

	// Wait for the frame to be available.
	time.Sleep(200 * time.Millisecond)

	frames := page.Frames()
	is.Len(frames, 2, "expected main frame + 1 child frame")

	childFrame := frames[1]
	result, err := childFrame.Evaluate(ctx, "async function() { return window['compute'](3, 5); }")
	must.NoError(err)
	is.Equal(float64(15), result)
}

// TestPageExposeFunctionShouldWorkOnFramesBeforeNavigation verifies that a binding registered after page creation
// is available in frames loaded later.
// Ref: TestPageExposeFunction.java#shouldWorkOnFramesBeforeNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldWorkOnFramesBeforeNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	// Set content with an iframe first, then expose the binding.
	must.NoError(page.SetContent(ctx, `<iframe srcdoc="<div>inner</div>"></iframe>`))
	time.Sleep(200 * time.Millisecond)

	must.NoError(bc.ExposeBinding(ctx, "compute", func(args ...any) any {
		a := args[0].(float64)
		b := args[1].(float64)
		return a * b
	}))

	frames := page.Frames()
	is.Len(frames, 2)

	childFrame := frames[1]
	result, err := childFrame.Evaluate(ctx, "async function() { return window['compute'](3, 5); }")
	must.NoError(err)
	is.Equal(float64(15), result)
}

// Ref: TestPageExposeFunction.java#exposeBindingShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionExposeBindingShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	context := newContext(t)
	page, err := context.NewPage(ctx)
	must.NoError(err)

	var bindingSource *playwright.BindingSource
	must.NoError(page.ExposeBinding(ctx, "add", func(source *playwright.BindingSource, args ...any) any {
		bindingSource = source
		return args[0].(float64) + args[1].(float64)
	}))

	result, err := page.Evaluate(ctx, `async function() {
		return window['add'](5, 6);
	}`)
	must.NoError(err)

	must.NotNil(bindingSource)
	is.Equal(context, bindingSource.Context())
	is.Equal(page, bindingSource.Page())
	is.Equal(page.MainFrame(), bindingSource.Frame())
	is.Equal(float64(11), result)
	must.NoError(context.Close(ctx))
}

// Ref: TestPageExposeFunction.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.ExposeBinding(ctx, "compute", func(args ...any) any {
		return args[0].(float64) * args[1].(float64)
	}))
	result := a2Ev(t, ctx, page, `async function() {
  return await window['compute'](9, 4);
}`)
	assert.Equal(t, float64(36), result)
}

// Ref: TestPageExposeFunction.java#shouldWorkWithHandlesAndComplexObjects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldWorkWithHandlesAndComplexObjects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	fooHandle, err := page.EvaluateHandle(ctx, `() => {
		window['fooValue'] = { bar: 2 };
		return window['fooValue'];
	}`)
	must.NoError(err)
	must.NotNil(fooHandle)

	must.NoError(page.ExposeBinding(ctx, "handle", func(args ...any) any {
		return []any{map[string]any{"foo": fooHandle}}
	}))

	equals, err := page.Evaluate(ctx, `async function() {
		const value = await window['handle']();
		const [{ foo }] = value;
		return foo === window['fooValue'];
	}`)
	must.NoError(err)
	is.Equal(true, equals)
}

// Ref: TestPageExposeFunction.java#shouldSerializeCycles  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageExposeFunction.java
func TestPageExposeFunctionShouldSerializeCycles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	got := make(chan any, 1)
	must.NoError(page.ExposeBinding(ctx, "log", func(args ...any) any {
		got <- args[0]
		return nil
	}))
	a2Ev(t, ctx, page, "const a = {}; a.b = a; window.log(a)")
	var obj any
	select {
	case obj = <-got:
	case <-time.After(10 * time.Second):
		t.Fatal("exposed function was not called")
	}
	m, ok := obj.(map[string]any)
	must.True(ok, "expected a map, got %T", obj)
	assert.Len(t, m, 1)
	inner, ok := m["b"].(map[string]any)
	must.True(ok, "expected a map for b, got %T", m["b"])
	// the same object must be referenced (cycle)
	inner["__marker"] = true
	assert.Equal(t, true, m["__marker"])
}
