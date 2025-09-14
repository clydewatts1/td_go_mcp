# td_go_mcp

A database-centric MCP (Model Context Protocol) server in Go that provides dynamic tools defined via YAML configurations.

## Features

- **Dynamic Tool Loading**: Tools are defined in YAML files and loaded at runtime
- **Database Integration**: Connect to databases via ODBC (default: Teradata DSN 'teradw')
- **SQL Template Processing**: Template-based SQL generation with parameter substitution
- **MCP Protocol**: Full stdio-based MCP server with `initialize`, `tools/list`, and `tools/call`
- **HTTP Server**: Health check and info endpoints
- **SQL Preview Mode**: Generate SQL without executing (add `"__preview": true` to tool calls)
- **Error Handling**: Comprehensive validation and error reporting

## Prerequisites

- Go 1.21+
- Windows PowerShell
- ODBC drivers for your target database (optional)

## Setup

1. **Initialize dependencies:**
```powershell
go mod tidy
```

2. **Configure database (optional):**
   Set environment variables or use default Teradata DSN:
```powershell
$env:DB_DSN="your_dsn_name"
# or
$env:DB_CONNECTION_STRING="your_connection_string"
```

## Project Structure

```
td_go_mcp/
├── cmd/
│   ├── mcp/         # MCP stdio server
│   └── server/      # HTTP server
├── internal/
│   ├── db/          # Database connection and config
│   ├── mcp/         # MCP protocol types and transport
│   ├── server/      # HTTP server handlers
│   └── tools/       # Tool definition loading and SQL processing
├── tools/           # YAML tool definitions
│   ├── count_records.yaml
│   ├── get_user_by_id.yaml
│   └── list_active_sessions.yaml
└── .vscode/         # VS Code configuration
```

## Tool Definition Format

Tools are defined in YAML files in the `tools/` directory:

```yaml
name: "example_tool"
description: "Example tool description"
sql_template: |
  SELECT * FROM users 
  WHERE id = '{{.user_id}}'
  {{#if .include_details}}
  AND status = 'active'
  {{/if}}
parameters:
  user_id:
    type: "string"
    description: "User ID to query"
  include_details:
    type: "boolean" 
    description: "Include detailed information"
    default: false
required: ["user_id"]
return_type: "object"
```

## Running the Servers

### MCP Server (stdio)
```powershell
go run ./cmd/mcp
```

### HTTP Server
```powershell
$env:PORT="8080"; go run ./cmd/server
```

## Gemini CLI Integration

To use this MCP server with the Gemini CLI, create or update your MCP configuration file:

### Windows Configuration
Create or edit `%APPDATA%\gemini\mcp_servers.json`:
```json
{
  "servers": {
    "td-go-mcp": {
      "command": "go",
      "args": ["run", "./cmd/mcp"],
      "cwd": "C:/Users/YOUR_USERNAME/Projects/td_go_mcp",
      "env": {
        "DB_DSN": "teradw"
      }
    }
  }
}
```

### Linux/macOS Configuration
Create or edit `~/.config/gemini/mcp_servers.json`:
```json
{
  "servers": {
    "td-go-mcp": {
      "command": "go",
      "args": ["run", "./cmd/mcp"],
      "cwd": "/path/to/your/td_go_mcp",
      "env": {
        "DB_DSN": "teradw"
      }
    }
  }
}
```

### Usage with Gemini CLI
After configuration, the tools will be available in Gemini CLI:
```bash
# Start Gemini CLI (it will automatically load MCP servers)
gemini

# Use the tools in conversation
"Can you count records in the users table?"
"Get user details for ID 12345"
"List all active sessions"
```

## VS Code Integration

The project includes VS Code configuration for MCP integration:

### Tasks Available:
- `go: tidy` - Install dependencies
- `go: build` - Build all packages
- `go: test` - Run tests
- `go: run server` - Start HTTP server (background)
- `go: run mcp` - Start MCP server (foreground)
- `test: http health` - Test HTTP health endpoint
- `test: mcp tools list` - Test MCP tools listing

### MCP Configuration:
The `.vscode/mcp.json` file configures VS Code to use this server:
```json
{
  "servers": {
    "td-go-mcp": {
      "type": "stdio",
      "command": "go",
      "args": ["run", "./cmd/mcp"]
    }
  }
}
```

## Testing

### Automated Testing
```powershell
# Run PowerShell test script
.\test-endpoints.ps1

# Or run Go tests
go test ./...
```

### Manual MCP Testing

1. **Initialize:**
```powershell
$req='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"manual"}}}'; $bytes=[System.Text.Encoding]::UTF8.GetBytes($req); $header='Content-Length: ' + $bytes.Length + "`r`n`r`n"; ($header + $req) | go run ./cmd/mcp
```

2. **List tools:**
```powershell
$req='{"jsonrpc":"2.0","id":2,"method":"tools/list"}'; $bytes=[System.Text.Encoding]::UTF8.GetBytes($req); $header='Content-Length: ' + $bytes.Length + "`r`n`r`n"; ($header + $req) | go run ./cmd/mcp
```

3. **Call tool (preview mode):**
```powershell
$req='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_user_by_id","arguments":{"user_id":"123","__preview":true}}}'; $bytes=[System.Text.Encoding]::UTF8.GetBytes($req); $header='Content-Length: ' + $bytes.Length + "`r`n`r`n"; ($header + $req) | go run ./cmd/mcp
```

4. **Call tool (execute - requires database):**
```powershell
$req='{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_user_by_id","arguments":{"user_id":"123"}}}'; $bytes=[System.Text.Encoding]::UTF8.GetBytes($req); $header='Content-Length: ' + $bytes.Length + "`r`n`r`n"; ($header + $req) | go run ./cmd/mcp
```

### Manual HTTP Testing

```powershell
# Health check
Invoke-WebRequest http://localhost:8080/healthz

# Server info
Invoke-WebRequest http://localhost:8080/
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_DRIVER` | Database driver | `odbc` |
| `DB_DSN` | ODBC Data Source Name | `teradw` |
| `DB_CONNECTION_STRING` | Full connection string | - |
| `DB_HOST` | Database host | - |
| `DB_PORT` | Database port | - |
| `DB_DATABASE` | Database name | - |
| `DB_USERNAME` | Username | - |
| `DB_PASSWORD` | Password | - |
| `PORT` | HTTP server port | `8080` |

## Adding New Tools

1. Create a new YAML file in the `tools/` directory
2. Define the tool schema (see example above)
3. Use Go template syntax in `sql_template` for parameter substitution
4. Restart the MCP server to load new tools

## Debugging

- Use VS Code launch configuration "Debug MCP (stdio)" to debug the MCP server
- Check VS Code Output panel for MCP server logs
- Use `__preview: true` in tool calls to see generated SQL without execution
- The HTTP server provides tool loading status at the root endpoint

## Development

To extend functionality:
- Add new tool YAML files in `tools/`
- Modify database configuration in `internal/db/config.go`
- Extend SQL template processing in `internal/tools/processor.go`
- Add new HTTP endpoints in `internal/server/server.go`