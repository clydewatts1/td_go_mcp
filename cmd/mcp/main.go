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
	slog.Info("[main] Starting td-go-mcp server")

	wd, err := os.Getwd()
	if err != nil {
		slog.Error("[main] Failed to get working directory", "err", err)
	} else {
		slog.Debug("[main] Working directory", "dir", wd)
	}

	logDirContents("tools")
	logDirContents("prompts")

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
	mcpServer := server.NewMCPServer("td-go-mcp", "0.2.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	slog.Info("[main] Registering tools and prompts with MCP server", "tools", len(loadedTools), "prompts", len(loadedPrompts))
	for _, toolDef := range loadedTools {
		slog.Debug("[main] Registering tool", "tool", toolDef.Name)
		addToolToServer(mcpServer, toolDef)
	}
	for _, promptDef := range loadedPrompts {
		slog.Debug("[main] Registering prompt", "prompt", promptDef.Name)
		addPromptToServer(mcpServer, promptDef)
	}

	slog.Info("[main] Starting MCP server with stdio transport...")
	if err := server.ServeStdio(mcpServer); err != nil {
		slog.Error("[main] Server error", "err", err)
		os.Exit(1)
	}
}

func logDirContents(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		slog.Error("[main] Failed to glob directory", "dir", dir, "err", err)
		return
	}
	if len(files) == 0 {
		slog.Warn("[main] Directory is empty or missing", "dir", dir)
	} else {
		slog.Info("[main] Directory contents", "dir", dir, "files", files)
	}
}
