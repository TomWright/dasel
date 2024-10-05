package lexer

import (
	"strings"
	"unicode"

	"github.com/tomwright/dasel/v3/internal/ptr"
)

type Tokenizer struct {
	i      int
	src    string
	srcLen int
}

func NewTokenizer(src string) *Tokenizer {
	return &Tokenizer{
		i:      0,
		src:    src,
		srcLen: len([]rune(src)),
	}
}

func (p *Tokenizer) Tokenize() (Tokens, error) {
	var tokens Tokens
	for {
		tok, err := p.Next()
		if err != nil {
			return nil, err
		}
		if tok.Kind == EOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens, nil
}

func (p *Tokenizer) peekRuneEqual(i int, to rune) bool {
	if i >= p.srcLen {
		return false
	}
	return rune(p.src[i]) == to
}

func (p *Tokenizer) peekRuneMatches(i int, fn func(rune) bool) bool {
	if i >= p.srcLen {
		return false
	}
	return fn(rune(p.src[i]))
}

func (p *Tokenizer) parseCurRune() (Token, error) {
	switch p.src[p.i] {
	case '.':
		if p.peekRuneEqual(p.i+1, '.') && p.peekRuneEqual(p.i+2, '.') {
			return NewToken(Spread, "...", p.i, 3), nil
		}
		return NewToken(Dot, ".", p.i, 1), nil
	case ',':
		return NewToken(Comma, ",", p.i, 1), nil
	case ':':
		return NewToken(Colon, ":", p.i, 1), nil
	case '[':
		return NewToken(OpenBracket, "[", p.i, 1), nil
	case ']':
		return NewToken(CloseBracket, "]", p.i, 1), nil
	case '(':
		return NewToken(OpenParen, "(", p.i, 1), nil
	case ')':
		return NewToken(CloseParen, ")", p.i, 1), nil
	case '{':
		return NewToken(OpenCurly, "{", p.i, 1), nil
	case '}':
		return NewToken(CloseCurly, "}", p.i, 1), nil
	case '*':
		return NewToken(Star, "*", p.i, 1), nil
	case '/':
		return NewToken(Slash, "/", p.i, 1), nil
	case '%':
		return NewToken(Percent, "%", p.i, 1), nil
	case '$':
		if p.peekRuneMatches(p.i+1, unicode.IsLetter) {
			pos := p.i + 1
			for pos < p.srcLen && (unicode.IsLetter(rune(p.src[pos])) || unicode.IsDigit(rune(p.src[pos]))) {
				pos++
			}
			return NewToken(Variable, p.src[p.i+1:pos], p.i, pos-p.i), nil
		}
		return NewToken(Dollar, "$", p.i, 1), nil
	case '=':
		if p.peekRuneEqual(p.i+1, '=') {
			return NewToken(Equal, "==", p.i, 2), nil
		}
		if p.peekRuneEqual(p.i+1, '~') {
			return NewToken(Like, "=~", p.i, 2), nil
		}
		return NewToken(Equals, "=", p.i, 1), nil
	case '+':
		if p.peekRuneEqual(p.i+1, '=') {
			return NewToken(IncrementBy, "+=", p.i, 2), nil
		}
		if p.peekRuneEqual(p.i+1, '+') {
			return NewToken(Increment, "++", p.i, 2), nil
		}
		return NewToken(Plus, "+", p.i, 1), nil
	case '-':
		if p.peekRuneEqual(p.i+1, '=') {
			return NewToken(DecrementBy, "-=", p.i, 2), nil
		}
		if p.peekRuneEqual(p.i+1, '-') {
			return NewToken(Decrement, "--", p.i, 2), nil
		}
		return NewToken(Dash, "-", p.i, 1), nil
	case '>':
		if p.peekRuneEqual(p.i+1, '=') {
			return NewToken(GreaterThanOrEqual, ">=", p.i, 2), nil
		}
		return NewToken(GreaterThan, ">", p.i, 1), nil
	case '<':
		if p.peekRuneEqual(p.i+1, '=') {
			return NewToken(LessThanOrEqual, "<>>=", p.i, 2), nil
		}
		return NewToken(LessThan, "<", p.i, 1), nil
	case '!':
		if p.peekRuneEqual(p.i+1, '=') {
			return NewToken(NotEqual, "!=", p.i, 2), nil
		}
		if p.peekRuneEqual(p.i+1, '~') {
			return NewToken(NotLike, "!~", p.i, 2), nil
		}
		return NewToken(Exclamation, "!", p.i, 1), nil
	case '&':
		if p.peekRuneEqual(p.i+1, '&') {
			return NewToken(And, "&&", p.i, 2), nil
		}
		return Token{}, &UnexpectedTokenError{
			Pos:   p.i,
			Token: rune(p.src[p.i]),
		}
	case '|':
		if p.peekRuneEqual(p.i+1, '|') {
			return NewToken(Or, "||", p.i, 2), nil
		}
		return Token{}, &UnexpectedTokenError{
			Pos:   p.i,
			Token: rune(p.src[p.i]),
		}
	case '"', '\'':
		pos := p.i
		buf := make([]rune, 0)
		pos++
		var escaped bool
		for pos < p.srcLen {
			if p.src[pos] == p.src[p.i] && !escaped {
				break
			}
			if escaped {
				escaped = false
				buf = append(buf, rune(p.src[pos]))
				pos++
				continue
			}
			if p.src[pos] == '\\' {
				pos++
				escaped = true
				continue
			}
			buf = append(buf, rune(p.src[pos]))
			pos++
		}
		res := NewToken(String, string(buf), p.i, pos+1-p.i)
		return res, nil
	default:
		pos := p.i

		matchStr := func(pos int, m string, caseInsensitive bool, kind TokenKind) *Token {
			l := len(m)
			if pos+(l-1) >= p.srcLen {
				return nil
			}
			other := p.src[pos : pos+l]
			if m == other || caseInsensitive && strings.EqualFold(m, other) {
				return ptr.To(NewToken(kind, other, pos, l))
			}
			return nil
		}

		if t := matchStr(pos, "null", true, Null); t != nil {
			return *t, nil
		}
		if t := matchStr(pos, "true", true, Bool); t != nil {
			return *t, nil
		}
		if t := matchStr(pos, "false", true, Bool); t != nil {
			return *t, nil
		}
		if t := matchStr(pos, "elseif", false, ElseIf); t != nil {
			return *t, nil
		}
		if t := matchStr(pos, "if", false, If); t != nil {
			return *t, nil
		}
		if t := matchStr(pos, "else", false, Else); t != nil {
			return *t, nil
		}

		if unicode.IsDigit(rune(p.src[pos])) {
			// Handle whole numbers
			for pos < p.srcLen && unicode.IsDigit(rune(p.src[pos])) {
				pos++
			}
			// Handle floats
			if pos < p.srcLen && p.src[pos] == '.' && pos+1 < p.srcLen && unicode.IsDigit(rune(p.src[pos+1])) {
				pos++
				for pos < p.srcLen && unicode.IsDigit(rune(p.src[pos])) {
					pos++
				}
			}
			return NewToken(Number, p.src[p.i:pos], p.i, pos-p.i), nil
		}

		if unicode.IsLetter(rune(p.src[pos])) {
			for pos < p.srcLen && (unicode.IsLetter(rune(p.src[pos])) || unicode.IsDigit(rune(p.src[pos]))) {
				pos++
			}
			return NewToken(Symbol, p.src[p.i:pos], p.i, pos-p.i), nil
		}

		return Token{}, &UnexpectedTokenError{
			Pos:   p.i,
			Token: rune(p.src[p.i]),
		}
	}
}

func (p *Tokenizer) Next() (Token, error) {
	if p.i >= len(p.src) {
		return NewToken(EOF, "", p.i, 0), nil
	}

	// Skip over whitespace
	for p.i < p.srcLen && unicode.IsSpace(rune(p.src[p.i])) {
		p.i++
	}

	t, err := p.parseCurRune()
	if err != nil {
		return Token{}, err
	}
	p.i += t.Len
	return t, nil
}
