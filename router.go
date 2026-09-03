package playwright

import (
	"context"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// routerEntry is a single registered route handler with optional call-count limit.
type routerEntry struct {
	pattern         string
	compiledPattern *regexp.Regexp // cached compile of pattern; nil if pattern is empty or invalid
	handler         RouteHandler
	times           int32 // 0 = unlimited; positive = max calls allowed
	count           int32 // atomic call counter
}

// routeRouter manages a LIFO stack of route handlers and dispatches them sequentially,
// supporting fallback chaining: a handler may call Route.Fallback to defer to the next handler.
// When all handlers fall back, onAllFallback is called if set; otherwise continueAsFallback is used.
type routeRouter struct {
	mu             sync.Mutex
	entries        []*routerEntry
	onAllFallback  func(*Route)  // optional: called when every registered handler falls back
	handlerTimeout time.Duration // max time to wait for a handler to resolve; 0 uses defaultHandlerTimeout
}

const defaultHandlerTimeout = 30 * time.Second

// add registers a new handler (appended; dispatch iterates in reverse for LIFO order).
// times=0 means unlimited calls.
func (r *routeRouter) add(pattern string, handler RouteHandler, times int) {
	entry := &routerEntry{pattern: pattern, handler: handler}
	if pattern != "" {
		if compiled, err := buildGlobRegexp(pattern); err == nil {
			entry.compiledPattern = compiled
		}
	}
	if times > 0 {
		entry.times = int32(times)
	}
	r.mu.Lock()
	r.entries = append(r.entries, entry)
	r.mu.Unlock()
}

// clear removes all handlers.
func (r *routeRouter) clear() {
	r.mu.Lock()
	r.entries = nil
	r.mu.Unlock()
}

// hasHandlers reports whether at least one handler is registered.
func (r *routeRouter) hasHandlers() bool {
	r.mu.Lock()
	n := len(r.entries)
	r.mu.Unlock()
	return n > 0
}

// dispatch runs the handler chain for the given route in a new goroutine.
// Handlers are tried in LIFO order. If a handler calls Route.Fallback, the next handler is tried.
// If all handlers fall back (or none is registered), the route is passed to the next routing layer
// via "continue" with isFallback=true.
func (r *routeRouter) dispatch(route *Route) {
	r.mu.Lock()
	entries := make([]*routerEntry, len(r.entries))
	copy(entries, r.entries)
	r.mu.Unlock()

	timeout := r.handlerTimeout
	if timeout <= 0 {
		timeout = defaultHandlerTimeout
	}

	go func() {
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			// Enforce call-count limit.
			if entry.times > 0 {
				n := atomic.AddInt32(&entry.count, 1)
				if n > entry.times {
					// Entry exhausted — remove it from the router to prevent unbounded growth.
					r.mu.Lock()
					for j, e := range r.entries {
						if e == entry {
							r.entries = append(r.entries[:j], r.entries[j+1:]...)
							break
						}
					}
					r.mu.Unlock()
					continue
				}
			}

			// Skip if the request URL does not match this entry's pattern.
			if entry.pattern != "" {
				req := route.Request()
				if req != nil {
					var matches bool
					if entry.compiledPattern != nil {
						matches = entry.compiledPattern.MatchString(req.URL())
					} else {
						matches = matchGlob(entry.pattern, req.URL())
					}
					if !matches {
						continue
					}
				}
			}

			// Arm the route with a done channel so the handler can signal completion.
			doneCh := make(chan bool, 1)
			route.setDoneCh(doneCh)

			// Run the handler; panic-safe.
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						slog.Error("goaneco-playwright: route handler panic", "panic", rec)
						// Signal "handled" so the dispatch loop doesn't hang.
						select {
						case doneCh <- false:
						default:
						}
					}
				}()
				entry.handler(route)
			}()

			// Wait for the handler to either handle the route or call Fallback.
			handlerTimer := time.NewTimer(timeout)
			select {
			case isFallback := <-doneCh:
				handlerTimer.Stop()
				if !isFallback {
					return // Route was handled (Fulfill/Continue/Abort). Stop chain.
				}
				// Handler called Fallback — continue to the next entry.
			case <-handlerTimer.C:
				slog.Error("goaneco-playwright: route handler did not resolve within timeout", "timeout", timeout)
				return
			}
		}

		// All handlers fell back or no handlers matched.
		r.mu.Lock()
		fb := r.onAllFallback
		r.mu.Unlock()
		if fb != nil {
			// Delegate to the next routing layer (e.g., page → context).
			fb(route)
			return
		}
		// Default: pass to the network via isFallback=true.
		if err := route.continueAsFallback(context.Background()); err != nil {
			slog.Error("goaneco-playwright: continueAsFallback error", "err", err)
		}
	}()
}

// buildGlobRegexp compiles a glob pattern into a regexp.
// ** matches any path segment sequence; * matches within a segment; ? matches one non-slash character.
// {png,jpg} expands to an alternation group (png|jpg).
func buildGlobRegexp(pattern string) (*regexp.Regexp, error) {
	var re strings.Builder
	re.WriteString("^")
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				re.WriteString(".*")
				i++
			} else {
				re.WriteString("[^/]*")
			}
		case '?':
			re.WriteString("[^/]")
		case '{':
			// Expand {a,b,c} into (a|b|c). Scan to matching '}'.
			end := strings.IndexByte(pattern[i+1:], '}')
			if end < 0 {
				// No closing brace — treat as literal.
				re.WriteString(`\{`)
			} else {
				alternatives := pattern[i+1 : i+1+end]
				re.WriteByte('(')
				parts := strings.Split(alternatives, ",")
				for j, part := range parts {
					if j > 0 {
						re.WriteByte('|')
					}
					for _, c := range part {
						if strings.ContainsRune(`.+()[]^$|\\`, c) {
							re.WriteByte('\\')
						}
						re.WriteRune(c)
					}
				}
				re.WriteByte(')')
				i += end + 1 // skip past '}'
			}
		case '.', '+', '(', ')', '[', ']', '^', '$', '|', '\\':
			re.WriteByte('\\')
			re.WriteByte(pattern[i])
		default:
			re.WriteByte(pattern[i])
		}
	}
	re.WriteString("$")
	return regexp.Compile(re.String())
}

// matchGlob returns true if the URL matches the glob pattern.
func matchGlob(pattern, url string) bool {
	compiled, err := buildGlobRegexp(pattern)
	if err != nil {
		return true // fail open: if pattern is invalid, let the handler run
	}
	return compiled.MatchString(url)
}
