package internal

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Version represents the KDL specification version.
type Version int

const (
	VersionUnknown Version = 0
	Version1       Version = 1
	Version2       Version = 2
)

// TokenType represents the type of a KDL token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNewline
	TokenIdentifier
	TokenQuotedString
	TokenRawString
	TokenMultiLineString
	TokenInteger
	TokenFloat
	TokenHexInt
	TokenOctalInt
	TokenBinaryInt
	TokenTrue
	TokenFalse
	TokenNull
	TokenInf
	TokenNegInf
	TokenNaN
	TokenOpenBrace
	TokenCloseBrace
	TokenOpenParen
	TokenCloseParen
	TokenEquals
	TokenSemicolon
	TokenSlashDash
)

// Token represents a single KDL token.
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Col     int
}

// Tokenizer is a KDL tokenizer supporting both v1 and v2.
type Tokenizer struct {
	input   []rune
	pos     int
	line    int
	col     int
	Version Version
	peeked  *Token
}

// NewTokenizer creates a new tokenizer for the given input.
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input:   []rune(input),
		pos:     0,
		line:    1,
		col:     1,
		Version: VersionUnknown,
	}
}

func (t *Tokenizer) peek() rune {
	if t.pos >= len(t.input) {
		return -1
	}
	return t.input[t.pos]
}

func (t *Tokenizer) peekAt(offset int) rune {
	p := t.pos + offset
	if p >= len(t.input) {
		return -1
	}
	return t.input[p]
}

func (t *Tokenizer) advance() rune {
	if t.pos >= len(t.input) {
		return -1
	}
	ch := t.input[t.pos]
	t.pos++
	if ch == '\n' {
		t.line++
		t.col = 1
	} else {
		t.col++
	}
	return ch
}

func (t *Tokenizer) error(msg string) error {
	return fmt.Errorf("kdl: %s at line %d, col %d", msg, t.line, t.col)
}

// PeekToken returns the next token without consuming it.
func (t *Tokenizer) PeekToken() (Token, error) {
	if t.peeked != nil {
		return *t.peeked, nil
	}
	tok, err := t.NextToken()
	if err != nil {
		return Token{}, err
	}
	t.peeked = &tok
	return tok, nil
}

// NextToken returns the next token from the input.
func (t *Tokenizer) NextToken() (Token, error) {
	if t.peeked != nil {
		tok := *t.peeked
		t.peeked = nil
		return tok, nil
	}
	return t.readToken()
}

func (t *Tokenizer) readToken() (Token, error) {
	t.skipWhitespaceAndComments()

	if t.pos >= len(t.input) {
		return Token{Type: TokenEOF, Line: t.line, Col: t.col}, nil
	}

	ch := t.peek()

	// Newline
	if ch == '\n' || ch == '\r' {
		return t.readNewline()
	}

	// Semicolon (node terminator)
	if ch == ';' {
		tok := Token{Type: TokenSemicolon, Value: ";", Line: t.line, Col: t.col}
		t.advance()
		return tok, nil
	}

	// Braces
	if ch == '{' {
		tok := Token{Type: TokenOpenBrace, Value: "{", Line: t.line, Col: t.col}
		t.advance()
		return tok, nil
	}
	if ch == '}' {
		tok := Token{Type: TokenCloseBrace, Value: "}", Line: t.line, Col: t.col}
		t.advance()
		return tok, nil
	}

	// Parens (type annotations)
	if ch == '(' {
		tok := Token{Type: TokenOpenParen, Value: "(", Line: t.line, Col: t.col}
		t.advance()
		return tok, nil
	}
	if ch == ')' {
		tok := Token{Type: TokenCloseParen, Value: ")", Line: t.line, Col: t.col}
		t.advance()
		return tok, nil
	}

	// Equals
	if ch == '=' {
		tok := Token{Type: TokenEquals, Value: "=", Line: t.line, Col: t.col}
		t.advance()
		return tok, nil
	}

	// Slash-dash
	if ch == '/' && t.peekAt(1) == '-' {
		tok := Token{Type: TokenSlashDash, Value: "/-", Line: t.line, Col: t.col}
		t.advance()
		t.advance()
		return tok, nil
	}

	// Line continuation (backslash followed by newline)
	if ch == '\\' {
		next := t.peekAt(1)
		if next == '\n' || next == '\r' {
			t.advance() // consume backslash
			t.skipNewline()
			t.skipWhitespaceAndComments()
			return t.readToken()
		}
	}

	// v2 keywords starting with #
	if ch == '#' {
		return t.readHash()
	}

	// v1 raw strings: r"..." or r#"..."#
	if ch == 'r' && t.peekAt(1) == '"' {
		return t.readV1RawString()
	}
	if ch == 'r' && t.peekAt(1) == '#' {
		return t.readV1RawString()
	}

	// Quoted strings
	if ch == '"' {
		return t.readQuotedString()
	}

	// Numbers
	if ch == '-' || ch == '+' || (ch >= '0' && ch <= '9') {
		return t.readNumber()
	}

	// Bare identifiers
	if isIdentStart(ch) {
		return t.readIdentifier()
	}

	return Token{}, t.error(fmt.Sprintf("unexpected character: %c", ch))
}

