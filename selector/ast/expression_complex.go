package ast

import "github.com/tomwright/dasel/v3/selector/lexer"

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

func ChainExprs(exprs ...Expr) Expr {
	if len(exprs) == 1 {
		return exprs[0]
	}
	return ChainedExpr{
		Exprs: exprs,
	}
}

func (ChainedExpr) expr() {}

type SpreadExpr struct{}

func (SpreadExpr) expr() {}

func IsSpreadExpr(e Expr) bool {
	_, ok := e.(SpreadExpr)
	return ok
}

type RangeExpr struct {
	Start Expr
	End   Expr
}

func (RangeExpr) expr() {}

type IndexExpr struct {
	Index Expr
}

func (IndexExpr) expr() {}

type PropertyExpr struct {
	Property Expr
}

func (PropertyExpr) expr() {}

type KeyValue struct {
	Key   Expr
	Value Expr
}

type ObjectExpr struct {
	Pairs []KeyValue
}

func (ObjectExpr) expr() {}

type MapExpr struct {
	Exprs Expressions
}

func (MapExpr) expr() {}
