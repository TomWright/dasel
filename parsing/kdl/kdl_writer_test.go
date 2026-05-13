package kdl

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

func newTestWriter(t *testing.T) parsing.Writer {
	t.Helper()
	w, err := newKDLWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func TestWriter_SimpleScalars(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("name", model.NewStringValue("Bob"))
	_ = val.SetMapKey("age", model.NewIntValue(76))

	w := newTestWriter(t)
	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, `name "Bob"`) {
		t.Errorf("expected name node, got:\n%s", result)
	}
	if !strings.Contains(result, "age 76") {
		t.Errorf("expected age node, got:\n%s", result)
	}
}

func TestWriter_BoolAndNull(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("active", model.NewBoolValue(true))
	_ = val.SetMapKey("empty", model.NewNullValue())

	w := newTestWriter(t)
	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, "active #true") {
		t.Errorf("expected bool node, got:\n%s", result)
	}
	if !strings.Contains(result, "empty") {
		t.Errorf("expected empty node, got:\n%s", result)
	}
}

func TestWriter_Slice(t *testing.T) {
	val := model.NewMapValue()
	plugins := model.NewSliceValue()
	_ = plugins.Append(model.NewStringValue("git"))
	_ = plugins.Append(model.NewStringValue("docker"))
	_ = val.SetMapKey("plugin", plugins)

	w := newTestWriter(t)
	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if strings.Count(result, "plugin") != 2 {
		t.Errorf("expected 2 plugin nodes, got:\n%s", result)
	}
}

func TestWriter_NestedMap(t *testing.T) {
	val := model.NewMapValue()
	server := model.NewMapValue()
	_ = server.SetMapKey("host", model.NewStringValue("localhost"))
	_ = server.SetMapKey("port", model.NewIntValue(80))
	_ = val.SetMapKey("server", server)

	w := newTestWriter(t)
	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, "server") {
		t.Errorf("expected server node, got:\n%s", result)
	}
	if !strings.Contains(result, `host="localhost"`) {
		t.Errorf("expected host property, got:\n%s", result)
	}
}

func TestWriter_CompactMode(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("name", model.NewStringValue("Bob"))

	opts := parsing.WriterOptions{Compact: true}
	w, err := newKDLWriter(opts)
	if err != nil {
		t.Fatal(err)
	}

	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if strings.Contains(result, "\n") {
		t.Errorf("compact mode should not contain newlines, got:\n%s", result)
	}
}

func TestWriter_RoundTrip(t *testing.T) {
	input := `name "Bob"
age 76
active true`

	r := newTestReader(t)
	val, err := r.Read([]byte(input))
	if err != nil {
		t.Fatal(err)
	}

	w := newTestWriter(t)
	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	// Re-read and verify
	val2, err := r.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	assertMapStringValue(t, val2, "name", "Bob")
	assertMapIntValue(t, val2, "age", 76)
	assertMapBoolValue(t, val2, "active", true)
}

func TestWriter_MapWithArgsAndChildren(t *testing.T) {
	// Build a map with $args, properties, and nested children
	val := model.NewMapValue()
	server := model.NewMapValue()

	args := model.NewSliceValue()
	_ = args.Append(model.NewIntValue(80))
	_ = server.SetMapKey("$args", args)

	_ = server.SetMapKey("host", model.NewStringValue("localhost"))

	tls := model.NewMapValue()
	_ = tls.SetMapKey("enabled", model.NewBoolValue(true))
	_ = server.SetMapKey("tls", tls)

	_ = val.SetMapKey("server", server)

	w := newTestWriter(t)
	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, "server 80") {
		t.Errorf("expected server with arg 80, got:\n%s", result)
	}
	if !strings.Contains(result, `host="localhost"`) {
		t.Errorf("expected host property, got:\n%s", result)
	}
	if !strings.Contains(result, "tls") {
		t.Errorf("expected tls child, got:\n%s", result)
	}
}

func TestWriter_V1Output(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("active", model.NewBoolValue(true))
	_ = val.SetMapKey("disabled", model.NewBoolValue(false))
	_ = val.SetMapKey("empty", model.NewNullValue())

	opts := parsing.DefaultWriterOptions()
	opts.Ext["kdl-version"] = "1"
	w, err := newKDLWriter(opts)
	if err != nil {
		t.Fatal(err)
	}

	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, "active true") {
		t.Errorf("expected v1 'active true', got:\n%s", result)
	}
	if !strings.Contains(result, "disabled false") {
		t.Errorf("expected v1 'disabled false', got:\n%s", result)
	}
	// null node has no args in v1 or v2, so just check it exists
	if !strings.Contains(result, "empty") {
		t.Errorf("expected 'empty' node, got:\n%s", result)
	}
}

func TestWriter_V2Output(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("active", model.NewBoolValue(true))
	_ = val.SetMapKey("disabled", model.NewBoolValue(false))

	opts := parsing.DefaultWriterOptions()
	opts.Ext["kdl-version"] = "2"
	w, err := newKDLWriter(opts)
	if err != nil {
		t.Fatal(err)
	}

	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, "active #true") {
		t.Errorf("expected v2 'active #true', got:\n%s", result)
	}
	if !strings.Contains(result, "disabled #false") {
		t.Errorf("expected v2 'disabled #false', got:\n%s", result)
	}
}

func TestWriter_DefaultIsV2(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("flag", model.NewBoolValue(true))

	w := newTestWriter(t)
	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, "#true") {
		t.Errorf("expected v2 output by default, got:\n%s", result)
	}
}

func TestWriter_InvalidVersion(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("x", model.NewBoolValue(true))

	opts := parsing.DefaultWriterOptions()
	opts.Ext["kdl-version"] = "3"
	w, err := newKDLWriter(opts)
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(val)
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
	if !strings.Contains(err.Error(), "unsupported output version") {
		t.Errorf("expected unsupported version error, got: %v", err)
	}
}

func TestWriter_V1RoundTrip(t *testing.T) {
	val := model.NewMapValue()
	_ = val.SetMapKey("name", model.NewStringValue("Bob"))
	_ = val.SetMapKey("active", model.NewBoolValue(true))
	_ = val.SetMapKey("count", model.NewIntValue(42))

	opts := parsing.DefaultWriterOptions()
	opts.Ext["kdl-version"] = "1"
	w, err := newKDLWriter(opts)
	if err != nil {
		t.Fatal(err)
	}

	data, err := w.Write(val)
	if err != nil {
		t.Fatal(err)
	}

	// Re-read and verify
	r := newTestReader(t)
	val2, err := r.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	assertMapStringValue(t, val2, "name", "Bob")
	assertMapBoolValue(t, val2, "active", true)
	assertMapIntValue(t, val2, "count", 42)
}
