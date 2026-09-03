//go:build e2e

package e2e

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
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

// -----------------------------------------------------------------------------
// Minimal RFC 6455 WebSocket server helpers
// -----------------------------------------------------------------------------

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// wsHandshake performs the WebSocket upgrade handshake.
// Returns the hijacked net.Conn and its buffered reader on success.
func wsHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, *bufio.ReadWriter, error) {
	key := r.Header.Get("Sec-Websocket-Key")
	if key == "" {
		http.Error(w, "missing Sec-WebSocket-Key", http.StatusBadRequest)
		return nil, nil, fmt.Errorf("missing Sec-WebSocket-Key")
	}
	h := sha1.New()
	io.WriteString(h, key+wsGUID)
	accept := base64.StdEncoding.EncodeToString(h.Sum(nil))

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijacking not supported", http.StatusInternalServerError)
		return nil, nil, fmt.Errorf("hijacking not supported")
	}
	conn, rw, err := hj.Hijack()
	if err != nil {
		return nil, nil, err
	}
	// Write the 101 Switching Protocols response directly on the raw conn.
	resp := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	if _, err = fmt.Fprint(conn, resp); err != nil {
		conn.Close()
		return nil, nil, err
	}
	return conn, rw, nil
}

// wsReadFrame reads one WebSocket frame and returns opcode + unmasked payload.
// Only handles text (1), binary (2), and close (8) frames.
func wsReadFrame(rw *bufio.ReadWriter) (opcode byte, payload []byte, err error) {
	header := make([]byte, 2)
	if _, err = io.ReadFull(rw, header); err != nil {
		return 0, nil, err
	}
	opcode = header[0] & 0x0F
	masked := (header[1] & 0x80) != 0
	payloadLen := int(header[1] & 0x7F)

	switch payloadLen {
	case 126:
		ext := make([]byte, 2)
		if _, err = io.ReadFull(rw, ext); err != nil {
			return 0, nil, err
		}
		payloadLen = int(binary.BigEndian.Uint16(ext))
	case 127:
		ext := make([]byte, 8)
		if _, err = io.ReadFull(rw, ext); err != nil {
			return 0, nil, err
		}
		payloadLen = int(binary.BigEndian.Uint64(ext))
	}

	var maskKey [4]byte
	if masked {
		if _, err = io.ReadFull(rw, maskKey[:]); err != nil {
			return 0, nil, err
		}
	}

	payload = make([]byte, payloadLen)
	if _, err = io.ReadFull(rw, payload); err != nil {
		return 0, nil, err
	}

	if masked {
		for i, b := range payload {
			payload[i] = b ^ maskKey[i%4]
		}
	}
	return opcode, payload, nil
}

// wsWriteTextFrame writes a single unmasked text frame to conn.
func wsWriteTextFrame(conn net.Conn, text string) error {
	payload := []byte(text)
	frame := make([]byte, 2+len(payload))
	frame[0] = 0x81 // FIN + text opcode
	frame[1] = byte(len(payload))
	copy(frame[2:], payload)
	_, err := conn.Write(frame)
	return err
}

// wsWriteCloseFrame sends a close frame with the given code (1000=normal).
func wsWriteCloseFrame(conn net.Conn, code uint16) error {
	frame := []byte{0x88, 2, byte(code >> 8), byte(code)}
	_, err := conn.Write(frame)
	return err
}

// wsServerEchoOnce serves a single WebSocket connection:
// waits for one text message, echoes "incoming", then closes.
func wsServerEchoOnce(w http.ResponseWriter, r *http.Request) {
	conn, rw, err := wsHandshake(w, r)
	if err != nil {
		return
	}
	defer conn.Close()

	// Read one frame (the client may or may not send one before we respond).
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck
	opcode, _, _ := wsReadFrame(rw)
	if opcode == 8 { // close
		return
	}
	_ = wsWriteTextFrame(conn, "incoming")
	// Give the client a moment to receive it.
	time.Sleep(200 * time.Millisecond)
}

// wsServerAcceptAndClose immediately closes after handshake (for close-event test).
func wsServerAcceptAndClose(w http.ResponseWriter, r *http.Request) {
	conn, rw, err := wsHandshake(w, r)
	if err != nil {
		return
	}
	defer conn.Close()
	// Wait for open event: browser sends a close frame after receiving open.
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck
	wsReadFrame(rw)                                       //nolint:errcheck
	wsWriteCloseFrame(conn, 1000)                         //nolint:errcheck
}

// wsServerEchoBack echoes each received text frame back, then closes.
func wsServerEchoBack(w http.ResponseWriter, r *http.Request) {
	conn, rw, err := wsHandshake(w, r)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second)) //nolint:errcheck
	for {
		opcode, _, err := wsReadFrame(rw)
		if err != nil || opcode == 8 {
			return
		}
		if opcode == 1 || opcode == 2 {
			_ = wsWriteTextFrame(conn, "incoming")
			return
		}
	}
}

// wsServerAcceptOnly upgrades and then stays open without sending anything.
func wsServerAcceptOnly(w http.ResponseWriter, r *http.Request) {
	conn, rw, err := wsHandshake(w, r)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) //nolint:errcheck
	wsReadFrame(rw)                                        //nolint:errcheck
}

// waitForWebSocket is a test helper: registers OnWebSocket, runs trigger, returns first ws.
func waitForWebSocket(t *testing.T, page *playwright.Page, trigger func()) *playwright.WebSocket {
	t.Helper()
	ch := make(chan *playwright.WebSocket, 1)
	cancel := page.OnWebSocket(func(ws *playwright.WebSocket) {
		select {
		case ch <- ws:
		default:
		}
	})
	defer cancel()

	trigger()

	select {
	case ws := <-ch:
		return ws
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for WebSocket")
		return nil
	}
}

// waitForClose blocks until the WebSocket fires its close event or times out.
func waitForClose(t *testing.T, ws *playwright.WebSocket) {
	t.Helper()
	closed := make(chan struct{}, 1)
	cancel := ws.OnClose(func() {
		select {
		case closed <- struct{}{}:
		default:
		}
	})
	defer cancel()

	if ws.IsClosed() {
		return
	}
	select {
	case <-closed:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for WebSocket close")
	}
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

// TestWebSocketShouldWork verifies that page.Evaluate can open a WebSocket,
// receive a message from the server, and return it.
// Ref: TestWebSocket.java#shouldWork
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
// Ref: TestWebSocket.java#shouldEmitCloseEvents
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
// Ref: TestWebSocket.java#shouldEmitFrameEvents
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

// TestWebSocketWaitForWebSocket verifies that page.WaitForWebSocket returns the
// WebSocket created during the trigger function.
// Ref: TestWebSocket.java#shouldNotHaveStrayErrorEvents
func TestWebSocketWaitForWebSocket(t *testing.T) {
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

// sortStrings sorts a string slice in place using insertion sort (avoids import of sort).
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}
