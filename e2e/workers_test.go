//go:build e2e

// E2E tests for: TestWorkers.java.
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

// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestPageOnWorker(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Serve a minimal service worker script.
	srv.SetRoute("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.addEventListener('install', () => self.skipWaiting());`))
	})
	srv.ServeWithBody("/sw-page", "text/html", `
		<script>
			if ('serviceWorker' in navigator) {
				navigator.serviceWorker.register('/sw.js');
			}
		</script>
	`)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	var mu sync.Mutex
	var workers []*playwright.Worker
	workerSeen := make(chan struct{}, 1)

	page.OnWorker(func(w *playwright.Worker) {
		mu.Lock()
		workers = append(workers, w)
		mu.Unlock()
		select {
		case workerSeen <- struct{}{}:
		default:
		}
	})

	err = page.Goto(ctx, srv.Prefix()+"/sw-page")
	must.NoError(err, "Goto failed")

	// Wait for worker to register (with timeout). Service workers may not work
	// in all configurations; skip rather than fail if none appears.
	select {
	case <-workerSeen:
		// at least one worker appeared
	case <-time.After(5 * time.Second):
		t.Skip("no service worker appeared within 5s")
	}

	mu.Lock()
	workerCount := len(workers)
	mu.Unlock()

	if workerCount == 0 {
		t.Error("expected at least one worker, got zero")
	}

	// Verify Page.Workers() is consistent with what OnWorker reported.
	pageWorkers := page.Workers()
	if len(pageWorkers) == 0 {
		t.Error("Page.Workers() returned empty slice after OnWorker fired")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestPageWorkersInitiallyEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/no-workers", "text/html", `<p>no workers here</p>`)
	err := page.Goto(ctx, srv.Prefix()+"/no-workers")
	must.NoError(err, "Goto failed")

	workers := page.Workers()
	if len(workers) != 0 {
		t.Errorf("expected 0 workers on a plain page, got %d", len(workers))
	}
}

// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	// Serve a dedicated worker script and the page that loads it.
	srv.SetRoute("/eval-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.myValue = 42;`))
	})
	srv.ServeWithBody("/eval-worker-page", "text/html",
		`<script>new Worker('/eval-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	err := page.Goto(ctx, srv.Prefix()+"/eval-worker-page")
	must.NoError(err, "Goto failed")

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	val, err := w.Evaluate(ctx, "self.myValue")
	must.NoError(err, "Worker.Evaluate failed")
	if val != float64(42) {
		t.Errorf("Worker.Evaluate = %v (%T), want float64(42)", val, val)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerEvaluateExpression(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/expr-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.add = (a, b) => a + b;`))
	})
	srv.ServeWithBody("/expr-worker-page", "text/html",
		`<script>new Worker('/expr-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	err := page.Goto(ctx, srv.Prefix()+"/expr-worker-page")
	must.NoError(err, "Goto failed")

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	// Evaluate a JS expression with an argument.
	val, err := w.Evaluate(ctx, "1 + 2")
	must.NoError(err, "Worker.Evaluate failed")
	if val != float64(3) {
		t.Errorf("Worker.Evaluate(1+2) = %v, want 3", val)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/url-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.x = 1;`))
	})
	srv.ServeWithBody("/url-worker-page", "text/html",
		`<script>new Worker('/url-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	err := page.Goto(ctx, srv.Prefix()+"/url-worker-page")
	must.NoError(err, "Goto failed")

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	url := w.URL()
	if !strings.HasSuffix(url, "/url-worker.js") {
		t.Errorf("Worker.URL() = %q, expected suffix '/url-worker.js'", url)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerEvaluateHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/handle-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.obj = {x: 7};`))
	})
	srv.ServeWithBody("/handle-worker-page", "text/html",
		`<script>new Worker('/handle-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	err := page.Goto(ctx, srv.Prefix()+"/handle-worker-page")
	must.NoError(err, "Goto failed")

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	handle, err := w.EvaluateHandle(ctx, "() => ({answer: 42})")
	must.NoError(err, "Worker.EvaluateHandle failed")
	must.NotNil(handle, "Worker.EvaluateHandle returned nil handle")

	// Retrieve the JSON value from the handle.
	val, err := handle.JSONValue(ctx)
	must.NoError(err, "JSHandle.JSONValue failed")
	m, ok := val.(map[string]any)
	if !ok {
		t.Fatalf("JSONValue result is not a map: %T %v", val, val)
	}
	if m["answer"] != float64(42) {
		t.Errorf("expected answer=42, got %v", m["answer"])
	}
}

// TestWorkerEvaluateWithArgument verifies Worker.Evaluate can receive arguments.
// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerEvaluateWithArgument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/arg-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.multiply = (a, b) => a * b;`))
	})
	srv.ServeWithBody("/arg-worker-page", "text/html",
		`<script>new Worker('/arg-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/arg-worker-page"))

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	val, err := w.Evaluate(ctx, "([a, b]) => a * b", []any{6, 7})
	must.NoError(err)
	is.Equal(float64(42), val)
}

// TestWorkerEvaluateReturnsObject verifies Worker.Evaluate can return complex objects.
// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerEvaluateReturnsObject(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/obj-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.getInfo = () => ({ name: "worker", version: 2 });`))
	})
	srv.ServeWithBody("/obj-worker-page", "text/html",
		`<script>new Worker('/obj-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/obj-worker-page"))

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	val, err := w.Evaluate(ctx, "() => ({ name: 'worker', version: 2 })")
	must.NoError(err)

	m, ok := val.(map[string]any)
	is.True(ok, "expected map result")
	is.Equal("worker", m["name"])
	is.Equal(float64(2), m["version"])
}

// TestWorkerURLMatchesScript verifies Worker.URL returns the script URL.
// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerURLMatchesScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/my-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.x = 1;`))
	})
	srv.ServeWithBody("/my-worker-page", "text/html",
		`<script>new Worker('/my-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/my-worker-page"))

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	is.Contains(w.URL(), "/my-worker.js")
}

// TestWorkerRemovedAfterNavigation verifies that workers are cleared after page navigation.
// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerRemovedAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/nav-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.x = 1;`))
	})
	srv.ServeWithBody("/nav-worker-page", "text/html",
		`<script>new Worker('/nav-worker.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/nav-worker-page"))

	select {
	case <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	is.NotEmpty(page.Workers(), "expected workers before navigation")

	// Navigate away — workers should be gone after new page loads.
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	is.Empty(page.Workers(), "expected no workers after navigating to plain page")
}

// TestPageWorkersEmpty verifies Workers returns empty slice when no workers.
// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestPageWorkersEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>no workers</div>`))

	workers := page.Workers()
	is.Empty(workers)
}

