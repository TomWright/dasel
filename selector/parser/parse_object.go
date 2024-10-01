package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseObject(p *Parser) (ast.Expr, error) {
	p.pushScope(scopeObject)
	defer p.popScope()

	if err := p.expect(lexer.OpenCurly); err != nil {
		return nil, err
	}
	p.advance()

	pairs := make([]ast.KeyValue, 0)

	for {
		if p.current().IsKind(lexer.CloseCurly) {
			break
		}

		if p.current().IsKind(lexer.Comma) {
			p.advance()
			continue
		}

		if p.current().IsKind(lexer.Spread) {
			p.advance()
			pairs = append(pairs, ast.KeyValue{
				Key:   ast.SpreadExpr{},
				Value: ast.SpreadExpr{},
			})
			if err := p.expect(lexer.Comma, lexer.CloseCurly); err != nil {
				return nil, err
			}
			continue
		}

		if p.current().IsKind(lexer.Symbol) && p.peek().IsKind(lexer.Comma, lexer.CloseCurly) {
			// if the next token is a comma or close curly, then it is a shorthand property
			pairs = append(pairs, ast.KeyValue{
				Key:   ast.StringExpr{Value: p.current().Value},
				Value: ast.PropertyExpr{Property: ast.StringExpr{Value: p.current().Value}},
			})
			p.advance()
			if err := p.expect(lexer.Comma, lexer.CloseCurly); err != nil {
				return nil, err
			}
			continue
		}

		key, err := p.parseExpression(bpDefault)
		if err != nil {
			return nil, err
		}

		// Attempt to simplify the key to a string expression.
		if prop, ok := key.(ast.PropertyExpr); ok {
			key = prop.Property
		}

		if err := p.expect(lexer.Equals); err != nil {
			return nil, err
		}
		p.advance()

		val, err := p.parseExpression(bpDefault)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, ast.KeyValue{
			Key:   key,
			Value: val,
		})
		if err := p.expect(lexer.Comma, lexer.CloseCurly); err != nil {
			return nil, err
		}
	}

	if err := p.expect(lexer.CloseCurly); err != nil {
		return nil, err
	}
	p.advance()

	return ast.ObjectExpr{
		Pairs: pairs,
	}, nil
}
