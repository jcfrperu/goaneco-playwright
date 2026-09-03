//go:build e2e

// Multi-frame E2E tests: iframes, frame navigation, frame evaluation.
// Migration of: TestFrame.java
package e2e

import (
	"strings"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// waitForFrame polls page.Frame(name) until it appears or the deadline is reached.
func waitForFrame(t *testing.T, page *playwright.Page, name string) *playwright.Frame {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if f := page.Frame(name); f != nil {
			return f
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("frame %q never attached within 5 seconds", name)
	return nil
}

func TestPageMainFrame(t *testing.T) {
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

func TestFrameEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<p id="el">hello frame</p>`)
	must.NoError(err, "SetContent failed")

	val, err := page.MainFrame().Evaluate(ctx, "document.getElementById('el').textContent")
	must.NoError(err, "frame.Evaluate failed")
	if val != "hello frame" {
		t.Errorf("Evaluate = %q, want 'hello frame'", val)
	}
}

func TestFrameLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button id="btn">click me</button>`)
	must.NoError(err, "SetContent failed")

	text, err := page.MainFrame().Locator("#btn").InnerText(ctx)
	must.NoError(err, "frame locator InnerText failed")
	if text != "click me" {
		t.Errorf("InnerText = %q, want 'click me'", text)
	}
}

func TestFrameSetContentAndContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	frame := page.MainFrame()
	err := frame.SetContent(ctx, `<p id="p">frame content</p>`)
	must.NoError(err, "frame.SetContent failed")

	content, err := frame.Content(ctx)
	must.NoError(err, "frame.Content failed")
	if !strings.Contains(content, "frame content") {
		t.Errorf("Content does not contain 'frame content': %s", content)
	}
}

func TestChildFrameAttachedAndName(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/inner", "text/html", `<title>Inner</title><p>inner</p>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'childframe';
		f.src = '`+srv.Prefix()+`/inner';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	child := waitForFrame(t, page, "childframe")
	if child.Name() != "childframe" {
		t.Errorf("frame.Name = %q, want 'childframe'", child.Name())
	}
}

func TestChildFrameEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/frame-eval", "text/html", `<p id="msg">iframe text</p>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'evalframe';
		f.src = '`+srv.Prefix()+`/frame-eval';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	fr := waitForFrame(t, page, "evalframe")
	err = fr.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	val, err := fr.Evaluate(ctx, "document.getElementById('msg').textContent")
	must.NoError(err, "frame.Evaluate failed")
	if val != "iframe text" {
		t.Errorf("frame.Evaluate = %q, want 'iframe text'", val)
	}
}

func TestPageFramesSlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/f1", "text/html", "<p>f1</p>")
	srv.ServeWithBody("/f2", "text/html", "<p>f2</p>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		['f1','f2'].forEach(function(name) {
			const f = document.createElement('iframe');
			f.name = name;
			f.src = '`+srv.Prefix()+`/' + name;
			document.body.appendChild(f);
		});
	}`)
	must.NoError(err, "inject iframes failed")

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if len(page.Frames()) >= 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if n := len(page.Frames()); n < 2 {
		t.Errorf("expected at least 2 child frames, got %d", n)
	}
}

func TestFrameGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/nav1", "text/html", "<title>Nav1</title>")
	srv.ServeWithBody("/nav2", "text/html", "<title>Nav2</title>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto parent failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'gotoframe';
		f.src = '`+srv.Prefix()+`/nav1';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	fr := waitForFrame(t, page, "gotoframe")
	err = fr.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	err = fr.Goto(ctx, srv.Prefix()+"/nav2")
	must.NoError(err, "frame.Goto failed")

	title, err := fr.Title(ctx)
	must.NoError(err, "frame.Title failed")
	if title != "Nav2" {
		t.Errorf("frame.Title after Goto = %q, want 'Nav2'", title)
	}
}

func TestFrameQuerySelector(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/qs", "text/html", `<button id="btn">press</button>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto parent failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'qsframe';
		f.src = '`+srv.Prefix()+`/qs';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	fr := waitForFrame(t, page, "qsframe")
	err = fr.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	el, err := fr.QuerySelector(ctx, "#btn")
	must.NoError(err, "frame.QuerySelector failed")
	must.NotNil(el, "QuerySelector returned nil, expected element")

	text, err := el.InnerText(ctx)
	must.NoError(err, "InnerText failed")
	if text != "press" {
		t.Errorf("InnerText = %q, want 'press'", text)
	}
}

func TestFrameIsDetached(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/iframe-det", "text/html", "<p>inner</p>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'detachme';
		f.src = '`+srv.Prefix()+`/iframe-det';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	frame := waitForFrame(t, page, "detachme")
	err = frame.WaitForLoadState(ctx, "load")
	must.NoError(err, "WaitForLoadState failed")

	if frame.IsDetached() {
		t.Error("frame should not be detached before removal")
	}

	_, err = page.Evaluate(ctx, `() => {
		document.querySelector('iframe[name="detachme"]').remove();
	}`)
	must.NoError(err, "remove iframe failed")

	// Poll until IsDetached() returns true (event propagates asynchronously).
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if frame.IsDetached() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Error("frame should be detached after removal from DOM")
}

func TestFrameByURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/byurl", "text/html", "<p>by url</p>")

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto parent failed")

	_, err = page.Evaluate(ctx, `() => {
		const f = document.createElement('iframe');
		f.name = 'urlframe';
		f.src = '`+srv.Prefix()+`/byurl';
		document.body.appendChild(f);
	}`)
	must.NoError(err, "inject iframe failed")

	_ = waitForFrame(t, page, "urlframe")

	deadline := time.Now().Add(3 * time.Second)
	targetURL := srv.Prefix() + "/byurl"
	for time.Now().Before(deadline) {
		if fr := page.FrameByURL(targetURL); fr != nil {
			if fr.Name() != "urlframe" {
				t.Errorf("FrameByURL name = %q, want 'urlframe'", fr.Name())
			}
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
	t.Errorf("FrameByURL(%q) returned nil", targetURL)
}

// TestFrameTitle verifies that Frame.Title returns the title of the document loaded in that frame.
// Ref: TestFrameShouldReturnFrameTitle
func TestFrameTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Outer page embeds title.html (/title.html has <title>Woof-Woof</title>)
	srv.ServeWithBody("/frame-title-outer", "text/html",
		`<iframe name="titleframe" src="`+srv.Prefix()+`/title.html"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-title-outer")
	must.NoError(err, "Goto failed")

	child := waitForFrame(t, page, "titleframe")
	must.NotNil(child, "child frame not found")

	title, err := child.Title(ctx)
	must.NoError(err, "Frame.Title failed")
	is.Equal("Woof-Woof", title, "child frame title should match iframe document title")
}

func TestFrameURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frame-url-inner", "text/html", `<p>inner</p>`)
	srv.ServeWithBody("/frame-url-outer", "text/html",
		`<iframe name="urlframe" src="`+srv.Prefix()+`/frame-url-inner"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-url-outer")
	must.NoError(err, "Goto failed")

	child := waitForFrame(t, page, "urlframe")
	url := child.URL()
	if !strings.HasSuffix(url, "/frame-url-inner") {
		t.Errorf("Frame.URL() = %q, expected suffix '/frame-url-inner'", url)
	}
}
