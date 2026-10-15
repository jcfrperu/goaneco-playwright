//go:build e2e

// E2E tests for: TestPageEventPopup.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

// TestPageEventPopupShouldWorkWithClickingTargetblank verifies that clicking a link with target="_blank" opens a popup,
// the opener is accessible from the popup, and the popup belongs to the correct frame.
// Ref: TestPageEventPopup.java#shouldWorkWithClickingTarget_blank  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPopup.java
func TestPageEventPopupShouldWorkWithClickingTargetblank(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)
	srv.ServeWithBody("/one-style.html", "text/html", `<!DOCTYPE html><html><head><title>Style</title></head><body></body></html>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<a target=_blank rel='opener' href='/one-style.html'>yo</a>`))

	popup := waitForPopup(t, page, func() {
		must.NoError(page.Locator("a").Click(ctx))
	})

	must.NotNil(popup)

	pageHasOpener, err := page.Evaluate(ctx, "() => !!window.opener")
	must.NoError(err)
	is.Equal(false, pageHasOpener, "page.opener should be null for the original tab")

	popupHasOpener, err := popup.Evaluate(ctx, "() => !!window.opener")
	must.NoError(err)
	is.Equal(true, popupHasOpener, "popup.opener should be non-null")

	// Verify the popup has a main frame (Go API does not expose Frame.Page())
	is.NotNil(popup.MainFrame(), "popup.MainFrame() should not be nil")
}

// --- TestPopup tests ---

// TestPopupWindowOpenerAccess verifies that window.opener is accessible from the popup
// when rel="opener" is set on the link.
// Ref: TestPageEventPopup.java#shouldWorkWithClickingTarget_blank
func TestPopupWindowOpenerAccess(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	is := assert.New(t)
	must := require.New(t)

	srv := testserver.New(t)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<button onclick="window.open('about:blank', '_blank')">open</button>`))

	popup := waitForPopup(t, page, func() {
		must.NoError(page.Locator("button").Click(ctx))
	})

	must.NotNil(popup)

	// The popup opened via window.open() should have window.opener set
	popupHasOpener, err := popup.Evaluate(ctx, "() => !!window.opener")
	must.NoError(err)
	is.Equal(true, popupHasOpener, "popup should have window.opener")
}

// TestPopupOnPopupEvent verifies that page.OnPopup fires with the correct popup Page
// when window.open() is called.
// Ref: TestPageEventPopup.java / TestPopup.java general
// Extra: no equivalent method in playwright-java; related class: TestPageEventPopup.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventPopup.java
func TestPopupOnPopupEvent(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	must := require.New(t)

	srv := testserver.New(t)
	srv.ServeWithBody("/popup-target.html", "text/html", `<title>Popup Target</title>`)

	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var receivedPopup *playwright.Page

	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		receivedPopup = p
		mu.Unlock()
	})
	defer cancel()

	_, err := page.Evaluate(ctx, "url => window.open(url)", srv.Prefix()+"/popup-target.html")
	must.NoError(err)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		p := receivedPopup
		mu.Unlock()
		if p != nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	mu.Lock()
	p := receivedPopup
	mu.Unlock()
	must.NotNil(p, "OnPopup should have fired with a non-nil page")
}
