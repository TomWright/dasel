package selfupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func executeCmd(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}

func fetchGitHubRelease(httpClient *http.Client, user string, repo string, tag string) (*Release, error) {
	if tag != "latest" {
		tag = fmt.Sprintf("tags/%s", tag)
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s", user, repo, tag)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request to get latest release: %w", err)
	}

	resp, err := httpClient.Do(req)
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

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
