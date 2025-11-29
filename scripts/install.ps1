Write-Host "migInstalling developer tools"

Write-Host "Installing task & air via go generate..."
go generate -tags tools ./tools

$GOPATH = go env GOPATH
$BIN = "$GOPATH\bin"

if (!(Test-Path $BIN)) {
    Write-Host "ERROR: GOPATH/bin not found: $BIN"
    exit 1
}

if ($env:PATH -notlike "*$BIN*") {
    Write-Host "Adding Go bin to PATH..."
    setx PATH "$env:PATH;$BIN" | Out-Null
} else {
    Write-Host "Go bin already in PATH."
}

$TASK = "$BIN\task.exe"

if (!(Test-Path $TASK)) {
    Write-Host "ERROR: task.exe not found in $BIN"
    exit 1
}

Write-Host "Running task init..."
& $TASK init

Write-Host "Installation complete"
