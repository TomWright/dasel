package cli_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/tomwright/dasel/v3/internal/cli"
)

func runDasel(args []string, in []byte) ([]byte, []byte, error) {
	stdOut := bytes.NewBuffer([]byte{})
	stdErr := bytes.NewBuffer([]byte{})

	cmd := cli.RootCmd()
	cmd.SetArgs(args)
	cmd.SetOut(stdOut)
	cmd.SetErr(stdErr)

	if in != nil {
		cmd.SetIn(bytes.NewReader(in))
	}

	err := cmd.Execute()
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
