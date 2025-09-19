# Test MCP tools list
Write-Host "Testing MCP tools list..." -ForegroundColor Green

$req = @"
{"jsonrpc":"2.0","id":1,"method":"tools/list"}
"@

$bytes = [System.Text.Encoding]::UTF8.GetBytes($req)
$header = "Content-Length: " + $bytes.Length + "`r`n`r`n"
$payload = $header + $req

try {
    $exe = Join-Path (Join-Path (Join-Path $PSScriptRoot '..') 'bin') 'mcp.exe'
    $exe = (Resolve-Path $exe).Path
    if (!(Test-Path $exe)) { throw "mcp.exe not found at $exe" }
    Write-Host "Using MCP exe: $exe" -ForegroundColor Gray
    $result = $payload | & "$exe"
    Write-Host "MCP Tools List Result:" -ForegroundColor Cyan
    Write-Host $result -ForegroundColor White
} catch {
    Write-Host "Test failed: $($_.Exception.Message)" -ForegroundColor Red
}