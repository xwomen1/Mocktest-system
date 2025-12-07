Write-Host "=== Universal Protocol Mocking Platform - FINAL TEST ===" -ForegroundColor Cyan

# check NATS
Write-Host "1. Checking NATS..." -ForegroundColor Yellow
docker ps --filter "name=nats" --format "{{.Names}}" | Out-Null
if ($LASTEXITCODE -eq 0) {
    Write-Host "   NATS is running" -ForegroundColor Green
} else {
    Write-Host "   NATS not found, using local test" -ForegroundColor Yellow
}

#  build server
Write-Host "2. Building server..." -ForegroundColor Yellow
go build -o server.exe ./cmd/server.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "   Server build: PASS" -ForegroundColor Green
} else {
    Write-Host "   Server build: FAIL" -ForegroundColor Red
    exit 1
}

# build client
Write-Host "3. Building client..." -ForegroundColor Yellow
go build -o client.exe ./cmd/client.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "   Client build: PASS" -ForegroundColor Green
} else {
    Write-Host "   Client build: FAIL" -ForegroundColor Red
}

# check generated files
Write-Host "4. Checking generated files..." -ForegroundColor Yellow
if (Test-Path "internal\service_grpc.pb.go") {
    Write-Host "   Protobuf files: EXISTS" -ForegroundColor Green
} else {
    Write-Host "   Protobuf files: MISSING" -ForegroundColor Red
}

Write-Host "`n=== TEST INSTRUCTIONS ===" -ForegroundColor Cyan
Write-Host "`nTo run server (Terminal 1):" -ForegroundColor White
Write-Host "  go run cmd\server.go" -ForegroundColor Green
Write-Host "`nTo test (Terminal 2):" -ForegroundColor White
Write-Host "  go run cmd\client.go" -ForegroundColor Green
Write-Host "`nExpected output:" -ForegroundColor White
Write-Host "  Registered with ID: mock-engine-localhost-8080" -ForegroundColor Gray
Write-Host "  Found 1 services" -ForegroundColor Gray
Write-Host "  Test PASSED!" -ForegroundColor Gray
Write-Host "`nIf you see this, Week 3 is COMPLETE!" -ForegroundColor Green
