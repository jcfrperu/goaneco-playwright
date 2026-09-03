//go:build e2e

// E2E tests for Page.FrameLocator — locating elements inside iframes.
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPageFrameLocatorLocator verifies that FrameLocator.Locator() finds elements inside an iframe.
func TestPageFrameLocatorLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Serve the inner iframe content
	srv.ServeWithBody("/inner-frame", "text/html", `<div id="greeting">Hello from iframe</div>`)

	// Serve the outer page with an iframe
	srv.ServeWithBody("/frame-host", "text/html",
		`<iframe id="my-frame" src="`+srv.Prefix()+`/inner-frame"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-host")
	must.NoError(err, "Goto failed")

	// Locate the element inside the iframe
	fl := page.FrameLocator("#my-frame")
	greeting := fl.Locator("#greeting")

	text, err := greeting.InnerText(ctx)
	must.NoError(err, "FrameLocator.Locator().InnerText() failed")
	if text != "Hello from iframe" {
		t.Errorf("expected 'Hello from iframe', got %q", text)
	}
}

// TestPageFrameLocatorGetByText verifies that FrameLocator.GetByText() works inside an iframe.
func TestPageFrameLocatorGetByText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frame-text", "text/html", `<p>Buy Now</p><p>Learn More</p>`)
	srv.ServeWithBody("/frame-host-text", "text/html",
		`<iframe id="tf" src="`+srv.Prefix()+`/frame-text"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-host-text")
	must.NoError(err, "Goto failed")

	fl := page.FrameLocator("#tf")
	el := fl.GetByText("Buy Now")

	visible, err := el.IsVisible(ctx)
	must.NoError(err, "FrameLocator.GetByText().IsVisible() failed")
	if !visible {
		t.Error("expected 'Buy Now' to be visible inside iframe")
	}
}

