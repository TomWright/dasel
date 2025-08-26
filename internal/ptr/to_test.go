package ptr_test

import (
	"github.com/tomwright/dasel/v3/internal/ptr"
	"testing"
)

func TestTo(t *testing.T) {
	a := 1
	if exp, got := a, *(ptr.To(a)); exp != a {
		t.Errorf("expected %d, got %d", exp, got)
	}
}
