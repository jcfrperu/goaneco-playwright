//go:build e2e

// E2E tests for: TestPageDialog.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageOnDialogAlert(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	dialogSeen := make(chan *playwright.Dialog, 1)
	page.OnDialog(func(d *playwright.Dialog) {
		dialogSeen <- d
		_ = d.Accept(ctx)
	})

	_, err = page.Evaluate(ctx, "alert('hello from alert')")
	must.NoError(err, "Evaluate alert failed")

	select {
	case d := <-dialogSeen:
		if d.Type() != "alert" {
			t.Errorf("dialog type = %q, want 'alert'", d.Type())
		}
		if d.Message() != "hello from alert" {
			t.Errorf("dialog message = %q, want 'hello from alert'", d.Message())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("dialog handler was never called")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageOnDialogConfirmDismiss(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Dismiss(ctx)
	})

	result, err := page.Evaluate(ctx, "confirm('proceed?')")
	must.NoError(err, "Evaluate confirm failed")
	if result != false {
		t.Errorf("confirm result = %v, want false (dismissed)", result)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageOnDialogPrompt(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Accept(ctx, "Juan")
	})

	result, err := page.Evaluate(ctx, "prompt('enter name')")
	must.NoError(err, "Evaluate prompt failed")
	if result != "Juan" {
		t.Errorf("prompt result = %q, want 'Juan'", result)
	}
}

