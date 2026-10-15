package playwright

import "log/slog"

// safeGo runs fn in a new goroutine and recovers from panics, logging them instead of
// crashing the whole process (and the user's test run) from a user-supplied callback.
func safeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Default().Error("goaneco-playwright: panic in user handler", "panic", r)
			}
		}()
		fn()
	}()
}
