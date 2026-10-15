//go:build e2e

// Shared helpers for ported playwright-java tests (batch a2).
package e2e

import (
	"context"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
	"sync"
)

const a2ResourcesDir = "../playwright-reference-projects/playwright-java/playwright/src/test/resources"

// a2ServeAll registers every file from the playwright-java test resources directory
// as a route of the test server (equivalent to the static file serving of Server.java).
func a2ServeAll(t testing.TB, srv *testserver.Server) {
	t.Helper()
	root, err := filepath.Abs(a2ResourcesDir)
	if err != nil {
		t.Fatalf("a2ServeAll: %v", err)
	}
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		route := "/" + filepath.ToSlash(rel)
		file := path
		srv.SetRoute(route, func(w http.ResponseWriter, r *http.Request) {
			data, err := os.ReadFile(file)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			ct := mime.TypeByExtension(filepath.Ext(file))
			if strings.HasSuffix(file, ".html") {
				ct = "text/html; charset=utf-8"
			}
			if ct == "" {
				ct = "application/octet-stream"
			}
			w.Header().Set("Content-Type", ct)
			w.Write(data)
		})
		return nil
	})
	if err != nil {
		t.Fatalf("a2ServeAll: %v", err)
	}
}

// a2Setup returns a standard test context, a test server serving all resources and a new page.
func a2Setup(t *testing.T) (context.Context, *testserver.Server, *playwright.Page) {
	t.Helper()
	ctx := testCtx(t)
	srv := testserver.New(t)
	a2ServeAll(t, srv)
	page := newPage(t)
	return ctx, srv, page
}

func a2Ptr[T any](v T) *T { return &v }

// attachFrame adds an iframe with given id and url to the page and returns its Frame.
func attachFrame(t *testing.T, ctx context.Context, page *playwright.Page, name, url string) *playwright.Frame {
	t.Helper()
	handle, err := page.EvaluateHandle(ctx, `async ({ frameId, url }) => {
		const frame = document.createElement('iframe');
		frame.src = url;
		frame.id = frameId;
		document.body.appendChild(frame);
		await new Promise(x => frame.onload = x);
		return frame;
	}`, map[string]any{"frameId": name, "url": url})
	require.NoError(t, err)
	el := handle.AsElement()
	require.NotNil(t, el)
	frame, err := el.ContentFrame(ctx)
	require.NoError(t, err)
	return frame
}

// a2NewPage creates a page in a new context created with the given options.
func a2NewPage(t *testing.T, opts *playwright.BrowserContextOptions) *playwright.Page {
	t.Helper()
	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx, opts)
	if err != nil {
		t.Fatalf("a2NewPage: NewContext: %v", err)
	}
	t.Cleanup(func() {
		cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(cctx)
	})
	page, err := bCtx.NewPage(ctx)
	if err != nil {
		t.Fatalf("a2NewPage: NewPage: %v", err)
	}
	return page
}

// a2AttachFrame attaches an iframe with the given id and url, and waits for it to load.
func a2AttachFrame(ctx context.Context, page *playwright.Page, id, url string) (*playwright.Frame, error) {
	_, err := page.Evaluate(ctx, `async ({frameId, url}) => {
		const frame = document.createElement('iframe');
		frame.src = url;
		frame.id = frameId;
		document.body.appendChild(frame);
		await new Promise(x => frame.onload = x);
		return frame;
	}`, map[string]any{"frameId": id, "url": url})
	if err != nil {
		return nil, err
	}
	frames := page.Frames()
	return frames[len(frames)-1], nil
}

// a2Ev evaluates an expression on the page and fails the test on error.
func a2Ev(t *testing.T, ctx context.Context, page *playwright.Page, expr string, arg ...any) any {
	t.Helper()
	v, err := page.Evaluate(ctx, expr, arg...)
	if err != nil {
		t.Fatalf("evaluate %q: %v", expr, err)
	}
	return v
}

// a2Click clicks the first element matching the selector (non-strict, like Page.click in Java).
func a2Click(ctx context.Context, page *playwright.Page, sel string, opts ...*playwright.LocatorClickOptions) error {
	return page.Locator(sel).First().Click(ctx, opts...)
}

// a2Fulfill fulfills a route with an HTML body.
func a2Fulfill(ctx context.Context, route *playwright.Route, body string) {
	_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body, ContentType: a2Ptr("text/html")})
}

// a2RouteIframe mirrors routeIframe from TestLocatorFrame.java.
func a2RouteIframe(ctx context.Context, page *playwright.Page) error {
	if err := page.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		a2Fulfill(ctx, route, "<iframe src='iframe.html' name='frame1'></iframe>")
	}); err != nil {
		return err
	}
	if err := page.Route(ctx, "**/iframe.html", func(route *playwright.Route) {
		a2Fulfill(ctx, route, `<html>
  <div>
    <button data-testid="buttonId">Hello iframe</button>
    <iframe src="iframe-2.html"></iframe>
  </div>
  <span>1</span>
  <span>2</span>
  <label for=target>Name</label><input id=target type=text placeholder=Placeholder title=Title alt=Alternative>
</html>`)
	}); err != nil {
		return err
	}
	return page.Route(ctx, "**/iframe-2.html", func(route *playwright.Route) {
		a2Fulfill(ctx, route, "<html><button>Hello nested iframe</button></html>")
	})
}

