# td-go-mcp-fresh# td-go-mcp-fresh



A fresh implementation of a database-centric MCP (Model Context Protocol) server in Go using the mark3labs/mcp-go library.A fresh implementation of a database-centric MCP (Model Context Protocol) server using the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) library.



## Overview## Features



This MCP server provides database query tools that can be used by MCP clients to interact with SQLite databases. It's designed as a clean, minimal implementation using modern Go practices and the official MCP Go SDK.- **Built with mark3labs/mcp-go**: Modern, well-maintained MCP library

- **Database Integration**: SQLite database with sample data

## Features- **Three Core Tools**:

  - `count_records` - Count records in tables with optional filtering

- **MCP Protocol Support**: Full Model Context Protocol implementation using mark3labs/mcp-go  - `get_user_by_id` - Retrieve user information by ID

- **Database Tools**: Three core database operation tools:  - `list_active_sessions` - List active user sessions

  - `count_records`: Count records in tables with optional filtering- **Pure Go**: No CGO dependencies, cross-platform compatible

  - `get_user_by_id`: Retrieve user information by ID- **Simple Architecture**: Clean, maintainable codebase

  - `list_active_sessions`: List all active user sessions

- **SQLite Backend**: Uses SQLite for development and testing## Prerequisites

- **JSON Schema Validation**: Proper input validation using JSON Schema

- **Error Handling**: Comprehensive error handling and validation- Go 1.21+ (tested with Go 1.24.0)



## Prerequisites## Quick Start



- Go 1.21 or later1. **Clone and build:**

- SQLite (automatically handled by the glebarez/sqlite driver)```bash

git checkout refactor/mcp-go-fresh-start

## Buildingcd fresh-start

go build -o mcp-fresh.exe

```bash```

go build .

```2. **Run the server:**

```bash

This creates a `td-go-mcp-fresh.exe` executable../mcp-fresh.exe

```

## Running

3. **Test with MCP client:**

The server runs as an MCP stdio server:```bash

# In another terminal

```bashecho '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./mcp-fresh.exe

./td-go-mcp-fresh.exe```

```

## Available Tools

## Testing

### count_records

Run the included test script:Count records in a specified table with optional filtering.



```bash**Parameters:**

./test-server.bat- `table_name` (string, required): Name of the table to count records from

```- `status` (string, optional): Optional status filter



This will test:**Example:**

- MCP protocol initialization```json

- Tool listing{

- All three database tools with sample data  "jsonrpc": "2.0",

  "id": 1,

## Database Schema  "method": "tools/call",

  "params": {

The server initializes with sample data including:    "name": "count_records",

    "arguments": {

- **users** table: id, name, email, status      "table_name": "users",

- **sessions** table: id, user_id, start_time, status      "status": "active"

    }

## Tools  }

}

### count_records```



Counts records in a specified table with optional filtering.### get_user_by_id

Retrieve user information by ID.

**Parameters:**

- `table_name` (string, required): Name of the table to count records from**Parameters:**

- `date_filter` (string, optional): Date filter in YYYY-MM-DD format- `user_id` (string, required): The unique identifier for the user

- `status` (string, optional): Status filter- `include_details` (boolean, optional): Include detailed user information



**Example:****Example:**

```json```json

{{

  "name": "count_records",  "jsonrpc": "2.0",

  "arguments": {  "id": 2,

    "table_name": "users",  "method": "tools/call",

    "status": "active"  "params": {

  }    "name": "get_user_by_id",

}    "arguments": {

```      "user_id": "1",

      "include_details": true

### get_user_by_id    }

  }

Retrieves user information by ID.}

```

**Parameters:**

- `user_id` (string, required): The unique identifier for the user### list_active_sessions

- `include_details` (boolean, optional): Include detailed user information (default: false)List all active user sessions.



**Example:****Parameters:**

```json- `limit` (integer, optional): Maximum number of sessions to return (default: 10)

{

  "name": "get_user_by_id",**Example:**

  "arguments": {```json

    "user_id": "1",{

    "include_details": true  "jsonrpc": "2.0",

  }  "id": 3,

}  "method": "tools/call",

```  "params": {

    "name": "list_active_sessions",

### list_active_sessions    "arguments": {

      "limit": 5

Lists all active user sessions.    }

  }

**Parameters:**}

- `limit` (integer, optional): Maximum number of sessions to return (default: 100)```

- `user_type` (string, optional): Filter by user type (not implemented in this version)

## Database Schema

**Example:**

```jsonThe server uses an in-memory SQLite database with the following schema:

{

  "name": "list_active_sessions",```sql

  "arguments": {-- Users table

    "limit": 10CREATE TABLE users (

  }    id INTEGER PRIMARY KEY,

}    name TEXT NOT NULL,

```    email TEXT UNIQUE,

    status TEXT DEFAULT 'active'

## Development);



### Dependencies-- Sessions table

CREATE TABLE sessions (

- `github.com/mark3labs/mcp-go`: MCP protocol implementation    id INTEGER PRIMARY KEY,

- `github.com/glebarez/sqlite`: SQLite database driver    user_id INTEGER,

    session_token TEXT,

### Project Structure    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)

```);

.```

├── main.go              # Main server implementation

├── go.mod              # Go module definition## Sample Data

├── go.sum              # Dependency checksums

└── test-server.bat     # Test scriptThe database is pre-populated with sample data:

```- 3 users (John Doe, Jane Smith, Bob Johnson)

- 2 active sessions

## License

## Architecture

This project is part of the td-go-mcp repository.
```
fresh-start/
├── main.go              # Main server implementation
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
└── README.md           # This documentation
```

## Key Differences from Original Implementation

1. **Library**: Uses `mark3labs/mcp-go` instead of custom MCP implementation
2. **Database**: Uses pure Go SQLite driver (no CGO)
3. **Architecture**: Simpler, more maintainable structure
4. **Dependencies**: Minimal external dependencies
5. **Tools**: Focused on core database operations

## Development

To extend the functionality:

1. **Add new tools:**
```go
mcpServer.AddTool(mcp.Tool{
    Name:        "new_tool",
    Description: "Description of new tool",
    InputSchema: mcp.ToolInputSchema{
        Type: "object",
        Properties: map[string]any{
            "param": map[string]any{
                "type": "string",
                "description": "Parameter description",
            },
        },
        Required: []string{"param"},
    },
}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Tool implementation
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            mcp.TextContent{
                Type: "text",
                Text: "Result",
            },
        },
    }, nil
})
```

2. **Modify database schema:**
   Update the `initDatabase` function with new tables/queries

3. **Add configuration:**
   Use environment variables or config files for database connections

## Testing

### Manual Testing
```bash
# Test tools listing
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./mcp-fresh.exe

# Test count_records
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"count_records","arguments":{"table_name":"users"}}}' | ./mcp-fresh.exe

# Test get_user_by_id
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_user_by_id","arguments":{"user_id":"1"}}}' | ./mcp-fresh.exe
```

### Integration with Gemini CLI

1. Configure Gemini CLI to use this server:
```bash
gemini mcp add td-go-mcp-fresh "path/to/mcp-fresh.exe"
```

2. Test in Gemini CLI:
```bash
gemini @td-go-mcp-fresh "count records in the users table"
```

## Migration from Original Implementation

This fresh implementation maintains the same core functionality as the original while providing:

- **Better maintainability** with established MCP library
- **Cross-platform compatibility** with pure Go SQLite
- **Simplified architecture** with less custom code
- **Future-proofing** with actively maintained dependencies

The original YAML-based tool configuration can be re-implemented as needed using the new architecture.