func (t *Tokenizer) readNewline() (Token, error) {
	tok := Token{Type: TokenNewline, Value: "\n", Line: t.line, Col: t.col}
	t.skipNewline()
	return tok, nil
}

func (t *Tokenizer) skipNewline() {
	switch t.peek() {
	case '\r':
		t.advance()
		if t.peek() == '\n' {
			t.advance()
		}
	case '\n':
		t.advance()
	}
}

func (t *Tokenizer) skipWhitespaceAndComments() {
	for t.pos < len(t.input) {
		ch := t.peek()

		// Regular whitespace (not newlines — those are tokens)
		if ch == ' ' || ch == '\t' || ch == '\u00A0' || ch == '\uFEFF' ||
			(ch > 0x2000 && ch <= 0x200A) || ch == 0x202F || ch == 0x205F || ch == 0x3000 {
			t.advance()
			continue
		}

		// Line comments
		if ch == '/' && t.peekAt(1) == '/' {
			t.advance()
			t.advance()
			for t.pos < len(t.input) && t.peek() != '\n' && t.peek() != '\r' {
				t.advance()
			}
			continue
		}

		// Block comments (nestable)
		if ch == '/' && t.peekAt(1) == '*' {
			t.advance()
			t.advance()
			depth := 1
			for t.pos < len(t.input) && depth > 0 {
				if t.peek() == '/' && t.peekAt(1) == '*' {
					t.advance()
					t.advance()
					depth++
				} else if t.peek() == '*' && t.peekAt(1) == '/' {
					t.advance()
					t.advance()
					depth--
				} else {
					t.advance()
				}
			}
			continue
		}

		// Line continuation
		if ch == '\\' {
			next := t.peekAt(1)
			if next == '\n' || next == '\r' {
				t.advance() // consume backslash
				t.skipNewline()
				continue
			}
		}

		break
	}
}

