package toml_test

import (
	"io/ioutil"
	"path/filepath"
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

func TestTomlReader_MoreCases(t *testing.T) {
	// simple array of ints
	t.Run("simple array", tomlReaderTest([]byte(`nums = [1, 2, 3]`), func() *model.Value {
		res := model.NewMapValue()
		s := model.NewSliceValue()
		_ = s.Append(model.NewIntValue(1))
		_ = s.Append(model.NewIntValue(2))
		_ = s.Append(model.NewIntValue(3))
		_ = res.SetMapKey("nums", s)
		return res
	}))

	// mixed type array
	t.Run("mixed array", tomlReaderTest([]byte(`mix = [1, "two", true]`), func() *model.Value {
		res := model.NewMapValue()
		s := model.NewSliceValue()
		_ = s.Append(model.NewIntValue(1))
		_ = s.Append(model.NewStringValue("two"))
		_ = s.Append(model.NewBoolValue(true))
		_ = res.SetMapKey("mix", s)
		return res
	}))

	// inline nested table and array
	t.Run("inline nested table", tomlReaderTest([]byte(`props = { sub = { a = 1 }, arr = [1,2] }`), func() *model.Value {
		res := model.NewMapValue()
		props := model.NewMapValue()
		sub := model.NewMapValue()
		_ = sub.SetMapKey("a", model.NewIntValue(1))
		_ = props.SetMapKey("sub", sub)
		arr := model.NewSliceValue()
		_ = arr.Append(model.NewIntValue(1))
		_ = arr.Append(model.NewIntValue(2))
		_ = props.SetMapKey("arr", arr)
		_ = res.SetMapKey("props", props)
		return res
	}))

	// quoted key with space
	t.Run("quoted key with space", tomlReaderTest([]byte(`"a b" = "val"`), func() *model.Value {
		res := model.NewMapValue()
		_ = res.SetMapKey("a b", model.NewStringValue("val"))
		return res
	}))

	// dotted+quoted mixture
	t.Run("dotted and quoted mixture", tomlReaderTest([]byte(`a."b.c".d = "x"`), func() *model.Value {
		res := model.NewMapValue()
		a := model.NewMapValue()
		b := model.NewMapValue()
		_ = b.SetMapKey("d", model.NewStringValue("x"))
		_ = a.SetMapKey("b.c", b)
		_ = res.SetMapKey("a", a)
		return res
	}))

	// negative integer
	t.Run("negative integer", tomlReaderTest([]byte(`n = -5`), func() *model.Value {
		res := model.NewMapValue()
		_ = res.SetMapKey("n", model.NewIntValue(-5))
		return res
	}))

	// scientific float
	t.Run("scientific float", tomlReaderTest([]byte(`f = 1e3`), func() *model.Value {
		res := model.NewMapValue()
		_ = res.SetMapKey("f", model.NewFloatValue(1000.0))
		return res
	}))

	// array of inline tables
	t.Run("array of inline tables", tomlReaderTest([]byte(`items = [{a = 1}, {a = 2}]`), func() *model.Value {
		res := model.NewMapValue()
		items := model.NewSliceValue()
		p1 := model.NewMapValue()
		_ = p1.SetMapKey("a", model.NewIntValue(1))
		_ = items.Append(p1)
		p2 := model.NewMapValue()
		_ = p2.SetMapKey("a", model.NewIntValue(2))
		_ = items.Append(p2)
		_ = res.SetMapKey("items", items)
		return res
	}))

	// nested tables using headers
	t.Run("nested table headers", tomlReaderTest([]byte(`[server]
ip = "127.0.0.1"
[server.db]
name = "maindb"`), func() *model.Value {
		res := model.NewMapValue()
		server := model.NewMapValue()
		_ = server.SetMapKey("ip", model.NewStringValue("127.0.0.1"))
		db := model.NewMapValue()
		_ = db.SetMapKey("name", model.NewStringValue("maindb"))
		_ = server.SetMapKey("db", db)
		_ = res.SetMapKey("server", server)
		return res
	}))
}

func TestTomlReader_QuotedKeys(t *testing.T) {
	// quoted single key with dot should be a single key
	t.Run("quoted single segment containing dot", tomlReaderTest([]byte(`"a.b" = 1`), func() *model.Value {
		res := model.NewMapValue()
		_ = res.SetMapKey("a.b", model.NewIntValue(1))
		return res
	}))

	// unquoted dotted key should create nested maps
	t.Run("unquoted dotted key creates nested maps", tomlReaderTest([]byte(`a.b = 2`), func() *model.Value {
		res := model.NewMapValue()
		a := model.NewMapValue()
		_ = a.SetMapKey("b", model.NewIntValue(2))
		_ = res.SetMapKey("a", a)
		return res
	}))

	// mixture: first segment unquoted, second quoted containing dot
	t.Run("mixed quoted segment", tomlReaderTest([]byte(`a."b.c" = 3`), func() *model.Value {
		res := model.NewMapValue()
		a := model.NewMapValue()
		_ = a.SetMapKey("b.c", model.NewIntValue(3))
		_ = res.SetMapKey("a", a)
		return res
	}))
}

func TestTomlReader_ComplexFile(t *testing.T) {
	dataPath := filepath.Join("testdata", "complex_example.toml")
	b, err := ioutil.ReadFile(dataPath)
	if err != nil {
		t.Fatalf("failed reading test data: %v", err)
	}

	r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error creating reader: %v", err)
	}

	val, err := r.Read(b)
	if err != nil {
		t.Fatalf("unexpected error reading toml: %v", err)
	}

	// spot check some keys
	owner, err := val.GetMapKey("owner")
	if err != nil {
		t.Fatalf("missing owner: %v", err)
	}
	name, err := owner.GetMapKey("name")
	if err != nil {
		t.Fatalf("missing owner.name: %v", err)
	}
	if got, _ := name.StringValue(); got != "Tom Preston-Werner" {
		t.Fatalf("unexpected owner.name: %s", got)
	}

	// quoted key
	qk, err := val.GetMapKey("quoted key")
	if err != nil {
		t.Fatalf("missing quoted key: %v", err)
	}
	if s, _ := qk.StringValue(); s != "quoted value" {
		t.Fatalf("unexpected quoted key value: %s", s)
	}

	// products array length
	products, err := val.GetMapKey("products")
	if err != nil {
		t.Fatalf("missing products: %v", err)
	}
	if l, _ := products.SliceLen(); l != 2 {
		t.Fatalf("expected 2 products, got %d", l)
	}

	// a.b quoted
	ab, err := val.GetMapKey("a.b")
	if err != nil {
		t.Fatalf("missing a.b: %v", err)
	}
	if i, _ := ab.IntValue(); i != 42 {
		t.Fatalf("unexpected a.b: %d", i)
	}

	// nested table header value
	db, err := val.GetMapKey("database")
	if err != nil {
		t.Fatalf("missing database: %v", err)
	}
	rep, err := db.GetMapKey("replica")
	if err != nil {
		t.Fatalf("missing database.replica: %v", err)
	}
	if n, _ := rep.GetMapKey("name"); n == nil {
		t.Fatalf("missing database.replica.name")
	}
}

