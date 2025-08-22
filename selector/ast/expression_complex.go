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
	// Property can resolve to a string or number.
	// If it resolves to a number, we expect to be reading from an array.
	// If it resolves to a string, we expect to be reading from a map.
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
	Expr Expr
}

func (MapExpr) expr() {}

type FilterExpr struct {
	Expr Expr
}

func (FilterExpr) expr() {}

type SearchExpr struct {
	Expr Expr
}

func (SearchExpr) expr() {}

type RecursiveDescentExpr struct {
	IsWildcard bool
	Expr       Expr
}

func (RecursiveDescentExpr) expr() {}

type SortByExpr struct {
	Expr       Expr
	Descending bool
}

func (SortByExpr) expr() {}

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

type AssignExpr struct {
	Variable VariableExpr
	Value    Expr
}

func (AssignExpr) expr() {}
