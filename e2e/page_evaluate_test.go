//go:build e2e

// E2E tests for: TestPageEvaluate.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// WaitForFunction returns a JSHandle for the result object.
	handle, err := page.WaitForFunction(ctx, "() => ({x: 10, y: 20})", nil)
	must.NoError(err)
	must.NotNil(handle, "WaitForFunction failed")

	// EvaluateHandle returns a new JSHandle from an expression applied to this handle.
	xHandle, err := handle.EvaluateHandle(ctx, "obj => obj.x")
	must.NoError(err, "JSHandle.EvaluateHandle failed")
	must.NotNil(xHandle, "EvaluateHandle returned nil")

	val, err := xHandle.JSONValue(ctx)
	must.NoError(err, "JSONValue on EvaluateHandle result failed")
	is.Equal(float64(10), val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetProperty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({answer: 42, label: 'hello'})", nil)
	must.NoError(err)
	must.NotNil(handle, "WaitForFunction failed")

	// GetProperty returns a JSHandle for the named property.
	prop, err := handle.GetProperty(ctx, "answer")
	must.NoError(err, "GetProperty('answer') failed")

	val, err := prop.JSONValue(ctx)
	must.NoError(err, "JSONValue on property failed")
	is.Equal(float64(42), val)

	// GetProperty on a string property.
	labelProp, err := handle.GetProperty(ctx, "label")
	must.NoError(err, "GetProperty('label') failed")
	labelVal, err := labelProp.JSONValue(ctx)
	must.NoError(err, "JSONValue on label failed")
	is.Equal("hello", labelVal)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleDispose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({x: 1})", nil)
	must.NoError(err)
	must.NotNil(handle, "WaitForFunction failed")

	// Dispose should release the handle without error.
	err = handle.Dispose(ctx)
	must.NoError(err, "JSHandle.Dispose failed")
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body></body>`))

	handle, err := page.WaitForFunction(ctx, `() => ({ foo: 'bar', count: 42 })`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)

	m, ok := val.(map[string]any)
	is.True(ok, "expected map, got %T: %v", val, val)
	is.Equal("bar", m["foo"])
	is.Equal(float64(42), m["count"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueForPrimitive(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body></body>`))

	handle, err := page.WaitForFunction(ctx, `() => 42`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(42), val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateWithHandleAsArg(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="foo">hello</div>`))

	handle, err := page.WaitForFunction(ctx, `() => document.querySelector('div#foo')`, nil)
	must.NoError(err)

	result, err := handle.Evaluate(ctx, `el => el.textContent`)
	must.NoError(err)
	is.Equal("hello", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateChained(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body></body>`))

	handle, err := page.WaitForFunction(ctx, `() => ({ key: 'value' })`, nil)
	must.NoError(err)

	result, err := handle.Evaluate(ctx, `obj => obj.key`)
	must.NoError(err)
	is.Equal("value", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetPropertyReturnsHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body></body>`))

	handle, err := page.WaitForFunction(ctx, `() => ({ greeting: 'hello', count: 5 })`, nil)
	must.NoError(err)

	propHandle, err := handle.GetProperty(ctx, "greeting")
	must.NoError(err)
	must.NotNil(propHandle)

	val, err := propHandle.JSONValue(ctx)
	must.NoError(err)
	is.Equal("hello", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleDisposeIsIdempotent(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<body></body>`))

	handle, err := page.WaitForFunction(ctx, `() => ({ value: 1 })`, nil)
	must.NoError(err)

	must.NoError(handle.Dispose(ctx))

	must.NoError(handle.Dispose(ctx))
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueNumber(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => 42", nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(42), val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateCreatesHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">content</div>`))

	handle, err := page.WaitForFunction(ctx, "() => document.getElementById('el')", nil)
	must.NoError(err)
	must.NotNil(handle)

	result, err := handle.Evaluate(ctx, `el => el.id`)
	must.NoError(err)
	is.Equal("el", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetPropertyReturnsValueExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({name: 'test', value: 99})", nil)
	must.NoError(err)
	must.NotNil(handle)

	nameProp, err := handle.GetProperty(ctx, "name")
	must.NoError(err)
	must.NotNil(nameProp)

	nameVal, err := nameProp.JSONValue(ctx)
	must.NoError(err)
	is.Equal("test", nameVal)

	valProp, err := handle.GetProperty(ctx, "value")
	must.NoError(err)

	numVal, err := valProp.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(99), numVal)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueString(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => 'just a string'", nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal("just a string", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateChainingExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({a: {b: 99}})", nil)
	must.NoError(err)
	must.NotNil(handle)

	innerHandle, err := handle.EvaluateHandle(ctx, `obj => obj.a`)
	must.NoError(err)
	must.NotNil(innerHandle)

	result, err := innerHandle.JSONValue(ctx)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal(float64(99), m["b"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueNullExtra(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => null`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Nil(val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueBoolean(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => true`, nil)
	must.NoError(err)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(true, val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateReturnsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => 42`, nil)
	must.NoError(err)

	result, err := handle.Evaluate(ctx, `val => val * 2`)
	must.NoError(err)
	is.Equal(float64(84), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetPropertyReturnsValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => ({ key: 'value' })`, nil)
	must.NoError(err)

	prop, err := handle.GetProperty(ctx, "key")
	must.NoError(err)
	must.NotNil(prop)

	val, err := prop.JSONValue(ctx)
	must.NoError(err)
	is.Equal("value", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateHandleChaining(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => ({ x: 10 })`, nil)
	must.NoError(err)

	child, err := handle.EvaluateHandle(ctx, `obj => obj`)
	must.NoError(err)
	must.NotNil(child)

	val, err := child.GetProperty(ctx, "x")
	must.NoError(err)
	must.NotNil(val)

	jsonVal, err := val.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(10), jsonVal)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateHandleChain(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">hello</div>`))

	handle, err := page.WaitForFunction(ctx, `() => document.getElementById('el')`, nil)
	must.NoError(err)
	must.NotNil(handle)

	inner, err := handle.EvaluateHandle(ctx, `el => el.firstChild`)
	must.NoError(err)
	must.NotNil(inner)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleDisposeDoesNotError(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => 1`, nil)
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.Dispose(ctx))
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateReturnsPrimitive(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => ({name: "alice", age: 30})`, nil)
	must.NoError(err)
	must.NotNil(handle)

	name, err := handle.GetProperty(ctx, "name")
	must.NoError(err)
	must.NotNil(name)

	val, err := name.JSONValue(ctx)
	must.NoError(err)
	is.Equal("alice", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetPropertyAge(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => ({age: 25})`, nil)
	must.NoError(err)
	must.NotNil(handle)

	prop, err := handle.GetProperty(ctx, "age")
	must.NoError(err)
	must.NotNil(prop)

	val, err := prop.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(25), val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueStringExtra4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	handle, err := page.WaitForFunction(ctx, `() => "hello world"`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal("hello world", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateReturnsBoolEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => true", nil)
	must.NoError(err)
	must.NotNil(handle)

	result, err := handle.Evaluate(ctx, `v => v === true`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateReturnsNumberEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => 99", nil)
	must.NoError(err)
	must.NotNil(handle)

	result, err := handle.Evaluate(ctx, `v => v + 1`)
	must.NoError(err)
	is.Equal(float64(100), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueArrayEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => [10, 20, 30]", nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	must.NotNil(val)
	arr := val.([]any)
	is.Len(arr, 3)
	is.Equal(float64(20), arr[1])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetPropertyFromObjectEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({ city: 'Lima' })", nil)
	must.NoError(err)
	must.NotNil(handle)

	prop, err := handle.GetProperty(ctx, "city")
	must.NoError(err)
	must.NotNil(prop)

	val, err := prop.JSONValue(ctx)
	must.NoError(err)
	is.Equal("Lima", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleDisposeNoErrorEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => window", nil)
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.Dispose(ctx))
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetPropertyEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, `() => ({ name: "test", value: 42 })`, nil)
	must.NoError(err)
	must.NotNil(handle)

	nameHandle, err := handle.GetProperty(ctx, "name")
	must.NoError(err)
	must.NotNil(nameHandle)

	val, err := nameHandle.JSONValue(ctx)
	must.NoError(err)
	is.Equal("test", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleGetPropertyBEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, `() => ({ a: 1, b: 99 })`, nil)
	must.NoError(err)
	must.NotNil(handle)

	bHandle, err := handle.GetProperty(ctx, "b")
	must.NoError(err)
	must.NotNil(bHandle)

	val, err := bHandle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(99), val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleJSONValueNumberEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, `() => 3.14`, nil)
	must.NoError(err)
	must.NotNil(handle)

	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.InDelta(3.14, val.(float64), 0.001)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestJSHandleEvaluateEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, `() => [1, 2, 3, 4, 5]`, nil)
	must.NoError(err)
	must.NotNil(handle)

	result, err := handle.Evaluate(ctx, `arr => arr.length`)
	must.NoError(err)
	is.Equal(float64(5), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestLocatorEvaluateAllEmptyReturnsEmptyArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>no items here</div>`))

	result, err := page.Locator(".no-match").EvaluateAll(ctx, `els => els.length`)
	must.NoError(err)
	is.Equal(float64(0), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateReturnNull(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => null`)
	must.NoError(err)
	is.Nil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateReturnUndefined(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => undefined`)
	must.NoError(err)
	is.Nil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateArithmeticExpression(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => 2 + 3`)
	must.NoError(err)
	is.Equal(float64(5), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateReturnsStringValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => 'hello world'`)
	must.NoError(err)
	is.Equal("hello world", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateReturnsObject(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => ({key: 'value', num: 42})`)
	must.NoError(err)

	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("value", m["key"])
	is.Equal(float64(42), m["num"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateReturnsArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [1, 2, 3]`)
	must.NoError(err)

	arr, ok := result.([]any)
	is.True(ok)
	is.Len(arr, 3)
	is.Equal(float64(1), arr[0])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateWithDOMElementArg(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="myDiv">hello</div>`))

	result, err := page.Evaluate(ctx, `() => document.getElementById('myDiv').textContent`)
	must.NoError(err)
	is.Equal("hello", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluatePromise(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => Promise.resolve(42)`)
	must.NoError(err)
	is.Equal(float64(42), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWorksWithComplexArgs(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `({a, b}) => a + b`, map[string]any{"a": 3, "b": 4})
	must.NoError(err)
	is.Equal(float64(7), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWorksWithArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `arr => arr.reduce((a, b) => a + b, 0)`, []int{1, 2, 3})
	must.NoError(err)
	is.Equal(float64(6), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateReturnsList(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => [1, 2, 3]`)
	must.NoError(err)
	must.NotNil(result)

	arr, ok := result.([]any)
	is.True(ok)
	is.Len(arr, 3)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateThrowingReturnsError(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	_, err := page.Evaluate(ctx, `() => { throw new Error('boom'); }`)
	is.Error(err)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateReturnNestedObject(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => ({ outer: { inner: 42 } })`)
	must.NoError(err)
	must.NotNil(result)

	m, ok := result.(map[string]any)
	is.True(ok)
	outer, ok := m["outer"].(map[string]any)
	is.True(ok)
	is.Equal(float64(42), outer["inner"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWithNullArg(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `arg => arg === null`, nil)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateManipulatesDOM(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"></div>`))

	_, err := page.Evaluate(ctx, `() => {
		const p = document.createElement('p');
		p.id = 'added';
		p.textContent = 'dynamically added';
		document.getElementById('container').appendChild(p);
	}`)
	must.NoError(err)

	text, err := page.Locator("#added").InnerText(ctx)
	must.NoError(err)
	is.Equal("dynamically added", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateReadsWindowProperty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	_, err := page.Evaluate(ctx, `() => { window.myProp = 'hello'; }`)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => window.myProp`)
	must.NoError(err)
	is.Equal("hello", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateReturnsUndefinedForMissingProp(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `() => window.noSuchProperty`)
	must.NoError(err)
	is.Nil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWithStringConcatenation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `(greeting) => greeting + ' world'`, "hello")
	must.NoError(err)
	is.Equal("hello world", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateCanQueryDocument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li id="a">a</li>
			<li id="b">b</li>
		</ul>
	`))

	result, err := page.Evaluate(ctx, `() => document.querySelectorAll('li').length`)
	must.NoError(err)
	is.Equal(float64(2), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWithFloatArg(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	result, err := page.Evaluate(ctx, `x => x * 2`, 3.14)
	must.NoError(err)
	is.InDelta(6.28, result, 0.001)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateReturnsArrayEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	val, err := page.Evaluate(ctx, `() => [1, 2, 3]`)
	must.NoError(err)
	slice, ok := val.([]any)
	is.True(ok)
	is.Len(slice, 3)
	is.Equal(float64(1), slice[0])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateObjectWithMixedTypes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	val, err := page.Evaluate(ctx, `() => ({name:"alice",age:30,active:true})`)
	must.NoError(err)
	m, ok := val.(map[string]any)
	is.True(ok)
	is.Equal("alice", m["name"])
	is.Equal(float64(30), m["age"])
	is.Equal(true, m["active"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateModifiesDOMAndReadsBack(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="out"></div>`))

	_, err := page.Evaluate(ctx, `() => { document.getElementById('out').textContent = 'modified'; }`)
	must.NoError(err)

	text, err := page.Locator("#out").InnerText(ctx)
	must.NoError(err)
	is.Equal("modified", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluatePromiseEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	val, err := page.Evaluate(ctx, `() => Promise.resolve(42)`)
	must.NoError(err)
	is.Equal(float64(42), val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateStringConcat(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	val, err := page.Evaluate(ctx, `() => "hello" + " " + "world"`)
	must.NoError(err)
	is.Equal("hello world", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateNestedFunction(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div></div>`))

	val, err := page.Evaluate(ctx, `() => {
		function double(n) { return n * 2; }
		return double(21);
	}`)
	must.NoError(err)
	is.Equal(float64(42), val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateNullReturnsNilEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => null`)
	must.NoError(err)
	is.Nil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateUndefinedReturnsNilEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => undefined`)
	must.NoError(err)
	is.Nil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateBooleanTrueEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => true`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateBooleanFalseEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => false`)
	must.NoError(err)
	is.Equal(false, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateObjectPropertyEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => ({ name: 'Claude', version: 4 })`)
	must.NoError(err)
	must.NotNil(result)
	m := result.(map[string]any)
	is.Equal("Claude", m["name"])
	is.Equal(float64(4), m["version"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateArrowFunctionWithArgEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `(n) => n * 5`, 7)
	must.NoError(err)
	is.Equal(float64(35), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateDateObjectEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => new Date(2024, 0, 1).getFullYear()`)
	must.NoError(err)
	is.Equal(float64(2024), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateRegexMatchEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => /hello/i.test('Hello World')`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateChainedCallsEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => "hello world".split(" ").length`)
	must.NoError(err)
	is.Equal(float64(2), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateArrayMapEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [1, 2, 3].map(n => n * 2)`)
	must.NoError(err)
	must.NotNil(result)
	arr := result.([]any)
	is.Len(arr, 3)
	is.Equal(float64(6), arr[2])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateTypeofCheckEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => typeof "hello"`)
	must.NoError(err)
	is.Equal("string", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateConditionalEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `(x) => x > 5 ? 'big' : 'small'`, 10)
	must.NoError(err)
	is.Equal("big", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateDOMQueryEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="item">A</li>
			<li class="item">B</li>
			<li class="item">C</li>
		</ul>
	`))

	result, err := page.Evaluate(ctx, `() => document.querySelectorAll('.item').length`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateDocumentReadyStateEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>Ready</div>`))

	result, err := page.Evaluate(ctx, `() => document.readyState`)
	must.NoError(err)
	is.Equal("complete", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateLocationHrefEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>page</div>`))

	result, err := page.Evaluate(ctx, `() => typeof window.location.href`)
	must.NoError(err)
	is.Equal("string", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateWindowInnerSizeEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>size</div>`))

	width, err := page.Evaluate(ctx, `() => window.innerWidth`)
	must.NoError(err)
	is.Greater(width.(float64), float64(0))

	height, err := page.Evaluate(ctx, `() => window.innerHeight`)
	must.NoError(err)
	is.Greater(height.(float64), float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateCreateElementEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"></div>`))

	_, err := page.Evaluate(ctx, `() => {
		const el = document.createElement('p');
		el.id = 'created';
		el.textContent = 'Created dynamically';
		document.getElementById('container').appendChild(el);
	}`)
	must.NoError(err)

	text, err := page.Locator("#created").InnerText(ctx)
	must.NoError(err)
	is.Equal("Created dynamically", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateReturnsMapEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => ({ key: "value", num: 1 })`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("value", m["key"])
	is.Equal(float64(1), m["num"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateNullReturnEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => null`)
	must.NoError(err)
	is.Nil(result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWithNestedObjectEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => ({ outer: { inner: 42 } })`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	is.True(ok)
	outer, ok := m["outer"].(map[string]any)
	is.True(ok)
	is.Equal(float64(42), outer["inner"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateArithmeticEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => 2 ** 10`)
	must.NoError(err)
	is.Equal(float64(1024), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateStringConcatEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => "hello" + " " + "world"`)
	must.NoError(err)
	is.Equal("hello world", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWithArgArrayEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `(arr) => arr.reduce((a,b) => a+b, 0)`, []int{1, 2, 3, 4})
	must.NoError(err)
	is.Equal(float64(10), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateRegexEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => /hello/.test('hello world')`)
	must.NoError(err)
	is.Equal(true, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSetTimeoutEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => new Promise(resolve => setTimeout(resolve, 10))`)
	must.NoError(err)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateStringMethodEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => 'Hello World'.toLowerCase()`)
	must.NoError(err)
	is.Equal("hello world", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateArrayMethodEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [1,2,3,4,5].filter(n => n % 2 === 0)`)
	must.NoError(err)
	arr, ok := result.([]any)
	is.True(ok)
	is.Equal(2, len(arr))
	is.Equal(float64(2), arr[0])
	is.Equal(float64(4), arr[1])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateTryCatchEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => {
		try {
			throw new Error('test error');
		} catch(e) {
			return 'caught: ' + e.message;
		}
	}`)
	must.NoError(err)
	is.Equal("caught: test error", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSetPropertyEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Original</div>`))

	_, err := page.Evaluate(ctx, `() => { document.getElementById('d').textContent = 'Modified'; }`)
	must.NoError(err)

	text, err := page.Locator("#d").InnerText(ctx)
	must.NoError(err)
	is.Equal("Modified", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateAddElementEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"></ul>`))

	_, err := page.Evaluate(ctx, `() => {
		var li = document.createElement('li');
		li.id = 'added';
		li.textContent = 'New item';
		document.getElementById('list').appendChild(li);
	}`)
	must.NoError(err)

	text, err := page.Locator("#added").InnerText(ctx)
	must.NoError(err)
	is.Equal("New item", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateRemoveElementEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div><span id="rm">Remove me</span></div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('rm').remove()`)
	must.NoError(err)

	count, err := page.Locator("#rm").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWindowPropertiesEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetViewportSize(ctx, 1024, 768))

	result, err := page.Evaluate(ctx, `() => ({ w: window.innerWidth, h: window.innerHeight })`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal(float64(1024), m["w"])
	is.Equal(float64(768), m["h"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateDocumentURLEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))

	result, err := page.Evaluate(ctx, `() => document.URL`)
	must.NoError(err)
	is.NotEmpty(result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSetAttributeEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').setAttribute('data-test', 'value123')`)
	must.NoError(err)

	attr, err := page.Locator("#d").GetAttribute(ctx, "data-test")
	must.NoError(err)
	must.NotNil(attr)
	is.Equal("value123", *attr)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateCountChildrenEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul id="list">
			<li>A</li>
			<li>B</li>
			<li>C</li>
		</ul>
	`))

	result, err := page.Evaluate(ctx, `() => document.getElementById('list').childElementCount`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateScrollPositionEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div style="height:3000px;"></div>`))
	_, err := page.Evaluate(ctx, `() => window.scrollTo(0, 500)`)
	must.NoError(err)

	scroll, err := page.Evaluate(ctx, `() => window.scrollY`)
	must.NoError(err)
	is.Equal(float64(500), scroll)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateLocalStorageEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p>page</p>`))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('key', 'value')`)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => localStorage.getItem('key')`)
	must.NoError(err)
	is.Equal("value", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateObjectKeysEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => Object.keys({a:1, b:2, c:3})`)
	must.NoError(err)
	arr, ok := result.([]any)
	is.True(ok)
	is.Equal(3, len(arr))
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateJSONParseEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => JSON.parse('{"key": "value"}')`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	is.True(ok)
	is.Equal("value", m["key"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateJSONStringifyEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => JSON.stringify({a: 1, b: "two"})`)
	must.NoError(err)
	str, ok := result.(string)
	is.True(ok)
	is.Contains(str, `"a":1`)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateMapFilterReduceEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [1,2,3,4,5].map(x => x*2).filter(x => x > 4).reduce((a,b) => a+b, 0)`)
	must.NoError(err)
	is.Equal(float64(24), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluatePromiseChainEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => Promise.resolve(42).then(v => v * 2)`)
	must.NoError(err)
	is.Equal(float64(84), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateAsyncAwaitEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `async () => {
		await new Promise(resolve => setTimeout(resolve, 10));
		return 'done';
	}`)
	must.NoError(err)
	is.Equal("done", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWindowLocationEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body>Page</body></html>`))

	href, err := page.Evaluate(ctx, `() => window.location.href`)
	must.NoError(err)
	is.NotEmpty(href)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateDocumentCookieEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body>Page</body></html>`))

	_, err := page.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateArrayMapEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [1, 2, 3].map(x => x * 2)`)
	must.NoError(err)
	arr, ok := result.([]interface{})
	is.True(ok)
	is.Len(arr, 3)
	is.Equal(float64(6), arr[2])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateObjectKeysEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => Object.keys({a: 1, b: 2, c: 3})`)
	must.NoError(err)
	keys, ok := result.([]interface{})
	is.True(ok)
	is.Len(keys, 3)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateTemplateStringEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => `Hello ${'World'}`")
	must.NoError(err)
	is.Equal("Hello World", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSetAttributeEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Content</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').setAttribute('data-test', 'hello')`)
	must.NoError(err)

	val, err := page.Locator("#d").GetAttribute(ctx, "data-test")
	must.NoError(err)
	is.Equal("hello", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateCreateElementEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="container"></div>`))

	_, err := page.Evaluate(ctx, `() => {
		const btn = document.createElement('button');
		btn.id = 'new-btn';
		btn.textContent = 'Created';
		document.getElementById('container').appendChild(btn);
	}`)
	must.NoError(err)

	text, err := page.Locator("#new-btn").TextContent(ctx)
	must.NoError(err)
	is.Equal("Created", text)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateRemoveElementEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="del">To be removed</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('del').remove()`)
	must.NoError(err)

	count, err := page.Locator("#del").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateToggleClassEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="visible">Content</div>`))

	_, err := page.Evaluate(ctx, `() => document.getElementById('d').classList.toggle('visible')`)
	must.NoError(err)

	cls, err := page.Evaluate(ctx, `() => document.getElementById('d').className`)
	must.NoError(err)
	is.Equal("", cls)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateStringPaddingEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => '5'.padStart(3, '0')`)
	must.NoError(err)
	is.Equal("005", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluatePromiseResolveEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => Promise.resolve(42)`)
	must.NoError(err)
	is.Equal(float64(42), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateAsyncAwaitEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `async () => { return 'async result'; }`)
	must.NoError(err)
	is.Equal("async result", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateDestructuringEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => { const [a, b] = [10, 20]; return a + b; }`)
	must.NoError(err)
	is.Equal(float64(30), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSpreadOperatorEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [...[1, 2, 3], ...[4, 5]].length`)
	must.NoError(err)
	is.Equal(float64(5), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateOptionalChainingEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => { const obj = null; return obj?.name ?? 'default'; }`)
	must.NoError(err)
	is.Equal("default", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSetLocalStorageEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('key', 'value123')`)
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => localStorage.getItem('key')`)
	must.NoError(err)
	is.Equal("value123", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateGetSessionStorageEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => sessionStorage.setItem('session-key', 'session-val')`)
	must.NoError(err)

	val, err := page.Evaluate(ctx, `() => sessionStorage.getItem('session-key')`)
	must.NoError(err)
	is.Equal("session-val", val)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateJSONParseEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => JSON.parse('{"name":"Alice","age":30}').name`)
	must.NoError(err)
	is.Equal("Alice", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateJSONStringifyEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => JSON.stringify({x: 1, y: 2})`)
	must.NoError(err)
	is.Equal(`{"x":1,"y":2}`, result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateRegexMatchEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => 'hello world'.match(/world/)[0]`)
	must.NoError(err)
	is.Equal("world", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSetIntervalEx17(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body></body></html>`))

	_, err := page.Evaluate(ctx, `() => {
		window.__ticks = 0;
		const id = setInterval(() => window.__ticks++, 10);
		setTimeout(() => clearInterval(id), 150);
	}`)
	must.NoError(err)

	must.NoError(page.WaitForTimeout(ctx, 200))

	ticks, err := page.Evaluate(ctx, `() => window.__ticks`)
	must.NoError(err)
	is.Greater(ticks.(float64), float64(0))
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateNullCoalescingEx17(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => { const x = null; return x ?? 'fallback'; }`)
	must.NoError(err)
	is.Equal("fallback", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateArrayFilterEx17(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [1, 2, 3, 4, 5].filter(x => x % 2 === 0)`)
	must.NoError(err)
	arr, ok := result.([]interface{})
	is.True(ok)
	is.Len(arr, 2)
	is.Equal(float64(2), arr[0])
	is.Equal(float64(4), arr[1])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateArrayReduceEx17(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [1, 2, 3, 4, 5].reduce((sum, x) => sum + x, 0)`)
	must.NoError(err)
	is.Equal(float64(15), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateStringTrimEx17(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => '  hello world  '.trim()`)
	must.NoError(err)
	is.Equal("hello world", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateMapToObjectEx18(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => { const m = new Map([['a', 1], ['b', 2]]); return Object.fromEntries(m); }`)
	must.NoError(err)
	m, ok := result.(map[string]interface{})
	is.True(ok)
	is.Equal(float64(1), m["a"])
	is.Equal(float64(2), m["b"])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSetSizeEx18(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => new Set([1, 2, 3, 2, 1]).size`)
	must.NoError(err)
	is.Equal(float64(3), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateDateToISOStringEx18(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => new Date('2025-01-15').toISOString().slice(0,10)`)
	must.NoError(err)
	is.Equal("2025-01-15", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateSymbolDescriptionEx18(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => Symbol('mySymbol').description`)
	must.NoError(err)
	is.Equal("mySymbol", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateWeakMapEx18(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => { const key = {}; const wm = new WeakMap(); wm.set(key, 42); return wm.get(key); }`)
	must.NoError(err)
	is.Equal(float64(42), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateProxyObjectEx19(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => {
		const target = { x: 10 };
		const proxy = new Proxy(target, {
			get(t, k) { return t[k] * 2; }
		});
		return proxy.x;
	}`)
	must.NoError(err)
	is.Equal(float64(20), result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateGeneratorEx19(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => {
		function* gen() { yield 1; yield 2; yield 3; }
		return [...gen()];
	}`)
	must.NoError(err)
	arr, ok := result.([]interface{})
	is.True(ok)
	is.Len(arr, 3)
	is.Equal(float64(2), arr[1])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateStringRepeatEx19(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => 'ab'.repeat(3)`)
	must.NoError(err)
	is.Equal("ababab", result)
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateArrayFlatEx19(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => [[1,2],[3,4],[5]].flat()`)
	must.NoError(err)
	arr, ok := result.([]interface{})
	is.True(ok)
	is.Len(arr, 5)
	is.Equal(float64(3), arr[2])
}

// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateObjectAssignEx19(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => Object.assign({a:1}, {b:2}, {c:3})`)
	must.NoError(err)
	m, ok := result.(map[string]interface{})
	is.True(ok)
	is.Equal(float64(1), m["a"])
	is.Equal(float64(3), m["c"])
}

// TestWaitForFunctionHandleReturnsHandle verifies WaitForFunction returns a non-nil JSHandle.
// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestWaitForFunctionHandleReturnsHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => window.location.href", nil)
	must.NoError(err)
	must.NotNil(handle)
}

// TestWaitForFunctionHandleCanEvaluateOnResult verifies handle result can be used for GetProperty.
// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestWaitForFunctionHandleCanEvaluateOnResult(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({x: 10, y: 20})", nil)
	must.NoError(err)
	must.NotNil(handle)

	xHandle, err := handle.GetProperty(ctx, "x")
	must.NoError(err)

	xVal, err := xHandle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(10), xVal)
}

