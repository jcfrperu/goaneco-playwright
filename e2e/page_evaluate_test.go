//go:build e2e

// Extended Page.Evaluate E2E tests.
// Migration of: TestPageEvaluate.java
package e2e

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
)

// TestPageEvaluateTransferNaN verifies NaN roundtrip via Evaluate.
// Ref: TestPageEvaluate.java#shouldTransferNaN
func TestPageEvaluateTransferNaN(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a", math.NaN())
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsNaN(f), "expected NaN, got %v (%T)", result, result)
}

// TestPageEvaluateTransferInfinity verifies +Infinity roundtrip.
// Ref: TestPageEvaluate.java#shouldTransferInfinity
func TestPageEvaluateTransferInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a", math.Inf(1))
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, 1), "expected +Inf, got %v (%T)", result, result)
}

// TestPageEvaluateTransferNegativeInfinity verifies -Infinity roundtrip.
// Ref: TestPageEvaluate.java#shouldTransferNegativeInfinity
func TestPageEvaluateTransferNegativeInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a", math.Inf(-1))
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, -1), "expected -Inf, got %v (%T)", result, result)
}

// TestPageEvaluateTransferNegativeZero verifies -0 roundtrip.
// Ref: TestPageEvaluate.java#shouldTransfer0
func TestPageEvaluateTransferNegativeZero(t *testing.T) {
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

// TestPageEvaluateReturnNaN verifies NaN return value from JS expression.
// Ref: TestPageEvaluate.java#shouldReturnNaN
func TestPageEvaluateReturnNaN(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => NaN")
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsNaN(f), "expected NaN, got %v (%T)", result, result)
}

// TestPageEvaluateReturnInfinity verifies Infinity return value.
// Ref: TestPageEvaluate.java#shouldReturnInfinity
func TestPageEvaluateReturnInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => Infinity")
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, 1), "expected +Inf, got %v (%T)", result, result)
}

// TestPageEvaluateReturnNegativeInfinity verifies -Infinity return.
// Ref: TestPageEvaluate.java#shouldReturnNegativeInfinity
func TestPageEvaluateReturnNegativeInfinity(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => -Infinity")
	must.NoError(err, "Evaluate failed")
	f, ok := result.(float64)
	must.Truef(ok && math.IsInf(f, -1), "expected -Inf, got %v (%T)", result, result)
}

// TestPageEvaluateRoundtripUnserializableValues verifies a map with NaN/Infinity/-Infinity/-0 is roundtripped.
// Ref: TestPageEvaluate.java#shouldRoundtripUnserializableValues
func TestPageEvaluateRoundtripUnserializableValues(t *testing.T) {
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

// TestPageEvaluateAwaitPromise verifies that an async JS expression is awaited.
// Ref: TestPageEvaluate.java#shouldAwaitPromise
func TestPageEvaluateAwaitPromise(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => Promise.resolve(8 * 7)")
	must.NoError(err, "Evaluate failed")
	v, ok := result.(float64)
	must.Truef(ok && v == 56, "expected 56, got %v (%T)", result, result)
}

// TestPageEvaluateRejectPromise verifies that a rejected JS promise propagates as an error.
// Ref: TestPageEvaluate.java#shouldRejectPromiseWithException
func TestPageEvaluateRejectPromise(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => not_existing_object.property")
	is.ErrorContains(err, "not_existing_object")
}

// TestPageEvaluateThrowsString verifies a thrown string propagates as an error.
// Ref: TestPageEvaluate.java#shouldSupportThrownStringsAsErrorMessages
func TestPageEvaluateThrowsString(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => { throw 'qwerty'; }")
	is.ErrorContains(err, "qwerty")
}

// TestPageEvaluateThrowsNumber verifies a thrown number propagates as an error.
// Ref: TestPageEvaluate.java#shouldSupportThrownNumbersAsErrorMessages
func TestPageEvaluateThrowsNumber(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "() => { throw 100500; }")
	is.ErrorContains(err, "100500")
}

// TestPageEvaluateModifyGlobal verifies that Evaluate can set and read global variables.
// Ref: TestPageEvaluate.java#shouldModifyGlobalEnvironment
func TestPageEvaluateModifyGlobal(t *testing.T) {
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

// TestPageEvaluateUnicodeChars verifies that Unicode keys in map arguments work correctly.
// Ref: TestPageEvaluate.java#shouldWorkWithUnicodeChars
func TestPageEvaluateUnicodeChars(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "a => a['中文字符']", map[string]any{"中文字符": 42})
	must.NoError(err, "Evaluate failed")
	v, ok := result.(float64)
	must.Truef(ok && v == 42, "expected 42, got %v (%T)", result, result)
}

