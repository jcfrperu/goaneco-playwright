//go:build e2e

// CDPSession tests are Chromium-only since CDP is a Chrome-specific protocol.
// Tests are skipped on non-Chromium browsers.
package e2e

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestCDPSessionWork(t *testing.T) {
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

func TestCDPSessionDetach(t *testing.T) {
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

func TestCDPSessionThrowNiceErrors(t *testing.T) {
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

func TestCDPSessionNotBreakPageClose(t *testing.T) {
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

func TestCDPSessionDetachOnPageClose(t *testing.T) {
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