// TestWaitForFunctionHandleDispose verifies Dispose can be called on handle.
// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestWaitForFunctionHandleDispose(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({})", nil)
	must.NoError(err)
	must.NotNil(handle)

	must.NoError(handle.Dispose(ctx))
}

// TestWaitForFunctionHandleReturnsDOMElement verifies WaitForFunction can wrap a DOM element.
// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestWaitForFunctionHandleReturnsDOMElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">text</div>`))

	handle, err := page.WaitForFunction(ctx, "() => document.getElementById('el')", nil)
	must.NoError(err)
	must.NotNil(handle)

	// Can evaluate on the element handle
	result, err := handle.Evaluate(ctx, `el => el.id`)
	must.NoError(err)
	is.Equal("el", result)
}

// TestWaitForFunctionHandleGetPropertyNested verifies GetProperty chain works.
// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestWaitForFunctionHandleGetPropertyNested(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.WaitForFunction(ctx, "() => ({a: {b: 42}})", nil)
	must.NoError(err)

	aHandle, err := handle.GetProperty(ctx, "a")
	must.NoError(err)
	must.NotNil(aHandle)

	bHandle, err := aHandle.GetProperty(ctx, "b")
	must.NoError(err)

	bVal, err := bHandle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(float64(42), bVal)
}

