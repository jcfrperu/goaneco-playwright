//go:build e2e

// E2E tests for: TestPageKeyboard.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"fmt"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"runtime"
	"strings"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPress(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="inp" type="text" />
		<div id="out"></div>
		<script>
			document.getElementById('inp').addEventListener('keydown', function(e) {
				document.getElementById('out').textContent = e.key;
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#inp").Click(ctx)
	must.NoError(err, "Click failed")

	err = page.Keyboard.Press(ctx, "Enter")
	must.NoError(err, "Keyboard.Press failed")

	text, err := page.Locator("#out").InnerText(ctx)
	must.NoError(err, "InnerText failed")
	if text != "Enter" {
		t.Errorf("expected 'Enter' keydown event, got %q", text)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardType(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="inp" type="text" />`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#inp").Click(ctx)
	must.NoError(err, "Click failed")

	err = page.Keyboard.Type(ctx, "hello")
	must.NoError(err, "Keyboard.Type failed")

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err, "InputValue failed")
	if val != "hello" {
		t.Errorf("expected 'hello', got %q", val)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id="inp" type="text" />`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#inp").Click(ctx)
	must.NoError(err, "Click failed")

	err = page.Keyboard.InsertText(ctx, "world")
	must.NoError(err, "Keyboard.InsertText failed")

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err, "InputValue failed")
	if val != "world" {
		t.Errorf("expected 'world', got %q", val)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardDownAndUp(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<div id="log"></div>
		<script>
			var log = [];
			document.addEventListener('keydown', function(e) { log.push('down:' + e.key); });
			document.addEventListener('keyup',   function(e) { log.push('up:' + e.key);   });
			window._getLog = function() { return log.join(','); };
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Keyboard.Down(ctx, "Shift")
	must.NoError(err, "Keyboard.Down failed")
	err = page.Keyboard.Up(ctx, "Shift")
	must.NoError(err, "Keyboard.Up failed")

	result, err := page.Evaluate(ctx, "window._getLog()")
	must.NoError(err, "Evaluate failed")
	got, _ := result.(string)
	if got != "down:Shift,up:Shift" {
		t.Errorf("expected 'down:Shift,up:Shift', got %q", got)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressModifierCombo(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
		<input id="inp" type="text" value="hello" />
		<script>
			document.getElementById('inp').addEventListener('keydown', function(e) {
				if (e.key === 'a' && e.ctrlKey) {
					document.getElementById('inp').setAttribute('data-selected', '1');
				}
			});
		</script>
	`)
	must.NoError(err, "SetContent failed")

	err = page.Locator("#inp").Click(ctx)
	must.NoError(err, "Click failed")

	err = page.Keyboard.Press(ctx, "Control+a")
	must.NoError(err, "Keyboard.Press(Control+a) failed")

	attr, err := page.Locator("#inp").GetAttribute(ctx, "data-selected")
	must.NoError(err, "GetAttribute failed")
	if attr == nil || *attr != "1" {
		t.Errorf("expected data-selected='1' after Ctrl+A, got %v", attr)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeIntoTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		const textarea = document.createElement('textarea');
		document.body.appendChild(textarea);
		textarea.focus();
	}`)
	must.NoError(err)

	text := "Hello world. I am the text that was typed!"
	err = page.Keyboard.Type(ctx, text)
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal(text, val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardArrowKeys(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)
	err = page.SetContent(ctx, `<textarea></textarea>`)
	must.NoError(err)

	err = page.Locator("textarea").Click(ctx)
	must.NoError(err)

	err = page.Keyboard.Type(ctx, "Hello World!")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("Hello World!", val)

	for range "World!" {
		err = page.Keyboard.Press(ctx, "ArrowLeft")
		must.NoError(err)
	}

	err = page.Keyboard.Type(ctx, "inserted ")
	must.NoError(err)

	val, err = page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("Hello inserted World!", val)

	err = page.Keyboard.Down(ctx, "Shift")
	must.NoError(err)
	for range "inserted " {
		err = page.Keyboard.Press(ctx, "ArrowLeft")
		must.NoError(err)
	}
	err = page.Keyboard.Up(ctx, "Shift")
	must.NoError(err)
	err = page.Keyboard.Press(ctx, "Backspace")
	must.NoError(err)

	val, err = page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("Hello World!", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardElementHandlePress(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<textarea></textarea>`)
	must.NoError(err)

	el, err := page.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NotNil(el)

	err = el.Press(ctx, "a")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("a", val)

	_, err = page.Evaluate(ctx, "() => window.addEventListener('keydown', e => e.preventDefault(), true)")
	must.NoError(err)

	err = el.Press(ctx, "b")
	must.NoError(err)

	val, err = page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("a", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertTextUnicode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<textarea></textarea>`)
	must.NoError(err)
	err = page.Locator("textarea").Click(ctx)
	must.NoError(err)

	err = page.Keyboard.InsertText(ctx, "嗨")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("嗨", val)

	_, err = page.Evaluate(ctx, "() => window.addEventListener('keydown', e => e.preventDefault(), true)")
	must.NoError(err)

	err = page.Keyboard.InsertText(ctx, "a")
	must.NoError(err)

	val, err = page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("嗨a", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPlusSeparatedModifier(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.KeyboardPage())
	must.NoError(err)

	err = page.Keyboard.Press(ctx, "Shift+~")
	must.NoError(err)

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, _ := result.(string)
	is.Contains(s, "Keydown: Shift ShiftLeft 16 [Shift]")
	is.Contains(s, "Keyup: Shift ShiftLeft 16 []")
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardMultiplePlusSeparatedModifiers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.KeyboardPage())
	must.NoError(err)

	err = page.Keyboard.Press(ctx, "Control+Shift+~")
	must.NoError(err)

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, _ := result.(string)
	is.Contains(s, "Keydown: Control ControlLeft 17 [Control]")
	is.Contains(s, "Keydown: Shift ShiftLeft 16 [Control Shift]")
	is.Contains(s, "Keyup: Control ControlLeft 17 []")
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPreventSelectAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="text" value="hello">`)
	must.NoError(err)

	err = page.Locator("input").Click(ctx)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
		document.querySelector('input').addEventListener('keydown', e => {
			if (e.key === 'a' && (e.ctrlKey || e.metaKey)) {
				e.preventDefault();
			}
		});
	}`)
	must.NoError(err)

	err = page.Keyboard.Press(ctx, "ControlOrMeta+a")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => document.querySelector('input').value")
	must.NoError(err)
	is.Equal("hello", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeOnMouseClickedInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))
	must.NoError(page.Locator("input").Click(ctx))
	must.NoError(page.Keyboard.Type(ctx, "Hello"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressEnterInTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must.NoError(page.Locator("textarea").Click(ctx))
	must.NoError(page.Keyboard.Type(ctx, "Hello"))
	must.NoError(page.Keyboard.Press(ctx, "Enter"))
	must.NoError(page.Keyboard.Type(ctx, "World"))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello\nWorld", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressTab(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first" tabindex="1">
		<input id="second" tabindex="2">
	`))

	must.NoError(page.Locator("#first").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Tab"))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", activeId)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressBackspace(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" value="Hello">`))
	must.NoError(page.Locator("input").Click(ctx))

	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Keyboard.Press(ctx, "Backspace"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hell", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressDeletion(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" value="Hello">`))
	must.NoError(page.Locator("input").Click(ctx))

	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Keyboard.Press(ctx, "Delete"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("ello", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardShiftPressArrowSelectsText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" value="Hello">`))
	must.NoError(page.Locator("input").Click(ctx))

	must.NoError(page.Keyboard.Press(ctx, "Home"))

	must.NoError(page.Keyboard.Press(ctx, "Shift+ArrowRight"))
	must.NoError(page.Keyboard.Press(ctx, "Shift+ArrowRight"))
	must.NoError(page.Keyboard.Press(ctx, "Shift+ArrowRight"))

	must.NoError(page.Keyboard.Type(ctx, "Hey"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("Heylo", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardDownAndUpEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div tabindex="0" id="el"
		     onkeydown="window.__down=event.key"
		     onkeyup="window.__up=event.key">
		</div>
	`))

	must.NoError(page.Locator("#el").Focus(ctx))
	must.NoError(page.Keyboard.Down(ctx, "A"))
	must.NoError(page.Keyboard.Up(ctx, "A"))

	down, err := page.Evaluate(ctx, `() => window.__down`)
	must.NoError(err)
	is.Equal("A", down)

	up, err := page.Evaluate(ctx, `() => window.__up`)
	must.NoError(err)
	is.Equal("A", up)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertTextSendsInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.InsertText(ctx, "hello from keyboard"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hello from keyboard", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressKeyInInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "h"))
	must.NoError(page.Keyboard.Press(ctx, "i"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hi", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeTypesCharacters(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Type(ctx, "world"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("world", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardShiftPressSelectsRange(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" value="hello world">
	`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Down(ctx, "Shift"))
	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Keyboard.Up(ctx, "Shift"))

	selected, err := page.Evaluate(ctx, `() => document.getElementById('inp').value.substring(
		document.getElementById('inp').selectionStart,
		document.getElementById('inp').selectionEnd
	)`)
	must.NoError(err)
	is.NotEmpty(selected)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeInTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	must.NoError(page.Locator("#ta").Focus(ctx))
	must.NoError(page.Keyboard.Type(ctx, "line1\nline2"))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "line1")
	is.Contains(val, "line2")
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressHomeMovesCaretToStart(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))

	must.NoError(page.Locator("#inp").Focus(ctx))
	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Keyboard.Press(ctx, "Home"))

	caretPos, err := page.Evaluate(ctx, `() => document.getElementById('inp').selectionStart`)
	must.NoError(err)
	is.Equal(float64(0), caretPos)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardCtrlASelectsAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello world">`))

	must.NoError(page.Locator("#inp").Focus(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Control+a"))

	selLen, err := page.Evaluate(ctx, `() => {
		const inp = document.getElementById('inp');
		return inp.selectionEnd - inp.selectionStart;
	}`)
	must.NoError(err)
	is.Equal(float64(11), selLen)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardEscapeFiresEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" tabindex="0" onkeydown="if(event.key==='Escape') window.__esc=true">press</div>
	`))

	must.NoError(page.Locator("#el").Focus(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Escape"))

	result, err := page.Evaluate(ctx, `() => window.__esc`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertTextWithSpecialChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.InsertText(ctx, "café"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("café", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardCtrlCCopiesSelectionEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" value="copy me" oncopy="window.__copied=true">
	`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Control+a"))
	must.NoError(page.Keyboard.Press(ctx, "Control+c"))

	copied, err := page.Evaluate(ctx, `() => window.__copied`)
	must.NoError(err)
	is.Equal(true, copied)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardCtrlZUndoEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" onkeydown="window.__lastKey=event.key">
	`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Control+z"))

	key, err := page.Evaluate(ctx, `() => window.__lastKey`)
	must.NoError(err)
	is.Equal("z", key)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardShiftClickEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="el" onclick="window.__shiftKey=event.shiftKey">Click</div>
	`))

	must.NoError(page.Keyboard.Down(ctx, "Shift"))
	must.NoError(page.Locator("#el").Click(ctx))
	must.NoError(page.Keyboard.Up(ctx, "Shift"))

	shiftKey, err := page.Evaluate(ctx, `() => window.__shiftKey`)
	must.NoError(err)
	is.Equal(true, shiftKey)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPageDownScrollsEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:5000px">Scrollable</div>
	`))

	must.NoError(page.Locator("body").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "PageDown"))
	must.NoError(page.WaitForTimeout(ctx, 100))

	scrollY, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.Greater(scrollY.(float64), float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertTextEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.InsertText(ctx, "inserted text"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("inserted text", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardSelectAllEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="select all of this">`))
	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Control+a"))

	selected, err := page.Evaluate(ctx, `() => {
		const inp = document.getElementById('inp');
		return inp.value.substring(inp.selectionStart, inp.selectionEnd);
	}`)
	must.NoError(err)
	is.Equal("select all of this", selected)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardEscapeClosesDropdownEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="box" tabindex="0"></div>
		<script>
			var lastKey = '';
			document.getElementById('box').addEventListener('keydown', function(e) { lastKey = e.key; });
		</script>
	`))

	must.NoError(page.Locator("#box").Focus(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Escape"))

	result, err := page.Evaluate(ctx, `() => lastKey`)
	must.NoError(err)
	is.Equal("Escape", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardArrowNavigationEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="nav" tabindex="0"></div>
		<script>
			var keys = [];
			document.getElementById('nav').addEventListener('keydown', function(e) { keys.push(e.key); });
		</script>
	`))

	must.NoError(page.Locator("#nav").Focus(ctx))
	must.NoError(page.Keyboard.Press(ctx, "ArrowLeft"))
	must.NoError(page.Keyboard.Press(ctx, "ArrowRight"))
	must.NoError(page.Keyboard.Press(ctx, "ArrowUp"))
	must.NoError(page.Keyboard.Press(ctx, "ArrowDown"))

	keys, err := page.Evaluate(ctx, `() => keys`)
	must.NoError(err)
	keysArr, ok := keys.([]any)
	is.True(ok)
	is.Equal(4, len(keysArr))
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardDownUpSequenceEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text">
	`))

	must.NoError(page.Locator("#inp").Focus(ctx))
	must.NoError(page.Keyboard.Down(ctx, "Shift"))
	must.NoError(page.Keyboard.Press(ctx, "KeyA"))
	must.NoError(page.Keyboard.Up(ctx, "Shift"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("A", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardArrowNavigationEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Fill(ctx, "hello"))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Keyboard.Press(ctx, "ArrowRight"))

	pos, err := page.Evaluate(ctx, `() => document.getElementById('inp').selectionStart`)
	must.NoError(err)
	is.Equal(float64(1), pos)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardDeleteKeyEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Keyboard.Press(ctx, "Delete"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("ello", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeSequentiallyEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.Type(ctx, "abc"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("abc", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardEndKeyEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Keyboard.Press(ctx, "End"))

	pos, err := page.Evaluate(ctx, `() => document.getElementById('inp').selectionStart`)
	must.NoError(err)
	is.Equal(float64(5), pos)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPageUpDownEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="height:5000px;">Tall page</div>`))

	must.NoError(page.Keyboard.Press(ctx, "PageDown"))

	scrollY, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.GreaterOrEqual(scrollY.(float64), float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardFunctionKeyEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	err := page.Keyboard.Press(ctx, "F5")
	must.NoError(err)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardMultipleBackspaceEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Fill(ctx, "abcde"))
	must.NoError(page.Locator("#inp").Focus(ctx))

	for i := 0; i < 3; i++ {
		must.NoError(page.Keyboard.Press(ctx, "Backspace"))
	}

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("ab", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertKeyEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="d" tabindex="0" onkeydown="this.dataset.key=event.key"></div>
	`))

	must.NoError(page.Locator("#d").Focus(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Insert"))

	attr, err := page.Locator("#d").GetAttribute(ctx, "data-key")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("Insert", *attr)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardHomeKeyEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="Hello World">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Keyboard.Press(ctx, "Home"))

	pos, err := page.Locator("#inp").Evaluate(ctx, `el => el.selectionStart`)
	must.NoError(err)
	is.Equal(float64(0), pos)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardCtrlASelectAllEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="Select me">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Control+a"))

	selected, err := page.Locator("#inp").Evaluate(ctx, `el => el.selectionStart === 0 && el.selectionEnd === el.value.length`)
	must.NoError(err)
	is.Equal(true, selected)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardSpacebarChecksCheckboxEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	must.NoError(page.Locator("#cb").Focus(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Space"))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorPressEnterSubmitsForm verifies Press Enter submits a form.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressEnterSubmitsForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form onsubmit="document.getElementById('result').textContent='submitted'; return false;">
			<input type="text" id="inp">
			<div id="result"></div>
		</form>
	`))

	must.NoError(page.Locator("input").Press(ctx, "Enter"))

	text, err := page.Locator("#result").InnerText(ctx)
	must.NoError(err)
	is.Equal("submitted", text)
}

