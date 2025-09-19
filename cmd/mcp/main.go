// ...existing code...
package main

import (
	"os"

	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/exp/slog"
)

// ...globals and init() moved to init.go...

func main() {
	// Set up logging to file is handled in init.go

	defer func() {
		if database != nil {
			database.Close()
		}
	}()

	mcpServer := server.NewMCPServer("td-go-mcp", "0.2.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	slog.Info("Registering tools and prompts with MCP server", "tools", len(loadedTools), "prompts", len(loadedPrompts))

	for _, toolDef := range loadedTools {
		addToolToServer(mcpServer, toolDef)
	}

	for _, promptDef := range loadedPrompts {
		addPromptToServer(mcpServer, promptDef)
	}

	slog.Info("Starting MCP server with stdio transport...")

	if err := server.ServeStdio(mcpServer); err != nil {
		slog.Error("Server error", "err", err)
		os.Exit(1)
	}
}
