package ast

type Expressions []Expr

type Expr interface {
	expr()
}
