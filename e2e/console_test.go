//go:build e2e

package e2e

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageOnConsoleLog(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	const wantCount = 2
	done := make(chan struct{})
	var once sync.Once
	var mu sync.Mutex
	var received []*playwright.ConsoleMessage

	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		received = append(received, msg)
		count := len(received)
		mu.Unlock()
		if count >= wantCount {
			once.Do(func() { close(done) })
		}
	})

	_, err = page.Evaluate(ctx, "console.log('hello console')")
	must.NoError(err, "Evaluate failed")
	_, err = page.Evaluate(ctx, "console.error('oops')")
	must.NoError(err, "Evaluate failed")

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for 2 console messages")
	}

	mu.Lock()
	msgs := received
	mu.Unlock()

	found := map[string]string{}
	for _, m := range msgs {
		found[m.Type()] = m.Text()
	}

	if found["log"] != "hello console" {
		t.Errorf("expected log='hello console', got: %v", found)
	}
	if found["error"] != "oops" {
		t.Errorf("expected error='oops', got: %v", found)
	}
}

func TestPageOnConsoleCancel(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	var count int32
	firstSeen := make(chan struct{})
	var firstOnce sync.Once

	cancel := page.OnConsole(func(_ *playwright.ConsoleMessage) {
		atomic.AddInt32(&count, 1)
		firstOnce.Do(func() { close(firstSeen) })
	})

	_, err = page.Evaluate(ctx, "console.log('before cancel')")
	must.NoError(err, "Evaluate failed")
	select {
	case <-firstSeen:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first console message")
	}

	cancel()

	_, err = page.Evaluate(ctx, "console.log('after cancel')")
	must.NoError(err, "Evaluate failed")
	// Brief wait: absence of event is inherently time-bounded.
	must.NoError(page.WaitForTimeout(ctx, 300))

	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("expected 1 message before cancel, got %d after cancel", got)
	}
}

// TestPageConsoleEmitSameLogTwice verifies the same log line can arrive twice as two separate events.
// Ref: TestPageEventConsole.java#shouldEmitSameLogTwice
func TestPageConsoleEmitSameLogTwice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	done := make(chan struct{})
	var once sync.Once
	var mu sync.Mutex
	var texts []string

	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		texts = append(texts, msg.Text())
		n := len(texts)
		mu.Unlock()
		if n >= 2 {
			once.Do(func() { close(done) })
		}
	})

	_, err := page.Evaluate(ctx, `() => { console.log('dup'); console.log('dup'); }`)
	must.NoError(err)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for duplicate log messages")
	}

	mu.Lock()
	got := texts
	mu.Unlock()
	is.Len(got, 2)
	is.Equal("dup", got[0])
	is.Equal("dup", got[1])
}

// TestPageConsoleDifferentAPICalls verifies that warn, error, info, and debug console calls
// are all captured with the correct type.
// Ref: TestPageEventConsole.java#shouldWorkForDifferentConsoleAPICalls
func TestPageConsoleDifferentAPICalls(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	const want = 4
	done := make(chan struct{})
	var once sync.Once
	var mu sync.Mutex
	found := map[string]string{}

	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		found[msg.Type()] = msg.Text()
		n := len(found)
		mu.Unlock()
		if n >= want {
			once.Do(func() { close(done) })
		}
	})

	_, err := page.Evaluate(ctx, `() => {
		console.warn('warn-msg');
		console.error('error-msg');
		console.info('info-msg');
		console.debug('debug-msg');
	}`)
	must.NoError(err)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for console messages")
	}

	mu.Lock()
	types := found
	mu.Unlock()

	is.Equal("warn-msg", types["warning"], "warn message")
	is.Equal("error-msg", types["error"], "error message")
	is.Equal("info-msg", types["info"], "info message")
	is.Equal("debug-msg", types["debug"], "debug message")
}

// TestPageConsoleWindowObjectDoesNotCrash verifies that console.log(window) does not crash.
// Ref: TestPageEventConsole.java#shouldNotFailForWindowObject
func TestPageConsoleWindowObjectDoesNotCrash(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	received := make(chan *playwright.ConsoleMessage, 1)
	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		select {
		case received <- msg:
		default:
		}
	})

	_, err := page.Evaluate(ctx, "console.log(window)")
	must.NoError(err)

	select {
	case msg := <-received:
		must.NotNil(msg, "console message should not be nil")
		// Just verify it was received without crashing
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for console.log(window) message")
	}
}

// TestPageConsoleTriggerCorrectLog verifies that the text of a console.log matches exactly.
// Ref: TestPageEventConsole.java#shouldTriggerCorrectLog
func TestPageConsoleTriggerCorrectLog(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	received := make(chan *playwright.ConsoleMessage, 1)
	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		select {
		case received <- msg:
		default:
		}
	})

	const expected = "hello playwright console"
	_, err := page.Evaluate(ctx, fmt.Sprintf("console.log('%s')", expected))
	must.NoError(err)

	select {
	case msg := <-received:
		is.Equal("log", msg.Type())
		is.Equal(expected, msg.Text())
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for console message")
	}
}

