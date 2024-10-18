package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
	"github.com/tomwright/dasel/v3/selector/parser"
)

func TestParser_Parse_HappyPath(t *testing.T) {
	type testCase struct {
		input    string
		expected ast.Expr
	}

	run := func(t *testing.T, tc testCase) func(*testing.T) {
		return func(t *testing.T) {
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
	}

	t.Run("branching", func(t *testing.T) {
		t.Run("two branches", run(t, testCase{
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
		}))
		t.Run("three branches", run(t, testCase{
			input: `branch("foo", "bar", "baz")`,
			expected: ast.BranchExprs(
				ast.StringExpr{Value: "foo"},
				ast.StringExpr{Value: "bar"},
				ast.StringExpr{Value: "baz"},
			),
		}))
	})

	t.Run("literal access", func(t *testing.T) {
		t.Run("string", run(t, testCase{
			input:    `"hello world"`,
			expected: ast.StringExpr{Value: "hello world"},
		}))
		t.Run("int", run(t, testCase{
			input:    "42",
			expected: ast.NumberIntExpr{Value: 42},
		}))
		t.Run("float", run(t, testCase{
			input:    "42.1",
			expected: ast.NumberFloatExpr{Value: 42.1},
		}))
		t.Run("whole number float", run(t, testCase{
			input:    "42f",
			expected: ast.NumberFloatExpr{Value: 42},
		}))
		t.Run("bool true lowercase", run(t, testCase{
			input:    "true",
			expected: ast.BoolExpr{Value: true},
		}))
		t.Run("bool true uppercase", run(t, testCase{
			input:    "TRUE",
			expected: ast.BoolExpr{Value: true},
		}))
		t.Run("bool true mixed case", run(t, testCase{
			input:    "TrUe",
			expected: ast.BoolExpr{Value: true},
		}))
		t.Run("bool false lowercase", run(t, testCase{
			input:    "false",
			expected: ast.BoolExpr{Value: false},
		}))
		t.Run("bool false uppercase", run(t, testCase{
			input:    "FALSE",
			expected: ast.BoolExpr{Value: false},
		}))
		t.Run("bool false mixed case", run(t, testCase{
			input:    "FaLsE",
			expected: ast.BoolExpr{Value: false},
		}))
	})

	t.Run("property access", func(t *testing.T) {
		t.Run("single property access", run(t, testCase{
			input:    "foo",
			expected: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
		}))
		t.Run("chained property access", run(t, testCase{
			input: "foo.bar",
			expected: ast.ChainExprs(
				ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				ast.PropertyExpr{Property: ast.StringExpr{Value: "bar"}},
			),
		}))
		t.Run("property access spread", run(t, testCase{
			input: "foo...",
			expected: ast.ChainExprs(
				ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				ast.SpreadExpr{},
			),
		}))
		t.Run("property access spread into property access", run(t, testCase{
			input: "foo....bar",
			expected: ast.ChainExprs(
				ast.ChainExprs(
					ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					ast.SpreadExpr{},
				),
				ast.PropertyExpr{Property: ast.StringExpr{Value: "bar"}},
			),
		}))
	})

	t.Run("array access", func(t *testing.T) {
		t.Run("root array", func(t *testing.T) {
			t.Run("index", run(t, testCase{
				input: "$this[1]",
				expected: ast.ChainExprs(
					ast.VariableExpr{Name: "this"},
					ast.IndexExpr{Index: ast.NumberIntExpr{Value: 1}},
				),
			}))
			t.Run("range", func(t *testing.T) {
				t.Run("start and end funcs", run(t, testCase{
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
				}))
				t.Run("start and end", run(t, testCase{
					input: "$this[5:10]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}, End: ast.NumberIntExpr{Value: 10}},
					),
				}))
				t.Run("start", run(t, testCase{
					input: "$this[5:]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}},
					),
				}))
				t.Run("end", run(t, testCase{
					input: "$this[:10]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.RangeExpr{End: ast.NumberIntExpr{Value: 10}},
					),
				}))
			})
			t.Run("spread", func(t *testing.T) {
				t.Run("standard", run(t, testCase{
					input: "$this...",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.SpreadExpr{},
					),
				}))
				t.Run("brackets", run(t, testCase{
					input: "$this[...]",
					expected: ast.ChainExprs(
						ast.VariableExpr{Name: "this"},
						ast.SpreadExpr{},
					),
				}))
			})
		})
		t.Run("property array", func(t *testing.T) {
			t.Run("index", run(t, testCase{
				input: "foo[1]",
				expected: ast.ChainExprs(
					ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
					ast.IndexExpr{Index: ast.NumberIntExpr{Value: 1}},
				),
			}))
			t.Run("range", func(t *testing.T) {
				t.Run("start and end funcs", run(t, testCase{
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
				}))
				t.Run("start and end", run(t, testCase{
					input: "foo[5:10]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}, End: ast.NumberIntExpr{Value: 10}},
					),
				}))
				t.Run("start", run(t, testCase{
					input: "foo[5:]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.RangeExpr{Start: ast.NumberIntExpr{Value: 5}},
					),
				}))
				t.Run("end", run(t, testCase{
					input: "foo[:10]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.RangeExpr{End: ast.NumberIntExpr{Value: 10}},
					),
				}))
			})
			t.Run("spread", func(t *testing.T) {
				t.Run("standard", run(t, testCase{
					input: "foo...",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.SpreadExpr{},
					),
				}))
				t.Run("brackets", run(t, testCase{
					input: "foo[...]",
					expected: ast.ChainExprs(
						ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
						ast.SpreadExpr{},
					),
				}))
			})
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("single property", run(t, testCase{
			input: "foo.map(x)",
			expected: ast.ChainExprs(
				ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}},
				ast.MapExpr{
					Expr: ast.PropertyExpr{Property: ast.StringExpr{Value: "x"}},
				},
			),
		}))
		t.Run("nested property", run(t, testCase{
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
		}))
	})

	t.Run("object", func(t *testing.T) {
		t.Run("get single property", run(t, testCase{
			input: "{foo}",
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}},
			}},
		}))
		t.Run("get multiple properties", run(t, testCase{
			input: "{foo, bar, baz}",
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}},
				{Key: ast.StringExpr{Value: "bar"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "bar"}}},
				{Key: ast.StringExpr{Value: "baz"}, Value: ast.PropertyExpr{Property: ast.StringExpr{Value: "baz"}}},
			}},
		}))
		t.Run("set single property", run(t, testCase{
			input: `{"foo":1}`,
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.NumberIntExpr{Value: 1}},
			}},
		}))
		t.Run("set multiple properties", run(t, testCase{
			input: `{foo: 1, bar: 2, baz: 3}`,
			expected: ast.ObjectExpr{Pairs: []ast.KeyValue{
				{Key: ast.StringExpr{Value: "foo"}, Value: ast.NumberIntExpr{Value: 1}},
				{Key: ast.StringExpr{Value: "bar"}, Value: ast.NumberIntExpr{Value: 2}},
				{Key: ast.StringExpr{Value: "baz"}, Value: ast.NumberIntExpr{Value: 3}},
			}},
		}))
		t.Run("combine get set", run(t, testCase{
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
		}))
	})

	t.Run("variables", func(t *testing.T) {
		t.Run("single variable", run(t, testCase{
			input:    `$foo`,
			expected: ast.VariableExpr{Name: "foo"},
		}))
		t.Run("variable passed to func", run(t, testCase{
			input:    `len($foo)`,
			expected: ast.CallExpr{Function: "len", Args: ast.Expressions{ast.VariableExpr{Name: "foo"}}},
		}))
	})

	t.Run("combinations and grouping", func(t *testing.T) {
		t.Run("string concat with grouping", run(t, testCase{
			input: `(foo.a) + (foo.b)`,
			expected: ast.BinaryExpr{
				Left:     ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "a"}}),
				Operator: lexer.Token{Kind: lexer.Plus, Value: "+", Pos: 8, Len: 1},
				Right:    ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "b"}}),
			},
		}))
		t.Run("string concat with nested properties", run(t, testCase{
			input: `foo.a + foo.b`,
			expected: ast.BinaryExpr{
				Left:     ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "a"}}),
				Operator: lexer.Token{Kind: lexer.Plus, Value: "+", Pos: 6, Len: 1},
				Right:    ast.ChainExprs(ast.PropertyExpr{Property: ast.StringExpr{Value: "foo"}}, ast.PropertyExpr{Property: ast.StringExpr{Value: "b"}}),
			},
		}))
	})

	t.Run("conditional", func(t *testing.T) {
		t.Run("if", run(t, testCase{
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
		}))
		t.Run("if elseif else", run(t, testCase{
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
		}))
		t.Run("if elseif elseif else", run(t, testCase{
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
		}))
	})
}
