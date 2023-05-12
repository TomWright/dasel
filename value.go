package dasel

import (
	"github.com/tomwright/dasel/v2/dencoding"
	"github.com/tomwright/dasel/v2/util"
	"reflect"
)

// Value is a wrapper around reflect.Value that adds some handy helper funcs.
type Value struct {
	reflect.Value
	setFn    func(value Value)
	deleteFn func()
	metadata map[string]interface{}
}

// ValueOf wraps value in a Value.
func ValueOf(value interface{}) Value {
	switch v := value.(type) {
	case Value:
		return v
	case reflect.Value:
		return Value{
			Value: v,
		}
	default:
		return Value{
			Value: reflect.ValueOf(value),
		}
	}
}

// Metadata returns the metadata with a key of key for v.
func (v Value) Metadata(key string) interface{} {
	if v.metadata == nil {
		return nil
	}
	if m, ok := v.metadata[key]; ok {
		return m
	}
	return nil
}

// WithMetadata sets the given value into the values metadata.
func (v Value) WithMetadata(key string, value interface{}) Value {
	if v.metadata == nil {
		v.metadata = map[string]interface{}{}
	}
	v.metadata[key] = value
	return v
}

// Interface returns the interface{} value of v.
func (v Value) Interface() interface{} {
	return v.Unpack().Interface()
}

// Len returns v's length.
func (v Value) Len() int {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Unpack().Len()
	case reflect.Bool:
		if v.Interface() == true {
			return 1
		} else {
			return 0
		}
	default:
		return len(util.ToString(v.Interface()))
	}
}

// String returns the string v's underlying value, as a string.
func (v Value) String() string {
	return v.Unpack().String()
}

// IsEmpty returns true is v represents an empty reflect.Value.
func (v Value) IsEmpty() bool {
	return isEmptyReflectValue(unpackReflectValue(v.Value))
}

func isEmptyReflectValue(v reflect.Value) bool {
	if (v == reflect.Value{}) {
		return true
	}
	return v.Kind() == reflect.String && v.Interface() == UninitialisedPlaceholder
}

// Kind returns the underlying type of v.
func (v Value) Kind() reflect.Kind {
	return v.Unpack().Kind()
}

func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for _, v := range kinds {
		if v == kind {
			return true
		}
	}
	return false
}

var dencodingMapType = reflect.TypeOf(&dencoding.Map{})

func isdencodingMap(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr && value.Type() == dencodingMapType
}

func unpackReflectValue(value reflect.Value, kinds ...reflect.Kind) reflect.Value {
	if len(kinds) == 0 {
		kinds = append(kinds, reflect.Ptr, reflect.Interface)
	}
	res := value
	for {
		if isdencodingMap(res) {
			return res
		}
		if !containsKind(kinds, res.Kind()) {
			return res
		}
		res = res.Elem()
	}
}

func (v Value) FirstAddressable() reflect.Value {
	res := v.Value
	for !res.CanAddr() {
		res = res.Elem()
	}
	return res
}

// Unpack returns the underlying reflect.Value after resolving any pointers or interface types.
func (v Value) Unpack(kinds ...reflect.Kind) reflect.Value {
	return unpackReflectValue(v.Value, kinds...)
}

func (v Value) Type() reflect.Type {
	return v.Unpack().Type()
}

// Set sets underlying value of v.
// Depends on setFn since the implementation can differ depending on how the Value was initialised.
func (v Value) Set(value Value) {
	if v.setFn != nil {
		v.setFn(value)
		return
	}
	panic("unable to set value with missing setFn")
}

// Delete deletes the current element.
// Depends on deleteFn since the implementation can differ depending on how the Value was initialised.
func (v Value) Delete() {
	if v.deleteFn != nil {
		v.deleteFn()
		return
	}
	panic("unable to delete value with missing deleteFn")
}

func (v Value) IsDencodingMap() bool {
	if v.Kind() != reflect.Ptr {
		return false
	}
	_, ok := v.Interface().(*dencoding.Map)
	return ok
}

func (v Value) dencodingMapIndex(key Value) Value {
	getValueByKey := func() reflect.Value {
		if !v.IsDencodingMap() {
			return reflect.Value{}
		}
		om := v.Interface().(*dencoding.Map)
		if v, ok := om.Get(key.Value.String()); !ok {
			return reflect.Value{}
		} else {
			return reflect.ValueOf(v)
		}
	}
	index := Value{
		Value: getValueByKey(),
		setFn: func(value Value) {
			// Note that we do not use Interface() here as it will dereference the received value.
			// Instead, we only dereference the interface type to receive the pointer.
			v.Interface().(*dencoding.Map).Set(key.Value.String(), value.Unpack(reflect.Interface).Interface())
		},
		deleteFn: func() {
			v.Interface().(*dencoding.Map).Delete(key.Value.String())
		},
	}
	return index.
		WithMetadata("key", key.Interface()).
		WithMetadata("parent", v)
}

