package internal

import (
	"testing"
)

func collectTokens(t *testing.T, input string) []Token {
	t.Helper()
	tok := NewTokenizer(input)
	var tokens []Token
	for {
		token, err := tok.NextToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		tokens = append(tokens, token)
		if token.Type == TokenEOF {
			break
		}
	}
	return tokens
}

func TestTokenizer_BareIdentifier(t *testing.T) {
	tokens := collectTokens(t, "node")
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != TokenIdentifier || tokens[0].Value != "node" {
		t.Errorf("expected identifier 'node', got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_QuotedString(t *testing.T) {
	tokens := collectTokens(t, `"hello world"`)
	if tokens[0].Type != TokenQuotedString || tokens[0].Value != "hello world" {
		t.Errorf("expected quoted string 'hello world', got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_StringEscapes(t *testing.T) {
	tokens := collectTokens(t, `"hello\nworld\t!"`)
	if tokens[0].Value != "hello\nworld\t!" {
		t.Errorf("expected escaped string, got %q", tokens[0].Value)
	}
}

func TestTokenizer_UnicodeEscape(t *testing.T) {
	tokens := collectTokens(t, `"\u{1F600}"`)
	if tokens[0].Value != "\U0001F600" {
		t.Errorf("expected unicode char, got %q", tokens[0].Value)
	}
}

func TestTokenizer_Integer(t *testing.T) {
	tokens := collectTokens(t, "42")
	if tokens[0].Type != TokenInteger || tokens[0].Value != "42" {
		t.Errorf("expected integer 42, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_NegativeInteger(t *testing.T) {
	tokens := collectTokens(t, "-7")
	if tokens[0].Type != TokenInteger || tokens[0].Value != "-7" {
		t.Errorf("expected integer -7, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_Float(t *testing.T) {
	tokens := collectTokens(t, "3.14")
	if tokens[0].Type != TokenFloat || tokens[0].Value != "3.14" {
		t.Errorf("expected float 3.14, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_FloatExponent(t *testing.T) {
	tokens := collectTokens(t, "1.5e10")
	if tokens[0].Type != TokenFloat || tokens[0].Value != "1.5e10" {
		t.Errorf("expected float 1.5e10, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_HexInt(t *testing.T) {
	tokens := collectTokens(t, "0xff")
	if tokens[0].Type != TokenHexInt || tokens[0].Value != "0xff" {
		t.Errorf("expected hex int, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_OctalInt(t *testing.T) {
	tokens := collectTokens(t, "0o77")
	if tokens[0].Type != TokenOctalInt || tokens[0].Value != "0o77" {
		t.Errorf("expected octal int, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_BinaryInt(t *testing.T) {
	tokens := collectTokens(t, "0b1010")
	if tokens[0].Type != TokenBinaryInt || tokens[0].Value != "0b1010" {
		t.Errorf("expected binary int, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_UnderscoreInNumbers(t *testing.T) {
	tokens := collectTokens(t, "1_000_000")
	if tokens[0].Type != TokenInteger || tokens[0].Value != "1000000" {
		t.Errorf("expected integer 1000000, got %v %q", tokens[0].Type, tokens[0].Value)
	}
}

func TestTokenizer_V1True(t *testing.T) {
	tok := NewTokenizer("true")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenTrue || token.Value != "true" {
		t.Errorf("expected v1 true, got %v %q", token.Type, token.Value)
	}
	if tok.Version != Version1 {
		t.Errorf("expected version 1, got %v", tok.Version)
	}
}

func TestTokenizer_V1False(t *testing.T) {
	tok := NewTokenizer("false")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenFalse || token.Value != "false" {
		t.Errorf("expected v1 false, got %v %q", token.Type, token.Value)
	}
}

func TestTokenizer_V1Null(t *testing.T) {
	tok := NewTokenizer("null")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenNull || token.Value != "null" {
		t.Errorf("expected v1 null, got %v %q", token.Type, token.Value)
	}
}

func TestTokenizer_V2True(t *testing.T) {
	tok := NewTokenizer("#true")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenTrue || token.Value != "#true" {
		t.Errorf("expected v2 true, got %v %q", token.Type, token.Value)
	}
	if tok.Version != Version2 {
		t.Errorf("expected version 2, got %v", tok.Version)
	}
}

func TestTokenizer_V2False(t *testing.T) {
	tok := NewTokenizer("#false")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenFalse || token.Value != "#false" {
		t.Errorf("expected v2 false, got %v %q", token.Type, token.Value)
	}
}

func TestTokenizer_V2Null(t *testing.T) {
	tok := NewTokenizer("#null")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenNull || token.Value != "#null" {
		t.Errorf("expected v2 null, got %v %q", token.Type, token.Value)
	}
}

func TestTokenizer_V2Inf(t *testing.T) {
	tok := NewTokenizer("#inf")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenInf {
		t.Errorf("expected inf, got %v", token.Type)
	}
}

func TestTokenizer_V2NegInf(t *testing.T) {
	tok := NewTokenizer("#-inf")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenNegInf {
		t.Errorf("expected -inf, got %v", token.Type)
	}
}

func TestTokenizer_V2NaN(t *testing.T) {
	tok := NewTokenizer("#nan")
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenNaN {
		t.Errorf("expected nan, got %v", token.Type)
	}
}

func TestTokenizer_V1RawString(t *testing.T) {
	tok := NewTokenizer(`r"hello world"`)
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenRawString || token.Value != "hello world" {
		t.Errorf("expected raw string 'hello world', got %v %q", token.Type, token.Value)
	}
	if tok.Version != Version1 {
		t.Errorf("expected version 1, got %v", tok.Version)
	}
}

func TestTokenizer_V1RawStringWithHashes(t *testing.T) {
	tok := NewTokenizer(`r#"hello "world""#`)
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenRawString || token.Value != `hello "world"` {
		t.Errorf("expected raw string, got %v %q", token.Type, token.Value)
	}
}

func TestTokenizer_V2RawString(t *testing.T) {
	tok := NewTokenizer(`#"hello world"#`)
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenRawString || token.Value != "hello world" {
		t.Errorf("expected raw string 'hello world', got %v %q", token.Type, token.Value)
	}
	if tok.Version != Version2 {
		t.Errorf("expected version 2, got %v", tok.Version)
	}
}

func TestTokenizer_V2RawStringDoubleHash(t *testing.T) {
	tok := NewTokenizer(`##"contains #"quotes"# inside"##`)
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenRawString || token.Value != `contains #"quotes"# inside` {
		t.Errorf("expected raw string, got %v %q", token.Type, token.Value)
	}
}

func TestTokenizer_MultiLineString(t *testing.T) {
	input := "\"\"\"\n    hello\n    world\n    \"\"\""
	tok := NewTokenizer(input)
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Type != TokenMultiLineString {
		t.Errorf("expected multi-line string, got %v", token.Type)
	}
	if token.Value != "hello\nworld" {
		t.Errorf("expected dedented content, got %q", token.Value)
	}
}

func TestTokenizer_Braces(t *testing.T) {
	tokens := collectTokens(t, "{ }")
	if tokens[0].Type != TokenOpenBrace {
		t.Errorf("expected open brace, got %v", tokens[0].Type)
	}
	if tokens[1].Type != TokenCloseBrace {
		t.Errorf("expected close brace, got %v", tokens[1].Type)
	}
}

func TestTokenizer_Equals(t *testing.T) {
	tokens := collectTokens(t, "key=val")
	if tokens[1].Type != TokenEquals {
		t.Errorf("expected equals, got %v", tokens[1].Type)
	}
}

func TestTokenizer_SlashDash(t *testing.T) {
	tokens := collectTokens(t, "/-")
	if tokens[0].Type != TokenSlashDash {
		t.Errorf("expected slash-dash, got %v", tokens[0].Type)
	}
}

func TestTokenizer_LineComment(t *testing.T) {
	tokens := collectTokens(t, "node // comment\nother")
	types := tokenTypes(tokens)
	expected := []TokenType{TokenIdentifier, TokenNewline, TokenIdentifier, TokenEOF}
	assertTokenTypes(t, types, expected)
}

func TestTokenizer_BlockComment(t *testing.T) {
	tokens := collectTokens(t, "node /* comment */ other")
	if tokens[0].Type != TokenIdentifier || tokens[0].Value != "node" {
		t.Errorf("expected 'node', got %v %q", tokens[0].Type, tokens[0].Value)
	}
	if tokens[1].Type != TokenIdentifier || tokens[1].Value != "other" {
		t.Errorf("expected 'other', got %v %q", tokens[1].Type, tokens[1].Value)
	}
}

func TestTokenizer_NestedBlockComment(t *testing.T) {
	tokens := collectTokens(t, "node /* outer /* inner */ outer */ other")
	if tokens[0].Value != "node" {
		t.Errorf("expected 'node', got %q", tokens[0].Value)
	}
	if tokens[1].Value != "other" {
		t.Errorf("expected 'other', got %q", tokens[1].Value)
	}
}

func TestTokenizer_LineContinuation(t *testing.T) {
	tokens := collectTokens(t, "node \\\n  arg")
	if tokens[0].Type != TokenIdentifier || tokens[0].Value != "node" {
		t.Errorf("expected 'node', got %v %q", tokens[0].Type, tokens[0].Value)
	}
	if tokens[1].Type != TokenIdentifier || tokens[1].Value != "arg" {
		t.Errorf("expected 'arg', got %v %q", tokens[1].Type, tokens[1].Value)
	}
}

func TestTokenizer_Semicolon(t *testing.T) {
	tokens := collectTokens(t, "a; b")
	types := tokenTypes(tokens)
	expected := []TokenType{TokenIdentifier, TokenSemicolon, TokenIdentifier, TokenEOF}
	assertTokenTypes(t, types, expected)
}

func TestTokenizer_TypeAnnotation(t *testing.T) {
	tokens := collectTokens(t, "(u8)123")
	types := tokenTypes(tokens)
	expected := []TokenType{TokenOpenParen, TokenIdentifier, TokenCloseParen, TokenInteger, TokenEOF}
	assertTokenTypes(t, types, expected)
}

func TestTokenizer_EscapeS(t *testing.T) {
	tok := NewTokenizer(`"\s"`)
	token, err := tok.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Value != " " {
		t.Errorf(`expected space from \s, got %q`, token.Value)
	}
	if tok.Version != Version2 {
		t.Errorf("expected version 2, got %v", tok.Version)
	}
}

func tokenTypes(tokens []Token) []TokenType {
	types := make([]TokenType, len(tokens))
	for i, t := range tokens {
		types[i] = t.Type
	}
	return types
}

func assertTokenTypes(t *testing.T, got, expected []TokenType) {
	t.Helper()
	if len(got) != len(expected) {
		t.Errorf("expected %d tokens, got %d", len(expected), len(got))
		return
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("token %d: expected %v, got %v", i, expected[i], got[i])
		}
	}
}
