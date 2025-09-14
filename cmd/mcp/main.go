package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"td_go_mcp/internal/db"
	"td_go_mcp/internal/mcp"
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
	t := mcp.NewTransport(os.Stdin, os.Stdout)
	log.SetOutput(os.Stderr)

	// Ensure database connection is closed on exit
	defer func() {
		if database != nil {
			database.Close()
		}
	}()

	for {
		msg, err := t.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			log.Printf("read error: %v", err)
			return
		}

		var req mcp.Request
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Printf("json err: %v", err)
			continue
		}

		switch req.Method {
		case "initialize":
			handleInitialize(req, t)
		case "tools/list":
			handleToolsList(req, t)
		case "tools/call":
			handleToolsCall(req, t)
		default:
			sendError(req, t, -32601, "Method not found")
		}
	}
}

func handleInitialize(req mcp.Request, t *mcp.Transport) {
	var params mcp.InitializeParams
	if req.Params != nil {
		_ = json.Unmarshal(*req.Params, &params)
	}

	res := mcp.InitializeResult{
		Capabilities: map[string]any{"tools": map[string]any{}},
	}
	res.ServerInfo.Name = "td-go-mcp"
	res.ServerInfo.Version = "0.2.0"

	out, _ := json.Marshal(mcp.Response{JSONRPC: "2.0", ID: req.ID, Result: res})
	_ = t.Write(out)
}

func handleToolsList(req mcp.Request, t *mcp.Transport) {
	var mcpTools []mcp.Tool

	for i := range loadedTools {
		mcpTool := loadedTools[i].ToMCPTool()
		mcpTools = append(mcpTools, mcp.Tool{
			Name:        mcpTool["name"].(string),
			Description: mcpTool["description"].(string),
			InputSchema: mcpTool["inputSchema"].(map[string]any),
		})
	}

	out, _ := json.Marshal(mcp.Response{JSONRPC: "2.0", ID: req.ID, Result: mcp.ToolsListResult{Tools: mcpTools}})
	_ = t.Write(out)
}

func handleToolsCall(req mcp.Request, t *mcp.Transport) {
	var params mcp.ToolCallParams
	if req.Params != nil {
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			sendError(req, t, -32602, "Invalid params: "+err.Error())
			return
		}
	}

	// Validate tool name
	if params.Name == "" {
		sendError(req, t, -32602, "Tool name is required")
		return
	}

	processor, exists := processors[params.Name]
	if !exists {
		sendError(req, t, -32601, "Tool not found: "+params.Name)
		return
	}

	// Ensure args is not nil
	if params.Args == nil {
		params.Args = make(map[string]interface{})
	}

	// Check if this is a preview request
	preview := false
	if p, ok := params.Args["__preview"]; ok {
		if b, ok := p.(bool); ok {
			preview = b
		}
		delete(params.Args, "__preview") // Remove from args for validation
	}

	// Validate parameters against tool schema
	if err := processor.ValidateParameters(params.Args); err != nil {
		sendError(req, t, -32602, "Parameter validation failed: "+err.Error())
		return
	}

	// Process SQL template
	sql, err := processor.ProcessTemplate(params.Args)
	if err != nil {
		sendError(req, t, -32603, "SQL template processing failed: "+err.Error())
		return
	}

	// Validate SQL is not empty
	if strings.TrimSpace(sql) == "" {
		sendError(req, t, -32603, "Generated SQL is empty")
		return
	}

	var result mcp.ToolCallResult
	if preview {
		// Return the generated SQL instead of executing it
		result.Content = []mcp.ToolContent{{
			Type: "text",
			Text: "Generated SQL:\n" + sql,
		}}
	} else {
		// Execute SQL against database and return results
		if database == nil {
			result.Content = []mcp.ToolContent{{
				Type: "text",
				Text: "Database connection not available. Use '__preview': true to see generated SQL.\n\nGenerated SQL:\n" + sql,
			}}
		} else {
			rows, err := database.ExecuteQuery(sql)
			if err != nil {
				sendError(req, t, -32603, "SQL execution failed: "+err.Error())
				return
			}

			// Convert results to JSON
			resultJSON, err := json.Marshal(map[string]interface{}{
				"rows":  rows,
				"count": len(rows),
				"sql":   sql,
			})
			if err != nil {
				sendError(req, t, -32603, "Failed to marshal results: "+err.Error())
				return
			}

			result.Content = []mcp.ToolContent{{
				Type: "text",
				Text: string(resultJSON),
			}}
		}
	}

	out, _ := json.Marshal(mcp.Response{JSONRPC: "2.0", ID: req.ID, Result: result})
	_ = t.Write(out)
}

func sendError(req mcp.Request, t *mcp.Transport, code int, message string) {
	out, _ := json.Marshal(mcp.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error:   &mcp.RespError{Code: code, Message: message},
	})
	_ = t.Write(out)
}
