package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
	"github.com/tomwright/dasel/v3/selector/parser"
)

type happyTestCase struct {
	input    string
	expected ast.Expr
}

func (tc happyTestCase) run(t *testing.T) {
	tokens, err := lexer.NewTokenizer(tc.input).Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parser.NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(tc.expected, got) {
		t.Errorf("unexpected result: %s", cmp.Diff(tc.expected, got))
	}
}

func TestParser_Parse_HappyPath(t *testing.T) {
	t.Run("branching", func(t *testing.T) {
		t.Run("two branches", happyTestCase{
			input: `branch("hello", len("world"))`,
			expected: ast.BranchExprs(
				ast.StringExpr{Value: "hello"},
				ast.ChainExprs(
					ast.CallExpr{
						Function: "len",
						Args:     ast.Expressions{ast.StringExpr{Value: "world"}},
					},
				),
			),
		}.run)
		t.Run("three branches", happyTestCase{
			input: `branch("foo", "bar", "baz")`,
			expected: ast.BranchExprs(
				ast.StringExpr{Value: "foo"},
				ast.StringExpr{Value: "bar"},
				ast.StringExpr{Value: "baz"},
			),
		}.run)
	})

	t.Run("literal access", func(t *testing.T) {
		t.Run("string", happyTestCase{
			input:    `"hello world"`,
			expected: ast.StringExpr{Value: "hello world"},
		}.run)
		t.Run("int", happyTestCase{
			input:    "42",
			expected: ast.NumberIntExpr{Value: 42},
		}.run)
		t.Run("float", happyTestCase{
			input:    "42.1",
			expected: ast.NumberFloatExpr{Value: 42.1},
		}.run)
		t.Run("whole number float", happyTestCase{
			input:    "42f",
			expected: ast.NumberFloatExpr{Value: 42},
		}.run)
		t.Run("bool true lowercase", happyTestCase{
			input:    "true",
			expected: ast.BoolExpr{Value: true},
		}.run)
		t.Run("bool true uppercase", happyTestCase{
			input:    "TRUE",
			expected: ast.BoolExpr{Value: true},
		}.run)
		t.Run("bool true mixed case", happyTestCase{
			input:    "TrUe",
			expected: ast.BoolExpr{Value: true},
		}.run)
		t.Run("bool false lowercase", happyTestCase{
			input:    "false",
			expected: ast.BoolExpr{Value: false},
		}.run)
		t.Run("bool false uppercase", happyTestCase{
			input:    "FALSE",
			expected: ast.BoolExpr{Value: false},
		}.run)
		t.Run("bool false mixed case", happyTestCase{
			input:    "FaLsE",
			expected: ast.BoolExpr{Value: false},
		}.run)
	})

	t.Run("property access", func(t *testing.T) {
		t.Run("single property access", happyTestCase{
			input:    "foo",
			expected: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
		}.run)
		t.Run("chained property access", happyTestCase{
			input: "foo.bar",
			expected: ast.ChainExprs(
				ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				ast.PropertyExpr{Property: ast.StringExpr{Value: "bar"}},
			),
		}.run)
		t.Run("property access spread", happyTestCase{
			input: "foo...",
			expected: ast.ChainExprs(
				ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				ast.SpreadExpr{},
			),
		}.run)
		t.Run("property access spread into property access", happyTestCase{
			input: "foo....bar",
			expected: ast.ChainExprs(
				ast.ChainExprs(
					ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					ast.SpreadExpr{},
				),
				ast.PropertyExpr{Property: ast.StringExpr{Value: "bar"}},
			),
		}.run)
	})

	t.Run("array access", func(t *testing.T) {
		t.Run("chained with filter", happyTestCase{
			input: "filter(name == \"foo\")[0]",
			expected: ast.ChainExprs(
				ast.FilterExpr{
					Expr: ast.BinaryExpr{
						Left: ast.PropertyExpr{Property: ast.StringExpr{Value: "name"}},
						Operator: lexer.Token{
							Kind:  lexer.Equal,
							Value: "==",
							Pos:   12,
							Len:   2,
						},
						Right: ast.StringExpr{Value: "foo"},
					},
				},
				ast.PropertyExpr{Property: ast.NumberIntExpr{Value: 0}},
			),
		}.run)
		t.Run("chained with map", happyTestCase{
			input: "map(foo)[0]",
			expected: ast.ChainExprs(
				ast.MapExpr{
					Expr: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				},
				ast.PropertyExpr{Property: ast.NumberIntExpr{Value: 0}},
			),
		}.run)
		t.Run("root array", func(t *testing.T) {
			t.Run("index", happyTestCase{
				input: "$this[1]",
				expected: ast.ChainExprs(
					ast.VariableExpr{Name: "this"},
					ast.PropertyExpr{Property: ast.NumberIntExpr{Value: 1}},
				),
			}.run)
			t.Run("range", func(t *testing.T) {
				t.Run("start and end funcs", happyTestCase{
					input: "$this[calcStart(1):calcEnd()]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.RangeExpr{
							Start: ast.CallExpr{
								Function: "calcStart",
								Args: ast.Expressions{
									ast.NumberIntExpr{Value: 1},
								},
							},
							End: ast.CallExpr{
								Function: "calcEnd",
							},
						},
					),
				}.run)
				t.Run("start and end", happyTestCase{
					input: "$this[5:10]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}, End: ast.NumberIntExpr{Value: 10}},
					),
				}.run)
				t.Run("start", happyTestCase{
					input: "$this[5:]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}},
					),
				}.run)
				t.Run("end", happyTestCase{
					input: "$this[:10]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.RangeExpr{End: ast.NumberIntExpr{Value: 10}},
					),
				}.run)
			})
			t.Run("spread", func(t *testing.T) {
				t.Run("standard", happyTestCase{
					input: "$this...",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.SpreadExpr{},
					),
				}.run)
				t.Run("brackets", happyTestCase{
					input: "$this[...]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.SpreadExpr{},
					),
				}.run)
			})
		})
		t.Run("property array", func(t *testing.T) {
			t.Run("index", happyTestCase{
				input: "foo[1]",
				expected: ast.ChainExprs(
					ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					ast.PropertyExpr{Property: ast.NumberIntExpr{Value: 1}},
				),
			}.run)
			t.Run("range", func(t *testing.T) {
				t.Run("start and end funcs", happyTestCase{
					input: "foo[calcStart(1):calcEnd()]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.RangeExpr{
							Start: ast.CallExpr{
								Function: "calcStart",
								Args: ast.Expressions{
									ast.NumberIntExpr{Value: 1},
								},
							},
							End: ast.CallExpr{
								Function: "calcEnd",
							},
						},
					),
				}.run)
				t.Run("start and end", happyTestCase{
					input: "foo[5:10]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}, End: ast.NumberIntExpr{Value: 10}},
					),
				}.run)
				t.Run("start", happyTestCase{
					input: "foo[5:]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}},
					),
				}.run)
				t.Run("end", happyTestCase{
					input: "foo[:10]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.RangeExpr{End: ast.NumberIntExpr{Value: 10}},
					),
				}.run)
			})
			t.Run("spread", func(t *testing.T) {
				t.Run("standard", happyTestCase{
					input: "foo...",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.SpreadExpr{},
					),
				}.run)
				t.Run("brackets", happyTestCase{
					input: "foo[...]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.SpreadExpr{},
					),
				}.run)
			})
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("single property", happyTestCase{
			input: "foo.map(x)",
			expected: ast.ChainExprs(
				ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				ast.MapExpr{
					Expr: ast.PropertyExpr{Property: ast.StringExpr{Value: "x"}},
				},
			),
		}.run)
		t.Run("nested property", happyTestCase{
			input: "foo.map(x.y)",
			expected: ast.ChainExprs(
				ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				ast.MapExpr{
					Expr: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "x"}},
						ast.PropertyExpr{Property: ast.StringExpr{Value: "y"}},
					),
				},
			),
		}.run)
	})

	t.Run("object", func(t *testing.T) {
		t.Run("get single property", happyTestCase{
			input: "{foo}",
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}},
			}},
		}.run)
		t.Run("get multiple properties", happyTestCase{
			input: "{foo, bar, baz}",
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}},
				{Key: ast.StringExpr{Value: "bar"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "bar"}}},
				{Key: ast.StringExpr{Value: "baz"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "baz"}}},
			}},
		}.run)
		t.Run("set single property", happyTestCase{
			input: `{"foo":1}`,
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.NumberIntExpr{Value: 1}},
			}},
		}.run)
		t.Run("set multiple properties", happyTestCase{
			input: `{foo: 1, bar: 2, baz: 3}`,
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.NumberIntExpr{Value: 1}},
				{Key: ast.StringExpr{Value: "bar"}, Value: ast.NumberIntExpr{Value: 2}},
				{Key: ast.StringExpr{Value: "baz"}, Value: ast.NumberIntExpr{Value: 3}},
			}},
		}.run)
		t.Run("combine get set", happyTestCase{
			input: `{
				...,
				nestedSpread...,
				foo,
				bar: 2,
				"baz": evalSomething(),
				"Name": "Tom",
			}`,
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.SpreadExpr{}, Value: nil},
				{Key: ast.SpreadExpr{}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "nestedSpread"}}},
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}},
				{Key: ast.StringExpr{Value: "bar"}, Value: ast.NumberIntExpr{Value: 2}},
				{Key: ast.StringExpr{Value: "baz"}, Value: ast.CallExpr{Function: "evalSomething"}},
				{Key: ast.StringExpr{Value: "Name"}, Value: ast.StringExpr{Value: "Tom"}},
			}},
		}.run)
	})

	t.Run("variables", func(t *testing.T) {
		t.Run("single variable", happyTestCase{
			input:    `$foo`,
			expected: ast.VariableExpr{Name: "foo"},
		}.run)
		t.Run("variable passed to func", happyTestCase{
			input:    `len($foo)`,
			expected: ast.CallExpr{Function: "len", Args: ast.Expressions{ast.VariableExpr{Name: "foo"}}},
		}.run)
	})

	t.Run("combinations and grouping", func(t *testing.T) {
		t.Run("string concat with grouping", happyTestCase{
			input: `(foo.a) + (foo.b)`,
			expected: ast.BinaryExpr{
				Left:     ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "a"}}),
				Operator: lexer.Token{Kind: lexer.Plus, Value: "+", Pos: 8, Len: 1},
				Right:    ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "b"}}),
			},
		}.run)
		t.Run("string concat with nested properties", happyTestCase{
			input: `foo.a + foo.b`,
			expected: ast.BinaryExpr{
				Left:     ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "a"}}),
				Operator: lexer.Token{Kind: lexer.Plus, Value: "+", Pos: 6, Len: 1},
				Right:    ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "b"}}),
			},
		}.run)
	})

	t.Run("conditional", func(t *testing.T) {
		t.Run("if", happyTestCase{
			input: `if (foo == 1) { "yes" } else { "no" }`,
			expected: ast.ConditionalExpr{
				Cond: ast.BinaryExpr{
					Left:     ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					Operator: lexer.Token{Kind: lexer.Equal, Value: "==", Pos: 8, Len: 2},
					Right:    ast.NumberIntExpr{Value: 1},
				},
				Then: ast.StringExpr{Value: "yes"},
				Else: ast.StringExpr{Value: "no"},
			},
		}.run)
		t.Run("if elseif else", happyTestCase{
			input: `if (foo == 1) { "yes" } elseif (foo == 2) { "maybe" } else { "no" }`,
			expected: ast.ConditionalExpr{
				Cond: ast.BinaryExpr{
					Left:     ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					Operator: lexer.Token{Kind: lexer.Equal, Value: "==", Pos: 8, Len: 2},
					Right:    ast.NumberIntExpr{Value: 1},
				},
				Then: ast.StringExpr{Value: "yes"},
				Else: ast.ConditionalExpr{
					Cond: ast.BinaryExpr{
						Left:     ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						Operator: lexer.Token{Kind: lexer.Equal, Value: "==", Pos: 36, Len: 2},
						Right:    ast.NumberIntExpr{Value: 2},
					},
					Then: ast.StringExpr{Value: "maybe"},
					Else: ast.StringExpr{Value: "no"},
				},
			},
		}.run)
		t.Run("if elseif elseif else", happyTestCase{
			input: `if (foo == 1) { "yes" } elseif (foo == 2) { "maybe" } elseif (foo == 3) { "probably" } else { "no" }`,
			expected: ast.ConditionalExpr{
				Cond: ast.BinaryExpr{
					Left:     ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					Operator: lexer.Token{Kind: lexer.Equal, Value: "==", Pos: 8, Len: 2},
					Right:    ast.NumberIntExpr{Value: 1},
				},
				Then: ast.StringExpr{Value: "yes"},
				Else: ast.ConditionalExpr{
					Cond: ast.BinaryExpr{
						Left:     ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						Operator: lexer.Token{Kind: lexer.Equal, Value: "==", Pos: 36, Len: 2},
						Right:    ast.NumberIntExpr{Value: 2},
					},
					Then: ast.StringExpr{Value: "maybe"},
					Else: ast.ConditionalExpr{
						Cond: ast.BinaryExpr{
							Left:     ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
							Operator: lexer.Token{Kind: lexer.Equal, Value: "==", Pos: 66, Len: 2},
							Right:    ast.NumberIntExpr{Value: 3},
						},
						Then: ast.StringExpr{Value: "probably"},
						Else: ast.StringExpr{Value: "no"},
					},
				},
			},
		}.run)
	})

	t.Run("coalesce", func(t *testing.T) {
		t.Run("chained on left side", happyTestCase{
			input: `foo ?? bar ?? baz`,
			expected: ast.BinaryExpr{
				Left: ast.BinaryExpr{
					Left:     ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					Operator: lexer.Token{Kind: lexer.DoubleQuestionMark, Value: "??", Pos: 4, Len: 2},
					Right:    ast.PropertyExpr{Property: ast.StringExpr{Value: "bar"}},
				},
				Operator: lexer.Token{Kind: lexer.DoubleQuestionMark, Value: "??", Pos: 11, Len: 2},
				Right:    ast.PropertyExpr{Property: ast.StringExpr{Value: "baz"}},
			},
		}.run)

		t.Run("chained nested on left side", happyTestCase{
			input: `nested.one ?? nested.two ?? nested.three ?? 10`,
			expected: ast.BinaryExpr{
				Left: ast.BinaryExpr{
					Left: ast.BinaryExpr{
						Left: ast.ChainExprs(
							ast.PropertyExpr{Property: ast.StringExpr{Value: "nested"}},
							ast.PropertyExpr{Property: ast.StringExpr{Value: "one"}},
						),
						Operator: lexer.Token{Kind: lexer.DoubleQuestionMark, Value: "??", Pos: 11, Len: 2},
						Right: ast.ChainExprs(
							ast.PropertyExpr{Property: ast.StringExpr{Value: "nested"}},
							ast.PropertyExpr{Property: ast.StringExpr{Value: "two"}},
						),
					},
					Operator: lexer.Token{Kind: lexer.DoubleQuestionMark, Value: "??", Pos: 25, Len: 2},
					Right: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "nested"}},
						ast.PropertyExpr{Property: ast.StringExpr{Value: "three"}},
					),
				},
				Operator: lexer.Token{Kind: lexer.DoubleQuestionMark, Value: "??", Pos: 41, Len: 2},
				Right:    ast.NumberIntExpr{Value: 10},
			},
		}.run)
	})
}
