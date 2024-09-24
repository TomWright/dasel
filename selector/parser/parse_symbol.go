package parser

import (
	"github.com/tomwright/dasel/v2/selector/ast"
	"github.com/tomwright/dasel/v2/selector/lexer"
)

// parseSquareBrackets parses square bracket array access.
// E.g. [0], [0:1], [0:], [:2]. [...]
func parseSquareBrackets(p *Parser) (ast.Expr, error) {
	// Handle index (from bracket)
	p.advance()

	// Spread [...]
	if p.current().IsKind(lexer.Dot) && p.peekN(1).IsKind(lexer.Dot) && p.peekN(2).IsKind(lexer.Dot) {
		p.advanceN(3)
		if err := p.expect(p.current(), lexer.CloseBracket); err != nil {
			return nil, err
		}
		p.advance()
		return &ast.CallExpr{
			Function: "all",
			Args:     ast.Expressions{},
		}, nil
	}

	// Range [1:2]
	if p.current().IsKind(lexer.Number) && p.peekN(1).IsKind(lexer.Colon) && p.peekN(2).IsKind(lexer.Number) {
		from, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.advance()
		to, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(p.current(), lexer.CloseBracket); err != nil {
			return nil, err
		}
		p.advance()
		return &ast.CallExpr{
			Function: "range",
			Args: ast.Expressions{
				from, to,
			},
		}, nil
	}

	// Range [:2]
	if p.current().IsKind(lexer.Colon) && p.peekN(1).IsKind(lexer.Number) {
		from := &ast.NumberIntExpr{Value: -1}
		p.advanceN(1)
		to, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(p.current(), lexer.CloseBracket); err != nil {
			return nil, err
		}
		p.advance()
		return &ast.CallExpr{
			Function: "range",
			Args: ast.Expressions{
				from, to,
			},
		}, nil
	}

	// Range [1:]
	if p.current().IsKind(lexer.Number) && p.peekN(1).IsKind(lexer.Colon) {
		from, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.advanceN(1)
		to := &ast.NumberIntExpr{Value: -1}
		if err := p.expect(p.current(), lexer.CloseBracket); err != nil {
			return nil, err
		}
		p.advance()
		return &ast.CallExpr{
			Function: "range",
			Args: ast.Expressions{
				from, to,
			},
		}, nil
	}

	// Array index [1]
	if err := p.expect(p.current(), lexer.Number); err != nil {
		return nil, err
	}
	index, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expect(p.current(), lexer.CloseBracket); err != nil {
		return nil, err
	}
	p.advance()

	return &ast.CallExpr{
		Function: "index",
		Args:     ast.Expressions{index},
	}, nil
}

func parseSymbol(p *Parser) (ast.Expr, error) {
	token := p.current()

	next := p.peek()

	// Handle functions
	if next.IsKind(lexer.OpenParen) {
		p.advanceN(2)
		args, err := parseArgs(p)
		if err != nil {
			return nil, err
		}
		return &ast.CallExpr{
			Function: token.Value,
			Args:     args,
		}, nil
	}

	// Handle index (before bracket)
	if next.IsKind(lexer.OpenBracket) {
		p.advance()
		return &ast.CallExpr{
			Function: "property",
			Args:     ast.Expressions{&ast.StringExpr{Value: token.Value}},
		}, nil
	}

	if next.IsKind(lexer.Dot, lexer.EOF, lexer.Comma) {
		p.advance()
		return &ast.CallExpr{
			Function: "property",
			Args:     []ast.Expr{&ast.StringExpr{Value: token.Value}},
		}, nil
	}

	return nil, &UnexpectedTokenError{
		Token: next,
	}
}

func parseArgs(p *Parser) ([]ast.Expr, error) {
	args := make([]ast.Expr, 0)
	for p.hasToken() {
		if p.current().IsKind(lexer.CloseParen) {
			p.advance()
			break
		}

		arg, err := p.parseExpression()
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
