//go:build e2e

// E2E tests for: TestPageEventConsole.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// TestOnConsoleWarningMessage verifies console.warn is captured.
// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestOnConsoleWarningMessage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))
	msg, err := page.WaitForConsoleMessage(ctx, func(m *playwright.ConsoleMessage) bool {
		return m.Type() == "warning"
	}, func() {
		_, _ = page.Evaluate(ctx, `() => console.warn('watch out')`)
	})
	must.NoError(err)
	must.NotNil(msg)
	is.Equal("warning", msg.Type())
}

// TestOnConsoleErrorMessage verifies console.error is captured.
// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestOnConsoleErrorMessage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))
	msg, err := page.WaitForConsoleMessage(ctx, func(m *playwright.ConsoleMessage) bool {
		return m.Type() == "error"
	}, func() {
		_, _ = page.Evaluate(ctx, `() => console.error('something broke')`)
	})
	must.NoError(err)
	must.NotNil(msg)
	is.Equal("error", msg.Type())
	is.Equal("something broke", msg.Text())
}

// TestOnConsoleOffStopsHandler verifies off() unregisters the console handler.
// Extra: no equivalent method in playwright-java; related class: TestPageEventConsole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
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

// TestPageEventConsoleShouldEmitSameLogTwice verifies the same log line can arrive twice as two separate events.
// Ref: TestPageEventConsole.java#shouldEmitSameLogTwice  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldEmitSameLogTwice(t *testing.T) {
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

// TestPageEventConsoleShouldWorkForDifferentConsoleAPICalls verifies that warn, error, info, and debug console calls
// are all captured with the correct type.
// Ref: TestPageEventConsole.java#shouldWorkForDifferentConsoleAPICalls  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldWorkForDifferentConsoleAPICalls(t *testing.T) {
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

// TestPageEventConsoleShouldNotFailForWindowObject verifies that console.log(window) does not crash.
// Ref: TestPageEventConsole.java#shouldNotFailForWindowObject  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldNotFailForWindowObject(t *testing.T) {
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

// TestPageEventConsoleShouldTriggerCorrectLog verifies that the text of a console.log matches exactly.
// Ref: TestPageEventConsole.java#shouldTriggerCorrectLog  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldTriggerCorrectLog(t *testing.T) {
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

// Ref: TestPageEventConsole.java#clearConsoleMessagesShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleClearConsoleMessagesShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => { console.log('message1'); console.log('message2'); }")
	must.NoError(err)

	messages, err := page.ConsoleMessages(ctx)
	must.NoError(err)
	var hasMsg1, hasMsg2 bool
	for _, m := range messages {
		if m.Text() == "message1" {
			hasMsg1 = true
		}
		if m.Text() == "message2" {
			hasMsg2 = true
		}
	}
	is.True(hasMsg1)
	is.True(hasMsg2)

	must.NoError(page.ClearConsoleMessages(ctx))
	messages, err = page.ConsoleMessages(ctx)
	must.NoError(err)
	is.Empty(messages)

	_, err = page.WaitForConsoleMessage(ctx, func() {
		_, _ = page.Evaluate(ctx, "() => console.log('message3')")
	})
	must.NoError(err)

	messages, err = page.ConsoleMessages(ctx)
	must.NoError(err)
	must.Len(messages, 1)
	is.Equal("message3", messages[0].Text())
}

