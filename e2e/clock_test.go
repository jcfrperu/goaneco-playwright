//go:build e2e

// E2E tests for BrowserContext.Clock() — fake timer control.
package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClockSetFixedTime verifies that SetFixedTime fixes Date.now() to the given value.
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

func TestClockSetsInitialTimestamp(t *testing.T) {
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

func TestClockReplacesGlobalSetTimeout(t *testing.T) {
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

func TestClockReplacesGlobalSetInterval(t *testing.T) {
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

func TestClockFakesDateConstructor(t *testing.T) {
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

func TestClockSetFixedTimeDoesNotFakeMethods(t *testing.T) {
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

func TestClockSetFixedTimeNotAffectedByManipulation(t *testing.T) {
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

func TestClockWhileRunningShouldProgressTime(t *testing.T) {
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

func TestClockWhileRunningShouldRunFor(t *testing.T) {
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

func TestClockWhileRunningShouldFastForward(t *testing.T) {
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

func TestClockWhileRunningShouldPause(t *testing.T) {
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
