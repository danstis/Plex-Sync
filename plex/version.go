package plex

import "strings"

// Version contains the version of the app, this contains the version with any pre-release metadata "1.2.3-prerelease.4"
var Version = "0.0.0-default"

// ShortVersion contains the version in the format "1.2.3"
var ShortVersion = strings.Split(Version, "-")[0]
