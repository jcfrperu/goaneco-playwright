package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// ViewportSize defines the browser window viewport dimensions.
type ViewportSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Geolocation specifies geographical coordinates.
type Geolocation struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Accuracy  *float64 `json:"accuracy,omitempty"`
}

// HttpCredentials represents HTTP basic/digest authentication credentials.
type HttpCredentials struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Origin   *string `json:"origin,omitempty"`
	Send     *string `json:"send,omitempty"`
}

// ScreenSize defines emulated screen dimensions.
type ScreenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ProxyOptions defines proxy configuration for a BrowserContext.
type ProxyOptions struct {
	Server   string  `json:"server"`
	Bypass   *string `json:"bypass,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

// LocalStorageEntry represents a key-value pair in localStorage.
type LocalStorageEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// OriginStorage represents localStorage key-values for a given origin.
type OriginStorage struct {
	Origin       string              `json:"origin"`
	LocalStorage []LocalStorageEntry `json:"localStorage"`
}

// SameSiteAttribute represents the SameSite policy for cookies.
type SameSiteAttribute string

const (
	SameSiteStrict SameSiteAttribute = "Strict"
	SameSiteLax    SameSiteAttribute = "Lax"
	SameSiteNone   SameSiteAttribute = "None"
)

// Cookie represents a cookie in the storage state or browser context.
type Cookie struct {
	Name     string            `json:"name"`
	Value    string            `json:"value"`
	URL      *string           `json:"url,omitempty"`
	Domain   *string           `json:"domain,omitempty"`
	Path     *string           `json:"path,omitempty"`
	Expires  *float64          `json:"expires,omitempty"`
	HTTPOnly *bool             `json:"httpOnly,omitempty"`
	Secure   *bool             `json:"secure,omitempty"`
	SameSite SameSiteAttribute `json:"sameSite,omitempty"`
}

// ClearCookiesOptions configures optional filters for clearing cookies.
type ClearCookiesOptions struct {
	Name      *string        `json:"name,omitempty"`
	NameRegex *regexp.Regexp `json:"-"`
	Domain    *string        `json:"domain,omitempty"`
	Path      *string        `json:"path,omitempty"`
}

// StorageState contains serialized cookies and localStorage entries to populate a BrowserContext.
type StorageState struct {
	Cookies []Cookie        `json:"cookies,omitempty"`
	Origins []OriginStorage `json:"origins,omitempty"`
}

// BrowserContextOptions provides complete configuration settings for Browser.NewContext.
type BrowserContextOptions struct {
	AcceptDownloads   *bool               `json:"acceptDownloads,omitempty"`
	BaseURL           *string             `json:"baseURL,omitempty"`
	BypassCSP         *bool               `json:"bypassCSP,omitempty"`
	ColorScheme       *string             `json:"colorScheme,omitempty"`
	DeviceScaleFactor *float64            `json:"deviceScaleFactor,omitempty"`
	ExtraHTTPHeaders  map[string]string   `json:"extraHTTPHeaders,omitempty"`
	ForcedColors      *string             `json:"forcedColors,omitempty"`
	Geolocation       *Geolocation        `json:"geolocation,omitempty"`
	HasTouch          *bool               `json:"hasTouch,omitempty"`
	HttpCredentials   *HttpCredentials    `json:"httpCredentials,omitempty"`
	IgnoreHTTPSErrors *bool               `json:"ignoreHTTPSErrors,omitempty"`
	IsMobile          *bool               `json:"isMobile,omitempty"`
	JavaScriptEnabled *bool               `json:"javaScriptEnabled,omitempty"`
	Locale            *string             `json:"locale,omitempty"`
	NoDefaultViewport *bool               `json:"noDefaultViewport,omitempty"`
	Offline           *bool               `json:"offline,omitempty"`
	Permissions       []string            `json:"permissions,omitempty"`
	Proxy             *ProxyOptions       `json:"proxy,omitempty"`
	ReducedMotion     *string             `json:"reducedMotion,omitempty"`
	Screen            *ScreenSize         `json:"screen,omitempty"`
	RecordVideo       *RecordVideoOptions `json:"-"`                        // client-side: passed to wire params separately
	ServiceWorkers    *string             `json:"serviceWorkers,omitempty"` // "allow" | "block"
	StorageState      *StorageState       `json:"storageState,omitempty"`
	StorageStatePath  *string             `json:"-"` // client-side: path to a JSON file containing StorageState
	StrictSelectors   *bool               `json:"strictSelectors,omitempty"`
	TimezoneID        *string             `json:"timezoneId,omitempty"`
	UserAgent         *string             `json:"userAgent,omitempty"`
	Viewport          *ViewportSize       `json:"viewport,omitempty"`
}

type browserContextNewContextParamsWire struct {
	AcceptDownloads   *string              `json:"acceptDownloads,omitempty"` // "accept"|"deny"|"internal-browser-default"
	BaseURL           *string              `json:"baseURL,omitempty"`
	BypassCSP         *bool                `json:"bypassCSP,omitempty"`
	ColorScheme       *string              `json:"colorScheme,omitempty"`
	DeviceScaleFactor *float64             `json:"deviceScaleFactor,omitempty"`
	ExtraHTTPHeaders  []protocol.NameValue `json:"extraHTTPHeaders,omitempty"`
	ForcedColors      *string              `json:"forcedColors,omitempty"`
	Geolocation       *Geolocation         `json:"geolocation,omitempty"`
	HasTouch          *bool                `json:"hasTouch,omitempty"`
	HttpCredentials   *HttpCredentials     `json:"httpCredentials,omitempty"`
	IgnoreHTTPSErrors *bool                `json:"ignoreHTTPSErrors,omitempty"`
	IsMobile          *bool                `json:"isMobile,omitempty"`
	JavaScriptEnabled *bool                `json:"javaScriptEnabled,omitempty"`
	Locale            *string              `json:"locale,omitempty"`
	NoDefaultViewport *bool                `json:"noDefaultViewport,omitempty"`
	Offline           *bool                `json:"offline,omitempty"`
	Permissions       []string             `json:"permissions,omitempty"`
	Proxy             *ProxyOptions        `json:"proxy,omitempty"`
	RecordVideo       *RecordVideoOptions  `json:"recordVideo,omitempty"`
	ReducedMotion     *string              `json:"reducedMotion,omitempty"`
	Screen            *ScreenSize          `json:"screen,omitempty"`
	ServiceWorkers    *string              `json:"serviceWorkers,omitempty"` // "allow" | "block"
	StorageState      *StorageState        `json:"storageState,omitempty"`
	StrictSelectors   *bool                `json:"strictSelectors,omitempty"`
	TimezoneID        *string              `json:"timezoneId,omitempty"`
	UserAgent         *string              `json:"userAgent,omitempty"`
	Viewport          *ViewportSize        `json:"viewport,omitempty"`
}

// BrowserContext represents an isolated browser session (equivalent to an incognito profile).
type BrowserContext struct {
	owner                  ChannelOwner
	initializer            json.RawMessage // stored for future usage or introspection
	browser                *Browser        // reference to the parent Browser
	apiRequestCtx          *APIRequestContext
	mu                     sync.RWMutex
	closed                 bool
	pages                  []*Page
	pagesByGUID            map[string]*Page                 // GUID → Page; prevents duplicate registration for popup pages
	pagesInFlight          map[string]struct{}              // GUIDs currently being constructed; prevents duplicate event subscriptions
	ctxRouter              routeRouter                      // manages BrowserContext-level route handler chain
	ctxRouteListenerID     connection.ListenerID            // single event listener for context "route" events
	ctxRouteListenerActive bool                             // true once the listener is registered
	ctxRoutePatterns       []string                         // accumulated glob patterns for network interception
	bindingListenerIDs     map[string]connection.ListenerID // binding name → listener ID; prevents duplicate dispatch
	nextHandlerID          int
	pageCloseHandlers      map[int]func(*Page) // called when any page in this context fires its "close" event
	pageLoadHandlers       map[int]func(*Page) // called when any page in this context reaches "load" state
	pageListenersByID      map[int]func(*Page)
	pageListenerID         connection.ListenerID
	pageListenerActive     bool
	subscriptionCounts     map[string]int // reference counts for updateSubscription IPC events
	testIdAttributeName    string         // configurable attribute for GetByTestId; default "data-testid"
}

// Request returns the APIRequestContext bound to this BrowserContext.
// Requests made through it share cookies and credentials with pages in the context.
func (c *BrowserContext) Request() *APIRequestContext {
	return c.apiRequestCtx
}

// pageInitializer is the data structure of the `initializer` field in the `__create__` event for Page.
type pageInitializer struct {
	MainFrame struct {
		GUID string `json:"guid"`
	} `json:"mainFrame"`
	Opener *struct {
		GUID string `json:"guid"`
	} `json:"opener,omitempty"`
	ViewportSize *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"viewportSize,omitempty"`
}

