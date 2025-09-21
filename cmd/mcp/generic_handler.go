package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"td_go_mcp/internal/tools"

	mcp "github.com/metoro-io/mcp-golang"
	"golang.org/x/exp/slog"
)

// executeTool is a generic function that handles the common logic for executing a tool.
func executeTool(processor *tools.SQLProcessor, params map[string]any, preview bool) (*mcp.ToolResponse, error) {
	if err := processor.ValidateParameters(params); err != nil {
		slog.Error("[generic_handler] Parameter validation failed", "error", err, "params", params)
		return nil, fmt.Errorf("parameter validation failed: %v", err)
	}

	sql, err := processor.ProcessTemplate(params)
	if err != nil {
		slog.Error("[generic_handler] SQL template processing failed", "error", err, "params", params)
		return nil, fmt.Errorf("SQL template processing failed: %v", err)
	}
	slog.Debug("[generic_handler] Generated SQL", "sql", sql)

	if strings.TrimSpace(sql) == "" {
		slog.Error("[generic_handler] Generated SQL is empty", "tool", processor.ToolName)
		return nil, fmt.Errorf("generated SQL is empty")
	}

	if preview {
		slog.Info("[generic_handler] Preview mode enabled, returning generated SQL")
		return mcp.NewToolResponse(mcp.NewTextContent("Generated SQL:\n" + sql)), nil
	}

	if database == nil {
		slog.Warn("[generic_handler] Database connection not available", "tool", processor.ToolName)
		return nil, fmt.Errorf("database connection not available. Use '__preview': true to see generated SQL.\n\nGenerated SQL:\n%s", sql)
	}

	rows, err := database.ExecuteQuery(sql)
	if err != nil {
		slog.Error("[generic_handler] SQL execution failed", "error", err, "sql", sql)
		return nil, fmt.Errorf("SQL execution failed: %v", err)
	}

	jsonResult, err := json.Marshal(map[string]interface{}{
		"rows":  rows,
		"count": len(rows),
		"sql":   sql,
	})
	if err != nil {
		slog.Error("[generic_handler] Failed to marshal results", "error", err)
		return nil, fmt.Errorf("failed to marshal results: %v", err)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
}