// TestPageEvaluateShouldTransferNaN verifies NaN roundtrip via Evaluate.
// Ref: TestPageEvaluate.java#shouldTransferNaN  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferNaN(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a", math.NaN())
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsNaN(f), "expected NaN, got %v (%T)", result, result)
}

// TestPageEvaluateShouldTransferInfinity verifies +Infinity roundtrip.
// Ref: TestPageEvaluate.java#shouldTransferInfinity  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a", math.Inf(1))
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, 1), "expected +Inf, got %v (%T)", result, result)
}

// TestPageEvaluateShouldTransferNegativeInfinity verifies -Infinity roundtrip.
// Ref: TestPageEvaluate.java#shouldTransferNegativeInfinity  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferNegativeInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a", math.Inf(-1))
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, -1), "expected -Inf, got %v (%T)", result, result)
}

// TestPageEvaluateShouldTransfer0 verifies -0 roundtrip.
// Ref: TestPageEvaluate.java#shouldTransfer0  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransfer0(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	negZero := math.Copysign(0, -1)
	result, err := page.Evaluate(ctx, "a => a", negZero)
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.Signbit(f), "expected -0, got %v (%T)", result, result)
}

// TestPageEvaluateShouldReturnNaN verifies NaN return value from JS expression.
// Ref: TestPageEvaluate.java#shouldReturnNaN  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturnNaN(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => NaN")
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsNaN(f), "expected NaN, got %v (%T)", result, result)
}

