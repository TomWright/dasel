package command_test

import (
	"bytes"
	"github.com/tomwright/dasel/internal/command"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestChangeDefaultCommand(t *testing.T) {
	cachedArgs := os.Args
	defer func() {
		os.Args = cachedArgs
	}()

	testArgs := func(in []string, exp []string, blacklistedArgs ...string) func(t *testing.T) {
		return func(t *testing.T) {
			os.Args = in

			cmd := command.NewRootCMD()
			command.ChangeDefaultCommand(cmd, "select", blacklistedArgs...)

			got := os.Args
			if !reflect.DeepEqual(exp, got) {
				t.Errorf("expected args %v, got %v", exp, got)
			}
		}
	}

	t.Run("ChangeToSelect", testArgs(
		[]string{"dasel", "-p", "json", ".name"},
		[]string{"dasel", "select", "-p", "json", ".name"},
	))

	t.Run("AlreadySelect", testArgs(
		[]string{"dasel", "select", "-p", "json", ".name"},
		[]string{"dasel", "select", "-p", "json", ".name"},
	))

	t.Run("AlreadyPut", testArgs(
		[]string{"dasel", "put", "-p", "json", "-t", "string", "name=Tom"},
		[]string{"dasel", "put", "-p", "json", "-t", "string", "name=Tom"},
	))

	t.Run("IgnoreBlacklisted", testArgs(
		[]string{"dasel", "-v"},
		[]string{"dasel", "-v"},
		"-v",
	))

	t.Run("IgnoreBlacklisted", testArgs(
		[]string{"dasel", "select", "-v"},
		[]string{"dasel", "select", "-v"},
		"-v",
	))
}

func expectErr(args []string, expErr string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		cmd.SetOut(outputBuffer)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if err == nil || !strings.Contains(err.Error(), expErr) {
			t.Errorf("unexpected error: %v", err)
			return
		}
	}
}

func expectErrFromInput(in string, args []string, expErr string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		cmd.SetIn(bytes.NewReader([]byte(in)))
		cmd.SetOut(outputBuffer)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if err == nil || !strings.Contains(err.Error(), expErr) {
			t.Errorf("unexpected error: %v: %s", err, outputBuffer.String())
			return
		}
	}
}

func expectOutput(in string, args []string, exp string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		cmd.SetIn(bytes.NewReader([]byte(in)))
		cmd.SetOut(outputBuffer)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got := outputBuffer.String()
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func expectOutputAndErr(args []string, expErr string, expOutput string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		cmd.SetOut(outputBuffer)
		cmd.SetArgs(args)

		err := cmd.Execute()

		gotErr := ""
		if err != nil {
			gotErr = err.Error()
		}

		if expErr != gotErr {
			t.Errorf("expected err %s, got %s", expErr, gotErr)
		}

		gotOutput := outputBuffer.String()
		if expOutput != gotOutput {
			t.Errorf("expected:\n%s\ngot:\n%s", expOutput, gotOutput)
		}
	}
}
