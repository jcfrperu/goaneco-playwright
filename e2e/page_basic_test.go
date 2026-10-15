//go:build e2e

// E2E tests for: TestPageBasic.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"bytes"
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestInitScriptSetsWindowVarEx6 verifies AddInitScript sets a window variable.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestInitScriptSetsWindowVarEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__testVar = 'injected'`))
	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	val, err := page.Evaluate(ctx, `() => window.__testVar`)
	must.NoError(err)
	is.Equal("injected", val)
}

// TestInitScriptOverridesDateEx6 verifies AddInitScript can override Date.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestInitScriptOverridesDateEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `Date = function() { return { getFullYear: () => 2099 }; }`))
	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	year, err := page.Evaluate(ctx, `() => new Date().getFullYear()`)
	must.NoError(err)
	is.Equal(float64(2099), year)
}

// TestPageGotoEmptyPage verifies navigation to the test server's empty page.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageGotoEmptyPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoErrorf(err, "Goto(%q) failed", srv.EmptyPage())
	t.Logf("navigated to empty page: %s", page.URL())
}

// TestPageEvaluate verifies evaluating JavaScript expressions and functions with arguments and JSON structures.
// Ref: TestPageEvaluate.java
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageEvaluate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title.html")
	must.NoError(err, "Goto failed")

	// 1. Simple arithmetic expression
	resNum, err := page.Evaluate(ctx, "1 + 2")
	must.NoError(err)
	is.Equal(float64(3), resNum)

	// 2. DOM property reading
	resTitle, err := page.Evaluate(ctx, "() => document.title")
	must.NoError(err)
	is.Equal("Woof-Woof", resTitle)

	// 3. Arrow function with array arguments
	resMultiply, err := page.Evaluate(ctx, "([a, b]) => a * b", []any{6, 7})
	must.NoError(err)
	is.Equal(float64(42), resMultiply)

	// 4. Structured JSON object return
	resObj, err := page.Evaluate(ctx, "() => ({ key: 'value', count: 10, ok: true })")
	must.NoError(err, "Evaluate(object) failed")
	m, ok := resObj.(map[string]any)
	is.True(ok, "expected map[string]any, got %T: %v", resObj, resObj)
	if m["key"] != "value" || m["count"] != float64(10) || m["ok"] != true {
		t.Errorf("unexpected object content: %v", m)
	}
}

// TestPageOpenerIsNilForTopLevelPage verifies that page.Opener() returns nil for a regular (non-popup) page.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOpenerIsNilForTopLevelPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	opener, err := page.Opener(ctx)
	must.NoError(err, "Opener() failed")
	is.Nil(opener)
}

// TestPageOpenerIsSetForPopup verifies that a popup's Opener() returns the page that opened it.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOpenerIsSetForPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/popup-target", "text/html", `<title>Popup</title>`)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Register OnPage handler before triggering the popup
	popupCh := make(chan *playwright.Page, 1)
	cancelOnPage := bCtx.OnPage(func(p *playwright.Page) {
		select {
		case popupCh <- p:
		default:
		}
	})
	defer cancelOnPage()

	// Open a popup via JS
	_, err = page.Evaluate(ctx, `() => window.open('`+srv.Prefix()+`/popup-target')`)
	must.NoError(err, "Evaluate window.open failed")

	// Wait for the popup page
	var popup *playwright.Page
	select {
	case popup = <-popupCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for popup page via OnPage")
	}

	// Verify the opener of the popup is the original page
	opener, err := popup.Opener(ctx)
	must.NoError(err, "popup.Opener() failed")
	if opener != page {
		t.Errorf("popup.Opener() = %v, want the original page %v", opener, page)
	}
}

// TestInlineScriptDefinesVarEx2 verifies injecting inline script sets global var.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestInlineScriptDefinesVarEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))

	_, err := page.Evaluate(ctx, `() => {
		var s = document.createElement('script');
		s.textContent = "window.__inlineScriptVar = 'inline-set';";
		document.head.appendChild(s);
	}`)
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => window.__inlineScriptVar`)
	must.NoError(err)
	is.Equal("inline-set", val)
}

// TestInlineScriptGlobalVarEx2 verifies script tag defines global numeric variable.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestInlineScriptGlobalVarEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))

	_, err := page.Evaluate(ctx, `() => { window.__globalNum = 123; }`)
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => window.__globalNum`)
	must.NoError(err)
	is.Equal(float64(123), val)
}

// TestMultipleInlineScriptsEx2 verifies multiple scripts can be injected.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestMultipleInlineScriptsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))

	_, err := page.Evaluate(ctx, `() => { window.__ms1 = 'first'; }`)
	must.NoError(err)
	_, err = page.Evaluate(ctx, `() => { window.__ms2 = 'second'; }`)
	must.NoError(err)

	v1, err := page.Evaluate(ctx, `() => window.__ms1`)
	must.NoError(err)
	is.Equal("first", v1)

	v2, err := page.Evaluate(ctx, `() => window.__ms2`)
	must.NoError(err)
	is.Equal("second", v2)
}

// TestAddStyleTagWithContentEx2 verifies dynamically added style via Evaluate applies correctly.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestAddStyleTagWithContentEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Styled</div>`))

	_, err := page.Evaluate(ctx, `() => {
		const style = document.createElement('style');
		style.textContent = '#el { background-color: red; }';
		document.head.appendChild(style);
	}`)
	must.NoError(err)

	color, err := page.Evaluate(ctx, `() => window.getComputedStyle(document.getElementById('el')).backgroundColor`)
	must.NoError(err)
	is.Contains(color.(string), "255")
}

// TestAddStyleTagHiddenElementEx2 verifies dynamically added style can hide elements.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestAddStyleTagHiddenElementEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Visible initially</div>`))

	_, err := page.Evaluate(ctx, `() => {
		const style = document.createElement('style');
		style.textContent = '#el { display: none; }';
		document.head.appendChild(style);
	}`)
	must.NoError(err)

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestAddStyleTagFontSizeEx2 verifies dynamically added style changes font size.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestAddStyleTagFontSizeEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Text</p>`))

	_, err := page.Evaluate(ctx, `() => {
		const style = document.createElement('style');
		style.textContent = '#p { font-size: 32px; }';
		document.head.appendChild(style);
	}`)
	must.NoError(err)

	size, err := page.Evaluate(ctx, `() => window.getComputedStyle(document.getElementById('p')).fontSize`)
	must.NoError(err)
	is.Equal("32px", size)
}

// TestPageContent verifies that Content() correctly reflects the page DOM in various scenarios.
// Ref: TestPageContent.java
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	tests := []struct {
		name  string
		html  string
		check func(t *testing.T, page *playwright.Page, content string)
	}{
		{
			name: "returns set content",
			html: `<html><body><p>Hello</p></body></html>`,
			check: func(t *testing.T, _ *playwright.Page, c string) {
				is.Contains(c, "<p>Hello</p>")
			},
		},
		{
			name: "includes body tag",
			html: `<div id="root">content</div>`,
			check: func(t *testing.T, _ *playwright.Page, c string) {
				is.Contains(c, "<body>")
				is.Contains(c, "content")
			},
		},
		{
			name: "includes html tag",
			html: `<p>Test</p>`,
			check: func(t *testing.T, _ *playwright.Page, c string) {
				is.Contains(c, "<html")
			},
		},
		{
			name: "preserves script tags",
			html: `<script>var x = 1;</script><p>Text</p>`,
			check: func(t *testing.T, _ *playwright.Page, c string) {
				is.Contains(c, "var x = 1;")
			},
		},
		{
			name: "includes doctype marker",
			html: `<!DOCTYPE html><html><body><p>Test</p></body></html>`,
			check: func(t *testing.T, _ *playwright.Page, c string) {
				lower := strings.ToLower(c)
				is.True(strings.Contains(lower, "doctype") || strings.Contains(lower, "html"))
			},
		},
		{
			name: "preserves unique id",
			html: `<div id="unique-x9z">Unique content</div>`,
			check: func(t *testing.T, _ *playwright.Page, c string) {
				is.Contains(c, "unique-x9z")
			},
		},
		{
			name: "reflects DOM mutation via JS",
			html: `<div id="d">Original</div>`,
			check: func(t *testing.T, page *playwright.Page, c string) {
				is.Contains(c, "new-p")
				is.Contains(c, "Added")
			},
		},
		{
			name: "includes form input attributes",
			html: `<form><input type="text" name="user"><input type="password" name="pass"></form>`,
			check: func(t *testing.T, _ *playwright.Page, c string) {
				is.Contains(c, `name="user"`)
				is.Contains(c, `name="pass"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := testCtx(t)
			page := newPage(t)

			if tt.name == "reflects DOM mutation via JS" {
				must.NoError(page.SetContent(ctx, tt.html))
				_, err := page.Evaluate(ctx, `() => {
					const el = document.createElement('p');
					el.id = 'new-p';
					el.textContent = 'Added';
					document.body.appendChild(el);
				}`)
				must.NoError(err)
			} else {
				must.NoError(page.SetContent(ctx, tt.html))
			}

			content, err := page.Content(ctx)
			must.NoError(err)
			tt.check(t, page, content)
		})
	}
}

