//go:build e2e

package e2e

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

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

func TestExposeBindingString(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	err := bCtx.ExposeBinding(ctx, "greet", func(args ...any) any {
		name, _ := args[0].(string)
		return "hello " + name
	})
	must.NoError(err, "ExposeBinding failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	result, err := page.Evaluate(ctx, `greet("world")`)
	must.NoError(err, "Evaluate failed")
	if result != "hello world" {
		t.Errorf("greet('world') = %v, want 'hello world'", result)
	}
}

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

// Ensure playwright import is used.
var _ *playwright.BrowserContext

func TestExposeBindingCallableFromInitScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	called := make(chan bool, 1)
	err := bCtx.ExposeBinding(ctx, "woof", func(args ...any) any {
		select {
		case called <- true:
		default:
		}
		return nil
	})
	must.NoError(err)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	err = page.AddInitScript(ctx, "window['woof']()")
	must.NoError(err)

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	select {
	case <-called:

	case <-ctx.Done():
		t.Fatal("timed out waiting for init script to call binding")
	}
}

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

func TestExposeBindingDuplicateThrows(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	err := bCtx.ExposeBinding(ctx, "foo", func(args ...any) any { return nil })
	must.NoError(err)

	err = bCtx.ExposeBinding(ctx, "foo", func(args ...any) any { return nil })
	is.Error(err, "expected error on duplicate binding name")
}

// TestExposeBindingReceivesArgs verifies binding receives multiple arguments.
// Ref: TestPageExposeBinding.java#shouldReceiveArgs
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
// Ref: TestPageExposeBinding.java#shouldReturnValue
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
// Ref: TestPageExposeBinding.java#shouldCallMultipleTimes
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
