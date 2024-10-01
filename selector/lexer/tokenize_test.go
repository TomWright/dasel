package lexer

import "testing"

func TestTokenizer_Parse(t *testing.T) {
	type testCase struct {
		in  string
		out []TokenKind
	}

	runTest := func(tc testCase) func(t *testing.T) {
		return func(t *testing.T) {
			tok := NewTokenizer(tc.in)
			tokens, err := tok.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tc.out) {
				t.Fatalf("unexpected number of tokens: %d", len(tokens))
			}
			for i := range tokens {
				if tokens[i].Kind != tc.out[i] {
					t.Errorf("unexpected token kind at position %d: exp %v, got %v", i, tc.out[i], tokens[i].Kind)
					return
				}
			}
		}
	}

	t.Run("variables", runTest(testCase{
		in: "$foo $bar123 $baz $",
		out: []TokenKind{
			Variable,
			Variable,
			Variable,
			Dollar,
		},
	}))

	t.Run("everything", runTest(testCase{
		in: "foo.bar.baz[1] != 42.123 || foo.bar.baz['hello'] == 42 && x == 'a\\'b' + false . .... asd... $name",
		out: []TokenKind{
			Symbol, Dot, Symbol, Dot, Symbol, OpenBracket, Number, CloseBracket, NotEqual, Number,
			Or,
			Symbol, Dot, Symbol, Dot, Symbol, OpenBracket, String, CloseBracket, Equal, Number,
			And,
			Symbol, Equal, String,
			Plus, Bool,
			Dot, Spread, Dot,
			Symbol, Spread,
			Variable,
		},
	}))

	tok := NewTokenizer("foo.bar.baz[1] != 42.123 || foo.bar.baz['hello'] == 42 && x == 'a\\'b' + false . .... asd... $name")
	tokens, err := tok.Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exp := []TokenKind{
		Symbol, Dot, Symbol, Dot, Symbol, OpenBracket, Number, CloseBracket, NotEqual, Number,
		Or,
		Symbol, Dot, Symbol, Dot, Symbol, OpenBracket, String, CloseBracket, Equal, Number,
		And,
		Symbol, Equal, String,
		Plus, Bool,
		Dot, Spread, Dot,
		Symbol, Spread,
		Variable,
	}
	if len(tokens) != len(exp) {
		t.Fatalf("unexpected number of tokens: %d", len(tokens))
	}

	for i := range tokens {
		if tokens[i].Kind != exp[i] {
			t.Errorf("unexpected token kind at position %d: exp %v, got %v", i, exp[i], tokens[i].Kind)
			return
		}
	}
}
