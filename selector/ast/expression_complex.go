package ast

import "github.com/tomwright/dasel/v3/selector/lexer"

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (BinaryExpr) expr() {}

type UnaryExpr struct {
	Operator lexer.Token
	Right    Expr
}

func (UnaryExpr) expr() {}

type CallExpr struct {
	Function string
	Args     Expressions
}

func (CallExpr) expr() {}

type ChainedExpr struct {
	Exprs Expressions
}

func ChainExprs(exprs ...Expr) Expr {
	if len(exprs) == 0 {
		return nil
	}
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

type ArrayExpr struct {
	Exprs Expressions
}

func (ArrayExpr) expr() {}

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

type FilterExpr struct {
	Expr Expr
}

func (FilterExpr) expr() {}

type VariableExpr struct {
	Name string
}

func (VariableExpr) expr() {}

type GroupExpr struct {
	Expr Expr
}

func (GroupExpr) expr() {}

type ConditionalExpr struct {
	Cond Expr
	Then Expr
	Else Expr
}

func (ConditionalExpr) expr() {}

type BranchExpr struct {
	Exprs []Expr
}

func (BranchExpr) expr() {}

func BranchExprs(exprs ...Expr) Expr {
	return BranchExpr{
		Exprs: exprs,
	}
}
