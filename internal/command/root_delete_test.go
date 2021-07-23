package command_test

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/internal/command"
	"io"
	"strings"
	"testing"
)

func TestRootCMD_Delete(t *testing.T) {
	t.Run("InvalidFile", expectErr(
		[]string{"delete", "-f", "bad.json", "-s", "x"},
		"could not open input file",
	))
	t.Run("MissingParser", expectErr(
		[]string{"delete", "-s", "x"},
		"parser flag required when reading from stdin",
	))

	t.Run("InvalidSingleSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"delete", "-p", "json", "-s", "[-]"},
		"invalid index: -",
	))
	t.Run("InvalidMultiSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"delete", "-p", "json", "-m", "-s", "[-]"},
		"invalid index: -",
	))
	t.Run("ValueNotFound", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"delete", "-p", "json", "-s", ".age"},
		"no value found for selector: .age",
	))
}

func deleteTest(in string, parser string, selector string, output string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return deleteTestCheck(in, parser, selector, func(out string) error {
		if out != output {
			return fmt.Errorf("expected %v, got %v", output, out)
		}
		return nil
	}, expErr, additionalArgs...)
}

func deleteTestCheck(in string, parser string, selector string, checkFn func(out string) error, expErr error, additionalArgs ...string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"delete", "-p", parser,
		}
		if additionalArgs != nil {
			args = append(args, additionalArgs...)
		}
		args = append(args, selector)

		cmd.SetOut(outputBuffer)
		cmd.SetIn(strings.NewReader(in))
		cmd.SetArgs(args)

		err := cmd.Execute()

		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}

		output, err := io.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		if err := checkFn(string(output)); err != nil {
			t.Errorf("unexpected output: %s", err)
		}
	}
}

func deleteTestFromFile(inputPath string, selector string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"select", "-f", inputPath, "-s", selector,
		}

		cmd.SetOut(outputBuffer)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}

		output, err := io.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		if out != string(output) {
			t.Errorf("expected result %v, got %v", out, string(output))
		}
	}
}

func TestRootCmd_Delete_JSON(t *testing.T) {
	t.Run("RootElement", deleteTest(`{
  "email": "tom@wright.com",
  "name": "Tom"
}`, "json", ".email", newline(`{
  "name": "Tom"
}`), nil))

	t.Run("Index", deleteTest(`{
  "colours": ["blue", "green", "red"],
  "name": "Tom"
}`, "json", ".colours.[1]", newline(`{
  "colours": [
    "blue",
    "red"
  ],
  "name": "Tom"
}`), nil))

	t.Run("AnyIndex", deleteTest(`{
  "colours": ["blue", "green", "red"],
  "name": "Tom"
}`, "json", ".colours.[*]", newline(`{
  "colours": [],
  "name": "Tom"
}`), nil, "-m"))

	t.Run("RootObject", deleteTest(`{
  "name": "Tom"
}`, "json", ".", newline(`{}`), nil))
	t.Run("RootObjectMulti", deleteTest(`{
  "name": "Tom"
}`, "json", ".", newline(`{}`), nil, "-m"))

	t.Run("RootArray", deleteTest(`[1, 2, 3]`, "json", ".", newline(`[]`), nil))
	t.Run("RootArrayMulti", deleteTest(`[1, 2, 3]`, "json", ".", newline(`[]`), nil, "-m"))

}
