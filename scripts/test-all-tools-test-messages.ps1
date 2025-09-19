#!/usr/bin/env powershell
# Test all MCP tools with test messages

Write-Host "Testing all MCP tools with test messages..."

$mcpExe = ".\bin\mcp.exe"
Write-Host "Using MCP exe: $mcpExe"

# Helper function to test a tool
function Test-Tool {
    param(
        [string]$toolName,
        [hashtable]$arguments
    )
    
    Write-Host "`n=== Testing $toolName ==="
    
    # Convert arguments to JSON
    $argsJson = $arguments | ConvertTo-Json -Compress
    
    # Create the tool call request
    $toolCallRequest = @{
        jsonrpc = "2.0"
        id = 1
        method = "tools/call"
        params = @{
            name = $toolName
            arguments = $arguments
        }
    } | ConvertTo-Json -Compress
    
    # Initialize request
    $initRequest = '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"1.0.0","capabilities":{"tools":{"listChanged":true}},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
    
    # Combine requests with proper headers
    $initBytes = [System.Text.Encoding]::UTF8.GetBytes($initRequest)
    $toolCallBytes = [System.Text.Encoding]::UTF8.GetBytes($toolCallRequest)
    
    $combinedRequest = "Content-Length: " + $initBytes.Length + "`r`n`r`n" + $initRequest + "Content-Length: " + $toolCallBytes.Length + "`r`n`r`n" + $toolCallRequest
    
    # Send request and get response
    $result = $combinedRequest | & $mcpExe 2>&1
    
    Write-Host "Response:"
    $result | Where-Object { $_ -like "*result*" -and $_ -notlike "*Parse error*" } | ForEach-Object {
        # Extract and pretty print JSON from the response
        if ($_ -match '\{"jsonrpc.*\}$') {
            $jsonResponse = $matches[0] | ConvertFrom-Json
            if ($jsonResponse.result -and $jsonResponse.result.content) {
                $textContent = $jsonResponse.result.content[0].text
                try {
                    $parsedContent = $textContent | ConvertFrom-Json
                    Write-Host ($parsedContent | ConvertTo-Json -Depth 10)
                } catch {
                    Write-Host $textContent
                }
            }
        }
    }
}

# Test count_records tool
Test-Tool -toolName "count_records" -arguments @{
    table_name = "users"
    date_filter = "2024-01-15"
    status = "active"
}

# Test get_user_by_id tool
Test-Tool -toolName "get_user_by_id" -arguments @{
    user_id = "12345"
    include_details = $false
}

# Test list_active_sessions tool  
Test-Tool -toolName "list_active_sessions" -arguments @{
    limit = 10
    user_type = "admin"
}

Write-Host "`n=== All tests completed! ==="