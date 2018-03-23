# Plex-Sync

[![Build status](https://ci.appveyor.com/api/projects/status/bkv4g7crykq7ibc2/branch/master?svg=true)](https://ci.appveyor.com/project/danstis/plex-sync/branch/master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/544fa06319c1471c8d6b0ef5589e4f30)](https://www.codacy.com/app/danstis/Plex-Sync?utm_source=github.com&utm_medium=referral&utm_content=danstis/Plex-Sync&utm_campaign=Badge_Grade)
[![Coverage Status](https://coveralls.io/repos/github/danstis/Plex-Sync/badge.svg)](https://coveralls.io/github/danstis/Plex-Sync)

Plex-Sync is a synchronisation tool for Plex. It syncs watched status (and media in future) between a remote and a local plex server.

MyPlex is used to enable communication with the Plex servers.

## Usage Instructions

1.  Download the package from the [releases section](https://github.com/danstis/Plex-Sync/releases/latest).
1.  Extract to a local folder.
1.  Populate the .\config\tvshows.txt file with a list of shows to sync watched status for. The titles should match what is listed in Plex, for example:
    ```txt
    Cops
    The Americans (2013)
    ```
1.  Run the Plex-Sync.exe file. _NOTE: If running on newer versions of Windows, you will need to allow the file to run when Smart Screen blocks it._
1.  Open the web interface on the default port 8085 [http://localhost:8085](http://localhost:8085)
1.  Browse to the settings tab to generate a new token, as well as configure your servers and application settings.

## Development

This is my first attempt at creating a project using Go, I am always interested in feedback and welcome contributions.

Want to contribute? Great!

To fix a bug or add an enhancement:

*   Fork the repo
*   Create a new branch ( `git checkout -b improve-feature` )
*   Make the appropriate changes in the files
*   Update the ChangeLog file with your additions, the format is based on [Keep a Changelog](http://keepachangelog.com/)
*   Update the Readme with any changes that are required
*   Commit your changes ( `git commit -am 'Improve feature'` )
*   Push the branch ( `git push origin improve-feature` )
*   Create a Pull Request

### Bug / Feature Request

If you find a bug, or want this tool to do something that it currently does not, please raise an issue [here](https://github.com/danstis/Plex-Sync/issues).

Please be detailed in the issue body.
