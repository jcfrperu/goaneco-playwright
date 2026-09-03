//go:build e2e

package e2e

import (
	"strings"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tagSelectorScript is a custom selector engine that maps selectors to HTML tag names.
const tagSelectorScript = `{
  create(root, target) {
    return target.nodeName;
  },
  query(root, selector) {
    return root.querySelector(selector);
  },
  queryAll(root, selector) {
    return Array.from(root.querySelectorAll(selector));
  }
}`

// TestSelectorsRegisterShouldWork verifies that custom selector engines work for
// evalOnSelector and evalOnSelectorAll, and that engine names are case-sensitive.
// Ref: TestSelectorsRegister.java#shouldWork
func TestSelectorsRegisterShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	must.NoError(bCtx.RegisterSelectorEngine(ctx, "tag", tagSelectorScript))
	must.NoError(bCtx.RegisterSelectorEngine(ctx, "tag2", tagSelectorScript))

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<div><span></span></div><div></div>`))

	nodeName, err := page.EvalOnSelector(ctx, "tag=DIV", "e => e.nodeName")
	must.NoError(err)
	is.Equal("DIV", nodeName)

	nodeName, err = page.EvalOnSelector(ctx, "tag=SPAN", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SPAN", nodeName)

	count, err := page.EvalOnSelectorAll(ctx, "tag=DIV", "es => es.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	nodeName, err = page.EvalOnSelector(ctx, "tag2=DIV", "e => e.nodeName")
	must.NoError(err)
	is.Equal("DIV", nodeName)

	nodeName, err = page.EvalOnSelector(ctx, "tag2=SPAN", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SPAN", nodeName)

	count, err = page.EvalOnSelectorAll(ctx, "tag2=DIV", "es => es.length")
	must.NoError(err)
	is.Equal(float64(2), count)

	// Selector engine names are case-sensitive.
	_, err = page.QuerySelector(ctx, "tAG=DIV")
	is.Error(err)
	is.True(strings.Contains(err.Error(), "Unknown engine"), "expected 'Unknown engine' in error, got: %s", err.Error())
}

// TestSelectorsRegisterInMainAndIsolatedWorld verifies that a selector engine registered
// without contentScript runs in the main world (reading window variables) while one
// registered with contentScript=true runs in an isolated world and cannot see window variables.
// Ref: TestSelectorsRegister.java#shouldWorkInMainAndIsolatedWorld
func TestSelectorsRegisterInMainAndIsolatedWorld(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	createDummySelector := `{
  create(root, target) { },
  query(root, selector) {
    return window['__answer'];
  },
  queryAll(root, selector) {
    return window['__answer'] ? [window['__answer'], document.body, document.documentElement] : [];
  }
}`

	must.NoError(bCtx.RegisterSelectorEngine(ctx, "main", createDummySelector))
	must.NoError(bCtx.RegisterSelectorEngine(
		ctx,
		"isolated",
		createDummySelector,
		&playwright.RegisterSelectorEngineOptions{ContentScript: boolPtr(true)},
	))

	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.SetContent(ctx, `<div><span><section></section></span></div>`))
	_, err = page.Evaluate(ctx, "() => window['__answer'] = document.querySelector('span')")
	must.NoError(err)

	// Main-world engine can read window['__answer'].
	nodeName, err := page.EvalOnSelector(ctx, "main=ignored", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SPAN", nodeName)

	nodeName, err = page.EvalOnSelector(ctx, "css=div >> main=ignored", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SPAN", nodeName)

	answerDefined, err := page.EvalOnSelectorAll(ctx, "main=ignored", "es => window['__answer'] !== undefined")
	must.NoError(err)
	is.Equal(true, answerDefined)

	filtered, err := page.EvalOnSelectorAll(ctx, "main=ignored", "es => es.filter(e => e).length")
	must.NoError(err)
	is.Equal(float64(3), filtered)

	// Isolated-world engine cannot see window['__answer'] from the main world.
	el, err := page.QuerySelector(ctx, "isolated=ignored")
	must.NoError(err)
	is.Nil(el)

	el, err = page.QuerySelector(ctx, "css=div >> isolated=ignored")
	must.NoError(err)
	is.Nil(el)

	// $$eval always runs in main world to avoid per-element adoption.
	answerDefined, err = page.EvalOnSelectorAll(ctx, "isolated=ignored", "es => window['__answer'] !== undefined")
	must.NoError(err)
	is.Equal(true, answerDefined)

	filtered, err = page.EvalOnSelectorAll(ctx, "isolated=ignored", "es => es.filter(e => e).length")
	must.NoError(err)
	is.Equal(float64(3), filtered)

	// Chaining: at least one main-world engine forces evaluation in main world.
	nodeName, err = page.EvalOnSelector(ctx, "main=ignored >> isolated=ignored", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SPAN", nodeName)

	nodeName, err = page.EvalOnSelector(ctx, "isolated=ignored >> main=ignored", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SPAN", nodeName)

	// Custom engine can be chained with built-in css engine.
	nodeName, err = page.EvalOnSelector(ctx, "main=ignored >> css=section", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SECTION", nodeName)
}

// TestSelectorsRegisterHandlesErrors verifies that invalid engine names and duplicate
// registrations produce descriptive errors.
// Ref: TestSelectorsRegister.java#shouldHandleErrors
func TestSelectorsRegisterHandlesErrors(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	createDummySelector := `{
  create(root, target) {
    return target.nodeName;
  },
  query(root, selector) {
    return root.querySelector("dummy");
  },
  queryAll(root, selector) {
    return Array.from(root.querySelectorAll("dummy"));
  }
}`

	// Unregistered engine produces an unknown engine error at query time.
	bCtx1 := newContext(t)
	page1, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	_, err = page1.QuerySelector(ctx, "neverregister=ignored")
	is.Error(err)
	is.True(strings.Contains(err.Error(), "Unknown engine"), "expected 'Unknown engine' in error, got: %s", err.Error())

	// Engine name containing invalid character '$' is rejected.
	bCtx2 := newContext(t)
	err = bCtx2.RegisterSelectorEngine(ctx, "$", createDummySelector)
	is.Error(err)
	is.True(strings.Contains(err.Error(), "Selector engine name may only contain"), "expected name-validation error, got: %s", err.Error())

	// Registering a duplicate name is rejected.
	bCtx3 := newContext(t)
	must.NoError(bCtx3.RegisterSelectorEngine(ctx, "dummy", createDummySelector))
	must.NoError(bCtx3.RegisterSelectorEngine(ctx, "duMMy", createDummySelector))
	err = bCtx3.RegisterSelectorEngine(ctx, "dummy", createDummySelector)
	is.Error(err)
	is.True(strings.Contains(err.Error(), "selector engine has been already registered"), "expected already-registered error, got: %s", err.Error())

	// Built-in engine name 'css' is rejected.
	bCtx4 := newContext(t)
	err = bCtx4.RegisterSelectorEngine(ctx, "css", createDummySelector)
	is.Error(err)
	is.True(strings.Contains(err.Error(), "predefined selector engine"), "expected predefined-engine error, got: %s", err.Error())
}

// Ref: TestSelectorsRegister.java#shouldWorkWithPath
func TestSelectorsRegisterShouldWorkWithPath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)

	// Java registers from a file containing a section selector engine.
	// Here we inline the equivalent JS.
	const sectionSelectorScript = `({
  create(root, target) {},
  query(root, selector) { return root.querySelector('section'); },
  queryAll(root, selector) { return Array.from(root.querySelectorAll('section')); }
})`
	must.NoError(bCtx.RegisterSelectorEngine(ctx, "foo", sectionSelectorScript))
	must.NoError(page.SetContent(ctx, `<section></section>`))

	v, err := page.EvalOnSelector(ctx, "foo=whatever", "e => e.nodeName")
	must.NoError(err)
	is.Equal("SECTION", v)
}

// TestSelectorsRegisterAlreadyRegisteredError verifies that registering the same engine
// name twice on a fresh context produces the expected error.
// Ref: TestSelectorsRegister.java#shouldThrowAlreadyRegisteredErrorWhenRegistering
func TestSelectorsRegisterAlreadyRegisteredError(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx := newContext(t)

	must.NoError(bCtx.RegisterSelectorEngine(ctx, "alreadyRegistered", tagSelectorScript))

	err := bCtx.RegisterSelectorEngine(ctx, "alreadyRegistered", tagSelectorScript)
	is.Error(err)
	is.True(strings.Contains(err.Error(), "selector engine has been already registered"), "expected already-registered error, got: %s", err.Error())
}
