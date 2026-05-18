$ErrorActionPreference = "Stop"

$baseUrl = $env:TMS_BASE_URL
if ([string]::IsNullOrWhiteSpace($baseUrl)) {
  $baseUrl = "http://localhost:8080"
}

Write-Host "Checking $baseUrl/health"
$health = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get
$health | ConvertTo-Json -Depth 5

Write-Host "Checking $baseUrl/ready"
$ready = Invoke-RestMethod -Uri "$baseUrl/ready" -Method Get
$ready | ConvertTo-Json -Depth 5