// readHash handles # prefixed tokens in v2 (and auto-detects version).
func (t *Tokenizer) readHash() (Token, error) {
	line, col := t.line, t.col

	// Count leading hashes for raw strings
	hashCount := 0
	p := t.pos
	for p < len(t.input) && t.input[p] == '#' {
		hashCount++
		p++
	}

	// v2 raw string: #"..."# or ##"..."##, etc.
	if p < len(t.input) && t.input[p] == '"' {
		t.setVersion(Version2)
		return t.readV2RawString(hashCount)
	}

	// v2 keywords: #true, #false, #null, #inf, #-inf, #nan
	if hashCount == 1 {
		t.advance() // consume #
		remaining := string(t.input[t.pos:])

		keywords := []struct {
			text    string
			tokType TokenType
		}{
			{"true", TokenTrue},
			{"false", TokenFalse},
			{"null", TokenNull},
			{"inf", TokenInf},
			{"nan", TokenNaN},
			{"-inf", TokenNegInf},
		}

		for _, kw := range keywords {
			if strings.HasPrefix(remaining, kw.text) {
				// Make sure it's not part of a longer identifier
				endPos := t.pos + len(kw.text)
				if endPos >= len(t.input) || !isIdentChar(t.input[endPos]) {
					for range kw.text {
						t.advance()
					}
					t.setVersion(Version2)
					return Token{Type: kw.tokType, Value: "#" + kw.text, Line: line, Col: col}, nil
				}
			}
		}

		return Token{}, t.error("unexpected # character")
	}

	return Token{}, t.error("unexpected # character")
}

func (t *Tokenizer) readV1RawString() (Token, error) {
	line, col := t.line, t.col
	t.advance() // consume 'r'

	if t.peek() == '#' {
		// r#"..."# or r##"..."## etc.
		hashCount := 0
		for t.peek() == '#' {
			hashCount++
			t.advance()
		}
		if t.peek() != '"' {
			return Token{}, t.error("expected '\"' after r and hashes")
		}
		t.advance() // consume opening "

		var sb strings.Builder
		for {
			if t.pos >= len(t.input) {
				return Token{}, t.error("unterminated raw string")
			}
			ch := t.advance()
			if ch == '"' {
				// Check for matching hashes
				matched := 0
				for matched < hashCount && t.peek() == '#' {
					matched++
					t.advance()
				}
				if matched == hashCount {
					t.setVersion(Version1)
					return Token{Type: TokenRawString, Value: sb.String(), Line: line, Col: col}, nil
				}
				sb.WriteRune('"')
				for i := 0; i < matched; i++ {
					sb.WriteRune('#')
				}
			} else {
				sb.WriteRune(ch)
			}
		}
	}

	// r"..."
	if t.peek() != '"' {
		// Not a raw string, it's an identifier starting with 'r'
		t.pos--
		t.col--
		return t.readIdentifier()
	}
	t.advance() // consume opening "

	var sb strings.Builder
	for {
		if t.pos >= len(t.input) {
			return Token{}, t.error("unterminated raw string")
		}
		ch := t.advance()
		if ch == '"' {
			t.setVersion(Version1)
			return Token{Type: TokenRawString, Value: sb.String(), Line: line, Col: col}, nil
		}
		sb.WriteRune(ch)
	}
}

func (t *Tokenizer) readV2RawString(hashCount int) (Token, error) {
	line, col := t.line, t.col

	// Consume the hashes
	for i := 0; i < hashCount; i++ {
		t.advance()
	}

	// Check for multi-line raw string: #"""..."""#
	if t.peek() == '"' && t.peekAt(1) == '"' && t.peekAt(2) == '"' {
		t.advance() // first "
		t.advance() // second "
		t.advance() // third "

		// Must be followed by newline
		if t.peek() == '\r' || t.peek() == '\n' {
			t.skipNewline()
		}

		return t.readMultiLineRawContent(hashCount, line, col)
	}

	// Single-line raw string: #"..."#
	t.advance() // consume opening "

	var sb strings.Builder
	for {
		if t.pos >= len(t.input) {
			return Token{}, t.error("unterminated raw string")
		}
		ch := t.advance()
		if ch == '"' {
			// Check for matching closing hashes
			matched := 0
			for matched < hashCount && t.peek() == '#' {
				matched++
				t.advance()
			}
			if matched == hashCount {
				return Token{Type: TokenRawString, Value: sb.String(), Line: line, Col: col}, nil
			}
			sb.WriteRune('"')
			for i := 0; i < matched; i++ {
				sb.WriteRune('#')
			}
		} else {
			sb.WriteRune(ch)
		}
	}
}

