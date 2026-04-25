package yaml_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/yaml"
)

type testCase struct {
	in     string
	assert func(t *testing.T, res *model.Value)
}

func (tc testCase) run(t *testing.T) {
	r, err := yaml.YAML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	res, err := r.Read([]byte(tc.in))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tc.assert(t, res)
}

type rwTestCase struct {
	in  string
	out string
}

func (tc rwTestCase) run(t *testing.T) {
	if tc.out == "" {
		tc.out = tc.in
	}
	r, err := yaml.YAML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	w, err := yaml.YAML.NewWriter(parsing.WriterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	res, err := r.Read([]byte(tc.in))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	out, err := w.Write(res)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !bytes.Equal([]byte(tc.out), out) {
		t.Errorf("unexpected output: %s", cmp.Diff(tc.out, string(out)))
	}
}

func TestYamlWriter_Compact(t *testing.T) {
	r, err := yaml.YAML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	opts := parsing.DefaultWriterOptions()
	opts.Compact = true
	w, err := yaml.YAML.NewWriter(opts)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	t.Run("map", func(t *testing.T) {
		res, err := r.Read([]byte("a: 1\nb: 2\n"))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		out, err := w.Write(res)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		expected := "{a: 1, b: 2}\n"
		if string(out) != expected {
			t.Errorf("expected %q, got %q", expected, string(out))
		}
	})

	t.Run("sequence", func(t *testing.T) {
		res, err := r.Read([]byte("- 1\n- 2\n- 3\n"))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		out, err := w.Write(res)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		expected := "[1, 2, 3]\n"
		if string(out) != expected {
			t.Errorf("expected %q, got %q", expected, string(out))
		}
	})

	t.Run("nested", func(t *testing.T) {
		res, err := r.Read([]byte("map:\n  key: value\nlist:\n  - item1\n  - item2\n"))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		out, err := w.Write(res)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		expected := "{map: {key: value}, list: [item1, item2]}\n"
		if string(out) != expected {
			t.Errorf("expected %q, got %q", expected, string(out))
		}
	})
}

func TestYamlValue_UnmarshalYAML(t *testing.T) {
	t.Run("simple key value", testCase{
		in: `name: Tom`,
		assert: func(t *testing.T, res *model.Value) {
			got, err := res.GetMapKey("name")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			gotStr, err := got.StringValue()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if gotStr != "Tom" {
				t.Errorf("unexpected value: %s", gotStr)
			}
		},
	}.run)

	t.Run("multi document", testCase{
		in: `name: Tom
---
name: Jerry`,
		assert: func(t *testing.T, res *model.Value) {
			a, err := res.GetSliceIndex(0)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			got, err := a.GetMapKey("name")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			gotStr, err := got.StringValue()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if gotStr != "Tom" {
				t.Errorf("unexpected value: %s", gotStr)
			}

			b, err := res.GetSliceIndex(1)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			got, err = b.GetMapKey("name")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			gotStr, err = got.StringValue()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if gotStr != "Jerry" {
				t.Errorf("unexpected value: %s", gotStr)
			}
		},
	}.run)

	t.Run("multi document", testCase{
		in: `name: Tom
---
name: Jerry`,
		assert: func(t *testing.T, res *model.Value) {
			a, err := res.GetSliceIndex(0)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			got, err := a.GetMapKey("name")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			gotStr, err := got.StringValue()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if gotStr != "Tom" {
				t.Errorf("unexpected value: %s", gotStr)
			}

			b, err := res.GetSliceIndex(1)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			got, err = b.GetMapKey("name")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			gotStr, err = got.StringValue()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if gotStr != "Jerry" {
				t.Errorf("unexpected value: %s", gotStr)
			}
		},
	}.run)

	t.Run("multi document", rwTestCase{
		in: `name: Tom
---
name: Jerry
`,
	}.run)

	t.Run("generic", rwTestCase{
		in: `str: foo
int: 1
float: 1.1
bool: true
map:
    key: value
list:
    - item1
    - item2
`,
	}.run)

	// This test is technically wrong because we're only supporting the alias on read and not write.
	t.Run("alias", rwTestCase{
		in: `name: &name Tom
name2: *name
`,
		out: `name: Tom
name2: Tom
`,
	}.run)

	t.Run("null read write", rwTestCase{
		in: `name: null
`,
		out: `name: null
`,
	}.run)

	t.Run("null document read write", rwTestCase{
		in: `null
`,
		out: `null
`,
	}.run)

	t.Run("base numbers", func(t *testing.T) {
		t.Run("standard", rwTestCase{
			in: `10
`,
			out: `10
`,
		}.run)

		t.Run("zero", rwTestCase{
			in: `0
`,
			out: `0
`,
		}.run)

		t.Run("negative", rwTestCase{
			in: `-42
`,
			out: `-42
`,
		}.run)

		t.Run("hex lowercase", rwTestCase{
			in: `0x10
`,
			out: `16
`,
		}.run)

		t.Run("hex uppercase letters", rwTestCase{
			in: `0xff
`,
			out: `255
`,
		}.run)

		t.Run("octal", rwTestCase{
			in: `0o10
`,
			out: `8
`,
		}.run)

		t.Run("binary", rwTestCase{
			in: `0b10
`,
			out: `2
`,
		}.run)

		t.Run("leading zero is decimal", rwTestCase{
			in: `010
`,
			out: `10
`,
		}.run)

		t.Run("hex in map", rwTestCase{
			in: `val: 0x10
`,
			out: `val: 16
`,
		}.run)

		t.Run("octal in map", rwTestCase{
			in: `val: 0o77
`,
			out: `val: 63
`,
		}.run)

		t.Run("mixed types in map", rwTestCase{
			in: `dec: 42
hex: 0xff
oct: 0o77
bin: 0b1010
`,
			out: `dec: 42
hex: 255
oct: 63
bin: 10
`,
		}.run)

		t.Run("positive sign", rwTestCase{
			in: `+42
`,
			out: `42
`,
		}.run)

		t.Run("positive hex", rwTestCase{
			in: `+0x10
`,
			out: `16
`,
		}.run)

		t.Run("positive octal", rwTestCase{
			in: `+0o10
`,
			out: `8
`,
		}.run)

		t.Run("positive binary", rwTestCase{
			in: `+0b10
`,
			out: `2
`,
		}.run)

		t.Run("negative hex", rwTestCase{
			in: `-0x10
`,
			out: `-16
`,
		}.run)

		t.Run("negative octal", rwTestCase{
			in: `-0o10
`,
			out: `-8
`,
		}.run)

		t.Run("negative binary", rwTestCase{
			in: `-0b10
`,
			out: `-2
`,
		}.run)

		t.Run("underscore decimal", rwTestCase{
			in: `1_000
`,
			out: `1000
`,
		}.run)

		t.Run("underscore hex", rwTestCase{
			in: `0xFF_FF
`,
			out: `65535
`,
		}.run)

		t.Run("underscore binary", rwTestCase{
			in: `0b1010_1010
`,
			out: `170
`,
		}.run)
	})

	t.Run("bounded yaml expansion", func(t *testing.T) {
		in := `a: &a ["lol","lol","lol","lol","lol","lol","lol","lol","lol"]
b: &b [*a,*a,*a,*a,*a,*a,*a,*a,*a]
c: &c [*b,*b,*b,*b,*b,*b,*b,*b,*b]
d: &d [*c,*c,*c,*c,*c,*c,*c,*c,*c]
e: &e [*d,*d,*d,*d,*d,*d,*d,*d,*d]
f: &f [*e,*e,*e,*e,*e,*e,*e,*e,*e]
g: &g [*f,*f,*f,*f,*f,*f,*f,*f,*f]
h: &h [*g,*g,*g,*g,*g,*g,*g,*g,*g]
i: &i [*h,*h,*h,*h,*h,*h,*h,*h,*h]
`

		reader, err := parsing.Format("yaml").NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		var gotErr error

		maxWaitTime := 10 * time.Second
		gotErrCh := make(chan error)
		go func() {
			_, gotErr = reader.Read([]byte(in))
			gotErrCh <- gotErr
		}()

		select {
		case gotErr = <-gotErrCh:
			if gotErr == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(gotErr, yaml.ErrYamlExpansionDepthExceeded) && !errors.Is(gotErr, yaml.ErrYamlExpansionBudgetExceeded) {
				t.Fatalf("unexpected error: %s", gotErr)
			}
		case <-time.After(maxWaitTime):
			t.Fatalf("expected error within %s, but did not get one", maxWaitTime)
		}
	})

	t.Run("alias metadata is preserved", func(t *testing.T) {
		r, err := yaml.YAML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		res, err := r.Read([]byte("name: &name Tom\nname2: *name\n"))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		aliasValue, err := res.GetMapKey("name2")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		got, ok := aliasValue.MetadataValue("yaml-alias")
		if !ok {
			t.Fatal("expected yaml-alias metadata to be set")
		}
		if got != "name" {
			t.Fatalf("unexpected yaml-alias metadata: %v", got)
		}
	})

	t.Run("yaml expansion depth boundary", func(t *testing.T) {
		reader, err := parsing.Format("yaml").NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		makeDoc := func(depth int) string {
			if depth == 0 {
				return "root0: &root0 [value]\nresult: *root0\n"
			}

			res := "root0: &root0 [value]\n"
			for i := 1; i <= depth; i++ {
				res += fmt.Sprintf("root%d: &root%d [*root%d]\n", i, i, i-1)
			}
			res += fmt.Sprintf("result: *root%d\n", depth)
			return res
		}

		depthLimit := 32
		for _, tc := range []struct {
			name    string
			depth   int
			wantErr bool
		}{
			{name: "within limit", depth: depthLimit - 2, wantErr: false},
			{name: "at limit", depth: depthLimit - 1, wantErr: false},
			{name: "over limit", depth: depthLimit, wantErr: true},
		} {
			t.Run(tc.name, func(t *testing.T) {
				res, err := reader.Read([]byte(makeDoc(tc.depth)))
				if tc.wantErr {
					if err == nil {
						t.Fatal("expected error, got nil")
					}
					if !errors.Is(err, yaml.ErrYamlExpansionDepthExceeded) {
						t.Fatalf("unexpected error: %v", err)
					}
					return
				}

				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				result, err := res.GetMapKey("result")
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				for i := 0; i <= tc.depth; i++ {
					result, err = result.GetSliceIndex(0)
					if err != nil {
						t.Fatalf("unexpected error at nested index %d: %s", i, err)
					}
				}
				got, err := result.StringValue()
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if got != "value" {
					t.Fatalf("unexpected value: %s", got)
				}
			})
		}
	})

	t.Run("yaml expansion budget boundary", func(t *testing.T) {
		reader, err := parsing.Format("yaml").NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		makeDoc := func(aliasCount int) string {
			res := "root: &root value\nitems:\n"
			for i := 0; i < aliasCount; i++ {
				res += "  - *root\n"
			}
			return res
		}

		res, err := reader.Read([]byte(makeDoc(1000)))
		if err != nil {
			t.Fatalf("unexpected error at budget boundary: %s", err)
		}
		items, err := res.GetMapKey("items")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		last, err := items.GetSliceIndex(999)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		got, err := last.StringValue()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if got != "value" {
			t.Fatalf("unexpected value: %s", got)
		}

		_, err = reader.Read([]byte(makeDoc(1001)))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, yaml.ErrYamlExpansionBudgetExceeded) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("yaml expansion budget resets per document", func(t *testing.T) {
		reader, err := parsing.Format("yaml").NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		makeDoc := func(aliasCount int) string {
			res := "root: &root value\nitems:\n"
			for i := 0; i < aliasCount; i++ {
				res += "  - *root\n"
			}
			return res
		}

		multiDoc := makeDoc(1000) + "---\n" + makeDoc(1000)
		res, err := reader.Read([]byte(multiDoc))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		first, err := res.GetSliceIndex(0)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		second, err := res.GetSliceIndex(1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		for _, doc := range []*model.Value{first, second} {
			items, err := doc.GetMapKey("items")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			last, err := items.GetSliceIndex(999)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			got, err := last.StringValue()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got != "value" {
				t.Fatalf("unexpected value: %s", got)
			}
		}
	})
}
