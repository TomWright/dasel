package toml_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/toml"
)

func tomlReaderTest(data []byte, exp func() *model.Value) func(*testing.T) {
	return func(t *testing.T) {
		exp := exp()
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := r.Read(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		matchResult, err := got.Equal(exp)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		matchResultBool, err := matchResult.BoolValue()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !matchResultBool {
			t.Errorf("expected\n%s\ngot\n%s", exp.String(), got.String())
		}
	}
}

func TestTomlReader_Read(t *testing.T) {
	t.Run("key value", func(t *testing.T) {
		t.Run(
			"string",
			tomlReaderTest([]byte(`foo = "Bar"`), func() *model.Value {
				res := model.NewMapValue()
				_ = res.SetMapKey("foo", model.NewStringValue("Bar"))
				return res
			}),
		)

		t.Run(
			"int",
			tomlReaderTest([]byte(`foo = 123`), func() *model.Value {
				res := model.NewMapValue()
				_ = res.SetMapKey("foo", model.NewIntValue(123))
				return res
			}),
		)

		t.Run(
			"float",
			tomlReaderTest([]byte(`foo = 12.3`), func() *model.Value {
				res := model.NewMapValue()
				_ = res.SetMapKey("foo", model.NewFloatValue(12.3))
				return res
			}),
		)

		t.Run(
			"true",
			tomlReaderTest([]byte(`foo = true`), func() *model.Value {
				res := model.NewMapValue()
				_ = res.SetMapKey("foo", model.NewBoolValue(true))
				return res
			}),
		)

		t.Run(
			"false",
			tomlReaderTest([]byte(`foo = false`), func() *model.Value {
				res := model.NewMapValue()
				_ = res.SetMapKey("foo", model.NewBoolValue(false))
				return res
			}),
		)
	})

	t.Run("inline table", tomlReaderTest([]byte(`props = { key1 = "value1", key2 = "value2" }`), func() *model.Value {
		res := model.NewMapValue()
		inlineTable := model.NewMapValue()
		_ = inlineTable.SetMapKey("key1", model.NewStringValue("value1"))
		_ = inlineTable.SetMapKey("key2", model.NewStringValue("value2"))
		_ = res.SetMapKey("props", inlineTable)
		return res
	}))

	t.Run("table", tomlReaderTest([]byte(`[props]
key1 = "value1"
key2 = "value2"`), func() *model.Value {
		res := model.NewMapValue()
		inlineTable := model.NewMapValue()
		_ = inlineTable.SetMapKey("key1", model.NewStringValue("value1"))
		_ = inlineTable.SetMapKey("key2", model.NewStringValue("value2"))
		_ = res.SetMapKey("props", inlineTable)
		return res
	}))

	t.Run("array table", tomlReaderTest([]byte(`[[products]]
name = "foo"
id = 1

[[products]]
name = "bar"
id = 2`), func() *model.Value {
		productsArray := model.NewSliceValue()
		product1 := model.NewMapValue()
		_ = product1.SetMapKey("name", model.NewStringValue("foo"))
		_ = product1.SetMapKey("id", model.NewIntValue(1))
		_ = productsArray.Append(product1)

		product2 := model.NewMapValue()
		_ = product2.SetMapKey("name", model.NewStringValue("bar"))
		_ = product2.SetMapKey("id", model.NewIntValue(2))
		_ = productsArray.Append(product2)

		res := model.NewMapValue()
		_ = res.SetMapKey("products", productsArray)
		return res
	}))

	//	dataBytes := []byte(`title = "TOML Example"
	//[owner]
	//name = "Tom Preston-Werner"
	//props = { key1 = "value1", key2 = "value2" }
	//
	//[[products]]
	//name = "Hammer"
	//sku = 738594937
	//
	//[[products]]
	//name = "Screwdriver"
	//sku = 12341234
	//`)
	//	//dob = 1979-05-27T07:32:00-08:00
	//	dataModel := model.NewMapValue()
	//	//parsedTime, err := time.Parse(time.RFC3339, "1979-05-27T07:32:00-08:00")
	//	//if err != nil {
	//	//	t.Fatalf("unexpected error: %v", err)
	//	//}
	//	ownerMap := model.NewMapValue()
	//	_ = ownerMap.SetMapKey("name", model.NewStringValue("Tom Preston-Werner"))
	//	//_ = ownerMap.SetMapKey("dob", model.NewValue(parsedTime))
	//	_ = dataModel.SetMapKey("title", model.NewStringValue("TOML Example"))
	//	_ = dataModel.SetMapKey("owner", ownerMap)
}
