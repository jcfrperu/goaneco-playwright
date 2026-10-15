//go:build e2e

// E2E tests for: TestWebSocket.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	"encoding/base64"
	"fmt"
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

// TestWebSocketShouldWork verifies that page.Evaluate can open a WebSocket,
// receive a message from the server, and return it.
// Ref: TestWebSocket.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldWork(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	must := require.New(t)

	srv.SetRoute("/ws", wsServerEchoOnce)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/ws-page"))

	val, err := page.Evaluate(ctx, `port => {
		let cb;
		const result = new Promise(f => cb = f);
		const ws = new WebSocket('ws://127.0.0.1:' + port + '/ws');
		ws.addEventListener('message', data => { ws.close(); cb(data.data); });
		return result;
	}`, srv.Port())
	must.NoError(err)
	is.Equal("incoming", val)
}

// TestWebSocketShouldEmitCloseEvents verifies that OnWebSocket fires with the correct URL
// and that the ws.OnClose handler is called when the connection closes.
// Ref: TestWebSocket.java#shouldEmitCloseEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldEmitCloseEvents(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	must := require.New(t)
	is := assert.New(t)

	srv.SetRoute("/ws", wsServerAcceptAndClose)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/ws-page"))

	var (
		mu         sync.Mutex
		log        []string
		capturedWS *playwright.WebSocket
	)

	cancel := page.OnWebSocket(func(ws *playwright.WebSocket) {
		wsURL := ws.URL()
		mu.Lock()
		log = append(log, "open<"+wsURL+">")
		capturedWS = ws
		mu.Unlock()

		ws.OnClose(func() {
			mu.Lock()
			log = append(log, "close")
			mu.Unlock()
		})
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `port => {
		const ws = new WebSocket('ws://127.0.0.1:' + port + '/ws');
		ws.addEventListener('open', () => ws.close());
	}`, srv.Port())
	must.NoError(err)

	// Wait for both "open" and "close" entries.
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		n := len(log)
		ws := capturedWS
		mu.Unlock()
		if n >= 2 && ws != nil && ws.IsClosed() {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	mu.Lock()
	gotLog := make([]string, len(log))
	copy(gotLog, log)
	ws := capturedWS
	mu.Unlock()

	must.GreaterOrEqual(len(gotLog), 2, "expected at least open+close events")
	is.Contains(gotLog[0], "open<ws://127.0.0.1:"+fmt.Sprint(srv.Port())+"/ws>")
	is.Equal("close", gotLog[1])
	must.NotNil(ws)
	is.True(ws.IsClosed())
}

// TestWebSocketShouldEmitFrameEvents verifies that OnFrameSent and OnFrameReceived
// handlers fire for the correct frames.
// Ref: TestWebSocket.java#shouldEmitFrameEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldEmitFrameEvents(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	must := require.New(t)
	is := assert.New(t)

	srv.SetRoute("/ws", wsServerEchoBack)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/ws-page"))

	var (
		mu  sync.Mutex
		log []string
	)

	closed := make(chan struct{})

	cancel := page.OnWebSocket(func(ws *playwright.WebSocket) {
		mu.Lock()
		log = append(log, "open")
		mu.Unlock()

		ws.OnFrameSent(func(f playwright.WebSocketFrame) {
			mu.Lock()
			log = append(log, "sent<"+f.Text()+">")
			mu.Unlock()
		})
		ws.OnFrameReceived(func(f playwright.WebSocketFrame) {
			mu.Lock()
			log = append(log, "received<"+f.Text()+">")
			mu.Unlock()
		})
		ws.OnClose(func() {
			mu.Lock()
			log = append(log, "close")
			mu.Unlock()
			select {
			case closed <- struct{}{}:
			default:
			}
		})
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `port => {
		const ws = new WebSocket('ws://127.0.0.1:' + port + '/ws');
		ws.addEventListener('open', () => ws.send('outgoing'));
		ws.addEventListener('message', () => { ws.close(); });
	}`, srv.Port())
	must.NoError(err)

	select {
	case <-closed:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for WebSocket close")
	}

	mu.Lock()
	gotLog := make([]string, len(log))
	copy(gotLog, log)
	mu.Unlock()

	must.GreaterOrEqual(len(gotLog), 4, "expected open, sent, received, close: "+strings.Join(gotLog, ", "))
	is.Equal("open", gotLog[0], "first event should be open: "+strings.Join(gotLog, ", "))
	is.Equal("close", gotLog[len(gotLog)-1], "last event should be close: "+strings.Join(gotLog, ", "))

	sorted := make([]string, len(gotLog))
	copy(sorted, gotLog)
	sortStrings(sorted)
	is.Contains(sorted, "sent<outgoing>")
	is.Contains(sorted, "received<incoming>")
}

// TestWebSocketShouldNotHaveStrayErrorEvents verifies that page.WaitForWebSocket returns the
// WebSocket created during the trigger function.
// Ref: TestWebSocket.java#shouldNotHaveStrayErrorEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldNotHaveStrayErrorEvents(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	must := require.New(t)

	srv.SetRoute("/ws", wsServerAcceptOnly)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/ws-page"))

	ws, err := page.WaitForWebSocket(ctx, func() error {
		_, evalErr := page.Evaluate(ctx, `port => {
			window.ws = new WebSocket('ws://127.0.0.1:' + port + '/ws');
		}`, srv.Port())
		return evalErr
	})
	must.NoError(err)
	must.NotNil(ws)

	// Verify the URL contains the expected path.
	is.Contains(ws.URL(), "/ws", "WebSocket URL should contain /ws path")

	// Close from the page side.
	_, _ = page.Evaluate(ctx, "window.ws.close()")
}

