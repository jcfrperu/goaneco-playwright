package playwright

import (
	"context"
	"fmt"
)

// waitForOne subscribes with subscribe, runs the optional trigger, and waits for the
// first value delivered or for ctx to end. The subscription is always cancelled on return.
// subscribe receives a non-blocking send function and returns the cancel function.
func waitForOne[T any](ctx context.Context, name string, subscribe func(send func(T)) func(), trigger func() error) (T, error) {
	var zero T
	ch := make(chan T, 1)
	unsub := subscribe(func(v T) {
		select {
		case ch <- v:
		default:
		}
	})
	defer unsub()

	if trigger != nil {
		if err := trigger(); err != nil {
			return zero, fmt.Errorf("%s: trigger failed: %w", name, err)
		}
	}

	select {
	case v := <-ch:
		return v, nil
	case <-ctx.Done():
		return zero, fmt.Errorf("%s: %w", name, ctx.Err())
	}
}
