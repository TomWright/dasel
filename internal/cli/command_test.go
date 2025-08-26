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
}