func (t *Tokenizer) readMultiLineRawContent(hashCount int, line, col int) (Token, error) {
	closing := `"""` + strings.Repeat("#", hashCount)

	var sb strings.Builder
	for {
		if t.pos >= len(t.input) {
			return Token{}, t.error("unterminated multi-line raw string")
		}
		ch := t.peek()
		if ch == '"' {
			// Try to match closing sequence
			match := true
			for i, c := range closing {
				if t.peekAt(i) != c {
					match = false
					break
				}
			}
			if match {
				for range closing {
					t.advance()
				}
				content := dedentMultiLineString(sb.String())
				return Token{Type: TokenMultiLineString, Value: content, Line: line, Col: col}, nil
			}
		}
		sb.WriteRune(t.advance())
	}
}

func (t *Tokenizer) readQuotedString() (Token, error) {
	line, col := t.line, t.col
	t.advance() // consume opening "

	// Check for multi-line string: """
	if t.peek() == '"' && t.peekAt(1) == '"' {
		t.advance() // second "
		t.advance() // third "

		// Must be followed by newline
		if t.peek() == '\r' || t.peek() == '\n' {
			t.skipNewline()
		} else {
			return Token{}, t.error("multi-line string opening must be followed by newline")
		}

		t.setVersion(Version2)
		return t.readMultiLineString(line, col)
	}

	var sb strings.Builder
	for {
		if t.pos >= len(t.input) {
			return Token{}, t.error("unterminated string")
		}
		ch := t.advance()
		switch ch {
		case '"':
			return Token{Type: TokenQuotedString, Value: sb.String(), Line: line, Col: col}, nil
		case '\\':
			escaped, err := t.readEscape()
			if err != nil {
				return Token{}, err
			}
			sb.WriteString(escaped)
		case '\n', '\r':
			return Token{}, t.error("unexpected newline in string (use multi-line string or \\n)")
		default:
			sb.WriteRune(ch)
		}
	}
}

func (t *Tokenizer) readMultiLineString(line, col int) (Token, error) {
	var sb strings.Builder
	for {
		if t.pos >= len(t.input) {
			return Token{}, t.error("unterminated multi-line string")
		}
		ch := t.peek()
		if ch == '"' && t.peekAt(1) == '"' && t.peekAt(2) == '"' {
			t.advance()
			t.advance()
			t.advance()
			content := dedentMultiLineString(sb.String())
			return Token{Type: TokenMultiLineString, Value: content, Line: line, Col: col}, nil
		}
		ch = t.advance()
		if ch == '\\' {
			escaped, err := t.readEscape()
			if err != nil {
				return Token{}, err
			}
			sb.WriteString(escaped)
		} else {
			sb.WriteRune(ch)
		}
	}
}

// dedentMultiLineString removes common leading whitespace from multi-line strings.
// The last line's whitespace determines the indentation to strip.
func dedentMultiLineString(s string) string {
	lines := strings.Split(s, "\n")

	if len(lines) == 0 {
		return ""
	}

	// The last line (before closing """) determines the indent
	lastLine := lines[len(lines)-1]
	indent := ""
	for _, ch := range lastLine {
		if ch == ' ' || ch == '\t' {
			indent += string(ch)
		} else {
			break
		}
	}

	// Strip that indent from all lines, and remove the last line if it's only whitespace
	var result []string
	for i, line := range lines {
		if i == len(lines)-1 {
			// Last line: only include if it has content after indent
			stripped := strings.TrimPrefix(line, indent)
			if stripped != "" {
				result = append(result, stripped)
			}
			continue
		}
		if line == "" {
			result = append(result, "")
		} else {
			result = append(result, strings.TrimPrefix(line, indent))
		}
	}

	return strings.Join(result, "\n")
}

