//go:build e2e

// E2E tests for: TestLocatorFrame.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPageFrameLocatorGetByText verifies that FrameLocator.GetByText() works inside an iframe.
// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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
// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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

// Extra: no equivalent method in playwright-java; related class: TestLocatorFrame.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
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

// TestLocatorFrameLocatorFrameLocatorShouldWorkForIframe verifies that Locator.FrameLocator scopes lookups into an iframe.
// Ref: TestLocatorFrame.java#locatorFrameLocatorShouldWorkForIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameLocatorFrameLocatorShouldWorkForIframe(t *testing.T) {
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

// Ref: TestLocatorFrame.java#shouldWorkForIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldWorkForIframe(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button := page.FrameLocator("iframe").Locator("button")
	must.NoError(button.WaitFor(ctx))
	text, err := button.InnerText(ctx)
	must.NoError(err)
	assert.Equal(t, "Hello iframe", text)
	must.NoError(playwright.Expect(button).ToHaveText(ctx, "Hello iframe"))
	must.NoError(button.Click(ctx))
}

// Ref: TestLocatorFrame.java#shouldWorkForNestedIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldWorkForNestedIframe(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button := page.FrameLocator("iframe").FrameLocator("iframe").Locator("button")
	must.NoError(button.WaitFor(ctx))
	text, err := button.InnerText(ctx)
	must.NoError(err)
	assert.Equal(t, "Hello nested iframe", text)
	must.NoError(playwright.Expect(button).ToHaveText(ctx, "Hello nested iframe"))
	must.NoError(button.Click(ctx))
}

// Ref: TestLocatorFrame.java#shouldWorkForAnd  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldWorkForAnd(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	locator := page.FrameLocator("iframe").Locator("button")
	must.NoError(playwright.Expect(locator).ToHaveText(ctx, "Hello iframe"))
	text, err := locator.InnerText(ctx)
	must.NoError(err)
	assert.Equal(t, "Hello iframe", text)
	spans := page.FrameLocator("iframe").Locator("span")
	must.NoError(playwright.Expect(spans).ToHaveCount(ctx, 2))
}

// Ref: TestLocatorFrame.java#shouldWaitForFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldWaitForFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	err := page.FrameLocator("iframe").Locator("span").Click(ctx, &playwright.LocatorClickOptions{Timeout: a2Ptr(300.0)})
	must.Error(err)
	assert.Contains(t, err.Error(), `.contentFrame()`)
}

// Ref: TestLocatorFrame.java#shouldWaitForFrame2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldWaitForFrame2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	a2Ev(t, ctx, page, "url => setTimeout(() => location.href = url, 300)", srv.EmptyPage())
	must.NoError(page.FrameLocator("iframe").Locator("button").Click(ctx))
}

// Ref: TestLocatorFrame.java#shouldNotWaitForFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldNotWaitForFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(playwright.Expect(page.FrameLocator("iframe").Locator("span")).ToBeHidden(ctx))
}

// Ref: TestLocatorFrame.java#shouldNotWaitForFrame2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldNotWaitForFrame2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(playwright.Expect(page.FrameLocator("iframe").Locator("span")).Not().ToBeVisible(ctx))
}

// Ref: TestLocatorFrame.java#shouldNotWaitForFrame3  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldNotWaitForFrame3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(playwright.Expect(page.FrameLocator("iframe").Locator("span")).ToHaveCount(ctx, 0))
}

// Ref: TestLocatorFrame.java#shouldClickInLazyIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldClickInLazyIframe(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Route(ctx, "**/iframe.html", func(route *playwright.Route) {
		a2Fulfill(ctx, route, "<html><button>Hello iframe</button></html>")
	}))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	a2Ev(t, ctx, page, `setTimeout(() => {
      const iframe = document.createElement('iframe');
      document.body.appendChild(iframe);
      // navigate iframe
      setTimeout(() => iframe.src = 'iframe.html', 500);
    }, 500);`)
	button := page.FrameLocator("iframe").Locator("button")
	must.NoError(button.Click(ctx))
	must.NoError(playwright.Expect(button).ToHaveText(ctx, "Hello iframe"))
	text, err := button.InnerText(ctx)
	must.NoError(err)
	assert.Equal(t, "Hello iframe", text)
}

const a2ReattachScript = `setTimeout(() => {
        document.querySelector('iframe').remove();
        setTimeout(() => {
          const iframe = document.createElement('iframe');
          iframe.src = 'iframe-2.html';
          document.body.appendChild(iframe);
        }, 500);
      }, 500);`

// Ref: TestLocatorFrame.java#waitForShouldSurviveFrameReattach  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameWaitForShouldSurviveFrameReattach(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button := page.FrameLocator("iframe").Locator("button:has-text('Hello nested iframe')")
	a2Ev(t, ctx, page, a2ReattachScript)
	must.NoError(button.WaitFor(ctx))
}

