package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseCount(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Count); err != nil {
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

	countExpr := ast.CountExpr{
		Expr: expr,
	}

	res, err := parseFollowingSymbol(p, countExpr)
	if err != nil {
		return nil, err
	}

	return res, nil
}
