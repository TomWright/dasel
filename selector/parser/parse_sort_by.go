package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseSortBy(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.SortBy); err != nil {
		return nil, err
	}
	p.advance()

	if err := p.expect(lexer.OpenParen); err != nil {
		return nil, err
	}
	p.advance()

	sortExpr, err := p.parseExpressions(
		lexer.TokenKinds(lexer.CloseParen, lexer.Comma),
		nil,
		true,
		bpDefault,
		false,
	)
	if err != nil {
		return nil, err
	}

	res := ast.SortByExpr{
		Expr:       sortExpr,
		Descending: false,
	}

	if p.current().IsKind(lexer.CloseParen) {
		p.advance()
		return res, nil
	}

	if err := p.expect(lexer.Comma); err != nil {
		return nil, err
	}
	p.advance()

	if err := p.expect(lexer.Asc, lexer.Desc); err != nil {
		return nil, err
	}

	if p.current().IsKind(lexer.Desc) {
		res.Descending = true
	}

	p.advance()
	if err := p.expect(lexer.CloseParen); err != nil {
		return nil, err
	}
	p.advance()

	return res, nil
}
