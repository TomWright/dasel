package toml_test

import (
	"os"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/toml"
)

func TestTomlWriter_RoundTripSimple(t *testing.T) {
	doc := []byte(`title = "TOML Example"
[owner]
name = "Tom Preston-Werner"

[database]
ports = [8001, 8001, 8002]
enabled = true
`)

	reader, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error creating reader: %v", err)
	}
	writer, err := toml.TOML.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error creating writer: %v", err)
	}

	v, err := reader.Read(doc)
	if err != nil {
		t.Fatalf("failed to read doc: %v", err)
	}

	out, err := writer.Write(v)
	if err != nil {
		t.Fatalf("failed to write doc: %v", err)
	}

	v2, err := reader.Read(out)
	if err != nil {
		t.Fatalf("failed to read generated doc: %v", err)
	}

	res, err := v.Equal(v2)
	if err != nil {
		t.Fatalf("failed to compare values: %v", err)
	}
	b, err := res.BoolValue()
	if err != nil {
		t.Fatalf("failed to get bool from equal result: %v", err)
	}
	if !b {
		t.Fatalf("round-trip value mismatch\norig:\n%s\nnew:\n%s", string(doc), string(out))
	}
}

func TestTomlWriter_OrderPreserved(t *testing.T) {
	doc := []byte("a = 1\nb = 2\nc = 3\n")
	reader, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error creating reader: %v", err)
	}
	writer, err := toml.TOML.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error creating writer: %v", err)
	}

	v, err := reader.Read(doc)
	if err != nil {
		t.Fatalf("failed to read doc: %v", err)
	}

	out, err := writer.Write(v)
	if err != nil {
		t.Fatalf("failed to write doc: %v", err)
	}

	// Ensure key order a, b, c appears in the output
	outStr := string(out)
	ia := strings.Index(outStr, "a =")
	ib := strings.Index(outStr, "b =")
	ic := strings.Index(outStr, "c =")
	if ia == -1 || ib == -1 || ic == -1 {
		t.Fatalf("expected keys missing in output: %s", outStr)
	}
	if ia >= ib || ib >= ic {
		t.Fatalf("expected order a,b,c in output; got:\n%s", outStr)
	}
}

func TestTomlWriter_ArrayOfTables_RoundTrip(t *testing.T) {
	doc := []byte(`[[products]]
name = "Hammer"
sku = 738594937

[[products]]
name = "Screwdriver"
sku = 12341234
`)

	reader, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error creating reader: %v", err)
	}
	writer, err := toml.TOML.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error creating writer: %v", err)
	}

	v, err := reader.Read(doc)
	if err != nil {
		t.Fatalf("failed to read doc: %v", err)
	}

	out, err := writer.Write(v)
	if err != nil {
		t.Fatalf("failed to write doc: %v", err)
	}

	v2, err := reader.Read(out)
	if err != nil {
		t.Fatalf("failed to read generated doc: %v", err)
	}

	res, err := v.Equal(v2)
	if err != nil {
		t.Fatalf("failed to compare values: %v", err)
	}
	b, err := res.BoolValue()
	if err != nil {
		t.Fatalf("failed to get bool from equal result: %v", err)
	}
	if !b {
		t.Fatalf("array-table round-trip mismatch\norig:\n%s\nnew:\n%s", string(doc), string(out))
	}
}

