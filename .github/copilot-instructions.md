# td_go_mcp Project Instructions

This is a database-centric MCP (Model Context Protocol) server in Go.

## Project Type
MCP server scaffold in Go (Windows/PowerShell). Uses stdio command `go run ./cmd/mcp`.

## Requirements
This MCP server connects to databases and performs operations based on MCP tool calls.
- YAML-based tool definitions in the `tools/` directory
- Database connection via ODBC (default: Teradata DSN 'CLEARSCAPE')  
- SQL template processing with parameter substitution
- Preview mode for SQL generation without execution
- Comprehensive error handling and validation

## Features Implemented
- MCP stdio server with `initialize`, `tools/list`, and `tools/call`
- Dynamic tool loading from YAML configurations
- Database integration with ODBC connectivity
- SQL template processing with Go templates
- Basic HTTP server with health and info endpoints
- VS Code tasks for building, testing, and running
- PowerShell testing scripts
- Comprehensive documentation

## Development Notes
- Work through each checklist item systematically
- Keep communication concise and focused
- Follow development best practices
