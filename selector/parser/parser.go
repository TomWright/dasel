package parser

import (
	"fmt"

	"github.com/tomwright/dasel/v2/selector/ast"
	"github.com/tomwright/dasel/v2/selector/lexer"
)

type Parser struct {
	tokens lexer.Tokens
	i      int
}

func NewParser(tokens lexer.Tokens) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (ast.Expr, error) {
	var expressions ast.Expressions
	for p.hasToken() {
		if p.current().IsKind(lexer.EOF) {
			break
		}
		if p.current().IsKind(lexer.Dot) {
			p.advance()
			continue
		}
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}
	if len(expressions) == 1 {
		return expressions[0], nil
	}
	return &ast.ChainedExpr{Exprs: expressions}, nil
}

func (p *Parser) parseExpression() (ast.Expr, error) {
	switch p.current().Kind {
	case lexer.String:
		return parseStringLiteral(p)
	case lexer.Number:
		return parseNumberLiteral(p)
	case lexer.Symbol:
		return parseSymbol(p)
	case lexer.OpenBracket:
		return parseSquareBrackets(p)
	case lexer.Bool:
		return parseBoolLiteral(p)
	default:
		return nil, &UnexpectedTokenError{
			Token: p.current(),
		}
	}
}

func (p *Parser) hasToken() bool {
	return p.i < len(p.tokens) && !p.tokens[p.i].IsKind(lexer.EOF)
}

func (p *Parser) hasTokenN(n int) bool {
	return p.i+n < len(p.tokens) && !p.tokens[p.i+n].IsKind(lexer.EOF)
}

func (p *Parser) current() lexer.Token {
	if p.hasToken() {
		return p.tokens[p.i]
	}
	return lexer.Token{Kind: lexer.EOF}
}

func (p *Parser) advance() lexer.Token {
	p.i++
	return p.current()
}

func (p *Parser) advanceN(n int) lexer.Token {
	p.i += n
	return p.current()
}

func (p *Parser) peek() lexer.Token {
	return p.peekN(1)
}

func (p *Parser) peekN(n int) lexer.Token {
	if p.i+n >= len(p.tokens) {
		return lexer.Token{Kind: lexer.EOF}
	}
	return p.tokens[p.i+n]
}

func (p *Parser) expect(t lexer.Token, kind ...lexer.TokenKind) error {
	if t.IsKind(kind...) {
		return nil
	}
	return &PositionalError{
		Position: t.Pos,
		Err:      fmt.Errorf("unexpected token: %v", t.Value),
	}
}
