package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
	"github.com/jcfrperu/goaneco-playwright/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockConn creates a connection whose transport responds to known methods via handlers.
// Unknown methods receive an empty success response so tests don't hang on implicit calls.
// Pass strict=true to make the test fail immediately if an unmapped method is called.
func newMockConn(t *testing.T, handlers map[string]func(id int, params json.RawMessage), strict ...bool) *connection.Connection {
	t.Helper()
	isStrict := len(strict) > 0 && strict[0]
	conn := connection.NewConnection()
	conn.SetTransportSend(func(payload []byte) error {
		var msg struct {
			ID     *int            `json:"id"`
			GUID   string          `json:"guid"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		require.NoError(t, json.Unmarshal(payload, &msg), "unmarshal IPC message")
		if msg.ID == nil {
			return nil
		}
		if h, ok := handlers[msg.Method]; ok {
			go h(*msg.ID, msg.Params)
		} else if isStrict {
			t.Errorf("strict mock: unexpected IPC method %q called with no handler", msg.Method)
		} else {
			// Default: empty success so the caller doesn't block.
			go conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{}}`, *msg.ID)))
		}
		return nil
	})
	return conn
}

func TestBrowserNewContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"newContext": func(id int, _ json.RawMessage) {
			resp := protocol.BrowserNewContextResponse{
				Context: struct {
					Guid string `json:"guid"`
				}{Guid: "mock-context-123"},
			}
			resultBytes, err := json.Marshal(resp)
			must.NoError(err)
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"BrowserContext","guid":"mock-context-123","initializer":{}}}`))
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":%s}`, id, resultBytes)))
		},
	})

	browser := &Browser{
		owner: ChannelOwner{conn: conn, guid: "browser-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	bCtx, err := browser.NewContext(ctx)
	must.NoError(err)
	is.Equal("mock-context-123", bCtx.owner.GUID())
}

func TestPageGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var gotURL string
	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"goto": func(id int, params json.RawMessage) {
			var req protocol.FrameGotoRequest
			must.NoError(json.Unmarshal(params, &req))
			gotURL = req.Url
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{}}`, id)))
		},
	})

	page := &Page{
		owner:     ChannelOwner{conn: conn, guid: "page-guid"},
		mainFrame: ChannelOwner{conn: conn, guid: "frame-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	must.NoError(page.Goto(ctx, "https://example.com"))
	is.Equal("https://example.com", gotURL)
}

func TestSendRequest_ServerError(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"someMethod": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(fmt.Sprintf(
				`{"id":%d,"error":{"error":{"name":"TimeoutError","message":"30000ms exceeded"}}}`, id,
			)))
		},
	})

	owner := ChannelOwner{conn: conn, guid: "test-guid"}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	_, err := owner.SendMessageRequest(ctx, "someMethod", nil)
	is.Error(err)
	is.ErrorContains(err, "TimeoutError")
	is.ErrorContains(err, "30000ms exceeded")
}

func TestClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	tests := []struct {
		name  string
		setup func(conn *connection.Connection) func(context.Context) error
	}{
		{
			name: "Browser",
			setup: func(conn *connection.Connection) func(context.Context) error {
				b := &Browser{owner: ChannelOwner{conn: conn, guid: "browser-guid"}, connected: true}
				return b.Close
			},
		},
		{
			name: "BrowserContext",
			setup: func(conn *connection.Connection) func(context.Context) error {
				bc := &BrowserContext{owner: ChannelOwner{conn: conn, guid: "context-guid"}}
				return bc.Close
			},
		},
		{
			name: "Page",
			setup: func(conn *connection.Connection) func(context.Context) error {
				p := &Page{owner: ChannelOwner{conn: conn, guid: "page-guid"}}
				return p.Close
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var conn *connection.Connection
			conn = newMockConn(t, map[string]func(int, json.RawMessage){
				"close": func(id int, _ json.RawMessage) {
					conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{}}`, id)))
				},
			})

			closeFn := tt.setup(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			t.Cleanup(cancel)

			must.NoError(closeFn(ctx))
		})
	}
}

func TestPlaywrightGetters(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	pw := &Playwright{
		initializer: playwrightInitializer{
			Chromium: browserTypeRef{GUID: "chrom-guid"},
			Firefox:  browserTypeRef{GUID: "ff-guid"},
			WebKit:   browserTypeRef{GUID: "wk-guid"},
		},
	}

	cr, err := pw.Chromium()
	must.NoError(err)
	is.Equal("chrom-guid", cr.owner.GUID())

	ff, err := pw.Firefox()
	must.NoError(err)
	is.Equal("ff-guid", ff.owner.GUID())

	wk, err := pw.WebKit()
	must.NoError(err)
	is.Equal("wk-guid", wk.owner.GUID())

	pw.initializer.Chromium.GUID = ""
	_, err = pw.Chromium()
	is.Error(err, "expected error for missing Chromium GUID")
}

func TestPageTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"title": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{"value":"Example Domain"}}`, id)))
		},
	})

	page := &Page{
		mainFrame: ChannelOwner{conn: conn, guid: "frame-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	title, err := page.Title(ctx)
	must.NoError(err)
	is.Equal("Example Domain", title)
}

func TestPageContentAndSetContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	const wantHTML = "<html><body><h1>Hello</h1></body></html>"
	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"setContent": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{}}`, id)))
		},
		"content": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{"value":%q}}`, id, wantHTML)))
		},
	})

	page := &Page{
		mainFrame: ChannelOwner{conn: conn, guid: "frame-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	must.NoError(page.SetContent(ctx, "<h1>Hello</h1>"))

	content, err := page.Content(ctx)
	must.NoError(err)
	is.Equal(wantHTML, content)
}