// Browser returns the Browser instance to which this context belongs.
func (c *BrowserContext) Browser() *Browser {
	return c.browser
}

// Pages returns a slice containing all active pages in this BrowserContext.
func (c *BrowserContext) Pages() []*Page {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]*Page, len(c.pages))
	copy(res, c.pages)
	return res
}

func (c *BrowserContext) removePage(target *Page) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, p := range c.pages {
		if p == target {
			copy(c.pages[i:], c.pages[i+1:])
			c.pages[len(c.pages)-1] = nil
			c.pages = c.pages[:len(c.pages)-1]
			break
		}
	}
	delete(c.pagesByGUID, target.owner.guid)
}

// pageByGUID returns the registered Page with the given GUID, or nil if not found.
func (c *BrowserContext) pageByGUID(guid string) *Page {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pagesByGUID[guid]
}

// pageFromGUID creates a Page object from a server-assigned page GUID.
// Used for pages created server-side (e.g., window.open popups).
func (c *BrowserContext) pageFromGUID(guid string) (*Page, error) {
	pageObj, ok := c.owner.conn.GetObject(guid)
	if !ok {
		return nil, fmt.Errorf("page object %q not found in connection", guid)
	}

	var pageInit pageInitializer
	initData := pageObj.Initializer()
	if len(initData) > 0 {
		if err := json.Unmarshal(initData, &pageInit); err != nil {
			return nil, fmt.Errorf("failed to parse page initializer: %w", err)
		}
	}
	if pageInit.MainFrame.GUID == "" {
		return nil, fmt.Errorf("page initializer for %q missing mainFrame GUID", guid)
	}

	var frameInit frameInitializer
	if frameObj, ok := c.owner.conn.GetObject(pageInit.MainFrame.GUID); ok {
		if raw := frameObj.Initializer(); len(raw) > 0 {
			_ = json.Unmarshal(raw, &frameInit) // best-effort; zero-value fallback if server data is malformed
		}
	}

	openerGUID := ""
	if pageInit.Opener != nil {
		openerGUID = pageInit.Opener.GUID
	}

	var viewportSize *ViewportSize
	if pageInit.ViewportSize != nil {
		viewportSize = &ViewportSize{
			Width:  pageInit.ViewportSize.Width,
			Height: pageInit.ViewportSize.Height,
		}
	}

	page := &Page{
		owner:          c.owner.child(guid),
		mainFrame:      c.owner.child(pageInit.MainFrame.GUID),
		frameInit:      frameInit,
		initializer:    initData,
		openerGUID:     openerGUID,
		browserContext: c,
		framesByGUID:   make(map[string]*Frame),
		viewportSize:   viewportSize,
	}
	page.Keyboard = &Keyboard{page: page}
	page.Mouse = &Mouse{page: page}
	return page, nil
}

