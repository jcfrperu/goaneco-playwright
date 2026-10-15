//go:generate go run ./cmd/gen/main.go

package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jcfrperu/goaneco-playwright/internal/driver"
)

type browserTypeRef struct {
	GUID string `json:"guid"`
}

// playwrightInitializer describes the initial state of the Playwright root object.
// Note: Other fields sent by the server (e.g. selectors, electron, android)
// are intentionally omitted for now.
type playwrightInitializer struct {
	Chromium browserTypeRef `json:"chromium"`
	Firefox  browserTypeRef `json:"firefox"`
	WebKit   browserTypeRef `json:"webkit"`
}

// Playwright represents the main entry point to the Playwright automation library.
type Playwright struct {
	owner       ChannelOwner
	driver      *driver.Driver
	initializer playwrightInitializer
}

// Run initializes the Node.js driver, performs the "initialize" handshake with Root, and sets up IPC.
// The caller's context controls the timeout for the initial handshake.
// The operational lifecycle of the driver is independent and controlled via Playwright.Stop().
func Run(ctx context.Context, cliPath string) (*Playwright, error) {
	if _, err := os.Stat(cliPath); err != nil {
		return nil, fmt.Errorf("playwright driver not found at %q: %w", cliPath, err)
	}

	// context.Background() — driver process must outlive the initialization context.
	d, err := driver.StartDriver(context.Background(), cliPath)
	if err != nil {
		return nil, fmt.Errorf("failed to start driver: %w", err)
	}
	// stopOnErr stops the driver on any initialization failure so the caller doesn't need to.
	var initFailed bool
	defer func() {
		if initFailed {
			_ = d.Stop()
		}
	}()

	rootOwner := ChannelOwner{
		conn: d.Conn(),
		guid: "",
	}
	initParams := map[string]string{
		"sdkLanguage": "javascript",
	}
	respBytes, err := rootOwner.SendMessageRequest(ctx, "initialize", initParams)
	if err != nil {
		initFailed = true
		return nil, fmt.Errorf("failed to initialize Playwright: %w", err)
	}

	var initResp struct {
		Playwright struct {
			Guid string `json:"guid"`
		} `json:"playwright"`
	}
	if err := json.Unmarshal(respBytes, &initResp); err != nil {
		initFailed = true
		return nil, fmt.Errorf("failed to parse initialize response: %w", err)
	}

	obj, ok := d.Conn().GetObject(initResp.Playwright.Guid)
	if !ok {
		var waitErr error
		obj, waitErr = d.Conn().WaitPlaywright(ctx)
		if waitErr != nil {
			initFailed = true
			return nil, fmt.Errorf("playwright object %q not found in connection registry: %w", initResp.Playwright.Guid, waitErr)
		}
	}

	pw := &Playwright{
		owner: ChannelOwner{
			conn: d.Conn(),
			guid: obj.GUID(),
		},
		driver: d,
	}

	if len(obj.Initializer()) > 0 {
		if err := json.Unmarshal(obj.Initializer(), &pw.initializer); err != nil {
			initFailed = true
			return nil, fmt.Errorf("failed to parse playwright initializer: %w", err)
		}
	}

	return pw, nil
}

// browserType creates a local BrowserType reference using the GUID provided by the server.
// Loads protocol initializer data (name, executablePath) from the connection registry.
func (p *Playwright) browserType(guid string) *BrowserType {
	var init browserTypeInitializer
	if p.owner.conn != nil {
		if obj, ok := p.owner.conn.GetObject(guid); ok {
			if raw := obj.Initializer(); len(raw) > 0 {
				_ = json.Unmarshal(raw, &init) // best-effort; zero-value fallback if server data is malformed
			}
		}
	}
	return &BrowserType{
		owner:       p.owner.child(guid),
		initializer: init,
	}
}

// Chromium returns the Chromium BrowserType instance.
func (p *Playwright) Chromium() (*BrowserType, error) {
	if p.initializer.Chromium.GUID == "" {
		return nil, fmt.Errorf("chromium is not available in this Playwright installation")
	}
	return p.browserType(p.initializer.Chromium.GUID), nil
}

// Firefox returns the Firefox BrowserType instance.
func (p *Playwright) Firefox() (*BrowserType, error) {
	if p.initializer.Firefox.GUID == "" {
		return nil, fmt.Errorf("firefox is not available in this Playwright installation")
	}
	return p.browserType(p.initializer.Firefox.GUID), nil
}

// WebKit returns the WebKit BrowserType instance.
func (p *Playwright) WebKit() (*BrowserType, error) {
	if p.initializer.WebKit.GUID == "" {
		return nil, fmt.Errorf("webkit is not available in this Playwright installation")
	}
	return p.browserType(p.initializer.WebKit.GUID), nil
}

// Stop safely terminates the driver process.
func (p *Playwright) Stop() error {
	return p.driver.Stop()
}
