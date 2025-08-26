package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseRecursiveDescent(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.RecursiveDescent); err != nil {
		return nil, err
	}
	p.advance()

	cur := p.current()

	res := ast.RecursiveDescentExpr{}

	var err error
	switch cur.Kind {
	case lexer.Star:
		res.IsWildcard = true
		p.advance()
	case lexer.OpenBracket:
		res.Expr, err = parseIndexSquareBrackets(p, true)
	default:
		res.Expr, err = parseSymbol(p, false, false)
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}