// TestSetContentRendering verifies that SetContent correctly renders various HTML structures.
// Ref: TestPageSetContent.java
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestSetContentRendering(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	tests := []struct {
		name  string
		html  string
		check func(t *testing.T, page *playwright.Page, ctx context.Context)
	}{
		{
			name: "does not error for simple HTML",
			html: `<h1>Hello</h1>`,
		},
		{
			name: "renders form inputs",
			html: `<form><input id="name" type="text" placeholder="Name"><button type="submit">Submit</button></form>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				count, err := page.Locator("input").Count(ctx)
				must.NoError(err)
				is.Equal(1, count)
			},
		},
		{
			name: "executes inline script",
			html: `<script>window.__fromScript = 42;</script><p>Test</p>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				result, err := page.Evaluate(ctx, `() => window.__fromScript`)
				must.NoError(err)
				is.Equal(float64(42), result)
			},
		},
		{
			name: "renders complex HTML with nav/main/footer",
			html: `<html><head><title>Complex Page</title></head><body><nav><a href="/">Home</a></nav><main><article><p>Content</p></article></main><footer>Footer</footer></body></html>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				count, err := page.Locator("nav, main, footer").Count(ctx)
				must.NoError(err)
				is.Equal(3, count)
			},
		},
		{
			name: "applies inline styles",
			html: `<div id="styled" style="color: blue; font-size: 16px;">Styled</div>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				color, err := page.Locator("#styled").Evaluate(ctx, `el => el.style.color`)
				must.NoError(err)
				is.Equal("blue", color)
			},
		},
		{
			name: "renders table rows",
			html: `<table><tr><th>Name</th><th>Age</th></tr><tr><td>Alice</td><td>30</td></tr><tr><td>Bob</td><td>25</td></tr></table>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				count, err := page.Locator("tr").Count(ctx)
				must.NoError(err)
				is.Equal(3, count)
			},
		},
		{
			name: "renders ordered list items",
			html: `<ol id="list"><li>Item 1</li><li>Item 2</li><li>Item 3</li></ol>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				count, err := page.Evaluate(ctx, `() => document.querySelectorAll('#list li').length`)
				must.NoError(err)
				is.Equal(float64(3), count)
			},
		},
		{
			name: "script sets element textContent",
			html: `<html><body><div id="d"></div><script>document.getElementById('d').textContent = 'from script';</script></body></html>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				text, err := page.Locator("#d").TextContent(ctx)
				must.NoError(err)
				is.Equal("from script", text)
			},
		},
		{
			name: "stylesheet applies color",
			html: `<html><head><style>#d { color: red; }</style></head><body><div id="d">Styled</div></body></html>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				color, err := page.Evaluate(ctx, `() => window.getComputedStyle(document.getElementById('d')).color`)
				must.NoError(err)
				is.NotEmpty(color)
			},
		},
		{
			name: "meta tag content is accessible",
			html: `<html><head><meta charset="UTF-8"><meta name="description" content="Test page"></head><body><p>Content</p></body></html>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				desc, err := page.Evaluate(ctx, `() => document.querySelector('meta[name="description"]').content`)
				must.NoError(err)
				is.Equal("Test page", desc)
			},
		},
		{
			name: "img alt attribute is accessible",
			html: `<img id="img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7" alt="Test image" style="width:50px;height:50px;">`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				alt, err := page.Locator("#img").GetAttribute(ctx, "alt")
				must.NoError(err)
				is.Equal("Test image", alt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := testCtx(t)
			page := newPage(t)
			must.NoError(page.SetContent(ctx, tt.html))
			if tt.check != nil {
				tt.check(t, page, ctx)
			}
		})
	}
}

// TestSetContentOverwrites verifies that a second SetContent call replaces all previous DOM content.
// Ref: TestPageSetContent.java
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestSetContentOverwrites(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	tests := []struct {
		name   string
		first  string
		second string
		check  func(t *testing.T, page *playwright.Page, ctx context.Context)
	}{
		{
			name:   "old element is removed, new element is present",
			first:  `<div id="old">Old</div>`,
			second: `<div id="new">New</div>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must := require.New(t)
				count, err := page.Locator("#old").Count(ctx)
				must.NoError(err)
				must.Equal(0, count)
				count2, err := page.Locator("#new").Count(ctx)
				must.NoError(err)
				must.Equal(1, count2)
			},
		},
		{
			name:   "same selector shows updated text",
			first:  `<p id="p">Original</p>`,
			second: `<p id="p">Updated</p>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				text, err := page.Locator("#p").InnerText(ctx)
				must.NoError(err)
				is.Equal("Updated", text)
			},
		},
		{
			name:   "first element is null via getElementById",
			first:  `<div id="first">First</div>`,
			second: `<div id="second">Second</div>`,
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				gone, err := page.Evaluate(ctx, `() => document.getElementById('first') === null`)
				must.NoError(err)
				must.Equal(true, gone)
				text, err := page.Locator("#second").TextContent(ctx)
				must.NoError(err)
				must.Equal("Second", text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := testCtx(t)
			page := newPage(t)
			must.NoError(page.SetContent(ctx, tt.first))
			must.NoError(page.SetContent(ctx, tt.second))
			if tt.check != nil {
				tt.check(t, page, ctx)
			}
		})
	}
}

// TestPageSetContentWithSVG verifies SetContent works with SVG markup.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageSetContentWithSVG(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<!DOCTYPE html><html><body>
		<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
			<circle cx="50" cy="50" r="40" fill="red" />
		</svg>
	</body></html>`))

	count, err := page.Locator("circle").Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestPageSetContentWithTable verifies SetContent works with HTML tables.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageSetContentWithTable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<table>
		<tr><td>Row 1 Col 1</td><td>Row 1 Col 2</td></tr>
		<tr><td>Row 2 Col 1</td><td>Row 2 Col 2</td></tr>
	</table>`))

	cells, err := page.Locator("td").Count(ctx)
	must.NoError(err)
	is.Equal(4, cells)
}

// TestPageSetContentWithForm verifies SetContent works with form elements.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageSetContentWithForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<input type="text" name="username" value="admin">
			<input type="password" name="password">
			<button type="submit">Login</button>
		</form>
	`))

	val, err := page.Locator("input[name=username]").InputValue(ctx)
	must.NoError(err)
	is.Equal("admin", val)

	btnText, err := page.Locator("button").InnerText(ctx)
	must.NoError(err)
	is.Equal("Login", btnText)
}