// TestPageFrameLocatorFill verifies that FrameLocator allows filling inputs inside an iframe.
func TestPageFrameLocatorFill(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frame-form", "text/html", `<input id="name" type="text" />`)
	srv.ServeWithBody("/frame-host-form", "text/html",
		`<iframe id="ff" src="`+srv.Prefix()+`/frame-form"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/frame-host-form")
	must.NoError(err, "Goto failed")

	fl := page.FrameLocator("#ff")
	inp := fl.Locator("#name")

	err = inp.Fill(ctx, "playwright")
	must.NoError(err, "FrameLocator.Locator().Fill() failed")

	val, err := inp.InputValue(ctx)
	must.NoError(err, "InputValue() failed")
	if val != "playwright" {
		t.Errorf("expected input value 'playwright', got %q", val)
	}
}

// TestFrameLocatorGetByLabel verifies FrameLocator.GetByLabel finds labeled inputs inside iframes.
func TestFrameLocatorGetByLabel(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/fl-label-inner", "text/html", `
		<label for="email">Email</label>
		<input id="email" type="email" />
	`)
	srv.ServeWithBody("/fl-label-outer", "text/html",
		`<iframe id="fl-lbl" src="`+srv.Prefix()+`/fl-label-inner"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/fl-label-outer")
	must.NoError(err, "Goto failed")

	fl := page.FrameLocator("#fl-lbl")
	err = fl.GetByLabel("Email").Fill(ctx, "test@example.com")
	must.NoError(err, "FrameLocator.GetByLabel().Fill() failed")
	val, err := fl.GetByLabel("Email").InputValue(ctx)
	must.NoError(err, "InputValue() failed")
	if val != "test@example.com" {
		t.Errorf("GetByLabel input value = %q, want 'test@example.com'", val)
	}
}

// TestFrameLocatorGetByRole verifies FrameLocator.GetByRole finds elements by ARIA role inside iframes.
func TestFrameLocatorGetByRole(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/fl-role-inner", "text/html", `<button>Submit</button>`)
	srv.ServeWithBody("/fl-role-outer", "text/html",
		`<iframe id="fl-role" src="`+srv.Prefix()+`/fl-role-inner"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/fl-role-outer")
	must.NoError(err, "Goto failed")

	fl := page.FrameLocator("#fl-role")
	visible, err := fl.GetByRole("button").IsVisible(ctx)
	must.NoError(err, "FrameLocator.GetByRole().IsVisible() failed")
	if !visible {
		t.Error("button inside iframe should be visible via GetByRole")
	}
}

// TestFrameLocatorGetByTestId verifies FrameLocator.GetByTestId finds elements by data-testid inside iframes.
func TestFrameLocatorGetByTestId(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/fl-tid-inner", "text/html", `<span data-testid="badge">42</span>`)
	srv.ServeWithBody("/fl-tid-outer", "text/html",
		`<iframe id="fl-tid" src="`+srv.Prefix()+`/fl-tid-inner"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/fl-tid-outer")
	must.NoError(err, "Goto failed")

	fl := page.FrameLocator("#fl-tid")
	text, err := fl.GetByTestId("badge").InnerText(ctx)
	must.NoError(err, "FrameLocator.GetByTestId().InnerText() failed")
	if text != "42" {
		t.Errorf("GetByTestId text = %q, want '42'", text)
	}
}

// TestFrameLocatorFrameLocator verifies nested FrameLocator for iframes within iframes.
func TestFrameLocatorFrameLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/fl-nested-inner", "text/html", `<p id="deep">deep content</p>`)
	srv.ServeWithBody("/fl-nested-mid", "text/html",
		`<iframe id="inner-frame" src="`+srv.Prefix()+`/fl-nested-inner"></iframe>`)
	srv.ServeWithBody("/fl-nested-outer", "text/html",
		`<iframe id="mid-frame" src="`+srv.Prefix()+`/fl-nested-mid"></iframe>`)

	err := page.Goto(ctx, srv.Prefix()+"/fl-nested-outer")
	must.NoError(err, "Goto failed")

	// page → #mid-frame → #inner-frame → #deep
	text, err := page.FrameLocator("#mid-frame").FrameLocator("#inner-frame").Locator("#deep").InnerText(ctx)
	must.NoError(err, "nested FrameLocator.InnerText() failed")
	if text != "deep content" {
		t.Errorf("nested FrameLocator text = %q, want 'deep content'", text)
	}
}

// TestLocatorFrameLocator verifies that Locator.FrameLocator scopes lookups into an iframe.
func TestLocatorFrameLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/loc-fl-inner", "text/html", `<h1 id="h">Locator Frame</h1>`)
	srv.ServeWithBody("/loc-fl-outer", "text/html", `
		<div id="wrapper">
			<iframe id="loc-fr" src="`+srv.Prefix()+`/loc-fl-inner"></iframe>
		</div>
	`)

	err := page.Goto(ctx, srv.Prefix()+"/loc-fl-outer")
	must.NoError(err, "Goto failed")

	// Use Locator.FrameLocator instead of Page.FrameLocator
	text, err := page.Locator("#wrapper").FrameLocator("#loc-fr").Locator("#h").InnerText(ctx)
	must.NoError(err, "Locator.FrameLocator().Locator().InnerText() failed")
	if text != "Locator Frame" {
		t.Errorf("Locator.FrameLocator text = %q, want 'Locator Frame'", text)
	}
}

func TestFrameLocatorScopedCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul><li>a</li><li>b</li><li>c</li></ul>
	`))

	count, err := page.MainFrame().Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

func TestFrameLocatorInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="title">Hello Frame</h1>`))

	text, err := page.MainFrame().Locator("#title").InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello Frame", text)
}

func TestFrameLocatorClickWorks(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button onclick="document.getElementById('r').textContent='ok'">go</button>
		<div id="r"></div>
	`))

	must.NoError(page.MainFrame().Locator("button").Click(ctx))

	text, err := page.MainFrame().Locator("#r").InnerText(ctx)
	must.NoError(err)
	is.Equal("ok", text)
}

func TestFrameLocatorFillSetsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	must.NoError(page.MainFrame().Locator("input").Fill(ctx, "frame fill"))

	val, err := page.MainFrame().Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("frame fill", val)
}

func TestFrameLocatorIsVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="visible">shown</div>
		<div id="hidden" style="display:none">hidden</div>
	`))

	visibleResult, err := page.MainFrame().Locator("#visible").IsVisible(ctx)
	must.NoError(err)
	is.True(visibleResult)

	hiddenResult, err := page.MainFrame().Locator("#hidden").IsVisible(ctx)
	must.NoError(err)
	is.False(hiddenResult)
}
