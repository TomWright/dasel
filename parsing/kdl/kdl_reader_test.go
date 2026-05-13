package kdl

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

func newTestReader(t *testing.T) parsing.Reader {
	t.Helper()
	r, err := newKDLReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestReader_ScalarNodes(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`name "Bob"
age 76
active true`))
	if err != nil {
		t.Fatal(err)
	}

	assertMapStringValue(t, val, "name", "Bob")
	assertMapIntValue(t, val, "age", 76)
	assertMapBoolValue(t, val, "active", true)
}

func TestReader_DuplicateNodes(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`plugin "git"
plugin "docker"`))
	if err != nil {
		t.Fatal(err)
	}

	plugins, err := val.GetMapKey("plugin")
	if err != nil {
		t.Fatal(err)
	}
	if plugins.Type() != model.TypeSlice {
		t.Fatalf("expected slice, got %s", plugins.Type())
	}
	length, err := plugins.SliceLen()
	if err != nil {
		t.Fatal(err)
	}
	if length != 2 {
		t.Fatalf("expected 2 plugins, got %d", length)
	}
}

func TestReader_NodeWithChildren(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`server 80 host="localhost" {
    tls true
}`))
	if err != nil {
		t.Fatal(err)
	}

	server, err := val.GetMapKey("server")
	if err != nil {
		t.Fatal(err)
	}

	// Check $args
	args, err := server.GetMapKey("$args")
	if err != nil {
		t.Fatal(err)
	}
	if args.Type() != model.TypeSlice {
		t.Fatalf("expected slice for $args, got %s", args.Type())
	}

	// Check host property
	assertMapStringValue(t, server, "host", "localhost")

	// Check tls child
	assertMapBoolValue(t, server, "tls", true)
}

func TestReader_EmptyNode(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`empty-node`))
	if err != nil {
		t.Fatal(err)
	}

	node, err := val.GetMapKey("empty-node")
	if err != nil {
		t.Fatal(err)
	}
	if node.Type() != model.TypeNull {
		t.Errorf("expected null, got %s", node.Type())
	}
}

func TestReader_EmptyDocument(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(``))
	if err != nil {
		t.Fatal(err)
	}
	if val.Type() != model.TypeMap {
		t.Errorf("expected map, got %s", val.Type())
	}
}

func TestReader_NestedChildren(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`a {
    b {
        c "deep"
    }
}`))
	if err != nil {
		t.Fatal(err)
	}

	a, err := val.GetMapKey("a")
	if err != nil {
		t.Fatal(err)
	}
	b, err := a.GetMapKey("b")
	if err != nil {
		t.Fatal(err)
	}
	assertMapStringValue(t, b, "c", "deep")
}

func TestReader_V2Keywords(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`node #true #false #null`))
	if err != nil {
		t.Fatal(err)
	}

	node, err := val.GetMapKey("node")
	if err != nil {
		t.Fatal(err)
	}
	// Should be a map with $args
	args, err := node.GetMapKey("$args")
	if err != nil {
		t.Fatal(err)
	}

	elem0, err := args.GetSliceIndex(0)
	if err != nil {
		t.Fatal(err)
	}
	b, err := elem0.BoolValue()
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Errorf("expected true, got false")
	}
}

func TestReader_NumberTypes(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`hex 0xff
octal 0o77
binary 0b1010
float 3.14`))
	if err != nil {
		t.Fatal(err)
	}

	assertMapIntValue(t, val, "hex", 255)
	assertMapIntValue(t, val, "octal", 63)
	assertMapIntValue(t, val, "binary", 10)
	assertMapFloatValue(t, val, "float", 3.14)
}

func TestReader_MultipleArguments(t *testing.T) {
	r := newTestReader(t)
	val, err := r.Read([]byte(`node 1 2 3`))
	if err != nil {
		t.Fatal(err)
	}

	node, err := val.GetMapKey("node")
	if err != nil {
		t.Fatal(err)
	}

	args, err := node.GetMapKey("$args")
	if err != nil {
		t.Fatal(err)
	}

	length, err := args.SliceLen()
	if err != nil {
		t.Fatal(err)
	}
	if length != 3 {
		t.Fatalf("expected 3 args, got %d", length)
	}
}

// Helpers

func assertMapStringValue(t *testing.T, val *model.Value, key, expected string) {
	t.Helper()
	v, err := val.GetMapKey(key)
	if err != nil {
		t.Fatalf("key %q: %v", key, err)
	}
	s, err := v.StringValue()
	if err != nil {
		t.Fatalf("key %q string: %v", key, err)
	}
	if s != expected {
		t.Errorf("key %q: expected %q, got %q", key, expected, s)
	}
}

func assertMapIntValue(t *testing.T, val *model.Value, key string, expected int64) {
	t.Helper()
	v, err := val.GetMapKey(key)
	if err != nil {
		t.Fatalf("key %q: %v", key, err)
	}
	n, err := v.IntValue()
	if err != nil {
		t.Fatalf("key %q int: %v", key, err)
	}
	if n != expected {
		t.Errorf("key %q: expected %d, got %d", key, expected, n)
	}
}

func assertMapFloatValue(t *testing.T, val *model.Value, key string, expected float64) {
	t.Helper()
	v, err := val.GetMapKey(key)
	if err != nil {
		t.Fatalf("key %q: %v", key, err)
	}
	f, err := v.FloatValue()
	if err != nil {
		t.Fatalf("key %q float: %v", key, err)
	}
	if f != expected {
		t.Errorf("key %q: expected %f, got %f", key, expected, f)
	}
}

func assertMapBoolValue(t *testing.T, val *model.Value, key string, expected bool) {
	t.Helper()
	v, err := val.GetMapKey(key)
	if err != nil {
		t.Fatalf("key %q: %v", key, err)
	}
	b, err := v.BoolValue()
	if err != nil {
		t.Fatalf("key %q bool: %v", key, err)
	}
	if b != expected {
		t.Errorf("key %q: expected %v, got %v", key, expected, b)
	}
}
