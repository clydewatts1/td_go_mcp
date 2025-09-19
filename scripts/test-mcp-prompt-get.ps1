# Test MCP prompt get functionality
Write-Host "Testing MCP prompt get..."

$mcpExe = "$PSScriptRoot\..\bin\mcp.exe"
if (-not (Test-Path $mcpExe)) {
    $mcpExe = "C:\Users\cw171001\OneDrive - Teradata\Documents\GitHub\td_go_mcp\td_go_mcp\bin\mcp.exe"
}

Write-Host "Using MCP exe: $mcpExe"

# Create the JSON-RPC request for prompts/get
$request = '{"jsonrpc":"2.0","id":2,"method":"prompts/get","params":{"name":"code_review_assistant","arguments":{"code":"function add(a, b) { return a + b; }","language":"javascript","focus_areas":"security, performance"}}}'
$bytes = [System.Text.Encoding]::UTF8.GetBytes($request)
$header = "Content-Length: " + $bytes.Length + "`r`n`r`n"
$fullRequest = $header + $request

Write-Host "MCP Prompt Get Result:"
$fullRequest | & $mcpExe