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
	t.Run("InvalidRootNode", expectErrFromInput(
		``,
		[]string{"delete", "-p", "json", "-s", ".age"},
		"no value found for selector: .age",
	))
	t.Run("InvalidRootNodeMulti", expectErrFromInput(
		``,
		[]string{"delete", "-p", "json", "-m", "-s", ".age"},
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

func TestRootCmd_Delete_JSON(t *testing.T) {
	t.Run("RootNodeObject", deleteTest(`{
  "email": "tom@wright.com",
  "name": "Tom"
}`, "json", ".", newline(`{}`), nil))

	t.Run("RootNodeObjectMulti", deleteTest(`{
  "email": "tom@wright.com",
  "name": "Tom"
}`, "json", ".", newline(`{}`), nil, "-m"))

	t.Run("RootNodeArray", deleteTest(`["a", "b", "c"]`, "json", ".",
		newline(`[]`), nil))

	t.Run("RootNodeArrayMulti", deleteTest(`["a", "b", "c"]`, "json", ".",
		newline(`[]`), nil, "-m"))

	t.Run("RootNodeUnknown", deleteTest(`false`, "json", ".",
		newline(`{}`), nil))

	t.Run("RootNodeUnknownMulti", deleteTest(`false`, "json", ".",
		newline(`{}`), nil, "-m"))

	t.Run("Property", deleteTest(`{
  "email": "tom@wright.com",
  "name": "Tom"
}`, "json", ".email", newline(`{
  "name": "Tom"
}`), nil))

	t.Run("PropertyCompact", deleteTest(`{
  "email": "tom@wright.com",
  "name": "Tom"
}`, "json", ".email", newline(`{"name":"Tom"}`), nil, "-c"))

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

	t.Run("DeleteStringEscapeHTMLOn", deleteTest(`{
  "name": "Tom",
  "user": "Tom <contact@tomwright.me>"
}
`, "json", `.name`, `{
  "user": "Tom \u003ccontact@tomwright.me\u003e"
}
`, nil, "--escape-html=true"))

	t.Run("DeleteStringEscapeHTMLOff", deleteTest(`{
  "name": "Tom",
  "user": "Tom <contact@tomwright.me>"
}
`, "json", `.name`, `{
  "user": "Tom <contact@tomwright.me>"
}
`, nil, "--escape-html=false"))
}
