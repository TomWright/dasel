package lexer

import (
	"fmt"
	"slices"
)

type TokenKind int

func TokenKinds(tk ...TokenKind) []TokenKind {
	return tk
}

const (
	EOF TokenKind = iota
	Symbol
	Comma
	Colon
	OpenBracket  // [
	CloseBracket // ]
	OpenCurly
	CloseCurly
	OpenParen
	CloseParen
	Equal    // ==
	Equals   // =
	NotEqual // !=
	And
	Or
	Like    // =~
	NotLike // !~
	String
	Number
	Bool
	Plus
	Increment
	IncrementBy
	Dash
	Decrement
	DecrementBy
	Star
	Slash
	Percent
	Dot
	Spread           // ...
	RecursiveDescent // ..
	Dollar
	Variable
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
	Exclamation
	Null
	If
	Else
	ElseIf
	Branch
	Map
	Filter
	Search
	RegexPattern
	SortBy
	Asc
	Desc
	QuestionMark
	DoubleQuestionMark
	Semicolon
)

type Tokens []Token

func (tt Tokens) Split(kind TokenKind) []Tokens {
	var res []Tokens
	var cur Tokens
	for _, t := range tt {
		if t.Kind == kind {
			if len(cur) > 0 {
				res = append(res, cur)
			}
			cur = nil
			continue
		}
		cur = append(cur, t)
	}
	if len(cur) > 0 {
		res = append(res, cur)
	}
	return res
}

type Token struct {
	Kind  TokenKind
	Value string
	Pos   int
	Len   int
}

func NewToken(kind TokenKind, value string, pos int, len int) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Pos:   pos,
		Len:   len,
	}
}

func (t Token) IsKind(kind ...TokenKind) bool {
	return slices.Contains(kind, t.Kind)
}

type UnexpectedTokenError struct {
	Pos   int
	Token rune
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("failed to tokenize: unexpected token: %s at position %d.", string(e.Token), e.Pos)
}

type UnexpectedEOFError struct {
	Pos int
}

func (e *UnexpectedEOFError) Error() string {
	return fmt.Sprintf("failed to tokenize: unexpected EOF at position %d.", e.Pos)
}
