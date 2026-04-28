package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseTernary(p *Parser, cond ast.Expr) (ast.Expr, error) {
	// consume '?'
	if err := p.expect(lexer.QuestionMark); err != nil {
		return nil, err
	}
	p.advance()

	// parse the "then" expression — stops at Colon since it has no binding power
	then, err := p.parseExpression(bpDefault)
	if err != nil {
		return nil, err
	}
	if then == nil {
		return nil, &UnexpectedTokenError{Token: p.current()}
	}

	// expect and consume ':'
	if err := p.expect(lexer.Colon); err != nil {
		return nil, err
	}
	p.advance()

	// parse the "else" expression
	els, err := p.parseExpression(bpDefault)
	if err != nil {
		return nil, err
	}
	if els == nil {
		return nil, &UnexpectedTokenError{Token: p.current()}
	}

	return ast.ConditionalExpr{Cond: cond, Then: then, Else: els}, nil
}
