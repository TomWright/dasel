package internal

import (
	"testing"
)

func TestParser_SimpleNode(t *testing.T) {
	doc, err := Parse("node")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}
	if doc.Nodes[0].Name != "node" {
		t.Errorf("expected name 'node', got %q", doc.Nodes[0].Name)
	}
}

func TestParser_NodeWithArguments(t *testing.T) {
	doc, err := Parse(`node "hello" 42 true`)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Nodes[0]
	if len(node.Arguments) != 3 {
		t.Fatalf("expected 3 arguments, got %d", len(node.Arguments))
	}
	if node.Arguments[0].Value != "hello" {
		t.Errorf("expected 'hello', got %v", node.Arguments[0].Value)
	}
	if node.Arguments[1].Value != int64(42) {
		t.Errorf("expected 42, got %v", node.Arguments[1].Value)
	}
	if node.Arguments[2].Value != true {
		t.Errorf("expected true, got %v", node.Arguments[2].Value)
	}
}

func TestParser_NodeWithProperties(t *testing.T) {
	doc, err := Parse(`node key="value" count=5`)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Nodes[0]
	if len(node.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(node.Properties))
	}
	if node.Properties[0].Key != "key" || node.Properties[0].Value.Value != "value" {
		t.Errorf("expected key='value', got %s=%v", node.Properties[0].Key, node.Properties[0].Value.Value)
	}
	if node.Properties[1].Key != "count" || node.Properties[1].Value.Value != int64(5) {
		t.Errorf("expected count=5, got %s=%v", node.Properties[1].Key, node.Properties[1].Value.Value)
	}
}

func TestParser_NodeWithChildren(t *testing.T) {
	doc, err := Parse("parent {\n  child1\n  child2\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}
	parent := doc.Nodes[0]
	if len(parent.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(parent.Children))
	}
	if parent.Children[0].Name != "child1" {
		t.Errorf("expected 'child1', got %q", parent.Children[0].Name)
	}
	if parent.Children[1].Name != "child2" {
		t.Errorf("expected 'child2', got %q", parent.Children[1].Name)
	}
}

func TestParser_MixedNode(t *testing.T) {
	doc, err := Parse(`server 80 host="localhost" {
    tls true
}`)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Nodes[0]
	if node.Name != "server" {
		t.Errorf("expected 'server', got %q", node.Name)
	}
	if len(node.Arguments) != 1 || node.Arguments[0].Value != int64(80) {
		t.Errorf("expected arg 80, got %v", node.Arguments)
	}
	if len(node.Properties) != 1 || node.Properties[0].Key != "host" {
		t.Errorf("expected prop host, got %v", node.Properties)
	}
	if len(node.Children) != 1 || node.Children[0].Name != "tls" {
		t.Errorf("expected child tls, got %v", node.Children)
	}
}

func TestParser_MultipleNodes(t *testing.T) {
	doc, err := Parse("a\nb\nc")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(doc.Nodes))
	}
}

func TestParser_SemicolonSeparated(t *testing.T) {
	doc, err := Parse("a; b; c")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(doc.Nodes))
	}
}

func TestParser_SlashDashNode(t *testing.T) {
	doc, err := Parse("a\n/- b\nc")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 2 {
		t.Fatalf("expected 2 nodes (b skipped), got %d", len(doc.Nodes))
	}
	if doc.Nodes[0].Name != "a" || doc.Nodes[1].Name != "c" {
		t.Errorf("expected a, c; got %q, %q", doc.Nodes[0].Name, doc.Nodes[1].Name)
	}
}

func TestParser_SlashDashArgument(t *testing.T) {
	doc, err := Parse(`node 1 /- 2 3`)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Nodes[0]
	if len(node.Arguments) != 2 {
		t.Fatalf("expected 2 args (2 skipped), got %d", len(node.Arguments))
	}
	if node.Arguments[0].Value != int64(1) || node.Arguments[1].Value != int64(3) {
		t.Errorf("expected 1, 3; got %v, %v", node.Arguments[0].Value, node.Arguments[1].Value)
	}
}

