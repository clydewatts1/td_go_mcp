# Logging Configuration for td_go_mcp

## Current Logging Behavior:
- MCP Server: Logs to stderr (terminal)
- HTTP Server: Logs to stdout (terminal)
- Test Client: Logs to stderr (terminal) with structured logging

## Test Client Logging:
```powershell
# Test client logs to stderr by default
go run ./cmd/td_go_mcp_test_client/main.go mcp_requests/tools_list.json

# Redirect test client logs
go run ./cmd/td_go_mcp_test_client/main.go mcp_requests/tools_list.json 2> logs/test-client.log
```

## Redirect Logs to Files:

### MCP Server Logs:
```powershell
# Stderr to file
go run ./cmd/mcp 2> logs/mcp-server.log

# Both stdout and stderr to file
go run ./cmd/mcp > logs/mcp-all.log 2>&1
```

### HTTP Server Logs:
```powershell
# Stdout to file
go run ./cmd/server > logs/http-server.log

# Both stdout and stderr to file
go run ./cmd/server > logs/http-all.log 2>&1
```

### Background Processes with Logging:
```powershell
# Start HTTP server in background with logs
Start-Job -ScriptBlock { go run ./cmd/server > logs/http-server.log 2>&1 }

# Start MCP server with logs (foreground)
go run ./cmd/mcp 2> logs/mcp-server.log
```

## Log File Locations:
- Create a `logs/` directory in your project root
- Use timestamped log files for multiple runs
- Consider log rotation for production use

## Advanced Logging Configuration:
The `internal/logging/` directory contains advanced logging configuration:
- Structured logging with slog
- Configurable log levels
- JSON and text format options
- Custom log handlers for different outputs

## Combined Logging Example:
```powershell
# Log all components to separate files
Start-Job -ScriptBlock { go run ./cmd/server > logs/http-server.log 2>&1 }
go run ./cmd/mcp 2> logs/mcp-server.log &
go run ./cmd/td_go_mcp_test_client/main.go mcp_requests/tools_list.json 2> logs/test-client.log
```

## Viewing Logs:
```powershell
# View live logs
Get-Content logs/mcp-server.log -Wait

# View recent logs
Get-Content logs/mcp-server.log -Tail 20
```