# Test script for fresh-start MCP server with Gemini CLI

# Backup original config
$originalConfig = "$env:APPDATA\gemini\mcp_servers.json"
$backupConfig = "$env:APPDATA\gemini\mcp_servers_backup.json"

if (Test-Path $originalConfig) {
    Copy-Item $originalConfig $backupConfig -Force
    Write-Host "Backed up original Gemini MCP config"
}

# Create temporary config for fresh-start server
$projectPath = "C:\Users\cw171001\OneDrive - Teradata\Documents\GitHub\td_go_mcp\td_go_mcp\fresh-start"
$configDir = Split-Path $originalConfig -Parent
if (!(Test-Path $configDir)) {
    New-Item -ItemType Directory -Path $configDir -Force | Out-Null
}

$freshStartConfig = @{
    servers = @{
        "td-go-mcp-fresh" = @{
            command = "$projectPath\td-go-mcp-fresh.exe"
            args    = @()
            cwd     = $projectPath
            env     = @{}
        }
    }
}

$freshStartConfig | ConvertTo-Json -Depth 10 | Set-Content -Path $originalConfig -Encoding UTF8 -Force
Write-Host "Created temporary Gemini MCP config for fresh-start server"

try {
    # Test 1: List available tools
    Write-Host "`n--- Testing tools/list ---"
    $toolsResult = gemini @td-go-mcp-fresh tools/list
    Write-Host "Tools available:"
    $toolsResult | ConvertTo-Json -Depth 10

    # Test 2: Count records in users table
    Write-Host "`n--- Testing count_records ---"
    $countResult = gemini @td-go-mcp-fresh count_records --table_name users
    Write-Host "Count result: $countResult"

    # Test 3: Get user by ID
    Write-Host "`n--- Testing get_user_by_id ---"
    $userResult = gemini @td-go-mcp-fresh get_user_by_id --user_id "1"
    Write-Host "User result: $userResult"

    # Test 4: List active sessions
    Write-Host "`n--- Testing list_active_sessions ---"
    $sessionsResult = gemini @td-go-mcp-fresh list_active_sessions
    Write-Host "Sessions result: $sessionsResult"

    Write-Host "`n✅ All tests completed successfully!"
}
catch {
    Write-Error "Test failed: $_"
}
finally {
    # Restore original config
    if (Test-Path $backupConfig) {
        Move-Item $backupConfig $originalConfig -Force
        Write-Host "Restored original Gemini MCP config"
    }
}