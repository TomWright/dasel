# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Nothing yet.

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

[unreleased]: https://github.com/TomWright/dasel/compare/v1.13.5...HEAD
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