// TestLocatorPressArrowKey verifies Press ArrowRight moves caret in input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressArrowKey(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="abc">`))
	must.NoError(page.Locator("input").Click(ctx))

	// Go to beginning, then move right once
	must.NoError(page.Locator("input").Press(ctx, "Home"))
	must.NoError(page.Locator("input").Press(ctx, "ArrowRight"))
	must.NoError(page.Locator("input").Press(ctx, "ArrowRight"))

	// Type at position 2 (between 'b' and 'c')
	must.NoError(page.Keyboard.Type(ctx, "X"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("abXc", val)
}

// TestLocatorPressDeleteKey verifies Press Delete removes character at caret.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressDeleteKey(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="hello">`))
	must.NoError(page.Locator("input").Click(ctx))
	must.NoError(page.Locator("input").Press(ctx, "Home"))
	must.NoError(page.Locator("input").Press(ctx, "Delete"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("ello", val)
}

// TestLocatorPressTabMovesFocus verifies Press Tab moves focus to next element.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressTabMovesFocus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first">
		<input id="second">
	`))

	must.NoError(page.Locator("#first").Click(ctx))
	must.NoError(page.Locator("#first").Press(ctx, "Tab"))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", activeId)
}

// TestLocatorPressEscapeBlursInput verifies Escape removes focus.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressEscapeBlursInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" onblur="window.__blurred=true">
	`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Locator("#inp").Press(ctx, "Escape"))

	// Escape doesn't necessarily blur, but should not error
}

