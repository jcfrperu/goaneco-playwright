package playwright

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// serializedValueRaw represents the wire format of SerializedValue in the Playwright protocol.
// Ref: packages/protocol/spec/structs.yml and Serialization.java
type serializedValueRaw struct {
	N *float64             `json:"n,omitempty"`
	B *bool                `json:"b,omitempty"`
	S *string              `json:"s,omitempty"`
	V *string              `json:"v,omitempty"`
	D *string              `json:"d,omitempty"`
	A []serializedValueRaw `json:"a,omitempty"`
	O []nameSerializedRaw  `json:"o,omitempty"`
	H *int                 `json:"h,omitempty"` // index into SerializedArgument.handles for object handle references
}

type nameSerializedRaw struct {
	K string             `json:"k"`
	V serializedValueRaw `json:"v"`
}

type serializedArgumentRaw struct {
	Value   serializedValueRaw `json:"value"`
	Handles []any              `json:"handles"`
}

// serializeArgument serializes a Go value into the SerializedArgument structure required by Playwright IPC.
// ElementHandle and JSHandle are passed as handle references (index into the handles array).
func serializeArgument(arg any) serializedArgumentRaw {
	if arg == nil {
		v := "undefined"
		return serializedArgumentRaw{
			Value:   serializedValueRaw{V: &v},
			Handles: []any{},
		}
	}
	switch v := arg.(type) {
	case *ElementHandle:
		idx := 0
		return serializedArgumentRaw{
			Value:   serializedValueRaw{H: &idx},
			Handles: []any{map[string]string{"guid": v.owner.GUID()}},
		}
	case *JSHandle:
		idx := 0
		return serializedArgumentRaw{
			Value:   serializedValueRaw{H: &idx},
			Handles: []any{map[string]string{"guid": v.owner.GUID()}},
		}
	default:
		return serializedArgumentRaw{
			Value:   serializeValue(arg),
			Handles: []any{},
		}
	}
}

func serializeValue(v any) serializedValueRaw {
	if v == nil {
		str := "null"
		return serializedValueRaw{V: &str}
	}

	switch val := v.(type) {
	case string:
		return serializedValueRaw{S: &val}
	case bool:
		return serializedValueRaw{B: &val}
	case int:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case int8:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case int16:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case int32:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case int64:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case uint:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case uint8:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case uint16:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case uint32:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case uint64:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case float32:
		f := float64(val)
		return serializedValueRaw{N: &f}
	case float64:
		if math.IsNaN(val) {
			v := "NaN"
			return serializedValueRaw{V: &v}
		}
		if math.IsInf(val, 1) {
			v := "Infinity"
			return serializedValueRaw{V: &v}
		}
		if math.IsInf(val, -1) {
			v := "-Infinity"
			return serializedValueRaw{V: &v}
		}
		if val == 0 && math.Signbit(val) {
			v := "-0"
			return serializedValueRaw{V: &v}
		}
		return serializedValueRaw{N: &val}
	case []any:
		arr := make([]serializedValueRaw, len(val))
		for i, item := range val {
			arr[i] = serializeValue(item)
		}
		return serializedValueRaw{A: arr}
	case []string:
		arr := make([]serializedValueRaw, len(val))
		for i, item := range val {
			str := item
			arr[i] = serializedValueRaw{S: &str}
		}
		return serializedValueRaw{A: arr}
	case map[string]any:
		obj := make([]nameSerializedRaw, 0, len(val))
		for k, item := range val {
			obj = append(obj, nameSerializedRaw{
				K: k,
				V: serializeValue(item),
			})
		}
		return serializedValueRaw{O: obj}
	default:
		// Use reflection to handle typed maps and structs.
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Map:
			if rv.Type().Key().Kind() == reflect.String {
				obj := make([]nameSerializedRaw, 0, rv.Len())
				for _, key := range rv.MapKeys() {
					obj = append(obj, nameSerializedRaw{
						K: key.String(),
						V: serializeValue(rv.MapIndex(key).Interface()),
					})
				}
				return serializedValueRaw{O: obj}
			}
		case reflect.Slice, reflect.Array:
			arr := make([]serializedValueRaw, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				arr[i] = serializeValue(rv.Index(i).Interface())
			}
			return serializedValueRaw{A: arr}
		case reflect.Struct:
			t := rv.Type()
			obj := make([]nameSerializedRaw, 0, t.NumField())
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if !f.IsExported() {
					continue
				}
				name := f.Tag.Get("json")
				if name == "-" {
					continue // explicitly excluded field
				}
				if name == "" {
					name = f.Name
				} else if comma := strings.Index(name, ","); comma >= 0 {
					name = name[:comma]
				}
				obj = append(obj, nameSerializedRaw{
					K: name,
					V: serializeValue(rv.Field(i).Interface()),
				})
			}
			return serializedValueRaw{O: obj}
		case reflect.Ptr:
			if rv.IsNil() {
				str := "null"
				return serializedValueRaw{V: &str}
			}
			return serializeValue(rv.Elem().Interface())
		}
		// Final fallback: marshal to JSON string.
		data, err := json.Marshal(v)
		if err != nil {
			str := fmt.Sprint(v)
			return serializedValueRaw{S: &str}
		}
		str := string(data)
		return serializedValueRaw{S: &str}
	}
}

