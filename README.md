# td_go_mcp# td_go_mcp



A database-centric MCP (Model Context Protocol) server in Go that provides dynamic tools defined via YAML configurations.A minimal Go HTTP server scaffold intended as a base for an MCP server. It exposes a root endpoint returning JSON and a `/healthz` endpoint.



## Features## Prerequisites

- Go 1.21+

- **Dynamic Tool Loading**: Tools are defined in YAML files and loaded at runtime- Windows PowerShell (default shell)

- **Database Integration**: Connect to databases via ODBC (default: Teradata DSN 'teradw')

- **SQL Template Processing**: Template-based SQL generation with parameter substitution## Setup

- **MCP Protocol**: Full stdio-based MCP server with `initialize`, `tools/list`, and `tools/call`- Initialize and download dependencies:

- **HTTP Server**: Health check and info endpoints

- **SQL Preview Mode**: Generate SQL without executing (add `"__preview": true` to tool calls)```powershell

- **Error Handling**: Comprehensive validation and error reporting# From the repo root

go mod tidy

## Prerequisites```



- Go 1.21+## Run

- Windows PowerShell- Start the server (default port 8080):

- ODBC drivers for your target database (optional)

```powershell

## Setup# PowerShell

$env:PORT="8080"; go run ./cmd/server

1. **Initialize dependencies:**```

```powershell

go mod tidyOpen http://localhost:8080 to see the JSON response.

```

## Test

2. **Configure database (optional):**```powershell

   Set environment variables or use default Teradata DSN:go test ./...

```powershell```

$env:DB_DSN="your_dsn_name"

# or## VS Code

$env:DB_CONNECTION_STRING="your_connection_string"- Use the provided tasks in `.vscode/tasks.json`:

```  - `go: tidy`, `go: build`, `go: test`, `go: run server`.

- Debug with the `Debug Server` launch configuration.

## Project Structure

## Next Steps (MCP)

```This repo now includes a minimal MCP stdio server at `cmd/mcp` with:

td_go_mcp/ - `initialize`

├── cmd/ - `tools/list` returning a single `ping` tool

│   ├── mcp/         # MCP stdio server - `tools/call` for `ping`

│   └── server/      # HTTP server

├── internal/VS Code MCP config points to it via `.vscode/mcp.json`.

│   ├── db/          # Database connection and config

│   ├── mcp/         # MCP protocol types and transportManual stdio test (PowerShell):

│   ├── server/      # HTTP server handlers```powershell

│   └── tools/       # Tool definition loading and SQL processing$req = '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"manual"}}}'

├── tools/           # YAML tool definitions$frame = "Content-Length: " + ($req.Length) + "\r\n\r\n" + $req

│   ├── count_records.yamlecho $frame | go run ./cmd/mcp

│   ├── get_user_by_id.yaml```

│   └── list_active_sessions.yaml

└── .vscode/         # VS Code configurationAdditional tool calls:

``````powershell

# List tools

## Tool Definition Format$req = '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'; $frame = "Content-Length: " + ($req.Length) + "`r`n`r`n" + $req; $frame | go run ./cmd/mcp



Tools are defined in YAML files in the `tools/` directory:# Call ping

$req = '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"ping","arguments":{"text":"hello"}}}'; $frame = "Content-Length: " + ($req.Length) + "`r`n`r`n" + $req; $frame | go run ./cmd/mcp

```yaml

name: "example_tool"# Call time

description: "Example tool description"$req = '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"time","arguments":{}}}'; $frame = "Content-Length: " + ($req.Length) + "`r`n`r`n" + $req; $frame | go run ./cmd/mcp

sql_template: |

  SELECT * FROM users # Call upper

  WHERE id = '{{.user_id}}'$req = '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"upper","arguments":{"text":"abc def"}}}'; $frame = "Content-Length: " + ($req.Length) + "`r`n`r`n" + $req; $frame | go run ./cmd/mcp

  {{#if .include_details}}

  AND status = 'active'# Call sum

  {{/if}}$req = '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"sum","arguments":{"numbers":[1,2,3.5]}}}'; $frame = "Content-Length: " + ($req.Length) + "`r`n`r`n" + $req; $frame | go run ./cmd/mcp

parameters:

  user_id:# Call uuid

    type: "string"$req = '{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"uuid","arguments":{}}}'; $frame = "Content-Length: " + ($req.Length) + "`r`n`r`n" + $req; $frame | go run ./cmd/mcp

    description: "User ID to query"```

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