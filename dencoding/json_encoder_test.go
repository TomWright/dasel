package dencoding_test

import (
	"bytes"
	"github.com/tomwright/dasel/v2/dencoding"
	"testing"
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