func (t *Tokenizer) readEscape() (string, error) {
	if t.pos >= len(t.input) {
		return "", t.error("unterminated escape sequence")
	}
	ch := t.advance()
	switch ch {
	case 'n':
		return "\n", nil
	case 'r':
		return "\r", nil
	case 't':
		return "\t", nil
	case '\\':
		return "\\", nil
	case '"':
		return "\"", nil
	case 'b':
		return "\b", nil
	case 'f':
		return "\f", nil
	case 's':
		// v2 escape for space
		t.setVersion(Version2)
		return " ", nil
	case 'u':
		return t.readUnicodeEscape()
	default:
		return "", t.error(fmt.Sprintf("unknown escape sequence: \\%c", ch))
	}
}

func (t *Tokenizer) readUnicodeEscape() (string, error) {
	if t.pos >= len(t.input) || t.peek() != '{' {
		return "", t.error("expected '{' in unicode escape")
	}
	t.advance() // consume {

	var hexStr strings.Builder
	for t.pos < len(t.input) && t.peek() != '}' {
		hexStr.WriteRune(t.advance())
	}
	if t.pos >= len(t.input) {
		return "", t.error("unterminated unicode escape")
	}
	t.advance() // consume }

	var codePoint int
	_, err := fmt.Sscanf(hexStr.String(), "%x", &codePoint)
	if err != nil {
		return "", t.error(fmt.Sprintf("invalid unicode escape: %s", hexStr.String()))
	}

	if !utf8.ValidRune(rune(codePoint)) {
		return "", t.error(fmt.Sprintf("invalid unicode code point: U+%04X", codePoint))
	}

	return string(rune(codePoint)), nil
}

func (t *Tokenizer) readNumber() (Token, error) {
	line, col := t.line, t.col
	var sb strings.Builder

	// Optional sign
	if t.peek() == '-' || t.peek() == '+' {
		sb.WriteRune(t.advance())
	}

	// Check if this is actually an identifier (e.g., just "-" followed by a letter, or bare sign)
	if sb.Len() > 0 && (t.pos >= len(t.input) || !isDigit(t.peek())) {
		// It's a bare identifier like "-tag", "--flag", "+.", or just "+" / "-"
		if t.pos < len(t.input) && isIdentChar(t.peek()) {
			for t.pos < len(t.input) && isIdentChar(t.peek()) {
				sb.WriteRune(t.advance())
			}
		}
		// A sign alone or sign+ident chars is a bare identifier
		return Token{Type: TokenIdentifier, Value: sb.String(), Line: line, Col: col}, nil
	}

	if t.pos >= len(t.input) {
		// Bare sign at EOF is an identifier
		if sb.Len() > 0 {
			return Token{Type: TokenIdentifier, Value: sb.String(), Line: line, Col: col}, nil
		}
		return Token{}, t.error("unexpected end of input in number")
	}

	// Check for hex, octal, binary
	if t.peek() == '0' && t.pos+1 < len(t.input) {
		next := t.peekAt(1)
		switch next {
		case 'x', 'X':
			return t.readHexNumber(sb.String(), line, col)
		case 'o', 'O':
			return t.readOctalNumber(sb.String(), line, col)
		case 'b', 'B':
			return t.readBinaryNumber(sb.String(), line, col)
		}
	}

	// Decimal integer or float
	isFloat := false
	hasE := false

	for t.pos < len(t.input) {
		ch := t.peek()
		if ch == '_' {
			t.advance() // skip underscores in numbers
			continue
		}
		if isDigit(ch) {
			sb.WriteRune(t.advance())
		} else if ch == '.' && !isFloat && !hasE {
			isFloat = true
			sb.WriteRune(t.advance())
		} else if (ch == 'e' || ch == 'E') && !hasE {
			isFloat = true
			hasE = true
			sb.WriteRune(t.advance())
			if t.pos < len(t.input) && (t.peek() == '+' || t.peek() == '-') {
				sb.WriteRune(t.advance())
			}
		} else {
			break
		}
	}

	if isFloat {
		return Token{Type: TokenFloat, Value: sb.String(), Line: line, Col: col}, nil
	}
	return Token{Type: TokenInteger, Value: sb.String(), Line: line, Col: col}, nil
}

