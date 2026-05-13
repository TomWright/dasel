package internal

import (
	"strings"
	"testing"
)

// Tests for converting between KDL v1 and v2 syntax.
// These verify that:
// 1. v1 input produces the same AST as equivalent v2 input
// 2. Generating from either AST produces valid v2 output
// 3. Round-tripping through parse→generate→parse preserves semantics

// --- v1 ↔ v2 equivalent AST tests ---

func TestConvert_BooleansSameAST(t *testing.T) {
	// v1: true/false, v2: #true/#false — same AST
	docV1, err := Parse("node true false")
	if err != nil {
		t.Fatal(err)
	}
	docV2, err := Parse("node #true #false")
	if err != nil {
		t.Fatal(err)
	}

	assertSameArgValues(t, docV1.Nodes[0], docV2.Nodes[0])
}

func TestConvert_NullSameAST(t *testing.T) {
	docV1, err := Parse("node null")
	if err != nil {
		t.Fatal(err)
	}
	docV2, err := Parse("node #null")
	if err != nil {
		t.Fatal(err)
	}

	assertSameArgValues(t, docV1.Nodes[0], docV2.Nodes[0])
}

func TestConvert_BooleanPropSameAST(t *testing.T) {
	docV1, err := Parse("node active=true disabled=false")
	if err != nil {
		t.Fatal(err)
	}
	docV2, err := Parse("node active=#true disabled=#false")
	if err != nil {
		t.Fatal(err)
	}

	assertSamePropValues(t, docV1.Nodes[0], docV2.Nodes[0])
}

func TestConvert_NullPropSameAST(t *testing.T) {
	docV1, err := Parse("node val=null")
	if err != nil {
		t.Fatal(err)
	}
	docV2, err := Parse("node val=#null")
	if err != nil {
		t.Fatal(err)
	}

	assertSamePropValues(t, docV1.Nodes[0], docV2.Nodes[0])
}

func TestConvert_RawStringSameAST(t *testing.T) {
	// v1: r"hello", v2: #"hello"#
	docV1, err := Parse(`node r"hello"`)
	if err != nil {
		t.Fatal(err)
	}
	docV2, err := Parse(`node #"hello"#`)
	if err != nil {
		t.Fatal(err)
	}

	assertSameArgValues(t, docV1.Nodes[0], docV2.Nodes[0])
}

func TestConvert_RawStringWithHashesSameAST(t *testing.T) {
	// v1: r#"contains "quotes""#, v2: ##"contains "quotes""##
	docV1, err := Parse(`node r#"contains "quotes""#`)
	if err != nil {
		t.Fatal(err)
	}
	docV2, err := Parse(`node ##"contains "quotes""##`)
	if err != nil {
		t.Fatal(err)
	}

	assertSameArgValues(t, docV1.Nodes[0], docV2.Nodes[0])
}

func TestConvert_RawStringNoEscapes(t *testing.T) {
	// Both v1 and v2 raw strings should NOT process backslash escapes
	docV1, err := Parse(`node r"hello\nworld"`)
	if err != nil {
		t.Fatal(err)
	}
	docV2, err := Parse(`node #"hello\nworld"#`)
	if err != nil {
		t.Fatal(err)
	}

	// Both should have the literal text "hello\nworld" (no actual newline)
	assertArgString(t, docV1.Nodes[0], 0, "hello\\nworld")
	assertArgString(t, docV2.Nodes[0], 0, "hello\\nworld")
	assertSameArgValues(t, docV1.Nodes[0], docV2.Nodes[0])
}

// --- Shared syntax produces identical AST ---

func TestConvert_SharedQuotedStringsSameAST(t *testing.T) {
	input := `node "hello world" "with\nnewline"`
	doc1, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	doc2, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	assertSameArgValues(t, doc1.Nodes[0], doc2.Nodes[0])
}

func TestConvert_SharedNumbersSameAST(t *testing.T) {
	input := "node 42 3.14 -10 0xff 0o77 0b1010 1_000"
	doc1, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	doc2, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	assertSameArgValues(t, doc1.Nodes[0], doc2.Nodes[0])
}

func TestConvert_SharedChildrenSameAST(t *testing.T) {
	input := "parent {\n    child1 \"arg\"\n    child2 42\n}"
	doc1, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	doc2, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc1.Nodes[0].Children) != len(doc2.Nodes[0].Children) {
		t.Fatalf("children count mismatch")
	}
}

// --- v1 input → v2 output (generate always outputs v2) ---

