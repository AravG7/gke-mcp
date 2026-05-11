param (
    [string]$Command = "help"
)

$BinaryName = "gke-mcp.exe"
$UIDir = "ui"

function Show-Help {
    Write-Host "GKE MCP Server - Available Commands (PowerShell)"
    Write-Host ""
    Write-Host "  build         Build the binary"
    Write-Host "  build-ui      Build the UI TypeScript code"
    Write-Host "  run           Build and run the server (stdio)"
    Write-Host "  run-http      Build and run the server (http on 8080)"
    Write-Host "  install       Install the binary to GOPATH/bin"
    Write-Host "  install-ui    Install the UI npm packages"
    Write-Host "  test          Run tests"
    Write-Host "  test-ui       Run UI tests"
    Write-Host "  clean         Remove build artifacts"
}

function Build-UI {
    Write-Host "Building UI..."
    Push-Location $UIDir
    npm run build
    Pop-Location
    Write-Host "[OK] UI built"
}

function Build {
    Build-UI
    Write-Host "Building $BinaryName..."
    go build -o $BinaryName .
    Write-Host "[OK] Built $BinaryName"
}

function Run-App {
    Build
    & .\$BinaryName
}

function Run-Http {
    Build
    & .\$BinaryName --server-mode http --server-port 8080
}

function Install-App {
    Write-Host "Installing $BinaryName..."
    go install .
    $GoPath = (go env GOPATH)
    Write-Host "[OK] Installed to $GoPath\bin\$BinaryName"
}

function Install-UI {
    Write-Host "Installing UI..."
    Push-Location $UIDir
    npm install
    Pop-Location
    Write-Host "[OK] Installed UI to $UIDir"
}

function Run-Tests {
    Write-Host "Running tests..."
    go test -v ./...
}

function Run-UITests {
    Write-Host "Running UI tests..."
    Push-Location $UIDir
    npm run test
    Pop-Location
}

function Clean {
    Write-Host "Cleaning up..."
    if (Test-Path $BinaryName) { Remove-Item $BinaryName }
    if (Test-Path coverage.out) { Remove-Item coverage.out }
    if (Test-Path coverage.html) { Remove-Item coverage.html }
    if (Test-Path dist) { Remove-Item -Recurse -Force dist }
    if (Test-Path ui/dist) { Remove-Item -Recurse -Force ui/dist }
    Write-Host "[OK] Cleaned"
}

switch ($Command) {
    "build" { Build }
    "build-ui" { Build-UI }
    "run" { Run-App }
    "run-http" { Run-Http }
    "install" { Install-App }
    "install-ui" { Install-UI }
    "test" { Run-Tests }
    "test-ui" { Run-UITests }
    "clean" { Clean }
    "help" { Show-Help }
    default { 
        Write-Host "Unknown command: $Command"
        Show-Help 
    }
}