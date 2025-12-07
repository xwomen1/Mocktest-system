Write-Host "Creating configuration files..." -ForegroundColor Green


New-Item -ItemType Directory -Path "configs\dev" -Force | Out-Null
New-Item -ItemType Directory -Path "configs\prod" -Force | Out-Null
New-Item -ItemType Directory -Path "configs\test" -Force | Out-Null

# dev config
$devConfig = @'
environment: development

server:
  host: "0.0.0.0"
  port: 50051
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

nats:
  url: "nats://localhost:4222"
  cluster_id: "upm-dev-cluster"
  client_id: "service-registry-dev"

logging:
  level: "debug"
  format: "json"
  output: "stdout"

registry:
  heartbeat_interval: "30s"
  heartbeat_timeout: "90s"
  load_balancing_strategy: "round_robin"
'@

$devConfig | Out-File configs\dev\config.yaml -Encoding ASCII
Write-Host " Created: configs\dev\config.yaml" -ForegroundColor Green

$testConfig = @'
environment: test
server:
  port: 50052
nats:
  url: "nats://localhost:4222"
logging:
  level: "error"
'@

$testConfig | Out-File configs\test\config.yaml -Encoding ASCII
Write-Host "Created: configs\test\config.yaml" -ForegroundColor Green

Write-Host "`nConfiguration files created successfully!" -ForegroundColor Green
Write-Host "Run examples with: .\scripts\run-examples.ps1" -ForegroundColor Cyan