// TestPageOnDialogCancel verifies that the cancel function returned by OnDialog
// stops the handler from receiving subsequent dialog events.
// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageOnDialogCancel(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Register first handler and wait for it to fire once, then cancel.
	firstSeen := make(chan struct{})
	var firstOnce sync.Once
	cancel := page.OnDialog(func(d *playwright.Dialog) {
		firstOnce.Do(func() { close(firstSeen) })
		_ = d.Accept(ctx)
	})

	_, err = page.Evaluate(ctx, "alert('first')")
	must.NoError(err, "Evaluate first alert failed")
	select {
	case <-firstSeen:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for first dialog")
	}

	cancel()

	// Register a second handler; only it should fire for the next dialog.
	var count int32
	done := make(chan struct{})
	var doneOnce sync.Once
	page.OnDialog(func(d *playwright.Dialog) {
		atomic.AddInt32(&count, 1)
		doneOnce.Do(func() { close(done) })
		_ = d.Accept(ctx)
	})

	_, err = page.Evaluate(ctx, "alert('second')")
	must.NoError(err, "Evaluate second alert failed")
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for second dialog")
	}

	time.Sleep(100 * time.Millisecond)

	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("expected second handler called once, got %d (canceled handler may have fired too)", got)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestDialogTypeIsAlert(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	typeCh := make(chan string, 1)
	page.OnDialog(func(d *playwright.Dialog) {
		typeCh <- d.Type()
		_ = d.Accept(ctx)
	})

	_, err = page.Evaluate(ctx, "alert('hi')")
	must.NoError(err, "Evaluate failed")

	select {
	case dt := <-typeCh:
		if dt != "alert" {
			t.Errorf("Dialog.Type() = %q, want 'alert'", dt)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for dialog")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestDialogTypeIsConfirm(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	typeCh := make(chan string, 1)
	page.OnDialog(func(d *playwright.Dialog) {
		typeCh <- d.Type()
		_ = d.Dismiss(ctx)
	})

	_, err = page.Evaluate(ctx, "confirm('are you sure?')")
	must.NoError(err, "Evaluate failed")

	select {
	case dt := <-typeCh:
		if dt != "confirm" {
			t.Errorf("Dialog.Type() = %q, want 'confirm'", dt)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for dialog")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestDialogTypeIsPrompt(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	typeCh := make(chan string, 1)
	page.OnDialog(func(d *playwright.Dialog) {
		typeCh <- d.Type()
		_ = d.Accept(ctx)
	})

	_, err = page.Evaluate(ctx, "prompt('enter value', 'default')")
	must.NoError(err, "Evaluate failed")

	select {
	case dt := <-typeCh:
		if dt != "prompt" {
			t.Errorf("Dialog.Type() = %q, want 'prompt'", dt)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for dialog")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestDialogDefaultValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	defaultCh := make(chan string, 1)
	page.OnDialog(func(d *playwright.Dialog) {
		defaultCh <- d.DefaultValue()
		_ = d.Accept(ctx)
	})

	_, err = page.Evaluate(ctx, "prompt('enter value', 'my-default')")
	must.NoError(err, "Evaluate failed")

	select {
	case dv := <-defaultCh:
		if dv != "my-default" {
			t.Errorf("Dialog.DefaultValue() = %q, want 'my-default'", dv)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for dialog")
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestDialogDefaultValueEmptyForAlert(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	defaultCh := make(chan string, 1)
	page.OnDialog(func(d *playwright.Dialog) {
		defaultCh <- d.DefaultValue()
		_ = d.Accept(ctx)
	})

	_, err = page.Evaluate(ctx, "alert('no default here')")
	must.NoError(err, "Evaluate failed")

	select {
	case dv := <-defaultCh:
		if dv != "" {
			t.Errorf("Dialog.DefaultValue() for alert = %q, want ''", dv)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for dialog")
	}
}

// TestOnDialogAlertAccepted verifies OnDialog can accept an alert.
// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestOnDialogAlertAccepted(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var dialogType string

	off := page.OnDialog(func(d *playwright.Dialog) {
		mu.Lock()
		dialogType = d.Type()
		mu.Unlock()
		_ = d.Accept(ctx)
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	_, _ = page.Evaluate(ctx, `() => alert('hello')`)

	mu.Lock()
	dt := dialogType
	mu.Unlock()

	is.Equal("alert", dt)
}

// TestOnDialogConfirmAccepted verifies OnDialog can accept a confirm.
// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestOnDialogConfirmAccepted(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	off := page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Accept(ctx)
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	result, err := page.Evaluate(ctx, `() => confirm('proceed?')`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestOnDialogConfirmDismissed verifies OnDialog can dismiss a confirm.
// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestOnDialogConfirmDismissed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	off := page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Dismiss(ctx)
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	result, err := page.Evaluate(ctx, `() => confirm('proceed?')`)
	must.NoError(err)
	is.Equal(false, result)
}

// TestOnDialogPromptWithValue verifies OnDialog accepts prompt with value.
// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestOnDialogPromptAcceptedWithValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	off := page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Accept(ctx, "typed value")
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	result, err := page.Evaluate(ctx, `() => prompt('enter value')`)
	must.NoError(err)
	is.Equal("typed value", result)
}

// TestOnDialogMessageText verifies Dialog.Message() returns prompt text.
// Extra: no equivalent method in playwright-java; related class: TestPageDialog.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestOnDialogMessageText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var dialogMsg string

	off := page.OnDialog(func(d *playwright.Dialog) {
		mu.Lock()
		dialogMsg = d.Message()
		mu.Unlock()
		_ = d.Dismiss(ctx)
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	_, _ = page.Evaluate(ctx, `() => alert('the message')`)

	mu.Lock()
	msg := dialogMsg
	mu.Unlock()

	is.Equal("the message", msg)
}

// TestPageDialogShouldAllowAcceptingPrompts verifies Accept(value) returns the entered value.
// Ref: TestPageDialog.java#shouldAllowAcceptingPrompts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldAllowAcceptingPrompts(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	page.OnDialog(func(d *playwright.Dialog) {
		is.Equal("prompt", d.Type())
		is.Equal("yes.", d.DefaultValue())
		_ = d.Accept(ctx, "answer!")
	})

	result, err := page.Evaluate(ctx, "prompt('question?', 'yes.')")
	must.NoError(err)
	is.Equal("answer!", result)
}

// TestPageDialogShouldDismissThePrompt verifies Dismiss returns null for prompt.
// Ref: TestPageDialog.java#shouldDismissThePrompt  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldDismissThePrompt(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Dismiss(ctx)
	})

	result, err := page.Evaluate(ctx, "prompt('question?')")
	must.NoError(err)
	is.Nil(result)
}

// TestPageDialogShouldAcceptTheConfirmPrompt verifies that Accept returns true for confirm.
// Ref: TestPageDialog.java#shouldAcceptTheConfirmPrompt  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldAcceptTheConfirmPrompt(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Accept(ctx)
	})

	result, err := page.Evaluate(ctx, "confirm('boolean?')")
	must.NoError(err)
	is.Equal(true, result)
}

// TestPageDialogShouldDismissTheConfirmPrompt verifies that Dismiss returns false for confirm.
// Ref: TestPageDialog.java#shouldDismissTheConfirmPrompt  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldDismissTheConfirmPrompt(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Dismiss(ctx)
	})

	result, err := page.Evaluate(ctx, "() => confirm('boolean?')")
	must.NoError(err)
	is.Equal(false, result)
}

// TestPageDialogShouldHandleMultipleAlerts verifies that multiple sequential alerts are all handled.
// Ref: TestPageDialog.java#shouldHandleMultipleAlerts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldHandleMultipleAlerts(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Accept(ctx)
	})

	err := page.SetContent(ctx, `<p>Hello World</p>
	<script>
		alert('Please dismiss this dialog');
		alert('Please dismiss this dialog');
		alert('Please dismiss this dialog');
	</script>`)
	must.NoError(err)

	text, err := page.Locator("p").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Hello World", *text)
}

