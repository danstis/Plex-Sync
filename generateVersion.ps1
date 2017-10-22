param(
	# Source Folder of the project
	[string]
	$SourceFolder = ".\"
)

$Version = Gitversion.exe | ConvertFrom-Json
$VersionFile = "$SourceFolder\plex\version.go"

Write-Host ("Updateing version to {0}" -f $Version.SemVer)
$Content = (Get-Content $VersionFile) -replace '(?<=Version = ").*?(?=")', $Version.SemVer 

[IO.File]::WriteAllLines($VersionFile, $Content) # Write the file in UFT8 (no BOM)
