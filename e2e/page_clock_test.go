//go:build e2e

// E2E tests for: TestPageClock.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestClockSetFixedTime verifies that SetFixedTime fixes Date.now() to the given value.
// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockSetFixedTime(t *testing.T) {
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
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Fix time to 2024-01-15T00:00:00Z (Unix ms: 1705276800000)
	const fixedMs = 1705276800000.0
	err = bCtx.Clock().SetFixedTime(ctx, fixedMs)
	must.NoError(err, "SetFixedTime failed")

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err, "Evaluate failed")
	got, _ := raw.(float64)
	if got != fixedMs {
		t.Errorf("Date.now() = %v, want %v", got, fixedMs)
	}
}

// TestClockInstallAndFastForward verifies that after clock Install, FastForward
// advances the fake clock so setTimeout callbacks fire immediately.
// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockInstallAndFastForward(t *testing.T) {
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
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Install fake clock starting at t=0
	err = bCtx.Clock().Install(ctx, 0)
	must.NoError(err, "Clock.Install failed")

	// Schedule a 5-second timeout in the page
	_, err = page.Evaluate(ctx, `() => {
		window._fired = false;
		setTimeout(() => { window._fired = true; }, 5000);
	}`)
	must.NoError(err, "Evaluate failed")

	// The timeout should not have fired yet
	firedRaw, err := page.Evaluate(ctx, "() => window._fired")
	must.NoError(err, "Evaluate(fired check) failed")
	if firedRaw == true {
		t.Error("setTimeout fired before FastForward — unexpected")
	}

	// Fast-forward 6 seconds
	err = bCtx.Clock().FastForward(ctx, 6000)
	must.NoError(err, "Clock.FastForward failed")

	// Now the timeout should have fired
	firedRaw, err = page.Evaluate(ctx, "() => window._fired")
	must.NoError(err, "Evaluate(fired check after ff) failed")
	if firedRaw != true {
		t.Error("setTimeout did not fire after FastForward(6000)")
	}
}

// TestClockRunFor verifies that RunFor advances the fake clock and fires pending timers.
// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunFor(t *testing.T) {
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
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	err = bCtx.Clock().Install(ctx, 0)
	must.NoError(err, "Clock.Install failed")

	_, err = page.Evaluate(ctx, `() => {
		window._fired = false;
		setTimeout(() => { window._fired = true; }, 500);
	}`)
	must.NoError(err, "Evaluate failed")

	err = bCtx.Clock().RunFor(ctx, 1000)
	must.NoError(err, "Clock.RunFor failed")

	raw, err := page.Evaluate(ctx, "() => window._fired")
	must.NoError(err, "Evaluate(fired) failed")
	if raw != true {
		t.Error("timer should have fired after Clock.RunFor(1000)")
	}
}

// TestClockPauseAt verifies that PauseAt moves the clock to the given absolute time.
// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockPauseAt(t *testing.T) {
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
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Install at a known start time (1000ms), then pause further in the future.
	const startMs = 1000.0
	const pauseMs = 5000.0
	err = bCtx.Clock().Install(ctx, startMs)
	must.NoError(err, "Clock.Install failed")
	err = bCtx.Clock().PauseAt(ctx, pauseMs)
	must.NoError(err, "Clock.PauseAt failed")

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err, "Evaluate failed")
	got, _ := raw.(float64)
	if got != pauseMs {
		t.Errorf("Date.now() after PauseAt = %v, want %v", got, pauseMs)
	}
}

// TestClockResume verifies that Resume re-enables the clock after PauseAt so timers fire.
// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockResume(t *testing.T) {
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
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	err = bCtx.Clock().Install(ctx, 1000)
	must.NoError(err, "Clock.Install failed")
	err = bCtx.Clock().PauseAt(ctx, 1000)
	must.NoError(err, "Clock.PauseAt failed")

	// Schedule a timer; it won't fire while clock is paused.
	_, err = page.Evaluate(ctx, `() => {
		window._resumed = false;
		setTimeout(() => { window._resumed = true; }, 100);
	}`)
	must.NoError(err, "Evaluate failed")

	// Advance by 500ms (fires the 100ms timer) and then resume real-time.
	err = bCtx.Clock().RunFor(ctx, 500)
	must.NoError(err, "Clock.RunFor failed")
	err = bCtx.Clock().Resume(ctx)
	must.NoError(err, "Clock.Resume failed")

	raw, err := page.Evaluate(ctx, "() => window._resumed")
	must.NoError(err, "Evaluate(resumed) failed")
	if raw != true {
		t.Error("timer should have fired after RunFor+Resume")
	}
}

