//go:build e2e

// E2E tests for: TestPageAddScriptTag.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestAddInitScriptSetsGlobalVariable verifies init script can set global variables.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptSetsGlobalVariable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.injectedVar = 'hello from init';`))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() => window.injectedVar`)
	must.NoError(err)
	is.Equal("hello from init", result)
}

// TestAddInitScriptRunsBeforePageScript verifies init script runs before inline page scripts.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptRunsBeforePageScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.order = window.order || []; window.order.push('init');`))

	must.NoError(page.SetContent(ctx, `
		<script>window.order = window.order || []; window.order.push('page');</script>
	`))

	result, err := page.Evaluate(ctx, `() => window.order`)
	must.NoError(err)

	arr, ok := result.([]any)
	is.True(ok)
	is.Len(arr, 2)
	is.Equal("init", arr[0])
	is.Equal("page", arr[1])
}

// TestAddInitScriptRunsOnEachNavigation verifies init script runs after each navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptRunsOnEachNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.initCount = (window.initCount || 0) + 1;`))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	count1, err := page.Evaluate(ctx, `() => window.initCount`)
	must.NoError(err)
	is.Equal(float64(1), count1)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	count2, err := page.Evaluate(ctx, `() => window.initCount`)
	must.NoError(err)
	is.Equal(float64(1), count2) // resets each navigation
}

// TestAddInitScriptCanModifyNavigator verifies init script can modify navigator.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptCanModifyNavigator(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `
		Object.defineProperty(navigator, 'language', {
			get: () => 'fr-FR'
		});
	`))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	lang, err := page.Evaluate(ctx, `() => navigator.language`)
	must.NoError(err)
	is.Equal("fr-FR", lang)
}

// TestAddInitScriptSetsGlobalBeforePageScript verifies init script runs before page scripts.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptSetsGlobalBeforePageScriptExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__initRan = true;`))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() => window.__initRan`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestAddInitScriptOverridesNavigator verifies init script can override navigator properties.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptOverridesNavigatorExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `
		Object.defineProperty(navigator, 'languages', {
			get: () => ['fr-FR'],
			configurable: true,
		});
	`))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	langs, err := page.Evaluate(ctx, `() => navigator.languages`)
	must.NoError(err)
	must.NotNil(langs)
	list, ok := langs.([]any)
	is.True(ok)
	is.Len(list, 1)
	is.Equal("fr-FR", list[0])
}

// TestAddInitScriptRunsOnEveryNavigation verifies init script persists across navigations.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptRunsOnEveryNavigationExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	counter := 0
	must.NoError(page.AddInitScript(ctx, `window.__count = (window.__count || 0) + 1;`))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	result1, err := page.Evaluate(ctx, `() => window.__count`)
	must.NoError(err)
	is.Equal(float64(1), result1)
	counter++

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	result2, err := page.Evaluate(ctx, `() => window.__count`)
	must.NoError(err)
	is.Equal(float64(1), result2) // resets each navigation
	counter++

	is.Equal(2, counter)
}

// TestAddInitScriptMultipleScripts verifies multiple init scripts run in order.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptMultipleScripts(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__log = [];`))
	must.NoError(page.AddInitScript(ctx, `window.__log.push('second');`))
	must.NoError(page.AddInitScript(ctx, `window.__log.push('third');`))
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	result, err := page.Evaluate(ctx, `() => window.__log`)
	must.NoError(err)
	must.NotNil(result)
	list, ok := result.([]any)
	is.True(ok)
	is.GreaterOrEqual(len(list), 2)
}

// TestAddInitScriptRunsBeforePageScript verifies init script runs before any page scripts.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptRunsBeforePageScript2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__initRan = true;`))

	srv.ServeWithBody("/init-check", "text/html", `
		<script>
			window.__initWasSet = window.__initRan === true;
		</script>
	`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/init-check"))

	result, err := page.Evaluate(ctx, `() => window.__initWasSet`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestAddInitScriptOverridesNavigatorProperty verifies init script can override navigator props.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptOverridesNavigatorProperty2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `
		Object.defineProperty(navigator, 'language', {
			get: () => 'es-MX',
		});
	`))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	lang, err := page.Evaluate(ctx, `() => navigator.language`)
	must.NoError(err)
	is.Equal("es-MX", lang)
}

// TestAddInitScriptPersistsAcrossNavigation verifies init script runs on every navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptPersistsAcrossNavigation2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__visits = (window.__visits || 0) + 1;`))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	v1, err := page.Evaluate(ctx, `() => window.__visits`)
	must.NoError(err)
	is.Equal(float64(1), v1)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	v2, err := page.Evaluate(ctx, `() => window.__visits`)
	must.NoError(err)
	is.Equal(float64(1), v2) // resets each navigation
}

