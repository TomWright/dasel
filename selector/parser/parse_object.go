package parser

import (
	"fmt"

	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseObject(p *Parser) (ast.Expr, error) {

	//p.parseExpressionsFromTo(
	//	lexer.OpenCurly,
	//	lexer.CloseCurly,
	//	lexer.TokenKinds(lexer.Comma),
	//	false,
	//	bpDefault,
	//)

	if err := p.expect(lexer.OpenCurly); err != nil {
		return nil, err
	}
	p.advance()

	pairs := make([]ast.KeyValue, 0)

	parseKeyValue := func() (ast.KeyValue, error) {
		var res ast.KeyValue
		k, err := p.parseExpression(bpDefault)
		if err != nil {
			return res, err
		}

		// Handle spread
		kSpread, isSpread := ast.LastAsType[ast.SpreadExpr](k)
		if isSpread {
			res.Key = kSpread
			res.Value = ast.RemoveLast(k)
			if err := p.expect(lexer.Comma, lexer.CloseCurly); err != nil {
				return res, err
			}
			return res, nil
		}

		kProp, kIsProp := ast.AsType[ast.PropertyExpr](k)
		if p.current().IsKind(lexer.Comma, lexer.CloseCurly) {
			if !kIsProp {
				return res, fmt.Errorf("invalid shorthand property")
			}
			res.Key = kProp.Property
			res.Value = kProp
			return res, nil
		}

		// Handle unquoted keys
		if kIsProp {
			if kStr, ok := ast.AsType[ast.StringExpr](kProp.Property); ok {
				k = kStr
			}
		}

		if err := p.expect(lexer.Colon); err != nil {
			return res, err
		}
		p.advance()

		v, err := p.parseExpression(bpDefault)
		if err != nil {
			return res, err
		}

		res.Key = k
		res.Value = v
		return res, nil
	}

	for !p.current().IsKind(lexer.CloseCurly) {
		kv, err := parseKeyValue()
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, kv)

		if err := p.expect(lexer.Comma, lexer.CloseCurly); err != nil {
			return nil, fmt.Errorf("expected end of object element: %w", err)
		}
		if p.current().IsKind(lexer.Comma) {
			p.advance()
		}
	}

	if err := p.expect(lexer.CloseCurly); err != nil {
		return nil, fmt.Errorf("expected end of object: %w", err)
	}
	p.advance()

	return ast.ObjectExpr{
		Pairs: pairs,
	}, nil
}
