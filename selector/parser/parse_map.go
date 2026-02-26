package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseMap(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Map); err != nil {
		return nil, err
	}
	p.advance()

	expr, err := p.parseExpressionsFromTo(
		lexer.OpenParen,
		lexer.CloseParen,
		[]lexer.TokenKind{},
		true,
		bpDefault,
	)
	if err != nil {
		return nil, err
	}

	mapExpr := ast.MapExpr{
		Expr: expr,
	}

	res, err := parseFollowingSymbol(p, mapExpr)
	if err != nil {
		return nil, err
	}

	return res, nil
}
