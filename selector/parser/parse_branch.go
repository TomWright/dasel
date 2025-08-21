package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseBranch(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Branch); err != nil {
		return nil, err
	}

	p.advance()
	if err := p.expect(lexer.OpenParen); err != nil {
		return nil, err
	}
	p.advance()

	expressions, err := p.parseExpressionsAsSlice(
		[]lexer.TokenKind{lexer.CloseParen},
		[]lexer.TokenKind{lexer.Comma},
		true,
		bpDefault,
		true,
	)
	if err != nil {
		return nil, err
	}

	return ast.BranchExpr{
		Exprs: expressions,
	}, nil
}
