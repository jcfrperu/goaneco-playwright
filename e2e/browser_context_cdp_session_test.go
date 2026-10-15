//go:build e2e

// E2E tests for: TestBrowserContextCDPSession.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"encoding/json"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Ref: TestBrowserContextCDPSession.java#shouldDetachSession  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldDetachSession(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err)

	_, err = session.Send(ctx, "Runtime.enable", nil)
	must.NoError(err)

	result, err := session.Send(ctx, "Runtime.evaluate", map[string]any{
		"expression":    "1 + 2",
		"returnByValue": true,
	})
	must.NoError(err)

	var resp struct {
		Result struct {
			Value any `json:"value"`
		} `json:"result"`
	}
	must.NoError(json.Unmarshal(result, &resp))
	is.Equal(float64(3), resp.Result.Value)

	err = session.Detach(ctx)
	must.NoError(err)

	_, err = session.Send(ctx, "Runtime.evaluate", map[string]any{
		"expression":    "1 + 2",
		"returnByValue": true,
	})
	is.Error(err, "expected error sending after Detach")
}

// Ref: TestBrowserContextCDPSession.java#shouldWorkWithMainFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldWorkWithMainFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<title>Main Frame</title><p>hello</p>`)
	must.NoError(err, "SetContent failed")

	frame := page.MainFrame()
	must.NotNil(frame, "MainFrame() returned nil")
	title, err := frame.Title(ctx)
	must.NoError(err, "frame.Title failed")
	if title != "Main Frame" {
		t.Errorf("frame.Title = %q, want 'Main Frame'", title)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCDPSession.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestCDPSessionSend(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	srv.ServeWithBody("/cdp", "text/html", `<title>CDP Test</title>`)
	err = page.Goto(ctx, srv.Prefix()+"/cdp")
	must.NoError(err, "Goto failed")

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err, "NewCDPSession failed")

	// Use Runtime.evaluate to run JS via CDP
	result, err := session.Send(ctx, "Runtime.evaluate", map[string]any{
		"expression":    "1 + 2",
		"returnByValue": true,
	})
	must.NoError(err, "CDP Runtime.evaluate failed")

	var resp struct {
		Result struct {
			Value any `json:"value"`
		} `json:"result"`
	}
	err = json.Unmarshal(result, &resp)
	must.NoError(err, "failed to parse CDP result")

	// JSON numbers decode as float64
	if resp.Result.Value != float64(3) {
		t.Errorf("CDP 1+2 = %v, want 3", resp.Result.Value)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCDPSession.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestCDPSessionOnClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	srv.ServeWithBody("/cdp-close", "text/html", `<p>close test</p>`)
	err = page.Goto(ctx, srv.Prefix()+"/cdp-close")
	must.NoError(err, "Goto failed")

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err, "NewCDPSession failed")

	closed := make(chan struct{}, 1)
	_ = session.OnClose(func() {
		select {
		case closed <- struct{}{}:
		default:
		}
	})

	// Detaching the session triggers the OnClose callback.
	err = session.Detach(ctx)
	must.NoError(err, "Detach failed")

	select {
	case <-closed:
		// OnClose was called — success.
	case <-time.After(3 * time.Second):
		t.Error("OnClose handler was not called after Detach")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextCDPSession.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestCDPSessionOn(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	srv.ServeWithBody("/cdp-event", "text/html", `<p>test</p>`)
	err = page.Goto(ctx, srv.Prefix()+"/cdp-event")
	must.NoError(err, "Goto failed")

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err, "NewCDPSession failed")

	// Just verify that registering an event handler and sending a command works
	cancel := session.On("Runtime.consoleAPICalled", func(params json.RawMessage) {
		_ = params // handler registered; we just verify no panic
	})
	defer cancel()

	// Enable Runtime domain to get console events
	_, err = session.Send(ctx, "Runtime.enable", nil)
	must.NoError(err, "Runtime.enable failed")

	err = session.Detach(ctx)
	must.NoError(err, "Detach failed")
}

// Ref: TestBrowserContextCDPSession.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err)

	_, err = session.Send(ctx, "Runtime.enable", nil)
	must.NoError(err)

	_, err = session.Send(ctx, "Runtime.evaluate", map[string]any{
		"expression": "window.foo = 'bar'",
	})
	must.NoError(err)

	foo, err := page.Evaluate(ctx, "window['foo']")
	must.NoError(err)
	is.Equal("bar", foo)
}

// Ref: TestBrowserContextCDPSession.java#shouldThrowNiceErrors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldThrowNiceErrors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err)

	_, err = session.Send(ctx, "ThisCommand.DoesNotExist", nil)
	is.Error(err)
	is.Contains(err.Error(), "ThisCommand.DoesNotExist")
}

// Ref: TestBrowserContextCDPSession.java#shouldNotBreakPageClose  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldNotBreakPageClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err)

	err = session.Detach(ctx)
	must.NoError(err)

	err = page.Close(ctx)
	must.NoError(err)
}

