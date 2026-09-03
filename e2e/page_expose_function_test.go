//go:build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExposeFunctionCallableFromJSEx2 verifies ExposeBinding can be called from JS.
// Ref: TestPageExposeFunction.java#shouldBeCallableFromJS
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
// Ref: TestPageExposeFunction.java#shouldReceiveArgs
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
// Ref: TestPageExposeFunction.java#shouldReturnString
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
// Ref: TestPageExposeFunction.java#shouldReturnNumber
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
// Ref: TestPageExposeFunction.java#shouldReturnObject
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
// Ref: TestPageExposeFunction.java#shouldBeCallableFromEventHandler
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
// Ref: TestPageExposeFunction.java#shouldConcatenateStrings
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

// TestExposeFunctionThrowsInPageContext verifies that a handler panic propagates to JS as a caught error.
// Ref: TestPageExposeFunction.java#shouldThrowExceptionInPageContext
func TestExposeFunctionThrowsInPageContext(t *testing.T) {
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

// TestExposeFunctionCallableFromInitScript verifies that a binding registered before addInitScript
// is callable from the init script after reload.
// Ref: TestPageExposeFunction.java#shouldBeCallableFromInsideAddInitScript
func TestExposeFunctionCallableFromInitScript(t *testing.T) {
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

// TestExposeFunctionSurvivesNavigation verifies that a binding remains available after navigation.
// Ref: TestPageExposeFunction.java#shouldSurviveNavigation
func TestExposeFunctionSurvivesNavigation(t *testing.T) {
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

// TestExposeFunctionWorksCrossOrigin verifies that a binding works after cross-origin navigation.
// Ref: TestPageExposeFunction.java#shouldWorkAfterCrossOriginNavigation
func TestExposeFunctionWorksCrossOrigin(t *testing.T) {
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

// TestExposeFunctionWithComplexObjects verifies that a binding can receive and return complex objects.
// Ref: TestPageExposeFunction.java#shouldWorkWithComplexObjects
func TestExposeFunctionWithComplexObjects(t *testing.T) {
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

// TestExposeFunctionDuplicateRegistrationFails verifies that registering the same binding name twice fails.
// Ref: TestPageExposeFunction.java#shouldThrowForDuplicateRegistrations
func TestExposeFunctionDuplicateRegistrationFails(t *testing.T) {
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

// TestExposeFunctionWorksInFrame verifies that a binding is available in child frames.
// Ref: TestPageExposeFunction.java#shouldWorkOnFrames
func TestExposeFunctionWorksInFrame(t *testing.T) {
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

// TestExposeFunctionWorksOnFrameBeforeNavigation verifies that a binding registered after page creation
// is available in frames loaded later.
// Ref: TestPageExposeFunction.java#shouldWorkOnFramesBeforeNavigation
func TestExposeFunctionWorksOnFrameBeforeNavigation(t *testing.T) {
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