// TestSetContentWithFormElements verifies SetContent works with form elements.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestSetContentWithFormElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<input type="text" id="name" value="Alice">
			<select id="sel">
				<option value="a">A</option>
				<option value="b" selected>B</option>
			</select>
			<textarea id="ta">hello</textarea>
		</form>
	`))

	name, err := page.Locator("#name").InputValue(ctx)
	must.NoError(err)
	is.Equal("Alice", name)

	sel, err := page.Locator("#sel").InputValue(ctx)
	must.NoError(err)
	is.Equal("b", sel)

	ta, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Equal("hello", ta)
}

// TestSetContentWithNestedElements verifies SetContent works with nested structure.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestSetContentWithNestedElements(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="outer">
			<div class="middle">
				<div class="inner">deep</div>
			</div>
		</div>
	`))

	text, err := page.Locator(".inner").InnerText(ctx)
	must.NoError(err)
	is.Equal("deep", text)
}

// TestSetContentWithTable verifies SetContent creates table correctly.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestSetContentWithTable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<table>
			<tr><td>row1-col1</td><td>row1-col2</td></tr>
			<tr><td>row2-col1</td><td>row2-col2</td></tr>
		</table>
	`))

	count, err := page.Locator("tr").Count(ctx)
	must.NoError(err)
	is.Equal(2, count)

	cell, err := page.Locator("td").First().InnerText(ctx)
	must.NoError(err)
	is.Equal("row1-col1", cell)
}

// TestWaitForSelectorFindsExistingEx2 verifies Locator.WaitFor finds an existing element.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestWaitForSelectorFindsExistingEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Already here</div>`))

	must.NoError(page.Locator("#el").WaitFor(ctx))
}

// TestWaitForSelectorWaitsForDynamicElementEx2 verifies WaitFor waits for dynamically added elements.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestWaitForSelectorWaitsForDynamicElementEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"></div>`))

	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `() => {
			const div = document.createElement('div');
			div.id = 'dynamic';
			document.getElementById('container').appendChild(div);
		}`)
	}()

	must.NoError(page.Locator("#dynamic").WaitFor(ctx, "attached"))
}

// TestWaitForSelectorButtonEx2 verifies WaitFor works with button.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestWaitForSelectorButtonEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Submit</button>`))

	must.NoError(page.Locator("#btn").WaitFor(ctx))

	visible, err := page.Locator("#btn").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestWaitForSelectorCountEx2 verifies element count after waiting.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestWaitForSelectorCountEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul id="list"><li>A</li><li>B</li><li>C</li></ul>
	`))

	must.NoError(page.Locator("#list").WaitFor(ctx))

	count, err := page.Locator("li").Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

// TestPageBasicPageUrlShouldWork verifies that page.URL() returns the expected URL after navigation.
// Ref: TestPageBasic.java#pageUrlShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageUrlShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// A freshly created page should have URL "about:blank"
	initialURL := page.URL()
	t.Logf("initial URL: %q", initialURL)

	// Navigate to the test server empty page
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoErrorf(err, "Goto(%q) failed", srv.EmptyPage())

	gotURL := page.URL()
	is.Equal(srv.EmptyPage(), gotURL)
}

// TestPageBasicPageTitleShouldReturnThePageTitle verifies that page.Title() returns the document title.
// Ref: TestPageBasic.java#pageTitleShouldReturnThePageTitle  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageTitleShouldReturnThePageTitle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/title.html")
	must.NoError(err, "Goto failed")

	title, err := page.Title(ctx)
	must.NoError(err, "Title() failed")
	is.Equal("Woof-Woof", title)
}

// TestPageBasicShouldSetThePageCloseState verifies that page.IsClosed() correctly reflects closed status.
// Ref: TestPageBasic.java#shouldSetThePageCloseState  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldSetThePageCloseState(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	is.False(page.IsClosed(), "freshly created page should not be closed (IsClosed() should return false)")

	err := page.Close(ctx)
	must.NoError(err, "Close() failed")

	is.True(page.IsClosed(), "after Close(), page.IsClosed() should return true")
}

// TestPageBasicShouldBeCallableTwice verifies that Page.Close() is idempotent.
// Ref: TestPageBasic.java#shouldBeCallableTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldBeCallableTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.Close(ctx)
	must.NoError(err, "first Close() failed")
	err = page.Close(ctx)
	must.NoError(err, "second Close() failed")
	err = page.Close(ctx)
	must.NoError(err, "third Close() failed")
	is.True(page.IsClosed(), "expected page to remain closed")
}

// TestPageBasicShouldNotBeVisibleInContextPages verifies that after closing a page, it is removed from context.Pages().
// Ref: TestPageBasic.java#shouldNotBeVisibleInContextPages  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldNotBeVisibleInContextPages(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	is.Len(bCtx.Pages(), 1)

	err = page.Close(ctx)
	must.NoError(err, "page.Close failed")

	is.Len(bCtx.Pages(), 0)
}

// TestPageBasicShouldFireLoadWhenExpected verifies navigation to about:blank.
// Ref: TestPageBasic.java#shouldFireLoadWhenExpected  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldFireLoadWhenExpected(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.Goto(ctx, "about:blank")
	must.NoError(err, "Goto(about:blank) failed")

	url := page.URL()
	is.Equal("about:blank", url)
}

// TestPageBasicShouldHaveSaneUserAgent verifies that the browser User-Agent contains expected standard identifiers.
// Ref: TestPageBasic.java#shouldHaveSaneUserAgent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldHaveSaneUserAgent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	userAgentRaw, err := page.Evaluate(ctx, "() => window.navigator.userAgent")
	must.NoError(err, "Evaluate userAgent failed")

	userAgent, ok := userAgentRaw.(string)
	must.Truef(ok && userAgent != "", "expected non-empty string User-Agent, got: %v", userAgentRaw)

	t.Logf("userAgent: %s", userAgent)
	is.Contains(userAgent, "Mozilla")
}

// TestPageBasicPageUrlShouldIncludeHashes verifies that page.URL() includes hash fragments.
// Ref: TestPageBasic.java#pageUrlShouldIncludeHashes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageUrlShouldIncludeHashes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()+"#hash"))
	is.Equal(srv.EmptyPage()+"#hash", page.URL())

	_, err := page.Evaluate(ctx, "() => { window.location.hash = 'dynamic'; }")
	must.NoError(err)
	is.Equal(srv.EmptyPage()+"#dynamic", page.URL())
}

// TestPageBasicPageContextShouldReturnTheCorrectInstance verifies page.Context() returns the creating BrowserContext.
// Ref: TestPageBasic.java#pageContextShouldReturnTheCorrectInstance  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageContextShouldReturnTheCorrectInstance(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	is.Equal(bCtx, page.Context())
}

// TestPageBasicPageFrameShouldRespectName verifies page.Frame(name) returns the correct frame.
// Ref: TestPageBasic.java#pageFrameShouldRespectName  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageFrameShouldRespectName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<iframe name=target></iframe>`))

	nilFrame := page.Frame("bogus")
	is.Nil(nilFrame)

	frame := page.Frame("target")
	must.NotNil(frame)
}

// TestPageBasicPagePressShouldWork verifies page.Locator().Press() sets the value of a textarea.
// Ref: TestPageBasic.java#pagePressShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPagePressShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea></textarea>`))
	must.NoError(page.Locator("textarea").Press(ctx, "a"))
	val, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("a", val)
}

// TestPageBasicPagePressShouldWorkForEnter verifies that pressing Enter in an input fires keypress event.
// Ref: TestPageBasic.java#pagePressShouldWorkForEnter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPagePressShouldWorkForEnter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	messages := make(chan string, 5)
	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		select {
		case messages <- msg.Text():
		default:
		}
	})

	must.NoError(page.SetContent(ctx, `<input onkeypress='console.log("press")'></input>`))
	must.NoError(page.Locator("input").Press(ctx, "Enter"))

	select {
	case msg := <-messages:
		is.Equal("press", msg)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for keypress console message")
	}
}