// TestPageWorkersIsSlice verifies Workers returns a slice (not nil).
// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestPageWorkersIsSlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	workers := page.Workers()
	// Should return a non-nil empty slice, not nil
	must.NotNil(workers)
}

// TestWorkerEvaluateHandleReturnsHandle verifies EvaluateHandle returns a non-nil handle.
// Extra: no equivalent method in playwright-java; related class: TestWorkers.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkerEvaluateHandleReturnsHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	srv.SetRoute("/handle-worker2.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(`self.data = [1, 2, 3];`))
	})
	srv.ServeWithBody("/handle-worker2-page", "text/html",
		`<script>new Worker('/handle-worker2.js');</script>`)

	page := newPage(t)
	workerCh := waitForWorker(page)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/handle-worker2-page"))

	var w *playwright.Worker
	select {
	case w = <-workerCh:
	case <-time.After(10 * time.Second):
		t.Fatal("dedicated worker did not appear within 10s")
	}

	handle, err := w.EvaluateHandle(ctx, "() => [1, 2, 3]")
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	arr, ok := val.([]any)
	is.True(ok, "expected array result")
	is.Len(arr, 3)
}

const a2NewWorkerScript = "() => new Worker(URL.createObjectURL(new Blob(['console.log(1)'], {type: 'application/javascript'})))"

// a2WaitForWorker runs the trigger and waits for the next worker created on the page.
func a2WaitForWorker(t *testing.T, ctx context.Context, page *playwright.Page, trigger func()) *playwright.Worker {
	t.Helper()
	ch := make(chan *playwright.Worker, 4)
	off := page.OnWorker(func(w *playwright.Worker) { ch <- w })
	defer off()
	trigger()
	select {
	case w := <-ch:
		return w
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for worker")
	}
	return nil
}

// Ref: TestWorkers.java#shouldEmitCreatedAndDestroyedEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldEmitCreatedAndDestroyedEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, _, page := a2Setup(t)

	worker := a2WaitForWorker(t, ctx, page, func() {
		a2Ev(t, ctx, page, "() => { window.worker = new Worker(URL.createObjectURL(new Blob(['1'], {type: 'application/javascript'}))); }")
	})
	workerThisObj, err := worker.EvaluateHandle(ctx, "() => this")
	must.NoError(err)
	closedWorker, err := worker.WaitForClose(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.worker.terminate()")
		return err
	})
	must.NoError(err)
	is.Same(worker, closedWorker)

	_, err = workerThisObj.GetProperty(ctx, "self")
	must.Error(err)
	must.True(strings.Contains(err.Error(), "Target page, context or browser has been closed") ||
		strings.Contains(err.Error(), "Worker was closed"), err.Error())
}

// Ref: TestWorkers.java#shouldReportConsoleLogs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldReportConsoleLogs(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	ch := make(chan string, 4)
	page.OnConsole(func(m *playwright.ConsoleMessage) { ch <- m.Text() })
	a2Ev(t, ctx, page, a2NewWorkerScript)
	select {
	case text := <-ch:
		assert.Equal(t, "1", text)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for console message")
	}
}

