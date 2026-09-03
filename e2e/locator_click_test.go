//go:build e2e

// Locator click E2E tests.
// Migration of: TestLocatorClick.java
package e2e

import (
	"context"
	"testing"

	"github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocatorClickWork verifies basic Locator.Click() works.
// Ref: TestLocatorClick.java#shouldWork
func TestLocatorClickWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))
	must.NoError(page.Locator("button").Click(ctx))

	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Clicked", result)
}

// TestLocatorClickWorkWithNodeRemoved verifies Click() works even after deleting window.Node.
// Ref: TestLocatorClick.java#shouldWorkWithNodeRemoved
func TestLocatorClickWorkWithNodeRemoved(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))

	_, err := page.Evaluate(ctx, "() => delete window['Node']")
	must.NoError(err)

	must.NoError(page.Locator("button").Click(ctx))

	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Clicked", result)
}

// TestLocatorClickWithShiftModifier verifies Click() sends shiftKey event when Shift modifier is used.
// Ref: TestLocatorClick.java#shouldSupportShiftClick (TestPageClick.java#shouldClickWithShiftKey)
func TestLocatorClickWithShiftModifier(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))

	// Override button listener to track shiftKey
	_, err := page.Evaluate(ctx, `() => {
		const button = document.querySelector('button');
		button.addEventListener('click', event => {
			window['result'] = event.shiftKey ? 'Shift' : 'NoShift';
		});
	}`)
	must.NoError(err)

	must.NoError(page.Locator("button").Click(ctx, &playwright.LocatorClickOptions{
		Modifiers: []string{"Shift"},
	}))

	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Shift", result)
}

// TestLocatorClickWithControlModifier verifies Click() sends ctrlKey event when Control modifier is used.
// Ref: TestLocatorClick.java#shouldSupportControlOrMetaModifier
func TestLocatorClickWithControlModifier(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))

	_, err := page.Evaluate(ctx, `() => {
		window['key'] = '';
		const button = document.querySelector('button');
		button.addEventListener('click', event => {
			window['key'] = event.ctrlKey ? 'control' : event.metaKey ? 'meta' : '';
		});
	}`)
	must.NoError(err)

	must.NoError(page.Locator("button").Click(ctx, &playwright.LocatorClickOptions{
		Modifiers: []string{"Control"},
	}))

	key, err := page.Evaluate(ctx, "() => window['key']")
	must.NoError(err)
	is.Equal("control", key)
}

// TestLocatorClickWithNoModifier verifies Click() has no modifier by default.
// Ref: TestLocatorClick.java#shouldClickWithNoModifiers
func TestLocatorClickWithNoModifier(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))

	_, err := page.Evaluate(ctx, `() => {
		const button = document.querySelector('button');
		button.addEventListener('click', event => {
			window['result'] = event.shiftKey ? 'Shift' : event.ctrlKey ? 'Control' : 'None';
		});
	}`)
	must.NoError(err)

	must.NoError(page.Locator("button").Click(ctx))

	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("None", result)
}

// TestLocatorDoubleClickButton verifies DblClick() fires double-click events.
// Ref: TestLocatorClick.java#shouldDoubleClickTheButton
func TestLocatorDoubleClickButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))

	_, err := page.Evaluate(ctx, `() => {
  window['double'] = false;
  const button = document.querySelector('button');
  button.addEventListener('dblclick', event => { window['double'] = true; });
}`)
	must.NoError(err)

	must.NoError(page.Locator("button").DblClick(ctx))

	doubled, err := page.Evaluate(ctx, "double")
	must.NoError(err)
	is.Equal(true, doubled)

	result, err := page.Evaluate(ctx, "result")
	must.NoError(err)
	is.Equal("Clicked", result)
}