func TestTomlWriter_MoreCases(t *testing.T) {
	cases := map[string][]byte{
		"simple array":           []byte("nums = [1, 2, 3]"),
		"mixed array":            []byte("mix = [1, \"two\", true]"),
		"inline nested table":    []byte("props = { sub = { a = 1 }, arr = [1,2] }"),
		"quoted key with space":  []byte("\"a b\" = \"val\""),
		"dotted and quoted mix":  []byte("a.\"b.c\".d = \"x\""),
		"negative integer":       []byte("n = -5"),
		"scientific float":       []byte("f = 1e3"),
		"array of inline tables": []byte("items = [{a = 1}, {a = 2} ]"),
		"nested table headers":   []byte("[server]\nip = \"127.0.0.1\"\n[server.db]\nname = \"maindb\""),
		"quoted single dot":      []byte("\"a.b\" = 1"),
		"unquoted dotted":        []byte("a.b = 2"),
		"mixed quoted segment":   []byte("a.\"b.c\" = 3"),
		"inline then explicit":   []byte("t = {a = 1}\n[t]\nb = 2"),
		"array trailing comma":   []byte("arr = [1,2,]"),
		"local date":             []byte("d = 1979-05-27"),
		"local time":             []byte("t = 07:32:00"),
		"local datetime":         []byte("dt = 1979-05-27T07:32:00"),
		"datetime with tz":       []byte("dt = 1979-05-27T07:32:00-08:00"),
		"multiline basic string": []byte("m = '''not used'''\n"),
	}

	reader, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error creating reader: %v", err)
	}
	writer, err := toml.TOML.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error creating writer: %v", err)
	}

	for name, src := range cases {
		src := src
		name := name
		t.Run(name, func(t *testing.T) {
			v, err := reader.Read(src)
			if err != nil {
				t.Fatalf("reader error for %s: %v", name, err)
			}

			out, err := writer.Write(v)
			if err != nil {
				t.Fatalf("writer error for %s: %v", name, err)
			}

			v2, err := reader.Read(out)
			if err != nil {
				t.Fatalf("reader error for generated output %s: %v", name, err)
			}

			res, err := v.Equal(v2)
			if err != nil {
				t.Fatalf("compare error for %s: %v", name, err)
			}
			b, err := res.BoolValue()
			if err != nil {
				t.Fatalf("bool extraction error for %s: %v", name, err)
			}
			if !b {
				t.Fatalf("round-trip mismatch for %s\norig:\n%s\nnew:\n%s", name, string(src), string(out))
			}
		})
	}

	// Complex example file round-trip
	t.Run("complex example file", func(t *testing.T) {
		//t.Skip("Multiline string formatting not yet preserved")
		dataPath := "testdata/complex_example.toml"
		b, err := os.ReadFile(dataPath)
		if err != nil {
			t.Fatalf("failed reading test data: %v", err)
		}

		v, err := reader.Read(b)
		if err != nil {
			t.Fatalf("reader error for complex file: %v", err)
		}

		out, err := writer.Write(v)
		if err != nil {
			t.Fatalf("writer error for complex file: %v", err)
		}

		v2, err := reader.Read(out)
		if err != nil {
			t.Fatalf("failed to re-read generated complex doc: %v", err)
		}

		res, err := v.Equal(v2)
		if err != nil {
			t.Fatalf("compare error for complex file: %v", err)
		}
		b2, err := res.BoolValue()
		if err != nil {
			t.Fatalf("bool extraction error for complex file: %v", err)
		}
		if !b2 {
			// Print original and written output for debugging.
			t.Fatalf("complex-example round-trip mismatch\n--- ORIGINAL ---\n%s\n--- WRITTEN ---\n%s\n", string(b), string(out))
		}
	})
}

func TestTomlWriter_StrictOutput(t *testing.T) {
	reader, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error creating reader: %v", err)
	}
	writer, err := toml.TOML.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error creating writer: %v", err)
	}

	tests := map[string]struct {
		src []byte
		exp string
	}{
		"ordered scalars": {
			src: []byte("a = 1\nb = 2\nc = 3\n"),
			exp: "a = 1\nb = 2\nc = 3\n",
		},
		"inline array": {
			src: []byte("nums = [1,2,3]"),
			exp: "nums = [1, 2, 3]\n",
		},
		"quoted key": {
			src: []byte("\"a b\" = \"val\""),
			exp: "'a b' = 'val'\n",
		},
		"array of tables": {
			src: []byte(`[[products]]
name = "Hammer"
sku = 738594937

[[products]]
name = "Screwdriver"
sku = 12341234
`),
			exp: `[[products]]
name = 'Hammer'
sku = 738594937

[[products]]
name = 'Screwdriver'
sku = 12341234
`,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			v, err := reader.Read(tc.src)
			if err != nil {
				t.Fatalf("reader error for %s: %v", name, err)
			}

			out, err := writer.Write(v)
			if err != nil {
				t.Fatalf("writer error for %s: %v", name, err)
			}

			got := string(out)
			if got != tc.exp {
				t.Fatalf("strict output mismatch for %s\nexpected:\n%s\n got:\n%s", name, tc.exp, got)
			}
		})
	}
}
