//go:build e2e

// E2E tests for: TestPageEmulateMedia.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaColorScheme(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	darkScheme := "dark"
	err := page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: &darkScheme})
	must.NoError(err, "EmulateMedia(dark) failed")

	isDark, err := page.Evaluate(ctx, "() => matchMedia('(prefers-color-scheme: dark)').matches")
	must.NoError(err, "Evaluate dark scheme failed")
	if isDark != true {
		t.Errorf("expected dark color scheme to match, got %v", isDark)
	}

	lightScheme := "light"
	err = page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: &lightScheme})
	must.NoError(err, "EmulateMedia(light) failed")

	isLight, err := page.Evaluate(ctx, "() => matchMedia('(prefers-color-scheme: light)').matches")
	must.NoError(err, "Evaluate light scheme failed")
	if isLight != true {
		t.Errorf("expected light color scheme to match, got %v", isLight)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaDefaultsToLight(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	isLight, err := page.Evaluate(ctx, "() => matchMedia('(prefers-color-scheme: light)').matches")
	must.NoError(err, "Evaluate failed")
	if isLight != true {
		t.Errorf("expected default color scheme to be light, got %v", isLight)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaNoOpts(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// EmulateMedia without options should not error.
	err := page.EmulateMedia(ctx)
	must.NoError(err, "EmulateMedia() with no opts failed")
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateForcedColors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	none := "none"
	active := "active"

	noneMatches, err := page.Evaluate(ctx, "() => matchMedia('(forced-colors: none)').matches")
	must.NoError(err)
	is.Equal(true, noneMatches)

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ForcedColors: &none}))
	noneMatches, err = page.Evaluate(ctx, "() => matchMedia('(forced-colors: none)').matches")
	must.NoError(err)
	is.Equal(true, noneMatches)
	activeMatches, err := page.Evaluate(ctx, "() => matchMedia('(forced-colors: active)').matches")
	must.NoError(err)
	is.Equal(false, activeMatches)

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ForcedColors: &active}))
	noneMatches, err = page.Evaluate(ctx, "() => matchMedia('(forced-colors: none)').matches")
	must.NoError(err)
	is.Equal(false, noneMatches)
	activeMatches, err = page.Evaluate(ctx, "() => matchMedia('(forced-colors: active)').matches")
	must.NoError(err)
	is.Equal(true, activeMatches)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaColorSchemeDarkUpdates(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	cs := "dark"
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{
		ColorScheme: &cs,
	}))

	result, err := page.Evaluate(ctx, `() => window.matchMedia('(prefers-color-scheme: dark)').matches`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaColorSchemeLightUpdates(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	cs := "light"
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{
		ColorScheme: &cs,
	}))

	result, err := page.Evaluate(ctx, `() => window.matchMedia('(prefers-color-scheme: light)').matches`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaPrintType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	mt := "print"
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{
		Media: &mt,
	}))

	result, err := page.Evaluate(ctx, `() => window.matchMedia('print').matches`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaScreenType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	mt := "screen"
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{
		Media: &mt,
	}))

	result, err := page.Evaluate(ctx, `() => window.matchMedia('screen').matches`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEmulateMedia.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestEmulateMediaNoPreference(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	cs := "no-preference"
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{
		ColorScheme: &cs,
	}))

	dark, err := page.Evaluate(ctx, `() => window.matchMedia('(prefers-color-scheme: dark)').matches`)
	must.NoError(err)
	is.Equal(false, dark)
}

// Ref: TestPageEmulateMedia.java#shouldEmulateType  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldEmulateType(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// Default is "screen".
	val, err := page.Evaluate(ctx, "() => matchMedia('screen').matches")
	must.NoError(err, "Evaluate screen failed")
	if val != true {
		t.Errorf("expected screen to match by default, got %v", val)
	}

	// Switch to print.
	printMedia := "print"
	err = page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{Media: &printMedia})
	must.NoError(err, "EmulateMedia(print) failed")

	screenMatch, _ := page.Evaluate(ctx, "() => matchMedia('screen').matches")
	printMatch, _ := page.Evaluate(ctx, "() => matchMedia('print').matches")

	if screenMatch != false {
		t.Errorf("after EmulateMedia(print): screen should not match, got %v", screenMatch)
	}
	if printMatch != true {
		t.Errorf("after EmulateMedia(print): print should match, got %v", printMatch)
	}
}

// Ref: TestPageEmulateMedia.java#shouldEmulateReducedMotion  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldEmulateReducedMotion(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	reduce := "reduce"
	err := page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ReducedMotion: &reduce})
	must.NoError(err, "EmulateMedia(reducedMotion=reduce) failed")

	matches, err := page.Evaluate(ctx, "() => matchMedia('(prefers-reduced-motion: reduce)').matches")
	must.NoError(err, "Evaluate failed")
	if matches != true {
		t.Errorf("expected reduced-motion to match, got %v", matches)
	}
}

