package parser

import (
	"fmt"

	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

type scope string

const (
	scopeRoot     scope = "root"
	scopeFuncArgs scope = "funcArgs"
	scopeArray    scope = "array"
	scopeObject   scope = "object"
	scopeMap      scope = "map"
)

type Parser struct {
	tokens lexer.Tokens
	i      int
	scopes []scope
}

func (p *Parser) pushScope(s scope) {
	p.scopes = append(p.scopes, s)
}

func (p *Parser) popScope() {
	p.scopes = p.scopes[:len(p.scopes)-1]
}

func (p *Parser) currentScope() scope {
	if len(p.scopes) == 0 {
		return scopeRoot
	}
	return p.scopes[len(p.scopes)-1]
}

func (p *Parser) endOfExpressionTokens() []lexer.TokenKind {
	switch p.currentScope() {
	case scopeRoot:
		return []lexer.TokenKind{lexer.EOF, lexer.Dot}
	case scopeFuncArgs:
		return []lexer.TokenKind{lexer.Comma, lexer.CloseParen}
	case scopeMap:
		return []lexer.TokenKind{lexer.Comma, lexer.CloseParen, lexer.Dot}
	case scopeArray:
		return []lexer.TokenKind{lexer.CloseBracket, lexer.Colon, lexer.Number, lexer.Symbol}
	case scopeObject:
		return []lexer.TokenKind{lexer.CloseCurly, lexer.Equals, lexer.Number, lexer.Symbol, lexer.Comma}
	default:
		return nil
	}
}

func (p *Parser) expectEndOfExpression() error {
	tokens := p.endOfExpressionTokens()
	if len(tokens) == 0 {
		return fmt.Errorf("no end of scope tokens found: %q", p.currentScope())
	}
	return p.expect(tokens...)
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
	switch len(expressions) {
	case 0:
		return nil, nil
	case 1:
		return expressions[0], nil
	default:
		return ast.ChainExprs(expressions...), nil
	}
}

func (p *Parser) parseExpression() (res ast.Expr, err error) {
	defer func() {
		if err == nil {
			err = p.expectEndOfExpression()
		}
	}()
	switch p.current().Kind {
	case lexer.String:
		return parseStringLiteral(p)
	case lexer.Number:
		return parseNumberLiteral(p)
	case lexer.Symbol:
		return parseSymbol(p)
	case lexer.OpenBracket:
		return parseSquareBrackets(p)
	case lexer.OpenCurly:
		return parseObject(p)
	case lexer.Bool:
		return parseBoolLiteral(p)
	case lexer.Spread:
		return parseSpread(p)
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

func (p *Parser) expect(kind ...lexer.TokenKind) error {
	t := p.current()
	if p.current().IsKind(kind...) {
		return nil
	}
	return &UnexpectedTokenError{
		Token: t,
	}
}

func (p *Parser) expectN(n int, kind ...lexer.TokenKind) error {
	t := p.peekN(n)
	if t.IsKind(kind...) {
		return nil
	}
	return &UnexpectedTokenError{
		Token: t,
	}
}
