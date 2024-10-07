package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseSymbol(p *Parser) (ast.Expr, error) {
	token := p.current()

	next := p.peek()

	// Handle functions
	if next.IsKind(lexer.OpenParen) {
		return parseFunc(p)
	}

	if next.IsKind(lexer.OpenBracket) {
		return parseArray(p)
	}

	prop := ast.PropertyExpr{
		Property: ast.StringExpr{Value: token.Value},
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
