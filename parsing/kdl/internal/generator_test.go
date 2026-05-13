package internal

import (
	"strings"
	"testing"
)

func TestGenerator_SimpleNode(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{Name: "node"},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(result) != "node" {
		t.Errorf("expected 'node', got %q", result)
	}
}

func TestGenerator_NodeWithArgs(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "node",
				Arguments: []*Value{
					{Value: "hello"},
					{Value: int64(42)},
					{Value: true},
				},
			},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	expected := `node "hello" 42 #true`
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_NodeWithProperties(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "node",
				Properties: []*Property{
					{Key: "key", Value: &Value{Value: "value"}},
					{Key: "count", Value: &Value{Value: int64(5)}},
				},
			},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	expected := `node key="value" count=5`
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_NodeWithChildren(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "parent",
				Children: []*Node{
					{Name: "child1"},
					{Name: "child2"},
				},
			},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	expected := "parent {\n    child1\n    child2\n}"
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_NullValue(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{Name: "node", Arguments: []*Value{{Value: nil}}},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(result) != "node #null" {
		t.Errorf("expected 'node #null', got %q", strings.TrimSpace(result))
	}
}

func TestGenerator_Float(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{Name: "node", Arguments: []*Value{{Value: float64(3.14)}}},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(result) != "node 3.14" {
		t.Errorf("expected 'node 3.14', got %q", strings.TrimSpace(result))
	}
}

func TestGenerator_TypeAnnotation(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "node",
				Type: "mytype",
				Arguments: []*Value{
					{Type: "u8", Value: int64(42)},
				},
			},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	expected := "(mytype)node (u8)42"
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_QuotedNodeName(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{Name: "my node"},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(result) != `"my node"` {
		t.Errorf("expected quoted name, got %q", strings.TrimSpace(result))
	}
}

func TestGenerator_CompactMode(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "parent",
				Children: []*Node{
					{Name: "child", Arguments: []*Value{{Value: int64(1)}}},
				},
			},
		},
	}
	opts := GenerateOptions{Compact: true}
	result, err := GenerateString(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	expected := "parent{child 1;}"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestGenerator_StringEscaping(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{Name: "node", Arguments: []*Value{{Value: "hello\nworld"}}},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	expected := `node "hello\nworld"`
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_RoundTrip(t *testing.T) {
	inputs := []string{
		`node "hello" 42`,
		"parent {\n    child 1\n}",
		`node key="value"`,
		`node #true #false #null`,
	}

	for _, input := range inputs {
		doc, err := Parse(input)
		if err != nil {
			t.Fatalf("parse %q: %v", input, err)
		}
		result, err := GenerateString(doc, DefaultGenerateOptions())
		if err != nil {
			t.Fatalf("generate %q: %v", input, err)
		}

		// Re-parse to verify
		doc2, err := Parse(result)
		if err != nil {
			t.Fatalf("re-parse %q: %v", result, err)
		}

		if len(doc.Nodes) != len(doc2.Nodes) {
			t.Errorf("round-trip %q: node count mismatch: %d vs %d", input, len(doc.Nodes), len(doc2.Nodes))
		}
	}
}

func TestGenerator_BoolFalse(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{Name: "node", Arguments: []*Value{{Value: false}}},
		},
	}
	result, err := GenerateString(doc, DefaultGenerateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(result) != "node #false" {
		t.Errorf("expected 'node #false', got %q", strings.TrimSpace(result))
	}
}

func TestGenerator_V1Output(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "node",
				Arguments: []*Value{
					{Value: true},
					{Value: false},
					{Value: nil},
				},
			},
		},
	}
	opts := DefaultGenerateOptions()
	opts.Version = Version1
	result, err := GenerateString(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	expected := "node true false null"
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_V2Output(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "node",
				Arguments: []*Value{
					{Value: true},
					{Value: false},
					{Value: nil},
				},
			},
		},
	}
	opts := DefaultGenerateOptions()
	opts.Version = Version2
	result, err := GenerateString(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	expected := "node #true #false #null"
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_V1DefaultsToV2(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{Name: "node", Arguments: []*Value{{Value: true}}},
		},
	}
	// Zero-value Version should default to v2
	opts := GenerateOptions{Indent: "    "}
	result, err := GenerateString(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(result) != "node #true" {
		t.Errorf("expected v2 output by default, got %q", strings.TrimSpace(result))
	}
}

func TestGenerator_V1Properties(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "node",
				Properties: []*Property{
					{Key: "active", Value: &Value{Value: true}},
					{Key: "empty", Value: &Value{Value: nil}},
				},
			},
		},
	}
	opts := DefaultGenerateOptions()
	opts.Version = Version1
	result, err := GenerateString(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	expected := "node active=true empty=null"
	if strings.TrimSpace(result) != expected {
		t.Errorf("expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestGenerator_V1Children(t *testing.T) {
	doc := &Document{
		Nodes: []*Node{
			{
				Name: "parent",
				Children: []*Node{
					{Name: "enabled", Arguments: []*Value{{Value: true}}},
					{Name: "data", Arguments: []*Value{{Value: nil}}},
				},
			},
		},
	}
	opts := DefaultGenerateOptions()
	opts.Version = Version1
	result, err := GenerateString(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "enabled true") {
		t.Errorf("expected v1 'enabled true' in output, got:\n%s", result)
	}
	if !strings.Contains(result, "data null") {
		t.Errorf("expected v1 'data null' in output, got:\n%s", result)
	}
}

func TestGenerator_V1RoundTrip(t *testing.T) {
	// Parse v1 input, generate as v1, re-parse — should match
	input := "node true false null"
	doc, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	opts := DefaultGenerateOptions()
	opts.Version = Version1
	output, err := GenerateString(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(output) != input {
		t.Errorf("expected %q, got %q", input, strings.TrimSpace(output))
	}
}
