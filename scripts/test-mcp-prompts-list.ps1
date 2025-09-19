# Test MCP prompts list
Write-Host "Testing MCP prompts list..."

$mcpExe = "$PSScriptRoot\..\bin\mcp.exe"
if (-not (Test-Path $mcpExe)) {
    $mcpExe = "C:\Users\cw171001\OneDrive - Teradata\Documents\GitHub\td_go_mcp\td_go_mcp\bin\mcp.exe"
}

Write-Host "Using MCP exe: $mcpExe"

# Create the JSON-RPC request for prompts/list
$request = '{"jsonrpc":"2.0","id":1,"method":"prompts/list"}'
$bytes = [System.Text.Encoding]::UTF8.GetBytes($request)
$header = "Content-Length: " + $bytes.Length + "`r`n`r`n"
$fullRequest = $header + $request

Write-Host "MCP Prompts List Result:"
$fullRequest | & $mcpExe