// registerPage adds a page to the context's tracking under the write lock.
// Returns false if the page was already registered (deduplication).
func (c *BrowserContext) registerPage(page *Page) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pagesByGUID == nil {
		c.pagesByGUID = make(map[string]*Page)
	}
	if _, exists := c.pagesByGUID[page.owner.guid]; exists {
		return false
	}
	c.pages = append(c.pages, page)
	c.pagesByGUID[page.owner.guid] = page
	return true
}

// OnPage registers a handler that fires whenever a new page is created in this context,
// including popup pages opened via window.open(). The returned function unregisters the handler.
func (c *BrowserContext) OnPage(fn func(*Page)) func() {
	c.mu.Lock()
	c.nextHandlerID++
	id := c.nextHandlerID
	if c.pageListenersByID == nil {
		c.pageListenersByID = make(map[int]func(*Page))
	}
	c.pageListenersByID[id] = fn
	if !c.pageListenerActive {
		c.pageListenerActive = true
		c.pageListenerID = c.owner.conn.OnEvent(c.owner.guid, "page", func(params json.RawMessage) {
			var event protocol.BrowserContextPageEvent
			if err := json.Unmarshal(params, &event); err != nil || event.Page.Guid == "" {
				return
			}

			c.mu.Lock()
			existing := c.pagesByGUID[event.Page.Guid]
			if existing == nil {
				if c.pagesInFlight == nil {
					c.pagesInFlight = make(map[string]struct{})
				}
				if _, inFlight := c.pagesInFlight[event.Page.Guid]; inFlight {
					c.mu.Unlock()
					return
				}
				c.pagesInFlight[event.Page.Guid] = struct{}{}
			}
			listeners := make([]func(*Page), 0, len(c.pageListenersByID))
			for _, cb := range c.pageListenersByID {
				listeners = append(listeners, cb)
			}
			c.mu.Unlock()

			if existing != nil {
				for _, cb := range listeners {
					cb := cb
					go cb(existing)
				}
				return
			}

			defer func() {
				c.mu.Lock()
				delete(c.pagesInFlight, event.Page.Guid)
				c.mu.Unlock()
			}()

			page, err := c.pageFromGUID(event.Page.Guid)
			if err != nil {
				slog.Default().Error("OnPage: failed to construct page from event", "guid", event.Page.Guid, "error", err)
				return
			}
			page.subscribeToNavigation()
			page.subscribeToVideo()
			page.subscribeToPageErrors()
			page.subscribeToPageClose()
			page.subscribeToPageLoad()
			page.subscribeToFrames()
			page.subscribeToWorkers()

			c.registerPage(page)
			for _, cb := range listeners {
				cb := cb
				go cb(page)
			}
		})
	}
	c.mu.Unlock()
	return func() {
		c.mu.Lock()
		delete(c.pageListenersByID, id)
		c.mu.Unlock()
	}
}

