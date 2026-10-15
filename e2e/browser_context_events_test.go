//go:build e2e

// E2E tests for: TestBrowserContextEvents.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"sync"
	"testing"
	"time"
)

// TestBrowserContextOnRequestFiredForPage verifies BrowserContext.OnRequest fires for page navigations.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextOnRequestFiredForPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var urls []string

	cancel := bc.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(urls)
	is.Contains(urls[0], "empty.html")
}

// TestBrowserContextOnResponseFiredForPage verifies BrowserContext.OnResponse fires for page navigations.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextOnResponseFiredForPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var statuses []int

	cancel := bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statuses = append(statuses, resp.Status())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(statuses)
	is.Equal(200, statuses[0])
}

// TestBrowserContextOnRequestFinishedFiredForPage verifies BrowserContext.OnRequestFinished fires.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextOnRequestFinishedFiredForPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	done := make(chan struct{}, 1)
	var mu sync.Mutex
	var gotURL string

	cancel := bc.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			gotURL = req.URL()
			mu.Unlock()
			select {
			case done <- struct{}{}:
			default:
			}
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnRequestFinished")
	}

	mu.Lock()
	defer mu.Unlock()
	is.Equal(srv.EmptyPage(), gotURL)
}

// TestBrowserContextOnPageLoadFiresAfterNavigation verifies BrowserContext.OnPageLoad fires.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextOnPageLoadFiresAfterNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	done := make(chan struct{}, 1)

	cancel := bc.OnPageLoad(func(p *playwright.Page) {
		select {
		case done <- struct{}{}:
		default:
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageLoad")
	}
}

// TestBrowserContextOnPageFired verifies BrowserContext.OnPage fires when a new page is created.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextOnPageFired(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	done := make(chan *playwright.Page, 1)
	cancel := bc.OnPage(func(p *playwright.Page) {
		select {
		case done <- p:
		default:
		}
	})
	defer cancel()

	// Open a popup — this creates a new page in the same context.
	_, err = page.Evaluate(ctx, "() => window.open('about:blank')")
	must.NoError(err)

	select {
	case newPage := <-done:
		must.NotNil(newPage)
		is.Equal(bc, newPage.Context())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPage event")
	}
}

// TestBrowserContextOnRequestFiredEx verifies OnRequest fires for navigation.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextOnRequestFiredEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/req-page", "text/html", `<html><body>Request page</body></html>`)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(c)
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var captured string
	off := bc.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		defer mu.Unlock()
		if captured == "" {
			captured = req.URL()
		}
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/req-page"))

	mu.Lock()
	u := captured
	mu.Unlock()

	is.NotEmpty(u)
	is.Contains(u, "/req-page")
}

// TestBrowserContextOnResponseFiredEx verifies OnResponse fires after request.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextEvents.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextOnResponseFiredEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/resp-page", "text/html", `<html><body>Response page</body></html>`)

	bc, err := globalBrowser.NewContext(ctx)
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bc.Close(c)
	})

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var respURL string
	off := bc.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		defer mu.Unlock()
		if respURL == "" {
			respURL = resp.URL()
		}
	})
	t.Cleanup(off)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/resp-page"))

	mu.Lock()
	u := respURL
	mu.Unlock()

	is.NotEmpty(u)
	is.Contains(u, "/resp-page")
}

// TestBrowserContextEventsConsoleEventShouldWorkInPopup verifies that console messages from a popup page
// are captured by a context-level console handler.
// Ref: TestBrowserContextEvents.java#consoleEventShouldWorkInPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsConsoleEventShouldWorkInPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	msgCh := make(chan *playwright.ConsoleMessage, 1)
	cancelConsole := bc.OnConsole(func(m *playwright.ConsoleMessage) {
		if m.Text() == "hello" {
			select {
			case msgCh <- m:
			default:
			}
		}
	})
	defer cancelConsole()

	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "const win = window.open(''); win.console.log('hello');")
		return err
	})
	must.NoError(err)
	must.NotNil(popup)

	select {
	case msg := <-msgCh:
		is.Equal("hello", msg.Text())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for console message from popup")
	}
}

