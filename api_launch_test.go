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

func TestBrowserTypeLaunch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var conn *connection.Connection // forward declaration for closure
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"launch": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"Browser","guid":"browser-guid-123","initializer":{}}}`))
			resp := protocol.BrowserTypeLaunchResponse{
				Browser: struct {
					Guid string `json:"guid"`
				}{Guid: "browser-guid-123"},
			}
			resultBytes, err := json.Marshal(resp)
			must.NoError(err)
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":%s}`, id, resultBytes)))
		},
	})

	bt := &BrowserType{
		owner: ChannelOwner{conn: conn, guid: "browser-type-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	browser, err := bt.Launch(ctx, nil)
	must.NoError(err)
	is.Equal("browser-guid-123", browser.owner.GUID())
}

func TestBrowserTypeLaunch_WithOptions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var receivedParams browserTypeLaunchParamsWire
	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"launch": func(id int, params json.RawMessage) {
			must.NoError(json.Unmarshal(params, &receivedParams))
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"Browser","guid":"browser-guid-opts","initializer":{}}}`))
			resp := protocol.BrowserTypeLaunchResponse{
				Browser: struct {
					Guid string `json:"guid"`
				}{Guid: "browser-guid-opts"},
			}
			resultBytes, err := json.Marshal(resp)
			must.NoError(err)
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":%s}`, id, resultBytes)))
		},
	})

	bt := &BrowserType{
		owner: ChannelOwner{conn: conn, guid: "browser-type-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	headless := false
	slowMo := 100.0
	channel := "chrome"
	executablePath := "/path/to/chrome"
	args := []string{"--no-sandbox", "--disable-gpu"}
	env := map[string]string{"DEBUG": "pw:api"}

	browser, err := bt.Launch(ctx, &BrowserTypeLaunchOptions{
		Headless:       &headless,
		SlowMo:         &slowMo,
		Channel:        &channel,
		ExecutablePath: &executablePath,
		Args:           args,
		Env:            env,
	})
	must.NoError(err)
	is.Equal("browser-guid-opts", browser.owner.GUID())

	must.NotNil(receivedParams.Headless)
	is.Equal(false, *receivedParams.Headless)
	must.NotNil(receivedParams.SlowMo)
	is.Equal(100.0, *receivedParams.SlowMo)
	must.NotNil(receivedParams.Channel)
	is.Equal("chrome", *receivedParams.Channel)
	must.NotNil(receivedParams.ExecutablePath)
	is.Equal("/path/to/chrome", *receivedParams.ExecutablePath)
	is.Len(receivedParams.Args, 2)
	is.Equal("--no-sandbox", receivedParams.Args[0])
	is.Len(receivedParams.Env, 1)
	is.Equal("DEBUG", receivedParams.Env[0].Name)
	is.Equal("pw:api", receivedParams.Env[0].Value)
}

func TestBrowserNewContext_WithOptions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var receivedParams browserContextNewContextParamsWire
	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"newContext": func(id int, params json.RawMessage) {
			must.NoError(json.Unmarshal(params, &receivedParams))
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"BrowserContext","guid":"bctx-guid-opts","initializer":{}}}`))
			resp := protocol.BrowserNewContextResponse{
				Context: struct {
					Guid string `json:"guid"`
				}{Guid: "bctx-guid-opts"},
			}
			resultBytes, err := json.Marshal(resp)
			must.NoError(err)
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":%s}`, id, resultBytes)))
		},
	})

	browser := &Browser{
		owner:     ChannelOwner{conn: conn, guid: "browser-guid"},
		connected: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	ua := "custom-agent"
	loc := "en-US"
	headers := map[string]string{"X-Test": "value"}
	vp := &ViewportSize{Width: 800, Height: 600}
	acceptDownloads := true
	creds := &HttpCredentials{Username: "user1", Password: "pwd"}
	storage := &StorageState{
		Origins: []OriginStorage{
			{Origin: "http://localhost", LocalStorage: []LocalStorageEntry{{Name: "k", Value: "v"}}},
		},
	}

	bCtx, err := browser.NewContext(ctx, &BrowserContextOptions{
		UserAgent:        &ua,
		Locale:           &loc,
		ExtraHTTPHeaders: headers,
		Viewport:         vp,
		AcceptDownloads:  &acceptDownloads,
		HttpCredentials:  creds,
		StorageState:     storage,
	})
	must.NoError(err)
	is.Equal("bctx-guid-opts", bCtx.owner.GUID())

	must.NotNil(receivedParams.UserAgent)
	is.Equal("custom-agent", *receivedParams.UserAgent)
	must.NotNil(receivedParams.Locale)
	is.Equal("en-US", *receivedParams.Locale)
	must.NotNil(receivedParams.Viewport)
	is.Equal(800, receivedParams.Viewport.Width)
	must.NotNil(receivedParams.AcceptDownloads)
	is.Equal("accept", *receivedParams.AcceptDownloads)
	must.NotNil(receivedParams.HttpCredentials)
	is.Equal("user1", receivedParams.HttpCredentials.Username)
	must.NotNil(receivedParams.StorageState)
	is.Len(receivedParams.StorageState.Origins, 1)
	is.Len(receivedParams.ExtraHTTPHeaders, 1)
	is.Equal("X-Test", receivedParams.ExtraHTTPHeaders[0].Name)
}

func TestBrowserNewPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"newContext": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"BrowserContext","guid":"ctx-123","initializer":{}}}`))
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{"context":{"guid":"ctx-123"}}}`, id)))
		},
		"newPage": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"Frame","guid":"frame-123","initializer":{}}}`))
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"Page","guid":"page-123","initializer":{"mainFrame":{"guid":"frame-123"}}}}`))
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{"page":{"guid":"page-123"}}}`, id)))
		},
	})

	browser := &Browser{
		owner: ChannelOwner{conn: conn, guid: "browser-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	page, err := browser.NewPage(ctx)
	must.NoError(err)
	is.Equal("page-123", page.owner.GUID())
	is.Equal("frame-123", page.mainFrame.GUID())
}

// TestBrowserNewPage_CleanupOnFailure verifies that if NewPage fails after
// creating a BrowserContext, Browser.NewPage sends "close" to prevent orphaned contexts.
func TestBrowserNewPage_CleanupOnFailure(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	closeSent := make(chan struct{}, 1)
	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"newContext": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"BrowserContext","guid":"ctx-fail-123","initializer":{}}}`))
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{"context":{"guid":"ctx-fail-123"}}}`, id)))
		},
		"newPage": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(fmt.Sprintf(
				`{"id":%d,"error":{"error":{"name":"Error","message":"page creation failed"}}}`, id,
			)))
		},
		"close": func(id int, _ json.RawMessage) {
			closeSent <- struct{}{}
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{}}`, id)))
		},
	})

	browser := &Browser{
		owner: ChannelOwner{conn: conn, guid: "browser-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	_, err := browser.NewPage(ctx)
	is.Error(err)

	select {
	case <-closeSent:
		// cleanup succeeded
	case <-ctx.Done():
		t.Error("expected 'close' to be sent for orphaned context, timed out")
	}
}

// TestBrowserContext_NewPage directly tests BrowserContext.NewPage().
func TestBrowserContext_NewPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	var conn *connection.Connection
	conn = newMockConn(t, map[string]func(int, json.RawMessage){
		"newPage": func(id int, _ json.RawMessage) {
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"Frame","guid":"frame-ctx-1","initializer":{}}}`))
			conn.Dispatch([]byte(`{"guid":"","method":"__create__","params":{"type":"Page","guid":"page-ctx-1","initializer":{"mainFrame":{"guid":"frame-ctx-1"}}}}`))
			conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"result":{"page":{"guid":"page-ctx-1"}}}`, id)))
		},
	})

	bCtx := &BrowserContext{
		owner: ChannelOwner{conn: conn, guid: "context-direct-guid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	is.Equal("page-ctx-1", page.owner.GUID())
	is.Equal("frame-ctx-1", page.mainFrame.GUID())
}
