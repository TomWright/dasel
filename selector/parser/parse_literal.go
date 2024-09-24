package parser

import (
	"strconv"
	"strings"

	"github.com/tomwright/dasel/v2/selector/ast"
)

func parseStringLiteral(p *Parser) (ast.Expr, error) {
	token := p.current()
	p.advance()
	return &ast.StringExpr{
		Value: token.Value,
	}, nil
}

func parseBoolLiteral(p *Parser) (ast.Expr, error) {
	token := p.current()
	p.advance()
	return &ast.BoolExpr{
		Value: strings.EqualFold(token.Value, "true"),
	}, nil
}

func parseNumberLiteral(p *Parser) (ast.Expr, error) {
	token := p.current()
	p.advance()
	if strings.Contains(token.Value, ".") {
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, err
		}
		return &ast.NumberFloatExpr{
			Value: value,
		}, nil
	}
	value, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &ast.NumberIntExpr{
		Value: value,
	}, nil
}
