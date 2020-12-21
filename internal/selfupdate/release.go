package selfupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Updater allows you to update a binary file from a latest github release.
type Updater interface {
	// FindLatestRelease returns the latest release on the given repo.
	FindLatestRelease(user string, repo string) (*Release, error)
	// DownloadAssetWithName downloads and returns a path to the asset with the given name.
	DownloadAssetWithName(release *Release, name string) (string, error)
	// DownloadAsset downloads and returns a path to the given asset.
	DownloadAsset(asset *ReleaseAsset) (string, error)
	// GetVersions returns the current version and the version of the given executable.
	GetVersions(path string) (*Version, *Version, error)
	// Replace replaces the current executable with the given executable.
	Replace(path string) error
}

// NewUpdater returns an Updater.
func NewUpdater(currentVersion string) Updater {
	return &stdUpdater{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		currentVersion: currentVersion,
	}
}

// Release is a github release.
type Release struct {
	URL     string          `json:"url"`
	Assets  []*ReleaseAsset `json:"assets"`
	Name    string          `json:"name"`
	TagName string          `json:"tag_name"`
}

// ReleaseAsset is an asset of a Release.
type ReleaseAsset struct {
	URL                string `json:"url"`
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type stdUpdater struct {
	httpClient     *http.Client
	currentVersion string
}

// FindLatestRelease returns the latest release on the given repo.
func (u *stdUpdater) FindLatestRelease(user string, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", user, repo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request to get latest release: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform request to get latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response bytes: %w", err)
	}

	release := &Release{}
	if err := json.Unmarshal(respBytes, release); err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	return release, nil
}

// DownloadAssetWithName downloads and returns a path to the asset with the given name.
func (u *stdUpdater) DownloadAssetWithName(release *Release, name string) (string, error) {
	for _, a := range release.Assets {
		if a.Name == name {
			return u.DownloadAsset(a)
		}
	}
	return "", fmt.Errorf("no asset found with name: %s", name)
}

// DownloadAsset downloads and returns a path to the given asset.
func (u *stdUpdater) DownloadAsset(asset *ReleaseAsset) (string, error) {
	path, err := filepath.Abs(fmt.Sprintf("./%s", asset.Name))
	if err != nil {
		return "", err
	}

	// Get the data
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	if err := os.Chmod(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not make executable: %w", err)
	}

	return path, nil
}

// GetVersions returns the current version and the version of the given executable.
func (u *stdUpdater) GetVersions(path string) (*Version, *Version, error) {
	versionOutput, err := exec.Command(path, "--version").Output()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get new version: %w", err)
	}

	current := versionFromString(u.currentVersion)
	latest := versionFromString(string(versionOutput))

	return current, latest, nil
}

// Replace replaces the current executable with the given executable.
func (u *stdUpdater) Replace(path string) error {
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get current executable path: %w", err)
	}

	if err := os.Rename(path, currentPath); err != nil {
		return fmt.Errorf("could not replace old executable: %w", err)
	}

	return nil
}
