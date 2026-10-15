package playwright

import (
	"context"
	"fmt"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// Clock provides control over the browser context's fake clock, enabling tests to
// manipulate time-dependent logic without waiting for real time to pass.
// Obtained via BrowserContext.Clock() or Page.Clock().
type Clock struct {
	context *BrowserContext
}

// Clock returns a Clock instance for this browser context.
func (c *BrowserContext) Clock() *Clock {
	return &Clock{context: c}
}

// Clock returns the Clock instance for this page's browser context.
func (p *Page) Clock() *Clock {
	if p.browserContext != nil {
		return p.browserContext.Clock()
	}
	return nil
}

func parseTimeArg(val any) (*float64, *string, error) {
	switch v := val.(type) {
	case string:
		return nil, &v, nil
	case float64:
		return protocol.Float(v), nil, nil
	case float32:
		return protocol.Float(float64(v)), nil, nil
	case int:
		return protocol.Float(float64(v)), nil, nil
	case int64:
		return protocol.Float(float64(v)), nil, nil
	case int32:
		return protocol.Float(float64(v)), nil, nil
	default:
		return nil, nil, fmt.Errorf("unsupported time value: %v (type %T)", val, val)
	}
}

// Install replaces the native browser clock with a fake one.
// timeVal is the initial time (milliseconds since epoch or string);
// omit it to use the current real time (an explicit 0 means the Unix epoch).
func (c *Clock) Install(ctx context.Context, timeVal ...any) error {
	req := protocol.BrowserContextClockInstallRequest{}
	if len(timeVal) > 0 && timeVal[0] != nil {
		num, str, err := parseTimeArg(timeVal[0])
		if err != nil {
			return fmt.Errorf("clock.install failed: %w", err)
		}
		req.TimeNumber = num
		req.TimeString = str
	}
	_, err := c.context.owner.SendMessageRequest(ctx, "clockInstall", req)
	if err != nil {
		return fmt.Errorf("clock.install failed: %w", err)
	}
	return nil
}

// FastForward advances the fake clock by the given number of milliseconds or string time (e.g. "01:50"),
// firing all timers that expire during that period without blocking.
func (c *Clock) FastForward(ctx context.Context, ticks any) error {
	num, str, err := parseTimeArg(ticks)
	if err != nil {
		return fmt.Errorf("clock.fastForward failed: %w", err)
	}
	req := protocol.BrowserContextClockFastForwardRequest{
		TicksNumber: num,
		TicksString: str,
	}
	_, err = c.context.owner.SendMessageRequest(ctx, "clockFastForward", req)
	if err != nil {
		return fmt.Errorf("clock.fastForward failed: %w", err)
	}
	return nil
}

// RunFor simulates the passage of time by the given milliseconds or string time (e.g. "01:00"),
// running all timers that fire during that period.
func (c *Clock) RunFor(ctx context.Context, ticks any) error {
	num, str, err := parseTimeArg(ticks)
	if err != nil {
		return fmt.Errorf("clock.runFor failed: %w", err)
	}
	req := protocol.BrowserContextClockRunForRequest{
		TicksNumber: num,
		TicksString: str,
	}
	_, err = c.context.owner.SendMessageRequest(ctx, "clockRunFor", req)
	if err != nil {
		return fmt.Errorf("clock.runFor failed: %w", err)
	}
	return nil
}

// PauseAt pauses the fake clock at the given time (milliseconds since epoch or string).
// No timers fire until Resume is called.
func (c *Clock) PauseAt(ctx context.Context, timeVal any) error {
	num, str, err := parseTimeArg(timeVal)
	if err != nil {
		return fmt.Errorf("clock.pauseAt failed: %w", err)
	}
	req := protocol.BrowserContextClockPauseAtRequest{
		TimeNumber: num,
		TimeString: str,
	}
	_, err = c.context.owner.SendMessageRequest(ctx, "clockPauseAt", req)
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
// (milliseconds since epoch or string) without installing a fake clock for timers.
func (c *Clock) SetFixedTime(ctx context.Context, timeVal any) error {
	num, str, err := parseTimeArg(timeVal)
	if err != nil {
		return fmt.Errorf("clock.setFixedTime failed: %w", err)
	}
	req := protocol.BrowserContextClockSetFixedTimeRequest{
		TimeNumber: num,
		TimeString: str,
	}
	_, err = c.context.owner.SendMessageRequest(ctx, "clockSetFixedTime", req)
	if err != nil {
		return fmt.Errorf("clock.setFixedTime failed: %w", err)
	}
	return nil
}

// SetSystemTime changes the system clock value to the given time (milliseconds since epoch or string)
// without changing how timers fire.
func (c *Clock) SetSystemTime(ctx context.Context, timeVal any) error {
	num, str, err := parseTimeArg(timeVal)
	if err != nil {
		return fmt.Errorf("clock.setSystemTime failed: %w", err)
	}
	req := protocol.BrowserContextClockSetSystemTimeRequest{
		TimeNumber: num,
		TimeString: str,
	}
	_, err = c.context.owner.SendMessageRequest(ctx, "clockSetSystemTime", req)
	if err != nil {
		return fmt.Errorf("clock.setSystemTime failed: %w", err)
	}
	return nil
}
