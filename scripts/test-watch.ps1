Write-Host "=== Testing Configuration Watch ===" -ForegroundColor Cyan

$currentDir = Get-Location
Write-Host "Current directory: $currentDir"

# check config file
$configPath = "configs\dev\config.yaml"
if (-not (Test-Path $configPath)) {
    Write-Host "Creating config file..." -ForegroundColor Yellow
    .\scripts\create-config.ps1
}

Write-Host "`nStarting config watcher..." -ForegroundColor Green
Write-Host "In another terminal, edit and save:" -ForegroundColor White
Write-Host "  $configPath" -ForegroundColor Gray
Write-Host "`nChanges to try:" -ForegroundColor White
Write-Host "  1. Change server.port to 6000" -ForegroundColor Gray
Write-Host "  2. Change logging.level to 'info'" -ForegroundColor Gray
Write-Host "  3. Change nats.url to 'nats://localhost:4223'" -ForegroundColor Gray
Write-Host "`nPress Ctrl+C in this terminal to stop watching" -ForegroundColor Yellow


go run examples/config-watch/main.go