// TestBrowserContextEventsConsoleEventShouldWorkInPopup2 verifies console messages from a javascript: URL popup.
// Ref: TestBrowserContextEvents.java#consoleEventShouldWorkInPopup2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsConsoleEventShouldWorkInPopup2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	msgCh := make(chan *playwright.ConsoleMessage, 1)
	cancelConsole := bc.OnConsole(func(m *playwright.ConsoleMessage) {
		if m.Type() == "log" && m.Text() == "hello" {
			select {
			case msgCh <- m:
			default:
			}
		}
	})
	defer cancelConsole()

	popupCh := make(chan *playwright.Page, 1)
	cancelPage := bc.OnPage(func(pg *playwright.Page) {
		select {
		case popupCh <- pg:
		default:
		}
	})
	defer cancelPage()

	_, err = page.Evaluate(ctx, `async () => {
		const win = window.open('javascript:console.log("hello")');
		await new Promise(f => setTimeout(f, 0));
		win.close();
	}`)
	must.NoError(err)

	select {
	case <-popupCh:
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for popup page")
	}

	select {
	case msg := <-msgCh:
		is.Equal("hello", msg.Text())
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for console message from popup2")
	}
}

// TestBrowserContextEventsConsoleEventShouldWorkInImmediatelyClosedPopup verifies console messages
// from a popup that opens and immediately closes.
// Ref: TestBrowserContextEvents.java#consoleEventShouldWorkInImmediatelyClosedPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsConsoleEventShouldWorkInImmediatelyClosedPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	msgCh := make(chan *playwright.ConsoleMessage, 1)
	cancelConsole := bc.OnConsole(func(m *playwright.ConsoleMessage) {
		if m.Text() == "hello" {
			select {
			case msgCh <- m:
			default:
			}
		}
	})
	defer cancelConsole()

	_, err = page.Evaluate(ctx, `async () => {
		const win = window.open();
		win.console.log('hello');
		win.close();
	}`)
	must.NoError(err)

	select {
	case msg := <-msgCh:
		is.Equal("hello", msg.Text())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for console message from immediately-closed popup")
	}
}

// TestBrowserContextEventsDialogEventShouldWorkInPopup verifies that dialog events from a popup page
// are captured and handled by a page-level dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkInPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsDialogEventShouldWorkInPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := bc.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	resultCh := make(chan any, 1)
	popup, err := page.WaitForPopup(ctx, func() error {
		go func() {
			result, _ := page.Evaluate(ctx, "() => { const win = window.open(''); return win.prompt('hey?'); }")
			select {
			case resultCh <- result:
			default:
			}
		}()
		return nil
	})
	must.NoError(err)
	must.NotNil(popup)

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from popup")
	}

	select {
	case result := <-resultCh:
		is.Equal("hello", result)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for evaluate result")
	}
}

// TestBrowserContextEventsDialogEventShouldWorkInPopup2 verifies that a dialog opened via javascript: URL
// is captured by a context-level dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkInPopup2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsDialogEventShouldWorkInPopup2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := bc.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	_, err = page.Evaluate(ctx, `window.open('javascript:prompt("hey?")');`)
	must.NoError(err)

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
		is.Nil(dialog.Page())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from javascript: URL popup")
	}
}

// TestBrowserContextEventsDialogEventShouldWorkInImmdiatelyClosedPopup verifies that a dialog from a popup
// that immediately closes is captured by a page-level dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkInImmdiatelyClosedPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsDialogEventShouldWorkInImmdiatelyClosedPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := bc.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	resultCh := make(chan any, 1)
	popup, err := page.WaitForPopup(ctx, func() error {
		go func() {
			result, _ := page.Evaluate(ctx, `async () => {
				const win = window.open();
				const result = win.prompt('hey?');
				win.close();
				return result;
			}`)
			select {
			case resultCh <- result:
			default:
			}
		}()
		return nil
	})
	must.NoError(err)
	must.NotNil(popup)

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from immediately-closed popup")
	}

	select {
	case result := <-resultCh:
		is.Equal("hello", result)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for evaluate result")
	}
}

// TestBrowserContextEventsDialogEventShouldWorkWithInlineScriptTag verifies that a dialog triggered by an
// inline <script> in a popup loaded via link click is captured by a context dialog handler.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkWithInlineScriptTag  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsDialogEventShouldWorkWithInlineScriptTag(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	srv.SetRoute("/popup.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`<script>window.result = prompt('hey?')</script>`))
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, "<a href='"+srv.Prefix()+"/popup.html' target=_blank>Click me</a>"))

	dialogCh := make(chan *playwright.Dialog, 1)
	cancelDialog := bc.OnDialog(func(d *playwright.Dialog) {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer bgCancel()
		_ = d.Accept(bgCtx, "hello")
		select {
		case dialogCh <- d:
		default:
		}
	})
	defer cancelDialog()

	popupCh := make(chan *playwright.Page, 1)
	cancelPage := bc.OnPage(func(pg *playwright.Page) {
		select {
		case popupCh <- pg:
		default:
		}
	})
	defer cancelPage()

	must.NoError(page.Locator("a").Click(ctx))

	var popup *playwright.Page
	select {
	case popup = <-popupCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for popup page")
	}

	select {
	case dialog := <-dialogCh:
		is.Equal("hey?", dialog.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for dialog from inline script")
	}

	// Wait for window.result to be set to "hello" in popup
	is.Eventually(func() bool {
		result, err := popup.Evaluate(ctx, "window.result")
		return err == nil && result == "hello"
	}, 5*time.Second, 100*time.Millisecond, "timed out waiting for window.result to be 'hello'")
}

