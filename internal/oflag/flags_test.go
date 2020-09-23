package oflag_test

import (
	"github.com/tomwright/dasel/internal/oflag"
	"testing"
)

func TestStringList(t *testing.T) {
	l := oflag.NewStringList()
	if exp, got := "[]", l.String(); exp != got {
		t.Errorf("expected %s, got %s", exp, got)
		return
	}

	if err := l.Set("a"); err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if exp, got := "[a]", l.String(); exp != got {
		t.Errorf("expected %s, got %s", exp, got)
		return
	}

	if err := l.Set("b"); err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if exp, got := "[a b]", l.String(); exp != got {
		t.Errorf("expected %s, got %s", exp, got)
		return
	}
}
