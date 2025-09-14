# td_go_mcp

Minimal Go MCP stdio server for VS Code AI.

## Features

- **MCP Protocol**: Full stdio-based MCP server with `initialize`, `tools/list`, and `tools/call`
- **Demo Tools**: `ping`, `time`, `upper`, `sum`, and `uuid` tools for testing
- **VS Code Integration**: Pre-configured for MCP client integration

## Prerequisites

- Go 1.21+

## Setup

1. **Initialize dependencies:**
```powershell
go mod tidy
```

## Project Structure

```
td_go_mcp/
├── cmd/
│   ├── mcp/          # MCP stdio server with demo tools  
│   └── server/       # HTTP server (legacy)
├── internal/
│   └── mcp/          # MCP protocol types and transport
└── .vscode/          # VS Code configuration
```

## Running the MCP Server

### MCP Server (stdio)
```powershell
go run ./cmd/mcp
```

## VS Code Integration

The project includes VS Code configuration for MCP integration:

### Tasks Available:
- `go: tidy` - Install dependencies  
- `go: test` - Run tests
- `go: run mcp` - Start MCP server (foreground)

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

### Manual MCP Testing

1. **Initialize:**
```bash
req='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"manual"}}}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp
```

2. **List tools:**
```bash  
req='{"jsonrpc":"2.0","id":2,"method":"tools/list"}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp
```

3. **Call tools:**
```bash
# ping tool
req='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"ping","arguments":{"text":"hello"}}}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp

# time tool  
req='{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"time","arguments":{}}}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp

# upper tool
req='{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"upper","arguments":{"text":"Go Mcp"}}}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp

# sum tool
req='{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"sum","arguments":{"numbers":[1,2,3.5]}}}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp

# uuid tool
req='{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"uuid","arguments":{}}}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp
```

## Available Tools

- **ping**: Echo back the input text prefixed with 'pong: '  
- **time**: Return current UTC time in RFC3339 format
- **upper**: Uppercase the provided text
- **sum**: Sum an array of numbers  
- **uuid**: Generate a UUID v4

## VS Code Usage

This MCP server is ready for VS Code integration. The configuration is already set up in `.vscode/mcp.json`:

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

### To use with VS Code MCP clients:

1. **Install an MCP-compatible VS Code extension** such as:
   - GitHub Copilot with MCP support
   - Claude for VS Code with MCP support
   - Or any other MCP client extension

2. **Verify the server works** by running manually:
```bash
# Test that server starts and responds
req='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test"}}}'; bytes=`echo -n "$req" | wc -c`; (printf "Content-Length: $bytes\r\n\r\n$req"; sleep 1) | go run ./cmd/mcp
```

3. **Reload VS Code** after installing MCP extensions:
   - Command Palette → "Developer: Reload Window"

4. **The MCP client should automatically discover** the `td-go-mcp` server and its tools:
   - ping: Echo with "pong:" prefix
   - time: Current UTC timestamp  
   - upper: Uppercase text
   - sum: Add array of numbers
   - uuid: Generate UUID v4

5. **Use tools through the MCP client** - specific usage depends on the extension.

## Development

To test the MCP server:
- Use VS Code task `go: run mcp` or debug with "Debug MCP (stdio)"
- Use manual commands above to test protocol directly
- Server supports standard MCP JSON-RPC over stdio with Content-Length framing