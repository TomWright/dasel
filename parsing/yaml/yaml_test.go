package yaml_test

import (
	"bytes"
	"testing"

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
}
