//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

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
