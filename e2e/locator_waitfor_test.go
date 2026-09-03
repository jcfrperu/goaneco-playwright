//go:build e2e

// Locator.WaitFor E2E tests.
// Migration of: TestLocatorWaitFor.java
package e2e

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorWaitForAttached verifies WaitFor("attached") resolves when element is added to DOM.
// Ref: TestLocatorWaitFor.java#shouldWorkWithAttached
func TestLocatorWaitForAttached(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	// Add a span after a short delay
	_, err := page.Evaluate(ctx, `() => setTimeout(() => {
		document.querySelector('div').innerHTML = '<span>hello</span>';
	}, 100)`)
	must.NoError(err)

	must.NoError(page.Locator("span").WaitFor(ctx, "attached"))
	text, err := page.Locator("span").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("hello", *text)
}

// TestLocatorWaitForDetached verifies WaitFor("detached") resolves when element is removed from DOM.
// Ref: TestLocatorWaitFor.java#shouldWorkWithDetached
func TestLocatorWaitForDetached(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>target</span></div>`))

	// Remove the span after a short delay
	_, err := page.Evaluate(ctx, `() => setTimeout(() => {
		document.querySelector('span').remove();
	}, 100)`)
	must.NoError(err)

	must.NoError(page.Locator("span").WaitFor(ctx, "detached"))
}

// TestLocatorWaitForDefaultIsVisible verifies WaitFor() with no args waits for visible.
// Ref: TestLocatorWaitFor.java#shouldWorkWithVisible (default)
func TestLocatorWaitForDefaultIsVisible(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="display:none">hidden</div>`))

	// Make visible after a short delay
	_, err := page.Evaluate(ctx, `() => setTimeout(() => {
		document.querySelector('div').style.display = '';
	}, 100)`)
	must.NoError(err)

	// WaitFor with no args should wait for visible
	must.NoError(page.Locator("div").WaitFor(ctx))
	visible, err := page.Locator("div").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorWaitForHiddenState verifies WaitFor("hidden") resolves when element becomes hidden.
// Ref: TestLocatorWaitFor.java#shouldWorkWithHidden
func TestLocatorWaitForHiddenState(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>visible</div>`))

	// Hide it after a short delay
	_, err := page.Evaluate(ctx, `() => setTimeout(() => {
		document.querySelector('div').style.display = 'none';
	}, 100)`)
	must.NoError(err)

	must.NoError(page.Locator("div").WaitFor(ctx, "hidden"))
	visible, err := page.Locator("div").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorWaitForAlreadyAttached verifies WaitFor("attached") resolves immediately if element is already in DOM.
// Ref: TestLocatorWaitFor.java#shouldWorkWithAlreadyAttached
func TestLocatorWaitForAlreadyAttached(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span>already here</span></div>`))

	// Should resolve immediately since span is already attached
	must.NoError(page.Locator("span").WaitFor(ctx, "attached"))
}

func TestLocatorWaitForVisibleAfterDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="display:none">content</div>
	`))

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => document.getElementById('el').style.display = 'block'`)
	}()

	must.NoError(page.Locator("#el").WaitFor(ctx, "visible"))
}

func TestLocatorWaitForHiddenAfterDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">visible</div>`))

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => document.getElementById('el').style.display = 'none'`)
	}()

	must.NoError(page.Locator("#el").WaitFor(ctx, "hidden"))
}

func TestLocatorWaitForAttachedAfterDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"></div>`))

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => {
			const el = document.createElement('span');
			el.id = 'new';
			el.textContent = 'hello';
			document.getElementById('container').appendChild(el);
		}`)
	}()

	must.NoError(page.Locator("#new").WaitFor(ctx, "attached"))

	text, err := page.Locator("#new").InnerText(ctx)
	must.NoError(err)
	is.Equal("hello", text)
}

func TestLocatorWaitForDetachedAfterDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">to remove</div>`))

	go func() {
		time.Sleep(100 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => document.getElementById('el').remove()`)
	}()

	must.NoError(page.Locator("#el").WaitFor(ctx, "detached"))
}

// TestLocatorWaitForVisibleEx2 verifies WaitFor waits until element is visible.
// Ref: TestLocatorWaitFor.java#shouldWaitForVisible
func TestLocatorWaitForVisibleEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" style="display:none">Content</div><script>setTimeout(() => { document.getElementById('el').style.display = 'block' }, 50)</script>`))

	must.NoError(page.Locator("#el").WaitFor(ctx))
}