func TestConvert_V1BooleansToV2Output(t *testing.T) {
	doc, err := Parse("node true false")
	if err != nil {
		t.Fatal(err)
	}
	out, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	// Output should use v2 syntax
	if !strings.Contains(out, "#true") || !strings.Contains(out, "#false") {
		t.Errorf("expected v2 booleans in output, got:\n%s", out)
	}
	if strings.Contains(out, " true") && !strings.Contains(out, "#true") {
		t.Errorf("bare 'true' in output (should be #true)")
	}
}

func TestConvert_V1NullToV2Output(t *testing.T) {
	doc, err := Parse("node null")
	if err != nil {
		t.Fatal(err)
	}
	out, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "#null") {
		t.Errorf("expected #null in output, got:\n%s", out)
	}
}

func TestConvert_V1RawStringToV2QuotedOutput(t *testing.T) {
	doc, err := Parse(`node r"hello world"`)
	if err != nil {
		t.Fatal(err)
	}
	out, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	// Generator outputs regular quoted strings, not raw
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted string in output, got:\n%s", out)
	}
}

// --- Full v1→v2 round-trip ---

func TestConvert_V1FullDocumentToV2(t *testing.T) {
	v1Input := `// A v1 document
name "Bob"
age 76
active true
empty null
server {
    host "localhost"
    port 8080
    tls false
}`

	// Parse as v1
	doc, err := Parse(v1Input)
	if err != nil {
		t.Fatal(err)
	}

	// Generate as v2
	v2Output, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}

	// v2 output should use #true/#false/#null
	if strings.Contains(v2Output, " true\n") && !strings.Contains(v2Output, "#true") {
		t.Errorf("expected #true in v2 output")
	}
	if strings.Contains(v2Output, " false\n") && !strings.Contains(v2Output, "#false") {
		t.Errorf("expected #false in v2 output")
	}
	if strings.Contains(v2Output, " null\n") && !strings.Contains(v2Output, "#null") {
		t.Errorf("expected #null in v2 output")
	}

	// Re-parse the v2 output and verify same data
	doc2, err := Parse(v2Output)
	if err != nil {
		t.Fatalf("failed to re-parse v2 output: %v\nOutput:\n%s", err, v2Output)
	}

	// Verify structure preserved
	if len(doc.Nodes) != len(doc2.Nodes) {
		t.Fatalf("node count mismatch: %d vs %d", len(doc.Nodes), len(doc2.Nodes))
	}

	for i, n := range doc.Nodes {
		if n.Name != doc2.Nodes[i].Name {
			t.Errorf("node %d: name mismatch: %q vs %q", i, n.Name, doc2.Nodes[i].Name)
		}
	}
}

// --- v2→v1 equivalent: parsing v2 output with v1 keywords ---

func TestConvert_V2OutputReparsesCorrectly(t *testing.T) {
	inputs := []string{
		"node #true #false #null",
		"node 42 3.14 0xff",
		`node "hello" key="val"`,
		"parent {\n    child 1\n}",
		"a\nb\nc",
	}

	for _, input := range inputs {
		doc, err := Parse(input)
		if err != nil {
			t.Fatalf("parse %q: %v", input, err)
		}

		out, err := GenerateString(doc, DefaultGenerateOptions())
		if err != nil {
			t.Fatalf("generate %q: %v", input, err)
		}

		doc2, err := Parse(out)
		if err != nil {
			t.Fatalf("reparse %q (output: %q): %v", input, out, err)
		}

		if len(doc.Nodes) != len(doc2.Nodes) {
			t.Errorf("round-trip %q: node count %d vs %d", input, len(doc.Nodes), len(doc2.Nodes))
		}

		for i := range doc.Nodes {
			if i >= len(doc2.Nodes) {
				break
			}
			assertSameArgValues(t, doc.Nodes[i], doc2.Nodes[i])
		}
	}
}

// --- Mixed v1/v2 features in same document ---

func TestConvert_MixedV1V2SyntaxShared(t *testing.T) {
	// Common syntax that works in both versions
	input := `node "string" 42 3.14 key="val" {
    child
}`
	doc, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}

	assertNodeName(t, doc.Nodes[0], "node")
	assertArgString(t, doc.Nodes[0], 0, "string")
	assertArgInt(t, doc.Nodes[0], 1, 42)
	assertArgFloat(t, doc.Nodes[0], 2, 3.14)
	assertPropString(t, doc.Nodes[0], "key", "val")
	assertChildCount(t, doc.Nodes[0], 1)
}

// --- Real-world config conversion tests ---

