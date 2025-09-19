package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"td_go_mcp/internal/db"
	"td_go_mcp/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	loadedTools   []tools.ToolDefinition
	loadedPrompts []tools.PromptDefinition
	processors    map[string]*tools.SQLProcessor
	database      *db.DB
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

	// Load prompts from YAML files
	loadedPrompts, err = tools.LoadPromptsFromDirectory("tools")
	if err != nil {
		log.Printf("Error loading prompts: %v", err)
		loadedPrompts = []tools.PromptDefinition{} // Continue with empty prompts
	}

	log.Printf("Loaded %d tools and %d prompts from YAML files", len(loadedTools), len(loadedPrompts))

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

func addToolToServer(mcpServer *server.MCPServer, toolDef tools.ToolDefinition) {
	// Convert tool definition to mcp.Tool
	mcpTool := convertToolDefinition(toolDef)

	// Create handler function for this tool
	handler := createToolHandler(toolDef)

	// Add tool to server
	mcpServer.AddTool(mcpTool, handler)
	log.Printf("Registered tool: %s", toolDef.Name)
}

func addPromptToServer(mcpServer *server.MCPServer, promptDef tools.PromptDefinition) {
	// Convert prompt definition to mcp.Prompt
	mcpPrompt := convertPromptDefinition(promptDef)

	// Create handler function for this prompt
	handler := createPromptHandler(promptDef)

	// Add prompt to server
	mcpServer.AddPrompt(mcpPrompt, handler)
	log.Printf("Registered prompt: %s", promptDef.Name)
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
			// Execute SQL against database or return test message
			if database == nil {
				// If no database connection and test message is configured, use test data
				if toolDef.ReturnTestMessage != "" {
					testData, err := loadTestMessage(toolDef.ReturnTestMessage)
					if err != nil {
						return mcp.NewToolResultText(fmt.Sprintf("Database connection not available and failed to load test data: %v\n\nGenerated SQL:\n%s", err, sql)), nil
					}

					// Return test data with metadata
					result := map[string]interface{}{
						"data":   testData,
						"source": "test_message",
						"file":   toolDef.ReturnTestMessage,
						"sql":    sql,
					}

					resultJSON, err := json.Marshal(result)
					if err != nil {
						return mcp.NewToolResultError(fmt.Sprintf("failed to marshal test results: %v", err)), nil
					}

					return mcp.NewToolResultText(string(resultJSON)), nil
				}

				// No test data available
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

// loadTestMessage loads test data from a JSON file
func loadTestMessage(filepath string) (interface{}, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test message file %s: %w", filepath, err)
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse test message JSON from %s: %w", filepath, err)
	}

	return result, nil
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

func convertPromptDefinition(promptDef tools.PromptDefinition) mcp.Prompt {
	// Start with basic prompt options
	opts := []mcp.PromptOption{
		mcp.WithPromptDescription(promptDef.Description),
	}

	// Convert parameters to prompt arguments
	if len(promptDef.Parameters) > 0 {
		for paramName, param := range promptDef.Parameters {
			argumentOpts := []mcp.ArgumentOption{
				mcp.ArgumentDescription(param.Description),
			}

			// Check if this parameter is required (we'll use a simple heuristic -
			// parameters without defaults are considered required)
			if param.Default == nil {
				argumentOpts = append(argumentOpts, mcp.RequiredArgument())
			}

			opts = append(opts, mcp.WithArgument(paramName, argumentOpts...))
		}
	}

	return mcp.NewPrompt(promptDef.Name, opts...)
}

func createPromptHandler(promptDef tools.PromptDefinition) func(context.Context, mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		log.Printf("Handling prompt request: %s", promptDef.Name)

		// Extract parameters from request
		args := req.Params.Arguments
		if args == nil {
			args = make(map[string]string)
		}

		// Process the prompt template by substituting parameters
		processedPrompt := promptDef.Prompt

		// Simple template processing - replace {{param}} with actual values
		for paramName, paramValue := range args {
			placeholder := fmt.Sprintf("{{%s}}", paramName)
			processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, paramValue)
		}

		// Replace any remaining placeholders with defaults if available
		for paramName, param := range promptDef.Parameters {
			placeholder := fmt.Sprintf("{{%s}}", paramName)
			if strings.Contains(processedPrompt, placeholder) {
				if param.Default != nil {
					if defaultStr, ok := param.Default.(string); ok {
						processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, defaultStr)
					}
				} else {
					// Remove unfilled placeholders
					processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, "")
				}
			}
		}

		// Create prompt messages
		messages := []mcp.PromptMessage{
			{
				Role:    mcp.RoleUser,
				Content: mcp.TextContent{Type: "text", Text: processedPrompt},
			},
		}

		return mcp.NewGetPromptResult(promptDef.Description, messages), nil
	}
}
