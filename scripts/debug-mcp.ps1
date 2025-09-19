# Debug MCP Communication
# This script will help debug what VS Code is sending to the MCP server

Write-Host "Starting MCP server debug session..." -ForegroundColor Green

# Start the MCP server and capture all input/output
$process = Start-Process -FilePath "go" -ArgumentList "run", "./cmd/mcp" -NoNewWindow -PassThru -RedirectStandardInput -RedirectStandardOutput -RedirectStandardError

Write-Host "MCP server started with PID: $($process.Id)" -ForegroundColor Cyan
Write-Host "You can now try connecting from VS Code..."
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Yellow

try {
    # Wait for the process to finish or be interrupted
    $process.WaitForExit()
} finally {
    if (!$process.HasExited) {
        $process.Kill()
    }
}

Write-Host "MCP server stopped." -ForegroundColor Red