func TestParser_SlashDashProperty(t *testing.T) {
	doc, err := Parse(`node a=1 /- b=2 c=3`)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Nodes[0]
	if len(node.Properties) != 2 {
		t.Fatalf("expected 2 props (b skipped), got %d", len(node.Properties))
	}
	if node.Properties[0].Key != "a" || node.Properties[1].Key != "c" {
		t.Errorf("expected a, c; got %q, %q", node.Properties[0].Key, node.Properties[1].Key)
	}
}

func TestParser_SlashDashChildren(t *testing.T) {
	doc, err := Parse("node /- {\n  child\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes[0].Children) != 0 {
		t.Errorf("expected no children (slashdashed), got %d", len(doc.Nodes[0].Children))
	}
}

func TestParser_TypeAnnotation(t *testing.T) {
	doc, err := Parse(`(mytype)node (u8)42`)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Nodes[0]
	if node.Type != "mytype" {
		t.Errorf("expected type 'mytype', got %q", node.Type)
	}
	if node.Arguments[0].Type != "u8" {
		t.Errorf("expected arg type 'u8', got %q", node.Arguments[0].Type)
	}
}

func TestParser_QuotedNodeName(t *testing.T) {
	doc, err := Parse(`"my node" "value"`)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Nodes[0].Name != "my node" {
		t.Errorf("expected 'my node', got %q", doc.Nodes[0].Name)
	}
}

func TestParser_NullValue(t *testing.T) {
	doc, err := Parse("node null")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Nodes[0].Arguments[0].Value != nil {
		t.Errorf("expected nil, got %v", doc.Nodes[0].Arguments[0].Value)
	}
}

func TestParser_V2Keywords(t *testing.T) {
	doc, err := Parse(`node #true #false #null`)
	if err != nil {
		t.Fatal(err)
	}
	args := doc.Nodes[0].Arguments
	if args[0].Value != true {
		t.Errorf("expected true, got %v", args[0].Value)
	}
	if args[1].Value != false {
		t.Errorf("expected false, got %v", args[1].Value)
	}
	if args[2].Value != nil {
		t.Errorf("expected nil, got %v", args[2].Value)
	}
}

func TestParser_HexOctalBinary(t *testing.T) {
	doc, err := Parse("node 0xff 0o77 0b1010")
	if err != nil {
		t.Fatal(err)
	}
	args := doc.Nodes[0].Arguments
	if args[0].Value != int64(255) {
		t.Errorf("expected 255, got %v", args[0].Value)
	}
	if args[1].Value != int64(63) {
		t.Errorf("expected 63, got %v", args[1].Value)
	}
	if args[2].Value != int64(10) {
		t.Errorf("expected 10, got %v", args[2].Value)
	}
}

func TestParser_VersionMarker(t *testing.T) {
	input := `/- kdl-version 2
node #true`
	doc, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	// The /- kdl-version 2 is consumed as a slashdashed node
	// The parser should still produce the node
	found := false
	for _, n := range doc.Nodes {
		if n.Name == "node" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'node' in document")
	}
}

func TestParser_EmptyDocument(t *testing.T) {
	doc, err := Parse("")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(doc.Nodes))
	}
}

func TestParser_CommentsOnly(t *testing.T) {
	doc, err := Parse("// comment\n/* block */")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(doc.Nodes))
	}
}

func TestParser_NestedChildren(t *testing.T) {
	doc, err := Parse("a {\n  b {\n    c\n  }\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}
	a := doc.Nodes[0]
	if len(a.Children) != 1 {
		t.Fatalf("expected 1 child of a, got %d", len(a.Children))
	}
	b := a.Children[0]
	if len(b.Children) != 1 {
		t.Fatalf("expected 1 child of b, got %d", len(b.Children))
	}
	if b.Children[0].Name != "c" {
		t.Errorf("expected 'c', got %q", b.Children[0].Name)
	}
}

func TestParser_NodeNoArgs(t *testing.T) {
	doc, err := Parse("empty-node")
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Nodes[0]
	if len(node.Arguments) != 0 {
		t.Errorf("expected 0 args, got %d", len(node.Arguments))
	}
	if len(node.Properties) != 0 {
		t.Errorf("expected 0 props, got %d", len(node.Properties))
	}
	if len(node.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(node.Children))
	}
}
