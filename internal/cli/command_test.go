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
}
