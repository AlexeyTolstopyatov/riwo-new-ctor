#
# Windows script for project building
# Enable PowerShell scripts invoking in Settings for run it
# from any ConHost instance
#
# Now the binary stores in riwo/statis/build because
# all static pages and web elements must store in static catalog
#

$ProjectPath = Get-Location
$OutputPath = "$(Get-Location)\build\main.wasm"
Set-Location $ProjectPath
$Env:GOOS = "js"
$Env:GOARCH = "wasm"
Write-Host "building..."

go build -o $OutputPath "$($ProjectPath)\cmd\main.go"

Write-Host "operation completed? $($?)"