// MapIndex returns the value associated with key in the map v.
// It returns the zero Value if no field was found.
func (v Value) MapIndex(key Value) Value {
	index := Value{
		Value: v.Unpack().MapIndex(key.Value),
		setFn: func(value Value) {
			v.Unpack().SetMapIndex(key.Value, value.Value)
		},
		deleteFn: func() {
			v.Unpack().SetMapIndex(key.Value, reflect.Value{})
		},
	}
	return index.
		WithMetadata("key", key.Interface()).
		WithMetadata("parent", v)
}

func (v Value) MapKeys() []Value {
	res := make([]Value, 0)
	for _, k := range v.Unpack().MapKeys() {
		res = append(res, Value{Value: k})
	}
	return res
}

// FieldByName returns the struct field with the given name.
// It returns the zero Value if no field was found.
func (v Value) FieldByName(name string) Value {
	return Value{
		Value: v.Unpack().FieldByName(name),
		setFn: func(value Value) {
			v.Unpack().FieldByName(name).Set(value.Value)
		},
		deleteFn: func() {
			field := v.Unpack().FieldByName(name)
			field.Set(reflect.New(field.Type()))
		},
	}.
		WithMetadata("key", name).
		WithMetadata("parent", v)
}

// NumField returns the number of fields in the struct v.
func (v Value) NumField() int {
	return v.Unpack().NumField()
}

// Index returns v's i'th element.
// It panics if v's Kind is not Array, Slice, or String or i is out of range.
func (v Value) Index(i int) Value {
	return Value{
		Value: v.Unpack().Index(i),
		setFn: func(value Value) {
			v.Unpack().Index(i).Set(value.Value)
		},
		deleteFn: func() {
			currentLen := v.Len()
			updatedSlice := reflect.MakeSlice(sliceInterfaceType, currentLen-1, v.Len()-1)
			// Rebuild the slice excluding the deleted element
			for indexToRead := 0; indexToRead < currentLen; indexToRead++ {
				indexToWrite := indexToRead
				if indexToRead == i {
					continue
				}
				if indexToRead > i {
					indexToWrite--
				}
				updatedSlice.Index(indexToWrite).Set(
					v.Index(indexToRead).Value,
				)
			}

			v.Unpack().Set(updatedSlice)
		},
	}.
		WithMetadata("key", i).
		WithMetadata("parent", v)
}

// Append appends an empty value to the end of the slice.
func (v Value) Append() Value {
	currentLen := v.Len()
	newLen := currentLen + 1
	updatedSlice := reflect.MakeSlice(sliceInterfaceType, newLen, newLen)
	// copy all existing elements into updatedSlice.
	// this leaves the last element empty.
	for i := 0; i < currentLen; i++ {
		updatedSlice.Index(i).Set(
			v.Index(i).Value,
		)
	}

	v.FirstAddressable().Set(updatedSlice)

	// Set the last element to uninitialised.
	updatedSlice.Index(currentLen).Set(
		v.Index(currentLen).asUninitialised().Value,
	)

	return v
}

var sliceInterfaceType = reflect.TypeOf([]interface{}{})
var mapStringInterfaceType = reflect.TypeOf(map[string]interface{}{})

var UninitialisedPlaceholder interface{} = "__dasel_not_found__"

func (v Value) asUninitialised() Value {
	v.Value = reflect.ValueOf(UninitialisedPlaceholder)
	return v
}

func (v Value) initEmptyMap() Value {
	emptyMap := reflect.MakeMap(mapStringInterfaceType)
	v.Set(Value{Value: emptyMap})
	v.Value = emptyMap
	return v
}

func (v Value) initEmptydencodingMap() Value {
	om := dencoding.NewMap()
	rom := reflect.ValueOf(om)
	v.Set(Value{Value: rom})
	v.Value = rom
	return v
}

func (v Value) initEmptySlice() Value {
	emptySlice := reflect.MakeSlice(sliceInterfaceType, 0, 0)

	addressableSlice := reflect.New(emptySlice.Type())
	addressableSlice.Elem().Set(emptySlice)

	v.Set(Value{Value: addressableSlice})
	v.Value = addressableSlice
	return v
}

