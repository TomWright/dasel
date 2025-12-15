package toml

import (
	"fmt"
	"reflect"
	"strconv"

	pkg "github.com/pelletier/go-toml/v2"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

var _ parsing.Writer = (*tomlWriter)(nil)

func newTOMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	// Default to double quotes unless overridden.
	w := &tomlWriter{}
	return w, nil
}

type tomlWriter struct {
}

// Write converts the dasel model.Value into Go values backed by dynamically
// generated struct types (reflect.StructOf) that preserve key ordering, then
// delegates to go-toml Marshal to produce canonical TOML output.
func (j *tomlWriter) Write(value *model.Value) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	var gv interface{}
	var err error

	if value.IsMap() {
		gv, err = buildGoValueForMap(value)
		if err != nil {
			return nil, fmt.Errorf("failed to construct go value: %w", err)
		}
	} else {
		// handle scalars, slices, etc. Use goTypeAndValue to get a concrete reflect.Value
		typ, rv, err := goTypeAndValue(value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert non-map top-level value: %w", err)
		}
		// For nil/zero interface, ensure we pass nil interface rather than a typed zero.
		if typ.Kind() == reflect.Interface && rv.IsZero() {
			gv = nil
		} else {
			gv = rv.Interface()
		}
	}

	outBytes, err := pkg.Marshal(gv)
	if err != nil {
		return nil, fmt.Errorf("toml marshal failed: %w", err)
	}

	// Ensure trailing newline for consistency with other format writers/tests.
	if len(outBytes) == 0 || outBytes[len(outBytes)-1] != '\n' {
		outBytes = append(outBytes, '\n')
	}

	return outBytes, nil
}

// buildGoValueForMap constructs a reflect.Value that is a struct with fields in
// the same order as keys in the provided map Value. The struct fields are
// tagged with `toml:"<original key>"` so go-toml will use the original key
// names when encoding.
func buildGoValueForMap(v *model.Value) (interface{}, error) {
	kvs, err := v.MapKeyValues()
	if err != nil {
		return nil, err
	}

	// Build struct fields in order
	fields := make([]reflect.StructField, 0, len(kvs))
	fieldValues := make([]reflect.Value, 0, len(kvs))

	for i, kv := range kvs {
		ft, fv, err := goTypeAndValue(kv.Value)
		if err != nil {
			return nil, fmt.Errorf("error converting key %q: %w", kv.Key, err)
		}

		// create exported field name (F0, F1...) and set tag to preserve toml key
		field := reflect.StructField{
			Name: "F" + strconv.Itoa(i),
			Type: ft,
			Tag:  reflect.StructTag(`toml:"` + kv.Key + `"`),
		}
		fields = append(fields, field)
		fieldValues = append(fieldValues, fv)
	}

	structType := reflect.StructOf(fields)
	val := reflect.New(structType).Elem()

	for i := range fieldValues {
		val.Field(i).Set(fieldValues[i])
	}

	return val.Interface(), nil
}

// goTypeAndValue returns a reflect.Type and a reflect.Value suitable for use as a
// struct field/type for the given model.Value. It converts nested maps into
// reflect.StructOf types recursively (preserving key order), converts slices to
// slices of appropriate element types when possible, or []interface{} otherwise.
func goTypeAndValue(v *model.Value) (reflect.Type, reflect.Value, error) {
	switch v.Type() {
	case model.TypeString:
		s, err := v.StringValue()
		if err != nil {
			return nil, reflect.Value{}, err
		}
		return reflect.TypeOf(""), reflect.ValueOf(s), nil
	case model.TypeInt:
		i, err := v.IntValue()
		if err != nil {
			return nil, reflect.Value{}, err
		}
		// use int64 to be safe
		return reflect.TypeOf(int64(0)), reflect.ValueOf(i), nil
	case model.TypeFloat:
		f, err := v.FloatValue()
		if err != nil {
			return nil, reflect.Value{}, err
		}
		return reflect.TypeOf(float64(0)), reflect.ValueOf(f), nil
	case model.TypeBool:
		b, err := v.BoolValue()
		if err != nil {
			return nil, reflect.Value{}, err
		}
		return reflect.TypeOf(true), reflect.ValueOf(b), nil
	case model.TypeNull:
		// represent null as interface{}(nil)
		return reflect.TypeOf((*interface{})(nil)).Elem(), reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem()), nil
	case model.TypeMap:
		// For nested maps, prefer map[string]interface{} values rather than
		// generating nested struct types. Using struct types for nested maps
		// causes go-toml to emit explicit table headers for those nested
		// structures which can change the document shape on round-trip. The
		// top-level map is still handled by buildGoValueForMap which generates
		// a struct to preserve ordering.
		kvs, err := v.MapKeyValues()
		if err != nil {
			return nil, reflect.Value{}, err
		}
		m := make(map[string]interface{})
		for _, kv := range kvs {
			_, rv, err := goTypeAndValue(kv.Value)
			if err != nil {
				return nil, reflect.Value{}, fmt.Errorf("error in nested key %q: %w", kv.Key, err)
			}
			m[kv.Key] = rv.Interface()
		}
		return reflect.TypeOf(map[string]interface{}{}), reflect.ValueOf(m), nil
	case model.TypeSlice:
		// Decide element types. If all elements are maps and have compatible keys,
		// build a struct type for elements and return a slice of that struct type.
		length, _ := v.SliceLen()
		if length == 0 {
			// empty slice -> use []interface{}
			elemType := reflect.TypeOf((*interface{})(nil)).Elem()
			return reflect.SliceOf(elemType), reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0), nil
		}

		// inspect first element
		first, _ := v.GetSliceIndex(0)
		if first.IsMap() {
			// Collect union of keys across elements in order of appearance
			seen := map[string]bool{}
			keys := make([]string, 0)
			_ = v.RangeSlice(func(_ int, item *model.Value) error {
				if !item.IsMap() {
					// mixed types -> fallback to []interface{}
					keys = nil
					return nil
				}
				kvs, _ := item.MapKeyValues()
				for _, kv := range kvs {
					if !seen[kv.Key] {
						seen[kv.Key] = true
						keys = append(keys, kv.Key)
					}
				}
				return nil
			})

			if keys != nil {
				// Build []map[string]interface{} to preserve keys and values without
				// forcing a generated struct type which can cause conversion issues.
				sliceMaps := make([]map[string]interface{}, 0, length)
				_ = v.RangeSlice(func(_ int, item *model.Value) error {
					m := map[string]interface{}{}
					if item.IsMap() {
						kvs, _ := item.MapKeyValues()
						for _, kv := range kvs {
							_, rv, err := goTypeAndValue(kv.Value)
							if err != nil {
								return err
							}
							m[kv.Key] = rv.Interface()
						}
					}
					sliceMaps = append(sliceMaps, m)
					return nil
				})
				return reflect.TypeOf([]map[string]interface{}{}), reflect.ValueOf(sliceMaps), nil
			}
		}

		// fallback: build []interface{}
		elems := make([]interface{}, 0, length)
		_ = v.RangeSlice(func(_ int, item *model.Value) error {
			gt, rv, err := goTypeAndValue(item)
			if err != nil {
				return err
			}
			// convert reflect.Value to interface{}
			elems = append(elems, rv.Interface())
			_ = gt
			return nil
		})
		return reflect.TypeOf([]interface{}{}), reflect.ValueOf(elems), nil
	default:
		// fallback to stringified interface
		s := fmt.Sprintf("%v", v.Interface())
		return reflect.TypeOf(""), reflect.ValueOf(s), nil
	}
}
