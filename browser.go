package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// browserInitializer contains protocol initializer data for Browser.
// Corresponds to the `initializer` field in the `__create__` event for Browser.
// Ref: packages/protocol/spec/browser.yml
type browserInitializer struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	BrowserName string `json:"browserName"`
}

// contextListener pairs a handler with a unique ID so it can be removed by cancel func.
type contextListener struct {
	id uint64
	fn func(*BrowserContext)
}

// Browser represents a local or remote browser instance.
type Browser struct {
	owner            ChannelOwner
	btype            *BrowserType // reference to the BrowserType that created it via Launch
	initializer      browserInitializer
	mu               sync.RWMutex
	connected        bool
	contexts         []*BrowserContext
	contextListeners []contextListener
	nextListenerID   uint64
}

// Version returns the browser version.
// The value comes from the protocol initializer sent in the __create__ event.
func (b *Browser) Version() string {
	return b.initializer.Version
}

// BrowserType returns the BrowserType that launched this browser.
// Returns nil if the Browser was created through a mechanism other than Launch.
func (b *Browser) BrowserType() *BrowserType {
	return b.btype
}

// IsConnected returns true if the browser is actively connected to the Playwright driver.
func (b *Browser) IsConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.connected
}

// OnContext registers a callback that is invoked whenever a new BrowserContext is created in this browser.
// The returned function unregisters the handler; calling it more than once is safe.
func (b *Browser) OnContext(fn func(*BrowserContext)) func() {
	b.mu.Lock()
	b.nextListenerID++
	id := b.nextListenerID
	b.contextListeners = append(b.contextListeners, contextListener{id: id, fn: fn})
	b.mu.Unlock()
	var once sync.Once
	return func() {
		once.Do(func() {
			b.mu.Lock()
			defer b.mu.Unlock()
			for i, l := range b.contextListeners {
				if l.id == id {
					b.contextListeners = append(b.contextListeners[:i], b.contextListeners[i+1:]...)
					return
				}
			}
		})
	}
}

// Contexts returns a slice of all active BrowserContext instances belonging to this Browser.
func (b *Browser) Contexts() []*BrowserContext {
	b.mu.RLock()
	defer b.mu.RUnlock()
	res := make([]*BrowserContext, len(b.contexts))
	copy(res, b.contexts)
	return res
}

func (b *Browser) removeContext(target *BrowserContext) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, c := range b.contexts {
		if c == target {
			b.contexts = append(b.contexts[:i], b.contexts[i+1:]...)
			break
		}
	}
}

// Close terminates the browser and all associated browser contexts.
// It is idempotent: later calls perform no action and return nil.
//
// Fail-fast design: local state (connected = false, contexts = nil) is updated immediately
// to prevent redundant concurrent calls. If the underlying IPC call fails due to driver disconnection,
// the Playwright server automatically cleans up remote resources on process termination.
func (b *Browser) Close(ctx context.Context) error {
	b.mu.Lock()
	if !b.connected {
		b.mu.Unlock()
		return nil
	}
	b.connected = false
	contexts := b.contexts
	b.contexts = nil
	b.mu.Unlock()

	for _, c := range contexts {
		c.mu.Lock()
		c.closed = true
		pages := c.pages
		c.pages = nil
		c.mu.Unlock()
		for _, p := range pages {
			p.markClosed()
		}
	}

	req := protocol.BrowserCloseRequest{}
	_, err := b.owner.SendMessageRequest(ctx, "close", req)
	if err != nil {
		return fmt.Errorf("failed to close browser: %w", err)
	}
	return nil
}

