package lexer

import "testing"

func TestTokenizer_Parse(t *testing.T) {
	tok := NewTokenizer("foo.bar.baz[1] != 42.123 || foo.bar.baz['hello'] == 42 && x == 'a\\'b' + false")
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
		Add, Bool,
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