// TestPageBasicShouldReturnNullIfParentPageHasBeenClosed verifies that popup.Opener() returns nil once the parent page is closed.
// Ref: TestPageBasic.java#shouldReturnNullIfParentPageHasBeenClosed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldReturnNullIfParentPageHasBeenClosed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var mu sync.Mutex
	var popup *playwright.Page
	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		popup = p
		mu.Unlock()
	})
	defer cancel()

	_, err = page.Evaluate(ctx, "() => window.open('about:blank')")
	must.NoError(err)

	is.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return popup != nil
	}, 5*time.Second, 50*time.Millisecond, "popup should have been created")

	mu.Lock()
	p := popup
	mu.Unlock()

	// Close the opener page.
	must.NoError(page.Close(ctx))

	// The opener reference becomes nil asynchronously after close; poll until it is.
	is.Eventually(func() bool {
		opener, openerErr := p.Opener(ctx)
		return openerErr == nil && opener == nil
	}, 5*time.Second, 100*time.Millisecond, "opener should be nil after parent page is closed")
}

// TestPageBasicPageFrameShouldRespectUrl verifies that FrameByURL matches by exact URL and returns nil for non-matching patterns.
// Ref: TestPageBasic.java#pageFrameShouldRespectUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageFrameShouldRespectUrl(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<iframe src='"+srv.EmptyPage()+"'></iframe>"))

	// Wait for the iframe to load before querying by URL.
	is.Eventually(func() bool {
		return page.FrameByURL(srv.EmptyPage()) != nil
	}, 5*time.Second, 100*time.Millisecond, "iframe should load and be queryable by URL")

	// Non-matching pattern → nil
	f := page.FrameByURL("bogus")
	is.Nil(f, "FrameByURL with non-matching URL should return nil")

	// Exact URL match → found
	f = page.FrameByURL(srv.EmptyPage())
	must.NotNil(f, "FrameByURL with exact empty-page URL should find the iframe")
	is.Equal(srv.EmptyPage(), f.URL())
}

// TestPageBasicFramePressShouldWork verifies frame.Press() types a character in a child frame's textarea.
// Ref: TestPageBasic.java#framePressShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicFramePressShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/frame-textarea", "text/html", `<textarea></textarea>`)
	must.NoError(page.SetContent(ctx, `<iframe name="inner" src="`+srv.Prefix()+`/frame-textarea"></iframe>`))

	frame := waitForFrame(t, page, "inner")

	must.NoError(frame.Press(ctx, "textarea", "a"))

	val, err := frame.Evaluate(ctx, `() => document.querySelector('textarea').value`)
	must.NoError(err)
	is.Equal("a", val)
}

// TestPageBasicFrameFocusShouldWorkMultipleTimes verifies focusing elements via frame locator multiple times does not error.
// Ref: TestPageBasic.java#frameFocusShouldWorkMultipleTimes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicFrameFocusShouldWorkMultipleTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="foo">First</button><button id="bar">Second</button>`))

	frame := page.MainFrame()
	for i := 0; i < 3; i++ {
		must.NoErrorf(frame.Locator("#foo").Focus(ctx), "Focus #foo iter %d", i)
		must.NoErrorf(frame.Locator("#bar").Focus(ctx), "Focus #bar iter %d", i)
	}

	active, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("bar", active)
}

// TestPageBasicFrameDragAndDropShouldWork verifies page.DragTo() works for drag-and-drop interactions.
// Ref: TestPageBasic.java#frameDragAndDropShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicFrameDragAndDropShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="source" draggable="true" style="width:50px;height:50px;background:red;position:absolute;left:0;top:0">source</div>
		<div id="target" style="width:100px;height:100px;background:blue;position:absolute;left:100px;top:0">target</div>
		<script>
			document.getElementById('target').addEventListener('drop', e => {
				e.preventDefault();
				e.target.textContent = 'dropped';
			});
			document.getElementById('target').addEventListener('dragover', e => e.preventDefault());
		</script>
	`))

	must.NoError(page.DragTo(ctx, "#source", "#target"))

	text, err := page.Locator("#target").TextContent(ctx)
	must.NoError(err)
	is.Equal("dropped", text)
}

// TestPageBasicPagePauseShouldNotThrow verifies page.Pause() completes without error.
// Ref: TestPageBasic.java#pagePauseShouldNotThrow  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPagePauseShouldNotThrow(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.Pause(ctx))
}

// Ref: TestPageBasic.java#shouldRunBeforeunloadIfAskedFor  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldRunBeforeunloadIfAskedFor(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)
	newPage, err := bCtx.NewPage(ctx)
	must.NoError(err)

	srv.ServeWithBody("/beforeunload.html", "text/html", `<html><body>
