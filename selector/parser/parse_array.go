package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseArray(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Symbol, lexer.Variable); err != nil {
		return nil, err
	}
	if err := p.expectN(1, lexer.OpenBracket); err != nil {
		return nil, err
	}
	token := p.current()
	p.advance()

	idx, err := parseSquareBrackets(p)
	if err != nil {
		return nil, err
	}

	var e ast.Expr

	switch {
	case token.IsKind(lexer.Variable):
		e = ast.VariableExpr{
			Name: token.Value,
		}
	case token.IsKind(lexer.Symbol):
		e = ast.PropertyExpr{
			Property: ast.StringExpr{Value: token.Value},
		}
	default:
		panic("unexpected token kind")
	}

	return ast.ChainExprs(
		e,
		idx,
	), nil
}

// parseSquareBrackets parses square bracket array access.
// E.g. [0], [0:1], [0:], [:2]
func parseSquareBrackets(p *Parser) (ast.Expr, error) {
	// Handle index (from bracket)
	if err := p.expect(lexer.OpenBracket); err != nil {
		return nil, err
	}

	p.advance()

	// Spread [...]
	if p.current().IsKind(lexer.Spread) {
		p.advance()
		if err := p.expect(lexer.CloseBracket); err != nil {
			return nil, err
		}
		p.advance()
		return ast.SpreadExpr{}, nil
	}

	var (
		start ast.Expr
		end   ast.Expr
		err   error
	)

	if p.current().IsKind(lexer.Colon) {
		p.advance()
		// We have no start index
		end, err = p.parseExpression(bpDefault)
		if err != nil {
			return nil, err
		}
		p.advance()
		return ast.RangeExpr{
			End: end,
		}, nil
	}

	start, err = p.parseExpression(bpDefault)
	if err != nil {
		return nil, err
	}

	if p.current().IsKind(lexer.CloseBracket) {
		// This is an index
		p.advance()
		return ast.IndexExpr{
			Index: start,
		}, nil
	}

	if err := p.expect(lexer.Colon); err != nil {
		return nil, err
	}

	p.advance()

	if p.current().IsKind(lexer.CloseBracket) {
		// There is no end
		p.advance()
		return ast.RangeExpr{
			Start: start,
		}, nil
	}

	end, err = p.parseExpression(bpDefault)
	if err != nil {
		return nil, err
	}

	p.advance()
	return ast.RangeExpr{
		Start: start,
		End:   end,
	}, nil
}
