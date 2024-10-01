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

	expressions := make([]ast.Expr, 0)

	var expr ast.Expr
	var err error
	var replaceLast bool
	for {
		if p.current().IsKind(lexer.CloseParen) {
			if len(expressions) == 0 {
				return nil, fmt.Errorf("expected at least one expression in map")
			}
			p.advance()
			break
		}

		if p.current().IsKind(lexer.Dot) {
			p.advance()
			continue
		}

		expr, replaceLast, err = p.parseExpression(expr)
		if err != nil {
			return nil, err
		}
		if replaceLast {
			expressions[len(expressions)-1] = expr
			continue
		}
		expressions = append(expressions, expr)
	}

	return ast.MapExpr{
		Exprs: expressions,
	}, nil
}