// a2RouteAmbiguous mirrors routeAmbiguous from TestLocatorFrame.java.
func a2RouteAmbiguous(ctx context.Context, page *playwright.Page) error {
	if err := page.Route(ctx, "**/empty.html", func(route *playwright.Route) {
		a2Fulfill(ctx, route, "<iframe src='iframe-1.html'></iframe>\n<iframe src='iframe-2.html'></iframe>\n<iframe src='iframe-3.html'></iframe>")
	}); err != nil {
		return err
	}
	return page.Route(ctx, "**/iframe-*", func(route *playwright.Route) {
		u, _ := url.Parse(route.Request().URL())
		a2Fulfill(ctx, route, "<html><button>Hello from "+strings.TrimPrefix(u.Path, "/")+"</button></html>")
	})
}

// a2Calls records the invocations of the exposed "stub" function.
type a2Calls struct {
	mu    sync.Mutex
	calls [][]any
}

func (c *a2Calls) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.calls)
}

func (c *a2Calls) Get(i int) []any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls[i]
}

func (c *a2Calls) All() [][]any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([][]any(nil), c.calls...)
}

// a2WaitCalls waits until n calls were recorded and then verifies there are exactly n.
func (c *a2Calls) WaitLen(t *testing.T, n int) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for c.Len() < n && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if n > 0 {
		time.Sleep(50 * time.Millisecond)
	}
	if got := c.Len(); got != n {
		t.Fatalf("expected %d calls, got %d", n, got)
	}
}

// a2ExposeStub exposes window.stub on the page, recording every call.
func a2ExposeStub(t *testing.T, ctx context.Context, page *playwright.Page) *a2Calls {
	t.Helper()
	calls := &a2Calls{}
	err := page.ExposeBinding(ctx, "stub", func(args ...any) any {
		calls.mu.Lock()
		defer calls.mu.Unlock()
		calls.calls = append(calls.calls, args)
		return nil
	})
	if err != nil {
		t.Fatalf("ExposeBinding: %v", err)
	}
	return calls
}

// a2ClockSetup creates a page with window.stub exposed and returns it with its clock.
func a2ClockSetup(t *testing.T) (context.Context, *testserver.Server, *playwright.Page, *a2Calls, *playwright.Clock) {
	t.Helper()
	ctx, srv, page := a2Setup(t)
	calls := a2ExposeStub(t, ctx, page)
	return ctx, srv, page, calls, page.Context().Clock()
}

// a2ClockPaused installs the fake clock at time 0 and pauses it at 1000 (setupRunForTest and friends).
func a2ClockPaused(t *testing.T, ctx context.Context, clock *playwright.Clock) {
	t.Helper()
	if err := clock.Install(ctx, 0); err != nil {
		t.Fatalf("clock.Install: %v", err)
	}
	if err := clock.PauseAt(ctx, 1000); err != nil {
		t.Fatalf("clock.PauseAt: %v", err)
	}
}

// a2TrackTapEvents mirrors trackEvents from TestTap.java.
func a2TrackTapEvents(ctx context.Context, target *playwright.ElementHandle) (*playwright.JSHandle, error) {
	return target.EvaluateHandle(ctx, `target => {
  const events = [];
  for (const event of [
    'mousedown', 'mouseenter', 'mouseleave', 'mousemove', 'mouseout', 'mouseover', 'mouseup', 'click',
    'pointercancel', 'pointerdown', 'pointerenter', 'pointerleave', 'pointermove', 'pointerout', 'pointerover', 'pointerup',
    'touchstart', 'touchend', 'touchmove', 'touchcancel'])
    target.addEventListener(event, () => events.push(event), false);
  return events;
}`)
}

// a2TouchPage creates a page in a context with touch enabled and a test server serving all resources.
func a2TouchPage(t *testing.T) (context.Context, *testserver.Server, *playwright.Page) {
	t.Helper()
	ctx := testCtx(t)
	srv := testserver.New(t)
	a2ServeAll(t, srv)
	return ctx, srv, a2NewPage(t, &playwright.BrowserContextOptions{HasTouch: a2Ptr(true)})
}

// a2CaptureLastKeyDown mirrors captureLastKeyDown from TestPageKeyboard.java, storing the event in window.__a2LastKeyDown.
func a2CaptureLastKeyDown(t *testing.T, ctx context.Context, page *playwright.Page) {
	t.Helper()
	a2Ev(t, ctx, page, `() => {
  const lastEvent = { repeat: false, location: -1, code: '', key: '', metaKey: false, keyIdentifier: 'unsupported' };
  document.addEventListener('keydown', e => {
    lastEvent.repeat = e.repeat;
    lastEvent.location = e.location;
    lastEvent.key = e.key;
    lastEvent.code = e.code;
    lastEvent.metaKey = e.metaKey;
    lastEvent.keyIdentifier = 'keyIdentifier' in e && e['keyIdentifier'];
  }, true);
  window.__a2LastKeyDown = lastEvent;
}`)
}