// TestPageDialogShouldHandleMultipleConfirms verifies that multiple sequential confirms are all handled.
// Ref: TestPageDialog.java#shouldHandleMultipleConfirms  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldHandleMultipleConfirms(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	page.OnDialog(func(d *playwright.Dialog) {
		_ = d.Accept(ctx)
	})

	err := page.SetContent(ctx, `<p>Hello World</p>
	<script>
		confirm('Please confirm me?');
		confirm('Please confirm me?');
		confirm('Please confirm me?');
	</script>`)
	must.NoError(err)

	text, err := page.Locator("p").TextContent(ctx)
	must.NoError(err)
	must.NotNil(text)
	is.Equal("Hello World", *text)
}

// TestPageDialogShouldAutoDismissThePromptWithoutListeners verifies prompt is auto-dismissed without a listener.
// Ref: TestPageDialog.java#shouldAutoDismissThePromptWithoutListeners  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldAutoDismissThePromptWithoutListeners(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => prompt('question?')")
	must.NoError(err)
	is.Nil(result)
}

// Ref: TestPageDialog.java#shouldFire  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldFire(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	var fired atomic.Bool
	page.OnDialog(func(d *playwright.Dialog) {
		assert.Equal(t, "alert", d.Type())
		assert.Equal(t, "", d.DefaultValue())
		assert.Equal(t, "yo", d.Message())
		fired.Store(true)
		assert.NoError(t, d.Accept(ctx))
	})
	a2Ev(t, ctx, page, "alert('yo')")
	assert.True(t, fired.Load())
}

// Ref: TestPageDialog.java#shouldBeAbleToCloseContextWithOpenAlert  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldBeAbleToCloseContextWithOpenAlert(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	var shown atomic.Bool
	page.OnDialog(func(d *playwright.Dialog) { shown.Store(true) })
	a2Ev(t, ctx, page, "() => setTimeout(() => alert('hello'), 0)")
	must.Eventually(shown.Load, 10*time.Second, 100*time.Millisecond)
	must.NoError(bCtx.Close(ctx))
}

// Ref: TestPageDialog.java#shouldAutoDismissTheAlertWithoutListeners  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldAutoDismissTheAlertWithoutListeners(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.SetContent(ctx, "<div onclick='window.alert(123); window._clicked=true'>Click me</div>"))
	must.NoError(a2Click(ctx, page, "div"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "window._clicked"))
}

// Ref: TestPageDialog.java#shouldFireDialogClosedWhenDialogIsAccepted  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldFireDialogClosedWhenDialogIsAccepted(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, _, page := a2Setup(t)

	var closed []*playwright.Dialog
	var mu sync.Mutex
	offClosed := page.OnDialogClosed(func(d *playwright.Dialog) {
		mu.Lock()
		closed = append(closed, d)
		mu.Unlock()
	})
	defer offClosed()

	var opened *playwright.Dialog
	offDialog := page.OnDialog(func(d *playwright.Dialog) {
		opened = d
		must.NoError(d.Accept(ctx))
	})
	defer offDialog()

	_, err := page.Evaluate(ctx, "alert('yo')")
	must.NoError(err)

	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(closed) == 1
	}, 10*time.Second, 50*time.Millisecond)

	mu.Lock()
	is.Equal(opened, closed[0])
	mu.Unlock()

	// Perform some roundtrips to ensure the event does not fire twice.
	_, err = page.Evaluate(ctx, "1")
	must.NoError(err)
	_, err = page.Evaluate(ctx, "1")
	must.NoError(err)

	mu.Lock()
	is.Equal(1, len(closed))
	mu.Unlock()
}

// Ref: TestPageDialog.java#shouldFireDialogClosedWhenDialogIsDismissed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldFireDialogClosedWhenDialogIsDismissed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, _, page := a2Setup(t)

	var closed *playwright.Dialog
	var mu sync.Mutex
	offClosed := page.OnDialogClosed(func(d *playwright.Dialog) {
		mu.Lock()
		closed = d
		mu.Unlock()
	})
	defer offClosed()

	offDialog := page.OnDialog(func(d *playwright.Dialog) {
		must.NoError(d.Dismiss(ctx))
	})
	defer offDialog()

	_, err := page.Evaluate(ctx, "confirm('boolean?')")
	must.NoError(err)

	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return closed != nil
	}, 10*time.Second, 50*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	is.Equal("confirm", closed.Type())
	is.Equal("boolean?", closed.Message())
}

// Ref: TestPageDialog.java#shouldFireDialogClosedForAutoDismissedDialogs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDialog.java
func TestPageDialogShouldFireDialogClosedForAutoDismissedDialogs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, _, page := a2Setup(t)

	var closed *playwright.Dialog
	var mu sync.Mutex
	offClosed := page.OnDialogClosed(func(d *playwright.Dialog) {
		mu.Lock()
		closed = d
		mu.Unlock()
	})
	defer offClosed()

	_, err := page.Evaluate(ctx, "alert('yo')")
	must.NoError(err)

	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return closed != nil
	}, 10*time.Second, 50*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	is.Equal("yo", closed.Message())
}