// TestPageEvaluateTransferArrays verifies that arrays are passed correctly.
// Ref: TestPageEvaluate.java#shouldTransferArrays + shouldTransferArraysAsArraysNotObjects
func TestPageEvaluateTransferArrays(t *testing.T) {
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

// TestPageEvaluateNullUndefinedFields verifies that null and undefined values serialize correctly.
// Ref: TestPageEvaluate.java#shouldProperlySerializeNullFields + shouldReturnUndefinedForObjectsWithSymbols
func TestPageEvaluateNullUndefinedFields(t *testing.T) {
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

// TestPageEvaluateReturnComplexObject verifies roundtrip of a complex map value.
// Ref: TestPageEvaluate.java#shouldReturnComplexObjects
func TestPageEvaluateReturnComplexObject(t *testing.T) {
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

// TestPageEvaluateStringExpressions verifies string expressions without function wrappers.
// Ref: TestPageEvaluate.java#shouldAcceptAString + shouldAcceptAStringWithSemiColons + shouldAcceptAStringWithComments
func TestPageEvaluateStringExpressions(t *testing.T) {
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

// TestPageEvaluateReturnUndefinedForNonSerializable verifies non-serializable values return nil.
// Ref: TestPageEvaluate.java#shouldReturnUndefinedForNonSerializableObjects
func TestPageEvaluateReturnUndefinedForNonSerializable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => () => {}")
	must.NoError(err, "Evaluate failed")
	is.Nil(result)
}

// TestPageEvaluateThrowTrickyError verifies that an error with a custom message propagates.
// Ref: TestPageEvaluate.java#shouldBeAbleToThrowATrickyError
func TestPageEvaluateThrowTrickyError(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "errorText => { throw new Error(errorText); }", "My error")
	is.ErrorContains(err, "My error")
}

// TestPageEvaluateElementHandleAsArg verifies that an ElementHandle can be passed as an argument.
// Ref: TestPageEvaluate.java#shouldAcceptElementHandleAsAnArgument
func TestPageEvaluateElementHandleAsArg(t *testing.T) {
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

// TestPageEvaluateWork verifies a simple JS expression evaluates correctly.
// Ref: TestPageEvaluate.java#shouldWork
func TestPageEvaluateWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => 7 * 3")
	must.NoError(err)
	v, ok := result.(float64)
	must.Truef(ok && v == 21, "expected 21, got %v (%T)", result, result)
}

// TestPageEvaluateRoundtripPromiseToValue verifies that Promise.resolve values roundtrip correctly.
// Ref: TestPageEvaluate.java#shouldRoundtripPromiseToValue
func TestPageEvaluateRoundtripPromiseToValue(t *testing.T) {
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

// TestPageEvaluateRoundtripPromiseToUnserializableValues verifies Promise.resolve with unserializable map values.
// Ref: TestPageEvaluate.java#shouldRoundtripPromiseToUnserializableValues
func TestPageEvaluateRoundtripPromiseToUnserializableValues(t *testing.T) {
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

// TestPageEvaluateSymbolFields verifies that Symbol-keyed properties and Symbol values serialize as null.
// Ref: TestPageEvaluate.java#shouldReturnUndefinedForObjectsWithSymbols
func TestPageEvaluateSymbolFields(t *testing.T) {
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

// TestPageEvaluateThrowsOnReload verifies that triggering a reload inside evaluate returns a navigation error.
// Ref: TestPageEvaluate.java#shouldThrowWhenEvaluationTriggersReload
func TestPageEvaluateThrowsOnReload(t *testing.T) {
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

// TestPageEvaluateInPageContext verifies that a page-loaded global var is accessible via Evaluate.
// Ref: TestPageEvaluate.java#shouldEvaluateInThePageContext
func TestPageEvaluateInPageContext(t *testing.T) {
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

// TestPageEvaluateOverwrittenPromise verifies that Evaluate works even when window.Promise is overwritten.
// Ref: TestPageEvaluate.java#shouldWorkWithOverwrittenPromise
func TestPageEvaluateOverwrittenPromise(t *testing.T) {
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

// TestPageEvaluateSerializeUndefinedFields verifies that undefined object fields serialize as null.
// Ref: TestPageEvaluate.java#shouldProperlySerializeUndefinedFields + shouldReturnUndefinedProperties
func TestPageEvaluateSerializeUndefinedFields(t *testing.T) {
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

// TestPageEvaluateReturnNegativeZero verifies that `() => -0` returns negative zero.
// Ref: TestPageEvaluate.java#shouldReturn0
func TestPageEvaluateReturnNegativeZero(t *testing.T) {
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

// TestPageEvaluateDisposedElementHandle verifies that using a disposed ElementHandle errors.
// Ref: TestPageEvaluate.java#shouldThrowIfUnderlyingElementWasDisposed
func TestPageEvaluateDisposedElementHandle(t *testing.T) {
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

// TestPageEvaluateSimulateUserGesture verifies that document.execCommand succeeds in evaluate context.
// Ref: TestPageEvaluate.java#shouldSimulateAUserGesture
func TestPageEvaluateSimulateUserGesture(t *testing.T) {
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

// TestPageEvaluateNavigationDoesNotThrow verifies that changing location.href during evaluate does not error.
// Ref: TestPageEvaluate.java#shouldNotThrowAnErrorWhenEvaluationDoesANavigation
func TestPageEvaluateNavigationDoesNotThrow(t *testing.T) {
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

// TestPageEvaluateSyncNavigationReturnsObject verifies that synchronous reload with object return works.
// Ref: TestPageEvaluate.java#shouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsAnObject
func TestPageEvaluateSyncNavigationReturnsObject(t *testing.T) {
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

// TestPageEvaluateSyncNavigationReturnsUndefined verifies that synchronous reload with undefined return works.
// Ref: TestPageEvaluate.java#shouldNotThrowAnErrorWhenEvaluationDoesASynchronousNavigationAndReturnsUndefined
func TestPageEvaluateSyncNavigationReturnsUndefined(t *testing.T) {
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

// TestPageEvaluateErrorInPromise verifies that an error thrown inside a Promise propagates.
// Ref: TestPageEvaluate.java#shouldThrowErrorWithDetailedInformationOnExceptionInsidePromise
func TestPageEvaluateErrorInPromise(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => new Promise(() => {
		throw new Error('Error in promise');
	})`)
	is.ErrorContains(err, "Error in promise")
}

// TestPageEvaluateWithJSONNull verifies that Evaluate works even when window.JSON is set to null.
// Ref: TestPageEvaluate.java#shouldWorkEvenWhenJSONIsSetToNull
func TestPageEvaluateWithJSONNull(t *testing.T) {
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

// TestPageEvaluateAwaitPromiseFromPopup verifies that a Promise from window.open can be awaited.
// Ref: TestPageEvaluate.java#shouldAwaitPromiseFromPopup
func TestPageEvaluateAwaitPromiseFromPopup(t *testing.T) {
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

// TestPageEvaluateNonStrictExpression verifies that undeclared variables work in non-strict mode.
// Ref: TestPageEvaluate.java#shouldWorkWithNonStrictExpressions
func TestPageEvaluateNonStrictExpression(t *testing.T) {
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

// TestPageEvaluateStrictModeThrows verifies that undeclared variables throw in use-strict mode.
// Ref: TestPageEvaluate.java#shouldRespectUseStrictExpression
func TestPageEvaluateStrictModeThrows(t *testing.T) {
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

// TestPageEvaluateNoLeakUtilityScript verifies that `this` in global context is window.
// Ref: TestPageEvaluate.java#shouldNotLeakUtilityScript
func TestPageEvaluateNoLeakUtilityScript(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, "() => this === window")
	must.NoError(err)
	is.Equal(true, result)
}

// TestPageEvaluateNoLeakHandles verifies that `handles` is not accessible as a global.
// Ref: TestPageEvaluate.java#shouldNotLeakHandles
func TestPageEvaluateNoLeakHandles(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, "handles.length")
	is.ErrorContains(err, "handles")
}

// TestPageEvaluateAliasWindowDocumentNode verifies that window, document, and body serialize as "ref: ..." strings.
// Ref: TestPageEvaluate.java#shouldAliasWindowDocumentAndNode
func TestPageEvaluateAliasWindowDocumentNode(t *testing.T) {
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

// TestPageEvaluateToJSONIgnored verifies that a toJSON method on an object is not called during evaluation.
// Ref: TestPageEvaluate.java#shouldNotUseToJSONWhenEvaluating
func TestPageEvaluateToJSONIgnored(t *testing.T) {
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

func TestEvaluateDateEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	result, err := page.Evaluate(ctx, `() => new Date('2020-01-01').getFullYear()`)
	must.NoError(err)
	is.Equal(float64(2020), result)
}

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

func TestEvaluateSetTimeoutEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `() => new Promise(resolve => setTimeout(resolve, 10))`)
	must.NoError(err)
}

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

func TestEvaluateDocumentCookieEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<html><body>Page</body></html>`))

	_, err := page.Evaluate(ctx, `() => document.cookie`)
	must.NoError(err)
}

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
// Ref: TestPageEvaluateHandle.java#shouldReturnHandle
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
// Ref: TestPageEvaluateHandle.java#shouldEvaluateOnResult
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
// Ref: TestPageEvaluateHandle.java#shouldDispose
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
// Ref: TestPageEvaluateHandle.java#shouldReturnDOMElement
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
// Ref: TestPageEvaluateHandle.java#shouldGetNestedProperty
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
