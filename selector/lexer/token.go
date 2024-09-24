package lexer

import (
	"fmt"
	"slices"
)

type TokenKind int

const (
	EOF TokenKind = iota
	Symbol
	Comma
	Colon
	OpenBracket
	CloseBracket
	OpenCurly
	CloseCurly
	OpenParen
	CloseParen
	Equal
	NotEqual
	Not
	And
	Or
	Like
	NotLike
	String
	Number
	Bool
	Add
	Increment
	IncrementBy
	Subtract
	Decrement
	DecrementBy
	Multiply
	Divide
	Modulus
	Dot
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

type PositionalError struct {
	Pos int
	Err error
}

func (e *PositionalError) Error() string {
	return fmt.Sprintf("%v. Position %d.", e.Pos, e.Err)
}

type UnexpectedTokenError struct {
	Pos   int
	Token rune
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("unexpected token: %s at position %d.", string(e.Token), e.Pos)
}
