// ...existing code...
package main

import (
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/exp/slog"
)

// ...globals and init() moved to init.go...

func main() {
	// Set up logging to file is handled in init.go

	// Diagnostic logging
	wd, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get working directory", "err", err)
	}
	slog.Info("Starting server", "working_directory", wd)

	// Log contents of tools and prompts directories
	logDirContents("tools")
	logDirContents("prompts")

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

func logDirContents(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		slog.Error("Failed to glob directory", "dir", dir, "err", err)
		return
	}
	slog.Info("Directory contents", "dir", dir, "files", files)
}