// TestBrowserContextEventsConsoleEventShouldWorkSmoke verifies BrowserContext.OnConsole receives page console messages.
// Ref: TestBrowserContextEvents.java#consoleEventShouldWorkSmoke  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsConsoleEventShouldWorkSmoke(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := a4NewContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	ch := make(chan *playwright.ConsoleMessage, 4)
	off := bCtx.OnConsole(func(m *playwright.ConsoleMessage) {
		select {
		case ch <- m:
		default:
		}
	})
	defer off()
	_, err = page.Evaluate(ctx, "console.log('hello')")
	must.NoError(err)

	select {
	case m := <-ch:
		is.Equal("hello", m.Text())
		is.Equal(page, m.Page())
	case <-time.After(5 * time.Second):
		t.Fatal("no console message received")
	}
}

// TestBrowserContextEventsDialogEventShouldWorkSmoke verifies BrowserContext.onDialog.
// Ref: TestBrowserContextEvents.java#dialogEventShouldWorkSmoke  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsDialogEventShouldWorkSmoke(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	dialogs := make(chan *playwright.Dialog, 1)
	off := bCtx.OnDialog(func(dialog *playwright.Dialog) {
		if err := dialog.Accept(ctx, "hello"); err != nil {
			t.Errorf("accept dialog: %v", err)
		}
		dialogs <- dialog
	})
	defer off()

	result, err := page.Evaluate(ctx, "prompt('hey?')")
	must.NoError(err)
	is.Equal("hello", result)
	select {
	case dialog := <-dialogs:
		is.Equal("hey?", dialog.Message())
		is.Equal(page, dialog.Page())
	case <-time.After(5 * time.Second):
		t.Fatal("no context dialog event received")
	}
}

// TestBrowserContextEventsDialogClosedEventShouldWork verifies the close event is shared by page and context listeners.
// Ref: TestBrowserContextEvents.java#dialogClosedEventShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsDialogClosedEventShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	dialogSeen := make(chan *playwright.Dialog, 1)
	contextClosed := make(chan *playwright.Dialog, 1)
	pageClosed := make(chan *playwright.Dialog, 1)
	offDialog := bCtx.OnDialog(func(dialog *playwright.Dialog) {
		if err := dialog.Accept(ctx, "hello"); err != nil {
			t.Errorf("accept dialog: %v", err)
		}
		dialogSeen <- dialog
	})
	defer offDialog()
	offContextClosed := bCtx.OnDialogClosed(func(dialog *playwright.Dialog) { contextClosed <- dialog })
	defer offContextClosed()
	offPageClosed := page.OnDialogClosed(func(dialog *playwright.Dialog) { pageClosed <- dialog })
	defer offPageClosed()

	result, err := page.Evaluate(ctx, "prompt('hey?')")
	must.NoError(err)
	is.Equal("hello", result)

	var opened *playwright.Dialog
	select {
	case opened = <-dialogSeen:
		is.Equal("hey?", opened.Message())
	case <-time.After(5 * time.Second):
		t.Fatal("dialog event not received")
	}
	select {
	case closed := <-contextClosed:
		is.Same(opened, closed)
	case <-time.After(5 * time.Second):
		t.Fatal("context dialog closed event not received")
	}
	select {
	case closed := <-pageClosed:
		is.Same(opened, closed)
	case <-time.After(5 * time.Second):
		t.Fatal("page dialog closed event not received")
	}
}

// TestBrowserContextEventsWeberrorEventShouldIncludeLocation verifies the weberror event includes a location.
// Ref: TestBrowserContextEvents.java#weberrorEventShouldIncludeLocation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsWeberrorEventShouldIncludeLocation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	srv.ServeWithBody("/error.js", "application/javascript", "\nfunction foo() {\n  throw new Error('boom');\n}\nfoo();\n")
	srv.ServeWithBody("/error.html", "text/html", "<script src=\"/error.js\"></script>")

	var mu sync.Mutex
	var gotError *playwright.WebError
	cancel := bCtx.OnWebError(func(we *playwright.WebError) {
		mu.Lock()
		gotError = we
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.Prefix()+"/error.html"))
	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return gotError != nil
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(gotError)
	must.NotNil(gotError.Location())
	is.Equal(srv.Prefix()+"/error.js", gotError.Location().URL)
	is.Equal(2, gotError.Location().Line)
	is.Greater(gotError.Location().Column, 0)
}