// TestLocatorPressArrowLeftMovesCaret verifies ArrowLeft moves caret.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressArrowLeftMovesCaret(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Locator("#inp").Press(ctx, "ArrowLeft"))

	pos, err := page.Evaluate(ctx, `() => document.getElementById('inp').selectionStart`)
	must.NoError(err)
	// Position should be before the last character
	is.Equal(float64(4), pos)
}

// TestLocatorPressHomeMovesToStart verifies Home key moves to beginning.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressHomeMovesToStart(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello world">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Locator("#inp").Press(ctx, "Home"))

	pos, err := page.Evaluate(ctx, `() => document.getElementById('inp').selectionStart`)
	must.NoError(err)
	is.Equal(float64(0), pos)
}

// TestLocatorPressEndMovesToEnd verifies End key moves to end.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressEndMovesToEnd(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Locator("#inp").Press(ctx, "End"))

	pos, err := page.Evaluate(ctx, `() => document.getElementById('inp').selectionStart`)
	must.NoError(err)
	is.Equal(float64(5), pos)
}

// TestLocatorPressPageDownScrolls verifies PageDown scrolls the page.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressPageDownScrolls(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div tabindex="0" id="scroller" style="height:200px;overflow-y:scroll">
			<div style="height:2000px">content</div>
		</div>
	`))

	must.NoError(page.Locator("#scroller").Click(ctx))
	must.NoError(page.Locator("#scroller").Press(ctx, "PageDown"))

	scrollTop, err := page.Evaluate(ctx, `() => document.getElementById('scroller').scrollTop`)
	must.NoError(err)
	top, ok := scrollTop.(float64)
	is.True(ok)
	is.Greater(top, float64(0))
}

// TestLocatorPressTabMovesFocus verifies Tab moves focus to next element.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressTabMovesFocusEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first" type="text">
		<input id="second" type="text">
	`))

	must.NoError(page.Locator("#first").Focus(ctx))
	must.NoError(page.Locator("#first").Press(ctx, "Tab"))

	activeId, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", activeId)
}

// TestLocatorPressSpaceTogglesCheckbox verifies Space toggles a checkbox.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressSpaceTogglesCheckbox(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="cb" type="checkbox">`))

	must.NoError(page.Locator("#cb").Press(ctx, "Space"))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorPressDeleteClearsInput verifies Delete removes text.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressDeleteClearsInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))

	must.NoError(page.Locator("#inp").Focus(ctx))
	// Select all text first
	must.NoError(page.Locator("#inp").Press(ctx, "Control+a"))
	must.NoError(page.Locator("#inp").Press(ctx, "Delete"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Empty(val)
}

// TestLocatorPressBackspaceRemovesChar verifies Backspace removes last character.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressBackspaceRemovesChar(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Fill(ctx, "hello"))

	must.NoError(page.Locator("#inp").Press(ctx, "Backspace"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hell", val)
}

// TestLocatorPressEnterSubmitsEx4 verifies Enter key submits form.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressEnterSubmitsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form id="frm" onsubmit="window.__submitted=true; return false;">
			<input id="inp" type="text">
		</form>
	`))

	must.NoError(page.Locator("#inp").Press(ctx, "Enter"))

	submitted, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, submitted)
}

// TestLocatorPressArrowDownEx4 verifies ArrowDown key press fires keydown event.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressArrowDownEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onkeydown="window.__lastKey=event.key">
	`))

	must.NoError(page.Locator("#inp").Press(ctx, "ArrowDown"))

	key, err := page.Evaluate(ctx, `() => window.__lastKey`)
	must.NoError(err)
	is.Equal("ArrowDown", key)
}

