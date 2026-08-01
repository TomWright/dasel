package plain_test

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/plain"
)

func newTestWriter(t *testing.T) parsing.Writer {
	t.Helper()
	w, err := plain.Plain.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error creating writer: %v", err)
	}
	return w
}

func TestPlainWriter_Scalars(t *testing.T) {
	tests := map[string]struct {
		in  *model.Value
		exp string
	}{
		"string":            {in: model.NewStringValue("world"), exp: "world\n"},
		"string with space": {in: model.NewStringValue("hello there"), exp: "hello there\n"},
		"string with quote": {in: model.NewStringValue(`say "hi"`), exp: "say \"hi\"\n"},
		"string that looks quoted": {
			in: model.NewStringValue(`'world'`), exp: "'world'\n",
		},
		"empty string":  {in: model.NewStringValue(""), exp: "\n"},
		"int":           {in: model.NewIntValue(123), exp: "123\n"},
		"negative int":  {in: model.NewIntValue(-5), exp: "-5\n"},
		"zero":          {in: model.NewIntValue(0), exp: "0\n"},
		"float":         {in: model.NewFloatValue(12.3), exp: "12.3\n"},
		"float no frac": {in: model.NewFloatValue(1000), exp: "1000\n"},
		"bool true":     {in: model.NewBoolValue(true), exp: "true\n"},
		"bool false":    {in: model.NewBoolValue(false), exp: "false\n"},
		// Null has no text form, so it writes an empty line rather than the
		// literal "null", which would be ambiguous with the string "null".
		"null": {in: model.NewNullValue(), exp: "\n"},
	}

	w := newTestWriter(t)

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			out, err := w.Write(tc.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(out); got != tc.exp {
				t.Fatalf("expected %q, got %q", tc.exp, got)
			}
		})
	}
}

func TestPlainWriter_RejectsNonScalars(t *testing.T) {
	w := newTestWriter(t)

	t.Run("map", func(t *testing.T) {
		v := model.NewMapValue()
		if err := v.SetMapKey("a", model.NewIntValue(1)); err != nil {
			t.Fatal(err)
		}
		_, err := w.Write(v)
		if err == nil {
			t.Fatal("expected an error writing a map, got none")
		}
		if !strings.Contains(err.Error(), "map") {
			t.Fatalf("expected the error to name the type, got %v", err)
		}
	})

	t.Run("slice", func(t *testing.T) {
		v := model.NewSliceValue()
		if err := v.Append(model.NewIntValue(1)); err != nil {
			t.Fatal(err)
		}
		_, err := w.Write(v)
		if err == nil {
			t.Fatal("expected an error writing a slice, got none")
		}
		if !strings.Contains(err.Error(), "array") {
			t.Fatalf("expected the error to name the type, got %v", err)
		}
	})

	t.Run("empty map", func(t *testing.T) {
		if _, err := w.Write(model.NewMapValue()); err == nil {
			t.Fatal("expected an error writing an empty map, got none")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		if _, err := w.Write(model.NewSliceValue()); err == nil {
			t.Fatal("expected an error writing an empty slice, got none")
		}
	})
}

// Plain output carries no structure, so there is nothing to read back.
func TestPlain_HasNoReader(t *testing.T) {
	if _, err := plain.Plain.NewReader(parsing.DefaultReaderOptions()); err == nil {
		t.Fatal("expected plain to have no registered reader")
	}
}