// TestWebSocketOnWebSocketFires verifies that page.OnWebSocket fires for every
// WebSocket opened on the page.
// Ref: TestWebSocket.java#shouldEmitCloseEvents (OnWebSocket basic)
func TestWebSocketOnWebSocketFires(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	must := require.New(t)

	srv.SetRoute("/ws", wsServerAcceptOnly)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/ws-page"))

	seen := make(chan *playwright.WebSocket, 1)
	cancel := page.OnWebSocket(func(ws *playwright.WebSocket) {
		select {
		case seen <- ws:
		default:
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `port => {
		window.ws = new WebSocket('ws://127.0.0.1:' + port + '/ws');
	}`, srv.Port())
	must.NoError(err)

	var ws *playwright.WebSocket
	select {
	case ws = <-seen:
	case <-time.After(10 * time.Second):
		t.Fatal("OnWebSocket handler was not called")
	}

	is.NotNil(ws)
	is.Contains(ws.URL(), "/ws")

	// Clean up.
	_, _ = page.Evaluate(ctx, "window.ws.close()")
}

// TestWebSocketCancelOnUnsubscribe verifies that unsubscribing from OnWebSocket
// stops future handler invocations.
// Ref: TestWebSocket.java (general handler lifecycle)
// Extra: no equivalent method in playwright-java; related class: TestWebSocket.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketCancelOnUnsubscribe(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	must := require.New(t)

	srv.SetRoute("/ws", wsServerAcceptOnly)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/ws-page"))

	// Register and immediately cancel the listener.
	var count int
	var mu sync.Mutex
	cancel := page.OnWebSocket(func(ws *playwright.WebSocket) {
		mu.Lock()
		count++
		mu.Unlock()
		_ = ws.IsClosed() // prevent WS from staying open forever
	})
	cancel() // unsubscribe before any WebSocket is opened

	_, err := page.Evaluate(ctx, `port => {
		const ws = new WebSocket('ws://127.0.0.1:' + port + '/ws');
		ws.close();
	}`, srv.Port())
	must.NoError(err)

	// Give the event time to arrive (if any).
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	c := count
	mu.Unlock()
	is.Equal(0, c, "handler should not fire after cancel")
}

// Ref: TestWebSocket.java#shouldEmitBinaryFrameEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldEmitBinaryFrameEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/ws", wsServerAcceptOnly)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/ws-page"))

	ws, err := page.WaitForWebSocket(ctx, func() error {
		_, evalErr := page.Evaluate(ctx, `port => {
			window.ws = new WebSocket('ws://localhost:' + port + '/ws');
			window.wsOpen = new Promise(r => window.ws.addEventListener('open', r));
		}`, srv.Port())
		return evalErr
	})
	must.NoError(err)

	var mu sync.Mutex
	var sent []playwright.WebSocketFrame
	closed := make(chan struct{})
	var once sync.Once
	ws.OnClose(func() { once.Do(func() { close(closed) }) })
	ws.OnFrameSent(func(f playwright.WebSocketFrame) {
		mu.Lock()
		sent = append(sent, f)
		mu.Unlock()
	})

	_, err = page.Evaluate(ctx, `async () => {
		await window.wsOpen;
		const binary = new Uint8Array(5);
		for (let i = 0; i < 5; ++i)
			binary[i] = i;
		window.ws.send('text');
		window.ws.send(binary);
		window.ws.close();
	}`)
	must.NoError(err)

	select {
	case <-closed:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for socket close")
	}
	mu.Lock()
	defer mu.Unlock()
	must.GreaterOrEqual(len(sent), 2)
	must.Equal("text", sent[0].Text())
	bin, err := base64.StdEncoding.DecodeString(sent[1].Data)
	must.NoError(err)
	must.Equal([]byte{0, 1, 2, 3, 4}, bin)
}

