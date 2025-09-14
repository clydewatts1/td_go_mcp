# Gemini CLI MCP Server Setup Script
# This script helps configure the td_go_mcp server for Gemini CLI

param(
    [string]$ConfigPath = "",
    [string]$ProjectPath = $PWD.Path
)

# Determine the config directory based on OS
if ($IsWindows -or $env:OS -eq "Windows_NT") {
    $defaultConfigDir = "$env:APPDATA\gemini"
} else {
    $defaultConfigDir = "$env:HOME/.config/gemini"
}

if ($ConfigPath -eq "") {
    $ConfigPath = Join-Path $defaultConfigDir "mcp_servers.json"
}

Write-Host "Setting up Gemini CLI MCP configuration..." -ForegroundColor Green
Write-Host "Project Path: $ProjectPath" -ForegroundColor Cyan
Write-Host "Config Path: $ConfigPath" -ForegroundColor Cyan

# Create config directory if it doesn't exist
$configDir = Split-Path $ConfigPath -Parent
if (!(Test-Path $configDir)) {
    New-Item -ItemType Directory -Path $configDir -Force | Out-Null
    Write-Host "Created config directory: $configDir" -ForegroundColor Yellow
}

# Create the MCP server configuration
$config = @{
    servers = @{
        "td-go-mcp" = @{
            command = "go"
            args = @("run", "./cmd/mcp")
            cwd = $ProjectPath.Replace('\', '/')
            env = @{
                DB_DSN = "teradw"
                DB_DRIVER = "odbc"
            }
        }
    }
}

# Check if config file exists and merge if needed
if (Test-Path $ConfigPath) {
    Write-Host "Existing config found, merging servers..." -ForegroundColor Yellow
    try {
        $existingConfig = Get-Content $ConfigPath | ConvertFrom-Json
        if ($existingConfig.servers) {
            $existingConfig.servers | Add-Member -NotePropertyName "td-go-mcp" -NotePropertyValue $config.servers."td-go-mcp" -Force
            $config = $existingConfig
        }
    }
    catch {
        Write-Host "Could not parse existing config, creating new one..." -ForegroundColor Yellow
    }
}

# Write the configuration
try {
    $config | ConvertTo-Json -Depth 10 | Set-Content $ConfigPath
    Write-Host "✓ MCP server configuration written to: $ConfigPath" -ForegroundColor Green
    
    Write-Host ""
    Write-Host "Setup complete! You can now use the td-go-mcp server with Gemini CLI." -ForegroundColor Green
    Write-Host ""
    Write-Host "Available tools:" -ForegroundColor Cyan
    Write-Host "- count_records: Count records in database tables"
    Write-Host "- get_user_by_id: Retrieve user details by ID"
    Write-Host "- list_active_sessions: List currently active sessions"
    Write-Host ""
    Write-Host "Example usage in Gemini CLI:" -ForegroundColor Cyan
    Write-Host "  'Can you count records in the users table?'"
    Write-Host "  'Get user details for ID 12345'"
    Write-Host "  'List all active sessions'"
}
catch {
    Write-Error "Failed to write configuration: $_"
    exit 1
}