func TestTomlReader_EdgeCases(t *testing.T) {
	// conflict: scalar then dotted key
	t.Run("scalar then dotted key conflict", func(t *testing.T) {
		src := []byte("a = 1\na.b = 2")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		_, err = r.Read(src)
		if err == nil {
			t.Fatalf("expected error for scalar then dotted key conflict, got nil")
		}
	})

	// array-table conflict: scalar then array-table
	t.Run("scalar then array-table conflict", func(t *testing.T) {
		src := []byte("a = 1\n[[a]]\nname = \"x\"")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		_, err = r.Read(src)
		if err == nil {
			t.Fatalf("expected error for scalar then array-table conflict, got nil")
		}
	})

	// repeated table headers should merge keys
	t.Run("repeated explicit table headers merge", func(t *testing.T) {
		src := []byte("[t]\na = 1\n[t]\nb = 2")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		m, err := val.GetMapKey("t")
		if err != nil {
			t.Fatalf("missing table t: %v", err)
		}
		if a, err := m.GetMapKey("a"); err != nil || a == nil {
			t.Fatalf("expected a in t: %v", err)
		}
		if b, err := m.GetMapKey("b"); err != nil || b == nil {
			t.Fatalf("expected b in t: %v", err)
		}
	})

	// inline table then explicit table should merge
	t.Run("inline table then explicit header merges", func(t *testing.T) {
		src := []byte("t = {a = 1}\n[t]\nb = 2")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		m, err := val.GetMapKey("t")
		if err != nil {
			t.Fatalf("missing table t: %v", err)
		}
		if a, err := m.GetMapKey("a"); err != nil || a == nil {
			t.Fatalf("expected a in t: %v", err)
		}
		if b, err := m.GetMapKey("b"); err != nil || b == nil {
			t.Fatalf("expected b in t: %v", err)
		}
	})

	// integer overflow
	t.Run("integer overflow", func(t *testing.T) {
		src := []byte("big = 9223372036854775808") // int64 max + 1
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		_, err = r.Read(src)
		if err == nil {
			t.Fatalf("expected error for integer overflow, got nil")
		}
	})

	// ensure arrays with trailing commas parse (already covered elsewhere but add explicit check)
	t.Run("array trailing comma parse", func(t *testing.T) {
		src := []byte("arr = [1,2,]")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		arr, err := val.GetMapKey("arr")
		if err != nil {
			t.Fatalf("missing arr: %v", err)
		}
		if l, _ := arr.SliceLen(); l != 2 {
			t.Fatalf("expected 2 items in arr, got %d", l)
		}
	})

	// ensure arrays of inline tables types preserved
	t.Run("array of inline tables preserves types", func(t *testing.T) {
		src := []byte("items = [{a = 1}, {a = 2}]")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		items, err := val.GetMapKey("items")
		if err != nil {
			t.Fatalf("missing items: %v", err)
		}
		if l, _ := items.SliceLen(); l != 2 {
			t.Fatalf("expected 2 items, got %d", l)
		}
		// first element has a=1
		it0, _ := items.GetSliceIndex(0)
		if a, err := it0.GetMapKey("a"); err != nil {
			t.Fatalf("missing a in first item: %v", err)
		} else if ai, _ := a.IntValue(); ai != 1 {
			t.Fatalf("unexpected a in first item: %d", ai)
		}
	})
}

