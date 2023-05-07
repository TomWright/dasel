package dencoding_test

import (
	"bytes"
	"github.com/tomwright/dasel/v2/dencoding"
	"io"
	"reflect"
	"testing"
)

func TestJSONDecoder_Decode(t *testing.T) {
	b := []byte(`{"x":1,"a":"hello"}{"x":2,"a":"there"}{"a":"Tom","x":3}`)
	dec := dencoding.NewJSONDecoder(bytes.NewReader(b))

	maps := make([]any, 0)
	for {
		var v any
		if err := dec.Decode(&v); err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("unexpected error: %v", err)
			return
		}
		maps = append(maps, v)
	}

	exp := [][]dencoding.KeyValue{
		{
			{Key: "x", Value: int64(1)},
			{Key: "a", Value: "hello"},
		},
		{
			{Key: "x", Value: int64(2)},
			{Key: "a", Value: "there"},
		},
		{
			{Key: "a", Value: "Tom"},
			{Key: "x", Value: int64(3)},
		},
	}

	got := make([][]dencoding.KeyValue, 0)
	for _, v := range maps {
		if m, ok := v.(*dencoding.Map); ok {
			got = append(got, m.KeyValues())
		}
	}

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