// deserializeValue reconstructs a native Go value from the protocol SerializedValue structure.
func deserializeValue(raw serializedValueRaw) (any, error) {
	if raw.S != nil {
		return *raw.S, nil
	}
	if raw.N != nil {
		return *raw.N, nil
	}
	if raw.B != nil {
		return *raw.B, nil
	}
	if raw.V != nil {
		switch *raw.V {
		case "undefined", "null":
			return nil, nil
		case "NaN":
			return math.NaN(), nil
		case "Infinity":
			return math.Inf(1), nil
		case "-Infinity":
			return math.Inf(-1), nil
		case "-0":
			return math.Copysign(0, -1), nil
		default:
			return *raw.V, nil
		}
	}
	if raw.D != nil {
		return *raw.D, nil
	}
	if raw.A != nil {
		list := make([]any, len(raw.A))
		for i, item := range raw.A {
			deserialized, err := deserializeValue(item)
			if err != nil {
				return nil, err
			}
			list[i] = deserialized
		}
		return list, nil
	}
	if raw.O != nil {
		m := make(map[string]any, len(raw.O))
		for _, entry := range raw.O {
			deserialized, err := deserializeValue(entry.V)
			if err != nil {
				return nil, err
			}
			m[entry.K] = deserialized
		}
		return m, nil
	}
	return nil, nil
}

// headersToNameValues converts a map of HTTP headers to protocol NameValue slice.
func headersToNameValues(headers map[string]string) []protocol.NameValue {
	if headers == nil {
		return nil
	}
	nv := make([]protocol.NameValue, 0, len(headers))
	for k, v := range headers {
		nv = append(nv, protocol.NameValue{
			Name:  k,
			Value: v,
		})
	}
	return nv
}

// nameValuesToHeaders converts a slice of protocol NameValue entries to a map of headers.
// Multiple occurrences of the same header (e.g. Set-Cookie) are joined by newline characters,
// matching the Playwright protocol specification.
func nameValuesToHeaders(nv []protocol.NameValue) map[string]string {
	if nv == nil {
		return nil
	}
	res := make(map[string]string, len(nv))
	for _, item := range nv {
		lowerName := strings.ToLower(item.Name)
		if existing, ok := res[lowerName]; ok {
			res[lowerName] = existing + "\n" + item.Value
		} else {
			res[lowerName] = item.Value
		}
	}
	return res
}

// getHeader searches a NameValue slice for a header by name (case-insensitive)
// and returns its concatenated value or empty string if not found.
func getHeader(nv []protocol.NameValue, name string) string {
	target := strings.ToLower(name)
	var matches []string
	for _, item := range nv {
		if strings.ToLower(item.Name) == target {
			matches = append(matches, item.Value)
		}
	}
	if len(matches) == 0 {
		return ""
	}
	return strings.Join(matches, "\n")
}
