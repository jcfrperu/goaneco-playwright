package playwright

import "time"

// Timeout constants used across IPC calls and subscription helpers.
const (
	// defaultActionTimeoutMs is the default timeout (ms) sent to Playwright server for most actions.
	defaultActionTimeoutMs = 30000.0

	// defaultLaunchTimeoutMs is the default timeout (ms) for browser launch operations.
	// Launch can be slow on resource-constrained systems so it uses a longer window than actions.
	defaultLaunchTimeoutMs = 120000.0

	// defaultBindingHandlerTimeoutMs is the timeout (ms) for user-provided expose-binding handlers.
	// Separate from defaultActionTimeoutMs so they can be tuned independently.
	defaultBindingHandlerTimeoutMs = 30000.0

	// defaultSubscriptionTimeout is the context timeout used when registering event subscriptions.
	defaultSubscriptionTimeout = 5 * time.Second

	// defaultBrowserCleanupTimeout is the context timeout used during browser cleanup on NewPage.
	defaultBrowserCleanupTimeout = 3 * time.Second
)
