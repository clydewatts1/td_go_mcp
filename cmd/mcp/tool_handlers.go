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
	slog.Debug("[tool_handlers] Adding tool to server", "tool", toolDef.Name)
	var handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
	if toolDef.Name == "connection_status" {
		slog.Debug("[tool_handlers] Using connection_statusHandler", "tool", toolDef.Name)
		handler = connectionStatusHandler
	} else if toolDef.Name == "glossary_resource" {
		slog.Debug("[tool_handlers] Using glossaryResourceHandler", "tool", toolDef.Name)
		handler = glossaryResourceHandler
	} else {
		slog.Debug("[tool_handlers] Using createToolHandler", "tool", toolDef.Name)
		handler = createToolHandler(toolDef)
	}
	mcpServer.AddTool(convertToolDefinition(toolDef), handler)
	slog.Info("[tool_handlers] Registered tool", "tool", toolDef.Name)
}

// Handler for glossary_resource tool
func glossaryResourceHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Info("[tool_handlers] Handling glossary_resource tool call")
	if len(loadedGlossaries) == 0 {
		slog.Warn("[tool_handlers] No glossary resources loaded")
		return mcp.NewToolResultError("No glossary resources loaded"), nil
	}
	resultJSON, err := json.Marshal(loadedGlossaries[0].Resource.Words)
	if err != nil {
		slog.Error("[tool_handlers] Failed to marshal glossary resource", "error", err)
		return mcp.NewToolResultError("failed to marshal glossary resource: " + err.Error()), nil
	}
	slog.Debug("[tool_handlers] Returning glossary resource JSON", "json", string(resultJSON))
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
		slog.Info("[tool_handlers] Handling tool call", "tool", toolDef.Name)
		processor, exists := processors[toolDef.Name]
		if !exists {
			slog.Error("[tool_handlers] Tool processor not found", "tool", toolDef.Name)
			return mcp.NewToolResultError(fmt.Sprintf("tool processor not found: %s", toolDef.Name)), nil
		}
		args := req.GetArguments()
		slog.Debug("[tool_handlers] Tool call arguments", "args", args)
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
			slog.Error("[tool_handlers] Parameter validation failed", "error", err, "params", params)
			return mcp.NewToolResultError(fmt.Sprintf("parameter validation failed: %v", err)), nil
		}
		sql, err := processor.ProcessTemplate(params)
		if err != nil {
			slog.Error("[tool_handlers] SQL template processing failed", "error", err, "params", params)
			return mcp.NewToolResultError(fmt.Sprintf("SQL template processing failed: %v", err)), nil
		}
		slog.Debug("[tool_handlers] Generated SQL", "sql", sql)
		if strings.TrimSpace(sql) == "" {
			slog.Error("[tool_handlers] Generated SQL is empty", "tool", toolDef.Name)
			return mcp.NewToolResultError("generated SQL is empty"), nil
		}
		if preview {
			slog.Info("[tool_handlers] Preview mode enabled, returning generated SQL")
			return mcp.NewToolResultText("Generated SQL:\n" + sql), nil
		} else {
			if database == nil {
				slog.Warn("[tool_handlers] Database connection not available, using test message or preview", "tool", toolDef.Name)
				if toolDef.ReturnTestMessage != "" {
					testData, err := loadTestMessage(toolDef.ReturnTestMessage)
					if err != nil {
						slog.Error("[tool_handlers] Failed to load test message", "file", toolDef.ReturnTestMessage, "error", err)
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
						slog.Error("[tool_handlers] Failed to marshal test results", "error", err)
						return mcp.NewToolResultError(fmt.Sprintf("failed to marshal test results: %v", err)), nil
					}
					slog.Debug("[tool_handlers] Returning test message result JSON", "json", string(resultJSON))
					return mcp.NewToolResultText(string(resultJSON)), nil
				}
				return mcp.NewToolResultText("Database connection not available. Use '__preview': true to see generated SQL.\n\nGenerated SQL:\n" + sql), nil
			}
			rows, err := database.ExecuteQuery(sql)
			if err != nil {
				slog.Error("[tool_handlers] SQL execution failed", "error", err, "sql", sql)
				return mcp.NewToolResultError(fmt.Sprintf("SQL execution failed: %v", err)), nil
			}
			resultJSON, err := json.Marshal(map[string]interface{}{
				"rows":  rows,
				"count": len(rows),
				"sql":   sql,
			})
			if err != nil {
				slog.Error("[tool_handlers] Failed to marshal results", "error", err)
				return mcp.NewToolResultError(fmt.Sprintf("failed to marshal results: %v", err)), nil
			}
			slog.Debug("[tool_handlers] Returning SQL execution result JSON", "json", string(resultJSON))
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