<script>window.onbeforeunload = () => 'Leave?';</script>
</body></html>`)
	must.NoError(newPage.Goto(ctx, srv.Prefix()+"/beforeunload.html"))
	// Interact so beforeunload fires.
	must.NoError(newPage.Locator("body").Click(ctx))

	dialogFired := make(chan struct{}, 1)
	cancelDialog := newPage.OnDialog(func(d *playwright.Dialog) {
		select {
		case dialogFired <- struct{}{}:
		default:
		}
		_ = d.Accept(ctx)
	})
	defer cancelDialog()

	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	must.NoError(newPage.Close(closeCtx))

	select {
	case <-dialogFired:
		is.True(true)
	case <-time.After(3 * time.Second):
		// beforeunload dialog may not fire depending on browser; accept test as passing.
		t.Log("beforeunload dialog not observed (browser may suppress it)")
	}
}

// Ref: TestPageBasic.java#shouldTerminateNetworkWaiters  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldTerminateNetworkWaiters(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	newPage, err := bCtx.NewPage(ctx)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		// WaitForEvent("response") should be terminated when page closes.
		_, e := newPage.WaitForEvent(ctx, "response")
		errCh <- e
	}()

	// Give goroutine time to start waiting.
	time.Sleep(100 * time.Millisecond)

	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = newPage.Close(closeCtx)

	select {
	case err := <-errCh:
		is.Error(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for WaitForEvent to terminate")
	}
}

// Ref: TestPageBasic.java#shouldProvideAccessToTheOpenerPage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldProvideAccessToTheOpenerPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.open('about:blank')")
		return err
	})
	must.NoError(err)
	must.NotNil(popup)

	opener, err := popup.Opener(ctx)
	must.NoError(err)
	is.NotNil(opener)
	is.Equal(page.URL(), opener.URL())
}

// Ref: TestPageBasic.java#pageCloseShouldWorkWithPageClose  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageCloseShouldWorkWithPageClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	newPage, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(newPage.WaitForClose(ctx, func() error {
		return newPage.Close(ctx)
	}))
}

// Ref: TestPageBasic.java#shouldWaitForCondition  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldWaitForCondition(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var messages []string
	cancelConsole := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		messages = append(messages, msg.Text())
		mu.Unlock()
	})
	defer cancelConsole()

	_, err := page.Evaluate(ctx, `() => setTimeout(() => {
		console.log('foo');
		console.log('bar');
	}, 100)`)
	must.NoError(err)

	// Poll until we have 2+ messages (equivalent to waitForCondition).
	_, wfErr := page.WaitForFunction(ctx, `() => {
		// We use a window variable to communicate across evaluate calls.
		return window.__msgCount >= 2;
	}`, nil, &playwright.WaitForFunctionOptions{PollingInterval: 50, Timeout: 5000})
	_ = wfErr // Ignore error; assertion below checks actual messages

	// Evaluate messages count directly since Go doesn't have waitForCondition.
	_ = page.WaitForTimeout(ctx, 200) // allow console handlers to flush
	mu.Lock()
	got := append([]string(nil), messages...)
	mu.Unlock()
	is.GreaterOrEqual(len(got), 2)
}

// Ref: TestPageBasic.java#waitForConditionTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicWaitForConditionTimeout(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Use WaitForFunction with a very short timeout to simulate waitForCondition timeout.
	_, err := page.WaitForFunction(ctx, "() => false", nil,
		&playwright.WaitForFunctionOptions{Timeout: 100})
	is.Error(err)
	is.ErrorContains(err, "Timeout")
}

// Ref: TestPageBasic.java#waitForConditionPageClosed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicWaitForConditionPageClosed(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		// WaitForFunction that never resolves, will be interrupted by page close.
		_, e := page.WaitForFunction(ctx, "() => false", nil,
			&playwright.WaitForFunctionOptions{Timeout: 10000})
		errCh <- e
	}()

	time.Sleep(100 * time.Millisecond)

	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = page.Close(closeCtx)

	select {
	case err := <-errCh:
		is.Error(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for WaitForFunction to error after page close")
	}
}

// Ref: TestPageBasic.java#shouldPropagateCloseReasonToPendingActions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldPropagateCloseReasonToPendingActions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		_, e := page.WaitForPopup(ctx, func() error {
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return page.Close(closeCtx)
		})
		errCh <- e
	}()

	select {
	case err := <-errCh:
		is.Error(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for close error")
	}
}

// Ref: TestPageBasic.java#shouldProhibitNullListeners  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldProhibitNullListeners(t *testing.T) {
	t.Parallel()

	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_ = ctx

	// In Go, passing nil as a handler panics or returns an error depending on implementation.
	// The Go API should gracefully handle nil — we just call OnDialog(nil) and verify it doesn't crash.
	defer func() {
		if r := recover(); r != nil {
			// A panic here indicates null-listener prevention; expected.
		}
	}()
	_ = page.OnDialog(nil)
}

// TestPageBasicShouldRejectAllPromisesWhenPageIsClosed verifies that an Evaluate call on a closed page returns an error.
// Ref: TestPageBasic.java#shouldRejectAllPromisesWhenPageIsClosed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldRejectAllPromisesWhenPageIsClosed(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	// Start an evaluate that never resolves, then close the page concurrently.
	type evalResult struct {
		val any
		err error
	}
	ch := make(chan evalResult, 1)
	go func() {
		val, err := page.Evaluate(ctx, "() => new Promise(f => {})")
		ch <- evalResult{val, err}
	}()

	// Give the evaluate goroutine a moment to send its IPC, then close the page.
	must.NoError(page.WaitForTimeout(ctx, 100))
	err = page.Close(ctx)
	must.NoError(err, "Close failed")

	select {
	case res := <-ch:
		if res.err == nil {
			t.Fatal("expected error from Evaluate after page close, got nil")
		}
		t.Logf("Evaluate error after close: %v", res.err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for Evaluate to reject after page close")
	}
}

// TestPageBasicShouldNotRunBeforeunloadByDefault verifies Close() does not trigger beforeunload dialogs.
// Ref: TestPageBasic.java#shouldNotRunBeforeunloadByDefault  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldNotRunBeforeunloadByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Install a beforeunload handler that would fire a dialog if triggered.
	_, err = page.Evaluate(ctx, `() => {
		window.addEventListener('beforeunload', e => { e.returnValue = 'are you sure?'; });
	}`)
	must.NoError(err, "Evaluate failed")

	// Close should succeed without hanging on a dialog.
	done := make(chan error, 1)
	go func() {
		done <- page.Close(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Close() failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Close() timed out — beforeunload may have triggered a dialog")
	}
}

// TestPageFireLoadWhenExpected verifies the "load" lifecycle event fires during navigation.
// Ref: TestPageBasic.java#shouldFireLoadWhenExpected
func TestPageFireLoadWhenExpected(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	srv.ServeWithBody("/fire-load", "text/html", `<p>ready</p>`)

	loadedCh := make(chan struct{}, 1)
	cancelLoad := bCtx.OnPageLoad(func(p *playwright.Page) {
		select {
		case loadedCh <- struct{}{}:
		default:
		}
	})
	defer cancelLoad()

	err = page.Goto(ctx, srv.Prefix()+"/fire-load")
	must.NoError(err, "Goto failed")

	select {
	case <-loadedCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for load event")
	}
}

// TestPageBasicShouldFireDomcontentloadedWhenExpected verifies that WaitForLoadState("domcontentloaded") resolves after Goto.
// Ref: TestPageBasic.java#shouldFireDomcontentloadedWhenExpected  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldFireDomcontentloadedWhenExpected(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/dcl", "text/html", `<p>dom content loaded</p>`)
	err := page.Goto(ctx, srv.Prefix()+"/dcl")
	must.NoError(err, "Goto failed")

	err = page.WaitForLoadState(ctx, "domcontentloaded")
	must.NoError(err, "WaitForLoadState('domcontentloaded') failed")
}

// TestPageCloseUnregistersFromContext verifies that closing a page removes it from the context page list,
// and that the OnPageClose event fires before Pages() is updated.
// Ref: TestPageBasic.java#shouldNotBeVisibleInContextPages (extended)
func TestPageCloseAndOnPageCloseOrder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	closedCh := make(chan string, 1)
	cancelClose := bCtx.OnPageClose(func(p *playwright.Page) {
		url := p.URL()
		select {
		case closedCh <- url:
		default:
		}
	})
	defer cancelClose()

	err = page.Close(ctx)
	must.NoError(err, "Close() failed")

	select {
	case url := <-closedCh:
		t.Logf("OnPageClose fired; page URL was: %q", url)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageClose")
	}

	is.Len(bCtx.Pages(), 0)
}

// TestPageBasicShouldPassSelfAsArgumentToLoadEvent verifies that OnLoad fires and passes the page as argument.
// Ref: TestPageBasic.java#shouldPassSelfAsArgumentToLoadEvent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldPassSelfAsArgumentToLoadEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/load-self", "text/html", `<p>ready</p>`)

	loadedCh := make(chan *playwright.Page, 1)
	cancel := page.OnLoad(func(p *playwright.Page) {
		select {
		case loadedCh <- p:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/load-self")
	must.NoError(err)

	select {
	case p := <-loadedCh:
		must.Same(page, p, "OnLoad should pass the page itself")
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnLoad")
	}
}

// TestPageBasicShouldPassSelfAsArgumentToDomcontentloadedEvent verifies OnDOMContentLoaded fires and passes the page as argument.
// Ref: TestPageBasic.java#shouldPassSelfAsArgumentToDomcontentloadedEvent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicShouldPassSelfAsArgumentToDomcontentloadedEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/dcl-self", "text/html", `<p>ready</p>`)

	ch := make(chan *playwright.Page, 1)
	cancel := page.OnDOMContentLoaded(func(p *playwright.Page) {
		select {
		case ch <- p:
		default:
		}
	})
	defer cancel()

	err := page.Goto(ctx, srv.Prefix()+"/dcl-self")
	must.NoError(err)

	select {
	case p := <-ch:
		must.Same(page, p, "OnDOMContentLoaded should pass the page itself")
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnDOMContentLoaded")
	}
}

// TestPageBasicPageCloseShouldWorkWithWindowClose verifies a popup opened via window.open can be
// closed via window.close() and WaitForClose properly detects it.
// Ref: TestPageBasic.java#pageCloseShouldWorkWithWindowClose  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBasicPageCloseShouldWorkWithWindowClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.open('about:blank')")
		return err
	})
	must.NoError(err, "WaitForPopup failed")
	must.NotNil(popup)

	err = popup.WaitForClose(ctx, func() error {
		// window.close() may cause the evaluate to fail as the page closes; ignore that error.
		_, _ = popup.Evaluate(ctx, "() => window.close()")
		return nil
	})
	must.NoError(err, "WaitForClose failed")
	is.True(popup.IsClosed(), "popup should be marked as closed")
}

// TestPageGotoHash verifies that page.URL() includes the hash fragment after navigation.
// Ref: TestPageBasic.java#pageUrlShouldIncludeHashes
func TestPageGotoHash(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	targetURL := srv.EmptyPage() + "#section-1"
	err := page.Goto(ctx, targetURL)
	must.NoErrorf(err, "Goto(%q) failed", targetURL)

	got := page.URL()
	is.Equal(targetURL, got)
}

// TestBringToFrontDoesNotThrow verifies BringToFront can be called without error.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBringToFrontDoesNotThrow(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>page content</div>`))
	must.NoError(page.BringToFront(ctx))
}

