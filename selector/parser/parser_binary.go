package parser

import "github.com/tomwright/dasel/v3/selector/ast"

func parseBinary(p *Parser, left ast.Expr) (ast.Expr, bool, error) {
	if err := p.expect(leftDenotationTokens...); err != nil {
		return nil, false, err
	}
	for {
		if !p.current().IsKind(leftDenotationTokens...) {
			break
		}

		token := p.current()
		p.advance()

		right, _, err := p.parseExpression(left)
		if err != nil {
			return nil, false, err
		}

		left = ast.BinaryExpr{
			Left:     left,
			Operator: token,
			Right:    right,
		}
	}

	return left, true, nil
}
