# If GitVersion is not installed:
#   choco install GitVersion.Portable

$Version = Gitversion.exe | ConvertFrom-Json
$VersionFile = ".\plex\version.go"

Write-Host ("Updateing version to {0}" -f $Version.SemVer)
(Get-Content $VersionFile) -replace '(?<=Version = ").*?(?=")', $Version.SemVer |  Out-File $VersionFile -Encoding utf8
