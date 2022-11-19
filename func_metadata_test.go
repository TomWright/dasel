package dasel

import (
	"reflect"
	"testing"
)

func TestMetadataFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"metadata()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "metadata",
			Args:     []string{},
		}),
	)

	t.Run("Type", func(t *testing.T) {
		orig := []interface{}{
			"abc", true, false, 1, 1.1, []interface{}{1},
		}
		ctx := newSelectContext(&orig, "all().metadata(type)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		exp := []interface{}{
			"string", "bool", "bool", "int", "float64", "slice",
		}
		got := s.Interfaces()

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})
}