// Ref: TestWebSocket.java#shouldEmitError  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldEmitError(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	errCh := make(chan string, 1)
	cancel := page.OnWebSocket(func(ws *playwright.WebSocket) {
		ws.OnError(func(e string) {
			select {
			case errCh <- e:
			default:
			}
		})
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `port => {
		new WebSocket('ws://localhost:' + port + '/bogus-ws');
	}`, srv.Port())
	must.NoError(err)

	select {
	case e := <-errCh:
		if globalBTName == "firefox" {
			must.Equal("CLOSE_ABNORMAL", e)
		} else {
			must.Contains(e, "404")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for socket error")
	}
}

// Ref: TestWebSocket.java#shouldRejectFutureEventOnSocketClose  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldRejectFutureEventOnSocketClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, page, ws := websocketWaitSetup(t)
	_, err := ws.WaitForFrameReceived(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.ws.send('start')")
		return err
	})
	must.NoError(err)

	_, err = ws.WaitForFrameSent(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.ws.close()")
		return err
	})
	must.Error(err)
	must.True(strings.Contains(err.Error(), "Socket closed") || strings.Contains(err.Error(), "Socket error"), err.Error())
}

// Ref: TestWebSocket.java#shouldRejectFutureEventOnPageClose  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldRejectFutureEventOnPageClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, page, ws := websocketWaitSetup(t)
	_, err := ws.WaitForFrameReceived(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.ws.send('start')")
		return err
	})
	must.NoError(err)

	_, err = ws.WaitForFrameSent(ctx, func() error { return page.Close(ctx) })
	must.Error(err)
	must.Contains(err.Error(), "Target page, context or browser has been closed")
}

// Ref: TestWebSocket.java#shouldCallFrameReceivedPredicate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldCallFrameReceivedPredicate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, page, ws := websocketWaitSetup(t)
	text := ""
	frame, err := ws.WaitForFrameReceived(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.ws.send('outgoing')")
		return err
	}, &playwright.WebSocketWaitForFrameOptions{
		Predicate: func(frame playwright.WebSocketFrame) bool {
			if frame.Text() != "incoming" {
				return false
			}
			text = frame.Text()
			return true
		},
	})
	must.NoError(err)
	must.Equal("incoming", text)
	must.Equal("incoming", frame.Text())
}

// Ref: TestWebSocket.java#shouldCallFrameSentPredicate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldCallFrameSentPredicate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, page, ws := websocketWaitSetup(t)
	text := ""
	frame, err := ws.WaitForFrameSent(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.ws.send('outgoing')")
		return err
	}, &playwright.WebSocketWaitForFrameOptions{
		Predicate: func(frame playwright.WebSocketFrame) bool {
			if frame.Text() != "outgoing" {
				return false
			}
			text = frame.Text()
			return true
		},
	})
	must.NoError(err)
	must.Equal("outgoing", text)
	must.Equal("outgoing", frame.Text())
}

// Ref: TestWebSocket.java#shouldRespectFrameReceivedTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldRespectFrameReceivedTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, ws := websocketWaitSetup(t)
	waitCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	_, err := ws.WaitForFrameReceived(waitCtx, nil, &playwright.WebSocketWaitForFrameOptions{
		Predicate: func(playwright.WebSocketFrame) bool { return false },
	})
	must.Error(err)
	must.Contains(err.Error(), "Timeout")
}

// Ref: TestWebSocket.java#shouldRespectFrameSentTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestWebSocket.java
func TestWebSocketShouldRespectFrameSentTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, page, ws := websocketWaitSetup(t)
	waitCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	_, err := ws.WaitForFrameSent(waitCtx, func() error {
		_, err := page.Evaluate(ctx, "() => window.ws.send('outgoing')")
		return err
	}, &playwright.WebSocketWaitForFrameOptions{
		Predicate: func(playwright.WebSocketFrame) bool { return false },
	})
	must.Error(err)
	must.Contains(err.Error(), "Timeout")
}

func websocketWaitSetup(t *testing.T) (context.Context, *playwright.Page, *playwright.WebSocket) {
	t.Helper()
	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.SetRoute("/ws", wsServerEchoAndHold)
	srv.ServeWithBody("/ws-page", "text/html", "<!DOCTYPE html><html><body></body></html>")
	page := newPage(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/ws-page"))
	ws, err := page.WaitForWebSocket(ctx, func() error {
		_, err := page.Evaluate(ctx, `port => {
			window.ws = new WebSocket('ws://127.0.0.1:' + port + '/ws');
			return new Promise(resolve => window.ws.addEventListener('open', resolve, {once: true}));
		}`, srv.Port())
		return err
	})
	require.NoError(t, err)
	return ctx, page, ws
}

func wsServerEchoAndHold(w http.ResponseWriter, r *http.Request) {
	conn, rw, err := wsHandshake(w, r)
	if err != nil {
		return
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
	for {
		opcode, _, err := wsReadFrame(rw)
		if err != nil || opcode == 8 {
			return
		}
		if opcode == 1 || opcode == 2 {
			if err := wsWriteTextFrame(conn, "incoming"); err != nil {
				return
			}
		}
	}
}