// TestPageEvaluateShouldReturnInfinity verifies Infinity return value.
// Ref: TestPageEvaluate.java#shouldReturnInfinity  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturnInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => Infinity")
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, 1), "expected +Inf, got %v (%T)", result, result)
}

// TestPageEvaluateShouldReturnNegativeInfinity verifies -Infinity return.
// Ref: TestPageEvaluate.java#shouldReturnNegativeInfinity  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturnNegativeInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => -Infinity")
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, -1), "expected -Inf, got %v (%T)", result, result)
}

// TestPageEvaluateShouldRoundtripUnserializableValues verifies a map with NaN/Infinity/-Infinity/-0 is roundtripped.
// Ref: TestPageEvaluate.java#shouldRoundtripUnserializableValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRoundtripUnserializableValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	input := map[string]any{
		"infinity":  math.Inf(1),
		"nInfinity": math.Inf(-1),
		"nan":       math.NaN(),
		"nZero":     math.Copysign(0, -1),
	}
	result, err := page.Evaluate(ctx, "v => v", input)
	must.NoError(err, "Evaluate failed")
	m, ok := result.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", result, result)
	if v, _ := m["infinity"].(float64); !math.IsInf(v, 1) {
		t.Errorf("infinity: expected +Inf, got %v", m["infinity"])
	}
	if v, _ := m["nInfinity"].(float64); !math.IsInf(v, -1) {
		t.Errorf("nInfinity: expected -Inf, got %v", m["nInfinity"])
	}
	if v, _ := m["nan"].(float64); !math.IsNaN(v) {
		t.Errorf("nan: expected NaN, got %v", m["nan"])
	}
	if v, _ := m["nZero"].(float64); !math.Signbit(v) {
		t.Errorf("nZero: expected -0, got %v", m["nZero"])
	}
}

