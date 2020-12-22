package command

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal/selfupdate"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func mustAbs(path string) string {
	res, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return res
}

func TestRootCMD_Update(t *testing.T) {
	t.Run("Successful", updateTestOutputEqual("v1.0.0",
		func(httpClient *http.Client, user string, repo string, tag string) (*selfupdate.Release, error) {
			if exp, got := "TomWright", user; exp != got {
				return nil, fmt.Errorf("exp user %s, got %s", exp, got)
			}
			if exp, got := "dasel", repo; exp != got {
				return nil, fmt.Errorf("exp repo %s, got %s", exp, got)
			}
			if exp, got := "latest", tag; exp != got {
				return nil, fmt.Errorf("exp tag %s, got %s", exp, got)
			}
			return &selfupdate.Release{
				Assets: []*selfupdate.ReleaseAsset{
					{
						Name:               "dasel_macos_amd64",
						BrowserDownloadURL: "asd",
					},
				},
				TagName: "v1.1.0",
			}, nil
		},
		func(url string, dest string) error {
			if exp, got := "asd", url; exp != got {
				return fmt.Errorf("exp url %s, got %s", exp, got)
			}
			if exp, got := mustAbs("./dasel_macos_amd64"), dest; exp != got {
				return fmt.Errorf("exp dest %s, got %s", exp, got)
			}
			return nil
		},
		func(name string, mode os.FileMode) error {
			if exp, got := mustAbs("./dasel_macos_amd64"), name; exp != got {
				return fmt.Errorf("exp name %s, got %s", exp, got)
			}
			if exp, got := os.ModePerm, mode; exp != got {
				return fmt.Errorf("exp mode %s, got %s", exp, got)
			}
			return nil
		},
		func(name string, arg ...string) ([]byte, error) {
			if exp, got := mustAbs("./dasel_macos_amd64"), name; exp != got {
				return nil, fmt.Errorf("exp name %s, got %s", exp, got)
			}
			if exp, got := []string{"--version"}, arg; !reflect.DeepEqual(exp, got) {
				return nil, fmt.Errorf("exp args %v, got %s", exp, got)
			}
			return []byte(`dasel version v1.1.0`), nil
		},
		func() (string, error) {
			return "/current", nil
		},
		func(src string, dst string) error {
			if exp, got := mustAbs("./dasel_macos_amd64"), src; exp != got {
				return fmt.Errorf("exp src %s, got %s", exp, got)
			}
			if exp, got := "/current", dst; exp != got {
				return fmt.Errorf("exp dst %s, got %s", exp, got)
			}
			return nil
		},
		func(path string) error {
			if exp, got := mustAbs("./dasel_macos_amd64"), path; exp != got {
				return fmt.Errorf("exp path %s, got %s", exp, got)
			}
			return nil
		},
		`Updating...
Current version: v1.0.0
Release version: v1.1.0
New version: v1.1.0
Successfully updated
`, nil))

	t.Run("SuccessfulDevelopment", updateTestOutputEqual("development",
		func(httpClient *http.Client, user string, repo string, tag string) (*selfupdate.Release, error) {
			if exp, got := "TomWright", user; exp != got {
				return nil, fmt.Errorf("exp user %s, got %s", exp, got)
			}
			if exp, got := "dasel", repo; exp != got {
				return nil, fmt.Errorf("exp repo %s, got %s", exp, got)
			}
			if exp, got := "latest", tag; exp != got {
				return nil, fmt.Errorf("exp tag %s, got %s", exp, got)
			}
			return &selfupdate.Release{
				Assets: []*selfupdate.ReleaseAsset{
					{
						Name:               "dasel_macos_amd64",
						BrowserDownloadURL: "asd",
					},
				},
				TagName: "v1.1.0",
			}, nil
		},
		func(url string, dest string) error {
			if exp, got := "asd", url; exp != got {
				return fmt.Errorf("exp url %s, got %s", exp, got)
			}
			if exp, got := mustAbs("./dasel_macos_amd64"), dest; exp != got {
				return fmt.Errorf("exp dest %s, got %s", exp, got)
			}
			return nil
		},
		func(name string, mode os.FileMode) error {
			if exp, got := mustAbs("./dasel_macos_amd64"), name; exp != got {
				return fmt.Errorf("exp name %s, got %s", exp, got)
			}
			if exp, got := os.ModePerm, mode; exp != got {
				return fmt.Errorf("exp mode %s, got %s", exp, got)
			}
			return nil
		},
		func(name string, arg ...string) ([]byte, error) {
			if exp, got := mustAbs("./dasel_macos_amd64"), name; exp != got {
				return nil, fmt.Errorf("exp name %s, got %s", exp, got)
			}
			if exp, got := []string{"--version"}, arg; !reflect.DeepEqual(exp, got) {
				return nil, fmt.Errorf("exp args %v, got %s", exp, got)
			}
			return []byte(`dasel version v1.1.0`), nil
		},
		func() (string, error) {
			return "/current", nil
		},
		func(src string, dst string) error {
			if exp, got := mustAbs("./dasel_macos_amd64"), src; exp != got {
				return fmt.Errorf("exp src %s, got %s", exp, got)
			}
			if exp, got := "/current", dst; exp != got {
				return fmt.Errorf("exp dst %s, got %s", exp, got)
			}
			return nil
		},
		func(path string) error {
			if exp, got := mustAbs("./dasel_macos_amd64"), path; exp != got {
				return fmt.Errorf("exp path %s, got %s", exp, got)
			}
			return nil
		},
		`Updating...
Current version: development
Release version: v1.1.0
New version: v1.1.0
Successfully updated
`, nil, "--dev"))

	t.Run("SkipDevelopment", updateTestOutputEqual("development",
		nil, nil, nil, nil, nil, nil, nil,
		``, ErrIgnoredDev))
}

