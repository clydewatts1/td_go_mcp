package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	mcp "github.com/metoro-io/mcp-golang"
	"golang.org/x/exp/slog"
)

func paramsToMap(params any) (map[string]any, error) {
	val := reflect.ValueOf(params)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("params must be a struct, got %T", params)
	}

	result := make(map[string]any)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		parts := strings.Split(jsonTag, ",")
		fieldName := parts[0]

		// Check if the field is empty and has omitempty tag
		if len(parts) > 1 && parts[1] == "omitempty" {
			if reflect.DeepEqual(val.Field(i).Interface(), reflect.Zero(field.Type).Interface()) {
				continue
			}
		}

		result[fieldName] = val.Field(i).Interface()
	}

	return result, nil
}

func listDatabasesHandler(ctx context.Context, args ListDatabasesArgs) (*mcp.ToolResponse, error) {
	slog.Info("[tool_handlers] Handling list_databases tool call")
	processor := processors["list_databases"]
	params, err := paramsToMap(args)
	if err != nil {
		return nil, err
	}
	preview, _ := params["__preview"].(bool)
	return executeTool(processor, params, preview)
}

func countRecordsHandler(ctx context.Context, args CountRecordsArgs) (*mcp.ToolResponse, error) {
	slog.Info("[tool_handlers] Handling count_records tool call")
	processor := processors["count_records"]
	params, err := paramsToMap(args)
	if err != nil {
		return nil, err
	}
	preview, _ := params["__preview"].(bool)
	return executeTool(processor, params, preview)
}

func getUserByIDHandler(ctx context.Context, args GetUserByIDArgs) (*mcp.ToolResponse, error) {
	slog.Info("[tool_handlers] Handling get_user_by_id tool call")
	processor := processors["get_user_by_id"]
	params, err := paramsToMap(args)
	if err != nil {
		return nil, err
	}
	preview, _ := params["__preview"].(bool)
	return executeTool(processor, params, preview)
}

func listActiveSessionsHandler(ctx context.Context, args ListActiveSessionsArgs) (*mcp.ToolResponse, error) {
	slog.Info("[tool_handlers] Handling list_active_sessions tool call")
	processor := processors["list_active_sessions"]
	params, err := paramsToMap(args)
	if err != nil {
		return nil, err
	}
	preview, _ := params["__preview"].(bool)
	return executeTool(processor, params, preview)
}

func listOfTablesHandler(ctx context.Context, args ListOfTablesArgs) (*mcp.ToolResponse, error) {
	slog.Info("[tool_handlers] Handling list_of_tables tool call")
	processor := processors["list_of_tables"]
	params, err := paramsToMap(args)
	if err != nil {
		return nil, err
	}
	preview, _ := params["__preview"].(bool)
	return executeTool(processor, params, preview)
}

// Handler for the 'connection_status' tool
func connectionStatusHandler(ctx context.Context) (*mcp.ToolResponse, error) {
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
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
}

// Handler for the 'glossary_resource' tool
func glossaryResourceHandler(ctx context.Context) (*mcp.ToolResponse, error) {
	slog.Info("[tool_handlers] Handling glossary_resource tool call")
	if len(loadedGlossaries) == 0 {
		slog.Warn("[tool_handlers] No glossary resources loaded")
		return nil, fmt.Errorf("no glossary resources loaded")
	}
	jsonResult, err := json.Marshal(loadedGlossaries[0].Resource.Words)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResponse(mcp.NewTextContent(string(jsonResult))), nil
}