// retainSubscription increments the reference count for a server-side subscription event.
// The first retain sends updateSubscription(enabled:true) to the Playwright server.
func (c *BrowserContext) retainSubscription(ctx context.Context, event string) {
	c.mu.Lock()
	if c.subscriptionCounts == nil {
		c.subscriptionCounts = make(map[string]int)
	}
	c.subscriptionCounts[event]++
	count := c.subscriptionCounts[event]
	c.mu.Unlock()
	if count == 1 {
		if _, err := c.owner.SendMessageRequest(ctx, "updateSubscription",
			protocol.BrowserContextUpdateSubscriptionRequest{Event: event, Enabled: true}); err != nil {
			slog.Default().Warn("retainSubscription: failed to enable", "event", event, "error", err)
		}
	}
}

// releaseSubscription decrements the reference count for a server-side subscription event.
// When the count reaches zero, updateSubscription(enabled:false) is sent to the Playwright server.
func (c *BrowserContext) releaseSubscription(event string) {
	c.mu.Lock()
	if c.subscriptionCounts == nil || c.subscriptionCounts[event] <= 0 {
		c.mu.Unlock()
		return // never subscribed; nothing to release
	}
	c.subscriptionCounts[event]--
	count := c.subscriptionCounts[event]
	c.mu.Unlock()

	if count == 0 {
		releaseCtx, cancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
		defer cancel()
		if _, err := c.owner.SendMessageRequest(releaseCtx, "updateSubscription",
			protocol.BrowserContextUpdateSubscriptionRequest{Event: event, Enabled: false}); err != nil {
			slog.Default().Warn("releaseSubscription: failed to disable", "event", event, "error", err)
		}
	}
}

// IsClosed returns true if the context has been closed.
func (c *BrowserContext) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// Close closes the browser context and all its associated pages.
// It is idempotent: later calls perform no action and return nil.
//
// Fail-fast design: local state (closed = true, pages = nil) is updated immediately
// to prevent redundant concurrent calls. If the underlying IPC call fails due to driver disconnection,
// the Playwright server automatically cleans up remote resources on process termination.
func (c *BrowserContext) Close(ctx context.Context) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	b := c.browser
	pages := c.pages
	c.pages = nil
	c.pagesByGUID = nil
	c.mu.Unlock()

	// Mark each owned page as closed so page.IsClosed() reflects context closure.
	for _, p := range pages {
		p.markClosed()
	}

	// Detach from parent Browser if present
	if b != nil {
		b.removeContext(c)
	}

	req := protocol.BrowserContextCloseRequest{}
	_, err := c.owner.SendMessageRequest(ctx, "close", req)
	if err != nil {
		return fmt.Errorf("failed to close context: %w", err)
	}
	return nil
}

// NewPage creates a new page in this browser context.
// Resolves the mainFrame GUID from the page initializer, then loads the Frame
// initializer so the initial URL is immediately available without additional IPC calls.
func (c *BrowserContext) NewPage(ctx context.Context) (*Page, error) {
	req := protocol.BrowserContextNewPageRequest{}
	result, err := c.owner.SendMessageRequest(ctx, "newPage", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	var resp protocol.BrowserContextNewPageResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse newPage response: %w", err)
	}

	if resp.Page.Guid == "" {
		return nil, fmt.Errorf("browserContext.newPage: server returned empty page GUID")
	}

	page, err := c.pageFromGUID(resp.Page.Guid)
	if err != nil {
		return nil, fmt.Errorf("browserContext.newPage: %w", err)
	}

	page.subscribeToNavigation()
	page.subscribeToVideo()
	page.subscribeToPageErrors()
	page.subscribeToPageClose()
	page.subscribeToPageLoad()
	page.subscribeToFrames()
	page.subscribeToWorkers()
	c.registerPage(page)
	return page, nil
}

// Cookies returns cookies matching specified URLs or all cookies in context if URLs are omitted.
func (c *BrowserContext) Cookies(ctx context.Context, urls ...string) ([]Cookie, error) {
	req := protocol.BrowserContextCookiesRequest{
		Urls: urls,
	}
	if req.Urls == nil {
		req.Urls = []string{}
	}
	result, err := c.owner.SendMessageRequest(ctx, "cookies", req)
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %w", err)
	}

	var resp protocol.BrowserContextCookiesResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse cookies response: %w", err)
	}

	cookies := make([]Cookie, len(resp.Cookies))
	for i, ck := range resp.Cookies {
		cookies[i] = cookieFromProtocol(ck)
	}
	return cookies, nil
}

type setNetworkCookieWire struct {
	Name     string             `json:"name"`
	Value    string             `json:"value"`
	Url      *string            `json:"url,omitempty"`
	Domain   *string            `json:"domain,omitempty"`
	Path     *string            `json:"path,omitempty"`
	Expires  *float64           `json:"expires,omitempty"`
	HttpOnly *bool              `json:"httpOnly,omitempty"`
	Secure   *bool              `json:"secure,omitempty"`
	SameSite *SameSiteAttribute `json:"sameSite,omitempty"`
}

