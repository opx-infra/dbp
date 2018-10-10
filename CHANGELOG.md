# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Calendar Versioning](https://calver.org/).

## [Unreleased]

## [18.10.1] - 2018/10/10
### Added
- OPX Installer build support (see [docs](https://opx-infra.github.io/dbp/))

## [18.10.0] - 2018/10/09
### Added
- Changelog
- Documentation and [docs site](https://opx-infra.github.io/dbp/)

### Changed
- Start using Calendar Versioning

## [0.7.1] - 2018/10/08
### Added
- `--no-extra-sources` flag for testing fully-local builds

### Changed
- Default OPX apt sources are used if no extra sources are specified
- Extra source files now expect a `.list` suffix
- Docker images are not pulled if the `--image` flag contains a colon

## [0.6.1] - 2018/10/08
### Changed
- Bump [gbp-docker](https://github.com/opx-infra/gbp-docker/tree/v1.0.5) version to v1.0.5, which fixes lintian issues

## [0.6.0] - 2018/10/08
### Added
- [Controlgraph](https://github.com/opx-infra/controlgraph) for resolving local build dependencies and the proper build order

### Removed
- [Builddepends](https://github.com/opx-infra/builddepends) support

[Unreleased]: https://github.com/opx-infra/dbp/compare/v18.10.1...HEAD
[18.10.1]: https://github.com/opx-infra/dbp/compare/v18.10.0...v18.10.1
[18.10.0]: https://github.com/opx-infra/dbp/compare/v0.7.1...v18.10.0
[0.7.1]: https://github.com/opx-infra/dbp/compare/v0.6.1...v0.7.1
[0.6.1]: https://github.com/opx-infra/dbp/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/opx-infra/dbp/compare/v0.5.5...v0.6.0