// TestClockSetSystemTime verifies that SetSystemTime changes the clock value.
// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockSetSystemTime(t *testing.T) {
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
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// Install fake clock first, then set system time
	err = bCtx.Clock().Install(ctx)
	must.NoError(err, "Clock.Install failed")

	const targetMs = 1700000000000.0
	err = bCtx.Clock().SetSystemTime(ctx, targetMs)
	must.NoError(err, "SetSystemTime failed")

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err, "Evaluate failed")
	got, _ := raw.(float64)
	// Allow small delta since the clock might tick slightly
	if got < targetMs || got > targetMs+1000 {
		t.Errorf("Date.now() = %v, expected ~%v", got, targetMs)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunForImmediately(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setTimeout(window.stub)")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 0))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(1), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunForInsufficientDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setTimeout(window.stub, 100)")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 11))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(0), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunForSufficientDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setTimeout(window.stub, 100)")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 100))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(1), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunForSimultaneousTimers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => { setTimeout(window.stub, 100); setTimeout(window.stub, 100); }")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 100))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(2), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunForMultipleSimultaneous(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => { setTimeout(window.stub, 100); setTimeout(window.stub, 100); setTimeout(window.stub, 99); setTimeout(window.stub, 100); }")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 100))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(4), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunForWaitsAfterTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setTimeout(window.stub, 150)")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 50))
	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(0), raw, "should not have fired yet at t=50")

	must.NoError(bCtx.Clock().RunFor(ctx, 100))
	raw, err = page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(1), raw, "should have fired after t=150 total")
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockRunForReturnsNowValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))
	must.NoError(bCtx.Clock().SetSystemTime(ctx, 0))

	must.NoError(bCtx.Clock().RunFor(ctx, 200))

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	is.Equal(float64(200), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockFastForwardSkipsTimers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setTimeout(() => { window.stub('should not be logged'); }, 1000)")
	must.NoError(err)

	must.NoError(bCtx.Clock().FastForward(ctx, 500))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(0), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockFastForwardPushesBackTime(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._fired = []; window.stub = function(ts) { window._fired.push(ts); }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setTimeout(() => { window.stub(Date.now()); }, 1000)")
	must.NoError(err)

	must.NoError(bCtx.Clock().FastForward(ctx, 2000))

	raw, err := page.Evaluate(ctx, "() => window._fired[0]")
	must.NoError(err)

	is.Equal(float64(3000), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockSetTimeoutReturnsId(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	raw, err := page.Evaluate(ctx, "() => setTimeout(window.stub, 1000)")
	must.NoError(err)
	id, ok := raw.(float64)
	is.True(ok && id > 0, "setTimeout should return a positive ID, got %v", raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockReplacesClearTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => { const to = setTimeout(window.stub, 1000); clearTimeout(to); }")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 1000))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(0), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockReplacesClearInterval(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => { const to = setInterval(window.stub, 500); clearInterval(to); }")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 1000))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(0), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockSetFixedTimeMultipleTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bCtx.Clock().SetFixedTime(ctx, 100))
	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	is.Equal(float64(100), raw)

	must.NoError(bCtx.Clock().SetFixedTime(ctx, 200))
	raw, err = page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	is.Equal(float64(200), raw)
}

// Extra: no equivalent method in playwright-java; related class: TestPageClock.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestClockSetFixedTimeAllowsFakeTimers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._timestamps = []; window.stub = function(ts) { window._timestamps.push(ts); }; }")
	must.NoError(err)

	must.NoError(bCtx.Clock().SetFixedTime(ctx, 100))
	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	is.Equal(float64(100), raw)

	must.NoError(bCtx.Clock().SetFixedTime(ctx, 200))
	_, err = page.Evaluate(ctx, "() => { setTimeout(() => window.stub(Date.now()), 0); }")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 0))

	raw, err = page.Evaluate(ctx, "() => window._timestamps[0]")
	must.NoError(err)
	is.Equal(float64(200), raw)
}

// Ref: TestPageClock.java#setsInitialTimestamp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockSetsInitialTimestamp(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))
	must.NoError(bCtx.Clock().SetSystemTime(ctx, 1400))

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	is.Equal(float64(1400), raw)
}

// Ref: TestPageClock.java#replacesGlobalSetTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockReplacesGlobalSetTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setTimeout(window.stub, 1000)")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 1000))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(1), raw)
}

