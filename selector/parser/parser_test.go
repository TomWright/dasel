package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v2/selector/ast"
	"github.com/tomwright/dasel/v2/selector/lexer"
	"github.com/tomwright/dasel/v2/selector/parser"
)

func TestParser_Parse(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected ast.Expr
	}{
		{
			name:  "single property",
			input: "foo",
			expected: &ast.CallExpr{
				Function: "property",
				Args:     ast.Expressions{&ast.StringExpr{Value: "foo"}},
			},
		},
		{
			name:  "chained properties",
			input: "foo.bar",
			expected: ast.ChainedExpr{
				Exprs: ast.Expressions{
					&ast.CallExpr{
						Function: "property",
						Args:     ast.Expressions{&ast.StringExpr{Value: "foo"}},
					},
					&ast.CallExpr{
						Function: "property",
						Args:     ast.Expressions{&ast.StringExpr{Value: "bar"}},
					},
				},
			},
		},
		{
			name:  "single function with no args",
			input: "all()",
			expected: &ast.CallExpr{
				Function: "all",
				Args:     ast.Expressions{},
			},
		},
		{
			name:  "single function with various args",
			input: "all(\"foo\", 'bar', false, TRUE, 123, 12.3, hello, funcOne(), funcTwo(1, 2, 3), asd[5])",
			expected: &ast.CallExpr{
				Function: "all",
				Args: ast.Expressions{
					&ast.StringExpr{Value: "foo"},
					&ast.StringExpr{Value: "bar"},
					&ast.BoolExpr{Value: false},
					&ast.BoolExpr{Value: true},
					&ast.NumberIntExpr{Value: 123},
					&ast.NumberFloatExpr{Value: 12.3},
					&ast.CallExpr{
						Function: "property",
						Args:     ast.Expressions{&ast.StringExpr{Value: "hello"}},
					},
					&ast.CallExpr{
						Function: "funcOne",
						Args:     ast.Expressions{},
					},
					&ast.CallExpr{
						Function: "funcTwo",
						Args: ast.Expressions{
							&ast.NumberIntExpr{Value: 1},
							&ast.NumberIntExpr{Value: 2},
							&ast.NumberIntExpr{Value: 3},
						},
					},
					&ast.CallExpr{
						Function: "property",
						Args:     ast.Expressions{&ast.StringExpr{Value: "asd"}},
					},
					&ast.CallExpr{
						Function: "index",
						Args: ast.Expressions{
							&ast.NumberIntExpr{Value: 5},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := lexer.NewTokenizer(tc.input).Tokenize()
			if err != nil {
				t.Fatal(err)
			}
			got, err := parser.NewParser(tokens).Parse()
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(tc.expected, got) {
				t.Fatalf("unexpected result: %s", cmp.Diff(tc.expected, got))
			}
		})
	}
}
