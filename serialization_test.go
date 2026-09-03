package playwright

import (
	"encoding/json"
	"testing"

	"github.com/jcfrperu/goaneco-playwright/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_headersToNameValues_PreservesAll (UNIT-HDR-01)
// Verifies that converting a map[string]string to []protocol.NameValue preserves all key-value pairs.
func Test_headersToNameValues_PreservesAll(t *testing.T) {
	is := assert.New(t)

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "text/html,application/xhtml+xml",
		"User-Agent":   "GoPlaywright/1.0",
	}

	nv := headersToNameValues(headers)
	is.Len(nv, 3)

	m := make(map[string]string)
	for _, item := range nv {
		m[item.Name] = item.Value
	}

	is.Equal("application/json", m["Content-Type"])
	is.Equal("text/html,application/xhtml+xml", m["Accept"])
	is.Equal("GoPlaywright/1.0", m["User-Agent"])

	// Nil input
	is.Nil(headersToNameValues(nil))
}

// Test_nameValuesToHeaders_MultiValueAndCaseInsensitive (UNIT-HDR-01 / UNIT-HDR-02)
// Verifies that converting []protocol.NameValue to map[string]string merges repeated headers (e.g. Set-Cookie)
// with newline delimiters and lowercases keys for standard header lookup.
func Test_nameValuesToHeaders_MultiValueAndCaseInsensitive(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	nv := []protocol.NameValue{
		{Name: "Accept", Value: "text/html"},
		{Name: "Set-Cookie", Value: "session=abc"},
		{Name: "Set-Cookie", Value: "theme=dark"},
		{Name: "X-Custom-Header", Value: "custom-val"},
	}

	headers := nameValuesToHeaders(nv)
	must.NotNil(headers)
	is.Equal("text/html", headers["accept"])
	is.Equal("session=abc\ntheme=dark", headers["set-cookie"])
	is.Equal("custom-val", headers["x-custom-header"])

	// Nil input
	is.Nil(nameValuesToHeaders(nil))
}

// Test_getHeader_CaseInsensitive (UNIT-HDR-01 / UNIT-HDR-02)
// Verifies case-insensitive header lookup in []protocol.NameValue.
func Test_getHeader_CaseInsensitive(t *testing.T) {
	is := assert.New(t)

	nv := []protocol.NameValue{
		{Name: "Content-Type", Value: "application/json"},
		{Name: "Set-Cookie", Value: "a=1"},
		{Name: "SET-COOKIE", Value: "b=2"},
	}

	is.Equal("application/json", getHeader(nv, "content-type"))
	is.Equal("application/json", getHeader(nv, "CONTENT-TYPE"))
	is.Equal("application/json", getHeader(nv, "Content-Type"))
	is.Equal("a=1\nb=2", getHeader(nv, "set-cookie"))
	is.Equal("", getHeader(nv, "non-existent"))
}

// TestSerialization_ExplicitZeroTimeout (UNIT-IPC-03 / UNIT-TIME-02)
// Verifies that when a timeout of 0.0 is explicitly passed, the serialized JSON contains
// {"timeout": 0} and is NOT omitted or converted to nil/default.
func TestSerialization_ExplicitZeroTimeout(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	zero := 0.0
	opts := &BrowserTypeLaunchOptions{
		Timeout: &zero,
	}

	data, err := json.Marshal(opts)
	must.NoError(err)

	var unmarshaled map[string]any
	err = json.Unmarshal(data, &unmarshaled)
	must.NoError(err)

	val, exists := unmarshaled["timeout"]
	is.True(exists, "explicit zero timeout must be preserved in JSON")
	is.Equal(float64(0), val)
}

// TestSerialization_NilTimeoutOmitted (UNIT-IPC-03 / UNIT-TIME-02)
// Verifies that when timeout is nil, the field is omitted so that defaults apply.
func TestSerialization_NilTimeoutOmitted(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	opts := &BrowserTypeLaunchOptions{
		Timeout: nil,
	}

	data, err := json.Marshal(opts)
	must.NoError(err)

	var unmarshaled map[string]any
	err = json.Unmarshal(data, &unmarshaled)
	must.NoError(err)

	_, exists := unmarshaled["timeout"]
	is.False(exists, "nil timeout must be omitted from JSON to use default")
}

// TestSerializeArgument_AndDeserialize_Roundtrip tests serializeArgument and deserializeValue across data types.
func TestSerializeArgument_AndDeserialize_Roundtrip(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	testCases := []struct {
		name     string
		input    any
		validate func(t *testing.T, result any)
	}{
		{
			name:  "string",
			input: "hello world",
			validate: func(t *testing.T, res any) {
				is.Equal("hello world", res)
			},
		},
		{
			name:  "integer",
			input: 42,
			validate: func(t *testing.T, res any) {
				is.Equal(float64(42), res)
			},
		},
		{
			name:  "float",
			input: 3.14159,
			validate: func(t *testing.T, res any) {
				is.InDelta(3.14159, res.(float64), 0.0001)
			},
		},
		{
			name:  "boolean true",
			input: true,
			validate: func(t *testing.T, res any) {
				is.Equal(true, res)
			},
		},
		{
			name:  "nil / undefined",
			input: nil,
			validate: func(t *testing.T, res any) {
				is.Nil(res)
			},
		},
		{
			name:  "slice of strings",
			input: []string{"a", "b", "c"},
			validate: func(t *testing.T, res any) {
				slice, ok := res.([]any)
				is.True(ok)
				is.Equal([]any{"a", "b", "c"}, slice)
			},
		},
		{
			name: "map with nested values",
			input: map[string]any{
				"name":    "playwright",
				"version": 1,
				"active":  true,
			},
			validate: func(t *testing.T, res any) {
				m, ok := res.(map[string]any)
				is.True(ok)
				is.Equal("playwright", m["name"])
				is.Equal(float64(1), m["version"])
				is.Equal(true, m["active"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serialized := serializeArgument(tc.input)
			deserialized, err := deserializeValue(serialized.Value)
			must.NoError(err)
			tc.validate(t, deserialized)
		})
	}
}
