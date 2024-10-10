package ast

import "regexp"

type NumberFloatExpr struct {
	Value float64
}

func (NumberFloatExpr) expr() {}

type NumberIntExpr struct {
	Value int64
}

func (NumberIntExpr) expr() {}

type StringExpr struct {
	Value string
}

func (StringExpr) expr() {}

type BoolExpr struct {
	Value bool
}

func (BoolExpr) expr() {}

type RegexExpr struct {
	Regex *regexp.Regexp
}

func (RegexExpr) expr() {}
