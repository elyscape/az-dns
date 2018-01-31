# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog] and this project adheres to
[Semantic Versioning].

## [Unreleased]

## [0.4.0] - 2018-01-29
### Added
- Some initial tests. More will be added as part of a planned refactor.

### Changed
- **BREAKING**: Due to changes in error handling, error messages are now
  prefixed with `Error: `.

### Fixed
- Non-relative record names with a trailing dot are now correctly converted to
  zone-relative record names.

## [0.3.0] - 2018-01-02
### Added
- Support for CAA records.

### Changed
- **BREAKING**: The root command is no longer exported.
- **BREAKING**: Status and error message are no longer in sentence case and no
  longer have trailing periods.
- Updated Azure SDK.

## [0.2.0] - 2017-12-06
### Added
- Changelog.
- Actual documentation in the README.

### Changed
- **BREAKING**: Renamed project to az-dns.
- Reformatted help text to use spaces instead of tabs.

## 0.1.0 - 2017-11-15
### Added
- `get`, `set`, and `clear` commands.

[Keep a Changelog]: http://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

[Unreleased]: https://github.com/elyscape/az-dns/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/elyscape/az-dns/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/elyscape/az-dns/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/elyscape/az-dns/compare/v0.1.0...v0.2.0