// cookieFromProtocol converts a protocol.NetworkCookie to the public Cookie type.
func cookieFromProtocol(nc protocol.NetworkCookie) Cookie {
	domain := nc.Domain
	path := nc.Path
	expires := nc.Expires
	httpOnly := nc.HttpOnly
	secure := nc.Secure
	var sameSite SameSiteAttribute
	if s, ok := nc.SameSite.(string); ok {
		sameSite = SameSiteAttribute(s)
	}
	return Cookie{
		Name:     nc.Name,
		Value:    nc.Value,
		Domain:   &domain,
		Path:     &path,
		Expires:  &expires,
		HTTPOnly: &httpOnly,
		Secure:   &secure,
		SameSite: sameSite,
	}
}

type browserContextAddCookiesParamsWire struct {
	Cookies []setNetworkCookieWire `json:"cookies"`
}

// AddCookies adds cookies to this browser context.
func (c *BrowserContext) AddCookies(ctx context.Context, cookies []Cookie) error {
	setCookies := make([]setNetworkCookieWire, len(cookies))
	for i, ck := range cookies {
		var sameSite *SameSiteAttribute
		if ck.SameSite != "" {
			s := ck.SameSite // local copy to avoid capturing loop variable address
			sameSite = &s
		}
		setCookies[i] = setNetworkCookieWire{
			Name:     ck.Name,
			Value:    ck.Value,
			Url:      ck.URL,
			Domain:   ck.Domain,
			Path:     ck.Path,
			Expires:  ck.Expires,
			HttpOnly: ck.HTTPOnly,
			Secure:   ck.Secure,
			SameSite: sameSite,
		}
	}

	req := browserContextAddCookiesParamsWire{
		Cookies: setCookies,
	}
	_, err := c.owner.SendMessageRequest(ctx, "addCookies", req)
	if err != nil {
		return fmt.Errorf("failed to add cookies: %w", err)
	}
	return nil
}

// ClearCookies clears cookies from the context, optionally filtered by name, domain, or path.
func (c *BrowserContext) ClearCookies(ctx context.Context, options ...*ClearCookiesOptions) error {
	var req protocol.BrowserContextClearCookiesRequest
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		req.Name = opt.Name
		req.Domain = opt.Domain
		req.Path = opt.Path
		if opt.NameRegex != nil {
			src, flags := regexpToJS(opt.NameRegex)
			req.NameRegexSource = &src
			req.NameRegexFlags = &flags
		}
	}
	_, err := c.owner.SendMessageRequest(ctx, "clearCookies", req)
	if err != nil {
		return fmt.Errorf("failed to clear cookies: %w", err)
	}
	return nil
}

// regexpToJS extracts the JS-compatible source and flags from a Go regexp.
func regexpToJS(re *regexp.Regexp) (source, flags string) {
	return extractRegexpInfo(re)
}

// SetExtraHTTPHeaders sets extra HTTP headers to be sent with every request in this context.
// Calling with an empty map removes all previously set extra headers.
func (c *BrowserContext) SetExtraHTTPHeaders(ctx context.Context, headers map[string]string) error {
	pairs := make([]protocol.NameValue, 0, len(headers))
	for k, v := range headers {
		pairs = append(pairs, protocol.NameValue{Name: k, Value: v})
	}
	req := protocol.BrowserContextSetExtraHTTPHeadersRequest{Headers: pairs}
	_, err := c.owner.SendMessageRequest(ctx, "setExtraHTTPHeaders", req)
	if err != nil {
		return fmt.Errorf("browserContext.setExtraHTTPHeaders failed: %w", err)
	}
	return nil
}

// GrantPermissions grants browser permissions to the given origin.
// Pass no origin to grant permissions to all origins.
func (c *BrowserContext) GrantPermissions(ctx context.Context, permissions []string, origin ...string) error {
	req := protocol.BrowserContextGrantPermissionsRequest{Permissions: permissions}
	if len(origin) > 0 {
		req.Origin = &origin[0]
	}
	_, err := c.owner.SendMessageRequest(ctx, "grantPermissions", req)
	if err != nil {
		return fmt.Errorf("browserContext.grantPermissions failed: %w", err)
	}
	return nil
}

// ClearPermissions resets all permission overrides for this context.
func (c *BrowserContext) ClearPermissions(ctx context.Context) error {
	_, err := c.owner.SendMessageRequest(ctx, "clearPermissions", protocol.BrowserContextClearPermissionsRequest{})
	if err != nil {
		return fmt.Errorf("browserContext.clearPermissions failed: %w", err)
	}
	return nil
}

