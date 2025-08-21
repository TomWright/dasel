package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseStringLiteral(p *Parser) (ast.Expr, error) {
	token := p.current()
	p.advance()
	return ast.StringExpr{
		Value: token.Value,
	}, nil
}

func parseBoolLiteral(p *Parser) (ast.Expr, error) {
	token := p.current()
	p.advance()
	return ast.BoolExpr{
		Value: strings.EqualFold(token.Value, "true"),
	}, nil
}

func parseSpread(p *Parser) (ast.Expr, error) {
	p.advance()
	return ast.SpreadExpr{}, nil
}

func parseNumberLiteral(p *Parser) (ast.Expr, error) {
	token := p.current()

	var negative bool
	if token.IsKind(lexer.Dash) {
		negative = true
		token = p.advance()
	}

	next := p.advance()
	switch {
	case next.IsKind(lexer.Symbol) && strings.EqualFold(next.Value, "f"):
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, err
		}
		p.advance()
		if negative {
			value = -value
		}
		return ast.NumberFloatExpr{
			Value: value,
		}, nil

	case strings.Contains(token.Value, "."):
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, err
		}
		if negative {
			value = -value
		}
		return ast.NumberFloatExpr{
			Value: value,
		}, nil

	default:
		value, err := strconv.ParseInt(token.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		if negative {
			value = -value
		}
		return ast.NumberIntExpr{
			Value: value,
		}, nil
	}
}

func parseRegexPattern(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.RegexPattern); err != nil {
		return nil, err
	}

	pattern := p.current()

	p.advance()

	comp, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regexp pattern: %w", err)
	}

	return ast.RegexExpr{
		Regex: comp,
	}, nil
}
