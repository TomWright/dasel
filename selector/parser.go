package selector

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
	"github.com/tomwright/dasel/v3/selector/parser"
)

func Parse(selector string) (ast.Expr, error) {
	tokens, err := lexer.NewTokenizer(selector).Tokenize()
	if err != nil {
		return nil, err
	}

	return parser.NewParser(tokens).Parse()
}