package dasel_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func runInstallScript(t *testing.T, targetOS, targetArch string, args ...string) (string, string, int) {
	t.Helper()
	cmdArgs := append([]string{"./install.sh"}, args...)
	cmd := exec.Command("/bin/sh", cmdArgs...)
	cmd.Dir = "."
	cmd.Env = append(cmd.Environ(),
		"TARGET_OS="+targetOS,
		"TARGET_ARCH="+targetArch,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run install.sh: %v", err)
		}
	}

	return stdout.String(), stderr.String(), exitCode
}

func TestInstallScript_ArchitectureDetection(t *testing.T) {
	testCases := []struct {
		name       string
		targetOS   string
		targetArch string
		expected   string
	}{
		{
			name:       "linux aarch64",
			targetOS:   "Linux",
			targetArch: "aarch64",
			expected:   "dasel_linux_arm64",
		},
		{
			name:       "linux arm64",
			targetOS:   "Linux",
			targetArch: "arm64",
			expected:   "dasel_linux_arm64",
		},
		{
			name:       "linux x86_64",
			targetOS:   "Linux",
			targetArch: "x86_64",
			expected:   "dasel_linux_amd64",
		},
		{
			name:       "linux amd64",
			targetOS:   "Linux",
			targetArch: "amd64",
			expected:   "dasel_linux_amd64",
		},
		{
			name:       "linux armv7l",
			targetOS:   "Linux",
			targetArch: "armv7l",
			expected:   "dasel_linux_arm32",
		},
		{
			name:       "linux arm",
			targetOS:   "Linux",
			targetArch: "arm",
			expected:   "dasel_linux_arm32",
		},
		{
			name:       "linux i386",
			targetOS:   "Linux",
			targetArch: "i386",
			expected:   "dasel_linux_386",
		},
		{
			name:       "linux i686",
			targetOS:   "Linux",
			targetArch: "i686",
			expected:   "dasel_linux_386",
		},
		{
			name:       "darwin arm64",
			targetOS:   "Darwin",
			targetArch: "arm64",
			expected:   "dasel_darwin_arm64",
		},
		{
			name:       "darwin aarch64",
			targetOS:   "Darwin",
			targetArch: "aarch64",
			expected:   "dasel_darwin_arm64",
		},
		{
			name:       "darwin x86_64",
			targetOS:   "Darwin",
			targetArch: "x86_64",
			expected:   "dasel_darwin_amd64",
		},
		{
			name:       "windows x86_64",
			targetOS:   "MINGW64_NT-10.0",
			targetArch: "x86_64",
			expected:   "dasel_windows_amd64.exe",
		},
		{
			name:       "windows i686",
			targetOS:   "MSYS_NT-10.0",
			targetArch: "i686",
			expected:   "dasel_windows_386.exe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := runInstallScript(t, tc.targetOS, tc.targetArch, "--target")
			if exitCode != 0 {
				t.Fatalf("unexpected exit code %d, stderr: %s", exitCode, stderr)
			}
			actual := strings.TrimSpace(stdout)
			if actual != tc.expected {
				t.Errorf("expected target %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestInstallScript_DryRunDetails(t *testing.T) {
	stdout, stderr, exitCode := runInstallScript(t, "Linux", "aarch64", "--dry-run")
	if exitCode != 0 {
		t.Fatalf("unexpected exit code %d, stderr: %s", exitCode, stderr)
	}
	out := stdout
	if !strings.Contains(out, "Target Asset: dasel_linux_arm64") {
		t.Errorf("dry-run missing target asset, got: %s", out)
	}
	if !strings.Contains(out, "Architecture: arm64") {
		t.Errorf("dry-run missing resolved arch, got: %s", out)
	}
	if !strings.Contains(out, "OS: linux") {
		t.Errorf("dry-run missing resolved os, got: %s", out)
	}
	if !strings.Contains(out, "Download URL: https://github.com/TomWright/dasel/releases/latest/download/dasel_linux_arm64") {
		t.Errorf("dry-run missing correct download url, got: %s", out)
	}
}

func TestInstallScript_VersionOption(t *testing.T) {
	stdout, stderr, exitCode := runInstallScript(t, "Linux", "aarch64", "--dry-run", "--version", "v3.11.2")
	if exitCode != 0 {
		t.Fatalf("unexpected exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Version: v3.11.2") {
		t.Errorf("dry-run missing version, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Download URL: https://github.com/TomWright/dasel/releases/download/v3.11.2/dasel_linux_arm64") {
		t.Errorf("dry-run missing versioned download url, got: %s", stdout)
	}
}

func TestInstallScript_BindirOption(t *testing.T) {
	stdout, stderr, exitCode := runInstallScript(t, "Linux", "aarch64", "--dry-run", "--bindir", "/opt/bin")
	if exitCode != 0 {
		t.Fatalf("unexpected exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Install Directory: /opt/bin") {
		t.Errorf("dry-run missing custom bindir, got: %s", stdout)
	}
}

func TestInstallScript_UnsupportedArch(t *testing.T) {
	_, stderr, exitCode := runInstallScript(t, "Linux", "mips64", "--target")
	if exitCode == 0 {
		t.Errorf("expected failure for unsupported arch, got exit code 0")
	}
	if !strings.Contains(stderr, "Unsupported architecture") {
		t.Errorf("expected unsupported architecture message, got: %s", stderr)
	}
}

func TestInstallScript_UnsupportedOS(t *testing.T) {
	_, stderr, exitCode := runInstallScript(t, "Solaris", "amd64", "--target")
	if exitCode == 0 {
		t.Errorf("expected failure for unsupported OS, got exit code 0")
	}
	if !strings.Contains(stderr, "Unsupported operating system") {
		t.Errorf("expected unsupported OS message, got: %s", stderr)
	}
}

func TestInstallScript_Help(t *testing.T) {
	stdout, _, exitCode := runInstallScript(t, "Linux", "x86_64", "--help")
	if exitCode != 0 {
		t.Errorf("expected exit code 0 for help, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage: install.sh") {
		t.Errorf("expected help output, got: %s", stdout)
	}
}
