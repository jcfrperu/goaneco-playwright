//go:build e2e

// E2E tests for JSHandle: EvaluateHandle, GetProperty, Dispose.
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