func (t *Tokenizer) readHexNumber(prefix string, line, col int) (Token, error) {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteRune(t.advance()) // 0
	sb.WriteRune(t.advance()) // x

	for t.pos < len(t.input) {
		ch := t.peek()
		if ch == '_' {
			t.advance()
			continue
		}
		if isHexDigit(ch) {
			sb.WriteRune(t.advance())
		} else {
			break
		}
	}
	return Token{Type: TokenHexInt, Value: sb.String(), Line: line, Col: col}, nil
}

func (t *Tokenizer) readOctalNumber(prefix string, line, col int) (Token, error) {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteRune(t.advance()) // 0
	sb.WriteRune(t.advance()) // o

	for t.pos < len(t.input) {
		ch := t.peek()
		if ch == '_' {
			t.advance()
			continue
		}
		if ch >= '0' && ch <= '7' {
			sb.WriteRune(t.advance())
		} else {
			break
		}
	}
	return Token{Type: TokenOctalInt, Value: sb.String(), Line: line, Col: col}, nil
}

func (t *Tokenizer) readBinaryNumber(prefix string, line, col int) (Token, error) {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteRune(t.advance()) // 0
	sb.WriteRune(t.advance()) // b

	for t.pos < len(t.input) {
		ch := t.peek()
		if ch == '_' {
			t.advance()
			continue
		}
		if ch == '0' || ch == '1' {
			sb.WriteRune(t.advance())
		} else {
			break
		}
	}
	return Token{Type: TokenBinaryInt, Value: sb.String(), Line: line, Col: col}, nil
}

func (t *Tokenizer) readIdentifier() (Token, error) {
	line, col := t.line, t.col
	var sb strings.Builder

	for t.pos < len(t.input) && isIdentChar(t.peek()) {
		sb.WriteRune(t.advance())
	}

	val := sb.String()

	// v1 keywords
	switch val {
	case "true":
		t.setVersion(Version1)
		return Token{Type: TokenTrue, Value: val, Line: line, Col: col}, nil
	case "false":
		t.setVersion(Version1)
		return Token{Type: TokenFalse, Value: val, Line: line, Col: col}, nil
	case "null":
		t.setVersion(Version1)
		return Token{Type: TokenNull, Value: val, Line: line, Col: col}, nil
	case "inf":
		return Token{Type: TokenInf, Value: val, Line: line, Col: col}, nil
	case "-inf":
		return Token{Type: TokenNegInf, Value: val, Line: line, Col: col}, nil
	case "nan":
		return Token{Type: TokenNaN, Value: val, Line: line, Col: col}, nil
	}

	return Token{Type: TokenIdentifier, Value: val, Line: line, Col: col}, nil
}

func (t *Tokenizer) setVersion(v Version) {
	if t.Version == VersionUnknown {
		t.Version = v
	}
}

// isIdentStart returns true if ch can start a bare identifier.
func isIdentStart(ch rune) bool {
	if ch <= 0 {
		return false
	}
	// Cannot start with a digit, sign, or various special chars
	if isDigit(ch) {
		return false
	}
	return isIdentChar(ch)
}

// isIdentChar returns true if ch can be part of a bare identifier.
func isIdentChar(ch rune) bool {
	if ch <= 0 {
		return false
	}
	// Disallowed in identifiers
	switch ch {
	case '\\', '/', '(', ')', '{', '}', '<', '>', ';', '[', ']', '=', ',', '"':
		return false
	}
	// No whitespace
	if unicode.IsSpace(ch) {
		return false
	}
	// No control characters
	if unicode.IsControl(ch) {
		return false
	}
	return true
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
