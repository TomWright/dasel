package command

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

// Runs the dasel root command.
// Returns stdout, stderr and error.
func runDasel(args []string, in []byte) ([]byte, []byte, error) {
	stdOut := bytes.NewBuffer([]byte{})
	stdErr := bytes.NewBuffer([]byte{})

	cmd := NewRootCMD()
	cmd.SetArgs(args)
	cmd.SetOut(stdOut)
	cmd.SetErr(stdErr)

	if in != nil {
		cmd.SetIn(bytes.NewReader(in))
	}

	err := cmd.Execute()
	return stdOut.Bytes(), stdErr.Bytes(), err
}

func runTest(args []string, in []byte, expStdOut []byte, expStdErr []byte, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		if expStdOut == nil {
			expStdOut = []byte{}
		}
		if expStdErr == nil {
			expStdErr = []byte{}
		}

		gotStdOut, gotStdErr, gotErr := runDasel(args, in)
		if expErr != gotErr && !errors.Is(expErr, gotErr) {
			t.Errorf("expected error %v, got %v", expErr, gotErr)
			return
		}

		if !reflect.DeepEqual(expStdErr, gotStdErr) {
			t.Errorf("expected stderr %s, got %s", string(expStdErr), string(gotStdErr))
		}

		if !reflect.DeepEqual(expStdOut, gotStdOut) {
			t.Errorf("expected stdout %s, got %s", string(expStdOut), string(gotStdOut))
		}
	}
}

var newlineBytes = []byte("\n")

func newline(input []byte) []byte {
	return append(input, newlineBytes...)
}
