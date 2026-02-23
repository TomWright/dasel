package execution_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

type testCase struct {
	in          *model.Value
	inFn        func() *model.Value
	s           string
	out         *model.Value
	outFn       func() *model.Value
	compareRoot bool
	opts        []execution.ExecuteOptionFn
}

func (tc testCase) run(t *testing.T) {
	in := tc.in
	if tc.inFn != nil {
		in = tc.inFn()
	}
	if in == nil {
		in = model.NewValue(nil)
	}
	exp := tc.out
	if tc.outFn != nil {
		exp = tc.outFn()
	}
	res, err := execution.ExecuteSelector(context.Background(), tc.s, in, execution.NewOptions(tc.opts...))
	if err != nil {
		t.Fatal(err)
	}

	if tc.compareRoot {
		res = in
	}

	equal, err := res.EqualTypeValue(exp)
	if err != nil {
		t.Fatal(err)
	}
	if !equal {
		t.Errorf("unexpected output:\nexp: %s\ngot: %s", exp.String(), res.String())
	}

	expMeta := exp.Metadata
	gotMeta := res.Metadata
	if !cmp.Equal(expMeta, gotMeta) {
		t.Errorf("unexpected output metadata: %v", cmp.Diff(expMeta, gotMeta))
	}
}

func TestExecuteSelector_HappyPath(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		inputMap := func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("title", "Mr").
					Set("age", int64(31)).
					Set("name", orderedmap.NewMap().
						Set("first", "Tom").
						Set("last", "Wright")),
			)
		}
		t.Run("property", testCase{
			in:  inputMap(),
			s:   `title`,
			out: model.NewStringValue("Mr"),
		}.run)
		t.Run("nested property", testCase{
			in:  inputMap(),
			s:   `name.first`,
			out: model.NewStringValue("Tom"),
		}.run)
		t.Run("concat with grouping", testCase{
			in:  inputMap(),
			s:   `title + " " + (name.first) + " " + (name.last)`,
			out: model.NewStringValue("Mr Tom Wright"),
		}.run)
		t.Run("concat", testCase{
			in:  inputMap(),
			s:   `title + " " + name.first + " " + name.last`,
			out: model.NewStringValue("Mr Tom Wright"),
		}.run)
		t.Run("add evaluated fields", testCase{
			in: inputMap(),
			s:  `{..., "over30": age > 30}`,
			outFn: func() *model.Value {
				return model.NewValue(
					orderedmap.NewMap().
						Set("title", "Mr").
						Set("age", int64(31)).
						Set("name", orderedmap.NewMap().
							Set("first", "Tom").
							Set("last", "Wright")).
						Set("over30", true),
				)
			},
		}.run)
	})

	t.Run("set", func(t *testing.T) {
		inputMap := func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("title", "Mr").
					Set("age", int64(31)).
					Set("name", orderedmap.NewMap().
						Set("first", "Tom").
						Set("last", "Wright")),
			)
		}
		inputSlice := func() *model.Value {
			return model.NewValue([]any{1, 2, 3})
		}

		t.Run("set property", testCase{
			in: inputMap(),
			s:  `title = "Mrs"`,
			outFn: func() *model.Value {
				res := inputMap()
				if err := res.SetMapKey("title", model.NewStringValue("Mrs")); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return res
			},
			compareRoot: true,
		}.run)

		t.Run("set index", testCase{
			in: inputSlice(),
			s:  `$this[1] = 4`,
			outFn: func() *model.Value {
				res := inputSlice()
				if err := res.SetSliceIndex(1, model.NewIntValue(4)); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return res
			},
			compareRoot: true,
		}.run)
	})
}
