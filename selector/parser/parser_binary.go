package parser

import "github.com/tomwright/dasel/v3/selector/ast"

func parseBinary(p *Parser, left ast.Expr) (ast.Expr, error) {
	if err := p.expect(leftDenotationTokens...); err != nil {
		return nil, err
	}
	operator := p.current()
	p.advance()
	right, err := p.parseExpression(getTokenBindingPower(operator.Kind))
	if err != nil {
		return nil, err
	}
	return ast.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}, nil
}