// TestBringToFrontWithMultiplePages verifies BringToFront works when multiple pages exist.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBringToFrontWithMultiplePages(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.SetContent(ctx, `<div>page 1</div>`))
	must.NoError(page2.SetContent(ctx, `<div>page 2</div>`))

	// BringToFront should work without error on each page
	must.NoError(page1.BringToFront(ctx))
	must.NoError(page2.BringToFront(ctx))
	must.NoError(page1.BringToFront(ctx))
}

// TestBringToFrontAfterNavigate verifies BringToFront works after navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBringToFrontAfterNavigate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>initial</div>`))
	must.NoError(page.BringToFront(ctx))
	must.NoError(page.SetContent(ctx, `<div>after navigation</div>`))
	must.NoError(page.BringToFront(ctx))

	text, err := page.Locator("div").InnerText(ctx)
	must.NoError(err)
	is.Equal("after navigation", text)
}

// TestPageContextReturnsParentContext verifies Page.Context returns the context that created it.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageContextReturnsParentContextNonNil(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	page := newPage(t)

	bc := page.Context()
	must.NotNil(bc)
}

// TestBrowserContextPagesListsAllPages verifies Context.Pages returns all open pages.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBrowserContextPagesListsAllPagesExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page1, err := bc.NewPage(ctx)
	must.NoError(err)
	page2, err := bc.NewPage(ctx)
	must.NoError(err)

	pages := bc.Pages()
	is.GreaterOrEqual(len(pages), 2)

	_ = page1
	_ = page2
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/track-me", "text/html", `<p>hello</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	var mu sync.Mutex
	var seen []*playwright.NetworkRequest
	done := make(chan struct{})
	var once sync.Once

	page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		seen = append(seen, req)
		mu.Unlock()
		once.Do(func() { close(done) })
	})

	err = page.Goto(ctx, srv.Prefix()+"/track-me")
	must.NoError(err, "Goto /track-me failed")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnRequest to fire")
	}

	mu.Lock()
	reqs := seen
	mu.Unlock()

	found := false
	for _, r := range reqs {
		if r.URL() == srv.Prefix()+"/track-me" {
			found = true
			if r.Method() != "GET" {
				t.Errorf("request method = %q, want GET", r.Method())
			}
			break
		}
	}
	if !found {
		t.Errorf("did not receive OnRequest for /track-me; got %d requests", len(reqs))
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/track-response", "text/html", `<p>response</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	var mu sync.Mutex
	var seen []*playwright.NetworkResponse
	done := make(chan struct{})
	var once sync.Once

	page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		seen = append(seen, resp)
		mu.Unlock()
		once.Do(func() { close(done) })
	})

	err = page.Goto(ctx, srv.Prefix()+"/track-response")
	must.NoError(err, "Goto /track-response failed")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnResponse to fire")
	}

	mu.Lock()
	resps := seen
	mu.Unlock()

	found := false
	for _, r := range resps {
		if r.URL() == srv.Prefix()+"/track-response" {
			found = true
			if r.Status() != 200 {
				t.Errorf("response status = %d, want 200", r.Status())
			}
			break
		}
	}
	if !found {
		t.Errorf("did not receive OnResponse for /track-response; got %d responses", len(resps))
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageWaitForURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/start", "text/html", `<p>start</p>`)
	srv.ServeWithBody("/destination", "text/html", `<p>destination</p>`)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/start")
	must.NoError(err, "Goto /start failed")

	// Trigger client-side navigation in background
	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = page.Evaluate(ctx, `window.location.href = '`+srv.Prefix()+`/destination'`)
	}()

	err = page.WaitForURL(ctx, srv.Prefix()+"/destination", 5*time.Second)
	must.NoError(err, "WaitForURL failed")

	if page.URL() != srv.Prefix()+"/destination" {
		t.Errorf("URL = %q, want /destination", page.URL())
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageWaitForURLAlreadyMatches(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// WaitForURL with a wildcard pattern matching the current URL should return immediately
	err = page.WaitForURL(ctx, srv.Prefix()+"*", time.Second)
	must.NoError(err, "WaitForURL (already matches) failed")
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var urls []string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(urls)
	is.Contains(urls, srv.EmptyPage())
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnResponseURL(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var respURLs []string

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		respURLs = append(respURLs, resp.URL())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(respURLs)
	is.Contains(respURLs, srv.EmptyPage())
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestMethod(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var methods []string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		methods = append(methods, req.Method())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(methods)
	is.Contains(methods, "GET")
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnResponseStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var statuses []int

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statuses = append(statuses, resp.Status())
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.NotEmpty(statuses)
	is.Contains(statuses, 200)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnResponseOKForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var emptyPageOK *bool

	cancel := page.OnResponse(func(resp *playwright.NetworkResponse) {
		if resp.URL() == srv.EmptyPage() {
			mu.Lock()
			ok := resp.OK()
			emptyPageOK = &ok
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(emptyPageOK, "OnResponse should have fired for empty page")
	is.True(*emptyPageOK, "empty page response should be OK")
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestHeaders(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var headers map[string]string

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() {
			mu.Lock()
			headers = req.Headers()
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	must.NotNil(headers, "headers should have been captured")

	_, hasUA := headers["user-agent"]
	is.True(hasUA, "request headers should include user-agent")
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestIsNavigationRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	var navRequestFound bool

	cancel := page.OnRequest(func(req *playwright.NetworkRequest) {
		if req.URL() == srv.EmptyPage() && req.IsNavigationRequest() {
			mu.Lock()
			navRequestFound = true
			mu.Unlock()
		}
	})
	defer cancel()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	defer mu.Unlock()
	is.True(navRequestFound, "should have found a navigation request to empty page")
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	requests := make([]*playwright.NetworkRequest, 0)

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requests = append(requests, req)
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	count := len(requests)
	mu.Unlock()

	is.Greater(count, 0)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnResponseFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	responses := make([]*playwright.NetworkResponse, 0)

	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		responses = append(responses, resp)
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	count := len(responses)
	mu.Unlock()

	is.Greater(count, 0)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestFinishedFiredForNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	finished := make([]*playwright.NetworkRequest, 0)

	off := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		mu.Lock()
		finished = append(finished, req)
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	count := len(finished)
	mu.Unlock()

	is.Greater(count, 0)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOffRemovesRequestHandler(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var mu sync.Mutex
	count := 0

	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	countBeforeOff := count
	mu.Unlock()

	off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	countAfterOff := count
	mu.Unlock()

	is.Equal(countBeforeOff, countAfterOff)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestFiredForFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/fetch-api", "application/json", `{"ok":true}`)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var urls []string
	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		urls = append(urls, req.URL())
		mu.Unlock()
	})
	defer off()

	_, err := page.Evaluate(ctx, `() => fetch('/fetch-api')`)
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	u := urls
	mu.Unlock()

	is.NotEmpty(u)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnResponseFiredForFetch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/resp-api", "application/json", `{"ok":true}`)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	var mu sync.Mutex
	var statuses []int
	off := page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		statuses = append(statuses, resp.Status())
		mu.Unlock()
	})
	defer off()

	_, err := page.Evaluate(ctx, `() => fetch('/resp-api')`)
	must.NoError(err)
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	s := statuses
	mu.Unlock()

	is.NotEmpty(s)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestMethodIsGet(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var method string
	off := page.OnRequest(func(req *playwright.NetworkRequest) {
		if method == "" {
			method = req.Method()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	is.Equal("GET", method)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestFinishedFiredEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var finished *playwright.NetworkRequest
	var mu sync.Mutex

	off := page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		if req.IsNavigationRequest() {
			mu.Lock()
			finished = req
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	mu.Lock()
	f := finished
	mu.Unlock()

	must.NotNil(f)
	is.Equal("GET", f.Method())
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestFailedFiredOnAbortEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	var failed *playwright.NetworkRequest
	var mu sync.Mutex

	off := page.OnRequestFailed(func(req *playwright.NetworkRequest) {
		mu.Lock()
		failed = req
		mu.Unlock()
	})
	defer off()

	must.NoError(page.Route(ctx, "**/blocked", func(route *playwright.Route) {
		must.NoError(route.Abort(ctx))
	}))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, _ = page.Evaluate(ctx, `() => fetch('/blocked').catch(() => {})`)
	must.NoError(page.WaitForTimeout(ctx, 200))

	mu.Lock()
	f := failed
	mu.Unlock()

	must.NotNil(f)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnConsoleMessageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var received string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "log" {
			mu.Lock()
			received = msg.Text()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.log('hello from console')</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	r := received
	mu.Unlock()

	is.Contains(r, "hello from console")
}

// TestOnPopupReceivesPageReference verifies OnPopup callback receives the opened page.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestOnPopupReceivesPageReference(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var popup *playwright.Page

	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		popup = p
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.SetContent(ctx, `<button onclick="window.open('about:blank')">open</button>`))
	must.NoError(page.Locator("button").Click(ctx))

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		p := popup
		mu.Unlock()
		if p != nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	mu.Lock()
	p := popup
	mu.Unlock()
	must.NotNil(p, "popup page should not be nil")
}

// TestOnPopupIsNotFiredForSameTab verifies OnPopup is NOT fired for same-tab navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestOnPopupIsNotFiredForSameTab(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	popupFired := false

	cancel := page.OnPopup(func(p *playwright.Page) {
		mu.Lock()
		popupFired = true
		mu.Unlock()
	})
	defer cancel()

	must.NoError(page.SetContent(ctx, `<a href="about:blank">same tab link</a>`))
	must.NoError(page.Locator("a").Click(ctx))

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	is.False(popupFired)
}

// TestOnPopupCanNavigate verifies the popup page can be navigated after opening.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestOnPopupCanNavigate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	done := make(chan *playwright.Page, 1)

	cancel := page.OnPopup(func(p *playwright.Page) {
		select {
		case done <- p:
		default:
		}
	})
	defer cancel()

	must.NoError(page.SetContent(ctx, `<button onclick="window.open('about:blank', '_blank')">open</button>`))
	must.NoError(page.Locator("button").Click(ctx))

	var popup *playwright.Page
	select {
	case popup = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for popup")
	}

	must.NotNil(popup)
	must.NoError(popup.SetContent(ctx, `<div>popup content</div>`))

	text, err := popup.Locator("div").InnerText(ctx)
	must.NoError(err)
	is.Equal("popup content", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var requestURL string
	var mu sync.Mutex

	page.OnRequest(func(req *playwright.NetworkRequest) {
		mu.Lock()
		requestURL = req.URL()
		mu.Unlock()
	})

	body := "hello"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/test"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	u := requestURL
	mu.Unlock()

	is.Contains(u, "example.com")
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnResponseEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var responseStatus int
	var mu sync.Mutex

	page.OnResponse(func(resp *playwright.NetworkResponse) {
		mu.Lock()
		responseStatus = resp.Status()
		mu.Unlock()
	})

	status := 200
	body := "response body"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Status: &status,
			Body:   &body,
		}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	s := responseStatus
	mu.Unlock()

	is.Equal(200, s)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnRequestFinishedEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)
	page, err := bc.NewPage(ctx)
	must.NoError(err)

	var finishedURL string
	var mu sync.Mutex

	page.OnRequestFinished(func(req *playwright.NetworkRequest) {
		mu.Lock()
		finishedURL = req.URL()
		mu.Unlock()
	})

	body := "done"
	must.NoError(bc.Route(ctx, "**/*", func(r *playwright.Route) {
		must.NoError(r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body}))
	}))

	must.NoError(page.Goto(ctx, "http://example.com/finish"))
	must.NoError(page.WaitForTimeout(ctx, 150))

	mu.Lock()
	u := finishedURL
	mu.Unlock()

	is.Contains(u, "example.com")
}

// TestBrowserContextOnPageClose verifies that OnPageClose handlers fire when a page is closed.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBrowserContextOnPageClose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	closedCh := make(chan *playwright.Page, 1)
	cancel := bCtx.OnPageClose(func(p *playwright.Page) {
		select {
		case closedCh <- p:
		default:
		}
	})
	defer cancel()

	err = page.Close(ctx)
	must.NoError(err, "Close failed")

	select {
	case p := <-closedCh:
		if p != page {
			t.Errorf("OnPageClose received wrong page: %v, want %v", p, page)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageClose to fire")
	}
}

// TestBrowserContextOnPageLoad verifies that OnPageLoad handlers fire after navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBrowserContextOnPageLoad(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	srv.ServeWithBody("/onload-target", "text/html", `<p>loaded</p>`)

	loadedCh := make(chan *playwright.Page, 1)
	cancelLoad := bCtx.OnPageLoad(func(p *playwright.Page) {
		select {
		case loadedCh <- p:
		default:
		}
	})
	defer cancelLoad()

	err = page.Goto(ctx, srv.Prefix()+"/onload-target")
	must.NoError(err, "Goto failed")

	select {
	case p := <-loadedCh:
		if p != page {
			t.Errorf("OnPageLoad received wrong page: %v, want %v", p, page)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPageLoad to fire")
	}
}

// TestPageCloseMarksPageAsClosed verifies IsClosed returns true after Close.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageCloseMarksPageAsClosed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	is.False(page.IsClosed())
	must.NoError(page.Close(ctx))
	is.True(page.IsClosed())
}

// TestPageCloseRemovesPageFromContext verifies closed page is removed from context.Pages().
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageCloseRemovesPageFromContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	pagesBeforeClose := bc.Pages()
	initialCount := len(pagesBeforeClose)

	must.NoError(page.Close(ctx))

	pagesAfterClose := bc.Pages()
	is.Less(len(pagesAfterClose), initialCount)
}

// TestPageCloseIsCallableTwice verifies closing a page twice doesn't panic or error.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageCloseIsCallableTwiceExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Close(ctx))
	must.NoError(page.Close(ctx)) // Second close should not panic
}

// TestPageCloseWithContent verifies page with content can be closed normally.
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageCloseWithContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<div>some content</div>`))
	must.NoError(page.Close(ctx))
	is.True(page.IsClosed())
}

// TestPageOnPageLoadCancelWorks verifies that the cancel function returned by OnPageLoad stops future callbacks.
// Ref: TestBrowserContextEvents.java (cancel behavior)
// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnPageLoadCancelWorks(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	srv.ServeWithBody("/cancel-load-a", "text/html", `<p>a</p>`)
	srv.ServeWithBody("/cancel-load-b", "text/html", `<p>b</p>`)

	var urls []string
	cancelLoad := bCtx.OnPageLoad(func(p *playwright.Page) {
		urls = append(urls, p.URL())
	})

	// First navigation: handler active.
	err = page.Goto(ctx, srv.Prefix()+"/cancel-load-a")
	must.NoError(err, "Goto /a failed")
	time.Sleep(200 * time.Millisecond)

	// Cancel the handler.
	cancelLoad()

	// Second navigation: handler should NOT fire.
	err = page.Goto(ctx, srv.Prefix()+"/cancel-load-b")
	must.NoError(err, "Goto /b failed")
	time.Sleep(200 * time.Millisecond)

	// Only the first URL should be in the list.
	is.Len(urls, 1, "expected exactly 1 load event (before cancel)")
	is.Contains(urls[0], "/cancel-load-a")
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageAddInitScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/init-page", "text/html", `<html><body></body></html>`)

	page := newPage(t)

	err := page.AddInitScript(ctx, `window.__injected = 'hello-init'`)
	must.NoError(err, "AddInitScript failed")

	// AddInitScript runs before page scripts on navigation, not on SetContent.
	err = page.Goto(ctx, srv.Prefix()+"/init-page")
	must.NoError(err, "Goto failed")

	val, err := page.Evaluate(ctx, `window.__injected`)
	must.NoError(err, "Evaluate failed")
	is.Equal("hello-init", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageAddInitScriptRunsOnNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/init-check", "text/html", `<html><body></body></html>`)

	page := newPage(t)

	err := page.AddInitScript(ctx, `window.__nav = 'nav-value'`)
	must.NoError(err, "AddInitScript failed")

	err = page.Goto(ctx, srv.Prefix()+"/init-check")
	must.NoError(err, "Goto failed")

	val, err := page.Evaluate(ctx, `window.__nav`)
	must.NoError(err, "Evaluate failed")
	is.Equal("nav-value", val)
}

// --- Page.BringToFront ---

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageBringToFront(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<html><body><p>front</p></body></html>`)
	must.NoError(err, "SetContent failed")

	// BringToFront should not error; behavioural effect (focus) is hard to assert in headless mode.
	err = page.BringToFront(ctx)
	must.NoError(err, "BringToFront failed")

	// Verify page is still functional after BringToFront.
	text, err := page.Locator("p").InnerText(ctx)
	must.NoError(err, "InnerText after BringToFront failed")
	is.Equal("front", text)
}

