package playwright

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerializeArgument_EmptyContainersKeepArrayAndObject(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	for name, arg := range map[string]any{
		"empty []any":       []any{},
		"empty []string":    []string{},
		"empty map":         map[string]any{},
		"nested empty list": []any{[]any{}, "a"},
	} {
		data, err := json.Marshal(serializeArgument(arg))
		must.NoError(err, name)
		is.NotContains(string(data), `"value":{}`, name)
	}

	data, err := json.Marshal(serializeArgument([]any{}))
	must.NoError(err)
	is.JSONEq(`{"value":{"a":[]},"handles":[]}`, string(data))

	data, err = json.Marshal(serializeArgument(map[string]any{}))
	must.NoError(err)
	is.JSONEq(`{"value":{"o":[]},"handles":[]}`, string(data))
}

func TestSerializeArgument_NestedHandlesGetIncrementalIndexes(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	eh := &ElementHandle{owner: ChannelOwner{guid: "handle@1"}}
	jh := &JSHandle{owner: ChannelOwner{guid: "handle@2"}}

	arg := serializeArgument([]any{eh, map[string]any{"x": jh}, "s"})
	must.Len(arg.Handles, 2)

	data, err := json.Marshal(arg)
	must.NoError(err)
	is.JSONEq(`{"value":{"a":[{"h":0},{"o":[{"k":"x","v":{"h":1}}]},{"s":"s"}]},
		"handles":[{"guid":"handle@1"},{"guid":"handle@2"}]}`, string(data))
}

func TestSerializeArgument_TopLevelHandle(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	arg := serializeArgument(&JSHandle{owner: ChannelOwner{guid: "handle@9"}})
	data, err := json.Marshal(arg)
	must.NoError(err)
	assert.JSONEq(t, `{"value":{"h":0},"handles":[{"guid":"handle@9"}]}`, string(data))
}
