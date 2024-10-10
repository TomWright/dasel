package dencoding_test

import (
	"bytes"
	"testing"

	"github.com/tomwright/dasel/v3/dencoding"
)

func TestJSONEncoder_Encode(t *testing.T) {
	orig := dencoding.NewMap().
		Set("c", "x").
		Set("b", "y").
		Set("a", "z")

	exp := `{
  "c": "x",
  "b": "y",
  "a": "z"
}
`

	gotBuffer := new(bytes.Buffer)

	encoder := dencoding.NewJSONEncoder(gotBuffer, dencoding.JSONEncodeIndent("", "  "))
	if err := encoder.Encode(orig); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	got := gotBuffer.String()

	if exp != got {
		t.Errorf("expected %s, got %s", exp, got)
	}
}