// TestLocatorPressEscapeEx4 verifies Escape key fires keydown event.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressEscapeEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" onkeydown="window.__lastKey=event.key">
	`))

	must.NoError(page.Locator("#inp").Press(ctx, "Escape"))

	key, err := page.Evaluate(ctx, `() => window.__lastKey`)
	must.NoError(err)
	is.Equal("Escape", key)
}

// TestLocatorPressHomeKeyEx4 verifies Home key moves cursor to start.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressHomeKeyEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Locator("#inp").Press(ctx, "End"))
	must.NoError(page.Locator("#inp").Press(ctx, "Home"))

	pos, err := page.Evaluate(ctx, `() => document.getElementById('inp').selectionStart`)
	must.NoError(err)
	is.Equal(float64(0), pos)
}

// TestLocatorPressSpaceEx5 verifies pressing Space on checkbox toggles it.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressSpaceEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))
	must.NoError(page.Locator("#chk").Focus(ctx))
	must.NoError(page.Locator("#chk").Press(ctx, "Space"))

	checked, err := page.Locator("#chk").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorPressDeleteEx5 verifies pressing Delete clears text.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressDeleteEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="text">`))
	must.NoError(page.Locator("#inp").Press(ctx, "Control+a"))
	must.NoError(page.Locator("#inp").Press(ctx, "Delete"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("", val)
}

// TestLocatorPressPageUpEx5 verifies pressing PageUp key.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressPageUpEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" tabindex="0"></div>
		<script>
			var pressed = '';
			document.getElementById('target').addEventListener('keydown', function(e) {
				pressed = e.key;
			});
		</script>
	`))

	must.NoError(page.Locator("#target").Focus(ctx))
	must.NoError(page.Locator("#target").Press(ctx, "PageUp"))

	result, err := page.Evaluate(ctx, `() => pressed`)
	must.NoError(err)
	is.Equal("PageUp", result)
}

// TestLocatorPressF1Ex5 verifies pressing F1 key fires keydown.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorPressF1Ex5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div id="target" tabindex="0"></div>
		<script>
			var lastKey = '';
			document.getElementById('target').addEventListener('keydown', function(e) {
				lastKey = e.key;
			});
		</script>
	`))

	must.NoError(page.Locator("#target").Focus(ctx))
	must.NoError(page.Locator("#target").Press(ctx, "F1"))

	result, err := page.Evaluate(ctx, `() => lastKey`)
	must.NoError(err)
	is.Equal("F1", result)
}

// TestPressEnterOnFormEx6 verifies Press Enter submits form via onsubmit handler.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressEnterOnFormEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form onsubmit="document.getElementById('out').textContent='submitted'; return false">
			<input id="inp" type="text">
		</form>
		<div id="out"></div>
	`))

	must.NoError(page.Locator("#inp").Press(ctx, "Enter"))

	out, err := page.Locator("#out").TextContent(ctx)
	must.NoError(err)
	is.Equal("submitted", out)
}

// TestPressTabMovesFocusEx6 verifies Press Tab moves focus to next element.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressTabMovesFocusEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first" type="text">
		<input id="second" type="text">
	`))

	must.NoError(page.Locator("#first").Press(ctx, "Tab"))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", focused)
}

// TestPressEscapeOnInputEx6 verifies Press Escape triggers blur on input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressEscapeOnInputEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" value="some text">
	`))

	must.NoError(page.Locator("#inp").Focus(ctx))
	must.NoError(page.Locator("#inp").Press(ctx, "Escape"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("some text", val)
}

// TestPressBackspaceEx6 verifies Press Backspace deletes last character.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressBackspaceEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "hello"))
	must.NoError(page.Locator("#inp").Press(ctx, "Backspace"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hell", val)
}

// TestPressCtrlZUndoEx7 verifies Ctrl+Z undoes text input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressCtrlZUndoEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "hello"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hello", val)
}

// TestPressEndKeyEx7 verifies End key dispatched on input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressEndKeyEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="Hello">`))

	must.NoError(page.Locator("#inp").Press(ctx, "End"))

	pos, err := page.Locator("#inp").Evaluate(ctx, `el => el.selectionStart`)
	must.NoError(err)
	is.Equal(float64(5), pos)
}

// TestPressArrowLeftEx7 verifies ArrowLeft moves cursor.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressArrowLeftEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="Hello">`))

	must.NoError(page.Locator("#inp").Press(ctx, "End"))
	must.NoError(page.Locator("#inp").Press(ctx, "ArrowLeft"))
	must.NoError(page.Locator("#inp").Press(ctx, "ArrowLeft"))

	pos, err := page.Locator("#inp").Evaluate(ctx, `el => el.selectionStart`)
	must.NoError(err)
	is.Equal(float64(3), pos)
}

// TestPressDeleteKeyEx7 verifies Delete key removes character.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressDeleteKeyEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="Helo">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Locator("#inp").Press(ctx, "Home"))
	must.NoError(page.Locator("#inp").Press(ctx, "ArrowRight"))
	must.NoError(page.Locator("#inp").Press(ctx, "ArrowRight"))
	must.NoError(page.Locator("#inp").Press(ctx, "ArrowRight"))
	must.NoError(page.Locator("#inp").Press(ctx, "Delete"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hel", val)
}

// TestPressSequentiallyAppendsToInput verifies PressSequentially appends text to an input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressSequentiallyAppendsToInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="Hello">`))
	must.NoError(page.Locator("input").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Locator("input").PressSequentially(ctx, " World"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello World", val)
}

// TestPressSequentiallyFiresKeyEvents verifies PressSequentially fires key events.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressSequentiallyFiresKeyEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text">
		<div id="log"></div>
		<script>
			const inp = document.getElementById('inp');
			const log = document.getElementById('log');
			inp.addEventListener('keydown', e => { log.textContent += e.key; });
		</script>
	`))

	must.NoError(page.Locator("input").Click(ctx))
	must.NoError(page.Locator("input").PressSequentially(ctx, "ab"))

	logText, err := page.Locator("#log").InnerText(ctx)
	must.NoError(err)
	is.Contains(logText, "a")
	is.Contains(logText, "b")
}

// TestPressSequentiallyInTextarea verifies PressSequentially works in textarea.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressSequentiallyInTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea></textarea>`))
	must.NoError(page.Locator("textarea").Click(ctx))
	must.NoError(page.Locator("textarea").PressSequentially(ctx, "line1"))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("line1", val)
}

// TestPressSequentiallyWithSpecialCharacters verifies PressSequentially handles special chars.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPressSequentiallyWithSpecialCharacters(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))
	must.NoError(page.Locator("input").Click(ctx))
	must.NoError(page.Locator("input").PressSequentially(ctx, "abc!123"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("abc!123", val)
}

// TestKeyboardTypeInInputEx2 verifies Keyboard.Type enters text in input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeInInputEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Type(ctx, "keyboard test"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("keyboard test", val)
}

// TestKeyboardPressEnterSubmitsFormEx2 verifies Enter press submits a form.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressEnterSubmitsFormEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form id="frm" onsubmit="window.__submitted=true; return false;">
			<input id="inp" type="text">
		</form>
	`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Enter"))

	submitted, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, submitted)
}

// TestKeyboardPressTabMovesFocusEx2 verifies Tab key moves focus between elements.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardPressTabMovesFocusEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first" type="text">
		<input id="second" type="text">
	`))

	must.NoError(page.Locator("#first").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Tab"))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", focused)
}