// Ref: TestPageClock.java#replacesGlobalSetInterval  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockReplacesGlobalSetInterval(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	_, err = page.Evaluate(ctx, "() => setInterval(window.stub, 500)")
	must.NoError(err)

	must.NoError(bCtx.Clock().RunFor(ctx, 1000))

	raw, err := page.Evaluate(ctx, "() => window._c")
	must.NoError(err)
	is.Equal(float64(2), raw)
}

// Ref: TestPageClock.java#fakesDateConstructor  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockFakesDateConstructor(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => { window._c = 0; window.stub = function() { window._c++; }; }")
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	raw, err := page.Evaluate(ctx, "() => new Date().getTime()")
	must.NoError(err)
	is.Equal(float64(1000), raw)
}

// Ref: TestPageClock.java#setFixedTimeDoesNotFakeMethods  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockSetFixedTimeDoesNotFakeMethods(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bCtx.Clock().SetFixedTime(ctx, 0))

	_, err = page.Evaluate(ctx, "() => new Promise(f => setTimeout(f, 1))")
	must.NoError(err)
}

// Ref: TestPageClock.java#setFixedTimeFixedTimeIsNotAffectedByClockManipulation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockSetFixedTimeFixedTimeIsNotAffectedByClockManipulation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bCtx.Clock().SetFixedTime(ctx, 100))
	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	is.Equal(float64(100), raw)

	must.NoError(bCtx.Clock().FastForward(ctx, 20))

	raw, err = page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	is.Equal(float64(100), raw, "setFixedTime should not be affected by fastForward")
}

// Ref: TestPageClock.java#whileRunningShouldProgressTime  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileRunningShouldProgressTime(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))

	_ = page.WaitForTimeout(ctx, 1000)

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	now, _ := raw.(float64)
	is.GreaterOrEqual(now, float64(1000))
	is.LessOrEqual(now, float64(2000))
}

// Ref: TestPageClock.java#whileRunningShouldRunFor  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileRunningShouldRunFor(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))

	must.NoError(bCtx.Clock().RunFor(ctx, 10000))

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	now, _ := raw.(float64)
	is.GreaterOrEqual(now, float64(10000))
	is.LessOrEqual(now, float64(11000))
}

// Ref: TestPageClock.java#whileRunningShouldFastForward  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileRunningShouldFastForward(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))

	must.NoError(bCtx.Clock().FastForward(ctx, 10000))

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	now, _ := raw.(float64)
	is.GreaterOrEqual(now, float64(10000))
	is.LessOrEqual(now, float64(11000))
}

// Ref: TestPageClock.java#whileRunningShouldPause  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileRunningShouldPause(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(bCtx.Clock().Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))

	must.NoError(bCtx.Clock().PauseAt(ctx, 1000))

	time.Sleep(300 * time.Millisecond)

	raw, err := page.Evaluate(ctx, "() => Date.now()")
	must.NoError(err)
	now, _ := raw.(float64)
	is.GreaterOrEqual(now, float64(0))
	is.LessOrEqual(now, float64(1000))
}

// Ref: TestPageClock.java#runForTriggersImmediatelyWithoutSpecifiedDelay  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForTriggersImmediatelyWithoutSpecifiedDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setTimeout(window.stub)")
	must.NoError(clock.RunFor(ctx, 0))
	calls.WaitLen(t, 1)
}

// Ref: TestPageClock.java#runForDoesNotTriggerWithoutSufficientDelay  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForDoesNotTriggerWithoutSufficientDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setTimeout(window.stub, 100)")
	must.NoError(clock.RunFor(ctx, 11))
	calls.WaitLen(t, 0)
}

// Ref: TestPageClock.java#runForTriggersAfterSufficientDelay  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForTriggersAfterSufficientDelay(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setTimeout(window.stub, 100)")
	must.NoError(clock.RunFor(ctx, 100))
	calls.WaitLen(t, 1)
}

// Ref: TestPageClock.java#runForTriggersSimultaneousTimers  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForTriggersSimultaneousTimers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => { setTimeout(window.stub, 100);setTimeout(window.stub, 100); }")
	must.NoError(clock.RunFor(ctx, 100))
	calls.WaitLen(t, 2)
}

