package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestBranch(t *testing.T) {
	t.Run("single branch", testCase{
		s: "branch(1)",
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			r.MarkAsBranch()
			if err := r.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("many branches", testCase{
		s: "branch(1, 1+1, 3/1, 123)",
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			r.MarkAsBranch()
			if err := r.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(3)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(123)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("spread into many branches", testCase{
		s: "[1,2,3].branch(...)",
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			r.MarkAsBranch()
			if err := r.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(3)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("chained branch set", testCase{
		s: "branch(1, 2, 3).$this=5",
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			r.MarkAsBranch()
			if err := r.Append(model.NewIntValue(5)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(5)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(5)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("chained branch math", testCase{
		s: "(branch(1, 2, 3)) * 2",
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			r.MarkAsBranch()
			if err := r.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(4)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(6)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("chained branch math using branched value", testCase{
		s: `branch({"x":1}, {"x":2}, {"x":3}).x * $this`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			r.MarkAsBranch()
			if err := r.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(4)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(9)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("map on branch", testCase{
		s: `branch([1], [2], [3]).map($this * 2).branch()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			r.MarkAsBranch()
			if err := r.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(4)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewIntValue(6)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
}