// TestLocatorWaitForHiddenEx2 verifies WaitFor with Hidden state waits until element hides.
// Ref: TestLocatorWaitFor.java#shouldWaitForHidden
func TestLocatorWaitForHiddenEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Visible</div><script>setTimeout(() => { document.getElementById('el').style.display = 'none' }, 50)</script>`))

	must.NoError(page.Locator("#el").WaitFor(ctx, "hidden"))
}

// TestLocatorWaitForAttachedEx2 verifies WaitFor with Attached state works for appended element.
// Ref: TestLocatorWaitFor.java#shouldWaitForAttached
func TestLocatorWaitForAttachedEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"></div><script>setTimeout(() => { const el = document.createElement('span'); el.id = 'dynamic'; document.getElementById('container').appendChild(el); }, 50)</script>`))

	must.NoError(page.Locator("#dynamic").WaitFor(ctx, "attached"))
}

// TestLocatorWaitForVisibleDelayedEx3 verifies WaitFor waits until element appears.
// Ref: TestLocatorWaitFor.java#shouldWaitForVisible
func TestLocatorWaitForVisibleDelayedEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="display:none;">Content</div>
		<script>
			setTimeout(function() {
				document.getElementById('el').style.display = '';
			}, 50);
		</script>
	`))

	must.NoError(page.Locator("#el").WaitFor(ctx, "visible"))

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorWaitForDetachedEx3 verifies WaitFor detects element removal.
// Ref: TestLocatorWaitFor.java#shouldWaitForDetached
func TestLocatorWaitForDetachedEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Remove me</div><script>setTimeout(() => { document.getElementById('el').remove() }, 50)</script>`))

	must.NoError(page.Locator("#el").WaitFor(ctx, "detached"))
}

// TestLocatorWaitForHiddenEx3 verifies WaitFor detects hidden state.
// Ref: TestLocatorWaitFor.java#shouldWaitForHidden
func TestLocatorWaitForHiddenEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Visible</div><script>setTimeout(() => { document.getElementById('el').style.display = 'none' }, 50)</script>`))

	must.NoError(page.Locator("#el").WaitFor(ctx, "hidden"))
}

// TestLocatorWaitForAlreadyPresent verifies WaitFor returns immediately when the element
// is already present in the DOM (no waiting required).
// Ref: TestLocatorWaitFor.java#shouldReturnImmediatelyIfVisible, #shouldWaitForButton,
//
//	#shouldWaitForInputAttached, #shouldWaitForParagraph, #shouldDefaultToVisible
func TestLocatorWaitForAlreadyPresent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	tests := []struct {
		name     string
		html     string
		selector string
		state    string // "visible", "attached", or "" for default
	}{
		{
			name:     "button already visible (default state)",
			html:     `<button id="btn">Already here</button>`,
			selector: "#btn",
			state:    "",
		},
		{
			name:     "div already visible (explicit visible)",
			html:     `<div id="el">Visible</div>`,
			selector: "#el",
			state:    "visible",
		},
		{
			name:     "button already visible (explicit visible)",
			html:     `<button id="btn">Ready</button>`,
			selector: "#btn",
			state:    "visible",
		},
		{
			name:     "input already attached",
			html:     `<input id="inp" type="text" value="ready">`,
			selector: "#inp",
			state:    "attached",
		},
		{
			name:     "paragraph already visible (default state)",
			html:     `<p id="p">Paragraph content</p>`,
			selector: "#p",
			state:    "",
		},
		{
			name:     "span already visible (default state)",
			html:     `<span id="s">Span text</span>`,
			selector: "#s",
			state:    "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			page := newPage(t)
			must.NoError(page.SetContent(testCtx(t), tc.html))
			must.NoError(page.Locator(tc.selector).WaitFor(testCtx(t), tc.state))
		})
	}
}
