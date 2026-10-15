//go:build e2e

// E2E tests for: TestBrowserContextExposeFunction.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

// Ref: TestBrowserContextExposeFunction.java#exposeBindingShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
func TestBrowserContextExposeFunctionExposeBindingShouldWork(t *testing.T) {
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

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
func TestBrowserContextAddInitScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	must.NoError(bCtx.AddInitScript(ctx, `window.__ctxInjected = 'from-context'`), "AddInitScript failed")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	must.NoError(page.Goto(ctx, srv.EmptyPage()), "Goto failed")

	val, err := page.Evaluate(ctx, `() => window.__ctxInjected`)
	must.NoError(err, "Evaluate failed")
	is.Equal("from-context", val, "context init script should inject window.__ctxInjected")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
func TestBrowserContextAddInitScriptAffectsAllPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)
	must.NoError(bCtx.AddInitScript(ctx, `window.__shared = 'shared-value'`), "AddInitScript failed")

	page1, err := bCtx.NewPage(ctx)
	must.NoError(err, "page1 NewPage failed")
	must.NoError(page1.Goto(ctx, srv.EmptyPage()), "page1 Goto failed")

	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "page2 NewPage failed")
	must.NoError(page2.Goto(ctx, srv.EmptyPage()), "page2 Goto failed")

	val1, err := page1.Evaluate(ctx, `() => window.__shared`)
	must.NoError(err, "page1 Evaluate failed")
	is.Equal("shared-value", val1, "page1 should have init script value")

	val2, err := page2.Evaluate(ctx, `() => window.__shared`)
	must.NoError(err, "page2 Evaluate failed")
	is.Equal("shared-value", val2, "page2 should have init script value")
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
func TestBrowserContextAddInitScriptEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	must.NoError(bc.AddInitScript(ctx, `window.__fromContext = 'hello';`))

	p1, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(p1.SetContent(ctx, `<div></div>`))

	result, err := p1.Evaluate(ctx, `() => window.__fromContext`)
	must.NoError(err)
	is.Equal("hello", result)
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
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

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextExposeFunction.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
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

// TestBrowserContextExposeFunctionShouldWork verifies context- and page-level exposed functions can be called from the page.
// Ref: TestBrowserContextExposeFunction.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
func TestBrowserContextExposeFunctionShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	num := func(v any) float64 { f, _ := v.(float64); return f }
	must.NoError(bCtx.ExposeBinding(ctx, "add", func(args ...any) any { return num(args[0]) + num(args[1]) }))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.ExposeBinding(ctx, "mul", func(args ...any) any { return num(args[0]) * num(args[1]) }))
	must.NoError(bCtx.ExposeBinding(ctx, "sub", func(args ...any) any { return num(args[0]) - num(args[1]) }))
	// The Go BindingHandler has no source argument, so the "addHandle" binding (source.frame().evaluateHandle) cannot be ported.
	result, err := page.Evaluate(ctx, "async () => ({ mul: await mul(9, 4), add: await add(9, 4), sub: await sub(9, 4) })")
	must.NoError(err)
	must.EqualValues(map[string]any{"mul": 36.0, "add": 13.0, "sub": 5.0}, result)
}

// TestBrowserContextExposeFunctionShouldThrowForDuplicateRegistrations verifies duplicate registrations are rejected.
// Ref: TestBrowserContextExposeFunction.java#shouldThrowForDuplicateRegistrations  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
func TestBrowserContextExposeFunctionShouldThrowForDuplicateRegistrations(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	nop := func(args ...any) any { return nil }
	must.NoError(bCtx.ExposeBinding(ctx, "foo", nop))
	must.NoError(bCtx.ExposeBinding(ctx, "bar", nop))
	err := bCtx.ExposeBinding(ctx, "foo", nop)
	must.Error(err)
	must.Contains(err.Error(), "Function \"foo\" has been already registered")

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	err = page.ExposeBinding(ctx, "foo", nop)
	must.Error(err)
	must.Contains(err.Error(), "Function \"foo\" has been already registered in the browser context")

	must.NoError(page.ExposeBinding(ctx, "baz", nop))
	err = bCtx.ExposeBinding(ctx, "baz", nop)
	must.Error(err)
	must.Contains(err.Error(), "Function \"baz\" has been already registered in one of the pages")
}

// TestBrowserContextExposeFunctionShouldBeCallableFromInsideAddInitScript verifies exposed functions are available to init scripts.
// Ref: TestBrowserContextExposeFunction.java#shouldBeCallableFromInsideAddInitScript  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextExposeFunction.java
func TestBrowserContextExposeFunctionShouldBeCallableFromInsideAddInitScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	var mu sync.Mutex
	var actual []any
	must.NoError(bCtx.ExposeBinding(ctx, "woof", func(args ...any) any {
		mu.Lock()
		actual = append(actual, args[0])
		mu.Unlock()
		return nil
	}))
	must.NoError(bCtx.AddInitScript(ctx, "window['woof']('context')"))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	_, err = page.Evaluate(ctx, "undefined")
	must.NoError(err)
	snapshot := func() []any { mu.Lock(); defer mu.Unlock(); return append([]any(nil), actual...) }
	must.Eventually(func() bool { return len(snapshot()) == 1 }, 5*time.Second, 20*time.Millisecond)
	must.Equal([]any{"context"}, snapshot())
	mu.Lock()
	actual = nil
	mu.Unlock()
	must.NoError(page.AddInitScript(ctx, "window['woof']('page')"))
	must.NoError(page.Reload(ctx))
	must.Eventually(func() bool { return len(snapshot()) == 2 }, 5*time.Second, 20*time.Millisecond)
	must.Equal([]any{"context", "page"}, snapshot())
}
