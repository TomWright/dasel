package dasel

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIsTruthy(t *testing.T) {

	type testDef struct {
		name string
		in   interface{}
		out  bool
	}

	baseData := []testDef{
		{"bool:true", true, true},
		{"bool:false", false, false},
		{"string:lowercaseTrue", "true", true},
		{"string:lowercaseFalse", "false", false},
		{"string:uppercaseTrue", "TRUE", true},
		{"string:uppercaseFalse", "FALSE", false},
		{"string:lowercaseYes", "yes", true},
		{"string:lowercaseNo", "no", false},
		{"string:uppercaseYes", "YES", true},
		{"string:lowercaseNo", "NO", false},
		{"[]byte:lowercaseTrue", []byte("true"), true},
		{"[]byte:lowercaseFalse", []byte("false"), false},
		{"[]byte:uppercaseTrue", []byte("TRUE"), true},
		{"[]byte:uppercaseFalse", []byte("FALSE"), false},
		{"[]byte:lowercaseYes", []byte("yes"), true},
		{"[]byte:lowercaseNo", []byte("no"), false},
		{"[]byte:uppercaseYes", []byte("YES"), true},
		{"[]byte:lowercaseNo", []byte("NO"), false},
		{"int:0", int(0), false},
		{"int8:0", int8(0), false},
		{"int16:0", int16(0), false},
		{"int32:0", int32(0), false},
		{"int64:0", int64(0), false},
		{"int:-1", int(-1), false},
		{"int8:-1", int8(-1), false},
		{"int16:-1", int16(-1), false},
		{"int32:-1", int32(-1), false},
		{"int64:-1", int64(-1), false},
		{"uint:0", uint(0), false},
		{"uint8:0", uint8(0), false},
		{"uint16:0", uint16(0), false},
		{"uint32:0", uint32(0), false},
		{"uint64:0", uint64(0), false},
		{"int:1", int(1), true},
		{"int8:1", int8(1), true},
		{"int16:1", int16(1), true},
		{"int32:1", int32(1), true},
		{"int64:1", int64(1), true},
		{"uint:1", uint(1), true},
		{"uint8:1", uint8(1), true},
		{"uint16:1", uint16(1), true},
		{"uint32:1", uint32(1), true},
		{"uint64:1", uint64(1), true},
		{"float32:0", float32(0), false},
		{"float64:0", float64(0), false},
		{"float32:-1", float32(-1), false},
		{"float64:-1", float64(-1), false},
		{"float32:1", float32(1), true},
		{"float64:1", float64(1), true},
		{"unhandled:[]string", []string{}, false},
	}

	testData := make([]testDef, 0)

	for _, td := range baseData {
		testData = append(
			testData,
			td,
			testDef{
				name: fmt.Sprintf("reflect.Value:%s", td.name),
				in:   reflect.ValueOf(td.in),
				out:  td.out,
			},
			testDef{
				name: fmt.Sprintf("dasel.Value:%s", td.name),
				in:   ValueOf(td.in),
				out:  td.out,
			},
		)
	}

	for _, test := range testData {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			if exp, got := tc.out, IsTruthy(tc.in); exp != got {
				t.Errorf("expected %v, got %v", exp, got)
			}
		})
	}
}
