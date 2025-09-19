package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"td_go_mcp/internal/db"
	"td_go_mcp/internal/tools"
)

var (
	loadedTools []tools.ToolDefinition
	processors  map[string]*tools.SQLProcessor
	database    *db.DB
)

func init() {
	var err error
	loadedTools, err = tools.LoadToolsFromDirectory("tools")
	if err != nil {
		log.Printf("Error loading tools: %v", err)
		loadedTools = []tools.ToolDefinition{} // Continue with empty tools
	}

	processors = make(map[string]*tools.SQLProcessor)
	for i := range loadedTools {
		processors[loadedTools[i].Name] = tools.NewSQLProcessor(loadedTools[i])
	}

	log.Printf("Loaded %d tools from YAML files", len(loadedTools))

	// Initialize database connection
	dbConfig := db.LoadConfig()
	database, err = db.Connect(dbConfig)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		log.Printf("Continuing without database - SQL preview mode only")
		database = nil
	} else {
		log.Printf("Database connection established successfully")
	}
}

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
	)

	log.Printf("Registering %d tools with MCP server", len(loadedTools))

	// Add tools to the MCP server
	for _, toolDef := range loadedTools {
		addToolToServer(mcpServer, toolDef)
	}

	log.Printf("Starting MCP server with stdio transport...")

	// Start the stdio server (this blocks until the server is stopped)
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func addToolToServer(mcpServer *server.MCPServer, toolDef tools.ToolDefinition) {
	// Convert tool definition to mcp.Tool
	mcpTool := convertToolDefinition(toolDef)
	
	// Create handler function for this tool
	handler := createToolHandler(toolDef)

	// Add tool to server
	mcpServer.AddTool(mcpTool, handler)
	log.Printf("Registered tool: %s", toolDef.Name)
}

func convertToolDefinition(toolDef tools.ToolDefinition) mcp.Tool {
	// Start with basic tool options
	opts := []mcp.ToolOption{
		mcp.WithDescription(toolDef.Description),
	}

	// Convert parameters to tool options
	for paramName, param := range toolDef.Parameters {
		switch param.Type {
		case "string":
			if contains(toolDef.Required, paramName) {
				opts = append(opts, mcp.WithString(paramName, mcp.Required(), mcp.Description(param.Description)))
			} else {
				paramOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
				if param.Default != nil {
					if defaultStr, ok := param.Default.(string); ok {
						paramOpts = append(paramOpts, mcp.DefaultString(defaultStr))
					}
				}
				opts = append(opts, mcp.WithString(paramName, paramOpts...))
			}
		case "integer", "number":
			if contains(toolDef.Required, paramName) {
				opts = append(opts, mcp.WithNumber(paramName, mcp.Required(), mcp.Description(param.Description)))
			} else {
				paramOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
				if param.Default != nil {
					if defaultNum, ok := param.Default.(float64); ok {
						paramOpts = append(paramOpts, mcp.DefaultNumber(defaultNum))
					}
				}
				opts = append(opts, mcp.WithNumber(paramName, paramOpts...))
			}
		case "boolean":
			if contains(toolDef.Required, paramName) {
				opts = append(opts, mcp.WithBoolean(paramName, mcp.Required(), mcp.Description(param.Description)))
			} else {
				paramOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
				if param.Default != nil {
					if defaultBool, ok := param.Default.(bool); ok {
						paramOpts = append(paramOpts, mcp.DefaultBool(defaultBool))
					}
				}
				opts = append(opts, mcp.WithBoolean(paramName, paramOpts...))
			}
		default:
			// Default to string type for unknown types
			if contains(toolDef.Required, paramName) {
				opts = append(opts, mcp.WithString(paramName, mcp.Required(), mcp.Description(param.Description)))
			} else {
				opts = append(opts, mcp.WithString(paramName, mcp.Description(param.Description)))
			}
		}
	}

	return mcp.NewTool(toolDef.Name, opts...)
}

func createToolHandler(toolDef tools.ToolDefinition) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Handling tool call: %s", toolDef.Name)

		processor, exists := processors[toolDef.Name]
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("tool processor not found: %s", toolDef.Name)), nil
		}

		// Extract parameters from request using the helper methods
		args := req.GetArguments()
		params := make(map[string]interface{})
		for paramName := range toolDef.Parameters {
			if value, exists := args[paramName]; exists {
				params[paramName] = value
			}
		}

		// Check if this is a preview request
		preview := false
		if args["__preview"] != nil {
			if b, ok := args["__preview"].(bool); ok {
				preview = b
			}
			delete(params, "__preview") // Remove from args for validation
		}

		// Validate parameters against tool schema
		if err := processor.ValidateParameters(params); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parameter validation failed: %v", err)), nil
		}

		// Process SQL template
		sql, err := processor.ProcessTemplate(params)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SQL template processing failed: %v", err)), nil
		}

		// Validate SQL is not empty
		if strings.TrimSpace(sql) == "" {
			return mcp.NewToolResultError("generated SQL is empty"), nil
		}

		if preview {
			// Return the generated SQL instead of executing it
			return mcp.NewToolResultText("Generated SQL:\n" + sql), nil
		} else {
			// Execute SQL against database and return results
			if database == nil {
				return mcp.NewToolResultText("Database connection not available. Use '__preview': true to see generated SQL.\n\nGenerated SQL:\n" + sql), nil
			}

			rows, err := database.ExecuteQuery(sql)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("SQL execution failed: %v", err)), nil
			}

			// Convert results to JSON
			resultJSON, err := json.Marshal(map[string]interface{}{
				"rows":  rows,
				"count": len(rows),
				"sql":   sql,
			})
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to marshal results: %v", err)), nil
			}

			return mcp.NewToolResultText(string(resultJSON)), nil
		}
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
