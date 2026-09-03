//go:build e2e

// BrowserContext Geolocation E2E tests.
// Migration of: TestBrowserContextGeolocation.java (cases)
package e2e

import (
	"strings"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBrowserContextSetGeolocationWithPermission verifies SetGeolocation updates location when permission granted.
// Ref: TestBrowserContextGeolocation.java#shouldSetGeolocation
func TestBrowserContextSetGeolocationWithPermission(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{
		Latitude:  48.8566,
		Longitude: 2.3522,
	}))

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	coords, err := page.Evaluate(ctx, `() => new Promise((resolve, reject) => {
		navigator.geolocation.getCurrentPosition(
			pos => resolve({ lat: pos.coords.latitude, lon: pos.coords.longitude }),
			err => reject(err)
		);
	})`)
	must.NoError(err)
	must.NotNil(coords)

	m, ok := coords.(map[string]any)
	is.True(ok)
	lat, ok := m["lat"].(float64)
	is.True(ok)
	is.InDelta(48.8566, lat, 0.001)
}

// TestBrowserContextSetGeolocationCanUpdate verifies geolocation can be changed.
// Ref: TestBrowserContextGeolocation.java#shouldUpdateGeolocation
func TestBrowserContextSetGeolocationCanUpdate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{
		Latitude:  10.0,
		Longitude: 20.0,
	}))

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{
		Latitude:  30.0,
		Longitude: 40.0,
	}))
}

// ---------------------------------------------------------------------------
// From geolocation_extra_test.go
// ---------------------------------------------------------------------------

// TestGeolocationShouldWork verifies SetGeolocation is reflected by getCurrentPosition.
// Ref: TestGeolocation.java#shouldWork
func TestGeolocationShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 10, Longitude: 10}))

	result, err := page.Evaluate(ctx, `() => new Promise(resolve =>
		navigator.geolocation.getCurrentPosition(pos =>
			resolve({ latitude: pos.coords.latitude, longitude: pos.coords.longitude })
		)
	)`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.InDelta(10.0, m["latitude"], 0.001)
	is.InDelta(10.0, m["longitude"], 0.001)
}

// TestGeolocationThrowWhenInvalidLongitude verifies SetGeolocation rejects longitude > 180.
// Ref: TestGeolocation.java#shouldThrowWhenInvalidLongitude
func TestGeolocationThrowWhenInvalidLongitude(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	bc := newContext(t)
	err := bc.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 10, Longitude: 200})
	is.Error(err, "longitude 200 should be rejected")
	is.ErrorContains(err, "longitude")
}

// TestGeolocationIsolateContexts verifies that geolocation is isolated per BrowserContext.
// Ref: TestGeolocation.java#shouldIsolateContexts
func TestGeolocationIsolateContexts(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc1, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc1.Close(ctx) })
	must.NoError(bc1.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 10, Longitude: 10}))
	page1, err := bc1.NewPage(ctx)
	must.NoError(err)
	must.NoError(page1.Goto(ctx, srv.EmptyPage()))

	bc2, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
		Geolocation: &playwright.Geolocation{Latitude: 20, Longitude: 20},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc2.Close(ctx) })
	page2, err := bc2.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	eval := func(p *playwright.Page) (float64, float64) {
		t.Helper()
		result, err := p.Evaluate(ctx, `() => new Promise(resolve =>
			navigator.geolocation.getCurrentPosition(pos =>
				resolve({ latitude: pos.coords.latitude, longitude: pos.coords.longitude })
			)
		)`)
		must.NoError(err)
		m, ok := result.(map[string]any)
		is.True(ok)
		lat := m["latitude"].(float64)
		lon := m["longitude"].(float64)
		return lat, lon
	}

	lat1, lon1 := eval(page1)
	is.InDelta(10.0, lat1, 0.001)
	is.InDelta(10.0, lon1, 0.001)

	lat2, lon2 := eval(page2)
	is.InDelta(20.0, lat2, 0.001)
	is.InDelta(20.0, lon2, 0.001)
}

// TestGeolocationUseContextOptions verifies geolocation set in NewContext options is applied.
// Ref: TestGeolocation.java#shouldUseContextOptions
func TestGeolocationUseContextOptions(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
		Geolocation: &playwright.Geolocation{Latitude: 10, Longitude: 10},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() => new Promise(resolve =>
		navigator.geolocation.getCurrentPosition(pos =>
			resolve({ latitude: pos.coords.latitude, longitude: pos.coords.longitude })
		)
	)`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.InDelta(10.0, m["latitude"], 0.001)
	is.InDelta(10.0, m["longitude"], 0.001)
}

// TestGeolocationWatchPositionNotified verifies watchPosition fires callbacks as geolocation changes.
// Ref: TestGeolocation.java#watchPositionShouldBeNotified
func TestGeolocationWatchPositionNotified(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Permissions: []string{"geolocation"},
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bc.Close(ctx) })

	page, err := bc.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 0, Longitude: 0}))

	msgCh := make(chan string, 50)
	off := page.OnConsole(func(msg *playwright.ConsoleMessage) {
		if msg.Type() == "log" {
			msgCh <- msg.Text()
		}
	})
	defer off()

	_, err = page.Evaluate(ctx, `() => {
		navigator.geolocation.watchPosition(pos => {
			const coords = pos.coords;
			console.log("lat=" + coords.latitude + " lng=" + coords.longitude);
		}, err => {});
	}`)
	must.NoError(err)

	waitForMsg := func(want string) {
		t.Helper()
		deadline := time.After(10 * time.Second)
		for {
			select {
			case msg := <-msgCh:
				if strings.Contains(msg, want) {
					return
				}
			case <-deadline:
				t.Fatalf("timed out waiting for console message containing %q", want)
			}
		}
	}

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 0, Longitude: 10}))
	waitForMsg("lat=0 lng=10")

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 20, Longitude: 30}))
	waitForMsg("lat=20 lng=30")

	must.NoError(bc.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 40, Longitude: 50}))
	waitForMsg("lat=40 lng=50")
}
