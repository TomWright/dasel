package parser

import (
	"slices"

	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
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

func (p *Parser) parseExpressionsAsSlice(
	breakOn []lexer.TokenKind,
	splitOn []lexer.TokenKind,
	requireExpressions bool,
	bp bindingPower,
	advanceOnBreak bool,
) (ast.Expressions, error) {
	var finalExpr ast.Expressions
	var current ast.Expressions
	for p.hasToken() {
		if p.current().IsKind(breakOn...) {
			if advanceOnBreak {
				p.advance()
			}
			break
		}
		if p.current().IsKind(splitOn...) {
			if requireExpressions && len(current) == 0 {
				return nil, &UnexpectedTokenError{Token: p.current()}
			}
			finalExpr = append(finalExpr, ast.ChainExprs(current...))
			current = nil
			p.advance()
			continue
		}
		expr, err := p.parseExpression(bp)
		if err != nil {
			return nil, err
		}
		current = append(current, expr)
	}

	if len(current) > 0 {
		finalExpr = append(finalExpr, ast.ChainExprs(current...))
	}

	if len(finalExpr) == 0 {
		return nil, nil
	}

	return finalExpr, nil
}

func (p *Parser) parseExpressions(
	breakOn []lexer.TokenKind,
	splitOn []lexer.TokenKind,
	requireExpressions bool,
	bp bindingPower,
	advanceOnBreak bool,
) (ast.Expr, error) {
	expressions, err := p.parseExpressionsAsSlice(breakOn, splitOn, requireExpressions, bp, advanceOnBreak)
	if err != nil {
		return nil, err
	}
	switch len(expressions) {
	case 0:
		return nil, nil
	default:
		return ast.ChainExprs(expressions...), nil
	}
}

func (p *Parser) Parse() (ast.Expr, error) {
	return p.parseExpressions([]lexer.TokenKind{lexer.EOF}, nil, true, bpDefault, true)
}

func (p *Parser) parseExpression(bp bindingPower) (left ast.Expr, err error) {
	if p.hasToken() && slices.Contains(rightDenotationTokens, p.current().Kind) {
		unary := ast.UnaryExpr{
			Operator: p.current(),
			Right:    nil,
		}
		p.advance()
		expr, err := p.parseExpression(getTokenBindingPower(unary.Operator.Kind))
		if err != nil {
			return nil, err
		}
		p.advance()
		unary.Right = expr
		left = unary
	}

	if !p.hasToken() {
		return
	}

	switch p.current().Kind {
	case lexer.String:
		left, err = parseStringLiteral(p)
	case lexer.Number:
		left, err = parseNumberLiteral(p)
	case lexer.Dash:
		if p.peek().Kind == lexer.Number {
			left, err = parseNumberLiteral(p)
		} else {
			return nil, &UnexpectedTokenError{
				Token: p.current(),
			}
		}
	case lexer.Symbol:
		left, err = parseSymbol(p, true)
	case lexer.OpenBracket:
		left, err = parseArray(p)
	case lexer.OpenCurly:
		left, err = parseObject(p)
	case lexer.Bool:
		left, err = parseBoolLiteral(p)
	case lexer.Spread:
		left, err = parseSpread(p)
	case lexer.Variable:
		left, err = parseVariable(p)
	case lexer.OpenParen:
		left, err = parseGroup(p)
	case lexer.If:
		left, err = parseIf(p)
	case lexer.Branch:
		left, err = parseBranch(p)
	case lexer.Map:
		left, err = parseMap(p)
	case lexer.Filter:
		left, err = parseFilter(p)
	case lexer.Search:
		left, err = parseSearch(p)
	case lexer.RecursiveDescent:
		left, err = parseRecursiveDescent(p)
	case lexer.RegexPattern:
		left, err = parseRegexPattern(p)
	case lexer.SortBy:
		left, err = parseSortBy(p)
	case lexer.Null:
		left = ast.NullExpr{}
		err = nil
		p.advance()
	default:
		return nil, &UnexpectedTokenError{
			Token: p.current(),
		}
	}

	if err != nil {
		return
	}

	toChain := ast.Expressions{left}
	// Ensure dot separated chains are parsed as a sequence of expressions
	if p.hasToken() && p.current().IsKind(lexer.Dot) {
		for p.hasToken() && p.current().IsKind(lexer.Dot) {
			p.advance()
			expr, err := p.parseExpression(bpUnary)
			if err != nil {
				return nil, err
			}
			toChain = append(toChain, expr)
		}
	}

	// Handle spread
	if p.hasToken() && p.current().IsKind(lexer.Spread) {
		expr, err := p.parseExpression(bpLiteral)
		if err != nil {
			return nil, err
		}
		toChain = append(toChain, expr)
	}

	if len(toChain) > 1 {
		left = ast.ChainExprs(toChain...)
	}

	// Handle binding powers
	for p.hasToken() && slices.Contains(leftDenotationTokens, p.current().Kind) && getTokenBindingPower(p.current().Kind) > bp {
		left, err = parseBinary(p, left)
		if err != nil {
			return
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