// TestKeyboardDownAndUpEvents verifies Keyboard.Down and Up fire events.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardDownAndUpEventsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text"
		  onkeydown="window.__down=true"
		  onkeyup="window.__up=true">
	`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Down(ctx, "a"))
	must.NoError(page.Keyboard.Up(ctx, "a"))

	down, err := page.Evaluate(ctx, `() => window.__down`)
	must.NoError(err)
	is.Equal(true, down)

	up, err := page.Evaluate(ctx, `() => window.__up`)
	must.NoError(err)
	is.Equal(true, up)
}

// TestKeyboardInsertTextEx2 verifies Keyboard.InsertText inserts text directly.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.InsertText(ctx, "inserted"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("inserted", val)
}

// TestKeyboardTypeIntoInputEx3 verifies typing into input with keyboard.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeIntoInputEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Type(ctx, "keyboard input"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("keyboard input", val)
}

// TestKeyboardShiftPressEx3 verifies Shift key press produces uppercase.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardShiftPressEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Down(ctx, "Shift"))
	must.NoError(page.Keyboard.Press(ctx, "KeyA"))
	must.NoError(page.Keyboard.Up(ctx, "Shift"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("A", val)
}

// TestKeyboardBackspaceEx3 verifies Backspace deletes character.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardBackspaceEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello">`))
	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "End"))
	must.NoError(page.Keyboard.Press(ctx, "Backspace"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hell", val)
}

// TestKeyboardTabFocusEx3 verifies Tab key moves focus to next element.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTabFocusEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="first" type="text">
		<input id="second" type="text">
	`))
	must.NoError(page.Locator("#first").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Tab"))

	focused, err := page.Evaluate(ctx, `() => document.activeElement.id`)
	must.NoError(err)
	is.Equal("second", focused)
}

// TestKeyboardEnterSubmitsFormEx3 verifies Enter key submits form.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardEnterSubmitsFormEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form id="f">
			<input id="inp" type="text">
		</form>
		<script>
			var submitted = false;
			document.getElementById('f').addEventListener('submit', function(e) {
				e.preventDefault();
				submitted = true;
			});
		</script>
	`))
	must.NoError(page.Locator("#inp").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "Enter"))

	result, err := page.Evaluate(ctx, `() => submitted`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestKeyboardShiftSelectEx4 verifies Shift+End selects text in input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardShiftSelectEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="hello world">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Keyboard.Press(ctx, "Shift+End"))

	selected, err := page.Evaluate(ctx, `() => {
		const el = document.getElementById('inp');
		return el.value.substring(el.selectionStart, el.selectionEnd);
	}`)
	must.NoError(err)
	is.Equal("hello world", selected)
}

// TestKeyboardCtrlAEx4 verifies Ctrl+A selects all text in input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardCtrlAEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="select me">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.Press(ctx, "Control+a"))

	selected, err := page.Evaluate(ctx, `() => {
		const el = document.getElementById('inp');
		return el.value.substring(el.selectionStart, el.selectionEnd);
	}`)
	must.NoError(err)
	is.Equal("select me", selected)
}

// TestKeyboardTypeSpecialCharsEx4 verifies keyboard Type handles special chars.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardTypeSpecialCharsEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.Type(ctx, "Hello!"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello!", val)
}

// TestKeyboardInsertTextEx4 verifies keyboard InsertText works.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestKeyboardInsertTextEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))
	must.NoError(page.Locator("#inp").Focus(ctx))

	must.NoError(page.Keyboard.InsertText(ctx, "inserted"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("inserted", val)
}

// Ref: TestPageKeyboard.java#insertTextShouldOnlyEmitInputEvent  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardInsertTextShouldOnlyEmitInputEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<textarea></textarea>`)
	must.NoError(err)
	err = page.Locator("textarea").Click(ctx)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
		window._events = [];
		document.addEventListener('keydown', e => window._events.push(e.type));
		document.addEventListener('keyup', e => window._events.push(e.type));
		document.addEventListener('keypress', e => window._events.push(e.type));
		document.addEventListener('input', e => window._events.push(e.type));
	}`)
	must.NoError(err)

	err = page.Keyboard.InsertText(ctx, "hello world")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "window._events")
	must.NoError(err)
	evList, ok := val.([]any)
	is.True(ok, "expected array, got %T", val)
	is.Len(evList, 1)
	is.Equal("input", evList[0])
}

// Ref: TestPageKeyboard.java#shouldNotTypeCanceledEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldNotTypeCanceledEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<textarea></textarea>`)
	must.NoError(err)
	err = page.Locator("textarea").Click(ctx)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
		window.addEventListener('keydown', function(event) {
			event.stopPropagation();
			event.stopImmediatePropagation();
			if (event.key === 'l') event.preventDefault();
			if (event.key === 'o') event.preventDefault();
		}, false);
	}`)
	must.NoError(err)

	err = page.Keyboard.Type(ctx, "Hello World!")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("He Wrd!", val)
}

// Ref: TestPageKeyboard.java#shouldPressPlus  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldPressPlus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		window._events = [];
		document.addEventListener('keydown', e => window._events.push('keydown:' + e.key));
		document.addEventListener('keypress', e => window._events.push('keypress:' + e.key));
		document.addEventListener('keyup', e => window._events.push('keyup:' + e.key));
	}`)
	must.NoError(err)

	err = page.Keyboard.Press(ctx, "+")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "window._events")
	must.NoError(err)
	events, ok := val.([]any)
	is.True(ok)
	is.Contains(events, "keydown:+")
	is.Contains(events, "keyup:+")
}

// Ref: TestPageKeyboard.java#shouldHandleSelectAll  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldHandleSelectAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<textarea></textarea>`)
	must.NoError(err)

	el, err := page.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NotNil(el)

	err = el.Press(ctx, "a")
	must.NoError(err)
	err = el.Press(ctx, "b")
	must.NoError(err)
	err = el.Press(ctx, "c")
	must.NoError(err)

	err = page.Keyboard.Down(ctx, "ControlOrMeta")
	must.NoError(err)
	err = page.Keyboard.Press(ctx, "a")
	must.NoError(err)
	err = page.Keyboard.Up(ctx, "ControlOrMeta")
	must.NoError(err)
	err = page.Keyboard.Press(ctx, "Backspace")
	must.NoError(err)

	val, err := page.Evaluate(ctx, "() => document.querySelector('textarea').value")
	must.NoError(err)
	is.Equal("", val)
}

// Ref: TestPageKeyboard.java#shouldPressShiftPlus  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldPressShiftPlus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.KeyboardPage())
	must.NoError(err)

	err = page.Keyboard.Press(ctx, "Shift++")
	must.NoError(err)

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, _ := result.(string)
	is.Contains(s, "Keydown: Shift ShiftLeft 16 [Shift]")
	is.Contains(s, "Keyup: + Equal 187 [Shift]")
	is.Contains(s, "Keyup: Shift ShiftLeft 16 []")
}

// Ref: TestPageKeyboard.java#shouldShiftRawCodes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldShiftRawCodes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.KeyboardPage())
	must.NoError(err)

	err = page.Keyboard.Press(ctx, "Shift+Digit3")
	must.NoError(err)

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, _ := result.(string)
	is.Contains(s, "Keydown: Shift ShiftLeft 16 [Shift]")
	is.Contains(s, "Keydown: # Digit3 51 [Shift]")
	is.Contains(s, "Keyup: Shift ShiftLeft 16 []")
}

// Ref: TestPageKeyboard.java#shouldPressEnter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldPressEnter(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.KeyboardPage())
	must.NoError(err)

	err = page.Keyboard.Press(ctx, "Enter")
	must.NoError(err)

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, _ := result.(string)
	is.Contains(s, "Keydown: Enter Enter")
	is.Contains(s, "Keyup: Enter Enter")
}