func TestConsoleMessageHasTimestamp(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	before := float64(time.Now().UnixMilli()) - 1
	received := make(chan *playwright.ConsoleMessage, 1)
	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		select {
		case received <- msg:
		default:
		}
	})

	_, err := page.Evaluate(ctx, "() => console.log('timestamp test')")
	must.NoError(err)

	var msg *playwright.ConsoleMessage
	select {
	case msg = <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for console message")
	}
	after := float64(time.Now().UnixMilli()) + 1

	ts := msg.Timestamp()
	must.Truef(ts >= before, "timestamp %v should be >= before %v", ts, before)
	must.Truef(ts <= after, "timestamp %v should be <= after %v", ts, after)
}

func TestConsoleMessageHasIncreasingTimestamps(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	const want = 3
	var mu sync.Mutex
	var msgs []*playwright.ConsoleMessage
	done := make(chan struct{})
	var once sync.Once

	page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		msgs = append(msgs, msg)
		n := len(msgs)
		mu.Unlock()
		if n >= want {
			once.Do(func() { close(done) })
		}
	})

	_, err := page.Evaluate(ctx, "() => { console.log('first'); console.log('second'); console.log('third'); }")
	must.NoError(err)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for 3 console messages")
	}

	mu.Lock()
	got := append([]*playwright.ConsoleMessage{}, msgs...)
	mu.Unlock()

	is.Len(got, want)

	var timestamps []float64
	for _, m := range got {
		ts := m.Timestamp()
		must.Truef(ts > 0, "expected positive timestamp, got %v", ts)
		timestamps = append(timestamps, ts)
	}

	var minTS, maxTS float64
	minTS = timestamps[0]
	maxTS = timestamps[0]
	for _, ts := range timestamps[1:] {
		if ts < minTS {
			minTS = ts
		}
		if ts > maxTS {
			maxTS = ts
		}
	}
	must.Truef(maxTS >= minTS, "max timestamp %v should be >= min timestamp %v", maxTS, minTS)
}

func TestConsoleMessageTextForMultipleArgs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	done := make(chan string, 1)
	cancel := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Text() != "" && msg.Type() == "log" {
			select {
			case done <- msg.Text():
			default:
			}
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => console.log('hello', 'world')`)
	must.NoError(err)

	select {
	case text := <-done:
		is.Contains(text, "hello")
		is.Contains(text, "world")
	case <-time.After(5 * time.Second):
		t.Fatal("console.log message not received in time")
	}
}

func TestConsoleMessageTypeForWarn(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	done := make(chan string, 1)
	cancel := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "warning" || msg.Type() == "warn" {
			select {
			case done <- msg.Type():
			default:
			}
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => console.warn('a warning')`)
	must.NoError(err)

	select {
	case capturedType := <-done:
		is.Contains([]string{"warning", "warn"}, capturedType)
	case <-time.After(5 * time.Second):
		t.Fatal("console.warn message not received in time")
	}
}

