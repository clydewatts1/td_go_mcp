// ...existing code...
package main

import (
	"fmt"
	"os"

	mcp "github.com/metoro-io/mcp-golang"
	"golang.org/x/exp/slog"

	"td_go_mcp/internal/transport/stdio"
)

// ...globals and init() moved to init.go...

func main() {
	slog.Info("[main] Starting td-go-mcp server")

	defer func() {
		slog.Info("[main] Shutting down server, closing database connection if open")
		if database != nil {
			err := database.Close()
			if err != nil {
				slog.Error("[main] Error closing database", "err", err)
			}
		}
	}()

	slog.Info("[main] Creating MCP server instance")

	stdioTransport := stdio.NewStdioServerTransport()
	mcpServer := mcp.NewServer(stdioTransport,
		mcp.WithName("td-go-mcp"),
		mcp.WithVersion("0.2.0"),
	)

	// Register tools
	toolHandlers := map[string]any{
		"list_databases":       listDatabasesHandler,
		"count_records":        countRecordsHandler,
		"get_user_by_id":       getUserByIDHandler,
		"list_active_sessions": listActiveSessionsHandler,
		"list_of_tables":       listOfTablesHandler,
		"connection_status":    connectionStatusHandler,
		"glossary_resource":    glossaryResourceHandler,
	}

	for _, toolDef := range loadedTools {
		if handler, ok := toolHandlers[toolDef.Name]; ok {
			err := mcpServer.RegisterTool(toolDef.Name, toolDef.Description, handler)
			if err != nil {
				slog.Error("failed to register tool", "tool", toolDef.Name, "err", err)
			}
		} else {
			slog.Warn("no handler found for tool", "tool", toolDef.Name)
		}
	}

	// Register prompts
	for i := range loadedPrompts {
		promptDef := loadedPrompts[i] // Create a new variable for the loop
		handler := createPromptHandler(promptDef)
		err := mcpServer.RegisterPrompt(promptDef.Name, promptDef.Description, handler)
		if err != nil {
			slog.Error("failed to register prompt", "prompt", promptDef.Name, "err", err)
		}
	}

	slog.Info("[main] Starting MCP server with stdio transport...")
	err := mcpServer.Serve()
	slog.Info("[main] mcpServer.Serve() returned", "err", err)
	if err != nil {
		slog.Error("[main] Server error", "err", err)
		os.Exit(1)
	}
	slog.Info("[main] Server finished. Press enter to exit.")
	fmt.Scanln()
}