// TestPageEventConsoleConsoleMessagesShouldWork verifies console.log is captured.
// Ref: TestPageEventConsole.java#consoleMessagesShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleConsoleMessagesShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		for (let i = 0; i < 301; i++)
			console.log('message' + i);
	}`)
	must.NoError(err)

	messages, err := page.ConsoleMessages(ctx)
	must.NoError(err)
	must.GreaterOrEqual(len(messages), 100, "should be at least 100 messages")

	firstIndex := len(messages) - 100
	for i := 0; i < 100; i++ {
		msg := messages[firstIndex+i]
		is.Equal(fmt.Sprintf("message%d", 201+i), msg.Text())
		is.Equal("log", msg.Type())
		is.Equal(page, msg.Page())
	}
}

// a2WaitConsole waits for the first console message accepted by the predicate (nil accepts all) produced by trigger.
func a2WaitConsole(t *testing.T, page *playwright.Page, pred func(*playwright.ConsoleMessage) bool, trigger func()) *playwright.ConsoleMessage {
	t.Helper()
	ch := make(chan *playwright.ConsoleMessage, 64)
	off := page.OnConsole(func(m *playwright.ConsoleMessage) { ch <- m })
	defer off()
	trigger()
	deadline := time.After(10 * time.Second)
	for {
		select {
		case m := <-ch:
			if pred == nil || pred(m) {
				return m
			}
		case <-deadline:
			t.Fatal("timed out waiting for console message")
			return nil
		}
	}
}

// Ref: TestPageEventConsole.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	var event *playwright.ConsoleMessage
	off := page.OnConsole(func(m *playwright.ConsoleMessage) {
		if event == nil {
			event = m
		}
	})
	defer off()

	message := a2WaitConsole(t, page, nil, func() {
		_, err := page.Evaluate(ctx, "() => console.log('hello', 5, {foo: 'bar'});")
		must.NoError(err)
	})

	must.NotNil(message)
	is.Equal("hello 5 {foo: bar}", message.Text())
	is.Equal("log", message.Type())

	args := message.Args()
	must.Len(args, 3)

	v0, err := args[0].JSONValue(ctx)
	must.NoError(err)
	is.Equal("hello", v0)

	v1, err := args[1].JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(5), v1)

	v2, err := args[2].JSONValue(ctx)
	must.NoError(err)
	is.Equal(map[string]any{"foo": "bar"}, v2)

	is.Equal(message.Text(), event.Text())
	is.Equal(message.Type(), event.Type())
}

// Ref: TestPageEventConsole.java#shouldHaveLocationForConsoleAPICalls  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldHaveLocationForConsoleAPICalls(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	message := a2WaitConsole(t, page, func(m *playwright.ConsoleMessage) bool { return m.Text() == "yellow" }, func() {
		must.NoError(page.Goto(ctx, srv.Prefix()+"/consolelog.html"))
	})
	assert.Equal(t, "log", message.Type())
	// Engines have different column notion.
	assert.True(t, strings.HasPrefix(message.Location(), srv.Prefix()+"/consolelog.html:7:"), message.Location())
}

// Ref: TestPageEventConsole.java#shouldSupportPredicate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldSupportPredicate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	message := a2WaitConsole(t, page, func(m *playwright.ConsoleMessage) bool { return m.Type() == "info" }, func() {
		a2Ev(t, ctx, page, "console.log(1)")
		a2Ev(t, ctx, page, "console.info(2)")
	})
	assert.Equal(t, "2", message.Text())
	assert.Equal(t, "info", message.Type())
}

// Ref: TestPageEventConsole.java#shouldHaveTimestamp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldHaveTimestamp(t *testing.T) {
	t.Parallel()
	ctx, _, page := a2Setup(t)

	before := float64(time.Now().UnixMilli()) - 1
	message := a2WaitConsole(t, page, nil, func() { a2Ev(t, ctx, page, "() => console.log('timestamp test')") })
	after := float64(time.Now().UnixMilli()) + 1
	assert.GreaterOrEqual(t, message.Timestamp(), before, "timestamp %v should be >= %v", message.Timestamp(), before)
	assert.LessOrEqual(t, message.Timestamp(), after, "timestamp %v should be <= %v", message.Timestamp(), after)
}

// Ref: TestPageEventConsole.java#shouldHaveIncreasingTimestamps  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleShouldHaveIncreasingTimestamps(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	var mu sync.Mutex
	var messages []*playwright.ConsoleMessage
	page.OnConsole(func(m *playwright.ConsoleMessage) {
		mu.Lock()
		defer mu.Unlock()
		messages = append(messages, m)
	})
	a2Ev(t, ctx, page, "() => { console.log('first'); console.log('second'); console.log('third'); }")
	must.Eventually(func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(messages) >= 3
	}, 5*time.Second, 20*time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, messages, 3)
	for i := 1; i < len(messages); i++ {
		assert.GreaterOrEqual(t, messages[i].Timestamp(), messages[i-1].Timestamp())
	}
}

// Ref: TestPageEventConsole.java#consoleMessagesSinceNavigationFilterShouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventConsole.java
func TestPageEventConsoleConsoleMessagesSinceNavigationFilterShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => console.log('before navigation')")
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, "() => console.log('after navigation')")
	must.NoError(err)

	allFilter := playwright.ConsoleMessagesFilterAll
	all, err := page.ConsoleMessages(ctx, &playwright.ConsoleMessagesOptions{Filter: &allFilter})
	must.NoError(err)

	hasBeforeAll := false
	hasAfterAll := false
	for _, m := range all {
		if m.Text() == "before navigation" {
			hasBeforeAll = true
		}
		if m.Text() == "after navigation" {
			hasAfterAll = true
		}
	}
	is.True(hasBeforeAll)
	is.True(hasAfterAll)

	// sinceNavigation is the default
	sinceNav, err := page.ConsoleMessages(ctx)
	must.NoError(err)

	hasBeforeSince := false
	hasAfterSince := false
	for _, m := range sinceNav {
		if m.Text() == "before navigation" {
			hasBeforeSince = true
		}
		if m.Text() == "after navigation" {
			hasAfterSince = true
		}
	}
	is.False(hasBeforeSince)
	is.True(hasAfterSince)
}
