package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseVariable(p *Parser) (ast.Expr, error) {
	token := p.current()

	next := p.peek()

	if next.IsKind(lexer.OpenBracket) {
		return parseIndex(p)
	}

	prop := ast.VariableExpr{
		Name: token.Value,
	}

	if next.IsKind(lexer.Spread) {
		p.advanceN(2)
		return ast.ChainExprs(
			prop,
			ast.SpreadExpr{},
		), nil
	}

	p.advance()
	return prop, nil
}
