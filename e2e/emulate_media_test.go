//go:build e2e

// Page.EmulateMedia E2E tests.
// Migration of: TestPageEmulateMedia.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

func TestEmulateMediaType(t *testing.T) {
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

func TestEmulateReducedMotion(t *testing.T) {
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

func TestEmulateMediaNoOpts(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// EmulateMedia without options should not error.
	err := page.EmulateMedia(ctx)
	must.NoError(err, "EmulateMedia() with no opts failed")
}

func TestEmulateMediaChangesActualColorsInCSS(t *testing.T) {
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

func TestEmulateContrast(t *testing.T) {
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
