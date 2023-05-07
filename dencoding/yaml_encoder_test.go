package dencoding_test

import (
	"bytes"
	"github.com/tomwright/dasel/v2/dencoding"
	"testing"
)

func TestYAMLEncoder_Encode(t *testing.T) {
	orig := dencoding.NewMap().
		Set("c", "x").
		Set("b", "y").
		Set("a", []any{"a", "c", "b"})

	exp := `c: x
b: y
a:
  - a
  - c
  - b
`

	gotBuffer := new(bytes.Buffer)

	encoder := dencoding.NewYAMLEncoder(gotBuffer, dencoding.YAMLEncodeIndent(2))
	if err := encoder.Encode(orig); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	got := gotBuffer.String()

	if exp != got {
		t.Errorf("expected %s, got %s", exp, got)
	}
}
