# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]


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

[unreleased]: https://github.com/danstis/Plex-Sync/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/danstis/Plex-Sync/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/danstis/Plex-Sync/compare/v0.0.1...v0.1.0
