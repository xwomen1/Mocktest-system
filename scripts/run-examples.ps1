Write-Host "=== Running Configuration Examples ===" -ForegroundColor Cyan

Write-Host "`n1. Basic Configuration Example" -ForegroundColor Yellow
cd examples/config-basic
go run .
cd ../..


Write-Host "`n2. Environment Variable Example" -ForegroundColor Yellow
cd examples/config-env
go run .
cd ../..

Write-Host "`n3. Configuration Watch Example" -ForegroundColor Yellow
if (Test-Path "configs\dev\config.yaml") {
    Write-Host "Config file found, starting watcher..." -ForegroundColor Green
    Write-Host "Open another terminal and modify configs/dev/config.yaml" -ForegroundColor White
    Write-Host "Change port to 6000 or log level to 'info'" -ForegroundColor White
    
    cd examples/config-watch
    go run .
    cd ../..
} else {
    Write-Host "Config file not found at configs\dev\config.yaml" -ForegroundColor Red
    Write-Host "Create it first with: .\scripts\create-config.ps1" -ForegroundColor Yellow
}

Write-Host "`n=== All examples completed ===" -ForegroundColor Green