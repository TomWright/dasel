package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseIfBody(p *Parser) (ast.Expr, error) {
	return p.parseExpressionsFromTo(lexer.OpenCurly, lexer.CloseCurly, []lexer.TokenKind{}, true, bpDefault)
}

func parseIfCondition(p *Parser) (ast.Expr, error) {
	return p.parseExpressionsFromTo(lexer.OpenParen, lexer.CloseParen, []lexer.TokenKind{}, true, bpDefault)
}

func parseIf(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.If); err != nil {
		return nil, err
	}
	p.advance()

	var exprs []*ast.ConditionalExpr

	process := func(parseCond bool) error {
		var err error
		var cond ast.Expr
		if parseCond {
			cond, err = parseIfCondition(p)
			if err != nil {
				return err
			}
		}

		body, err := parseIfBody(p)
		if err != nil {
			return err
		}

		exprs = append(exprs, &ast.ConditionalExpr{
			Cond: cond,
			Then: body,
		})

		return nil
	}

	if err := process(true); err != nil {
		return nil, err
	}

	for p.current().IsKind(lexer.ElseIf) {
		p.advance()

		if err := process(true); err != nil {
			return nil, err
		}
	}

	if p.current().IsKind(lexer.Else) {
		p.advance()

		body, err := parseIfBody(p)
		if err != nil {
			return nil, err
		}
		exprs[len(exprs)-1].Else = body
	}

	for i := len(exprs) - 1; i >= 0; i-- {
		if i > 0 {
			exprs[i-1].Else = *exprs[i]
		}
	}

	return *exprs[0], nil
}

func (p *Parser) parseExpressionsFromTo(
	from lexer.TokenKind,
	to lexer.TokenKind,
	splitOn []lexer.TokenKind,
	requireExpressions bool,
	bp bindingPower,
) (ast.Expr, error) {
	if err := p.expect(from); err != nil {
		return nil, err
	}
	p.advance()

	t, err := p.parseExpressions(
		[]lexer.TokenKind{to},
		splitOn,
		requireExpressions,
		bp,
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}
