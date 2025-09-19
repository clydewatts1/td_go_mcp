#!/usr/bin/env powershell
# Test MCP tool call with test messages

Write-Host "Testing MCP tool call with test messages..."

# Get the MCP executable path
$mcpExe = ".\bin\mcp.exe"
if (-Not (Test-Path $mcpExe)) {
    Write-Host "Building MCP server first..."
    go build -o bin/mcp.exe ./cmd/mcp
}

Write-Host "Using MCP exe: $mcpExe"

# Test count_records tool call (should use test data)
$countRecordsRequest = '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
        "name": "count_records",
        "arguments": {
            "table_name": "users",
            "date_filter": "2024-01-15",
            "status": "active"
        }
    }
}'

Write-Host "`nTesting count_records tool call..."
$countRecordsBytes = [System.Text.Encoding]::UTF8.GetBytes($countRecordsRequest)
$countRecordsHeader = "Content-Length: " + $countRecordsBytes.Length + "`r`n`r`n"
$countRecordsResult = ($countRecordsHeader + $countRecordsRequest) | & $mcpExe

Write-Host "Count Records Result:"
Write-Host $countRecordsResult

# Test get_user_by_id tool call (should use test data)
$getUserRequest = '{
    "jsonrpc": "2.0", 
    "id": 2,
    "method": "tools/call",
    "params": {
        "name": "get_user_by_id",
        "arguments": {
            "user_id": "12345",
            "include_details": false
        }
    }
}'

Write-Host "`nTesting get_user_by_id tool call..."
$getUserBytes = [System.Text.Encoding]::UTF8.GetBytes($getUserRequest)
$getUserHeader = "Content-Length: " + $getUserBytes.Length + "`r`n`r`n"
$getUserResult = ($getUserHeader + $getUserRequest) | & $mcpExe

Write-Host "Get User Result:"
Write-Host $getUserResult

# Test list_active_sessions tool call (should use test data)
$listSessionsRequest = '{
    "jsonrpc": "2.0",
    "id": 3, 
    "method": "tools/call",
    "params": {
        "name": "list_active_sessions",
        "arguments": {
            "limit": 5,
            "user_type": "user"
        }
    }
}'

Write-Host "`nTesting list_active_sessions tool call..."
$listSessionsBytes = [System.Text.Encoding]::UTF8.GetBytes($listSessionsRequest)
$listSessionsHeader = "Content-Length: " + $listSessionsBytes.Length + "`r`n`r`n"
$listSessionsResult = ($listSessionsHeader + $listSessionsRequest) | & $mcpExe

Write-Host "List Sessions Result:"
Write-Host $listSessionsResult

Write-Host "`nTest message functionality completed!"