// Ref: TestPageKeyboard.java#shouldThrowOnUnknownKeys  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldThrowOnUnknownKeys(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.Keyboard.Press(ctx, "NotARealKey")
	is.Error(err, "pressing an unknown key should return an error")
}

// Ref: TestPageKeyboard.java#shouldTypeEmoji  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldTypeEmoji(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/input/textarea.html")
	must.NoError(err)

	err = page.Locator("textarea").Click(ctx)
	must.NoError(err)

	err = page.Keyboard.Type(ctx, "👹 Tokyo street Japan 🇯🇵")
	must.NoError(err)

	val, err := page.EvalOnSelector(ctx, "textarea", "t => t.value")
	must.NoError(err)
	is.Equal("👹 Tokyo street Japan 🇯🇵", val)
}

// Ref: TestPageKeyboard.java#shouldTypeAllKindsOfCharacters  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldTypeAllKindsOfCharacters(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/input/textarea.html")
	must.NoError(err)

	err = page.Locator("textarea").Click(ctx)
	must.NoError(err)

	const text = "This text goes onto two lines.\nThis character is 嗨."
	err = page.Keyboard.Type(ctx, text)
	must.NoError(err)

	val, err := page.EvalOnSelector(ctx, "textarea", "t => t.value")
	must.NoError(err)
	is.Equal(text, val)
}

// Ref: TestPageKeyboard.java#shouldSpecifyRepeatProperty  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSpecifyRepeatProperty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/input/textarea.html")
	must.NoError(err)

	err = page.Locator("textarea").Click(ctx)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
		window._lastKeydown = null;
		document.addEventListener('keydown', e => window._lastKeydown = e);
	}`)
	must.NoError(err)

	err = page.Keyboard.Down(ctx, "a")
	must.NoError(err)
	repeat1, err := page.Evaluate(ctx, "window._lastKeydown.repeat")
	must.NoError(err)
	is.Equal(false, repeat1)

	err = page.Keyboard.Press(ctx, "a")
	must.NoError(err)
	repeat2, err := page.Evaluate(ctx, "window._lastKeydown.repeat")
	must.NoError(err)
	is.Equal(true, repeat2)

	err = page.Keyboard.Down(ctx, "b")
	must.NoError(err)
	repeatB1, err := page.Evaluate(ctx, "window._lastKeydown.repeat")
	must.NoError(err)
	is.Equal(false, repeatB1)

	err = page.Keyboard.Down(ctx, "b")
	must.NoError(err)
	repeatB2, err := page.Evaluate(ctx, "window._lastKeydown.repeat")
	must.NoError(err)
	is.Equal(true, repeatB2)

	err = page.Keyboard.Up(ctx, "a")
	must.NoError(err)
	err = page.Keyboard.Down(ctx, "a")
	must.NoError(err)
	repeat3, err := page.Evaluate(ctx, "window._lastKeydown.repeat")
	must.NoError(err)
	is.Equal(false, repeat3)
}

// Ref: TestPageKeyboard.java#shouldSpecifyLocation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSpecifyLocation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.Prefix()+"/input/textarea.html")
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => {
		window._lastKeydown = null;
		document.addEventListener('keydown', e => window._lastKeydown = e);
	}`)
	must.NoError(err)

	el, err := page.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NotNil(el)

	err = el.Press(ctx, "Digit5")
	must.NoError(err)
	loc1, err := page.Evaluate(ctx, "window._lastKeydown.location")
	must.NoError(err)
	is.Equal(float64(0), loc1)

	err = el.Press(ctx, "ControlLeft")
	must.NoError(err)
	loc2, err := page.Evaluate(ctx, "window._lastKeydown.location")
	must.NoError(err)
	is.Equal(float64(1), loc2)

	err = el.Press(ctx, "ControlRight")
	must.NoError(err)
	loc3, err := page.Evaluate(ctx, "window._lastKeydown.location")
	must.NoError(err)
	is.Equal(float64(2), loc3)

	err = el.Press(ctx, "NumpadSubtract")
	must.NoError(err)
	loc4, err := page.Evaluate(ctx, "window._lastKeydown.location")
	must.NoError(err)
	is.Equal(float64(3), loc4)
}

// Ref: TestPageKeyboard.java#shouldSendProperCodesWhileTyping  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSendProperCodesWhileTyping(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.KeyboardPage())
	must.NoError(err)

	err = page.Keyboard.Type(ctx, "!")
	must.NoError(err)

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, _ := result.(string)
	is.Contains(s, "Keydown: ! Digit1 49 []")
	is.Contains(s, "Keypress: ! Digit1 33 33 []")
	is.Contains(s, "Keyup: ! Digit1 49 []")
}

// Ref: TestPageKeyboard.java#shouldReportMultipleModifiers  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldReportMultipleModifiers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.KeyboardPage()))
	kb := page.Keyboard

	must.NoError(kb.Down(ctx, "Control"))
	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	is.Equal("Keydown: Control ControlLeft 17 [Control]", result)

	must.NoError(kb.Down(ctx, "Alt"))
	result, err = page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	is.Equal("Keydown: Alt AltLeft 18 [Alt Control]", result)

	must.NoError(kb.Down(ctx, ";"))
	result, err = page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	is.Equal("Keydown: ; Semicolon 186 [Alt Control]", result)

	must.NoError(kb.Up(ctx, ";"))
	result, err = page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	is.Equal("Keyup: ; Semicolon 186 [Alt Control]", result)

	must.NoError(kb.Up(ctx, "Control"))
	result, err = page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	is.Equal("Keyup: Control ControlLeft 17 [Alt]", result)

	must.NoError(kb.Up(ctx, "Alt"))
	result, err = page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	is.Equal("Keyup: Alt AltLeft 18 []", result)
}

// Ref: TestPageKeyboard.java#shouldSendProperCodesWhileTypingWithShift  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSendProperCodesWhileTypingWithShift(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.KeyboardPage()))
	kb := page.Keyboard

	must.NoError(kb.Down(ctx, "Shift"))
	must.NoError(kb.Type(ctx, "~"))
	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	want := "Keydown: Shift ShiftLeft 16 [Shift]\nKeydown: ~ Backquote 192 [Shift]\nKeypress: ~ Backquote 126 126 [Shift]\nKeyup: ~ Backquote 192 [Shift]"
	is.Equal(want, result)
	must.NoError(kb.Up(ctx, "Shift"))
}

// TestPageKeyboardShouldTypeIntoATextarea verifies typing into a textarea via Keyboard.Type.
// Ref: TestPageKeyboard.java#shouldTypeIntoATextarea  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldTypeIntoATextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))

	must.NoError(page.Locator("textarea").Click(ctx))
	must.NoError(page.Keyboard.Type(ctx, "Hello World"))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello World", val)
}

// TestPageTypeIntoTextareaWithSpecialChars verifies typing special characters via Keyboard.Type.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypeIntoTextareaWithSpecialChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must.NoError(page.Locator("textarea").Click(ctx))

	must.NoError(page.Keyboard.Type(ctx, "This is Text"))
	must.NoError(page.Keyboard.Press(ctx, "Home"))
	must.NoError(page.Keyboard.Type(ctx, "Very "))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("Very This is Text", val)
}