// TestLocatorClickBehaviors groups the common click interaction scenarios into subtests.
func TestLocatorClickBehaviors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	tests := []struct {
		name  string
		html  string
		act   func(t *testing.T, page *playwright.Page, ctx context.Context)
		check func(t *testing.T, page *playwright.Page, ctx context.Context)
	}{
		{
			name: "no modifier by default",
			html: `<button id="btn" onclick="window.__shiftHeld=event.shiftKey">Click me</button>`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("#btn").Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				result, err := page.Evaluate(ctx, `() => window.__shiftHeld`)
				must.NoError(err)
				is.Equal(false, result)
			},
		},
		{
			name: "checks checkbox",
			html: `<input id="cb" type="checkbox">`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("#cb").Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				checked, err := page.Locator("#cb").IsChecked(ctx)
				must.NoError(err)
				is.True(checked)
			},
		},
		{
			name: "checks radio button",
			html: `<input type="radio" name="choice" id="r1" value="a"><input type="radio" name="choice" id="r2" value="b">`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("#r2").Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				checked, err := page.Locator("#r2").IsChecked(ctx)
				must.NoError(err)
				is.True(checked)
			},
		},
		{
			name: "fires click event via onclick attribute",
			html: `<div id="el" style="width:100px;height:50px" onclick="window.__clicked=true">Click target</div>`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("#el").Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				result, err := page.Evaluate(ctx, `() => window.__clicked`)
				must.NoError(err)
				is.Equal(true, result)
			},
		},
		{
			name: "clicks last locator match",
			html: `<button>One</button><button>Two</button><button onclick="window.__lastClicked=true">Three</button>`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("button").Last().Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				result, err := page.Evaluate(ctx, `() => window.__lastClicked`)
				must.NoError(err)
				is.Equal(true, result)
			},
		},
		{
			name: "navigates to anchor via link",
			html: `<a id="link" href="#anchor">Go to anchor</a>`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("#link").Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				url, err := page.Evaluate(ctx, `() => window.location.hash`)
				must.NoError(err)
				is.Equal("#anchor", url)
			},
		},
		{
			name: "increments counter on each click",
			html: `<button id="btn">+1</button><script>var count = 0; document.getElementById('btn').addEventListener('click', function() { count++; });</script>`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("#btn").Click(ctx))
				must.NoError(page.Locator("#btn").Click(ctx))
				must.NoError(page.Locator("#btn").Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				result, err := page.Evaluate(ctx, `() => count`)
				must.NoError(err)
				is.Equal(float64(3), result)
			},
		},
		{
			name: "focuses input element",
			html: `<input id="inp" type="text">`,
			act: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				must.NoError(page.Locator("#inp").Click(ctx))
			},
			check: func(t *testing.T, page *playwright.Page, ctx context.Context) {
				focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
				must.NoError(err)
				is.Equal("inp", focused)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := testCtx(t)
			page := newPage(t)
			must.NoError(page.SetContent(ctx, tt.html))
			tt.act(t, page, ctx)
			if tt.check != nil {
				tt.check(t, page, ctx)
			}
		})
	}
}

func TestLocatorClickSubmitsFormEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form id="f">
			<input type="text" name="user">
			<button type="submit" id="sub">Submit</button>
		</form>
		<script>
			var submitted = false;
			document.getElementById('f').addEventListener('submit', function(e) {
				e.preventDefault();
				submitted = true;
			});
		</script>
	`))

	must.NoError(page.Locator("#sub").Click(ctx))

	result, err := page.Evaluate(ctx, `() => submitted`)
	must.NoError(err)
	is.Equal(true, result)
}

func TestLocatorClickHidesElementEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="panel">Visible</div>
		<button id="toggle" onclick="document.getElementById('panel').style.display='none'">Hide</button>
	`))

	must.NoError(page.Locator("#toggle").Click(ctx))

	visible, err := page.Locator("#panel").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

func TestLocatorClickAddClassEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el">Element</div>
		<button id="btn" onclick="document.getElementById('el').classList.add('active')">Activate</button>
	`))

	must.NoError(page.Locator("#btn").Click(ctx))

	hasClass, err := page.Evaluate(ctx, `() => document.getElementById('el').classList.contains('active')`)
	must.NoError(err)
	is.Equal(true, hasClass)
}

func TestLocatorClickUpdatesTextEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p id="msg">Before</p>
		<button id="btn" onclick="document.getElementById('msg').textContent='After'">Change</button>
	`))

	must.NoError(page.Locator("#btn").Click(ctx))

	text, err := page.Locator("#msg").InnerText(ctx)
	must.NoError(err)
	is.Equal("After", text)
}

func TestClickLinkEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<a id="a" href="#" onclick="document.getElementById('out').textContent='clicked'">Link</a>
		<span id="out"></span>
	`))

	must.NoError(page.Locator("#a").Click(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("clicked", out)
}

func TestClickToggleCheckboxEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	must.NoError(page.Locator("#chk").Click(ctx))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)

	must.NoError(page.Locator("#chk").Click(ctx))

	checked, err = page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

func TestClickSelectOptionEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel" onchange="document.getElementById('out').textContent=this.value">
			<option value="x">X</option>
			<option value="y">Y</option>
		</select>
		<span id="out"></span>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "y")
	must.NoError(err)

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("y", out)
}

func TestClickImageEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img id="img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			style="width:50px;height:50px;"
			onclick="document.getElementById('out').textContent='clicked'">
		<span id="out"></span>
	`))

	must.NoError(page.Locator("#img").Click(ctx))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("clicked", out)
}

func TestClickDivTriggersEventEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" style="width:100px;height:100px;background:#eee;"
			onclick="this.textContent='div clicked'">Click div</div>
	`))

	must.NoError(page.Locator("#d").Click(ctx))

	text, err := page.Locator("#d").TextContent(ctx)
	must.NoError(err)
	is.Equal("div clicked", text)
}

func TestClickFormResetEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form>
			<input id="inp" type="text" value="initial">
			<button type="reset" id="reset">Reset</button>
		</form>
	`))

	must.NoError(page.Locator("#inp").Fill(ctx, "changed"))
	must.NoError(page.Locator("#reset").Click(ctx))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("initial", val)
}

func TestClickExpandsDetailsEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<details id="det">
			<summary id="sum">Click to expand</summary>
			<p id="content">Hidden content</p>
		</details>
	`))

	must.NoError(page.Locator("#sum").Click(ctx))

	open, err := page.Locator("#det").Evaluate(ctx, `el => el.open`)
	must.NoError(err)
	is.Equal(true, open)
}

func TestClickChangesTabIndexEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div>
			<button id="tab1" onclick="document.getElementById('tab1').classList.add('active')">Tab 1</button>
			<button id="tab2" onclick="document.getElementById('tab2').classList.add('active')">Tab 2</button>
		</div>
	`))

	must.NoError(page.Locator("#tab2").Click(ctx))

	classes, err := page.Locator("#tab2").Evaluate(ctx, `el => el.className`)
	must.NoError(err)
	s, ok := classes.(string)
	is.True(ok)
	is.Contains(s, "active")
}

func TestClickNumberInputEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="n" type="number" value="5">`))

	must.NoError(page.Locator("#n").Click(ctx))

	focused, err := page.Locator("#n").Evaluate(ctx, `el => document.activeElement === el`)
	must.NoError(err)
	is.Equal(true, focused)
}

// TestLocatorDblClickFiresClickTwice verifies DblClick fires two click events.
// Ref: TestLocatorDblClick.java#shouldFireDoubleClickEvent
func TestLocatorDblClickFiresClickTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="counter" onclick="this.dataset.clicks = (parseInt(this.dataset.clicks||0)+1)">click me</div>
	`))

	must.NoError(page.Locator("#counter").DblClick(ctx))

	clicks, err := page.Evaluate(ctx, `() => document.getElementById('counter').dataset.clicks`)
	must.NoError(err)
	is.Equal("2", clicks)
}

// TestLocatorDblClickFiresDblClickEvent verifies DblClick fires the dblclick DOM event.
// Ref: TestLocatorDblClick.java#shouldFireDblClickEvent
func TestLocatorDblClickFiresDblClickEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" ondblclick="this.dataset.dblclicked='yes'">double-click</div>
	`))

	must.NoError(page.Locator("#target").DblClick(ctx))

	result, err := page.Evaluate(ctx, `() => document.getElementById('target').dataset.dblclicked`)
	must.NoError(err)
	is.Equal("yes", result)
}

// TestLocatorDblClickOnButton verifies DblClick works on a button element.
// Ref: TestLocatorDblClick.java#shouldWorkOnButton
func TestLocatorDblClickOnButton(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" ondblclick="this.textContent='double-clicked'">click me</button>
	`))

	must.NoError(page.Locator("button").DblClick(ctx))

	text, err := page.Locator("button").InnerText(ctx)
	must.NoError(err)
	is.Equal("double-clicked", text)
}

// TestDblClickOnInput verifies DblClick selects word in text input.
// Ref: TestLocatorDblClick.java#shouldSelectWordInInput
func TestDblClickOnInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" value="hello world">`))

	must.NoError(page.Locator("#inp").DblClick(ctx))

	// After dblclick on input, some text should be selected
	selected, err := page.Evaluate(ctx, `() => {
		const inp = document.getElementById('inp');
		return inp.value.substring(inp.selectionStart, inp.selectionEnd);
	}`)
	must.NoError(err)
	is.NotEmpty(selected)
}