func TestTomlReader_TimeStrings(t *testing.T) {
	// Local date
	t.Run("local date string", func(t *testing.T) {
		src := []byte("d = 1979-05-27")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error parsing: %v", err)
		}
		v, err := val.GetMapKey("d")
		if err != nil {
			t.Fatalf("missing key d: %v", err)
		}
		s, err := v.StringValue()
		if err != nil {
			t.Fatalf("value not string: %v", err)
		}
		if s != "1979-05-27" {
			t.Fatalf("expected %q got %q", "1979-05-27", s)
		}
	})

	// Local time
	t.Run("local time string", func(t *testing.T) {
		src := []byte("t = 07:32:00")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error parsing: %v", err)
		}
		v, err := val.GetMapKey("t")
		if err != nil {
			t.Fatalf("missing key t: %v", err)
		}
		s, err := v.StringValue()
		if err != nil {
			t.Fatalf("value not string: %v", err)
		}
		if s != "07:32:00" {
			t.Fatalf("expected %q got %q", "07:32:00", s)
		}
	})

	// Local date-time
	t.Run("local datetime string", func(t *testing.T) {
		src := []byte("dt = 1979-05-27T07:32:00")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error parsing: %v", err)
		}
		v, err := val.GetMapKey("dt")
		if err != nil {
			t.Fatalf("missing key dt: %v", err)
		}
		s, err := v.StringValue()
		if err != nil {
			t.Fatalf("value not string: %v", err)
		}
		if s != "1979-05-27T07:32:00" {
			t.Fatalf("expected %q got %q", "1979-05-27T07:32:00", s)
		}
	})

	// DateTime with timezone (RFC3339)
	t.Run("datetime with tz string", func(t *testing.T) {
		src := []byte("dt = 1979-05-27T07:32:00-08:00")
		r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error creating reader: %v", err)
		}
		val, err := r.Read(src)
		if err != nil {
			t.Fatalf("unexpected error parsing: %v", err)
		}
		v, err := val.GetMapKey("dt")
		if err != nil {
			t.Fatalf("missing key dt: %v", err)
		}
		s, err := v.StringValue()
		if err != nil {
			t.Fatalf("value not string: %v", err)
		}
		if s != "1979-05-27T07:32:00-08:00" {
			t.Fatalf("expected %q got %q", "1979-05-27T07:32:00-08:00", s)
		}
	})
}
