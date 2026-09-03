package playwright

import (
	"context"
	"fmt"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// Clock provides control over the browser context's fake clock, enabling tests to
// manipulate time-dependent logic without waiting for real time to pass.
// Obtained via BrowserContext.Clock().
type Clock struct {
	context *BrowserContext
}

// Clock returns a Clock instance for this browser context.
func (c *BrowserContext) Clock() *Clock {
	return &Clock{context: c}
}

// Install replaces the native browser clock with a fake one.
// timeMs is the initial time in milliseconds since the Unix epoch;
// omit or pass 0 to use the current real time.
func (c *Clock) Install(ctx context.Context, timeMs ...float64) error {
	req := protocol.BrowserContextClockInstallRequest{}
	if len(timeMs) > 0 && timeMs[0] != 0 {
		req.TimeNumber = protocol.Float(timeMs[0])
	}
	_, err := c.context.owner.SendMessageRequest(ctx, "clockInstall", req)
	if err != nil {
		return fmt.Errorf("clock.install failed: %w", err)
	}
	return nil
}

// FastForward advances the fake clock by the given number of milliseconds,
// firing all timers that expire during that period without blocking.
func (c *Clock) FastForward(ctx context.Context, ms float64) error {
	req := protocol.BrowserContextClockFastForwardRequest{
		TicksNumber: protocol.Float(ms),
	}
	_, err := c.context.owner.SendMessageRequest(ctx, "clockFastForward", req)
	if err != nil {
		return fmt.Errorf("clock.fastForward failed: %w", err)
	}
	return nil
}

// RunFor simulates the passage of time by the given milliseconds,
// running all timers that fire during that period.
func (c *Clock) RunFor(ctx context.Context, ms float64) error {
	req := protocol.BrowserContextClockRunForRequest{
		TicksNumber: protocol.Float(ms),
	}
	_, err := c.context.owner.SendMessageRequest(ctx, "clockRunFor", req)
	if err != nil {
		return fmt.Errorf("clock.runFor failed: %w", err)
	}
	return nil
}

// PauseAt pauses the fake clock at the given time (milliseconds since epoch).
// No timers fire until Resume is called.
func (c *Clock) PauseAt(ctx context.Context, timeMs float64) error {
	req := protocol.BrowserContextClockPauseAtRequest{
		TimeNumber: protocol.Float(timeMs),
	}
	_, err := c.context.owner.SendMessageRequest(ctx, "clockPauseAt", req)
	if err != nil {
		return fmt.Errorf("clock.pauseAt failed: %w", err)
	}
	return nil
}

// Resume resumes a previously paused clock, allowing timers to fire normally.
func (c *Clock) Resume(ctx context.Context) error {
	_, err := c.context.owner.SendMessageRequest(ctx, "clockResume", protocol.BrowserContextClockResumeRequest{})
	if err != nil {
		return fmt.Errorf("clock.resume failed: %w", err)
	}
	return nil
}

// SetFixedTime fixes both Date.now() and new Date() to return the given time
// (milliseconds since epoch) without installing a fake clock for timers.
func (c *Clock) SetFixedTime(ctx context.Context, timeMs float64) error {
	req := protocol.BrowserContextClockSetFixedTimeRequest{
		TimeNumber: protocol.Float(timeMs),
	}
	_, err := c.context.owner.SendMessageRequest(ctx, "clockSetFixedTime", req)
	if err != nil {
		return fmt.Errorf("clock.setFixedTime failed: %w", err)
	}
	return nil
}

// SetSystemTime changes the system clock value to the given time (milliseconds since epoch)
// without changing how timers fire.
func (c *Clock) SetSystemTime(ctx context.Context, timeMs float64) error {
	req := protocol.BrowserContextClockSetSystemTimeRequest{
		TimeNumber: protocol.Float(timeMs),
	}
	_, err := c.context.owner.SendMessageRequest(ctx, "clockSetSystemTime", req)
	if err != nil {
		return fmt.Errorf("clock.setSystemTime failed: %w", err)
	}
	return nil
}
