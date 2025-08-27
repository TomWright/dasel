package lexer_test

import (
	"errors"
	"testing"

	"github.com/tomwright/dasel/v3/selector/lexer"
)

type testCase struct {
	in  string
	out []lexer.TokenKind
}

func (tc testCase) run(t *testing.T) {
	tok := lexer.NewTokenizer(tc.in)
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

type errTestCase struct {
	in    string
	match func(error) bool
}

func (tc errTestCase) run(t *testing.T) {
	tok := lexer.NewTokenizer(tc.in)
	tokens, err := tok.Tokenize()
	if !tc.match(err) {
		t.Errorf("unexpected error, got %v", err)
	}
	if tokens != nil {
		t.Errorf("unexpected tokens: %v", tokens)
	}
}

// nolint:unused
func matchUnexpectedError(r rune, p int) func(error) bool {
	return func(err error) bool {
		var e *lexer.UnexpectedTokenError
		if !errors.As(err, &e) {
			return false
		}

		return e.Token == r && e.Pos == p
	}
}

func matchUnexpectedEOFError(p int) func(error) bool {
	return func(err error) bool {
		var e *lexer.UnexpectedEOFError
		if !errors.As(err, &e) {
			return false
		}

		return e.Pos == p
	}
}

func TestTokenizer_Parse(t *testing.T) {
	t.Run("variables", testCase{
		in: "$foo $bar123 $baz $ $quietLoudSCREAM $_hello $hello_world",
		out: []lexer.TokenKind{
			lexer.Variable,
			lexer.Variable,
			lexer.Variable,
			lexer.Dollar,
			lexer.Variable,
			lexer.Variable,
			lexer.Variable,
		},
	}.run)

	t.Run("if", testCase{
		in: `if elseif else`,
		out: []lexer.TokenKind{
			lexer.If,
			lexer.ElseIf,
			lexer.Else,
		},
	}.run)

	t.Run("regex", testCase{
		in: `r/asd/ r/hello there/`,
		out: []lexer.TokenKind{
			lexer.RegexPattern,
			lexer.RegexPattern,
		},
	}.run)

	t.Run("sort by", testCase{
		in: `sortBy(foo, asc)`,
		out: []lexer.TokenKind{
			lexer.SortBy,
			lexer.OpenParen,
			lexer.Symbol,
			lexer.Comma,
			lexer.Asc,
			lexer.CloseParen,
		},
	}.run)

	t.Run("recursive descent", func(t *testing.T) {
		t.Run("key", testCase{
			in: `..foo`,
			out: []lexer.TokenKind{
				lexer.RecursiveDescent,
				lexer.Symbol,
			},
		}.run)
		t.Run("index", testCase{
			in: `..[1]`,
			out: []lexer.TokenKind{
				lexer.RecursiveDescent,
				lexer.OpenBracket,
				lexer.Number,
				lexer.CloseBracket,
			},
		}.run)
		t.Run("wildcard", testCase{
			in: `..*`,
			out: []lexer.TokenKind{
				lexer.RecursiveDescent,
				lexer.Star,
			},
		}.run)
	})

	t.Run("everything", testCase{
		in: "foo.bar.baz[1] != 42.123 || foo.b_a_r.baz['hello'] == 42 && x == 'a\\'b' + false true . .... asd... $name null",
		out: []lexer.TokenKind{
			lexer.Symbol, lexer.Dot, lexer.Symbol, lexer.Dot, lexer.Symbol, lexer.OpenBracket, lexer.Number, lexer.CloseBracket, lexer.NotEqual, lexer.Number,
			lexer.Or,
			lexer.Symbol, lexer.Dot, lexer.Symbol, lexer.Dot, lexer.Symbol, lexer.OpenBracket, lexer.String, lexer.CloseBracket, lexer.Equal, lexer.Number,
			lexer.And,
			lexer.Symbol, lexer.Equal, lexer.String,
			lexer.Plus, lexer.Bool, lexer.Bool,
			lexer.Dot, lexer.Spread, lexer.Dot,
			lexer.Symbol, lexer.Spread,
			lexer.Variable, lexer.Null,
		},
	}.run)

	t.Run("unhappy", func(t *testing.T) {
		t.Run("unfinished double quote", errTestCase{
			in:    `"hello`,
			match: matchUnexpectedEOFError(6),
		}.run)
		t.Run("unfinished single quote", errTestCase{
			in:    `'hello`,
			match: matchUnexpectedEOFError(6),
		}.run)
	})
}
