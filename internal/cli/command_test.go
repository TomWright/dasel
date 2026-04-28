package cli_test

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/tomwright/dasel/v3/internal/cli"
)

func runDasel(args []string, in []byte) ([]byte, []byte, error) {
	stdOut := bytes.NewBuffer([]byte{})
	stdErr := bytes.NewBuffer([]byte{})
	stdIn := bytes.NewReader(in)

	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = append([]string{"dasel", "query"}, args...)

	_, err := cli.Run(stdIn, stdOut, stdErr)

	return stdOut.Bytes(), stdErr.Bytes(), err
}

type testCase struct {
	args   []string
	in     []byte
	stdout []byte
	stderr []byte
	err    error
}

func runTest(tc testCase) func(t *testing.T) {
	return func(t *testing.T) {
		if tc.stdout == nil {
			tc.stdout = []byte{}
		}
		if tc.stderr == nil {
			tc.stderr = []byte{}
		}

		gotStdOut, gotStdErr, gotErr := runDasel(tc.args, tc.in)
		if !errors.Is(gotErr, tc.err) && !errors.Is(tc.err, gotErr) {
			t.Errorf("expected error %v, got %v", tc.err, gotErr)
			return
		}

		if !reflect.DeepEqual(tc.stderr, gotStdErr) {
			t.Errorf("expected stderr %s, got %s", string(tc.stderr), string(gotStdErr))
		}

		if !reflect.DeepEqual(tc.stdout, gotStdOut) {
			t.Errorf("expected stdout %s, got %s", string(tc.stdout), string(gotStdOut))
		}
	}
}