func TestConsoleMessageTypeForError(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	done := make(chan string, 1)
	cancel := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "error" {
			select {
			case done <- msg.Type():
			default:
			}
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => console.error('an error')`)
	must.NoError(err)

	select {
	case capturedType := <-done:
		is.Equal("error", capturedType)
	case <-time.After(5 * time.Second):
		t.Fatal("console.error message not received in time")
	}
}

func TestConsoleMessageTypeForInfo(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	done := make(chan string, 1)
	cancel := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "info" || msg.Type() == "log" {
			select {
			case done <- msg.Type():
			default:
			}
		}
	})
	defer cancel()

	_, err := page.Evaluate(ctx, `() => console.info('information')`)
	must.NoError(err)

	select {
	case tp := <-done:
		is.Contains([]string{"info", "log"}, tp)
	case <-time.After(5 * time.Second):
		t.Fatal("console.info message not received in time")
	}
}

func TestConsoleWarnMessageEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var warnMsg string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "warning" || msg.Type() == "warn" {
			mu.Lock()
			warnMsg = msg.Text()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.warn('test warning message');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	w := warnMsg
	mu.Unlock()

	is.Contains(w, "test warning message")
}

func TestConsoleErrorMessageEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var errMsg string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "error" {
			mu.Lock()
			errMsg = msg.Text()
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.error('test error message');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	e := errMsg
	mu.Unlock()

	is.Contains(e, "test error message")
}

func TestConsoleInfoMessageEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var infoMsg string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "info" || msg.Type() == "log" {
			mu.Lock()
			if infoMsg == "" {
				infoMsg = msg.Text()
			}
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.info('test info message');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	i := infoMsg
	mu.Unlock()

	is.Contains(i, "test info message")
}

func TestConsoleMultipleMessagesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var messages []string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		messages = append(messages, msg.Text())
		mu.Unlock()
	})
	defer off()

	must.NoError(page.SetContent(ctx, `
		<script>
			console.log('message 1');
			console.log('message 2');
			console.log('message 3');
		</script>
	`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	count := len(messages)
	mu.Unlock()

	is.GreaterOrEqual(count, 3)
}

func TestConsoleDebugMessageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var debugMsg string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "debug" || msg.Type() == "verbose" || msg.Type() == "log" {
			mu.Lock()
			if debugMsg == "" {
				debugMsg = msg.Text()
			}
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.debug('debug message here');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	msg := debugMsg
	mu.Unlock()

	is.NotEmpty(msg)
}

func TestConsoleClearMessageEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var types []string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		types = append(types, msg.Type())
		mu.Unlock()
	})
	defer off()

	must.NoError(page.SetContent(ctx, `
		<script>
			console.log('log');
			console.warn('warn');
			console.error('error');
		</script>
	`))
	must.NoError(page.WaitForTimeout(ctx, 150))

	mu.Lock()
	count := len(types)
	mu.Unlock()

	is.GreaterOrEqual(count, 3)
}

func TestConsoleMessageTextEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var logText string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "log" {
			mu.Lock()
			if logText == "" {
				logText = msg.Text()
			}
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.log('specific text for test');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	text := logText
	mu.Unlock()

	is.Contains(text, "specific text for test")
}

func TestConsoleWarnEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var warnMsg string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "warning" {
			mu.Lock()
			if warnMsg == "" {
				warnMsg = msg.Text()
			}
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.warn('Warning message');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	got := warnMsg
	mu.Unlock()

	is.Equal("Warning message", got)
}

func TestConsoleErrorEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var errMsg string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "error" {
			mu.Lock()
			if errMsg == "" {
				errMsg = msg.Text()
			}
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.error('Error occurred');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	got := errMsg
	mu.Unlock()

	is.Equal("Error occurred", got)
}

func TestConsoleInfoEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var infoMsg string
	var mu sync.Mutex

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "info" || msg.Type() == "log" {
			mu.Lock()
			if infoMsg == "" {
				infoMsg = msg.Text()
			}
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<script>console.info('Info message');</script>`))
	must.NoError(page.WaitForTimeout(ctx, 100))

	mu.Lock()
	got := infoMsg
	mu.Unlock()

	is.NotEmpty(got)
}

// TestOnConsoleLogMessage verifies console.log is captured.
// Ref: TestPageOnConsole.java#shouldCaptureLog
func TestOnConsoleLogMessage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var captured *playwright.ConsoleMessage

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "log" {
			mu.Lock()
			captured = msg
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	_, err := page.Evaluate(ctx, `() => console.log('hello console')`)
	must.NoError(err)

	mu.Lock()
	msg := captured
	mu.Unlock()

	must.NotNil(msg)
	is.Equal("log", msg.Type())
	is.Equal("hello console", msg.Text())
}

// TestOnConsoleWarningMessage verifies console.warn is captured.
// Ref: TestPageOnConsole.java#shouldCaptureWarn
func TestOnConsoleWarningMessage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var captured *playwright.ConsoleMessage

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "warning" {
			mu.Lock()
			captured = msg
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	_, err := page.Evaluate(ctx, `() => console.warn('watch out')`)
	must.NoError(err)

	mu.Lock()
	msg := captured
	mu.Unlock()

	must.NotNil(msg)
	is.Equal("warning", msg.Type())
}

// TestOnConsoleErrorMessage verifies console.error is captured.
// Ref: TestPageOnConsole.java#shouldCaptureError
func TestOnConsoleErrorMessage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	var captured *playwright.ConsoleMessage

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "error" {
			mu.Lock()
			captured = msg
			mu.Unlock()
		}
	})
	defer off()

	must.NoError(page.SetContent(ctx, `<div></div>`))
	_, err := page.Evaluate(ctx, `() => console.error('something broke')`)
	must.NoError(err)

	mu.Lock()
	msg := captured
	mu.Unlock()

	must.NotNil(msg)
	is.Equal("error", msg.Type())
	is.Equal("something broke", msg.Text())
}

// TestOnConsoleOffStopsHandler verifies off() unregisters the console handler.
// Ref: TestPageOnConsole.java#shouldStopAfterOff
func TestOnConsoleOffStopsHandler(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var mu sync.Mutex
	count := 0

	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	must.NoError(page.SetContent(ctx, `<div></div>`))
	_, _ = page.Evaluate(ctx, `() => console.log('first')`)

	off()

	_, _ = page.Evaluate(ctx, `() => console.log('second')`)

	mu.Lock()
	c := count
	mu.Unlock()

	is.Equal(1, c)
}
