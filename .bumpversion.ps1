# Install GitVersion if it is not already installed
if (!(Get-Command GitVersion)) {
	# Check if Chocolately is installed, if not warn user
	if (!(Get-Command choco)) {
		Write-Output "Please ensure Chocolatey is installed. [https://www.chocolatey.org/]"
		exit
	}
	choco install gitversion.portable -pre -y
}

$Version = GitVersion.exe | ConvertFrom-Json
$FilePath = ".\plex\version.go"
(Get-Content $FilePath) -replace '(?<=Version = ").*?(?=")', $Version.SemVer | Set-Content -Path $FilePath -Encoding UTF8
Write-Output $Version.SemVer