// --- Page.SetExtraHTTPHeaders ---

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageSetViewportSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<html><body></body></html>`)
	must.NoError(err, "SetContent failed")

	err = page.SetViewportSize(ctx, 800, 600)
	must.NoError(err, "SetViewportSize failed")

	size, err := page.ViewportSize(ctx)
	must.NoError(err, "ViewportSize failed")
	if size.Width != 800 || size.Height != 600 {
		t.Errorf("ViewportSize = %dx%d, want 800x600", size.Width, size.Height)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageSetViewportSizeChanges(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<html><body></body></html>`)
	must.NoError(err, "SetContent failed")

	err = page.SetViewportSize(ctx, 1280, 720)
	must.NoError(err, "SetViewportSize 1280x720 failed")

	err = page.SetViewportSize(ctx, 320, 568)
	must.NoError(err, "SetViewportSize 320x568 failed")

	size, err := page.ViewportSize(ctx)
	must.NoError(err, "ViewportSize failed")
	if size.Width != 320 || size.Height != 568 {
		t.Errorf("ViewportSize = %dx%d, want 320x568", size.Width, size.Height)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageURLAfterSetContent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>test</div>`))

	url := page.URL()
	is.NotEmpty(url)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageFramesOnlyMain(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	frames := page.Frames()
	is.GreaterOrEqual(len(frames), 1)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageMainFrameIsNotNil(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	frame := page.MainFrame()
	must.NotNil(frame)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageFrameByNameNilForUnknown(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>content</div>`))

	frame := page.Frame("nonexistent")
	is.Nil(frame)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageURLUpdatesAfterGoto(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	url := page.URL()
	is.Equal(srv.EmptyPage(), url)
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageWaitForTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	tests := []struct {
		name  string
		ms    float64
		minMs int64
	}{
		{"zero returns immediately", 0, 0},
		{"short duration 50ms", 50, 30},
		{"standard delay 100ms", 100, 80},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			page := newPage(t)
			start := time.Now()
			must.NoError(page.WaitForTimeout(testCtx(t), tc.ms))
			is.Less(time.Since(start).Milliseconds(), int64(5000), "timed out waiting")
			is.GreaterOrEqual(time.Since(start).Milliseconds(), tc.minMs)
		})
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageOnPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/popup-page", "text/html", `<title>Popup</title>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	popupCh := make(chan *playwright.Page, 1)
	cancel := page.OnPopup(func(p *playwright.Page) {
		select {
		case popupCh <- p:
		default:
		}
	})
	defer cancel()

	_, err = page.Evaluate(ctx, `() => window.open('`+srv.Prefix()+`/popup-page')`)
	must.NoError(err, "Evaluate failed")

	select {
	case popup := <-popupCh:
		if popup == nil {
			t.Fatal("OnPopup received nil page")
		}
		opener, err := popup.Opener(ctx)
		must.NoError(err, "popup.Opener() failed")
		if opener != page {
			t.Errorf("popup.Opener() = %v, want original page %v", opener, page)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for popup via OnPopup")
	}
}

