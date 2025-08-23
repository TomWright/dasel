package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseVariable(p *Parser) (ast.Expr, error) {
	token := p.current()

	prop := ast.VariableExpr{
		Name: token.Value,
	}

	p.advance()

	//if p.current().IsKind(lexer.Equals) {
	//	return parseAssignment(p, prop)
	//}

	res, err := parseFollowingSymbol(p, prop)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func parseAssignment(p *Parser, variableExpr ast.VariableExpr) (ast.Expr, error) {
	if err := p.expect(lexer.Equals); err != nil {
		return nil, err
	}
	p.advance()

	res := ast.AssignExpr{
		Variable: variableExpr,
	}

	valueExpr, err := p.parseExpressions(lexer.TokenKinds(lexer.Semicolon, lexer.CloseParen), nil, true, bpAssignment, true)
	if err != nil {
		return nil, err
	}
	res.Value = valueExpr

	return res, nil
}