// SetGeolocation sets the geolocation for this context.
// Pass nil to clear any previously set geolocation.
func (c *BrowserContext) SetGeolocation(ctx context.Context, geo *Geolocation) error {
	var geolocation any
	if geo != nil {
		geolocation = geo
	}
	req := protocol.BrowserContextSetGeolocationRequest{Geolocation: geolocation}
	_, err := c.owner.SendMessageRequest(ctx, "setGeolocation", req)
	if err != nil {
		return fmt.Errorf("browserContext.setGeolocation failed: %w", err)
	}
	return nil
}

// SetOffline toggles offline mode for the entire context.
func (c *BrowserContext) SetOffline(ctx context.Context, offline bool) error {
	req := protocol.BrowserContextSetOfflineRequest{Offline: offline}
	_, err := c.owner.SendMessageRequest(ctx, "setOffline", req)
	if err != nil {
		return fmt.Errorf("browserContext.setOffline failed: %w", err)
	}
	return nil
}

// StorageState returns a snapshot of the current cookies and localStorage for this context.
func (c *BrowserContext) StorageState(ctx context.Context) (*StorageState, error) {
	result, err := c.owner.SendMessageRequest(ctx, "storageState", protocol.BrowserContextStorageStateRequest{})
	if err != nil {
		return nil, fmt.Errorf("browserContext.storageState failed: %w", err)
	}
	var resp protocol.BrowserContextStorageStateResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse browserContext.storageState response: %w", err)
	}

	state := &StorageState{
		Cookies: make([]Cookie, 0, len(resp.Cookies)),
		Origins: make([]OriginStorage, 0, len(resp.Origins)),
	}
	for _, pc := range resp.Cookies {
		state.Cookies = append(state.Cookies, cookieFromProtocol(pc))
	}
	for _, po := range resp.Origins {
		entries := make([]LocalStorageEntry, 0, len(po.LocalStorage))
		for _, nv := range po.LocalStorage {
			entries = append(entries, LocalStorageEntry{Name: nv.Name, Value: nv.Value})
		}
		state.Origins = append(state.Origins, OriginStorage{
			Origin:       po.Origin,
			LocalStorage: entries,
		})
	}
	return state, nil
}

// BrowserContextRouteOptions specifies optional parameters for BrowserContext.Route.
type BrowserContextRouteOptions struct {
	// Times limits how many times the handler is invoked. After Times calls the
	// handler is automatically skipped (subsequent requests go through normally).
	Times *int
}

// Route registers a network interception handler for all pages in this context.
// Multiple calls stack handlers in LIFO order; a handler may call Route.Fallback to defer
// to the next registered handler.
func (c *BrowserContext) Route(ctx context.Context, pattern string, handler RouteHandler, opts ...*BrowserContextRouteOptions) error {
	c.mu.Lock()
	found := false
	for _, existing := range c.ctxRoutePatterns {
		if existing == pattern {
			found = true
			break
		}
	}
	if !found {
		c.ctxRoutePatterns = append(c.ctxRoutePatterns, pattern)
	}
	allPatterns := make([]any, len(c.ctxRoutePatterns))
	for i, pat := range c.ctxRoutePatterns {
		allPatterns[i] = map[string]string{"glob": pat}
	}
	c.mu.Unlock()

	req := protocol.BrowserContextSetNetworkInterceptionPatternsRequest{
		Patterns: allPatterns,
	}
	_, err := c.owner.SendMessageRequest(ctx, "setNetworkInterceptionPatterns", req)
	if err != nil {
		return fmt.Errorf("browserContext.route: failed to set interception patterns: %w", err)
	}

	times := 0
	if len(opts) > 0 && opts[0] != nil && opts[0].Times != nil {
		times = *opts[0].Times
	}

	c.mu.Lock()
	c.ctxRouter.add(pattern, handler, times)
	if !c.ctxRouteListenerActive {
		c.ctxRouteListenerID = c.owner.conn.OnEvent(c.owner.guid, "route", func(params json.RawMessage) {
			var event protocol.BrowserContextRouteEvent
			if err := json.Unmarshal(params, &event); err != nil {
				return
			}
			if event.Route.Guid == "" {
				return
			}
			route := &Route{owner: c.owner.child(event.Route.Guid)}
			c.ctxRouter.dispatch(route)
		})
		c.ctxRouteListenerActive = true
	}
	c.mu.Unlock()
	return nil
}

// Unroute removes all network interception patterns from this context
// and deregisters the route event listener.
func (c *BrowserContext) Unroute(ctx context.Context) error {
	c.mu.Lock()
	c.ctxRouter.clear()
	c.ctxRoutePatterns = nil
	listenerID := c.ctxRouteListenerID
	active := c.ctxRouteListenerActive
	c.ctxRouteListenerActive = false
	c.mu.Unlock()

	if active {
		c.owner.conn.OffEvent(c.owner.guid, "route", listenerID)
	}

	req := protocol.BrowserContextSetNetworkInterceptionPatternsRequest{Patterns: []any{}}
	_, err := c.owner.SendMessageRequest(ctx, "setNetworkInterceptionPatterns", req)
	if err != nil {
		return fmt.Errorf("browserContext.unroute failed: %w", err)
	}
	return nil
}

