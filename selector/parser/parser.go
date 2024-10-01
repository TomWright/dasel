package parser

import (
	"fmt"
	"slices"

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
		return append([]lexer.TokenKind{lexer.EOF, lexer.Dot}, leftDenotationTokens...)
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
	var expr ast.Expr
	var err error
	var replaceLast bool
	for p.hasToken() {
		if p.current().IsKind(lexer.EOF) {
			break
		}
		if p.current().IsKind(lexer.Dot) {
			p.advance()
			continue
		}
		expr, replaceLast, err = p.parseExpression(expr)
		if err != nil {
			return nil, err
		}
		if replaceLast {
			expressions[len(expressions)-1] = expr
			continue
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

func (p *Parser) parseExpression(last ast.Expr) (res ast.Expr, replaceLast bool, err error) {
	defer func() {
		if err == nil {
			err = p.expectEndOfExpression()
		}
	}()

	if last != nil && slices.Contains(leftDenotationTokens, p.current().Kind) {
		res, replaceLast, err = parseBinary(p, last)
		return
	}

	switch p.current().Kind {
	case lexer.String:
		res, err = parseStringLiteral(p)
	case lexer.Number:
		res, err = parseNumberLiteral(p)
	case lexer.Symbol:
		res, err = parseSymbol(p)
	case lexer.OpenBracket:
		res, err = parseSquareBrackets(p)
	case lexer.OpenCurly:
		res, err = parseObject(p)
	case lexer.Bool:
		res, err = parseBoolLiteral(p)
	case lexer.Spread:
		res, err = parseSpread(p)
	case lexer.Variable:
		res, err = parseVariable(p)
	default:
		return nil, false, &UnexpectedTokenError{
			Token: p.current(),
		}
	}
	return
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

func (p *Parser) previous() lexer.Token {
	i := p.i - 1
	if i > 0 && i < len(p.tokens) {
		return p.tokens[i]
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