// TestDblClickFiresBothClickEvents verifies DblClick fires two click-type events.
// Ref: TestLocatorDblClick.java#shouldFireTwoClickEvents
func TestDblClickFiresBothClickEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px"
		     onclick="window.__clicks=(window.__clicks||0)+1">dbl click</div>
	`))

	must.NoError(page.Locator("#el").DblClick(ctx))

	clicks, err := page.Evaluate(ctx, `() => window.__clicks`)
	must.NoError(err)
	is.Equal(float64(2), clicks)
}

// TestDblClickWithMultipleElements verifies DblClick works on first when multiple match.
// Ref: TestLocatorDblClick.java#shouldWorkOnFirst
func TestDblClickWithFirstLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button ondblclick="window.__firstDbl=true">first</button>
		<button>second</button>
	`))

	must.NoError(page.Locator("button").First().DblClick(ctx))

	result, err := page.Evaluate(ctx, `() => window.__firstDbl`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestDblClickFiresEventEx3 verifies DblClick fires dblclick event.
// Ref: TestLocatorDblClick.java#shouldFireDblClickEvent
func TestDblClickFiresEventEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" style="width:100px;height:100px;">DblClick</div>
		<script>
			var count = 0;
			document.getElementById('el').addEventListener('dblclick', function() { count++; });
		</script>
	`))

	must.NoError(page.Locator("#el").DblClick(ctx))

	result, err := page.Evaluate(ctx, `() => count`)
	must.NoError(err)
	is.Equal(float64(1), result)
}

// TestDblClickOnImageEx3 verifies DblClick on image element.
// Ref: TestLocatorDblClick.java#shouldWorkOnImage
func TestDblClickOnImageEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<img id="img" src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
			style="width:50px;height:50px;" ondblclick="window.__imgDbl=true">
	`))

	must.NoError(page.Locator("#img").DblClick(ctx))

	result, err := page.Evaluate(ctx, `() => window.__imgDbl`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestDblClickCounterIncreasesEx3 verifies DblClick click count works.
// Ref: TestLocatorDblClick.java#shouldIncreaseClickCount
func TestDblClickCounterIncreasesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<button id="btn" style="width:100px;height:50px;">Button</button>
		<script>
			var dblCount = 0;
			document.getElementById('btn').addEventListener('dblclick', function() { dblCount++; });
		</script>
	`))

	must.NoError(page.Locator("#btn").DblClick(ctx))
	must.NoError(page.Locator("#btn").DblClick(ctx))

	result, err := page.Evaluate(ctx, `() => dblCount`)
	must.NoError(err)
	is.Equal(float64(2), result)
}

// Ref: TestLocatorClick.java#shouldSupportCotrolOrMetaModifier
func TestLocatorClickSupportCotrolOrMetaModifier(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/title.html", "text/html", "<title>Title Page</title><body>Title Page</body>")
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, "<a href='"+srv.Prefix()+"/title.html'>Go</a>"))

	pageCh := make(chan *playwright.Page, 1)
	cancelPage := bCtx.OnPage(func(p *playwright.Page) {
		pageCh <- p
	})
	defer cancelPage()

	must.NoError(page.GetByText("Go").Click(ctx, &playwright.LocatorClickOptions{
		Modifiers: []string{"ControlOrMeta"},
	}))

	var newPage *playwright.Page
	select {
	case newPage = <-pageCh:
	case <-ctx.Done():
		t.Fatal("timed out waiting for new page")
	}
	is.NotNil(newPage)
	// Allow navigation to complete.
	_ = newPage.WaitForLoadState(ctx, "load")
	is.Contains(newPage.URL(), "/title.html")
}

// Ref: TestLocatorClick.java#shouldClickWithTweenedMouseMovement
func TestLocatorClickWithTweenedMouseMovement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.ButtonPage()))

	steps := 10
	must.NoError(page.Locator("button").Click(ctx, &playwright.LocatorClickOptions{
		Steps: &steps,
	}))

	result, err := page.Evaluate(ctx, "() => window['result']")
	must.NoError(err)
	is.Equal("Clicked", result)
}

// Ref: TestLocatorClick.java#shouldNotScrollWhenScrollIsNone
func TestLocatorClickNotScrollWhenScrollIsNone(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Place button in viewport so scroll=none doesn't prevent click
	must.NoError(page.SetContent(ctx, `<button id="btn" onclick="window.__clicked=true" style="position:fixed;top:10px;left:10px;">Click</button>`))

	scroll := "none"
	must.NoError(page.Locator("#btn").Click(ctx, &playwright.LocatorClickOptions{
		Scroll: &scroll,
	}))

	result, err := page.Evaluate(ctx, "() => window.__clicked")
	must.NoError(err)
	is.Equal(true, result)
}

// Ref: TestLocatorClick.java#shouldClickInViewportElementWhenScrollIsNone
func TestLocatorClickInViewportElementWhenScrollIsNone(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Element is already in viewport; scroll=none should still allow clicking it
	must.NoError(page.SetContent(ctx, `
		<button id="inview" onclick="window.__inviewClicked=true"
			style="position:fixed;top:50px;left:50px;">In viewport</button>
	`))

	scroll := "none"
	must.NoError(page.Locator("#inview").Click(ctx, &playwright.LocatorClickOptions{
		Scroll: &scroll,
	}))

	result, err := page.Evaluate(ctx, "() => window.__inviewClicked")
	must.NoError(err)
	is.Equal(true, result)
}

// TestDblClickOnSpanEx4 verifies DblClick works on span elements.
// Ref: TestLocatorDblClick.java#shouldDblClickSpan
func TestDblClickOnSpanEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span id="sp" ondblclick="this.textContent='double clicked'">Click me</span>
	`))

	must.NoError(page.Locator("#sp").DblClick(ctx))

	text, err := page.Locator("#sp").TextContent(ctx)
	must.NoError(err)
	is.Equal("double clicked", text)
}
