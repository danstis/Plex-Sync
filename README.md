# Plex-Sync

[![Build status](https://ci.appveyor.com/api/projects/status/bkv4g7crykq7ibc2/branch/master?svg=true)](https://ci.appveyor.com/project/danstis/plex-sync/branch/master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/544fa06319c1471c8d6b0ef5589e4f30)](https://www.codacy.com/app/danstis/Plex-Sync?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=danstis/Plex-Sync&amp;utm_campaign=Badge_Grade)

Plex-Sync is a synchronisation tool for Plex. It syncs watched status (and media in future) between a remote and a local plex server.

MyPlex is used to enable communication with the Plex servers.

## Usage Instructions

1. Download the package from the releases section [here](https://github.com/danstis/Plex-Sync/releases).
1. Extract to a local folder.
1. Update the .\config\config.toml file with the details of your servers:

	```
	[general]
	interval = 600 			# Interval in seconds to perform sync cycle.
	webserverport = 80		# Port for the Web Interface to listen on.
	logfile = "logs/plex-sync.log" 	# Logfile for sync operations
	webserverlogfile = "logs/plex-sync-webserver.log" 	#Logfile for the webserver
	maxlogsize = 20 		# Max logfile Size in MB
	maxlogcount = 5 		# Max number of log backups
	maxlogage = 1 			# Max age of each logfile

	[localServer]
	name = "MyServer"		# Name of your local server, this is the value from MyPlex.
	hostname = "localhost"		# The DNS Hostname or IP Address of your local Plex server.
	port = 32400			# Port used to connect to your local Plex server.
	usessl = false			# Defines if SSL should be used to connect to the local server.

	[remoteServer]
	name = "MyRemote"		# Name of your remote server, this is the value from MyPlex.
	hostname = "server.domain.com"	# The DNS Hostname or IP Address of your remote Plex server.
	port = 32400 			# Port used to connect to your remote Plex server.
	usessl = true			# Defines if SSL should be used to connect to the remote server.

	[webui]
	cacheLifetime = 5	        # Days before cached thumbnails are refreshed.

	```
1. Populate the tvshows.txt file with a list of shows to sync watched status for. The titles should match what is listed in Plex, for example:
	```
	Cops
	The Americans (2013)
	```
1. Run the Plex-Sync.exe file. *NOTE: If running on newer versions of Windows, you will need to allow the file to run when Smart Screen blocks it.*

## Development

This is my first attempt at creating a project using Go, I am always interested in feedback and welcome contributions.

Want to contribute? Great!

To fix a bug or add an enhancement:

* Fork the repo
* Create a new branch ( `git checkout -b improve-feature` )
* Make the appropriate changes in the files
* Update the ChangeLog file with your additions, the format is based on [Keep a Changelog](http://keepachangelog.com/)
* Update the Readme with any changes that are required
* Commit your changes ( `git commit -am 'Improve feature'` )
* Push the branch ( `git push origin improve-feature` )
* Create a Pull Request

### Bug / Feature Request

If you find a bug, or want this tool to do something that it currently does not, please raise an issue [here](https://github.com/danstis/Plex-Sync/issues).

Please be detailed in the issue body.