// TestPageEvaluateShouldAwaitPromise verifies that an async JS expression is awaited.
// Ref: TestPageEvaluate.java#shouldAwaitPromise  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAwaitPromise(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => Promise.resolve(8 * 7)")
	must.NoError(err, "Evaluate failed")
	v, ok := result.(float64)
	must.Truef(ok && v == 56, "expected 56, got %v (%T)", result, result)
}

// TestPageEvaluateShouldRejectPromiseWithException verifies that a rejected JS promise propagates as an error.
// Ref: TestPageEvaluate.java#shouldRejectPromiseWithException  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRejectPromiseWithException(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => not_existing_object.property")
	is.ErrorContains(err, "not_existing_object")
}

// TestPageEvaluateShouldSupportThrownStringsAsErrorMessages verifies a thrown string propagates as an error.
// Ref: TestPageEvaluate.java#shouldSupportThrownStringsAsErrorMessages  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldSupportThrownStringsAsErrorMessages(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => { throw 'qwerty'; }")
	is.ErrorContains(err, "qwerty")
}

// TestPageEvaluateShouldSupportThrownNumbersAsErrorMessages verifies a thrown number propagates as an error.
// Ref: TestPageEvaluate.java#shouldSupportThrownNumbersAsErrorMessages  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldSupportThrownNumbersAsErrorMessages(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => { throw 100500; }")
	is.ErrorContains(err, "100500")
}

// TestPageEvaluateShouldModifyGlobalEnvironment verifies that Evaluate can set and read global variables.
// Ref: TestPageEvaluate.java#shouldModifyGlobalEnvironment  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldModifyGlobalEnvironment(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => window['globalVar'] = 123")
	must.NoError(err, "set global")
	result, err := page.Evaluate(ctx, "globalVar")
	must.NoError(err, "read global")
	v, ok := result.(float64)
	must.Truef(ok && v == 123, "expected 123, got %v (%T)", result, result)
}

// TestPageEvaluateShouldWorkWithUnicodeChars verifies that Unicode keys in map arguments work correctly.
// Ref: TestPageEvaluate.java#shouldWorkWithUnicodeChars  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkWithUnicodeChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a['中文字符']", map[string]any{"中文字符": 42})
	must.NoError(err, "Evaluate failed")
	v, ok := result.(float64)
	must.Truef(ok && v == 42, "expected 42, got %v (%T)", result, result)
}

// TestPageEvaluateShouldTransferArrays verifies that arrays are passed correctly.
// Ref: TestPageEvaluate.java#shouldTransferArrays  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferArrays(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a", []any{1, 2, 3})
	must.NoError(err, "Evaluate failed")
	arr, ok := result.([]any)
	must.Truef(ok && len(arr) == 3, "expected slice of 3, got %v (%T)", result, result)

	// Verify it's treated as an array in JS too.
	isArray, err := page.Evaluate(ctx, "a => Array.isArray(a)", []any{1, 2, 3})
	must.NoError(err, "isArray check failed")
	b, ok := isArray.(bool)
	must.Truef(ok && b, "expected true for Array.isArray, got %v", isArray)
}

// TestPageEvaluateShouldProperlySerializeNullFields verifies that null and undefined values serialize correctly.
// Ref: TestPageEvaluate.java#shouldProperlySerializeNullFields  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldProperlySerializeNullFields(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	// null field
	result, err := page.Evaluate(ctx, "() => ({ a: null })")
	must.NoError(err, "Evaluate null field failed")
	m, ok := result.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", result, result)
	if v, exists := m["a"]; !exists || v != nil {
		t.Errorf("expected {a: nil}, got %v", m)
	}

	// null argument
	res2, err := page.Evaluate(ctx, "x => x", nil)
	must.NoError(err, "Evaluate null arg failed")
	is.Nil(res2)
}

// TestPageEvaluateShouldReturnComplexObjects verifies roundtrip of a complex map value.
// Ref: TestPageEvaluate.java#shouldReturnComplexObjects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturnComplexObjects(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	input := map[string]any{"foo": "bar!"}
	result, err := page.Evaluate(ctx, "a => a", input)
	must.NoError(err, "Evaluate failed")
	m, ok := result.(map[string]any)
	must.Truef(ok && m["foo"] == "bar!", "expected {foo: bar!}, got %v (%T)", result, result)
}

// TestPageEvaluateShouldAcceptAString verifies string expressions without function wrappers.
// Ref: TestPageEvaluate.java#shouldAcceptAString  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAcceptAString(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	r1, err := page.Evaluate(ctx, "1 + 2")
	must.NoError(err, "1+2")
	if v, ok := r1.(float64); !ok || v != 3 {
		t.Errorf("expected 3, got %v", r1)
	}

	r2, err := page.Evaluate(ctx, "1 + 5;")
	must.NoError(err, "1+5 with semicolon")
	if v, ok := r2.(float64); !ok || v != 6 {
		t.Errorf("expected 6, got %v", r2)
	}

	r3, err := page.Evaluate(ctx, "2 + 5;\n// do some math!")
	must.NoError(err, "2+5 with comment")
	if v, ok := r3.(float64); !ok || v != 7 {
		t.Errorf("expected 7, got %v", r3)
	}
}

// TestPageEvaluateShouldReturnUndefinedForNonSerializableObjects verifies non-serializable values return nil.
// Ref: TestPageEvaluate.java#shouldReturnUndefinedForNonSerializableObjects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturnUndefinedForNonSerializableObjects(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => () => {}")
	must.NoError(err, "Evaluate failed")
	is.Nil(result)
}

// TestPageEvaluateShouldBeAbleToThrowATrickyError verifies that an error with a custom message propagates.
// Ref: TestPageEvaluate.java#shouldBeAbleToThrowATrickyError  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldBeAbleToThrowATrickyError(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "errorText => { throw new Error(errorText); }", "My error")
	is.ErrorContains(err, "My error")
}

// TestPageEvaluateShouldAcceptElementHandleAsAnArgument verifies that an ElementHandle can be passed as an argument.
// Ref: TestPageEvaluate.java#shouldAcceptElementHandleAsAnArgument  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAcceptElementHandleAsAnArgument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "goto failed")
	err = page.SetContent(ctx, "<section>42</section>")
	must.NoError(err, "setContent failed")
	elem, err := page.QuerySelector(ctx, "section")
	must.NoError(err, "querySelector failed")
	must.NotNil(elem, "expected non-nil element")
	result, err := page.Evaluate(ctx, "e => e.textContent", elem)
	must.NoError(err, "Evaluate with element handle failed")
	is.Equal("42", result)
}

// TestPageEvaluateShouldWork verifies a simple JS expression evaluates correctly.
// Ref: TestPageEvaluate.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => 7 * 3")
	must.NoError(err)
	v, ok := result.(float64)
	must.Truef(ok && v == 21, "expected 21, got %v (%T)", result, result)
}

// TestPageEvaluateShouldRoundtripPromiseToValue verifies that Promise.resolve values roundtrip correctly.
// Ref: TestPageEvaluate.java#shouldRoundtripPromiseToValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRoundtripPromiseToValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	r1, err := page.Evaluate(ctx, "value => Promise.resolve(value)", nil)
	must.NoError(err)
	is.Nil(r1)

	r2, err := page.Evaluate(ctx, "value => Promise.resolve(value)", math.Inf(1))
	must.NoError(err)
	f, ok := r2.(float64)
	must.Truef(ok && math.IsInf(f, 1), "expected +Inf, got %v (%T)", r2, r2)

	r3, err := page.Evaluate(ctx, "value => Promise.resolve(value)", math.Copysign(0, -1))
	must.NoError(err)
	f, ok = r3.(float64)
	must.Truef(ok && math.Signbit(f), "expected -0, got %v (%T)", r3, r3)
}

// TestPageEvaluateShouldRoundtripPromiseToUnserializableValues verifies Promise.resolve with unserializable map values.
// Ref: TestPageEvaluate.java#shouldRoundtripPromiseToUnserializableValues  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRoundtripPromiseToUnserializableValues(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	input := map[string]any{
		"infinity":  math.Inf(1),
		"nInfinity": math.Inf(-1),
		"nan":       math.NaN(),
		"nZero":     math.Copysign(0, -1),
	}
	result, err := page.Evaluate(ctx, "value => Promise.resolve(value)", input)
	must.NoError(err)
	m, ok := result.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", result, result)
	if v, _ := m["infinity"].(float64); !math.IsInf(v, 1) {
		t.Errorf("infinity: expected +Inf, got %v", m["infinity"])
	}
	if v, _ := m["nInfinity"].(float64); !math.IsInf(v, -1) {
		t.Errorf("nInfinity: expected -Inf, got %v", m["nInfinity"])
	}
	if v, _ := m["nan"].(float64); !math.IsNaN(v) {
		t.Errorf("nan: expected NaN, got %v", m["nan"])
	}
	if v, _ := m["nZero"].(float64); !math.Signbit(v) {
		t.Errorf("nZero: expected -0, got %v", m["nZero"])
	}
}

