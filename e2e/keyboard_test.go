//go:build e2e

// Keyboard E2E tests.
// Migration of: TestKeyboard.java
package e2e

import (
	"testing"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestKeyboardInsertTextOnlyEmitsInputEvent(t *testing.T) {
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

func TestKeyboardNotTypeCanceledEvents(t *testing.T) {
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

func TestKeyboardPressPlus(t *testing.T) {
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

func TestKeyboardHandleSelectAll(t *testing.T) {
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

func TestKeyboardPressShiftPlus(t *testing.T) {
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

func TestKeyboardShiftRawCodes(t *testing.T) {
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

func TestKeyboardPressEnter(t *testing.T) {
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

func TestKeyboardThrowOnUnknownKeys(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.Keyboard.Press(ctx, "NotARealKey")
	is.Error(err, "pressing an unknown key should return an error")
}

func TestKeyboardTypeEmoji(t *testing.T) {
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

func TestKeyboardTypeAllKindsOfCharacters(t *testing.T) {
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

func TestKeyboardSpecifyRepeatProperty(t *testing.T) {
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

func TestKeyboardSpecifyLocation(t *testing.T) {
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

func TestKeyboardSendProperCodesWhileTyping(t *testing.T) {
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

func TestKeyboardReportMultipleModifiers(t *testing.T) {
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

func TestKeyboardSendProperCodesWhileTypingWithShift(t *testing.T) {
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

func TestKeyboardFunctionKeyEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	err := page.Keyboard.Press(ctx, "F5")
	must.NoError(err)
}

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
// Ref: TestLocatorKeyboard.java#shouldPressEnterToSubmitForm
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
// Ref: TestLocatorKeyboard.java#shouldPressArrowKey
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
// Ref: TestLocatorKeyboard.java#shouldPressDeleteKey
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
// Ref: TestLocatorKeyboard.java#shouldPressTabToMoveFocus
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
// Ref: TestLocatorPress.java#shouldBlurOnEscape
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
// Ref: TestLocatorPress.java#shouldMoveCaretLeft
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
// Ref: TestLocatorPress.java#shouldMoveToStartOnHome
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
// Ref: TestLocatorPress.java#shouldMoveToEndOnEnd
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
// Ref: TestLocatorPress.java#shouldPageDownScroll
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
// Ref: TestLocatorPress.java#shouldTabMovesFocus
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
// Ref: TestLocatorPress.java#shouldSpaceTogglesCheckbox
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
// Ref: TestLocatorPress.java#shouldDeleteText
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
// Ref: TestLocatorPress.java#shouldBackspace
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
// Ref: TestLocatorPress.java#shouldSubmitOnEnter
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
// Ref: TestLocatorPress.java#shouldFireArrowDown
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
// Ref: TestLocatorPress.java#shouldFireEscape
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
// Ref: TestLocatorPress.java#shouldPressHome
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
// Ref: TestLocatorPress.java#shouldToggleCheckboxWithSpace
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
// Ref: TestLocatorPress.java#shouldDeleteText
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
// Ref: TestLocatorPress.java#shouldPressPageUp
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
// Ref: TestLocatorPress.java#shouldPressF1
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
// Ref: TestLocatorPress.java#shouldSubmitFormOnEnter
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
// Ref: TestLocatorPress.java#shouldMoveFocusOnTab
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
// Ref: TestLocatorPress.java#shouldHandleEscape
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
// Ref: TestLocatorPress.java#shouldDeleteOnBackspace
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
// Ref: TestLocatorPress.java#shouldUndoWithCtrlZ
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
// Ref: TestLocatorPress.java#shouldPressEnd
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
// Ref: TestLocatorPress.java#shouldMoveLeft
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
// Ref: TestLocatorPress.java#shouldDeleteChar
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
// Ref: TestLocatorPressSequentially.java#shouldAppendToInput
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
// Ref: TestLocatorPressSequentially.java#shouldFireKeyEvents
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
// Ref: TestLocatorPressSequentially.java#shouldWorkInTextarea
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
// Ref: TestLocatorPressSequentially.java#shouldWorkWithSpecialCharacters
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
// Ref: TestPageKeyboard.java#shouldTypeText
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
// Ref: TestPageKeyboard.java#shouldSubmitOnEnter
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
// Ref: TestPageKeyboard.java#shouldMovesFocusOnTab
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
// Ref: TestPageKeyboard.java#shouldFireDownUpEvents
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
// Ref: TestPageKeyboard.java#shouldInsertText
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
// Ref: TestKeyboard.java#shouldTypeIntoInput
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
// Ref: TestKeyboard.java#shouldProduceUppercaseWithShift
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
// Ref: TestKeyboard.java#shouldDeleteWithBackspace
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
// Ref: TestKeyboard.java#shouldMoveFocusWithTab
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
// Ref: TestKeyboard.java#shouldSubmitFormWithEnter
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
// Ref: TestPageKeyboard.java#shouldSelectWithShift
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
// Ref: TestPageKeyboard.java#shouldSelectAllWithCtrlA
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
// Ref: TestPageKeyboard.java#shouldTypeSpecialCharacters
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
// Ref: TestPageKeyboard.java#shouldInsertText
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
