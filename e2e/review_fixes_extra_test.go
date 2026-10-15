//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"sync"
	"testing"
	"time"
)

// Observaciones 1.3: handles nested inside arrays/objects must be serialized as handle references.
func TestEvaluateNestedHandleArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="a">hello</div><div id="b">world</div>`))
	a, err := page.QuerySelector(ctx, "#a")
	must.NoError(err)
	b, err := page.QuerySelector(ctx, "#b")
	must.NoError(err)

	res, err := page.Evaluate(ctx, "args => args.list[0].textContent + ' ' + args.second.textContent",
		map[string]any{"list": []any{a}, "second": b})
	must.NoError(err)
	is.Equal("hello world", res)
}

// Observaciones 1.4: console events must reach the handler in emission order.
func TestOnConsoleEventsInOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	const n = 200
	var mu sync.Mutex
	var got []string
	cancel := page.OnConsole(func(m *playwright.ConsoleMessage) {
		mu.Lock()
		got = append(got, m.Text())
		mu.Unlock()
	})
	defer cancel()

	_, err := page.Evaluate(ctx, fmt.Sprintf("() => { for (let i = 0; i < %d; i++) console.log('m' + i); }", n))
	must.NoError(err)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		l := len(got)
		mu.Unlock()
		if l >= n {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	mu.Lock()
	defer mu.Unlock()
	must.Len(got, n)
	for i, text := range got {
		if text != fmt.Sprintf("m%d", i) {
			t.Fatalf("message at position %d was %q", i, text)
		}
	}
}

// Observaciones 1.5: a panic inside a user handler must not crash the test process.
func TestOnConsolePanicInHandlerIsRecovered(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var wg sync.WaitGroup
	wg.Add(1)
	var once sync.Once
	cancel := page.OnConsole(func(m *playwright.ConsoleMessage) {
		once.Do(wg.Done)
		panic("handler panic")
	})
	defer cancel()

	_, err := page.Evaluate(ctx, "() => console.log('x')")
	must.NoError(err)
	wg.Wait()

	// The connection must still work afterwards.
	res, err := page.Evaluate(ctx, "() => 1 + 1")
	must.NoError(err)
	assert.Equal(t, float64(2), res)
}

// Observaciones 2.1: SetDefaultTimeout must bound actions that do not set their own timeout.
func TestSetDefaultTimeoutAppliesToLocatorActions(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>nothing to click</div>`))
	must.NoError(page.SetDefaultTimeout(ctx, 500))

	start := time.Now()
	err := page.Locator("#missing").Click(ctx)
	must.Error(err)
	assert.True(t, strings.Contains(err.Error(), "500ms"), "unexpected error: %v", err)
	assert.Less(t, time.Since(start), 10*time.Second)
}

// Observaciones 2.1: BrowserContext.SetDefaultTimeout is inherited by its pages.
func TestBrowserContextDefaultTimeoutInheritedByPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bctx := newContext(t)
	bctx.SetDefaultTimeout(500)
	page, err := bctx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<div>nothing to click</div>`))
	err = page.Locator("#missing").Click(ctx)
	must.Error(err)
	assert.True(t, strings.Contains(err.Error(), "500ms"), "unexpected error: %v", err)
}
