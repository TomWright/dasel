package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseAll(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.All); err != nil {
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

	allExpr := ast.AllExpr{
		Expr: expr,
	}

	res, err := parseFollowingSymbol(p, allExpr)
	if err != nil {
		return nil, err
	}

	return res, nil
}
