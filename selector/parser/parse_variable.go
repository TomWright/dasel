package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
)

func parseVariable(p *Parser) (ast.Expr, error) {
	token := p.current()

	prop := ast.VariableExpr{
		Name: token.Value,
	}

	p.advance()

	res, err := parseFollowingSymbol(p, prop)
	if err != nil {
		return nil, err
	}

	return res, nil
}