// Ref: TestLocatorFrame.java#clickShouldSurviveFrameReattach  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameClickShouldSurviveFrameReattach(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button := page.FrameLocator("iframe").Locator("button:has-text('Hello nested iframe')")
	a2Ev(t, ctx, page, a2ReattachScript)
	must.NoError(button.Click(ctx))
}

// Ref: TestLocatorFrame.java#clickShouldSurviveIframeNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameClickShouldSurviveIframeNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button := page.FrameLocator("iframe").Locator("button:has-text('Hello nested iframe')")
	a2Ev(t, ctx, page, `setTimeout(() => {
        document.querySelector('iframe').src = 'iframe-2.html';
      }, 500);`)
	must.NoError(button.Click(ctx))
}

// Ref: TestLocatorFrame.java#shouldNonWorkForNonFrame  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameShouldNonWorkForNonFrame(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.SetContent(ctx, "<div></div>"))
	must.NoError(page.SetDefaultTimeout(ctx, 3000))
	err := page.FrameLocator("div").Locator("button").WaitFor(ctx)
	must.Error(err)
	assert.Contains(t, err.Error(), "<div></div>")
	assert.Contains(t, err.Error(), "<iframe> was expected")
}

// Ref: TestLocatorFrame.java#locatorFrameLocatorShouldThrowOnAmbiguity  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameLocatorFrameLocatorShouldThrowOnAmbiguity(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteAmbiguous(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button := page.Locator("body").FrameLocator("iframe").Locator("button")
	err := button.WaitFor(ctx)
	must.Error(err)
	assert.Contains(t, err.Error(), `strict mode violation:`)
	assert.Contains(t, err.Error(), `resolved to 3 elements`)
}

// Ref: TestLocatorFrame.java#locatorFrameLocatorShouldNotThrowOnFirstLastNth  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameLocatorFrameLocatorShouldNotThrowOnFirstLastNth(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteAmbiguous(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button1 := page.Locator("body").FrameLocator("iframe").First().Locator("button")
	must.NoError(playwright.Expect(button1).ToHaveText(ctx, "Hello from iframe-1.html"))
	button2 := page.Locator("body").FrameLocator("iframe").Nth(1).Locator("button")
	must.NoError(playwright.Expect(button2).ToHaveText(ctx, "Hello from iframe-2.html"))
	button3 := page.Locator("body").FrameLocator("iframe").Last().Locator("button")
	must.NoError(playwright.Expect(button3).ToHaveText(ctx, "Hello from iframe-3.html"))
}

// Ref: TestLocatorFrame.java#getByCoverage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameGetByCoverage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	button1 := page.FrameLocator("iframe").GetByRole(playwright.AriaRoleButton)
	button2 := page.FrameLocator("iframe").GetByText("Hello")
	button3 := page.FrameLocator("iframe").GetByTestId("buttonId")
	must.NoError(playwright.Expect(button1).ToHaveText(ctx, "Hello iframe"))
	must.NoError(playwright.Expect(button2).ToHaveText(ctx, "Hello iframe"))
	must.NoError(playwright.Expect(button3).ToHaveText(ctx, "Hello iframe"))
	input1 := page.FrameLocator("iframe").GetByLabel("Name")
	must.NoError(playwright.Expect(input1).ToHaveValue(ctx, ""))
	input2 := page.FrameLocator("iframe").GetByPlaceholder("Placeholder")
	must.NoError(playwright.Expect(input2).ToHaveValue(ctx, ""))
	input3 := page.FrameLocator("iframe").GetByAltText("Alternative")
	must.NoError(playwright.Expect(input3).ToHaveValue(ctx, ""))
	input4 := page.FrameLocator("iframe").GetByTitle("Title")
	must.NoError(playwright.Expect(input4).ToHaveValue(ctx, ""))
}

// Ref: TestLocatorFrame.java#locatorContentFrameShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameLocatorContentFrameShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	locator := page.Locator("iframe")
	frameLocator := locator.ContentFrame()
	button := frameLocator.Locator("button")
	text, err := button.InnerText(ctx)
	must.NoError(err)
	is.Equal("Hello iframe", text)
	must.NoError(playwright.Expect(button).ToHaveText(ctx, "Hello iframe"))
	must.NoError(button.Click(ctx))
}

// Ref: TestLocatorFrame.java#frameLocatorOwnerShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestLocatorFrame.java
func TestLocatorFrameFrameLocatorOwnerShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(a2RouteIframe(ctx, page))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	frameLocator := page.FrameLocator("iframe")
	locator := frameLocator.Owner()
	must.NoError(playwright.Expect(locator).ToBeVisible(ctx))
	name, err := locator.GetAttribute(ctx, "name")
	must.NoError(err)
	must.NotNil(name)
	is.Equal("frame1", *name)
}