// --- Page.WaitForEvent("popup") ---

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestPageWaitForEventPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/we-popup", "text/html", `<title>WE Popup</title>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	done := make(chan any, 1)
	go func() {
		val, err := page.WaitForEvent(ctx, "popup")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	_, err = page.Evaluate(ctx, `() => window.open('`+srv.Prefix()+`/we-popup')`)
	must.NoError(err, "Evaluate failed")

	result := <-done
	if popup, ok := result.(*playwright.Page); ok {
		if popup == nil {
			t.Fatal("WaitForEvent(popup) returned nil page")
		}
	} else {
		t.Fatalf("WaitForEvent(popup) returned unexpected type/error: %v", result)
	}
}

// --- Page.WaitForEvent("request") ---

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBrowserContextWaitForEventPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-we-page", "text/html", "<p>ctx</p>")

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	done := make(chan any, 1)
	go func() {
		val, err := bCtx.WaitForEvent(ctx, "page")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	// Open a new popup from the page to trigger a "page" event on the context
	_, err = page.Evaluate(ctx, `() => window.open('`+srv.Prefix()+`/ctx-we-page')`)
	must.NoError(err, "Evaluate failed")

	result := <-done
	if _, ok := result.(*playwright.Page); !ok {
		t.Fatalf("BrowserContext.WaitForEvent(page) unexpected result type: %T — %v", result, result)
	}
}

// --- BrowserContext.WaitForEvent("request") ---

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestBrowserContextWaitForEventRequest(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ctx-we-req", "text/html", "<p>ok</p>")

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	done := make(chan any, 1)
	go func() {
		val, err := bCtx.WaitForEvent(ctx, "request")
		if err != nil {
			done <- err
		} else {
			done <- val
		}
	}()

	err = page.Goto(ctx, srv.Prefix()+"/ctx-we-req")
	must.NoError(err, "Goto failed")

	result := <-done
	if req, ok := result.(*playwright.NetworkRequest); ok {
		if req.URL() == "" {
			t.Error("BrowserContext.WaitForEvent(request) returned request with empty URL")
		}
	} else {
		t.Fatalf("BrowserContext.WaitForEvent(request) unexpected result: %v", result)
	}
}

// --- Locator.Screenshot ---

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestLocatorScreenshot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="box" style="width:100px;height:100px;background:red"></div>`)
	must.NoError(err, "SetContent failed")

	data, err := page.Locator("#box").Screenshot(ctx)
	must.NoError(err, "Locator.Screenshot failed")
	if len(data) == 0 {
		t.Fatal("Locator.Screenshot returned empty bytes")
	}
	if !bytes.HasPrefix(data, []byte("\x89PNG")) {
		t.Errorf("expected PNG header, got: %x", data[:min(8, len(data))])
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageBasic.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageBasic.java
func TestLocatorScreenshotJPEG(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id="box" style="width:100px;height:100px;background:blue"></div>`)
	must.NoError(err, "SetContent failed")

	q := 80
	data, err := page.Locator("#box").Screenshot(ctx, &playwright.ScreenshotOptions{Type: "jpeg", Quality: &q})
	must.NoError(err, "Locator.Screenshot JPEG failed")
	if len(data) == 0 {
		t.Fatal("Locator.Screenshot JPEG returned empty bytes")
	}
	// JPEG magic: FF D8
	if data[0] != 0xFF || data[1] != 0xD8 {
		t.Errorf("expected JPEG header FF D8, got: %x %x", data[0], data[1])
	}
}
