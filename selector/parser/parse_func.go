package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseFunc(p *Parser) (ast.Expr, error) {
	p.pushScope(scopeFuncArgs)
	defer p.popScope()

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

func parseArgs(p *Parser) ([]ast.Expr, error) {
	args := make([]ast.Expr, 0)
	for p.hasToken() {
		if p.current().IsKind(lexer.CloseParen) {
			p.advance()
			break
		}

		arg, _, err := p.parseExpression(nil)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.current().IsKind(lexer.Comma) {
			p.advance()
		}
	}
	return args, nil
}
