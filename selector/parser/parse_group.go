package parser

import (
	"fmt"

	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseGroup(p *Parser) (ast.Expr, error) {
	p.pushScope(scopeGroup)
	defer p.popScope()
	if err := p.expect(lexer.OpenParen); err != nil {
		return nil, err
	}
	p.advance() // skip the open paren

	expressions := ast.Expressions{}
	for {
		if p.current().Kind == lexer.CloseParen {
			break
		}

		expr, err := p.parseExpression(bpDefault)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	if err := p.expect(lexer.CloseParen); err != nil {
		return nil, err
	}
	p.advance() // skip the close paren

	if len(expressions) == 0 {
		return nil, fmt.Errorf("group expression must contain at least one expression")
	}

	return ast.ChainExprs(expressions...), nil
}
