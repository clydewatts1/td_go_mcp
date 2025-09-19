package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"td_go_mcp/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/exp/slog"
)

// Place tool handler functions here, e.g. createToolHandler(...)

func addToolToServer(mcpServer *server.MCPServer, toolDef tools.ToolDefinition) {
	var handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
	if toolDef.Name == "connection_status" {
		handler = connectionStatusHandler
	} else if toolDef.Name == "glossary_resource" {
		handler = glossaryResourceHandler
	} else {
		handler = createToolHandler(toolDef)
	}
	mcpServer.AddTool(convertToolDefinition(toolDef), handler)
	slog.Info("Registered tool", "tool", toolDef.Name)
}

// Handler for glossary_resource tool
func glossaryResourceHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Info("Handling glossary_resource tool call")
	if len(loadedGlossaries) == 0 {
		return mcp.NewToolResultError("No glossary resources loaded"), nil
	}
	// For now, just return the first loaded glossary
	resultJSON, err := json.Marshal(loadedGlossaries[0].Resource.Words)
	if err != nil {
		return mcp.NewToolResultError("failed to marshal glossary resource: " + err.Error()), nil
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func convertToolDefinition(toolDef tools.ToolDefinition) mcp.Tool {
	opts := []mcp.ToolOption{
		mcp.WithDescription(toolDef.Description),
	}
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
		slog.Info("Handling tool call", "tool", toolDef.Name)
		processor, exists := processors[toolDef.Name]
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("tool processor not found: %s", toolDef.Name)), nil
		}
		args := req.GetArguments()
		params := make(map[string]interface{})
		for paramName := range toolDef.Parameters {
			if value, exists := args[paramName]; exists {
				params[paramName] = value
			}
		}
		preview := false
		if args["__preview"] != nil {
			if b, ok := args["__preview"].(bool); ok {
				preview = b
			}
			delete(params, "__preview")
		}
		if err := processor.ValidateParameters(params); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parameter validation failed: %v", err)), nil
		}
		sql, err := processor.ProcessTemplate(params)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SQL template processing failed: %v", err)), nil
		}
		if strings.TrimSpace(sql) == "" {
			return mcp.NewToolResultError("generated SQL is empty"), nil
		}
		if preview {
			return mcp.NewToolResultText("Generated SQL:\n" + sql), nil
		} else {
			if database == nil {
				if toolDef.ReturnTestMessage != "" {
					testData, err := loadTestMessage(toolDef.ReturnTestMessage)
					if err != nil {
						return mcp.NewToolResultText(fmt.Sprintf("Database connection not available and failed to load test data: %v\n\nGenerated SQL:\n%s", err, sql)), nil
					}
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
				return mcp.NewToolResultText("Database connection not available. Use '__preview': true to see generated SQL.\n\nGenerated SQL:\n" + sql), nil
			}
			rows, err := database.ExecuteQuery(sql)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("SQL execution failed: %v", err)), nil
			}
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func connectionStatusHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	status := "not connected"
	dsn := ""
	dbType := ""
	errMsg := ""

	if database != nil {
		dsn = database.DSN()
		dbType = database.Type()
		err := database.Ping()
		if err == nil {
			status = "connected"
		} else {
			status = "error"
			errMsg = err.Error()
		}
	}

	result := map[string]interface{}{
		"status": status,
		"dsn":    dsn,
		"type":   dbType,
		"error":  errMsg,
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("failed to marshal connection status result: " + err.Error()), nil
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}
