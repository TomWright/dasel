package ast

type Program struct {
	Statements []Statement
}

type Statement struct {
	Expressions Expr
}

type Expressions []Expr

type Expr interface {
	expr()
}

func IsType[T Expr](e Expr) bool {
	_, ok := AsType[T](e)
	return ok
}

func AsType[T Expr](e Expr) (T, bool) {
	v, ok := e.(T)
	return v, ok
}

func LastAsType[T Expr](e Expr) (T, bool) {
	return AsType[T](Last(e))
}

func Last(e Expr) Expr {
	if v, ok := e.(ChainedExpr); ok {
		return v.Exprs[len(v.Exprs)-1]
	}
	return e
}

func RemoveLast(e Expr) Expr {
	var res Expressions
	if v, ok := e.(ChainedExpr); ok {
		res = v.Exprs[0 : len(v.Exprs)-1]
	}
	return ChainExprs(res...)
}
