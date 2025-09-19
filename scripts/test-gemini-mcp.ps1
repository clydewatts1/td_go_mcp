# Test Gemini MCP connection
Write-Host "Testing Gemini MCP connection..." -ForegroundColor Green

# Test initialize request
$initReq = @"
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{},"clientInfo":{"name":"gemini","version":"0.5.3"}}}
"@

$bytes = [System.Text.Encoding]::UTF8.GetBytes($initReq)
$header = "Content-Length: " + $bytes.Length + "`r`n`r`n"
$payload = $header + $initReq

try {
    $exe = Join-Path (Join-Path (Join-Path $PSScriptRoot '..') 'bin') 'mcp.exe'
    $exe = (Resolve-Path $exe).Path
    if (!(Test-Path $exe)) { throw "mcp.exe not found at $exe" }
    Write-Host "Using MCP exe: $exe" -ForegroundColor Gray
    
    # Set environment variables
    $env:DB_DSN = "teradw"
    $env:DB_DRIVER = "odbc"
    $env:LOG_LEVEL = "info"
    
    $result = $payload | & "$exe"
    Write-Host "Initialize Result:" -ForegroundColor Cyan
    Write-Host $result -ForegroundColor White
    
    # Now test tools/list
    Write-Host "`nTesting tools/list after initialize..." -ForegroundColor Green
    $toolsReq = @"
{"jsonrpc":"2.0","id":2,"method":"tools/list"}
"@
    
    $bytes2 = [System.Text.Encoding]::UTF8.GetBytes($toolsReq)
    $header2 = "Content-Length: " + $bytes2.Length + "`r`n`r`n"
    $payload2 = $header2 + $toolsReq
    
    $result2 = $payload2 | & "$exe"
    Write-Host "Tools List Result:" -ForegroundColor Cyan
    Write-Host $result2 -ForegroundColor White
    
} catch {
    Write-Host "Test failed: $($_.Exception.Message)" -ForegroundColor Red
}