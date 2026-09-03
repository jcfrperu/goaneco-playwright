package playwright

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// browserTypeInitializer contains protocol initializer data for BrowserType.
// Corresponds to the `initializer` field in the `__create__` event for BrowserType.
type browserTypeInitializer struct {
	ExecutablePath string `json:"executablePath"`
	Name           string `json:"name"`
}

// BrowserTypeLaunchOptions provides complete configuration settings for BrowserType.Launch.
type BrowserTypeLaunchOptions struct {
	Args            []string          `json:"args,omitempty"`
	Channel         *string           `json:"channel,omitempty"`
	ChromiumSandbox *bool             `json:"chromiumSandbox,omitempty"`
	Devtools        *bool             `json:"devtools,omitempty"`
	DownloadsPath   *string           `json:"downloadsPath,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	ExecutablePath  *string           `json:"executablePath,omitempty"`
	// FirefoxUserPrefs is Firefox-only. json:"-" prevents accidental JSON serialization of this
	// public field; it is copied into the internal wire struct (firefoxUserPrefs) during Launch.
	FirefoxUserPrefs map[string]any `json:"-"`
	HandleSIGHUP     *bool          `json:"handleSIGHUP,omitempty"`
	HandleSIGINT     *bool          `json:"handleSIGINT,omitempty"`
	HandleSIGTERM    *bool          `json:"handleSIGTERM,omitempty"`
	Headless         *bool          `json:"headless,omitempty"`
	SlowMo           *float64       `json:"slowMo,omitempty"`
	Timeout          *float64       `json:"timeout,omitempty"`
	TracesDir        *string        `json:"tracesDir,omitempty"`
}

type browserTypeLaunchParamsWire struct {
	Args             []string             `json:"args,omitempty"`
	Channel          *string              `json:"channel,omitempty"`
	ChromiumSandbox  *bool                `json:"chromiumSandbox,omitempty"`
	Devtools         *bool                `json:"devtools,omitempty"`
	DownloadsPath    *string              `json:"downloadsPath,omitempty"`
	Env              []protocol.NameValue `json:"env,omitempty"`
	ExecutablePath   *string              `json:"executablePath,omitempty"`
	FirefoxUserPrefs map[string]any       `json:"firefoxUserPrefs,omitempty"`
	HandleSIGHUP     *bool                `json:"handleSIGHUP,omitempty"`
	HandleSIGINT     *bool                `json:"handleSIGINT,omitempty"`
	HandleSIGTERM    *bool                `json:"handleSIGTERM,omitempty"`
	Headless         *bool                `json:"headless,omitempty"`
	SlowMo           *float64             `json:"slowMo,omitempty"`
	Timeout          *float64             `json:"timeout,omitempty"`
	TracesDir        *string              `json:"tracesDir,omitempty"`
}

// BrowserType provides methods to launch and connect to browser instances.
type BrowserType struct {
	owner       ChannelOwner
	initializer browserTypeInitializer
}

// Name returns the browser name (e.g. "chromium", "firefox", "webkit").
// The value originates from the protocol initializer.
func (b *BrowserType) Name() string {
	return b.initializer.Name
}

// ExecutablePath returns the path to the browser executable binary.
// The value originates from the protocol initializer.
func (b *BrowserType) ExecutablePath() string {
	return b.initializer.ExecutablePath
}

// Launch launches a new browser instance using this BrowserType.
// Options can be omitted to use default settings, or a pointer to BrowserTypeLaunchOptions
// to customize launching (headless, args, slowMo, timeout, channel, executablePath, etc.).
func (b *BrowserType) Launch(ctx context.Context, opts ...*BrowserTypeLaunchOptions) (*Browser, error) {
	var options *BrowserTypeLaunchOptions
	if len(opts) > 0 {
		options = opts[0]
	}
	dt := defaultLaunchTimeoutMs
	wire := browserTypeLaunchParamsWire{
		Timeout: &dt,
	}

	if options != nil {
		wire.Args = options.Args
		wire.Channel = options.Channel
		wire.ChromiumSandbox = options.ChromiumSandbox
		wire.Devtools = options.Devtools
		wire.DownloadsPath = options.DownloadsPath
		wire.ExecutablePath = options.ExecutablePath
		wire.FirefoxUserPrefs = options.FirefoxUserPrefs
		wire.HandleSIGHUP = options.HandleSIGHUP
		wire.HandleSIGINT = options.HandleSIGINT
		wire.HandleSIGTERM = options.HandleSIGTERM
		wire.Headless = options.Headless
		wire.SlowMo = options.SlowMo
		wire.TracesDir = options.TracesDir

		if options.Timeout != nil {
			wire.Timeout = options.Timeout
		}

		if len(options.Env) > 0 {
			wire.Env = make([]protocol.NameValue, 0, len(options.Env))
			for k, v := range options.Env {
				wire.Env = append(wire.Env, protocol.NameValue{
					Name:  k,
					Value: v,
				})
			}
		}
	}

	respBytes, err := b.owner.SendMessageRequest(ctx, "launch", wire)
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	var resp protocol.BrowserTypeLaunchResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal launch response: %w", err)
	}

	if resp.Browser.Guid == "" {
		return nil, fmt.Errorf("server returned empty browser GUID")
	}

	// Retrieve browser initializer from connection registry
	browserObj, ok := b.owner.conn.GetObject(resp.Browser.Guid)
	if !ok {
		return nil, fmt.Errorf("browser object not found in connection after launch")
	}

	// Parse browser initializer for Version(), Name(), etc.
	var bInit browserInitializer
	if raw := browserObj.Initializer(); len(raw) > 0 {
		if err := json.Unmarshal(raw, &bInit); err != nil {
			return nil, fmt.Errorf("failed to parse browser initializer: %w", err)
		}
	}

	return &Browser{
		owner:       b.owner.child(resp.Browser.Guid),
		btype:       b,
		initializer: bInit,
		connected:   true,
	}, nil
}

// ConnectOverCDPOptions specifies optional parameters for BrowserType.ConnectOverCDP.
type ConnectOverCDPOptions struct {
	// Headers is a map of additional HTTP headers sent with the WebSocket handshake.
	Headers map[string]string
	// SlowMo slows down Playwright operations by the specified amount of milliseconds.
	SlowMo *float64
}

// ConnectOverCDP connects to an existing browser instance over the Chrome DevTools Protocol.
// endpointURL is the CDP WebSocket endpoint URL (e.g. "http://localhost:9222/").
func (b *BrowserType) ConnectOverCDP(
	ctx context.Context,
	endpointURL string,
	opts ...*ConnectOverCDPOptions,
) (*Browser, error) {
	req := protocol.BrowserTypeConnectOverCDPRequest{
		EndpointURL: &endpointURL,
		Headers:     []protocol.NameValue{},
	}

	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		req.SlowMo = o.SlowMo

		if len(o.Headers) > 0 {
			req.Headers = make([]protocol.NameValue, 0, len(o.Headers))
			for k, v := range o.Headers {
				req.Headers = append(req.Headers, protocol.NameValue{
					Name:  k,
					Value: v,
				})
			}
		}
	}

	respBytes, err := b.owner.SendMessageRequest(ctx, "connectOverCDP", req)
	if err != nil {
		return nil, fmt.Errorf("failed to connectOverCDP: %w", err)
	}

	var resp protocol.BrowserTypeConnectOverCDPResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connectOverCDP response: %w", err)
	}

	if resp.Browser.Guid == "" {
		return nil, fmt.Errorf("connectOverCDP: server returned empty browser GUID")
	}

	var bInit browserInitializer
	if obj, ok := b.owner.conn.GetObject(resp.Browser.Guid); ok {
		if raw := obj.Initializer(); len(raw) > 0 {
			_ = json.Unmarshal(raw, &bInit)
		}
	}

	return &Browser{
		owner:       b.owner.child(resp.Browser.Guid),
		btype:       b,
		initializer: bInit,
		connected:   true,
	}, nil
}

// LaunchPersistentContextOptions configures BrowserType.LaunchPersistentContext.
type LaunchPersistentContextOptions struct {
	// Launch options (subset of BrowserTypeLaunchOptions)
	Args             []string
	Channel          *string
	Headless         *bool
	SlowMo           *float64
	ExecutablePath   *string
	FirefoxUserPrefs map[string]any
	// Context options (subset of BrowserContextOptions)
	AcceptDownloads   *bool
	BaseURL           *string
	BypassCSP         *bool
	ColorScheme       *string
	Contrast          *string
	DeviceScaleFactor *float64
	ExtraHTTPHeaders  map[string]string
	ForcedColors      *string
	Geolocation       *Geolocation
	HasTouch          *bool
	HttpCredentials   *HttpCredentials
	IgnoreHTTPSErrors *bool
	IsMobile          *bool
	JavaScriptEnabled *bool
	Locale            *string
	Offline           *bool
	Permissions       []string
	ReducedMotion     *string
	TimezoneID        *string
	UserAgent         *string
	Viewport          *ViewportSize
}

// launchPersistentContextWire is the flat wire format for LaunchPersistentContext.
// All launch params and context params are sent at the top level; the Playwright
// protocol's $mixin notation in the schema is documentation-only, not wire format.
type launchPersistentContextWire struct {
	// Launch params
	Args             []string       `json:"args,omitempty"`
	Channel          *string        `json:"channel,omitempty"`
	ChromiumSandbox  *bool          `json:"chromiumSandbox,omitempty"`
	ExecutablePath   *string        `json:"executablePath,omitempty"`
	FirefoxUserPrefs map[string]any `json:"firefoxUserPrefs,omitempty"`
	Headless         *bool          `json:"headless,omitempty"`
	SlowMo           *float64       `json:"slowMo,omitempty"`
	Timeout          float64        `json:"timeout"`
	// Context params
	AcceptDownloads   *string              `json:"acceptDownloads,omitempty"` // "accept"|"deny"|"internal-browser-default"
	BaseURL           *string              `json:"baseURL,omitempty"`
	BypassCSP         *bool                `json:"bypassCSP,omitempty"`
	ColorScheme       *string              `json:"colorScheme,omitempty"`
	Contrast          *string              `json:"contrast,omitempty"`
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
	Offline           *bool                `json:"offline,omitempty"`
	Permissions       []string             `json:"permissions,omitempty"`
	ReducedMotion     *string              `json:"reducedMotion,omitempty"`
	TimezoneID        *string              `json:"timezoneId,omitempty"`
	UserAgent         *string              `json:"userAgent,omitempty"`
	Viewport          *ViewportSize        `json:"viewport,omitempty"`
	// Required
	UserDataDir string `json:"userDataDir"`
}

// LaunchPersistentContext launches a browser instance with a persistent user profile
// stored at userDataDir. Returns a Browser and its associated BrowserContext.
func (b *BrowserType) LaunchPersistentContext(ctx context.Context, userDataDir string, opts ...*LaunchPersistentContextOptions) (*Browser, *BrowserContext, error) {
	wire := launchPersistentContextWire{
		Timeout:     defaultLaunchTimeoutMs,
		UserDataDir: userDataDir,
	}

	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		// Launch options
		wire.Args = o.Args
		wire.Channel = o.Channel
		wire.FirefoxUserPrefs = o.FirefoxUserPrefs
		wire.Headless = o.Headless
		wire.ExecutablePath = o.ExecutablePath
		wire.SlowMo = o.SlowMo
		// Context options
		wire.BaseURL = o.BaseURL
		wire.BypassCSP = o.BypassCSP
		wire.ColorScheme = o.ColorScheme
		wire.Contrast = o.Contrast
		wire.DeviceScaleFactor = o.DeviceScaleFactor
		wire.ForcedColors = o.ForcedColors
		wire.Geolocation = o.Geolocation
		wire.HasTouch = o.HasTouch
		wire.HttpCredentials = o.HttpCredentials
		wire.IgnoreHTTPSErrors = o.IgnoreHTTPSErrors
		wire.IsMobile = o.IsMobile
		wire.JavaScriptEnabled = o.JavaScriptEnabled
		wire.Locale = o.Locale
		wire.Offline = o.Offline
		wire.Permissions = o.Permissions
		wire.ReducedMotion = o.ReducedMotion
		wire.TimezoneID = o.TimezoneID
		wire.UserAgent = o.UserAgent
		wire.Viewport = o.Viewport

		if o.AcceptDownloads != nil {
			if *o.AcceptDownloads {
				accept := "accept"
				wire.AcceptDownloads = &accept
			} else {
				deny := "deny"
				wire.AcceptDownloads = &deny
			}
		}

		if len(o.ExtraHTTPHeaders) > 0 {
			wire.ExtraHTTPHeaders = make([]protocol.NameValue, 0, len(o.ExtraHTTPHeaders))
			for k, v := range o.ExtraHTTPHeaders {
				wire.ExtraHTTPHeaders = append(wire.ExtraHTTPHeaders, protocol.NameValue{
					Name:  k,
					Value: v,
				})
			}
		}
	}

	respBytes, err := b.owner.SendMessageRequest(ctx, "launchPersistentContext", wire)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to launchPersistentContext: %w", err)
	}

	var resp protocol.BrowserTypeLaunchPersistentContextResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal launchPersistentContext response: %w", err)
	}

	if resp.Browser.Guid == "" || resp.Context.Guid == "" {
		return nil, nil, fmt.Errorf("launchPersistentContext: server returned empty GUIDs")
	}

	// Retrieve browser initializer
	var bInit browserInitializer
	if obj, ok := b.owner.conn.GetObject(resp.Browser.Guid); ok {
		if raw := obj.Initializer(); len(raw) > 0 {
			_ = json.Unmarshal(raw, &bInit) // best-effort; zero-value fallback if server data is malformed
		}
	}

	browser := &Browser{
		owner:       b.owner.child(resp.Browser.Guid),
		btype:       b,
		initializer: bInit,
		connected:   true,
	}

	// Retrieve context initializer
	var ctxInit json.RawMessage
	if obj, ok := b.owner.conn.GetObject(resp.Context.Guid); ok {
		ctxInit = obj.Initializer()
	}

	bCtx := &BrowserContext{
		owner:       b.owner.child(resp.Context.Guid),
		initializer: ctxInit,
		browser:     browser,
	}

	// Extract the APIRequestContext GUID from the BrowserContext initializer (same as NewContext does).
	var bcInit struct {
		RequestContext struct {
			GUID string `json:"guid"`
		} `json:"requestContext"`
	}
	if len(ctxInit) > 0 {
		if err := json.Unmarshal(ctxInit, &bcInit); err == nil && bcInit.RequestContext.GUID != "" {
			bCtx.apiRequestCtx = &APIRequestContext{
				owner: b.owner.child(bcInit.RequestContext.GUID),
			}
		}
	}

	browser.mu.Lock()
	browser.contexts = append(browser.contexts, bCtx)
	browser.mu.Unlock()

	return browser, bCtx, nil
}