// TestAddInitScriptSetsWindowPropertyEx4 verifies init script can set window properties.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptSetsWindowPropertyEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__testProp = 'initialized';`))

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => window.__testProp`)
	must.NoError(err)
	is.Equal("initialized", result)
}

// TestAddInitScriptModifiesNavigatorEx4 verifies init script can modify navigator.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptModifiesNavigatorEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `
		Object.defineProperty(navigator, 'platform', {
			get: () => 'Win64',
		});
	`))

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	platform, err := page.Evaluate(ctx, `() => navigator.platform`)
	must.NoError(err)
	is.Equal("Win64", platform)
}

// TestAddInitScriptWithArrayEx4 verifies init script can initialize arrays on window.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptWithArrayEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__log = [];`))

	must.NoError(page.SetContent(ctx, `<script>window.__log.push('entry1');</script>`))

	result, err := page.Evaluate(ctx, `() => window.__log.length`)
	must.NoError(err)
	is.Equal(float64(1), result)
}

// TestAddInitScriptRunsBeforePageScriptEx5 verifies init script runs before page scripts.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptRunsBeforePageScriptEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__initOrder = 'init';`))

	srv.ServeWithBody("/order", "text/html", `
		<script>window.__initOrder += '-page';</script>
	`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/order"))

	result, err := page.Evaluate(ctx, `() => window.__initOrder`)
	must.NoError(err)
	is.Equal("init-page", result)
}

// TestAddInitScriptRunsOnEveryNavigationEx5 verifies init script persists across navigations.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptRunsOnEveryNavigationEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__nav = true;`))

	srv.ServeWithBody("/nav1", "text/html", `<p>Nav 1</p>`)
	srv.ServeWithBody("/nav2", "text/html", `<p>Nav 2</p>`)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/nav1"))
	r1, err := page.Evaluate(ctx, `() => window.__nav`)
	must.NoError(err)
	is.Equal(true, r1)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/nav2"))
	r2, err := page.Evaluate(ctx, `() => window.__nav`)
	must.NoError(err)
	is.Equal(true, r2)
}

// TestAddInitScriptMultipleScriptsEx5 verifies multiple init scripts all run.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptMultipleScriptsEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.AddInitScript(ctx, `window.__s1 = true;`))
	must.NoError(page.AddInitScript(ctx, `window.__s2 = true;`))

	srv.ServeWithBody("/multi", "text/html", `<p>Multi</p>`)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/multi"))

	r1, err := page.Evaluate(ctx, `() => window.__s1`)
	must.NoError(err)
	is.Equal(true, r1)

	r2, err := page.Evaluate(ctx, `() => window.__s2`)
	must.NoError(err)
	is.Equal(true, r2)
}

// TestAddInitScriptRunsBeforePageScriptEx verifies init script runs before page scripts.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptRunsBeforePageScriptEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/init-order", "text/html", `
		<html>
		<head><script>window.__pageOrder = (window.__initOrder || 0) + 1;</script></head>
		<body></body>
		</html>
	`)

	must.NoError(page.AddInitScript(ctx, `window.__initOrder = 0;`))
	must.NoError(page.Goto(ctx, srv.Prefix()+"/init-order"))

	result, err := page.Evaluate(ctx, `() => window.__pageOrder`)
	must.NoError(err)
	is.Equal(float64(1), result)
}

// TestAddInitScriptMultipleEx verifies multiple init scripts all execute.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptMultipleEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/multi-init", "text/html", `<html><body></body></html>`)

	must.NoError(page.AddInitScript(ctx, `window.__a = 1;`))
	must.NoError(page.AddInitScript(ctx, `window.__b = 2;`))
	must.NoError(page.AddInitScript(ctx, `window.__c = window.__a + window.__b;`))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/multi-init"))

	result, err := page.Evaluate(ctx, `() => window.__c`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

// TestAddInitScriptNavigationPersistsEx verifies init script runs on every navigation.
// Extra: no equivalent method in playwright-java; related class: TestPageAddScriptTag.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageAddScriptTag.java
func TestAddInitScriptNavigationPersistsEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	srv := testserver.New(t)

	srv.ServeWithBody("/nav1", "text/html", `<html><body>Nav 1</body></html>`)
	srv.ServeWithBody("/nav2", "text/html", `<html><body>Nav 2</body></html>`)

	must.NoError(page.AddInitScript(ctx, `window.__init = 'initialized';`))

	must.NoError(page.Goto(ctx, srv.Prefix()+"/nav1"))
	r1, err := page.Evaluate(ctx, `() => window.__init`)
	must.NoError(err)
	is.Equal("initialized", r1)

	must.NoError(page.Goto(ctx, srv.Prefix()+"/nav2"))
	r2, err := page.Evaluate(ctx, `() => window.__init`)
	must.NoError(err)
	is.Equal("initialized", r2)
}
