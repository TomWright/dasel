package parser

import (
	"fmt"

	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseMap(p *Parser) (ast.Expr, error) {
	p.pushScope(scopeMap)
	defer p.popScope()

	if err := p.expect(lexer.Symbol); err != nil {
		return nil, err
	}
	if p.current().Value != "map" {
		return nil, fmt.Errorf("expected map but got %q", p.current().Value)
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
		bpCall,
	)
	if err != nil {
		return nil, err
	}

	return ast.MapExpr{
		Exprs: expressions,
	}, nil
}
