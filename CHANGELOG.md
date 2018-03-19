# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

<!-- ## [Unreleased]

*   Fixed|New|Changed - ([#1](https://github.com/danstis/rmstale/issues/1)) Description
-->

## [v0.6.0]

*   Fixed - ([#105](https://github.com/danstis/rmstale/issues/105)) Caching of thumbnail images.
*   Fixed - ([#104](https://github.com/danstis/rmstale/issues/104)) HTTPS setting not populating.
*   New - Create API for web interface calls.
*   Changed - ([#70](https://github.com/danstis/rmstale/issues/70)) Move settings into Database.
*   Changed - Refactor logging to new logger package.
*   Changed - ([#57](https://github.com/danstis/rmstale/issues/57)) Update web interface settings page to enable population of all settings.
*   Changed - Index images retreved from API.

## [v0.5.2]

*   Fixed - ([#89](https://github.com/danstis/rmstale/issues/89)) Token file not generated if cache folder does not exist.
*   Changed - Changelog formatting.

## [v0.5.1]

*   Fixed - ([#87](https://github.com/danstis/Plex-Sync/issues/87)) Update to return message if watched item has already been scrobbled on remote server.
*   Changed - Move config into own package.

## [v0.5.0]

*   Changed - ([#63](https://github.com/danstis/Plex-Sync/issues/63)) Add version number to console on startup.
*   Changed - ([#58](https://github.com/danstis/Plex-Sync/issues/58)) Changed text list of selected shows to thumbnail images from destination server.
*   Changed - ([#61](https://github.com/danstis/Plex-Sync/issues/61)) Modify build to control versioning.
*   Changed - ([#74](https://github.com/danstis/Plex-Sync/issues/74)) Moved tvshows.txt file to config folder.
*   Changed - ([#76](https://github.com/danstis/Plex-Sync/issues/76)) Update display of shows on homepage to be responsive.
*   Fixed - Cleanup code errors.

## [v0.4.1]

*   Fixed - ([#59](https://github.com/danstis/Plex-Sync/issues/59)) Handle returned errors for the GetToken method.

## [v0.4.0]

*   Changed - ([#50](https://github.com/danstis/Plex-Sync/issues/50)) Refresh shows when reloading the WebUI homepage.
*   Changed - ([#55](https://github.com/danstis/Plex-Sync/issues/55)) Ignore pre-release tags on header of web interface.
*   Changed - ([#54](https://github.com/danstis/Plex-Sync/issues/54)) Move full version number into the settings page.

## [v0.3.4]

*   New - Script to bump version using GitVersion.
*   New - Go generate script to update version.
*   Changed - Moved builds from TravisCI to Appveyor.

## [v0.3.1]

*   Fixed - Sort selected shows on Home page.

## [v0.3.0]

*   New - Show selected shows on home page.
*   New - Version to WebUI pages.

## [v0.2.0]

*   New - Ability to remove a cached token from the web interface.
*   New - Logging to file.
*   New - Logging to web interface.
*   New - Web interface for management of the application.
*   Changed - After token generation direct user to the settings page.
*   Changed - Moved generation of token file into web interface.
*   Changed - Update logging to go to file and StdOut.
*   Changed - Update release to remove all .go files from subdirectories.

## [v0.1.0]

*   New - Add Readme.
*   New - Usage instructions to the Readme file.
*   New - Cache token in local token file.
*   New - Enable syncing of Show watched status from a local server to a remote server.
*   New - Support for MyPlex Account.
*   Fixed - App would attempt to sync even if a token was not obtained.
*   Fixed - Spaces in TV Show names cause errors.

[unreleased]: https://github.com/danstis/Plex-Sync/compare/v0.6.0...HEAD
[v0.6.0]: https://github.com/danstis/Plex-Sync/compare/v0.5.2...v0.6.0
[v0.5.2]: https://github.com/danstis/Plex-Sync/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/danstis/Plex-Sync/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/danstis/Plex-Sync/compare/v0.4.1...v0.5.0
[v0.4.1]: https://github.com/danstis/Plex-Sync/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/danstis/Plex-Sync/compare/v0.3.4...v0.4.0
[v0.3.4]: https://github.com/danstis/Plex-Sync/compare/v0.3.1...v0.3.4
[v0.3.1]: https://github.com/danstis/Plex-Sync/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/danstis/Plex-Sync/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/danstis/Plex-Sync/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/danstis/Plex-Sync/compare/v0.0.1...v0.1.0
