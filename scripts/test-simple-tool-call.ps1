#!/usr/bin/env powershell
# Test MCP tool call with test messages - simplified version

Write-Host "Testing MCP tool call with test messages (simplified)..."

# Get the MCP executable path  
$mcpExe = ".\bin\mcp.exe"
if (-Not (Test-Path $mcpExe)) {
    Write-Host "Building MCP server first..."
    go build -o bin/mcp.exe ./cmd/mcp
}

Write-Host "Using MCP exe: $mcpExe"

# Initialize request
$initRequest = '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"1.0.0","capabilities":{"tools":{"listChanged":true}},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'

# Tool call request for count_records
$toolCallRequest = '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"count_records","arguments":{"table_name":"users","date_filter":"2024-01-15","status":"active"}}}'

# Combine requests with proper headers
$initBytes = [System.Text.Encoding]::UTF8.GetBytes($initRequest)
$toolCallBytes = [System.Text.Encoding]::UTF8.GetBytes($toolCallRequest)

$combinedRequest = "Content-Length: " + $initBytes.Length + "`r`n`r`n" + $initRequest + "Content-Length: " + $toolCallBytes.Length + "`r`n`r`n" + $toolCallRequest

Write-Host "`nSending combined request..."
$result = $combinedRequest | & $mcpExe

Write-Host "MCP Server Response:"
Write-Host $result

Write-Host "`nTest completed!"