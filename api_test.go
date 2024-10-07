package dasel_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3"
	"testing"
)

type modifyTestCase struct {
	selector string
	in       any
	value    any
	exp      any
}

func (tc modifyTestCase) run(t *testing.T) {
	if err := dasel.Modify(&tc.in, "[1]", 4); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cmp.Equal(tc.exp, tc.in) {
		t.Errorf("unexpected result: %s", cmp.Diff(tc.exp, tc.in))
	}
}

func TestModify(t *testing.T) {
	t.Run("int over int", modifyTestCase{
		selector: "[1]",
		in:       []int{1, 2, 3},
		value:    4,
		exp:      []int{1, 4, 3},
	}.run)
}
