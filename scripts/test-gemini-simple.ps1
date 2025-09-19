# Test with shorter timeout and proper format
Write-Host "Testing Gemini MCP with timeout fix..." -ForegroundColor Green

# Test with proper prompt format and shorter response
Write-Host "Test: Using positional arguments (recommended format)..." -ForegroundColor Cyan

try {
    # Set environment variables
    $env:DB_DSN = "teradw"
    $env:DB_DRIVER = "odbc"
    $env:LOG_LEVEL = "info"
    
    Write-Host "Executing: gemini 'What tools can you use?'" -ForegroundColor Gray
    
    # Use the new positional argument format (not deprecated -p flag)
    $result = gemini 'What tools can you use?' 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: Gemini responded" -ForegroundColor Green
        Write-Host $result -ForegroundColor White
    } else {
        Write-Host "Exit code: $LASTEXITCODE" -ForegroundColor Yellow
        Write-Host $result -ForegroundColor White
    }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nTest: Direct tool usage..." -ForegroundColor Cyan
try {
    Write-Host "Executing: gemini 'Count records in users table using preview mode'" -ForegroundColor Gray
    $result2 = gemini 'Count records in users table using preview mode' 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: Tool call worked" -ForegroundColor Green
        Write-Host $result2 -ForegroundColor White
    } else {
        Write-Host "Exit code: $LASTEXITCODE" -ForegroundColor Yellow  
        Write-Host $result2 -ForegroundColor White
    }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
}