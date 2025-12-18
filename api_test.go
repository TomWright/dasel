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