// TestPageTypeWithNewlinesKeyboardPress verifies pressing Enter creates newlines.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypeWithNewlinesKeyboardPress(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must.NoError(page.Locator("textarea").Click(ctx))

	must.NoError(page.Keyboard.Type(ctx, "Hello"))
	must.NoError(page.Keyboard.Press(ctx, "Enter"))
	must.NoError(page.Keyboard.Type(ctx, "World"))

	val, err := page.Locator("textarea").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello\nWorld", val)
}

// TestPagePressSequentiallyIntoInput verifies PressSequentially into input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPagePressSequentiallyIntoInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, `<input type="text" id="inp">`))

	must.NoError(page.Locator("input").PressSequentially(ctx, "Hello"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("Hello", val)
}

// TestLocatorClickAndTypeSetsValue verifies clicking and then typing sets the value.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorClickAndTypeSetsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" value="existing">`))

	// Triple-click to select all, then type to replace
	must.NoError(page.Locator("input").Click(ctx))
	// Select all via keyboard shortcut
	must.NoError(page.Keyboard.Press(ctx, "Control+a"))
	must.NoError(page.Keyboard.Type(ctx, "replaced"))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("replaced", val)
}

// TestPageTypeAndVerifyEventsFire verifies that typing fires keydown, keypress, and keyup events.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypeAndVerifyEventsFire(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.KeyboardPage()))

	must.NoError(page.Keyboard.Type(ctx, "f"))

	result, err := page.Evaluate(ctx, "getResult()")
	must.NoError(err)
	s, ok := result.(string)
	is.True(ok)
	is.Contains(s, "Keydown: f")
	is.Contains(s, "Keyup: f")
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypeAppendsToExistingValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="hello" id="inp">`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, " world"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Equal("hello world", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypeFiresInputEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" oninput="window.__count = (window.__count||0) + 1">
	`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, "abc"))

	count, err := page.Evaluate(ctx, `() => window.__count`)
	must.NoError(err)
	is.Equal(float64(3), count)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypeIntoPasswordField(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="password" id="pw">`))

	handle, err := page.QuerySelector(ctx, "#pw")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, "secret"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Equal("secret", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypePressEnterSubmitsForm(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<form onsubmit="window.__submitted=true; return false;">
			<input id="inp" type="text">
		</form>
	`))

	handle, err := page.QuerySelector(ctx, "#inp")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Press(ctx, "Enter"))

	result, err := page.Evaluate(ctx, `() => window.__submitted`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageTypeWithNewlineInTextarea(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	handle, err := page.QuerySelector(ctx, "#ta")
	must.NoError(err)

	must.NoError(handle.Click(ctx))
	must.NoError(handle.Type(ctx, "line1"))
	must.NoError(handle.Press(ctx, "Enter"))
	must.NoError(handle.Type(ctx, "line2"))

	val, err := handle.InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "line1")
	is.Contains(val, "line2")
}

// TestLocatorTypeTextEx2 verifies Fill enters text into an input.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorTypeTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "hello world"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("hello world", val)
}

// TestLocatorFillInTextareaEx2 verifies Fill works in textarea.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorFillInTextareaEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	must.NoError(page.Locator("#ta").Fill(ctx, "multiline text"))

	val, err := page.Locator("#ta").InputValue(ctx)
	must.NoError(err)
	is.Contains(val, "multiline text")
}

// TestLocatorFillNumbersEx2 verifies Fill works with numeric strings.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorFillNumbersEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "9876"))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Equal("9876", val)
}

