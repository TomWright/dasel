package parser

import (
	"fmt"

	"github.com/tomwright/dasel/v3/selector/lexer"
)

type PositionalError struct {
	Position int
	Err      error
}

func (e *PositionalError) Error() string {
	return fmt.Sprintf("%v. Position %d.", e.Err, e.Position)
}

type UnexpectedTokenError struct {
	Token lexer.Token
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("unexpected token %v %q at position %d.", e.Token.Kind, e.Token.Value, e.Token.Pos)
}