func updateTestOutputEqual(currentVersion string,
	fetchReleaseFn func(httpClient *http.Client, user string, repo string, tag string) (*selfupdate.Release, error),
	downloadFileFn func(url string, dest string) error,
	chmodFn func(name string, mode os.FileMode) error,
	executeCmdFn func(name string, arg ...string) ([]byte, error),
	executableFn func() (string, error),
	renameFn func(src string, dst string) error,
	removeFn func(path string) error,
	exp string, expErr error, additionalArgs ...string) func(t *testing.T) {

	return updateTestCheck(currentVersion, fetchReleaseFn, downloadFileFn, chmodFn, executeCmdFn, executableFn,
		renameFn, removeFn, func(out string) error {
			if exp != out {
				return fmt.Errorf("expected output %s, got %s", exp, out)
			}
			return nil
		}, expErr, additionalArgs...)
}

func updateTestCheck(currentVersion string,
	fetchReleaseFn func(httpClient *http.Client, user string, repo string, tag string) (*selfupdate.Release, error),
	downloadFileFn func(url string, dest string) error,
	chmodFn func(name string, mode os.FileMode) error,
	executeCmdFn func(name string, arg ...string) ([]byte, error),
	executableFn func() (string, error),
	renameFn func(src string, dst string) error,
	removeFn func(path string) error,
	checkFn func(out string) error, expErr error, additionalArgs ...string) func(t *testing.T) {

	updater := selfupdate.NewUpdater(currentVersion)

	updater.FetchReleaseFn = fetchReleaseFn
	updater.DownloadFileFn = downloadFileFn
	updater.ChmodFn = chmodFn
	updater.ExecuteCmdFn = executeCmdFn
	updater.ExecutableFn = executableFn
	updater.RenameFn = renameFn
	updater.RemoveFn = removeFn

	newRootCmd := func() *cobra.Command {
		root := NewRootCMD()

		for _, c := range root.Commands() {
			if c.Use == "update" {
				root.RemoveCommand(c)
			}
		}

		root.AddCommand(updateCommand(updater))

		return root
	}

	return func(t *testing.T) {
		cmd := newRootCmd()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"update",
		}
		if additionalArgs != nil {
			args = append(args, additionalArgs...)
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

		if expErr != nil || err != nil {
			return
		}

		output, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		if err := checkFn(string(output)); err != nil {
			t.Errorf("unexpected output: %s", err)
		}
	}
}