// Ref: TestPageEmulateMedia.java#shouldChangeTheActualColorsInCss  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldChangeTheActualColorsInCss(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<style>
@media (prefers-color-scheme: dark) {
  div { background: black; color: white; }
}
@media (prefers-color-scheme: light) {
  div { background: white; color: black; }
}
</style>
<div>Hello</div>`))

	getBackground := func() string {
		v, err := page.EvalOnSelector(ctx, "div", "div => window.getComputedStyle(div).backgroundColor")
		must.NoError(err)
		s, _ := v.(string)
		return s
	}

	light := "light"
	dark := "dark"

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: &light}))
	is.Equal("rgb(255, 255, 255)", getBackground())

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: &dark}))
	is.Equal("rgb(0, 0, 0)", getBackground())

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: &light}))
	is.Equal("rgb(255, 255, 255)", getBackground())
}

// Ref: TestPageEmulateMedia.java#shouldEmulateContrast  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldEmulateContrast(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	noPreference := "no-preference"
	more := "more"

	noPrefMatches, err := page.Evaluate(ctx, "matchMedia('(prefers-contrast: no-preference)').matches")
	must.NoError(err)
	is.Equal(true, noPrefMatches)

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{Contrast: &noPreference}))
	noPrefMatches, err = page.Evaluate(ctx, "matchMedia('(prefers-contrast: no-preference)').matches")
	must.NoError(err)
	is.Equal(true, noPrefMatches)
	moreMatches, err := page.Evaluate(ctx, "matchMedia('(prefers-contrast: more)').matches")
	must.NoError(err)
	is.Equal(false, moreMatches)

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{Contrast: &more}))
	noPrefMatches, err = page.Evaluate(ctx, "matchMedia('(prefers-contrast: no-preference)').matches")
	must.NoError(err)
	is.Equal(false, noPrefMatches)
	moreMatches, err = page.Evaluate(ctx, "matchMedia('(prefers-contrast: more)').matches")
	must.NoError(err)
	is.Equal(true, moreMatches)
}

func a2MatchMedia(t *testing.T, ctx context.Context, p interface {
	Evaluate(context.Context, string, ...any) (any, error)
}, scheme string) any {
	t.Helper()
	v, err := p.Evaluate(ctx, "() => matchMedia('(prefers-color-scheme: "+scheme+")').matches")
	if err != nil {
		t.Fatalf("matchMedia: %v", err)
	}
	return v
}

// Ref: TestPageEmulateMedia.java#shouldEmulateSchemeWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldEmulateSchemeWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: a2Ptr("light")}))
	assert.Equal(t, true, a2MatchMedia(t, ctx, page, "light"))
	assert.Equal(t, false, a2MatchMedia(t, ctx, page, "dark"))
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: a2Ptr("dark")}))
	assert.Equal(t, true, a2MatchMedia(t, ctx, page, "dark"))
	assert.Equal(t, false, a2MatchMedia(t, ctx, page, "light"))
}

// Ref: TestPageEmulateMedia.java#shouldDefaultToLight  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldDefaultToLight(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, _, page := a2Setup(t)

	assert.Equal(t, true, a2MatchMedia(t, ctx, page, "light"))
	assert.Equal(t, false, a2MatchMedia(t, ctx, page, "dark"))
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: a2Ptr("dark")}))
	assert.Equal(t, true, a2MatchMedia(t, ctx, page, "dark"))
	assert.Equal(t, false, a2MatchMedia(t, ctx, page, "light"))
	must.NoError(page.EmulateMedia(ctx, &playwright.EmulateMediaOptions{ColorScheme: a2Ptr("no-override")}))
	assert.Equal(t, false, a2MatchMedia(t, ctx, page, "dark"))
	assert.Equal(t, true, a2MatchMedia(t, ctx, page, "light"))
}

// Ref: TestPageEmulateMedia.java#shouldWorkDuringNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldWorkDuringNavigation(t *testing.T) {
	t.Skip("not implemented: test body is an empty TODO (requires async API) in the Java reference")
}

// Ref: TestPageEmulateMedia.java#shouldWorkInPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldWorkInPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)

	for _, tc := range []struct{ scheme, light, dark string }{{"dark", "false", "true"}, {"light", "true", "false"}} {
		page := a2NewPage(t, &playwright.BrowserContextOptions{ColorScheme: a2Ptr(tc.scheme)})
		must.NoError(page.Goto(ctx, srv.EmptyPage()))
		popup, err := page.WaitForPopup(ctx, func() error {
			_, err := page.Evaluate(ctx, "url => { window.open(url); }", srv.EmptyPage())
			return err
		})
		must.NoError(err)
		assert.Equal(t, tc.light == "true", a2MatchMedia(t, ctx, popup, "light"), tc.scheme)
		assert.Equal(t, tc.dark == "true", a2MatchMedia(t, ctx, popup, "dark"), tc.scheme)
	}
}

// Ref: TestPageEmulateMedia.java#shouldWorkInCrossProcessIframe  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEmulateMedia.java
func TestPageEmulateMediaShouldWorkInCrossProcessIframe(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := a2NewPage(t, &playwright.BrowserContextOptions{ColorScheme: a2Ptr("dark")})

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	frame, err := a2AttachFrame(ctx, page, "frame1", srv.CrossProcessPrefix()+"/empty.html")
	must.NoError(err)
	res, err := frame.Evaluate(ctx, "() => matchMedia('(prefers-color-scheme: dark)').matches")
	must.NoError(err)
	assert.Equal(t, true, res)
}