func TestConvert_ZellijConfigV1ToV2(t *testing.T) {
	// Simplified Zellij-like config in v1 style
	v1Config := `keybinds {
    normal {
        bind "Alt" "h" {
            action "MoveFocusOrTab" "Left"
        }
        bind "Alt" "l" {
            action "MoveFocusOrTab" "Right"
        }
    }
}
theme "dracula"
default_shell "zsh"
pane_frames true`

	doc, err := Parse(v1Config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate v2
	v2Output, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}

	// Verify v2 output uses #true
	if !strings.Contains(v2Output, "#true") {
		t.Errorf("expected #true in v2 output for pane_frames")
	}

	// Re-parse and verify round-trip
	doc2, err := Parse(v2Output)
	if err != nil {
		t.Fatalf("re-parse failed: %v\nOutput:\n%s", err, v2Output)
	}

	if len(doc.Nodes) != len(doc2.Nodes) {
		t.Fatalf("node count mismatch: %d vs %d", len(doc.Nodes), len(doc2.Nodes))
	}
}

func TestConvert_PackageConfigV1ToV2(t *testing.T) {
	v1Config := `package {
    name "my-app"
    version "1.0.0"
    description "A great app"
}
dependencies {
    dep "react" "^18.0.0"
    dep "typescript" "^5.0.0"
}
build {
    minify true
    sourcemap false
    target "es2022"
}`

	doc, err := Parse(v1Config)
	if err != nil {
		t.Fatal(err)
	}

	v2Output, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}

	// Should contain v2 booleans
	if !strings.Contains(v2Output, "#true") || !strings.Contains(v2Output, "#false") {
		t.Errorf("expected v2 booleans, got:\n%s", v2Output)
	}

	// Round-trip
	doc2, err := Parse(v2Output)
	if err != nil {
		t.Fatalf("re-parse failed: %v", err)
	}

	assertNodeCount(t, &Document{Nodes: doc2.Nodes}, 3)
}

// --- Version marker conversion ---

func TestConvert_V1MarkerDocumentToV2(t *testing.T) {
	v1Input := "/- kdl-version 1\nnode true false null"
	doc, err := Parse(v1Input)
	if err != nil {
		t.Fatal(err)
	}

	v2Output, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}

	// Verify v2 keywords in output
	if !strings.Contains(v2Output, "#true") {
		t.Errorf("expected #true in converted output")
	}
	if !strings.Contains(v2Output, "#false") {
		t.Errorf("expected #false in converted output")
	}
	if !strings.Contains(v2Output, "#null") {
		t.Errorf("expected #null in converted output")
	}
}

func TestConvert_V2MarkerDocumentRoundTrip(t *testing.T) {
	v2Input := "/- kdl-version 2\nnode #true #false #null"
	doc, err := Parse(v2Input)
	if err != nil {
		t.Fatal(err)
	}

	v2Output, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}

	doc2, err := Parse(v2Output)
	if err != nil {
		t.Fatal(err)
	}

	assertSameArgValues(t, doc.Nodes[0], doc2.Nodes[0])
}

// --- Helpers ---

func assertSameArgValues(t *testing.T, n1, n2 *Node) {
	t.Helper()
	if len(n1.Arguments) != len(n2.Arguments) {
		t.Fatalf("arg count mismatch: %d vs %d (node %q vs %q)",
			len(n1.Arguments), len(n2.Arguments), n1.Name, n2.Name)
	}
	for i := range n1.Arguments {
		v1 := n1.Arguments[i].Value
		v2 := n2.Arguments[i].Value
		if v1 == nil && v2 == nil {
			continue
		}
		if v1 == nil || v2 == nil {
			t.Errorf("arg %d: value mismatch: %v vs %v", i, v1, v2)
			continue
		}
		// Compare stringified values for simplicity with float edge cases
		if !valuesEqual(v1, v2) {
			t.Errorf("arg %d: value mismatch: %v (%T) vs %v (%T)", i, v1, v1, v2, v2)
		}
	}
}

func assertSamePropValues(t *testing.T, n1, n2 *Node) {
	t.Helper()
	if len(n1.Properties) != len(n2.Properties) {
		t.Fatalf("prop count mismatch: %d vs %d", len(n1.Properties), len(n2.Properties))
	}
	for i := range n1.Properties {
		p1 := n1.Properties[i]
		p2 := n2.Properties[i]
		if p1.Key != p2.Key {
			t.Errorf("prop %d: key mismatch: %q vs %q", i, p1.Key, p2.Key)
			continue
		}
		if !valuesEqual(p1.Value.Value, p2.Value.Value) {
			t.Errorf("prop %q: value mismatch: %v vs %v", p1.Key, p1.Value.Value, p2.Value.Value)
		}
	}
}

func valuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Direct comparison handles string, int64, bool
	return a == b
}
