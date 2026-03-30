package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseFilter(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Filter); err != nil {
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

	filterExpr := ast.FilterExpr{
		Expr: expr,
	}

	res, err := parseFollowingSymbol(p, filterExpr)
	if err != nil {
		return nil, err
	}

	return res, nil
}
