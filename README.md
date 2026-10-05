[![Gitbook](https://badges.aleen42.com/src/gitbook_1.svg)](https://daseldocs.tomwright.me)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomwright/dasel/v3)](https://goreportcard.com/report/github.com/tomwright/dasel/v3)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel/v3)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![GitHub Downloads](https://img.shields.io/github/downloads/TomWright/dasel/total)
![Homebrew Formula Downloads](https://img.shields.io/homebrew/installs/dy/dasel?label=brew%20installs)
![GitHub License](https://img.shields.io/github/license/TomWright/dasel)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/TomWright/dasel?label=latest%20release)](https://github.com/TomWright/dasel/releases/latest)
[![Homebrew tag (latest by date)](https://img.shields.io/homebrew/v/dasel)](https://formulae.brew.sh/formula/dasel)

<div align="center">
    <img src="./daselgopher.png" alt="Dasel mascot" width="250"/>
</div>

# Dasel

Dasel (short for **Data-Select**) is a command-line tool and library for querying, modifying, and transforming data structures such as JSON, YAML, TOML, XML, CSV, and KDL.

It provides a consistent, powerful syntax to traverse and update data — making it useful for developers, DevOps, and data wrangling tasks.

---

## Features

* **Multi-format support**: JSON, YAML, TOML, XML, CSV, HCL, INI, KDL.
* **Unified query syntax**: Access data in any format with the same selectors.
* **Query & search**: Extract values, lists, or structures with intuitive syntax.
* **Modify in place**: Update, insert, or delete values directly in structured files.
* **Convert between formats**: Seamlessly transform data from JSON → YAML, TOML → JSON, etc.
* **Script-friendly**: Simple CLI integration for shell scripts and pipelines.
* **Library support**: Import and use in Go projects.

---

## Installation

### Automated Script

Download and install the appropriate binary for your system architecture:

```sh
curl -sSLf https://raw.githubusercontent.com/TomWright/dasel/master/install.sh | sh
```

Pass options through environment variables:

```sh
# Custom installation directory
BINDIR=~/.local/bin curl -sSLf https://raw.githubusercontent.com/TomWright/dasel/master/install.sh | sh

# Specific release version
VERSION=v3.11.2 curl -sSLf https://raw.githubusercontent.com/TomWright/dasel/master/install.sh | sh
```

### Manual Binary Download

Download prebuilt binaries directly from GitHub releases using `curl`:

#### Linux (aarch64 / arm64, e.g. Amazon Linux 2023, Ubuntu ARM)

```sh
curl -sSLf https://github.com/TomWright/dasel/releases/latest/download/dasel_linux_arm64 -o dasel && chmod +x dasel
sudo mv ./dasel /usr/local/bin/dasel
```

Using the GitHub API:

```sh
curl -sSLf "$(curl -sSLf https://api.github.com/repos/TomWright/dasel/releases/latest | grep browser_download_url | grep linux_arm64 | grep -v '\.gz' | cut -d\" -f 4)" -L -o dasel && chmod +x dasel
sudo mv ./dasel /usr/local/bin/dasel
```

#### Linux (x86_64 / amd64)

```sh
curl -sSLf https://github.com/TomWright/dasel/releases/latest/download/dasel_linux_amd64 -o dasel && chmod +x dasel
sudo mv ./dasel /usr/local/bin/dasel
```

Using the GitHub API:

```sh
curl -sSLf "$(curl -sSLf https://api.github.com/repos/TomWright/dasel/releases/latest | grep browser_download_url | grep linux_amd64 | grep -v '\.gz' | cut -d\" -f 4)" -L -o dasel && chmod +x dasel
sudo mv ./dasel /usr/local/bin/dasel
```

#### macOS (Apple Silicon / arm64)

```sh
curl -sSLf https://github.com/TomWright/dasel/releases/latest/download/dasel_darwin_arm64 -o dasel && chmod +x dasel
sudo mv ./dasel /usr/local/bin/dasel
```

#### macOS (Intel / amd64)

```sh
curl -sSLf https://github.com/TomWright/dasel/releases/latest/download/dasel_darwin_amd64 -o dasel && chmod +x dasel
sudo mv ./dasel /usr/local/bin/dasel
```

### Release Asset Target Mapping

Each release provides standalone binaries and gzip archives for supported operating systems and CPU architectures:

| Operating System | Architecture | Kernel / Machine (`uname -m`) | Release Binary Asset |
| :--- | :--- | :--- | :--- |
| Linux | ARM64 | `aarch64`, `arm64` | `dasel_linux_arm64` |
| Linux | AMD64 | `x86_64`, `amd64` | `dasel_linux_amd64` |
| Linux | ARM32 | `armv7l`, `armv6l`, `arm` | `dasel_linux_arm32` |
| Linux | 386 | `i386`, `i686` | `dasel_linux_386` |
| macOS | ARM64 | `arm64`, `aarch64` | `dasel_darwin_arm64` |
| macOS | AMD64 | `x86_64` | `dasel_darwin_amd64` |
| Windows | AMD64 | `x86_64`, `amd64` | `dasel_windows_amd64.exe` |
| Windows | 386 | `i386`, `i686` | `dasel_windows_386.exe` |

### Architecture Verification

After downloading, verify the binary format and version:

```sh
file /usr/local/bin/dasel
dasel --version
```

### Homebrew (macOS/Linux)

```sh
brew install dasel
```

### Go Install

```sh
go install github.com/tomwright/dasel/v3/cmd/dasel@master
```

### Prebuilt Binaries

Prebuilt binaries are available on the [Releases](https://github.com/TomWright/dasel/releases) page for Linux, macOS, and Windows.

### Documentation

See the [installation docs](https://daseldocs.tomwright.me/getting-started/installation) for package manager alternatives and container images.

---

## Shell Completion

Dasel supports tab-completion for Bash, Zsh, Fish, and PowerShell. Generate and load the completion script for your shell:

```sh
# Bash
source <(dasel completion bash)

# Zsh
source <(dasel completion zsh)

# Fish
dasel completion fish | source

# PowerShell
dasel completion powershell | Out-String | Invoke-Expression
```

See the [shell completion docs](https://daseldocs.tomwright.me/getting-started/shell-completion) for permanent installation instructions.

## Man Page

Generate and view the man page:

```sh
dasel man | man -l -
```

---

## Basic Usage

### Selecting Values

By default, Dasel evaluates the final selector and prints the result.

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json 'foo.bar'
# Output: "baz"
```

### Modifying Values

Update values inline:

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json 'foo.bar = "bong"'
# Output: "bong"
```

Use `--root` to output the full document after modification:

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json --root 'foo.bar = "bong"'
# Output:
{
  "foo": {
    "bar": "bong"
  }
}
```

Update values based on previous value:

```sh
echo '[1,2,3,4,5]' | dasel -i json --root 'each($this = $this*2)'
# Output:
[
    2,
    4,
    6,
    8,
    10
]
```

### Format Conversion

```sh
cat data.json | dasel -i json -o yaml
```

### Recursive Descent (`..`)

Searches all nested objects and arrays for a matching key or index.

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json '..bar'
# Output:
[
    "baz"
]

```

### Search (`search`)

Finds all values matching a condition anywhere in the structure.

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json 'search(bar == "baz")'
# Output:
[
    {
        "bar": "baz"
    }
]

```

---

## Documentation

Full documentation is available at [daseldocs.tomwright.me](https://daseldocs.tomwright.me).

---

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

---

## License

MIT License. See [LICENSE](./LICENSE) for details.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/TomWright/dasel.svg)](https://starchart.cc/TomWright/dasel)
