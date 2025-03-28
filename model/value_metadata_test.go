package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestValue_IsBranch(t *testing.T) {
	val := model.NewNullValue()
	if exp, got := false, val.IsBranch(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	val.MarkAsBranch()
	if exp, got := true, val.IsBranch(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestValue_IsSpread(t *testing.T) {
	val := model.NewNullValue()
	if exp, got := false, val.IsSpread(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	val.MarkAsSpread()
	if exp, got := true, val.IsSpread(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
