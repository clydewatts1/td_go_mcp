// ...existing code...
package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

// ...globals and init() moved to init.go...

func main() {
	// Set up logging to stderr so it doesn't interfere with MCP protocol on stdout
	log.SetOutput(os.Stderr)

	// Ensure database connection is closed on exit
	defer func() {
		if database != nil {
			database.Close()
		}
	}()

	// Create MCP server
	mcpServer := server.NewMCPServer("td-go-mcp", "0.2.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	log.Printf("Registering %d tools and %d prompts with MCP server", len(loadedTools), len(loadedPrompts))

	// Add tools to the MCP server
	for _, toolDef := range loadedTools {
		addToolToServer(mcpServer, toolDef)
	}

	// Add prompts to the MCP server
	for _, promptDef := range loadedPrompts {
		addPromptToServer(mcpServer, promptDef)
	}

	log.Printf("Starting MCP server with stdio transport...")

	// Start the stdio server (this blocks until the server is stopped)
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
