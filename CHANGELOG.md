# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Changed

- ([#63](https://github.com/danstis/Plex-Sync/issues/63)) Add version number to console on startup.
- ([#58](https://github.com/danstis/Plex-Sync/issues/58)) Changed text list of selected shows to thumbnail images from destination server.
- ([#61](https://github.com/danstis/Plex-Sync/issues/61)) Modify build to control versioning.
- ([#74](https://github.com/danstis/Plex-Sync/issues/74)) Moved tvshows.txt file to config folder.

### Fixed

- Cleanup code errors.

## [v0.4.1]

### Fixed

- ([#59](https://github.com/danstis/Plex-Sync/issues/59)) Handle returned errors for the GetToken method.

## [v0.4.0]

### Changed

- ([#50](https://github.com/danstis/Plex-Sync/issues/50)) Refresh shows when reloading the WebUI homepage.
- ([#55](https://github.com/danstis/Plex-Sync/issues/55)) Ignore pre-release tags on header of web interface.
- ([#54](https://github.com/danstis/Plex-Sync/issues/54)) Move full version number into the settings page.

## [v0.3.4]

### Added

- Script to bump version using GitVersion.
- Go generate script to update version.

### Changed

- Moved builds from TravisCI to Appveyor.

## [v0.3.1]

### Fixed

- Sort selected shows on Home page.

## [v0.3.0]

### Added

- Show selected shows on home page.
- Version to WebUI pages.

## [v0.2.0]

### Added

- Ability to remove a cached token from the web interface.
- Logging to file.
- Logging to web interface.
- Web interface for management of the application.

### Changed

- After token generation direct user to the settings page.
- Moved generation of token file into web interface.
- Update logging to go to file and StdOut.
- Update release to remove all .go files from subdirectories.

## [v0.1.0]

### Added

- Add Readme.
- Add usage instructions to the Readme file.
- Cache token in local token file.
- Enable syncing of Show watched status from a local server to a remote server.
- Support for MyPlex Account.

### Fixed

- App would attempt to sync even if a token was not obtained.
- Spaces in TV Show names cause errors.

[unreleased]: https://github.com/danstis/Plex-Sync/compare/v0.4.1...HEAD
[v0.4.1]: https://github.com/danstis/Plex-Sync/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/danstis/Plex-Sync/compare/v0.3.4...v0.4.0
[v0.3.4]: https://github.com/danstis/Plex-Sync/compare/v0.3.1...v0.3.4
[v0.3.1]: https://github.com/danstis/Plex-Sync/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/danstis/Plex-Sync/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/danstis/Plex-Sync/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/danstis/Plex-Sync/compare/v0.0.1...v0.1.0