// TestBrowserContextEventsFrameNavigatedEventShouldWork verifies the frame navigated event.
// Ref: TestBrowserContextEvents.java#frameNavigatedEventShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsFrameNavigatedEventShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var navigated *playwright.Frame
	cancel := bCtx.OnFrameNavigated(func(f *playwright.Frame) {
		mu.Lock()
		navigated = f
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return navigated != nil
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	is.Equal(page.MainFrame(), navigated)
	is.Equal(srv.EmptyPage(), navigated.URL())
}

// TestBrowserContextEventsFrameAttachedEventShouldWork verifies the frame attached event.
// Ref: TestBrowserContextEvents.java#frameAttachedEventShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsFrameAttachedEventShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	var mu sync.Mutex
	var attached *playwright.Frame
	cancel := bCtx.OnFrameAttached(func(f *playwright.Frame) {
		mu.Lock()
		attached = f
		mu.Unlock()
	})
	defer cancel()

	_, err = page.Evaluate(ctx, `() => {
		const iframe = document.createElement('iframe');
		iframe.src = 'about:blank';
		document.body.appendChild(iframe);
	}`)
	must.NoError(err)
	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return attached != nil
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(attached)
	is.Equal(page.MainFrame(), attached.ParentFrame())
}

// TestBrowserContextEventsFrameDetachedEventShouldWork verifies the frame detached event.
// Ref: TestBrowserContextEvents.java#frameDetachedEventShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsFrameDetachedEventShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, `() => {
		const iframe = document.createElement('iframe');
		iframe.id = 'x';
		iframe.src = 'about:blank';
		document.body.appendChild(iframe);
	}`)
	must.NoError(page.Locator("iframe").WaitFor(ctx))

	var mu sync.Mutex
	var detached *playwright.Frame
	cancel := bCtx.OnFrameDetached(func(f *playwright.Frame) {
		mu.Lock()
		detached = f
		mu.Unlock()
	})
	defer cancel()

	_, err = page.Evaluate(ctx, "() => document.getElementById('x').remove()")
	must.NoError(err)
	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return detached != nil
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(detached)
	is.Equal(page.MainFrame(), detached.ParentFrame())
}

// TestBrowserContextEventsPageErrorEventShouldWork is not ported.
// Ref: TestBrowserContextEvents.java#pageErrorEventShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsPageErrorEventShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var gotError *playwright.WebError
	cancel := bCtx.OnWebError(func(we *playwright.WebError) {
		mu.Lock()
		gotError = we
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.SetContent(ctx, "<script>throw new Error('boom')</script>"))
	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return gotError != nil
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(gotError)
	is.Equal(page, gotError.Page())
	is.Contains(gotError.Error(), "boom")
}

// TestBrowserContextEventsPageLoadEventShouldWork verifies BrowserContext.OnPageLoad fires for a page navigation.
// Ref: TestBrowserContextEvents.java#pageLoadEventShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsPageLoadEventShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	ch := make(chan *playwright.Page, 4)
	off := bCtx.OnPageLoad(func(p *playwright.Page) {
		select {
		case ch <- p:
		default:
		}
	})
	defer off()
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	select {
	case p := <-ch:
		must.Equal(page, p)
	case <-time.After(10 * time.Second):
		t.Fatal("page load event not received")
	}
}

// TestBrowserContextEventsPageCloseEventShouldWork verifies BrowserContext.OnPageClose fires when a page closes.
// Ref: TestBrowserContextEvents.java#pageCloseEventShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextEvents.java
func TestBrowserContextEventsPageCloseEventShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	newPage, err := bCtx.NewPage(ctx)
	must.NoError(err)
	ch := make(chan *playwright.Page, 4)
	off := bCtx.OnPageClose(func(p *playwright.Page) {
		select {
		case ch <- p:
		default:
		}
	})
	defer off()
	must.NoError(newPage.Close(ctx))
	select {
	case p := <-ch:
		must.Equal(newPage, p)
	case <-time.After(10 * time.Second):
		t.Fatal("page close event not received")
	}
}
