# VS Code MCP Server Installation Guide

## Current Issue
VS Code is starting the MCP server but not communicating with it properly.

## Troubleshooting Steps

### 1. Check VS Code MCP Integration
VS Code may require MCP servers to be configured globally rather than per-workspace.

### 2. Global MCP Configuration
Create: `%APPDATA%\Code\User\mcp.json`
```json
{
  "servers": {
    "td-go-mcp": {
      "type": "stdio",
      "command": "go",
      "args": ["run", "./cmd/mcp"],
      "cwd": "C:/Users/cw171001/OneDrive - Teradata/Documents/GitHub/td_go_mcp/td_go_mcp",
      "env": {
        "DB_DSN": "CLEARSCAPE",
        "DB_DRIVER": "odbc"
      }
    }
  }
}
```

### 3. Language Model Tools Settings
In VS Code User Settings:
```json
{
  "languageModelTools.mcpServers": {
    "td-go-mcp": {
      "command": "go",
      "args": ["run", "./cmd/mcp"],
      "cwd": "C:/Users/cw171001/OneDrive - Teradata/Documents/GitHub/td_go_mcp/td_go_mcp",
      "env": {
        "DB_DSN": "CLEARSCAPE",
        "DB_DRIVER": "odbc"
      }
    }
  }
}
```

### 4. Required Extensions
Make sure these extensions are installed:
- GitHub Copilot
- MCP related extensions

### 5. Restart Steps
1. Close VS Code completely
2. Restart VS Code
3. Check Output panel for "Language Model Tools" or "MCP" logs

## Verify in VS Code

Follow these steps to confirm the MCP server is connected and usable from VS Code:

1) Reload the window
  - Press Ctrl+Shift+P → "Developer: Reload Window".

2) Confirm the MCP server command and logs
  - View → Output → select "Language Model Tools" (or "GitHub Copilot") from the dropdown.
  - You should see the server launching this command:
    - `${workspaceFolder}/bin/mcp.exe`
  - Expected log lines from the server:
    - "MCP server started, waiting for initialize request..."
    - After the client connects: "Received initialize request ..." and "Initialize response sent successfully".

3) Trigger a tool from chat
  - In Copilot Chat, enter: "List the available MCP tools."
  - Optionally, try a preview call: "Call get_user_by_id with preview, user_id=123."

4) If it still says "Waiting for server to respond to initialize"
  - Check `%APPDATA%\Code\User\mcp.json` for conflicts or wrong command/type.
  - Ensure `.vscode/mcp.json` and `.vscode/settings.json` both point to `${workspaceFolder}/bin/mcp.exe` with an absolute `cwd`.
  - Re-run workspace test tasks to verify the server works outside the client:
    - Run task: `test: mcp tools list`.
  - Review the Output panel for any path or permission errors launching `mcp.exe`.

Notes:
- For non-preview SQL execution, configure an ODBC DSN named `CLEARSCAPE` or adjust `DB_DSN`.
- Preview mode is available by adding `"__preview": true` to tool call arguments.