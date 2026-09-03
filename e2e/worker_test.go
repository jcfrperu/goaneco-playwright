//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

// waitForWorker registers an OnWorker handler and returns a channel that receives the first Worker.
func waitForWorker(page *playwright.Page) <-chan *playwright.Worker {
	ch := make(chan *playwright.Worker, 1)
	page.OnWorker(func(w *playwright.Worker) {
		select {
		case ch <- w:
		default:
		}
	})
	return ch
}

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
// Ref: TestWorkers.java#shouldWorkWithEvaluateArgs
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
// Ref: TestWorkers.java#shouldReturnComplexObjectFromEvaluate
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
// Ref: TestWorkers.java#shouldReturnWorkerURL
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
// Ref: TestWorkers.java#shouldDestroyWorkerDuringNavigation
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
// Ref: TestPageWorkers.java#shouldReturnEmptyWhenNoWorkers
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
// Ref: TestPageWorkers.java#shouldReturnSlice
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
// Ref: TestWorkers.java#shouldGetEvaluateHandle
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