// OnPageClose registers a handler that fires when any page in this context closes.
// Returns a function that cancels the registration.
func (c *BrowserContext) OnPageClose(fn func(*Page)) func() {
	c.mu.Lock()
	if c.pageCloseHandlers == nil {
		c.pageCloseHandlers = make(map[int]func(*Page))
	}
	id := c.nextHandlerID
	c.nextHandlerID++
	c.pageCloseHandlers[id] = fn
	c.mu.Unlock()
	return func() {
		c.mu.Lock()
		delete(c.pageCloseHandlers, id)
		c.mu.Unlock()
	}
}

// OnPageLoad registers a handler that fires when any page in this context reaches "load" state.
// Returns a function that cancels the registration.
func (c *BrowserContext) OnPageLoad(fn func(*Page)) func() {
	c.mu.Lock()
	if c.pageLoadHandlers == nil {
		c.pageLoadHandlers = make(map[int]func(*Page))
	}
	id := c.nextHandlerID
	c.nextHandlerID++
	c.pageLoadHandlers[id] = fn
	c.mu.Unlock()
	return func() {
		c.mu.Lock()
		delete(c.pageLoadHandlers, id)
		c.mu.Unlock()
	}
}

// callPageCloseHandlers invokes all registered OnPageClose handlers for the given page.
func (c *BrowserContext) callPageCloseHandlers(p *Page) {
	c.mu.RLock()
	handlers := make([]func(*Page), 0, len(c.pageCloseHandlers))
	for _, h := range c.pageCloseHandlers {
		handlers = append(handlers, h)
	}
	c.mu.RUnlock()
	for _, h := range handlers {
		h := h
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Default().Error("goaneco-playwright: page handler panic", "panic", r)
				}
			}()
			h(p)
		}()
	}
}

// callPageLoadHandlers invokes all registered OnPageLoad handlers for the given page.
func (c *BrowserContext) callPageLoadHandlers(p *Page) {
	c.mu.RLock()
	handlers := make([]func(*Page), 0, len(c.pageLoadHandlers))
	for _, h := range c.pageLoadHandlers {
		handlers = append(handlers, h)
	}
	c.mu.RUnlock()
	for _, h := range handlers {
		h := h
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Default().Error("goaneco-playwright: page handler panic", "panic", r)
				}
			}()
			h(p)
		}()
	}
}

// OnRequest registers a handler that fires for every network request in this context.
// The returned function cancels the listener; it is idempotent (safe to call multiple times).
func (c *BrowserContext) OnRequest(handler func(*NetworkRequest)) func() {
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	c.retainSubscription(subCtx, "request")
	subCancel()
	process := func(params json.RawMessage) {
		var event struct {
			Request protocol.Request `json:"request"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		raw := c.owner.Initializer(event.Request.Guid)
		if req := networkRequestFrom(c.owner, event.Request.Guid, raw); req != nil {
			go handler(req)
		}
	}
	id := c.owner.conn.OnEvent(c.owner.guid, "request", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			c.owner.conn.OffEvent(c.owner.guid, "request", id)
			c.releaseSubscription("request")
		})
	}
}

// OnResponse registers a handler that fires for every network response in this context.
// The returned function cancels the listener; it is idempotent (safe to call multiple times).
func (c *BrowserContext) OnResponse(handler func(*NetworkResponse)) func() {
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	c.retainSubscription(subCtx, "response")
	subCancel()
	process := func(params json.RawMessage) {
		var event struct {
			Response protocol.Response `json:"response"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		raw := c.owner.Initializer(event.Response.Guid)
		if resp := networkResponseFrom(c.owner, event.Response.Guid, raw); resp != nil {
			go handler(resp)
		}
	}
	id := c.owner.conn.OnEvent(c.owner.guid, "response", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			c.owner.conn.OffEvent(c.owner.guid, "response", id)
			c.releaseSubscription("response")
		})
	}
}

// OnRequestFinished registers a handler that fires when a network request completes.
// The returned function cancels the listener; it is idempotent (safe to call multiple times).
func (c *BrowserContext) OnRequestFinished(handler func(*NetworkRequest)) func() {
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	c.retainSubscription(subCtx, "requestFinished")
	subCancel()
	process := func(params json.RawMessage) {
		var event protocol.BrowserContextRequestFinishedEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		raw := c.owner.Initializer(event.Request.Guid)
		if req := networkRequestFrom(c.owner, event.Request.Guid, raw); req != nil {
			go handler(req)
		}
	}
	id := c.owner.conn.OnEvent(c.owner.guid, "requestFinished", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			c.owner.conn.OffEvent(c.owner.guid, "requestFinished", id)
			c.releaseSubscription("requestFinished")
		})
	}
}