// NewContext creates a new isolated browser context.
// Accepts optional BrowserContextOptions for configuring viewport, userAgent, locale, geolocation, etc.
func (b *Browser) NewContext(ctx context.Context, options ...*BrowserContextOptions) (*BrowserContext, error) {
	var wire browserContextNewContextParamsWire

	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		if opt.AcceptDownloads != nil {
			s := "deny"
			if *opt.AcceptDownloads {
				s = "accept"
			}
			wire.AcceptDownloads = &s
		}
		wire.BaseURL = opt.BaseURL
		wire.BypassCSP = opt.BypassCSP
		wire.ColorScheme = opt.ColorScheme
		wire.DeviceScaleFactor = opt.DeviceScaleFactor
		wire.ForcedColors = opt.ForcedColors
		wire.Geolocation = opt.Geolocation
		wire.HasTouch = opt.HasTouch
		wire.HttpCredentials = opt.HttpCredentials
		wire.IgnoreHTTPSErrors = opt.IgnoreHTTPSErrors
		wire.IsMobile = opt.IsMobile
		wire.JavaScriptEnabled = opt.JavaScriptEnabled
		wire.Locale = opt.Locale
		wire.NoDefaultViewport = opt.NoDefaultViewport
		wire.Offline = opt.Offline
		wire.Permissions = opt.Permissions
		wire.Proxy = opt.Proxy
		wire.ReducedMotion = opt.ReducedMotion
		wire.Screen = opt.Screen
		wire.StorageState = opt.StorageState
		if opt.StorageStatePath != nil && wire.StorageState == nil {
			data, err := os.ReadFile(*opt.StorageStatePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read StorageStatePath %q: %w", *opt.StorageStatePath, err)
			}
			var ss StorageState
			if err := json.Unmarshal(data, &ss); err != nil {
				return nil, fmt.Errorf("failed to parse StorageStatePath %q: %w", *opt.StorageStatePath, err)
			}
			wire.StorageState = &ss
		}
		wire.ServiceWorkers = opt.ServiceWorkers
		wire.StrictSelectors = opt.StrictSelectors
		wire.TimezoneID = opt.TimezoneID
		wire.UserAgent = opt.UserAgent
		wire.Viewport = opt.Viewport
		if opt.RecordVideo != nil {
			wire.RecordVideo = opt.RecordVideo
		}

		if len(opt.ExtraHTTPHeaders) > 0 {
			wire.ExtraHTTPHeaders = make([]protocol.NameValue, 0, len(opt.ExtraHTTPHeaders))
			for k, v := range opt.ExtraHTTPHeaders {
				wire.ExtraHTTPHeaders = append(wire.ExtraHTTPHeaders, protocol.NameValue{
					Name:  k,
					Value: v,
				})
			}
		}
	}

	result, err := b.owner.SendMessageRequest(ctx, "newContext", wire)
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	var resp protocol.BrowserNewContextResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse newContext response: %w", err)
	}

	if resp.Context.Guid == "" {
		return nil, fmt.Errorf("browser.newContext: server returned empty context GUID")
	}

	contextObj, ok := b.owner.conn.GetObject(resp.Context.Guid)
	if !ok {
		return nil, fmt.Errorf("browser.newContext: context object not found in connection registry")
	}

	bCtx := &BrowserContext{
		owner:       b.owner.child(resp.Context.Guid),
		initializer: contextObj.Initializer(),
		browser:     b,
	}

	// Extract the requestContext GUID from the BrowserContext initializer.
	// The Playwright server creates an APIRequestContext for each BrowserContext
	// and embeds its GUID under the "requestContext" key.
	var bcInit struct {
		RequestContext struct {
			GUID string `json:"guid"`
		} `json:"requestContext"`
	}
	if len(contextObj.Initializer()) > 0 {
		if err := json.Unmarshal(contextObj.Initializer(), &bcInit); err == nil && bcInit.RequestContext.GUID != "" {
			bCtx.apiRequestCtx = &APIRequestContext{
				owner: b.owner.child(bcInit.RequestContext.GUID),
			}
		}
	}

	b.mu.Lock()
	b.contexts = append(b.contexts, bCtx)
	listeners := make([]contextListener, len(b.contextListeners))
	copy(listeners, b.contextListeners)
	b.mu.Unlock()

	for _, l := range listeners {
		if l.fn != nil {
			l.fn(bCtx)
		}
	}

	return bCtx, nil
}

// NewPage is a convenience method that creates a new browser context and then a new page.
// If page creation fails, the context is closed to prevent leaks on the server.
func (b *Browser) NewPage(ctx context.Context, options ...*BrowserContextOptions) (*Page, error) {
	bCtx, err := b.NewContext(ctx, options...)
	if err != nil {
		return nil, err
	}
	page, err := bCtx.NewPage(ctx)
	if err != nil {
		// Intentionally uses Background: the caller's ctx may already be cancelled (e.g. on timeout),
		// but we still need to close the temporary context to avoid resource leaks.
		cleanupCtx, cancel := context.WithTimeout(context.Background(), defaultBrowserCleanupTimeout)
		defer cancel()
		_ = bCtx.Close(cleanupCtx) // best-effort cleanup; NewPage error takes precedence
		return nil, err
	}
	return page, nil
}