func makeAddressableSlice(value reflect.Value) reflect.Value {
	if !unpackReflectValue(value, reflect.Ptr).CanAddr() {
		unpacked := unpackReflectValue(value)

		emptySlice := reflect.MakeSlice(unpacked.Type(), unpacked.Len(), unpacked.Len())

		for i := 0; i < unpacked.Len(); i++ {
			emptySlice.Index(i).Set(makeAddressable(unpacked.Index(i)))
		}

		addressableSlice := reflect.New(emptySlice.Type())
		addressableSlice.Elem().Set(emptySlice)

		return addressableSlice
	} else {
		// Make contained values addressable
		unpacked := unpackReflectValue(value)
		for i := 0; i < unpacked.Len(); i++ {
			unpacked.Index(i).Set(makeAddressable(unpacked.Index(i)))
		}

		return value
	}
}

func makeAddressableMap(value reflect.Value) reflect.Value {
	if !unpackReflectValue(value, reflect.Ptr).CanAddr() {
		unpacked := unpackReflectValue(value)

		emptyMap := reflect.MakeMap(unpacked.Type())

		for _, key := range unpacked.MapKeys() {
			emptyMap.SetMapIndex(key, makeAddressable(unpacked.MapIndex(key)))
		}

		addressableMap := reflect.New(emptyMap.Type())
		addressableMap.Elem().Set(emptyMap)

		return addressableMap
	} else {
		// Make contained values addressable
		unpacked := unpackReflectValue(value)

		for _, key := range unpacked.MapKeys() {
			unpacked.SetMapIndex(key, makeAddressable(unpacked.MapIndex(key)))
		}

		return value
	}
}

func makeAddressable(value reflect.Value) reflect.Value {
	unpacked := unpackReflectValue(value)

	if isdencodingMap(unpacked) {
		om := value.Interface().(*dencoding.Map)
		for _, kv := range om.KeyValues() {
			var val any
			if v := deref(reflect.ValueOf(kv.Value)); v.IsValid() {
				val = makeAddressable(v).Interface()
			} else {
				val = nil
			}
			om.Set(kv.Key, val)
		}
		return value
	}

	switch unpacked.Kind() {
	case reflect.Slice:
		return makeAddressableSlice(value)
	case reflect.Map:
		return makeAddressableMap(value)
	default:
		return value
	}
}

func derefSlice(value reflect.Value) reflect.Value {
	unpacked := unpackReflectValue(value)

	res := reflect.MakeSlice(unpacked.Type(), unpacked.Len(), unpacked.Len())

	for i := 0; i < unpacked.Len(); i++ {
		if v := deref(unpacked.Index(i)); v.IsValid() {
			res.Index(i).Set(v)
		}
	}

	return res
}

func derefMap(value reflect.Value) reflect.Value {
	unpacked := unpackReflectValue(value)

	res := reflect.MakeMap(unpacked.Type())

	for _, key := range unpacked.MapKeys() {
		if v := deref(unpacked.MapIndex(key)); v.IsValid() {
			res.SetMapIndex(key, v)
		} else {
			res.SetMapIndex(key, reflect.ValueOf(new(any)))
		}
	}

	return res
}

func deref(value reflect.Value) reflect.Value {
	unpacked := unpackReflectValue(value)

	if isdencodingMap(unpacked) {
		om := value.Interface().(*dencoding.Map)
		for _, kv := range om.KeyValues() {
			if v := deref(reflect.ValueOf(kv.Value)); v.IsValid() {
				om.Set(kv.Key, v.Interface())
			} else {
				om.Set(kv.Key, nil)
			}
		}
		return value
	}

	switch unpacked.Kind() {
	case reflect.Slice:
		return derefSlice(value)
	case reflect.Map:
		return derefMap(value)
	default:
		return unpackReflectValue(value)
	}
}

// Values represents a list of Value's.
type Values []Value

// Interfaces returns the interface values for the underlying values stored in v.
func (v Values) Interfaces() []interface{} {
	res := make([]interface{}, 0)
	for _, val := range v {
		res = append(res, val.Interface())
	}
	return res
}

//func (v Values) initEmptyMaps() Values {
//	res := make(Values, len(v))
//	for k, value := range v {
//		if value.IsEmpty() {
//			res[k] = value.initEmptyMap()
//		} else {
//			res[k] = value
//		}
//	}
//	return res
//}

func (v Values) initEmptydencodingMaps() Values {
	res := make(Values, len(v))
	for k, value := range v {
		if value.IsEmpty() {
			res[k] = value.initEmptydencodingMap()
		} else {
			res[k] = value
		}
	}
	return res
}

func (v Values) initEmptySlices() Values {
	res := make(Values, len(v))
	for k, value := range v {
		if value.IsEmpty() {
			res[k] = value.initEmptySlice()
		} else {
			res[k] = value
		}
	}
	return res
}
