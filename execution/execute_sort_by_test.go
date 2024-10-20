package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncSortBy(t *testing.T) {
	runSortTests := func(in func() *model.Value, outAsc func() *model.Value, outDesc func() *model.Value) func(*testing.T) {
		return func(t *testing.T) {
			t.Run("asc default", testCase{
				inFn:  in,
				s:     `sortBy($this)`,
				outFn: outAsc,
			}.run)
			t.Run("asc", testCase{
				inFn:  in,
				s:     `sortBy($this, asc)`,
				outFn: outAsc,
			}.run)
			t.Run("desc", testCase{
				inFn:  in,
				s:     `sortBy($this, desc)`,
				outFn: outDesc,
			}.run)
		}
	}

	t.Run("int", runSortTests(
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewIntValue(2)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(1)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(4)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(3)); err != nil {
				t.Fatal(err)
			}
			return res
		},
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewIntValue(1)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(2)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(3)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(4)); err != nil {
				t.Fatal(err)
			}
			return res
		},
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewIntValue(4)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(3)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(2)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewIntValue(1)); err != nil {
				t.Fatal(err)
			}
			return res
		},
	))

	t.Run("float", runSortTests(
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewFloatValue(2.23)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(2)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(5.123)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(4.2)); err != nil {
				t.Fatal(err)
			}
			return res
		},
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewFloatValue(2)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(2.23)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(4.2)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(5.123)); err != nil {
				t.Fatal(err)
			}
			return res
		},
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewFloatValue(5.123)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(4.2)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(2.23)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewFloatValue(2)); err != nil {
				t.Fatal(err)
			}
			return res
		},
	))
	t.Run("string", runSortTests(
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewStringValue("def")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("abc")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("cde")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("bcd")); err != nil {
				t.Fatal(err)
			}
			return res
		},
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewStringValue("abc")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("bcd")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("cde")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("def")); err != nil {
				t.Fatal(err)
			}
			return res
		},
		func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewStringValue("def")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("cde")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("bcd")); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewStringValue("abc")); err != nil {
				t.Fatal(err)
			}
			return res
		},
	))
}
