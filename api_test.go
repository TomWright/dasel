package dasel_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3"
)

type modifyTestCase struct {
	selector string
	in       any
	value    any
	exp      any
	count    int
}

func (tc modifyTestCase) run(t *testing.T) {
	count, err := dasel.Modify(t.Context(), &tc.in, tc.selector, tc.value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != tc.count {
		t.Errorf("unexpected count: %d", count)
	}
	if !cmp.Equal(tc.exp, tc.in) {
		t.Errorf("unexpected result: %s", cmp.Diff(tc.exp, tc.in))
	}
}

func TestQuery(t *testing.T) {
	t.Run("basic query", func(t *testing.T) {
		inputData := map[string]any{
			"users": []map[string]any{
				{"name": "Alice", "age": 30},
				{"name": "Bob", "age": 25},
			},
		}
		results, count, err := dasel.Query(t.Context(), inputData, "users.map(name)...")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 2 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []string{"Alice", "Bob"}
		for i, r := range results {
			strVal, err := r.StringValue()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strVal != exp[i] {
				t.Errorf("unexpected result at index %d: %s", i, strVal)
			}
		}
	})
}

func TestSelect(t *testing.T) {
	t.Run("basic select", func(t *testing.T) {
		inputData := map[string]any{
			"users": []map[string]any{
				{"name": "Alice", "age": 30},
				{"name": "Bob", "age": 25},
			},
		}
		result, count, err := dasel.Select(t.Context(), inputData, "users.map(name)...")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 2 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"Alice", "Bob"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
}

func TestTernary(t *testing.T) {
	t.Run("true literal", func(t *testing.T) {
		result, count, err := dasel.Select(t.Context(), nil, `true ? "yes" : "no"`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"yes"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("false literal", func(t *testing.T) {
		result, count, err := dasel.Select(t.Context(), nil, `false ? "yes" : "no"`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"no"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("condition with comparison", func(t *testing.T) {
		inputData := map[string]any{"score": 85}
		result, count, err := dasel.Select(t.Context(), inputData, `score >= 70 ? "pass" : "fail"`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"pass"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("condition with comparison false", func(t *testing.T) {
		inputData := map[string]any{"score": 50}
		result, count, err := dasel.Select(t.Context(), inputData, `score >= 70 ? "pass" : "fail"`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"fail"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("nested ternary", func(t *testing.T) {
		result, count, err := dasel.Select(t.Context(), nil, `true ? (false ? "a" : "b") : "c"`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"b"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("ternary returns number", func(t *testing.T) {
		result, count, err := dasel.Select(t.Context(), nil, `true ? 42 : 0`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{int64(42)}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("ternary with arithmetic", func(t *testing.T) {
		result, count, err := dasel.Select(t.Context(), nil, `true ? 1 + 2 : 3 + 4`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{int64(3)}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("ternary with property lookup", func(t *testing.T) {
		inputData := map[string]any{"active": true, "name": "Alice"}
		result, count, err := dasel.Select(t.Context(), inputData, `active ? name : "unknown"`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"Alice"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("ternary in map", func(t *testing.T) {
		inputData := map[string]any{
			"items": []any{
				map[string]any{"val": 10},
				map[string]any{"val": 3},
				map[string]any{"val": 7},
			},
		}
		result, count, err := dasel.Select(t.Context(), inputData, `items.map(val >= 5 ? "big" : "small")...`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 3 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"big", "small", "big"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
	t.Run("ternary with logical operators", func(t *testing.T) {
		inputData := map[string]any{"a": 5, "b": 10}
		result, count, err := dasel.Select(t.Context(), inputData, `a > 1 && b > 5 ? "both" : "nope"`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Errorf("unexpected count: %d", count)
		}
		exp := []any{"both"}
		if !cmp.Equal(exp, result) {
			t.Errorf("unexpected result: %s", cmp.Diff(exp, result))
		}
	})
}

func TestModify(t *testing.T) {
	t.Run("index", func(t *testing.T) {
		t.Run("int over int", modifyTestCase{
			selector: "$this[1]",
			in:       []int{1, 2, 3},
			value:    4,
			exp:      []int{1, 4, 3},
			count:    1,
		}.run)
		t.Run("string over int", modifyTestCase{
			selector: "$this[1]",
			in:       []any{1, 2, 3},
			value:    "4",
			exp:      []any{1, "4", 3},
			count:    1,
		}.run)
	})
}