// TestPageEvaluateShouldReturnUndefinedForObjectsWithSymbols verifies that Symbol-keyed properties and Symbol values serialize as null.
// Ref: TestPageEvaluate.java#shouldReturnUndefinedForObjectsWithSymbols  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturnUndefinedForObjectsWithSymbols(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	r1, err := page.Evaluate(ctx, "() => [Symbol('foo4')]")
	must.NoError(err)
	arr, ok := r1.([]any)
	must.Truef(ok && len(arr) == 1 && arr[0] == nil, "expected [nil], got %v (%T)", r1, r1)

	r2, err := page.Evaluate(ctx, `() => {
		const a = {};
		a[Symbol('foo4')] = 42;
		return a;
	}`)
	must.NoError(err)
	m2, ok := r2.(map[string]any)
	must.Truef(ok && len(m2) == 0, "expected empty map, got %v (%T)", r2, r2)

	r3, err := page.Evaluate(ctx, `() => {
		return { foo: [{ a: Symbol('foo4') }] };
	}`)
	must.NoError(err)
	m3, ok := r3.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", r3, r3)
	foo, _ := m3["foo"].([]any)
	must.NotNil(foo)
	inner, _ := foo[0].(map[string]any)
	must.NotNil(inner)
	if v, exists := inner["a"]; !exists || v != nil {
		t.Errorf("expected {foo: [{a: nil}]}, got %v", r3)
	}
}

// TestPageEvaluateShouldThrowWhenEvaluationTriggersReload verifies that triggering a reload inside evaluate returns a navigation error.
// Ref: TestPageEvaluate.java#shouldThrowWhenEvaluationTriggersReload  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldThrowWhenEvaluationTriggersReload(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		location.reload();
		return new Promise(() => { });
	}`)
	is.ErrorContains(err, "navigation")
}

// TestPageEvaluateShouldEvaluateInThePageContext verifies that a page-loaded global var is accessible via Evaluate.
// Ref: TestPageEvaluate.java#shouldEvaluateInThePageContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldEvaluateInThePageContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/global-var.html", "text/html; charset=utf-8",
		`<!DOCTYPE html><html><head><script>window.globalVar = 123;</script></head><body></body></html>`)
	err := page.Goto(ctx, srv.Prefix()+"/global-var.html")
	must.NoError(err)
	result, err := page.Evaluate(ctx, "globalVar")
	must.NoError(err)
	v, ok := result.(float64)
	must.Truef(ok && v == 123, "expected 123, got %v (%T)", result, result)
}

// TestPageEvaluateShouldWorkWithOverwrittenPromise verifies that Evaluate works even when window.Promise is overwritten.
// Ref: TestPageEvaluate.java#shouldWorkWithOverwrittenPromise  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkWithOverwrittenPromise(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		const originalPromise = window.Promise;
		class Promise2 {
			static all(arg) { return wrap(originalPromise.all(arg)); }
			static race(arg) { return wrap(originalPromise.race(arg)); }
			static resolve(arg) { return wrap(originalPromise.resolve(arg)); }
			constructor(f) { this._promise = new originalPromise(f); }
			then(f, r) { return wrap(this._promise.then(f, r)); }
			catch(f) { return wrap(this._promise.catch(f)); }
			finally(f) { return wrap(this._promise.finally(f)); }
		}
		const wrap = p => {
			const result = new Promise2(() => { });
			result._promise = p;
			return result;
		};
		window.Promise = Promise2;
		window['__Promise2'] = Promise2;
	}`)
	must.NoError(err)

	r1, err := page.Evaluate(ctx, `() => {
		const p = Promise.all([Promise.race([]), new Promise(() => { }).then(() => { })]);
		return p instanceof window['__Promise2'];
	}`)
	must.NoError(err)
	is.Equal(true, r1)

	r2, err := page.Evaluate(ctx, "() => Promise.resolve(42)")
	must.NoError(err)
	v, ok := r2.(float64)
	must.Truef(ok && v == 42, "expected 42, got %v (%T)", r2, r2)
}

// TestPageEvaluateShouldProperlySerializeUndefinedFields verifies that undefined object fields serialize as null.
// Ref: TestPageEvaluate.java#shouldProperlySerializeUndefinedFields  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldProperlySerializeUndefinedFields(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => ({ a: undefined })")
	must.NoError(err)
	m, ok := result.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", result, result)
	if v, exists := m["a"]; !exists || v != nil {
		t.Errorf("expected {a: nil}, got %v", m)
	}
}

// TestPageEvaluateShouldReturn0 verifies that `() => -0` returns negative zero.
// Ref: TestPageEvaluate.java#shouldReturn0  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturn0(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => -0")
	must.NoError(err)
	f, ok := result.(float64)
	must.Truef(ok && math.Signbit(f), "expected -0, got %v (%T)", result, result)
}

// TestPageEvaluateNonSerializableWindow verifies that evaluating `window` returns "ref: <Window>".
// Ref: TestPageEvaluate.java#shouldReturnUndefinedForNonSerializableObjects (window case)
func TestPageEvaluateNonSerializableWindow(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => window")
	must.NoError(err)
	is.Equal("ref: <Window>", result)
}

// TestPageEvaluateShouldThrowIfUnderlyingElementWasDisposed verifies that using a disposed ElementHandle errors.
// Ref: TestPageEvaluate.java#shouldThrowIfUnderlyingElementWasDisposed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldThrowIfUnderlyingElementWasDisposed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<section>39</section>"))
	elem, err := page.QuerySelector(ctx, "section")
	must.NoError(err)
	must.NotNil(elem)
	err = elem.Dispose(ctx)
	must.NoError(err)

	_, err = page.Evaluate(ctx, "e => e.textContent", elem)
	is.ErrorContains(err, "no object with guid")
}

// TestPageEvaluateShouldSimulateAUserGesture verifies that document.execCommand succeeds in evaluate context.
// Ref: TestPageEvaluate.java#shouldSimulateAUserGesture  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldSimulateAUserGesture(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => {
		document.body.appendChild(document.createTextNode('test'));
		document.execCommand('selectAll');
		return document.execCommand('copy');
	}`)
	must.NoError(err)
	is.Equal(true, result)
}

// TestPageEvaluateShouldNotThrowAnErrorWhenEvaluationDoesANavigation verifies that changing location.href during evaluate does not error.
// Ref: TestPageEvaluate.java#shouldNotThrowAnErrorWhenEvaluationDoesANavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldNotThrowAnErrorWhenEvaluationDoesANavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/one-style.html", "text/html; charset=utf-8",
		`<!DOCTYPE html><html><head><style>div { color: red; }</style></head><body></body></html>`)
	err := page.Goto(ctx, srv.Prefix()+"/one-style.html")
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => {
		window.location.href = '/empty.html';
		return [42];
	}`)
	must.NoError(err)
	arr, ok := result.([]any)
	must.Truef(ok && len(arr) == 1, "expected [42], got %v (%T)", result, result)
	if v, ok := arr[0].(float64); !ok || v != 42 {
		t.Errorf("expected [42], got %v", arr)
	}
}

// TestPageEvaluateShouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsAnObject verifies that synchronous reload with object return works.
// Ref: TestPageEvaluate.java#shouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsAnObject  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsAnObject(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	if globalBTName == "webkit" {
		t.Skip("fixme: disabled for webkit")
	}
	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => {
		window.location.reload();
		return { a: 42 };
	}`)
	must.NoError(err)
	m, ok := result.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", result, result)
	if v, _ := m["a"].(float64); v != 42 {
		t.Errorf("expected {a: 42}, got %v", m)
	}
}

// TestPageEvaluateShouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsUndefined verifies that synchronous reload with undefined return works.
// Ref: TestPageEvaluate.java#shouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsUndefined  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsUndefined(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => {
		window.location.reload();
		return undefined;
	}`)
	must.NoError(err)
	is.Nil(result)
}

