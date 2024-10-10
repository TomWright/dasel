package dencoding_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/tomwright/dasel/v3/dencoding"
)

func TestYAMLDecoder_Decode(t *testing.T) {

	t.Run("Basic", func(t *testing.T) {

		b := []byte(`
x: 1
a: hello
---
x: 2
a: there
---
a: Tom
x: 3
---`)
		dec := dencoding.NewYAMLDecoder(bytes.NewReader(b))

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
	})

	// https://github.com/TomWright/dasel/issues/278
	t.Run("Issue278", func(t *testing.T) {
		b := []byte(`
key1: [value1,value2,value3,value4,value5]
key2: value6
`)
		dec := dencoding.NewYAMLDecoder(bytes.NewReader(b))

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
			dencoding.NewMap().
				Set("key1", []any{"value1", "value2", "value3", "value4", "value5"}).
				Set("key2", "value6"),
		}

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("YamlAliases", func(t *testing.T) {
		b := []byte(`foo: &foofoo
  bar: 1
  baz: &baz "baz"
spam:
  ham: "eggs"
  bar: 0
  <<: *foofoo
  baz: "bazbaz"

baz: *baz
`)

		dec := dencoding.NewYAMLDecoder(bytes.NewReader(b))

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

		exp := dencoding.NewMap().
			Set("foo", dencoding.NewMap().
				Set("bar", int64(1)).
				Set("baz", "baz")).
			Set("spam", dencoding.NewMap().
				Set("ham", "eggs").
				Set("bar", int64(1)).
				Set("baz", "bazbaz")).
			Set("baz", "baz")

		if len(got) != 1 {
			t.Errorf("expected result len of %d, got %d", 1, len(got))
			return
		}

		gotMap, ok := got[0].(*dencoding.Map)
		if !ok {
			t.Errorf("expected result to be of type %T, got %T", exp, got[0])
			return
		}

		if !reflect.DeepEqual(exp, gotMap) {
			t.Errorf("expected %v, got %v", exp, gotMap)
		}
	})

}
