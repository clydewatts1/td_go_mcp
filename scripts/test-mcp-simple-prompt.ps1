# Test MCP prompt get for simple prompt without parameters
Write-Host "Testing MCP prompt get for simple prompt..."

$mcpExe = "$PSScriptRoot\..\bin\mcp.exe"
if (-not (Test-Path $mcpExe)) {
    $mcpExe = "C:\Users\cw171001\OneDrive - Teradata\Documents\GitHub\td_go_mcp\td_go_mcp\bin\mcp.exe"
}

Write-Host "Using MCP exe: $mcpExe"

# Create the JSON-RPC request for prompts/get
$request = '{"jsonrpc":"2.0","id":3,"method":"prompts/get","params":{"name":"daily_standup"}}'
$bytes = [System.Text.Encoding]::UTF8.GetBytes($request)
$header = "Content-Length: " + $bytes.Length + "`r`n`r`n"
$fullRequest = $header + $request

Write-Host "MCP Prompt Get Result (Simple Prompt):"
$fullRequest | & $mcpExe