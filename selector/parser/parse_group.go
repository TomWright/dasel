package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseGroup(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.OpenParen); err != nil {
		return nil, err
	}
	p.advance() // skip the open paren

	return p.parseExpressions(
		[]lexer.TokenKind{lexer.CloseParen},
		[]lexer.TokenKind{},
		true,
		bpDefault,
	)
}
