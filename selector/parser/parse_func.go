package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseFunc(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Symbol); err != nil {
		return nil, err
	}
	if err := p.expectN(1, lexer.OpenParen); err != nil {
		return nil, err
	}

	token := p.current()

	p.advanceN(2)
	args, err := parseArgs(p)
	if err != nil {
		return nil, err
	}
	return ast.CallExpr{
		Function: token.Value,
		Args:     args,
	}, nil
}

func parseArgs(p *Parser) (ast.Expressions, error) {
	return p.parseExpressionsAsSlice(
		[]lexer.TokenKind{lexer.CloseParen},
		[]lexer.TokenKind{lexer.Comma},
		false,
		bpCall,
		true,
	)
}
