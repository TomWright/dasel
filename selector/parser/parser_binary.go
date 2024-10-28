package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
)

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

	//if l, ok := left.(ast.BinaryExpr); ok && l.Operator.Kind == lexer.DoubleQuestionMark {
	//	if r, ok := right.(ast.BinaryExpr); ok && r.Operator.Kind == lexer.DoubleQuestionMark {
	//		return ast.BinaryExpr{
	//			Left:     l.Left,
	//			Operator: l.Operator,
	//			Right: ast.BinaryExpr{
	//				Left:     l.Right,
	//				Operator: r.Operator,
	//				Right:    r.Right,
	//			},
	//		}, nil
	//	}
	//}
	//
	//if r, ok := right.(ast.BinaryExpr); ok && r.Operator.Kind == lexer.DoubleQuestionMark {
	//	return ast.BinaryExpr{
	//		Left: ast.BinaryExpr{
	//			Left:     left,
	//			Operator: operator,
	//			Right:    r.Left,
	//		},
	//		Operator: r.Operator,
	//		Right:    r.Right,
	//	}, nil
	//}
	//
	//if l, ok := left.(ast.BinaryExpr); ok && l.Operator.Kind == lexer.DoubleQuestionMark {
	//	return ast.BinaryExpr{
	//		Left:     l.Left,
	//		Operator: l.Operator,
	//		Right: ast.BinaryExpr{
	//			Left:     l.Right,
	//			Operator: operator,
	//			Right:    right,
	//		},
	//	}, nil
	//}

	return ast.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}, nil
}
