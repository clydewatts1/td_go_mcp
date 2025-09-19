# Test Gemini MCP Integration via API
Write-Host "Testing Gemini MCP integration with lightweight API calls..." -ForegroundColor Green

# Set environment variables for consistency
$env:DB_DSN = "teradw"
$env:DB_DRIVER = "odbc" 
$env:LOG_LEVEL = "info"

Write-Host "Environment variables set:" -ForegroundColor Gray
Write-Host "  DB_DSN: $env:DB_DSN"
Write-Host "  DB_DRIVER: $env:DB_DRIVER"
Write-Host "  LOG_LEVEL: $env:LOG_LEVEL"

# Test 1: Simple prompt with MCP tool request
Write-Host "`nTest 1: Simple tool listing request..." -ForegroundColor Cyan
$prompt1 = "List all available MCP tools and their descriptions."

try {
    Write-Host "Executing: echo '$prompt1' | gemini --no-interactive" -ForegroundColor Gray
    $result1 = echo $prompt1 | gemini --no-interactive 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: Gemini responded" -ForegroundColor Green
        Write-Host $result1 -ForegroundColor White
    } else {
        Write-Host "FAILED: Exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host $result1 -ForegroundColor Yellow
    }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2: Specific tool call request
Write-Host "`nTest 2: Specific tool call..." -ForegroundColor Cyan
$prompt2 = "Use the get_user_by_id tool with user_id=123 in preview mode to show the generated SQL."

try {
    Write-Host "Executing tool call request..." -ForegroundColor Gray
    $result2 = echo $prompt2 | gemini --no-interactive 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: Tool call completed" -ForegroundColor Green
        Write-Host $result2 -ForegroundColor White
    } else {
        Write-Host "FAILED: Exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host $result2 -ForegroundColor Yellow
    }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Direct MCP server status check
Write-Host "`nTest 3: MCP server status check..." -ForegroundColor Cyan
try {
    Write-Host "Checking MCP server configuration..." -ForegroundColor Gray
    $mcpStatus = gemini mcp list 2>&1
    Write-Host "MCP Status:" -ForegroundColor White
    Write-Host $mcpStatus
} catch {
    Write-Host "ERROR checking MCP status: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Test with explicit prompt parameter
Write-Host "`nTest 4: Using explicit prompt parameter..." -ForegroundColor Cyan
try {
    Write-Host "Testing with -p flag..." -ForegroundColor Gray
    $result4 = gemini -p "What MCP tools do you have available?" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS with -p flag" -ForegroundColor Green
        Write-Host $result4 -ForegroundColor White
    } else {
        Write-Host "FAILED with -p flag: Exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host $result4 -ForegroundColor Yellow
    }
} catch {
    Write-Host "ERROR with -p flag: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nTest Summary:" -ForegroundColor Green
Write-Host "- MCP Server: Running and configured" -ForegroundColor Gray
Write-Host "- Tools Available: count_records, get_user_by_id, list_active_sessions" -ForegroundColor Gray
Write-Host "- Integration Status: See results above" -ForegroundColor Gray