Write-Host "=== VS Code MCP Server Diagnostic ===" -ForegroundColor Green
Write-Host ""

# 1) Manual MCP initialize test
Write-Host "1) Testing MCP initialize via stdio..." -ForegroundColor Cyan
$req = '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"diagnostic"}}}'
$bytes = [System.Text.Encoding]::UTF8.GetBytes($req)
$header = "Content-Length: " + $bytes.Length + "`r`n`r`n"
$response = ($header + $req) | go run ./cmd/mcp 2>&1
$ok = $false
if ($LASTEXITCODE -eq 0 -and $response -match '"result"') {
  Write-Host "[OK] MCP server responded to initialize" -ForegroundColor Green
  $ok = $true
} else {
  Write-Host "[FAIL] MCP server initialize failed (exit $LASTEXITCODE)" -ForegroundColor Red
  Write-Host $response -ForegroundColor Yellow
}

# 2) Check workspace configs
Write-Host ""
Write-Host "2) Checking workspace config..." -ForegroundColor Cyan
$wsMcp = ".vscode/mcp.json"
$wsSettings = ".vscode/settings.json"
if (Test-Path $wsMcp) {
  Write-Host "[OK] Found .vscode/mcp.json" -ForegroundColor Green
} else {
  Write-Host "[WARN] Missing .vscode/mcp.json (optional)" -ForegroundColor Yellow
}
if (Test-Path $wsSettings) {
  $settings = Get-Content $wsSettings -Raw
  if ($settings -match "languageModelTools.mcpServers") {
    Write-Host "[OK] settings.json has languageModelTools.mcpServers" -ForegroundColor Green
  } else {
    Write-Host "[WARN] settings.json missing languageModelTools.mcpServers (optional)" -ForegroundColor Yellow
  }
} else {
  Write-Host "[WARN] Missing .vscode/settings.json (optional)" -ForegroundColor Yellow
}

# 3) Create user-level mcp.json
Write-Host ""
Write-Host "3) Creating user-level mcp.json..." -ForegroundColor Cyan
$userDir = Join-Path $env:APPDATA "Code\User"
if (!(Test-Path $userDir)) {
  New-Item -ItemType Directory -Path $userDir -Force | Out-Null
}
$userMcpPath = Join-Path $userDir "mcp.json"
$cwd = (Get-Location).Path -replace '\\','/'
$jsonObj = @{
  servers = @{
    "td-go-mcp" = @{
      type = "stdio"
      command = "go"
      args = @("run","./cmd/mcp")
      cwd = $cwd
      env = @{ DB_DSN = "teradw"; DB_DRIVER = "odbc" }
    }
  }
}
$json = $jsonObj | ConvertTo-Json -Depth 8
$json | Set-Content -Path $userMcpPath -Encoding UTF8
if (Test-Path $userMcpPath) {
  Write-Host "[OK] Wrote $userMcpPath" -ForegroundColor Green
} else {
  Write-Host "[FAIL] Failed to write $userMcpPath" -ForegroundColor Red
}

Write-Host ""
Write-Host "Next: Reload VS Code, then check Output -> Language Model Tools / Copilot for MCP logs." -ForegroundColor Yellow