@echo off
REM Test script for the fresh MCP server
echo Testing fresh MCP server...

REM Test initialize request
echo {"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}} | .\td-go-mcp-fresh.exe

echo.
echo Testing tools/list...
echo {"jsonrpc":"2.0","id":2,"method":"tools/list"} | .\td-go-mcp-fresh.exe

echo.
echo Testing count_records tool...
echo {"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"count_records","arguments":{"table_name":"users"}}} | .\td-go-mcp-fresh.exe

echo.
echo Testing get_user_by_id tool...
echo {"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_user_by_id","arguments":{"user_id":"1"}}} | .\td-go-mcp-fresh.exe

echo.
echo Testing list_active_sessions tool...
echo {"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"list_active_sessions"}} | .\td-go-mcp-fresh.exe

echo.
echo Test completed.