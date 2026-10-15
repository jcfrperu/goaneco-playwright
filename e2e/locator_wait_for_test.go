//go:build e2e

// E2E tests for: TestLocatorWaitFor.java.
package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestLocatorWaitForShouldWorkWithVisible verifies WaitFor() with no args waits for visible.
// Ref: TestLocatorWaitFor.java#shouldWorkWithVisible (default)
func TestLocatorWaitForShouldWorkWithVisible(t *testing.T) {
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

// TestLocatorWaitForShouldReturnImmediatelyIfVisible verifies WaitFor returns immediately when the element
// is already present in the DOM (no waiting required).
// Ref: TestLocatorWaitFor.java#shouldReturnImmediatelyIfVisible, #shouldWaitForButton,
//
//	#shouldWaitForInputAttached, #shouldWaitForParagraph, #shouldDefaultToVisible
func TestLocatorWaitForShouldReturnImmediatelyIfVisible(t *testing.T) {
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
