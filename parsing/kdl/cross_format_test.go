package kdl_test

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/json"
	"github.com/tomwright/dasel/v3/parsing/kdl"
	"github.com/tomwright/dasel/v3/parsing/yaml"
)

func mustReader(t *testing.T, f parsing.Format) parsing.Reader {
	t.Helper()
	r, err := f.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func mustWriter(t *testing.T, f parsing.Format) parsing.Writer {
	t.Helper()
	w, err := f.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatal(err)
	}
	return w
}

// --- KDL → JSON ---

func TestKDLToJSON_Scalars(t *testing.T) {
	kdlData := []byte(`name "Bob"
age 76
active true
score 9.5`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"name": "Bob"`)
	assertContains(t, result, `"age": 76`)
	assertContains(t, result, `"active": true`)
	assertContains(t, result, `"score": 9.5`)
}

func TestKDLToJSON_NestedChildren(t *testing.T) {
	kdlData := []byte(`server {
    host "localhost"
    port 8080
}`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"server"`)
	assertContains(t, result, `"host": "localhost"`)
	assertContains(t, result, `"port": 8080`)
}

func TestKDLToJSON_DuplicateNodes(t *testing.T) {
	kdlData := []byte(`plugin "git"
plugin "docker"
plugin "tmux"`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"plugin"`)
	assertContains(t, result, `"git"`)
	assertContains(t, result, `"docker"`)
	assertContains(t, result, `"tmux"`)
}

func TestKDLToJSON_NullAndBool(t *testing.T) {
	kdlData := []byte(`enabled #true
disabled #false
empty #null`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"enabled": true`)
	assertContains(t, result, `"disabled": false`)
	assertContains(t, result, `"empty": null`)
}

func TestKDLToJSON_NodeWithArgsAndProps(t *testing.T) {
	kdlData := []byte(`server 80 host="localhost" {
    tls #true
}`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"server"`)
	assertContains(t, result, `"$args"`)
	assertContains(t, result, `80`)
	assertContains(t, result, `"host": "localhost"`)
	assertContains(t, result, `"tls": true`)
}

func TestKDLToJSON_EmptyNode(t *testing.T) {
	kdlData := []byte(`marker`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"marker": null`)
}

func TestKDLToJSON_DeeplyNested(t *testing.T) {
	kdlData := []byte(`a {
    b {
        c "deep"
    }
}`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"a"`)
	assertContains(t, result, `"b"`)
	assertContains(t, result, `"c": "deep"`)
}

func TestKDLToJSON_NumberTypes(t *testing.T) {
	kdlData := []byte(`hex 0xff
octal 0o77
binary 0b1010
big 1_000_000`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"hex": 255`)
	assertContains(t, result, `"octal": 63`)
	assertContains(t, result, `"binary": 10`)
	assertContains(t, result, `"big": 1000000`)
}

// --- JSON → KDL ---

func TestJSONToKDL_SimpleObject(t *testing.T) {
	jsonData := []byte(`{
    "name": "Bob",
    "age": 76,
    "active": true
}
`)

	val, err := mustReader(t, json.JSON).Read(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(kdlData)
	assertContains(t, result, `name "Bob"`)
	assertContains(t, result, `age 76`)
	assertContains(t, result, `active #true`)
}

func TestJSONToKDL_NestedObject(t *testing.T) {
	jsonData := []byte(`{
    "server": {
        "host": "localhost",
        "port": 8080
    }
}
`)

	val, err := mustReader(t, json.JSON).Read(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(kdlData)
	assertContains(t, result, "server")
	assertContains(t, result, `host="localhost"`)
	assertContains(t, result, `port=8080`)
}

func TestJSONToKDL_Array(t *testing.T) {
	jsonData := []byte(`{
    "plugins": [
        "git",
        "docker"
    ]
}
`)

	val, err := mustReader(t, json.JSON).Read(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(kdlData)
	if strings.Count(result, "plugins") != 2 {
		t.Errorf("expected 2 'plugins' nodes, got:\n%s", result)
	}
	assertContains(t, result, `"git"`)
	assertContains(t, result, `"docker"`)
}

func TestJSONToKDL_NullAndBool(t *testing.T) {
	jsonData := []byte(`{
    "enabled": true,
    "disabled": false,
    "empty": null
}
`)

	val, err := mustReader(t, json.JSON).Read(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(kdlData)
	assertContains(t, result, "enabled #true")
	assertContains(t, result, "disabled #false")
	assertContains(t, result, "empty")
}

// --- KDL → YAML ---

func TestKDLToYAML_Scalars(t *testing.T) {
	kdlData := []byte(`name "Bob"
age 76
active true`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	yamlData, err := mustWriter(t, yaml.YAML).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(yamlData)
	assertContains(t, result, "name: Bob")
	assertContains(t, result, "age: 76")
	assertContains(t, result, "active: true")
}

func TestKDLToYAML_Nested(t *testing.T) {
	kdlData := []byte(`database {
    host "db.example.com"
    port 5432
}`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	yamlData, err := mustWriter(t, yaml.YAML).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(yamlData)
	assertContains(t, result, "database:")
	assertContains(t, result, "host: db.example.com")
	assertContains(t, result, "port: 5432")
}

// --- YAML → KDL ---

func TestYAMLToKDL_SimpleMap(t *testing.T) {
	yamlData := []byte(`name: Bob
age: 76
active: true
`)

	val, err := mustReader(t, yaml.YAML).Read(yamlData)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(kdlData)
	assertContains(t, result, `name "Bob"`)
	assertContains(t, result, `age 76`)
	assertContains(t, result, `active #true`)
}

func TestYAMLToKDL_NestedMap(t *testing.T) {
	yamlData := []byte(`database:
  host: db.example.com
  port: 5432
`)

	val, err := mustReader(t, yaml.YAML).Read(yamlData)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(kdlData)
	assertContains(t, result, "database")
	assertContains(t, result, `host="db.example.com"`)
	assertContains(t, result, `port=5432`)
}

// --- Round-trips via intermediate model ---

func TestRoundTrip_KDLToJSONToKDL(t *testing.T) {
	kdlInput := []byte(`name "Bob"
age 76
active #true`)

	// KDL → model
	val, err := mustReader(t, kdl.KDL).Read(kdlInput)
	if err != nil {
		t.Fatal(err)
	}

	// model → JSON
	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	// JSON → model
	val2, err := mustReader(t, json.JSON).Read(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	// model → KDL
	kdlOutput, err := mustWriter(t, kdl.KDL).Write(val2)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the KDL output can be re-read and has same values
	val3, err := mustReader(t, kdl.KDL).Read(kdlOutput)
	if err != nil {
		t.Fatal(err)
	}

	assertModelString(t, val3, "name", "Bob")
	assertModelInt(t, val3, "age", 76)
	assertModelBool(t, val3, "active", true)
}

func TestRoundTrip_JSONToKDLToJSON(t *testing.T) {
	jsonInput := []byte(`{
    "name": "Alice",
    "count": 42,
    "enabled": false
}
`)

	// JSON → model
	val, err := mustReader(t, json.JSON).Read(jsonInput)
	if err != nil {
		t.Fatal(err)
	}

	// model → KDL
	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	// KDL → model
	val2, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	// model → JSON
	jsonOutput, err := mustWriter(t, json.JSON).Write(val2)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonOutput)
	assertContains(t, result, `"name": "Alice"`)
	assertContains(t, result, `"count": 42`)
	assertContains(t, result, `"enabled": false`)
}

func TestRoundTrip_YAMLToKDLToYAML(t *testing.T) {
	yamlInput := []byte(`name: Alice
count: 42
enabled: false
`)

	val, err := mustReader(t, yaml.YAML).Read(yamlInput)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	val2, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	yamlOutput, err := mustWriter(t, yaml.YAML).Write(val2)
	if err != nil {
		t.Fatal(err)
	}

	result := string(yamlOutput)
	assertContains(t, result, "name: Alice")
	assertContains(t, result, "count: 42")
	assertContains(t, result, "enabled: false")
}

func TestRoundTrip_KDLToYAMLToKDL(t *testing.T) {
	kdlInput := []byte(`title "My App"
version 3
debug #false`)

	val, err := mustReader(t, kdl.KDL).Read(kdlInput)
	if err != nil {
		t.Fatal(err)
	}

	yamlData, err := mustWriter(t, yaml.YAML).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	val2, err := mustReader(t, yaml.YAML).Read(yamlData)
	if err != nil {
		t.Fatal(err)
	}

	kdlOutput, err := mustWriter(t, kdl.KDL).Write(val2)
	if err != nil {
		t.Fatal(err)
	}

	val3, err := mustReader(t, kdl.KDL).Read(kdlOutput)
	if err != nil {
		t.Fatal(err)
	}

	assertModelString(t, val3, "title", "My App")
	assertModelInt(t, val3, "version", 3)
	assertModelBool(t, val3, "debug", false)
}

// --- Real-world config conversions ---

func TestKDLToJSON_ZellijLikeConfig(t *testing.T) {
	kdlData := []byte(`keybinds {
    normal {
        bind "Ctrl+h" {
            action "MoveFocus" "Left"
        }
    }
}
theme "dracula"
pane_frames true`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := mustWriter(t, json.JSON).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(jsonData)
	assertContains(t, result, `"keybinds"`)
	assertContains(t, result, `"theme": "dracula"`)
	assertContains(t, result, `"pane_frames": true`)
}

func TestJSONToKDL_PackageConfig(t *testing.T) {
	jsonData := []byte(`{
    "name": "my-app",
    "version": "1.0.0",
    "private": true,
    "dependencies": {
        "react": "^18.0.0",
        "typescript": "^5.0.0"
    }
}
`)

	val, err := mustReader(t, json.JSON).Read(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	kdlData, err := mustWriter(t, kdl.KDL).Write(val)
	if err != nil {
		t.Fatal(err)
	}

	result := string(kdlData)
	assertContains(t, result, `name "my-app"`)
	assertContains(t, result, `version "1.0.0"`)
	assertContains(t, result, "private #true")
	assertContains(t, result, "dependencies")
}

// --- CLI-level cross-format (KDL added to existing patterns) ---

func TestKDLToJSON_SelectScalar(t *testing.T) {
	kdlData := []byte(`greeting "hello"
count 42`)

	val, err := mustReader(t, kdl.KDL).Read(kdlData)
	if err != nil {
		t.Fatal(err)
	}

	// Query-like: access a specific key
	v, err := val.GetMapKey("greeting")
	if err != nil {
		t.Fatal(err)
	}
	s, err := v.StringValue()
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello" {
		t.Errorf("expected 'hello', got %q", s)
	}

	v2, err := val.GetMapKey("count")
	if err != nil {
		t.Fatal(err)
	}
	n, err := v2.IntValue()
	if err != nil {
		t.Fatal(err)
	}
	if n != 42 {
		t.Errorf("expected 42, got %d", n)
	}
}

// --- Helpers ---

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected output to contain %q, got:\n%s", needle, haystack)
	}
}

func assertModelString(t *testing.T, val *model.Value, key, expected string) {
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

func assertModelInt(t *testing.T, val *model.Value, key string, expected int64) {
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

func assertModelBool(t *testing.T, val *model.Value, key string, expected bool) {
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
