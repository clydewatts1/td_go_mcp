# PowerShell test scripts for td_go_mcp

# Test HTTP endpoints
Write-Host "Testing HTTP endpoints..." -ForegroundColor Green

# Test health endpoint
try {
    $health = Invoke-WebRequest -Uri "http://localhost:8080/healthz" -UseBasicParsing
    Write-Host "Health Check: $($health.StatusCode) - $($health.Content)" -ForegroundColor Cyan
} catch {
    Write-Host "Health Check Failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Make sure HTTP server is running: go run ./cmd/server" -ForegroundColor Yellow
}

# Test root endpoint
try {
    $root = Invoke-WebRequest -Uri "http://localhost:8080/" -UseBasicParsing
    Write-Host "Root Endpoint: $($root.StatusCode)" -ForegroundColor Cyan
    $rootJson = $root.Content | ConvertFrom-Json
    Write-Host "Server Info: $($rootJson.name) v$($rootJson.version) - $($rootJson.tools_loaded) tools loaded" -ForegroundColor Cyan
} catch {
    Write-Host "Root Endpoint Failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nTesting MCP endpoints..." -ForegroundColor Green

# Helper function to send MCP request
function Send-MCPRequest {
    param([string]$JsonRequest)
    
    $bytes = [System.Text.Encoding]::UTF8.GetBytes($JsonRequest)
    $header = "Content-Length: " + $bytes.Length + "`r`n`r`n"
    $payload = $header + $JsonRequest
    
    try {
        $result = $payload | go run ./cmd/mcp 2>$null
        return $result
    } catch {
        Write-Host "MCP Request Failed: $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

# Test MCP initialize
$initReq = '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"powershell-test"}}}'
$initResult = Send-MCPRequest -JsonRequest $initReq
if ($initResult) {
    Write-Host "Initialize: OK" -ForegroundColor Cyan
} else {
    Write-Host "Initialize: FAILED" -ForegroundColor Red
}

# Test MCP tools/list
$toolsListReq = '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
$toolsResult = Send-MCPRequest -JsonRequest $toolsListReq
if ($toolsResult) {
    Write-Host "Tools List: OK" -ForegroundColor Cyan
    # Extract tool count from response
    if ($toolsResult -match '"tools":\[.*?\]') {
        $toolsJson = ($toolsResult | ConvertFrom-Json).result.tools
        Write-Host "Available Tools: $($toolsJson.Count)" -ForegroundColor Cyan
        foreach ($tool in $toolsJson) {
            Write-Host "  - $($tool.name): $($tool.description)" -ForegroundColor Gray
        }
    }
} else {
    Write-Host "Tools List: FAILED" -ForegroundColor Red
}

# Test MCP tool call with preview
$toolCallReq = '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_user_by_id","arguments":{"user_id":"123","__preview":true}}}'
$toolCallResult = Send-MCPRequest -JsonRequest $toolCallReq
if ($toolCallResult) {
    Write-Host "Tool Call (Preview): OK" -ForegroundColor Cyan
} else {
    Write-Host "Tool Call (Preview): FAILED" -ForegroundColor Red
}

Write-Host "`nTest completed!" -ForegroundColor Green
Write-Host "To run individual tests:" -ForegroundColor Yellow
Write-Host "  HTTP Health: Invoke-WebRequest http://localhost:8080/healthz" -ForegroundColor Gray
Write-Host "  MCP Tools: Use VS Code task 'test: mcp tools list'" -ForegroundColor Gray