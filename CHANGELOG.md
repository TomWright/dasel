# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Maps are now ordered internally.
- JSON and YAML maps maintain ordering on read/write.
- `all()` func now works with strings.
- `index()` func now works with strings.

### Fixed

- Multi-document output should now be displayed correctly.
- Index shorthand selector now works with multiple indexes.
- Null values are now correctly handled.

## [v2.2.0] - 2023-04-17

### Added

- `keys()` function.

## [v2.1.2] - 2023-03-27

### Added

- Join function.
- String function.

### Fixed

- Null error caused by null values in arrays. See [PR 307](https://github.com/TomWright/dasel/pull/307).

## [v2.1.1] - 2023-01-19

### Fixed

- Changed go module to `github.com/tomwright/dasel/v2` to ensure it works correctly with go modules.

## [v2.1.0] - 2023-01-11

### Added

- Ability to jump to a parent x levels up with `parent(x)`. Defaults to 1 level.

## [v2.0.2] - 2022-12-07

### Fixed

- Argument parsing issue that caused files to be written to the wrong place. See [discussion 268](https://github.com/TomWright/dasel/discussions/268).

## [v2.0.1] - 2022-12-07

### Added

- `float` type in `put` command.

### Fixed

- Output values are now correctly de-referenced. This fixed issues with encoded values not appearing correctly.
- Escape characters in selector strings now work as expected.

## [v2.0.0] - 2022-12-02

See [documentation](https://daseldocs.tomwright.me) for all changes.

- Selector syntax 

## [v1.27.3] - 2022-10-18

### Fixed

- The compact flag now works with the XML parser.

## [v1.27.2] - 2022-10-18

### Fixed

- Help text for select and delete commands now contain all available parsers.
- Errors now implement the `Is` interface so they are easier to use from go.
- Floats are now formatted in decimal format instead of scientific notification when writing to CSV ([Issue 245](https://github.com/TomWright/dasel/issues/245), [Issue 229](https://github.com/TomWright/dasel/issues/229))

## [v1.27.1] - 2022-09-28

### Fixed

- Improved selector comparison parsing to allow matching on values containing special characters.

## [v1.27.0] - 2022-09-26

### Added

- New `value-file` flag allows you to `put` values read from a file ([Issue 246](https://github.com/TomWright/dasel/issues/246))

## [v1.26.1] - 2022-08-24

### Fixed

- Make the completion command available for use ([Issue 216](https://github.com/TomWright/dasel/issues/216))
- Make the `__complete` command available for use

## [v1.26.0] - 2022-07-09

### Added

- Search optional selector - `(#:key=value)`

## [v1.25.1] - 2022-06-29

### Added

- Pre-commit hooks for validate command.

## [v1.25.0] - 2022-06-26

### Added

- Support for struct type usage in go package.
- Validate command.

## [v1.24.3] - 2022-04-23

### Added

- Gzip compressed binaries on releases.

## [v1.24.2] - 2022-04-22

### Fixed

- Update a package to avoid a High Vulnerability in golang.org/x/crypto with CVE ID [CVE-2022-27191](https://github.com/advisories/GHSA-8c26-wmh5-6g9v)

## [v1.24.1] - 2022-03-28

### Changed

- `storage` package has been moved outside the `internal` package.

### Fixed

- New funcs added in `v1.24.0` can now be used as expected since you can now access the `storage.ReadWriteOption`.

## [v1.24.0] - 2022-03-18

### Added

- `Node.NewFromFile` func to load a root node from a file.
- `Node.NewFromReader` func to load a root node from an `io.Reader`.
- `Node.WriteToFile` func to write results to a file.
- `Node.Write` func to write results to an `io.Writer`.

## [v1.23.0] - 2022-03-10

### Fixed

- Update github.com/pelletier/go-toml to consume fix for https://github.com/TomWright/dasel/issues/191.

### Added

- Sprig functions to output formatter template.

## [v1.22.1] - 2021-11-09

### Fixed

- Cleaned up error output

## [v1.22.0] - 2021-11-09

### Added

- Type selector `[@]`.

### Fixed

- Errors are now written to stderr as expected.

## [v1.21.2] - 2021-10-21

### Added

- Linux arm32 build target.

## [v1.21.1] - 2021-09-30

### Changed
- `--escape-html` flag now defaults to false.

## [v1.21.0] - 2021-09-29

### Added
- `--escape-html` flag.

### Fixed
- `put document` and `put object` are now aware of the `--merge-input-documents` flag.

## [v1.20.1] - 2021-09-28

### Added

- `buster-slim` and `alpine` tags to built docker images.

### Fixed

- Different encodings in XML files are now [handled as expected](https://github.com/TomWright/dasel/issues/164).

## [v1.20.0] - 2021-08-30

### Added

- `-v`, `--value` flag to workaround [dash issue](https://github.com/TomWright/dasel/issues/117).

### Fixed

- Fixed an issue in which unicode characters could cause issues when parsing selectors.

## [v1.19.0] - 2021-08-14

### Added

- `--colour`,`--color` flag to enable colourised output in select command.

## [v1.18.0] - 2021-08-11

### Added

- `--format` flag to `select` command.

## [v1.17.0] - 2021-08-08

### Added

- Support for `!=` comparison operator in dynamic and search selectors.
- Support for `-`/`keyValue` key in dynamic selectors.

## [v1.16.1] - 2021-08-02

### Fixed

- Fixed a bug that stopped the delete command editing files in place.

## [v1.16.0] - 2021-08-01

### Added

- Delete command.

## [v1.15.0] - 2021-05-06

### Added

- `--merge-input-documents` flag.

### Changed

- Optional `noupdater` build tag to disable the self-update command.

### Fixed

- Empty XML documents are now parsed correctly.
  - https://github.com/TomWright/dasel/issues/131

## [v1.14.1] - 2021-04-15

### Added

- arm64 build support.

## [v1.14.0] - 2021-04-11

### Added

- `.[#]` length selector.
- `>` comparison operator.
- `>=` comparison operator.
- `<` comparison operator.
- `<=` comparison operator.

## [v1.13.6] - 2021-03-29

### Changed

- Development versions of dasel will now include more specific version information where possible.

### Fixed

- Fix an issue that stopped dasel being able to output CSV documents when parsed from JSON. 

## [v1.13.5] - 2021-03-22

### Fixed

- Empty map values are now initialised as `map[string]interface{}` rather than `map[interface{}]interface{}`.

## [v1.13.4] - 2021-03-11

### Fixed

- Empty document input is now treated different in select and put commands.
  - https://github.com/TomWright/dasel/issues/99
  - https://github.com/TomWright/dasel/issues/102

## [v1.13.3] - 2021-03-05

### Fixed

- Blank YAML and CSV input is now treated as an empty document.

### Changed

- Blank JSON input is now treated as an empty document.

## [v1.13.2] - 2021-02-25

### Changed

- Improved information provided in `UnsupportedTypeForSelector` errors.
- Upgrade to go 1.16.

### Fixed

- Make sure the `-n`,`--null` flag has an effect in multi-select queries.

## [v1.13.1] - 2021-02-18

### Fixed

- Added `CGO_ENABLED=0` build flag to ensure linux_amd64 builds are statically linked.

## [v1.13.0] - 2021-02-11

### Added

- `--length` flag to select command.

## [v1.12.2] - 2021-01-05

### Fixed

- Fix a bug that stopped the write parser being properly detected when writing to the input file.

## [v1.12.1] - 2021-01-05

### Changed

- Build workflows now updated to run on ubuntu-latest and use a matrix to build assets for `linux`, `darwin` and
`windows` for both `amd64` and `386`.

### Fixed

- Release asset for macos/darwin is now named `dasel_darwin_amd64` instead of `dasel_macos_amd64`.
- Self-updater now identifies `dev` version as development.

## [v1.12.0] - 2021-01-02

### Added

- Add `-c`, `--compact` flag to remove pretty-print formatting from JSON output.
- Defined `storage.IndentOption(indent string) ReadWriteOption`.
- Defined `storage.PrettyPrintOption(enabled bool) ReadWriteOption`.

### Changed

- Changed `storage.Parser` funcs to allow the passing of `...ReadWriteOption`.

## [v1.11.0] - 2020-12-22

### Added

- Benchmark info now contains graphs.
- `update` command to self-update dasel.

### Changed

- Benchmark info now directly compares dasel, jq and yq.

## [v1.10.0] - 2020-12-19

### Added

- Add `dasel put document` command.
- Benchmark information.

### Fixed

- `-r`,`--read` and `-w`,`--write` flags are now used in `dasel put object`.
- Fix issues that occurred when writing to the root node.

### Changed

- Command names and descriptions.

## [v1.9.1] - 2020-12-12

### Fixed

- Stopped parsing XML entities in strings.

## [v1.9.0] - 2020-12-12

### Added

- Add keys/index selector in multi queries.
- Add `-n`,`--null` flag.

## [v1.8.0] - 2020-12-01

### Added

- Add ability to use `ANY_INDEX` (`[*]`) and `DYNAMIC` (`(x=y)`) selectors on maps/objects.

## [v1.7.0] - 2020-11-30

### Added

- Add `-r`,`--read` and `-w`,`--write` flags to specifically choose input/output parsers. This allows you to convert data between formats.

## [v1.6.2] - 2020-11-18

### Added

- Add support for multi-document JSON files.

## [v1.6.1] - 2020-11-17

### Changed

- Remove some validation on `dasel put object` to allow you to put empty objects.

## [v1.6.0] - 2020-11-17

### Added

- Add search selector to allow recursive searching from the current node.

## [v1.5.1] - 2020-11-14

### Fixed

- Fixed an issue that stopped new values being saved.

## [v1.5.0] - 2020-11-12

### Added

- Add ability to use `\` as an escape character in selectors.

## [v1.4.1] - 2020-11-11

### Fixed

- Fix an issue when parsing dynamic selectors.

## [v1.4.0] - 2020-11-08

### Added

- Add `-m`,`--multiple` flag to deal with multi-value queries.
- Add `ANY_INDEX` or `[*]` selector.
- Add `NextMultiple` property to the `Node` struct - this is used when processing multi-value queries.
- Add `Node.QueryMultiple` func.
- Add `Node.PutMultiple` func.

## [v1.3.0] - 2020-11-08

### Added

- Add support for CSV files.

## [v1.2.0] - 2020-11-07

### Added

- Add support for multi-document YAML files.
- Add CodeQL step in github actions.

### Changed

- Docker image is now pushed to ghcr instead of github packages.

## [v1.1.0] - 2020-11-01

### Added

- Add sub-selector support in dynamic selectors.

## [v1.0.4] - 2020-10-30

### Added

- Add `--plain` flag to tell dasel to output un-formatted values.

## [v1.0.3] - 2020-10-29

### Changed

- Command output is now followed by a newline.

## [v1.0.2] - 2020-10-28

### Added

- Docker image is now built and pushed when a new release is tagged.

## [v1.0.1] - 2020-10-28

### Added

- Add support for XML.

### Changed

- Add `-` as an alias for `stdin`/`stdout` in `--file` and `--output` flags.
- Selector can now be given as the first argument making the flag itself optional.
- `select` is now the default command.

## [v1.0.0] - 2020-10-27

### Added

- Add lots of tests.
- Add docs.
- Got accepted to go-awesome.

## [v0.0.5] - 2020-09-27

### Added

- Add support for TOML.

## [v0.0.4] - 2020-09-27

### Added

- Ability to check against the node value in a dynamic selector.
- Code coverage.

### Changed

- Use reflection instead of fixed type checks.

## [v0.0.3] - 2020-09-24

### Changed

- Use reflection instead of fixed type checks.
- Extract commands into their own functions to make them testable.

## [v0.0.2] - 2020-09-23

### Added

- Add ability to pipe data in/out of dasel.
- Add dasel put command.

## [v0.0.1] - 2020-09-22

### Added

- Everything!

[unreleased]: https://github.com/TomWright/dasel/compare/v2.2.0...HEAD
[v2.1.2]: https://github.com/TomWright/dasel/compare/v2.1.2...v2.2.0
[v2.1.2]: https://github.com/TomWright/dasel/compare/v2.1.1...v2.1.2
[v2.1.1]: https://github.com/TomWright/dasel/compare/v2.1.0...v2.1.1
[v2.1.0]: https://github.com/TomWright/dasel/compare/v2.0.2...v2.1.0
[v2.0.2]: https://github.com/TomWright/dasel/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/TomWright/dasel/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/TomWright/dasel/compare/v1.27.3...v2.0.0
[v1.27.3]: https://github.com/TomWright/dasel/compare/v1.27.2...v1.27.3
[v1.27.2]: https://github.com/TomWright/dasel/compare/v1.27.1...v1.27.2
[v1.27.1]: https://github.com/TomWright/dasel/compare/v1.27.0...v1.27.1
[v1.27.0]: https://github.com/TomWright/dasel/compare/v1.26.1...v1.27.0
[v1.26.1]: https://github.com/TomWright/dasel/compare/v1.26.0...v1.26.1
[v1.26.0]: https://github.com/TomWright/dasel/compare/v1.25.1...v1.26.0
[v1.25.1]: https://github.com/TomWright/dasel/compare/v1.25.0...v1.25.1
[v1.25.0]: https://github.com/TomWright/dasel/compare/v1.24.3...v1.25.0
[v1.24.3]: https://github.com/TomWright/dasel/compare/v1.24.2...v1.24.3
[v1.24.2]: https://github.com/TomWright/dasel/compare/v1.24.1...v1.24.2
[v1.24.1]: https://github.com/TomWright/dasel/compare/v1.24.0...v1.24.1
[v1.24.0]: https://github.com/TomWright/dasel/compare/v1.23.0...v1.24.0
[v1.23.0]: https://github.com/TomWright/dasel/compare/v1.22.1...v1.23.0
[v1.22.1]: https://github.com/TomWright/dasel/compare/v1.22.0...v1.22.1
[v1.22.0]: https://github.com/TomWright/dasel/compare/v1.21.2...v1.22.0
[v1.21.2]: https://github.com/TomWright/dasel/compare/v1.21.1...v1.21.2
[v1.21.1]: https://github.com/TomWright/dasel/compare/v1.21.0...v1.21.1
[v1.21.0]: https://github.com/TomWright/dasel/compare/v1.20.1...v1.21.0
[v1.20.1]: https://github.com/TomWright/dasel/compare/v1.20.0...v1.20.1
[v1.20.0]: https://github.com/TomWright/dasel/compare/v1.19.0...v1.20.0
[v1.19.0]: https://github.com/TomWright/dasel/compare/v1.18.0...v1.19.0
[v1.18.0]: https://github.com/TomWright/dasel/compare/v1.17.0...v1.18.0
[v1.17.0]: https://github.com/TomWright/dasel/compare/v1.16.1...v1.17.0
[v1.16.1]: https://github.com/TomWright/dasel/compare/v1.16.0...v1.16.1
[v1.16.0]: https://github.com/TomWright/dasel/compare/v1.15.0...v1.16.0
[v1.15.0]: https://github.com/TomWright/dasel/compare/v1.14.1...v1.15.0
[v1.14.1]: https://github.com/TomWright/dasel/compare/v1.14.0...v1.14.1
[v1.14.0]: https://github.com/TomWright/dasel/compare/v1.13.6...v1.14.0
[v1.13.6]: https://github.com/TomWright/dasel/compare/v1.13.5...v1.13.6
[v1.13.5]: https://github.com/TomWright/dasel/compare/v1.13.4...v1.13.5
[v1.13.4]: https://github.com/TomWright/dasel/compare/v1.13.3...v1.13.4
[v1.13.3]: https://github.com/TomWright/dasel/compare/v1.13.2...v1.13.3
[v1.13.2]: https://github.com/TomWright/dasel/compare/v1.13.1...v1.13.2
[v1.13.1]: https://github.com/TomWright/dasel/compare/v1.13.0...v1.13.1
[v1.13.0]: https://github.com/TomWright/dasel/compare/v1.12.2...v1.13.0
[v1.12.2]: https://github.com/TomWright/dasel/compare/v1.12.1...v1.12.2
[v1.12.1]: https://github.com/TomWright/dasel/compare/v1.12.0...v1.12.1
[v1.12.0]: https://github.com/TomWright/dasel/compare/v1.11.0...v1.12.0
[v1.11.0]: https://github.com/TomWright/dasel/compare/v1.10.0...v1.11.0
[v1.10.0]: https://github.com/TomWright/dasel/compare/v1.9.1...v1.10.0
[v1.9.1]: https://github.com/TomWright/dasel/compare/v1.9.0...v1.9.1
[v1.9.0]: https://github.com/TomWright/dasel/compare/v1.8.0...v1.9.0
[v1.8.0]: https://github.com/TomWright/dasel/compare/v1.7.0...v1.8.0
[v1.7.0]: https://github.com/TomWright/dasel/compare/v1.6.2...v1.7.0
[v1.6.2]: https://github.com/TomWright/dasel/compare/v1.6.1...v1.6.2
[v1.6.1]: https://github.com/TomWright/dasel/compare/v1.6.0...v1.6.1
[v1.6.0]: https://github.com/TomWright/dasel/compare/v1.5.1...v1.6.0
[v1.5.1]: https://github.com/TomWright/dasel/compare/v1.5.0...v1.5.1
[v1.5.0]: https://github.com/TomWright/dasel/compare/v1.4.1...v1.5.0
[v1.4.1]: https://github.com/TomWright/dasel/compare/v1.4.0...v1.4.1
[v1.4.0]: https://github.com/TomWright/dasel/compare/v1.3.0...v1.4.0
[v1.3.0]: https://github.com/TomWright/dasel/compare/v1.2.0...v1.3.0
[v1.1.0]: https://github.com/TomWright/dasel/compare/v1.0.4...v1.1.0
[v1.0.4]: https://github.com/TomWright/dasel/compare/v1.0.3...v1.0.4
[v1.0.3]: https://github.com/TomWright/dasel/compare/v1.0.2...v1.0.3
[v1.0.2]: https://github.com/TomWright/dasel/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/TomWright/dasel/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/TomWright/dasel/compare/v0.0.5...v1.0.0
[v0.0.5]: https://github.com/TomWright/dasel/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/TomWright/dasel/compare/v0.0.3...v0.0.4
[v0.0.3]: https://github.com/TomWright/dasel/compare/v0.0.2...v0.0.3
[v0.0.2]: https://github.com/TomWright/dasel/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/TomWright/dasel/releases/tag/v0.0.1
