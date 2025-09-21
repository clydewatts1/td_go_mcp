# td_go_mcp

A database-centric MCP (Model Context Protocol) server in Go.

## Fresh Start Implementation

The current implementation is located in the [resh-start/](./fresh-start/) directory.

This is a clean, minimal MCP server implementation using the mark3labs/mcp-go library, featuring:

- **MCP Protocol Support**: Full Model Context Protocol implementation
- **Database Tools**: Three core database operation tools (count_records, get_user_by_id, list_active_sessions)
- **SQLite Backend**: Uses SQLite for development and testing
- **JSON Schema Validation**: Proper input validation
- **Error Handling**: Comprehensive error handling and validation

## Getting Started

See the [fresh-start README](./fresh-start/README.md) for detailed setup and usage instructions.

## Building and Testing

`ash
cd fresh-start
go build .
./test-server.bat
`

## Directory Structure

- resh-start/ - Clean MCP server implementation using mark3labs/mcp-go
- Legacy directories have been removed for clarity