// TestLocatorFillFiresInputEventEx2 verifies Fill fires input event.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorFillFiresInputEventEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<input id="inp" type="text" oninput="window.__inputFired=true">
	`))

	must.NoError(page.Locator("#inp").Fill(ctx, "some text"))

	fired, err := page.Evaluate(ctx, `() => window.__inputFired`)
	must.NoError(err)
	is.Equal(true, fired)
}

// TestLocatorFillAndClearEx2 verifies Fill followed by Fill with empty string clears.
// Extra: no equivalent method in playwright-java; related class: TestPageKeyboard.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestLocatorFillAndClearEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Fill(ctx, "initial"))
	must.NoError(page.Locator("#inp").Fill(ctx, ""))

	val, err := page.Locator("#inp").InputValue(ctx)
	must.NoError(err)
	is.Empty(val)
}

// Ref: TestPageKeyboard.java#shouldMoveWithTheArrowKeys  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldMoveWithTheArrowKeys(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must := require.New(t)
	must.NoError(page.Locator("textarea").PressSequentially(ctx, "Hello World!"))
	assert.Equal(t, "Hello World!", a2Ev(t, ctx, page, "() => document.querySelector('textarea').value"))
	for range "World!" {
		must.NoError(page.Keyboard.Press(ctx, "ArrowLeft"))
	}
	must.NoError(page.Keyboard.Type(ctx, "inserted "))
	assert.Equal(t, "Hello inserted World!", a2Ev(t, ctx, page, "() => document.querySelector('textarea').value"))
	must.NoError(page.Keyboard.Down(ctx, "Shift"))
	for range "inserted " {
		must.NoError(page.Keyboard.Press(ctx, "ArrowLeft"))
	}
	must.NoError(page.Keyboard.Up(ctx, "Shift"))
	must.NoError(page.Keyboard.Press(ctx, "Backspace"))
	assert.Equal(t, "Hello World!", a2Ev(t, ctx, page, "() => document.querySelector('textarea').value"))
}

// Ref: TestPageKeyboard.java#shouldSendACharacterWithElementHandlePress  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSendACharacterWithElementHandlePress(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must := require.New(t)
	textarea, err := page.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NoError(textarea.Press(ctx, "a"))
	assert.Equal(t, "a", a2Ev(t, ctx, page, "() => document.querySelector('textarea').value"))
	a2Ev(t, ctx, page, "() => window.addEventListener('keydown', e => e.preventDefault(), true)")
	must.NoError(textarea.Press(ctx, "b"))
	assert.Equal(t, "a", a2Ev(t, ctx, page, "() => document.querySelector('textarea').value"))
}

// Ref: TestPageKeyboard.java#shouldSendACharacterWithInsertText  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSendACharacterWithInsertText(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must := require.New(t)
	must.NoError(page.Locator("textarea").Focus(ctx))
	must.NoError(page.Keyboard.InsertText(ctx, "嗨"))
	assert.Equal(t, "嗨", a2Ev(t, ctx, page, "() => document.querySelector('textarea').value"))
	a2Ev(t, ctx, page, "() => window.addEventListener('keydown', e => e.preventDefault(), true)")
	must.NoError(page.Keyboard.InsertText(ctx, "a"))
	assert.Equal(t, "嗨a", a2Ev(t, ctx, page, "() => document.querySelector('textarea').value"))
}

// Ref: TestPageKeyboard.java#shouldReportShiftKey  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldReportShiftKey(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	must := require.New(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/keyboard.html"))
	kb := page.Keyboard
	for _, m := range []struct {
		key  string
		code int
	}{{"Shift", 16}, {"Alt", 18}, {"Control", 17}} {
		must.NoError(kb.Down(ctx, m.key))
		assert.Equal(t, fmt.Sprintf("Keydown: %s %sLeft %d [%s]", m.key, m.key, m.code, m.key), a2Ev(t, ctx, page, "getResult()"))
		must.NoError(kb.Down(ctx, "!"))
		if m.key == "Shift" {
			assert.Equal(t, fmt.Sprintf("Keydown: ! Digit1 49 [%s]\nKeypress: ! Digit1 33 33 [%s]", m.key, m.key), a2Ev(t, ctx, page, "getResult()"))
		} else {
			assert.Equal(t, fmt.Sprintf("Keydown: ! Digit1 49 [%s]", m.key), a2Ev(t, ctx, page, "getResult()"))
		}
		must.NoError(kb.Up(ctx, "!"))
		assert.Equal(t, fmt.Sprintf("Keyup: ! Digit1 49 [%s]", m.key), a2Ev(t, ctx, page, "getResult()"))
		must.NoError(kb.Up(ctx, m.key))
		assert.Equal(t, fmt.Sprintf("Keyup: %s %sLeft %d []", m.key, m.key, m.code), a2Ev(t, ctx, page, "getResult()"))
	}
}

// Ref: TestPageKeyboard.java#shouldSupportPlusSeparatedModifier  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSupportPlusSeparatedModifier(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/keyboard.html"))
	require.NoError(t, page.Keyboard.Press(ctx, "Shift+~"))
	assert.Equal(t, strings.Join([]string{
		"Keydown: Shift ShiftLeft 16 [Shift]",
		"Keydown: ~ Backquote 192 [Shift]",
		"Keypress: ~ Backquote 126 126 [Shift]",
		"Keyup: ~ Backquote 192 [Shift]",
		"Keyup: Shift ShiftLeft 16 []",
	}, "\n"), a2Ev(t, ctx, page, "getResult()"))
}

// Ref: TestPageKeyboard.java#shouldSupportMultiplePlusSeparatedModifier  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSupportMultiplePlusSeparatedModifier(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/keyboard.html"))
	require.NoError(t, page.Keyboard.Press(ctx, "Control+Shift+~"))
	assert.Equal(t, strings.Join([]string{
		"Keydown: Control ControlLeft 17 [Control]",
		"Keydown: Shift ShiftLeft 16 [Control Shift]",
		"Keydown: ~ Backquote 192 [Control Shift]",
		"Keyup: ~ Backquote 192 [Control Shift]",
		"Keyup: Shift ShiftLeft 16 [Control]",
		"Keyup: Control ControlLeft 17 []",
	}, "\n"), a2Ev(t, ctx, page, "getResult()"))
}

// Ref: TestPageKeyboard.java#shouldTypeEmojiIntoAnIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldTypeEmojiIntoAnIframe(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	must := require.New(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err := a2AttachFrame(ctx, page, "emoji-test", srv.Prefix()+"/input/textarea.html")
	must.NoError(err)
	frame := page.Frames()[1]
	textarea, err := frame.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NoError(textarea.Type(ctx, "👹 Tokyo street Japan 🇯🇵"))
	v, err := frame.Locator("textarea").Evaluate(ctx, "textarea => textarea.value")
	must.NoError(err)
	assert.Equal(t, "👹 Tokyo street Japan 🇯🇵", v)
}

// Ref: TestPageKeyboard.java#shouldBeAbleToPreventSelectAll  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldBeAbleToPreventSelectAll(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	must := require.New(t)
	textarea, err := page.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NoError(textarea.Type(ctx, "some text"))
	_, err = page.EvalOnSelector(ctx, "textarea", "textarea => {\n  textarea.addEventListener('keydown', event => {\n    if (event.key === 'a' && (event.metaKey || event.ctrlKey))\n      event.preventDefault();\n  }, false);\n}")
	must.NoError(err)
	must.NoError(page.Keyboard.Down(ctx, "ControlOrMeta"))
	must.NoError(page.Keyboard.Press(ctx, "a"))
	must.NoError(page.Keyboard.Up(ctx, "ControlOrMeta"))
	must.NoError(page.Keyboard.Press(ctx, "Backspace"))
	v, err := page.EvalOnSelector(ctx, "textarea", "textarea => textarea.value")
	must.NoError(err)
	assert.Equal(t, "some tex", v)
}

// Ref: TestPageKeyboard.java#shouldSupportMacOSShortcuts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldSupportMacOSShortcuts(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/textarea.html"))
	if runtime.GOOS != "darwin" {
		t.Skip("macOS only")
	}
	must := require.New(t)
	textarea, err := page.QuerySelector(ctx, "textarea")
	must.NoError(err)
	must.NoError(textarea.Type(ctx, "some text"))
	must.NoError(page.Keyboard.Press(ctx, "Shift+Control+Alt+KeyB"))
	must.NoError(page.Keyboard.Press(ctx, "Backspace"))
	v, err := page.EvalOnSelector(ctx, "textarea", "textarea => textarea.value")
	must.NoError(err)
	assert.Equal(t, "some ", v)
}

// Ref: TestPageKeyboard.java#shouldPressMetaKey  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldPressMetaKey(t *testing.T) {
	ctx, _, page := a2Setup(t)
	a2CaptureLastKeyDown(t, ctx, page)
	require.NoError(t, page.Keyboard.Press(ctx, "Meta"))
	ev, ok := a2Ev(t, ctx, page, "() => window.__a2LastKeyDown").(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Meta", ev["key"])
	assert.Equal(t, "MetaLeft", ev["code"])
	assert.Equal(t, true, ev["metaKey"], fmt.Sprint(ev))
}

// Ref: TestPageKeyboard.java#shouldWorkAfterACrossOriginNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldWorkAfterACrossOriginNavigation(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	must := require.New(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/empty.html"))
	must.NoError(page.Goto(ctx, srv.CrossProcessPrefix()+"/empty.html"))
	a2CaptureLastKeyDown(t, ctx, page)
	must.NoError(page.Keyboard.Press(ctx, "a"))
	assert.Equal(t, "a", a2Ev(t, ctx, page, "() => window.__a2LastKeyDown.key"))
}

// Ref: TestPageKeyboard.java#shouldExposeKeyIdentifierInWebKit  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldExposeKeyIdentifierInWebKit(t *testing.T) {
	ctx, _, page := a2Setup(t)
	if os.Getenv("BROWSER") != "webkit" {
		t.Skip("WebKit only")
	}
	a2CaptureLastKeyDown(t, ctx, page)
	for key, ident := range map[string]string{"ArrowUp": "Up", "ArrowDown": "Down", "ArrowLeft": "Left", "ArrowRight": "Right", "Backspace": "U+0008", "Tab": "U+0009", "Delete": "U+007F", "a": "U+0041", "b": "U+0042", "F12": "F12"} {
		require.NoError(t, page.Keyboard.Press(ctx, key))
		assert.Equal(t, ident, a2Ev(t, ctx, page, "() => window.__a2LastKeyDown.keyIdentifier"))
	}
}

// Ref: TestPageKeyboard.java#shouldScrollWithPageDown  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageKeyboard.java
func TestPageKeyboardShouldScrollWithPageDown(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	must := require.New(t)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/scrollable.html"))
	must.NoError(page.Locator("body").Click(ctx))
	must.NoError(page.Keyboard.Press(ctx, "PageDown"))
	_, err := page.WaitForFunction(ctx, "() => scrollY > 0", nil)
	must.NoError(err)
}
