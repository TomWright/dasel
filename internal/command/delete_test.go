package command

import (
	"testing"
)

func TestDeleteCommand(t *testing.T) {

	t.Run("DeleteMapField", runTest(
		[]string{"delete", "-r", "json", "--pretty=false", "x"},
		[]byte(`{"x":1,"y":2}`),
		newline([]byte(`{"y":2}`)),
		nil,
		nil,
	))

	t.Run("DeleteNestedMapField", runTest(
		[]string{"delete", "-r", "json", "--pretty=false", "x.y"},
		[]byte(`{"x":{"x":1,"y":2},"y":{"x":1,"y":2}}`),
		newline([]byte(`{"x":{"x":1},"y":{"x":1,"y":2}}`)),
		nil,
		nil,
	))

	t.Run("DeleteSliceIndex", runTest(
		[]string{"delete", "-r", "json", "--pretty=false", "[1]"},
		[]byte(`[0,1,2]`),
		newline([]byte(`[0,2]`)),
		nil,
		nil,
	))

	t.Run("DeletedNestedSliceIndex", runTest(
		[]string{"delete", "-r", "json", "--pretty=false", "users.[1]"},
		[]byte(`{"users":[0,1,2]}`),
		newline([]byte(`{"users":[0,2]}`)),
		nil,
		nil,
	))

	t.Run("CheckIndentionForJSON", runTest(
		[]string{"delete", "-r", "json", "--indent", "6", "--pretty=true", "x.y"},
		[]byte(`{"x":{"x":1,"y":2}}`),
		newline([]byte("{\n      \"x\": {\n            \"x\": 1\n      }\n}")),
		nil,
		nil,
	))

	t.Run("CheckIndentionForYAML", runTest(
		[]string{"delete", "-r", "json", "-w", "yaml", "--indent", "6", "--pretty=true", "x.y"},
		[]byte(`{"x":{"x":1,"y":2}}`),
		newline([]byte("x:\n      x: 1")),
		nil,
		nil,
	))

	t.Run("CheckIndentionForTOML", runTest(
		[]string{"delete", "-r", "json", "-w", "toml", "--indent", "6", "--pretty=true", "x.y"},
		[]byte(`{"x":{"x":1,"y":2}}`),
		newline([]byte("[x]\n      x = 1")),
		nil,
		nil,
	))
}