// Ref: TestBrowserContextCDPSession.java#shouldDetachWhenPageCloses  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldDetachWhenPageCloses(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)

	session, err := bCtx.NewCDPSession(ctx, page)
	must.NoError(err)

	closed := make(chan struct{}, 1)
	_ = session.OnClose(func() {
		select {
		case closed <- struct{}{}:
		default:
		}
	})

	err = page.Close(ctx)
	must.NoError(err)

	select {
	case <-closed:

	case <-time.After(3 * time.Second):
		t.Error("CDPSession OnClose was not called after page.Close()")
	}
}

// Ref: TestBrowserContextCDPSession.java#shouldSendEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldSendEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	session, err := page.Context().NewCDPSession(ctx, page)
	must.NoError(err)
	_, err = session.Send(ctx, "Network.enable", nil)
	must.NoError(err)

	var mu sync.Mutex
	count := 0
	cancel := session.On("Network.requestWillBeSent", func(params json.RawMessage) {
		var ev struct {
			Type string `json:"type"`
		}
		if json.Unmarshal(params, &ev) == nil && ev.Type == "Document" {
			mu.Lock()
			count++
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return count >= 1
	}, 5*time.Second, 20*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	must.Equal(1, count)
}

// Ref: TestBrowserContextCDPSession.java#shouldThrowIfTargetIsPartOfMain  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldThrowIfTargetIsPartOfMain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, `<iframe src="about:blank"></iframe>`))

	must.Eventually(func() bool { return len(page.Frames()) == 2 }, 5*time.Second, 20*time.Millisecond)
	frames := page.Frames()
	must.Len(frames, 2)
	is.Equal("about:blank", frames[1].URL())
	_, err := page.Context().NewCDPSession(ctx, frames[1])
	must.Error(err)
	is.Contains(err.Error(), "This frame does not have a separate CDP session, it is a part of the parent frame's session")
}

// Ref: TestBrowserContextCDPSession.java#shouldAddMultipleEventListeners  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldAddMultipleEventListeners(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	session, err := page.Context().NewCDPSession(ctx, page)
	must.NoError(err)
	_, err = session.Send(ctx, "Network.enable", nil)
	must.NoError(err)

	var count atomic.Int32
	listener := func(params json.RawMessage) {
		var ev struct {
			Type string `json:"type"`
		}
		if json.Unmarshal(params, &ev) == nil && ev.Type == "Document" {
			count.Add(1)
		}
	}
	c1 := session.On("Network.requestWillBeSent", listener)
	defer c1()
	c2 := session.On("Network.requestWillBeSent", listener)
	defer c2()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.Eventually(func() bool { return count.Load() >= 2 }, 5*time.Second, 20*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	must.Equal(int32(2), count.Load())
}

// Ref: TestBrowserContextCDPSession.java#shouldEmitEventForEachCDPEvent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldEmitEventForEachCDPEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	session, err := page.Context().NewCDPSession(ctx, page)
	must.NoError(err)
	_, err = session.Send(ctx, "Network.enable", nil)
	must.NoError(err)

	events := make(chan json.RawMessage, 64)
	off := session.On("event", func(event json.RawMessage) {
		select {
		case events <- event:
		default:
		}
	})
	defer off()
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var requestURL string
	observed := 0
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	for requestURL == "" {
		select {
		case raw := <-events:
			observed++
			var event struct {
				Method string `json:"method"`
				Params struct {
					Request struct {
						URL string `json:"url"`
					} `json:"request"`
				} `json:"params"`
			}
			must.NoError(json.Unmarshal(raw, &event))
			if event.Method == "Network.requestWillBeSent" && event.Params.Request.URL == srv.EmptyPage() {
				requestURL = event.Params.Request.URL
			}
		case <-timer.C:
			t.Fatal("no Network.requestWillBeSent event received")
		}
	}
	is.Greater(observed, 0)
	is.Equal(srv.EmptyPage(), requestURL)
}

// Ref: TestBrowserContextCDPSession.java#shouldRemoveEventListeners  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextCDPSession.java
func TestBrowserContextCDPSessionShouldRemoveEventListeners(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName != "chromium" {
		t.Skip("CDPSession is Chromium-only")
	}
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	session, err := page.Context().NewCDPSession(ctx, page)
	must.NoError(err)
	_, err = session.Send(ctx, "Network.enable", nil)
	must.NoError(err)

	var count atomic.Int32
	listener := func(params json.RawMessage) {
		var ev struct {
			Type string `json:"type"`
		}
		if json.Unmarshal(params, &ev) == nil && ev.Type == "Document" {
			count.Add(1)
		}
	}
	off1 := session.On("Network.requestWillBeSent", listener)
	off2 := session.On("Network.requestWillBeSent", listener)
	defer off2()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.Eventually(func() bool { return count.Load() >= 2 }, 5*time.Second, 20*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	must.Equal(int32(2), count.Load())

	off1()
	count.Store(0)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.Eventually(func() bool { return count.Load() >= 1 }, 5*time.Second, 20*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	must.Equal(int32(1), count.Load())
}