// TestPageEvaluateShouldThrowErrorWithDetailedInformationOnExceptionInsidePromise verifies that an error thrown inside a Promise propagates.
// Ref: TestPageEvaluate.java#shouldThrowErrorWithDetailedInformationOnExceptionInsidePromise  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldThrowErrorWithDetailedInformationOnExceptionInsidePromise(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => new Promise(() => {
		throw new Error('Error in promise');
	})`)
	is.ErrorContains(err, "Error in promise")
}

// TestPageEvaluateShouldWorkEvenWhenJSONIsSetToNull verifies that Evaluate works even when window.JSON is set to null.
// Ref: TestPageEvaluate.java#shouldWorkEvenWhenJSONIsSetToNull  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkEvenWhenJSONIsSetToNull(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => { window.JSON.stringify = null; window.JSON = null; }")
	must.NoError(err)
	result, err := page.Evaluate(ctx, "() => ({ abc: 123 })")
	must.NoError(err)
	m, ok := result.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", result, result)
	if v, _ := m["abc"].(float64); v != 123 {
		t.Errorf("expected {abc: 123}, got %v", m)
	}
}

// TestPageEvaluateShouldAwaitPromiseFromPopup verifies that a Promise from window.open can be awaited.
// Ref: TestPageEvaluate.java#shouldAwaitPromiseFromPopup  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAwaitPromiseFromPopup(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	err := page.Goto(ctx, srv.EmptyPage())
	must.NoError(err)
	result, err := page.Evaluate(ctx, `() => {
		const win = window.open('about:blank');
		return new win['Promise'](f => f(42));
	}`)
	must.NoError(err)
	v, ok := result.(float64)
	must.Truef(ok && v == 42, "expected 42, got %v (%T)", result, result)
}

// TestPageEvaluateShouldWorkWithNonStrictExpressions verifies that undeclared variables work in non-strict mode.
// Ref: TestPageEvaluate.java#shouldWorkWithNonStrictExpressions  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkWithNonStrictExpressions(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => {
		y = 3.14;
		return y;
	}`)
	must.NoError(err)
	v, ok := result.(float64)
	must.Truef(ok && v == 3.14, "expected 3.14, got %v (%T)", result, result)
}

// TestPageEvaluateShouldRespectUseStrictExpression verifies that undeclared variables throw in use-strict mode.
// Ref: TestPageEvaluate.java#shouldRespectUseStrictExpression  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRespectUseStrictExpression(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => {
		'use strict';
		// @ts-ignore
		variableY = 3.14;
		// @ts-ignore
		return variableY;
	}`)
	is.ErrorContains(err, "variableY")
}

// TestPageEvaluateShouldNotLeakUtilityScript verifies that `this` in global context is window.
// Ref: TestPageEvaluate.java#shouldNotLeakUtilityScript  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldNotLeakUtilityScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => this === window")
	must.NoError(err)
	is.Equal(true, result)
}

// TestPageEvaluateShouldNotLeakHandles verifies that `handles` is not accessible as a global.
// Ref: TestPageEvaluate.java#shouldNotLeakHandles  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldNotLeakHandles(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "handles.length")
	is.ErrorContains(err, "handles")
}

// TestPageEvaluateShouldAliasWindowDocumentAndNode verifies that window, document, and body serialize as "ref: ..." strings.
// Ref: TestPageEvaluate.java#shouldAliasWindowDocumentAndNode  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAliasWindowDocumentAndNode(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "[window, document, document.body]")
	must.NoError(err)
	arr, ok := result.([]any)
	must.Truef(ok && len(arr) == 3, "expected slice of 3, got %v (%T)", result, result)
	is.Equal("ref: <Window>", arr[0])
	is.Equal("ref: <Document>", arr[1])
	is.Equal("ref: <Node>", arr[2])
}

// TestPageEvaluateShouldNotUseToJSONWhenEvaluating verifies that a toJSON method on an object is not called during evaluation.
// Ref: TestPageEvaluate.java#shouldNotUseToJSONWhenEvaluating  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldNotUseToJSONWhenEvaluating(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => ({ toJSON: () => 'string', data: 'data' })")
	must.NoError(err)
	m, ok := result.(map[string]any)
	must.Truef(ok, "expected map, got %T: %v", result, result)
	is.Equal("data", m["data"])
	if m["toJSON"] != nil {
		_, isMap := m["toJSON"].(map[string]any)
		must.Truef(isMap, "toJSON should be nil or empty map, got %T: %v", m["toJSON"], m["toJSON"])
	}
}

// Ref: TestPageEvaluate.java#shouldEvaluateDate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldEvaluateDate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => new Date('2020-01-01').getFullYear()`)
	must.NoError(err)
	is.Equal(float64(2020), result)
}

// Observaciones 1.2: empty arrays/objects as arguments must be accepted by the driver.
// Extra: no equivalent method in playwright-java; related class: TestPageEvaluate.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestEvaluateEmptyArrayAndObjectArguments(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	res, err := page.Evaluate(ctx, "a => a.length", []any{})
	must.NoError(err)
	is.Equal(float64(0), res)

	res, err = page.Evaluate(ctx, "o => Object.keys(o).length", map[string]any{})
	must.NoError(err)
	is.Equal(float64(0), res)

	res, err = page.Evaluate(ctx, "a => a[0].length + a[1]", []any{[]any{}, "a"})
	must.NoError(err)
	is.Equal("0a", res)
}

// Ref: TestPageEvaluate.java#shouldTransferPositive0  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferPositive0(t *testing.T) {
	ctx, _, page := a2Setup(t)
	assert.Equal(t, true, a2Ev(t, ctx, page, "a => Object.is(a, 0)", 0.0))
	v := a2Ev(t, ctx, page, "a => a", 0.0)
	f, ok := v.(float64)
	require.True(t, ok, "got %T", v)
	assert.True(t, math.IsInf(1/f, 1))
}

// Ref: TestPageEvaluate.java#shouldTransferArraysAsArraysNotObjects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferArraysAsArraysNotObjects(t *testing.T) {
	ctx, _, page := a2Setup(t)
	assert.Equal(t, true, a2Ev(t, ctx, page, "a => Array.isArray(a)", []any{1, 2, 3}))
}

// Ref: TestPageEvaluate.java#shouldTransferBigint  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferBigint(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	expected, ok := new(big.Int).SetString("123456789012345678901234567890", 10)
	must.True(ok)
	res1, err := page.Evaluate(ctx, "() => 123456789012345678901234567890n")
	must.NoError(err)
	is.Equal(expected, res1)

	res2, err := page.Evaluate(ctx, "a => a", expected)
	must.NoError(err)
	is.Equal(expected, res2)
}

// Ref: TestPageEvaluate.java#shouldTransferMapsAsEmptyObjects  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferMapsAsEmptyObjects(t *testing.T) {
	t.Skip("not implemented: Java body is empty (Not applicable)")
}

// Ref: TestPageEvaluate.java#shouldWorkRightAfterFramenavigated  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkRightAfterFramenavigated(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, srv, page := a2Setup(t)

	evalCh := make(chan any, 1)
	page.OnFrameNavigated(func(f *playwright.Frame) {
		val, err := f.Evaluate(ctx, "() => 6 * 7")
		if err == nil {
			select {
			case evalCh <- val:
			default:
			}
		}
	})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	select {
	case val := <-evalCh:
		is.EqualValues(42, val)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for frameNavigated evaluation")
	}
}

// Ref: TestPageEvaluate.java#shouldWorkRightAfterACrossOriginNavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkRightAfterACrossOriginNavigation(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, srv, page := a2Setup(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	evalCh := make(chan any, 1)
	page.OnFrameNavigated(func(f *playwright.Frame) {
		val, err := f.Evaluate(ctx, "() => 6 * 7")
		if err == nil {
			select {
			case evalCh <- val:
			default:
			}
		}
	})
	must.NoError(page.Goto(ctx, srv.CrossProcessPrefix()+"/empty.html"))
	select {
	case val := <-evalCh:
		is.EqualValues(42, val)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for frameNavigated evaluation")
	}
}

// Ref: TestPageEvaluate.java#shouldWorkFromInsideAnExposedFunction  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkFromInsideAnExposedFunction(t *testing.T) {
	ctx, _, page := a2Setup(t)
	err := page.ExposeBinding(ctx, "callController", func(args ...any) any {
		res, err := page.Evaluate(ctx, "({ a, b }) => a * b", map[string]any{"a": args[0], "b": args[1]})
		if err != nil {
			return err.Error()
		}
		return res
	})
	require.NoError(t, err)
	res := a2Ev(t, ctx, page, "async function() {\n  return await window['callController'](9, 3);\n}")
	assert.EqualValues(t, 27, res)
}

// Ref: TestPageEvaluate.java#shouldProperlySerializeUndefinedArguments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldProperlySerializeUndefinedArguments(t *testing.T) {
	t.Skip("not implemented: Java body is empty (Not applicable)")
}

