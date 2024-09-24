package ast

import "github.com/tomwright/dasel/v2/selector/lexer"

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (BinaryExpr) expr() {}

type CallExpr struct {
	Function string
	Args     Expressions
}

func (CallExpr) expr() {}

type ChainedExpr struct {
	Exprs Expressions
}

func (ChainedExpr) expr() {}
