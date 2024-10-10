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
	if err := p.expect(lexer.OpenParen); err != nil {
		return nil, err
	}
	p.advance()

	expressions, err := p.parseExpressionsAsSlice(
		[]lexer.TokenKind{lexer.CloseParen},
		[]lexer.TokenKind{},
		true,
		bpDefault,
	)
	if err != nil {
		return nil, err
	}

	return ast.MapExpr{
		Exprs: expressions,
	}, nil
}

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

	return ast.FilterExpr{
		Expr: expr,
	}, nil
}

func parseBranch(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Branch); err != nil {
		return nil, err
	}

	p.advance()
	if err := p.expect(lexer.OpenParen); err != nil {
		return nil, err
	}
	p.advance()

	expressions, err := p.parseExpressionsAsSlice(
		[]lexer.TokenKind{lexer.CloseParen},
		[]lexer.TokenKind{lexer.Comma},
		true,
		bpDefault,
	)
	if err != nil {
		return nil, err
	}

	return ast.BranchExpr{
		Exprs: expressions,
	}, nil
}