// Ref: TestPageEvaluate.java#shouldReturnUndefinedProperties  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldReturnUndefinedProperties(t *testing.T) {
	ctx, _, page := a2Setup(t)
	assert.Equal(t, map[string]any{"a": nil}, a2Ev(t, ctx, page, "() => ({ a: undefined })"))
}

// Ref: TestPageEvaluate.java#shouldProperlySerializeNullArguments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldProperlySerializeNullArguments(t *testing.T) {
	ctx, _, page := a2Setup(t)
	assert.Nil(t, a2Ev(t, ctx, page, "x => x", nil))
}

// Ref: TestPageEvaluate.java#shouldAcceptAStringWithSemiColons  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAcceptAStringWithSemiColons(t *testing.T) {
	ctx, _, page := a2Setup(t)
	assert.EqualValues(t, 6, a2Ev(t, ctx, page, "1 + 5;"))
}

// Ref: TestPageEvaluate.java#shouldAcceptAStringWithComments  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAcceptAStringWithComments(t *testing.T) {
	ctx, _, page := a2Setup(t)
	assert.EqualValues(t, 7, a2Ev(t, ctx, page, "2 + 5;\n// do some math!"))
}

// Ref: TestPageEvaluate.java#shouldThrowANiceErrorAfterANavigation  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldThrowANiceErrorAfterANavigation(t *testing.T) {
	t.Skip("not implemented: Page.WaitForNavigation")
}

// Ref: TestPageEvaluate.java#shouldTransfer100MbOfDataFromPageToNodeJs  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransfer100MbOfDataFromPageToNodeJs(t *testing.T) {
	ctx, _, page := a2Setup(t)
	a := a2Ev(t, ctx, page, "() => Array(100 * 1024 * 1024 + 1).join('a')")
	str, ok := a.(string)
	require.True(t, ok, "got %T", a)
	assert.Equal(t, 100*1024*1024, len(str))
	assert.Equal(t, 0, len(strings.Trim(str, "a")))
}

// Ref: TestPageEvaluate.java#shouldWorkWithNewFunctionAndCSP  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkWithNewFunctionAndCSP(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "script-src "+srv.Prefix())
		w.Header().Set("Content-Type", "text/html")
	})
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/empty.html"))
	assert.Equal(t, true, a2Ev(t, ctx, page, "() => new Function('return true')()"))
}

// Ref: TestPageEvaluate.java#shouldWorkWithCSP  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkWithCSP(t *testing.T) {
	ctx, srv, page := a2Setup(t)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "script-src 'self'")
		w.Header().Set("Content-Type", "text/html")
	})
	require.NoError(t, page.Goto(ctx, srv.EmptyPage()))
	assert.EqualValues(t, 4, a2Ev(t, ctx, page, "() => 2 + 2"))
}

// Ref: TestPageEvaluate.java#shouldEvaluateException  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldEvaluateException(t *testing.T) {
	ctx, _, page := a2Setup(t)
	result := a2Ev(t, ctx, page, "() => {\n  return (function functionOnStack() {\n    return new Error('error message');\n  })();\n}")
	msg := fmt.Sprint(result)
	assert.Contains(t, msg, "Error: error message")
	assert.Contains(t, msg, "functionOnStack")
}

// Ref: TestPageEvaluate.java#shouldEvaluateException2  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldEvaluateException2(t *testing.T) {
	ctx, _, page := a2Setup(t)
	result := a2Ev(t, ctx, page, "new Error('error message')")
	assert.Contains(t, fmt.Sprint(result), "Error: error message")
}

// Ref: TestPageEvaluate.java#shouldPassExceptionArgument  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldPassExceptionArgument(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	e := errors.New("error message")
	raw, err := page.Evaluate(ctx, "e => ({ message: e.message, name: e.name, stack: e.stack })", e)
	must.NoError(err)
	received, ok := raw.(map[string]any)
	must.True(ok)
	is.Equal("error message", received["message"])
	is.Equal("Exception", received["name"])
	is.NotNil(received["stack"])
	is.Contains(fmt.Sprint(received["stack"]), "TestPageEvaluateShouldPassExceptionArgument")
}

// Ref: TestPageEvaluate.java#shouldRoundtripDate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRoundtripDate(t *testing.T) {
	ctx, _, page := a2Setup(t)
	date := time.Date(2020, 5, 27, 1, 31, 38, 506000000, time.UTC)
	result := a2Ev(t, ctx, page, "date => date", date)
	assert.Equal(t, "2020-05-27T01:31:38.506Z", fmt.Sprint(result))
}

// Ref: TestPageEvaluate.java#shouldRoundtripRegex  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRoundtripRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	regex := regexp.MustCompile("(?im)hello")
	result, err := page.Evaluate(ctx, "regex => regex", regex)
	must.NoError(err)
	re, ok := result.(*regexp.Regexp)
	must.True(ok, "expected *regexp.Regexp, got %T", result)
	is.Equal(regex.String(), re.String())
}

// Ref: TestPageEvaluate.java#shouldJsonValueDate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldJsonValueDate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.EvaluateHandle(ctx, "() => ({ date: new Date('2020-05-27T01:31:38.506Z') })")
	must.NoError(err)
	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	m, ok := val.(map[string]any)
	must.True(ok)
	expected := time.Date(2020, 5, 27, 1, 31, 38, 506000000, time.UTC)
	is.Equal(expected, m["date"])
}

// Ref: TestPageEvaluate.java#shouldEvaluateUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldEvaluateUrl(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	res, err := page.Evaluate(ctx, "() => ({ url: new URL('https://example.com/') })")
	must.NoError(err)
	m, ok := res.(map[string]any)
	must.True(ok)
	expected, err := url.Parse("https://example.com/")
	must.NoError(err)
	is.Equal(expected, m["url"])
}

// Ref: TestPageEvaluate.java#shouldRoundtripUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRoundtripUrl(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	u, err := url.Parse("https://example.com/")
	must.NoError(err)
	res, err := page.Evaluate(ctx, "url => url", u)
	must.NoError(err)
	got, ok := res.(*url.URL)
	must.True(ok, "expected *url.URL, got %T", res)
	is.Equal(u.String(), got.String())
}

// Ref: TestPageEvaluate.java#shouldRoundtripComplexUrl  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldRoundtripComplexUrl(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	u, err := url.Parse("https://user:password@www.contoso.com:80/Home/Index.htm?q1=v1&q2=v2#FragmentName")
	must.NoError(err)
	res, err := page.Evaluate(ctx, "url => url", u)
	must.NoError(err)
	got, ok := res.(*url.URL)
	must.True(ok, "expected *url.URL, got %T", res)
	is.Equal(u.String(), got.String())
}

// Ref: TestPageEvaluate.java#shouldNotUseToJSONInJsonValue  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldNotUseToJSONInJsonValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	handle, err := page.EvaluateHandle(ctx, "() => ({ toJSON: () => 'string', data: 'data' })")
	must.NoError(err)
	val, err := handle.JSONValue(ctx)
	must.NoError(err)
	is.Equal(map[string]any{"data": "data", "toJSON": map[string]any{}}, val)
}

// Ref: TestPageEvaluate.java#shouldWorkForCircularObject  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldWorkForCircularObject(t *testing.T) {
	ctx, _, page := a2Setup(t)
	result := a2Ev(t, ctx, page, "() => {\n  const a = {};\n  a.b = a;\n  return a;\n}")
	m, ok := result.(map[string]any)
	require.True(t, ok, "got %T", result)
	assert.Len(t, m, 1)
	assert.NotNil(t, m["b"])
}

// Ref: TestPageEvaluate.java#shouldAcceptParameter  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldAcceptParameter(t *testing.T) {
	ctx, _, page := a2Setup(t)
	now := time.Now().UTC().Truncate(time.Millisecond)
	result := a2Ev(t, ctx, page, "p => p", now)
	assert.Equal(t, now.Format("2006-01-02T15:04:05.000Z"), fmt.Sprint(result))
}

// Ref: TestPageEvaluate.java#shouldTransferMaps  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferMaps(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	res, err := page.Evaluate(ctx, "() => new Map([[1, { test: 42n }]])")
	must.NoError(err)
	is.Equal(map[string]any{}, res)
}

// Ref: TestPageEvaluate.java#shouldTransferSets  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEvaluate.java
func TestPageEvaluateShouldTransferSets(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)

	res, err := page.Evaluate(ctx, "() => new Set([1, { test: 42n }])")
	must.NoError(err)
	is.Equal(map[string]any{}, res)
}