// Ref: TestWorkers.java#shouldHaveJSHandlesForConsoleLogs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldHaveJSHandlesForConsoleLogs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, _, page := a2Setup(t)

	messages := make(chan *playwright.ConsoleMessage, 1)
	page.OnConsole(func(message *playwright.ConsoleMessage) { messages <- message })
	a2Ev(t, ctx, page, "() => new Worker(URL.createObjectURL(new Blob(['console.log(1,2,3,this)'], {type: 'application/javascript'})))")

	var message *playwright.ConsoleMessage
	select {
	case message = <-messages:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for worker console message")
	}
	if globalBTName == "firefox" {
		is.Equal("1 2 3 JSHandle@object", message.Text())
	} else {
		is.Equal("1 2 3 DedicatedWorkerGlobalScope", message.Text())
	}
	args := message.Args()
	must.Len(args, 4)
	origin, err := args[3].GetProperty(ctx, "origin")
	must.NoError(err)
	value, err := origin.JSONValue(ctx)
	must.NoError(err)
	is.Equal("null", value)
}

// Ref: TestWorkers.java#shouldEvaluate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	worker := a2WaitForWorker(t, ctx, page, func() { a2Ev(t, ctx, page, a2NewWorkerScript) })
	res, err := worker.Evaluate(ctx, "1+1")
	must.NoError(err)
	assert.Equal(t, float64(2), res)
}

// Ref: TestWorkers.java#shouldReportErrors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldReportErrors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	a2Ev(t, ctx, page, "() => new Worker(URL.createObjectURL(new Blob([`\n  setTimeout(() => {\n    // Do a console.log just to check that we do not confuse it with an error.\n    console.log('hey');\n    throw new Error('this is my error');\n  })\n`], {type: 'application/javascript'})))")
	deadline := time.Now().Add(30 * time.Second)
	var errs []string
	for len(errs) == 0 {
		must.True(time.Now().Before(deadline), "Timed out")
		must.NoError(page.WaitForTimeout(ctx, 100))
		var err error
		errs, err = page.PageErrors(ctx)
		must.NoError(err)
	}
	assert.Contains(t, errs[0], "this is my error")
}

// Ref: TestWorkers.java#shouldClearUponCrossProcessNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldClearUponCrossProcessNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	worker := a2WaitForWorker(t, ctx, page, func() {
		a2Ev(t, ctx, page, "() => new Worker(URL.createObjectURL(new Blob(['console.log(1)'], {type: 'application/javascript'})))")
	})
	destroyed := make(chan struct{}, 1)
	worker.OnClose(func(closedWorker *playwright.Worker) {
		destroyed <- struct{}{}
	})
	must.NoError(page.Goto(ctx, srv.CrossProcessPrefix()+"/empty.html"))
	select {
	case <-destroyed:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for worker close event")
	}
	must.Empty(page.Workers())
}

// Ref: TestWorkers.java#shouldReportNetworkActivityOnWorkerCreation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldReportNetworkActivityOnWorkerCreation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/one-style.css", "text/css", "body { color: red; }")
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	targetURL := srv.Prefix() + "/one-style.css"

	var req *playwright.NetworkRequest
	resp, err := page.WaitForResponse(ctx, targetURL, func() {
		var waitErr error
		req, waitErr = page.WaitForRequest(ctx, targetURL, func() {
			_, evalErr := page.Evaluate(ctx, `url => new Worker(URL.createObjectURL(new Blob([`+"`"+`
				fetch('${url}').then(response => response.text()).then(console.log);
			`+"`"+`], {type: 'application/javascript'})))`, targetURL)
			must.NoError(evalErr)
		})
		must.NoError(waitErr)
		must.NotNil(req)
		is.Equal(targetURL, req.URL())
	})
	must.NoError(err)
	must.NotNil(resp)
	is.Equal(req, resp.Request())
	is.True(resp.OK())
}

// Ref: TestWorkers.java#shouldFormatNumberUsingContextLocale  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldFormatNumberUsingContextLocale(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := a2NewPage(t, &playwright.BrowserContextOptions{Locale: a2Ptr("ru-RU")})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	worker := a2WaitForWorker(t, ctx, page, func() { a2Ev(t, ctx, page, a2NewWorkerScript) })
	res, err := worker.Evaluate(ctx, "() => (10000.20).toLocaleString()")
	must.NoError(err)
	assert.Equal(t, "10\u00A0000,2", res)
}

// Ref: TestWorkers.java#shouldReportConsoleEventOnTheWorker  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWorkers.java
func TestWorkersShouldReportConsoleEventOnTheWorker(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx, _, page := a2Setup(t)

	worker, err := page.WaitForWorker(ctx, func() {
		_, _ = page.Evaluate(ctx, "() => { window.worker = new Worker(URL.createObjectURL(new Blob(['42'], {type: 'application/javascript'}))); }")
	})
	must.NoError(err)
	must.NotNil(worker)

	var message2 *playwright.ConsoleMessage
	var message3 *playwright.ConsoleMessage

	offPage := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		message2 = msg
	})
	defer offPage()

	offCtx := page.Context().OnConsole(func(msg *playwright.ConsoleMessage) {
		message3 = msg
	})
	defer offCtx()

	message1, err := worker.WaitForConsoleMessage(ctx, func() {
		_, _ = worker.Evaluate(ctx, "() => console.log('hello from worker')")
	})
	must.NoError(err)
	must.NotNil(message1)

	is.Equal("hello from worker", message1.Text())
	must.Same(message1, message2)
	must.Same(message1, message3)
}
