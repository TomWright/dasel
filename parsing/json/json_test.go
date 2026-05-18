package json_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/json"
)

func TestJson(t *testing.T) {
	doc := []byte(`{
    "string": "foo",
    "int": 1,
    "float": 1.1,
    "bool": true,
    "null": null,
    "array": [
        1,
        2,
        3
    ],
    "object": {
        "key": "value"
    }
}
`)
	reader, err := json.JSON.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatal(err)
	}
	writer, err := json.JSON.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatal(err)
	}

	value, err := reader.Read(doc)
	if err != nil {
		t.Fatal(err)
	}

	newDoc, err := writer.Write(value)
	if err != nil {
		t.Fatal(err)
	}

	if string(doc) != string(newDoc) {
		t.Fatalf("expected %s, got %s...\n%s", string(doc), string(newDoc), cmp.Diff(string(doc), string(newDoc)))
	}
}

func TestJsonCompact(t *testing.T) {
	doc := []byte(`{
    "string": "foo",
    "int": 1,
    "float": 1.1,
    "bool": true,
    "null": null,
    "array": [
        1,
        2,
        3
    ],
    "object": {
        "key": "value"
    }
}
`)
	reader, err := json.JSON.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatal(err)
	}
	opts := parsing.DefaultWriterOptions()
	opts.Compact = true
	writer, err := json.JSON.NewWriter(opts)
	if err != nil {
		t.Fatal(err)
	}

	value, err := reader.Read(doc)
	if err != nil {
		t.Fatal(err)
	}

	newDoc, err := writer.Write(value)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"string":"foo","int":1,"float":1.1,"bool":true,"null":null,"array":[1,2,3],"object":{"key":"value"}}` + "\n"
	if string(newDoc) != expected {
		t.Fatalf("expected %s, got %s...\n%s", expected, string(newDoc), cmp.Diff(expected, string(newDoc)))
	}
}

func TestNDJSON(t *testing.T) {
	newReader := func(t *testing.T) parsing.Reader {
		t.Helper()
		r, err := json.JSON.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	newWriter := func(t *testing.T, compact bool) parsing.Writer {
		t.Helper()
		opts := parsing.DefaultWriterOptions()
		opts.Compact = compact
		w, err := json.JSON.NewWriter(opts)
		if err != nil {
			t.Fatal(err)
		}
		return parsing.MultiDocumentWriter(w)
	}

	// readWrite is a helper that reads input, writes it back via MultiDocumentWriter,
	// and compares the output to the expected string.
	readWrite := func(t *testing.T, input, expected string, compact bool) {
		t.Helper()
		reader := newReader(t)
		writer := newWriter(t, compact)

		value, err := reader.Read([]byte(input))
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		out, err := writer.Write(value)
		if err != nil {
			t.Fatalf("write error: %v", err)
		}

		if string(out) != expected {
			t.Fatalf("unexpected output:\n%s", cmp.Diff(expected, string(out)))
		}
	}

	// assertValue is a helper that reads input and runs an assertion function against
	// the parsed model.Value.
	assertValue := func(t *testing.T, input string, assert func(t *testing.T, v *model.Value)) {
		t.Helper()
		reader := newReader(t)
		value, err := reader.Read([]byte(input))
		if err != nil {
			t.Fatalf("read error: %v", err)
		}
		assert(t, value)
	}

	t.Run("single object unchanged", func(t *testing.T) {
		input := "{\"name\":\"Tom\"}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if v.IsBranch() {
				t.Fatal("single document should not be a branch")
			}
			got, err := v.GetMapKey("name")
			if err != nil {
				t.Fatal(err)
			}
			s, err := got.StringValue()
			if err != nil {
				t.Fatal(err)
			}
			if s != "Tom" {
				t.Fatalf("expected Tom, got %s", s)
			}
		})
	})

	t.Run("single array unchanged", func(t *testing.T) {
		input := "[1,2,3]\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if v.IsBranch() {
				t.Fatal("single document should not be a branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 3 {
				t.Fatalf("expected length 3, got %d", length)
			}
		})
	})

	t.Run("single scalar unchanged", func(t *testing.T) {
		input := "42\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if v.IsBranch() {
				t.Fatal("single document should not be a branch")
			}
			i, err := v.IntValue()
			if err != nil {
				t.Fatal(err)
			}
			if i != 42 {
				t.Fatalf("expected 42, got %d", i)
			}
		})
	})

	t.Run("two objects", func(t *testing.T) {
		input := "{\"name\":\"Tom\"}\n{\"name\":\"Jim\"}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("multiple documents should be a branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}

			for i, expected := range []string{"Tom", "Jim"} {
				doc, err := v.GetSliceIndex(i)
				if err != nil {
					t.Fatal(err)
				}
				got, err := doc.GetMapKey("name")
				if err != nil {
					t.Fatal(err)
				}
				s, err := got.StringValue()
				if err != nil {
					t.Fatal(err)
				}
				if s != expected {
					t.Fatalf("doc %d: expected %s, got %s", i, expected, s)
				}
			}
		})
	})

	t.Run("three objects", func(t *testing.T) {
		input := "{\"a\":1}\n{\"b\":2}\n{\"c\":3}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 3 {
				t.Fatalf("expected 3 documents, got %d", length)
			}
		})
	})

	t.Run("mixed types", func(t *testing.T) {
		input := "{\"name\":\"Tom\"}\n[1,2,3]\n\"hello\"\n42\ntrue\nnull\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 6 {
				t.Fatalf("expected 6 documents, got %d", length)
			}

			// Check object
			doc0, _ := v.GetSliceIndex(0)
			name, _ := doc0.GetMapKey("name")
			s, _ := name.StringValue()
			if s != "Tom" {
				t.Fatalf("doc 0: expected Tom, got %s", s)
			}

			// Check array
			doc1, _ := v.GetSliceIndex(1)
			arrLen, _ := doc1.SliceLen()
			if arrLen != 3 {
				t.Fatalf("doc 1: expected length 3, got %d", arrLen)
			}

			// Check string
			doc2, _ := v.GetSliceIndex(2)
			s2, _ := doc2.StringValue()
			if s2 != "hello" {
				t.Fatalf("doc 2: expected hello, got %s", s2)
			}

			// Check int
			doc3, _ := v.GetSliceIndex(3)
			i, _ := doc3.IntValue()
			if i != 42 {
				t.Fatalf("doc 3: expected 42, got %d", i)
			}

			// Check bool
			doc4, _ := v.GetSliceIndex(4)
			b, _ := doc4.BoolValue()
			if !b {
				t.Fatal("doc 4: expected true")
			}

			// Check null
			doc5, _ := v.GetSliceIndex(5)
			if doc5.Type() != model.TypeNull {
				t.Fatalf("doc 5: expected null, got %s", doc5.Type())
			}
		})
	})

	t.Run("objects with no trailing newline", func(t *testing.T) {
		input := "{\"a\":1}\n{\"b\":2}"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}
		})
	})

	t.Run("extra whitespace between documents", func(t *testing.T) {
		input := "{\"a\":1}\n\n\n{\"b\":2}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}
		})
	})

	t.Run("pretty printed objects on separate lines", func(t *testing.T) {
		input := "{\n    \"name\": \"Tom\"\n}\n{\n    \"name\": \"Jim\"\n}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}
		})
	})

	t.Run("nested objects", func(t *testing.T) {
		input := "{\"user\":{\"name\":\"Tom\",\"age\":30}}\n{\"user\":{\"name\":\"Jim\",\"age\":25}}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}

			doc0, _ := v.GetSliceIndex(0)
			user0, _ := doc0.GetMapKey("user")
			name0, _ := user0.GetMapKey("name")
			s0, _ := name0.StringValue()
			if s0 != "Tom" {
				t.Fatalf("expected Tom, got %s", s0)
			}
			age0, _ := user0.GetMapKey("age")
			a0, _ := age0.IntValue()
			if a0 != 30 {
				t.Fatalf("expected 30, got %d", a0)
			}

			doc1, _ := v.GetSliceIndex(1)
			user1, _ := doc1.GetMapKey("user")
			name1, _ := user1.GetMapKey("name")
			s1, _ := name1.StringValue()
			if s1 != "Jim" {
				t.Fatalf("expected Jim, got %s", s1)
			}
		})
	})

	t.Run("empty input", func(t *testing.T) {
		assertValue(t, "", func(t *testing.T, v *model.Value) {
			if v.IsBranch() {
				t.Fatal("empty input should not be a branch")
			}
			if v.Type() != model.TypeNull {
				t.Fatalf("expected null, got %s", v.Type())
			}
		})
	})

	t.Run("whitespace only input", func(t *testing.T) {
		assertValue(t, "   \n\n  \n", func(t *testing.T, v *model.Value) {
			if v.IsBranch() {
				t.Fatal("whitespace-only input should not be a branch")
			}
			if v.Type() != model.TypeNull {
				t.Fatalf("expected null, got %s", v.Type())
			}
		})
	})

	t.Run("two scalar strings", func(t *testing.T) {
		input := "\"hello\"\n\"world\"\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			doc0, _ := v.GetSliceIndex(0)
			s0, _ := doc0.StringValue()
			if s0 != "hello" {
				t.Fatalf("expected hello, got %s", s0)
			}
			doc1, _ := v.GetSliceIndex(1)
			s1, _ := doc1.StringValue()
			if s1 != "world" {
				t.Fatalf("expected world, got %s", s1)
			}
		})
	})

	t.Run("two scalar ints", func(t *testing.T) {
		input := "1\n2\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			doc0, _ := v.GetSliceIndex(0)
			i0, _ := doc0.IntValue()
			if i0 != 1 {
				t.Fatalf("expected 1, got %d", i0)
			}
			doc1, _ := v.GetSliceIndex(1)
			i1, _ := doc1.IntValue()
			if i1 != 2 {
				t.Fatalf("expected 2, got %d", i1)
			}
		})
	})

	t.Run("two arrays", func(t *testing.T) {
		input := "[1,2]\n[3,4]\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, _ := v.SliceLen()
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}

			doc0, _ := v.GetSliceIndex(0)
			l0, _ := doc0.SliceLen()
			if l0 != 2 {
				t.Fatalf("doc 0: expected length 2, got %d", l0)
			}

			doc1, _ := v.GetSliceIndex(1)
			el, _ := doc1.GetSliceIndex(0)
			i, _ := el.IntValue()
			if i != 3 {
				t.Fatalf("doc 1[0]: expected 3, got %d", i)
			}
		})
	})

	t.Run("round trip compact two objects", func(t *testing.T) {
		input := "{\"name\":\"Tom\"}\n{\"name\":\"Jim\"}\n"
		// Each document gets a trailing \n from the writer, plus \n separator between docs.
		expected := "{\"name\":\"Tom\"}\n\n{\"name\":\"Jim\"}\n"
		readWrite(t, input, expected, true)
	})

	t.Run("round trip compact single object", func(t *testing.T) {
		input := "{\"name\":\"Tom\"}\n"
		expected := "{\"name\":\"Tom\"}\n"
		readWrite(t, input, expected, true)
	})

	t.Run("round trip pretty two objects", func(t *testing.T) {
		input := "{\"name\":\"Tom\"}\n{\"name\":\"Jim\"}\n"
		expected := "{\n    \"name\": \"Tom\"\n}\n\n{\n    \"name\": \"Jim\"\n}\n"
		readWrite(t, input, expected, false)
	})

	t.Run("round trip compact mixed types", func(t *testing.T) {
		input := "{\"a\":1}\n[1,2]\n\"str\"\n42\ntrue\nnull\n"
		expected := "{\"a\":1}\n\n[1,2]\n\n\"str\"\n\n42\n\ntrue\n\nnull\n"
		readWrite(t, input, expected, true)
	})

	t.Run("objects separated by spaces", func(t *testing.T) {
		input := "{\"a\":1}   {\"b\":2}"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, _ := v.SliceLen()
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}
		})
	})

	t.Run("documents with float values", func(t *testing.T) {
		input := "{\"val\":1.5}\n{\"val\":2.5}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			doc0, _ := v.GetSliceIndex(0)
			val0, _ := doc0.GetMapKey("val")
			f0, _ := val0.FloatValue()
			if f0 != 1.5 {
				t.Fatalf("expected 1.5, got %f", f0)
			}
			doc1, _ := v.GetSliceIndex(1)
			val1, _ := doc1.GetMapKey("val")
			f1, _ := val1.FloatValue()
			if f1 != 2.5 {
				t.Fatalf("expected 2.5, got %f", f1)
			}
		})
	})

	t.Run("documents with null values", func(t *testing.T) {
		input := "{\"a\":null}\n{\"b\":null}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			doc0, _ := v.GetSliceIndex(0)
			val0, _ := doc0.GetMapKey("a")
			if val0.Type() != model.TypeNull {
				t.Fatalf("expected null, got %s", val0.Type())
			}
		})
	})

	t.Run("documents with boolean values", func(t *testing.T) {
		input := "{\"ok\":true}\n{\"ok\":false}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			doc0, _ := v.GetSliceIndex(0)
			val0, _ := doc0.GetMapKey("ok")
			b0, _ := val0.BoolValue()
			if !b0 {
				t.Fatal("expected true")
			}
			doc1, _ := v.GetSliceIndex(1)
			val1, _ := doc1.GetMapKey("ok")
			b1, _ := val1.BoolValue()
			if b1 {
				t.Fatal("expected false")
			}
		})
	})

	t.Run("documents with empty objects", func(t *testing.T) {
		input := "{}\n{}\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, _ := v.SliceLen()
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}
		})
	})

	t.Run("documents with empty arrays", func(t *testing.T) {
		input := "[]\n[]\n"
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, _ := v.SliceLen()
			if length != 2 {
				t.Fatalf("expected 2 documents, got %d", length)
			}
		})
	})

	t.Run("large number of documents", func(t *testing.T) {
		input := ""
		for i := 0; i < 100; i++ {
			input += "{\"i\":" + string(rune('0'+i%10)) + "}\n"
		}
		// Use a simpler approach - just build manually
		input = ""
		for i := 0; i < 100; i++ {
			input += "{\"i\":1}\n"
		}
		assertValue(t, input, func(t *testing.T, v *model.Value) {
			if !v.IsBranch() {
				t.Fatal("expected branch")
			}
			length, err := v.SliceLen()
			if err != nil {
				t.Fatal(err)
			}
			if length != 100 {
				t.Fatalf("expected 100 documents, got %d", length)
			}
		})
	})
}
