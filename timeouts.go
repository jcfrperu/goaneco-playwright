package playwright

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"
)

// Timeout constants used across IPC calls and subscription helpers.
const (
	// defaultActionTimeoutMs is a sentinel stored in request params when the caller did not give an
	// explicit timeout. ChannelOwner.sendMessage replaces it with the effective timeout (value set via
	// SetDefaultTimeout / SetDefaultNavigationTimeout on the page or context, or fallbackTimeoutMs).
	defaultActionTimeoutMs = -1.0

	// fallbackTimeoutMs is Playwright's built-in default timeout (ms) when none was configured.
	fallbackTimeoutMs = 30000.0

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

// timeoutScope holds the user-configured default timeouts of a Page, BrowserContext or Frame.
// A scope without values defers to its parent.
type timeoutScope struct {
	mu         sync.RWMutex
	timeout    *float64
	navTimeout *float64
	parentGUID string
}

// timeoutScopes maps object GUIDs to their scope. Entries are only created when the user calls
// SetDefaultTimeout / SetDefaultNavigationTimeout, so the map stays empty in the common case.
var timeoutScopes sync.Map // string -> *timeoutScope

func getTimeoutScope(guid string) *timeoutScope {
	if v, ok := timeoutScopes.Load(guid); ok {
		return v.(*timeoutScope)
	}
	return nil
}

func ensureTimeoutScope(guid, parentGUID string) *timeoutScope {
	v, _ := timeoutScopes.LoadOrStore(guid, &timeoutScope{parentGUID: parentGUID})
	s := v.(*timeoutScope)
	s.mu.Lock()
	if s.parentGUID == "" {
		s.parentGUID = parentGUID
	}
	s.mu.Unlock()
	return s
}

func (s *timeoutScope) setTimeout(ms float64) {
	s.mu.Lock()
	s.timeout = &ms
	s.mu.Unlock()
}

func (s *timeoutScope) setNavTimeout(ms float64) {
	s.mu.Lock()
	s.navTimeout = &ms
	s.mu.Unlock()
}

// navigationMethods are the protocol methods governed by the default navigation timeout.
var navigationMethods = map[string]bool{
	"goto": true, "reload": true, "goBack": true, "goForward": true,
	"setContent": true, "waitForNavigation": true, "waitForURL": true,
}

// effectiveTimeout resolves the timeout (ms) for a method sent to this object: it walks from the
// object up through its parents (sub-frame -> parent frame -> main frame -> page -> context) and
// returns the first configured value, falling back to Playwright's 30s default.
func (c *ChannelOwner) effectiveTimeout(method string) float64 {
	nav := navigationMethods[method]
	guid := c.guid
	for depth := 0; depth < 16 && guid != ""; depth++ {
		next := ""
		if s := getTimeoutScope(guid); s != nil {
			s.mu.RLock()
			if nav && s.navTimeout != nil {
				v := *s.navTimeout
				s.mu.RUnlock()
				return v
			}
			if s.timeout != nil {
				v := *s.timeout
				s.mu.RUnlock()
				return v
			}
			next = s.parentGUID
			s.mu.RUnlock()
		}
		if next == "" {
			next = c.parentFrameGUID(guid)
		}
		guid = next
	}
	return fallbackTimeoutMs
}

// parentFrameGUID returns the parent frame GUID from a Frame's initializer, or "" otherwise.
func (c *ChannelOwner) parentFrameGUID(guid string) string {
	if c.conn == nil {
		return ""
	}
	raw := c.Initializer(guid)
	if len(raw) == 0 {
		return ""
	}
	var init struct {
		ParentFrame *struct {
			GUID string `json:"guid"`
		} `json:"parentFrame"`
	}
	if err := json.Unmarshal(raw, &init); err != nil || init.ParentFrame == nil {
		return ""
	}
	return init.ParentFrame.GUID
}

var (
	unsetTimeoutMid  = []byte(`"timeout":-1,`)
	unsetTimeoutLast = []byte(`"timeout":-1}`)
)

// resolveUnsetTimeouts replaces the "timeout":-1 sentinel in marshaled params with the effective timeout.
func (c *ChannelOwner) resolveUnsetTimeouts(method string, params []byte) []byte {
	if !bytes.Contains(params, unsetTimeoutMid) && !bytes.Contains(params, unsetTimeoutLast) {
		return params
	}
	ms, _ := json.Marshal(c.effectiveTimeout(method))
	params = bytes.ReplaceAll(params, unsetTimeoutMid, append(append([]byte(`"timeout":`), ms...), ','))
	params = bytes.ReplaceAll(params, unsetTimeoutLast, append(append([]byte(`"timeout":`), ms...), '}'))
	return params
}
