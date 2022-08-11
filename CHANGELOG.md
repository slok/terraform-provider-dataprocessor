# Changelog

## [Unreleased]

### Added

- Added `check_fs` plugin example.
- Added `data_structure_transformation` plugin example.
- Added `filtering` plugin example.
- Added `simple_validation` plugin example.
- Added `remote_plugin` plugin example.
- Added `complexvalidation` plugin example.

### Fixed

- Error messages.
- `yq` processor panic/errors because of race conditions when using concurrency.

## [v0.3.0] - 2022-08-07

### Added

- `Go plugins v1` data source.

## [v0.2.0] - 2022-08-06

### Added

- `YQ` data source.

### Changed

- `JQ` data source `query` arg has been renamed to `expression`.

## [v0.1.0] - 2022-08-06

### Added

- Bootstrap provider.
- `JQ` data source.
- Released on Terraform registry.

[unreleased]: https://github.com/slok/terraform-provider-dataprocessor/compare/v0.3.0...HEAD
[v0.3.0]: https://github.com/slok/terraform-provider-dataprocessor/releases/tag/v0.3.0
[v0.2.0]: https://github.com/slok/terraform-provider-dataprocessor/releases/tag/v0.2.0
[v0.1.0]: https://github.com/slok/terraform-provider-dataprocessor/releases/tag/v0.1.0
