package playwright

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

type serializedRegExpRaw struct {
	P string `json:"p"`
	F string `json:"f,omitempty"`
}

type serializedErrorRaw struct {
	M string `json:"m"`
	N string `json:"n"`
	S string `json:"s,omitempty"`
}

// serializedValueRaw represents the wire format of SerializedValue in the Playwright protocol.
// Ref: packages/protocol/spec/structs.yml and Serialization.java
type serializedValueRaw struct {
	N   *float64             `json:"n,omitempty"`
	B   *bool                `json:"b,omitempty"`
	S   *string              `json:"s,omitempty"`
	V   *string              `json:"v,omitempty"`
	D   *string              `json:"d,omitempty"`
	U   *string              `json:"u,omitempty"`
	BI  *string              `json:"bi,omitempty"`
	R   *serializedRegExpRaw `json:"r,omitempty"`
	E   *serializedErrorRaw  `json:"e,omitempty"`
	A   []serializedValueRaw `json:"a,omitempty"`
	O   []nameSerializedRaw  `json:"o,omitempty"`
	H   *int                 `json:"h,omitempty"` // index into SerializedArgument.handles for object handle references
	ID  *int                 `json:"id,omitempty"`
	Ref *int                 `json:"ref,omitempty"`
}

// MarshalJSON always emits "a" and "o" when the slices are non-nil, even if empty.
// With a plain omitempty tag, empty arrays/objects would be dropped and the driver
// would fail with "Attempting to deserialize unexpected value".
func (r serializedValueRaw) MarshalJSON() ([]byte, error) {
	type wire struct {
		N   *float64              `json:"n,omitempty"`
		B   *bool                 `json:"b,omitempty"`
		S   *string               `json:"s,omitempty"`
		V   *string               `json:"v,omitempty"`
		D   *string               `json:"d,omitempty"`
		U   *string               `json:"u,omitempty"`
		BI  *string               `json:"bi,omitempty"`
		R   *serializedRegExpRaw  `json:"r,omitempty"`
		E   *serializedErrorRaw   `json:"e,omitempty"`
		A   *[]serializedValueRaw `json:"a,omitempty"`
		O   *[]nameSerializedRaw  `json:"o,omitempty"`
		H   *int                  `json:"h,omitempty"`
		ID  *int                  `json:"id,omitempty"`
		Ref *int                  `json:"ref,omitempty"`
	}
	w := wire{N: r.N, B: r.B, S: r.S, V: r.V, D: r.D, U: r.U, BI: r.BI, R: r.R, E: r.E, H: r.H, ID: r.ID, Ref: r.Ref}
	if r.A != nil {
		w.A = &r.A
	}
	if r.O != nil {
		w.O = &r.O
	}
	return json.Marshal(w)
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
	handles := []any{}
	value := serializeValueH(arg, &handles)
	return serializedArgumentRaw{Value: value, Handles: handles}
}

// serializeValue serializes a value that cannot contain handles.
func serializeValue(v any) serializedValueRaw {
	handles := []any{}
	return serializeValueH(v, &handles)
}

