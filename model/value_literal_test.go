package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestValue_IsNull(t *testing.T) {
	v := model.NewNullValue()
	if !v.IsNull() {
		t.Fatalf("expected value to be null")
	}
}
