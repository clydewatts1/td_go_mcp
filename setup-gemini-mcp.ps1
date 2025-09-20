# Gemini CLI MCP Server Setup Script
param(
    [string]$ConfigPath = (Join-Path $env:APPDATA "gemini\mcp_servers.json"),
    [string]$ProjectPath = $PWD.Path
)

$configDir = Split-Path $ConfigPath -Parent
if (!(Test-Path $configDir)) {
    New-Item -ItemType Directory -Path $configDir -Force | Out-Null
}

$serverConfig = @{
    command = "go"
    args    = @("run", "./cmd/mcp")
    cwd     = $ProjectPath.Replace('\', '/')
    env     = @{
        DB_DSN                 = "CLEARSCAPE"
        DB_DRIVER              = "odbc"
        TD_GO_MCP_PROJECT_ROOT = $ProjectPath.Replace('\', '/')
    }
}

$config = @{}
if (Test-Path $ConfigPath) {
    $content = Get-Content $ConfigPath -Raw
    if (-not [string]::IsNullOrWhiteSpace($content)) {
        $config = $content | ConvertFrom-Json -ErrorAction SilentlyContinue
    }
}

if ($null -eq $config.servers) {
    $config | Add-Member -MemberType NoteProperty -Name 'servers' -Value @{}
}

$config.servers['td-go-mcp'] = $serverConfig
$config.servers['td-go-mcp'] = $serverConfig

$config | ConvertTo-Json -Depth 10 | Set-Content -Path $ConfigPath -Encoding UTF8 -Force