// OnRequestFailed registers a handler that fires when a network request fails.
// The returned function cancels the listener; it is idempotent (safe to call multiple times).
func (c *BrowserContext) OnRequestFailed(handler func(*NetworkRequest)) func() {
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	c.retainSubscription(subCtx, "requestFailed")
	subCancel()
	process := func(params json.RawMessage) {
		var event protocol.BrowserContextRequestFailedEvent
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}
		raw := c.owner.Initializer(event.Request.Guid)
		if req := networkRequestFrom(c.owner, event.Request.Guid, raw); req != nil {
			go handler(req)
		}
	}
	id := c.owner.conn.OnEvent(c.owner.guid, "requestFailed", process)
	var once sync.Once
	return func() {
		once.Do(func() {
			c.owner.conn.OffEvent(c.owner.guid, "requestFailed", id)
			c.releaseSubscription("requestFailed")
		})
	}
}

// WaitForEvent waits for a context-level event to fire and returns its payload.
// Supported events: "page" (*Page), "request" (*NetworkRequest), "response" (*NetworkResponse).
func (c *BrowserContext) WaitForEvent(ctx context.Context, event string) (any, error) {
	ch := make(chan any, 1)
	send := func(v any) {
		select {
		case ch <- v:
		default:
		}
	}

	var unsub func()
	switch event {
	case "page":
		unsub = c.OnPage(func(pg *Page) { send(pg) })
	case "request":
		unsub = c.OnRequest(func(r *NetworkRequest) { send(r) })
	case "response":
		unsub = c.OnResponse(func(r *NetworkResponse) { send(r) })
	default:
		return nil, fmt.Errorf("browserContext.WaitForEvent: unsupported event %q", event)
	}
	defer unsub()

	select {
	case val := <-ch:
		return val, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("browserContext.WaitForEvent(%q): %w", event, ctx.Err())
	}
}

// AddInitScript injects a JavaScript script that runs before any page script on every
// navigation in all pages of this context.
func (b *BrowserContext) AddInitScript(ctx context.Context, script string) error {
	_, err := b.owner.SendMessageRequest(ctx, "addInitScript", protocol.BrowserContextAddInitScriptRequest{Source: script})
	if err != nil {
		return fmt.Errorf("browserContext.addInitScript failed: %w", err)
	}
	return nil
}

// SetTestIdAttribute changes the attribute name used by GetByTestId locators.
// The default is "data-testid". Calling this in a @AfterEach/cleanup resets it.
func (c *BrowserContext) SetTestIdAttribute(ctx context.Context, attributeName string) error {
	req := protocol.BrowserContextSetTestIdAttributeNameRequest{TestIdAttributeName: attributeName}
	_, err := c.owner.SendMessageRequest(ctx, "setTestIdAttributeName", req)
	if err != nil {
		return fmt.Errorf("browserContext.setTestIdAttributeName failed: %w", err)
	}
	c.mu.Lock()
	c.testIdAttributeName = attributeName
	c.mu.Unlock()
	return nil
}

// TestIdAttributeName returns the attribute name used by GetByTestId locators.
// Defaults to "data-testid" if SetTestIdAttribute has not been called.
func (c *BrowserContext) TestIdAttributeName() string {
	c.mu.RLock()
	attr := c.testIdAttributeName
	c.mu.RUnlock()
	if attr == "" {
		return "data-testid"
	}
	return attr
}

// RegisterSelectorEngineOptions specifies options for BrowserContext.RegisterSelectorEngine.
type RegisterSelectorEngineOptions struct {
	// ContentScript controls whether the selector engine runs in the page's content script context.
	// When true, the engine runs in an isolated world; when false (the default), it runs in the main world.
	ContentScript *bool
}

// RegisterSelectorEngine registers a custom selector engine for use within this BrowserContext.
// name is the selector engine name (only [a-zA-Z0-9_] characters are allowed).
// script is a JavaScript string that evaluates to a selector engine object with
// create, query, and queryAll methods.
func (b *BrowserContext) RegisterSelectorEngine(
	ctx context.Context,
	name string,
	script string,
	opts ...*RegisterSelectorEngineOptions,
) error {
	req := protocol.BrowserContextRegisterSelectorEngineRequest{
		SelectorEngine: protocol.SelectorEngine{
			Name:   name,
			Source: script,
		},
	}
	if len(opts) > 0 && opts[0] != nil && opts[0].ContentScript != nil {
		req.SelectorEngine.ContentScript = opts[0].ContentScript
	}
	_, err := b.owner.SendMessageRequest(ctx, "registerSelectorEngine", req)
	if err != nil {
		return fmt.Errorf("browserContext.RegisterSelectorEngine: %w", err)
	}
	return nil
}