// Ref: TestPageClock.java#runForTriggersMultipleSimultaneousTimers  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForTriggersMultipleSimultaneousTimers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => { setTimeout(window.stub, 100);setTimeout(window.stub, 100);setTimeout(window.stub, 99);setTimeout(window.stub, 100); }")
	must.NoError(clock.RunFor(ctx, 100))
	calls.WaitLen(t, 4)
}

// Ref: TestPageClock.java#runForWaitsAfterSetTimeoutWasCalled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForWaitsAfterSetTimeoutWasCalled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setTimeout(window.stub, 150)")
	must.NoError(clock.RunFor(ctx, 50))
	calls.WaitLen(t, 0)
	must.NoError(clock.RunFor(ctx, 100))
	calls.WaitLen(t, 1)
}

// Ref: TestPageClock.java#runForTriggersEventWhenSomeThrow  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForTriggersEventWhenSomeThrow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, `() => {
  setTimeout(() => {throw new Error(); }, 100);
  setTimeout(window.stub, 120);
}`)
	must.Error(clock.RunFor(ctx, 120))
	calls.WaitLen(t, 1)
}

// Ref: TestPageClock.java#runForCreatesUpdatedDateWhileTicking  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForCreatesUpdatedDateWhileTicking(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	must.NoError(clock.SetSystemTime(ctx, 0))
	a2Ev(t, ctx, page, "() => setInterval(() => { window.stub(new Date().getTime()); }, 10)")
	must.NoError(clock.RunFor(ctx, 100))
	calls.WaitLen(t, 10)
	var expected [][]any
	for i := 1; i <= 10; i++ {
		expected = append(expected, []any{float64(i * 10)})
	}
	assert.Equal(t, expected, calls.All())
}

// Ref: TestPageClock.java#runForPasses8Seconds  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForPasses8Seconds(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setInterval(window.stub, 4000)")
	must.NoError(clock.RunFor(ctx, "08"))
	calls.WaitLen(t, 2)
}

// Ref: TestPageClock.java#runForPasses1Minute  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForPasses1Minute(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setInterval(window.stub, 6000)")
	must.NoError(clock.RunFor(ctx, "01:00"))
	calls.WaitLen(t, 10)
}

// Ref: TestPageClock.java#runForPasses2Hours34MinutesAnd10Seconds  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForPasses2Hours34MinutesAnd10Seconds(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setInterval(window.stub, 10000)")
	must.NoError(clock.RunFor(ctx, "02:34:10"))
	calls.WaitLen(t, 925)
}

// Ref: TestPageClock.java#runForThrowsForInvalidFormat  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForThrowsForInvalidFormat(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setInterval(window.stub, 10000)")
	must.Error(clock.RunFor(ctx, "12:02:34:10"))
	calls.WaitLen(t, 0)
}

// Ref: TestPageClock.java#runForReturnsTheCurrentNowValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockRunForReturnsTheCurrentNowValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, _, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	must.NoError(clock.SetSystemTime(ctx, 0))
	must.NoError(clock.RunFor(ctx, 200))
	assert.Equal(t, float64(200), a2Ev(t, ctx, page, "() => Date.now()"))
}

// Ref: TestPageClock.java#fastForwardIgnoresTimersWhichWouldntBeRun  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockFastForwardIgnoresTimersWhichWouldntBeRun(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setTimeout(() => { window.stub('should not be logged'); }, 1000)")
	must.NoError(clock.FastForward(ctx, 500))
	calls.WaitLen(t, 0)
}

// Ref: TestPageClock.java#fastForwardPushesBackExecutionTimeForSkippedTimers  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockFastForwardPushesBackExecutionTimeForSkippedTimers(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setTimeout(() => { window.stub(Date.now()); }, 1000)")
	must.NoError(clock.FastForward(ctx, 2000))
	calls.WaitLen(t, 1)
	assert.Equal(t, float64(1000+2000), calls.Get(0)[0])
}

// Ref: TestPageClock.java#fastForwardSupportsStringTimeArguments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockFastForwardSupportsStringTimeArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => setTimeout(() => { window.stub(Date.now()); }, 100000)")
	must.NoError(clock.FastForward(ctx, "01:50"))
	calls.WaitLen(t, 1)
	assert.Equal(t, float64(1000+110000), calls.Get(0)[0])
}

// Ref: TestPageClock.java#globalFakeSetTimeoutShouldReturnId  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockGlobalFakeSetTimeoutShouldReturnId(t *testing.T) {
	t.Parallel()
	ctx, _, page, _, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	to := a2Ev(t, ctx, page, "() => setTimeout(window.stub, 1000)")
	assert.Greater(t, to.(float64), float64(0))
}

// Ref: TestPageClock.java#replacesGlobalClearTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockReplacesGlobalClearTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => { const to = setTimeout(window.stub, 1000); clearTimeout(to); }")
	must.NoError(clock.RunFor(ctx, 1000))
	calls.WaitLen(t, 0)
}

// Ref: TestPageClock.java#replacesGlobalClearInterval  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockReplacesGlobalClearInterval(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	a2ClockPaused(t, ctx, clock)
	a2Ev(t, ctx, page, "() => { const to = setInterval(window.stub, 500); clearInterval(to); }")
	must.NoError(clock.RunFor(ctx, 1000))
	calls.WaitLen(t, 0)
}

// Ref: TestPageClock.java#replacesGlobalPerformanceNow  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockReplacesGlobalPerformanceNow(t *testing.T) {
	t.Skip("not implemented: test body is empty in the Java reference")
}

// Ref: TestPageClock.java#replacesGlobalPerformanceTimeOrigin  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockReplacesGlobalPerformanceTimeOrigin(t *testing.T) {
	t.Skip("not implemented: test body is empty in the Java reference")
}

func a2PopupDate() float64 {
	return float64(time.Date(2015, 9, 25, 0, 0, 0, 0, time.Local).UnixMilli())
}

// Ref: TestPageClock.java#shouldTickAfterPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockShouldTickAfterPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, _, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	now := a2PopupDate()
	must.NoError(clock.PauseAt(ctx, now))
	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.open('about:blank')")
		return err
	})
	must.NoError(err)
	assert.Equal(t, now, a2Ev(t, ctx, popup, "() => Date.now()"))
	must.NoError(clock.RunFor(ctx, 1000))
	assert.Equal(t, now+1000, a2Ev(t, ctx, popup, "() => Date.now()"))
}

// Ref: TestPageClock.java#shouldTickBeforePopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockShouldTickBeforePopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, _, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	now := a2PopupDate()
	must.NoError(clock.PauseAt(ctx, now))
	must.NoError(clock.RunFor(ctx, 1000))
	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "() => window.open('about:blank')")
		return err
	})
	must.NoError(err)
	assert.Equal(t, now+1000, a2Ev(t, ctx, popup, "() => Date.now()"))
}

// Ref: TestPageClock.java#shouldRunTimeBeforePopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockShouldRunTimeBeforePopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page, _, _ := a2ClockSetup(t)

	srv.ServeWithBody("/popup.html", "text/html", "<script>window.time = Date.now()</script>")
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForTimeout(ctx, 2000))
	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "url => window.open(url)", srv.Prefix()+"/popup.html")
		return err
	})
	must.NoError(err)
	must.NoError(popup.WaitForURL(ctx, srv.Prefix()+"/popup.html"))
	popupTime := a2Ev(t, ctx, popup, "time").(float64)
	assert.GreaterOrEqual(t, popupTime, float64(2000), "popupTime = %v", popupTime)
}

// Ref: TestPageClock.java#shouldNotRunTimeBeforePopupOnPause  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockShouldNotRunTimeBeforePopupOnPause(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, srv, page, _, clock := a2ClockSetup(t)

	srv.ServeWithBody("/popup.html", "text/html", "<script>window.time = Date.now()</script>")
	a2ClockPaused(t, ctx, clock)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.WaitForTimeout(ctx, 2000))
	popup, err := page.WaitForPopup(ctx, func() error {
		_, err := page.Evaluate(ctx, "url => window.open(url)", srv.Prefix()+"/popup.html")
		return err
	})
	must.NoError(err)
	must.NoError(popup.WaitForURL(ctx, srv.Prefix()+"/popup.html"))
	assert.Equal(t, float64(1000), a2Ev(t, ctx, popup, "time"))
}

// Ref: TestPageClock.java#setFixedTimeAllowsSettingTimeMultipleTimes  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockSetFixedTimeAllowsSettingTimeMultipleTimes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, _, clock := a2ClockSetup(t)

	must.NoError(clock.SetFixedTime(ctx, 100))
	assert.Equal(t, float64(100), a2Ev(t, ctx, page, "() => Date.now()"))
	must.NoError(clock.SetFixedTime(ctx, 200))
	assert.Equal(t, float64(200), a2Ev(t, ctx, page, "() => Date.now()"))
}

// Ref: TestPageClock.java#setFixedTimeAllowsInstallingFakeTimersAfterSettingTime  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockSetFixedTimeAllowsInstallingFakeTimersAfterSettingTime(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	must.NoError(clock.SetFixedTime(ctx, 100))
	assert.Equal(t, float64(100), a2Ev(t, ctx, page, "() => Date.now()"))
	must.NoError(clock.SetFixedTime(ctx, 200))
	a2Ev(t, ctx, page, "async () => { setTimeout(() => window.stub(Date.now()), 0); }")
	must.NoError(clock.RunFor(ctx, 0))
	calls.WaitLen(t, 1)
	assert.Equal(t, float64(200), calls.Get(0)[0])
}

// Ref: TestPageClock.java#whileRunningShouldFastForwardTo  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileRunningShouldFastForwardTo(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, _, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))
	must.NoError(clock.FastForward(ctx, 10000))
	now := a2Ev(t, ctx, page, "() => Date.now()").(float64)
	assert.True(t, now >= 10000 && now <= 11000, "now = %v", now)
}

// Ref: TestPageClock.java#whileRunningShouldPauseAndFastForward  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileRunningShouldPauseAndFastForward(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, _, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))
	must.NoError(clock.PauseAt(ctx, 1000))
	must.NoError(clock.FastForward(ctx, 1000))
	assert.Equal(t, float64(2000), a2Ev(t, ctx, page, "() => Date.now()"))
}

// Ref: TestPageClock.java#whileRunningShouldSetSystemTimeOnPause  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileRunningShouldSetSystemTimeOnPause(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, _, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))
	must.NoError(clock.PauseAt(ctx, 1000))
	assert.Equal(t, float64(1000), a2Ev(t, ctx, page, "() => Date.now()"))
}

// Ref: TestPageClock.java#whileOnPauseFastForwardShouldNotRunNestedImmediate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileOnPauseFastForwardShouldNotRunNestedImmediate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))
	must.NoError(clock.PauseAt(ctx, 1000))
	a2Ev(t, ctx, page, "() => { setTimeout(() => { window.stub('outer'); setTimeout(() => window.stub('inner'), 0); }, 1000); }")
	must.NoError(clock.FastForward(ctx, 1000))
	calls.WaitLen(t, 1)
	assert.Equal(t, "outer", calls.Get(0)[0])
	must.NoError(clock.FastForward(ctx, 1))
	calls.WaitLen(t, 2)
	assert.Equal(t, "inner", calls.Get(1)[0])
}

// Ref: TestPageClock.java#whileOnPauseRunForShouldNotRunNestedImmediate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileOnPauseRunForShouldNotRunNestedImmediate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))
	must.NoError(clock.PauseAt(ctx, 1000))
	a2Ev(t, ctx, page, "() => { setTimeout(() => { window.stub('outer'); setTimeout(() => window.stub('inner'), 0); }, 1000); }")
	must.NoError(clock.RunFor(ctx, 1000))
	calls.WaitLen(t, 1)
	assert.Equal(t, "outer", calls.Get(0)[0])
	must.NoError(clock.RunFor(ctx, 1))
	calls.WaitLen(t, 2)
	assert.Equal(t, "inner", calls.Get(1)[0])
}

// Ref: TestPageClock.java#whileOnPauseRunForShouldNotRunNestedImmediateFromMicrotask  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockWhileOnPauseRunForShouldNotRunNestedImmediateFromMicrotask(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page, calls, clock := a2ClockSetup(t)

	must.NoError(clock.Install(ctx, 0))
	must.NoError(page.Goto(ctx, "data:text/html,"))
	must.NoError(clock.PauseAt(ctx, 1000))
	a2Ev(t, ctx, page, "() => { setTimeout(() => { window.stub('outer'); void Promise.resolve().then(() => setTimeout(() => window.stub('inner'), 0)); }, 1000); }")
	must.NoError(clock.RunFor(ctx, 1000))
	calls.WaitLen(t, 1)
	assert.Equal(t, "outer", calls.Get(0)[0])
	must.NoError(clock.RunFor(ctx, 1))
	calls.WaitLen(t, 2)
	assert.Equal(t, "inner", calls.Get(1)[0])
}

// Ref: TestPageClock.java#shouldThrowForInvalidDate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageClock.java
func TestPageClockShouldThrowForInvalidDate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, _, _, clock := a2ClockSetup(t)

	err := clock.SetSystemTime(ctx, "invalid")
	must.Error(err)
	assert.Contains(t, err.Error(), "Invalid date: invalid")
}
