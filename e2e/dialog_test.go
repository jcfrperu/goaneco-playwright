//go:build e2e

package e2e

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

// TestDialogAcceptPromptWithValue verifies Accept(value) returns the entered value.
// Ref: TestPageDialog.java#shouldAllowAcceptingPrompts
func TestDialogAcceptPromptWithValue(t *testing.T) {
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

// TestDialogDismissPromptReturnsNull verifies Dismiss returns null for prompt.
// Ref: TestPageDialog.java#shouldDismissThePrompt
func TestDialogDismissPromptReturnsNull(t *testing.T) {
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

// TestDialogAcceptConfirm verifies that Accept returns true for confirm.
// Ref: TestPageDialog.java#shouldAcceptTheConfirmPrompt
func TestDialogAcceptConfirm(t *testing.T) {
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

// TestDialogDismissConfirm verifies that Dismiss returns false for confirm.
// Ref: TestPageDialog.java#shouldDismissTheConfirmPrompt
func TestDialogDismissConfirm(t *testing.T) {
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

// TestDialogMultipleAlerts verifies that multiple sequential alerts are all handled.
// Ref: TestPageDialog.java#shouldHandleMultipleAlerts
func TestDialogMultipleAlerts(t *testing.T) {
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
	is.Equal("Hello World", text)
}

// TestDialogMultipleConfirms verifies that multiple sequential confirms are all handled.
// Ref: TestPageDialog.java#shouldHandleMultipleConfirms
func TestDialogMultipleConfirms(t *testing.T) {
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
	is.Equal("Hello World", text)
}

// TestDialogAutoDismissWithoutListeners verifies prompt is auto-dismissed without a listener.
// Ref: TestPageDialog.java#shouldAutoDismissThePromptWithoutListeners
func TestDialogAutoDismissWithoutListeners(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => prompt('question?')")
	must.NoError(err)
	is.Nil(result)
}

// TestOnDialogAlertAccepted verifies OnDialog can accept an alert.
// Ref: TestPageOnDialog.java#shouldAcceptAlert
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
// Ref: TestPageOnDialog.java#shouldAcceptConfirm
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
// Ref: TestPageOnDialog.java#shouldDismissConfirm
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
// Ref: TestPageOnDialog.java#shouldAcceptPromptWithValue
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
// Ref: TestPageOnDialog.java#shouldReturnMessage
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
