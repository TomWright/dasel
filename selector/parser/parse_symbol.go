package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

// parseFollowingSymbols deals with the expressions following symbols/variables, e.g.
// $this[0][1]['name']
// foo['bar']['baz'][1]
func parseFollowingSymbol(p *Parser, prev ast.Expr) (ast.Expr, error) {
	res := ast.Expressions{prev}

	for p.hasToken() {
		if p.current().IsKind(lexer.Spread) {
			p.advanceN(1)
			res = append(res, ast.SpreadExpr{})
			continue
		}

		// String based indexes
		if p.current().IsKind(lexer.OpenBracket) {

			if p.peekN(1).IsKind(lexer.Spread) && p.peekN(2).IsKind(lexer.CloseBracket) {
				p.advanceN(3)
				res = append(res, ast.SpreadExpr{})
				continue
			}

			if p.peekN(1).IsKind(lexer.Star) && p.peekN(2).IsKind(lexer.CloseBracket) {
				p.advanceN(3)
				res = append(res, ast.SpreadExpr{})
				continue
			}

			e, err := parseIndexSquareBrackets(p, false)
			if err != nil {
				return nil, err
			}
			switch ex := e.(type) {
			case ast.RangeExpr:
				res = append(res, ex)
			case ast.IndexExpr:
				// Convert this to a property expr. This property executor deals
				// with maps + arrays.
				res = append(res, ast.PropertyExpr{
					Property: ex.Index,
				})
			}

			//e, err := p.parseExpressionsFromTo(lexer.OpenBracket, lexer.CloseBracket, nil, true, bpDefault)
			//if err != nil {
			//	return nil, err
			//}
			//res = append(res, ast.PropertyExpr{
			//	Property: e,
			//})
			continue
		}

		break
	}

	return ast.ChainExprs(res...), nil
}

func parseSymbol(p *Parser, withFollowing bool, allowFunc bool) (ast.Expr, error) {
	token := p.current()

	next := p.peek()

	// Handle functions
	if next.IsKind(lexer.OpenParen) && allowFunc {
		return parseFunc(p)
	}

	prop := ast.PropertyExpr{
		Property: ast.StringExpr{Value: token.Value},
	}

	p.advance()

	if withFollowing {
		res, err := parseFollowingSymbol(p, prop)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return prop, nil
}
