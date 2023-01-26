package dencoding_test

import (
	"bytes"
	"github.com/tomwright/dasel/v2/dencoding"
	"io"
	"reflect"
	"testing"
)

func TestTOMLDecoder_Decode(t *testing.T) {

	t.Run("KeyValue", func(t *testing.T) {
		b := []byte(`x = 1
a = 'hello'`)
		dec := dencoding.NewTOMLDecoder(bytes.NewReader(b))

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

		exp := []any{
			map[string]any{
				"x": int64(1),
				"a": "hello",
			},
		}

		got := maps

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("Table", func(t *testing.T) {
		b := []byte(`
[user]
name = "Tom"
age = 29
`)
		dec := dencoding.NewTOMLDecoder(bytes.NewReader(b))

		got := make([]any, 0)
		for {
			var v any
			if err := dec.Decode(&v); err != nil {
				if err == io.EOF {
					break
				}
				t.Errorf("unexpected error: %v", err)
				return
			}
			got = append(got, v)
		}

		exp := []any{
			map[string]any{
				"user": map[string]any{
					"age":  int64(29),
					"name": "Tom",
				},
			},
		}

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
}
