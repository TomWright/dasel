package selfupdate

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	gitHubUsername = "TomWright"
	gitHubRepo     = "dasel"
)

// NewUpdater returns an Updater.
func NewUpdater(currentVersion string) *Updater {
	return &Updater{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		currentVersion: currentVersion,

		FetchReleaseFn: fetchGitHubRelease,
		DownloadFileFn: downloadFile,
		ChmodFn:        os.Chmod,
		ExecuteCmdFn:   executeCmd,
		ExecutableFn:   os.Executable,
		RenameFn:       os.Rename,
		RemoveFn:       os.Remove,
	}
}

// Release is a github release.
type Release struct {
	URL     string          `json:"url"`
	Assets  []*ReleaseAsset `json:"assets"`
	Name    string          `json:"name"`
	TagName string          `json:"tag_name"`
}

// FindAssetForSystem searches returns the asset for the given OS and arch.
func (r *Release) FindAssetForSystem(os string, arch string) *ReleaseAsset {
	ext := ""
	if os == "windows" {
		ext = ".exe"
	}
	matches := []string{fmt.Sprintf("dasel_%s_%s%s", os, arch, ext)}
	if os == "darwin" {
		matches = append(matches, fmt.Sprintf("dasel_%s_%s%s", "macos", arch, ext))
	}

	for _, a := range r.Assets {
		for _, m := range matches {
			if a.Name == m {
				return a
			}
		}
	}

	return nil
}

// Version returns the version of the release.
func (r *Release) Version() *Version {
	return versionFromString(r.TagName)
}

// ReleaseAsset is an asset of a Release.
type ReleaseAsset struct {
	URL                string `json:"url"`
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Updater provides functionality that allows you to update a binary from github.
type Updater struct {
	httpClient     *http.Client
	currentVersion string

	// FetchReleaseFn is used to fetch github release information.
	FetchReleaseFn func(httpClient *http.Client, user string, repo string, tag string) (*Release, error)
	// DownloadFileFn is used to download a file.
	DownloadFileFn func(url string, dest string) error
	// ChmodFn is used to change the permissions of a file.
	ChmodFn func(name string, mode os.FileMode) error
	// ExecuteCmdFn executes the given command and returns the output.
	ExecuteCmdFn func(name string, arg ...string) ([]byte, error)
	// ExecutableFn is used to return the path of the current executable.
	ExecutableFn func() (string, error)
	// RenameFn is used to rename a file.
	RenameFn func(src string, dst string) error
	// RemoveFn is used to remove a file.
	RemoveFn func(path string) error
}

// FindLatestRelease returns the latest release on the given repo.
func (u *Updater) FindLatestRelease() (*Release, error) {
	return u.FetchReleaseFn(u.httpClient, gitHubUsername, gitHubRepo, "latest")
}

// DownloadAsset downloads and returns a path to the given asset.
func (u *Updater) DownloadAsset(asset *ReleaseAsset) (string, error) {
	path, err := filepath.Abs(fmt.Sprintf("./%s", asset.Name))
	if err != nil {
		return "", err
	}

	if err := u.DownloadFileFn(asset.BrowserDownloadURL, path); err != nil {
		return "", fmt.Errorf("could not download file: %w", err)
	}

	if err := u.ChmodFn(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not make executable: %w", err)
	}

	return path, nil
}

// CurrentVersion returns the current version.
func (u *Updater) CurrentVersion() *Version {
	return versionFromString(u.currentVersion)
}

// GetVersion returns the version of the given executable.
func (u *Updater) GetVersion(path string) (*Version, error) {
	versionOutput, err := u.ExecuteCmdFn(path, "--version")
	if err != nil {
		return nil, fmt.Errorf("could not get new version: %w", err)
	}

	latest := versionFromString(string(versionOutput))

	return latest, nil
}

// Replace replaces the current executable with the given executable.
func (u *Updater) Replace(path string) error {
	currentPath, err := u.ExecutableFn()
	if err != nil {
		return fmt.Errorf("cannot get current executable path: %w", err)
	}

	if err := u.RenameFn(path, currentPath); err != nil {
		return fmt.Errorf("could not replace old executable: %w", err)
	}

	return nil
}

// CleanUp cleans up the given path.
// This should be deferred once a download has completed.
func (u *Updater) CleanUp(path string) {
	_ = u.RemoveFn(path)
}
