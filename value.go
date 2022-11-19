package dasel

import (
	"fmt"
	"log"
	"reflect"
)

var deletePlaceholder = reflect.ValueOf("__dasel.delete_placeholder__")

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
	if v.IsEmpty() {
		return nil
	}
	unpacked := v.Unpack()
	if !unpacked.CanInterface() {
		return nil
	}
	return unpacked.Interface()
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
		return len(fmt.Sprint(v.Interface()))
	}
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

// IsDeletePlaceholder returns true is v represents a delete placeholder.
func (v Value) IsDeletePlaceholder() bool {
	return unpackReflectValue(v.Value) == deletePlaceholder
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

func unpackReflectValue(value reflect.Value, kinds ...reflect.Kind) reflect.Value {
	if len(kinds) == 0 {
		kinds = append(kinds, reflect.Ptr, reflect.Interface)
	}
	res := value
	for containsKind(kinds, res.Kind()) {
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
	log.Println("unable to set value with missing setFn")
}

// Delete deletes the current element.
// Depends on deleteFn since the implementation can differ depending on how the Value was initialised.
func (v Value) Delete() {
	if v.deleteFn != nil {
		v.deleteFn()
		return
	}
	log.Println("unable to delete value with missing deleteFn")
}

// MapIndex returns the value associated with key in the map v.
// It returns the zero Value if no field was found.
func (v Value) MapIndex(key Value) Value {
	return Value{
		Value: v.Unpack().MapIndex(key.Value),
		setFn: func(value Value) {
			v.Unpack().SetMapIndex(key.Value, value.Value)
		},
		deleteFn: func() {
			v.Unpack().SetMapIndex(key.Value, reflect.Value{})
		},
		metadata: map[string]interface{}{
			"type":   unpackReflectValue(v.Unpack().MapIndex(key.Value)).Kind().String(),
			"key":    key.Interface(),
			"parent": v,
		},
	}
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
		metadata: map[string]interface{}{
			"type":   unpackReflectValue(v.Unpack().FieldByName(name)).Kind().String(),
			"key":    name,
			"parent": v,
		},
	}
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

			v.Unpack(reflect.Ptr).Set(updatedSlice)
		},
		metadata: map[string]interface{}{
			"type":   unpackReflectValue(v.Unpack().Index(i)).Kind().String(),
			"key":    i,
			"parent": v,
		},
	}
}

// Append appends an empty value to the end of the slice.
func (v Value) Append() Value {
	emptyElement := reflect.ValueOf(UninitialisedPlaceholder)
	updatedSlice := reflect.Append(v.Unpack(), emptyElement)

	unpackedPtr := v.Unpack(reflect.Ptr)
	if unpackedPtr.CanSet() {
		unpackedPtr.Set(updatedSlice)
		return v
	}

	unpackedInterface := v.Unpack()
	if unpackedInterface.CanSet() {
		unpackedInterface.Set(updatedSlice)
		return v
	}

	if v.Value.CanSet() {
		v.Value.Set(updatedSlice)
		return v
	}

	panic("cannot find addressable element in slice")
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

func (v Value) initEmptySlice() Value {
	emptySlice := reflect.MakeSlice(sliceInterfaceType, 0, 0)

	addressableSlice := reflect.New(emptySlice.Type())
	addressableSlice.Elem().Set(emptySlice)

	v.Set(Value{Value: addressableSlice})
	v.Value = addressableSlice
	return v
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

func (v Values) initEmptyMaps() Values {
	res := make(Values, len(v))
	for k, value := range v {
		if value.IsEmpty() {
			res[k] = value.initEmptyMap()
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