func TestRun(t *testing.T) {
	t.Run("complex set", func(t *testing.T) {
		t.Run("set nested with spread", runTest(testCase{
			args: []string{"-i", "json", "-o", "json", "--root", `user = {user..., name: {"first": $this.user.name, "last": "Doe"}}`},
			in:   []byte(`{"user": {"name": "John"}}`),
			stdout: []byte(`{
    "user": {
        "name": {
            "first": "John",
            "last": "Doe"
        }
    }
}
`),
			stderr: nil,
			err:    nil,
		}))
		t.Run("set nested", runTest(testCase{
			args: []string{"-i", "json", "-o", "json", "--root", `user.name = {"first": user.name, "last": "Doe"}`},
			in:   []byte(`{"user": {"name": "John"}}`),
			stdout: []byte(`{
    "user": {
        "name": {
            "first": "John",
            "last": "Doe"
        }
    }
}
`),
			stderr: nil,
			err:    nil,
		}))
		t.Run("set nested with localised group", runTest(testCase{
			args: []string{"-i", "json", "-o", "json", "--root", `user.(name = {"first": name, "last": "Doe"})`},
			in:   []byte(`{"user": {"name": "John"}}`),
			stdout: []byte(`{
    "user": {
        "name": {
            "first": "John",
            "last": "Doe"
        }
    }
}
`),
			stderr: nil,
			err:    nil,
		}))
		t.Run("set recursive descent", func(t *testing.T) {

			t.Run("property", runTest(testCase{
				args: []string{"-i", "json", "-o", "json", "--root", `$root..x.each($this = $this+1)`},
				in:   []byte(`[{"x":1},{"x":2},{"x":3}]`),
				stdout: []byte(`[
    {
        "x": 2
    },
    {
        "x": 3
    },
    {
        "x": 4
    }
]
`),
				stderr: nil,
				err:    nil,
			}))

			t.Run("index", runTest(testCase{
				args: []string{"-i", "json", "-o", "json", "--root", `$root..[1].each($this = $this+1)`},
				in:   []byte(`[ {"x":[1,2,3]} , {"y":[4,5,6]} , {"z":[7,8,9]} ]`),
				stdout: []byte(`[
    {
        "x": [
            1,
            3,
            3
        ]
    },
    {
        "y": [
            4,
            6,
            6
        ]
    },
    {
        "z": [
            7,
            9,
            9
        ]
    }
]
`),
				stderr: nil,
				err:    nil,
			}))

			t.Run("wildcard", runTest(testCase{
				args: []string{"-i", "json", "-o", "json", "--root", `$root..*.each($this = 4)`},
				in:   []byte(`[{"x":1},{"x":2},{"x":3}]`),
				stdout: []byte(`[
    {
        "x": 4
    },
    {
        "x": 4
    },
    {
        "x": 4
    }
]
`),
				stderr: nil,
				err:    nil,
			}))

		})
		t.Run("create object with empty stdin", runTest(testCase{
			args: []string{`{"name":"Tom"}`},
			in:   []byte{},
			stdout: []byte(`{
    "name": "Tom"
}
`),
			stderr: nil,
			err:    nil,
		}))
	})
	t.Run("set search", runTest(testCase{
		args: []string{"-i", "json", "-o", "json", "--root", `search(has("x")).each(x = x+1)`},
		in:   []byte(`[{"x":1},{"x":2},{"x":3}]`),
		stdout: []byte(`[
    {
        "x": 2
    },
    {
        "x": 3
    },
    {
        "x": 4
    }
]
`),
		stderr: nil,
		err:    nil,
	}))
	t.Run("recursive descent", func(t *testing.T) {
		t.Run("wildcard", runTest(testCase{
			args: []string{"-i", "json", `..*`},
			in: []byte(`{
  "user": {
    "name": "Alice",
    "roles": ["admin", "editor"],
    "meta": {
      "active": true,
      "score": 42
    }
  },
  "tags": ["x", "y"],
  "count": 10
}`),
			stdout: []byte(`[
    "Alice",
    "admin",
    "editor",
    true,
    42,
    "x",
    "y",
    10
]
`),
			stderr: nil,
			err:    nil,
		}))

		t.Run("property", runTest(testCase{
			args: []string{"-i", "json", `..name`},
			in: []byte(`{
  "user": {
    "name": "Alice",
    "roles": ["admin", "editor"],
    "meta": {
      "active": true,
      "score": 42
    }
  },
  "tags": ["x", "y"],
  "count": 10
}`),
			stdout: []byte(`[
    "Alice"
]
`),
			stderr: nil,
			err:    nil,
		}))

		t.Run("property2", runTest(testCase{
			args: []string{"-i", "json", `..name`},
			in:   []byte(`[{"name":"Tom"}, {"name":"Jim"}, {"foo": "Bar"}]`),
			stdout: []byte(`[
    "Tom",
    "Jim"
]
`),
			stderr: nil,
			err:    nil,
		}))

		t.Run("index", runTest(testCase{
			args: []string{"-i", "json", `..[0]`},
			in: []byte(`{
  "user": {
    "name": "Alice",
    "roles": ["admin", "editor"],
    "meta": {
      "active": true,
      "score": 42
    }
  },
  "tags": ["x", "y"],
  "count": 10
}`),
			stdout: []byte(`[
    "admin",
    "x"
]
`),
			stderr: nil,
			err:    nil,
		}))
	})
	t.Run("any", func(t *testing.T) {
		t.Run("true result", runTest(testCase{
			args:   []string{"-i", "json", `users.any(age > 25)`},
			in:     []byte(`{"users": [{"name": "Alice", "age": 20}, {"name": "Bob", "age": 30}]}`),
			stdout: []byte("true\n"),
			stderr: nil,
			err:    nil,
		}))
		t.Run("false result", runTest(testCase{
			args:   []string{"-i", "json", `users.any(age > 50)`},
			in:     []byte(`{"users": [{"name": "Alice", "age": 20}, {"name": "Bob", "age": 30}]}`),
			stdout: []byte("false\n"),
			stderr: nil,
			err:    nil,
		}))
	})
	t.Run("all", func(t *testing.T) {
		t.Run("true result", runTest(testCase{
			args:   []string{"-i", "json", `users.all(age > 18)`},
			in:     []byte(`{"users": [{"name": "Alice", "age": 20}, {"name": "Bob", "age": 30}]}`),
			stdout: []byte("true\n"),
			stderr: nil,
			err:    nil,
		}))
		t.Run("false result", runTest(testCase{
			args:   []string{"-i", "json", `users.all(age > 25)`},
			in:     []byte(`{"users": [{"name": "Alice", "age": 20}, {"name": "Bob", "age": 30}]}`),
			stdout: []byte("false\n"),
			stderr: nil,
			err:    nil,
		}))
	})
	t.Run("ternary", func(t *testing.T) {
		t.Run("true literal", runTest(testCase{
			args:   []string{`true ? "yes" : "no"`},
			in:     []byte{},
			stdout: []byte("\"yes\"\n"),
			err:    nil,
		}))
		t.Run("false literal", runTest(testCase{
			args:   []string{`false ? "yes" : "no"`},
			in:     []byte{},
			stdout: []byte("\"no\"\n"),
			err:    nil,
		}))
		t.Run("with json data true", runTest(testCase{
			args:   []string{"-i", "json", `age >= 18 ? "adult" : "minor"`},
			in:     []byte(`{"age": 25}`),
			stdout: []byte("\"adult\"\n"),
			err:    nil,
		}))
		t.Run("with json data false", runTest(testCase{
			args:   []string{"-i", "json", `age >= 18 ? "adult" : "minor"`},
			in:     []byte(`{"age": 10}`),
			stdout: []byte("\"minor\"\n"),
			err:    nil,
		}))
		t.Run("nested ternary", runTest(testCase{
			args:   []string{`true ? (false ? "a" : "b") : "c"`},
			in:     []byte{},
			stdout: []byte("\"b\"\n"),
			err:    nil,
		}))
		t.Run("with comparison", runTest(testCase{
			args:   []string{"-i", "json", `score >= 90 ? "A" : (score >= 80 ? "B" : "C")`},
			in:     []byte(`{"score": 85}`),
			stdout: []byte("\"B\"\n"),
			err:    nil,
		}))
		t.Run("ternary in map", runTest(testCase{
			args: []string{"-i", "json", `items.map(val >= 5 ? "big" : "small")`},
			in:   []byte(`{"items": [{"val": 10}, {"val": 3}, {"val": 7}]}`),
			stdout: []byte(`[
    "big",
    "small",
    "big"
]
`),
			err: nil,
		}))
		t.Run("ternary returns number", runTest(testCase{
			args:   []string{`true ? 42 : 0`},
			in:     []byte{},
			stdout: []byte("42\n"),
			err:    nil,
		}))
		t.Run("ternary with arithmetic", runTest(testCase{
			args:   []string{`true ? 1 + 2 : 3 + 4`},
			in:     []byte{},
			stdout: []byte("3\n"),
			err:    nil,
		}))
		t.Run("ternary with logical operators", runTest(testCase{
			args:   []string{"-i", "json", `a > 1 && b > 5 ? "both" : "nope"`},
			in:     []byte(`{"a": 5, "b": 10}`),
			stdout: []byte("\"both\"\n"),
			err:    nil,
		}))
		t.Run("chained property in condition", runTest(testCase{
			args:   []string{"-i", "json", `user.active ? "yes" : "no"`},
			in:     []byte(`{"user": {"active": true}}`),
			stdout: []byte("\"yes\"\n"),
			err:    nil,
		}))
		t.Run("ternary with json output", runTest(testCase{
			args: []string{"-i", "json", "-o", "json", `active ? {status: "on"} : {status: "off"}`},
			in:   []byte(`{"active": true}`),
			stdout: []byte(`{
    "status": "on"
}
`),
			err: nil,
		}))
		t.Run("ternary returns true", runTest(testCase{
			args:   []string{`true ? true : false`},
			in:     []byte{},
			stdout: []byte("true\n"),
			err:    nil,
		}))
	})
	t.Run("compact", func(t *testing.T) {
		t.Run("json", runTest(testCase{
			args:   []string{"-i", "json", "-o", "json", "--compact"},
			in:     []byte(`{"name": "Tom", "age": 30}`),
			stdout: []byte("{\"name\":\"Tom\",\"age\":30}\n"),
			err:    nil,
		}))
		t.Run("json nested", runTest(testCase{
			args:   []string{"-i", "json", "-o", "json", "--compact"},
			in:     []byte(`{"user": {"name": "Tom"}, "items": [1, 2, 3]}`),
			stdout: []byte("{\"user\":{\"name\":\"Tom\"},\"items\":[1,2,3]}\n"),
			err:    nil,
		}))
	})
	t.Run("merge", func(t *testing.T) {
		t.Run("deep", runTest(testCase{
			args: []string{"-i", "json", "-o", "json",
				`merge({"a": {"x": 1, "y": 2}}, {"a": {"y": 3, "z": 4}})`},
			in: []byte{},
			stdout: []byte("{\n    \"a\": {\n        \"x\": 1,\n        \"y\": 3,\n        \"z\": 4\n    }\n}\n"),
		}))
	})
	t.Run("count", func(t *testing.T) {
		t.Run("some match", runTest(testCase{
			args:   []string{"-i", "json", `users.count(age > 25)`},
			in:     []byte(`{"users": [{"name": "Alice", "age": 20}, {"name": "Bob", "age": 30}, {"name": "Charlie", "age": 35}]}`),
			stdout: []byte("2\n"),
			stderr: nil,
			err:    nil,
		}))
		t.Run("none match", runTest(testCase{
			args:   []string{"-i", "json", `users.count(age > 50)`},
			in:     []byte(`{"users": [{"name": "Alice", "age": 20}, {"name": "Bob", "age": 30}]}`),
			stdout: []byte("0\n"),
			stderr: nil,
			err:    nil,
		}))
	})
}
