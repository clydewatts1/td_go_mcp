# Build Windows .exe binaries for MCP stdio server and HTTP server
param()

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot  # repo root
Set-Location $root

$bin = Join-Path $root 'bin'
if (-not (Test-Path $bin)) { New-Item -ItemType Directory -Path $bin | Out-Null }

Write-Host "Building mcp.exe..." -ForegroundColor Cyan
& go build -o (Join-Path $bin 'mcp.exe') './cmd/mcp'

Write-Host "Building server.exe..." -ForegroundColor Cyan
& go build -o (Join-Path $bin 'server.exe') './cmd/server'

Write-Host "Built binaries:" -ForegroundColor Green
Get-ChildItem $bin | Select-Object Name,Length,LastWriteTime | Format-Table -AutoSize
