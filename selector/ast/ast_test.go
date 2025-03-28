package ast

import "testing"

// TestExpr_expr tests the expr method of all the types in the ast package.
// Note that this doesn't actually do anything and is just forcing test coverage.
// The expr func only exists for type safety with the Expr interface.
func TestExpr_expr(t *testing.T) {
	NumberFloatExpr{}.expr()
	NumberIntExpr{}.expr()
	StringExpr{}.expr()
	BoolExpr{}.expr()
	BinaryExpr{}.expr()
	UnaryExpr{}.expr()
	CallExpr{}.expr()
	ChainedExpr{}.expr()
	SpreadExpr{}.expr()
	RangeExpr{}.expr()
	IndexExpr{}.expr()
	ArrayExpr{}.expr()
	PropertyExpr{}.expr()
	ObjectExpr{}.expr()
	MapExpr{}.expr()
	VariableExpr{}.expr()
	GroupExpr{}.expr()
	ConditionalExpr{}.expr()
	BranchExpr{}.expr()
}