// serializeValueH serializes v, appending every ElementHandle/JSHandle found at any depth
// to handles and referencing it by its index ({"h": n}).
func serializeValueH(v any, handles *[]any) serializedValueRaw {
	if v == nil {
		str := "null"
		return serializedValueRaw{V: &str}
	}

	switch val := v.(type) {
	case *ElementHandle:
		if val == nil {
			str := "null"
			return serializedValueRaw{V: &str}
		}
		idx := len(*handles)
		*handles = append(*handles, map[string]string{"guid": val.owner.GUID()})
		return serializedValueRaw{H: &idx}
	case *JSHandle:
		if val == nil {
			str := "null"
			return serializedValueRaw{V: &str}
		}
		idx := len(*handles)
		*handles = append(*handles, map[string]string{"guid": val.owner.GUID()})
		return serializedValueRaw{H: &idx}
	case string:
		return serializedValueRaw{S: &val}
	case bool:
		return serializedValueRaw{B: &val}
	case *url.URL:
		if val == nil {
			str := "null"
			return serializedValueRaw{V: &str}
		}
		uStr := val.String()
		return serializedValueRaw{U: &uStr}
	case url.URL:
		uStr := val.String()
		return serializedValueRaw{U: &uStr}
	case *time.Time:
		if val == nil {
			str := "null"
			return serializedValueRaw{V: &str}
		}
		dStr := val.UTC().Format("2006-01-02T15:04:05.000Z")
		return serializedValueRaw{D: &dStr}
	case time.Time:
		dStr := val.UTC().Format("2006-01-02T15:04:05.000Z")
		return serializedValueRaw{D: &dStr}
	case *big.Int:
		if val == nil {
			str := "null"
			return serializedValueRaw{V: &str}
		}
		biStr := val.String()
		return serializedValueRaw{BI: &biStr}
	case big.Int:
		biStr := val.String()
		return serializedValueRaw{BI: &biStr}
	case *regexp.Regexp:
		if val == nil {
			str := "null"
			return serializedValueRaw{V: &str}
		}
		p, f := extractJsRegex(val.String())
		return serializedValueRaw{R: &serializedRegExpRaw{P: p, F: f}}
	case regexp.Regexp:
		p, f := extractJsRegex(val.String())
		return serializedValueRaw{R: &serializedRegExpRaw{P: p, F: f}}
	case error:
		msg := val.Error()
		stack := string(debug.Stack())
		return serializedValueRaw{E: &serializedErrorRaw{M: msg, N: "Exception", S: stack}}
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
			arr[i] = serializeValueH(item, handles)
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
				V: serializeValueH(item, handles),
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
						V: serializeValueH(rv.MapIndex(key).Interface(), handles),
					})
				}
				return serializedValueRaw{O: obj}
			}
		case reflect.Slice, reflect.Array:
			arr := make([]serializedValueRaw, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				arr[i] = serializeValueH(rv.Index(i).Interface(), handles)
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
					V: serializeValueH(rv.Field(i).Interface(), handles),
				})
			}
			return serializedValueRaw{O: obj}
		case reflect.Ptr:
			if rv.IsNil() {
				str := "null"
				return serializedValueRaw{V: &str}
			}
			return serializeValueH(rv.Elem().Interface(), handles)
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
	refs := make(map[int]any)
	return deserializeValueInner(raw, refs)
}

func deserializeValueInner(raw serializedValueRaw, refs map[int]any) (any, error) {
	if raw.Ref != nil {
		if val, ok := refs[*raw.Ref]; ok {
			return val, nil
		}
	}
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
	if raw.U != nil {
		parsed, err := url.Parse(*raw.U)
		if err != nil {
			return *raw.U, nil
		}
		return parsed, nil
	}
	if raw.D != nil {
		t, err := time.Parse(time.RFC3339Nano, *raw.D)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05.000Z", *raw.D)
		}
		if err != nil {
			return *raw.D, nil
		}
		return t, nil
	}
	if raw.BI != nil {
		bi, ok := new(big.Int).SetString(*raw.BI, 10)
		if !ok {
			return *raw.BI, nil
		}
		return bi, nil
	}
	if raw.R != nil {
		prefix := ""
		if raw.R.F != "" {
			prefix = "(?" + raw.R.F + ")"
		}
		re, err := regexp.Compile(prefix + raw.R.P)
		if err != nil {
			re, _ = regexp.Compile(raw.R.P)
		}
		if re != nil {
			return re, nil
		}
		return raw.R.P, nil
	}
	if raw.E != nil {
		return errors.New(raw.E.M), nil
	}
	if raw.A != nil {
		list := make([]any, len(raw.A))
		if raw.ID != nil {
			refs[*raw.ID] = list
		}
		for i, item := range raw.A {
			deserialized, err := deserializeValueInner(item, refs)
			if err != nil {
				return nil, err
			}
			list[i] = deserialized
		}
		return list, nil
	}
	if raw.O != nil {
		m := make(map[string]any, len(raw.O))
		if raw.ID != nil {
			refs[*raw.ID] = m
		}
		for _, entry := range raw.O {
			deserialized, err := deserializeValueInner(entry.V, refs)
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

func extractJsRegex(expr string) (string, string) {
	if strings.HasPrefix(expr, "(?") {
		if idx := strings.Index(expr, ")"); idx > 2 {
			flags := expr[2:idx]
			valid := true
			var jsFlags strings.Builder
			for _, r := range flags {
				if r == 'i' || r == 'm' || r == 's' {
					jsFlags.WriteRune(r)
				} else {
					valid = false
					break
				}
			}
			if valid {
				return expr[idx+1:], jsFlags.String()
			}
		}
	}
	return expr, ""
}
