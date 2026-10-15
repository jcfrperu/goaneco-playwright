//go:build e2e

// E2E tests for: TestGeolocation.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"strings"
	"testing"
	"time"
)

// TestBrowserContextSetGeolocationCanUpdate verifies geolocation can be changed.
// Extra: no equivalent method in playwright-java; related class: TestGeolocation.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
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
// Ref: TestGeolocation.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
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

// TestGeolocationShouldThrowWhenInvalidLongitude verifies SetGeolocation rejects longitude > 180.
// Ref: TestGeolocation.java#shouldThrowWhenInvalidLongitude  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
func TestGeolocationShouldThrowWhenInvalidLongitude(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)

	bc := newContext(t)
	err := bc.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 10, Longitude: 200})
	is.Error(err, "longitude 200 should be rejected")
	is.ErrorContains(err, "longitude")
}

// TestGeolocationShouldIsolateContexts verifies that geolocation is isolated per BrowserContext.
// Ref: TestGeolocation.java#shouldIsolateContexts  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
func TestGeolocationShouldIsolateContexts(t *testing.T) {
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

// TestGeolocationShouldUseContextOptions verifies geolocation set in NewContext options is applied.
// Ref: TestGeolocation.java#shouldUseContextOptions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
func TestGeolocationShouldUseContextOptions(t *testing.T) {
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

// TestGeolocationWatchPositionShouldBeNotified verifies watchPosition fires callbacks as geolocation changes.
// Ref: TestGeolocation.java#watchPositionShouldBeNotified  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
func TestGeolocationWatchPositionShouldBeNotified(t *testing.T) {
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

// TestBrowserContextOptionsGeolocation verifies geolocation and permissions in BrowserContextOptions.
// Ref: TestGeolocation.java#shouldIsolateContexts
func TestBrowserContextOptionsGeolocation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		Geolocation: &playwright.Geolocation{
			Latitude:  48.8584,
			Longitude: 2.2945,
		},
		Permissions: []string{"geolocation"},
	})
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	coordsRaw, err := page.Evaluate(ctx, `() => new Promise(resolve => {
		navigator.geolocation.getCurrentPosition(pos => {
			resolve({ latitude: pos.coords.latitude, longitude: pos.coords.longitude });
		});
	})`)
	must.NoError(err, "Evaluate geolocation failed")

	coords, ok := coordsRaw.(map[string]any)
	is.True(ok, "expected map[string]any")

	lat, okLat := coords["latitude"].(float64)
	lng, okLng := coords["longitude"].(float64)
	const tolerance = 0.0001
	if !okLat || !okLng || math.Abs(lat-48.8584) > tolerance || math.Abs(lng-2.2945) > tolerance {
		t.Errorf("geolocation coords = {lat: %v, lng: %v}, want approx {48.8584, 2.2945}", lat, lng)
	}
}

// TestGeolocationShouldNotModifyPassedDefaultOptionsObject verifies SetGeolocation does not mutate the options passed to NewContext.
// Ref: TestGeolocation.java#shouldNotModifyPassedDefaultOptionsObject  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
func TestGeolocationShouldNotModifyPassedDefaultOptionsObject(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	geo := &playwright.Geolocation{Latitude: 10, Longitude: 10}
	options := &playwright.BrowserContextOptions{Geolocation: geo}
	bCtx := a4NewContext(t, options)
	must.NoError(bCtx.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 20, Longitude: 20}))
	must.Equal(geo, options.Geolocation)
	must.Equal(playwright.Geolocation{Latitude: 10, Longitude: 10}, *options.Geolocation)
}

// TestGeolocationShouldUseContextOptionsForPopup verifies popups use the geolocation set on the context.
// Ref: TestGeolocation.java#shouldUseContextOptionsForPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestGeolocation.java
func TestGeolocationShouldUseContextOptionsForPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/geolocation.html", "text/html", `<script>
window['geolocationPromise'] = new Promise(resolve => {
  navigator.geolocation.getCurrentPosition(position => {
    resolve({latitude: position.coords.latitude, longitude: position.coords.longitude});
  });
});
</script>`)
	bCtx := a4NewContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(bCtx.GrantPermissions(ctx, []string{"geolocation"}))
	must.NoError(bCtx.SetGeolocation(ctx, &playwright.Geolocation{Latitude: 10, Longitude: 10}))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	popup, err := page.WaitForPopup(ctx, func() error {
		_, e := page.Evaluate(ctx, "url => window['_popup'] = window.open(url)", srv.Prefix()+"/geolocation.html")
		return e
	})
	must.NoError(err)
	must.NoError(popup.WaitForLoadState(ctx))
	geolocation, err := popup.Evaluate(ctx, "window['geolocationPromise']")
	must.NoError(err)
	must.EqualValues(map[string]any{"longitude": 10.0, "latitude": 10.0}, geolocation)
}
