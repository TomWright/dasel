# dasel

[![Gitbook](https://badges.aleen42.com/src/gitbook_1.svg)](https://daseldocs.tomwright.me)
[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/dasel)](https://goreportcard.com/report/github.com/TomWright/dasel)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![GitHub All Releases Downloads](https://img.shields.io/github/downloads/TomWright/dasel/total)
![GitHub License](https://img.shields.io/github/license/TomWright/dasel)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/TomWright/dasel?label=latest%20release)](https://github.com/TomWright/dasel/releases/latest)
[![Homebrew tag (latest by date)](https://img.shields.io/homebrew/v/dasel)](https://formulae.brew.sh/formula/dasel)

Dasel (short for data-selector) allows you to query and modify data structures using selector strings.

Comparable to [jq](https://github.com/stedolan/jq) / [yq](https://github.com/kislyuk/yq), but supports JSON, YAML, TOML, XML and CSV with zero runtime dependencies.

## One tool to rule them all

Say good bye to learning new tools just to work with a different data format.

Dasel uses a standard selector syntax no matter the data format. This means that once you learn how to use dasel you immediately have the ability to query/modify any of the supported data types without any additional tools or effort.

![Update Kubernetes Manifest](update_kubernetes.gif)

## Table of contents

- [Dasel](#dasel)
- [One tool to rule them all](#one-tool-to-rule-them-all)
- [Quickstart](#quickstart)
- [Completion](#completion)
- [Issue vs discussion](#issue-vs-discussion)
- [Features](#features)
- [Table of contents](#table-of-contents)
- [Documentation](#documentation)
- [Playground](#playground)
- [Benchmarks](#benchmarks)
- [Pre-Commit](#pre-commit)

## Quickstart

Dasel is available on [homebrew](https://daseldocs.tomwright.me/installation#homebrew), [ASDF](https://daseldocs.tomwright.me/installation#asdf), [scoop](https://daseldocs.tomwright.me/installation#scoop), [docker](https://daseldocs.tomwright.me/installation#docker), [Nix](https://daseldocs.tomwright.me/installation#nix) or as [compiled binaries](https://daseldocs.tomwright.me/installation#manual) from the [latest release](https://github.com/TomWright/dasel/releases/latest).

```bash
brew install dasel
```

You can also install a [development version](https://daseldocs.tomwright.me/installation#development-version) with:

```bash
go install github.com/tomwright/dasel/v2/cmd/dasel@master
```

For more information see the [installation documentation](https://daseldocs.tomwright.me/installation).

### Select

```bash
echo '{"name": "Tom"}' | dasel -r json 'name'
"Tom"
```

See [select documentation](https://daseldocs.tomwright.me/commands/select).

### Convert json to yaml

```bash
echo '{"name": "Tom"}' | dasel -r json -w yaml
name: Tom
```

See [select documentation](https://daseldocs.tomwright.me/commands/select).

### Put

```bash
echo '{"name": "Tom"}' | dasel put -r json -t string -v 'contact@tomwright.me' 'email'
{
  "email": "contact@tomwright.me",
  "name": "Tom"
}
```

See [put documentation](https://daseldocs.tomwright.me/commands/put).

### Delete

```bash
echo '{
  "email": "contact@tomwright.me",
  "name": "Tom"
}' | dasel delete -r json '.email'
{
  "name": "Tom"
}
```

See [delete documentation](https://daseldocs.tomwright.me/commands/delete).

## Completion

If you want to use completion from the terminal you can do the following (using zsh in this example):

Add the following to `~/.zshrc` and reload your terminal.
```bash
export fpath=(~/zsh/site-functions $fpath)
mkdir -p ~/zsh/site-functions
dasel completion zsh > ~/zsh/site-functions/_dasel
compinit
```

## Pre-Commit

Add `dasel` hooks to `.pre-commit-config.yaml` file

```yaml
- repo: https://github.com/TomWright/dasel
  rev: v1.25.1
  hooks:
    - id: dasel-validate
```

for a native execution of dasel, or use:

- `dasel-validate-docker` pre-commit hook for executing dasel using the official [Docker images](https://daseldocs.tomwright.me/installation#docker)
- `dasel-validate-bin` pre-commit hook for executing dasel using the official [binary](https://daseldocs.tomwright.me/installation)

## Issue vs Discussion

I have enabled [discussions](https://github.com/TomWright/dasel/discussions) on this repository.

I am aware there may be some confusion when deciding where you should communicate when reporting issues, asking questions or raising feature requests so this section aims to help us align on that.

Please [raise an issue](https://github.com/TomWright/dasel/issues) if:

- You find a bug.
- You have a feature request and can clearly describe your request.

Please [open a discussion](https://github.com/TomWright/dasel/discussions) if:

- You have a question.
- You're not sure how to achieve something with dasel.
- You have an idea but don't quite know how you would like it to work.
- You have achieved something cool with dasel and want to show it off.
- Anything else!

## Features

- [Query/select data from structured data files](https://daseldocs.tomwright.me/commands/select).
- [Update data in structured data files](https://daseldocs.tomwright.me/commands/put).
- Create data files.
- [Supports multiple data formats/types](https://daseldocs.tomwright.me/supported-file-formats).
- [Convert between data formats/types](https://daseldocs.tomwright.me/examples/change-file-format).
- Uses a [standard query/selector syntax](https://daseldocs.tomwright.me/functions/selector-overview) across all data formats.
- Zero runtime dependencies.
- [Available on Linux, Mac and Windows](https://daseldocs.tomwright.me/installation).
- Available to [import and use in your own projects](https://pkg.go.dev/github.com/tomwright/dasel/v2).
- [Run via Docker](https://daseldocs.tomwright.me/installation#docker).
- [Faster than jq/yq](#benchmarks).
- [Pre-commit hooks](#pre-commit).

## Documentation

The official dasel docs can be found at [daseldocs.tomwright.me](https://daseldocs.tomwright.me).

## Playground

You can test out dasel commands using the [playground](https://dasel.tomwright.me).

Source code for the playground can be found at [github.com/TomWright/daselplayground](https://github.com/TomWright/daselplayground).

## Benchmarks

In my tests dasel has been up to 3x faster than jq and 15x faster than yq.

See the [benchmark directory](./benchmark/README.md).

## Stargazers over time

[![Stargazers over time](https://starchart.cc/TomWright/dasel.svg)](https://starchart.cc/TomWright/dasel)
